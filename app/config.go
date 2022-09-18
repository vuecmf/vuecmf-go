package app

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type upload struct {
	FileSize int    `yaml:"file_size"`
	FileExt  string `yaml:"file_ext"`
	Dir      string `yaml:"dir"`
	Url      string `yaml:"url"`
}

type AppConfig struct {
	Upload *upload
}

var appConf *AppConfig

// Config 读取app应用配置
func Config() *AppConfig {
	confContent, err := os.Open("config/app.yaml")
	if err != nil {
		log.Fatal("无法读取应用配置文件app.yaml")
	}

	err = yaml.NewDecoder(confContent).Decode(&appConf)
	if err != nil {
		log.Fatal("数据库配置文件解析错误！")
	}
	return appConf
}
