package models

import (
	"fmt"
	"ginchat/utils"
	"time"

	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model
	Name          string
	PassWord      string
	Phone         string `valid:"matches(^1[3-9]{1}\\d{9}$)"`
	Email         string `valid:"email"`
	Identity      string
	ClientIP      string
	ClientPort    string
	Salt          string
	LoginTime     time.Time
	HeartbeatTime time.Time
	LoginOutTime  time.Time
	IsLogOut      bool
	DeviceInfo    string
}

func NewUserBasic() (u UserBasic) {
	u.HeartbeatTime = time.Now()
	u.LoginOutTime = time.Now()
	u.LoginTime = time.Now()
	return
}

func (u *UserBasic) TableName() string {
	return "user_basic"
}

func FindUserByNameAndPwd(name string, password string) UserBasic {
	user := UserBasic{}
	db := utils.GetDB()
	db.Where("name = ? and pass_word = ?", name, password).First(&user)
	return user
}

func FindUserByName(name string) UserBasic {
	db := utils.GetDB()
	user := NewUserBasic()
	db.Where("name=?", name).First(&user)
	return user
}

func FindUserByPhone(phone string) UserBasic {
	db := utils.GetDB()
	user := NewUserBasic()
	db.Where("phone=?", phone).First(&user)
	return user
}

func FindUserByEmail(email string) UserBasic {
	db := utils.GetDB()
	user := NewUserBasic()
	db.Where("email=?", email).First(&user)
	return user
}

func GetUserList() []*UserBasic {
	data := make([]*UserBasic, 10)
	db := utils.GetDB()
	db.Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	return data
}

func CreateUser(user UserBasic) *gorm.DB {
	db := utils.GetDB()
	return db.Create(&user)
}

func DeleteUser(user UserBasic) *gorm.DB {
	db := utils.GetDB()
	return db.Delete(&user)
}

func UpdateUser(user UserBasic) *gorm.DB {
	db := utils.GetDB()
	return db.Model(&user).Updates(UserBasic{Name: user.Name, PassWord: user.PassWord, Phone: user.Phone, Email: user.Email})
}
