//+----------------------------------------------------------------------
// | Copyright (c) 2024 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: tulihua2004@126.com
// +----------------------------------------------------------------------

package service

import (
	"errors"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/vuecmf/vuecmf-go/v3/app"
	"github.com/vuecmf/vuecmf-go/v3/app/vuecmf/helper"
	"github.com/vuecmf/vuecmf-go/v3/app/vuecmf/model"
	"gorm.io/gorm"
	"log"
	"strconv"
	"strings"
	"sync"
)

// AuthService 权限管理
type AuthService struct {
	Enforcer *casbin.Enforcer
}

var authOnce sync.Once
var authInstance *AuthService

// Auth 获取授权组件实例
func Auth() *AuthService {
	authOnce.Do(func() {
		var enf *casbin.Enforcer
		a, err := gormadapter.NewAdapterByDBWithCustomTable(app.Db, &model.Rules{}, TableName("rules"))
		if err != nil {
			log.Fatalln("初始化权限异常：" + err.Error())
			return
		}

		enf, err = casbin.NewEnforcer("config/tauthz-rbac-model.conf", a)
		if err != nil {
			log.Fatalln("读取权限配置文件异常：" + err.Error())
			return
		}
		authInstance = &AuthService{
			Enforcer: enf,
		}
	})

	return authInstance
}

// AddRolesForUser 给指定用户添加角色
//
//	参数：
//		username 用户名
//		roleIdList 角色ID列表
func (svc *AuthService) AddRolesForUser(username string, roleIdList []int) (bool, error) {
	err := app.Db.Transaction(func(tx *gorm.DB) error {
		//先清除移除历史角色
		_, err := svc.DelAllRolesForUser(username)
		if err != nil {
			return err
		}

		rolesList := Roles().GetRoleNameList(roleIdList)
		appNameList := AppConfig().GetAuthAppList()
		for _, appName := range appNameList {
			for _, roleName := range rolesList {
				_, err2 := svc.Enforcer.AddRoleForUser(username, roleName, appName)
				if err2 != nil {
					err = err2
					break
				}
			}
		}

		return err
	})

	if err != nil {
		return false, err
	}
	return true, nil
}

// DelRolesForUser 批量删除指定用户下的角色
//
//	参数：
//		username 用户名
//		roles 角色列表
func (svc *AuthService) DelRolesForUser(username string, roles []string) (bool, error) {
	err := app.Db.Transaction(func(tx *gorm.DB) error {
		appNameList := AppConfig().GetAuthAppList()
		for _, appName := range appNameList {
			for _, role := range roles {
				_, err2 := svc.Enforcer.DeleteRoleForUser(username, role, appName)
				if err2 != nil {
					return err2
				}
			}
		}

		return nil
	})

	if err != nil {
		return false, err
	}
	return true, nil
}

// DelAllRolesForUser 删除用户的所有角色
//
//	参数：
//		username 用户名
func (svc *AuthService) DelAllRolesForUser(username string) (bool, error) {
	err := app.Db.Transaction(func(tx *gorm.DB) error {
		appNameList := AppConfig().GetAuthAppList()
		for _, appName := range appNameList {
			_, err2 := svc.Enforcer.DeleteRolesForUser(username, appName)
			if err2 != nil {
				return err2
			}
		}

		return nil
	})

	if err != nil {
		return false, err
	}
	return true, nil
}

// AddUsersForRole 给角色分配用户
//
//	参数：
//		role 角色名
//		username 用户名列表
func (svc *AuthService) AddUsersForRole(role string, username []string) (bool, error) {
	err := app.Db.Transaction(func(tx *gorm.DB) error {
		//先取出角色下原有所有用户
		oldUsers, err := svc.GetUsers(role)
		if err != nil {
			return err
		}
		//取出需要删除的用户
		delUserList := helper.MinusStrList(oldUsers, username)

		roleArr := []string{role}

		appNameList := AppConfig().GetAuthAppList()
		for _, appName := range appNameList {
			//删除用户
			for _, user := range delUserList {
				_, err2 := svc.Enforcer.DeleteRoleForUser(user, role, appName)
				if err2 != nil {
					return err2
				}
			}

			//添加用户
			for _, user := range username {
				_, err2 := svc.Enforcer.AddRolesForUser(user, roleArr, appName)
				if err2 != nil {
					return err2
				}
			}
		}

		return nil
	})

	if err != nil {
		return false, err
	}
	return true, nil
}

