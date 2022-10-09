package model

// LoginForm 登录表单
type LoginForm struct {
	LoginName     string `json:"login_name" form:"login_name" binding:"required" required_tips:"登录名不能为空"`
	Password      string `json:"password" form:"password" binding:"required" required_tips:"密码不能为空"`
	LastLoginTime JSONTime
	LastLoginIp   string
}

// DataLoginForm 提交的表单数据
type DataLoginForm struct {
	Data *LoginForm `json:"data" form:"data" binding:"required" required_tips:"参数data不能为空"`
}

// LogoutForm 登出表单
type LogoutForm struct {
	Token string `json:"token" form:"token"`
}

// DataLogoutForm 提交的登出表单数据
type DataLogoutForm struct {
	Data *LogoutForm `json:"data" form:"data" binding:"required" required_tips:"参数data不能为空"`
}
