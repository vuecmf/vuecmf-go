//+----------------------------------------------------------------------
// | Copyright (c) 2023 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: vuecmf <tulihua2004@126.com>
// +----------------------------------------------------------------------

// Package model 模型
package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
	"time"
)

// DataBatchForm 批量导入数据 提交的表单数据
type DataBatchForm struct {
	Data string `json:"data" form:"data" binding:"required" required_tips:"参数data不能为空"`
}

// IdForm 存储ID值
type IdForm struct {
	Id uint `json:"id" form:"id"`
}

// DataIdForm 根据ID获取详情
type DataIdForm struct {
	Data *IdForm `json:"data" form:"data" binding:"required" required_tips:"参数data不能为空"`
}

// IdListForm 存储ID列表
type IdListForm struct {
	IdList string `json:"id_list" form:"id_list"`
}

//DataIdListForm 根据ID批量删除
type DataIdListForm struct {
	Data *IdListForm `json:"data" form:"data" binding:"required" required_tips:"参数data不能为空"`
}

// DropdownForm 存储Dropdown数据
type DropdownForm struct {
	TableName       string `json:"table_name" form:"table_name"`
	ModelId         uint   `json:"model_id" form:"model_id"`
	RelationModelId uint   `json:"relation_model_id" form:"relation_model_id"`
	AppId           uint   `json:"app_id" form:"app_id"`
	LinkageFieldId  uint   `json:"linkage_field_id" form:"linkage_field_id"`
	ModelFieldId    uint   `json:"model_field_id" form:"model_field_id"`
}

// DataDropdownForm 获取Dropdown数据所传form
type DataDropdownForm struct {
	Data *DropdownForm `json:"data" form:"data" binding:"required" required_tips:"参数data不能为空"`
}

// GetError 获取form中错误提示信息
func GetError(errs error, f interface{}) error {
	fData := reflect.ValueOf(f).Elem().FieldByName("Data")
	fDataType := reflect.TypeOf(fData.Interface())

	errList, ok := errs.(validator.ValidationErrors)
	if !ok {
		return errs
	}

	for _, fieldError := range errList {
		fieldName := fieldError.Field()      //获取验证的字段名
		tagKey := fieldError.Tag() + "_tips" //错误提示的tag key

		//data参数检测
		if fieldName == "Data" {
			dataField, dfExist := reflect.TypeOf(f).Elem().FieldByName("Data")
			if dfExist {
				tagTips := dataField.Tag.Get(tagKey)
				if tagTips != "" {
					return errors.New(tagTips)
				}
			}
		}

		//data内部参数检测
		field, exist := fDataType.Elem().FieldByName(fieldName) //根据字段名获取表单结构体中的字段
		if exist {
			tagTips := field.Tag.Get(tagKey)
			if tagTips != "" {
				return errors.New(field.Tag.Get(tagKey)) //获取结构字段中设置的tag标签信息
			}
		}
		return errors.New(fieldName + " " + fieldError.Tag())
	}
	return errs
}

//时间格式化
const (
	TimeFormat = "2006-01-02 15:04:05"
	DateFormat = "2006-01-02"
)

//JSONTime JSON时间格式化
type JSONTime struct {
	time.Time
}

//MarshalJSON 格式化时间后输出JSON
func (t JSONTime) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}
	formatted := fmt.Sprintf("\"%s\"", t.Format(TimeFormat))
	return []byte(formatted), nil
}

//UnmarshalJSON 解析JSON中时间后 转换成Time类型
func (t *JSONTime) UnmarshalJSON(data []byte) error {
	if len(data) == 2 {
		*t = JSONTime{Time: time.Time{}}
		return nil
	}
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now, err := time.ParseInLocation("\""+TimeFormat+"\"", string(data), loc)
	*t = JSONTime{Time: now}
	return err
}

// Value 写入数据库时的值检查
func (t JSONTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan DB查询时，检测是否为Time类型
func (t *JSONTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = JSONTime{Time: value}
		return nil
	}
	return fmt.Errorf("无法解析的时间 %v ", v)
}

//String 将Time类型格式化输出为字符串
func (t JSONTime) String() string {
	return t.Format(TimeFormat)
}

//JSONDate JSON时间格式化
type JSONDate struct {
	time.Time
}

//MarshalJSON 格式化时间后输出JSON
func (t JSONDate) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}
	formatted := fmt.Sprintf("\"%s\"", t.Format(DateFormat))
	return []byte(formatted), nil
}

//UnmarshalJSON 解析JSON中时间后 转换成Time类型
func (t *JSONDate) UnmarshalJSON(data []byte) error {
	if len(data) == 2 {
		*t = JSONDate{Time: time.Time{}}
		return nil
	}
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now, err := time.ParseInLocation("\""+DateFormat+"\"", string(data), loc)
	*t = JSONDate{Time: now}
	return err
}

// Value 写入数据库时的值检查
func (t JSONDate) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan DB查询时，检测是否为Time类型
func (t *JSONDate) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = JSONDate{Time: value}
		return nil
	}
	return fmt.Errorf("无法解析的时间 %v ", v)
}

//String 将Time类型格式化输出为字符串
func (t JSONDate) String() string {
	return t.Format(DateFormat)
}
