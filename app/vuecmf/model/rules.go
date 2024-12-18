//+----------------------------------------------------------------------
// | Copyright (c) 2024 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: tulihua2004@126.com
// +----------------------------------------------------------------------

package model

// Rules 权限规则 模型结构
type Rules struct {
	ID    uint   `json:"id" form:"id"  gorm:"column:id;primaryKey;autoIncrement;size:32;not null;comment:自增ID"`
	Ptype string `json:"ptype" form:"ptype" gorm:"column:ptype;size:4;uniqueIndex:unique_index;not null;default:'';comment:类型：g=组或角色，p=策略"`
	V0    string `json:"v0" form:"v0" gorm:"column:v0;size:64;uniqueIndex:unique_index;not null;default:'';comment:对应定义的sub(用户名或角色名)"`
	V1    string `json:"v1" form:"v1" gorm:"column:v1;size:64;uniqueIndex:unique_index;not null;default:'';comment:对应定义的dom(角色或应用名)"`
	V2    string `json:"v2" form:"v2" gorm:"column:v2;size:64;uniqueIndex:unique_index;not null;default:'';comment:对应定义的obj(应用名或控制器名)"`
	V3    string `json:"v3" form:"v3" gorm:"column:v3;size:64;uniqueIndex:unique_index;not null;default:'';comment:对应定义的act(动作名称)"`
	V4    string `json:"v4" form:"v4" gorm:"column:v4;size:64;uniqueIndex:unique_index;not null;default:'';comment:预留，暂用不到"`
	V5    string `json:"v5" form:"v5" gorm:"column:v5;size:64;uniqueIndex:unique_index;not null;default:'';comment:预留，暂用不到"`
	V6    string `json:"v6" form:"v6" gorm:"column:v6;size:64;uniqueIndex:unique_index;not null;default:'';comment:预留，暂用不到"`
	V7    string `json:"v7" form:"v7" gorm:"column:v7;size:64;uniqueIndex:unique_index;not null;default:'';comment:预留，暂用不到"`
}
