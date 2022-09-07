package model

import "time"

// LoginForm 登录表单
type LoginForm struct {
	LoginName     string    `json:"login_name" form:"login_name" binding:"required" tips:"登录名不能为空"`
	Password      string    `json:"password" form:"password" binding:"required" tips:"密码不能为空"`
	LastLoginTime time.Time `time_format:"2006-01-02 15:04:05"`
	LastLoginIp   string
}

// DataLoginForm 提交的表单数据
type DataLoginForm struct {
	Data *LoginForm `json:"data" form:"data"`
}

// LogoutForm 登出表单
type LogoutForm struct {
	Token string `json:"token" form:"token"`
}

// DataLogoutForm 提交的登出表单数据
type DataLogoutForm struct {
	Data *LogoutForm `json:"data" form:"data"`
}
