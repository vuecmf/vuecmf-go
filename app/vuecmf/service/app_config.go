package service

import "github.com/vuecmf/vuecmf-go/app/vuecmf/model"

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
func (m *appConfigService) GetAppList() []string {
	var appList []string
	Db.Table(NS.TableName("app_config")).Select("app_name").
		Where("type = 20").
		Where("status = 10").Find(&appList)
	return appList
}

//GetAuthAppList 获取需要授权的应用列表
func (m *appConfigService) GetAuthAppList() []string {
	var appList []string
	Db.Table(NS.TableName("app_config")).Select("app_name").
		Where("auth_enable = 10").
		Where("status = 10").Find(&appList)
	return appList
}

//GetFullAppList 获取所有可用的应用列表
func (m *appConfigService) GetFullAppList() []*model.AppConfig {
	var res []*model.AppConfig
	Db.Table(NS.TableName("app_config")).Select("app_name, login_enable, auth_enable, exclusion_url").
		Where("status = 10").Find(&res)
	return res
}

type modelList struct {
	Key      uint   `json:"key"`
	Label    string `json:"label"`
	Disabled bool   `json:"disabled"`
}

// GetAllModels 获取所有模型列表
func (m *appConfigService) GetAllModels() interface{} {
	var result []modelList
	Db.Table(NS.TableName("model_config")).Select("id `key`, label, false disabled").
		Where("status = 10").
		Find(&result)
	return result
}

//GetModels 根据应用名称获取模型ID
func (m *appConfigService) GetModels(appName string) []int {
	var res []int
	Db.Table(NS.TableName("model_config")).Select("id").
		Where("role_name in ?", roleNameList).
		Where("status = 10").
		Find(&res)
	return res
}
