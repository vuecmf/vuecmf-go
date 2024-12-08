# 欢迎使用VueCMF

VueCMF内容管理快速开发框架v3(Go版后端API)

## VueCMF是什么？
VueCMF是一款完全开源免费的内容管理快速开发框架。采用前后端分离模式搭建，v3版本前端使用vue3、Element Plus和TypeScript构建，后端API的PHP版基于ThinkPHP8开发，Go版基于Gin开发。可用于快速开发CMS、CRM、WMS、OMS、ERP等管理系统，开发简单、高效易用，极大减少系统的开发周期和研发成本！甚至不用写一行代码使用VueCMF就能设计出功能强大的后台管理系统。

VueCMF开发框架主要有以下功能：

+ 系统授权（多级管理员、多级角色、多级权限）
+ 应用管理
+ 模型配置（字段、索引、动作、表单）
+ 菜单配置

## 使用文档

+ [使用手册(http://www.vuecmf.com)](http://www.vuecmf.com/)


## 前端：
github https://github.com/vuecmf/vuecmf-web

gitee https://gitee.com/emei/vuecmf-web


## 后端：

#### PHP版:
github https://github.com/vuecmf/vuecmf-php

gitee https://gitee.com/emei/vuecmf-php

#### Go版:
github https://github.com/vuecmf/vuecmf-go

gitee https://gitee.com/emei/vuecmf-go


## 环境要求
* MySQL >= 5.7
* Go >= 1.19

## 安装govuecmf命令行工具

~~~
go install github.com/vuecmf/vuecmf-go/v3/tools/govuecmf
~~~

注意：**以下操作均在命令行中执行**

## 安装

创建新项目

~~~
mkdir myproject
cd myproject
govuecmf init myproject
~~~


## 初始化数据

修改config/database.yaml文件中数据库连接配置

然后执行如下操作，进行数据初始化

```
govuecmf migrate init
```
更多命令操作，可执行如下，查看帮助
```
govuecmf -h
```

## 调试与编译
调试
~~~
go run .
~~~
编译
~~~
go build
~~~

## 启动项目
直接执行已编译好的可执行文件即可
~~~
./myproject
~~~




