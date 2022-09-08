package app

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/model"
)

func Auth() (*casbin.Enforcer, error) {
	db := Db("default")
	a, err := gormadapter.NewAdapterByDBWithCustomTable(db, &model.Rules{})
	if err != nil {
		return nil, err
	}
	return casbin.NewEnforcer("../config/tauthz-rbac-model.conf", a)
}
