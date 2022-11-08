// Package service
//+----------------------------------------------------------------------
// | Copyright (c) 2022 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: vuecmf <tulihua2004@126.com>
// +----------------------------------------------------------------------
package service

import (
	"github.com/vuecmf/vuecmf-go/app/vuecmf/model"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

// modelFieldService modelField服务结构
type modelFieldService struct {
	*BaseService
}

var modelField *modelFieldService

// ModelField 获取modelField服务实例
func ModelField() *modelFieldService {
	if modelField == nil {
		modelField = &modelFieldService{}
	}
	return modelField
}

// Create 创建单条或多条数据, 成功返回影响行数
func (ser *modelFieldService) Create(data *model.ModelField) (int64, error) {
	err := Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(data).Error; err != nil {
			return err
		}
		return Make().AddField(data, tx)
	})
	if err != nil {
		return 0, err
	}
	return 1, nil
}

// Update 更新数据, 成功返回影响行数
func (ser *modelFieldService) Update(data *model.ModelField) (int64, error) {
	var oldFieldName string
	Db.Table(NS.TableName("model_field")).Select("field_name").
		Where("id = ?", data.Id).Find(&oldFieldName)

	err := Db.Transaction(func(tx *gorm.DB) error {
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
func (ser *modelFieldService) Delete(id uint, model *model.ModelField) (int64, error) {
	err := Db.Transaction(func(tx *gorm.DB) error {
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
func (ser *modelFieldService) DeleteBatch(idList string, model *model.ModelField) (int64, error) {
	idArr := strings.Split(idList, ",")
	err := Db.Transaction(func(tx *gorm.DB) error {
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

// fieldInfo 列表表字段信息
type fieldInfo struct {
	FieldId  int    `json:"field_id"`
	Prop     string `json:"prop"`
	Type     string `json:"type"`
	Label    string `json:"label"`
	Width    int    `json:"width"`
	Length   int    `json:"length"`
	Show     bool   `json:"show"`
	Fixed    bool   `json:"fixed"`
	Filter   bool   `json:"filter"`
	Tooltip  string `json:"tooltip"`
	ModelId  int    `json:"model_id"`
	Sortable bool   `json:"sortable"`
}

// GetFieldInfo 根据模型ID获取对应的字段信息
func (ser *modelFieldService) GetFieldInfo(modelId int) []fieldInfo {
	var list []fieldInfo

	Db.Table(NS.TableName("model_field")).Select(
		"id field_id,"+
			"field_name prop,"+
			"type,"+
			"label,"+
			"column_width width,"+
			"length,"+
			"if(is_show = 10,true, false) `show`,"+
			"if(is_fixed = 10,true, false) fixed,"+
			"if(is_filter = 10,true, false) `filter`,"+
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
func (ser *modelFieldService) getFilterFields(tableName string) []string {
	var filterFields []string
	Db.Table(NS.TableName("model_field")+" MF").Select("field_name").
		Joins("left join "+NS.TableName("model_config")+" MC on MF.model_id = MC.id").
		Where("MF.is_filter = 10").
		Where("MF.type in (?)", []string{"char", "varchar"}).
		Where("MC.table_name = ?", tableName).
		Limit(50).Find(&filterFields)

	return filterFields
}
