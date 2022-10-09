package model

// ModelFormLinkage 模型表单联动 模型结构
type ModelFormLinkage struct {
	Id              uint   `json:"id" form:"id"  gorm:"column:id;primaryKey;autoIncrement;size:32;not null;comment:自增ID"`
	ModelId         uint   `json:"model_id" form:"model_id" binding:"required" required_tips:"请选择" gorm:"column:model_id;size:32;not null;default:0;comment:所属模型ID"`
	ModelFieldId    uint   `json:"model_field_id" form:"model_field_id" binding:"required" required_tips:"请选择" gorm:"column:model_field_id;size:32;uniqueIndex:unique_index;not null;default:0;comment:模型字段ID"`
	LinkageFieldId  uint   `json:"linkage_field_id" form:"linkage_field_id" binding:"required" required_tips:"请选择" gorm:"column:linkage_field_id;size:32;uniqueIndex:unique_index;not null;default:0;comment:联动表单的字段ID"`
	LinkageActionId uint   `json:"linkage_action_id" form:"linkage_action_id" binding:"required" required_tips:"请选择" gorm:"column:linkage_action_id;size:32;not null;default:0;comment:获取联动表单数据的动作ID"`
	Status          uint16 `json:"status" form:"status"  gorm:"column:status;size:8;not null;default:10;comment:状态：10=开启，20=禁用"`
}

// DataModelFormLinkageForm 提交的表单数据
type DataModelFormLinkageForm struct {
	Data *ModelFormLinkage `json:"data" form:"data" binding:"required" required_tips:"参数data不能为空"`
}
