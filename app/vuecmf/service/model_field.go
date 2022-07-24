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
	*base
}

var modelField *modelFieldService

// ModelField 获取modelField服务实例
func ModelField() *modelFieldService {
	if modelField == nil {
		modelField = &modelFieldService{}
	}
	return modelField
}

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

func (ser *modelFieldService) GetFieldInfo(modelId int) interface{} {
	var list []fieldInfo

	db.Table(ns.TableName("model_field")).Select(
		"id field_id,"+
			"field_name prop,"+
			"label,"+
			"column_width width,"+
			"length,"+
			"if(is_show = 10,) show,"+
			"is_fixed fixed,"+
			"is_filter filter,"+
			"note tooltip,"+
			"model_id,"+
			"true sortable").
		Where("model_id = ?", modelId).
		Where("status = 10").
		Order("sort_num").
		Find(&list)

	return list
}
