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

type FieldOption struct {
	Base
}

func init() {
	fieldoption := &FieldOption{}
	fieldoption.TableName = "fieldoption"
	fieldoption.Model = &model.FieldOption{}
	route.Register(fieldoption, "POST", "vuecmf")
}

// Index 列表页
func (ctrl *FieldOption) Index(c *gin.Context) {
	listParams := &helper.DataListParams{}
	common(c, listParams, func() (interface{}, error) {
		return service.FieldOption().List(listParams)
	})
}

// Save 新增/更新 单条数据
func (ctrl *FieldOption) Save(c *gin.Context) {
	data := &model.DataFieldOptionForm{}
	common(c, data, func() (interface{}, error) {
		if data.Data.Id == 0 {
			return service.FieldOption().Create(data.Data)
		} else {
			return service.FieldOption().Update(data.Data)
		}
	})
}

// Saveall 批量添加多条数据
func (ctrl *FieldOption) Saveall(c *gin.Context) {
	data := &model.DataBatchForm{}
	common(c, data, func() (interface{}, error) {
		var dataBatch []model.FieldOption
		err := json.Unmarshal([]byte(data.Data), &dataBatch)
		if err != nil {
			return nil, err
		}
		return service.FieldOption().Create(dataBatch)
	})
}
