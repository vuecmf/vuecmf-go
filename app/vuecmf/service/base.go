//+----------------------------------------------------------------------
// | Copyright (c) 2023 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: vuecmf <tulihua2004@126.com>
// +----------------------------------------------------------------------

// Package service 服务
package service

import (
	"errors"
	"fmt"
	"github.com/vuecmf/vuecmf-go/app"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/helper"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/model"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"runtime"
	"strconv"
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
	if Conf != nil {
		Db = app.Db(strings.ToLower(Conf.Env))
	}

	if Db != nil {
		NS = Db.NamingStrategy
	}
}

// FullModelFields 模型的所有字段相关信息（字段、表单、字段选项、字段关联、表单验证规则）
type FullModelFields struct {
	FieldInfo    []FieldInfo                           `json:"field_info"`
	FormInfo     []formInfo                            `json:"form_info"`
	FieldOption  map[string][]*helper.ModelFieldOption `json:"field_option"`
	RelationInfo *modelRelationInfo                    `json:"relation_info"`
	FormRules    interface{}                           `json:"form_rules"`
	ModelId      int                                   `json:"model_id"`
}

//GetErrMsg 获取异常信息
func GetErrMsg(err error) string {
	if Conf.Debug == false {
		return err.Error()
	}

	if err != nil {
		prefix := "[error] "
		res := ""
		red := string([]byte{27, 91, 51, 49, 109})
		reset := string([]byte{27, 91, 48, 109})

		for i := 1; i < 6; i++ {
			pc, file, line, rs := runtime.Caller(i)

			if rs {
				errMsg := fmt.Sprintf("%s 在 %s ，文件 %s 第%d行; ", err.Error(), runtime.FuncForPC(pc).Name(), file, line)
				if i == 2 {
					res = prefix + errMsg
				}
				//红色显示打印
				fmt.Println(red, prefix+errMsg, reset)
			}
		}
		return res
	}
	return ""
}

// CommonList 公共列表 服务方法
//	参数：
// 		modelData 存储列表数据
//		tableName 表名
//		filterFields 需要模糊查询的字段列表
//		params 查询条件参数
//		isSuper 是否为超级管理员
func (b *BaseService) CommonList(modelData interface{}, tableName string, filterFields []string, params *helper.DataListParams, isSuper int) (interface{}, error) {
	if params.Data.Action == "getField" {
		return b.GetFieldList(tableName, params.Data.Filter, isSuper)
	} else {
		return helper.Page(tableName, filterFields, Db, NS).Filter(modelData, params.Data)
	}
}

