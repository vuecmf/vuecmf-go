package vuecmf

import (
	"github.com/gin-gonic/gin"
	"github.com/vuecmf/vuecmf-go/app/controller"
)

// Router 路由配置
func Router(r *gin.Engine) {
	//控制器实例
	adminCtrl := controller.AdminController{}
	makeCtrl := controller.MakeController{}

	r.POST("/vuecmf/make/model", makeCtrl.Model)

	//列表
	r.POST("/vuecmf/:controller", adminCtrl.Index)
	r.POST("/vuecmf/:controller/", adminCtrl.Index)
	r.POST("/vuecmf/:controller/index", adminCtrl.Index)

	//保存单条数据

}