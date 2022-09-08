// Package service
//+----------------------------------------------------------------------
// | Copyright (c) 2022 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: vuecmf <tulihua2004@126.com>
// +----------------------------------------------------------------------
package service

// modelFieldService modelField服务结构
type modelFieldService struct {
	*baseService
}

var modelField *modelFieldService

// ModelField 获取modelField服务实例
func ModelField() *modelFieldService {
	if modelField == nil {
		modelField = &modelFieldService{}
	}
	return modelField
}

// fieldInfo 列表表字段信息
type fieldInfo struct {
	FieldId  int    `json:"field_id"`
	Prop     string `json:"prop"`
	Label    string `json:"label"`
	Width    int    `json:"width"`
	Length   int    `json:"length"`
	Show     bool   `json:"show"`
	Fixed    bool   `json:"fixed"`
	Filter   bool   `json:"filter"`
	Tooltip  string `json:"tooltip"`
	ModelId  int    `json:"model_id"`
	Sortable bool   `json:"sortable"`
}

// GetFieldInfo 根据模型ID获取对应的字段信息
func (ser *modelFieldService) GetFieldInfo(modelId int) []fieldInfo {
	var list []fieldInfo

	db.Table(ns.TableName("model_field")).Select(
		"id field_id,"+
			"field_name prop,"+
			"label,"+
			"column_width width,"+
			"length,"+
			"if(is_show = 10,true, false) `show`,"+
			"if(is_fixed = 10,true, false) fixed,"+
			"if(is_filter = 10,true, false) `filter`,"+
			"note tooltip,"+
			"model_id,"+
			"true sortable").
		Where("model_id = ?", modelId).
		Where("status = 10").
		Order("sort_num").
		Find(&list)

	return list
}

// getFilterFields 根据表名获取该表需要模糊查询的字段
func (ser *modelFieldService) getFilterFields(tableName string) []string {
	var filterFields []string
	db.Table(ns.TableName("model_field")+" MF").Select("field_name").
		Joins("left join "+ns.TableName("model_config")+" MC on MF.model_id = MC.id").
		Where("MF.is_filter = 10").
		Where("MF.type in (?)", []string{"char", "varchar"}).
		Where("MC.table_name = ?", tableName).
		Limit(50).Find(&filterFields)

	return filterFields
}
