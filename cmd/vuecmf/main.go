package main

import (
	"flag"
	"fmt"
	"github.com/vuecmf/vuecmf-go/app/vuecmf"
)

func main()  {
	action := flag.String("a", "version", "要执行的动作名称，可选值(version|migrate|make)，默认version")
	atype := flag.String("t", "", "要执行动作对应的类型值")
	name := flag.String("n", "", "类型值对应的名称")
	help := flag.Bool("h", false, "帮助文档")

	//自定义帮助文档
	flag.Usage = func() {
		fmt.Println(`
vuecmf工具使用介绍：
1、数据库迁移
  用法: vuecmf -a migrate -t [type]
  选项值：
    type 要执行的操作，可选值有 init(数据库初始化)、up(升级数据库)、down(回滚数据库)
  例如：
    初始化数据库：
    vuecmf -a migrate init
    升级数据库：
    vuecmf -a migrate up
    回滚数据库：
    vuecmf -a migrate down

2、代码生成
  用法: vuecmf -a make -t [type] -n [name]
  选项值：
    type 要生成模块类型，可选值有 app(应用)、controller(控制器)、model(模型)、service(服务)
    name 要生成类型对应的名称
  例如：
    生成应用：
    vuecmf -a make -t app -n demo
    生成控制器：
    vuecmf -a make -t controller -n user
    生成控制器：
    vuecmf -a make -t model -n user
    生成服务：
    vuecmf -a make -t service -n user

3、查看vuecmf框架当前版本号
  vuecmf -a version

4、查看帮助文档
  vuecmf -h

Options: 
    `)
		flag.PrintDefaults()
	}

	// 解析命令行参数
	flag.Parse()

	if *help {
		flag.Usage()
	} else {
		switch *action {
		case "migrate":
			Migrate(*atype)
		case "make":
			Make(*atype, *name)
		case "version":
			fmt.Printf("vuecmf version \"%v\"\n", vuecmf.Version)
		default:
			fmt.Println("不支持的动作类型！查看帮助文档, 请执行 vuecmf -h")
		}
	}

}