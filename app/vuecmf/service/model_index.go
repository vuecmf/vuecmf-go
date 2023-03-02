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
	"gorm.io/gorm"
	"strconv"
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
// 参数：
//		data 需保存的数据
func (s *modelIndexService) Create(data *model.ModelIndex) (int64, error) {
	err := Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(data).Error; err != nil {
			return err
		}
		return Make().AddIndex(data, tx)
	})

	if err != nil {
		return 0, err
	}
	return 1, nil
}

// Update 更新数据, 成功返回影响行数
// 参数：
//		data 需更新的数据
func (s *modelIndexService) Update(data *model.ModelIndex) (int64, error) {
	err := Db.Transaction(func(tx *gorm.DB) error {
		//删除原索引
		if err := Make().DelIndex(data.Id, tx); err != nil {
			return err
		}
		//更新索引数据
		if err := tx.Updates(data).Error; err != nil {
			return err
		}
		//添加新索引
		return Make().AddIndex(data, tx)
	})

	if err != nil {
		return 0, err
	}
	return 1, nil
}

// Delete 根据ID删除数据
// 参数：
//		id 需删除的ID
// 		model 模型实例
func (s *modelIndexService) Delete(id uint, model interface{}) (int64, error) {
	err := Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(model, id).Error; err != nil {
			return err
		}
		return Make().DelIndex(id, tx)
	})
	if err != nil {
		return 0, err
	}
	return 1, nil
}

// DeleteBatch 根据ID删除数据， 多个用英文逗号分隔
// 参数：
//		idList 需删除的ID列表
// 		model 模型实例
func (s *modelIndexService) DeleteBatch(idList string, model interface{}) (int64, error) {
	idArr := strings.Split(idList, ",")
	err := Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(model, idArr).Error; err != nil {
			return err
		}
		for _, id := range idArr {
			delId, _ := strconv.Atoi(id)
			if err := Make().DelIndex(uint(delId), tx); err != nil {
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
