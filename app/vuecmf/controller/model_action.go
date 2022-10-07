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
	modelAction.SaveForm = &model.DataModelActionForm{}
	modelAction.FilterFields = []string{"label", "api_path", "action_type"}

	route.Register(modelAction, "POST", "vuecmf")
}

// GetApiMap 获取API映射的路径
func (ser *ModelAction) GetApiMap(c *gin.Context) {
	dataApiMapForm := &model.DataApiMapForm{}
	common(c, dataApiMapForm, func() (interface{}, error) {
		apiPath := service.ModelAction().GetApiMap(dataApiMapForm.Data.TableName, dataApiMapForm.Data.ActionType, dataApiMapForm.Data.AppId)
		return apiPath, nil
	})
}

// GetActionList 获取所有模型的动作列表
func (ser *ModelAction) GetActionList(c *gin.Context) {
	dataActionListForm := &model.DataActionListForm{}
	common(c, dataActionListForm, func() (interface{}, error) {
		return service.ModelAction().GetActionList(dataActionListForm.Data.RoleName, dataActionListForm.Data.AppName)
	})
}
