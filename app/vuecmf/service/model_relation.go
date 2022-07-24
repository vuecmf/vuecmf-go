// Package service
//+----------------------------------------------------------------------
// | Copyright (c) 2022 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: vuecmf <tulihua2004@126.com>
// +----------------------------------------------------------------------
package service

// modelRelationService modelRelation服务结构
type modelRelationService struct {
	*base
}

var modelRelation *modelRelationService

// ModelRelation 获取modelRelation服务实例
func ModelRelation() *modelRelationService {
	if modelRelation == nil {
		modelRelation = &modelRelationService{}
	}
	return modelRelation
}