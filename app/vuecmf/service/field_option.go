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
	"fmt"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/helper"
	"strconv"
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
type resFieldOption struct {
	FieldId     int
	OptionValue string
	OptionLabel string
}

/*type option struct {
	OptionValue string
	OptionLabel string
}*/

// GetFieldOptions 根据模型ID获取模型的字段选项列表
func (ser *fieldOptionService) GetFieldOptions(modelId int, tableName string, isTree bool, labelFieldName string) (map[string][]*helper.ModelFieldOption, error) {
	var list = make(map[string][]*helper.ModelFieldOption)
	var result []*resFieldOption

	Db.Table(NS.TableName("field_option")).
		Select("model_field_id field_id, option_value, option_label").
		Where("model_id = ?", modelId).
		Where("status = 10").
		Find(&result)

	for _, val := range result {
		/*if list[strconv.Itoa(val.FieldId)] == nil {
			list[strconv.Itoa(val.FieldId)] = &helper.ModelFieldOption{
				Value: val.OptionValue,
				Label: val.OptionLabel,
			}
		}*/
		list[strconv.Itoa(val.FieldId)] = append(list[strconv.Itoa(val.FieldId)], &helper.ModelFieldOption{
			Value: val.OptionValue,
			Label: val.OptionLabel,
		})

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

		//tree := map[string]string{}
		var tree []*helper.ModelFieldOption
		tree = helper.FormatTree(tree, Db, NS.TableName(tableName), "id", 0, labelFieldName, "pid", orderField, 1)
		list[strconv.Itoa(pidFieldId)] = tree

		fmt.Println("tree===", tree)

	}

	return list, nil
}
