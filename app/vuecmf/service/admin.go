package service

import (
	"github.com/vuecmf/vuecmf-go/app/vuecmf/helper"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/model"
)

// AdminService admin服务结构
type AdminService struct {
	*base
	TableName string 
}


// List 获取列表数据
// 		参数：params 查询参数
func (service *AdminService) List(params *helper.DataListParams) interface{} {
	var adminList []model.Admin

	//modelId := service.ModelConfig().GetModelId(tableName)

	page := &helper.Page{
		Model: adminList,
		TableName: service.TableName,
		DbConf: "default",
	}

	return page.Filter(params)

}

var admin *AdminService

func Admin() *AdminService {
	if admin == nil {
		admin = &AdminService{TableName: "admin"}
	}
	return admin
}
