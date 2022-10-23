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
)

type AppConfig struct {
	Base
}

func init() {
	appConfig := &AppConfig{}
	appConfig.TableName = "app_config"
	appConfig.Model = &model.AppConfig{}
	appConfig.ListData = &[]model.AppConfig{}
	appConfig.FilterFields = []string{"app_name", "exclusion_url"}

	route.Register(appConfig, "POST", "vuecmf")
}

// Save 新增/更新 单条数据
func (ctrl *AppConfig) Save(c *gin.Context) {
	saveForm := &model.DataAppConfigForm{}
	common(c, saveForm, func() (interface{}, error) {
		if saveForm.Data.Id == uint(0) {
			return service.Base().Create(saveForm.Data)
		} else {
			return service.Base().Update(saveForm.Data)
		}
	})
}
