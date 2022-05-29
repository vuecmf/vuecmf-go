package router

import (
	"github.com/gin-gonic/gin"
	"github.com/vuecmf/vuecmf-go/app/middleware"
	"github.com/vuecmf/vuecmf-go/app/router/vuecmf"
)

func Load() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CheckAuth)

	//加入vuecmf应用的路由配置
	vuecmf.Router(r)
	//在此加入路由配置


	return r
}