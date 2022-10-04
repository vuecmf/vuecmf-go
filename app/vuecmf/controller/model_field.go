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
	modelField.ListData = &[]model.ModelField{}
	modelField.SaveForm = &model.DataModelFieldForm{}
	modelField.FilterFields = []string{"field_name", "label", "type", "note", "default_value"}

	route.Register(modelField, "POST", "vuecmf")
}
