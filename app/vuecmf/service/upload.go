//+----------------------------------------------------------------------
// | Copyright (c) 2024 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: tulihua2004@126.com
// +----------------------------------------------------------------------

package service

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vuecmf/vuecmf-go/v3/app"
	"github.com/vuecmf/vuecmf-go/v3/app/vuecmf/helper"
	"github.com/vuecmf/vuecmf-go/v3/app/vuecmf/model"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// UploadService upload服务结构
type UploadService struct {
	*BaseService
}

var uploadOnce sync.Once
var upload *UploadService

// Upload 获取upload服务实例
func Upload() *UploadService {
	uploadOnce.Do(func() {
		upload = &UploadService{
			BaseService: &BaseService{
				"upload",
				&model.Upload{},
				&[]model.Upload{},
				[]string{""},
			},
		}
	})
	return upload
}

type UploadRuleRow struct {
	RuleType  string
	RuleValue string
	ErrorTips string
}

// GetFileMimeType 获取上传文件的MIME类型
//
//	参数：
//		fileHeader 文件头信息
func (svc *UploadService) GetFileMimeType(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}

	buf := make([]byte, 512)
	_, err = file.Read(buf)
	if err != nil {
		return "", errors.New("文件读取失败")
	}

	return http.DetectContentType(buf), nil
}

// UploadFile 文件上传
//
//	参数：
//		fieldName 字段名
//		ctx gin.Context上下文
func (svc *UploadService) UploadFile(fieldName string, ctx *gin.Context) (map[string]string, error) {
	var uploadRules []*UploadRuleRow

	var fileSize int
	var fileExt string
	//var isFile bool
	var isImage bool
	var fileMime string

	fileHeader, err := ctx.FormFile("file")

	if err != nil {
		return nil, errors.New(fieldName + "|上传异常：" + err.Error())
	}

	currentFileMime, err2 := svc.GetFileMimeType(fileHeader)

	if err2 != nil {
		return nil, errors.New(fieldName + "|上传异常：" + err2.Error())
	}

	DbTable("model_form_rules", "vmfr").Select("rule_type, rule_value, error_tips").
		Joins("left join "+TableName("model_form")+" vmf on vmfr.model_form_id = vmf.id").
		Joins("left join "+TableName("model_field")+" vmf2 on vmf.model_field_id = vmf2.id").
		Where("rule_type in ('file','image','fileExt','fileMime','fileSize')").
		Where("vmfr.status = 10").
		Where("vmf.status = 10").
		Where("vmf2.status = 10").
		Where("vmf2.field_name = ?", fieldName).
		Find(&uploadRules)

	if len(uploadRules) != 0 {
		for _, row := range uploadRules {
			switch row.RuleType {
			case "file":
				//isFile = true
			case "image":
				isImage = true
			case "fileExt":
				fileExt = row.RuleValue
			case "fileSize":
				fileSize, _ = strconv.Atoi(row.RuleValue)
			case "fileMime":
				fileMime = row.RuleValue
			}
		}
	}

	if fileSize == 0 {
		fileSize = app.Cfg.Upload.AllowFileSize
	}

	if fileExt == "" {
		fileExt = app.Cfg.Upload.AllowFileType
	}

	if fileMime == "" {
		fileMime = app.Cfg.Upload.AllowFileMime
	}

	//文件类型检测
	if helper.InSlice(currentFileMime, strings.Split(fileMime, ",")) == false {
		return nil, errors.New(fieldName + "|上传异常：不支持该文件类型 " + currentFileMime)
	}

	uploadUrl := app.Cfg.Upload.Url
	fileName := fileHeader.Filename
	currentFileExt := helper.GetFileExt(fileName)

	//文件扩展名检测
	if helper.InSlice(currentFileExt, strings.Split(fileExt, ",")) == false {
		return nil, errors.New(fieldName + "|上传异常：不支持该文件类型 " + currentFileExt)
	}

	currentFileBaseName := helper.GetFileBaseName(fileName)
	codeByte := md5.Sum([]byte(currentFileBaseName))
	newFileName := fmt.Sprintf("%x", codeByte)
	currentTime := time.Now().Format("20060102")

	uid := strconv.Itoa(helper.InterfaceToInt(app.Request(ctx).GetCtxVal("uid")))

	saveDir := app.Cfg.Upload.Dir + "/" + uid + "/" + currentTime + "/"

	_, err = os.Stat(saveDir)
	if err != nil {
		err = os.MkdirAll(saveDir, 0666)
		if err != nil {
			return nil, errors.New(fieldName + "|上传异常：创建文件夹失败！" + err.Error())
		}
	}

	dst := saveDir + newFileName + "." + currentFileExt
	err = ctx.SaveUploadedFile(fileHeader, dst)

	if err != nil {
		return nil, errors.New(fieldName + "|上传异常：文件上传失败！" + err.Error())
	}

	if isImage == true {
		//缩放图像文件
		if app.Cfg.Upload.Image.ResizeEnable == true {
			err = helper.Img().Load(dst).Resize(
				dst,
				app.Cfg.Upload.Image.ImageWidth,
				app.Cfg.Upload.Image.ImageHeight,
				app.Cfg.Upload.Image.KeepRatio,
				app.Cfg.Upload.Image.FillBackground,
				app.Cfg.Upload.Image.CenterAlign,
				app.Cfg.Upload.Image.Crop)
		}

		//给图像添加水印
		if app.Cfg.Water.Enable == true {
			fontList := []app.FontInfo{app.Cfg.Water.Conf}
			err = helper.Img().Load(dst).FontWater(fontList)
			if err != nil {
				return nil, errors.New(fieldName + "|上传异常：添加水印失败！" + err.Error())
			}
		}
	}

	var res = make(map[string]string)
	res["field_name"] = fieldName
	res["url"] = uploadUrl + dst
	res["path"] = dst
	res["file_name"] = fileName

	return res, err

}
