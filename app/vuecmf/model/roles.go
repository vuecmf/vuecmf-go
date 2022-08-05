package model



// Roles 角色 模型结构
type Roles struct {
	Base
	RoleName string `json:"role_name" gorm:"column:role_name;size:64;uniqueIndex:unique_index;not null;default:;comment:用户的角色名称"`
	AppName string `json:"app_name" gorm:"column:app_name;size:64;uniqueIndex:unique_index;not null;default:;comment:角色所属应用名称"`
	Pid uint `json:"pid" gorm:"column:pid;size:11;not null;default:0;comment:父级ID"`
	IdPath string `json:"id_path" gorm:"column:id_path;size:255;not null;default:;comment:角色ID层级路径"`
	Remark string `json:"remark" gorm:"column:remark;size:255;not null;default:;comment:角色的备注信息"`
	
	Children *RolesTree `json:"children" gorm:"-"`
}


var rolesModel *Roles

// RolesModel 获取Roles模型实例
func RolesModel() *Roles {
	if rolesModel == nil {
		rolesModel = &Roles{}
	}
	return rolesModel
}

type RolesTree []*Roles

// ToTree 将列表数据转换树形结构
func (m *Roles) ToTree(data []*Roles) RolesTree {
	treeData := make(map[uint]*Roles)
	for _, val := range data {
		treeData[val.Id] = val
	}

	var treeList RolesTree

	for _, item := range treeData {
		if item.Pid == 0 {
			treeList = append(treeList, item)
			continue
		}
		if pItem, ok := treeData[item.Pid]; ok {
			if pItem.Children == nil {
				children := RolesTree{item}
				pItem.Children = &children
				continue
			}
			*pItem.Children = append(*pItem.Children, item)
		}
	}

	return treeList

}
