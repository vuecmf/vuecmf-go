package model



// Menu 菜单 模型结构
type Menu struct {
	SortNum uint `json:"sort_num" form:"sort_num"  gorm:"column:sort_num;size:11;not null;default:0;comment:菜单的排列顺序(小在前)"`
	Status uint `json:"status" form:"status"  gorm:"column:status;size:4;not null;default:10;comment:状态：10=开启，20=禁用"`
	Title string `json:"title" form:"title" binding:"required" required_tips:"菜单标题必填" gorm:"column:title;size:64;not null;default:;comment:菜单标题"`
	Id uint `json:"id" form:"id"  gorm:"column:id;primaryKey;autoIncrement;size:11;not null;default:0;comment:自增ID"`
	Icon string `json:"icon" form:"icon"  gorm:"column:icon;size:32;not null;default:;comment:菜单图标"`
	Pid uint `json:"pid" form:"pid"  gorm:"column:pid;size:11;not null;default:0;comment:父级ID"`
	ModelId uint `json:"model_id" form:"model_id"  gorm:"column:model_id;size:11;not null;default:0;comment:模型ID"`
	Type uint `json:"type" form:"type"  gorm:"column:type;size:4;not null;default:20;comment:类型：10=内置，20=扩展"`
	
	Children *MenuTree `json:"children" gorm:"-"`
}

// DataMenuForm 提交的表单数据
type DataMenuForm struct {
    Data *Menu `json:"data" form:"data"`
}


var menuModel *Menu

// MenuModel 获取Menu模型实例
func MenuModel() *Menu {
	if menuModel == nil {
		menuModel = &Menu{}
	}
	return menuModel
}

type MenuTree []*Menu

// ToTree 将列表数据转换树形结构
func (m *Menu) ToTree(data []*Menu) MenuTree {
	treeData := make(map[uint]*Menu)
	for _, val := range data {
		treeData[val.Id] = val
	}

	var treeList MenuTree

	for _, item := range treeData {
		if item.Pid == 0 {
			treeList = append(treeList, item)
			continue
		}
		if pItem, ok := treeData[item.Pid]; ok {
			if pItem.Children == nil {
				children := MenuTree{item}
				pItem.Children = &children
				continue
			}
			*pItem.Children = append(*pItem.Children, item)
		}
	}

	return treeList

}
