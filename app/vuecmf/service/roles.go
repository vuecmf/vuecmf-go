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
	"github.com/vuecmf/vuecmf-go/app/vuecmf/helper"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/model"
)

// rolesService roles服务结构
type rolesService struct {
	*baseService
	TableName string
}

var roles *rolesService

// Roles 获取roles服务实例
func Roles() *rolesService {
	if roles == nil {
		roles = &rolesService{TableName: "roles"}
	}
	return roles
}

// List 获取列表数据
// 		参数：params 查询参数
func (ser *rolesService) List(params *helper.DataListParams) (interface{}, error) {
	if params.Data.Action == "getField" {
		//拉取列表的字段信息
		return ser.getFieldList(ser.TableName, params.Data.Filter)
	} else {
		//拉取列表的数据
		var rolesList []*model.Roles
		var res = make(map[string]interface{})

		ser.getList(&rolesList, ser.TableName, params)

		//转换成树形列表
		tree := model.RolesModel().ToTree(rolesList)
		res["data"] = tree
		return res, nil
	}
}

// AddUsers 给角色分配用户
func (ser *rolesService) AddUsers(roleName string, userIdList []string, appName string) (interface{}, error) {
	if appName == "" {
		appName = "vuecmf"
	}

	if len(userIdList) == 0 {
		//若传入的为空，则先查出该角色下原有用户列表，然后全部删除
		userList, err := Auth().GetUsers(roleName, appName)
		if err != nil {
			return nil, errors.New("该角色(" + roleName + ")没有分配任务用户")
		}
		return Auth().DelUsersForRole(roleName, userList, appName)
	}

	var userList []string
	db.Table(ns.TableName("admin")).Select("username").
		Where("id in ?", userIdList).
		Where("status = 10").Find(&userList)

	return Auth().AddUsersForRole(roleName, userList, appName)
}

// GetUsers 获取角色下所有用户的ID
func (ser *rolesService) GetUsers(roleName string, appName string) (interface{}, error) {
	if appName == "" {
		appName = "vuecmf"
	}

	userList, err := Auth().GetUsers(roleName, appName)
	if err != nil {
		return nil, errors.New("该角色(" + roleName + ")没有分配任务用户")
	}

	var userIdList []string
	db.Table(ns.TableName("admin")).Select("id").
		Where("username in ?", userList).
		Where("status = 10").Find(&userIdList)

	return userIdList, nil
}

// GetAllUsers 获取所有用户
func (ser *rolesService) GetAllUsers() (interface{}, error) {
	type row struct {
		Key      uint   `json:"key"`
		Label    string `json:"label"`
		Disabled bool   `json:"disabled"`
	}

	var res []*row

	db.Table(ns.TableName("admin")).Select("id `key`, username label, false disabled").
		Where("status = 10").
		Where("is_super != 10").Find(&res)

	return res, nil
}
