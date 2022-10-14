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
	"github.com/vuecmf/vuecmf-go/app/vuecmf/helper"
	"github.com/vuecmf/vuecmf-go/app/vuecmf/model"
	"strconv"
)

// menuService menu服务结构
type menuService struct {
	*BaseService
	TableName string
}

var menu *menuService

// Menu 获取menu服务实例
func Menu() *menuService {
	if menu == nil {
		menu = &menuService{TableName: "menu"}
	}
	return menu
}

// GetIdPath 获取父级ID的ID路径
func (ser *menuService) GetIdPath(pid uint) string {
	var pidIdPath string
	Db.Table(NS.TableName(ser.TableName)).Select("id_path").Where("id = ?", pid).Find(&pidIdPath)
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
	Title 	 string
	PathName string
}

// GetPathName 获取父级ID的path路径
func (ser *menuService) GetPathName(pid uint) string {
	pidPathName := ""
	var mi menuInfo
	Db.Table(NS.TableName(ser.TableName)).Select("title, path_name").Where("id = ?", pid).Find(&mi)
	if pid > 0 {
		if mi.PathName == "" {
			pidPathName = mi.Title
		} else {
			pidPathName += "," + mi.Title
		}
	}
	return pidPathName
}

// Create 创建单条或多条数据, 成功返回影响行数
func (ser *menuService) Create(data *model.Menu) (int64, error) {
	data.IdPath = ser.GetIdPath(data.Pid)
	data.PathName = ser.GetPathName(data.Pid)
	res := Db.Create(data)
	return res.RowsAffected, res.Error
}

// Update 更新数据, 成功返回影响行数
func (ser *menuService) Update(data *model.Menu) (int64, error) {
	data.IdPath = ser.GetIdPath(data.Pid)
	data.PathName = ser.GetPathName(data.Pid)
	res := Db.Updates(data)
	return res.RowsAffected, res.Error
}


// List 获取列表数据
// 		参数：params 查询参数
func (ser *menuService) List(params *helper.DataListParams) (interface{}, error) {
	if params.Data.Action == "getField" {
		//拉取列表的字段信息
		return ser.getFieldList(ser.TableName, params.Data.Filter)
	} else {
		//拉取列表的数据
		var menuList []*model.Menu
		var res = make(map[string]interface{})

		ser.getList(&menuList, ser.TableName, params)

		//转换成树形列表
		tree := model.MenuModel().ToTree(menuList)
		res["data"] = tree
		return res, nil
	}
}

// Nav 获取用户的后台导航菜单
func (ser *menuService) Nav(username string, isSuper interface{}) (interface{}, error) {
	var err error
	//先取不需要授权的应用下的所有动作ID
	apiIdList := ModelAction().GetNotAuthActionIds()

	//再取出需要授权的应用下有权限的动作ID
	appList := AppConfig().GetAuthAppList()
	for _, appName := range appList {
		idList, err := Auth().GetPermissions(username, isSuper, appName)
		if err != nil {
			break
		}

		for _, val := range idList {
			apiIdList = append(apiIdList, val...)
		}
	}

	if err != nil {
		return nil, err
	}

	var res = make(map[string]interface{})
	res["api_maps"] = ModelAction().GetAllApiMap(apiIdList)

	modelIdList := ModelAction().GetModelIdListById(apiIdList)

	var NavMenuList []*model.NavMenu
	ser.getNavMenu(&NavMenuList, modelIdList, apiIdList)
	menuList := model.MenuModel().ToNavTree(NavMenuList)
	res["nav_menu"] = menuList

	return res, nil
}

func (ser *menuService) getNavMenu(dataList interface{}, modelIdList []string, apiIdList []string) {
	Db.Table(NS.TableName(ser.TableName)+" vm").
		Select("concat('m',vm.id) mid, vm.id, vm.pid, vm.id_path id_path_str, vm.path_name path_name_str, vm.title, vm.icon, vm.model_id, vmc.table_name, vmc.component_tpl, vmc.search_field_id, vmc.is_tree, vma.action_type default_action_type, vm.app_id, AC.app_name").
		Joins("left join "+NS.TableName("model_config")+" vmc on vmc.id = vm.model_id and vmc.status = 10").
		Joins("left join "+NS.TableName("model_action")+" vma on vmc.default_action_id = vma.id and vma.status = 10 and vma.id in ?", apiIdList).
		Joins("left join " + NS.TableName("app_config") + " AC on vm.app_id = AC.id").
		Where("vm.status = 10").
		Where("vm.model_id in ?", modelIdList).
		Order("vm.sort_num").Find(dataList)
}
