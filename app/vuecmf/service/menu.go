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

// menuService menu服务结构
type menuService struct {
	*baseService
	TableName string
}

var menu *menuService

// Menu 获取menu服务实例
func Menu() *menuService {
	if menu == nil {
		menu = &menuService{TableName: "menu"}
	}
	return menu
}

// List 获取列表数据
// 		参数：params 查询参数
func (ser *menuService) List(params *helper.DataListParams) (interface{}, error) {
	if params.Data.Action == "getField" {
		//拉取列表的字段信息
		return ser.getFieldList(ser.TableName, params.Data.Filter)
	} else {
		//拉取列表的数据
		var menuList []*model.Menu
		var res = make(map[string]interface{})

		ser.getList(&menuList, ser.TableName, params)

		//转换成树形列表
		tree := model.MenuModel().ToTree(menuList)
		res["data"] = tree
		return res, nil
	}
}
