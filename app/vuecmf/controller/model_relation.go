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

type ModelRelation struct {
	Base
}

func init() {
	modelRelation := &ModelRelation{}
	modelRelation.TableName = "model_relation"
	modelRelation.Model = &model.ModelRelation{}
	modelRelation.ListData = &[]model.ModelRelation{}
	modelRelation.FilterFields = []string{"relation_show_field_id"}

	route.Register(modelRelation, "POST", "vuecmf")
}

// Save 新增/更新 单条数据
func (ctrl *ModelRelation) Save(c *gin.Context) {
	saveForm := &model.DataModelRelationForm{}
	Common(c, saveForm, func() (interface{}, error) {
		if saveForm.Data.Id == uint(0) {
			return service.Base().Create(saveForm.Data)
		} else {
			return service.Base().Update(saveForm.Data)
		}
	})
}
