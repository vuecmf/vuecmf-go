//+----------------------------------------------------------------------
// | Copyright (c) 2023 http://www.vuecmf.com All rights reserved.
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

// ModelConfig 模型配置管理
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
		if saveForm.Data.Id == uint(0) {
			return service.ModelConfig().Create(saveForm.Data)
		} else {
			return service.ModelConfig().Update(saveForm.Data)
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
