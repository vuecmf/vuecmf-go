//+----------------------------------------------------------------------
// | Copyright (c) 2023 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: vuecmf <tulihua2004@126.com>
// +----------------------------------------------------------------------

package service

import (
	"github.com/vuecmf/vuecmf-go/app/vuecmf/model"
)

// appConfigService appConfig服务结构
type appConfigService struct {
	*BaseService
}

var appConfig *appConfigService

// AppConfig 获取appConfig服务实例
func AppConfig() *appConfigService {
	if appConfig == nil {
		appConfig = &appConfigService{}
	}
	return appConfig
}

//GetAppList 获取扩展应用列表
func (s *appConfigService) GetAppList() []string {
	var appList []string
	Db.Table(NS.TableName("app_config")).Select("app_name").
		Where("type = 20").
		Where("status = 10").Find(&appList)
	return appList
}

//GetAuthAppList 获取需要授权的应用列表
func (s *appConfigService) GetAuthAppList() []string {
	var appList []string
	Db.Table(NS.TableName("app_config")).Select("app_name").
		Where("auth_enable = 10").
		Where("status = 10").Find(&appList)
	return appList
}

//GetFullAppList 获取所有可用的应用列表
func (s *appConfigService) GetFullAppList() map[string]*model.AppConfig {
	var ac []*model.AppConfig
	var res = make(map[string]*model.AppConfig)
	Db.Table(NS.TableName("app_config")).Select("app_name, login_enable, auth_enable, exclusion_url").
		Where("status = 10").Find(&ac)

	for _, v := range ac {
		res[v.AppName] = v
	}

	return res
}

//GetAppModelCount 获取指定应用的模型数量
//参数：
//		appId 应用ID
func (s *appConfigService) GetAppModelCount(appId uint) int64 {
	var res int64
	Db.Table(NS.TableName("model_config")).Where("app_id = ?", appId).Count(&res)
	return res
}

//GetAppListByModelId 根据模型ID获取应用列表
//参数：
//		modelId 模型ID
func (s *appConfigService) GetAppListByModelId(modelId uint) []string {
	var res []string
	Db.Table(NS.TableName("app_config")+" AC").Select("app_name").
		Joins("left join "+NS.TableName("model_config")+" MC ON MC.app_id = AC.id").
		Where("MC.id = ?", modelId).
		Where("AC.status = 10").
		Where("MC.status = 10").
		Group("app_name").Find(&res)
	return res
}

//GetAppListByTableName 根据表名获取应用列表
//参数：
//		tableName 表名
func (s *appConfigService) GetAppListByTableName(tableName string) []string {
	var res []string
	Db.Table(NS.TableName("app_config")+" AC").Select("app_name").
		Joins("left join "+NS.TableName("model_config")+" MC ON MC.app_id = AC.id").
		Where("MC.table_name = ?", tableName).
		Where("AC.status = 10").
		Where("MC.status = 10").
		Group("app_name").Find(&res)
	return res
}

//GetAppNameById 根据应用ID获取对应的应用名称
//参数：
//		appId 应用ID
func (s *appConfigService) GetAppNameById(appId uint) string {
	var res string
	Db.Table(NS.TableName("app_config")).Select("app_name").
		Where("id = ?", appId).
		Where("status = 10").
		Find(&res)
	return res
}
