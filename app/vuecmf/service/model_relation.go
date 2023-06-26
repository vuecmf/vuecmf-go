//+----------------------------------------------------------------------
// | Copyright (c) 2023 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: vuecmf <tulihua2004@126.com>
// +----------------------------------------------------------------------

package service

import (
	"github.com/vuecmf/vuecmf-go/app/vuecmf/helper"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

// modelRelationService modelRelation服务结构
type modelRelationService struct {
	*BaseService
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
	Options     interface{} `json:"options"`
	Linkage     interface{} `json:"linkage"`
	FullOptions interface{} `json:"full_options"`
}

type fieldInfoST struct {
	FieldId         int
	LinkageFieldId  int
	LinkageActionId int
}

type actionInfoST struct {
	ActionType, TableName string
}

//modelRelationService 联动关联字段信息
//	参数：
//		modelId 模型ID
func (ser *modelRelationService) getRelationLinkage(modelId int) interface{} {
	var fieldInfo []fieldInfoST
	var result = make(map[int]map[int]map[string]string)

	//先取出有关联表的字段及关联信息
	Db.Table(NS.TableName("model_form_linkage")).
		Select("model_field_id field_id, linkage_field_id, linkage_action_id").
		Where("model_id = ?", modelId).
		Where("status = 10").
		Find(&fieldInfo)

	for _, val := range fieldInfo {
		var actionInfo actionInfoST
		Db.Table(NS.TableName("model_action")+" MA").
			Select("MA.action_type, MC.table_name").
			Joins("LEFT JOIN "+NS.TableName("model_config")+" MC ON MA.model_id = MC.id").
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
	FieldId             int    //字段ID
	RelationModelId     int    //关联模型ID
	RelationTableName   string //关联模型的表名
	RelationFieldName   string //关联模型的字段名
	RelationShowFieldId string //需关联显示的字段ID,多个逗号分隔
	RelationFilter      string //关联筛选条件
}

// 关联的字段选项信息
type relationOptions struct {
	Label     string
	FieldName string
}

// getRelationOptions 关联模型的数据列表
//	参数：
//		modelId 模型ID
//		filter 筛选条件
//		db  菜单下拉的db实例
func (ser *modelRelationService) getRelationOptions(modelId int, filter map[string]interface{}, db *gorm.DB) map[string]interface{} {
	var reFieldInfo []relationFieldInfo
	var result = make(map[string]interface{})

	//先取出有关联表的字段及关联信息
	Db.Table(NS.TableName("model_relation")+" VMR").
		Select("model_field_id field_id, relation_model_id, MC.table_name relation_table_name, VMF.field_name relation_field_name, relation_show_field_id, relation_filter").
		Joins("LEFT JOIN "+NS.TableName("model_field")+" VMF ON VMF.id = VMR.relation_field_id").
		Joins("LEFT JOIN "+NS.TableName("model_config")+" MC ON MC.id = VMR.relation_model_id").
		Where("VMR.relation_field_id != 0").
		Where("VMR.model_id = ?", modelId).
		Where("VMR.status = 10").
		Where("VMF.status = 10").
		Find(&reFieldInfo)

	modelTableName := ModelConfig().GetModelTableName(modelId)

	for _, val := range reFieldInfo {
		var options []*helper.ModelFieldOption

		isTree := ModelConfig().IsTree(val.RelationModelId)
		if isTree {
			//若关联的模型是目录树的、则下拉选项需要格式化树型结构
			options = helper.FormatTree(options, db, db.NamingStrategy.TableName(val.RelationTableName), filter, "id", 0, "title", "pid", "sort_num", 1)

		} else {
			var showFieldNameArr []string
			Db.Table(NS.TableName("model_field")).
				Select("field_name").
				Where("id IN ?", strings.Split(val.RelationShowFieldId, ",")).
				Where("status = 10").Find(&showFieldNameArr)

			var reOptions []relationOptions

			if modelTableName == "model_form_rules" && val.RelationTableName == "model_form" && helper.InSlice("model_field_id", showFieldNameArr) && helper.InSlice("type", showFieldNameArr) {
				query := Db.Table(NS.TableName(val.RelationTableName) + " " + val.RelationTableName).
					Select("concat(F.field_name,\"(\",F.label,\")-\",FP.option_label) label, " + val.RelationTableName + "." + val.RelationFieldName + " field_name").
					Joins("LEFT JOIN " + NS.TableName("model_field") + " F ON F.id = " + val.RelationTableName + ".model_field_id and F.status = 10").
					Joins("LEFT JOIN " + NS.TableName("field_option") + " FP ON FP.option_value = " + val.RelationTableName + ".type and FP.status = 10").
					Where(val.RelationTableName + ".status = 10").
					Where("F.status = 10").
					Where("FP.status = 10")

				if filter["model_id"] != nil {
					query = query.Where(val.RelationTableName+".model_id = ?", filter["model_id"])
				}

				if val.RelationFilter != "" {
					query = query.Where(val.RelationFilter)
				}

				query.Find(&reOptions)

			} else {
				showFieldStr := "id"
				if len(showFieldNameArr) > 0 {
					showFieldStr = showFieldNameArr[0]
					showFieldNameArr = helper.SliceRemove(showFieldNameArr, 0)
					if len(showFieldNameArr) > 0 {
						showFieldStr = "concat(" + showFieldStr + ",'('," + strings.Join(showFieldNameArr, ",'-',") + ",')')"
					}
				}

				query := db.Table(db.NamingStrategy.TableName(val.RelationTableName) + " " + val.RelationTableName).
					Select(showFieldStr + " label," + val.RelationFieldName + " field_name").
					Where("status = 10")

				if filter != nil && (val.RelationTableName == "model_field") {
					for field, filterVal := range filter {
						query = query.Where(field+" = ?", filterVal)
					}
				}

				if val.RelationFilter != "" {
					query = query.Where(val.RelationFilter)
				}

				query.Find(&reOptions)
			}

			for _, item := range reOptions {
				options = append(options, &helper.ModelFieldOption{
					Value: item.FieldName,
					Label: item.Label,
				})
			}

		}

		//关联模型的数据列表，供表单中下拉框中使用
		result[strconv.Itoa(val.FieldId)] = options
	}

	return result
}

// GetRelationInfo 获取模型的关联信息
//	参数：
//		modelId 模型ID
//		filter 筛选条件
//		db  菜单下拉的db实例
func (ser *modelRelationService) GetRelationInfo(modelId int, filter map[string]interface{}, db *gorm.DB) *modelRelationInfo {
	mri := &modelRelationInfo{}

	//供表单中与之相关的下拉框联动变化
	mri.Linkage = ser.getRelationLinkage(modelId)

	//供表单中下拉框中使用
	mri.Options = ser.getRelationOptions(modelId, filter, db)

	//供列表及搜索表单下拉框中使用
	delete(filter, "model_id")
	mri.FullOptions = ser.getRelationOptions(modelId, filter, db)

	return mri

}
