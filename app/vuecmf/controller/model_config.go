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
    modelConfig.listData = &[]model.ModelConfig{}
    modelConfig.saveForm = &model.DataModelConfigForm{}

    route.Register(modelConfig, "POST", "vuecmf")
}
