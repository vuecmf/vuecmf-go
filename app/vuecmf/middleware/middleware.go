package middleware

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	"github.com/vuecmf/vuecmf-go/app"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/model"
)

var middlewares = make(map[string]map[string]func(ctx *gin.Context))

//GetMiddleWares 获取所有中间件
func GetMiddleWares() map[string]map[string]func(ctx *gin.Context) {
	if middlewares["vuecmf"] == nil {
		middlewares["vuecmf"] = map[string]func(ctx *gin.Context){}
	}

	//vuecmf应用 登录验证
	middlewares["vuecmf"]["login"] = func(ctx *gin.Context) {
		defer func() {
			resp := app.Response{Context: ctx}
			if err := recover(); err != nil {
				resp.SendFailure("请求失败", err)
				ctx.Abort()
			}
		}()

		token := ctx.Request.Header.Get("token")
		fmt.Println("path=", ctx.Request.URL)
		fmt.Println("token=", token)
		//panic("access deny")
	}

	//vuecmf应用 权限验证
	middlewares["vuecmf"]["auth"] = func(ctx *gin.Context) {
		fmt.Println("开始权限验证")
	}

	return middlewares
}


func Test() {
	db := app.Db{}

	conn := db.Connect()

	// Initialize a Gorm adapter and use it in a Casbin enforcer:
	// The adapter will use an existing gorm.DB instnace.
	a, _ := gormadapter.NewAdapterByDBWithCustomTable(conn.Db, &model.Rules{})
	e, _ := casbin.NewEnforcer("examples/rbac_model.conf", a)

	// Load the policy from DB.
	e.LoadPolicy()

	// Check the permission.
	e.Enforce("alice", "data1", "read")

	// Modify the policy.
	// e.AddPolicy(...)
	// e.RemovePolicy(...)

	// Save the policy back to DB.
	e.SavePolicy()
}