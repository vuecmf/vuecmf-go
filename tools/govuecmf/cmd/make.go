package cmd

import (
	"errors"
	"fmt"
	"github.com/vuecmf/vuecmf-go/v3/app"
	"github.com/vuecmf/vuecmf-go/v3/app/vuecmf/service"
	"os"
	"os/exec"
	"strings"
)

// InitProject 初始化项目
func InitProject(projectName string) {
	//执行go mod 初始化项目
	// 检查 go.mod 文件是否存在
	if _, err := os.Stat("go.mod"); err != nil {
		if os.IsNotExist(err) {
			// go.mod 文件不存在，执行 go mod init
			cmd := exec.Command("go", "mod", "init", projectName)
			if err = cmd.Run(); err != nil {
				fmt.Println("项目初始化失败！" + err.Error())
				return
			}
			fmt.Println("项目初始化成功！")
		}
	}

	//生成app目录、home 应用模块
	if err := createApp("home", projectName, true); err != nil {
		fmt.Println(err.Error())
	}

	//生成config目录、相关配置文件
	if err := createConf(projectName); err != nil {
		fmt.Println(err.Error())
	}

	//生成migrations目录、升级数据库SQL的JSON文件
	if err := createDir("migrations", "数据库升级目录", "存放数据库升级SQL的json文件"); err != nil {
		fmt.Println(err.Error())
	}

	//生成static目录， 静态文件
	if err := createDir("static", "静态文件目录", "存放网页相关静态文件（css、js、images等）"); err != nil {
		fmt.Println(err.Error())
	}

	//生成uploads目录、上传文件
	if err := createDir("uploads", "上传文件目录", "存放上传的相关文件"); err != nil {
		fmt.Println(err.Error())
	}

	//生成项目根目录下启动文件
	if err := createRunFile(projectName); err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("创建启动文件main.go完成！")

	//生成vuecmf模块视图
	if err := createVuecmfView(); err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("正在下载项目依赖包，请稍后...")

	//执行go mod tidy 自动下载依赖包
	cmd := exec.Command("go", "mod", "tidy")
	if err := cmd.Run(); err != nil {
		fmt.Println("无法下载相关依赖包，项目初始化失败！" + err.Error())
		tips := `请检查GOPROXY的地址是否可访问，如默认的地址https://proxy.golang.org无法访问，可以设置为镜像地址https://goproxy.cn，
具体在命令行中操作如下：
go env -w GOPROXY=https://goproxy.cn,direct
执行上面操作后，再重新执行项目初始化： 
govuecmf init `
		fmt.Println(tips + projectName)
		return
	}

	fmt.Println(`欢迎使用
 _    __           ________  _________
| |  / /_  _____  / ____/  |/  / ____/
| | / / / / / _ \/ /   / /|_/ / /_    
| |/ / /_/ /  __/ /___/ /  / / __/    
|___/\__,_/\___/\____/_/  /_/_/       
                                      

项目初始化完成！

1、数据库初始化
修改config/database.yaml中数据库连接配置，然后执行下面命令初始化数据
govuecmf migrate init

更多使用，欢迎查看帮助文档
govuecmf -h

2、调试项目
go run .

3、编译项目
go build .

`)
}

