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

// ModelRelationController 模型关联管理
type ModelRelationController struct {
	BaseController
	Svc *service.ModelRelationService
}

var modelRelationController *ModelRelationController
var modelRelationCtrlOnce sync.Once

// ModelRelation 获取控制器实例
func ModelRelation() *ModelRelationController {
	modelRelationCtrlOnce.Do(func() {
		modelRelationController = &ModelRelationController{
			Svc: service.ModelRelation(),
		}
	})
	return modelRelationController
}

// Action 控制器入口
func (ctrl ModelRelationController) Action() gin.HandlerFunc {
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
func (ctrl ModelRelationController) save(c *gin.Context) (int64, error) {
	var params *model.DataModelRelationForm
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
