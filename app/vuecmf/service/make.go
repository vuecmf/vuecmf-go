package service

import (
	"github.com/vuecmf/vuecmf-go/app"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/helper"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/model"
	"io/ioutil"
	"strconv"
	"strings"
)

// MakeService make服务结构
type MakeService struct {
	*base
}

//Model 功能：生成模型文件
//		参数：tableName string 表名（不带表前缀）
func (service *MakeService) Model(tableName string) bool {
	var result []model.ModelField

	//查出需要生成模型表的字段相关信息
	db := app.Db("default")
	db.Table("model_field MF").Db.
		Select("MF.*").
		Joins("left join " + db.Conf.Prefix + "model_config MC on MF.model_id = MC.id").
		Where("MF.field_name NOT IN('id','status')").
		Where("MC.table_name = ?", tableName).Scan(&result)

	//读取模型模板文件
	tplContent, err := ioutil.ReadFile("app/vuecmf/make/stubs/model.stub")
	if err != nil {
		//fmt.Println("读取model模板失败")
		panic("读取model模板失败")
	}

	modelContent := ""
	hasTime := false

	//模型字段信息处理
	for _, value := range result{
		notNull := ""
		defaultVal := ""
		size := ""
		autoCreateTime := ""
		fieldType := "string"
		uniqueIndex := ""

		if value.Type == "timestamp" {
			hasTime = true
			fieldType = "time.Time"
		}else if value.Type == "int" || value.Type == "bigint" {
			fieldType = "int"
		}else if value.Type == "smallint" {
			fieldType = "int16"
		}else if value.Type == "tinyint"{
			fieldType = "int8"
		}else if value.Type == "float" {
			fieldType = "float32"
		}else if value.Type == "double" || value.Type == "decimal" {
			fieldType = "float64"
		}

		if value.IsNull == 20 {
			notNull = "not null;"
		}

		if value.FieldName == "update_time" || value.FieldName == "last_login_time" || value.DefaultValue == "CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" {
			autoCreateTime = "autoCreateTime;autoUpdateTime;"
		}else if value.DefaultValue == "CURRENT_TIMESTAMP" {
			autoCreateTime = "autoCreateTime;"
		}else{
			defaultVal = "default:" + value.DefaultValue + ";"
			size = "size:" + strconv.Itoa(value.Length) + ";"
		}

		//字段唯一索引处理
		modelIndexId := 0
		id := strconv.Itoa(int(value.Id))
		db.Table("model_index").Db.Select("id").
			Where("model_field_id = ? or model_field_id like ? or model_field_id like ?", id, id + ",%", "%," + id).
			Find(&modelIndexId)

		if modelIndexId > 0 {
			uniqueIndex = "uniqueIndex:unique_index;"
		}

		modelContent += helper.UnderToCamel(value.FieldName) + " " + fieldType + " `json:\"" + value.FieldName +
			"\" gorm:\"column:" + value.FieldName + ";" + size + uniqueIndex + notNull + autoCreateTime + defaultVal +
			"comment:"+ value.Note +"\"`\n\t"
	}

	modelLabel := ""
	db.Table("model_config").Db.Select("label").
		Where("table_name = ?", tableName).Find(&modelLabel)

	//替换模板文件中内容
	txt := string(tplContent)
	txt = strings.Replace(txt, "{{.comment}}", modelLabel, -1)
	txt = strings.Replace(txt, "{{.model_name}}", helper.UnderToCamel(tableName), -1)

	if hasTime == true {
		txt = strings.Replace(txt,"{{.import}}", "import \"time\"", -1)
	}else{
		txt = strings.Replace(txt,"{{.import}}", "", -1)
	}

	txt = strings.Replace(txt,"{{.body}}", modelContent, -1)

	err = ioutil.WriteFile("app/vuecmf/model/"+ tableName +".go",[]byte(txt), 0666)

	if err != nil {
		return false
	}

	return true
}