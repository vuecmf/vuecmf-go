//+----------------------------------------------------------------------
// | Copyright (c) 2023 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: vuecmf <tulihua2004@126.com>
// +----------------------------------------------------------------------

package app

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

// uploadImage 图像文件上传配置
type uploadImage struct {
	ResizeEnable   bool `yaml:"resize_enable"`   //是否缩放图片
	ImageWidth     int  `yaml:"image_width"`     //上传的图片裁切后的宽度
	ImageHeight    int  `yaml:"image_height"`    //上传的图片裁切后的高度
	KeepRatio      bool `yaml:"keep_ratio"`      //是否保持等比例缩放
	FillBackground int  `yaml:"fill_background"` //填充的背景颜色 0 - 255 （R、G、B）的值共一个数值， 0 = 透明背景， 255 = 白色背景
	CenterAlign    bool `yaml:"center_align"`    //是否以图片的中心来进行等比缩放
	Crop           bool `yaml:"crop"`            //是否裁切图片
}

// upload 上传配置
type upload struct {
	AllowFileSize int         `yaml:"allow_file_size"` //允许上传的最大文件，单位M
	AllowFileType string      `yaml:"allow_file_type"` //支持上传的文件类型
	AllowFileMime string      `yaml:"allow_file_mime"` //支持上传文件的MIME类型
	Dir           string      `yaml:"dir"`             //文件保存目录
	Url           string      `yaml:"url"`             //文件访问链接域名
	Image         uploadImage `yaml:"image"`           //图像文件上传配置
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

//跨域配置
type crossDomain struct {
	Enable        bool   `yaml:"enable"`         //是否开启跨域请求
	AllowedOrigin string `yaml:"allowed_origin"` //允许请求的来源， 例如 http://www.vuecmf.com
}

// Conf 应用配置
type Conf struct {
	Module      string       `yaml:"module"`       //项目名称，与go.mod中module保持一致
	Env         string       `yaml:"env"`          //当前运行环境， dev 开发环境，test 测试环境，prod 生产环境
	Debug       bool         `yaml:"debug"`        //是否开启调试模式
	ServerHost  string		 `yaml:"server_host"`  //服务器地址或域名
	ServerPort  string       `yaml:"server_port"`  //服务运行的端口
	CrossDomain *crossDomain `yaml:"cross_domain"` //跨域请求配置
	Upload      *upload      `yaml:"upload"`       //上传配置
	Water       *Water       `yaml:"water"`        //水印配置
	StaticDir   string		 `yaml:"static_dir"`   //静态资源目录
}

var appConf *Conf

// Config 读取app应用配置
func Config() *Conf {
	if appConf != nil {
		return appConf
	}

	confContent, err := os.Open("config/app.yaml")
	if err != nil {
		//log.Fatal("无法读取应用配置文件app.yaml")
		return nil
	}

	err = yaml.NewDecoder(confContent).Decode(&appConf)
	if err != nil {
		log.Fatal("应用配置文件解析错误！")
	}
	return appConf
}
