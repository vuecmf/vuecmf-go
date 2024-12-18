//+----------------------------------------------------------------------
// | Copyright (c) 2024 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: tulihua2004@126.com
// +----------------------------------------------------------------------

package model

// ModelField 模型字段 模型结构
type ModelField struct {
	Id              uint   `json:"id" form:"id"  gorm:"column:id;primaryKey;autoIncrement;size:32;not null;comment:自增ID"`
	FieldName       string `json:"field_name" form:"field_name" binding:"required" required_tips:"字段名称必填" gorm:"column:field_name;size:64;uniqueIndex:unique_index;not null;default:'';comment:表的字段名称"`
	Label           string `json:"label" form:"label" binding:"required" required_tips:"字段中文名必填" gorm:"column:label;size:64;not null;default:'';comment:表的字段中文名称"`
	ModelId         uint   `json:"model_id" form:"model_id" binding:"required" required_tips:"请选择" gorm:"column:model_id;size:32;uniqueIndex:unique_index;not null;default:0;comment:所属模型ID"`
	Type            string `json:"type" form:"type" binding:"required" required_tips:"请选择" gorm:"column:type;size:20;not null;default:'';comment:表的字段类型"`
	Length          uint   `json:"length" form:"length" binding:"number" number_tips:"请输入数字" gorm:"column:length;size:32;not null;default:0;comment:表的字段长度"`
	DecimalLength   uint16 `json:"decimal_length" form:"decimal_length"  gorm:"column:decimal_length;size:8;not null;default:0;comment:表的字段为decimal类型时的小数位数"`
	IsNull          uint16 `json:"is_null" form:"is_null"  gorm:"column:is_null;size:8;not null;default:10;comment:是否为空：10=是，20=否"`
	Note            string `json:"note" form:"note"  gorm:"column:note;size:255;not null;default:'';comment:表的字段备注说明"`
	DefaultValue    string `json:"default_value" form:"default_value"  gorm:"column:default_value;size:255;not null;default:'';comment:数据默认值"`
	IsAutoIncrement uint16 `json:"is_auto_increment" form:"is_auto_increment"  gorm:"column:is_auto_increment;size:8;not null;default:20;comment:是否自动递增：10=是，20=否"`
	IsLabel         uint16 `json:"is_label" form:"is_label"  gorm:"column:is_label;size:8;not null;default:20;comment:是否为标题字段：10=是，20=否"`
	IsSigned        uint16 `json:"is_signed" form:"is_signed"  gorm:"column:is_signed;size:8;not null;default:10;comment:是否可为负数：10=是，20=否"`
	IsShow          uint16 `json:"is_show" form:"is_show"  gorm:"column:is_show;size:8;not null;default:10;comment:默认列表中显示：10=显示，20=不显示"`
	IsFixed         uint16 `json:"is_fixed" form:"is_fixed"  gorm:"column:is_fixed;size:8;not null;default:20;comment:默认列表中固定：10=固定，20=不固定"`
	ColumnWidth     uint   `json:"column_width" form:"column_width"  gorm:"column:column_width;size:32;not null;default:150;comment:列表中默认显示宽度：0表示不限"`
	IsFilter        uint16 `json:"is_filter" form:"is_filter"  gorm:"column:is_filter;size:8;not null;default:10;comment:列表中是否可为筛选条件：10=是，20=否"`
	IsCode          uint16 `json:"is_code" form:"is_code"  gorm:"column:is_code;size:8;not null;default:20;comment:是否显示文本源码，10=是，20=否"`
	SortNum         uint   `json:"sort_num" form:"sort_num"  gorm:"column:sort_num;size:32;not null;default:0;comment:表单/列表中字段的排列顺序(小在前)"`
	Status          uint16 `json:"status" form:"status"  gorm:"column:status;size:8;not null;default:10;comment:状态：10=开启，20=禁用"`
}

// DataModelFieldForm 提交的表单数据
type DataModelFieldForm struct {
	Data *ModelField `json:"data" form:"data" binding:"required" required_tips:"参数data不能为空"`
}
