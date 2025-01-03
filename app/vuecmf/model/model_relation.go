//+----------------------------------------------------------------------
// | Copyright (c) 2024 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: tulihua2004@126.com
// +----------------------------------------------------------------------

package model

// ModelRelation 模型关联 模型结构
type ModelRelation struct {
	Id                  uint   `json:"id" form:"id"  gorm:"column:id;primaryKey;autoIncrement;size:32;not null;comment:自增ID"`
	ModelId             uint   `json:"model_id" form:"model_id" binding:"required" required_tips:"请选择" gorm:"column:model_id;size:32;not null;default:0;comment:所属模型ID"`
	ModelFieldId        uint   `json:"model_field_id" form:"model_field_id" binding:"required" required_tips:"请选择" gorm:"column:model_field_id;size:32;uniqueIndex:unique_index;not null;default:0;comment:模型字段ID"`
	RelationModelId     uint   `json:"relation_model_id" form:"relation_model_id" binding:"required" required_tips:"请选择" gorm:"column:relation_model_id;size:32;not null;default:0;comment:关联模型ID"`
	RelationFieldId     uint   `json:"relation_field_id" form:"relation_field_id" binding:"required" required_tips:"请选择" gorm:"column:relation_field_id;size:32;uniqueIndex:unique_index;not null;default:0;comment:关联模型字段ID"`
	RelationShowFieldId string `json:"relation_show_field_id" form:"relation_show_field_id" binding:"required" required_tips:"请选择" gorm:"column:relation_show_field_id;size:255;not null;default:'';comment:关联模型显示字段ID,多个逗号分隔，全部用*"`
	RelationFilter      string `json:"relation_filter" form:"relation_filter"  gorm:"column:relation_filter;size:255;not null;default:'';comment:关联过滤条件"`
	Status              uint16 `json:"status" form:"status"  gorm:"column:status;size:8;not null;default:10;comment:状态：10=开启，20=禁用"`
}

// DataModelRelationForm 提交的表单数据
type DataModelRelationForm struct {
	Data *ModelRelation `json:"data" form:"data" binding:"required" required_tips:"参数data不能为空"`
}
