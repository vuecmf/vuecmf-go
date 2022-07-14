package controller

import "github.com/vuecmf/vuecmf-go/app"

type base struct {
	Request app.Request
	Response app.Response
}

// ListParams 列表参数
type ListParams struct {
	Keywords 	string `json:"keywords" form:"keywords"`
	OrderField 	string `json:"order_field" form:"order_field"`
	OrderSort 	string `json:"order_sort" form:"order_sort"`
	Page     	int `json:"page" form:"page"`
	PageSize 	int `json:"page_size" form:"page_size"`
	Filter 		map[string]interface{}
}

// DataListParams 列表参数 用data包裹一下
type DataListParams struct {
	Data *ListParams `json:"data" form:"data"`
}





