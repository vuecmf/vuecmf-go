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

// page 列表结构
type page struct {
	tableName string   //模型对应表名
	db        *gorm.DB //数据库连接实例
	ns        schema.Namer
}

// Filter 列表过滤器
//	参数：
//		model 模型实例
//		params POST请求传递的参数
func (p *page) Filter(model interface{}, params *DataListParams) interface{} {
	if params.Data == nil {
		panic("请求参数data不能为空")
	}
	data := params.Data

	if data.Action == "getField" {

	}

	if data.PageSize == 0 {
		data.PageSize = 20
	}

	if data.Page == 0 {
		data.Page = 1
	}

	if data.OrderField == "" {
		data.OrderField = "id"
	}

	//查询出该表中的可过滤的字段
	var filterFields []string
	p.db.Table(p.ns.TableName("model_field")+" MF").Select("field_name").
		Joins("left join "+p.ns.TableName("model_config")+" MC on MF.model_id = MC.id").
		Where("MF.is_filter = 10").
		Where("MF.type in (?)", []string{"char", "varchar"}).
		Where("MC.table_name = ?", p.tableName).
		Limit(50).Find(&filterFields)

	offset := (data.Page - 1) * data.PageSize
	query := p.db.Table(p.ns.TableName(p.tableName)).Offset(offset).Limit(data.PageSize)

	if data.Keywords != "" {
		kw := data.Keywords + "%"
		for k, field := range filterFields {
			field = strings.Trim(field, " ")
			if k == 0 {
				query = query.Where(field+" LIKE ?", kw)
			} else {
				query = query.Or(field+" LIKE ?", kw)
			}
		}

	} else if len(data.Filter) > 0 {
		query = query.Where(data.Filter)
	}

	query.Order(data.OrderField + " " + data.OrderSort).Find(&model)
	return model
}

var pInstances = make(map[string]*page)

// Page 获取列表分页实例
//	 参数：
//		tableName	模型对应的表名
//		db			gorm的DB实例指针
//		ns			gorm数据库相关信息接口
func Page(tableName string, db *gorm.DB, ns schema.Namer) *page {
	p, ok := pInstances[tableName]
	if ok == false {
		pInstances[tableName] = &page{
			tableName: tableName,
			db:        db,
			ns:        ns,
		}
		return pInstances[tableName]
	} else {
		return p
	}
}
