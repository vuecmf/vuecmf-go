package model



// Upload  模型结构
type Upload struct {
	
}

// DataUploadForm 提交的表单数据
type DataUploadForm struct {
    Data *Upload `json:"data" form:"data"`
}