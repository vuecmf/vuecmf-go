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

type ModelField struct {
	Base
}

func init() {
	modelField := &ModelField{}
	modelField.TableName = "model_field"
	modelField.Model = &model.ModelField{}
	modelField.ListData = &[]model.ModelField{}
	modelField.FilterFields = []string{"field_name", "label", "type", "note", "default_value"}

	route.Register(modelField, "POST", "vuecmf")
}

// Save 新增/更新 单条数据
func (ctrl *ModelField) Save(c *gin.Context) {
	saveForm := &model.DataModelFieldForm{}
	common(c, saveForm, func() (interface{}, error) {
		if saveForm.Data.Id == uint(0) {
			return service.Base().Create(saveForm.Data)
		} else {
			return service.Base().Update(saveForm.Data)
		}
	})
}
