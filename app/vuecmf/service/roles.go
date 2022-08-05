// Package service
//+----------------------------------------------------------------------
// | Copyright (c) 2022 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: vuecmf <tulihua2004@126.com>
// +----------------------------------------------------------------------
package service

import (
	"github.com/vuecmf/vuecmf-go/app/vuecmf/helper"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/model"
)

// rolesService roles服务结构
type rolesService struct {
	*base
	TableName string
}

var roles *rolesService

// Roles 获取roles服务实例
func Roles() *rolesService {
	if roles == nil {
		roles = &rolesService{TableName: "roles"}
	}
	return roles
}

// List 获取列表数据
// 		参数：params 查询参数
func (ser *rolesService) List(params *helper.DataListParams) interface{} {
	if params.Data.Action == "getField" {
		//拉取列表的字段信息
		return ser.getFieldList(ser.TableName, params.Data.Filter)
	}else{
		//拉取列表的数据
		var rolesList []*model.Roles
		var res = make(map[string]interface{})

		ser.getList(&rolesList, ser.TableName, params)

		//转换成树形列表
		tree := model.RolesModel().ToTree(rolesList)
		res["data"] = tree
		return res
	}
}
