// Package controller
//+----------------------------------------------------------------------
// | Copyright (c) 2022 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: vuecmf <tulihua2004@126.com>
// +----------------------------------------------------------------------
package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/vuecmf/vuecmf-go/app/route"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/helper"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/model"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/service"
)

type Roles struct {
	Base
}

func init() {
	roles := &Roles{}
	roles.TableName = "roles"
	roles.Model = &model.Roles{}
	roles.ListData = &[]model.Roles{}
	roles.FilterFields = []string{"role_name", "app_name", "id_path", "remark"}

	route.Register(roles, "POST", "vuecmf")
}

// Index 列表页
func (ctrl *Roles) Index(c *gin.Context) {
	listParams := &helper.DataListParams{}
	Common(c, listParams, func() (interface{}, error) {
		return service.Roles().List(listParams)
	})
}

// Save 新增/更新 单条数据
func (ctrl *Roles) Save(c *gin.Context) {
	saveForm := &model.DataRolesForm{}
	Common(c, saveForm, func() (interface{}, error) {
		if saveForm.Data.Id == uint(0) {
			return service.Roles().Create(saveForm.Data)
		} else {
			return service.Roles().Update(saveForm.Data)
		}
	})
}

// AddUsers 给角色分配用户
func (ctrl *Roles) AddUsers(c *gin.Context) {
	roleUsersForm := &model.DataRoleUsersForm{}
	Common(c, roleUsersForm, func() (interface{}, error) {
		return service.Roles().AddUsers(roleUsersForm.Data.RoleName, roleUsersForm.Data.UseridList)
	})
}

// DelUsers 删除角色下的用户
func (ctrl *Roles) DelUsers(c *gin.Context) {
	roleUsersForm := &model.DataRoleUsersForm{}
	Common(c, roleUsersForm, func() (interface{}, error) {
		return service.Auth().DelUsersForRole(roleUsersForm.Data.RoleName, roleUsersForm.Data.UseridList)
	})
}

// AddPermission 给角色分配权限项
func (ctrl *Roles) AddPermission(c *gin.Context) {
	permissionForm := &model.DataPermissionForm{}
	Common(c, permissionForm, func() (interface{}, error) {
		return service.Auth().AddPermission(permissionForm.Data.RoleName, permissionForm.Data.ActionId)
	})
}

// DelPermission 删除角色的权限项
func (ctrl *Roles) DelPermission(c *gin.Context) {
	permissionForm := &model.DataPermissionForm{}
	Common(c, permissionForm, func() (interface{}, error) {
		return service.Auth().DelPermission(permissionForm.Data.RoleName, permissionForm.Data.ActionId)
	})
}

// GetUsers 获取角色下所有用户的ID
func (ctrl *Roles) GetUsers(c *gin.Context) {
	roleForm := &model.DataRoleForm{}
	Common(c, roleForm, func() (interface{}, error) {
		return service.Roles().GetUsers(roleForm.Data.RoleName)
	})
}

// GetPermission 获取角色的所有权限项
func (ctrl *Roles) GetPermission(c *gin.Context) {
	roleForm := &model.DataRoleForm{}
	Common(c, roleForm, func() (interface{}, error) {
		return service.Auth().GetPermissions(roleForm.Data.RoleName, nil)
	})
}

// GetAllUsers 获取所有用户
func (ctrl *Roles) GetAllUsers(c *gin.Context) {
	Common(c, nil, func() (interface{}, error) {
		return service.Roles().GetAllUsers()
	})
}
