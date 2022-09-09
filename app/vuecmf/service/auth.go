package service

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/vuecmf/vuecmf-go/app"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/model"
	"gorm.io/gorm"
	"strings"
)

type auth struct {
	Enforcer *casbin.Enforcer
}

var au *auth

// Auth 获取授权组件实例
func Auth() *auth {
	if au == nil {
		db := app.Db("default")
		a, err := gormadapter.NewAdapterByDBWithCustomTable(db, &model.Rules{})
		var enf *casbin.Enforcer
		if err == nil {
			enf, err = casbin.NewEnforcer("../config/tauthz-rbac-model.conf", a)
			if err != nil{
				enf = nil
			}
		}
		au = &auth{
			Enforcer: enf,
		}
	}
	return au
}

// AddRolesForUser 给指定用户添加角色
func (au *auth) AddRolesForUser(username string, roles []string, appName string) (bool, error) {
	err := db.Transaction(func(tx *gorm.DB) error {
		//先清除原角色记录
		/*_, err2 := au.Enforcer.DeleteRolesForUser(username, appName)
		if err2 != nil {
			return err2
		}*/
		//再添加新角色
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
	err := db.Transaction(func(tx *gorm.DB) error {
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


// AddUsersForRole 给角色分配用户
func (au *auth) AddUsersForRole(role string, username []string, appName string) (bool, error) {
	err := db.Transaction(func(tx *gorm.DB) error {
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
	err := db.Transaction(func(tx *gorm.DB) error {
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
	db.Table(ns.TableName("model_action")).Select("api_path").
		Where("id in ?", actionIdArr).
		Where("status = 10").
		Find(&actionPathArr)

	err := db.Transaction(func(tx *gorm.DB) error {
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
			if arr[2] != "" {
				action = arr[2]
			}
			_,err2 = au.Enforcer.AddPermissionForUser(userOrRole, appName, controller, action)
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
func (au *auth) DelPermission(userOrRole string, actionIdList string) (bool, error)  {
	actionIdArr := strings.Split(actionIdList, ",")
	var actionPathArr []string
	db.Table(ns.TableName("model_action")).Select("api_path").
		Where("id in ?", actionIdArr).
		Where("status = 10").
		Find(&actionPathArr)

	err := db.Transaction(func(tx *gorm.DB) error {
		//再解析出路径中的控制器及动作，并分配权限
		for _, path := range actionPathArr {
			arr := strings.Split(strings.Trim(path, "/"), "/")
			if len(arr) < 2 {
				continue
			}
			appName := arr[0]
			controller := arr[1]
			action := "index"
			if arr[2] != "" {
				action = arr[2]
			}
			_,err2 := au.Enforcer.DeletePermissionForUser(userOrRole, appName, controller, action)
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


func (au *auth) GetPermission() {

}


func (au *auth) GetUsers() {

}


func (au *auth) GetRoles() {

}


