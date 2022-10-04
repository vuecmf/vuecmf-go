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
	"os"
	"strconv"
	"strings"
)

// makeService make服务结构
type makeService struct {
	*BaseService
}

// formRow 表单字段信息
type formRow struct {
	FieldName string   //字段名
	Gorm      string   //gorm表信息
	Type      string   //字段类型
	IsSigned  string   //是否为负数
	Rules     []string //验证规则
	ErrTips   []string //错误提示语
}

// formRules 表单字段验证规则
type formRules struct {
	RuleType  string `json:"rule_type"`  //验证规则类型
	RuleValue string `json:"rule_value"` //规则值
	ErrorTips string `json:"error_tips"` //错误提示
}

//Model 功能：生成模型代码文件
//		参数：tableName string 表名（不带表前缀）
func (makeSer *makeService) Model(tableName string, appName string) error {
	if appName == "" {
		appName = "vuecmf"
	}

	var result []model.ModelField

	//查出需要生成模型表的字段相关信息
	Db.Table(NS.TableName("model_field")+" VMF").
		Select("VMF.*").
		Joins("left join "+NS.TableName("model_config")+" MC on VMF.model_id = MC.id").
		Where("VMF.status = 10").
		Where("MC.status = 10").
		Where("MC.table_name = ?", tableName).
		Order("VMF.sort_num").
		Find(&result)

	txt := `package model

{{.import}}

// {{.model_name}} {{.comment}} 模型结构
type {{.model_name}} struct {
	{{.body}}
}

// Data{{.model_name}}Form 提交的表单数据
type Data{{.model_name}}Form struct {
    Data *{{.model_name}} |json:"data" form:"data"|
}`
	modelConf := ModelConfig().GetModelConfig(tableName)
	if modelConf.IsTree == true {
		txt = `package model

{{.import}}

// {{.model_name}} {{.comment}} 模型结构
type {{.model_name}} struct {
	{{.body}}
	Children *{{.model_name}}Tree |json:"children" gorm:"-"|
}

// Data{{.model_name}}Form 提交的表单数据
type Data{{.model_name}}Form struct {
    Data *{{.model_name}} |json:"data" form:"data"|
}


var {{.model_value}}Model *{{.model_name}}

// {{.model_name}}Model 获取{{.model_name}}模型实例
func {{.model_name}}Model() *{{.model_name}} {
	if {{.model_value}}Model == nil {
		{{.model_value}}Model = &{{.model_name}}{}
	}
	return {{.model_value}}Model
}

type {{.model_name}}Tree []*{{.model_name}}

// ToTree 将列表数据转换树形结构
func (m *{{.model_name}}) ToTree(data []*{{.model_name}}) {{.model_name}}Tree {
	treeData := make(map[uint]*{{.model_name}})
	for _, val := range data {
		treeData[val.Id] = val
	}

	var treeList {{.model_name}}Tree

	for _, item := range treeData {
		if item.Pid == 0 {
			treeList = append(treeList, item)
			continue
		}
		if pItem, ok := treeData[item.Pid]; ok {
			if pItem.Children == nil {
				children := {{.model_name}}Tree{item}
				pItem.Children = &children
				continue
			}
			*pItem.Children = append(*pItem.Children, item)
		}
	}

	return treeList

}
`
	}

	txt = strings.Replace(txt, "|", "`", -1)

	var formList []*formRow
	ruleMaps := getRuleMaps()

	dbType := app.DbConf().Connect["default"].Type

	for _, value := range result {
		fr := &formRow{}

		//gorm 处理
		notNull := ""
		defaultVal := ""
		size := ""
		autoCreateTime := ""
		uniqueIndex := ""
		autoIncrement := ""

		if value.IsNull == 20 {
			notNull = "not null;"
		}

		if value.IsAutoIncrement == 10 {
			autoIncrement = "primaryKey;autoIncrement;"
		}

		if value.FieldName == "update_time" || value.FieldName == "last_login_time" || value.DefaultValue == "CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" {
			autoCreateTime = "autoCreateTime;autoUpdateTime;"
		} else if value.DefaultValue == "CURRENT_TIMESTAMP" {
			autoCreateTime = "autoCreateTime;"
		} else {
			//字段默认值
			switch {
			case value.IsAutoIncrement != 10 && value.DefaultValue != "":
				defaultVal = "default:" + value.DefaultValue + ";"
			case value.IsAutoIncrement != 10 && value.DefaultValue == "":
				defaultVal = "default:'';"
			}

			//针对MYSQL整型类型字段长度处理
			if strings.ToLower(dbType) == "mysql" && (value.Type == "int" || value.Type == "bigint" || value.Type == "smallint" || value.Type == "tinyint") {
				switch {
				case value.Length <= 4:
					value.Length = 8
				case value.Length <= 6:
					value.Length = 16
				case value.Length <= 9:
					value.Length = 24
				case value.Length <= 11:
					value.Length = 32
				default:
					value.Length = 64
				}
			}

			size = "size:" + strconv.Itoa(int(value.Length)) + ";"
		}

		//字段唯一索引处理
		modelIndexId := 0
		id := strconv.Itoa(int(value.Id))
		Db.Table(NS.TableName("model_index")).Select("id").
			Where("model_field_id = ? or model_field_id like ? or model_field_id like ?", id, id+",%", "%,"+id).
			Find(&modelIndexId)

		if modelIndexId > 0 {
			uniqueIndex = "uniqueIndex:unique_index;"
		}
		gormCnf := " gorm:\"column:" + value.FieldName + ";" + autoIncrement + size + uniqueIndex + notNull + autoCreateTime + defaultVal +
			"comment:" + value.Note + "\""

		fr.FieldName = value.FieldName                  //字段名
		fr.Gorm = gormCnf                               //gorm表信息
		fr.Type = value.Type                            //字段类型
		fr.IsSigned = strconv.Itoa(int(value.IsSigned)) //是否为负数

		rules := ""
		var formRulesList []formRules

		Db.Table(NS.TableName("model_form")+" VMF").
			Select("VMFR.rule_type, VMFR.rule_value, VMFR.error_tips").
			Joins("left join "+NS.TableName("model_form_rules")+" VMFR on VMF.id = VMFR.model_form_id").
			Where("VMF.status = 10").
			Where("VMFR.status = 10").
			Where("VMF.model_field_id = ?", value.Id).
			Find(&formRulesList)

		//数据验证规则
		for _, rule := range formRulesList {
			if ruleMaps[rule.RuleType] != "" {
				switch ruleMaps[rule.RuleType] {
				case "eq", "gt", "lt", "gte", "lte", "datetime", "len", "max", "min", "required_if", "required_with", "required_without":
					rules = ruleMaps[rule.RuleType]
					if rule.RuleValue != "" {
						rules += "=" + rule.RuleValue
					}
				default:
					rules = ruleMaps[rule.RuleType]
				}

				fr.Rules = append(fr.Rules, rules) //验证规则
				//错误提示语句
				fr.ErrTips = append(fr.ErrTips, ruleMaps[rule.RuleType]+"_tips:\""+rule.ErrorTips+"\"")

			}
		}

		formList = append(formList, fr)
	}

	modelContent := ""
	hasTime := false

	//模型字段信息处理
	for _, row := range formList {
		fieldType := "string"
		timeFormat := ""

		switch row.Type {
		case "timestamp", "datetime", "date":
			hasTime = true
			fieldType = "time.Time"
			timeFormat = "time_format:\"2006-01-02 15:04:05\" "
			if row.Type == "date" {
				timeFormat = "time_format:\"2006-01-02\" "
			}
		case "int", "bigint":
			fieldType = "int"
			if row.IsSigned == "20" {
				fieldType = "uint"
			}
		case "smallint":
			fieldType = "int16"
			if row.IsSigned == "20" {
				fieldType = "uint16"
			}
		case "tinyint":
			fieldType = "int8"
			if row.IsSigned == "20" {
				fieldType = "uint8"
			}
		case "float":
			fieldType = "float32"
		case "double", "decimal":
			fieldType = "float64"
		}

		modelContent += helper.UnderToCamel(row.FieldName) + " " + fieldType + " `json:\"" + row.FieldName +
			"\" form:\"" + row.FieldName + "\" " + timeFormat
		if len(row.Rules) > 0 {
			modelContent += "binding:\"" + strings.Join(row.Rules, ",") + "\" " + strings.Join(row.ErrTips, " ")
		}
		modelContent += row.Gorm + "`\n\t"
	}

	//获取模型标签名称
	modelLabel := ""
	Db.Table(NS.TableName("model_config")).Select("label").
		Where("table_name = ?", tableName).Find(&modelLabel)

	modelName := helper.UnderToCamel(tableName)
	modelValue := strings.ToLower(modelName)

	//替换模板文件中内容
	txt = strings.Replace(txt, "{{.app_name}}", appName, -1)
	txt = strings.Replace(txt, "{{.comment}}", modelLabel, -1)
	txt = strings.Replace(txt, "{{.model_name}}", modelName, -1)
	txt = strings.Replace(txt, "{{.model_value}}", modelValue, -1)

	if hasTime == true {
		txt = strings.Replace(txt, "{{.import}}", "import \"time\"", -1)
	} else {
		txt = strings.Replace(txt, "{{.import}}", "", -1)
	}

	txt = strings.Replace(txt, "{{.body}}", modelContent, -1)
	err := os.WriteFile("app/"+appName+"/model/"+tableName+".go", []byte(txt), 0666)
	return err
}

