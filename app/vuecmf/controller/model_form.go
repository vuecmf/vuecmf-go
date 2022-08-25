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

type ModelForm struct {
    Base
}

func init() {
    modelform := &ModelForm{}
	modelform.TableName = "modelform"
	modelform.Model = &model.ModelForm{}
	route.Register(modelform, "POST", "vuecmf")
}

// Index 列表页
func (ctrl *ModelForm) Index(c *gin.Context) {
    listParams := &helper.DataListParams{}
	common(c, listParams, func() (interface{}, error) {
		var result []model.ModelForm
        return service.Base().CommonList(result, ctrl.TableName, listParams)
	})
}

// Save 新增/更新 单条数据
func (ctrl *ModelForm) Save(c *gin.Context) {
	data := &model.DataModelFormForm{}
	common(c, data, func() (interface{}, error) {
		if data.Data.Id == 0 {
			return service.ModelForm().Create(data.Data)
		} else {
			return service.ModelForm().Update(data.Data)
		}
	})
}

// Saveall 批量添加多条数据
func (ctrl *ModelForm) Saveall(c *gin.Context) {
	data := &model.DataBatchForm{}
	common(c, data, func() (interface{}, error) {
		var dataBatch []model.ModelForm
		err := json.Unmarshal([]byte(data.Data), &dataBatch)
		if err != nil {
			return nil, err
		}
		return service.ModelForm().Create(dataBatch)
	})
}

