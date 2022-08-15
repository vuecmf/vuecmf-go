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
	"fmt"
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

// fullModelFields 模型的所有字段相关信息（字段、表单、字段选项、字段关联、表单验证规则）
type fullModelFields struct {
	FieldInfo    []fieldInfo               `json:"field_info"`
	FormInfo     []formInfo                `json:"form_info"`
	FieldOption  map[int]map[string]string `json:"field_option"`
	RelationInfo *modelRelationInfo        `json:"relation_info"`
	FormRules    interface{}               `json:"form_rules"`
	ModelId      int                       `json:"model_id"`
}

// commonList 公共列表 服务方法
func (b *base) commonList(modelData interface{}, tableName string, params *helper.DataListParams) (interface{}, error) {
	if params.Data.Action == "getField" {
		return b.getFieldList(tableName, params.Data.Filter)
	} else {
		return helper.Page(tableName, db, ns).Filter(modelData, params)
	}
}

// getFieldList 根据表名获取对应所有字段信息
func (b *base) getFieldList(tableName string, filter map[string]interface{}) (*fullModelFields, error) {
	modelConf := ModelConfig().GetModelConfig(tableName)
	modelId := modelConf.ModelId
	fieldInfo := ModelField().GetFieldInfo(modelId)                                                              //模型的字段信息
	formInfo := ModelForm().GetFormInfo(modelId)                                                                 //模型的表单信息
	relationInfo := ModelRelation().GetRelationInfo(modelId, filter)
	formRulesInfo := ModelFormRules().GetRuleListForForm(modelId)
	fieldOption, err := FieldOption().GetFieldOptions(modelId, tableName, modelConf.IsTree, modelConf.LabelFieldName) //模型的关联信息

	if err != nil {
		return nil, err
	}

	return &fullModelFields{
		FieldInfo:    fieldInfo,
		FormInfo:     formInfo,
		FieldOption:  fieldOption,
		RelationInfo: relationInfo,
		FormRules:    formRulesInfo,
		ModelId:      modelId,
	}, nil
}

// getList 根据表名获取对应列表数据(无分页列表数据，如树型列表)
//	参数：
//		dataList  需要填充的列表数据
//		tableName 表名
//		params    过滤条件
func (b *base) getList(dataList interface{}, tableName string, params *helper.DataListParams) {
	query := db.Table(ns.TableName(tableName)).Select("*").Where("status = 10")

	modelConf := ModelConfig().GetModelConfig(tableName)

	if params.Data.Keywords != "" {
		query = query.Where(modelConf.LabelFieldName+" like ?", "%"+params.Data.Keywords+"%")
	}

	orderField := "sort_num"
	if tableName == "roles" {
		orderField = ""
	}

	if orderField != "" {
		query = query.Order(orderField)
	}

	query.Find(dataList)
}

// Create 创建单条数据
func (b *base) Create(data interface{}) (interface{}, error) {
	res := db.Create(data)
	return data, res.Error
}

func (b *base) Update(data interface{}) (interface{}, error) {
	res := db.Updates(data)
	fmt.Println("影响条数：", res.RowsAffected)
	return data, res.Error
}
