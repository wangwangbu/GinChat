package models

import "gorm.io/gorm"

// 用户关系
type GroupBasic struct {
	gorm.Model
	Name    string
	OwnerID uint
	Icon    string
	Type    int
	Desc    string
}

func (table *GroupBasic) TableName() string {
	return "group_basic"
}
