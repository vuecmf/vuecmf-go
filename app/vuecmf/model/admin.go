//+----------------------------------------------------------------------
// | Copyright (c) 2024 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: tulihua2004@126.com
// +----------------------------------------------------------------------

package model

import (
	"github.com/vuecmf/vuecmf-go/v3/app/vuecmf/helper"
	"gorm.io/gorm"
)

// Admin 管理员 模型结构
type Admin struct {
	Id            uint     `json:"id" form:"id"  gorm:"column:id;primaryKey;autoIncrement;size:32;not null;comment:自增ID"`
	Username      string   `json:"username" form:"username" binding:"required" required_tips:"用户名必填" gorm:"column:username;size:32;unique;not null;default:'';comment:用户名"`
	Password      string   `json:"password" form:"password"  gorm:"column:password;size:255;not null;default:'';comment:密码"`
	Email         string   `json:"email" form:"email" binding:"required,email" required_tips:"邮箱必填" email_tips:"邮箱输入有误" gorm:"column:email;size:64;unique;not null;default:'';comment:邮箱"`
	Mobile        string   `json:"mobile" form:"mobile" binding:"required" required_tips:"手机必填" gorm:"column:mobile;size:32;unique;not null;default:'';comment:手机"`
	IsSuper       uint16   `json:"is_super" form:"is_super"  gorm:"column:is_super;size:8;not null;default:20;comment:超级管理员：10=是，20=否"`
	RegTime       JSONTime `json:"reg_time" form:"reg_time" gorm:"type:timestamp;column:reg_time;not null;autoCreateTime;default:CURRENT_TIMESTAMP;comment:注册时间"`
	RegIp         string   `json:"reg_ip" form:"reg_ip"  gorm:"column:reg_ip;size:24;not null;default:'';comment:注册IP"`
	LastLoginTime JSONTime `json:"last_login_time" form:"last_login_time" gorm:"type:timestamp;column:last_login_time;not null;autoCreateTime;autoUpdateTime;default:CURRENT_TIMESTAMP;comment:最后登录时间"`
	LastLoginIp   string   `json:"last_login_ip" form:"last_login_ip"  gorm:"column:last_login_ip;size:24;not null;default:'';comment:最后登录IP"`
	UpdateTime    JSONTime `json:"update_time" form:"update_time" gorm:"type:timestamp;column:update_time;not null;autoCreateTime;autoUpdateTime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间"`
	Token         string   `json:"token" form:"token"  gorm:"column:token;size:255;not null;default:'';comment:api访问token"`
	Pid           uint     `json:"pid" form:"pid"  gorm:"column:pid;size:32;not null;default:0;comment:父级用户ID"`
	Status        uint16   `json:"status" form:"status"  gorm:"column:status;size:8;not null;default:10;comment:状态：10=开启，20=禁用"`
}

// DataAdminForm 提交的表单数据
type DataAdminForm struct {
	Data *Admin `json:"data" form:"data" binding:"required" required_tips:"参数data不能为空"`
}

// BeforeSave 数据更新前处理
func (m *Admin) BeforeSave(tx *gorm.DB) error {
	var err error
	//如果有填写密码，则更新加密密码
	if m.Password != "" && len(m.Password) >= 4 && len(m.Password) <= 32 {
		m.Password, err = helper.PasswordHash(m.Password)
	}
	return err
}
