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
	"github.com/vuecmf/vuecmf-go/app/route"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/model"
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
