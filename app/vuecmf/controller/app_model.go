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

type AppModel struct {
	Base
}

func init() {
	appModel := &AppModel{}
	appModel.TableName = "app_model"
	appModel.Model = &model.AppModel{}
	appModel.ListData = &[]model.AppModel{}
	appModel.FilterFields = []string{""}

	route.Register(appModel, "POST", "vuecmf")
}

// Save 新增/更新 单条数据
func (ctrl *AppModel) Save(c *gin.Context) {
	saveForm := &model.DataAppModelForm{}
	common(c, saveForm, func() (interface{}, error) {
		if saveForm.Data.Id == uint(0) {
			return service.Base().Create(saveForm.Data)
		} else {
			return service.Base().Update(saveForm.Data)
		}
	})
}

// GetModelList 获取应用下所有模型
func (ctrl *AppModel) GetModelList(c *gin.Context) {
	modelListForm := &model.DataModelListForm{}
	common(c, modelListForm, func() (interface{}, error) {
		return service.AppModel().GetModelList(modelListForm.Data.AppId)
	})
}



