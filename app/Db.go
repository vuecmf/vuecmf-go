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

// Db Db结构
type Db struct {
	Db *gorm.DB
	Conf DatabaseConf
}

// Connect 连接数据库
func (conn *Db) Connect(confName ...string) (tx *Db) {
	confContent, err := os.Open("config/database.yaml")
	if err != nil {
		panic("无法读取数据库配置文件database.yaml")
	}

	databaseConf := make(map[string]DatabaseConf)

	err = yaml.NewDecoder(confContent).Decode(&databaseConf)
	if err != nil {
		panic("数据库配置文件解析错误！")
	}

	conf := databaseConf["default"]

	conn.Conf = conf
	
	if 0 != len(confName) {
		conf = databaseConf[confName[0]]
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

	conn.Db = db

	tx = conn
	
	return
}

// Table 查询表
func (conn *Db) Table(tableName string) (tx *Db) {
	conn.Db = conn.Db.Table(conn.Conf.Prefix + tableName)
	tx = conn
	return
}

// Paginate 分页 page=当前页码，options = 每页显示条数
func (conn *Db) Paginate(page int, options ...int) (db *gorm.DB) {
	fmt.Println("page=", page, " len=", len(options))

	pageSize := 20  //默认每页显示20条
	if len(options) > 0 {
		pageSize = options[0]
	}

	offset := (page - 1) * pageSize

	db = conn.Db.Offset(offset).Limit(pageSize)
	return
}