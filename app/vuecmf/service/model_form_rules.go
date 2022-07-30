// Package service
//+----------------------------------------------------------------------
// | Copyright (c) 2022 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: vuecmf <tulihua2004@126.com>
// +----------------------------------------------------------------------
package service

// modelFormRulesService modelFormRules服务结构
type modelFormRulesService struct {
	*base
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

func (ser *modelFormRulesService) GetRuleListForForm (modelId int) interface{} {
	var data []ruleListFormST

	db.Table(ns.TableName("model_form_rules") + " VMFR").
	  	Select('VMF2.field_name, rule_type, rule_value, error_tips').
		Joins("LEFT JOIN model_form VMF ON VMFR.model_form_id = VMF.id").
		Joins("INNER JOIN model_field VMF2 ON VMF.model_field_id = VMF2.id").
		Where("rule_type IN ?", []string{"require","length","date","email","integer","number","regex","float","array","url"}).
		Where('VMFR.model_id = ?', modelId).
		Where("VMFR.status = 10").
		Where("VMF.status = 10").
		Where("VMF2.status = 10").
		Find(&data)

	result := make(map[string]map[int]map[string]string)

	for key, val := range data {
		switch val.RuleType {
		case "require":
			result[val.FieldName][key]["required"] = "true"
			result[val.FieldName][key]["message"] =  val.ErrorTips
			result[val.FieldName][key]["trigger"] = "blur"
		case "length":
		$arr = explode(',', $val->rule_value);
		$result[$val->field_name][] = [
		'min' => intval($arr[0]) ?? 0,
		'max' => intval($arr[1]) ?? 0,
		'message' => $val->error_tips,
		'trigger' => 'blur'
		];
		break;
		case 'date':
		case 'email':
		case 'integer':
		case 'number':
		case 'regex':
		case 'float':
		case 'array':
		case 'url':
		$val->rule_type == 'regex' && $val->rule_type = 'regexp';
		$result[$val->field_name][] = [
		'type' => $val->rule_type,
		'required' => true,
		'message' => $val->error_tips,
		'trigger' => ['blur', 'change']
		];
		break;
		}
	}

	return result
}
