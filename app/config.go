package app

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

// upload 上传配置
type upload struct {
	AllowFileSize int    `yaml:"allow_file_size"` //允许上传的最大文件，单位M
	AllowFileType string `yaml:"allow_file_type"` //支持上传的文件类型
	AllowFileMime string `yaml:"allow_file_mime"` //支持上传文件的MIME类型
	Dir           string `yaml:"dir"`             //文件保存目录
	Url           string `yaml:"url"`             //文件访问链接域名
	ImageWidth    int    `yaml:"image_width"`     //上传的图片裁切后的宽度
	ImageHeight   int    `yaml:"image_height"`    //上传的图片裁切后的高度
}

// FontInfo 水印文字配置信息
type FontInfo struct {
	Size     float64 `yaml:"size"`     //文字大小
	Message  string  `yaml:"message"`  //文字内容
	Position int     `yaml:"position"` //文字存放位置
	Dx       int     `yaml:"dx"`       //文字x轴留白距离
	Dy       int     `yaml:"dy"`       //文字y轴留白距离
	R        uint8   `yaml:"r"`        //文字颜色值RGBA中的R值
	G        uint8   `yaml:"g"`        //文字颜色值RGBA中的G值
	B        uint8   `yaml:"b"`        //文字颜色值RGBA中的B值
	A        uint8   `yaml:"a"`        //文字颜色值RGBA中的A值
}

// Water 图片上传水印设置
type Water struct {
	Enable    bool     `yaml:"enable"`     //是否启用水印
	WaterFont string   `yaml:"water_font"` //水印字体路径
	Conf      FontInfo `yaml:"conf"`       //水印配置
}

// Config 应用配置
type Config struct {
	Upload *upload `yaml:"upload"`
	Water  *Water  `yaml:"water"`
}

var appConf *Config

// Conf 读取app应用配置
func Conf() *Config {
	confContent, err := os.Open("config/app.yaml")
	if err != nil {
		log.Fatal("无法读取应用配置文件app.yaml")
	}

	err = yaml.NewDecoder(confContent).Decode(&appConf)
	if err != nil {
		log.Fatal("应用配置文件解析错误！")
	}
	return appConf
}
