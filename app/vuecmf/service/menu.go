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
)

// menuService menu服务结构
type menuService struct {
	*baseService
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
func (ser *menuService) Nav(username string, isSuper interface{}, appName string) (interface{}, error) {
	if appName == "" {
		appName = "vuecmf"
	}

	idList, err := Auth().GetPermissions(username, isSuper, appName)
	if err != nil {
		return nil, err
	}

	var apiIdList []string
	for _, val := range idList {
		apiIdList = append(apiIdList, val...)
	}

	var res = make(map[string]interface{})
	res["api_maps"] = ModelAction().GetAllApiMap(apiIdList)

	modelIdList := ModelAction().GetModelIdListById(apiIdList)

	/*var menuList = make(map[string]interface{})
	var idPath  []string
	var pathName []string
	ser.getNav(menuList, modelIdList, apiIdList, 0, idPath, pathName)*/

	var NavMenuList []*model.NavMenu
	ser.getNavMenu(&NavMenuList, modelIdList, apiIdList)
	menuList := model.MenuModel().ToNavTree(NavMenuList)
	res["nav_menu"] = menuList

	return res, nil
}


func (ser *menuService) getNavMenu(dataList interface{}, modelIdList []string, apiIdList []string) {
	db.Table(ns.TableName(ser.TableName) + " vm").
		Select("concat('m',vm.id) mid, vm.id, vm.pid, vm.title, vm.icon, vm.model_id, vmc.table_name, vmc.component_tpl, vmc.search_field_id, vmc.is_tree, vma.action_type default_action_type").
		Joins("left join " + ns.TableName("model_config") + " vmc on vmc.id = vm.model_id and vmc.status = 10").
		Joins("left join " + ns.TableName("model_action") + " vma on vmc.default_action_id = vma.id and vma.status = 10 and vma.id in ?", apiIdList).
		Where("vm.status = 10").
		Where("vm.model_id in ?", modelIdList).
		Order("vm.sort_num").Find(dataList)
}



/*type menuItem struct {
	Mid string
	Id int
	Pid int
	Title string
	Icon string
	ModelId int
	IdPath []string
	PathName []string
	Children interface{}
	modelCfg
}

type modelCfg struct {
	TableName string
	DefaultActionType string
	ComponentTpl string
	SearchFieldId string
	IsTree int
}

// getNav 根据模型ID/动作ID 获取对应导航菜单列表
func (ser *menuService) getNav(menuList map[string]interface{}, modelIdList []string, apiIdList []string, pid int, idPath []string, pathName []string) interface{}{
	var res []menuItem
	db.Table(ns.TableName(ser.TableName)).Select("concat('m',id) mid, id, pid, title, icon, model_id").
		Where("pid = ?", pid).
		Where("status = 10").
		Where("model_id in ?", modelIdList).
		Order("sort_num").Find(&res)




	for _, val := range res {
		if len(idPath) == 0 {
			val.IdPath = append(val.IdPath, val.Mid)
			val.PathName = append(val.PathName, val.Title)
		} else {
			val.IdPath = idPath
			val.PathName = pathName
			val.IdPath = append(val.IdPath, val.Mid)
			val.PathName = append(val.PathName, val.Title)
		}

		menuList[val.Mid] = make(map[string]interface{})
		child := ser.getNav(menuList[val.Mid], modelIdList, apiIdList, val.Id, val.IdPath, val.PathName)
		if child != nil {
			val.Children = child
		} else {
			var mc modelCfg
			db.Table(ns.TableName("model_config") + " MC").
				Select("MC.table_name, MA.action_type default_action_type, MC.component_tpl, MC.search_field_id, MC.is_tree").
				Joins("left join " + ns.TableName("model_action") + " MA on MC.default_action_id = MA.id").
				Where("MC.id = ?", val.ModelId).
				Where("MA.id in ?", apiIdList).
				Where("MC.status = 10").
				Where("MA.status = 10").Find(&mc)

			if mc.TableName != "" {
				val.modelCfg = mc
			} else {
				continue
			}
		}

		menuList[val.Mid] = val
	}
	return menuList
}*/