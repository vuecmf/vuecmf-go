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
	"strconv"
	"strings"
)

// rolesService roles服务结构
type rolesService struct {
	*BaseService
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

// GetIdPath 获取父级ID的ID路径
func (ser *rolesService) GetIdPath(pid uint) string {
	var pidIdPath string
	Db.Table(NS.TableName("roles")).Select("id_path").Where("id = ?", pid).Find(&pidIdPath)
	if pid > 0 {
		pidIdPath += "," + strconv.Itoa(int(pid))
	}
	return pidIdPath
}

// Create 创建单条或多条数据, 成功返回影响行数
func (ser *rolesService) Create(data *model.Roles) (int64, error) {
	data.IdPath = ser.GetIdPath(data.Pid)
	res := Db.Create(data)
	return res.RowsAffected, res.Error
}

// Update 更新数据, 成功返回影响行数
func (ser *rolesService) Update(data *model.Roles) (int64, error) {
	var oldRoleName string
	Db.Table(NS.TableName("roles")).Select("role_name").
		Where("id = ?", data.Id).Find(&oldRoleName)
	if oldRoleName != "" && oldRoleName != data.RoleName {
		if _, err := Auth().Enforcer.DeleteRole(oldRoleName); err != nil {
			return 0, err
		}
	}
	data.IdPath = ser.GetIdPath(data.Pid)
	res := Db.Updates(data)
	return res.RowsAffected, res.Error
}

// Delete 根据ID删除数据
func (ser *rolesService) Delete(id uint, model *model.Roles) (int64, error) {
	var roleName string
	Db.Table(NS.TableName("roles")).Select("role_name").
		Where("id = ?", id).Find(&roleName)
	if _, err := Auth().Enforcer.DeleteRole(roleName); err != nil {
		return 0, err
	}

	res := Db.Delete(model, id)
	return res.RowsAffected, res.Error
}

// DeleteBatch 根据ID删除数据， 多个用英文逗号分隔
func (ser *rolesService) DeleteBatch(idList string, model *model.Roles) (int64, error) {
	idArr := strings.Split(idList, ",")
	for _, id := range idArr {
		var roleName string
		Db.Table(NS.TableName("roles")).Select("role_name").
			Where("id = ?", id).Find(&roleName)
		if _, err := Auth().Enforcer.DeleteRole(roleName); err != nil {
			return 0, err
		}
	}

	res := Db.Delete(model, idArr)
	return res.RowsAffected, res.Error
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
	Db.Table(NS.TableName("admin")).Select("username").
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
	Db.Table(NS.TableName("admin")).Select("id").
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

	Db.Table(NS.TableName("admin")).Select("id `key`, username label, false disabled").
		Where("status = 10").
		Where("is_super != 10").Find(&res)

	return res, nil
}
