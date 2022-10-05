package model



// AppModel 应用与模型关系 模型结构
type AppModel struct {
	Id uint `json:"id" form:"id"  gorm:"column:id;primaryKey;autoIncrement;size:32;not null;comment:自增ID"`
	AppId uint `json:"app_id" form:"app_id" binding:"required" required_tips:"请选择" gorm:"column:app_id;size:32;uniqueIndex:unique_index;not null;default:0;comment:应用配置ID"`
	ModelId uint `json:"model_id" form:"model_id" binding:"required" required_tips:"请选择" gorm:"column:model_id;size:32;uniqueIndex:unique_index;not null;default:0;comment:模型ID"`
	Status uint16 `json:"status" form:"status"  gorm:"column:status;size:8;not null;default:10;comment:状态：10=开启，20=禁用"`
	
}

// DataAppModelForm 提交的表单数据
type DataAppModelForm struct {
    Data *AppModel `json:"data" form:"data"`
}