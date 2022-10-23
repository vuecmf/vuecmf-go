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
	"github.com/vuecmf/vuecmf-go/app"
	"github.com/vuecmf/vuecmf-go/app/route"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/model"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/service"
	"time"
)

type Admin struct {
	Base
}

func init() {
	admin := &Admin{}
	admin.TableName = "admin"
	admin.Model = &model.Admin{}
	admin.ListData = &[]model.Admin{}
	admin.FilterFields = []string{"username", "email", "mobile", "token"}

	admin.AppName = "vuecmf"

	route.Register(admin, "POST", admin.AppName)
}

// Save 新增/更新 单条数据
func (ctrl *Admin) Save(c *gin.Context) {
	saveForm := &model.DataAdminForm{}
	common(c, saveForm, func() (interface{}, error) {
		if saveForm.Data.Id == uint(0) {
			return service.Base().Create(saveForm.Data)
		} else {
			return service.Base().Update(saveForm.Data)
		}
	})
}

// Login 用户登录
func (ctrl *Admin) Login(c *gin.Context) {
	dataLoginForm := &model.DataLoginForm{}
	common(c, dataLoginForm, func() (interface{}, error) {
		dataLoginForm.Data.LastLoginIp = c.ClientIP()
		dataLoginForm.Data.LastLoginTime = model.JSONTime{Time: time.Now()}
		return service.Admin(ctrl.AppName).Login(dataLoginForm.Data)
	})
}

// Logout 退出登录
func (ctrl *Admin) Logout(c *gin.Context) {
	dataLogoutForm := &model.DataLogoutForm{}
	common(c, dataLogoutForm, func() (interface{}, error) {
		return service.Admin(ctrl.AppName).Logout(dataLogoutForm.Data)
	})
}

// AddRole 给用户添加角色
func (ctrl *Admin) AddRole(c *gin.Context) {
	dataAddRoleForm := &model.DataAddRoleForm{}
	common(c, dataAddRoleForm, func() (interface{}, error) {
		if dataAddRoleForm.Data.AppName == "" {
			dataAddRoleForm.Data.AppName = "vuecmf"
		}
		if len(dataAddRoleForm.Data.RoleIdList) == 0 {
			//如果角色列表为空，即表示用户没有角色，则清空用户所有角色
			return service.Auth().DelAllRolesForUser(dataAddRoleForm.Data.Username, dataAddRoleForm.Data.AppName)
		} else {
			return service.Auth().AddRolesForUser(dataAddRoleForm.Data.Username, dataAddRoleForm.Data.RoleIdList, dataAddRoleForm.Data.AppName)
		}
	})
}

// DelRole 删除用户的角色
func (ctrl *Admin) DelRole(c *gin.Context) {
	dataDelRoleForm := &model.DataDelRoleForm{}
	common(c, dataDelRoleForm, func() (interface{}, error) {
		if dataDelRoleForm.Data.AppName == "" {
			dataDelRoleForm.Data.AppName = "vuecmf"
		}
		return service.Auth().DelRolesForUser(dataDelRoleForm.Data.Username, dataDelRoleForm.Data.RoleName, dataDelRoleForm.Data.AppName)
	})
}

// AddPermission 分配角色的权限
func (ctrl *Admin) AddPermission(c *gin.Context) {
	dataPermissionForm := &model.DataPermissionForm{}
	common(c, dataPermissionForm, func() (interface{}, error) {
		return service.Auth().AddPermission(dataPermissionForm.Data.RoleName, dataPermissionForm.Data.ActionId)
	})
}

// DelPermission 删除角色的权限
func (ctrl *Admin) DelPermission(c *gin.Context) {
	dataPermissionForm := &model.DataPermissionForm{}
	common(c, dataPermissionForm, func() (interface{}, error) {
		return service.Auth().DelPermission(dataPermissionForm.Data.RoleName, dataPermissionForm.Data.ActionId)
	})
}

// GetPermission 获取角色的权限
func (ctrl *Admin) GetPermission(c *gin.Context) {
	dataRoleForm := &model.DataRoleForm{}
	common(c, dataRoleForm, func() (interface{}, error) {
		isSuper := app.Request(c).GetCtxVal("is_super")
		return service.Auth().GetPermissions(dataRoleForm.Data.RoleName, isSuper, dataRoleForm.Data.AppName)
	})
}

// GetAllRoles 获取所有角色
func (ctrl *Admin) GetAllRoles(c *gin.Context) {
	common(c, nil, func() (interface{}, error) {
		res := service.Auth().GetAllRoles()
		return res, nil
	})
}

// GetRoles 获取用户的所有角色
func (ctrl *Admin) GetRoles(c *gin.Context) {
	dataUsernameForm := &model.DataUsernameForm{}
	common(c, dataUsernameForm, func() (interface{}, error) {
		return service.Auth().GetRoles(dataUsernameForm.Data.Username, dataUsernameForm.Data.AppName)
	})
}

// SetUserPermission 分配用户的权限
func (ctrl *Admin) SetUserPermission(c *gin.Context) {
	dataUserPermissionForm := &model.DataUserPermissionForm{}
	common(c, dataUserPermissionForm, func() (interface{}, error) {
		return service.Auth().AddPermission(
			dataUserPermissionForm.Data.Username,
			dataUserPermissionForm.Data.ActionId,
		)
	})
}

// GetUserPermission 获取用户的权限
func (ctrl *Admin) GetUserPermission(c *gin.Context) {
	dataUsernameForm := &model.DataUsernameForm{}
	common(c, dataUsernameForm, func() (interface{}, error) {
		//isSuper := app.Request(c).GetCtxVal("is_super")
		return service.Auth().GetPermissions(
			dataUsernameForm.Data.Username,
			//isSuper,
			20,
			dataUsernameForm.Data.AppName,
		)
	})
}
