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
	"net/http"
	"strings"
)

var middlewares = make(map[string]map[string]func(ctx *gin.Context))

//GetMiddleWares 获取所有中间件
func GetMiddleWares() map[string]map[string]func(ctx *gin.Context) {
	//先取出所有应用列表
	appList := service.AppConfig().GetFullAppList()
	for _, ai := range appList {
		if middlewares[ai.AppName] == nil {
			middlewares[ai.AppName] = map[string]func(ctx *gin.Context){}
		}

		fmt.Println("DDD==", ai.CrossDomainEnable)

		//访问权限验证
		middlewares[ai.AppName]["auth"] = func(ctx *gin.Context) {
			fmt.Println("eeee==", ai.CrossDomainEnable)

			defer func() {
				if err := recover(); err != nil {
					app.Response(ctx).SendFailure("请求失败", err)
					ctx.Abort()
				}
			}()

			ctx.Header("Access-Control-Allow-Origin", "*")
			ctx.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token, x-token")
			ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
			ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
			ctx.Header("Access-Control-Allow-Credentials", "true")
			if ctx.Request.Method == "OPTIONS" {
				ctx.AbortWithStatus(http.StatusNoContent)
			}



			//跨域请求设置
			if ai.CrossDomainEnable == 10 {
				ctx.Header("Access-Control-Allow-Origin", "*")
				ctx.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token, x-token")
				ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
				ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
				ctx.Header("Access-Control-Allow-Credentials", "true")
				if ctx.Request.Method == "OPTIONS" {
					ctx.AbortWithStatus(http.StatusNoContent)
				}
			}

			path := strings.ToLower(ctx.Request.URL.String())
			pathArr := strings.Split(path, "/")
			routeApp := "index"
			routeController := "index"
			routeAction := "index"

			switch  {
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

			if routeAction == "login" {
				return
			}


			//登录验证
			if ai.LoginEnable == 10 {
				token := ctx.Request.Header.Get("token")
				adm, err := service.Admin(ai.AppName).IsLogin(token, ctx.ClientIP())

				//权限验证
				if err == nil && ai.AuthEnable == 10 {
					ctx.Set("is_super", adm.IsSuper)
					res, err  := service.Auth().Enforcer.Enforce(adm.Username, routeApp, routeController, routeAction)
					if err == nil && res == false {
						err = errors.New("您没有访问权限！")
					}
				}

				if err != nil {
					app.Response(ctx).SendFailure(err.Error(), err)
					ctx.Abort()
					return
				}

			}

		}
	}

	return middlewares
}
