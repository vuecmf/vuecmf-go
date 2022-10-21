// Package service
//+----------------------------------------------------------------------
// | Copyright (c) 2022 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: vuecmf <tulihua2004@126.com>
// +----------------------------------------------------------------------
package service

import (
	"errors"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/model"
	"gorm.io/gorm"
	"log"
	"strings"
)

// auth 权限管理
type auth struct {
	Enforcer *casbin.Enforcer
}

var authInstance *auth

// Auth 获取授权组件实例
func Auth() *auth {
	if authInstance == nil {
		var enf *casbin.Enforcer
		a, err := gormadapter.NewAdapterByDBWithCustomTable(Db, &model.Rules{}, NS.TableName("rules"))
		if err != nil {
			log.Fatalln("初始化权限异常：" + err.Error())
			return nil
		}

		enf, err = casbin.NewEnforcer("config/tauthz-rbac-model.conf", a)
		if err != nil {
			log.Fatalln("读取权限配置文件异常：" + err.Error())
			return nil
		}
		authInstance = &auth{
			Enforcer: enf,
		}
	}
	return authInstance
}

// AddRolesForUser 给指定用户添加角色
func (au *auth) AddRolesForUser(username string, roles []string, appName string) (bool, error) {
	if appName == "" {
		appName = "vuecmf"
	}

	err := Db.Transaction(func(tx *gorm.DB) error {
		_, err2 := au.Enforcer.AddRolesForUser(username, roles, appName)
		return err2
	})

	if err != nil {
		return false, err
	}
	return true, nil
}