//Service 功能：生成服务代码文件
//		  参数：tableName string 表名（不带表前缀）
func (makeSer *makeService) Service(tableName string, appName string) error {
	if appName == "" {
		appName = "vuecmf"
	}

	serviceMethod := helper.UnderToCamel(tableName)
	nameArr := []rune(serviceMethod)
	nameArr[0] += 32
	serviceName := string(nameArr)

	txt := `package service

{{.import_base}}


// {{.service_name}}Service {{.service_name}}服务结构
type {{.service_name}}Service struct {
	{{.extend_base}}
}

var {{.service_name}} *{{.service_name}}Service

// {{.service_method}} 获取{{.service_name}}服务实例
func {{.service_method}}() *{{.service_name}}Service {
	if {{.service_name}} == nil {
		{{.service_name}} = &{{.service_name}}Service{}
	}
	return {{.service_name}}
}`
	modelConf := ModelConfig().GetModelConfig(tableName)
	if modelConf.IsTree == true {
		txt = `package service

import (
	"github.com/vuecmf/vuecmf-go/app/vuecmf/helper"
	"github.com/vuecmf/vuecmf-go/app/{{.app_name}}/model"
	{{.import_base2}}
)

// {{.service_name}}Service {{.service_name}}服务结构
type {{.service_name}}Service struct {
	{{.extend_base}}
	TableName string
}

var {{.service_name}} *{{.service_name}}Service

// {{.service_method}} 获取{{.service_name}}服务实例
func {{.service_method}}() *{{.service_name}}Service {
	if {{.service_name}} == nil {
		{{.service_name}} = &{{.service_name}}Service{TableName: "{{.service_name}}"}
	}
	return {{.service_name}}
}

// List 获取列表数据
// 		参数：params 查询参数
func (ser *{{.service_name}}Service) List(params *helper.DataListParams) (interface{} , error) {
	if params.Data.Action == "getField" {
		//拉取列表的字段信息
		return ser.getFieldList(ser.TableName, params.Data.Filter)
	}else{
		//拉取列表的数据
		var {{.service_name}}List []*model.{{.service_method}}
		var res = make(map[string]interface{})

		ser.getList(&{{.service_name}}List, ser.TableName, params)

		//转换成树形列表
		tree := model.{{.service_method}}Model().ToTree({{.service_name}}List)
		res["data"] = tree
		return res, nil
	}
}
`
	}

	importBase := ""
	importBase2 := ""
	extendBase := "*BaseService"
	if appName != "vuecmf" {
		importBase = "import \"github.com/vuecmf/vuecmf-go/app/vuecmf/service\""
		importBase2 = "\"github.com/vuecmf/vuecmf-go/app/vuecmf/service\""
		extendBase = "*service.BaseService"
	}

	txt = strings.Replace(txt, "{{.app_name}}", appName, -1)
	txt = strings.Replace(txt, "{{.service_name}}", serviceName, -1)
	txt = strings.Replace(txt, "{{.service_method}}", serviceMethod, -1)
	txt = strings.Replace(txt, "{{.import_base}}", importBase, -1)
	txt = strings.Replace(txt, "{{.import_base2}}", importBase2, -1)
	txt = strings.Replace(txt, "{{.extend_base}}", extendBase, -1)

	err := os.WriteFile("app/"+appName+"/service/"+tableName+".go", []byte(txt), 0666)
	return err
}

