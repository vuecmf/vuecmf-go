package form

import "time"

// Admin 管理员 表单结构
type Admin struct {
	Id int `json:"id" form:"id" `
	RegTime time.Time `json:"reg_time" form:"reg_time" time_format:"2006-01-02 15:04:05" `
	UpdateTime time.Time `json:"update_time" form:"update_time" time_format:"2006-01-02 15:04:05" `
	Token string `json:"token" form:"token" `
	Status int `json:"status" form:"status" `
	Email string `json:"email" form:"email" binding:"required,email" required_tips:"邮箱必填" email_tips:"邮箱输入有误"`
	Mobile string `json:"mobile" form:"mobile" binding:"required" required_tips:"手机必填"`
	IsSuper int `json:"is_super" form:"is_super" `
	RegIp int `json:"reg_ip" form:"reg_ip" `
	LastLoginTime time.Time `json:"last_login_time" form:"last_login_time" time_format:"2006-01-02 15:04:05" `
	LastLoginIp int `json:"last_login_ip" form:"last_login_ip" `
	Username string `json:"username" form:"username" binding:"required,min=4,max=32" required_tips:"用户名必填" min_tips:"用户名长度为4到32个字符" max_tips:"用户名长度为4到32个字符"`
	Password string `json:"password" form:"password" `
	
}

type DataAdminForm struct {
    Data *Admin `json:"data" form:"data"`
}