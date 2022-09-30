package main

import (
	"errors"
	"fmt"
	"github.com/vuecmf/vuecmf-go/app"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/model"
)


// Migrator 数据库迁移
func Migrator(aType string){
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

/**
sqlType := "bigint"
	switch {
	case field.Size <= 8:
		sqlType = "tinyint"
	case field.Size <= 16:
		sqlType = "smallint"
	case field.Size <= 24:
		sqlType = "mediumint"
	case field.Size <= 32:
		sqlType = "int"
	}
 */


// initDb 数据库初始化
func initDb() error {
	db := app.Db("demo")

	//创建管理员表
	err := db.Migrator().DropTable(&model.Admin{})
	if err != nil {
		return errors.New("删除历史表"+ db.NamingStrategy.TableName("admin") +"失败:" + err.Error())
	}
	err = db.Set("gorm:table_options","ENGINE=InnoDB COLLATE=utf8mb4_unicode_ci COMMENT='系统--管理员表'").
		AutoMigrate(&model.Admin{})
	if err != nil {
		return errors.New("创建表"+ db.NamingStrategy.TableName("admin") +"失败:" + err.Error())
	}






	





	return err
}

// up 数据库升级
func up() error {
	return nil
}

// down 数据库回滚
func down() error {
	return nil
}