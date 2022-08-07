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
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vuecmf/vuecmf-go/app/route"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/form"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/helper"
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

func (ctrl *Admin) Login(c *gin.Context) {
	loginForm := &form.LoginForm{}
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
