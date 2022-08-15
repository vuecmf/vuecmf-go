package model



// ModelConfig 模型配置 模型结构
type ModelConfig struct {
	Id uint `json:"id" form:"id"  gorm:"column:id;primaryKey;autoIncrement;size:11;not null;default:0;comment:自增ID"`
	SearchFieldId string `json:"search_field_id" form:"search_field_id"  gorm:"column:search_field_id;size:255;not null;default:;comment:搜索字段ID，多个用逗号分隔"`
	Status uint `json:"status" form:"status"  gorm:"column:status;size:4;not null;default:10;comment:状态：10=开启，20=禁用"`
	TableName string `json:"table_name" form:"table_name" binding:"required" required_tips:"表名必填" gorm:"column:table_name;size:64;uniqueIndex:unique_index;not null;default:;comment:模型对应的表名(不含表前缘)"`
	Label string `json:"label" form:"label" binding:"required" required_tips:"模型标签必填" gorm:"column:label;size:64;not null;default:;comment:模型标签"`
	ComponentTpl string `json:"component_tpl" form:"component_tpl" binding:"required" required_tips:"请选择" gorm:"column:component_tpl;size:255;not null;default:;comment:组件模板"`
	DefaultActionId uint `json:"default_action_id" form:"default_action_id"  gorm:"column:default_action_id;size:11;not null;default:0;comment:默认动作ID"`
	Type uint `json:"type" form:"type"  gorm:"column:type;size:4;not null;default:20;comment:类型：10=内置，20=扩展"`
	IsTree uint `json:"is_tree" form:"is_tree"  gorm:"column:is_tree;size:4;not null;default:20;comment:是否为目录树：10=是，20=否"`
	Remark string `json:"remark" form:"remark"  gorm:"column:remark;size:255;not null;default:;comment:模型对应表的备注"`
	
}

// DataModelConfigForm 提交的表单数据
type DataModelConfigForm struct {
    Data *ModelConfig `json:"data" form:"data"`
}