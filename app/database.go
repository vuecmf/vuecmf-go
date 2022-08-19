// Package app
//+----------------------------------------------------------------------
// | Copyright (c) 2022 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: vuecmf <tulihua2004@126.com>
// +----------------------------------------------------------------------
package app

import (
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

// databaseConf 数据库配置
type databaseConf struct {
	Type            string `yaml:"type"`
	Host            string `yaml:"host"`
	Port            string `yaml:"port"`
	User            string `yaml:"user"`
	Password        string `yaml:"password"`
	Database        string `yaml:"database"`
	Charset         string `yaml:"charset"`
	Prefix          string `yaml:"prefix"`
	MaxIdleConnNums int    `yaml:"max_idle_conn_nums"` //默认
	MaxOpenConnNums int    `yaml:"max_open_conn_nums"`
	ConnMaxLifetime int64  `yaml:"conn_max_lifetime"`
	Debug           bool   `yaml:"debug"`
}

// database database结构
type database struct {
	db *gorm.DB
}

var conf = make(map[string]databaseConf)

// Connect 连接数据库
func (conn *database) connect(confName string) *gorm.DB {
	_, isExist := conf[confName]
	if isExist {
		return conn.db
	}

	confContent, err := os.Open("config/database.yaml")
	if err != nil {
		log.Fatal("无法读取数据库配置文件database.yaml")
	}

	err = yaml.NewDecoder(confContent).Decode(&conf)
	if err != nil {
		log.Fatal("数据库配置文件解析错误！")
	}

	cfg, ok := conf[confName]
	if ok == false {
		log.Fatal("数据库配置（" + confName + "）不存在")
	}

	dsn := cfg.User + ":" + cfg.Password + "@tcp(" + cfg.Host +
		":" + cfg.Port + ")/" + cfg.Database + "?charset=" + cfg.Charset +
		"&parseTime=True&loc=Local"

	db, err2 := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   cfg.Prefix, // 表前缀
			SingularTable: true,       // 使用单数表名
		},
	})

	if err2 != nil {
		log.Fatal("数据库连接失败！")
	}

	if cfg.Debug {
		db = db.Debug()
	}

	sqlDB, err3 := db.DB()
	if err3 != nil {
		log.Fatal("获取SQL DB 失败")
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConnNums)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConnNums)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime * int64(time.Minute)))

	conn.db = db
	return conn.db
}

var db *database

// Db 获取数据库连接
//    参数：confName 数据库配置名称
func Db(confName string) *gorm.DB {
	if db == nil {
		db = &database{}
	}
	return db.connect(confName)
}
