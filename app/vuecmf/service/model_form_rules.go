//+----------------------------------------------------------------------
// | Copyright (c) 2024 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: tulihua2004@126.com
// +----------------------------------------------------------------------

package service

import (
	"github.com/vuecmf/vuecmf-go/v3/app"
	"github.com/vuecmf/vuecmf-go/v3/app/vuecmf/model"
	"gorm.io/gorm"
	"strings"
	"sync"
)

// ModelFormRulesService modelFormRules服务结构
type ModelFormRulesService struct {
	*BaseService
}

var modelFormRulesOnce sync.Once
var modelFormRules *ModelFormRulesService

// ModelFormRules 获取modelFormRules服务实例
func ModelFormRules() *ModelFormRulesService {
	modelFormRulesOnce.Do(func() {
		modelFormRules = &ModelFormRulesService{
			BaseService: &BaseService{
				"model_form_rules",
				&model.ModelFormRules{},
				&[]model.ModelFormRules{},
				[]string{"rule_type", "rule_value", "error_tips"},
			},
		}
	})
	return modelFormRules
}

type modelInfo struct {
	TableName string
	AppName   string
}

// Delete 根据ID删除数据
//
//	参数：
//		id 需删除的ID
//		model 模型实例
func (svc *ModelFormRulesService) Delete(id uint) (int64, error) {
	var res modelInfo
	DbTable("model_form_rules", "R").Select("C.table_name, A.app_name").
		Joins("left join "+TableName("model_config")+" C on R.model_id = C.id").
		Joins("left join "+TableName("app_config")+" A on C.app_id = A.id").
		Where("R.id = ?", id).
		Where("C.status = 10").
		Where("A.status = 10").
		Find(&res)
	err := app.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&svc.Model, id).Error; err != nil {
			return err
		}
		return Make().Model(res.TableName, res.AppName)
	})
	if err != nil {
		return 0, err
	}
	return 1, nil
}

type ruleListFormST struct {
	FieldName, RuleType, RuleValue, ErrorTips, FieldType, FormType string
}

// GetRuleListForForm 根据模型ID获取对应的表单验证规则
//
//	参数：
//		modelId 模型ID
func (svc *ModelFormRulesService) GetRuleListForForm(modelId int) interface{} {
	var data []ruleListFormST

	DbTable("model_form_rules", "VMFR").
		Select("VMF2.field_name, VMF2.type field_type, rule_type, rule_value, error_tips, VMF.type form_type").
		Joins("LEFT JOIN "+TableName("model_form")+" VMF ON VMFR.model_form_id = VMF.id").
		Joins("INNER JOIN "+TableName("model_field")+" VMF2 ON VMF.model_field_id = VMF2.id").
		Where("rule_type IN ?", []string{"require", "length", "date", "email", "integer", "number", "regex", "float", "array", "url"}).
		Where("VMFR.model_id = ?", modelId).
		Where("VMFR.status = 10").
		Where("VMF.status = 10").
		Where("VMF2.status = 10").
		Find(&data)

	result := make(map[string][]map[string]interface{})

	for _, val := range data {
		var fieldType string

		switch val.FieldType {
		case "bigint", "int", "smallint", "tinyint":
			fieldType = "integer"
		case "decimal", "double", "float":
			fieldType = "float"
		case "date", "datetime", "timestamp":
			fieldType = "date"
		default:
			fieldType = "string"
			if val.FormType == "select_mul" || val.FormType == "upload_image" || val.FormType == "upload_file" {
				fieldType = "array"
			}
		}

		rule := make(map[string]interface{})

		//验证数据类型
		rule["type"] = fieldType

		switch val.RuleType {
		case "require":
			rule["required"] = true
			rule["message"] = val.ErrorTips
			rule["trigger"] = "blur"
		case "length":
			arr := strings.Split(val.RuleValue, ",")
			rule["min"] = arr[0]
			rule["max"] = arr[1]
			rule["message"] = val.ErrorTips
			rule["trigger"] = "blur"
		case "date", "email", "integer", "number", "regex", "float", "array", "url":
			if val.RuleType == "regex" {
				val.RuleType = "regexp"
			}
			rule["type"] = val.RuleType
			rule["required"] = "true"
			rule["message"] = val.ErrorTips
			rule["trigger"] = "blur, change"
		}

		result[val.FieldName] = append(result[val.FieldName], rule)

	}

	return result
}
