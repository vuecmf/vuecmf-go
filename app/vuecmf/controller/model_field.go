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
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/vuecmf/vuecmf-go/app/route"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/helper"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/model"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/service"
)

type ModelField struct {
	Base
}

func init() {
	modelfield := &ModelField{}
	modelfield.TableName = "modelfield"
	modelfield.Model = &model.ModelField{}
	route.Register(modelfield, "POST", "vuecmf")
}

// Index 列表页
func (ctrl *ModelField) Index(c *gin.Context) {
	listParams := &helper.DataListParams{}
	common(c, listParams, func() (interface{}, error) {
		return service.ModelField().List(listParams)
	})
}

// Save 新增/更新 单条数据
func (ctrl *ModelField) Save(c *gin.Context) {
	data := &model.DataModelFieldForm{}
	common(c, data, func() (interface{}, error) {
		if data.Data.Id == 0 {
			return service.ModelField().Create(data.Data)
		} else {
			return service.ModelField().Update(data.Data)
		}
	})
}

// Saveall 批量添加多条数据
func (ctrl *ModelField) Saveall(c *gin.Context) {
	data := &model.DataBatchForm{}
	common(c, data, func() (interface{}, error) {
		var dataBatch []model.ModelField
		err := json.Unmarshal([]byte(data.Data), &dataBatch)
		if err != nil {
			return nil, err
		}
		return service.ModelField().Create(dataBatch)
	})
}
