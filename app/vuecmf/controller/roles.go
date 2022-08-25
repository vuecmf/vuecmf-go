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

type Roles struct {
    Base
}

func init() {
    roles := &Roles{}
	roles.TableName = "roles"
	roles.Model = &model.Roles{}
	route.Register(roles, "POST", "vuecmf")
}

// Index 列表页
func (ctrl *Roles) Index(c *gin.Context) {
    listParams := &helper.DataListParams{}
	common(c, listParams, func() (interface{}, error) {
        return service.Roles().List(listParams)
	})
}

// Save 新增/更新 单条数据
func (ctrl *Roles) Save(c *gin.Context) {
	data := &model.DataRolesForm{}
	common(c, data, func() (interface{}, error) {
		if data.Data.Id == 0 {
			return service.Roles().Create(data.Data)
		} else {
			return service.Roles().Update(data.Data)
		}
	})
}

// Saveall 批量添加多条数据
func (ctrl *Roles) Saveall(c *gin.Context) {
	data := &model.DataBatchForm{}
	common(c, data, func() (interface{}, error) {
		var dataBatch []model.Roles
		err := json.Unmarshal([]byte(data.Data), &dataBatch)
		if err != nil {
			return nil, err
		}
		return service.Roles().Create(dataBatch)
	})
}
