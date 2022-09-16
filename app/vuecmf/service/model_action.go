// Package service
//+----------------------------------------------------------------------
// | Copyright (c) 2022 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: vuecmf <tulihua2004@126.com>
// +----------------------------------------------------------------------
package service

// modelActionService modelAction服务结构
type modelActionService struct {
	*baseService
}

// GetAllApiMap 根据动作ID获取 api 路径映射
func (ser *modelActionService) GetAllApiMap(apiIdList []string) interface{} {
	var res = make(map[string]map[string]string)

	type actionList struct {
		TableName string
		ActionType string
		ApiPath string
	}

	var ac []actionList

	db.Table(ns.TableName("model_action") + " ma").Select("mc.table_name, ma.action_type, ma.api_path").
		Joins("inner join " + ns.TableName("model_config") + " mc on ma.model_id = mc.id").
		Where("ma.id in ?", apiIdList).
		Find(&ac)

	for _, val := range ac {
		if res[val.TableName] == nil {
			res[val.TableName] = map[string]string{}
		}
		res[val.TableName][val.ActionType] = val.ApiPath
	}

	return res
}

// GetModelIdListById 根据动作ID获取对应模型ID列表
func (ser *modelActionService) GetModelIdListById(apiIdList []string) []string {
	var res []string
	db.Table(ns.TableName("model_action")).Select("model_id").
		Where("id in ?", apiIdList).
		Group("model_id").
		Find(&res)
	res = append(res, "0")
	return res
}



var modelAction *modelActionService

// ModelAction 获取modelAction服务实例
func ModelAction() *modelActionService {
	if modelAction == nil {
		modelAction = &modelActionService{}
	}
	return modelAction
}