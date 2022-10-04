package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/vuecmf/vuecmf-go/app/vuecmf/controller"

	"github.com/vuecmf/vuecmf-go/app/route"
	"log"
)

func main() {
	engine := gin.Default()

	//初始化路由
	route.InitRoute(engine)

	err := engine.Run(":8080")
	if err != nil {
		log.Fatal("服务启动失败！", err)
	}

}
