package model



// Roles 角色 模型结构
type Roles struct {
	Id uint `json:"id" form:"id"  gorm:"column:id;primaryKey;autoIncrement;size:32;not null;comment:自增ID"`
	RoleName string `json:"role_name" form:"role_name" binding:"required" required_tips:"角色名称必填" gorm:"column:role_name;size:64;uniqueIndex:unique_index;not null;default:;comment:用户的角色名称"`
	AppName string `json:"app_name" form:"app_name"  gorm:"column:app_name;size:64;uniqueIndex:unique_index;not null;default:;comment:角色所属应用名称"`
	Pid uint `json:"pid" form:"pid"  gorm:"column:pid;size:32;not null;default:0;comment:父级ID"`
	IdPath string `json:"id_path" form:"id_path"  gorm:"column:id_path;size:255;not null;default:;comment:角色ID层级路径"`
	Remark string `json:"remark" form:"remark"  gorm:"column:remark;size:255;not null;default:;comment:角色的备注信息"`
	Status uint `json:"status" form:"status"  gorm:"column:status;size:8;not null;default:10;comment:状态：10=开启，20=禁用"`
	
	Children *RolesTree `json:"children" gorm:"-"`
}

// DataRolesForm 提交的表单数据
type DataRolesForm struct {
    Data *Roles `json:"data" form:"data"`
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
