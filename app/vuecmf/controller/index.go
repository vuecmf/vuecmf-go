//+----------------------------------------------------------------------
// | Copyright (c) 2023 http://www.vuecmf.com All rights reserved.
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
)

// Index 首页
type Index struct {
}

func init() {
	route.Register(&Index{}, "GET|POST", "vuecmf")
}

func (ctrl *Index) Index(c *gin.Context) {
	app.Response(c).SendHtml("vuecmf/index/index.html", gin.H{
		"welcome": "Welcome to VueCMF V2.5",
	})
}

func (ctrl *Index) Success(c *gin.Context) {
	type user struct {
		Name string
		Age  int
	}

	app.Response(c).SendSuccess("success", &user{
		Name: "Zhang san",
		Age:  18,
	}, 0)
}

func (ctrl *Index) Fail(c *gin.Context) {
	app.Response(c).SendFailure("fail", "", 500)
}
