// Package app
//+----------------------------------------------------------------------
// | Copyright (c) 2022 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: vuecmf <tulihua2004@126.com>
// +----------------------------------------------------------------------
package app

import (
	"errors"
	//"errors"
	"github.com/gin-gonic/gin"
	"strings"
)

// request 定义request结构体，加入gin上下文
type request struct {
	*gin.Context
}

// Input 获取客户端GET/POST请求及header的所有输入数据
func (r *request) Input(param string, bindParam interface{}) error {
	param = strings.ToLower(param)

	var err error

	switch param {
	case "post":
		err = r.ShouldBind(bindParam)
	case "get":
		err = r.ShouldBindQuery(bindParam)
	case "header":
		err = r.ShouldBindHeader(bindParam)
	default:
		err = errors.New("输入参数有误！只支持post, get, header")
	}

	return err
}

// Get 获取GET请求参数
func (r *request) Get(fieldName string) string {
	return r.Query(fieldName)
}

// Post 获取formData方式的POST请求参数
func (r *request) Post(fieldName string) string {
	return r.PostForm(fieldName)
}

// Header 获取头信息中数据
func (r *request) Header(fieldName string) string {
	return r.GetHeader(fieldName)
}

// GetCtxVal 根据key获取中间件中传入的数据
func (r *request) GetCtxVal(key string) interface{} {
	val, exist := r.Context.Get(key)
	if exist == false {
		return nil
	}

	switch val.(type) {
	case uint:
		return int(val.(uint))
	case int8:
		return int(val.(int8))
	case uint8:
		return int(val.(uint8))
	case int16:
		return int(val.(int16))
	case uint16:
		return int(val.(uint16))
	case int32:
		return int(val.(int32))
	case uint32:
		return int(val.(uint32))
	case int64:
		return int(val.(int64))
	case uint64:
		return int(val.(uint64))
	case float32:
		return val.(float32)
	case float64:
		return val.(float64)
	case string:
		return val.(string)
	}
	return val
}

// Request 获取请求实例
func Request(ctx *gin.Context) *request {
	return &request{
		Context: ctx,
	}
}
