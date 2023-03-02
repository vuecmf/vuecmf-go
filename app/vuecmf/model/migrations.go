//+----------------------------------------------------------------------
// | Copyright (c) 2023 http://www.vuecmf.com All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://github.com/vuecmf/vuecmf-go/blob/master/LICENSE )
// +----------------------------------------------------------------------
// | Author: vuecmf <tulihua2004@126.com>
// +----------------------------------------------------------------------

package model

import "time"

// Migrations  数据迁移 模型结构
type Migrations struct {
	Version       int64     `json:"version" form:"version" gorm:"column:version;size:64;not null;default:0;comment:版本号"`
	MigrationName string    `json:"migration_name" form:"migration_name" gorm:"column:migration_name;size:100;default:'';comment:迁移名称"`
	StartTime     time.Time `json:"start_time" form:"start_time" gorm:"column:start_time;not null;autoCreateTime;comment:开始时间"`
	EndTime       time.Time `json:"end_time" form:"end_time" gorm:"column:end_time;not null;autoCreateTime;comment:结束时间"`
	Breakpoint    uint8     `json:"breakpoint" form:"breakpoint" gorm:"column:breakpoint;size:1;not null;default:0;comment:断点"`
}
