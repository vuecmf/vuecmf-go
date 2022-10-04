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
	"github.com/vuecmf/vuecmf-go/app/route"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/model"
)

type ModelFormLinkage struct {
	Base
}

func init() {
	modelFormLinkage := &ModelFormLinkage{}
	modelFormLinkage.TableName = "model_form_linkage"
	modelFormLinkage.Model = &model.ModelFormLinkage{}
	modelFormLinkage.ListData = &[]model.ModelFormLinkage{}
	modelFormLinkage.SaveForm = &model.DataModelFormLinkageForm{}
	modelFormLinkage.FilterFields = []string{""}

	route.Register(modelFormLinkage, "POST", "vuecmf")
}
