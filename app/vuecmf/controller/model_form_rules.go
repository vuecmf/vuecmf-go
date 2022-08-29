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

type ModelFormRules struct {
    Base
}

func init() {
	modelFormRules := &ModelFormRules{}
    modelFormRules.TableName = "model_form_rules"
    modelFormRules.Model = &model.ModelFormRules{}
    modelFormRules.listData = &[]model.ModelFormRules{}
    modelFormRules.saveForm = &model.DataModelFormRulesForm{}

    route.Register(modelFormRules, "POST", "vuecmf")
}
