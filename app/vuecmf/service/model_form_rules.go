// Package service
//+----------------------------------------------------------------------
// | Copyright (c) 2022 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: vuecmf <tulihua2004@126.com>
// +----------------------------------------------------------------------
package service

import "strings"

// modelFormRulesService modelFormRules服务结构
type modelFormRulesService struct {
	*BaseService
}

var modelFormRules *modelFormRulesService

// ModelFormRules 获取modelFormRules服务实例
func ModelFormRules() *modelFormRulesService {
	if modelFormRules == nil {
		modelFormRules = &modelFormRulesService{}
	}
	return modelFormRules
}

type ruleListFormST struct {
	FieldName, RuleType, RuleValue, ErrorTips, FieldType string
}

// GetRuleListForForm 根据模型ID获取对应的表单验证规则
func (ser *modelFormRulesService) GetRuleListForForm(modelId int) interface{} {
	var data []ruleListFormST

	Db.Table(NS.TableName("model_form_rules")+" VMFR").
		Select("VMF2.field_name, VMF2.type field_type, rule_type, rule_value, error_tips").
		Joins("LEFT JOIN "+NS.TableName("model_form")+" VMF ON VMFR.model_form_id = VMF.id").
		Joins("INNER JOIN "+NS.TableName("model_field")+" VMF2 ON VMF.model_field_id = VMF2.id").
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
		case "bigint","int","smallint","tinyint":
			fieldType = "integer"
		case "decimal","double","float":
			fieldType = "number"
		case "date","datetime","timestamp":
			fieldType = "date"
		default:
			fieldType = "string"
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
