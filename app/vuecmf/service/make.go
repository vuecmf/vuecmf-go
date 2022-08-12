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

// formRules 表单规则
type formRules struct {
	FieldName string `json:"field_name"` //字段名
	FieldType string `json:"field_type"` //字段类型
	RuleType  string `json:"rule_type"`  //验证规则类型
	RuleValue string `json:"rule_value"` //规则值
	ErrorTips string `json:"error_tips"` //错误提示
}

//Form 功能：生成表单代码文件
//		参数：tableName string 表名（不带表前缀）
func (makeSer *makeService) Form(tableName string) bool {
	var result []formRules

	//查出需要生成模型表的字段相关信息
	db.Table(ns.TableName("model_field")+" VMF").
		Select("VMF.field_name, VMF.type field_type, VMFR.rule_type, VMFR.rule_value, VMFR.error_tips").
		Joins("left join "+ns.TableName("model_form")+" VMF2 on VMF2.model_field_id = VMF.id").
		Joins("left join "+ns.TableName("model_form_rules")+" VMFR on VMF2.id = VMFR.model_form_id").
		Joins("left join "+ns.TableName("model_config")+" MC on VMF.model_id = MC.id").
		Where("VMF.status = 10").
		Where("MC.status = 10").
		Where("MC.table_name = ?", tableName).Find(&result)

	//读取模型模板文件
	tplContent, err := ioutil.ReadFile("app/vuecmf/make/stubs/form.stub")
	if err != nil {
		panic("读取form模板失败")
	}

	formList := make(map[string]map[string][]string)
	ruleMaps := getRuleMaps()

	for _, value := range result {
		if formList[value.FieldName] == nil {
			formList[value.FieldName] = map[string][]string{}
		}

		formList[value.FieldName]["type"] = []string{value.FieldType} //字段类型
		rules := ""
		if ruleMaps[value.RuleType] != "" {
			switch ruleMaps[value.RuleType] {
			case "eq", "gt", "lt", "gte", "lte", "datetime", "len", "max", "min", "required_if", "required_with", "required_without":
				rules = ruleMaps[value.RuleType]
				if value.RuleValue != "" {
					rules += "=" + value.RuleValue
				}
			default:
				rules = ruleMaps[value.RuleType]
			}

			//验证规则
			formList[value.FieldName]["rules"] = append(formList[value.FieldName]["rules"], rules)
			//错误提示语句
			formList[value.FieldName]["tips"] = append(formList[value.FieldName]["tips"], ruleMaps[value.RuleType]+"_tips:\""+value.ErrorTips+"\"")
		}
	}

	formContent := ""
	hasTime := false

	//表单验证字段信息处理
	for fieldName, value := range formList {
		fieldType := "string"
		timeFormat := ""

		switch value["type"][0] {
		case "timestamp", "datetime", "date":
			hasTime = true
			fieldType = "time.Time"
			timeFormat = "time_format:\"2006-01-02 15:04:05\" "
			if value["type"][0] == "date" {
				timeFormat = "time_format:\"2006-01-02\" "
			}
		case "int", "bigint":
			fieldType = "int"
		case "smallint":
			fieldType = "int16"
		case "tinyint":
			fieldType = "int8"
		case "float":
			fieldType = "float32"
		case "double", "decimal":
			fieldType = "float64"
		}

		formContent += helper.UnderToCamel(fieldName) + " " + fieldType + " `json:\"" + fieldName +
			"\" form:\"" + fieldName + "\" " + timeFormat
		if len(value["rules"]) > 0 {
			formContent += "binding:\"" + strings.Join(value["rules"], ",") + "\" " + strings.Join(value["tips"], " ")
		}
		formContent += "`\n\t"
	}

	//获取模型标签名称
	formLabel := ""
	db.Table(ns.TableName("model_config")).Select("label").
		Where("table_name = ?", tableName).Find(&formLabel)

	formName := helper.UnderToCamel(tableName)

	//替换模板文件中内容
	txt := string(tplContent)
	txt = strings.Replace(txt, "{{.comment}}", formLabel, -1)
	txt = strings.Replace(txt, "{{.form_name}}", formName, -1)

	if hasTime == true {
		txt = strings.Replace(txt, "{{.import}}", "import \"time\"", -1)
	} else {
		txt = strings.Replace(txt, "{{.import}}", "", -1)
	}

	txt = strings.Replace(txt, "{{.body}}", formContent, -1)

	err = ioutil.WriteFile("app/vuecmf/form/"+tableName+"_form.go", []byte(txt), 0666)

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

var ruleMaps = make(map[string]string)

// getRuleMaps 获取验证规则映射， 兼容PHP中的验证规则名称
func getRuleMaps() map[string]string {
	ruleMaps["="] = "eq"
	ruleMaps[">"] = "gt"
	ruleMaps["<"] = "lt"
	ruleMaps[">="] = "gte"
	ruleMaps["<="] = "lte"
	ruleMaps["alpha"] = "alpha"
	ruleMaps["alphaNum"] = "alphanum"
	ruleMaps["boolean"] = "boolean"
	ruleMaps["lower"] = "lowercase"
	ruleMaps["upper"] = "uppercase"
	ruleMaps["integer"] = "numeric"
	ruleMaps["number"] = "number"
	ruleMaps["date"] = "datetime"
	ruleMaps["email"] = "email"
	ruleMaps["file"] = "file"
	ruleMaps["ip"] = "ip"
	ruleMaps["macAddr"] = "mac"
	ruleMaps["max"] = "max"
	ruleMaps["min"] = "min"
	ruleMaps["require"] = "required"
	ruleMaps["requireIf"] = "required_if"
	ruleMaps["requireWith"] = "required_with"
	ruleMaps["requireWithout"] = "required_without"
	ruleMaps["unique"] = "unique"
	ruleMaps["url"] = "url"
	ruleMaps["zip"] = "postcode_iso3166_alpha2"
	ruleMaps["len"] = "len"
	//ruleMaps["regex"] = "regex"

	return ruleMaps
}
