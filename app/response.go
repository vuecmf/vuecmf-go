package app

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

type IResponse interface {
	SendJson(code int, msg string, data interface{})
	SendSuccess(msg string, data interface{}, code ...int)
	SendFailure(msg string, data interface{}, code ...int)
	Json(w http.ResponseWriter, data any) error
}

// Response 定义Response，加入gin上下文
type Response struct {
	*gin.Context
}

// SendJson 输出JSON内容到客户端
func (r *Response) SendJson(code int, msg string, data interface{}) {
	r.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg": msg,
		"data": data,
	})
}

// SendSuccess 成功返回信息到客户端
func (r *Response) SendSuccess(msg string, data interface{}, code ...int) {
	codeNum := 0
	if 0 != len(code) {
		codeNum = code[0]
	}
	r.SendJson(codeNum, msg, data)
}

// SendFailure 失败返回信息到客户端
func (r *Response) SendFailure(msg string, data interface{}, code ...int) {
	codeNum := 500
	if 0 != len(code) {
		codeNum = code[0]
	}
	r.SendJson(codeNum, msg, data)
}


func Json(w http.ResponseWriter, data any) error {
	header := w.Header()
	header["Content-Type"] = []string{"application/json; charset=utf-8"}
	content, _ := json.Marshal(data)

	w.WriteHeader(200)
	_, err := w.Write(content)
	return err
}
