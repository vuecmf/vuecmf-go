package route

import (
	"github.com/gin-gonic/gin"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/helper"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/middleware"
	"net/http"
	"reflect"
	"strings"
)

//路由映射列表
var routes = make(map[string]map[string]map[string]reflect.Value)

// Register 注册控制器路由
func Register(ctrl interface{}, method, appName string) {
	method = strings.ToUpper(method)
	appName = helper.CamelToUnder(appName)
	methodArr := strings.Split(method, "|")

	refV := reflect.ValueOf(ctrl)
	actionList := refV.NumMethod()
	ctrlName := reflect.TypeOf(ctrl).Elem().Name()
	ctrlName = helper.CamelToUnder(ctrlName)

	for _, methodName := range methodArr {
		if routes[appName] == nil {
			routes[appName] = map[string]map[string]reflect.Value{}
		}

		if routes[appName][ctrlName] == nil {
			routes[appName][ctrlName] = map[string]reflect.Value{}
		}

		for i := 0; i < actionList; i++ {
			actionName := refV.Type().Method(i).Name
			actionName = helper.CamelToUnder(actionName)
			routes[appName][ctrlName][actionName+":"+methodName] = refV.Method(i)
		}
	}
}

var Engine *gin.Engine

// InitRoute 初始化路由列表
func InitRoute(eng *gin.Engine) {
	Engine = eng
	//表单上传文件最大5M
	Engine.MaxMultipartMemory = 5 << 20

	//上传目录 静态文件服务
	Engine.StaticFS("/uploads", http.Dir("uploads"))

	//静态文件目录（css、js、gif、jpg等）
	Engine.StaticFS("/static", http.Dir("static"))

	//加载模板目录
	Engine.LoadHTMLGlob("views/**/**/*")

	//获取所有中间件
	mw := middleware.GetMiddleWares()

	for groupName, ctrl := range routes {
		engine := Engine.Group("/" + groupName + "/")

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

				indexUrl := ""
				indexUrl2 := ""

				if arr[0] == "index" {
					indexUrl = "/" + ctrlName + "/"
					indexUrl2 = "/" + ctrlName
				}

				if ctrlName == "index" {
					indexUrl2 = "/"
				}

				if arr[1] == "GET" {
					engine.GET(url, getHandle(method))
					if indexUrl != "" {
						engine.GET(indexUrl, getHandle(method))
						engine.GET(indexUrl2, getHandle(method))
					}
				} else if arr[1] == "POST" {
					engine.POST(url, getHandle(method))
					if indexUrl != "" {
						engine.POST(indexUrl, getHandle(method))
						engine.POST(indexUrl2, getHandle(method))
					}
				}
			}
		}
	}
}

// getHandle 获取路由执行的方法
func getHandle(method reflect.Value) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		args := make([]reflect.Value, 1)
		args[0] = reflect.ValueOf(ctx)
		method.Call(args)
	}
}
