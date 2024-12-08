//+----------------------------------------------------------------------
// | Copyright (c) 2024 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: tulihua2004@126.com
// +----------------------------------------------------------------------

package service

import (
	"github.com/vuecmf/vuecmf-go/v3/app"
	"github.com/vuecmf/vuecmf-go/v3/app/vuecmf/model"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"sync"
)

// ModelFieldService modelField服务结构
type ModelFieldService struct {
	*BaseService
}

var modelFieldOnce sync.Once
var modelField *ModelFieldService

// ModelField 获取modelField服务实例
func ModelField() *ModelFieldService {
	modelFieldOnce.Do(func() {
		modelField = &ModelFieldService{
			BaseService: &BaseService{
				"model_field",
				&model.ModelField{},
				&[]model.ModelField{},
				[]string{"field_name", "label", "type", "note", "default_value"},
			},
		}
	})
	return modelField
}

// Create 创建单条或多条数据, 成功返回插入的ID
//
//	参数：
//		data 需保存的数据
func (svc *ModelFieldService) Create(data *model.ModelField) (int64, error) {
	err := app.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(data).Error; err != nil {
			return err
		}
		return Make().AddField(data, tx)
	})
	if err != nil {
		return 0, err
	}
	return int64(data.Id), nil
}

// Update 更新数据, 成功返回影响行数
//
//	参数：
//		data 需更新的数据
func (svc *ModelFieldService) Update(data *model.ModelField) (int64, error) {
	var oldFieldName string
	DbTable("model_field").Select("field_name").
		Where("id = ?", data.Id).Find(&oldFieldName)

	err := app.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Updates(data).Error; err != nil {
			return err
		}
		if err := Make().RenameField(data, oldFieldName, tx); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return 0, err
	}
	return 1, nil
}

// Delete 根据ID删除数据
//
//	参数：
//		id 需删除的ID
//		model 模型实例
func (svc *ModelFieldService) Delete(id uint, model *model.ModelField) (int64, error) {
	err := app.Db.Transaction(func(tx *gorm.DB) error {
		tx.Model(model).Where("id = ?", id).Find(&model)
		if err := tx.Delete(model, id).Error; err != nil {
			return err
		}
		return Make().DelField(model, tx)
	})
	if err != nil {
		return 0, err
	}
	return 1, nil
}

// DeleteBatch 根据ID删除数据， 多个用英文逗号分隔
//
//	参数：
//		idList 需删除的ID列表
//		model 模型实例
func (svc *ModelFieldService) DeleteBatch(idList string, model *model.ModelField) (int64, error) {
	idArr := strings.Split(idList, ",")
	err := app.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(model, idArr).Error; err != nil {
			return err
		}
		for _, id := range idArr {
			mid, _ := strconv.Atoi(id)
			model.Id = uint(mid)
			if err := Make().DelField(model, tx); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return 0, err
	}
	return int64(len(idArr)), nil

}

// FieldInfo 列表表字段信息
type FieldInfo struct {
	FieldId  int    `json:"field_id"`
	Prop     string `json:"prop"`
	Type     string `json:"type"`
	Label    string `json:"label"`
	Width    int    `json:"width"`
	Length   int    `json:"length"`
	Show     bool   `json:"show"`
	Fixed    bool   `json:"fixed"`
	Filter   bool   `json:"filter"`
	Code     bool   `json:"code"`
	Tooltip  string `json:"tooltip"`
	ModelId  int    `json:"model_id"`
	Sortable bool   `json:"sortable"`
}

// GetFieldInfo 根据模型ID获取对应的字段信息
//
//	参数：
//		modelId 模型ID
func (svc *ModelFieldService) GetFieldInfo(modelId int) []FieldInfo {
	var list []FieldInfo

	DbTable("model_field").Select(
		"id field_id,"+
			"field_name prop,"+
			"type,"+
			"label,"+
			"column_width width,"+
			"length,"+
			"if(is_show = 10,true, false) `show`,"+
			"if(is_fixed = 10,true, false) fixed,"+
			"if(is_filter = 10,true, false) `filter`,"+
			"if(is_code = 10,true, false) `code`,"+
			"note tooltip,"+
			"model_id,"+
			"true sortable").
		Where("model_id = ?", modelId).
		Where("status = 10").
		Order("sort_num").
		Find(&list)

	return list
}

// getFilterFields 根据表名获取该表需要模糊查询的字段
//
//	参数：
//		tableName 表名
func (svc *ModelFieldService) getFilterFields(tableName string) []string {
	var filterFields []string
	DbTable("model_field", "MF").Select("field_name").
		Joins("left join "+TableName("model_config")+" MC on MF.model_id = MC.id").
		Where("MF.is_filter = 10").
		Where("MF.type in (?)", []string{"char", "varchar"}).
		Where("MC.table_name = ?", tableName).
		Limit(50).Find(&filterFields)

	return filterFields
}

// GetFieldId 根据字段名获取对应字段ID
//
//	参数：
//		fieldName 字段名
//		modelId 模型ID
func (svc *ModelFieldService) GetFieldId(fieldName string, modelId uint) uint {
	var id uint
	DbTable("model_field").Select("id").
		Where("field_name = ?", fieldName).
		Where("model_id = ?", modelId).
		Where("status = 10").
		Find(&id)
	return id
}
