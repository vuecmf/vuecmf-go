package main

import (
	"flag"
	"fmt"
	"github.com/vuecmf/vuecmf-go/v3/app"
	"github.com/vuecmf/vuecmf-go/v3/tools/govuecmf/cmd"
	"os"
	"strconv"
	"strings"
)

func PrintError(msg string) {
	fmt.Println("\033[31m", msg, "\033[0m")
}

func PrintSuccess(msg string) {
	fmt.Println("\033[32m", msg, "\033[0m")
}

func main() {
	args := os.Args

	//自定义帮助文档
	flag.Usage = func() {
		fmt.Println(`
vuecmf工具使用介绍：
1、项目初始化
  用法：govuecmf init [project]
  选项值：
    project 项目名称
  例如：
    govuecmf init myweb

2、数据库迁移
  用法： govuecmf migrate [type] -v [version]
  选项值：
    type 要执行的操作，可选值有 init(数据库初始化)、up(升级数据库)、down(回滚数据库)
    version 当执行down操作时，需指定回滚的版本号
  例如：
    初始化数据库：
    govuecmf migrate init
    升级数据库：
    govuecmf migrate up
    回滚数据库：
    govuecmf migrate down -v 20221002101042

3、代码生成
  用法： govuecmf make [type] -n [name] -m [app_module]
  选项值：
    type 要生成模块类型，可选值有 app(应用模块)、controller(控制器)、model(模型)、service(服务)
    name 要生成类型对应的名称
    app_module 应用模块名称
  例如：
    生成demo应用模块：
    govuecmf make app -n demo
    在demo模块下生成user控制器：
    govuecmf make controller -n user -m demo
    在demo模块下生成user模型：
    govuecmf make model -n user -m demo
    在demo模块下生成user服务：
    govuecmf make service -n user -m demo

4、查看vuecmf框架当前版本号
  govuecmf

5、查看帮助文档
  govuecmf -h

    `)

	}

	if len(args) == 1 {
		fmt.Println("vuecmf version v" + app.Version)
	} else {
		param1 := strings.ToLower(args[1])
		switch param1 {
		case "-h": // 显示帮助文档
			flag.Usage()

		case "init": // 初始化项目
			if len(args) == 2 {
				PrintError("项目初始化失败！请输入项目名称")
				return
			}
			cmd.InitProject(args[2])

		case "migrate": //
			if len(args) == 2 {
				PrintError("操作执行失败！缺少数据库迁移操作, 仅支持init|up|down 例如: govuecmf migrate up, 更多请输入govuecmf -h")
				return
			}

			version := 0
			if args[2] == "down" {
				if len(args) == 5 && args[3] == "-v" {
					version, _ = strconv.Atoi(args[4])
				} else {
					PrintError("操作执行失败！缺少参数-v, 例如: govuecmf migrate down -v 20221002101042, 更多请输入govuecmf -h")
					return
				}
			}
			msg, err := cmd.Migrator(args[2], version)
			if err != nil {
				PrintError(err.Error())
			} else {
				if args[2] == "init" {
					PrintSuccess("管理员表中写入的一条初始数据 password(密码) = 123456")
				}

				PrintSuccess("恭喜您，" + msg + "操作执行成功! ^_^ ")
			}

		case "make":
			if len(args) == 2 {
				PrintError("操作执行失败！缺少要生成的类型, 仅支持app|controller|model|service 例如: govuecmf make app -n demo, 更多请输入govuecmf -h")
				return
			}

			if len(args) != 5 && len(args) != 7 {
				PrintError("操作执行失败！缺少参数, 例如: govuecmf make app -n demo, 更多请输入govuecmf -h")
				return
			}

			name := ""
			module := ""

			if len(args) == 5 {
				if args[3] == "-n" {
					name = strings.Trim(args[4], " ")
				}
				if name == "" {
					PrintError("操作执行失败！缺少输入参数-n的值, 例如: govuecmf make app -n demo, 更多请输入govuecmf -h")
					return
				}
			}

			if len(args) == 7 {
				if args[3] == "-n" {
					name = strings.Trim(args[4], " ")
				} else if args[3] == "-m" {
					module = strings.Trim(args[4], " ")
				}

				if args[5] == "-m" {
					module = strings.Trim(args[6], " ")
				} else if args[5] == "-n" {
					name = strings.Trim(args[6], " ")
				}

				if name == "" {
					PrintError("操作执行失败！缺少输入参数-n的值, 例如: govuecmf make app -n demo, 更多请输入govuecmf -h")
					return
				}

				if module == "" {
					PrintError("操作执行失败！缺少输入参数-m的值, 例如: govuecmf make app -n demo, 更多请输入govuecmf -h")
					return
				}

			}

			err := cmd.Make(args[2], name, module)
			if err != nil {
				PrintError(err.Error())
			} else {
				PrintSuccess("操作执行成功! ^_^ ")

			}

		default:
			PrintError("不支持的动作类型！查看帮助文档, 请执行 govuecmf -h")
		}

	}
}
