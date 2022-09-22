// Package service
//+----------------------------------------------------------------------
// | Copyright (c) 2022 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: vuecmf <tulihua2004@126.com>
// +----------------------------------------------------------------------
package service

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vuecmf/vuecmf-go/app"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/helper"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// uploadService upload服务结构
type uploadService struct {
	*baseService
}

var upload *uploadService

// Upload 获取upload服务实例
func Upload() *uploadService {
	if upload == nil {
		upload = &uploadService{}
	}
	return upload
}

type uploadRuleRow struct {
	RuleType  string
	RuleValue string
	ErrorTips string
}

// GetFileMimeType 获取上传文件的MIME类型
func (ser *uploadService) GetFileMimeType(fileHeader *multipart.FileHeader) (string, error) {
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
func (ser *uploadService) UploadFile(fieldName string, ctx *gin.Context) (interface{}, error) {
	var uploadRules []*uploadRuleRow

	var fileSize int
	var fileExt string
	//var isFile bool
	//var isImage bool
	var fileMime string

	fileHeader, err := ctx.FormFile("file")

	if err != nil {
		return nil, errors.New(fieldName + "|上传异常：" + err.Error())
	}

	currentFileMime, err2 := ser.GetFileMimeType(fileHeader)

	if err2 != nil {
		return nil, errors.New(fieldName + "|上传异常：" + err2.Error())
	}

	db.Table(ns.TableName("model_form_rules")+" vmfr").Select("rule_type, rule_value, error_tips").
		Joins("left join "+ns.TableName("model_form")+" vmf on vmfr.model_form_id = vmf.id").
		Joins("left join "+ns.TableName("model_field")+" vmf2 on vmf.model_field_id = vmf2.id").
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
				//isImage = true
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
		fileSize = config.Upload.AllowFileSize
	}

	if fileExt == "" {
		fileExt = config.Upload.AllowFileType
	}

	if fileMime == "" {
		fileMime = config.Upload.AllowFileMime
	}

	if helper.InSlice(currentFileMime, strings.Split(fileMime, ",")) == false {
		return nil, errors.New(fieldName + "|上传异常：不支持该文件类型 " + currentFileMime)
	}

	uploadUrl := config.Upload.Url
	fileName := fileHeader.Filename
	currentFileExt := helper.GetFileExt(fileName)

	if helper.InSlice(currentFileExt, strings.Split(fileExt, ",")) == false {
		return nil, errors.New(fieldName + "|上传异常：不支持该文件类型 " + currentFileExt)
	}

	currentFileBaseName := helper.GetFileBaseName(fileName)
	codeByte := md5.Sum([]byte(currentFileBaseName))
	newFileName := fmt.Sprintf("%x", codeByte)
	currentTime := time.Now().Format("20060102")

	saveDir := config.Upload.Dir + "/" + currentTime + "/"

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

	if  config.Water.Enable == true {
		fontList := []app.FontInfo{ config.Water.Conf }
		err = helper.Img().FontWater(dst, fontList)
		if err != nil {
			return nil, errors.New(fieldName + "|上传异常：添加水印失败！" + err.Error())
		}
	}

	helper.Img().Resize(dst)

	var res = make(map[string]string)
	res["field_name"] = fieldName
	res["url"] = uploadUrl + dst

	return res, nil

}
