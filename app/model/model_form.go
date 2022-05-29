package model

// ModelForm 模型表单 模型结构
type ModelForm struct {
	Base
	ModelId      int    `json:"model_id" gorm:"column:model_id;size:11;not null;default:0;comment:所属模型ID"`
	ModelFieldId int    `json:"model_field_id" gorm:"column:model_field_id;size:11;uniqueIndex:unique_index;not null;default:0;comment:模型字段ID"`
	Type         string `json:"type" gorm:"column:type;size:32;not null;default:;comment:表单控件类型"`
	DefaultValue string `json:"default_value" gorm:"column:default_value;size:255;not null;default:;comment:表单控件默认值"`
	IsDisabled   uint8  `json:"is_disabled" gorm:"column:is_disabled;size:4;not null;default:20;comment:添加/编辑表单中是否禁用： 10=是，20=否"`
	SortNum      int    `json:"sort_num" gorm:"column:sort_num;size:11;not null;default:0;comment:菜单的排列顺序(小在前)"`
}
