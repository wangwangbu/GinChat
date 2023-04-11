package models

import (
	"fmt"
	"ginchat/utils"

	"gorm.io/gorm"
)

// 用户关系
type Contact struct {
	gorm.Model
	OwnerID  uint // 用户 ID
	TargetId uint // 对应 ID
	Type     int  // 关系类型 1好友 2群组
	Desc     string
}

func (table *Contact) TableName() string {
	return "contact"
}

func SearchFriends(userId uint) []UserBasic {
	contacts := make([]Contact, 0)
	objIds := make([]uint64, 0)
	db := utils.GetDB()
	db.Where("owner_id = ? and type = 1", userId).Find(&contacts)
	for _, v := range contacts {
		fmt.Println(v)
		objIds = append(objIds, uint64(v.TargetId))
	}
	users := make([]UserBasic, 0)
	db.Where("id in ?", objIds).Find(&users)
	return users
}