package app

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type upload struct {
	AllowFileSize  int    `yaml:"allow_file_size"`  //允许上传的最大文件，单位M
	AllowFileType  string `yaml:"allow_file_type"`  //支持上传的文件类型
	AllowFileMime  string `yaml:"allow_file_mime"`  //支持上传文件的MIME类型
	Dir            string `yaml:"dir"`              //文件保存目录
	Url            string `yaml:"url"`              //文件访问链接域名
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
