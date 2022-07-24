// Package service
//+----------------------------------------------------------------------
// | Copyright (c) 2022 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: vuecmf <tulihua2004@126.com>
// +----------------------------------------------------------------------
package service

// fieldOptionService fieldOption服务结构
type fieldOptionService struct {
	*base
}

var fieldOption *fieldOptionService

// FieldOption 获取fieldOption服务实例
func FieldOption() *fieldOptionService {
	if fieldOption == nil {
		fieldOption = &fieldOptionService{}
	}
	return fieldOption
}