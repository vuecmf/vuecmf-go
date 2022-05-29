package model



// ModelRelation 模型关联 模型结构
type ModelRelation struct {
	Base
	ModelId int `json:"model_id" gorm:"column:model_id;size:11;not null;default:0;comment:所属模型ID"`
	ModelFieldId int `json:"model_field_id" gorm:"column:model_field_id;size:11;uniqueIndex:unique_index;not null;default:0;comment:模型字段ID"`
	RelationModelId int `json:"relation_model_id" gorm:"column:relation_model_id;size:11;not null;default:0;comment:关联模型ID"`
	RelationFieldId int `json:"relation_field_id" gorm:"column:relation_field_id;size:11;uniqueIndex:unique_index;not null;default:0;comment:关联模型字段ID"`
	RelationShowFieldId string `json:"relation_show_field_id" gorm:"column:relation_show_field_id;size:255;not null;default:;comment:关联模型显示字段ID,多个逗号分隔，全部用*"`
	
}