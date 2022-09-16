package model

type addRoleForm struct {
	Username   string   `json:"username" form:"username" binding:"required" required_tips:"用户名不能为空"`
	RoleIdList []string `json:"role_id_list" form:"role_id_list"`
	AppName    string   `json:"app_name" form:"app_name"`
}

// DataAddRoleForm 添加角色表单
type DataAddRoleForm struct {
	Data *addRoleForm `json:"data" form:"data"`
}

type delRoleForm struct {
	Username string   `json:"username" form:"username" binding:"required" required_tips:"用户名不能为空"`
	RoleName []string `json:"role_name" form:"role_name" binding:"required" required_tips:"角色不能为空"`
	AppName  string   `json:"app_name" form:"app_name"`
}

// DataDelRoleForm 用户名表单
type DataDelRoleForm struct {
	Data *delRoleForm `json:"data" form:"data"`
}

type permissionForm struct {
	RoleName string `json:"role_name" form:"role_name" binding:"required" required_tips:"角色不能为空"`
	ActionId string `json:"action_id" form:"action_id" binding:"required" required_tips:"请选择功能项"`
}

// DataPermissionForm 添加/删除角色的权限表单
type DataPermissionForm struct {
	Data *permissionForm `json:"data" form:"data"`
}

type roleForm struct {
	RoleName string `json:"role_name" form:"role_name" binding:"required" required_tips:"角色不能为空"`
	AppName  string `json:"app_name" form:"app_name"`
}

// DataRoleForm 角色权限表单
type DataRoleForm struct {
	Data *roleForm `json:"data" form:"data"`
}

type usernameForm struct {
	Username string `json:"username" form:"username" binding:"required" required_tips:"用户名不能为空"`
	AppName  string `json:"app_name" form:"app_name"`
}

// DataUsernameForm 用户权限表单
type DataUsernameForm struct {
	Data *usernameForm `json:"data" form:"data"`
}

type userPermissionForm struct {
	Username string `json:"username" form:"username" binding:"required" required_tips:"用户名不能为空"`
	ActionId string `json:"action_id" form:"action_id" binding:"required" required_tips:"请选择功能项"`
}

// DataUserPermissionForm 添加用户的权限表单
type DataUserPermissionForm struct {
	Data *userPermissionForm `json:"data" form:"data"`
}
