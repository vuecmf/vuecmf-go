// Package helper
//+----------------------------------------------------------------------
// | Copyright (c) 2022 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: vuecmf <tulihua2004@126.com>
// +----------------------------------------------------------------------
package helper

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strings"
)

// ListParams 列表参数
type listParams struct {
	Keywords   string                 `json:"keywords" form:"keywords"`       //搜索关键字
	OrderField string                 `json:"order_field" form:"order_field"` //列表排序字段
	OrderSort  string                 `json:"order_sort" form:"order_sort"`   //字段排序方式 （desc 倒序, asc升序）
	Page       int                    `json:"page" form:"page"`               //列表当前页码
	PageSize   int                    `json:"page_size" form:"page_size"`     //列表每页显示条数
	Action     string                 `json:"action" form:"action"`           //请求动作
	Filter     map[string]interface{} `json:"filter" form:"filter"`           //精确多字段过滤查询
}

// DataListParams 列表参数 用data包裹一下
type DataListParams struct {
	Data *listParams `json:"data" form:"data"`
}

// page 分页结构体
type page struct {
	tableName    string   //模型对应表名
	db           *gorm.DB //数据库连接实例
	ns           schema.Namer
	filterFields []string //需要模糊查询的字段
}

// result 存放分页列表返回结果
type result struct {
	Data  interface{} `json:"data"`
	Total int64       `json:"total"`
}

// Filter 列表过滤器
//	参数：
//		model 模型实例
//		params POST请求传递的参数
func (p *page) Filter(model interface{}, params *DataListParams) (*result, error) {
	if params.Data == nil {
		return nil, errors.New("请求参数data不能为空")
	}
	data := params.Data

	if data.PageSize == 0 {
		data.PageSize = 20
	}

	if data.Page == 0 {
		data.Page = 1
	}

	if data.OrderField == "" {
		data.OrderField = "id"
	}

	offset := (data.Page - 1) * data.PageSize
	query := p.db.Table(p.ns.TableName(p.tableName)).Offset(offset).Limit(data.PageSize)
	totalQuery := p.db.Table(p.ns.TableName(p.tableName))

	if data.Keywords != "" {
		kw := data.Keywords + "%"
		for k, field := range p.filterFields {
			field = strings.Trim(field, " ")
			if k == 0 {
				query = query.Where(field+" LIKE ?", kw)
				totalQuery = totalQuery.Where(field+" LIKE ?", kw)
			} else {
				query = query.Or(field+" LIKE ?", kw)
				totalQuery = totalQuery.Or(field+" LIKE ?", kw)
			}
		}

	} else if len(data.Filter) > 0 {
		//过滤掉空值
		for k, v := range data.Filter {
			switch val := v.(type) {
			case string:
				if val == "" {
					delete(data.Filter, k)
				}
			case []string:
				if len(val) == 0 {
					delete(data.Filter, k)
				}
			case []interface{}:
				if len(val) == 0 {
					delete(data.Filter, k)
				}
			case []int:
				if len(val) == 0 {
					delete(data.Filter, k)
				}
			}
		}

		query = query.Where(data.Filter)
		totalQuery = totalQuery.Where(data.Filter)
	}

	query.Order(data.OrderField + " " + data.OrderSort).Find(model)

	var total int64
	totalQuery.Count(&total)

	res := &result{
		Data:  model,
		Total: total,
	}

	return res, nil
}

var pInstances = make(map[string]*page)

// Page 获取列表分页实例
//	 参数：
//		tableName	模型对应的表名
//		db			gorm的DB实例指针
//		ns			gorm数据库相关信息接口
func Page(tableName string, filterFields []string, db *gorm.DB, ns schema.Namer) *page {
	p, ok := pInstances[tableName]
	if ok == false {
		pInstances[tableName] = &page{
			tableName:    tableName,
			db:           db,
			ns:           ns,
			filterFields: filterFields,
		}
		return pInstances[tableName]
	} else {
		return p
	}
}
