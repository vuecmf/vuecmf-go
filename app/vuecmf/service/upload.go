// Package service
//+----------------------------------------------------------------------
// | Copyright (c) 2022 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: vuecmf <tulihua2004@126.com>
// +----------------------------------------------------------------------
package service

// uploadService upload服务结构
type uploadService struct {
	*baseService
}

var upload *uploadService

// Upload 获取upload服务实例
func Upload() *uploadService {
	if upload == nil {
		upload = &uploadService{}
	}
	return upload
}