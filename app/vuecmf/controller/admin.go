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
	"strconv"
	"strings"
	"sync"
	"time"
)

// AdminController 管理员管理
type AdminController struct {
	BaseController
	Svc *service.AdminService
}

var adminController *AdminController
var adminCtrlOnce sync.Once

// Admin 获取admin控制器实例
func Admin() *AdminController {
	adminCtrlOnce.Do(func() {
		adminController = &AdminController{
			Svc: service.Admin(),
		}
	})
	return adminController
}

// Action admin控制器入口
func (ctrl AdminController) Action() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var res any

		switch GetActionName(c) {
		case "":
			res, err = ctrl.index(c)
		case "save":
			res, err = ctrl.save(c)
		case "login": // 用户登录
			res, err = ctrl.login(c)
		case "logout": // 退出登录
			res, err = ctrl.logout(c)
		case "add_role":
			res, err = ctrl.addRole(c)
		case "del_role":
			res, err = ctrl.delRole(c)
		case "add_permission":
			res, err = ctrl.addPermission(c)
		case "del_permission":
			res, err = ctrl.delPermission(c)
		case "get_permission":
			res, err = ctrl.getPermission(c)
		case "get_all_roles":
			res, err = ctrl.getAllRoles(c)
		case "get_roles":
			res, err = ctrl.getRoles(c)
		case "set_user_permission":
			res, err = ctrl.setUserPermission(c)
		case "get_user_permission":
			res, err = ctrl.getUserPermission(c)
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
func (ctrl AdminController) index(c *gin.Context) (any, error) {
	isSuper := MGet(c, "is_super").(uint16)
	var params *helper.DataListParams
	err := Post(c, &params)
	if err != nil {
		return nil, err
	}

	if isSuper != 10 {
		//非超级管理员，只能看自己和自己创建的账号
		uid := strconv.Itoa(int(MGet(c, "uid").(uint)))
		condition := make(map[string]interface{})
		condition["pid"] = uid
		params.Data.Filter["id"] = uid
		params.Data.Filter["or"] = condition
	}
	return ctrl.Svc.CommonList(params, isSuper)
}

// save 新增/更新 单条数据
func (ctrl AdminController) save(c *gin.Context) (int64, error) {
	var res int64
	var err error
	uid := MGet(c, "uid").(uint)
	isSuper := MGet(c, "is_super").(uint16)

	var params *model.DataAdminForm
	err = Post(c, &params)
	if err != nil {
		return 0, err
	}

	if params.Data.Id == uint(0) {
		newId, err := ctrl.Svc.Create(params.Data)

		if isSuper == 10 {
			params.Data.Pid = newId
		} else {
			params.Data.Pid = uid
			//非超级管理员，添加账号时，自动加上主账号作为前缀
			userInfo := ctrl.Svc.GetUser(uid)
			if !strings.HasPrefix(params.Data.Username, userInfo.Username+".") {
				params.Data.Username = userInfo.Username + "." + params.Data.Username
			}
		}

		if err == nil {
			res, err = ctrl.Svc.Update(params.Data)
		}

	} else {
		if isSuper != 10 {
			//非超级管理员，添加账号时，自动加上主账号作为前缀
			userInfo := ctrl.Svc.GetUser(uid)
			parentInfo := ctrl.Svc.GetUser(userInfo.Pid)
			if uid != userInfo.Pid && !strings.HasPrefix(params.Data.Username, parentInfo.Username+".") {
				params.Data.Username = parentInfo.Username + "." + params.Data.Username
			}
		}
		res, err = ctrl.Svc.Update(params.Data)
	}
	return res, err
}

// login 用户登录
func (ctrl AdminController) login(c *gin.Context) (any, error) {
	var params *model.DataLoginForm
	err := Post(c, &params)
	if err != nil {
		return nil, err
	}

	params.Data.LastLoginIp = c.ClientIP()
	params.Data.LastLoginTime = model.JSONTime{Time: time.Now()}
	return ctrl.Svc.Login(params.Data)
}

// logout 退出登录
func (ctrl AdminController) logout(c *gin.Context) (bool, error) {
	var params *model.DataLogoutForm
	err := Post(c, &params)
	if err != nil {
		return false, err
	}
	return ctrl.Svc.Logout(params.Data)
}

// addRole 给用户添加角色
func (ctrl AdminController) addRole(c *gin.Context) (bool, error) {
	var params *model.DataAddRoleForm
	err := Post(c, &params)
	if err != nil {
		return false, err
	}

	if len(params.Data.RoleIdList) == 0 {
		//如果角色列表为空，即表示用户没有角色，则清空用户所有角色
		return service.Auth().DelAllRolesForUser(params.Data.Username)
	} else {
		return service.Auth().AddRolesForUser(params.Data.Username, params.Data.RoleIdList)
	}
}

// delRole 删除用户的角色
func (ctrl AdminController) delRole(c *gin.Context) (bool, error) {
	var params *model.DataDelRoleForm
	err := Post(c, &params)
	if err != nil {
		return false, err
	}

	return service.Auth().DelRolesForUser(params.Data.Username, params.Data.RoleName)
}

// addPermission 分配角色的权限
func (ctrl AdminController) addPermission(c *gin.Context) (bool, error) {
	var params *model.DataPermissionForm
	err := Post(c, &params)
	if err != nil {
		return false, err
	}

	return service.Auth().AddPermission(params.Data.RoleName, params.Data.ActionId)
}

// delPermission 删除角色的权限
func (ctrl AdminController) delPermission(c *gin.Context) (bool, error) {
	var params *model.DataPermissionForm
	err := Post(c, &params)
	if err != nil {
		return false, err
	}

	return service.Auth().DelPermission(params.Data.RoleName, params.Data.ActionId)

}

// getPermission 获取角色的权限
func (ctrl AdminController) getPermission(c *gin.Context) (map[string][]string, error) {
	isSuper := MGet(c, "is_super").(uint16)
	var params *model.DataRoleForm
	err := Post(c, &params)
	if err != nil {
		return nil, err
	}
	return service.Auth().GetPermissions(params.Data.RoleName, isSuper)

}

// getAllRoles 获取所有角色
func (ctrl AdminController) getAllRoles(c *gin.Context) (any, error) {
	isSuper := MGet(c, "is_super").(uint16)
	var params *model.DataRoleForm
	err := Post(c, &params)
	if err != nil {
		return nil, err
	}
	res := service.Auth().GetAllRoles(params.Data.RoleName, isSuper)
	return res, nil
}

// getRoles 获取用户的所有角色
func (ctrl AdminController) getRoles(c *gin.Context) ([]int, error) {
	var params *model.DataUsernameForm
	err := Post(c, &params)
	if err != nil {
		return nil, err
	}
	return service.Auth().GetRoles(params.Data.Username)
}

// setUserPermission 分配用户的权限
func (ctrl AdminController) setUserPermission(c *gin.Context) (bool, error) {
	var params *model.DataUserPermissionForm
	err := Post(c, &params)
	if err != nil {
		return false, err
	}
	return service.Auth().AddPermission(params.Data.Username, params.Data.ActionId)
}

// getUserPermission 获取用户的权限
func (ctrl AdminController) getUserPermission(c *gin.Context) (map[string][]string, error) {
	var params *model.DataUsernameForm
	err := Post(c, &params)
	if err != nil {
		return nil, err
	}
	return service.Auth().GetPermissions(params.Data.Username, 0)
}
