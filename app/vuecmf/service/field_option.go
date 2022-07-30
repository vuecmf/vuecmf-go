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

// 模型的字段选项
type modelFieldOption struct {
	FieldId int
	OptionValue string
	OptionLabel string
}

// GetFieldOptions 根据模型ID获取模型的字段选项列表
func (ser *fieldOptionService) GetFieldOptions(modelId int) map[int]map[string]string {
	var list = make(map[int]map[string]string)
	var result []modelFieldOption

	db.Table(ns.TableName("field_option")).
		Select("model_field_id field_id, option_value, option_label").
		Where("model_id = ?", modelId).
		Where("status = 10").
		Find(&result)

	for _, val := range result {
		if list[val.FieldId] == nil {
			list[val.FieldId] = map[string]string{}
		}
		list[val.FieldId][val.OptionValue] = val.OptionLabel
	}
	return list
}