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
	"github.com/vuecmf/vuecmf-go/app/vuecmf/model"
	"io/ioutil"
	"strconv"
	"strings"
)

// makeService make服务结构
type makeService struct {
	*base
}

//Model 功能：生成模型代码文件
//		参数：tableName string 表名（不带表前缀）
func (makeSer *makeService) Model(tableName string) bool {
	var result []model.ModelField

	//查出需要生成模型表的字段相关信息
	db.Table(ns.TableName("model_field")+" MF").
		Select("MF.*").
		Joins("left join "+ns.TableName("model_config")+" MC on MF.model_id = MC.id").
		Where("MF.field_name NOT IN('id','status')").
		Where("MC.table_name = ?", tableName).Scan(&result)

	tplFile := "model.stub"
	modelConf := ModelConfig().GetModelConfig(tableName)
	if modelConf.IsTree == true {
		tplFile = "tree_model.stub"
	}

	//读取模型模板文件
	tplContent, err := ioutil.ReadFile("app/vuecmf/make/stubs/" + tplFile)
	if err != nil {
		panic("读取model模板失败")
	}

	modelContent := ""
	hasTime := false

	//模型字段信息处理
	for _, value := range result {
		notNull := ""
		defaultVal := ""
		size := ""
		autoCreateTime := ""
		fieldType := "string"
		uniqueIndex := ""

		if value.Type == "timestamp" {
			hasTime = true
			fieldType = "time.Time"
		} else if value.Type == "int" || value.Type == "bigint" {
			fieldType = "int"
			if modelConf.IsTree == true {
				fieldType = "uint"
			}
		} else if value.Type == "smallint" {
			fieldType = "int16"
		} else if value.Type == "tinyint" {
			fieldType = "int8"
		} else if value.Type == "float" {
			fieldType = "float32"
		} else if value.Type == "double" || value.Type == "decimal" {
			fieldType = "float64"
		}

		if value.IsNull == 20 {
			notNull = "not null;"
		}

		if value.FieldName == "update_time" || value.FieldName == "last_login_time" || value.DefaultValue == "CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" {
			autoCreateTime = "autoCreateTime;autoUpdateTime;"
		} else if value.DefaultValue == "CURRENT_TIMESTAMP" {
			autoCreateTime = "autoCreateTime;"
		} else {
			defaultVal = "default:" + value.DefaultValue + ";"
			size = "size:" + strconv.Itoa(value.Length) + ";"
		}

		//字段唯一索引处理
		modelIndexId := 0
		id := strconv.Itoa(int(value.Id))
		db.Table(ns.TableName("model_index")).Select("id").
			Where("model_field_id = ? or model_field_id like ? or model_field_id like ?", id, id+",%", "%,"+id).
			Find(&modelIndexId)

		if modelIndexId > 0 {
			uniqueIndex = "uniqueIndex:unique_index;"
		}

		modelContent += helper.UnderToCamel(value.FieldName) + " " + fieldType + " `json:\"" + value.FieldName +
			"\" gorm:\"column:" + value.FieldName + ";" + size + uniqueIndex + notNull + autoCreateTime + defaultVal +
			"comment:" + value.Note + "\"`\n\t"
	}

	//获取模型标签名称
	modelLabel := ""
	db.Table(ns.TableName("model_config")).Select("label").
		Where("table_name = ?", tableName).Find(&modelLabel)

	modelName := helper.UnderToCamel(tableName)
	modelValue := strings.ToLower(modelName)

	//替换模板文件中内容
	txt := string(tplContent)
	txt = strings.Replace(txt, "{{.comment}}", modelLabel, -1)
	txt = strings.Replace(txt, "{{.model_name}}", modelName, -1)
	txt = strings.Replace(txt, "{{.model_value}}", modelValue, -1)

	if hasTime == true {
		txt = strings.Replace(txt, "{{.import}}", "import \"time\"", -1)
	} else {
		txt = strings.Replace(txt, "{{.import}}", "", -1)
	}

	txt = strings.Replace(txt, "{{.body}}", modelContent, -1)

	err = ioutil.WriteFile("app/vuecmf/model/"+tableName+".go", []byte(txt), 0666)

	if err != nil {
		return false
	}

	return true
}

//Service 功能：生成服务代码文件
//		  参数：tableName string 表名（不带表前缀）
func (makeSer *makeService) Service(tableName string) bool {
	serviceMethod := helper.UnderToCamel(tableName)
	nameArr := []rune(serviceMethod)
	nameArr[0] += 32
	serviceName := string(nameArr)

	tplFile := "service.stub"
	modelConf := ModelConfig().GetModelConfig(tableName)
	if modelConf.IsTree == true {
		tplFile = "tree_service.stub"
	}

	tplContent, err := ioutil.ReadFile("app/vuecmf/make/stubs/" + tplFile)
	if err != nil {
		panic("读取service模板失败")
	}

	txt := string(tplContent)
	txt = strings.Replace(txt, "{{.service_name}}", serviceName, -1)
	txt = strings.Replace(txt, "{{.service_method}}", serviceMethod, -1)

	err = ioutil.WriteFile("app/vuecmf/service/"+tableName+".go", []byte(txt), 0666)

	if err != nil {
		return false
	}

	return true
}

