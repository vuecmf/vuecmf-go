package model



// ModelFormRules 模型表单验证 模型结构
type ModelFormRules struct {
	Base
	ModelId int `json:"model_id" gorm:"column:model_id;size:11;not null;default:0;comment:所属模型ID"`
	ModelFormId int `json:"model_form_id" gorm:"column:model_form_id;size:11;not null;default:0;comment:模型表单ID"`
	RuleType string `json:"rule_type" gorm:"column:rule_type;size:32;not null;default:;comment:表单验证类型"`
	RuleValue string `json:"rule_value" gorm:"column:rule_value;size:255;not null;default:;comment:表单验证规则"`
	ErrorTips string `json:"error_tips" gorm:"column:error_tips;size:255;not null;default:;comment:表单验证不通过的错误提示信息"`
	
}