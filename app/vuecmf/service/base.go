//+----------------------------------------------------------------------
// | Copyright (c) 2024 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: tulihua2004@126.com
// +----------------------------------------------------------------------

// Package service 服务
package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vuecmf/vuecmf-go/v3/app"
	"github.com/vuecmf/vuecmf-go/v3/app/vuecmf/helper"
	"github.com/vuecmf/vuecmf-go/v3/app/vuecmf/model"
	"gorm.io/gorm"
	"runtime"
	"strconv"
	"strings"
)

var CacheUser = "vuecmf_user"

type BaseService struct {
	TableName    string      //表名称
	Model        interface{} //表对应的模型实例
	ListData     interface{} //存储列表结果
	FilterFields []string    //支持模糊查询的字段
}

// DbTable 获取数据库表连接实例
//
//	参数：
//	tableName 数据库表名(不含表前缀)
func DbTable(tableName string, alias ...string) *gorm.DB {
	if len(alias) > 0 {
		return app.Db.Table(TableName(tableName) + " " + alias[0])
	}
	return app.Db.Table(TableName(tableName))
}

// TableName 获取完整数据库表名
//
//	参数：
//	tableName 数据库表名(不含表前缀)
func TableName(tableName string) string {
	return app.NS.TableName(tableName)
}

