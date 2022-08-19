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
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vuecmf/vuecmf-go/app"
	"github.com/vuecmf/vuecmf-go/app/route"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/helper"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/model"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/service"
)

type Admin struct {
}

func init() {
	route.Register(&Admin{}, "GET|POST", "vuecmf")
}

// Index 列表页
func (ctrl *Admin) Index(c *gin.Context) {
	app.Cache().Set("hello", "123456")

	listParams := &helper.DataListParams{}
	common(c, listParams, func() (interface{}, error) {
		return service.Admin().List(listParams)
	})
}

// Save 新增/更新 单条数据
func (ctrl *Admin) Save(c *gin.Context) {
	data := &model.DataAdminForm{}

	common(c, data, func() (interface{}, error) {
		if data.Data.Id == 0 {
			return service.Admin().Create(data.Data)
		} else {
			return service.Admin().Update(data.Data)
		}
	})
}


// Saveall 批量添加多条数据
func (ctrl *Admin) Saveall(c *gin.Context) {
	data := &model.DataBatchForm{}
	common(c, data, func() (interface{}, error) {
		var dataBatch []model.Admin
		err := json.Unmarshal([]byte(data.Data), &dataBatch)
		if err != nil {
			return nil, err
		}
		return service.Admin().Create(dataBatch)
	})
}

// Detail 根据ID获取详情
func (ctrl *Admin) Detail(c *gin.Context) {
	data := &model.DataIdForm{}
	common(c, data, func() (interface{}, error) {
		var result model.Admin
		err := service.Admin().Detail(data.Id, &result)
		return result, err
	})
}

// Delete 根据ID删除单条数据
func (ctrl *Admin) Delete(c *gin.Context) {
	data := &model.DataIdForm{}
	common(c, data, func() (interface{}, error) {
		return service.Admin().Delete(data.Id, &model.Admin{})
	})
}

// Deletebatch 根据ID列表批量删除多条数据
func (ctrl *Admin) Deletebatch(c *gin.Context) {
	data := &model.DataIdListForm{}
	common(c, data, func() (interface{}, error) {
		return service.Admin().DeleteBatch(data.IdList, &model.Admin{})
	})
}


func (ctrl *Admin) Dropdown(c *gin.Context) {

}




func (ctrl *Admin) Login(c *gin.Context) {
	login := &model.LoginForm{
		Username: "haha",
		Password: "123456",
	}

	_ = app.Cache().Set("user", login)

	var loginRes model.LoginForm
	_ = app.Cache().Get("user", &loginRes)

	fmt.Println("loginRes = ", loginRes)

	var str string
	app.Cache().Get("hello", &str)
	fmt.Println("str cache = ", str)


	loginForm := &model.LoginForm{}
	common(c, loginForm, func() (interface{}, error) {
		fmt.Println(c.Get("bbb"))

		return loginForm, nil
	})

	//获取输入内容
	//app.Request(c).Input("post", &loginForm)

	//fmt.Println(c.Get("bbb"))

	//输出内容
	//app.Response(c).SendSuccess("提交成功", loginForm)

}
