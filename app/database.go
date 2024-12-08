//+----------------------------------------------------------------------
// | Copyright (c) 2024 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: tulihua2004@126.com
// +----------------------------------------------------------------------

package app

import (
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var Db *gorm.DB     //数据库连接实例
var NS schema.Namer //数据库名称服务
var Cfg *Conf       // 配置信息实例
var once sync.Once

func init() {
	//初始化项目时，不返回数据库连接
	args := os.Args
	if (len(args) > 1 && (strings.ToLower(args[1]) == "init" || strings.ToLower(args[1]) == "-h")) ||
		(len(args) == 1 && (strings.HasSuffix(args[0], "vuecmf.exe") || strings.HasSuffix(args[0], "vuecmf"))) {
		return
	}

	Cfg = Config()
	if Cfg != nil {
		Db = Connect(strings.ToLower(Cfg.Env))
	}

	if Db != nil {
		NS = Db.NamingStrategy
	}
}

// connConf 数据库配置
type connConf struct {
	Type                   string `yaml:"type"`                     //数据库类型
	Host                   string `yaml:"host"`                     //数据库地址
	Port                   string `yaml:"port"`                     //端口
	User                   string `yaml:"user"`                     //用户名
	Password               string `yaml:"password"`                 //密码
	Database               string `yaml:"database"`                 //数据库名称
	Charset                string `yaml:"charset"`                  //字符编码
	Prefix                 string `yaml:"prefix"`                   //表前缀
	MaxIdleConnNums        int    `yaml:"max_idle_conn_nums"`       //设置空闲连接池中连接的最大数量
	MaxOpenConnNums        int    `yaml:"max_open_conn_nums"`       //设置打开数据库连接的最大数量
	ConnMaxLifetime        int64  `yaml:"conn_max_lifetime"`        //设置了连接可复用的最大时间，单位：分钟
	SkipDefaultTransaction bool   `yaml:"skip_default_transaction"` //是否禁用默认事务, 若禁用默认事务 只在需要时使用事务 性能会提升30%+
	Debug                  bool   `yaml:"debug"`                    //是否开启调试模式，开启后，控制台会打印所执行的SQL语句
}

// DatabaseConf 数据库配置
type DatabaseConf struct {
	Connect map[string]connConf `yaml:"connect"`
}

var DbCfg *DatabaseConf

// DbConf 读取数据库配置信息
func DbConf() *DatabaseConf {
	confContent, err := os.Open("config/database.yaml")
	if err != nil {
		log.Fatal("无法读取数据库配置文件database.yaml")
	}

	err = yaml.NewDecoder(confContent).Decode(&DbCfg)
	if err != nil {
		log.Fatal("数据库配置文件解析错误！")
	}
	return DbCfg
}

var conn = make(map[string]*gorm.DB)

// Connect 连接数据库
//
//	参数：
//	confName 数据库配置名称
func Connect(confName string) *gorm.DB {
	once.Do(func() {
		DbCfg = DbConf()
	})

	_, isExist := DbCfg.Connect[confName]
	if isExist && conn[confName] != nil {
		return conn[confName]
	}

	cfg, ok := DbCfg.Connect[confName]
	if ok == false {
		log.Fatal("数据库配置（" + confName + "）不存在")
	}

	dsn := cfg.User + ":" + cfg.Password + "@tcp(" + cfg.Host +
		":" + cfg.Port + ")/" + cfg.Database + "?charset=" + cfg.Charset +
		"&parseTime=True&loc=Local"

	var err2 error
	conn[confName], err2 = gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: cfg.SkipDefaultTransaction, //是否禁用默认事务
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   cfg.Prefix, // 表前缀
			SingularTable: true,       // 使用单数表名
		},
	})

	if err2 != nil {
		log.Fatal("数据库连接失败！")
	}

	if cfg.Debug {
		conn[confName] = conn[confName].Debug()
	}

	sqlDB, err3 := conn[confName].DB()
	if err3 != nil {
		log.Fatal("获取SQL DB 失败")
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConnNums)

	// SetMaxOpenConns 设置打开数据库连接的最大数量
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConnNums)

	// SetConnMaxLifetime 设置了连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime * int64(time.Minute)))

	return conn[confName]
}
