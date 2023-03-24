//+----------------------------------------------------------------------
// | Copyright (c) 2023 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: vuecmf <tulihua2004@126.com>
// +----------------------------------------------------------------------

package service

import (
	"errors"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/helper"
	"gorm.io/gorm"
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

// GetFieldOptions 根据模型ID获取模型的字段选项列表
//	参数：
// 		modelId 模型ID
//		tableName 表名
//		isTree 是否为目录树
//		labelFieldName 需要显示为标签的字段
//		filter 筛选条件
//		db  菜单下拉的db实例
func (ser *fieldOptionService) GetFieldOptions(modelId int, tableName string, isTree bool, labelFieldName string, filter map[string]interface{}, db *gorm.DB) (map[string][]*helper.ModelFieldOption, error) {
	var list = make(map[string][]*helper.ModelFieldOption)
	var result []*resFieldOption

	Db.Table(NS.TableName("field_option")).
		Select("model_field_id field_id, option_value, if((option_value REGEXP '[0-9]') = 1 , option_label, concat(option_value,' (',option_label, ')')) option_label").
		Where("model_id = ?", modelId).
		Where("status = 10").
		Find(&result)

	for _, val := range result {
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

		var tree []*helper.ModelFieldOption
		tree = helper.FormatTree(tree, db, db.NamingStrategy.TableName(tableName), filter, "id", 0, labelFieldName, "pid", orderField, 1)
		list[strconv.Itoa(pidFieldId)] = tree

	}

	return list, nil
}
