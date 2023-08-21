//+----------------------------------------------------------------------
// | Copyright (c) 2023 http://www.vuecmf.com All rights reserved.
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
	"github.com/vuecmf/vuecmf-go/app/vuecmf/helper"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/model"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/service"
	"strconv"
	"strings"
	"time"
)

// Admin 管理员管理
type Admin struct {
	Base
}

func init() {
	admin := &Admin{}
	admin.TableName = "admin"
	admin.Model = &model.Admin{}
	admin.ListData = &[]model.Admin{}
	admin.FilterFields = []string{"username", "email", "mobile", "token"}

	route.Register(admin, "POST", "vuecmf")
}

// Index 列表页
func (ctrl *Admin) Index(c *gin.Context) {
	listParams := &helper.DataListParams{}
	Common(c, listParams, func() (interface{}, error) {
		isSuper := helper.InterfaceToInt(app.Request(c).GetCtxVal("is_super"))
		if isSuper != 10 {
			//非超级管理员，只能看自己和自己创建的账号
			uid := strconv.Itoa(app.Request(c).GetCtxVal("uid").(int))
			condition := make(map[string]interface{})
			condition["pid"] = uid

			listParams.Data.Filter["id"] = uid
			listParams.Data.Filter["or"] = condition
		}

		return service.Base().CommonList(ctrl.ListData, ctrl.TableName, ctrl.FilterFields, listParams, isSuper)
	})
}

// Save 新增/更新 单条数据
func (ctrl *Admin) Save(c *gin.Context) {
	saveForm := &model.DataAdminForm{}
	Common(c, saveForm, func() (interface{}, error) {
		uid := uint(helper.InterfaceToInt(app.Request(c).GetCtxVal("uid")))
		isSuper := app.Request(c).GetCtxVal("is_super")

		if saveForm.Data.Id == uint(0) {
			newId, err := service.Admin().Create(saveForm.Data)

			if isSuper == 10 {
				saveForm.Data.Pid = newId
			} else {
				saveForm.Data.Pid = uid
				//非超级管理员，添加账号时，自动加上主账号作为前缀
				userInfo := service.Admin().GetUser(uid)
				if !strings.HasPrefix(saveForm.Data.Username, userInfo.Username+".") {
					saveForm.Data.Username = userInfo.Username + "." + saveForm.Data.Username
				}
			}
			if err != nil {
				return newId, err
			}

			return service.Base().Update(saveForm.Data)
		} else {
			if isSuper != 10 {
				//非超级管理员，添加账号时，自动加上主账号作为前缀
				userInfo := service.Admin().GetUser(uid)
				parentInfo := service.Admin().GetUser(userInfo.Pid)
				if uid != userInfo.Pid && !strings.HasPrefix(saveForm.Data.Username, parentInfo.Username+".") {
					saveForm.Data.Username = parentInfo.Username + "." + saveForm.Data.Username
				}
			}

			return service.Base().Update(saveForm.Data)
		}
	})
}

// Login 用户登录
func (ctrl *Admin) Login(c *gin.Context) {
	dataLoginForm := &model.DataLoginForm{}
	Common(c, dataLoginForm, func() (interface{}, error) {
		dataLoginForm.Data.LastLoginIp = c.ClientIP()
		dataLoginForm.Data.LastLoginTime = model.JSONTime{Time: time.Now()}
		return service.Admin().Login(dataLoginForm.Data)
	})
}

// Logout 退出登录
func (ctrl *Admin) Logout(c *gin.Context) {
	dataLogoutForm := &model.DataLogoutForm{}
	Common(c, dataLogoutForm, func() (interface{}, error) {
		return service.Admin().Logout(dataLogoutForm.Data)
	})
}

// AddRole 给用户添加角色
func (ctrl *Admin) AddRole(c *gin.Context) {
	dataAddRoleForm := &model.DataAddRoleForm{}
	Common(c, dataAddRoleForm, func() (interface{}, error) {
		if len(dataAddRoleForm.Data.RoleIdList) == 0 {
			//如果角色列表为空，即表示用户没有角色，则清空用户所有角色
			return service.Auth().DelAllRolesForUser(dataAddRoleForm.Data.Username)
		} else {
			return service.Auth().AddRolesForUser(dataAddRoleForm.Data.Username, dataAddRoleForm.Data.RoleIdList)
		}
	})
}

// DelRole 删除用户的角色
func (ctrl *Admin) DelRole(c *gin.Context) {
	dataDelRoleForm := &model.DataDelRoleForm{}
	Common(c, dataDelRoleForm, func() (interface{}, error) {
		return service.Auth().DelRolesForUser(dataDelRoleForm.Data.Username, dataDelRoleForm.Data.RoleName)
	})
}

// AddPermission 分配角色的权限
func (ctrl *Admin) AddPermission(c *gin.Context) {
	dataPermissionForm := &model.DataPermissionForm{}
	Common(c, dataPermissionForm, func() (interface{}, error) {
		return service.Auth().AddPermission(dataPermissionForm.Data.RoleName, dataPermissionForm.Data.ActionId)
	})
}

// DelPermission 删除角色的权限
func (ctrl *Admin) DelPermission(c *gin.Context) {
	dataPermissionForm := &model.DataPermissionForm{}
	Common(c, dataPermissionForm, func() (interface{}, error) {
		return service.Auth().DelPermission(dataPermissionForm.Data.RoleName, dataPermissionForm.Data.ActionId)
	})
}

// GetPermission 获取角色的权限
func (ctrl *Admin) GetPermission(c *gin.Context) {
	dataRoleForm := &model.DataRoleForm{}
	Common(c, dataRoleForm, func() (interface{}, error) {
		isSuper := app.Request(c).GetCtxVal("is_super")
		return service.Auth().GetPermissions(dataRoleForm.Data.RoleName, isSuper)
	})
}

// GetAllRoles 获取所有角色
func (ctrl *Admin) GetAllRoles(c *gin.Context) {
	dataRoleForm := &model.DataRoleForm{}
	Common(c, dataRoleForm, func() (interface{}, error) {
		isSuper := app.Request(c).GetCtxVal("is_super")
		res := service.Auth().GetAllRoles(dataRoleForm.Data.RoleName, isSuper)
		return res, nil
	})
}

// GetRoles 获取用户的所有角色
func (ctrl *Admin) GetRoles(c *gin.Context) {
	dataUsernameForm := &model.DataUsernameForm{}
	Common(c, dataUsernameForm, func() (interface{}, error) {
		return service.Auth().GetRoles(dataUsernameForm.Data.Username)
	})
}

// SetUserPermission 分配用户的权限
func (ctrl *Admin) SetUserPermission(c *gin.Context) {
	dataUserPermissionForm := &model.DataUserPermissionForm{}
	Common(c, dataUserPermissionForm, func() (interface{}, error) {
		return service.Auth().AddPermission(
			dataUserPermissionForm.Data.Username,
			dataUserPermissionForm.Data.ActionId,
		)
	})
}

// GetUserPermission 获取用户的权限
func (ctrl *Admin) GetUserPermission(c *gin.Context) {
	dataUsernameForm := &model.DataUsernameForm{}
	Common(c, dataUsernameForm, func() (interface{}, error) {
		return service.Auth().GetPermissions(
			dataUsernameForm.Data.Username,
			nil,
		)
	})
}
