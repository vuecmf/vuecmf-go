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
	"github.com/vuecmf/vuecmf-go/app/vuecmf/helper"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/model"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/service"
)

type Menu struct {
    Base
}

func init() {
	menu := &Menu{}
    menu.TableName = "menu"
    menu.Model = &model.Menu{}
    menu.listData = &[]model.Menu{}
    menu.saveForm = &model.DataMenuForm{}
    menu.filterFields = []string{"title","icon"}

    route.Register(menu, "POST", "vuecmf")
}

// Index 列表页
func (ctrl *Menu) Index(c *gin.Context) {
    listParams := &helper.DataListParams{}
	common(c, listParams, func() (interface{}, error) {
        return service.Menu().List(listParams)
	})
}

