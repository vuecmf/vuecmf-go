// Package service
//+----------------------------------------------------------------------
// | Copyright (c) 2022 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: vuecmf <tulihua2004@126.com>
// +----------------------------------------------------------------------
package service

// modelFormService modelForm服务结构
type modelFormService struct {
	*base
}

var modelForm *modelFormService

// ModelForm 获取modelForm服务实例
func ModelForm() *modelFormService {
	if modelForm == nil {
		modelForm = &modelFormService{}
	}
	return modelForm
}