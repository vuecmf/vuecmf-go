//+----------------------------------------------------------------------
// | Copyright (c) 2024 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: tulihua2004@126.com
// +----------------------------------------------------------------------

package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/vuecmf/vuecmf-go/v3/app"
	"github.com/vuecmf/vuecmf-go/v3/app/vuecmf/service"
	"sync"
)

// MakeController 代码生成控制器
type MakeController struct {
}

var makeController *MakeController
var makeCtrlOnce sync.Once

// Make 获取代码生成控制器实例
func Make() *MakeController {
	makeCtrlOnce.Do(func() {
		makeController = &MakeController{}
	})
	return makeController
}

// Before 路由前置拦截器
func (ctrl MakeController) Before() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO 加入前置业务处理逻辑

		c.Next()
	}
}

// After 路由后置拦截器
func (ctrl MakeController) After() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO 加入后置业务处理逻辑
	}
}

// Action 动作
func (ctrl MakeController) Action() gin.HandlerFunc {
	return func(c *gin.Context) {
		actionName := c.Param("action")
		switch actionName {
		case "model":
			ctrl.model(c)
		case "service":
			ctrl.service(c)
		case "controller":
			ctrl.controller(c)
		default:
			app.Response(c).SendFailure("无效的action", nil)
		}
	}
}

// model 生成模型代码文件
func (ctrl MakeController) model(c *gin.Context) {
	tableName := app.Request(c).Get("table_name")
	err := service.Make().Model(tableName, "vuecmf")

	if err != nil {
		app.Response(c).SendFailure("模型代码文件生成失败!"+err.Error(), nil)
	} else {
		app.Response(c).SendSuccess("模型代码文件生成成功", nil)
	}
}

// service 生成服务代码文件
func (ctrl MakeController) service(c *gin.Context) {
	tableName := app.Request(c).Get("table_name")
	err := service.Make().Service(tableName, "vuecmf")

	if err != nil {
		app.Response(c).SendFailure("服务代码文件生成失败!"+err.Error(), nil)
	} else {
		app.Response(c).SendSuccess("服务代码文件生成成功", nil)
	}
}

// controller 生成服务代码文件
func (ctrl MakeController) controller(c *gin.Context) {
	tableName := app.Request(c).Get("table_name")
	err := service.Make().Controller(tableName, "vuecmf")

	if err != nil {
		app.Response(c).SendFailure("控制器代码文件生成失败!"+err.Error(), nil)
	} else {
		app.Response(c).SendSuccess("控制器代码文件生成成功", nil)
	}
}