// 创建项目配置目录及文件
func createConf(projectName string) error {
	confDir := "config"
	if _, err := os.Stat(confDir); err != nil {
		if err = os.MkdirAll(confDir, 0666); err != nil {
			return errors.New("创建配置目录" + confDir + "失败！" + err.Error())
		}
	}

	if err := createAppConf(confDir, projectName); err != nil {
		return err
	}

	dbFile := `connect:
  #开发环境
  dev:
    type: mysql       #数据库类型
    host: 127.0.0.1   #数据库地址
    port: 3306        #端口
    user: root        #用户名
    password: 123456  #密码
    database: vuecmf  #数据库名称
    charset: utf8     #字符编码
    prefix: vuecmf_   #表前缀
    max_idle_conn_nums: 10  #设置空闲连接池中连接的最大数量
    max_open_conn_nums: 100 #设置打开数据库连接的最大数量
    conn_max_lifetime: 2  #设置了连接可复用的最大时间，单位：分钟
    skip_default_transaction: true #是否禁用默认事务, 若禁用默认事务 只在需要时使用事务 性能会提升30%+
    debug: true   #是否开启调试模式，开启后，控制台会打印所执行的SQL语句

  #测试环境
  test:
    type: mysql       #数据库类型
    host: 127.0.0.1   #数据库地址
    port: 3306        #端口
    user: root        #用户名
    password: 123456  #密码
    database: vuecmf  #数据库名称
    charset: utf8     #字符编码
    prefix: vuecmf_   #表前缀
    max_idle_conn_nums: 10  #设置空闲连接池中连接的最大数量
    max_open_conn_nums: 100 #设置打开数据库连接的最大数量
    conn_max_lifetime: 2  #设置了连接可复用的最大时间，单位：分钟
    skip_default_transaction: true #是否禁用默认事务, 若禁用默认事务 只在需要时使用事务 性能会提升30%+
    debug: true   #是否开启调试模式，开启后，控制台会打印所执行的SQL语句

  #生产环境
  prod:
    type: mysql       #数据库类型
    host: 127.0.0.1   #数据库地址
    port: 3306        #端口
    user: root        #用户名
    password: 123456  #密码
    database: vuecmf  #数据库名称
    charset: utf8     #字符编码
    prefix: vuecmf_   #表前缀
    max_idle_conn_nums: 10  #设置空闲连接池中连接的最大数量
    max_open_conn_nums: 100 #设置打开数据库连接的最大数量
    conn_max_lifetime: 2  #设置了连接可复用的最大时间，单位：分钟
    skip_default_transaction: true #是否禁用默认事务, 若禁用默认事务 只在需要时使用事务 性能会提升30%+
    debug: false   #是否开启调试模式，开启后，控制台会打印所执行的SQL语句
`
	if err := os.WriteFile(confDir+"/database.yaml", []byte(dbFile), 0666); err != nil {
		return errors.New("创建数据库配置文件database.yaml失败！" + err.Error())
	}

	authFile := `#请求的校验参数定义，参数个数必须与验证时的传入的个数一致，对应数据表中 ptype=p 的 v0,v1,v2,v3 字段一一对应
#例如 $res = Enforcer::enforce('lily', 'appname' , 'controller', 'action');
[request_definition]
r = sub, dom, obj, act

#策略定义，对应数据表中 ptype=p 的 v0,v1,v2,v3 字段一一对应， 且数据表中对应的值不能为空
[policy_definition]
p = sub, dom, obj, act

#分组和角色定义，对应数据表中 ptype=g, 且数据表对应的值不能为空
[role_definition]
g = _, _, _

[policy_effect]
e = some(where (p.eft == allow))

#匹配规则，g代表分组或角色，括号里面的必须与数据表中 ptype=g的 v0, v1, v2 字段一一对应
[matchers]
m = g(r.sub, p.sub, r.dom) && r.dom == p.dom && r.obj == p.obj && r.act == p.act`
	if err := os.WriteFile(confDir+"/tauthz-rbac-model.conf", []byte(authFile), 0666); err != nil {
		return errors.New("创建权限配置文件tauthz-rbac-model.conf失败！" + err.Error())
	}

	fmt.Println("创建创建项目配置目录及文件完成！")

	return nil
}

// 创建应用配置文件
func createAppConf(confDir, projectName string) error {
	appFile := `#项目名称
module: "{{.module_name}}" #项目名称，与go.mod中module保持一致

#当前运行环境
env: "dev"  #当前运行环境， dev 开发环境，test 测试环境，prod 生产环境

#创建或更新模型时是否重新生成入口文件main.go, 默认为true
update_main: true

#是否开启调试模式
debug: true

#服务运行的端口
server_port: "8080"

#是否允许跨域请求
cross_domain:
  enable: true  #是否开启跨域请求
  allowed_origin: "http://localhost:8081"  #允许请求的来源; 多个用英文逗号分隔; 例如 http://www.vuecmf.com

#静态资源目录
static_dir: "static"

#上传配置
upload:
  allow_file_size: 5        #允许上传的最大文件，单位M
  allow_file_type: "gif,jpg,jpeg,png,bmp,tif,txt,csv,xls,xlsx,doc,docx,zip,rar,gz,vsd,mdb,pdf,rmvb,flv,mp4,mp3,mpg,wmv,wav,avi,mid,ini,wps,mov,dbx,pst,ram"       #支持上传除图片外的文件类型
  allow_file_mime: "image/gif,image/jpeg,image/png,application/zip,application/octet-stream,text/plain; charset=utf-8,application/pdf"
  dir: "uploads"            #文件保存目录
  url: "http://localhost:8080/"  #文件访问链接域名
  image:
    resize_enable: true #是否缩放图片
    image_width: 600   #上传的图片裁切后的宽度
    image_height: 600  #上传的图片裁切后的高度
    keep_ratio: true #是否保持等比例缩放
    fill_background: 255 #填充的背景颜色 0 - 255 （R、G、B）的值共一个数值， 0 = 透明背景， 255 = 白色背景
    center_align: true #是否以图片的中心来进行等比缩放
    crop: true #是否裁切图片

#水印配置
water:
  enable: false
  water_font: "config/simhei.ttf"   #水印字体文件
  conf:
    size: 72.0 #水印文字大小
    message: "vuecmf" #水印文本内容
    position: 4  #水印位置， 0=左上角， 1=右上角，2=左下角，3=右下角，4=中间
    dx: 0  #文字x轴留白距离
    dy: 0  #文字y轴留白距离
    r: 0   #文字颜色值RGBA中的R值 0 - 255
    g: 0   #文字颜色值RGBA中的G值 0 - 255
    b: 0   #文字颜色值RGBA中的B值 0 - 255
    a: 50  #文字颜色值RGBA中的A值，即透明度 0 - 100

`
	appFile = strings.Replace(appFile, "{{.module_name}}", projectName, -1)
	_ = os.WriteFile(confDir+"/app.yaml", []byte(appFile), 0666)

	return nil
}

