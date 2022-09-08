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
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/vuecmf/vuecmf-go/app"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/helper"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/model"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/service"
	"reflect"
)

type Base struct {
	TableName string      //表名称
	Model     interface{} //表对应的模型实例
	listData  interface{} //存储列表结果
	saveForm  interface{} //保存单条数据的form
	filterFields []string //支持模糊查询的字段
	AppName   string      //当前应用标识
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
		return service.Base().CommonList(ctrl.listData, ctrl.TableName, ctrl.filterFields, listParams)
	})
}

// Save 新增/更新 单条数据
func (ctrl *Admin) Save(c *gin.Context) {
	common(c, ctrl.saveForm, func() (interface{}, error) {
		data := reflect.ValueOf(ctrl.saveForm).Elem().FieldByName("Data").Interface()
		id := reflect.ValueOf(data).Elem().FieldByName("Id").Interface()

		if id == uint(0) {
			return service.Base().Create(data)
		} else {
			return service.Base().Update(data)
		}
	})
}

// Saveall 批量添加多条数据
func (ctrl *Admin) Saveall(c *gin.Context) {
	data := &model.DataBatchForm{}
	common(c, data, func() (interface{}, error) {
		//var dataBatch []model.Admin
		err := json.Unmarshal([]byte(data.Data), &ctrl.listData)
		if err != nil {
			return nil, err
		}
		return service.Base().Create(ctrl.listData)
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

func (ctrl *Base) Cache(c *gin.Context) {
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
