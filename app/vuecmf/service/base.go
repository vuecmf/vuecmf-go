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
	"github.com/vuecmf/vuecmf-go/app/vuecmf/model"
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
	FieldInfo    []fieldInfo               `json:"field_info"`
	FormInfo     []formInfo                `json:"form_info"`
	FieldOption  map[int]map[string]string `json:"field_option"`
	RelationInfo *modelRelationInfo        `json:"relation_info"`
	FormRules    interface{}               `json:"form_rules"`
	ModelId      int                       `json:"model_id"`
}

// commonList 公共列表 服务方法
func (b *base) commonList(modelData interface{}, tableName string, params *helper.DataListParams) interface{} {
	modelConf := ModelConfig().GetModelConfig(tableName)

	if params.Data.Action == "getField" {
		modelId := modelConf.ModelId
		fieldInfo := ModelField().GetFieldInfo(modelId)       //模型的字段信息
		formInfo := ModelForm().GetFormInfo(modelId)          //模型的表单信息
		fieldOption := FieldOption().GetFieldOptions(modelId) //模型的关联信息
		relationInfo := ModelRelation().GetRelationInfo(modelId, params.Data.Filter)
		formRulesInfo := ModelFormRules().GetRuleListForForm(modelId)

		//目录树列表中 父级字段处理
		if modelConf.IsTree {
			if modelConf.LabelFieldName == "" {
				panic("模型还没有设置标题字段")
			}
			orderField := "sort_num"
			if tableName == "roles" {
				orderField = ""
			}

			var pidFieldId int
			db.Table(ns.TableName("model_field")).
				Select("id").
				Where("field_name = 'pid'").
				Where("model_id = ?", modelId).
				Limit(1).Find(&pidFieldId)

			tree := map[string]string{}
			helper.FormatTree(tree, db, ns.TableName(tableName), "id", 0, modelConf.LabelFieldName, "pid", orderField, 1)
			fieldOption[pidFieldId] = tree
		}

		return &fullModelFields{
			FieldInfo:    fieldInfo,
			FormInfo:     formInfo,
			FieldOption:  fieldOption,
			RelationInfo: relationInfo,
			FormRules:    formRulesInfo,
			ModelId:      modelId,
		}
	} else if modelConf.IsTree == true {
		//列表数据（目录树形式）
		orderField := "sort_num"
		if tableName == "roles" {
			orderField = ""
		}

		//先查询出所有数据,  // https://blog.csdn.net/LW1314QS/article/details/124517399
		var dataForTree []model.Menu


		//然后将数据格式化目录树


		var res = make(map[string]interface{})
		var tree []model.MenuTree



		res["data"] = helper.TreeList(tree, db, ns.TableName(tableName), 0, params.Data.Keywords, "pid", modelConf.LabelFieldName, orderField)
		return res
	} else {
		return helper.Page(tableName, db, ns).Filter(modelData, params)
	}
}
