//+----------------------------------------------------------------------
// | Copyright (c) 2024 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: tulihua2004@126.com
// +----------------------------------------------------------------------

// Package controller 控制器
package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/vuecmf/vuecmf-go/v3/app"
	"github.com/vuecmf/vuecmf-go/v3/app/vuecmf/helper"
	"github.com/vuecmf/vuecmf-go/v3/app/vuecmf/model"
	"github.com/vuecmf/vuecmf-go/v3/app/vuecmf/service"
	"strconv"
	"strings"
	"sync"
)

// BaseController 基础结构
type BaseController struct {
}

var baseController *BaseController
var baseCtrlOnce sync.Once

// Base 获取控制器实例
func Base() *BaseController {
	baseCtrlOnce.Do(func() {
		baseController = &BaseController{}
	})
	return baseController
}

// IController 控制器接口
type IController interface {
	Before() gin.HandlerFunc // 前置拦截器
	Action() gin.HandlerFunc // 动作
	After() gin.HandlerFunc  // 后置拦截器
}

// MGet 获取中间件或控制器间传的数据
//
// 参数：
//
//	c gin上下文
//	key 传递数据的key
func MGet(c *gin.Context, key string) any {
	res, _ := c.Get(key)
	if key != "error" && res == nil {
		return ""
	}
	return res
}

// GetActionName 获取控制器间的当前动作名
//
// 参数：
//
//	c gin上下文
func GetActionName(c *gin.Context) string {
	return strings.TrimPrefix(c.Param("action"), "/")
}

// Get 获取get提交的数据
//
// 参数：
//
//	c gin上下文
//	params 提交的数据
func Get(c *gin.Context, params interface{}) error {
	err := app.Request(c).Input("get", params)

	if err != nil {
		var reason string
		if err.Error() == "EOF" {
			reason = "参数为空"
		} else {
			reason = service.GetErrMsg(model.GetError(err, params))
		}
		err = errors.New("请求失败：" + reason)
	}
	return err
}

// Post 获取POST提交的数据
//
// 参数：
//
//	c gin上下文
//	params 提交的数据
func Post(c *gin.Context, params interface{}) error {
	err := app.Request(c).Input("post", params)
	if err != nil {
		var reason string
		if err.Error() == "EOF" {
			reason = "参数为空"
		} else {
			reason = service.GetErrMsg(model.GetError(err, params))
		}
		err = errors.New("请求失败：" + reason)
	}
	return err
}

// Before 路由前置拦截器
func (ctrl BaseController) Before() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

// After 路由后置拦截器
func (ctrl BaseController) After() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := MGet(c, "error")
		if err != nil {
			app.Response(c).SendFailure(err.(error).Error(), nil)
			return
		}

		res := MGet(c, "result")
		app.Response(c).SendSuccess("请求成功", res)
	}
}

// Action 控制器入口
func (ctrl BaseController) Action(c *gin.Context, svc *service.BaseService) (any, error) {
	switch GetActionName(c) {
	case "": //列表页
		var isSuper uint16
		isSuperVal := MGet(c, "is_super")
		switch isSuperVal.(type) {
		case uint16:
			isSuper = isSuperVal.(uint16)
		case string:
			res, _ := strconv.Atoi(isSuperVal.(string))
			isSuper = uint16(res)
		}

		var params *helper.DataListParams
		err := Post(c, &params)
		if err != nil {
			return nil, err
		}
		return svc.CommonList(params, isSuper)
	case "save_all": //批量添加多条数据
		var params *model.DataBatchForm
		err := Post(c, &params)
		if err != nil {
			return nil, err
		}
		return svc.CreateAll(params.Data)
	case "detail": //根据ID获取详情
		var params *model.DataIdForm
		err := Post(c, &params)
		if err != nil {
			return nil, err
		}
		return svc.Detail(params.Data.Id)
	case "delete": //根据ID删除单条数据
		var params *model.DataIdForm
		err := Post(c, &params)
		if err != nil {
			return nil, err
		}
		return svc.Delete(params.Data.Id)
	case "delete_batch": //根据ID列表批量删除多条数据
		var params *model.DataIdListForm
		err := Post(c, &params)
		if err != nil {
			return nil, err
		}
		return svc.DeleteBatch(params.Data.IdList)
	case "dropdown": //下拉列表数据
		var params *model.DataDropdownForm
		err := Post(c, &params)
		if err != nil {
			return nil, err
		}
		return svc.Dropdown(params.Data)
	default:
		app.Response(c).NotFound()
		return nil, nil
	}
}
