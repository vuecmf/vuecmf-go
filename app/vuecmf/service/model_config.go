//+----------------------------------------------------------------------
// | Copyright (c) 2024 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: tulihua2004@126.com
// +----------------------------------------------------------------------

package service

import (
	"github.com/vuecmf/vuecmf-go/v3/app"
	"github.com/vuecmf/vuecmf-go/v3/app/vuecmf/model"
	"gorm.io/gorm"
	"strings"
	"sync"
)

// ModelConfigService modelConfig服务结构
type ModelConfigService struct {
	*BaseService
}

var modelConfigOnce sync.Once
var modelConfig *ModelConfigService

// ModelConfig 获取modelConfig服务实例
func ModelConfig() *ModelConfigService {
	modelConfigOnce.Do(func() {
		modelConfig = &ModelConfigService{
			BaseService: &BaseService{
				"model_config",
				&model.ModelConfig{},
				&[]model.ModelConfig{},
				[]string{"table_name", "label", "component_tpl", "remark"},
			},
		}
	})
	return modelConfig
}

// Create 创建单条或多条数据, 成功返回影响行数
//
//	参数：
//		data 需保存的数据
func (svc *ModelConfigService) Create(data *model.ModelConfig) (int64, error) {
	//初始化模型相关数据
	if err := Make().BuildModel(data); err != nil {
		return 0, err
	}
	return 1, nil
}

type modelConfigInfo struct {
	TableName string
	AppName   string
}

// Update 更新数据, 成功返回影响行数
//
//	参数：
//		data 需更新的数据
func (svc *ModelConfigService) Update(data *model.ModelConfig) (int64, error) {
	var oldModel modelConfigInfo
	DbTable("model_config", "MC").Select("MC.table_name, AC.app_name").
		Joins("left join "+TableName("app_config")+" AC on MC.app_id = AC.id").
		Where("MC.id = ?", data.Id).
		Where("MC.status = 10").
		Where("AC.status = 10").
		Find(&oldModel)

	err := app.Db.Transaction(func(tx *gorm.DB) error {
		// 若更新时，修改了表名，则相应修改数据库表名
		if oldModel.TableName != "" && oldModel.TableName != data.TableName {
			if err := tx.Migrator().RenameTable(TableName(oldModel.TableName), TableName(data.TableName)); err != nil {
				return err
			}
			//清除原表相关代码文件，重新生成新的代码文件
			if err := Make().RemoveAll(oldModel.TableName); err != nil {
				return err
			}
			if err := Make().MakeAll(data.TableName); err != nil {
				return err
			}
		} else {
			//否则只更新模型
			appName := AppConfig().GetAppNameById(data.AppId)
			if err := Make().Model(data.TableName, appName); err != nil {
				return err
			}
		}

		//更新动作表中的api_path
		appName := AppConfig().GetAppNameById(data.AppId)
		oldPath := "/" + oldModel.AppName + "/" + oldModel.TableName
		newPath := "/" + appName + "/" + data.TableName
		tx.Table(TableName("model_action")).
			Where("model_id = ?", data.Id).
			Where("status = 10").Update("api_path", gorm.Expr("replace(api_path,?,?)", oldPath, newPath))

		return tx.Updates(data).Error
	})

	if err != nil {
		return 0, err
	}
	return 1, nil
}

// Delete 根据ID删除数据
//
//	参数：
//		id 需删除的ID
//		model 模型实例
func (svc *ModelConfigService) Delete(id uint, model *model.ModelConfig) (int64, error) {
	err := app.Db.Transaction(func(tx *gorm.DB) error {
		model.Id = id
		model.TableName = svc.GetModelTableName(int(id))
		//清除相关数据
		if err := Make().RemoveModelData(model); err != nil {
			return err
		}
		if err := tx.Delete(model, id).Error; err != nil {
			return err
		}
		//清除相关代码文件
		return Make().RemoveAll(model.TableName)
	})

	if err != nil {
		return 0, err
	}

	return 1, err
}

// DeleteBatch 根据ID删除数据， 多个用英文逗号分隔
//
//	参数：
//		idList 需删除的ID列表
//		modelInstace 模型实例
func (svc *ModelConfigService) DeleteBatch(idList string, modelInstace *model.ModelConfig) (int64, error) {
	idArr := strings.Split(idList, ",")
	err := app.Db.Transaction(func(tx *gorm.DB) error {
		var modelList []*model.ModelConfig
		DbTable("model_config").Select("id,table_name").
			Where("id in ?", idArr).
			Where("status = 10").Find(&modelList)

		if err := tx.Delete(modelInstace, idArr).Error; err != nil {
			return err
		}
		for _, mc := range modelList {
			//清除相关数据
			if err := Make().RemoveModelData(mc); err != nil {
				return err
			}
			//清除相关代码文件
			if err := Make().RemoveAll(mc.TableName); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	return 1, err
}

// GetModelId 根据表名获取模型ID
//
//	参数：
//		tableName 表名
func (svc *ModelConfigService) GetModelId(tableName string) int {
	var modelId int
	DbTable("model_config").Select("id").
		Where("table_name = ?", tableName).
		Where("status = 10").
		Limit(1).Find(&modelId)
	return modelId
}

// GetModelTableName 根据模型ID获取模型对应表名
//
//	参数：
//		modelId 模型ID
func (svc *ModelConfigService) GetModelTableName(modelId int) string {
	var tableName string
	DbTable("model_config").Select("table_name").
		Where("id = ?", modelId).
		Limit(1).Find(&tableName)
	return tableName
}

// IsTree 根据模型ID判断是否为目录树
//
//	参数：
//		modelId 模型ID
func (svc *ModelConfigService) IsTree(modelId int) bool {
	var isTree int
	DbTable("model_config").Select("is_tree").
		Where("id = ?", modelId).
		Limit(1).Find(&isTree)
	return isTree == 10
}

type modelConf struct {
	TableName      string
	IsTree         bool
	ModelId        int
	LabelFieldName string
}

// GetModelConfig 根据模型表名获取模型的配置信息
//
//	参数：
//		tableName 表名
func (svc *ModelConfigService) GetModelConfig(tableName string) modelConf {
	var mc modelConf
	DbTable("model_config").
		Select("table_name, if(is_tree = 10, true, false) is_tree, id model_id, '' label_field_name").
		Where("status = 10").
		Where("table_name = ?", tableName).
		Limit(1).
		Find(&mc)

	var labelFieldName string
	DbTable("model_field").
		Select("field_name").
		Where("status = 10").
		Where("model_id = ?", mc.ModelId).
		Where("is_label = 10").
		Find(&labelFieldName)

	mc.LabelFieldName = labelFieldName

	return mc
}
