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
	"errors"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/model"
	"strings"
)

// modelActionService modelAction服务结构
type modelActionService struct {
	*BaseService
}

// Update 更新数据, 成功返回影响行数
func (ser *modelActionService) Update(data *model.ModelAction) (int64, error) {
	//清除相关权限项
	var oldApiPath string
	Db.Table(NS.TableName("model_action")).Select("api_path").
		Where("id = ?", data.Id).Find(&oldApiPath)
	if oldApiPath != "" && oldApiPath != data.ApiPath {
		arr := strings.Split(oldApiPath, "/")
		if len(arr) == 2 {
			if _, err := Auth().Enforcer.DeletePermission(arr[0], arr[1], "index"); err != nil {
				return 0, err
			}
		} else if len(arr) == 3 {
			if _, err := Auth().Enforcer.DeletePermission(arr[0], arr[1], arr[2]); err != nil {
				return 0, err
			}
		}
	}

	res := Db.Updates(data)
	return res.RowsAffected, res.Error
}

// Delete 根据ID删除数据
func (ser *modelActionService) Delete(id uint, model *model.ModelAction) (int64, error) {
	//清除相关权限项
	var apiPath string
	Db.Table(NS.TableName("model_action")).Select("api_path").
		Where("id = ?", id).Find(&apiPath)
	if apiPath != "" {
		arr := strings.Split(apiPath, "/")
		if len(arr) == 2 {
			if _, err := Auth().Enforcer.DeletePermission(arr[0], arr[1], "index"); err != nil {
				return 0, err
			}
		} else if len(arr) == 3 {
			if _, err := Auth().Enforcer.DeletePermission(arr[0], arr[1], arr[2]); err != nil {
				return 0, err
			}
		}
	}

	res := Db.Delete(model, id)
	return res.RowsAffected, res.Error
}

// DeleteBatch 根据ID删除数据， 多个用英文逗号分隔
func (ser *modelActionService) DeleteBatch(idList string, model *model.ModelAction) (int64, error) {
	idArr := strings.Split(idList, ",")

	//清除相关权限项
	for _, id := range idArr {
		var apiPath string
		Db.Table(NS.TableName("model_action")).Select("api_path").
			Where("id = ?", id).Find(&apiPath)
		if apiPath != "" {
			arr := strings.Split(apiPath, "/")
			if len(arr) == 2 {
				if _, err := Auth().Enforcer.DeletePermission(arr[0], arr[1], "index"); err != nil {
					return 0, err
				}
			} else if len(arr) == 3 {
				if _, err := Auth().Enforcer.DeletePermission(arr[0], arr[1], arr[2]); err != nil {
					return 0, err
				}
			}
		}
	}

	res := Db.Delete(model, idArr)
	return res.RowsAffected, res.Error
}

// GetAllApiMap 根据动作ID获取 api 路径映射
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
func (ser *modelActionService) GetApiMap(tableName string, actionType string) string {
	var apiPath string
	Db.Table(NS.TableName("model_action")+" ma").Select("api_path").
		Joins("inner join "+NS.TableName("model_config")+" mc on ma.model_id = mc.id").
		Where("mc.table_name = ?", tableName).
		Where("ma.action_type = ?", actionType).
		Where("ma.status = 10").
		Where("mc.status = 10").
		Find(&apiPath)

	return apiPath
}

// GetActionList 获取所有模型的动作列表
func (ser *modelActionService) GetActionList(roleName string, appName string) (interface{}, error) {
	var res = make(map[string]map[uint]string)
	type row struct {
		Id    uint
		Label string
	}
	if roleName != "" {
		//若传入角色名称，则只取当前角色的父级角色拥有的权限
		var pid uint
		Db.Table(NS.TableName("roles")).
			Select("pid").
			Where("role_name = ?", roleName).
			Where("status = 10").Find(&pid)

		if pid == 0 {
			return nil, errors.New("当前角色没有父级角色")
		}

		//父级角色
		var pidRoleName string
		Db.Table(NS.TableName("roles")).
			Select("role_name").
			Where("id = ?", pid).
			Where("status = 10").Find(&pidRoleName)

		if pidRoleName == "" {
			return nil, errors.New("没有获取到父级角色名称")
		}

		if appName == "" {
			appName = "vuecmf"
		}

		perList, err := Auth().GetPermissions(pidRoleName, nil, appName)
		if err != nil {
			return nil, err
		}

		for modelName, actionIdList := range perList {
			var actionRes []row
			Db.Table(NS.TableName("model_action")).Select("id, label").
				Where("id in ?", actionIdList).
				Where("status = 10").Find(&actionRes)
			res[modelName] = map[uint]string{}
			for _, ac := range actionRes {
				res[modelName][ac.Id] = ac.Label
			}
		}

		return res, nil
	}

	//否则，获取所有权限列表
	var modelListRes []row
	Db.Table(NS.TableName("model_config")).Select("id, label").
		Where("status = 10").Find(&modelListRes)
	for _, mc := range modelListRes {
		var maList []row
		Db.Table(NS.TableName("model_action")).Select("id, label").
			Where("model_id = ?", mc.Id).
			Where("status = 10").Find(&maList)
		res[mc.Label] = map[uint]string{}
		for _, ac := range maList {
			res[mc.Label][ac.Id] = ac.Label
		}
	}

	return res, nil
}

var modelAction *modelActionService

// ModelAction 获取modelAction服务实例
func ModelAction() *modelActionService {
	if modelAction == nil {
		modelAction = &modelActionService{}
	}
	return modelAction
}
