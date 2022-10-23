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
	"strings"
)

type ModelIndex struct {
	Base
}

func init() {
	modelIndex := &ModelIndex{}
	modelIndex.TableName = "model_index"
	modelIndex.Model = &model.ModelIndex{}
	modelIndex.ListData = &[]model.ModelIndex{}
	modelIndex.FilterFields = []string{"index_type"}

	route.Register(modelIndex, "POST", "vuecmf")
}

// Save 新增/更新 单条数据
func (ctrl *ModelIndex) Save(c *gin.Context) {
	saveForm := &model.DataModelIndexForm{}
	common(c, saveForm, func() (interface{}, error) {
		saveData := &model.ModelIndex{}
		saveData.Id = saveForm.Data.Id
		saveData.ModelId = saveForm.Data.ModelId
		saveData.ModelFieldId = strings.Join(saveForm.Data.ModelFieldId, ",")
		saveData.IndexType = saveForm.Data.IndexType
		saveData.Status = saveForm.Data.Status

		if saveForm.Data.Id == uint(0) {
			return service.ModelIndex().Create(saveData)
		} else {
			return service.ModelIndex().Update(saveData)
		}
	})
}
