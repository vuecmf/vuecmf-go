package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/vuecmf/vuecmf-go/app/route"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/model"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/service"
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

// GetModelList 获取应用下所有模型
func (ctrl *AppModel) GetModelList(c *gin.Context) {
	modelListForm := &model.DataModelListForm{}
	common(c, modelListForm, func() (interface{}, error) {
		return service.AppModel().GetModelList(modelListForm.Data.AppId)
	})
}
