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
	"sync"
)

// IndexController 首页
type IndexController struct {
}

type user struct {
	Name string
	Age  int
}

var indexController *IndexController
var indexCtrlOnce sync.Once

func Index() *IndexController {
	indexCtrlOnce.Do(func() {
		indexController = &IndexController{}
	})
	return indexController
}

// Before 路由前置拦截器
func (ctrl IndexController) Before() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO 加入前置业务处理逻辑

		c.Next()
	}
}

// After 路由后置拦截器
func (ctrl IndexController) After() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO 加入后置业务处理逻辑
	}
}

// Action 控制器入口
func (ctrl IndexController) Action() gin.HandlerFunc {
	return func(c *gin.Context) {
		switch GetActionName(c) {
		case "success":
			ctrl.success(c)
		case "fail":
			ctrl.fail(c)
		default:
			ctrl.index(c)
		}

		c.Next()

	}
}

func (ctrl IndexController) index(c *gin.Context) {
	app.Response(c).SendHtml("vuecmf/index/index.html", gin.H{
		"welcome": "Welcome to VueCMF V3",
	})
}

func (ctrl IndexController) success(c *gin.Context) {
	app.Response(c).SendSuccess("success", &user{
		Name: "Zhang san",
		Age:  18,
	}, 0)
}

func (ctrl IndexController) fail(c *gin.Context) {
	app.Response(c).SendFailure("fail", "", 500)
}
