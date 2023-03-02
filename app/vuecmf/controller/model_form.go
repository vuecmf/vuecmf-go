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

// ModelForm 模型表单管理
type ModelForm struct {
	Base
}

func init() {
	modelForm := &ModelForm{}
	modelForm.TableName = "model_form"
	modelForm.Model = &model.ModelForm{}
	modelForm.ListData = &[]model.ModelForm{}
	modelForm.FilterFields = []string{"type", "default_value"}

	route.Register(modelForm, "POST", "vuecmf")
}

// Save 新增/更新 单条数据
func (ctrl *ModelForm) Save(c *gin.Context) {
	saveForm := &model.DataModelFormForm{}
	Common(c, saveForm, func() (interface{}, error) {
		if saveForm.Data.Id == uint(0) {
			return service.Base().Create(saveForm.Data)
		} else {
			return service.Base().Update(saveForm.Data)
		}
	})
}
