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
	FieldName, RuleType, RuleValue, ErrorTips string
}

// GetRuleListForForm 根据模型ID获取对应的表单验证规则
func (ser *modelFormRulesService) GetRuleListForForm(modelId int) interface{} {
	var data []ruleListFormST

	Db.Table(NS.TableName("model_form_rules")+" VMFR").
		Select("VMF2.field_name, rule_type, rule_value, error_tips").
		Joins("LEFT JOIN model_form VMF ON VMFR.model_form_id = VMF.id").
		Joins("INNER JOIN model_field VMF2 ON VMF.model_field_id = VMF2.id").
		Where("rule_type IN ?", []string{"require", "length", "date", "email", "integer", "number", "regex", "float", "array", "url"}).
		Where("VMFR.model_id = ?", modelId).
		Where("VMFR.status = 10").
		Where("VMF.status = 10").
		Where("VMF2.status = 10").
		Find(&data)

	result := make(map[string]map[int]map[string]string)

	for key, val := range data {
		switch val.RuleType {
		case "require":
			result[val.FieldName][key]["required"] = "true"
			result[val.FieldName][key]["message"] = val.ErrorTips
			result[val.FieldName][key]["trigger"] = "blur"
		case "length":
			arr := strings.Split(val.RuleValue, ",")
			result[val.FieldName][key]["min"] = arr[0]
			result[val.FieldName][key]["max"] = arr[1]
			result[val.FieldName][key]["message"] = val.ErrorTips
			result[val.FieldName][key]["trigger"] = "blur"
		case "date", "email", "integer", "number", "regex", "float", "array", "url":
			if val.RuleType == "regex" {
				val.RuleType = "regexp"
			}
			result[val.FieldName][key]["type"] = val.RuleType
			result[val.FieldName][key]["required"] = "true"
			result[val.FieldName][key]["message"] = val.ErrorTips
			result[val.FieldName][key]["trigger"] = "blur, change"
		}
	}

	return result
}
