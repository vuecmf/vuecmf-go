package model



// ModelIndex 模型索引 模型结构
type ModelIndex struct {
	ModelId uint `json:"model_id" form:"model_id" binding:"required" required_tips:"请选择" gorm:"column:model_id;size:11;not null;default:0;comment:所属模型ID"`
	ModelFieldId uint `json:"model_field_id" form:"model_field_id" binding:"required" required_tips:"请选择" gorm:"column:model_field_id;size:11;not null;default:0;comment:模型字段ID"`
	IndexType string `json:"index_type" form:"index_type" binding:"required" required_tips:"请选择" gorm:"column:index_type;size:32;not null;default:NORMAL;comment:索引类型： PRIMARY=主键，NORMAL=常规，UNIQUE=唯一，FULLTEXT=全文"`
	Id uint `json:"id" form:"id"  gorm:"column:id;primaryKey;autoIncrement;size:11;not null;default:0;comment:自增ID"`
	Status uint `json:"status" form:"status"  gorm:"column:status;size:4;not null;default:10;comment:状态：10=开启，20=禁用"`
	
}

// DataModelIndexForm 提交的表单数据
type DataModelIndexForm struct {
    Data *ModelIndex `json:"data" form:"data"`
}