// Package middleware
//+----------------------------------------------------------------------
// | Copyright (c) 2022 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: vuecmf <tulihua2004@126.com>
// +----------------------------------------------------------------------
package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vuecmf/vuecmf-go/app"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/service"
	"strings"
)

var middlewares = make(map[string]func(ctx *gin.Context))

//GetMiddleWares 获取所有中间件
func GetMiddleWares() map[string]func(ctx *gin.Context) {
	//先取出所有应用列表
	appList := service.AppConfig().GetFullAppList()

	//访问权限验证
	middlewares["auth"] = func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				err2 := errors.New(fmt.Sprintf("%s", err))
				app.Response(ctx).SendFailure("请求失败", service.GetErrMsg(err2), 1003)
				ctx.Abort()
			}
		}()

		path := strings.ToLower(ctx.Request.URL.String())
		tmpArr := strings.Split(path, "?")
		pathArr := strings.Split(tmpArr[0], "/")
		routeApp := "index"
		routeController := "index"
		routeAction := "index"

		switch {
		case len(pathArr) == 2 && pathArr[1] != "":
			routeApp = pathArr[1]
		case len(pathArr) == 3:
			routeApp = pathArr[1]
			if pathArr[2] != "" {
				routeController = pathArr[2]
			}
		case len(pathArr) == 4:
			routeApp = pathArr[1]
			routeController = pathArr[2]
			if pathArr[3] != "" {
				routeAction = pathArr[3]
			}
		case len(pathArr) > 4:
			routeApp = pathArr[1]
			routeController = pathArr[2]
			routeAction = pathArr[3]
		}

		//过滤设置了排除验证的URL
		flag := false
		exclusionUrlArr := strings.Split(strings.ToLower(strings.Replace(appList[routeApp].ExclusionUrl, " ", "", -1)), ",")
		checkUrl := "/" + routeApp + "/" + routeController + "/" + routeAction
		for _, exUrl := range exclusionUrlArr {
			if exUrl == checkUrl || exUrl + "/index" == checkUrl {
				flag = true
				break
			}
		}
		if flag {
			return
		}

		//登录验证
		if appList[routeApp].LoginEnable == 10 {
			token := ctx.Request.Header.Get("token")
			adm, err := service.Admin().IsLogin(token, ctx.ClientIP())

			//权限验证
			if err == nil {
				ctx.Set("is_super", adm.IsSuper)
				if adm.IsSuper == 10 {
					return
				}

				if appList[routeApp].AuthEnable == 10 {
					res, err := service.Auth().Enforcer.Enforce(adm.Username, routeApp, routeController, routeAction)
					if err == nil && res == false {
						err = errors.New("您没有访问权限！")
					}
				}
			}

			if err != nil {
				app.Response(ctx).SendFailure(err.Error(), err, 1003)
				ctx.Abort()
				return
			}
		}
	}

	return middlewares
}