//Controller 功能：生成控制器代码文件
//		  参数：tableName string 表名（不带表前缀）
func (makeSer *makeService) Controller(tableName string, appName string) error {
	if appName == "" {
		appName = "vuecmf"
	}

	controllerName := helper.UnderToCamel(tableName)
	ctrlValName := helper.ToFirstLower(controllerName)

	//查询模型是否有需要模糊查询的字段
	filterFields := ModelField().getFilterFields(tableName)
	filterFieldStr := "\"" + strings.Join(filterFields, "\",\"") + "\""

	moduleName := app.AppConfig().Module

	txt := `package controller

import (
	"github.com/vuecmf/vuecmf-go/app/route"
	"{{.module_name}}/app/{{.app_name}}/model"
	{{.import_base}}
)

type {{.controller_name}} struct {
    {{.extend_base}}
}

func init() {
	{{.controller_var_name}} := &{{.controller_name}}{}
    {{.controller_var_name}}.TableName = "{{.table_name}}"
    {{.controller_var_name}}.Model = &model.{{.controller_name}}{}
    {{.controller_var_name}}.ListData = &[]model.{{.controller_name}}{}
    {{.controller_var_name}}.SaveForm = &model.Data{{.controller_name}}Form{}
    {{.controller_var_name}}.FilterFields = []string{{{.filter_fields}}}

    route.Register({{.controller_var_name}}, "POST", "{{.app_name}}")
}
`
	modelConf := ModelConfig().GetModelConfig(tableName)
	if modelConf.IsTree == true {
		txt = `package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/vuecmf/vuecmf-go/app/route"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/helper"
	"{{.module_name}}/app/{{.app_name}}/model"
	"{{.module_name}}/app/{{.app_name}}/service"
	{{.import_base}}
)

type {{.controller_name}} struct {
    {{.extend_base}}
}

func init() {
	{{.controller_var_name}} := &{{.controller_name}}{}
    {{.controller_var_name}}.TableName = "{{.controller_var_name}}"
    {{.controller_var_name}}.Model = &model.{{.controller_name}}{}
    {{.controller_var_name}}.ListData = &[]model.{{.controller_name}}{}
    {{.controller_var_name}}.SaveForm = &model.Data{{.controller_name}}Form{}
    {{.controller_var_name}}.FilterFields = []string{{{.filter_fields}}}

    route.Register({{.controller_var_name}}, "POST", "{{.app_name}}")
}

// Index 列表页
func (ctrl *{{.controller_name}}) Index(c *gin.Context) {
    listParams := &helper.DataListParams{}
	common(c, listParams, func() (interface{}, error) {
        return service.{{.controller_name}}().List(listParams)
	})
}

`
	}

	importBase := ""
	extendBase := "Base"
	if appName != "vuecmf" {
		importBase = "\"github.com/vuecmf/vuecmf-go/app/vuecmf/controller\""
		extendBase = "controller.Base"
	}

	txt = strings.Replace(txt, "{{.module_name}}", moduleName, -1)
	txt = strings.Replace(txt, "{{.app_name}}", appName, -1)
	txt = strings.Replace(txt, "{{.controller_name}}", controllerName, -1)
	txt = strings.Replace(txt, "{{.controller_var_name}}", ctrlValName, -1)
	txt = strings.Replace(txt, "{{.table_name}}", tableName, -1)
	txt = strings.Replace(txt, "{{.filter_fields}}", filterFieldStr, -1)
	txt = strings.Replace(txt, "{{.import_base}}", importBase, -1)
	txt = strings.Replace(txt, "{{.extend_base}}", extendBase, -1)

	err := os.WriteFile("app/"+appName+"/controller/"+tableName+".go", []byte(txt), 0666)
	return err
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
