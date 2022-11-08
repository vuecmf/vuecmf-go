// Package service
//+----------------------------------------------------------------------
// | Copyright (c) 2022 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: vuecmf <tulihua2004@126.com>
// +----------------------------------------------------------------------
package service

import (
	"github.com/vuecmf/vuecmf-go/app/vuecmf/model"
	"gorm.io/gorm"
	"strings"
)

// modelConfigService modelConfig服务结构
type modelConfigService struct {
	*BaseService
}

var modelConfig *modelConfigService

// ModelConfig 获取modelConfig服务实例
func ModelConfig() *modelConfigService {
	if modelConfig == nil {
		modelConfig = &modelConfigService{}
	}
	return modelConfig
}

// Create 创建单条或多条数据, 成功返回影响行数
func (s *modelConfigService) Create(data *model.ModelConfig) (int64, error) {
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
func (s *modelConfigService) Update(data *model.ModelConfig) (int64, error) {
	var oldModel modelConfigInfo
	Db.Table(NS.TableName("model_config")+" MC").Select("MC.table_name, AC.app_name").
		Joins("left join "+NS.TableName("app_config")+" AC on MC.app_id = AC.id").
		Where("MC.id = ?", data.Id).
		Where("MC.status = 10").
		Where("AC.status = 10").
		Find(&oldModel)

	err := Db.Transaction(func(tx *gorm.DB) error {
		// 若更新时，修改了表名，则相应修改数据库表名
		if oldModel.TableName != "" && oldModel.TableName != data.TableName {
			if err := tx.Migrator().RenameTable(NS.TableName(oldModel.TableName), NS.TableName(data.TableName)); err != nil {
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
		tx.Table(NS.TableName("model_action")).
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
func (s *modelConfigService) Delete(id uint, model *model.ModelConfig) (int64, error) {
	err := Db.Transaction(func(tx *gorm.DB) error {
		model.Id = id
		model.TableName = s.GetModelTableName(int(id))
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

	if err = Make().UpdateRunFile(); err != nil {
		return 0, err
	}
	return 1, err
}

// DeleteBatch 根据ID删除数据， 多个用英文逗号分隔
func (s *modelConfigService) DeleteBatch(idList string, modelInstace *model.ModelConfig) (int64, error) {
	idArr := strings.Split(idList, ",")
	err := Db.Transaction(func(tx *gorm.DB) error {
		var modelList []*model.ModelConfig
		Db.Table(NS.TableName("model_config")).Select("id,table_name").
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
func (s *modelConfigService) GetModelId(tableName string) int {
	var modelId int
	Db.Table(NS.TableName("model_config")).Select("id").
		Where("table_name = ?", tableName).
		Where("status = 10").
		Limit(1).Find(&modelId)
	return modelId
}

// GetModelTableName 根据模型ID获取模型对应表名
func (s *modelConfigService) GetModelTableName(modelId int) string {
	var tableName string
	Db.Table(NS.TableName("model_config")).Select("table_name").
		Where("id = ?", modelId).
		Limit(1).Find(&tableName)
	return tableName
}

// IsTree 根据模型ID判断是否为目录树
func (s *modelConfigService) IsTree(modelId int) bool {
	var isTree int
	Db.Table(NS.TableName("model_config")).Select("is_tree").
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
func (s *modelConfigService) GetModelConfig(tableName string) modelConf {
	var mc modelConf
	Db.Table(NS.TableName("model_config")).
		Select("table_name, if(is_tree = 10, true, false) is_tree, id model_id, '' label_field_name").
		Where("status = 10").
		Where("table_name = ?", tableName).
		Limit(1).
		Find(&mc)

	var labelFieldName string
	Db.Table(NS.TableName("model_field")).
		Select("field_name").
		Where("status = 10").
		Where("model_id = ?", mc.ModelId).
		Where("is_label = 10").
		Find(&labelFieldName)

	mc.LabelFieldName = labelFieldName

	return mc
}
