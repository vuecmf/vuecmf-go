package model



// ModelFormLinkage 模型表单联动 模型结构
type ModelFormLinkage struct {
	Base
	ModelId int `json:"model_id" gorm:"column:model_id;size:11;not null;default:0;comment:所属模型ID"`
	ModelFieldId int `json:"model_field_id" gorm:"column:model_field_id;size:11;uniqueIndex:unique_index;not null;default:0;comment:模型字段ID"`
	LinkageFieldId int `json:"linkage_field_id" gorm:"column:linkage_field_id;size:11;uniqueIndex:unique_index;not null;default:0;comment:联动表单的字段ID"`
	LinkageActionId int `json:"linkage_action_id" gorm:"column:linkage_action_id;size:11;not null;default:0;comment:获取联动表单数据的动作ID"`
	
}