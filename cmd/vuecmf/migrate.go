package main

import (
	"errors"
	"fmt"
	"github.com/vuecmf/vuecmf-go/app"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/model"
)


// Migrate 数据库迁移
func Migrate(aType string){
	var err error
	switch aType {
	case "init":
		err = initDb()
	case "up":
		err = up()
	case "down":
		err = down()
	default:
		if aType == "" {
			err = errors.New("参数-t不能为空")
		} else {
			err = errors.New("不支持的选项类型！仅支持init|up|down")
		}
	}

	if err != nil {
		fmt.Println("数据库迁移操作执行失败！" + err.Error())
	} else {
		fmt.Println("恭喜您，数据库迁移操作执行成功! ^_^ ")
	}

}

// initDb 数据库初始化
func initDb() error {
	db := app.Db("demo")
	err := db.AutoMigrate(&model.Admin{})



	if err != nil {
		fmt.Println("数据库初始化失败！" + err.Error())
	}


	return nil
}

// up 数据库升级
func up() error {
	return nil
}

// down 数据库回滚
func down() error {
	return nil
}