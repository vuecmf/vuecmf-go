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
	*base
}

func init() {
	route.Register(&Make{}, "GET|POST", "vuecmf")
}

// Model 生成模型方法
func (ctrl *Make) Model(c *gin.Context) {
	tableName := app.Request(c).Get("table_name")
	makeRes := service.Make().Model(tableName)

	if makeRes {
		app.Response(c).SendSuccess("模型生成成功", nil)
	} else {
		app.Response(c).SendFailure("模型生成失败", nil)
	}

}
