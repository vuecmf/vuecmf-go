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
	//"errors"
	"github.com/gin-gonic/gin"
	"strings"
)

// request 定义request结构体，加入gin上下文
type request struct {
	*gin.Context
}

// Input 获取客户端GET/POST请求及header的所有输入数据
func (r *request) Input(param string, bindParam interface{}) {
	param = strings.ToLower(param)

	switch param {
	case "post":
		_ = r.ShouldBind(bindParam)
	case "get":
		_ = r.ShouldBindQuery(bindParam)
	case "header":
		_ = r.ShouldBindHeader(bindParam)
	default:
		panic("输入参数有误！只支持post, get, header")
	}

	return
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

// Request 获取请求实例
func Request(ctx *gin.Context) *request {
	return &request{
		Context: ctx,
	}
}
