package model

import "time"

// Admin 管理员 模型结构
type Admin struct {
	Base
	Username string `json:"username" gorm:"column:username;size:32;uniqueIndex:unique_index;not null;default:;comment:用户名"`
	Password string `json:"password" gorm:"column:password;size:255;not null;default:;comment:密码"`
	Email string `json:"email" gorm:"column:email;size:64;uniqueIndex:unique_index;not null;default:;comment:邮箱"`
	Mobile string `json:"mobile" gorm:"column:mobile;size:32;uniqueIndex:unique_index;not null;default:;comment:手机"`
	IsSuper uint8 `json:"is_super" gorm:"column:is_super;size:4;not null;default:20;comment:超级管理员：10=是，20=否"`
	RegTime time.Time `json:"reg_time" gorm:"column:reg_time;not null;autoCreateTime;comment:注册时间"`
	RegIp string `json:"reg_ip" gorm:"column:reg_ip;size:24;not null;default:;comment:注册IP"`
	LastLoginTime time.Time `json:"last_login_time" gorm:"column:last_login_time;not null;autoCreateTime;autoUpdateTime;comment:最后登录时间"`
	LastLoginIp string `json:"last_login_ip" gorm:"column:last_login_ip;size:24;not null;default:;comment:最后登录IP"`
	UpdateTime time.Time `json:"update_time" gorm:"column:update_time;not null;autoCreateTime;autoUpdateTime;comment:更新时间"`
	Token string `json:"token" gorm:"column:token;size:255;not null;default:;comment:api访问token"`
	
}