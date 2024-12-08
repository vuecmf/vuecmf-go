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

// AppConfigController 应用配置管理
type AppConfigController struct {
	BaseController
	Svc *service.AppConfigService
}

var appConfigController *AppConfigController
var appConfigCtrlOnce sync.Once

// AppConfig 获取控制器实例
func AppConfig() *AppConfigController {
	appConfigCtrlOnce.Do(func() {
		appConfigController = &AppConfigController{
			Svc: service.AppConfig(),
		}
	})
	return appConfigController
}

// Action 控制器入口
func (ctrl AppConfigController) Action() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var res any

		switch GetActionName(c) {
		case "save":
			res, err = ctrl.save(c)
		case "delete":
			res, err = ctrl.delete(c)
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

// Save 新增/更新 单条数据
func (ctrl AppConfigController) save(c *gin.Context) (int64, error) {
	var params *model.DataAppConfigForm
	err := Post(c, &params)
	if err != nil {
		return 0, err
	}
	if params.Data.Id == uint(0) {
		//创建应用目录
		if err = service.Make().CreateApp(params.Data.AppName, app.Config().Module); err != nil {
			return 0, err
		}
		return ctrl.Svc.Create(params.Data)
	} else {
		//更新应用目录
		if err = service.Make().RenameApp(params.Data.Id, params.Data.AppName); err != nil {
			return 0, err
		}
		return ctrl.Svc.Update(params.Data)
	}
}

// delete 根据ID删除单条数据
func (ctrl AppConfigController) delete(c *gin.Context) (int64, error) {
	var params *model.DataIdForm
	err := Post(c, &params)
	if err != nil {
		return 0, err
	}
	//先检查应用下是否存在模型，若存在则不允许删除
	if num := ctrl.Svc.GetAppModelCount(params.Data.Id); num > 0 {
		return 0, errors.New("不允许删除有分配模型的应用！")
	}
	if err = service.Make().RemoveApp(params.Data.Id); err != nil {
		return 0, err
	}
	return ctrl.Svc.Delete(params.Data.Id)
}
