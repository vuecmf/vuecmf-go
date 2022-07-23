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
	"github.com/vuecmf/vuecmf-go/app"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/helper"
)

type base struct {
}

// commonIndex 公共列表页方法
func commonIndex(c *gin.Context, fun func(listParams *helper.DataListParams) interface{}) {
	defer func() {
		if err := recover(); err != nil {
			app.Response(c).SendFailure("拉取失败：", err)
		}
	}()

	listParams := helper.DataListParams{}
	app.Request(c).Input("post", &listParams)

	list := fun(&listParams)

	app.Response(c).SendSuccess("拉取成功", list)
}
