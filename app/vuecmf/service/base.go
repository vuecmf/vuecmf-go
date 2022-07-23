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
	"github.com/vuecmf/vuecmf-go/app"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var db *gorm.DB
var ns schema.Namer

type base struct {
}

func init() {
	db = app.Db("default")
	ns = db.NamingStrategy
}
