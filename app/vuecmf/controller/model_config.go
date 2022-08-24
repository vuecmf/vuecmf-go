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

type ModelConfig struct {
	Base
}

func init() {
	modelconfig := &ModelConfig{}
	modelconfig.TableName = "modelconfig"
	modelconfig.Model = &model.ModelConfig{}
	route.Register(modelconfig, "POST", "vuecmf")
}

// Index 列表页
func (ctrl *ModelConfig) Index(c *gin.Context) {
	listParams := &helper.DataListParams{}
	common(c, listParams, func() (interface{}, error) {
		return service.ModelConfig().List(listParams)
	})
}

// Save 新增/更新 单条数据
func (ctrl *ModelConfig) Save(c *gin.Context) {
	data := &model.DataModelConfigForm{}
	common(c, data, func() (interface{}, error) {
		if data.Data.Id == 0 {
			return service.ModelConfig().Create(data.Data)
		} else {
			return service.ModelConfig().Update(data.Data)
		}
	})
}

// Saveall 批量添加多条数据
func (ctrl *ModelConfig) Saveall(c *gin.Context) {
	data := &model.DataBatchForm{}
	common(c, data, func() (interface{}, error) {
		var dataBatch []model.ModelConfig
		err := json.Unmarshal([]byte(data.Data), &dataBatch)
		if err != nil {
			return nil, err
		}
		return service.ModelConfig().Create(dataBatch)
	})
}
