// Package service
//+----------------------------------------------------------------------
// | Copyright (c) 2022 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: vuecmf <tulihua2004@126.com>
// +----------------------------------------------------------------------
package service

// modelActionService modelAction服务结构
type modelActionService struct {
	*baseService
}

var modelAction *modelActionService

// ModelAction 获取modelAction服务实例
func ModelAction() *modelActionService {
	if modelAction == nil {
		modelAction = &modelActionService{}
	}
	return modelAction
}