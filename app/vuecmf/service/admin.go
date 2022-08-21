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

// adminService admin服务结构
type adminService struct {
	*baseService
	TableName string
}

// List 获取列表数据
// 		参数：params 查询参数
func (ser *adminService) List(params *helper.DataListParams) (interface{}, error) {
	var adminList []model.Admin
	return ser.commonList(adminList, ser.TableName, params)
}

var admin *adminService

// Admin 获取admin服务实例
func Admin() *adminService {
	if admin == nil {
		admin = &adminService{TableName: "admin"}
	}
	return admin
}
