package model

// Upload  模型结构
type Upload struct {
	Status uint `json:"status" form:"status"  gorm:"column:status;size:4;not null;default:10;comment:状态：10=开启，20=禁用"`
	Id     uint `json:"id" form:"id"  gorm:"column:id;primaryKey;autoIncrement;size:11;not null;default:0;comment:自增ID"`
}

// DataUploadForm 提交的表单数据
type DataUploadForm struct {
	Data *Upload `json:"data" form:"data"`
}
