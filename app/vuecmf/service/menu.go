//+----------------------------------------------------------------------
// | Copyright (c) 2024 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: tulihua2004@126.com
// +----------------------------------------------------------------------

package service

import (
	"github.com/vuecmf/vuecmf-go/v3/app"
	"github.com/vuecmf/vuecmf-go/v3/app/vuecmf/helper"
	"github.com/vuecmf/vuecmf-go/v3/app/vuecmf/model"
	"strconv"
	"sync"
)

// MenuService menu服务结构
type MenuService struct {
	*BaseService
}

var menuOnce sync.Once
var menu *MenuService

// Menu 获取menu服务实例
func Menu() *MenuService {
	menuOnce.Do(func() {
		menu = &MenuService{
			BaseService: &BaseService{
				"menu",
				&model.Menu{},
				&[]model.Menu{},
				[]string{"title", "icon"},
			},
		}
	})
	return menu
}

// GetIdPath 获取父级ID的ID路径
//
//	参数：
//		pid 父级ID
func (svc *MenuService) GetIdPath(pid uint) string {
	var pidIdPath string
	DbTable(svc.TableName).Select("id_path").Where("id = ?", pid).Find(&pidIdPath)
	if pid > 0 {
		if pidIdPath == "" {
			pidIdPath = "m" + strconv.Itoa(int(pid))
		} else {
			pidIdPath += ",m" + strconv.Itoa(int(pid))
		}
	}
	return pidIdPath
}

type menuInfo struct {
	Title    string
	PathName string
}

// GetPathName 获取父级ID的path路径
//
//	参数：
//		pid 父级ID
//		title 标题
func (svc *MenuService) GetPathName(pid uint, title string) string {
	pidPathName := ""
	var parent menuInfo
	DbTable(svc.TableName).Select("title, path_name").Where("id = ?", pid).Find(&parent)
	if pid > 0 {
		if parent.PathName == "" {
			pidPathName = parent.Title + "," + title
		} else {
			pidPathName = parent.PathName + "," + title
		}
	}
	return pidPathName
}

// Create 创建单条或多条数据, 成功返回影响行数
//
//	参数：
//		data 需保存的数据
func (svc *MenuService) Create(data *model.Menu) (int64, error) {
	data.IdPath = svc.GetIdPath(data.Pid)
	data.PathName = svc.GetPathName(data.Pid, data.Title)
	res := app.Db.Create(data)
	return res.RowsAffected, res.Error
}

// Update 更新数据, 成功返回影响行数
//
//	参数：
//		data 需更新的数据
func (svc *MenuService) Update(data *model.Menu) (int64, error) {
	data.IdPath = svc.GetIdPath(data.Pid)
	data.PathName = svc.GetPathName(data.Pid, data.Title)
	res := app.Db.Updates(data)
	return res.RowsAffected, res.Error
}

// List 获取列表数据
//
//	参数：
//		params 查询参数
func (svc *MenuService) List(params *helper.DataListParams) (interface{}, error) {
	if params.Data.Action == "getField" {
		//拉取列表的字段信息
		return svc.GetFieldList(svc.TableName, params.Data.Filter, 10)
	} else {
		//拉取列表的数据
		var menuList []*model.Menu
		var res = make(map[string]interface{})

		svc.GetList(&menuList, svc.TableName, params.Data)

		//转换成树形列表
		tree := model.MenuModel().ToTree(menuList)
		res["data"] = tree
		return res, nil
	}
}

// Nav 获取用户的后台导航菜单
//
//	参数：
//		username 用户名
//		isSuper 是否为超级管理员
func (svc *MenuService) Nav(username string, isSuper uint16) (interface{}, error) {
	var err error
	//先取不需要授权的应用下的所有动作ID
	apiIdList := ModelAction().GetNotAuthActionIds()

	//再取出需要授权的应用下有权限的动作ID
	idList, err := Auth().GetPermissions(username, isSuper)
	if err != nil {
		return nil, err
	}
	for _, val := range idList {
		apiIdList = append(apiIdList, val...)
	}

	var res = make(map[string]interface{})
	res["api_maps"] = ModelAction().GetAllApiMap(apiIdList)

	modelIdList := ModelAction().GetModelIdListById(apiIdList)

	NavMenuList, err := svc.getNavMenu(modelIdList, apiIdList)
	if err != nil {
		return nil, err
	}

	menuList := model.MenuModel().ToNavTree(NavMenuList)
	res["nav_menu"] = menuList

	return res, nil
}

func (svc *MenuService) getNavMenu(modelIdList []string, apiIdList []string) ([]*model.NavMenu, error) {
	var dataList []*model.NavMenu
	err := DbTable(svc.TableName, "vm").
		Select("concat('m',vm.id) mid, vm.id, vm.pid, vm.id_path id_path_str, vm.path_name path_name_str, vm.title, vm.icon, vm.model_id, vmc.table_name, vmc.component_tpl, vmc.search_field_id, vmc.is_tree, vma.action_type default_action_type, vm.app_id, AC.app_name").
		Joins("left join "+TableName("model_config")+" vmc on vmc.id = vm.model_id and vmc.status = 10").
		Joins("left join "+TableName("model_action")+" vma on vmc.default_action_id = vma.id and vma.status = 10 and vma.id in ?", apiIdList).
		Joins("left join "+TableName("app_config")+" AC on vm.app_id = AC.id").
		Where("vm.status = 10").
		Where("vm.model_id in ?", modelIdList).
		Order("vm.sort_num").Find(&dataList).Error
	return dataList, err
}
