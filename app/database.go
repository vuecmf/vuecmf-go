package app

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"os"
)

// DatabaseConf 数据库配置
type DatabaseConf struct {
	Type string `yaml:"type"`
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	User string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Charset string `yaml:"charset"`
	Prefix string `yaml:"prefix"`
	Debug bool `yaml:"debug"`
}

// database database结构
type database struct {
	Db *gorm.DB
	Conf DatabaseConf
}

// Connect 连接数据库
func (conn *database) Connect(confName string) *database {
	confContent, err := os.Open("config/database.yaml")
	if err != nil {
		panic("无法读取数据库配置文件database.yaml")
	}

	databaseConf := make(map[string]DatabaseConf)

	err = yaml.NewDecoder(confContent).Decode(&databaseConf)
	if err != nil {
		panic("数据库配置文件解析错误！")
	}

	conf, ok := databaseConf[confName]
	if ok == false {
		panic("数据库配置（"+ confName +"）不存在")
	}

	dsn := conf.User + ":" + conf.Password + "@tcp(" + conf.Host +
		":" + conf.Port + ")/" + conf.Database + "?charset=" + conf.Charset +
		"&parseTime=True&loc=Local"

  	db, err2 := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: conf.Prefix,   	// 表前缀
			SingularTable: true, 					// 使用单数表名
		},
	})

	if err2 != nil {
		panic("数据库连接失败！")
	}

	if conf.Debug {
		db = db.Debug()
	}

	conn.Conf = conf
	conn.Db = db
	return conn
}

// Table 查询表
func (conn *database) Table(tableName string) *database {
	conn.Db = conn.Db.Table(conn.Conf.Prefix + tableName)
	return conn
}

// Paginate 分页 page=当前页码，options = 每页显示条数
func (conn *database) Paginate(page int, options ...int) (db *gorm.DB) {
	fmt.Println("page=", page, " len=", len(options))

	pageSize := 20  //默认每页显示20条
	if len(options) > 0 {
		pageSize = options[0]
	}

	offset := (page - 1) * pageSize

	db = conn.Db.Offset(offset).Limit(pageSize)
	return
}

var db *database

// Db 获取数据库连接
//    参数：confName 数据库配置名称
func Db(confName string) *database{
	//if db == nil {
		db = &database{}
	//}
	return db.Connect(confName)
}