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

// ModelFormLinkage 模型表单联动管理
type ModelFormLinkage struct {
	Base
}

func init() {
	modelFormLinkage := &ModelFormLinkage{}
	modelFormLinkage.TableName = "model_form_linkage"
	modelFormLinkage.Model = &model.ModelFormLinkage{}
	modelFormLinkage.ListData = &[]model.ModelFormLinkage{}
	modelFormLinkage.FilterFields = []string{""}

	route.Register(modelFormLinkage, "POST", "vuecmf")
}

// Save 新增/更新 单条数据
func (ctrl *ModelFormLinkage) Save(c *gin.Context) {
	saveForm := &model.DataModelFormLinkageForm{}
	Common(c, saveForm, func() (interface{}, error) {
		if saveForm.Data.Id == uint(0) {
			return service.Base().Create(saveForm.Data)
		} else {
			return service.Base().Update(saveForm.Data)
		}
	})
}
