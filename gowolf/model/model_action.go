package model

import (
	"github.com/vuecmf/vuecmf-go/app"
)

type ModelAction struct {
	Id  int `gorm:"primaryKey"`
	Label string
	ApiPath string
	ModelId int `gorm:"index:action_type"`
	ActionType string `gorm:"index:action_type"`
	Status int8
}


func GetList() (list []*ModelAction) {
	//var list []ModelAction

	app.Db().First(&list)

	


	return
}