package model

// LoginForm 登录表单
type LoginForm struct {
	Username string `json:"username" form:"username" binding:"required" tips:"用户名不能为空"`
	Password string `json:"password" form:"password" binding:"required" tips:"密码不能为空"`
}
