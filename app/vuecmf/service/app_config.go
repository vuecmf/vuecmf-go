package service

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
