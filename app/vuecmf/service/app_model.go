package service

import (
	"github.com/vuecmf/vuecmf-go/app/vuecmf/model"
	"strings"
)

// appModelService appModel服务结构
type appModelService struct {
	*BaseService
}

var appModel *appModelService

// AppModel 获取appModel服务实例
func AppModel() *appModelService {
	if appModel == nil {
		appModel = &appModelService{}
	}
	return appModel
}

// Create 创建单条或多条数据, 成功返回影响行数
func (s *appModelService) Create(data *model.AppModel) (int64, error) {
	res := Db.Create(data)
	if err := Make().MakeAppModel(data.AppId, data.ModelId); err != nil {
		return 0, err
	}
	return res.RowsAffected, res.Error
}

// Update 更新数据, 成功返回影响行数
func (s *appModelService) Update(data *model.AppModel) (int64, error) {
	type Old struct {
		AppId   uint
		ModelId uint
	}
	var oldRes Old
	Db.Table(NS.TableName("app_model")).Select("app_id, model_id").
		Where("id = ?", data.Id).Find(&oldRes)

	res := Db.Updates(data)

	if oldRes.AppId != data.AppId || oldRes.ModelId != data.ModelId {
		//清除原有的相关文件
		if err := Make().RemoveAppModel(oldRes.AppId, oldRes.ModelId); err != nil {
			return 0, err
		}
		//生成新的相关文件
		if err := Make().MakeAppModel(data.AppId, data.ModelId); err != nil {
			return 0, err
		}
	}

	return res.RowsAffected, res.Error
}

// Delete 根据ID删除数据
func (s *appModelService) Delete(id uint, model *model.AppModel) (int64, error) {
	type Old struct {
		AppId   uint
		ModelId uint
	}
	var oldRes Old
	Db.Table(NS.TableName("app_model")).Select("app_id, model_id").
		Where("id = ?", id).Find(&oldRes)

	if err := Make().RemoveAppModel(oldRes.AppId, oldRes.ModelId); err != nil {
		return 0, err
	}

	res := Db.Delete(model, id)
	return res.RowsAffected, res.Error
}

// DeleteBatch 根据ID删除数据， 多个用英文逗号分隔
func (s *appModelService) DeleteBatch(idList string, model *model.AppModel) (int64, error) {
	idArr := strings.Split(idList, ",")
	type Old struct {
		AppId   uint
		ModelId uint
	}

	for _, id := range idArr {
		var oldRes Old
		Db.Table(NS.TableName("app_model")).Select("app_id, model_id").
			Where("id = ?", id).Find(&oldRes)

		if err := Make().RemoveAppModel(oldRes.AppId, oldRes.ModelId); err != nil {
			return 0, err
		}
	}

	res := Db.Delete(model, idArr)
	return res.RowsAffected, res.Error
}

// GetAppList 根据表名获取对应应用列表
func (s *appModelService) GetAppList(tableName string) []string {
	var appList []string
	Db.Table(NS.TableName("app_model")+" AM").Select("A.app_name").
		Joins("left join "+NS.TableName("app_config")+" A on AM.app_id = A.id").
		Joins("left join "+NS.TableName("model_config")+" M on AM.model_id = M.id").
		Where("M.table_name = ?", tableName).
		Where("A.app_name != 'vuecmf'").
		Where("AM.status = 10").
		Where("A.status = 10").
		Where("M.status = 10").
		Group("A.app_name").
		Find(&appList)
	return appList
}

// GetAppListByModelId 根据模型ID获取对应应用列表
func (s *appModelService) GetAppListByModelId(modelId uint) []string {
	var appList []string
	Db.Table(NS.TableName("app_model")+" AM").Select("A.app_name").
		Joins("left join "+NS.TableName("app_config")+" A on AM.app_id = A.id").
		Where("AM.model_id = ?", modelId).
		Where("A.app_name != 'vuecmf'").
		Where("AM.status = 10").
		Where("A.status = 10").
		Group("A.app_name").
		Find(&appList)
	return appList
}
