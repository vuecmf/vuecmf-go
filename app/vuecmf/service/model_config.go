package service

import "github.com/vuecmf/vuecmf-go/app"

// ModelConfigService model_config服务结构
type ModelConfigService struct {
	*base
}

// GetModelId 根据表名获取模型ID
func (service *ModelConfigService) GetModelId(tableName string) int {
	var modelId int
	app.Db("default").Table("model_config").Db.Select("id").
		Where("table_name = ?", tableName).Limit(1).Find(&modelId)
	return modelId
}

var modelConfigService *ModelConfigService

func ModelConfig() *ModelConfigService {
	if modelConfigService == nil {
		modelConfigService = &ModelConfigService{}
	}
	return modelConfigService
}