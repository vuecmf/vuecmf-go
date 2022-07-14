package model

// ModelConfig 模型配置 模型结构
type ModelConfig struct {
	Base
	TableName       string `json:"table_name" gorm:"column:table_name;size:64;uniqueIndex:unique_index;not null;default:;comment:模型对应的表名(不含表前缘)"`
	Label           string `json:"label" gorm:"column:label;size:64;not null;default:;comment:模型标签"`
	ComponentTpl    string `json:"component_tpl" gorm:"column:component_tpl;size:255;not null;default:;comment:组件模板"`
	DefaultActionId int    `json:"default_action_id" gorm:"column:default_action_id;size:11;not null;default:0;comment:默认动作ID"`
	SearchFieldId   string `json:"search_field_id" gorm:"column:search_field_id;size:255;not null;default:;comment:搜索字段ID，多个用逗号分隔"`
	Type            uint8  `json:"type" gorm:"column:type;size:4;not null;default:20;comment:类型：10=内置，20=扩展"`
	IsTree          uint8  `json:"is_tree" gorm:"column:is_tree;size:4;not null;default:20;comment:是否为目录树：10=是，20=否"`
	Remark          string `json:"remark" gorm:"column:remark;size:255;not null;default:;comment:模型对应表的备注"`
}
