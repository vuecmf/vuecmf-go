//+----------------------------------------------------------------------
// | Copyright (c) 2023 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: vuecmf <tulihua2004@126.com>
// +----------------------------------------------------------------------

package service

import (
	"errors"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/model"
	"gorm.io/gorm"
	"strings"
)

// modelActionService modelAction服务结构
type modelActionService struct {
	*BaseService
}

// Create 创建单条或多条数据, 成功返回影响行数
//	参数：
//		data 需保存的数据
func (ser *modelActionService) Create(data *model.ModelAction) (int64, error) {
	var num int64
	Db.Table(NS.TableName("model_action")).
		Where("model_id = ?", data.ModelId).
		Where("action_type = ?", data.ActionType).
		Count(&num)
	if num > 0 {
		return 0, errors.New("动作类型（" + data.ActionType + "）已存在")
	}
	res := Db.Create(data)
	return res.RowsAffected, res.Error
}

// Update 更新数据, 成功返回影响行数
//	参数：
//		data 需更新的数据
func (ser *modelActionService) Update(data *model.ModelAction) (int64, error) {
	var old model.ModelAction
	Db.Table(NS.TableName("model_action")).
		Where("id = ?", data.Id).
		Find(&old)

	var num int64
	Db.Table(NS.TableName("model_action")).
		Where("model_id = ?", data.ModelId).
		Where("action_type = ?", data.ActionType).
		Count(&num)

	if num > 0 && old.ActionType != data.ActionType {
		return 0, errors.New("动作类型（" + data.ActionType + "）已存在")
	}

	//清除相关权限项
	err := Db.Transaction(func(tx *gorm.DB) error {
		if old.ApiPath != "" && old.ApiPath != data.ApiPath {
			arr := strings.Split(old.ApiPath, "/")
			if len(arr) == 2 {
				if _, err := Auth().Enforcer.DeletePermission(arr[0], arr[1], "index"); err != nil {
					return err
				}
			} else if len(arr) == 3 {
				if _, err := Auth().Enforcer.DeletePermission(arr[0], arr[1], arr[2]); err != nil {
					return err
				}
			}
		}

		return tx.Updates(data).Error
	})

	if err != nil {
		return 0, err
	}
	return 1, nil
}

// Delete 根据ID删除数据
//	参数：
//		id 需删除的ID
// 		model 模型实例
func (ser *modelActionService) Delete(id uint, model *model.ModelAction) (int64, error) {
	//清除相关权限项
	var apiPath string
	Db.Table(NS.TableName("model_action")).Select("api_path").
		Where("id = ?", id).Find(&apiPath)

	err := Db.Transaction(func(tx *gorm.DB) error {
		if apiPath != "" {
			arr := strings.Split(apiPath, "/")
			if len(arr) == 2 {
				if _, err := Auth().Enforcer.DeletePermission(arr[0], arr[1], "index"); err != nil {
					return err
				}
			} else if len(arr) == 3 {
				if _, err := Auth().Enforcer.DeletePermission(arr[0], arr[1], arr[2]); err != nil {
					return err
				}
			}
		}

		return tx.Delete(model, id).Error
	})

	if err != nil {
		return 0, err
	}
	return 1, nil
}

// DeleteBatch 根据ID删除数据， 多个用英文逗号分隔
//	参数：
//		idList 需删除的ID列表
// 		model 模型实例
func (ser *modelActionService) DeleteBatch(idList string, model *model.ModelAction) (int64, error) {
	idArr := strings.Split(idList, ",")

	err := Db.Transaction(func(tx *gorm.DB) error {
		//清除相关权限项
		for _, id := range idArr {
			var apiPath string
			tx.Table(NS.TableName("model_action")).Select("api_path").
				Where("id = ?", id).Find(&apiPath)
			if apiPath != "" {
				arr := strings.Split(apiPath, "/")
				if len(arr) == 2 {
					if _, err := Auth().Enforcer.DeletePermission(arr[0], arr[1], "index"); err != nil {
						return err
					}
				} else if len(arr) == 3 {
					if _, err := Auth().Enforcer.DeletePermission(arr[0], arr[1], arr[2]); err != nil {
						return err
					}
				}
			}
		}

		return Db.Delete(model, idArr).Error
	})

	if err != nil {
		return 0, err
	}
	return int64(len(idArr)), nil
}

// GetAllApiMap 根据动作ID获取 api 路径映射
//	参数：
//		apiIdList api ID列表
func (ser *modelActionService) GetAllApiMap(apiIdList []string) interface{} {
	var res = make(map[string]map[string]string)

	type actionList struct {
		TableName  string
		ActionType string
		ApiPath    string
	}

	var ac []actionList

	Db.Table(NS.TableName("model_action")+" ma").Select("mc.table_name, ma.action_type, ma.api_path").
		Joins("inner join "+NS.TableName("model_config")+" mc on ma.model_id = mc.id").
		Where("ma.id in ?", apiIdList).
		Where("ma.status = 10").
		Where("mc.status = 10").
		Find(&ac)

	for _, val := range ac {
		if res[val.TableName] == nil {
			res[val.TableName] = map[string]string{}
		}
		res[val.TableName][val.ActionType] = val.ApiPath
	}

	return res
}

