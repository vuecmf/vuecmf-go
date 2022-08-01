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
	"github.com/vuecmf/vuecmf-go/app/vuecmf/helper"
	"strconv"
	"strings"
)

// modelRelationService modelRelation服务结构
type modelRelationService struct {
	*base
}

var modelRelation *modelRelationService

// ModelRelation 获取modelRelation服务实例
func ModelRelation() *modelRelationService {
	if modelRelation == nil {
		modelRelation = &modelRelationService{}
	}
	return modelRelation
}

type modelRelationInfo struct {
	Options interface{} `json:"options"`
	Linkage interface{} `json:"linkage"`
	FullOptions interface{} `json:"full_options"`
}

type fieldInfoST struct {
	FieldId int
	LinkageFieldId int
	LinkageActionId int
}

type actionInfoST struct {
	ActionType, TableName string
}

//modelRelationService 联动关联字段信息
func (ser *modelRelationService) getRelationLinkage (modelId int) interface{} {
	var fieldInfo []fieldInfoST
	var result = make(map[int]map[int]map[string]string)

	//先取出有关联表的字段及关联信息
	db.Table(ns.TableName("model_form_linkage")).
		Select("model_field_id field_id, linkage_field_id, linkage_action_id").
		Where("model_id = ?", modelId).
		Where("status = 10").
		Find(&fieldInfo)

	for _, val := range fieldInfo {
		var actionInfo actionInfoST
		db.Table(ns.TableName("model_action") + " MA").
			Select("MA.action_type, MC.table_name").
			Joins("LEFT JOIN " + ns.TableName("model_config") + " MC ON MA.model_id = MC.id").
			Where("MA.id = ?", val.LinkageActionId).
			Where("MA.status = 10").
			Where("MC.status = 10").
			Find(&actionInfo)

		//联动关联字段信息, 供表单中与之相关的下拉框联动变化
		if result[val.FieldId] == nil {
			result[val.FieldId] = make(map[int]map[string]string)
		}
		if result[val.FieldId][val.LinkageFieldId] == nil {
			result[val.FieldId][val.LinkageFieldId] = make(map[string]string)
		}

		result[val.FieldId][val.LinkageFieldId]["relation_field_id"] = strconv.Itoa(val.LinkageFieldId)
		result[val.FieldId][val.LinkageFieldId]["action_table_name"] = actionInfo.TableName
		result[val.FieldId][val.LinkageFieldId]["action_type"] = actionInfo.ActionType
	}

	return result
}

// 关联模型的字段信息
type relationFieldInfo struct {
	FieldId int  //字段ID
	RelationModelId int  //关联模型ID
	RelationTableName string //关联模型的表名
	RelationFieldName string //关联模型的字段名
	RelationShowFieldId string  //需关联显示的字段ID,多个逗号分隔
}

// 关联的字段选项信息
type relationOptions struct {
	Label string
	FieldName string
}

// getRelationOptions 关联模型的数据列表
func (ser *modelRelationService) getRelationOptions (modelId int, filter map[string]interface{}) interface{} {
	var fieldInfo []relationFieldInfo
	result := map[int]map[string]string{}
	options := map[string]string{}

	//先取出有关联表的字段及关联信息
	db.Table(ns.TableName("model_relation") + " VMR").
		Select("model_field_id field_id, relation_model_id, MC.table_name relation_table_name, VMF.field_name relation_field_name, relation_show_field_id").
		Joins("LEFT JOIN " + ns.TableName("model_field") + " VMF ON VMF.id = VMR.relation_field_id").
		Joins("LEFT JOIN " + ns.TableName("model_config") + " MC ON MC.id = VMR.relation_model_id").
		Where("VMR.relation_field_id != 0").
		Where("VMR.model_id = ?", modelId).
		Where("VMR.status = 10").
		Where("VMF.status = 10").
		Find(&fieldInfo)

	modelTableName := ModelConfig().GetModelTableName(modelId)

	for _, val := range fieldInfo {
		isTree := ModelConfig().IsTree(val.RelationModelId)
		if isTree {
			//若关联的模型是目录树的、则下拉选项需要格式化树型结构
			helper.FormatTree(options, db, ns.TableName(val.RelationTableName), "id", 0, "title", "pid", "sort_num", 1)

		} else {
			var showFieldNameArr []string
			db.Table(ns.TableName("model_field")).
				Select("field_name").
				Where("id IN ?", strings.Split(val.RelationShowFieldId, ",")).
				Where("status = 10").Find(&showFieldNameArr)

			var reOptions []relationOptions

			if modelTableName == "model_form_rules" && val.RelationTableName == "model_form" && helper.InSlice("model_field_id", showFieldNameArr)  && helper.InSlice("type", showFieldNameArr) {
				query := db.Table(ns.TableName(val.RelationTableName) + " A").
					Select("concat(F.field_name,\"(\",F.label,\")-\",FP.option_label) label, A." + val.RelationFieldName + " field_name").
					Joins("LEFT JOIN " + ns.TableName("model_field") + " F ON F.id = A.model_field_id and F.status = 10").
					Joins("LEFT JOIN " + ns.TableName("field_option") + " FP ON FP.option_value = A.type and FP.status = 10").
					Where("A.status = 10").
					Where("F.status = 10").
					Where("FP.status = 10")

				if filter["model_id"] != nil {
					query.Where("A.model_id = ?", filter["model_id"])
				}

				query.Find(&reOptions)

			}else{
				showFieldStr := "id"
				if len(showFieldNameArr) > 0 {
					showFieldStr = showFieldNameArr[0]
					helper.SliceRemove(showFieldNameArr, 0)
					if len(showFieldNameArr) > 0 {
						showFieldStr = "concat(" + showFieldStr + ",'('," + strings.Join(showFieldNameArr,",'-',") + ",')')"
					}
				}

				query := db.Table(ns.TableName(val.RelationTableName)).
					Select(showFieldStr + " label," + val.RelationFieldName + " field_name").
					Where("status = 10")

				if filter != nil && helper.InSlice(val.RelationTableName, []string{"model_field","model_action"}) {
					for field, filterVal := range filter {
						if field == "model_id" {
							//取出所关联的模型ID
							var modelIdArr []int
							db.Table(ns.TableName("model_relation")).Select("relation_model_id").
								Where("model_id = ?", filterVal).Find(&modelIdArr)
							modelIdArr = append(modelIdArr, helper.InterfaceToInt(filterVal))
							query = query.Where(field + " IN ?", modelIdArr)

						} else {
							query = query.Where(field + " = ?", filterVal)
						}
					}
				}

				query.Find(&reOptions)
			}

			for _, item := range reOptions {
				options[item.FieldName] = item.Label
			}

		}

		//关联模型的数据列表，供表单中下拉框中使用
		result[val.FieldId] = options

	}

	return result
}

// GetRelationInfo 获取模型的关联信息
func (ser *modelRelationService) GetRelationInfo (modelId int, filter map[string]interface{}) *modelRelationInfo {
	mri := &modelRelationInfo{}

	//供表单中与之相关的下拉框联动变化
	mri.Linkage = ser.getRelationLinkage(modelId)

	//供表单中下拉框中使用
	mri.Options = ser.getRelationOptions(modelId, filter)

	//供列表及搜索表单下拉框中使用
	delete(filter, "model_id")
	mri.FullOptions = ser.getRelationOptions(modelId, filter)

	return mri

}