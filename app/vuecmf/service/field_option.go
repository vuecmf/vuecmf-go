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
	"errors"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/helper"
)

// fieldOptionService fieldOption服务结构
type fieldOptionService struct {
	*BaseService
	TableName string
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
	FieldId     int
	OptionValue string
	OptionLabel string
}

// GetFieldOptions 根据模型ID获取模型的字段选项列表
func (ser *fieldOptionService) GetFieldOptions(modelId int, tableName string, isTree bool, labelFieldName string) (map[int]map[string]string, error) {
	var list = make(map[int]map[string]string)
	var result []modelFieldOption

	Db.Table(NS.TableName("field_option")).
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

	//目录树列表中 父级字段处理
	if isTree {
		if labelFieldName == "" {
			return nil, errors.New("模型还没有设置标题字段")
		}
		orderField := "sort_num"
		if tableName == "roles" {
			orderField = ""
		}

		var pidFieldId int
		Db.Table(NS.TableName("model_field")).
			Select("id").
			Where("field_name = 'pid'").
			Where("model_id = ?", modelId).
			Limit(1).Find(&pidFieldId)

		tree := map[string]string{}
		helper.FormatTree(tree, Db, NS.TableName(tableName), "id", 0, labelFieldName, "pid", orderField, 1)
		list[pidFieldId] = tree
	}

	return list, nil
}
