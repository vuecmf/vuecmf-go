// Package helper
//+----------------------------------------------------------------------
// | Copyright (c) 2022 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: vuecmf <tulihua2004@126.com>
// +----------------------------------------------------------------------
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
func ToFirstUpper(str string) string {
	strArr := []rune(str)
	strArr[0] -= 32
	return string(strArr)
}

// ToFirstLower 字符串首字母转小写
func ToFirstLower(str string) string {
	strArr := []rune(str)
	strArr[0] += 32
	return string(strArr)
}

// UnderToCamel 下横线转驼峰风格
func UnderToCamel(str string) string {
	str = strings.Replace(str, "_", " ", -1)
	str = cases.Title(language.Und).String(str)
	str = strings.Replace(str, " ", "", -1)
	return str
}

// CamelToUnder 驼峰转下横线
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

// InSlice 判断字符串是否在指定的切片中
func InSlice(item string, items []string) bool {
	for _, val := range items {
		if val == item {
			return true
		}
	}
	return false
}

// SliceRemove 删除字符串切片中元素
func SliceRemove(slice []string, index int) []string {
	return append(slice[:index], slice[index+1:]...)
}

// PasswordHash 加密密码
func PasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// PasswordVerify 验证密码是否正确
func PasswordVerify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GetFileExt 获取文件名的扩展名
func GetFileExt(fileName string) string {
	arr := strings.Split(fileName, ".")
	return strings.ToLower(arr[len(arr)-1])
}

// GetFileBaseName 获取不包含扩展名的文件名称
func GetFileBaseName(fileName string) string {
	arr := strings.Split(fileName, ".")
	if len(arr) > 1 {
		return strings.Join(SliceRemove(arr, len(arr)-1), ".")
	} else {
		return fileName
	}
}

// GetRandomString 生成图片名字
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

// InterfaceToInt interface类型转换成int
func InterfaceToInt(val interface{}) int {
	var res int
	switch val.(type) {
	case uint:
		res = int(val.(uint))
		break
	case int8:
		res = int(val.(int8))
		break
	case uint8:
		res = int(val.(uint8))
		break
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

type TreeRes struct {
	Id    int    //主键值
	Label string //标题
}

// FormatTree 格式化下拉树型列表
//	参数：
// 		tree map				存储返回的结果
// 		tableName string        表名
//		pk string				主键字段名称
// 		pid int                 父级ID
// 		label string            标题字段名
// 		pidField string         父级字段名
// 		orderField string       排序字段名
// 		level int               层级数
//	返回值：map
func FormatTree(tree map[string]string, db *gorm.DB, tableName string, pk string, pid int, label string, pidField string, orderField string, level int) {
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
	if orderField != "" {
		model.Order(orderField)
	}
	model.Find(&treeResList)

	for key, val := range treeResList {
		prefix := strings.Repeat("┊ ", level-1)
		db.Table(tableName).Where(pidField+" = ?", val.Id).Where("status = 10").Count(&childTotal)
		if childTotal > 0 || key != len(treeResList)-1 {
			prefix += "┊┈ "
		} else {
			prefix += "└─ "
		}
		tree[strconv.Itoa(val.Id)] = prefix + val.Label
		FormatTree(tree, db, tableName, pk, val.Id, label, pidField, orderField, level+1)
	}
}

/*func TreeList(tree []model.MenuTree, db *gorm.DB, tableName string, pid int, keywords string, pidField string, searchField string, orderField string) []model.MenuTree {
	query := db.Table(tableName).Select("*").Where(pidField+" = ?", pid)
	if keywords != "" {
		query.Where(searchField+" like ?", "%"+keywords+"%")
	}
	if orderField != "" {
		query.Order(orderField)
	}
	var res = make([]model.Menu, 0)
	query.Find(&res)

	fmt.Println(res)

	for key, val := range res {


		child := TreeList(tree, db, tableName, int(val.Id), keywords, pidField, searchField, orderField)

		if child != nil {
			treeItem.Children = child
		}


	}
	fmt.Println("res=", res)

	return tree
}*/
