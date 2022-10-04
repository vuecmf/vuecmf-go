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
	"github.com/vuecmf/vuecmf-go/app/vuecmf/service"
)

// Make 代码生成控制器
type Make struct {
}

func init() {
	route.Register(&Make{}, "GET|POST", "vuecmf")
}

// Model 生成模型代码文件
func (ctrl *Make) Model(c *gin.Context) {
	tableName := app.Request(c).Get("table_name")
	err := service.Make().Model(tableName, "vuecmf")

	if err != nil {
		app.Response(c).SendFailure("模型代码文件生成失败!"+err.Error(), nil)
	} else {
		app.Response(c).SendSuccess("模型代码文件生成成功", nil)
	}
}

// Service 生成服务代码文件
func (ctrl *Make) Service(c *gin.Context) {
	tableName := app.Request(c).Get("table_name")
	err := service.Make().Service(tableName, "vuecmf")

	if err != nil {
		app.Response(c).SendFailure("服务代码文件生成失败!"+err.Error(), nil)
	} else {
		app.Response(c).SendSuccess("服务代码文件生成成功", nil)
	}
}

// Controller 生成服务代码文件
func (ctrl *Make) Controller(c *gin.Context) {
	tableName := app.Request(c).Get("table_name")
	err := service.Make().Controller(tableName, "vuecmf")

	if err != nil {
		app.Response(c).SendFailure("控制器代码文件生成失败!"+err.Error(), nil)
	} else {
		app.Response(c).SendSuccess("控制器代码文件生成成功", nil)
	}
}
