//+----------------------------------------------------------------------
// | Copyright (c) 2024 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: tulihua2004@126.com
// +----------------------------------------------------------------------

package service

import (
	"errors"
	"github.com/vuecmf/vuecmf-go/v3/app"
	"github.com/vuecmf/vuecmf-go/v3/app/vuecmf/model"
	"gorm.io/gorm"
	"strings"
	"sync"
)

// ModelActionService modelAction服务结构
type ModelActionService struct {
	*BaseService
}

// Create 创建单条或多条数据, 成功返回影响行数
//
//	参数：
//		data 需保存的数据
func (svc *ModelActionService) Create(data *model.ModelAction) (int64, error) {
	var num int64
	DbTable("model_action").
		Where("model_id = ?", data.ModelId).
		Where("action_type = ?", data.ActionType).
		Count(&num)
	if num > 0 {
		return 0, errors.New("动作类型（" + data.ActionType + "）已存在")
	}
	res := app.Db.Create(data)
	return res.RowsAffected, res.Error
}

// Update 更新数据, 成功返回影响行数
//
//	参数：
//		data 需更新的数据
func (svc *ModelActionService) Update(data *model.ModelAction) (int64, error) {
	var old model.ModelAction
	DbTable("model_action").
		Where("id = ?", data.Id).
		Find(&old)

	var num int64
	DbTable("model_action").
		Where("model_id = ?", data.ModelId).
		Where("action_type = ?", data.ActionType).
		Count(&num)

	if num > 0 && old.ActionType != data.ActionType {
		return 0, errors.New("动作类型（" + data.ActionType + "）已存在")
	}

	//清除相关权限项
	err := app.Db.Transaction(func(tx *gorm.DB) error {
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
//
//	参数：
//		id 需删除的ID
//		model 模型实例
func (svc *ModelActionService) Delete(id uint, model *model.ModelAction) (int64, error) {
	//清除相关权限项
	var apiPath string
	DbTable("model_action").Select("api_path").
		Where("id = ?", id).Find(&apiPath)

	err := app.Db.Transaction(func(tx *gorm.DB) error {
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
//
//	参数：
//		idList 需删除的ID列表
//		model 模型实例
func (svc *ModelActionService) DeleteBatch(idList string, model *model.ModelAction) (int64, error) {
	idArr := strings.Split(idList, ",")

	err := app.Db.Transaction(func(tx *gorm.DB) error {
		//清除相关权限项
		for _, id := range idArr {
			var apiPath string
			tx.Table(TableName("model_action")).Select("api_path").
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

		return app.Db.Delete(model, idArr).Error
	})

	if err != nil {
		return 0, err
	}
	return int64(len(idArr)), nil
}

// GetAllApiMap 根据动作ID获取 api 路径映射
//
//	参数：
//		apiIdList api ID列表
func (svc *ModelActionService) GetAllApiMap(apiIdList []string) interface{} {
	var res = make(map[string]map[string]string)

	type actionList struct {
		TableName  string
		ActionType string
		ApiPath    string
	}

	var ac []actionList

	DbTable("model_action", "ma").Select("mc.table_name, ma.action_type, ma.api_path").
		Joins("inner join "+TableName("model_config")+" mc on ma.model_id = mc.id").
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
//
//	参数：
//		apiIdList api ID列表
func (svc *ModelActionService) GetModelIdListById(apiIdList []string) []string {
	var res []string
	DbTable("model_action").Select("model_id").
		Where("id in ?", apiIdList).
		Where("status = 10").
		Group("model_id").
		Find(&res)
	res = append(res, "0")
	return res
}

// GetApiMap 获取API映射的路径
//
//	参数：
//		tableName 表名
//		actionType 动作类型
//		appId 应用ID
func (svc *ModelActionService) GetApiMap(tableName string, actionType string, appId uint) string {
	if appId == 0 {
		appId = 1
	}
	var apiPath string
	DbTable("model_action", "ma").Select("api_path").
		Joins("inner join "+TableName("model_config")+" mc on ma.model_id = mc.id").
		Where("mc.table_name = ?", tableName).
		Where("ma.action_type = ?", actionType).
		Where("mc.app_id = ?", appId).
		Where("ma.status = 10").
		Where("mc.status = 10").
		Find(&apiPath)

	return apiPath
}

// GetNotAuthActionIds 获取无需授权应用下的所有动作ID
func (svc *ModelActionService) GetNotAuthActionIds() []string {
	var res []string
	DbTable("model_action", "MA").Select("MA.id").
		Joins("left join " + TableName("model_config") + " MC on MC.id = MA.model_id").
		Joins("left join " + TableName("app_config") + " AC on MC.app_id = AC.id").
		Where("AC.auth_enable = 20").
		Where("MC.status = 10").
		Where("AC.status = 10").
		Where("MA.status = 10").Find(&res)
	return res
}

// GetActionList 获取所有模型的动作列表
//
//	参数：
//		dataForm 表单参数
func (svc *ModelActionService) GetActionList(dataForm *model.DataActionListForm) (interface{}, error) {
	var res = make(map[string]map[string]string)
	type row struct {
		Id         string
		Label      string
		ModelLabel string
	}

	//父级角色
	var pidRoleName string

	if dataForm.Data.RoleName != "" {
		//若传入角色名称，则只取当前角色的父级角色拥有的权限
		var pid uint
		DbTable("roles").
			Select("pid").
			Where("role_name = ?", dataForm.Data.RoleName).
			Where("status = 10").Find(&pid)
		if pid > 0 {
			DbTable("roles").
				Select("role_name").
				Where("id = ?", pid).
				Where("status = 10").Find(&pidRoleName)
		}

		if pidRoleName == "" {
			return res, nil
		}

	} else if dataForm.Data.Username != "" {
		userInfo := Admin().GetUserByUsername(dataForm.Data.Username)
		pidUserInfo := Admin().GetUser(userInfo.Pid)
		roleArr, err := Auth().GetRolesForUser(pidUserInfo.Username)
		if err == nil {
			//多角色的，暂只取一个角色
			pidRoleName = roleArr[0]
		}

		if pidRoleName == "" {
			return res, nil
		}
	}

	if pidRoleName != "" {
		perList, err := Auth().GetPermissions(pidRoleName, 0)
		if err != nil {
			return nil, err
		}

		for modelName, actionIdList := range perList {
			var actionRes []row
			DbTable("model_action").Select("id, label").
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

	//否则，获取所有权限列表(排除关闭权限验证的应用)
	var actionRes []row
	DbTable("model_action", "MA").Select("MA.id, MC.label model_label,  MA.label").
		Joins("left join " + TableName("model_config") + " MC on MC.id = MA.model_id").
		Joins("left join " + TableName("app_config") + " AC on MC.app_id = AC.id").
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

// GetListActionIdByModelId 获取模型的列表动作ID
//
//	参数：
//		modelId 模型ID
func (svc *ModelActionService) GetListActionIdByModelId(modelId uint) uint {
	var actionId uint
	DbTable("model_action").Select("id").
		Where("model_id = ?", modelId).
		Where("action_type = 'list'").
		Where("status = 10").
		Find(&actionId)
	return actionId
}

var modelActionOnce sync.Once
var modelAction *ModelActionService

// ModelAction 获取modelAction服务实例
func ModelAction() *ModelActionService {
	modelActionOnce.Do(func() {
		modelAction = &ModelActionService{
			BaseService: &BaseService{
				"model_action",
				&model.ModelAction{},
				&[]model.ModelAction{},
				[]string{"label", "api_path", "action_type"},
			},
		}
	})
	return modelAction
}