// 创建目录
func createDir(dirName string, comment string, readme string) error {
	if _, err := os.Stat(dirName); err != nil {
		if err = os.MkdirAll(dirName, 0666); err != nil {
			return errors.New("创建" + comment + dirName + "失败！" + err.Error())
		}
	}

	if err := os.WriteFile(dirName+"/README.md", []byte(readme), 0666); err != nil {
		return errors.New("创建" + comment + "中文件README.md失败！" + err.Error())
	}

	fmt.Println("创建" + comment + dirName + "完成！")

	return nil
}

// 创建启动文件
func createRunFile(projectName string) error {
	mainFile := `package main

import (
	"github.com/gin-gonic/gin"
	"github.com/vuecmf/vuecmf-go/v3/app"
	sysRoute "github.com/vuecmf/vuecmf-go/v3/app/route"
	"log"
	"{{.module_name}}/app/route"
)

func main() {
	cfg := app.Config()

	if cfg.Debug == false {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.Default()

	//注册路由
	sysRoute.Register(engine, cfg, route.Config())

	err := engine.Run(":" + cfg.ServerPort)
	if err != nil {
		log.Fatal("服务启动失败！", err)
	}

}
`
	mainFile = strings.Replace(mainFile, "{{.module_name}}", projectName, -1)
	if err := os.WriteFile("main.go", []byte(mainFile), 0666); err != nil {
		return errors.New("创建启动文件main.go失败！" + err.Error())
	}
	return nil
}

// 创建vuecmf视图层文件
func createVuecmfView() error {
	//创建视图层目录
	viewDir := "views/vuecmf"
	if _, err := os.Stat(viewDir); err != nil {
		if err = os.MkdirAll(viewDir, 0666); err != nil {
			return errors.New("创建视图层目录views失败！" + err.Error())
		}
	}

	//创建index视图模板
	viewDir = viewDir + "/index"
	if _, err := os.Stat(viewDir); err != nil {
		if err = os.MkdirAll(viewDir, 0666); err != nil {
			return errors.New("创建index视图层目录失败！" + err.Error())
		}
	}
	tpl := `{{ define "vuecmf/index/index.html" }}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>欢迎使用VueCMF快速开发框架</title>
    <style>
        *{ padding: 0; margin: 0; }
        div{ padding: 4px 48px;}
        a{color:#2E5CD5;cursor: pointer;text-decoration: none}
        a:hover{text-decoration:underline; }
        body{ background: #fff; font-family: "Century Gothic","Microsoft yahei"; color: #333;font-size:18px;}
        h1{ font-size: 100px; font-weight: normal; margin-bottom: 12px; }
        p{ line-height: 1.6em; font-size: 42px }
    </style>
</head>
<body>
<div style="padding: 24px 48px;"> <h1>:) </h1><p>{{ .welcome }}<br/></p><span style="font-size:25px;">[ Powered by <a href="http://www.vuecmf.com/" target="_blank">vuecmf</a> ]</span><script src="https://hm.baidu.com/hm.js?74079f71bcec1421dd89f7c08ed21d68"></script></div>
</body>
</html>
{{ end }}
`
	if err := os.WriteFile(viewDir+"/index.html", []byte(tpl), 0666); err != nil {
		return errors.New("创建index视图模板失败！" + err.Error())
	}
	return nil
}

// Make 生成代码文件
func Make(aType, name string, appName string) error {
	var err error
	switch aType {
	case "app":
		err = createApp(name, "", false)
	case "controller":
		err = createController(name, appName)
	case "model":
		err = createModel(name, appName)
	case "service":
		err = createService(name, appName)
	}

	return err
}

// 创建应用模块
func createApp(appName, projectName string, isInit bool) error {
	//创建路由配置文件
	moduleName := projectName
	if !isInit {
		moduleName = app.Config().Module
	}

	err := service.Make().CreateApp(appName, moduleName)
	if err != nil {
		fmt.Println("创建应用"+appName+"失败！", err.Error())
		return err
	}

	fmt.Println("创建应用" + appName + "完成！")

	return nil
}

// 创建应用模块的控制器
func createController(ctrlName string, appName string) error {
	err := service.Make().Controller(ctrlName, appName)
	if err != nil {
		err = errors.New("创建控制器" + ctrlName + "失败！" + err.Error())
	}
	return err
}

// 创建应用模块的模型
func createModel(modelName string, appName string) error {
	err := service.Make().Model(modelName, appName)
	if err != nil {
		err = errors.New("创建模型" + modelName + "失败！" + err.Error())
	}
	return err
}

// 创建应用模块的服务
func createService(serviceName string, appName string) error {
	err := service.Make().Service(serviceName, appName)
	if err != nil {
		err = errors.New("创建服务" + serviceName + "失败！" + err.Error())
	}
	return err
}
