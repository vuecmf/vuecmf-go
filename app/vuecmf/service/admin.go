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
	"github.com/vuecmf/vuecmf-go/app/vuecmf"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/helper"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/model"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// adminService admin服务结构
type adminService struct {
	*BaseService
}

var admin *adminService

// Admin 获取admin服务实例
func Admin() *adminService {
	if admin == nil {
		admin = &adminService{}
	}
	return admin
}

type LoginRes struct {
	Token  string            `json:"token"`
	User   map[string]string `json:"user"`
	Server map[string]string `json:"server"`
}

//IsLogin 验证是否登录
func (ser *adminService) IsLogin(token string, loginIp string) (*model.Admin, error) {
	if token == "" {
		return nil, errors.New("您还没有登录，请先登录！")
	}

	var adm *model.Admin
	if err := Db.Table(NS.TableName("admin")).Select("username, password, is_super, last_login_time, status").
		Where("token = ?", token).Find(&adm).Error; err != nil {
		return nil, errors.New("验证是否登录IsLogin异常：" + err.Error())
	}
	if adm == nil {
		return nil, errors.New("您还没有登录或登录已失效，请重新登录！")
	}
	if adm.Status == 20 {
		return nil, errors.New("登录账号已禁用！")
	}

	codeByte := md5.Sum([]byte(adm.Username + adm.Password + loginIp + time.Now().Format(model.DateFormat)))
	newToken := fmt.Sprintf("%x", codeByte)

	if token != newToken {
		return nil, errors.New("登录已失效，请重新登录！")
	}

	return adm, nil
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

	Db.Table(NS.TableName("admin")).
		Where("username = ? or email = ? or mobile = ?", loginForm.LoginName, loginForm.LoginName, loginForm.LoginName).
		Where("status = 10").
		Find(&adminInfo)

	//设置登录失败次数，超过则不允许登录，两小时后重试
	if !helper.PasswordVerify(loginForm.Password, adminInfo.Password) {
		_ = app.Cache().Set(loginTimesCacheKey, loginErrTimes+1)
		return nil, errors.New("错误的登录名称或密码！请检查是否输入有误。")
	}

	codeByte := md5.Sum([]byte(adminInfo.Username + adminInfo.Password + loginForm.LastLoginIp + loginForm.LastLoginTime.Format(model.DateFormat)))
	token := fmt.Sprintf("%x", codeByte)

	res := Db.Updates(&model.Admin{
		Id:            adminInfo.Id,
		LastLoginTime: loginForm.LastLoginTime,
		LastLoginIp:   loginForm.LastLoginIp,
		Token:         token,
	})

	if res.Error != nil {
		return nil, errors.New("登录出现异常！请稍后重试。" + res.Error.Error())
	}

	var mysqlVersion string
	Db.Raw("select version() as v").Scan(&mysqlVersion)

	role := "超级管理员"
	if adminInfo.IsSuper != 10 {
		roleArr, err := Auth().GetRolesForUser(adminInfo.Username)
		if err != nil {
			return nil, errors.New("获取角色失败。" + err.Error())
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
	result.Server["version"] = vuecmf.Version
	result.Server["os"] = runtime.GOOS
	result.Server["software"] = "Gin"
	result.Server["mysql"] = mysqlVersion
	result.Server["upload_max_size"] = strconv.Itoa(Conf.Upload.AllowFileSize) + "M"

	return result, nil
}

// Logout 用户退出登录
func (ser *adminService) Logout(logoutForm *model.LogoutForm) (bool, error) {
	if logoutForm.Token == "" {
		return false, errors.New("token不能为空")
	}

	Db.Table(NS.TableName("admin")).Where("token = ?", logoutForm.Token).
		Where("status = 10").Update("token", "")

	//清除系统缓存
	_ = app.Cache().Del(CacheUser)

	return true, nil
}

//GetUserNames 根据用户ID获取用户名
func (ser *adminService) GetUserNames(userIdList []int) []string {
	var res []string
	Db.Table(NS.TableName("admin")).Select("username").
		Where("id in ?", userIdList).
		Find(&res)
	return res
}

//GetUser 根据用户ID获取用户信息
func (ser *adminService) GetUser(userId uint) model.Admin {
	var res model.Admin
	Db.Table(NS.TableName("admin")).Select("*").
		Where("id = ?", userId).
		Find(&res)
	return res
}
