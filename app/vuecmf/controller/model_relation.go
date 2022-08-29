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

type ModelRelation struct {
    Base
}

func init() {
	modelRelation := &ModelRelation{}
    modelRelation.TableName = "model_relation"
    modelRelation.Model = &model.ModelRelation{}
    modelRelation.listData = &[]model.ModelRelation{}
    modelRelation.saveForm = &model.DataModelRelationForm{}

    route.Register(modelRelation, "POST", "vuecmf")
}
