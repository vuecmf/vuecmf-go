package model



// Menu 菜单 模型结构
type Menu struct {
	Base
	Title string `json:"title" gorm:"column:title;size:64;not null;default:;comment:菜单标题"`
	Icon string `json:"icon" gorm:"column:icon;size:32;not null;default:;comment:菜单图标"`
	Pid int `json:"pid" gorm:"column:pid;size:11;not null;default:0;comment:父级ID"`
	ModelId int `json:"model_id" gorm:"column:model_id;size:11;not null;default:0;comment:模型ID"`
	Type int16 `json:"type" gorm:"column:type;size:4;not null;default:20;comment:类型：10=内置，20=扩展"`
	SortNum int `json:"sort_num" gorm:"column:sort_num;size:11;not null;default:0;comment:菜单的排列顺序(小在前)"`
	Children *MenuTree `json:"children"`
}

type MenuTree []*Menu

