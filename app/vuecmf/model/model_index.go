//+----------------------------------------------------------------------
// | Copyright (c) 2024 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: tulihua2004@126.com
// +----------------------------------------------------------------------

package model

// ModelIndex 模型索引 模型结构
type ModelIndex struct {
	Id           uint   `json:"id" form:"id"  gorm:"column:id;primaryKey;autoIncrement;size:32;not null;comment:自增ID"`
	ModelId      uint   `json:"model_id" form:"model_id" binding:"required" required_tips:"请选择" gorm:"column:model_id;size:32;not null;default:0;comment:所属模型ID"`
	ModelFieldId string `json:"model_field_id" form:"model_field_id" binding:"required" required_tips:"请选择" gorm:"column:model_field_id;size:255;not null;default:'';comment:模型字段ID"`
	IndexType    string `json:"index_type" form:"index_type" binding:"required" required_tips:"请选择" gorm:"column:index_type;size:32;not null;default:'NORMAL';comment:索引类型： PRIMARY=主键，NORMAL=常规，UNIQUE=唯一，FULLTEXT=全文"`
	Status       uint16 `json:"status" form:"status"  gorm:"column:status;size:8;not null;default:10;comment:状态：10=开启，20=禁用"`
}

// DataModelIndexForm 提交的表单数据
type DataModelIndexForm struct {
	Data *ModelIndex `json:"data" form:"data" binding:"required" required_tips:"参数data不能为空"`
}
