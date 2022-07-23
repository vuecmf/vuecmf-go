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
	*base
}

// GetModelId 根据表名获取模型ID
func (s *modelConfigService) GetModelId(tableName string) int {
	var modelId int
	db.Model(ns.TableName("model_config")).Select("id").
		Where("table_name = ?", tableName).Limit(1).Find(&modelId)
	return modelId
}

var modelConfig *modelConfigService

// ModelConfig 获取modelConfig服务实例
func ModelConfig() *modelConfigService {
	if modelConfig == nil {
		modelConfig = &modelConfigService{}
	}
	return modelConfig
}
