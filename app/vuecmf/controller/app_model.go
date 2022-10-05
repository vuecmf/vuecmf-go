package controller

import (
	"github.com/vuecmf/vuecmf-go/app/route"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/model"
	
)

type AppModel struct {
    Base
}

func init() {
	appModel := &AppModel{}
    appModel.TableName = "app_model"
    appModel.Model = &model.AppModel{}
    appModel.ListData = &[]model.AppModel{}
    appModel.SaveForm = &model.DataAppModelForm{}
    appModel.FilterFields = []string{""}

    route.Register(appModel, "POST", "vuecmf")
}
