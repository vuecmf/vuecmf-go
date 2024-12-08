//+----------------------------------------------------------------------
// | Copyright (c) 2024 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: tulihua2004@126.com
// +----------------------------------------------------------------------

package service

import (
	"github.com/vuecmf/vuecmf-go/v3/app/vuecmf/model"
	"sync"
)

// AppConfigService appConfig服务结构
type AppConfigService struct {
	*BaseService
}

var appConfigOnce sync.Once
var appConfig *AppConfigService

// AppConfig 获取appConfig服务实例
func AppConfig() *AppConfigService {
	appConfigOnce.Do(func() {
		appConfig = &AppConfigService{
			BaseService: &BaseService{
				"app_config",
				&model.AppConfig{},
				&[]model.AppConfig{},
				[]string{"app_name", "exclusion_url"},
			},
		}
	})
	return appConfig
}

// GetAppList 获取扩展应用列表
func (svc *AppConfigService) GetAppList() []string {
	var appList []string
	DbTable("app_config").Select("app_name").
		Where("type = 20").
		Where("status = 10").Find(&appList)
	return appList
}

// GetAuthAppList 获取需要授权的应用列表
func (svc *AppConfigService) GetAuthAppList() []string {
	var appList []string
	DbTable("app_config").Select("app_name").
		Where("auth_enable = 10").
		Where("status = 10").Find(&appList)
	return appList
}

// GetFullAppList 获取所有可用的应用列表
func (svc *AppConfigService) GetFullAppList() map[string]*model.AppConfig {
	var ac []*model.AppConfig
	var res = make(map[string]*model.AppConfig)
	DbTable("app_config").Select("app_name, login_enable, auth_enable, exclusion_url").
		Where("status = 10").Find(&ac)

	for _, v := range ac {
		res[v.AppName] = v
	}

	return res
}

// GetAppModelCount 获取指定应用的模型数量
//
//	参数：
//		appId 应用ID
func (svc *AppConfigService) GetAppModelCount(appId uint) int64 {
	var res int64
	DbTable("model_config").Where("app_id = ?", appId).Count(&res)
	return res
}

// GetAppListByModelId 根据模型ID获取应用列表
//
//	参数：
//		modelId 模型ID
func (svc *AppConfigService) GetAppListByModelId(modelId uint) []string {
	var res []string
	DbTable("app_config", "AC").Select("app_name").
		Joins("left join "+TableName("model_config")+" MC ON MC.app_id = AC.id").
		Where("MC.id = ?", modelId).
		Where("AC.status = 10").
		Where("MC.status = 10").
		Group("app_name").Find(&res)
	return res
}

// GetAppListByTableName 根据表名获取应用列表
//
//	参数：
//		tableName 表名
func (svc *AppConfigService) GetAppListByTableName(tableName string) []string {
	var res []string
	DbTable("app_config", "AC").Select("app_name").
		Joins("left join "+TableName("model_config")+" MC ON MC.app_id = AC.id").
		Where("MC.table_name = ?", tableName).
		Where("AC.status = 10").
		Where("MC.status = 10").
		Group("app_name").Find(&res)
	return res
}

// GetAppNameById 根据应用ID获取对应的应用名称
//
//	参数：
//		appId 应用ID
func (svc *AppConfigService) GetAppNameById(appId uint) string {
	var res string
	DbTable("app_config").Select("app_name").
		Where("id = ?", appId).
		Where("status = 10").
		Find(&res)
	return res
}