// DelUsersForRole 批量删除指定角色下的用户
//
//	参数：
//		role 角色名
//		userIdList 用户ID列表
func (svc *AuthService) DelUsersForRole(role string, userIdList []int) (bool, error) {
	if len(userIdList) == 0 {
		return false, errors.New("该角色(" + role + ")没有分配用户")
	}

	username := Admin().GetUserNames(userIdList)

	err := app.Db.Transaction(func(tx *gorm.DB) error {
		appNameList := AppConfig().GetAuthAppList()
		for _, appName := range appNameList {
			for _, user := range username {
				_, err2 := svc.Enforcer.DeleteRoleForUser(user, role, appName)
				if err2 != nil {
					return err2
				}
			}
		}

		return nil
	})

	if err != nil {
		return false, err
	}
	return true, nil
}

// AddPermission 根据动作ID 给用户或角色分配权限
//
//	参数：
//		userOrRole 用户名或角色名
//		actionIdList 动作ID列表
func (svc *AuthService) AddPermission(userOrRole string, actionIdList string) (bool, error) {
	actionIdArr := strings.Split(actionIdList, ",")
	var actionPathArr []string
	DbTable("model_action").Select("api_path").
		Where("id in ?", actionIdArr).
		Where("status = 10").
		Find(&actionPathArr)

	err := app.Db.Transaction(func(tx *gorm.DB) error {
		//先清空原有权限
		_, err2 := svc.Enforcer.DeletePermissionsForUser(userOrRole)
		if err2 != nil {
			return err2
		}

		//再解析出路径中的控制器及动作，并分配权限
		for _, path := range actionPathArr {
			arr := strings.Split(strings.Trim(path, "/"), "/")
			if len(arr) < 2 {
				continue
			}
			appName := arr[0]
			controller := arr[1]
			action := "index"
			if len(arr) >= 3 && arr[2] != "" {
				action = arr[2]
			}
			_, err2 = svc.Enforcer.AddPermissionForUser(userOrRole, appName, controller, action)
			if err2 != nil {
				return err2
			}
		}
		return nil
	})

	if err != nil {
		return false, err
	}
	return true, nil
}

// DelPermission 根据动作ID 删除用户或角色的权限
//
//	参数：
//		userOrRole 用户名或角色名
//		actionIdList 动作ID列表
func (svc *AuthService) DelPermission(userOrRole string, actionIdList string) (bool, error) {
	actionIdArr := strings.Split(actionIdList, ",")
	var actionPathArr []string
	DbTable("model_action").Select("api_path").
		Where("id in ?", actionIdArr).
		Where("status = 10").
		Find(&actionPathArr)

	err := app.Db.Transaction(func(tx *gorm.DB) error {
		//再解析出路径中的控制器及动作，并分配权限
		for _, path := range actionPathArr {
			arr := strings.Split(strings.Trim(path, "/"), "/")
			if len(arr) < 2 {
				continue
			}
			appName := arr[0]
			controller := arr[1]
			action := "index"

			if len(arr) >= 3 && arr[2] != "" {
				action = arr[2]
			}
			_, err2 := svc.Enforcer.DeletePermissionForUser(userOrRole, appName, controller, action)
			if err2 != nil {
				return err2
			}
		}

		return nil
	})

	if err != nil {
		return false, err
	}
	return true, nil
}

