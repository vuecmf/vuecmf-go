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
	"strings"
)

// modelIndexService modelIndex服务结构
type modelIndexService struct {
	*BaseService
}

var modelIndex *modelIndexService

// ModelIndex 获取modelIndex服务实例
func ModelIndex() *modelIndexService {
	if modelIndex == nil {
		modelIndex = &modelIndexService{}
	}
	return modelIndex
}

// Create 创建单条或多条数据, 成功返回影响行数
func (s *modelIndexService) Create(data *model.ModelIndex) (int64, error) {
	res := Db.Create(data)
	if err := Make().AddIndex(data); err != nil {
		return 0, err
	}
	return res.RowsAffected, res.Error
}

// Update 更新数据, 成功返回影响行数
func (s *modelIndexService) Update(data *model.ModelIndex) (int64, error) {
	//删除原索引
	var oldModelFieldId string
	Db.Table(NS.TableName("model_index")).Select("model_field_id").
		Where("id = ?", data.Id).Find(&oldModelFieldId)
	if err := Make().DelIndex(oldModelFieldId, data.ModelId); err != nil {
		return 0, err
	}

	res := Db.Updates(data)

	//添加新索引
	if err := Make().AddIndex(data); err != nil {
		return 0, err
	}

	return res.RowsAffected, res.Error
}

// Delete 根据ID删除数据
func (s *modelIndexService) Delete(id uint, model interface{}) (int64, error) {
	type Res struct {
		ModelFieldId string
		ModelId      uint
	}
	var rs Res
	Db.Table(NS.TableName("model_index")).Select("model_field_id, model_id").
		Where("id = ?", id).Find(&rs)

	if err := Make().DelIndex(rs.ModelFieldId, rs.ModelId); err != nil {
		return 0, err
	}

	res := Db.Delete(model, id)
	return res.RowsAffected, res.Error
}

// DeleteBatch 根据ID删除数据， 多个用英文逗号分隔
func (s *modelIndexService) DeleteBatch(idList string, model interface{}) (int64, error) {
	idArr := strings.Split(idList, ",")

	type Res struct {
		ModelFieldId string
		ModelId      uint
	}
	for _, id := range idArr {
		var rs Res
		Db.Table(NS.TableName("model_index")).Select("model_field_id, model_id").
			Where("id = ?", id).Find(&rs)

		if err := Make().DelIndex(rs.ModelFieldId, rs.ModelId); err != nil {
			return 0, err
		}
	}

	res := Db.Delete(model, idArr)
	return res.RowsAffected, res.Error
}
