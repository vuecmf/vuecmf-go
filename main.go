package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vuecmf/vuecmf-go/app/route"
	_ "github.com/vuecmf/vuecmf-go/app/vuecmf/controller"
)


func main() {
	engine := gin.Default()

	//加入权限验证
	//engine.Use(middleware.CheckAuth)

	//初始路由
	route.InitRoute(engine)

	err := engine.Run(":8080")
	if err != nil {
		fmt.Println("服务启动失败！", err)
	}
}
