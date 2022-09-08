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

type ModelIndex struct {
    Base
}

func init() {
	modelIndex := &ModelIndex{}
    modelIndex.TableName = "model_index"
    modelIndex.Model = &model.ModelIndex{}
    modelIndex.listData = &[]model.ModelIndex{}
    modelIndex.saveForm = &model.DataModelIndexForm{}
    modelIndex.filterFields = []string{"index_type"}

    route.Register(modelIndex, "POST", "vuecmf")
}
