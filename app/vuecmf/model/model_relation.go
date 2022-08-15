package model



// ModelRelation 模型关联 模型结构
type ModelRelation struct {
	ModelId uint `json:"model_id" form:"model_id" binding:"required" required_tips:"请选择" gorm:"column:model_id;size:11;not null;default:0;comment:所属模型ID"`
	ModelFieldId uint `json:"model_field_id" form:"model_field_id" binding:"required" required_tips:"请选择" gorm:"column:model_field_id;size:11;uniqueIndex:unique_index;not null;default:0;comment:模型字段ID"`
	RelationModelId uint `json:"relation_model_id" form:"relation_model_id" binding:"required" required_tips:"请选择" gorm:"column:relation_model_id;size:11;not null;default:0;comment:关联模型ID"`
	RelationFieldId uint `json:"relation_field_id" form:"relation_field_id" binding:"required" required_tips:"请选择" gorm:"column:relation_field_id;size:11;uniqueIndex:unique_index;not null;default:0;comment:关联模型字段ID"`
	RelationShowFieldId string `json:"relation_show_field_id" form:"relation_show_field_id" binding:"required" required_tips:"请选择" gorm:"column:relation_show_field_id;size:255;not null;default:;comment:关联模型显示字段ID,多个逗号分隔，全部用*"`
	Id uint `json:"id" form:"id"  gorm:"column:id;primaryKey;autoIncrement;size:11;not null;default:0;comment:自增ID"`
	Status uint `json:"status" form:"status"  gorm:"column:status;size:4;not null;default:10;comment:状态：10=开启，20=禁用"`
	
}

// DataModelRelationForm 提交的表单数据
type DataModelRelationForm struct {
    Data *ModelRelation `json:"data" form:"data"`
}