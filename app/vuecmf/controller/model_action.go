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
	"github.com/vuecmf/vuecmf-go/app/route"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/model"
)

type ModelAction struct {
    Base
}

func init() {
	modelAction := &ModelAction{}
    modelAction.TableName = "model_action"
    modelAction.Model = &model.ModelAction{}
    modelAction.listData = &[]model.ModelAction{}
    modelAction.saveForm = &model.DataModelActionForm{}
    modelAction.filterFields = []string{"label","api_path","action_type"}

    route.Register(modelAction, "POST", "vuecmf")
}
