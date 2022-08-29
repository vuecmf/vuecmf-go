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

type ModelForm struct {
    Base
}

func init() {
	modelForm := &ModelForm{}
    modelForm.TableName = "model_form"
    modelForm.Model = &model.ModelForm{}
    modelForm.listData = &[]model.ModelForm{}
    modelForm.saveForm = &model.DataModelFormForm{}

    route.Register(modelForm, "POST", "vuecmf")
}
