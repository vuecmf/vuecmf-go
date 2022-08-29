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

type FieldOption struct {
    Base
}

func init() {
	fieldOption := &FieldOption{}
    fieldOption.TableName = "field_option"
    fieldOption.Model = &model.FieldOption{}
    fieldOption.listData = &[]model.FieldOption{}
    fieldOption.saveForm = &model.DataFieldOptionForm{}

    route.Register(fieldOption, "POST", "vuecmf")
}
