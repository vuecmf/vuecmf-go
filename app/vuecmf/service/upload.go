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
	"github.com/vuecmf/vuecmf-go/app/vuecmf/helper"
	"strconv"
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

func (ser *uploadService) UploadFile(fieldName string, ctx *gin.Context) (interface{}, error) {
	var uploadRules []*uploadRuleRow

	var fileSize int
	var fileExt string
	var isFile bool
	var isImage bool
	var fileMime string

	file, err := ctx.FormFile("file")

	if err != nil {
		return nil, errors.New(fieldName + "|上传异常：" + err.Error())
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
				isFile = true
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
		fileSize = config.Upload.FileSize
	}

	if fileExt == "" {
		fileExt = config.Upload.FileExt
	}

	uploadDir := "./" + config.Upload.Dir
	uploadUrl := config.Upload.Url
	fileName := file.Filename
	currentFileExt := helper.GetFileExt(fileName)

	currentFileBaseName := helper.GetFileBaseName(fileName)
	codeByte := md5.Sum([]byte(currentFileBaseName))
	newFileName := fmt.Sprintf("%x", codeByte)
	currentTime := time.Now().Format("20060102")

	dst := uploadDir + "/" + currentTime + "/" + newFileName
	err = ctx.SaveUploadedFile(file, dst)

}
