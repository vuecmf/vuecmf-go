// Package service
//+----------------------------------------------------------------------
// | Copyright (c) 2022 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: vuecmf <tulihua2004@126.com>
// +----------------------------------------------------------------------
package service

// adminService admin服务结构
type adminService struct {
	*baseService
}

var admin *adminService

// Admin 获取admin服务实例
func Admin() *adminService {
	if admin == nil {
		admin = &adminService{}
	}
	return admin
}
