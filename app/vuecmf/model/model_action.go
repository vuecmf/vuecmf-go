package model

// ModelAction 模型动作 模型结构
type ModelAction struct {
	Label      string `json:"label" form:"label" binding:"required" required_tips:"动作标签必填" gorm:"column:label;size:64;not null;default:;comment:动作标签"`
	ApiPath    string `json:"api_path" form:"api_path" binding:"required" required_tips:"后端请求地址必填" gorm:"column:api_path;size:255;not null;default:;comment:后端请求地址"`
	ModelId    uint   `json:"model_id" form:"model_id" binding:"required" required_tips:"请选择" gorm:"column:model_id;size:11;uniqueIndex:unique_index;not null;default:0;comment:所属模型ID"`
	ActionType string `json:"action_type" form:"action_type" binding:"required" required_tips:"请选择" gorm:"column:action_type;size:32;uniqueIndex:unique_index;not null;default:;comment:动作类型"`
	Id         uint   `json:"id" form:"id"  gorm:"column:id;primaryKey;autoIncrement;size:11;not null;default:0;comment:自增ID"`
	Status     uint   `json:"status" form:"status"  gorm:"column:status;size:4;not null;default:10;comment:状态：10=开启，20=禁用"`
}

// DataModelActionForm 提交的表单数据
type DataModelActionForm struct {
	Data *ModelAction `json:"data" form:"data"`
}

type apiMapForm struct {
	TableName  string `json:"table_name" form:"table_name" binding:"required" required_tips:"表名(table_name)不能为空"`
	ActionType string `json:"action_type" form:"action_type" binding:"required" required_tips:"动作类型(action_type)不能为空"`
}

type DataApiMapForm struct {
	Data *apiMapForm `json:"data" form:"data"`
}

type actionListForm struct {
	RoleName string `json:"role_name" form:"role_name"`
	AppName  string `json:"app_name" form:"app_name"`
}

// DataActionListForm 获取所有模型的动作列表表单
type DataActionListForm struct {
	Data *actionListForm `json:"data" form:"data"`
}
