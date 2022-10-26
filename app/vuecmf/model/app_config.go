package model

// AppConfig 应用配置 模型结构
type AppConfig struct {
	Id                uint   `json:"id" form:"id"  gorm:"column:id;primaryKey;autoIncrement;size:32;not null;comment:自增ID"`
	AppName           string `json:"app_name" form:"app_name" binding:"required" required_tips:"应用名称必填" gorm:"column:app_name;size:64;uniqueIndex:unique_index;not null;default:'';comment:应用名称"`
	LoginEnable       uint16 `json:"login_enable" form:"login_enable"  gorm:"column:login_enable;size:8;not null;default:10;comment:是否开启登录验证, 10=是，20=否"`
	AuthEnable        uint16 `json:"auth_enable" form:"auth_enable"  gorm:"column:auth_enable;size:8;not null;default:10;comment:是否开启权限验证, 10=是，20=否"`
	ExclusionUrl      string `json:"exclusion_url" form:"exclusion_url" gorm:"column:exclusion_url;size:2000;not null;default:'';comment:排除验证的URL,多个用英文逗号分隔"`
	Type              uint16 `json:"type" form:"type"  gorm:"column:type;size:8;not null;default:20;comment:类型：10=内置，20=扩展"`
	Status            uint16 `json:"status" form:"status"  gorm:"column:status;size:8;not null;default:10;comment:状态：10=开启，20=禁用"`
}

// DataAppConfigForm 提交的表单数据
type DataAppConfigForm struct {
	Data *AppConfig `json:"data" form:"data" binding:"required" required_tips:"参数data不能为空"`
}

//appNameForm 应用名表单
type appNameForm struct {
	AppName string  `json:"app_name" form:"app_name" binding:"required" required_tips:"应用名不能为空"`
}

//DataAppNameForm 提交应用名表单
type DataAppNameForm struct {
	Data *appNameForm `json:"data" form:"data" binding:"required" required_tips:"参数data不能为空"`
}
