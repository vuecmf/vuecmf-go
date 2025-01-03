//+----------------------------------------------------------------------
// | Copyright (c) 2024 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: tulihua2004@126.com
// +----------------------------------------------------------------------

package service

import (
	"github.com/vuecmf/vuecmf-go/v3/app/vuecmf/model"
	"sync"
)

// ModelFormService modelForm服务结构
type ModelFormService struct {
	*BaseService
}

var modelFormOnce sync.Once
var modelForm *ModelFormService

// ModelForm 获取modelForm服务实例
func ModelForm() *ModelFormService {
	modelFormOnce.Do(func() {
		modelForm = &ModelFormService{
			BaseService: &BaseService{
				"model_form",
				&model.ModelForm{},
				&[]model.ModelForm{},
				[]string{"type", "default_value"},
			},
		}
	})
	return modelForm
}

// formInfo 表单字段信息
type formInfo struct {
	FieldId      int    `json:"field_id"`
	FieldName    string `json:"field_name"`
	Label        string `json:"label"`
	Type         string `json:"type"`
	DefaultValue string `json:"default_value"`
	Placeholder  string `json:"placeholder"`
	IsDisabled   bool   `json:"is_disabled"`
	IsEdit       bool   `json:"is_edit"`
	SortNum      int    `json:"sort_num"`
}

// GetFormInfo 根据模型ID获取模型的表单信息
//
//	参数：
//		modelId 模型ID
//		isSuper 是否为超级管理员
func (svc *ModelFormService) GetFormInfo(modelId int, isSuper uint16) []formInfo {
	var list []formInfo

	query := DbTable("model_form", "VMF").
		Select("VMF.model_field_id field_id, VMF.`type`, VMF.default_value,VMF.placeholder, if(VMF.is_disabled = 10, 1, 0) is_disabled, if(VMF.is_edit = 10, 1, 0) is_edit, VMF.sort_num, VMF2.field_name, VMF2.label").
		Joins("inner join "+TableName("model_field")+" VMF2 ON VMF.model_field_id = VMF2.id").
		Where("VMF.model_id = ?", modelId).
		Where("VMF.status = 10").
		Where("VMF2.status = 10")

	if isSuper != 10 {
		query = query.Where("VMF2.field_name != 'is_super'")
	}

	query.Order("VMF.sort_num").Find(&list)

	return list

}

// DelByFieldId 根据字段ID删除
//
//	参数：
//		fieldId 字段ID
func (svc *ModelFormService) DelByFieldId(fieldId uint) error {
	res := DbTable("model_form").Delete(&model.ModelForm{ModelFieldId: fieldId})
	return res.Error
}
