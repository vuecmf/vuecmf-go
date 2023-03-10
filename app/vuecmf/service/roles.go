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
//	参数：
//		pid 父级ID
func (ser *rolesService) GetIdPath(pid uint) string {
	var pidIdPath string
	Db.Table(NS.TableName("roles")).Select("id_path").Where("id = ?", pid).Find(&pidIdPath)
	if pid > 0 {
		if pidIdPath == "" {
			pidIdPath = strconv.Itoa(int(pid))
		} else {
			pidIdPath += "," + strconv.Itoa(int(pid))
		}
	}
	return pidIdPath
}

// Create 创建单条或多条数据, 成功返回影响行数
//	参数：
//		data 需保存的数据
func (ser *rolesService) Create(data *model.Roles) (int64, error) {
	data.IdPath = ser.GetIdPath(data.Pid)
	res := Db.Create(data)
	return res.RowsAffected, res.Error
}

// Update 更新数据, 成功返回影响行数
//	参数：
//		data 需更新的数据
func (ser *rolesService) Update(data *model.Roles) (int64, error) {
	var oldRoleName string
	Db.Table(NS.TableName("roles")).Select("role_name").
		Where("id = ?", data.Id).Find(&oldRoleName)

	data.IdPath = ser.GetIdPath(data.Pid)
	res := Db.Updates(data)

	if oldRoleName != "" && oldRoleName != data.RoleName {
		if err := Auth().UpdateRoles(oldRoleName, data.RoleName); err != nil {
			return 0, err
		}
	}

	return res.RowsAffected, res.Error
}

// Delete 根据ID删除数据
//	参数：
//		id 需删除的ID
// 		model 模型实例
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
//	参数：
//		idList 需删除的ID列表
// 		model 模型实例
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
//	参数：
//		params 查询参数
func (ser *rolesService) List(params *helper.DataListParams) (interface{}, error) {
	if params.Data.Action == "getField" {
		//拉取列表的字段信息
		return ser.GetFieldList(ser.TableName, params.Data.Filter, 10)
	} else {
		//拉取列表的数据
		var rolesList []*model.Roles
		var res = make(map[string]interface{})

		ser.GetList(&rolesList, ser.TableName, params.Data)

		//转换成树形列表
		tree := model.RolesModel().ToTree(rolesList)
		res["data"] = tree
		return res, nil
	}
}

// AddUsers 给角色分配用户
//	参数：
//		roleName 角色名
// 		userIdList 用户ID列表
func (ser *rolesService) AddUsers(roleName string, userIdList []int) (interface{}, error) {
	if len(userIdList) == 0 {
		//若传入的为空，则先查出该角色下原有用户列表，然后全部删除
		oldUserIdList, err := ser.GetUsers(roleName)
		if err != nil {
			return nil, errors.New("该角色(" + roleName + ")没有分配用户")
		}
		return Auth().DelUsersForRole(roleName, oldUserIdList)
	}

	var userList []string
	Db.Table(NS.TableName("admin")).Select("username").
		Where("id in ?", userIdList).Find(&userList)

	return Auth().AddUsersForRole(roleName, userList)
}

// GetUsers 获取角色下所有用户的ID
//	参数：
//		roleName 角色名
func (ser *rolesService) GetUsers(roleName string) ([]int, error) {
	userList, err := Auth().GetUsers(roleName)
	if err != nil {
		return nil, errors.New("该角色(" + roleName + ")没有分配任务用户")
	}

	var userIdList []int
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

//GetRoleNameList 根据角色ID获取角色名称
//	参数：
//		roleIdList 角色ID列表
func (ser *rolesService) GetRoleNameList(roleIdList []int) []string {
	var res []string
	Db.Table(NS.TableName("roles")).Select("role_name").Where("id in ?", roleIdList).Find(&res)
	return res
}

//GetRoleIdList 根据角色名称获取角色ID
//	参数：
//		roleNameList 角色名称列表
func (ser *rolesService) GetRoleIdList(roleNameList []string) []int {
	var res []int
	Db.Table(NS.TableName("roles")).Select("id").
		Where("role_name in ?", roleNameList).
		Where("status = 10").
		Find(&res)
	return res
}
