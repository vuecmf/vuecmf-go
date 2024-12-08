//+----------------------------------------------------------------------
// | Copyright (c) 2024 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: tulihua2004@126.com
// +----------------------------------------------------------------------

package model

// ModelForm 模型表单 模型结构
type ModelForm struct {
	Id           uint   `json:"id" form:"id"  gorm:"column:id;primaryKey;autoIncrement;size:32;not null;comment:自增ID"`
	ModelId      uint   `json:"model_id" form:"model_id" binding:"required" required_tips:"请选择" gorm:"column:model_id;size:32;not null;default:0;comment:所属模型ID"`
	ModelFieldId uint   `json:"model_field_id" form:"model_field_id" binding:"required" required_tips:"请选择" gorm:"column:model_field_id;size:32;uniqueIndex:unique_index;not null;default:0;comment:模型字段ID"`
	Type         string `json:"type" form:"type" binding:"required" required_tips:"请选择" gorm:"column:type;size:32;not null;default:'';comment:表单控件类型"`
	DefaultValue string `json:"default_value" form:"default_value"  gorm:"column:default_value;size:255;not null;default:'';comment:表单控件默认值"`
	Placeholder  string `json:"placeholder" form:"placeholder"  gorm:"column:placeholder;size:255;not null;default:'';comment:表单提示信息"`
	IsDisabled   uint16 `json:"is_disabled" form:"is_disabled"  gorm:"column:is_disabled;size:8;not null;default:20;comment:添加/编辑表单中是否禁用： 10=是，20=否"`
	IsEdit       uint16 `json:"is_edit" form:"is_edit"  gorm:"column:is_edit;size:8;not null;default:10;comment:可编辑： 10=是，20=否"`
	SortNum      uint   `json:"sort_num" form:"sort_num"  gorm:"column:sort_num;size:32;not null;default:0;comment:菜单的排列顺序(小在前)"`
	Status       uint16 `json:"status" form:"status"  gorm:"column:status;size:8;not null;default:10;comment:状态：10=开启，20=禁用"`
}

// DataModelFormForm 提交的表单数据
type DataModelFormForm struct {
	Data *ModelForm `json:"data" form:"data" binding:"required" required_tips:"参数data不能为空"`
}
