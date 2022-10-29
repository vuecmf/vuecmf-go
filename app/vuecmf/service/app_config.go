package service

import (
	"github.com/vuecmf/vuecmf-go/app/vuecmf/model"
	"gorm.io/gorm"
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

type modelList struct {
	Key      uint   `json:"key"`
	Label    string `json:"label"`
	Disabled bool   `json:"disabled"`
}

// GetAllModels 获取所有模型列表
func (s *appConfigService) GetAllModels() interface{} {
	var result []modelList
	Db.Table(NS.TableName("model_config")).Select("id `key`, label, false disabled").
		Where("status = 10").
		Find(&result)
	return result
}

//GetModels 根据应用名称获取模型ID
func (s *appConfigService) GetModels(appId uint) ([]int, error) {
	var res []int
	Db.Table(NS.TableName("app_model")).Select("model_id").
		Where("app_id = ?", appId).
		Where("status = 10").
		Find(&res)
	return res, nil
}

//DelAllModelsForApp 清空应用下所有模型
func (s *appConfigService) DelAllModelsForApp(appId uint) (bool, error) {
	//取出要删除的所有模型ID
	var modelIdList []uint
	Db.Table(NS.TableName("app_model")).Select("model_id").Where("app_id = ?", appId).Find(&modelIdList)
	//清除代码文件
	for _, modelId := range modelIdList {
		if err := Make().RemoveAppModel(appId, modelId); err != nil {
			return false, err
		}
	}
	//删除应用下的模型
	Db.Where("app_id = ?", appId).Delete(&model.AppModel{})
	return true, nil
}

//AddModelsForApp 为应用分配模型
func (s *appConfigService) AddModelsForApp(appId uint, modelIdList []uint) (bool, error) {
	err := Db.Transaction(func(tx *gorm.DB) error {
		//先取应用原有的模型列表
		var oldModelIdList []uint
		Db.Table(NS.TableName("app_model")).Select("model_id").Where("app_id = ?", appId).Find(&oldModelIdList)

		//需要删除的模型
		var delModelIdList []uint
		//需要添加的模型
		var addAppModelList []model.AppModel

		for _, oldModelId := range oldModelIdList {
			flag := false
			for _, modelId := range modelIdList {
				if oldModelId == modelId {
					flag = true
					break
				}
			}
			if flag == false {
				delModelIdList = append(delModelIdList, oldModelId)
			}
		}

		for _, modelId := range modelIdList {
			flag := false
			for _, oldModelId := range oldModelIdList {
				if modelId == oldModelId {
					flag = true
					break
				}
			}
			if flag == false {
				addAppModelList = append(addAppModelList, model.AppModel{
					AppId:   appId,
					ModelId: modelId,
				})
			}
		}

		if len(delModelIdList) > 0 {
			//删除原有模型在当前配置的不存在的模型
			if err := tx.Where("app_id = ?", appId).Where("model_id in ?", delModelIdList).Delete(&model.AppModel{}).Error; err != nil {
				return err
			}

			//清除代码文件
			for _, modelId := range delModelIdList {
				if err := Make().RemoveAppModel(appId, modelId); err != nil {
					return err
				}
			}
		}

		if len(addAppModelList) > 0 {
			//添加当前在原有的模型中不存在的模型
			if err := tx.Create(addAppModelList).Error; err != nil {
				return err
			}

			//生成新的代码文件
			for _, item := range addAppModelList {
				if err := Make().MakeAppModel(item.AppId, item.ModelId); err != nil {
					return err
				}
			}
		}

		return nil
	})

	return true, err
}

//GetAppModelCount 获取指定应用的模型数量
func (s *appConfigService) GetAppModelCount(appId uint) int64 {
	var res int64
	Db.Table(NS.TableName("app_model")).Where("app_id = ?", appId).Count(&res)
	return res
}
