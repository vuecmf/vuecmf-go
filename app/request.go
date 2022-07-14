package app

import (
	//"errors"
	"github.com/gin-gonic/gin"
	"strings"
)

type IRequest interface {
	Input(param string, bindParam interface{})
	Get(fieldName string) string
	Post(fieldName string) string
	Header(fieldName string) string
}

// Request 定义Request，加入gin上下文
type Request struct {
	*gin.Context
}

// Input 获取客户端GET/POST请求及header的所有输入数据
func (r *Request) Input(param string, bindParam interface{}) {
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
		//err = errors.New("输入参数有误！只支持post, get, header")
	}

	return
}

// Get 获取GET请求参数
func (r *Request) Get(fieldName string) string {
	return r.Query(fieldName)
}

// Post 获取formData方式的POST请求参数
func (r *Request) Post(fieldName string) string {
	return r.PostForm(fieldName)
}

// Header 获取头信息中数据
func (r *Request) Header(fieldName string) string {
	return r.GetHeader(fieldName)
}