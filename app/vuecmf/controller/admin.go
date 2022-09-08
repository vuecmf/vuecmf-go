// Package controller
//+----------------------------------------------------------------------
// | Copyright (c) 2022 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: vuecmf <tulihua2004@126.com>
// +----------------------------------------------------------------------
package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/vuecmf/vuecmf-go/app/route"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/model"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/service"
	"time"
)

type Admin struct {
    Base
}

func init() {
	admin := &Admin{}
	admin.TableName = "admin"
	admin.Model = &model.Admin{}
	admin.listData = &[]model.Admin{}
	admin.saveForm = &model.DataAdminForm{}
	admin.filterFields = []string{"username","email","mobile","token"}
	admin.AppName = "vuecmf"

	route.Register(admin, "POST", admin.AppName)
}

// Login 用户登录
func (ctrl *Admin) Login(c *gin.Context) {
	dataLoginForm := &model.DataLoginForm{}
	common(c, dataLoginForm, func() (interface{}, error) {
		dataLoginForm.Data.LastLoginIp = c.ClientIP()
		dataLoginForm.Data.LastLoginTime = time.Now()
		return service.Admin(ctrl.TableName, ctrl.AppName).Login(dataLoginForm.Data)
	})
}

func (ctrl *Admin) Logout(c *gin.Context) {
	dataLogoutForm := &model.DataLogoutForm{}
	common(c, dataLogoutForm, func() (interface{}, error) {
		return service.Admin(ctrl.TableName, ctrl.AppName).Logout(dataLogoutForm.Data)
	})
}

func (ctrl *Admin) AddRole(c *gin.Context) {

}

func (ctrl *Admin) DelRole(c *gin.Context) {

}

func (ctrl *Admin) AddPermission(c *gin.Context) {

}

func (ctrl *Admin) DelPermission(c *gin.Context) {

}

func (ctrl *Admin) GetPermission(c *gin.Context) {

}

func (ctrl *Admin) GetAllRoles(c *gin.Context) {

}

func (ctrl *Admin) GetRoles(c *gin.Context) {

}

func (ctrl *Admin) SetUserPermission(c *gin.Context) {

}

func (ctrl *Admin) GetUserPermission(c *gin.Context) {

}
