//+----------------------------------------------------------------------
// | Copyright (c) 2023 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: vuecmf <tulihua2004@126.com>
// +----------------------------------------------------------------------

package service

// modelFormLinkageService modelFormLinkage服务结构
type modelFormLinkageService struct {
	*BaseService
}

var modelFormLinkage *modelFormLinkageService

// ModelFormLinkage 获取modelFormLinkage服务实例
func ModelFormLinkage() *modelFormLinkageService {
	if modelFormLinkage == nil {
		modelFormLinkage = &modelFormLinkageService{}
	}
	return modelFormLinkage
}
