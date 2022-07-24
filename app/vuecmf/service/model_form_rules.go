// Package service
//+----------------------------------------------------------------------
// | Copyright (c) 2022 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: vuecmf <tulihua2004@126.com>
// +----------------------------------------------------------------------
package service

// modelFormRulesService modelFormRules服务结构
type modelFormRulesService struct {
	*base
}

var modelFormRules *modelFormRulesService

// ModelFormRules 获取modelFormRules服务实例
func ModelFormRules() *modelFormRulesService {
	if modelFormRules == nil {
		modelFormRules = &modelFormRulesService{}
	}
	return modelFormRules
}