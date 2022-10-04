package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/vuecmf/vuecmf-go/app"
	"github.com/vuecmf/vuecmf-go/app/route"
)

type Index struct {
}

func init() {
	route.Register(&Index{}, "GET|POST", "vuecmf")
}

func (ctrl *Index) Index(c *gin.Context) {
	app.Response(c).SendHtml("vuecmf/index/index.html", gin.H{
		"welcome": "Welcome to VueCMF V2.",
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
