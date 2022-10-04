// Package service
//+----------------------------------------------------------------------
// | Copyright (c) 2022 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: vuecmf <tulihua2004@126.com>
// +----------------------------------------------------------------------
package service

// modelIndexService modelIndex服务结构
type modelIndexService struct {
	*BaseService
}

var modelIndex *modelIndexService

// ModelIndex 获取modelIndex服务实例
func ModelIndex() *modelIndexService {
	if modelIndex == nil {
		modelIndex = &modelIndexService{}
	}
	return modelIndex
}
