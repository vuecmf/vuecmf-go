package app

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
)


// ToFirstUpper 字符串首字母转大写
func ToFirstUpper(str string) string {
	str_arr := []rune(str)
	str_arr[0] -= 32
	return string(str_arr)
}

// UnderToCamel 下横线转驼峰风格
func UnderToCamel(str string) string {
	str = strings.Replace(str,"_"," ", -1)
	str = cases.Title(language.Und).String(str)
	str = strings.Replace(str, " ","", -1)
	return str
}

// CamelToUnder 驼峰转下横线
func CamelToUnder(str string) string {
	var output []rune

	for i,c := range str {
		if i == 0 && c < 91 {
			output = append(output, c + 32)
		}else if i > 0 && c < 91 {
			output = append(output, 95)
			output = append(output, c + 32)
		}else{
			output = append(output, c)
		}
	}

	return string(output)
}