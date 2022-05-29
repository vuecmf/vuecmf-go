package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/vuecmf/vuecmf-go/app"
	"github.com/vuecmf/vuecmf-go/app/service"
)

type AdminController struct {
	*BaseController
}


// LoginForm 登录表单
type LoginForm struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" `
}


func (ctrl *AdminController) Index(c *gin.Context){
	req := app.Request{Context: c}

	listParams := DataListParams{}
	req.Input("post", &listParams)

	adminService := service.AdminService{}

	adminList := adminService.List(
		listParams.Data.Filter,
		listParams.Data.Keywords,
		listParams.Data.Page,
		listParams.Data.PageSize,
		listParams.Data.OrderField,
		listParams.Data.OrderSort,
		)

	resp := app.Response{Context: c}
	resp.SendSuccess("拉取成功", adminList)
}


func (ctrl *BaseController) Login(c *gin.Context) {
	loginForm := LoginForm{Username: "aaa",Password: "123456"}

	req := app.Request{Context: c}
	resp := app.Response{Context: c}

	//获取输入内容
	req.Input("post", &loginForm)

	//输出内容
	resp.SendSuccess("提交成功", loginForm)

}