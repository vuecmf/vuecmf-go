//+----------------------------------------------------------------------
// | Copyright (c) 2023 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: vuecmf <tulihua2004@126.com>
// +----------------------------------------------------------------------

package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/vuecmf/vuecmf-go/app/route"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/model"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/service"
)

// FieldOption 字段选项管理
type FieldOption struct {
	Base
}

func init() {
	fieldOption := &FieldOption{}
	fieldOption.TableName = "field_option"
	fieldOption.Model = &model.FieldOption{}
	fieldOption.ListData = &[]model.FieldOption{}
	fieldOption.FilterFields = []string{"option_value", "option_label"}

	route.Register(fieldOption, "POST", "vuecmf")
}

// Save 新增/更新 单条数据
func (ctrl *FieldOption) Save(c *gin.Context) {
	saveForm := &model.DataFieldOptionForm{}
	Common(c, saveForm, func() (interface{}, error) {
		if saveForm.Data.Id == uint(0) {
			return service.Base().Create(saveForm.Data)
		} else {
			return service.Base().Update(saveForm.Data)
		}
	})
}
