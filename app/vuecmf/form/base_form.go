package form

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
)

type Former interface {

}

// GetError 获取form中错误提示信息
func GetError(errs error, f interface{}) string {
	//fType := reflect.TypeOf(f)
	fVal := reflect.ValueOf(f)
	for _, fieldError := range errs.(validator.ValidationErrors) {
		fieldName := fieldError.Field()                  //获取验证的字段名

		params := fVal.Elem().FieldByName("Data")
		fmt.Println("parames = ", params)


		//field, exist := fType.Elem().FieldByName(fieldName) //根据字段名获取表单结构体中的字段

		paramsType := reflect.TypeOf(params)

		field, exist := paramsType.FieldByName(fieldName)


		fmt.Println("num=", paramsType.NumField())
		for i := 0; i < paramsType.NumField(); i++ {
			fieldType := paramsType.Field(i)
			fmt.Println(fieldType.Name, fieldType.Tag)
		}

		if exist {
			return field.Tag.Get(fieldError.Tag() + "_tips") //获取结构字段中设置的tag标签信息
		}
		return fieldName + " " + fieldError.Tag()
	}
	return errs.Error()
}
