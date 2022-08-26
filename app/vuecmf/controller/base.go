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
	"github.com/vuecmf/vuecmf-go/app/vuecmf/helper"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/model"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/service"
)

type Base struct {
	TableName string
	Model     interface{}
	listRes   interface{}
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

// Index 列表页
func (ctrl *Base) Index(c *gin.Context) {
	listParams := &helper.DataListParams{}
	common(c, listParams, func() (interface{}, error) {
		//var result []model.Admin
		return service.Base().CommonList(ctrl.listRes, ctrl.TableName, listParams)
	})
}


// Detail 根据ID获取详情
func (ctrl *Base) Detail(c *gin.Context) {
	data := &model.DataIdForm{}
	common(c, data, func() (interface{}, error) {
		err := service.Base().Detail(data.Data.Id, ctrl.Model)
		return ctrl.Model, err
	})
}

// Delete 根据ID删除单条数据
func (ctrl *Base) Delete(c *gin.Context) {
	data := &model.DataIdForm{}
	common(c, data, func() (interface{}, error) {
		return service.Base().Delete(data.Data.Id, ctrl.Model)
	})
}

// Deletebatch 根据ID列表批量删除多条数据
func (ctrl *Admin) Deletebatch(c *gin.Context) {
	data := &model.DataIdListForm{}
	common(c, data, func() (interface{}, error) {
		return service.Base().DeleteBatch(data.Data.IdList, ctrl.Model)
	})
}

// Dropdown 下拉列表数据
func (ctrl *Base) Dropdown(c *gin.Context) {
	data := &model.DataDropdownForm{}
	common(c, data, func() (interface{}, error) {
		return service.Base().Dropdown(data.Data, ctrl.TableName)
	})
}


func (ctrl *Base) Cache(c *gin.Context){
	//app.Cache().Set("hello", "123456")

	/*login := &model.LoginForm{
		Username: "haha",
		Password: "123456",
	}

	_ = app.Cache().Set("user", login)

	var loginRes model.LoginForm
	_ = app.Cache().Get("user", &loginRes)

	fmt.Println("loginRes = ", loginRes)

	var str string
	app.Cache().Get("hello", &str)
	fmt.Println("str cache = ", str)*/
}
