//+----------------------------------------------------------------------
// | Copyright (c) 2024 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: tulihua2004@126.com
// +----------------------------------------------------------------------

package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/vuecmf/vuecmf-go/v3/app"
	"github.com/vuecmf/vuecmf-go/v3/app/vuecmf/model"
	"github.com/vuecmf/vuecmf-go/v3/app/vuecmf/service"
	"sync"
)

// UploadController 文件上传
type UploadController struct {
	BaseController
	Svc *service.UploadService
}

var uploadController *UploadController
var uploadCtrlOnce sync.Once

// Upload 文件上传控制器实例
func Upload() *UploadController {
	uploadCtrlOnce.Do(func() {
		uploadController = &UploadController{
			Svc: service.Upload(),
		}
	})
	return uploadController
}

// Action 控制器入口
func (ctrl UploadController) Action() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var res any

		switch GetActionName(c) {
		case "save":
			res, err = ctrl.save(c)
		case "":
			res, err = ctrl.index(c)
		default:
			res, err = ctrl.BaseController.Action(c, ctrl.Svc.BaseService)
		}

		if err != nil {
			c.Set("error", err)
		} else {
			c.Set("result", res)
		}

		c.Next()
	}
}

// save 新增/更新 单条数据
func (ctrl UploadController) save(c *gin.Context) (int64, error) {
	var params *model.DataUploadForm
	err := Post(c, &params)
	if err != nil {
		return 0, err
	}
	if params.Data.Id == uint(0) {
		return ctrl.Svc.Create(params.Data)
	} else {
		return ctrl.Svc.Update(params.Data)
	}
}

// index 文件上传
func (ctrl UploadController) index(c *gin.Context) (any, error) {
	fieldName := app.Request(c).Post("field_name")
	if fieldName == "" {
		return nil, errors.New("上传字段名(field_name)不能为空")
	}
	return ctrl.Svc.UploadFile(fieldName, c)
}
