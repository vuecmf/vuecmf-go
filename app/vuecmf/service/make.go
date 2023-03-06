//+----------------------------------------------------------------------
// | Copyright (c) 2023 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: vuecmf <tulihua2004@126.com>
// +----------------------------------------------------------------------

package service

import (
	"encoding/json"
	"errors"
	"github.com/vuecmf/vuecmf-go/app"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/helper"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/model"
	"gorm.io/gorm"
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

//Model 生成模型代码文件
//	参数：
//		tableName string 表名（不带表前缀）
//		appName string 应用名称
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
    Data *{{.model_name}} 'json:"data" form:"data" binding:"required" required_tips:"参数data不能为空"'
}`
	modelConf := ModelConfig().GetModelConfig(tableName)
	if modelConf.IsTree == true {
		txt = `package model

{{.import}}

// {{.model_name}} {{.comment}} 模型结构
type {{.model_name}} struct {
	{{.body}}
	Children *{{.model_name}}Tree 'json:"children" gorm:"-"'
}

// Data{{.model_name}}Form 提交的表单数据
type Data{{.model_name}}Form struct {
    Data *{{.model_name}} 'json:"data" form:"data" binding:"required" required_tips:"参数data不能为空"'
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
	idList := make([]uint, 0, len(data))
	for _, val := range data {
		treeData[val.Id] = val
		idList = append(idList, val.Id)
	}

	var treeList {{.model_name}}Tree

	for _, id := range idList {
		if treeData[id].Pid == 0 || treeData[treeData[id].Pid] == nil {
			treeList = append(treeList, treeData[id])
			continue
		}
		if pItem, ok := treeData[treeData[id].Pid]; ok {
			if pItem.Children == nil {
				children := {{.model_name}}Tree{treeData[id]}
				pItem.Children = &children
				continue
			}
			*pItem.Children = append(*pItem.Children, treeData[id])
		}
	}

	return treeList

}
`
	}

	txt = strings.Replace(txt, "'", "`", -1)

	var formList []*formRow
	ruleMaps := getRuleMaps()

	dbType := app.DbConf().Connect[Conf.Env].Type

	for _, value := range result {
		fr := &formRow{}

		//gorm 处理
		notNull := ""
		defaultVal := ""
		size := ""
		autoCreateTime := ""
		uniqueIndex := ""
		autoIncrement := ""
		columnType := ""

		if value.IsNull == 20 {
			notNull = "not null;"
		}

		if value.IsAutoIncrement == 10 {
			autoIncrement = "primaryKey;autoIncrement;"
		}

		if value.Type == "timestamp" {
			columnType = "type:timestamp;"
		}

		if value.FieldName == "update_time" || value.FieldName == "last_login_time" || value.DefaultValue == "CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" {
			autoCreateTime = "autoCreateTime;autoUpdateTime;"
			defaultVal = "default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;"
		} else if value.DefaultValue == "CURRENT_TIMESTAMP" {
			autoCreateTime = "autoCreateTime;"
			defaultVal = "default:CURRENT_TIMESTAMP;"
		} else {
			//字段默认值
			switch {
			case value.IsAutoIncrement != 10 && (value.Type == "varchar" || value.Type == "char" || value.Type == "text" || value.Type == "mediumtext" || value.Type == "longtext"):
				defaultVal = "default:'" + value.DefaultValue + "';"
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
			Where("index_type = 'UNIQUE'").
			Find(&modelIndexId)

		if modelIndexId > 0 {
			uniqueIndex = "uniqueIndex:unique_index;"
		}
		gormCnf := " gorm:\"" + columnType + "column:" + value.FieldName + ";" + autoIncrement + size + uniqueIndex + notNull + autoCreateTime + defaultVal +
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

		switch row.Type {
		case "timestamp", "datetime":
			hasTime = true
			fieldType = "model.JSONTime"
		case "date":
			hasTime = true
			fieldType = "model.JSONDate"
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
			"\" form:\"" + row.FieldName + "\" "
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
		txt = strings.Replace(txt, "{{.import}}", "import \"github.com/vuecmf/vuecmf-go/app/vuecmf/model\"", -1)
	} else {
		txt = strings.Replace(txt, "{{.import}}", "", -1)
	}

	txt = strings.Replace(txt, "{{.body}}", modelContent, -1)
	err := os.WriteFile("app/"+appName+"/model/"+tableName+".go", []byte(txt), 0666)
	return err
}

//Service 生成服务代码文件
//	参数：
//		tableName string 表名（不带表前缀）
//		appName string 应用名称
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
	moduleName := ""
	modelConf := ModelConfig().GetModelConfig(tableName)
	if modelConf.IsTree == true {
		moduleName = app.Config().Module
		txt = `package service

import (
	"github.com/vuecmf/vuecmf-go/app/vuecmf/helper"
	"{{.module_name}}/app/{{.app_name}}/model"
	{{.import_base2}}
	"strconv"
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
		{{.service_name}} = &{{.service_name}}Service{TableName: "{{.table_name}}"}
	}
	return {{.service_name}}
}

// GetIdPath 获取父级ID的ID路径
func (ser *{{.service_name}}Service) GetIdPath(pid uint) string {
	var pidIdPath string
	service.Db.Table(service.NS.TableName("{{.table_name}}")).Select("id_path").Where("id = ?", pid).Find(&pidIdPath)
	if pid > 0 {
		if pidIdPath == "" {
			pidIdPath = strconv.Itoa(int(pid))
		} else {
			pidIdPath += "," + strconv.Itoa(int(pid))
		}
	}
	return pidIdPath
}

// Create 创建单条或多条数据, 成功返回影响行数
func (ser *{{.service_name}}Service) Create(data *model.{{.service_method}}) (int64, error) {
	data.IdPath = ser.GetIdPath(data.Pid)
	res := service.Db.Create(data)
	return res.RowsAffected, res.Error
}

// Update 更新数据, 成功返回影响行数
func (ser *{{.service_name}}Service) Update(data *model.{{.service_method}}) (int64, error) {
	data.IdPath = ser.GetIdPath(data.Pid)
	res := service.Db.Updates(data)
	return res.RowsAffected, res.Error
}

// List 获取列表数据
// 		参数：params 查询参数
func (ser *{{.service_name}}Service) List(params *helper.DataListParams) (interface{} , error) {
	if params.Data.Action == "getField" {
		//拉取列表的字段信息
		return ser.GetFieldList(ser.TableName, params.Data.Filter, 10)
	}else{
		//拉取列表的数据
		var {{.service_name}}List []*model.{{.service_method}}
		var res = make(map[string]interface{})

		ser.GetList(&{{.service_name}}List, ser.TableName, params)

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

	txt = strings.Replace(txt, "{{.module_name}}", moduleName, -1)
	txt = strings.Replace(txt, "{{.app_name}}", appName, -1)
	txt = strings.Replace(txt, "{{.table_name}}", tableName, -1)
	txt = strings.Replace(txt, "{{.service_name}}", serviceName, -1)
	txt = strings.Replace(txt, "{{.service_method}}", serviceMethod, -1)
	txt = strings.Replace(txt, "{{.import_base}}", importBase, -1)
	txt = strings.Replace(txt, "{{.import_base2}}", importBase2, -1)
	txt = strings.Replace(txt, "{{.extend_base}}", extendBase, -1)

	err := os.WriteFile("app/"+appName+"/service/"+tableName+".go", []byte(txt), 0666)
	return err
}

//Controller 生成控制器代码文件
//	参数：
//		tableName string 表名（不带表前缀）
//		appName string 应用名称
func (makeSer *makeService) Controller(tableName string, appName string) error {
	if appName == "" {
		appName = "vuecmf"
	}

	controllerName := helper.UnderToCamel(tableName)
	ctrlValName := helper.ToFirstLower(controllerName)

	//查询模型是否有需要模糊查询的字段
	filterFields := ModelField().getFilterFields(tableName)
	filterFieldStr := "\"" + strings.Join(filterFields, "\",\"") + "\""

	moduleName := app.Config().Module

	txt := `package controller

import (
	"github.com/gin-gonic/gin"
	"{{.module_name}}/app/{{.app_name}}/model"
	"{{.module_name}}/app/{{.app_name}}/service"
	"github.com/vuecmf/vuecmf-go/app/route"
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
    {{.controller_var_name}}.FilterFields = []string{{{.filter_fields}}}

    route.Register({{.controller_var_name}}, "POST", "{{.app_name}}")
}

// Save 新增/更新 单条数据
func (ctrl *{{.controller_name}}) Save(c *gin.Context) {
	saveForm := &model.Data{{.controller_name}}Form{}
	controller.Common(c, saveForm, func() (interface{}, error) {
		if saveForm.Data.Id == uint(0) {
			return service.{{.controller_name}}().Create(saveForm.Data)
		} else {
			return service.{{.controller_name}}().Update(saveForm.Data)
		}
	})
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
    {{.controller_var_name}}.TableName = "{{.table_name}}"
    {{.controller_var_name}}.Model = &model.{{.controller_name}}{}
    {{.controller_var_name}}.ListData = &[]model.{{.controller_name}}{}
    {{.controller_var_name}}.FilterFields = []string{{{.filter_fields}}}

    route.Register({{.controller_var_name}}, "POST", "{{.app_name}}")
}

// Index 列表页
func (ctrl *{{.controller_name}}) Index(c *gin.Context) {
    listParams := &helper.DataListParams{}
	controller.Common(c, listParams, func() (interface{}, error) {
        return service.{{.controller_name}}().List(listParams)
	})
}

// Save 新增/更新 单条数据
func (ctrl *{{.controller_name}}) Save(c *gin.Context) {
	saveForm := &model.Data{{.controller_name}}Form{}
	controller.Common(c, saveForm, func() (interface{}, error) {
		if saveForm.Data.Id == uint(0) {
			return service.{{.controller_name}}().Create(saveForm.Data)
		} else {
			return service.{{.controller_name}}().Update(saveForm.Data)
		}
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

	if err := makeSer.UpdateRunFile(); err != nil {
		return err
	}

	return os.WriteFile("app/"+appName+"/controller/"+tableName+".go", []byte(txt), 0666)
}

//RemoveModel 删除模型代码文件
//	参数：
//		tableName string 表名（不带表前缀）
//		appName string 应用名称
func (makeSer *makeService) RemoveModel(tableName string, appName string) error {
	pathName := "app/" + appName + "/model/" + tableName + ".go"
	//文件不存在的就直接返回
	if _, err := os.Stat(pathName); err != nil {
		return nil
	}
	return os.Remove(pathName)
}

// RemoveService 删除服务代码文件
//	参数：
//		tableName string 表名（不带表前缀）
//		appName string 应用名称
func (makeSer *makeService) RemoveService(tableName string, appName string) error {
	pathName := "app/" + appName + "/service/" + tableName + ".go"
	//文件不存在的就直接返回
	if _, err := os.Stat(pathName); err != nil {
		return nil
	}
	return os.Remove(pathName)
}

// RemoveController 删除控制器代码文件
//	参数：
//		tableName string 表名（不带表前缀）
//		appName string 应用名称
func (makeSer *makeService) RemoveController(tableName string, appName string) error {
	if err := makeSer.UpdateRunFile(); err != nil {
		return err
	}

	pathName := "app/" + appName + "/controller/" + tableName + ".go"
	//文件不存在的就直接返回
	if _, err := os.Stat(pathName); err != nil {
		return nil
	}

	return os.Remove(pathName)
}

// UpdateRunFile 更新启动文件
func (makeSer *makeService) UpdateRunFile() error {
	if app.Config().UpdateMain == false {
		return nil
	}

	projectName := app.Config().Module
	mainFile := `package main

import (
	"github.com/gin-gonic/gin"
	"github.com/vuecmf/vuecmf-go/app"
	"github.com/vuecmf/vuecmf-go/app/route"
	_ "github.com/vuecmf/vuecmf-go/app/vuecmf/controller"
	"log"
${import_package}
)

func main() {
	engine := gin.Default()

	//初始化路由
	route.InitRoute(engine)

	err := engine.Run(":" + app.Config().ServerPort)
	if err != nil {
		log.Fatal("服务启动失败！", err)
	}

}
`
	var appList []string
	Db.Table(NS.TableName("app_config") + " A").Select("app_name").
		Joins("left join " + NS.TableName("model_config") + " MC on A.id = MC.app_id").
		Where("A.status = 10").
		Where("A.app_name != 'vuecmf'").
		Where("MC.status = 10").
		Group("app_name").Find(&appList)

	importPackage := ""
	for _, appName := range appList {
		importPackage = importPackage + "\t_ \"" + projectName + "/app/" + appName + "/controller\"\n"
	}

	mainFile = strings.Replace(mainFile, "${import_package}", importPackage, -1)
	if err := os.WriteFile("main.go", []byte(mainFile), 0666); err != nil {
		return errors.New("更新启动文件main.go失败！" + err.Error())
	}
	return nil
}

// RemoveAll 删除表名相关所有控制器、模型及服务
//	参数：
//		tableName 表名
func (makeSer *makeService) RemoveAll(tableName string) error {
	var err error
	//先根据tableName查出所有相关的模型、服务及控制器，然后全部删除
	appList := AppConfig().GetAppListByTableName(tableName)

	if len(appList) > 0 {
		for _, appName := range appList {
			if err = makeSer.RemoveController(tableName, appName); err != nil {
				return err
			}
			if err = makeSer.RemoveModel(tableName, appName); err != nil {
				return err
			}
			if err = makeSer.RemoveService(tableName, appName); err != nil {
				return err
			}
		}
	}
	return nil
}

// MakeAll 根据表名生成相关的所有控制器、模型及服务
//	参数：
//		tableName 表名
func (makeSer *makeService) MakeAll(tableName string) error {
	var err error
	//先根据tableName查出所有相关的模型、服务及控制器，然后生成
	appList := AppConfig().GetAppListByTableName(tableName)

	if len(appList) > 0 {
		for _, appName := range appList {
			if err = makeSer.Controller(tableName, appName); err != nil {
				break
			}
			if err = makeSer.Model(tableName, appName); err != nil {
				break
			}
			if err = makeSer.Service(tableName, appName); err != nil {
				break
			}
		}
	}
	return nil
}

// MakeAppModel 根据应用ID及模型ID生成对应代码文件
//	参数：
//		appId 应用ID
//		tableName 表名
func (makeSer *makeService) MakeAppModel(appId uint, tableName string) error {
	appName := AppConfig().GetAppNameById(appId)
	if appName == "" {
		return errors.New("没有找到应用名称")
	}

	if err := makeSer.Controller(tableName, appName); err != nil {
		return err
	}
	if err := makeSer.Model(tableName, appName); err != nil {
		return err
	}
	if err := makeSer.Service(tableName, appName); err != nil {
		return err
	}
	return nil
}

// RemoveAppModel 根据应用ID及模型ID删除对应代码文件
//	参数：
//		appId 应用ID
//		modelId 模型ID
func (makeSer *makeService) RemoveAppModel(appId, modelId uint) error {
	var appName string
	Db.Table(NS.TableName("app_config")).Select("app_name").
		Where("id = ?", appId).
		Where("status = 10").Find(&appName)

	if appName == "vuecmf" {
		return nil
	}

	if appName == "" {
		return errors.New("没有找到应用名称")
	}

	var tableName string
	Db.Table(NS.TableName("model_config")).Select("table_name").
		Where("id = ?", modelId).
		Where("status = 10").Find(&tableName)
	if tableName == "" {
		return errors.New("没有找到模型(" + strconv.Itoa(int(modelId)) + ")对应的表名")
	}

	//更新菜单中使用的模型
	Db.Table(NS.TableName("menu")).
		Where("app_id = ?", appId).
		Where("model_id = ?", modelId).
		Update("model_id", 0)

	if err := makeSer.RemoveController(tableName, appName); err != nil {
		return err
	}
	if err := makeSer.RemoveModel(tableName, appName); err != nil {
		return err
	}
	if err := makeSer.RemoveService(tableName, appName); err != nil {
		return err
	}
	return nil
}

// BuildModel 生成模型相关数据
//	参数：
//		mc 模型配置实例
func (makeSer *makeService) BuildModel(mc *model.ModelConfig) error {
	var baseTable interface{}
	var insertDataJson string
	if mc.IsTree == 10 {
		type BaseTable struct {
			Id      uint   `json:"id" form:"id"  gorm:"column:id;primaryKey;autoIncrement;size:32;not null;comment:自增ID"`
			Title   string `json:"title" form:"title" binding:"required" required_tips:"标题必填" gorm:"column:title;size:64;not null;default:'';comment:标题"`
			Pid     uint   `json:"pid" form:"pid"  gorm:"column:pid;size:32;not null;default:0;comment:父级ID"`
			IdPath  string `json:"id_path" form:"id_path"  gorm:"column:id_path;size:255;not null;default:'';comment:ID层级路径"`
			SortNum uint   `json:"sort_num" form:"sort_num"  gorm:"column:sort_num;size:32;not null;default:0;comment:菜单的排列顺序(小在前)"`
			Status  uint16 `json:"status" form:"status"  gorm:"column:status;size:8;not null;default:10;comment:状态：10=开启，20=禁用"`
		}

		baseTable = &BaseTable{}

		//写入数据
		insertDataJson = `[
    {
        "field_name": "id",
        "label": "ID",
        "model_id": {$model_id},
        "type": "int",
        "length": 11,
        "decimal_length": 0,
        "is_null": 20,
        "note": "自增ID",
        "default_value": "0",
        "is_auto_increment": 10,
        "is_label": 20,
        "is_signed": 20,
        "is_show": 10,
        "is_fixed": 20,
        "column_width": 100,
        "is_filter": 20,
        "sort_num": 0,
        "status": 10
    },
    {
        "field_name": "title",
        "label": "标题",
        "model_id": {$model_id},
        "type": "varchar",
        "length": 64,
        "decimal_length": 0,
        "is_null": 20,
        "note": "标题",
        "default_value": "",
        "is_auto_increment": 20,
        "is_label": 10,
        "is_signed": 20,
        "is_show": 10,
        "is_fixed": 20,
        "column_width": 0,
        "is_filter": 10,
        "sort_num": 60,
        "status": 10
    },
    {
        "field_name": "pid",
        "label": "父级",
        "model_id": {$model_id},
        "type": "int",
        "length": 11,
        "decimal_length": 0,
        "is_null": 20,
        "note": "父级ID",
        "default_value": "0",
        "is_auto_increment": 20,
        "is_label": 20,
        "is_signed": 20,
        "is_show": 10,
        "is_fixed": 20,
        "column_width": 150,
        "is_filter": 10,
        "sort_num": 9996,
        "status": 10
    },
    {
        "field_name": "id_path",
        "label": "层级路径",
        "model_id": {$model_id},
        "type": "varchar",
        "length": 255,
        "decimal_length": 0,
        "is_null": 20,
        "note": "ID层级路径",
        "default_value": "",
        "is_auto_increment": 20,
        "is_label": 20,
        "is_signed": 20,
        "is_show": 10,
        "is_fixed": 20,
        "column_width": 150,
        "is_filter": 10,
        "sort_num": 9997,
        "status": 10
    },
    {
        "field_name": "sort_num",
        "label": "排序",
        "model_id": {$model_id},
        "type": "int",
        "length": 11,
        "decimal_length": 0,
        "is_null": 20,
        "note": "排列顺序(小在前)",
        "default_value": "0",
        "is_auto_increment": 20,
        "is_label": 20,
        "is_signed": 20,
        "is_show": 10,
        "is_fixed": 20,
        "column_width": 100,
        "is_filter": 10,
        "sort_num": 9998,
        "status": 10
    },
    {
        "field_name": "status",
        "label": "状态",
        "model_id": {$model_id},
        "type": "int",
        "length": 4,
        "decimal_length": 0,
        "is_null": 20,
        "note": "状态：10=开启，20=禁用",
        "default_value": "10",
        "is_auto_increment": 20,
        "is_label": 20,
        "is_signed": 20,
        "is_show": 10,
        "is_fixed": 20,
        "column_width": 100,
        "is_filter": 10,
        "sort_num": 9999,
        "status": 10
    }
]`

	} else {
		type BaseTable struct {
			Id     uint   `json:"id" form:"id"  gorm:"column:id;primaryKey;autoIncrement;size:32;not null;comment:自增ID"`
			Status uint16 `json:"status" form:"status"  gorm:"column:status;size:8;not null;default:10;comment:状态：10=开启，20=禁用"`
		}
		baseTable = &BaseTable{}

		insertDataJson = `[
    {
        "field_name": "id",
        "label": "ID",
        "model_id": {$model_id},
        "type": "int",
        "length": 11,
        "decimal_length": 0,
        "is_null": 20,
        "note": "自增ID",
        "default_value": "0",
        "is_auto_increment": 10,
        "is_label": 20,
        "is_signed": 20,
        "is_show": 10,
        "is_fixed": 20,
        "column_width": 100,
        "is_filter": 20,
        "sort_num": 0,
        "status": 10
    },
    {
        "field_name": "status",
        "label": "状态",
        "model_id": {$model_id},
        "type": "int",
        "length": 4,
        "decimal_length": 0,
        "is_null": 20,
        "note": "状态：10=开启，20=禁用",
        "default_value": "10",
        "is_auto_increment": 20,
        "is_label": 20,
        "is_signed": 20,
        "is_show": 10,
        "is_fixed": 20,
        "column_width": 100,
        "is_filter": 10,
        "sort_num": 9999,
        "status": 10
    }
]`

	}

	return Db.Transaction(func(tx *gorm.DB) error {
		//创建模型配置数据
		if err := tx.Create(mc).Error; err != nil {
			return err
		}

		//创建表
		if err := tx.Set("gorm:table_options", "ENGINE=InnoDB COLLATE=utf8mb4_unicode_ci COMMENT='"+mc.Remark+"'").AutoMigrate(&baseTable); err != nil {
			return errors.New("创建基础表" + NS.TableName(mc.TableName) + "失败:" + err.Error())
		}
		//将表重命名为需要创建的表名称
		if err := tx.Migrator().RenameTable(NS.TableName("base_table"), NS.TableName(mc.TableName)); err != nil {
			return errors.New("创建表" + NS.TableName(mc.TableName) + "失败:" + err.Error())
		}

		//写入数据
		insertDataJson = strings.Replace(insertDataJson, "{$model_id}", strconv.Itoa(int(mc.Id)), -1)
		var insertData []model.ModelField
		if err := json.Unmarshal([]byte(insertDataJson), &insertData); err != nil {
			return err
		}
		if err := tx.Create(insertData).Error; err != nil {
			return err
		}

		//添加字段选项
		insertDataJson = `[
    {
        "model_id": {$model_id},
        "model_field_id": {$model_field_id},
        "option_value": "10",
        "option_label": "开启"
    },
    {
        "model_id": {$model_id},
        "model_field_id": {$model_field_id},
        "option_value": "20",
        "option_label": "禁用"
    }
]`
		var modelFieldId string
		tx.Table(NS.TableName("model_field")).Select("id").
			Where("model_id = ?", mc.Id).
			Where("field_name = 'status'").
			Where("status = 10").Find(&modelFieldId)

		insertDataJson = strings.Replace(insertDataJson, "{$model_id}", strconv.Itoa(int(mc.Id)), -1)
		insertDataJson = strings.Replace(insertDataJson, "{$model_field_id}", modelFieldId, -1)
		var fieldOptionData []model.FieldOption
		if err := json.Unmarshal([]byte(insertDataJson), &fieldOptionData); err != nil {
			return err
		}
		if err := tx.Create(fieldOptionData).Error; err != nil {
			return err
		}

		//添加动作信息
		insertDataJson = `[
    {
        "label": "{$label}管理列表",
        "api_path": "/{$app_name}/{$table_name}",
		"model_id": {$model_id},
        "action_type": "list"
    },
    {
        "label": "保存{$label}",
        "api_path": "/{$app_name}/{$table_name}/save",
		"model_id": {$model_id},
        "action_type": "save"
    },
	{
        "label": "删除{$label}",
        "api_path": "/{$app_name}/{$table_name}/delete",
		"model_id": {$model_id},
        "action_type": "delete"
    },
	{
        "label": "{$label}下拉列表",
        "api_path": "/{$app_name}/{$table_name}/dropdown",
		"model_id": {$model_id},
        "action_type": "dropdown"
    },
	{
        "label": "批量保存{$label}",
        "api_path": "/{$app_name}/{$table_name}/save_all",
		"model_id": {$model_id},
        "action_type": "save_all"
    }
]`
		appName := AppConfig().GetAppNameById(mc.AppId)

		insertDataJson = strings.Replace(insertDataJson, "{$label}", mc.Label, -1)
		insertDataJson = strings.Replace(insertDataJson, "{$app_name}", appName, -1)
		insertDataJson = strings.Replace(insertDataJson, "{$table_name}", mc.TableName, -1)
		insertDataJson = strings.Replace(insertDataJson, "{$model_id}", strconv.Itoa(int(mc.Id)), -1)
		var modelActionData []model.ModelAction
		if err := json.Unmarshal([]byte(insertDataJson), &modelActionData); err != nil {
			return err
		}
		if err := tx.Create(modelActionData).Error; err != nil {
			return err
		}

		//设置模型的默认动作
		listActionId := ModelAction().GetListActionIdByModelId(mc.Id)
		if err := tx.Table(NS.TableName("model_config")).
			Where("id = ?", mc.Id).
			Update("default_action_id", listActionId).Error; err != nil {
			return err
		}

		//生成代码文件
		return Make().MakeAppModel(mc.AppId, mc.TableName)
	})
}

// RemoveModelData 删除模型相关的所有数据
//	参数：
//		mc 模型配置实例
func (makeSer *makeService) RemoveModelData(mc *model.ModelConfig) error {
	//根据动作表找到对应权限项，清除rules表相关信息
	var actionList []string
	Db.Table(NS.TableName("model_action")).Select("api_path").
		Where("model_id = ?", mc.Id).
		Where("status = 10").Find(&actionList)

	return Db.Transaction(func(tx *gorm.DB) error {
		if len(actionList) > 0 {
			for _, path := range actionList {
				arr := strings.Split(strings.Trim(path, "/"), "/")
				if len(arr) < 2 {
					continue
				}
				appName := arr[0]
				ctrl := arr[1]
				action := "index"
				if len(arr) > 2 {
					action = arr[2]
				}
				if err := tx.Where("v1 = ?", appName).
					Where("v2 = ?", ctrl).
					Where("v3 = ?", action).
					Delete(&model.Rules{}).Error; err != nil {
					return err
				}
			}
		}

		//清除字段信息
		if err := tx.Where("model_id = ?", mc.Id).Delete(&model.ModelField{}).Error; err != nil {
			return err
		}

		//清除索引信息
		if err := tx.Where("model_id = ?", mc.Id).Delete(&model.ModelIndex{}).Error; err != nil {
			return err
		}

		//清除字段选项
		if err := tx.Where("model_id = ?", mc.Id).Delete(&model.FieldOption{}).Error; err != nil {
			return err
		}

		//清除关联表信息
		if err := tx.Where("model_id = ?", mc.Id).Delete(&model.ModelRelation{}).Error; err != nil {
			return err
		}

		//清除动作信息
		if err := tx.Where("model_id = ?", mc.Id).Delete(&model.ModelAction{}).Error; err != nil {
			return err
		}

		//清除表单信息
		if err := tx.Where("model_id = ?", mc.Id).Delete(&model.ModelForm{}).Error; err != nil {
			return err
		}

		//清除表单校验规则信息
		if err := tx.Where("model_id = ?", mc.Id).Delete(&model.ModelFormRules{}).Error; err != nil {
			return err
		}

		//清除表单联动信息
		if err := tx.Where("model_id = ?", mc.Id).Delete(&model.ModelFormLinkage{}).Error; err != nil {
			return err
		}

		//清除菜单信息
		if err := tx.Where("model_id = ?", mc.Id).Delete(&model.Menu{}).Error; err != nil {
			return err
		}

		//删除模型对应表及相关数据
		return tx.Migrator().DropTable(NS.TableName(mc.TableName))
	})
}

// UpdateModel 根据模型ID更新模型文件
//	参数：
//		modelId 模型ID
func (makeSer *makeService) UpdateModel(modelId uint) error {
	var err error
	appList := AppConfig().GetAppListByModelId(modelId)
	tableName := ModelConfig().GetModelTableName(int(modelId))
	if len(appList) > 0 {
		for _, appName := range appList {
			if err = makeSer.Model(tableName, appName); err != nil {
				break
			}
		}
	}
	return err
}

//GetFieldSql 获取字段相关操作SQL
//	参数：
//		mf 模型字段实例
//		ac 动作名称 可选值：add  modify  del
//		oldFieldName 原字段名
func (makeSer *makeService) GetFieldSql(mf *model.ModelField, ac string, oldFieldName string) (string, error) {
	tableName := ModelConfig().GetModelTableName(int(mf.ModelId))

	if ac == "del" {
		return "ALTER TABLE `" + NS.TableName(tableName) + "` DROP " + mf.FieldName, nil
	}

	fieldLen := ""
	signed := ""
	isNull := ""

	if mf.Length > 0 {
		switch mf.Type {
		case "char", "varchar", "tinyint", "smallint", "int", "bigint":
			fieldLen = "(" + strconv.Itoa(int(mf.Length)) + ")"
		case "float", "double", "decimal":
			fieldLen = "(" + strconv.Itoa(int(mf.Length)) + ", " + strconv.Itoa(int(mf.DecimalLength)) + ")"
		}
	}
	if mf.IsSigned == 20 && (mf.Type == "tinyint" || mf.Type == "smallint" || mf.Type == "int" || mf.Type == "bigint" || mf.Type == "float" || mf.Type == "double" || mf.Type == "decimal") {
		signed = " unsigned "
	}

	if mf.Type != "text" && mf.Type != "mediumtext" && mf.Type != "longtext" {
		if mf.IsNull == 10 {
			isNull = " DEFAULT NULL "
		} else {
			isNull = " NOT NULL "
			defVal := ""
			if mf.Type == "datetime" || mf.Type == "timestamp" {
				if strings.HasPrefix(mf.FieldName, "create") || strings.HasPrefix(mf.FieldName, "add") {
					defVal = " DEFAULT CURRENT_TIMESTAMP "
				} else {
					defVal = " DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP "
				}
			} else if mf.DefaultValue == "" {
				switch mf.Type {
				case "char", "varchar":
					defVal = " DEFAULT '' "
				case "tinyint", "smallint", "int", "bigint", "float", "double", "decimal":
					defVal = " DEFAULT '0' "
				}
			} else {
				defVal = " DEFAULT '" + mf.DefaultValue + "' "
			}
			isNull += defVal
		}
	}

	comment := ""
	if mf.Note != "" {
		comment = " COMMENT '" + mf.Note + "' "
	} else if mf.Label != "" {
		comment = " COMMENT '" + mf.Label + "' "
	}

	acSql := ""
	if ac == "add" {
		acSql = " ADD `" + mf.FieldName + "` "
	} else if ac == "modify" {
		if oldFieldName == mf.FieldName {
			acSql = " MODIFY `" + mf.FieldName + "` "
		} else {
			acSql = " CHANGE `" + oldFieldName + "` `" + mf.FieldName + "` "
		}
	} else {
		return "", errors.New("参数ac只能为add或modify")
	}

	sql := "ALTER TABLE `" + NS.TableName(tableName) + "`" + acSql + mf.Type + fieldLen + signed + isNull + comment
	return sql, nil
}

// AddField 添加字段并更新模型文件
//	参数：
//		mf 模型字段实例
//		tx gorm.DB实例
func (makeSer *makeService) AddField(mf *model.ModelField, tx *gorm.DB) error {
	sql, err := makeSer.GetFieldSql(mf, "add", "")
	if err != nil {
		return err
	}
	if err = tx.Exec(sql).Error; err != nil {
		return err
	}
	//更新所有相关的模型文件
	return makeSer.UpdateModel(mf.ModelId)
}

// RenameField 添加字段并更新模型文件
//	参数：
//		mf 模型字段实例
//		oldFieldName 原字段名
//		tx gorm.DB实例
func (makeSer *makeService) RenameField(mf *model.ModelField, oldFieldName string, tx *gorm.DB) error {
	sql, err := makeSer.GetFieldSql(mf, "modify", oldFieldName)
	if err != nil {
		return err
	}
	if err = tx.Exec(sql).Error; err != nil {
		return err
	}
	//更新所有相关的模型文件
	return makeSer.UpdateModel(mf.ModelId)
}

// DelField 删除字段并更新模型文件
//	参数：
//		mf 模型字段实例
//		tx gorm.DB实例
func (makeSer *makeService) DelField(mf *model.ModelField, tx *gorm.DB) error {
	sql, err := makeSer.GetFieldSql(mf, "del", "")
	if err != nil {
		return err
	}
	if err = tx.Exec(sql).Error; err != nil {
		return err
	}

	//更新所有相关的模型文件
	return makeSer.UpdateModel(mf.ModelId)
}

// AddIndex 添加索引 并更新模型文件
//	参数：
//		mi 模型索引实例
//		tx gorm.DB实例
func (makeSer *makeService) AddIndex(mi *model.ModelIndex, tx *gorm.DB) error {
	if mi.ModelFieldId != "" {
		tableName := ModelConfig().GetModelTableName(int(mi.ModelId))
		indexType := mi.IndexType
		if indexType == "NORMAL" {
			indexType = ""
		}
		var fieldNameList []string
		tx.Table(NS.TableName("model_field")).Select("field_name").
			Where("id in ?", strings.Split(mi.ModelFieldId, ",")).
			Find(&fieldNameList)
		indexName := "idx_" + strings.Join(fieldNameList, "_")
		indexCol := "`" + strings.Join(fieldNameList, "`, `") + "`"

		sql := "ALTER TABLE `" + NS.TableName(tableName) + "` ADD " + indexType + " INDEX `" + indexName + "`(" + indexCol + ") USING BTREE"
		if err := tx.Exec(sql).Error; err != nil {
			return err
		}
		//更新所有相关的模型文件
		return makeSer.UpdateModel(mi.ModelId)
	}
	return nil
}

type ModelIndexRes struct {
	ModelFieldId string
	ModelId      uint
}

// DelIndex 删除索引 并更新模型文件
//	参数：
//		id 模型索引ID
//		tx gorm.DB实例
func (makeSer *makeService) DelIndex(id uint, tx *gorm.DB) error {
	var rs ModelIndexRes
	Db.Table(NS.TableName("model_index")).Select("model_field_id, model_id").
		Where("id = ?", id).Find(&rs)

	tableName := ModelConfig().GetModelTableName(int(rs.ModelId))
	if rs.ModelFieldId != "" {
		var fieldNameList []string
		Db.Table(NS.TableName("model_field")).Select("field_name").
			Where("id in ?", strings.Split(rs.ModelFieldId, ",")).
			Find(&fieldNameList)
		indexName := "idx_" + strings.Join(fieldNameList, "_")
		sql := "ALTER TABLE `" + NS.TableName(tableName) + "` DROP INDEX " + indexName

		if err := tx.Exec(sql).Error; err != nil {
			return err
		}
	}
	//更新所有相关的模型文件
	return makeSer.UpdateModel(rs.ModelId)
}

//CreateApp 创建应用相关目录
//	参数：
//		appName 应用名称
func (makeSer *makeService) CreateApp(appName string) error {
	//先创建目录
	appDir := "app/" + appName
	if _, err := os.Stat(appDir); err != nil {
		if err = os.MkdirAll(appDir, 0666); err != nil {
			return errors.New("创建应用" + appName + "失败！" + err.Error())
		}
	}

	//创建控制器目录
	controllerDir := appDir + "/controller"
	if _, err := os.Stat(controllerDir); err != nil {
		if err = os.MkdirAll(controllerDir, 0666); err != nil {
			return errors.New("创建控制器controller失败！" + err.Error())
		}
	}

	//创建服务层目录
	serviceDir := appDir + "/service"
	if _, err := os.Stat(serviceDir); err != nil {
		if err = os.MkdirAll(serviceDir, 0666); err != nil {
			return errors.New("创建服务层目录service失败！" + err.Error())
		}
	}

	//创建模型层目录
	modelDir := appDir + "/model"
	if _, err := os.Stat(modelDir); err != nil {
		if err = os.MkdirAll(modelDir, 0666); err != nil {
			return errors.New("创建模型层目录model失败！" + err.Error())
		}
	}

	//创建视图层目录
	viewDir := "views/" + appName
	if _, err := os.Stat(viewDir); err != nil {
		if err = os.MkdirAll(viewDir, 0666); err != nil {
			return errors.New("创建视图层目录views失败！" + err.Error())
		}
	}
	return nil
}

//RenameApp 重命名应用名称
//	参数：
//		appId 应用ID
//		newAppName 新应用名称
func (makeSer *makeService) RenameApp(appId uint, newAppName string) error {
	var oldAppName string
	Db.Table(NS.TableName("app_config")).Select("app_name").
		Where("id = ?", appId).Find(&oldAppName)

	if oldAppName == newAppName {
		return nil
	}

	//重命名应用目录
	oldAppDir := "app/" + oldAppName
	newAppDir := "app/" + newAppName
	//判断新目录是否已经存在
	if _, err := os.Stat(newAppDir); err == nil {
		return errors.New("应用目录" + newAppDir + "已存在!")
	}
	//重命名应用目录
	if _, err := os.Stat(oldAppDir); err == nil {
		if err = os.Rename(oldAppDir, newAppDir); err != nil {
			return errors.New("应用" + oldAppDir + "更新失败！" + err.Error())
		}
	}

	return nil
}

//RemoveApp 移除应用
//	参数：
//		appId 应用ID
func (makeSer *makeService) RemoveApp(appId uint) error {
	var appName string
	Db.Table(NS.TableName("app_config")).Select("app_name").
		Where("id = ?", appId).Find(&appName)

	appDir := "app/" + appName

	if _, err := os.Stat(appDir); err == nil {
		if err = os.RemoveAll(appDir); err != nil {
			return errors.New("应用" + appName + "移除失败！" + err.Error())
		}
	} else {
		return errors.New("应用" + appName + "移除失败！" + err.Error())
	}
	return nil
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