// GetFieldList 根据表名获取对应所有字段信息
//	参数：
//		tableName 表名
//		filter 查询条件参数
//		isSuper 是否为超级管理员
func (b *BaseService) GetFieldList(tableName string, filter map[string]interface{}, isSuper int) (*FullModelFields, error) {
	modelCfg := ModelConfig().GetModelConfig(tableName)
	modelId := modelCfg.ModelId
	fieldInfoList := ModelField().GetFieldInfo(modelId)       //模型的字段信息
	formInfoList := ModelForm().GetFormInfo(modelId, isSuper) //模型的表单信息
	relationInfoList := ModelRelation().GetRelationInfo(modelId, filter, Db)
	formRulesInfoList := ModelFormRules().GetRuleListForForm(modelId)
	fieldOptionList, err := FieldOption().GetFieldOptions(modelId, tableName, modelCfg.IsTree, modelCfg.LabelFieldName, filter, Db) //模型的关联信息

	if err != nil {
		return nil, err
	}

	return &FullModelFields{
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
func (b *BaseService) GetList(dataList interface{}, tableName string, data *helper.ListParams) {
	query := Db.Table(NS.TableName(tableName)).Select("*").Where("status = 10")

	modelCfg := ModelConfig().GetModelConfig(tableName)

	if data.Keywords != "" {
		arr := strings.Split(data.Keywords, ",")
		var conditionArr []string
		params := make(map[string]interface{})
		for k, v := range arr {
			key := strconv.Itoa(k)
			conditionArr = append(conditionArr, modelCfg.LabelFieldName+" like @kw_"+key)
			params["kw_"+key] = v + "%"
		}

		query = query.Where(strings.Join(conditionArr, " or "), params)
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
//	参数：
// 		data 保存的数据
func (b *BaseService) Create(data interface{}) (int64, error) {
	res := Db.Create(data)
	return res.RowsAffected, res.Error
}

// Update 更新数据, 成功返回影响行数
//	参数：
// 		data 更新的数据
func (b *BaseService) Update(data interface{}) (int64, error) {
	res := Db.Updates(data)
	return res.RowsAffected, res.Error
}

// Detail 根据ID获取详情
//	参数：
// 		id 查询ID
//		result 存储结果数据
func (b *BaseService) Detail(id uint, result interface{}) error {
	res := Db.First(&result, id)
	return res.Error
}

// Delete 根据ID删除数据
//	参数：
// 		id 需要删除的ID
//		model 模型实例
func (b *BaseService) Delete(id uint, model interface{}) (int64, error) {
	res := Db.Delete(model, id)
	return res.RowsAffected, res.Error
}

// DeleteBatch 根据ID删除数据， 多个用英文逗号分隔
//	参数：
// 		idList 需要删除的ID列表
//		model 模型实例
func (b *BaseService) DeleteBatch(idList string, model interface{}) (int64, error) {
	idArr := strings.Split(idList, ",")
	res := Db.Delete(model, idArr)
	return res.RowsAffected, res.Error
}

type DropdownList struct {
	Value uint   `json:"value"`
	Label string `json:"label"`
}

type FieldOptionDropdownList struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// Dropdown 获取模型的下拉列表
//	参数：
// 		form 模型的下拉筛选条件表单
//		modelName 模型名称
func (b *BaseService) Dropdown(form *model.DropdownForm, modelName string) (interface{}, error) {
	var labelFieldIdList []string
	var re model.ModelRelation

	if form.RelationModelId > 0 {
		form.ModelId = form.RelationModelId
	} else if form.TableName != "" {
		Db.Table(NS.TableName("model_config")).Select("id").
			Where("table_name = ?", form.TableName).
			Find(&form.ModelId)
	} else if form.LinkageFieldId > 0 {
		//联动模型的下拉
		Db.Table(NS.TableName("model_relation")).Select("relation_model_id, relation_show_field_id").
			Where("model_field_id = ?", form.LinkageFieldId).
			Where("status = 10").
			Find(&re)
		if re.RelationModelId == 0 {
			return nil, errors.New("联动字段没有设置模型关联")
		}
		form.ModelId = re.RelationModelId
	}

	if form.LinkageFieldId == 0 && form.ModelFieldId == 0 {
		Db.Table(NS.TableName("model_relation")).Select("relation_show_field_id").
			Where("relation_model_id = ?", form.ModelId).
			Where("status = 10").
			Limit(1).
			Find(&re)
	}

	labelField := "id"
	valueField := "id"

	fieldType := ""

	var labelFieldList []string
	if modelName == "model_field" {
		labelFieldList = []string{"field_name", "label"}
	} else if modelName == "model_action" {
		labelFieldList = []string{"action_type", "label"}
	} else if modelName == "field_option" {
		labelFieldList = []string{"option_label"}
		valueField = "option_value"
		if form.ModelFieldId > 0 {
			Db.Table(NS.TableName("model_field")).Select("type fieldType").
				Where("id = ?", form.ModelFieldId).
				Where("status = 10").
				Limit(1).Find(&fieldType)
		}

	} else {
		labelQuery := Db.Table(NS.TableName("model_field")).Select("field_name").Where("status = 10")

		if re.RelationShowFieldId != "" {
			//字段有模型关联的，下拉列表中显示的字段 从模型关联中的设置的显示字段
			labelFieldIdList = strings.Split(re.RelationShowFieldId, ",")
			labelQuery = labelQuery.Where("id in ?", labelFieldIdList)
		} else {
			//没有模型模型关联的，下拉列表中直接显示标题字段
			modelId := ModelConfig().GetModelId(modelName)
			labelQuery = labelQuery.Where("is_label = 10").Where("model_id = ?", modelId)
		}

		labelQuery.Find(&labelFieldList)
	}

	if len(labelFieldList) > 0 {
		labelField = labelFieldList[0]
		labelFieldList = helper.SliceRemove(labelFieldList, 0)
		if len(labelFieldList) > 0 {
			labelField = "concat(" + labelField + ",'('," + strings.Join(labelFieldList, ",'-',") + ",')')"
		}
	}

	query := Db.Table(NS.TableName(modelName)).Select(labelField + " label, " + valueField + " value").Where("status = 10")
	if modelName == "model_config" {
		query = query.Where("app_id = ?", form.AppId)
	} else if form.ModelFieldId > 0 {
		query = query.Where("model_field_id = ?", form.ModelFieldId)
	} else {
		query = query.Where("model_id = ?", form.ModelId)
	}

	if fieldType == "varchar" || fieldType == "char" {
		var result []FieldOptionDropdownList
		query.Find(&result)
		return result, nil
	} else {
		var result []DropdownList
		query.Find(&result)
		return result, nil
	}

}

var base *BaseService

// Base 获取BaseService服务实例
func Base() *BaseService {
	if base == nil {
		base = &BaseService{}
	}
	return base
}
