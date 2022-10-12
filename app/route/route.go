package route

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vuecmf/vuecmf-go/app"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/helper"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/middleware"
	"net/http"
	"reflect"
	"strings"
	"time"
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

// InitRoute 初始化路由列表
func InitRoute(eng *gin.Engine) {
	cfg := app.Config()

	//表单上传文件最大5M
	eng.MaxMultipartMemory = int64(cfg.Upload.AllowFileSize) << 20

	//上传目录 静态文件服务
	eng.StaticFS("/uploads", http.Dir("uploads"))

	//静态文件目录（css、js、gif、jpg等）
	eng.StaticFS("/static", http.Dir("static"))

	//加载模板目录
	eng.LoadHTMLGlob("views/**/**/*")

	//跨域设置
	if cfg.CrossDomain.Enable {
		allowOrigins := strings.Split(strings.Replace(cfg.CrossDomain.AllowedOrigin," ","", -1), ",")
		var newAllowOrigins []string
		for _, v := range allowOrigins {
			newAllowOrigins = append(newAllowOrigins, strings.Trim(v, "/"))
		}
		eng.Use(cors.New(cors.Config{
			AllowOrigins:     newAllowOrigins,
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
			AllowHeaders:     []string{"Origin","Content-Length", "Content-Type","AccessToken","X-CSRF-Token", "Authorization", "token"},
			ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
			AllowCredentials: true,
			MaxAge: 12 * time.Hour,
		}))
	}

	//注册中间件及路由
	regRouteAndMiddleware(eng)

}

//regRouteAndMiddleware 注册中间件及路由
func regRouteAndMiddleware(eng *gin.Engine) {
	//获取所有中间件
	mw := middleware.GetMiddleWares()

	for groupName, ctrl := range routes {
		rg := eng.Group("/" + groupName + "/")

		//加入中间件
		if mw[groupName] != nil {
			for _, handle := range mw[groupName] {
				rg.Use(handle)
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
					rg.GET(url, getHandle(method))
					if indexUrl != "" {
						rg.GET(indexUrl, getHandle(method))
						rg.GET(indexUrl2, getHandle(method))
					}
				} else if arr[1] == "POST" {
					rg.POST(url, getHandle(method))
					if indexUrl != "" {
						rg.POST(indexUrl, getHandle(method))
						rg.POST(indexUrl2, getHandle(method))
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
