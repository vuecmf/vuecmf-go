//+----------------------------------------------------------------------
// | Copyright (c) 2024 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: tulihua2004@126.com
// +----------------------------------------------------------------------

package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/vuecmf/vuecmf-go/v3/app/vuecmf/helper"
	"github.com/vuecmf/vuecmf-go/v3/app/vuecmf/model"
	"github.com/vuecmf/vuecmf-go/v3/app/vuecmf/service"
	"sync"
)

// RolesController 角色管理
type RolesController struct {
	BaseController
	Svc *service.RolesService
}

var rolesController *RolesController
var rolesCtrlOnce sync.Once

// Roles 获取控制器实例
func Roles() *RolesController {
	rolesCtrlOnce.Do(func() {
		rolesController = &RolesController{
			Svc: service.Roles(),
		}
	})
	return rolesController
}

// Action 控制器入口
func (ctrl RolesController) Action() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var res any

		switch GetActionName(c) {
		case "":
			res, err = ctrl.index(c)
		case "save":
			res, err = ctrl.save(c)
		case "addUsers":
			res, err = ctrl.addUsers(c)
		case "delUsers":
			res, err = ctrl.delUsers(c)
		case "addPermission":
			res, err = ctrl.addPermission(c)
		case "delPermission":
			res, err = ctrl.delPermission(c)
		case "getUsers":
			res, err = ctrl.getUsers(c)
		case "getPermission":
			res, err = ctrl.getPermission(c)
		case "getAllUsers":
			res, err = ctrl.getAllUsers(c)
		default:
			res, err = ctrl.BaseController.Action(c, ctrl.Svc.BaseService)
		}

		if err != nil {
			c.Set("error", err)
		} else {
			c.Set("result", res)
		}

		c.Next()
	}
}

// index 列表页
func (ctrl RolesController) index(c *gin.Context) (any, error) {
	var params *helper.DataListParams
	err := Post(c, &params)
	if err != nil {
		return nil, err
	}
	return ctrl.Svc.List(params)
}

// save 新增/更新 单条数据
func (ctrl RolesController) save(c *gin.Context) (int64, error) {
	var params *model.DataRolesForm
	err := Post(c, &params)
	if err != nil {
		return 0, err
	}
	if params.Data.Id == uint(0) {
		return ctrl.Svc.Create(params.Data)
	} else {
		return ctrl.Svc.Update(params.Data)
	}
}

// addUsers 给角色分配用户
func (ctrl RolesController) addUsers(c *gin.Context) (any, error) {
	var params *model.DataRoleUsersForm
	err := Post(c, &params)
	if err != nil {
		return nil, err
	}
	return ctrl.Svc.AddUsers(params.Data.RoleName, params.Data.UseridList)
}

// delUsers 删除角色下的用户
func (ctrl RolesController) delUsers(c *gin.Context) (bool, error) {
	var params *model.DataRoleUsersForm
	err := Post(c, &params)
	if err != nil {
		return false, err
	}
	return service.Auth().DelUsersForRole(params.Data.RoleName, params.Data.UseridList)
}

// addPermission 给角色分配权限项
func (ctrl RolesController) addPermission(c *gin.Context) (bool, error) {
	var params *model.DataPermissionForm
	err := Post(c, &params)
	if err != nil {
		return false, err
	}
	return service.Auth().AddPermission(params.Data.RoleName, params.Data.ActionId)
}

// delPermission 删除角色的权限项
func (ctrl RolesController) delPermission(c *gin.Context) (bool, error) {
	var params *model.DataPermissionForm
	err := Post(c, &params)
	if err != nil {
		return false, err
	}
	return service.Auth().DelPermission(params.Data.RoleName, params.Data.ActionId)
}

// getUsers 获取角色下所有用户的ID
func (ctrl RolesController) getUsers(c *gin.Context) ([]int, error) {
	var params *model.DataRoleForm
	err := Post(c, &params)
	if err != nil {
		return nil, err
	}
	return ctrl.Svc.GetUsers(params.Data.RoleName)
}

// getPermission 获取角色的所有权限项
func (ctrl RolesController) getPermission(c *gin.Context) (map[string][]string, error) {
	var params *model.DataRoleForm
	err := Post(c, &params)
	if err != nil {
		return nil, err
	}
	return service.Auth().GetPermissions(params.Data.RoleName, 0)
}

// getAllUsers 获取所有用户
func (ctrl RolesController) getAllUsers(c *gin.Context) (any, error) {
	isSuper := MGet(c, "is_super").(uint16)
	var params *model.DataUsernameForm
	err := Post(c, &params)
	if err != nil {
		return nil, err
	}
	return ctrl.Svc.GetAllUsers(params.Data.Username, isSuper)
}
