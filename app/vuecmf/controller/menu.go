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
)

// Menu 菜单管理
type Menu struct {
	Base
}

func init() {
	menu := &Menu{}
	menu.TableName = "menu"
	menu.Model = &model.Menu{}
	menu.ListData = &[]model.Menu{}
	menu.FilterFields = []string{"title", "icon"}

	route.Register(menu, "POST", "vuecmf")
}

// Index 列表页
func (ctrl *Menu) Index(c *gin.Context) {
	listParams := &helper.DataListParams{}
	Common(c, listParams, func() (interface{}, error) {
		return service.Menu().List(listParams)
	})
}

// Save 新增/更新 单条数据
func (ctrl *Menu) Save(c *gin.Context) {
	saveForm := &model.DataMenuForm{}
	Common(c, saveForm, func() (interface{}, error) {
		if saveForm.Data.Id == uint(0) {
			return service.Menu().Create(saveForm.Data)
		} else {
			return service.Menu().Update(saveForm.Data)
		}
	})
}

// Nav 获取用户的导航菜单列表
func (ctrl *Menu) Nav(c *gin.Context) {
	dataUsernameForm := &model.DataUsernameForm{}
	Common(c, dataUsernameForm, func() (interface{}, error) {
		isSuper := app.Request(c).GetCtxVal("is_super")
		return service.Menu().Nav(dataUsernameForm.Data.Username, isSuper)
	})
}
