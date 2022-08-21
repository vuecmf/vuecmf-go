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
	"github.com/vuecmf/vuecmf-go/app/vuecmf/model"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/service"
)

type Base struct {
	TableName string
	Model     interface{}
}

// common 控制器公共入口方法
func common(c *gin.Context, formParams interface{}, fun func() (interface{}, error)) {
	defer func() {
		if err := recover(); err != nil {
			app.Response(c).SendFailure("请求异常", err)
		}
	}()

	err := app.Request(c).Input("post", formParams)

	if err != nil {
		app.Response(c).SendFailure("请求失败："+model.GetError(err, formParams), nil)
		return
	}

	list, err := fun()

	if err != nil {
		app.Response(c).SendFailure("请求失败："+err.Error(), nil)
		return
	}

	app.Response(c).SendSuccess("请求成功", list)
}

// Detail 根据ID获取详情
func (ctrl *Base) Detail(c *gin.Context) {
	data := &model.DataIdForm{}
	common(c, data, func() (interface{}, error) {
		err := service.Base().Detail(data.Data.Id, ctrl.Model)
		return ctrl.Model, err
	})
}

// Dropdown 下拉列表数据
func (ctrl *Base) Dropdown(c *gin.Context) {
	data := &model.DataDropdownForm{}
	common(c, data, func() (interface{}, error) {
		return service.Base().Dropdown(data.Data, ctrl.TableName)
	})
}
