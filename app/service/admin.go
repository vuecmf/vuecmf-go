package service

import (
	"github.com/vuecmf/vuecmf-go/app"
	"github.com/vuecmf/vuecmf-go/app/model"
)

// AdminService admin服务结构
type AdminService struct {
	*BaseService
}

// List 功能：获取列表数据
// 		参数：
//			filter 		map 	过滤查询
//			keywords 	string 	关键词
//			page 		int		当前页码
//			pageSize 	int 	每页显示条数
//			orderField 	string 	排序字段名
//			orderSort 	string 	排序方式（desc 倒序, asc升序）
func (service *AdminService) List(filter map[string]interface{}, keywords string, page int, pageSize int, orderField string, orderSort string) []model.Admin{

	var adminList []model.Admin

	if pageSize == 0 {
		pageSize = 20
	}

	if page == 0 {
		page = 1
	}

	if orderField == "" {
		orderField = "id"
	}

	db := app.Db{}
	query := db.Connect().Table("admin").Paginate(page, pageSize)

	if keywords != "" {
		kw := "%" + keywords + "%"
		query = query.Where("username LIKE ?", kw).
					  Or("email LIKE ?", kw).
					  Or("mobile LIKE ?", kw)

	}else if len(filter) > 0 {
		query = query.Where(filter)
	}

	query.Order(orderField + " " + orderSort).Find(&adminList)

	return adminList

}