// GetErrMsg 获取异常信息
func GetErrMsg(err error) string {
	if app.Cfg.Debug == false {
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

// FullModelFields 模型的所有字段相关信息（字段、表单、字段选项、字段关联、表单验证规则）
type FullModelFields struct {
	FieldInfo    []FieldInfo                           `json:"field_info"`
	FormInfo     []formInfo                            `json:"form_info"`
	FieldOption  map[string][]*helper.ModelFieldOption `json:"field_option"`
	RelationInfo *modelRelationInfo                    `json:"relation_info"`
	FormRules    interface{}                           `json:"form_rules"`
	ModelId      int                                   `json:"model_id"`
}

type DropdownList struct {
	Value uint   `json:"value"`
	Label string `json:"label"`
}

type FieldOptionDropdownList struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// CommonList 公共列表 服务方法
//
//	参数：
//		params 查询条件参数
//		isSuper 是否为超级管理员
func (svc *BaseService) CommonList(params *helper.DataListParams, isSuper uint16) (interface{}, error) {
	if params.Data.Action == "getField" {
		return svc.GetFieldList(svc.TableName, params.Data.Filter, isSuper)
	} else {
		return helper.Page(svc.TableName, svc.FilterFields, app.Db, app.NS).Filter(svc.ListData, params.Data)
	}
}

// GetFieldList 根据表名获取对应所有字段信息
//
//	参数：
//		tableName 表名
//		filter 查询条件参数
//		isSuper 是否为超级管理员
func (svc *BaseService) GetFieldList(tableName string, filter map[string]interface{}, isSuper uint16, relationFilter ...map[string]map[string]interface{}) (*FullModelFields, error) {
	modelCfg := ModelConfig().GetModelConfig(tableName)
	modelId := modelCfg.ModelId
	fieldInfoList := ModelField().GetFieldInfo(modelId)       //模型的字段信息
	formInfoList := ModelForm().GetFormInfo(modelId, isSuper) //模型的表单信息

	reFilter := make(map[string]map[string]interface{})
	if 0 != len(relationFilter) {
		reFilter = relationFilter[0]
	}

	relationInfoList := ModelRelation().GetRelationInfo(modelId, filter, app.Db, reFilter)
	formRulesInfoList := ModelFormRules().GetRuleListForForm(modelId)
	fieldOptionList, err := FieldOption().GetFieldOptions(modelId, tableName, modelCfg.IsTree, modelCfg.LabelFieldName, filter, app.Db) //模型的关联信息

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
//
//	参数：
//		dataList  需要填充的列表数据
//		tableName 表名
//		params    过滤条件
func (svc *BaseService) GetList(dataList interface{}, tableName string, data *helper.ListParams) {
	query := DbTable(tableName).Select("*").Where("status = 10")

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

// Create 创建单条数据, 成功返回影响行数
//
//	参数：
//		data 保存的数据
func (svc *BaseService) Create(data interface{}) (int64, error) {
	res := app.Db.Create(data)
	return res.RowsAffected, res.Error
}

// CreateAll 创建多条数据, 成功返回影响行数
//
//	参数：
//		data 保存的JSON字符串数据
func (svc *BaseService) CreateAll(data string) (int64, error) {
	err := json.Unmarshal([]byte(data), &svc.ListData)
	if err != nil {
		return 0, err
	}
	res := app.Db.Create(svc.ListData)
	return res.RowsAffected, res.Error
}

// Update 更新数据, 成功返回影响行数
//
//	参数：
//		data 更新的数据
func (svc *BaseService) Update(data interface{}) (int64, error) {
	res := app.Db.Updates(data)
	return res.RowsAffected, res.Error
}

// Detail 根据ID获取详情
//
//	参数：
//		id 查询ID
func (svc *BaseService) Detail(id uint) (interface{}, error) {
	res := app.Db.First(&svc.Model, id)
	return svc.Model, res.Error
}

// Delete 根据ID删除数据
//
//	参数：
//		id 需要删除的ID
func (svc *BaseService) Delete(id uint) (int64, error) {
	res := app.Db.Delete(svc.Model, id)
	return res.RowsAffected, res.Error
}

// DeleteBatch 根据ID删除数据， 多个用英文逗号分隔
//
//	参数：
//		idList 需要删除的ID列表
//		model 模型实例
func (svc *BaseService) DeleteBatch(idList string) (int64, error) {
	idArr := strings.Split(idList, ",")
	res := app.Db.Delete(svc.Model, idArr)
	return res.RowsAffected, res.Error
}

// Dropdown 获取模型的下拉列表
//
//	参数：
//		postData 模型的下拉筛选条件表单
//		modelName 模型名称
func (svc *BaseService) Dropdown(postData *model.DropdownForm) (interface{}, error) {
	var labelFieldIdList []string
	var re model.ModelRelation

	if postData.RelationModelId > 0 {
		postData.ModelId = postData.RelationModelId
	} else if postData.TableName != "" {
		DbTable("model_config").Select("id").
			Where("table_name = ?", postData.TableName).
			Find(&postData.ModelId)
	} else if postData.LinkageFieldId > 0 {
		//联动模型的下拉
		DbTable("model_relation").Select("relation_model_id, relation_show_field_id").
			Where("model_field_id = ?", postData.LinkageFieldId).
			Where("status = 10").
			Find(&re)
		if re.RelationModelId == 0 {
			return nil, errors.New("联动字段没有设置模型关联")
		}
		postData.ModelId = re.RelationModelId
	}

	if postData.LinkageFieldId == 0 && postData.ModelFieldId == 0 {
		DbTable("model_relation").Select("relation_show_field_id").
			Where("relation_model_id = ?", postData.ModelId).
			Where("status = 10").
			Limit(1).
			Find(&re)
	}

	labelField := "id"
	valueField := "id"

	fieldType := ""

	var labelFieldList []string
	if svc.TableName == "model_field" {
		labelFieldList = []string{"field_name", "label"}
	} else if svc.TableName == "model_action" {
		labelFieldList = []string{"action_type", "label"}
	} else if svc.TableName == "field_option" {
		labelFieldList = []string{"option_label"}
		valueField = "option_value"
		if postData.ModelFieldId > 0 {
			DbTable("model_field").Select("type fieldType").
				Where("id = ?", postData.ModelFieldId).
				Where("status = 10").
				Limit(1).Find(&fieldType)
		}

	} else {
		labelQuery := DbTable("model_field").Select("field_name").Where("status = 10")

		if re.RelationShowFieldId != "" {
			//字段有模型关联的，下拉列表中显示的字段 从模型关联中的设置的显示字段
			labelFieldIdList = strings.Split(re.RelationShowFieldId, ",")
			labelQuery = labelQuery.Where("id in ?", labelFieldIdList)
		} else {
			//没有模型模型关联的，下拉列表中直接显示标题字段
			modelId := ModelConfig().GetModelId(svc.TableName)
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

	query := DbTable(svc.TableName).Select(labelField + " label, " + valueField + " value").Where("status = 10")
	if svc.TableName == "model_config" {
		query = query.Where("app_id = ?", postData.AppId)
	} else if postData.ModelFieldId > 0 {
		query = query.Where("model_field_id = ?", postData.ModelFieldId)
	} else {
		query = query.Where("model_id = ?", postData.ModelId)
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
