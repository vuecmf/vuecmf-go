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
	"github.com/vuecmf/vuecmf-go/app/vuecmf/service"
)

type Menu struct {
	*base
}

func init() {
	route.Register(&Menu{}, "GET|POST", "vuecmf")
}

// Index 列表页
func (ctrl *Menu) Index(c *gin.Context) {
	commonIndex(c, func(listParams *helper.DataListParams) interface{} {
		return service.Menu().List(listParams)
	})
}

