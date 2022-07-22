package helper

import (
	"fmt"
	"strings"
	"github.com/vuecmf/vuecmf-go/app"
)

// ListParams 列表参数
type listParams struct {
	Keywords 	string `json:"keywords" form:"keywords"`		//搜索关键字
	OrderField 	string `json:"order_field" form:"order_field"`	//列表排序字段
	OrderSort 	string `json:"order_sort" form:"order_sort"`	//字段排序方式 （desc 倒序, asc升序）
	Page     	int `json:"page" form:"page"`					//列表当前页码
	PageSize 	int `json:"page_size" form:"page_size"`			//列表每页显示条数
	Filter 		map[string]interface{}							//精确多字段过滤查询
}

// DataListParams 列表参数 用data包裹一下
type DataListParams struct {
	Data *listParams `json:"data" form:"data"`
}

// Page 列表结构
type Page struct {
	Model interface{}  //数据模型
	TableName string   //模型对应表名
	DbConf string
}

// Filter 列表过滤器
func(p *Page) Filter(params *DataListParams) interface{} {
	if params.Data == nil {
		panic("请求参数data不能为空")
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

	db := app.Db(p.DbConf)


	//查询出该表中的可过滤的字段
	var filterFields []string
	db.Table("model_field MF").Db.Select("field_name").
		Joins("left join " + db.Conf.Prefix + "model_config MC on MF.model_id = MC.id").
		Where("MF.is_filter = 10").
		Where("MF.type in (?)", []string{"char","varchar"}).
		Where("MC.table_name = ?", p.TableName).
		Limit(50).Find(&filterFields)

	fmt.Println("fields=", filterFields)


	offset := (data.Page - 1) * data.PageSize

	query := db.Db.Model(&p.Model).Offset(offset).Limit(data.PageSize)

	if data.Keywords != "" {
		kw := "%" + data.Keywords + "%"
		for k, field := range filterFields {
			field = strings.Trim(field," ")
			if k == 0 {
				query = query.Where(field + " LIKE ?", kw)
			}else{
				query = query.Or(field + " LIKE ?", kw)
			}
		}

	}else if len(data.Filter) > 0 {
		query = query.Where(data.Filter)
	}

	fmt.Println(data.OrderField + " " + data.OrderSort)

	query.Order(data.OrderField + " " + data.OrderSort).Find(&p.Model)

	return p.Model
}




