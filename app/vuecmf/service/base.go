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

var Db *gorm.DB     //数据库连接实例
var NS schema.Namer //数据库名称服务
var Conf *app.Conf  //系统配置实例

var CacheUser = "vuecmf_user"

type BaseService struct {
}

func init() {
	Conf = app.Config()
	Db = app.Db(strings.ToLower(Conf.Env))
	if Db != nil {
		NS = Db.NamingStrategy
	}
}

// fullModelFields 模型的所有字段相关信息（字段、表单、字段选项、字段关联、表单验证规则）
type fullModelFields struct {
	FieldInfo    []fieldInfo                           `json:"field_info"`
	FormInfo     []formInfo                            `json:"form_info"`
	FieldOption  map[string][]*helper.ModelFieldOption `json:"field_option"`
	RelationInfo *modelRelationInfo                    `json:"relation_info"`
	FormRules    interface{}                           `json:"form_rules"`
	ModelId      int                                   `json:"model_id"`
}

// CommonList 公共列表 服务方法
func (b *BaseService) CommonList(modelData interface{}, tableName string, filterFields []string, params *helper.DataListParams) (interface{}, error) {
	if params.Data.Action == "getField" {
		return b.GetFieldList(tableName, params.Data.Filter)
	} else {
		return helper.Page(tableName, filterFields, Db, NS).Filter(modelData, params)
	}
}

// GetFieldList 根据表名获取对应所有字段信息
func (b *BaseService) GetFieldList(tableName string, filter map[string]interface{}) (*fullModelFields, error) {
	modelCfg := ModelConfig().GetModelConfig(tableName)
	modelId := modelCfg.ModelId
	fieldInfoList := ModelField().GetFieldInfo(modelId) //模型的字段信息
	formInfoList := ModelForm().GetFormInfo(modelId)    //模型的表单信息
	relationInfoList := ModelRelation().GetRelationInfo(modelId, filter)
	formRulesInfoList := ModelFormRules().GetRuleListForForm(modelId)
	fieldOptionList, err := FieldOption().GetFieldOptions(modelId, tableName, modelCfg.IsTree, modelCfg.LabelFieldName) //模型的关联信息

	if err != nil {
		return nil, err
	}

	return &fullModelFields{
		FieldInfo:    fieldInfoList,
		FormInfo:     formInfoList,
		FieldOption:  fieldOptionList,
		RelationInfo: relationInfoList,
		FormRules:    formRulesInfoList,
		ModelId:      modelId,
	}, nil
}

// GetList 根据表名获取对应列表数据(无分页列表数据，如树型列表)
//	参数：
//		dataList  需要填充的列表数据
//		tableName 表名
//		params    过滤条件
func (b *BaseService) GetList(dataList interface{}, tableName string, params *helper.DataListParams) {
	query := Db.Table(NS.TableName(tableName)).Select("*").Where("status = 10")

	modelCfg := ModelConfig().GetModelConfig(tableName)

	data := params.Data

	if data.Keywords != "" {
		query = query.Where(modelCfg.LabelFieldName+" like ?", "%"+data.Keywords+"%")
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
func (b *BaseService) Create(data interface{}) (int64, error) {
	res := Db.Create(data)
	return res.RowsAffected, res.Error
}

// Update 更新数据, 成功返回影响行数
func (b *BaseService) Update(data interface{}) (int64, error) {
	res := Db.Updates(data)
	return res.RowsAffected, res.Error
}

// Detail 根据ID获取详情
func (b *BaseService) Detail(id uint, result interface{}) error {
	res := Db.First(&result, id)
	return res.Error
}

// Delete 根据ID删除数据
func (b *BaseService) Delete(id uint, model interface{}) (int64, error) {
	res := Db.Delete(model, id)
	return res.RowsAffected, res.Error
}

// DeleteBatch 根据ID删除数据， 多个用英文逗号分隔
func (b *BaseService) DeleteBatch(idList string, model interface{}) (int64, error) {
	idArr := strings.Split(idList, ",")
	res := Db.Delete(model, idArr)
	return res.RowsAffected, res.Error
}

type DropdownList struct {
	Value uint   `json:"value"`
	Label string `json:"label"`
}

// Dropdown 获取模型的下拉列表
func (b *BaseService) Dropdown(form *model.DropdownForm, modelName string) (interface{}, error) {
	if form.RelationModelId > 0 {
		form.ModelId = form.RelationModelId
	} else if form.TableName != "" {
		Db.Table(NS.TableName("model_config")).Select("id").
			Where("table_name = ?", form.TableName).
			Find(&form.ModelId)
	}

	modelId := ModelConfig().GetModelId(modelName)
	var labelFieldList []string
	Db.Table(NS.TableName("model_field")).Select("field_name").
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

	query := Db.Table(NS.TableName(modelName)).Select(labelField+" label, id value").Where("status = 10");
	if modelName == "model_config" {
		query = query.Where("app_id = ?", form.AppId)
	}else{
		query = query.Where("model_id = ?", form.ModelId)
	}
	query.Find(&result)

	return result, nil

}

var base *BaseService

// Base 获取BaseService服务实例
func Base() *BaseService {
	if base == nil {
		base = &BaseService{}
	}
	return base
}
