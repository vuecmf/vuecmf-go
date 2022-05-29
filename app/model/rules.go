package model

// Rules 权限规则 模型结构
type Rules struct {
	ID    uint   `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Ptype string `json:"ptype" gorm:"column:ptype;size:4;uniqueIndex:unique_index;not null;default:;comment:类型：g=组或角色，p=策略"`
	V0    string `json:"v0" gorm:"column:v0;size:64;uniqueIndex:unique_index;not null;default:;comment:对应定义的sub(用户名或角色名)"`
	V1    string `json:"v1" gorm:"column:v1;size:64;uniqueIndex:unique_index;not null;default:;comment:对应定义的dom(角色或应用名)"`
	V2    string `json:"v2" gorm:"column:v2;size:64;uniqueIndex:unique_index;not null;default:;comment:对应定义的obj(应用名或控制器名)"`
	V3    string `json:"v3" gorm:"column:v3;size:64;uniqueIndex:unique_index;not null;default:;comment:对应定义的act(动作名称)"`
	V4    string `json:"v4" gorm:"column:v4;size:64;uniqueIndex:unique_index;not null;default:;comment:预留，暂用不到"`
	V5    string `json:"v5" gorm:"column:v5;size:64;uniqueIndex:unique_index;not null;default:;comment:预留，暂用不到"`
	
}