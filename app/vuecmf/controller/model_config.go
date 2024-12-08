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

// ModelConfigController 模型配置管理
type ModelConfigController struct {
	BaseController
	Svc *service.ModelConfigService
}

var modelConfigController *ModelConfigController
var modelConfigCtrlOnce sync.Once

// ModelConfig 获取控制器实例
func ModelConfig() *ModelConfigController {
	modelConfigCtrlOnce.Do(func() {
		modelConfigController = &ModelConfigController{
			Svc: service.ModelConfig(),
		}
	})
	return modelConfigController
}

// Action 控制器入口
func (ctrl ModelConfigController) Action() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var res any

		switch GetActionName(c) {
		case "save":
			res, err = ctrl.save(c)
		case "delete":
			res, err = ctrl.delete(c)
		case "delete_batch":
			res, err = ctrl.deleteBatch(c)
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
func (ctrl ModelConfigController) save(c *gin.Context) (int64, error) {
	var params *model.DataModelConfigForm
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
func (ctrl ModelConfigController) delete(c *gin.Context) (int64, error) {
	var params *model.DataIdForm
	err := Post(c, &params)
	if err != nil {
		return 0, err
	}
	return ctrl.Svc.Delete(params.Data.Id, &model.ModelConfig{})
}

// deleteBatch 根据ID列表批量删除多条数据
func (ctrl ModelConfigController) deleteBatch(c *gin.Context) (int64, error) {
	var params *model.DataIdListForm
	err := Post(c, &params)
	if err != nil {
		return 0, err
	}
	return ctrl.Svc.DeleteBatch(params.Data.IdList, &model.ModelConfig{})
}
