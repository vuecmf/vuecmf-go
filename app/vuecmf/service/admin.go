//+----------------------------------------------------------------------
// | Copyright (c) 2024 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: tulihua2004@126.com
// +----------------------------------------------------------------------

package service

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/vuecmf/vuecmf-go/v3/app"
	"github.com/vuecmf/vuecmf-go/v3/app/vuecmf/helper"
	"github.com/vuecmf/vuecmf-go/v3/app/vuecmf/model"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// AdminService admin服务结构
type AdminService struct {
	*BaseService
}

var adminOnce sync.Once
var admin *AdminService

// Admin 获取admin服务实例
func Admin() *AdminService {
	adminOnce.Do(func() {
		admin = &AdminService{
			BaseService: &BaseService{
				"admin",
				&model.Admin{},
				&[]model.Admin{},
				[]string{"username", "email", "mobile", "token"},
			},
		}
	})
	return admin
}

type LoginRes struct {
	Token  string            `json:"token"`
	User   map[string]string `json:"user"`
	Server map[string]string `json:"server"`
}

// Create 创建单条或多条数据, 成功返回影响行数
//
//	参数：
//		data 需要保存的数据
func (svc *AdminService) Create(data *model.Admin) (uint, error) {
	res := app.Db.Create(&data)
	return data.Id, res.Error
}

// Update 更新数据, 成功返回影响行数
//
//	参数：
//		data 需要更新的数据
func (svc *AdminService) Update(data *model.Admin) (int64, error) {
	//如果修改用户名，则更新权限中用户名
	var oldUserName string
	DbTable("admin").Select("username").
		Where("id = ?", data.Id).Find(&oldUserName)
	err := Auth().UpdateUser(oldUserName, data.Username)
	if err != nil {
		return 0, err
	}

	res := app.Db.Updates(data)
	return res.RowsAffected, res.Error
}

// IsLogin 验证是否登录
//
//	参数：
//		token 验证token
//		loginIp 登录IP
func (svc *AdminService) IsLogin(token string, loginIp string) (*model.Admin, error) {
	if token == "" {
		return nil, errors.New("您还没有登录，请先登录！")
	}

	var adm *model.Admin
	if err := DbTable("admin").Select("id, username, password, is_super, last_login_time, status").
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
//
//	参数：
//		loginForm 登录传入的表单数据
func (svc *AdminService) Login(loginForm *model.LoginForm) (interface{}, error) {
	loginTimesCacheKey := "vuecmf:login_err_times:" + loginForm.LoginName
	var loginErrTimes int
	_ = app.Cache().Get(loginTimesCacheKey, &loginErrTimes)
	if loginErrTimes > 5 {
		return nil, errors.New("连续登录失败已超过6次，请过两小时后重试！")
	}

	var adminInfo model.Admin

	DbTable("admin").
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

	res := app.Db.Updates(&model.Admin{
		Id:            adminInfo.Id,
		LastLoginTime: loginForm.LastLoginTime,
		LastLoginIp:   loginForm.LastLoginIp,
		Token:         token,
	})

	if res.Error != nil {
		return nil, errors.New("登录出现异常！请稍后重试。" + res.Error.Error())
	}

	var mysqlVersion string
	app.Db.Raw("select version() as v").Scan(&mysqlVersion)

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
	result.Server["version"] = app.Version
	result.Server["os"] = runtime.GOOS
	result.Server["software"] = "Gin"
	result.Server["mysql"] = mysqlVersion
	result.Server["upload_max_size"] = strconv.Itoa(app.Cfg.Upload.AllowFileSize) + "M"

	return result, nil
}

// Logout 用户退出登录
//
//	参数：
//		logoutForm 退出传入的表单数据
func (svc *AdminService) Logout(logoutForm *model.LogoutForm) (bool, error) {
	if logoutForm.Token == "" {
		return false, errors.New("token不能为空")
	}

	DbTable("admin").Where("token = ?", logoutForm.Token).
		Where("status = 10").Update("token", "")

	//清除系统缓存
	_ = app.Cache().Del(CacheUser)

	return true, nil
}

// GetUserNames 根据用户ID获取用户名
//
//	参数：
//		userIdList 用户ID列表
func (svc *AdminService) GetUserNames(userIdList []int) []string {
	var res []string
	DbTable("admin").Select("username").
		Where("id in ?", userIdList).
		Find(&res)
	return res
}

// GetUser 根据用户ID获取用户信息
//
//	参数：
//		userId 用户ID
func (svc *AdminService) GetUser(userId uint) model.Admin {
	var res model.Admin
	DbTable("admin").Select("*").
		Where("id = ?", userId).
		Find(&res)
	return res
}

// GetUserByUsername 根据用户名获取用户信息
//
//	参数：
//		username 用户名
func (svc *AdminService) GetUserByUsername(username string) model.Admin {
	var res model.Admin
	DbTable("admin").Select("*").
		Where("username = ?", username).
		Find(&res)
	return res
}
