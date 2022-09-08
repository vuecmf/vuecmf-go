// Package controller
//+----------------------------------------------------------------------
// | Copyright (c) 2022 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: vuecmf <tulihua2004@126.com>
// +----------------------------------------------------------------------
package controller

import (
	"github.com/vuecmf/vuecmf-go/app/route"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/model"
)

type Admin struct {
    Base
}

func init() {
	admin := &Admin{}
    admin.TableName = "admin"
    admin.Model = &model.Admin{}
    admin.listData = &[]model.Admin{}
    admin.saveForm = &model.DataAdminForm{}
    admin.filterFields = []string{"username","email","mobile","token"}

    route.Register(admin, "POST", "vuecmf")
}
