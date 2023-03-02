//+----------------------------------------------------------------------
// | Copyright (c) 2023 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: vuecmf <tulihua2004@126.com>
// +----------------------------------------------------------------------

package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// response 定义response结构体，加入gin上下文
type response struct {
	*gin.Context
}

//SendHtml 输出并渲染网页
//参数：
//	tplName 模板名
func (r *response) SendHtml(tplName string, obj any) {
	r.HTML(http.StatusOK, tplName, obj)
}

//SendText 输出文本
//参数：
//	msg 输出的文本内容
func (r *response) SendText(msg string) {
	r.String(http.StatusOK, "%s", msg)
}

//SendJson 输出JSON内容到客户端
//参数：
//	code 响应码
//	msg 消息提示内容
//	data 返回的内容
func (r *response) SendJson(code int, msg string, data interface{}) {
	r.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

//SendSuccess 成功返回信息到客户端
//参数：
//	msg 消息提示内容
//	data 返回的内容
//	code 成功的响应码，默认0
func (r *response) SendSuccess(msg string, data interface{}, code ...int) {
	codeNum := 0
	if 0 != len(code) {
		codeNum = code[0]
	}
	r.SendJson(codeNum, msg, data)
}

//SendFailure 失败返回信息到客户端
//参数：
//	msg 消息提示内容
//	data 返回的内容
//	code 失败的响应码，默认500
func (r *response) SendFailure(msg string, data interface{}, code ...int) {
	codeNum := 500
	if 0 != len(code) {
		codeNum = code[0]
	}
	r.SendJson(codeNum, msg, data)
}

// Response 获取response实例
func Response(ctx *gin.Context) *response {
	return &response{
		Context: ctx,
	}
}

/*func Json(w http.ResponseWriter, data any) error {
	header := w.Header()
	header["Content-Type"] = []string{"application/json; charset=utf-8"}
	content, _ := json.Marshal(data)

	w.WriteHeader(200)
	_, err := w.Write(content)
	return err
}*/
