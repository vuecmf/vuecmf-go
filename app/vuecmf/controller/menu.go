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

// MenuController 菜单管理
type MenuController struct {
	BaseController
	Svc *service.MenuService
}

var menuController *MenuController
var menuCtrlOnce sync.Once

// Menu 获取Menu控制器实例
func Menu() *MenuController {
	menuCtrlOnce.Do(func() {
		menuController = &MenuController{
			Svc: service.Menu(),
		}
	})
	return menuController
}

// Action 控制器入口
func (ctrl MenuController) Action() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var res any

		switch GetActionName(c) {
		case "":
			res, err = ctrl.index(c)
		case "save":
			res, err = ctrl.save(c)
		case "nav":
			res, err = ctrl.nav(c)
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
func (ctrl MenuController) index(c *gin.Context) (any, error) {
	var params *helper.DataListParams
	err := Post(c, &params)
	if err != nil {
		return nil, err
	}
	return ctrl.Svc.List(params)
}

// save 新增/更新 单条数据
func (ctrl MenuController) save(c *gin.Context) (int64, error) {
	var params *model.DataMenuForm
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

// nav 获取用户的导航菜单列表
func (ctrl MenuController) nav(c *gin.Context) (any, error) {
	isSuper := MGet(c, "is_super").(uint16)
	var params *model.DataUsernameForm
	err := Post(c, &params)
	if err != nil {
		return nil, err
	}
	return ctrl.Svc.Nav(params.Data.Username, isSuper)
}
