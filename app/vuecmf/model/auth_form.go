//+----------------------------------------------------------------------
// | Copyright (c) 2024 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: tulihua2004@126.com
// +----------------------------------------------------------------------

package model

type addRoleForm struct {
	Username   string `json:"username" form:"username" binding:"required" required_tips:"用户名不能为空"`
	RoleIdList []int  `json:"role_id_list" form:"role_id_list"`
}

// DataAddRoleForm 添加角色表单
type DataAddRoleForm struct {
	Data *addRoleForm `json:"data" form:"data" binding:"required" required_tips:"参数data不能为空"`
}

type delRoleForm struct {
	Username string   `json:"username" form:"username" binding:"required" required_tips:"用户名不能为空"`
	RoleName []string `json:"role_name" form:"role_name" binding:"required" required_tips:"角色不能为空"`
}

// DataDelRoleForm 删除角色表单
type DataDelRoleForm struct {
	Data *delRoleForm `json:"data" form:"data" binding:"required" required_tips:"参数data不能为空"`
}

type permissionForm struct {
	RoleName string `json:"role_name" form:"role_name" binding:"required" required_tips:"角色不能为空"`
	ActionId string `json:"action_id" form:"action_id" binding:"required" required_tips:"请选择功能项"`
}

// DataPermissionForm 添加/删除角色的权限表单
type DataPermissionForm struct {
	Data *permissionForm `json:"data" form:"data" binding:"required" required_tips:"参数data不能为空"`
}

type roleForm struct {
	RoleName string `json:"role_name" form:"role_name" binding:"required" required_tips:"角色不能为空"`
}

// DataRoleForm 角色权限表单
type DataRoleForm struct {
	Data *roleForm `json:"data" form:"data" binding:"required" required_tips:"参数data不能为空"`
}

type usernameForm struct {
	Username string `json:"username" form:"username" binding:"required" required_tips:"用户名不能为空"`
}

// DataUsernameForm 用户权限表单
type DataUsernameForm struct {
	Data *usernameForm `json:"data" form:"data" binding:"required" required_tips:"参数data不能为空"`
}

type userPermissionForm struct {
	Username string `json:"username" form:"username" binding:"required" required_tips:"用户名不能为空"`
	ActionId string `json:"action_id" form:"action_id" binding:"required" required_tips:"请选择功能项"`
}

// DataUserPermissionForm 添加用户的权限表单
type DataUserPermissionForm struct {
	Data *userPermissionForm `json:"data" form:"data" binding:"required" required_tips:"参数data不能为空"`
}
