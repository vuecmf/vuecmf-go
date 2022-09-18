// Package controller
//+----------------------------------------------------------------------
// | Copyright (c) 2022 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: vuecmf <tulihua2004@126.com>
// +----------------------------------------------------------------------
package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/vuecmf/vuecmf-go/app"
	"github.com/vuecmf/vuecmf-go/app/route"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/model"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/service"
)

type Upload struct {
	Base
}

func init() {
	upload := &Upload{}
	upload.TableName = "upload"
	upload.Model = &model.Upload{}
	upload.listData = &[]model.Upload{}
	upload.saveForm = &model.DataUploadForm{}
	upload.filterFields = []string{""}

	route.Register(upload, "POST", "vuecmf")
}

// Index 文件上传
func (ctrl *Upload) Index(c *gin.Context) {
	common(c, nil, func() (interface{}, error) {
		fieldName := app.Request(c).Post("field_name")
		if fieldName == "" {
			return nil, errors.New("上传字段名(field_name)不能为空")
		}
		return service.Upload().UploadFile(fieldName, c)
	})
}
