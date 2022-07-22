package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/vuecmf/vuecmf-go/app"
	"github.com/vuecmf/vuecmf-go/app/route"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/form"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/helper"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/service"
)

type Admin struct {
	*base
}

func init(){
	route.Register(&Admin{}, "GET|POST", "vuecmf")
}

func (ctrl *Admin) Index(c *gin.Context){
	req := app.Request{Context: c}
	resp := app.Response{Context: c}

	defer func() {
		if err := recover(); err != nil {
			resp.SendFailure("拉取失败：", err)
		}
	}()

	listParams := helper.DataListParams{}
	req.Input("post", &listParams)

	adminList := service.Admin().List(&listParams)

	resp.SendSuccess("拉取成功", adminList)

}


func (ctrl *Admin) Login(c *gin.Context) {
	loginForm := form.LoginForm{Username: "aaa",Password: "123456"}

	req := app.Request{Context: c}
	resp := app.Response{Context: c}

	//获取输入内容
	req.Input("post", &loginForm)

	//输出内容
	resp.SendSuccess("提交成功", loginForm)

}