//Controller 功能：生成控制器代码文件
//		  参数：tableName string 表名（不带表前缀）
func (makeSer *makeService) Controller(tableName string) bool {
	controllerName := helper.UnderToCamel(tableName)

	tplContent, err := ioutil.ReadFile("app/vuecmf/make/stubs/controller.stub")
	if err != nil {
		panic("读取controller模板失败")
	}

	txt := string(tplContent)
	txt = strings.Replace(txt, "{{.controller_name}}", controllerName, -1)

	err = ioutil.WriteFile("app/vuecmf/controller/"+tableName+".go", []byte(txt), 0666)

	if err != nil {
		return false
	}

	return true
}

//Form 功能：生成表单代码文件
//		参数：tableName string 表名（不带表前缀）
func (makeSer *makeService) Form(tableName string) bool {
	var result []model.ModelField

	//查出需要生成模型表的字段相关信息
	db.Table(ns.TableName("model_field")+" MF").
		Select("MF.*").
		Joins("left join "+ns.TableName("model_config")+" MC on MF.model_id = MC.id").
		Where("MF.field_name NOT IN('id','status')").
		Where("MC.table_name = ?", tableName).Scan(&result)

	//读取模型模板文件
	tplContent, err := ioutil.ReadFile("app/vuecmf/make/stubs/form.stub")
	if err != nil {
		panic("读取form模板失败")
	}

	formContent := ""
	hasTime := false

	//模型字段信息处理
	for _, value := range result {
		notNull := ""
		defaultVal := ""
		size := ""
		autoCreateTime := ""
		fieldType := "string"
		uniqueIndex := ""

		if value.Type == "timestamp" {
			hasTime = true
			fieldType = "time.Time"
		} else if value.Type == "int" || value.Type == "bigint" {
			fieldType = "int"

		} else if value.Type == "smallint" {
			fieldType = "int16"
		} else if value.Type == "tinyint" {
			fieldType = "int8"
		} else if value.Type == "float" {
			fieldType = "float32"
		} else if value.Type == "double" || value.Type == "decimal" {
			fieldType = "float64"
		}

		if value.IsNull == 20 {
			notNull = "not null;"
		}

		if value.FieldName == "update_time" || value.FieldName == "last_login_time" || value.DefaultValue == "CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" {
			autoCreateTime = "autoCreateTime;autoUpdateTime;"
		} else if value.DefaultValue == "CURRENT_TIMESTAMP" {
			autoCreateTime = "autoCreateTime;"
		} else {
			defaultVal = "default:" + value.DefaultValue + ";"
			size = "size:" + strconv.Itoa(value.Length) + ";"
		}

		//字段唯一索引处理
		modelIndexId := 0
		id := strconv.Itoa(int(value.Id))
		db.Table(ns.TableName("model_index")).Select("id").
			Where("model_field_id = ? or model_field_id like ? or model_field_id like ?", id, id+",%", "%,"+id).
			Find(&modelIndexId)

		if modelIndexId > 0 {
			uniqueIndex = "uniqueIndex:unique_index;"
		}

		formContent += helper.UnderToCamel(value.FieldName) + " " + fieldType + " `json:\"" + value.FieldName +
			"\" gorm:\"column:" + value.FieldName + ";" + size + uniqueIndex + notNull + autoCreateTime + defaultVal +
			"comment:" + value.Note + "\"`\n\t"
	}

	//获取模型标签名称
	modelLabel := ""
	db.Table(ns.TableName("model_config")).Select("label").
		Where("table_name = ?", tableName).Find(&modelLabel)

	formName := helper.UnderToCamel(tableName)

	//替换模板文件中内容
	txt := string(tplContent)
	txt = strings.Replace(txt, "{{.comment}}", modelLabel, -1)
	txt = strings.Replace(txt, "{{.form_name}}", formName, -1)

	if hasTime == true {
		txt = strings.Replace(txt, "{{.import}}", "import \"time\"", -1)
	} else {
		txt = strings.Replace(txt, "{{.import}}", "", -1)
	}

	txt = strings.Replace(txt, "{{.body}}", formContent, -1)

	err = ioutil.WriteFile("app/vuecmf/form/"+tableName+".go", []byte(txt), 0666)

	if err != nil {
		return false
	}

	return true
}

var makeSer *makeService

// Make 获取make服务实例
func Make() *makeService {
	if makeSer == nil {
		makeSer = &makeService{}
	}
	return makeSer
}
