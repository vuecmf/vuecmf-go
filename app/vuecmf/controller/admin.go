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
	listParams := &helper.DataListParams{}
	common(c, listParams, func() (interface{}, error) {
		return service.Admin().List(listParams)
	})
}

// Save 保存单条数据
func (ctrl *Admin) Save(c *gin.Context) {
	data := &model.DataAdminForm{}

	app.Cache().Set("hello", []byte("world 123"))

	common(c, data, func() (interface{}, error) {
		if data.Data.Id == 0 {
			return service.Admin().Create(data.Data)
		} else {
			res, _ := app.Cache().Get("hello")
			fmt.Println(string(res))
			return service.Admin().Update(data.Data)
		}
	})
}

func (ctrl *Admin) Login(c *gin.Context) {
	login := &model.LoginForm{
		Username: "haha",
		Password: "123456",
	}

	loginBt,_ := json.Marshal(login)

	app.Cache().Set("user", loginBt)

	userBt, _ := app.Cache().Get("user")

	var login2 model.LoginForm
	json.Unmarshal(userBt, &login2)
	fmt.Println("login2 = ",login2)



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
