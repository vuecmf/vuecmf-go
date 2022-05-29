package model

// Base 基础模型
type Base struct {
	Id uint `json:"id" gorm:"column:id;primaryKey;autoIncrement;size:11;comment:ID"`
	Status uint8 `json:"status" gorm:"column:status;size:4;not null;default:10;comment:状态：10=开启，20=禁用"`
}