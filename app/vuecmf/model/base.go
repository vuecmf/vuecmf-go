package model

import (
	"github.com/go-playground/validator/v10"
	"reflect"
)

// DataBatchForm 批量导入数据 提交的表单数据
type DataBatchForm struct {
	Data string `json:"data" form:"data"`
}

type IdForm struct {
	Id uint `json:"id" form:"id"`
}

// DataIdForm 根据ID获取详情
type DataIdForm struct {
	Data *IdForm `json:"data" form:"data"`
}

type IdListForm struct {
	IdList string `json:"id_list" form:"id_list"`
}

//DataIdListForm 根据ID批量删除
type DataIdListForm struct {
	Data *IdListForm `json:"data" form:"data"`
}

type DropdownForm struct {
	TableName       string `json:"table_name" form:"table_name"`
	ModelId         uint   `json:"model_id" form:"model_id"`
	RelationModelId uint   `json:"relation_model_id" form:"relation_model_id"`
}

// DataDropdownForm 获取Dropdown数据所传form
type DataDropdownForm struct {
	Data *DropdownForm `json:"data" form:"data"`
}

// GetError 获取form中错误提示信息
func GetError(errs error, f interface{}) string {
	fData := reflect.ValueOf(f).Elem().FieldByName("Data")
	fDataType := reflect.TypeOf(fData.Interface())

	for _, fieldError := range errs.(validator.ValidationErrors) {
		fieldName := fieldError.Field()      //获取验证的字段名
		tagKey := fieldError.Tag() + "_tips" //错误提示的tag key

		//data参数检测
		dataField, dfExist := reflect.TypeOf(f).Elem().FieldByName("Data")
		if dfExist {
			tagTips := dataField.Tag.Get(tagKey)
			if tagTips != "" {
				return tagTips
			}
		}

		//data内部参数检测
		field, exist := fDataType.Elem().FieldByName(fieldName) //根据字段名获取表单结构体中的字段
		if exist {
			tagTips := field.Tag.Get(tagKey)
			if tagTips != "" {
				return field.Tag.Get(tagKey) //获取结构字段中设置的tag标签信息
			}
		}
		return fieldName + " " + fieldError.Tag()
	}
	return errs.Error()
}