// GetPermissions 获取(用户或角色)所有权限ID列表
//
//	参数：
//		userOrRole 用户名或角色名
//		isSuper 是否为超级管理员
func (svc *AuthService) GetPermissions(userOrRole string, isSuper uint16) (map[string][]string, error) {
	if userOrRole == "" {
		return nil, errors.New("用户或角色不能为空")
	}

	var res = make(map[string][]string)

	type action struct {
		Id    string
		Label string
	}

	if isSuper == 10 {
		//超级管理员拥有所有权限
		var actionList []action

		DbTable("model_action", "MA").Select("MA.id, MC.label").
			Joins("left join " + TableName("model_config") + " MC on MA.model_id = MC.id").
			Where("MA.status = 10").
			Where("MC.status = 10").
			Find(&actionList)

		for _, item := range actionList {
			res[item.Label] = append(res[item.Label], item.Id)
		}

	} else {
		appList := AppConfig().GetAuthAppList()
		for _, appName := range appList {
			data, err := svc.Enforcer.GetImplicitPermissionsForUser(userOrRole, appName)
			if err != nil {
				return nil, err
			}

			var pathList []string //API请求地址列表
			n := 0

			for _, val := range data {
				pathList = append(pathList, "/"+val[1]+"/"+val[2]+"/"+val[3])
				if val[3] == "index" {
					pathList = append(pathList, "/"+val[1]+"/"+val[2])
				}
				n++
				if n%100 == 0 {
					var actionList []action
					DbTable("model_action", "MA").Select("MA.id, MC.label").
						Joins("left join "+TableName("model_config")+" MC ON MA.model_id = MC.id").
						Joins("left join "+TableName("app_config")+" AC on MC.app_id = AC.id").
						Where("AC.app_name = ?", appName).
						Where("MA.api_path in ?", pathList).
						Where("MA.status = 10").
						Where("MC.status = 10").
						Where("AC.status = 10").
						Find(&actionList)

					for _, item := range actionList {
						res[item.Label] = append(res[item.Label], item.Id)
					}
					pathList = nil
				}
			}

			if pathList != nil {
				var actionList []action
				DbTable("model_action", "MA").Select("MA.id, MC.label").
					Joins("left join "+TableName("model_config")+" MC ON MA.model_id = MC.id").
					Joins("left join "+TableName("app_config")+" AC on MC.app_id = AC.id").
					Where("AC.app_name = ?", appName).
					Where("MA.api_path in ?", pathList).
					Where("MA.status = 10").
					Where("MC.status = 10").
					Where("AC.status = 10").
					Find(&actionList)

				for _, item := range actionList {
					res[item.Label] = append(res[item.Label], item.Id)
				}
			}
		}

	}

	return res, nil
}

// GetPermissionsForModelLabel 获取指定模型的权限ID列表
/*func (au *auth) GetPermissionsForModelLabel(userOrRole string, isSuper interface{}, modelLabel string) ([]string, error) {
	res, err := au.GetPermissions(userOrRole, isSuper)
	if err != nil {
		return nil, err
	}
	return res[modelLabel], nil
}*/

// GetUsers 获取指定角色下所有用户
//
//	参数：
//		role 角色名
func (svc *AuthService) GetUsers(role string) ([]string, error) {
	if role == "" {
		return nil, errors.New("角色不能为空")
	}
	return svc.Enforcer.GetUsersForRole(role, "vuecmf")
}

// GetRoles 获取指定用户名下所有角色
//
//	参数：
//		username 用户名
func (svc *AuthService) GetRoles(username string) ([]int, error) {
	if username == "" {
		return nil, errors.New("用户名不能为空")
	}
	roleNameList, err := svc.Enforcer.GetRolesForUser(username, "vuecmf")
	if err != nil {
		return nil, err
	}
	return Roles().GetRoleIdList(roleNameList), nil
}

type roleList struct {
	Key      uint   `json:"key"`
	Label    string `json:"label"`
	Disabled bool   `json:"disabled"`
}

// GetAllRoles 获取所有角色列表
func (svc *AuthService) GetAllRoles(roleName string, isSuper uint16) interface{} {
	var result []roleList
	query := DbTable("roles").Select("id `key`, role_name label, false disabled").
		Where("status = 10")

	if isSuper != 10 && roleName != "" {
		var pid int
		DbTable("roles").Select("id").
			Where("status = 10").
			Where("role_name = ?", roleName).
			Find(&pid)
		pidStr := strconv.Itoa(pid)
		query.Where(" id = ? or pid = ? or id_path like ? or id_path like ? or id_path like ?", pid, pid, pidStr+",%", "%,"+pidStr+",%", "%,"+pidStr)
	}

	query.Find(&result)
	return result
}

// GetRolesForUser 获取指定用户下所有角色名称
//
//	参数：
//		userName 用户名
func (svc *AuthService) GetRolesForUser(userName string) ([]string, error) {
	return svc.Enforcer.GetRolesForUser(userName, "vuecmf")
}

// UpdateRoles 更新权限的角色名称
//
//	参数：
//		oldRoleName 原角色名
//		newRoleName 新角色名
func (svc *AuthService) UpdateRoles(oldRoleName, newRoleName string) error {
	DbTable("rules").Where("ptype = 'p'").
		Where("v0 = ?", oldRoleName).
		Update("v0", newRoleName)

	DbTable("rules").Where("ptype = 'g'").
		Where("v1 = ?", oldRoleName).
		Update("v1", newRoleName)
	return nil
}

// UpdateUser 更新权限的用户名
//
//	参数：
//		oldUserName 原用户名
//		newUserName 新用户名
func (svc *AuthService) UpdateUser(oldUserName, newUserName string) error {
	res := DbTable("rules").Where("v0 = ?", oldUserName).
		Update("v0", newUserName)

	return res.Error
}
