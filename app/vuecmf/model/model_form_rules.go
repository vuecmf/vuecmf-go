package model



// ModelFormRules 模型表单验证 模型结构
type ModelFormRules struct {
	Id uint `json:"id" form:"id"  gorm:"column:id;primaryKey;autoIncrement;size:32;not null;comment:自增ID"`
	ModelId uint `json:"model_id" form:"model_id" binding:"required" required_tips:"请选择" gorm:"column:model_id;size:32;not null;default:0;comment:所属模型ID"`
	ModelFormId uint `json:"model_form_id" form:"model_form_id" binding:"required" required_tips:"请选择" gorm:"column:model_form_id;size:32;not null;default:0;comment:模型表单ID"`
	RuleType string `json:"rule_type" form:"rule_type" binding:"required" required_tips:"请选择" gorm:"column:rule_type;size:32;not null;default:'';comment:表单验证类型"`
	RuleValue string `json:"rule_value" form:"rule_value"  gorm:"column:rule_value;size:255;not null;default:'';comment:表单验证规则"`
	ErrorTips string `json:"error_tips" form:"error_tips"  gorm:"column:error_tips;size:255;not null;default:'';comment:表单验证不通过的错误提示信息"`
	Status uint16 `json:"status" form:"status"  gorm:"column:status;size:8;not null;default:10;comment:状态：10=开启，20=禁用"`
	
}

// DataModelFormRulesForm 提交的表单数据
type DataModelFormRulesForm struct {
    Data *ModelFormRules `json:"data" form:"data"`
}