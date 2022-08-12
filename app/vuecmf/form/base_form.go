package form

import (
	"github.com/go-playground/validator/v10"
	"reflect"
)

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
