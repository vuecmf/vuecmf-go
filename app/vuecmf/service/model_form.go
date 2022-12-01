// Package service
//+----------------------------------------------------------------------
// | Copyright (c) 2022 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: vuecmf <tulihua2004@126.com>
// +----------------------------------------------------------------------
package service

import "github.com/vuecmf/vuecmf-go/app/vuecmf/model"

// modelFormService modelForm服务结构
type modelFormService struct {
	*BaseService
}

var modelForm *modelFormService

// ModelForm 获取modelForm服务实例
func ModelForm() *modelFormService {
	if modelForm == nil {
		modelForm = &modelFormService{}
	}
	return modelForm
}

// formInfo 表单字段信息
type formInfo struct {
	FieldId      int    `json:"field_id"`
	FieldName    string `json:"field_name"`
	Label        string `json:"label"`
	Type         string `json:"type"`
	DefaultValue string `json:"default_value"`
	IsDisabled   bool   `json:"is_disabled"`
	SortNum      int    `json:"sort_num"`
}

// GetFormInfo 根据模型ID获取模型的表单信息
func (ser *modelFormService) GetFormInfo(modelId int) []formInfo {
	var list []formInfo

	Db.Table(NS.TableName("model_form")+" VMF").
		Select("VMF.model_field_id field_id, VMF.`type`, VMF.default_value, if(VMF.is_disabled = 10, 1, 0) is_disabled, VMF.sort_num, VMF2.field_name, VMF2.label").
		Joins("inner join "+NS.TableName("model_field")+" VMF2 ON VMF.model_field_id = VMF2.id").
		Where("VMF.model_id = ?", modelId).
		Where("VMF.status = 10").
		Where("VMF2.status = 10").
		Order("VMF.sort_num").
		Find(&list)

	return list

}

//DelByFieldId 根据字段ID删除
func (ser *modelFormService) DelByFieldId(fieldId uint) error {
	res := Db.Table(NS.TableName("model_form")).Delete(&model.ModelForm{ModelFieldId: fieldId})
	return res.Error
}
