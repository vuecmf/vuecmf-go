package model



// ModelIndex 模型索引 模型结构
type ModelIndex struct {
	Base
	ModelId int `json:"model_id" gorm:"column:model_id;size:11;not null;default:0;comment:所属模型ID"`
	ModelFieldId int `json:"model_field_id" gorm:"column:model_field_id;size:11;not null;default:0;comment:模型字段ID"`
	IndexType string `json:"index_type" gorm:"column:index_type;size:32;not null;default:NORMAL;comment:索引类型： PRIMARY=主键，NORMAL=常规，UNIQUE=唯一，FULLTEXT=全文"`
	
}