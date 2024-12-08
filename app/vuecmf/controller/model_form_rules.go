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

// ModelFormRulesController 模型表单验证规则
type ModelFormRulesController struct {
	BaseController
	Svc *service.ModelFormRulesService
}

var modelFormRulesController *ModelFormRulesController
var modelFormRulesCtrlOnce sync.Once

// ModelFormRules 获取控制器实例
func ModelFormRules() *ModelFormRulesController {
	modelFormRulesCtrlOnce.Do(func() {
		modelFormRulesController = &ModelFormRulesController{
			Svc: service.ModelFormRules(),
		}
	})
	return modelFormRulesController
}

// Action 控制器入口
func (ctrl ModelFormRulesController) Action() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var res any

		switch GetActionName(c) {
		case "save":
			res, err = ctrl.save(c)
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
func (ctrl ModelFormRulesController) save(c *gin.Context) (int64, error) {
	var params *model.DataModelFormRulesForm
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

// delete 根据ID删除单条数据
func (ctrl ModelFormRulesController) delete(c *gin.Context) (int64, error) {
	var params *model.DataIdForm
	err := Post(c, &params)
	if err != nil {
		return 0, err
	}
	return ctrl.Svc.Delete(params.Data.Id)
}
