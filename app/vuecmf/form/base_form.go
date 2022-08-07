package form

import (
	"github.com/go-playground/validator/v10"
	"reflect"
)

// GetError 获取form中错误提示信息
func GetError(errs error, f interface{}) string {
	ft := reflect.TypeOf(f)
	for _, fieldError := range errs.(validator.ValidationErrors) {
		fieldName := fieldError.Field()                  //获取验证的字段名
		field, exist := ft.Elem().FieldByName(fieldName) //根据字段名获取表单结构体中的字段
		if exist {
			return field.Tag.Get("tips") //获取结构字段中设置的tag标签信息
		}
		return fieldName + " " + fieldError.Tag()
	}
	return errs.Error()
}
