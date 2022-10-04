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

type ModelConfig struct {
	Base
}

func init() {
	modelConfig := &ModelConfig{}
	modelConfig.TableName = "model_config"
	modelConfig.Model = &model.ModelConfig{}
	modelConfig.ListData = &[]model.ModelConfig{}
	modelConfig.SaveForm = &model.DataModelConfigForm{}
	modelConfig.FilterFields = []string{"table_name", "label", "component_tpl", "remark"}

	route.Register(modelConfig, "POST", "vuecmf")
}
