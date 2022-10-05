package controller

import (
	"github.com/vuecmf/vuecmf-go/app/route"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/model"
	
)

type AppConfig struct {
    Base
}

func init() {
	appConfig := &AppConfig{}
    appConfig.TableName = "app_config"
    appConfig.Model = &model.AppConfig{}
    appConfig.ListData = &[]model.AppConfig{}
    appConfig.SaveForm = &model.DataAppConfigForm{}
    appConfig.FilterFields = []string{"app_name"}

    route.Register(appConfig, "POST", "vuecmf")
}
