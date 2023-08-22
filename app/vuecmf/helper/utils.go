//+----------------------------------------------------------------------
// | Copyright (c) 2023 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: vuecmf <tulihua2004@126.com>
// +----------------------------------------------------------------------

// Package helper 助手工具
package helper

import (
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// ToFirstUpper 字符串首字母转大写
//	参数：
// 		str 需要转换的字符串
func ToFirstUpper(str string) string {
	strArr := []rune(str)
	strArr[0] -= 32
	return string(strArr)
}

// ToFirstLower 字符串首字母转小写
//	参数：
// 		str 需要转换的字符串
func ToFirstLower(str string) string {
	strArr := []rune(str)
	strArr[0] += 32
	return string(strArr)
}

// UnderToCamel 下横线转驼峰风格
//	参数：
// 		str 需要转换的字符串
func UnderToCamel(str string) string {
	str = strings.Replace(str, "_", " ", -1)
	str = cases.Title(language.Und).String(str)
	str = strings.Replace(str, " ", "", -1)
	return str
}

// CamelToUnder 驼峰转下横线
//	参数：
// 		str 需要转换的字符串
func CamelToUnder(str string) string {
	var output []rune

	for i, c := range str {
		if i == 0 && c < 91 && c > 64 {
			output = append(output, c+32)
		} else if i > 0 && c < 91 && c > 64 {
			output = append(output, 95)
			output = append(output, c+32)
		} else {
			output = append(output, c)
		}
	}

	return string(output)
}

//StrSliceToMap 字符串切片转map
func StrSliceToMap(items []string) map[string]struct{} {
	res := make(map[string]struct{}, len(items))
	for _, v := range items {
		res[v] = struct{}{}
	}
	return res
}

// InSlice 判断字符串是否在指定的切片中
//	参数：
// 		item 需要判断的字符串
// 		items 指定的字符串切片
func InSlice(item string, items []string) bool {
	m := StrSliceToMap(items)
	_, ok := m[item]
	return ok
}

// SliceRemove 删除字符串切片中元素
//	参数：
// 		slice 指定的字符串切片
// 		index 需要删除的切片索引
func SliceRemove(slice []string, index int) []string {
	return append(slice[:index], slice[index+1:]...)
}

// PasswordHash 加密密码
//	参数：
// 		password 需要加密的密码
func PasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// PasswordVerify 验证密码是否正确
//	参数：
// 		password 需要验证的密码
// 		hash 哈希值
func PasswordVerify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GetFileExt 获取文件名的扩展名
//	参数：
// 		fileName 文件名
func GetFileExt(fileName string) string {
	arr := strings.Split(fileName, ".")
	return strings.ToLower(arr[len(arr)-1])
}

// GetFileBaseName 获取不包含扩展名的文件名称
//	参数：
// 		fileName 文件名
func GetFileBaseName(fileName string) string {
	arr := strings.Split(fileName, ".")
	if len(arr) > 1 {
		return strings.Join(SliceRemove(arr, len(arr)-1), ".")
	} else {
		return fileName
	}
}

// GetRandomString 生成图片名字
//	参数：
// 		length 需要生成的字符串名字长度
func GetRandomString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	bytesLen := len(bytes)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(bytesLen)])
	}
	return string(result)
}

//SetString 对[]string 类型的切片进行元素唯一化处理
//	参数：
// 		arr 需要处理的字符串切片
func SetString(arr []string) []string {
	newArr := make([]string, 0)

	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
			}

		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return newArr
}

// InterfaceToInt interface类型转换成int
//	参数：
// 		val 需要转换的值
func InterfaceToInt(val interface{}) int {
	var res int
	switch val.(type) {
	case uint:
		res = int(val.(uint))
	case int8:
		res = int(val.(int8))
	case uint8:
		res = int(val.(uint8))
	case int16:
		res = int(val.(int16))
	case uint16:
		res = int(val.(uint16))
	case int32:
		res = int(val.(int32))
	case uint32:
		res = int(val.(uint32))
	case int64:
		res = int(val.(int64))
	case uint64:
		res = int(val.(uint64))
	case float32:
		res = int(val.(float32))
	case float64:
		res = int(val.(float64))
	case string:
		res, _ = strconv.Atoi(val.(string))
	default:
		res = val.(int)
	}
	return res
}

// TreeRes 存储目录树结果
type TreeRes struct {
	Id    int    //主键值
	Label string //标题
}

// ModelFieldOption 存储模型字段选项
type ModelFieldOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// FormatTree 格式化下拉树型列表
//	参数：
// 		tree map				存储返回的结果
// 		tableName string        表名
//		filter string			过滤条件
//		pk string				主键字段名称
// 		pid int                 父级ID
// 		label string            标题字段名
// 		pidField string         父级字段名
// 		orderField string       排序字段名
// 		level int               层级数
//	返回值：map
func FormatTree(tree []*ModelFieldOption, db *gorm.DB, tableName string, filter map[string]interface{}, pk string, pid int, label string, pidField string, orderField string, level int) []*ModelFieldOption {
	//参数为空的，设置默认值
	if label == "" {
		label = "title"
	}
	if pidField == "" {
		pidField = "pid"
	}

	var treeResList []TreeRes
	var childTotal int64

	model := db.Table(tableName).Select(label+" label,"+pk+" id").
		Where(pidField+" = ?", pid).
		Where("status = 10")

	if filter != nil {
		model = model.Where(filter)
	}
	if orderField != "" {
		model = model.Order(orderField)
	}

	model.Find(&treeResList)

	for key, val := range treeResList {
		prefix := strings.Repeat("┊ ", level-1)

		totalQuery := db.Table(tableName).Where(pidField+" = ?", val.Id).Where("status = 10")
		if filter != nil {
			totalQuery = totalQuery.Where(filter)
		}
		totalQuery.Count(&childTotal)

		if childTotal > 0 || key != len(treeResList)-1 {
			prefix += "┊┈ "
		} else {
			prefix += "└─ "
		}
		//tree[strconv.Itoa(val.Id)] = prefix + val.Label
		tree = append(tree, &ModelFieldOption{
			Value: strconv.Itoa(val.Id),
			Label: prefix + val.Label,
		})

		tree = FormatTree(tree, db, tableName, filter, pk, val.Id, label, pidField, orderField, level+1)
	}

	return tree
}

//MinusStrList 字符串差集
func MinusStrList(strA, strB []string) []string {
	var res []string
	tmp := make(map[string]bool)
	for _, str := range strA {
		if _, ok := tmp[str]; !ok {
			tmp[str] = true
		}
	}

	for _, str := range strB {
		if _, ok := tmp[str]; ok {
			delete(tmp, str)
		}
	}

	for k, _ := range tmp {
		res = append(res, k)
	}
	return res
}
