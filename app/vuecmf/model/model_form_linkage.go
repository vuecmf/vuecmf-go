package model



// ModelFormLinkage 模型表单联动 模型结构
type ModelFormLinkage struct {
	LinkageActionId uint `json:"linkage_action_id" form:"linkage_action_id" binding:"required" required_tips:"请选择" gorm:"column:linkage_action_id;size:11;not null;default:0;comment:获取联动表单数据的动作ID"`
	Id uint `json:"id" form:"id"  gorm:"column:id;primaryKey;autoIncrement;size:11;not null;default:0;comment:自增ID"`
	Status uint `json:"status" form:"status"  gorm:"column:status;size:4;not null;default:10;comment:状态：10=开启，20=禁用"`
	ModelId uint `json:"model_id" form:"model_id" binding:"required" required_tips:"请选择" gorm:"column:model_id;size:11;not null;default:0;comment:所属模型ID"`
	ModelFieldId uint `json:"model_field_id" form:"model_field_id" binding:"required" required_tips:"请选择" gorm:"column:model_field_id;size:11;uniqueIndex:unique_index;not null;default:0;comment:模型字段ID"`
	LinkageFieldId uint `json:"linkage_field_id" form:"linkage_field_id" binding:"required" required_tips:"请选择" gorm:"column:linkage_field_id;size:11;uniqueIndex:unique_index;not null;default:0;comment:联动表单的字段ID"`
	
}

// DataModelFormLinkageForm 提交的表单数据
type DataModelFormLinkageForm struct {
    Data *ModelFormLinkage `json:"data" form:"data"`
}