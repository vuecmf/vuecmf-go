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

type ModelIndex struct {
    Base
}

func init() {
    modelindex := &ModelIndex{}
	modelindex.TableName = "modelindex"
	modelindex.Model = &model.ModelIndex{}
	route.Register(modelindex, "POST", "vuecmf")
}

// Index 列表页
func (ctrl *ModelIndex) Index(c *gin.Context) {
    listParams := &helper.DataListParams{}
	common(c, listParams, func() (interface{}, error) {
		var result []model.ModelIndex
        return service.Base().CommonList(result, ctrl.TableName, listParams)
	})
}

// Save 新增/更新 单条数据
func (ctrl *ModelIndex) Save(c *gin.Context) {
	data := &model.DataModelIndexForm{}
	common(c, data, func() (interface{}, error) {
		if data.Data.Id == 0 {
			return service.ModelIndex().Create(data.Data)
		} else {
			return service.ModelIndex().Update(data.Data)
		}
	})
}

// Saveall 批量添加多条数据
func (ctrl *ModelIndex) Saveall(c *gin.Context) {
	data := &model.DataBatchForm{}
	common(c, data, func() (interface{}, error) {
		var dataBatch []model.ModelIndex
		err := json.Unmarshal([]byte(data.Data), &dataBatch)
		if err != nil {
			return nil, err
		}
		return service.ModelIndex().Create(dataBatch)
	})
}

