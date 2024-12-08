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

// ModelActionController 模型动作
type ModelActionController struct {
	BaseController
	Svc *service.ModelActionService
}

var modelActionController *ModelActionController
var modelActionCtrlOnce sync.Once

// ModelAction 获取控制器实例
func ModelAction() *ModelActionController {
	modelActionCtrlOnce.Do(func() {
		modelActionController = &ModelActionController{
			Svc: service.ModelAction(),
		}
	})
	return modelActionController
}

// Action 控制器入口
func (ctrl ModelActionController) Action() gin.HandlerFunc {
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
		case "get_api_map":
			res, err = ctrl.getApiMap(c)
		case "get_action_list":
			res, err = ctrl.getActionList(c)
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
func (ctrl ModelActionController) save(c *gin.Context) (int64, error) {
	var params *model.DataModelActionForm
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
func (ctrl ModelActionController) delete(c *gin.Context) (int64, error) {
	var params *model.DataIdForm
	err := Post(c, &params)
	if err != nil {
		return 0, err
	}

	return ctrl.Svc.Delete(params.Data.Id, &model.ModelAction{})
}

// deleteBatch 根据ID列表批量删除多条数据
func (ctrl ModelActionController) deleteBatch(c *gin.Context) (int64, error) {
	var params *model.DataIdListForm
	err := Post(c, &params)
	if err != nil {
		return 0, err
	}

	return ctrl.Svc.DeleteBatch(params.Data.IdList, &model.ModelAction{})
}

// getApiMap 获取API映射的路径
func (ctrl ModelActionController) getApiMap(c *gin.Context) (string, error) {
	var params *model.DataApiMapForm
	err := Post(c, &params)
	if err != nil {
		return "", err
	}

	apiPath := ctrl.Svc.GetApiMap(params.Data.TableName, params.Data.ActionType, params.Data.AppId)
	return apiPath, nil
}

// getActionList 获取所有模型的动作列表
func (ctrl ModelActionController) getActionList(c *gin.Context) (any, error) {
	var params *model.DataActionListForm
	err := Post(c, &params)
	if err != nil {
		return "", err
	}

	return ctrl.Svc.GetActionList(params)
}
