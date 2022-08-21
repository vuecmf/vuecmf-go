// Package service
//+----------------------------------------------------------------------
// | Copyright (c) 2022 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: vuecmf <tulihua2004@126.com>
// +----------------------------------------------------------------------
package service

// modelConfigService modelConfig服务结构
type modelConfigService struct {
	*baseService
}

var modelConfig *modelConfigService

// ModelConfig 获取modelConfig服务实例
func ModelConfig() *modelConfigService {
	if modelConfig == nil {
		modelConfig = &modelConfigService{}
	}
	return modelConfig
}

// GetModelId 根据表名获取模型ID
func (s *modelConfigService) GetModelId(tableName string) int {
	var modelId int
	db.Table(ns.TableName("model_config")).Select("id").
		Where("table_name = ?", tableName).
		Where("status = 10").
		Limit(1).Find(&modelId)
	return modelId
}

// GetModelTableName 根据模型ID获取模型对应表名
func (s *modelConfigService) GetModelTableName(modelId int) string {
	var tableName string
	db.Table(ns.TableName("model_config")).Select("table_name").
		Where("id = ?", modelId).
		Limit(1).Find(&tableName)
	return tableName
}

// IsTree 根据模型ID判断是否为目录树
func (s *modelConfigService) IsTree(modelId int) bool {
	var isTree int
	db.Table(ns.TableName("model_config")).Select("is_tree").
		Where("id = ?", modelId).
		Limit(1).Find(&isTree)
	return isTree == 10
}

type modelConf struct {
	TableName      string
	IsTree         bool
	ModelId        int
	LabelFieldName string
}

// GetModelConfig 根据模型表名获取模型的配置信息
func (s *modelConfigService) GetModelConfig(tableName string) modelConf {
	var modelConfig modelConf
	db.Table(ns.TableName("model_config")).
		Select("table_name, if(is_tree = 10, true, false) is_tree, id model_id, '' label_field_name").
		Where("status = 10").
		Where("table_name = ?", tableName).
		Limit(1).
		Find(&modelConfig)

	var labelFieldName string
	db.Table(ns.TableName("model_field")).
		Select("field_name").
		Where("status = 10").
		Where("model_id = ?", modelConfig.ModelId).
		Where("is_label = 10").
		Find(&labelFieldName)

	modelConfig.LabelFieldName = labelFieldName

	return modelConfig
}
