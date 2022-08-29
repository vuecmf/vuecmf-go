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

type ModelField struct {
    Base
}

func init() {
	modelField := &ModelField{}
    modelField.TableName = "model_field"
    modelField.Model = &model.ModelField{}
    modelField.listData = &[]model.ModelField{}
    modelField.saveForm = &model.DataModelFieldForm{}

    route.Register(modelField, "POST", "vuecmf")
}
