package route

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vuecmf/vuecmf-go/v3/app"
	"github.com/vuecmf/vuecmf-go/v3/app/vuecmf/controller"
	"github.com/vuecmf/vuecmf-go/v3/app/vuecmf/middleware"
	"net/http"
	"strings"
	"time"
)

// Route 路由
type Route struct {
	Path       string                 // 路由路径
	Controller controller.IController //控制器实例
}

// RoutesGroup 路由分组
type RoutesGroup struct {
	GroupName string
	Get       []Route
	Post      []Route
}

// MiddlewareGroup 中间件分组
type MiddlewareGroup struct {
	GroupName  string
	Middleware func(ctx *gin.Context)
}

// getMiddleware 获取中间件
func getMiddleware(middlewareGroup []MiddlewareGroup, groupName string) func(ctx *gin.Context) {
	for _, v := range middlewareGroup {
		if v.GroupName == groupName {
			return v.Middleware
		}
	}
	return nil
}

// Register 注册路由
func Register(eng *gin.Engine, cfg *app.Conf, middlewareGroup []MiddlewareGroup, routesGroup []RoutesGroup) {

	//表单上传文件最大5M
	eng.MaxMultipartMemory = int64(cfg.Upload.AllowFileSize) << 20

	//上传目录 静态文件服务
	eng.StaticFS("/uploads", http.Dir(cfg.Upload.Dir))

	//静态文件目录（css、js、gif、jpg等）
	eng.StaticFS("/static", http.Dir(cfg.StaticDir))

	//加载模板目录
	eng.LoadHTMLGlob("views/**/**/*")

	//跨域设置
	if cfg.CrossDomain.Enable {
		allowOrigins := strings.Split(strings.Replace(cfg.CrossDomain.AllowedOrigin, " ", "", -1), ",")
		var newAllowOrigins []string
		for _, v := range allowOrigins {
			newAllowOrigins = append(newAllowOrigins, strings.Trim(v, "/"))
		}
		eng.Use(cors.New(cors.Config{
			AllowOrigins:     newAllowOrigins,
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "AccessToken", "X-CSRF-Token", "Authorization", "token"},
			ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
	}

	//获取所有中间件
	eng.Use(middleware.GetMiddleWare())

	appMiddleware := getMiddleware(middlewareGroup, "/")
	if appMiddleware != nil {
		eng.Use(appMiddleware)
	}

	sysRoutesGroup := config()

	//注册Get路由请求
	for _, group := range sysRoutesGroup {
		appGroup := eng.Group(group.GroupName)
		{
			for _, route := range group.Get {
				appGroup.GET(route.Path, route.Controller.Before(), route.Controller.Action(), route.Controller.After())
			}
		}
	}

	//注册用户自定义路由Get请求
	for _, group := range routesGroup {
		mw := getMiddleware(middlewareGroup, group.GroupName)
		appGroup := eng.Group(group.GroupName)
		if mw != nil {
			appGroup.Use(mw)
		}
		{
			for _, route := range group.Get {
				appGroup.GET(route.Path, route.Controller.Before(), route.Controller.Action(), route.Controller.After())
			}
		}
	}

	//注册Post路由请求
	for _, group := range sysRoutesGroup {
		appGroup := eng.Group(group.GroupName)
		{
			for _, route := range group.Post {
				appGroup.POST(route.Path, route.Controller.Before(), route.Controller.Action(), route.Controller.After())
			}
		}
	}

	//注册用户自定义路由POST请求
	for _, group := range routesGroup {
		mw := getMiddleware(middlewareGroup, group.GroupName)
		appGroup := eng.Group(group.GroupName)
		if mw != nil {
			appGroup.Use(mw)
		}
		{
			for _, route := range group.Post {
				appGroup.POST(route.Path, route.Controller.Before(), route.Controller.Action(), route.Controller.After())
			}
		}
	}
}
