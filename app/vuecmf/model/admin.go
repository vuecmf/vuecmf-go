package model

import "time"

// Admin 管理员 模型结构
type Admin struct {
	Id uint `json:"id" form:"id"  gorm:"column:id;primaryKey;autoIncrement;size:11;not null;comment:自增ID"`
	Username string `json:"username" form:"username" binding:"required" required_tips:"用户名必填" gorm:"column:username;size:32;uniqueIndex:unique_index;not null;default:;comment:用户名"`
	Password string `json:"password" form:"password"  gorm:"column:password;size:255;not null;default:;comment:密码"`
	Email string `json:"email" form:"email" binding:"required,email" required_tips:"邮箱必填" email_tips:"邮箱输入有误" gorm:"column:email;size:64;uniqueIndex:unique_index;not null;default:;comment:邮箱"`
	Mobile string `json:"mobile" form:"mobile" binding:"required" required_tips:"手机必填" gorm:"column:mobile;size:32;uniqueIndex:unique_index;not null;default:;comment:手机"`
	IsSuper uint `json:"is_super" form:"is_super"  gorm:"column:is_super;size:4;not null;default:20;comment:超级管理员：10=是，20=否"`
	RegTime time.Time `json:"reg_time" form:"reg_time" time_format:"2006-01-02 15:04:05"  gorm:"column:reg_time;not null;autoCreateTime;comment:注册时间"`
	RegIp string `json:"reg_ip" form:"reg_ip"  gorm:"column:reg_ip;size:24;not null;default:;comment:注册IP"`
	LastLoginTime time.Time `json:"last_login_time" form:"last_login_time" time_format:"2006-01-02 15:04:05"  gorm:"column:last_login_time;not null;autoCreateTime;autoUpdateTime;comment:最后登录时间"`
	LastLoginIp string `json:"last_login_ip" form:"last_login_ip"  gorm:"column:last_login_ip;size:24;not null;default:;comment:最后登录IP"`
	UpdateTime time.Time `json:"update_time" form:"update_time" time_format:"2006-01-02 15:04:05"  gorm:"column:update_time;not null;autoCreateTime;autoUpdateTime;comment:更新时间"`
	Token string `json:"token" form:"token"  gorm:"column:token;size:255;not null;default:;comment:api访问token"`
	Status uint `json:"status" form:"status"  gorm:"column:status;size:4;not null;default:10;comment:状态：10=开启，20=禁用"`
	
}

// DataAdminForm 提交的表单数据
type DataAdminForm struct {
    Data *Admin `json:"data" form:"data"`
}