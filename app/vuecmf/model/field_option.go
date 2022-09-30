package model



// FieldOption 字段选项 模型结构
type FieldOption struct {
	Id uint `json:"id" form:"id"  gorm:"column:id;primaryKey;autoIncrement;size:32;not null;comment:自增ID"`
	ModelId uint `json:"model_id" form:"model_id" binding:"required" required_tips:"请选择" gorm:"column:model_id;size:32;not null;default:0;comment:所属模型ID"`
	ModelFieldId uint `json:"model_field_id" form:"model_field_id" binding:"required" required_tips:"请选择" gorm:"column:model_field_id;size:32;uniqueIndex:unique_index;not null;default:0;comment:模型字段ID"`
	OptionValue string `json:"option_value" form:"option_value" binding:"required" required_tips:"选项值必填" gorm:"column:option_value;size:64;uniqueIndex:unique_index;not null;default:;comment:选项值"`
	OptionLabel string `json:"option_label" form:"option_label" binding:"required" required_tips:"选项标签必填" gorm:"column:option_label;size:255;not null;default:;comment:选项标签"`
	Status uint `json:"status" form:"status"  gorm:"column:status;size:8;not null;default:10;comment:状态：10=开启，20=禁用"`
	Type uint `json:"type" form:"type"  gorm:"column:type;size:8;not null;default:20;comment:类型：10=内置，20=扩展"`
	
}

// DataFieldOptionForm 提交的表单数据
type DataFieldOptionForm struct {
    Data *FieldOption `json:"data" form:"data"`
}