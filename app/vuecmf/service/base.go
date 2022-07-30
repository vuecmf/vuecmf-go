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
	"github.com/vuecmf/vuecmf-go/app"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/helper"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var db *gorm.DB
var ns schema.Namer

type base struct {
}

func init() {
	db = app.Db("default")
	ns = db.NamingStrategy
}

type fullModelFields struct {
	FieldInfo []fieldInfo `json:"field_info"`
	FormInfo []formInfo `json:"form_info"`
	FieldOption map[int]map[string]string `json:"field_option"`
	RelationInfo *modelRelationInfo `json:"relation_info"`
}


// commonList 公共列表 服务方法
func (b *base) commonList(model interface{}, tableName string, params *helper.DataListParams) interface{}{
	if params.Data.Action == "getField" {
		modelId := ModelConfig().GetModelId(tableName)
		fieldInfo := ModelField().GetFieldInfo(modelId) //模型的字段信息
		formInfo := ModelForm().GetFormInfo(modelId)  //模型的表单信息
		fieldOption := FieldOption().GetFieldOptions(modelId) //模型的关联信息
		relationInfo := ModelRelation().GetRelationInfo(modelId, params.Data.Filter)


		return &fullModelFields{
			FieldInfo: fieldInfo,
			FormInfo: formInfo,
			FieldOption: fieldOption,
			RelationInfo: relationInfo,
		}
	}else{
		return helper.Page(tableName, db, ns).Filter(model, params)
	}
}
