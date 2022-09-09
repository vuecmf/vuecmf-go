// Package service
//+----------------------------------------------------------------------
// | Copyright (c) 2022 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: vuecmf <tulihua2004@126.com>
// +----------------------------------------------------------------------
package service

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/vuecmf/vuecmf-go/app"
	"github.com/vuecmf/vuecmf-go/app/route"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/helper"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/model"
	"runtime"
	"strconv"
	"strings"
)

// adminService admin服务结构
type adminService struct {
	*baseService
	TableName string
	AppName   string
}

var admin *adminService

// Admin 获取admin服务实例
func Admin(tableName string, appName string) *adminService {
	if admin == nil {
		admin = &adminService{
			TableName: ns.TableName(tableName),
			AppName:   appName,
		}
	}
	return admin
}

type LoginRes struct {
	Token  string            `json:"token"`
	User   map[string]string `json:"user"`
	Server map[string]string `json:"server"`
}

// Login 用户登录
func (ser *adminService) Login(loginForm *model.LoginForm) (interface{}, error) {
	loginTimesCacheKey := "vuecmf:login_err_times:" + loginForm.LoginName
	var loginErrTimes int
	_ = app.Cache().Get(loginTimesCacheKey, &loginErrTimes)
	if loginErrTimes > 5 {
		return nil, errors.New("连续登录失败已超过6次，请过两小时后重试！")
	}

	var adminInfo model.Admin

	db.Table(ser.TableName).
		Where("username = ? or email = ? or mobile = ?", loginForm.LoginName, loginForm.LoginName, loginForm.LoginName).
		Where("status = 10").
		Find(&adminInfo)

	//设置登录失败次数，超过则不允许登录，两小时后重试
	if !helper.PasswordVerify(loginForm.Password, adminInfo.Password) {
		_ = app.Cache().Set(loginTimesCacheKey, loginErrTimes+1)
		return nil, errors.New("错误的登录名称或密码！请检查是否输入有误。")
	}

	codeByte := md5.Sum([]byte(adminInfo.Username + adminInfo.Password + loginForm.LastLoginIp))
	token := fmt.Sprintf("%x", codeByte)

	res := db.Table(ser.TableName).Where("id = ?", adminInfo.Id).
		Updates(model.Admin{
			LastLoginTime: loginForm.LastLoginTime,
			LastLoginIp:   loginForm.LastLoginIp,
			Token:         token,
		})

	if res.Error != nil {
		return nil, errors.New("登录出现异常！请稍后重试。" + res.Error.Error())
	}

	var mysqlVersion string
	db.Raw("select version() as v").Scan(&mysqlVersion)

	role := "超级管理员"
	if adminInfo.IsSuper != 10 {
		auth := Auth().Enforcer
		if auth != nil {
			return nil, errors.New("获取角色失败。")
		}
		roleArr, err2 := auth.GetRolesForUser(adminInfo.Username, strings.ToLower(ser.AppName))
		if err2 != nil {
			return nil, errors.New("获取角色失败。" + err2.Error())
		}
		role = strings.Join(roleArr, ",")
	}

	var result LoginRes
	result.Token = token
	result.User = map[string]string{}
	result.User["username"] = adminInfo.Username
	result.User["role"] = role
	result.User["last_login_time"] = loginForm.LastLoginTime.String()
	result.User["last_login_ip"] = loginForm.LastLoginIp
	result.Server = map[string]string{}
	result.Server["version"] = "2.0.0"
	result.Server["os"] = runtime.GOOS
	result.Server["software"] = "Gin"
	result.Server["mysql"] = mysqlVersion
	result.Server["upload_max_size"] = strconv.FormatInt(route.Engine.MaxMultipartMemory/1024/1024, 10) + "M"

	return result, nil
}

// Logout 用户退出登录
func (ser *adminService) Logout(logoutForm *model.LogoutForm) (bool, error) {
	if logoutForm.Token == "" {
		return false, errors.New("token不能为空")
	}

	db.Table(ser.TableName).Where("token = ?", logoutForm.Token).
		Where("status = 10").Update("token", "")

	//清除系统缓存
	_ = app.Cache().Del(CacheUser)

	return true, nil
}