// GetModelIdListById 根据动作ID获取对应模型ID列表
//	参数：
//		apiIdList api ID列表
func (ser *modelActionService) GetModelIdListById(apiIdList []string) []string {
	var res []string
	Db.Table(NS.TableName("model_action")).Select("model_id").
		Where("id in ?", apiIdList).
		Where("status = 10").
		Group("model_id").
		Find(&res)
	res = append(res, "0")
	return res
}

// GetApiMap 获取API映射的路径
//	参数：
//		tableName 表名
//		actionType 动作类型
//		appId 应用ID
func (ser *modelActionService) GetApiMap(tableName string, actionType string, appId uint) string {
	if appId == 0 {
		appId = 1
	}
	var apiPath string
	Db.Table(NS.TableName("model_action")+" ma").Select("api_path").
		Joins("inner join "+NS.TableName("model_config")+" mc on ma.model_id = mc.id").
		Where("mc.table_name = ?", tableName).
		Where("ma.action_type = ?", actionType).
		Where("mc.app_id = ?", appId).
		Where("ma.status = 10").
		Where("mc.status = 10").
		Find(&apiPath)

	return apiPath
}

// GetNotAuthActionIds 获取无需授权应用下的所有动作ID
func (ser *modelActionService) GetNotAuthActionIds() []string {
	var res []string
	Db.Table(NS.TableName("model_action") + " MA").Select("MA.id").
		Joins("left join " + NS.TableName("model_config") + " MC on MC.id = MA.model_id").
		Joins("left join " + NS.TableName("app_config") + " AC on MC.app_id = AC.id").
		Where("AC.auth_enable = 20").
		Where("MC.status = 10").
		Where("AC.status = 10").
		Where("MA.status = 10").Find(&res)
	return res
}

// GetActionList 获取所有模型的动作列表
//	参数：
//		roleName 角色名
func (ser *modelActionService) GetActionList(roleName string) (interface{}, error) {
	var res = make(map[string]map[string]string)
	type row struct {
		Id         string
		Label      string
		ModelLabel string
	}

	if roleName != "" {
		//若传入角色名称，则只取当前角色的父级角色拥有的权限
		var pid uint
		Db.Table(NS.TableName("roles")).
			Select("pid").
			Where("role_name = ?", roleName).
			Where("status = 10").Find(&pid)

		if pid > 0 {
			//父级角色
			var pidRoleName string
			Db.Table(NS.TableName("roles")).
				Select("role_name").
				Where("id = ?", pid).
				Where("status = 10").Find(&pidRoleName)

			if pidRoleName == "" {
				return nil, errors.New("没有获取到父级角色名称")
			}

			perList, err := Auth().GetPermissions(pidRoleName, nil)
			if err != nil {
				return nil, err
			}

			for modelName, actionIdList := range perList {
				var actionRes []row
				Db.Table(NS.TableName("model_action")).Select("id, label").
					Where("id in ?", actionIdList).
					Where("status = 10").Find(&actionRes)
				if res[modelName] == nil {
					res[modelName] = make(map[string]string)
				}
				for _, ac := range actionRes {
					res[modelName][ac.Id] = ac.Label
				}
			}

			return res, nil
		}
	}

	//否则，获取所有权限列表(排除关闭权限验证的应用)
	var actionRes []row
	Db.Table(NS.TableName("model_action")+" MA").Select("MA.id, MC.label model_label,  MA.label").
		Joins("left join "+NS.TableName("model_config")+" MC on MC.id = MA.model_id").
		Joins("left join "+NS.TableName("app_config")+" AC on MC.app_id = AC.id").
		Where("AC.auth_enable = 10").
		Where("MC.status = 10").
		Where("AC.status = 10").
		Where("MA.status = 10").Find(&actionRes)

	for _, ac := range actionRes {
		if res[ac.ModelLabel] == nil {
			res[ac.ModelLabel] = make(map[string]string)
		}
		res[ac.ModelLabel][ac.Id] = ac.Label
	}

	return res, nil
}

//GetListActionIdByModelId 获取模型的列表动作ID
//	参数：
//		modelId 模型ID
func (ser *modelActionService) GetListActionIdByModelId(modelId uint) uint {
	var actionId uint
	Db.Table(NS.TableName("model_action")).Select("id").
		Where("model_id = ?", modelId).
		Where("action_type = 'list'").
		Where("status = 10").
		Find(&actionId)
	return actionId
}

var modelAction *modelActionService

// ModelAction 获取modelAction服务实例
func ModelAction() *modelActionService {
	if modelAction == nil {
		modelAction = &modelActionService{}
	}
	return modelAction
}