// DelRolesForUser 批量删除指定用户下的角色
func (au *auth) DelRolesForUser(username string, roles []string, appName string) (bool, error) {
	if appName == "" {
		appName = "vuecmf"
	}

	err := Db.Transaction(func(tx *gorm.DB) error {
		for _, role := range roles {
			_, err2 := au.Enforcer.DeleteRoleForUser(username, role, appName)
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

// DelAllRolesForUser 删除用户的所有角色
func (au *auth) DelAllRolesForUser(username string, appName string) (bool, error) {
	if appName == "" {
		appName = "vuecmf"
	}

	err := Db.Transaction(func(tx *gorm.DB) error {
		_, err2 := au.Enforcer.DeleteRolesForUser(username, appName)
		if err2 != nil {
			return err2
		}
		return nil
	})

	if err != nil {
		return false, err
	}
	return true, nil
}

// AddUsersForRole 给角色分配用户
func (au *auth) AddUsersForRole(role string, username []string, appName string) (bool, error) {
	if appName == "" {
		appName = "vuecmf"
	}

	err := Db.Transaction(func(tx *gorm.DB) error {
		roleArr := []string{role}
		for _, user := range username {
			_, err2 := au.Enforcer.AddRolesForUser(user, roleArr, appName)
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

// DelUsersForRole 批量删除指定角色下的用户
func (au *auth) DelUsersForRole(role string, username []string, appName string) (bool, error) {
	if appName == "" {
		appName = "vuecmf"
	}

	err := Db.Transaction(func(tx *gorm.DB) error {
		for _, user := range username {
			_, err2 := au.Enforcer.DeleteRoleForUser(user, role, appName)
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

// AddPermission 根据动作ID 给用户或角色分配权限
func (au *auth) AddPermission(userOrRole string, actionIdList string) (bool, error) {
	actionIdArr := strings.Split(actionIdList, ",")
	var actionPathArr []string
	Db.Table(NS.TableName("model_action")).Select("api_path").
		Where("id in ?", actionIdArr).
		Where("status = 10").
		Find(&actionPathArr)

	err := Db.Transaction(func(tx *gorm.DB) error {
		//先清空原有权限
		_, err2 := au.Enforcer.DeletePermissionsForUser(userOrRole)
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
			_, err2 = au.Enforcer.AddPermissionForUser(userOrRole, appName, controller, action)
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
func (au *auth) DelPermission(userOrRole string, actionIdList string) (bool, error) {
	actionIdArr := strings.Split(actionIdList, ",")
	var actionPathArr []string
	Db.Table(NS.TableName("model_action")).Select("api_path").
		Where("id in ?", actionIdArr).
		Where("status = 10").
		Find(&actionPathArr)

	err := Db.Transaction(func(tx *gorm.DB) error {
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
			_, err2 := au.Enforcer.DeletePermissionForUser(userOrRole, appName, controller, action)
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
func (au *auth) GetPermissions(userOrRole string, isSuper interface{}, appName string) (map[string][]string, error) {
	if userOrRole == "" {
		return nil, errors.New("用户或角色不能为空")
	}

	if appName == "" {
		appName = "vuecmf"
	}

	var res = make(map[string][]string)

	type action struct {
		Id    string
		Label string
	}

	if isSuper == 10 {
		//超级管理员拥有所有权限
		var actionList []action

		Db.Table(NS.TableName("model_action")+" MA").Select("MA.id, MC.label").
			Joins("left join "+NS.TableName("model_config")+" MC on MA.model_id = MC.id").
			Joins("left join "+NS.TableName("app_config")+" AC on MA.app_id = AC.id").
			Where("AC.app_name = ?", appName).
			Where("MA.status = 10").
			Where("MC.status = 10").
			Find(&actionList)

		for _, item := range actionList {
			res[item.Label] = append(res[item.Label], item.Id)
		}

	} else {
		data, err := au.Enforcer.GetImplicitPermissionsForUser(userOrRole, appName)
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
				Db.Table(NS.TableName("model_action")+" MA").Select("MA.id, MC.label").
					Joins("left join "+NS.TableName("model_config")+" MC ON MA.model_id = MC.id").
					Joins("left join "+NS.TableName("app_config")+" AC on MA.app_id = AC.id").
					Where("AC.app_name = ?", appName).
					Where("MA.api_path in ?", pathList).
					Where("MA.status = 10").
					Where("MC.status = 10").
					Find(&actionList)

				for _, item := range actionList {
					res[item.Label] = append(res[item.Label], item.Id)
				}
				pathList = nil
			}
		}

		if pathList != nil {
			var actionList []action
			Db.Table(NS.TableName("model_action")+" MA").Select("MA.id, MC.label").
				Joins("left join "+NS.TableName("model_config")+" MC ON MA.model_id = MC.id").
				Joins("left join "+NS.TableName("app_config")+" AC on MA.app_id = AC.id").
				Where("AC.app_name = ?", appName).
				Where("MA.api_path in ?", pathList).
				Where("MA.status = 10").
				Where("MC.status = 10").
				Find(&actionList)

			for _, item := range actionList {
				res[item.Label] = append(res[item.Label], item.Id)
			}
		}
	}

	return res, nil
}

// GetPermissionsForModelLabel 获取指定模型的权限ID列表
func (au *auth) GetPermissionsForModelLabel(userOrRole string, isSuper interface{}, modelLabel string, appName string) ([]string, error) {
	if appName == "" {
		appName = "vuecmf"
	}
	res, err := au.GetPermissions(userOrRole, isSuper, appName)
	if err != nil {
		return nil, err
	}
	return res[modelLabel], nil
}

// GetUsers 获取指定角色下所有用户
func (au *auth) GetUsers(role string, appName string) ([]string, error) {
	if role == "" {
		return nil, errors.New("角色不能为空")
	}
	if appName == "" {
		appName = "vuecmf"
	}
	return au.Enforcer.GetUsersForRole(role, appName)
}

// GetRoles 获取指定用户名下所有角色
func (au *auth) GetRoles(username string, appName string) ([]string, error) {
	if username == "" {
		return nil, errors.New("用户名不能为空")
	}
	if appName == "" {
		appName = "vuecmf"
	}
	return au.Enforcer.GetRolesForUser(username, appName)
}

type roleList struct {
	Key      uint   `json:"key"`
	Label    string `json:"label"`
	Disabled bool   `json:"disabled"`
}

// GetAllRoles 获取所有角色列表
func (au *auth) GetAllRoles() interface{} {
	var result []roleList
	Db.Table(NS.TableName("roles")).Select("id `key`, role_name label, false disabled").
		Where("status = 10").
		Find(&result)
	return result
}
