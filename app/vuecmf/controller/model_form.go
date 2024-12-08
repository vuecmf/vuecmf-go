//+----------------------------------------------------------------------
// | Copyright (c) 2024 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: tulihua2004@126.com
// +----------------------------------------------------------------------

package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/vuecmf/vuecmf-go/v3/app/vuecmf/model"
	"github.com/vuecmf/vuecmf-go/v3/app/vuecmf/service"
	"sync"
)

// ModelFormController 模型表单管理
type ModelFormController struct {
	BaseController
	Svc *service.ModelFormService
}

var modelFormController *ModelFormController
var modelFormCtrlOnce sync.Once

// ModelForm 获取控制器实例
func ModelForm() *ModelFormController {
	modelFormCtrlOnce.Do(func() {
		modelFormController = &ModelFormController{
			Svc: service.ModelForm(),
		}
	})
	return modelFormController
}

// Action 控制器入口
func (ctrl ModelFormController) Action() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var res any

		switch GetActionName(c) {
		case "save":
			res, err = ctrl.save(c)
		default:
			res, err = ctrl.BaseController.Action(c, ctrl.Svc.BaseService)
		}

		//将处理结果传入到After后置拦截器中统一处理
		if err != nil {
			c.Set("error", err)
		} else {
			c.Set("result", res)
		}
		c.Next()
	}
}

// save 新增/更新 单条数据
func (ctrl ModelFormController) save(c *gin.Context) (int64, error) {
	var params *model.DataModelFormForm
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
