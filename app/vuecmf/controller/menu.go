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
	"encoding/json"
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
	route.Register(menu, "POST", "vuecmf")
}

// Index 列表页
func (ctrl *Menu) Index(c *gin.Context) {
    listParams := &helper.DataListParams{}
	common(c, listParams, func() (interface{}, error) {
        return service.Menu().List(listParams)
	})
}

// Save 新增/更新 单条数据
func (ctrl *Menu) Save(c *gin.Context) {
	data := &model.DataMenuForm{}
	common(c, data, func() (interface{}, error) {
		if data.Data.Id == 0 {
			return service.Menu().Create(data.Data)
		} else {
			return service.Menu().Update(data.Data)
		}
	})
}

// Saveall 批量添加多条数据
func (ctrl *Menu) Saveall(c *gin.Context) {
	data := &model.DataBatchForm{}
	common(c, data, func() (interface{}, error) {
		var dataBatch []model.Menu
		err := json.Unmarshal([]byte(data.Data), &dataBatch)
		if err != nil {
			return nil, err
		}
		return service.Menu().Create(dataBatch)
	})
}

