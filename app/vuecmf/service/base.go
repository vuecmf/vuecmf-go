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
	"strings"
)

var db *gorm.DB
var ns schema.Namer
var CacheUser = "vuecmf_user"
var config *app.AppConfig

type baseService struct {
}

func init() {
	db = app.Db("default")
	ns = db.NamingStrategy
	config = app.Config()
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

// CommonList 公共列表 服务方法
func (b *baseService) CommonList(modelData interface{}, tableName string, filterFields []string, params *helper.DataListParams) (interface{}, error) {
	if params.Data.Action == "getField" {
		return b.getFieldList(tableName, params.Data.Filter)
	} else {
		return helper.Page(tableName, filterFields, db, ns).Filter(modelData, params)
	}
}

// getFieldList 根据表名获取对应所有字段信息
func (b *baseService) getFieldList(tableName string, filter map[string]interface{}) (*fullModelFields, error) {
	modelConf := ModelConfig().GetModelConfig(tableName)
	modelId := modelConf.ModelId
	fieldInfo := ModelField().GetFieldInfo(modelId) //模型的字段信息
	formInfo := ModelForm().GetFormInfo(modelId)    //模型的表单信息
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
func (b *baseService) getList(dataList interface{}, tableName string, params *helper.DataListParams) {
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

// Create 创建单条或多条数据, 成功返回影响行数
func (b *baseService) Create(data interface{}) (int64, error) {
	res := db.Create(data)
	return res.RowsAffected, res.Error
}

// Update 更新数据, 成功返回影响行数
func (b *baseService) Update(data interface{}) (int64, error) {
	res := db.Updates(data)
	return res.RowsAffected, res.Error
}

// Detail 根据ID获取详情
func (b *baseService) Detail(id uint, result interface{}) error {
	res := db.First(&result, id)
	return res.Error
}

// Delete 根据ID删除数据
func (b *baseService) Delete(id uint, model interface{}) (int64, error) {
	res := db.Delete(model, id)
	return res.RowsAffected, res.Error
}

// DeleteBatch 根据ID删除数据， 多个用英文逗号分隔
func (b *baseService) DeleteBatch(idList string, model interface{}) (int64, error) {
	idArr := strings.Split(idList, ",")
	res := db.Delete(model, idArr)
	return res.RowsAffected, res.Error
}

type DropdownList struct {
	Id    uint   `json:"id"`
	Label string `json:"label"`
}

// Dropdown 获取模型的下拉列表
func (b *baseService) Dropdown(form *model.DropdownForm, modelName string) (interface{}, error) {
	if form.RelationModelId > 0 {
		form.ModelId = form.RelationModelId
	}
	if form.TableName == "" && form.ModelId == 0 {
		return nil, nil
	}

	if form.TableName != "" {
		db.Table(ns.TableName("model_config")).Select("id").
			Where("table_name = ?", form.TableName).
			Where("status = 10").
			Find(&form.ModelId)
	}

	modelId := ModelConfig().GetModelId(modelName)
	var labelFieldList []string
	db.Table(ns.TableName("model_field")).Select("field_name").
		Where("model_id = ?", modelId).
		Where("is_label = 10").
		Where("status = 10").
		Find(&labelFieldList)

	labelField := "id"

	if len(labelFieldList) > 0 {
		labelField = labelFieldList[0]
		labelFieldList = helper.SliceRemove(labelFieldList, 0)
		if len(labelFieldList) > 0 {
			labelField = "concat(" + labelField + ",'('," + strings.Join(labelFieldList, ",'-',") + ",')')"
		}
	}

	var result []DropdownList

	db.Table(ns.TableName(modelName)).Select(labelField+" label, id").
		Where("model_id = ?", form.ModelId).
		Where("status = 10").
		Find(&result)

	return result, nil

}

var base *baseService

// Base 获取baseService服务实例
func Base() *baseService {
	if base == nil {
		base = &baseService{}
	}
	return base
}
