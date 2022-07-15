package route

import (
	"github.com/gin-gonic/gin"
	"github.com/vuecmf/vuecmf-go/app"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/middleware"
	"reflect"
	"strings"
)

//路由映射列表
var routes = make(map[string]map[string]map[string]reflect.Value)

// Register 注册控制器路由
func Register(ctrl interface{}, method, appName string) {
	method = strings.ToUpper(method)
	appName = app.CamelToUnder(appName)
	methodArr := strings.Split(method, "|")

	refV := reflect.ValueOf(ctrl)
	actionList := refV.NumMethod()
	ctrlName := reflect.TypeOf(ctrl).Elem().Name()
	ctrlName = app.CamelToUnder(ctrlName)

	for _,methodName := range methodArr {
		if routes[appName] == nil {
			routes[appName] = map[string]map[string]reflect.Value{}
			if routes[appName][ctrlName] == nil {
				routes[appName][ctrlName] = map[string]reflect.Value{}
			}
		}

		for i := 0; i < actionList; i++ {
			actionName := refV.Type().Method(i).Name
			actionName = app.CamelToUnder(actionName)
			if routes[appName][ctrlName] != nil {
				routes[appName][ctrlName][actionName + ":" + methodName] = refV.Method(i)
			}
		}
	}
}


// InitRoute 初始化路由列表
func InitRoute(engine *gin.Engine){
	mw := middleware.GetMiddleWares()

	for groupName, ctrl := range routes {
		engine := engine.Group("/" + groupName + "/")

		//加入中间件
		if mw[groupName] != nil {
			for _, handle := range mw[groupName] {
				engine.Use(handle)
			}
		}

		//注册路由
		for ctrlName, action := range ctrl {
			for actionName, method := range action {
				arr := strings.Split(actionName, ":")
				url := "/" + ctrlName + "/" + arr[0]

				if arr[1] == "GET" {
					engine.GET(url, getHandle(method))
				}else if arr[1] == "POST" {
					engine.POST(url, getHandle(method))
				}
			}
		}
	}
}


// getHandle 获取路由执行的方法
func getHandle(method reflect.Value) gin.HandlerFunc{
	return func(ctx *gin.Context) {
		args := make([]reflect.Value, 1)
		args[0] = reflect.ValueOf(ctx)
		method.Call(args)
	}
}