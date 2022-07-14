package model



// ModelAction 模型动作 模型结构
type ModelAction struct {
	Base
	Label string `json:"label" gorm:"column:label;size:64;not null;default:;comment:动作标签"`
	ApiPath string `json:"api_path" gorm:"column:api_path;size:255;not null;default:;comment:后端请求地址"`
	ModelId int `json:"model_id" gorm:"column:model_id;size:11;uniqueIndex:unique_index;not null;default:0;comment:所属模型ID"`
	ActionType string `json:"action_type" gorm:"column:action_type;size:32;uniqueIndex:unique_index;not null;default:;comment:动作类型"`
	
}