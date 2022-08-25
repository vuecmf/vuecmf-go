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

type ModelFormRules struct {
    Base
}

func init() {
    modelformrules := &ModelFormRules{}
	modelformrules.TableName = "model_form_rules"
	modelformrules.Model = &model.ModelFormRules{}
	route.Register(modelformrules, "POST", "vuecmf")
}

// Index 列表页
func (ctrl *ModelFormRules) Index(c *gin.Context) {
    listParams := &helper.DataListParams{}
	common(c, listParams, func() (interface{}, error) {
		var result []model.ModelFormRules
        return service.Base().CommonList(result, ctrl.TableName, listParams)
	})
}

// Save 新增/更新 单条数据
func (ctrl *ModelFormRules) Save(c *gin.Context) {
	data := &model.DataModelFormRulesForm{}
	common(c, data, func() (interface{}, error) {
		if data.Data.Id == 0 {
			return service.ModelFormRules().Create(data.Data)
		} else {
			return service.ModelFormRules().Update(data.Data)
		}
	})
}

// Saveall 批量添加多条数据
func (ctrl *ModelFormRules) Saveall(c *gin.Context) {
	data := &model.DataBatchForm{}
	common(c, data, func() (interface{}, error) {
		var dataBatch []model.ModelFormRules
		err := json.Unmarshal([]byte(data.Data), &dataBatch)
		if err != nil {
			return nil, err
		}
		return service.ModelFormRules().Create(dataBatch)
	})
}

