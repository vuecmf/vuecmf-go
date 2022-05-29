package model

import "time"

// Migrations  数据迁移 模型结构
type Migrations struct {
	Base
	Version       int64     `json:"version" gorm:"column:version;size:20;not null;default:0;comment:版本号"`
	MigrationName string    `json:"migration_name" gorm:"column:migration_name;size:100;default:;comment:迁移名称"`
	StartTime     time.Time `json:"start_time" gorm:"column:start_time;not null;autoCreateTime;comment:开始时间"`
	EndTime       time.Time `json:"end_time" gorm:"column:end_time;not null;autoCreateTime;comment:结束时间"`
	Breakpoint    uint8     `json:"breakpoint" gorm:"column:breakpoint;size:1;not null;default:0;comment:断点"`
}
