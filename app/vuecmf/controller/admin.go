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
	"github.com/vuecmf/vuecmf-go/app/route"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/helper"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/model"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/service"
)

type Admin struct {
	Base
}

func init() {
	admin := &Admin{}
	admin.TableName = "admin"
	admin.Model = &model.Admin{}
	admin.listRes = []model.Admin{}
	route.Register(admin, "GET|POST", "vuecmf")
}

// Index 列表页
func (ctrl *Admin) Index(c *gin.Context) {
	listParams := &helper.DataListParams{}
	common(c, listParams, func() (interface{}, error) {
		var result []model.Admin
		return service.Base().CommonList(result, ctrl.TableName, listParams)
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

func (ctrl *Admin) Login(c *gin.Context) {
	loginForm := &model.LoginForm{}
	common(c, loginForm, func() (interface{}, error) {
		//fmt.Println(c.Get("bbb"))

		return loginForm, nil
	})

	//获取输入内容
	//app.Request(c).Input("post", &loginForm)

	//fmt.Println(c.Get("bbb"))

	//输出内容
	//app.Response(c).SendSuccess("提交成功", loginForm)

}
