package model



// FieldOption 字段选项 模型结构
type FieldOption struct {
	Base
	ModelId int `json:"model_id" gorm:"column:model_id;size:11;not null;default:0;comment:所属模型ID"`
	ModelFieldId int `json:"model_field_id" gorm:"column:model_field_id;size:11;uniqueIndex:unique_index;not null;default:0;comment:模型字段ID"`
	OptionValue string `json:"option_value" gorm:"column:option_value;size:64;uniqueIndex:unique_index;not null;default:;comment:选项值"`
	OptionLabel string `json:"option_label" gorm:"column:option_label;size:255;not null;default:;comment:选项标签"`
	Type uint8 `json:"type" gorm:"column:type;size:4;not null;default:20;comment:类型：10=内置，20=扩展"`
	
}