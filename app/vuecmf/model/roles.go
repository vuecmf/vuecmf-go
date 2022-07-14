package model



// Roles 角色 模型结构
type Roles struct {
	Base
	RoleName string `json:"role_name" gorm:"column:role_name;size:64;uniqueIndex:unique_index;not null;default:;comment:用户的角色名称"`
	AppName string `json:"app_name" gorm:"column:app_name;size:64;uniqueIndex:unique_index;not null;default:;comment:角色所属应用名称"`
	Pid int `json:"pid" gorm:"column:pid;size:11;not null;default:0;comment:父级ID"`
	IdPath string `json:"id_path" gorm:"column:id_path;size:255;not null;default:;comment:角色ID层级路径"`
	Remark string `json:"remark" gorm:"column:remark;size:255;not null;default:;comment:角色的备注信息"`
	
}