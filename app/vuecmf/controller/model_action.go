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

type ModelAction struct {
	Base
}

func init() {
	modelAction := &ModelAction{}
	modelAction.TableName = "model_action"
	modelAction.Model = &model.ModelAction{}
	modelAction.ListData = &[]model.ModelAction{}
	modelAction.FilterFields = []string{"label", "api_path", "action_type"}

	route.Register(modelAction, "POST", "vuecmf")
}

// Save 新增/更新 单条数据
func (ctrl *ModelAction) Save(c *gin.Context) {
	saveForm := &model.DataModelActionForm{}
	Common(c, saveForm, func() (interface{}, error) {
		if saveForm.Data.Id == uint(0) {
			return service.Base().Create(saveForm.Data)
		} else {
			return service.Base().Update(saveForm.Data)
		}
	})
}

// GetApiMap 获取API映射的路径
func (ctrl *ModelAction) GetApiMap(c *gin.Context) {
	dataApiMapForm := &model.DataApiMapForm{}
	Common(c, dataApiMapForm, func() (interface{}, error) {
		apiPath := service.ModelAction().GetApiMap(dataApiMapForm.Data.TableName, dataApiMapForm.Data.ActionType, dataApiMapForm.Data.AppId)
		return apiPath, nil
	})
}

// GetActionList 获取所有模型的动作列表
func (ctrl *ModelAction) GetActionList(c *gin.Context) {
	dataActionListForm := &model.DataActionListForm{}
	Common(c, dataActionListForm, func() (interface{}, error) {
		return service.ModelAction().GetActionList(dataActionListForm.Data.RoleName)
	})
}
