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

// ModelFormRules 模型表单验证规则
type ModelFormRules struct {
	Base
}

func init() {
	modelFormRules := &ModelFormRules{}
	modelFormRules.TableName = "model_form_rules"
	modelFormRules.Model = &model.ModelFormRules{}
	modelFormRules.ListData = &[]model.ModelFormRules{}
	modelFormRules.FilterFields = []string{"rule_type", "rule_value", "error_tips"}

	route.Register(modelFormRules, "POST", "vuecmf")
}

// Save 新增/更新 单条数据
func (ctrl *ModelFormRules) Save(c *gin.Context) {
	saveForm := &model.DataModelFormRulesForm{}
	Common(c, saveForm, func() (interface{}, error) {
		if saveForm.Data.Id == uint(0) {
			return service.Base().Create(saveForm.Data)
		} else {
			return service.Base().Update(saveForm.Data)
		}
	})
}

// Delete 根据ID删除单条数据
func (ctrl *ModelFormRules) Delete(c *gin.Context) {
	data := &model.DataIdForm{}
	Common(c, data, func() (interface{}, error) {
		return service.ModelFormRules().Delete(data.Data.Id, ctrl.Model)
	})
}
