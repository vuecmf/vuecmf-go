package form

// AdminForm 管理员 表单结构
type AdminForm struct {
	Username string `json:"username" form:"username" binding:"required,len=32" required_tips:"用户名必填" len_tips:"用户名长度为4到32个字符"`
	Email string `json:"email" form:"email" binding:"required,email" required_tips:"邮箱必填" email_tips:"邮箱输入有误"`
	Mobile string `json:"mobile" form:"mobile" binding:"required" required_tips:"手机必填"`
	
}

type DataAdminForm struct {
    Data AdminForm `json:"data" form:"data"`
}

