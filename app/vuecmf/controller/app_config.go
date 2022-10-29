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
	"errors"
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
	Common(c, saveForm, func() (interface{}, error) {
		if saveForm.Data.Id == uint(0) {
			//创建应用目录
			if err := service.Make().CreateApp(saveForm.Data.AppName); err != nil {
				return nil, err
			}
			return service.Base().Create(saveForm.Data)
		} else {
			//更新应用目录
			if err := service.Make().RenameApp(saveForm.Data.Id, saveForm.Data.AppName); err != nil {
				return nil, err
			}
			return service.Base().Update(saveForm.Data)
		}
	})
}

// Delete 根据ID删除单条数据
func (ctrl *AppConfig) Delete(c *gin.Context) {
	data := &model.DataIdForm{}
	Common(c, data, func() (interface{}, error) {
		//先检查应用下是否存在模型，若存在则不允许删除
		if num := service.AppConfig().GetAppModelCount(data.Data.Id); num > 0 {
			return nil, errors.New("不允许删除有分配模型的应用！")
		}
		//移除应用相关目录
		if err := service.Make().RemoveApp(data.Data.Id); err != nil {
			return nil, err
		}
		return service.Base().Delete(data.Data.Id, ctrl.Model)
	})
}

// GetAllModels 获取所有模型
func (ctrl *AppConfig) GetAllModels(c *gin.Context) {
	Common(c, nil, func() (interface{}, error) {
		res := service.AppConfig().GetAllModels()
		return res, nil
	})
}

// GetModels 获取应用的所有模型
func (ctrl *AppConfig) GetModels(c *gin.Context) {
	dataAppIdForm := &model.DataAppIdForm{}
	Common(c, dataAppIdForm, func() (interface{}, error) {
		return service.AppConfig().GetModels(dataAppIdForm.Data.AppId)
	})
}

// AddModel 给应用添加模型
func (ctrl *AppConfig) AddModel(c *gin.Context) {
	dataAddModelForm := &model.DataAddModelForm{}
	Common(c, dataAddModelForm, func() (interface{}, error) {
		if len(dataAddModelForm.Data.ModelIdList) == 0 {
			//如果模型列表为空，即表示应用没有模型，则清空应用的所有模型
			return service.AppConfig().DelAllModelsForApp(dataAddModelForm.Data.AppId)
		} else {
			return service.AppConfig().AddModelsForApp(dataAddModelForm.Data.AppId, dataAddModelForm.Data.ModelIdList)
		}
	})
}
