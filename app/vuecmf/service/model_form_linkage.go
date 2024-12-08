//+----------------------------------------------------------------------
// | Copyright (c) 2024 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: tulihua2004@126.com
// +----------------------------------------------------------------------

package service

import (
	"github.com/vuecmf/vuecmf-go/v3/app/vuecmf/model"
	"sync"
)

// ModelFormLinkageService modelFormLinkage服务结构
type ModelFormLinkageService struct {
	*BaseService
}

var modelFormLinkageOnce sync.Once
var modelFormLinkage *ModelFormLinkageService

// ModelFormLinkage 获取modelFormLinkage服务实例
func ModelFormLinkage() *ModelFormLinkageService {
	modelFormLinkageOnce.Do(func() {
		modelFormLinkage = &ModelFormLinkageService{
			BaseService: &BaseService{
				"model_form_linkage",
				&model.ModelFormLinkage{},
				&[]model.ModelFormLinkage{},
				[]string{""},
			},
		}
	})
	return modelFormLinkage
}
