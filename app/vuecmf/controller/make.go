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

func init(){
	route.Register(&Make{}, "GET|POST", "vuecmf")
}

// Model 生成模型方法
func (ctrl *Make) Model(c *gin.Context) {
	req := app.Request{Context: c}
	tableName := req.Get("table_name")

	makeService := service.MakeService{}
	makeRes := makeService.Model(tableName)

	resp := app.Response{Context: c}

	if makeRes {
		resp.SendSuccess("模型生成成功", nil)
	} else {
		resp.SendFailure("模型生成失败", nil)
	}

}
