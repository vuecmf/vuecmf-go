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
	"strconv"
	"strings"
)

type ModelConfig struct {
	Base
}

func init() {
	modelConfig := &ModelConfig{}
	modelConfig.TableName = "model_config"
	modelConfig.Model = &model.ModelConfig{}
	modelConfig.ListData = &[]model.ModelConfig{}
	modelConfig.FilterFields = []string{"table_name", "label", "component_tpl", "remark"}

	route.Register(modelConfig, "POST", "vuecmf")
}

// Save 新增/更新 单条数据
func (ctrl *ModelConfig) Save(c *gin.Context) {
	saveForm := &model.DataModelConfigForm{}
	Common(c, saveForm, func() (interface{}, error) {
		var searchFieldIdList []string
		if len(saveForm.Data.SearchFieldId) > 0 {
			for _, id := range saveForm.Data.SearchFieldId {
				searchFieldIdList = append(searchFieldIdList, strconv.Itoa(id))
			}
		}

		saveData := &model.ModelConfig{}
		saveData.Id = saveForm.Data.Id
		saveData.AppId = saveForm.Data.AppId
		saveData.TableName = saveForm.Data.TableName
		saveData.Label = saveForm.Data.Label
		saveData.ComponentTpl = saveForm.Data.ComponentTpl
		saveData.DefaultActionId = saveForm.Data.DefaultActionId
		saveData.SearchFieldId = strings.Join(searchFieldIdList, ",")
		saveData.Type = saveForm.Data.Type
		saveData.IsTree = saveForm.Data.IsTree
		saveData.Remark = saveForm.Data.Remark
		saveData.Status = saveForm.Data.Status

		if saveForm.Data.Id == uint(0) {
			return service.ModelConfig().Create(saveData)
		} else {
			return service.ModelConfig().Update(saveData)
		}
	})
}

// Delete 根据ID删除单条数据
func (ctrl *ModelConfig) Delete(c *gin.Context) {
	data := &model.DataIdForm{}
	Common(c, data, func() (interface{}, error) {
		return service.ModelConfig().Delete(data.Data.Id, &model.ModelConfig{})
	})
}

// DeleteBatch 根据ID列表批量删除多条数据
func (ctrl *ModelConfig) DeleteBatch(c *gin.Context) {
	data := &model.DataIdListForm{}
	Common(c, data, func() (interface{}, error) {
		return service.ModelConfig().DeleteBatch(data.Data.IdList, &model.ModelConfig{})
	})
}
