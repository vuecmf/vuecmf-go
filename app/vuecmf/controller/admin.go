package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/vuecmf/vuecmf-go/app"
	"github.com/vuecmf/vuecmf-go/app/route"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/service"
)

type Admin struct {
	*base
}

func init(){
	route.Register(&Admin{}, "GET|POST", "vuecmf")
}

// LoginForm 登录表单
type LoginForm struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" `
}


func (ctrl *Admin) Index(c *gin.Context){
	req := app.Request{Context: c}
	resp := app.Response{Context: c}

	defer func() {
		if err := recover(); err != nil {
			resp.SendFailure("拉取失败：", err)
		}
	}()

	listParams := DataListParams{}
	req.Input("post", &listParams)

	adminService := service.AdminService{}

	if listParams.Data == nil {
		panic("请求参数data不能为空")
	}

	adminList := adminService.List(
		listParams.Data.Filter,
		listParams.Data.Keywords,
		listParams.Data.Page,
		listParams.Data.PageSize,
		listParams.Data.OrderField,
		listParams.Data.OrderSort,
		)

	resp.SendSuccess("拉取成功", adminList)
	//app.Json(c.Writer, adminList)
}


func (ctrl *Admin) Login(c *gin.Context) {
	loginForm := LoginForm{Username: "aaa",Password: "123456"}

	req := app.Request{Context: c}
	resp := app.Response{Context: c}

	//获取输入内容
	req.Input("post", &loginForm)

	//输出内容
	resp.SendSuccess("提交成功", loginForm)

}