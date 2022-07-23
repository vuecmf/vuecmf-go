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
	"github.com/vuecmf/vuecmf-go/app/route"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/form"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/helper"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/service"
)

type Admin struct {
	*base
}

func init() {
	route.Register(&Admin{}, "GET|POST", "vuecmf")
}

// Index 列表页
func (ctrl *Admin) Index(c *gin.Context) {
	commonIndex(c, func(listParams *helper.DataListParams) interface{} {
		return service.Admin().List(listParams)
	})
}

func (ctrl *Admin) Login(c *gin.Context) {
	loginForm := form.LoginForm{Username: "aaa", Password: "123456"}

	//获取输入内容
	app.Request(c).Input("post", &loginForm)

	//输出内容
	app.Response(c).SendSuccess("提交成功", loginForm)

}
