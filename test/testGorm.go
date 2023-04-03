package main

import (
	// "fmt"
	"ginchat/models"
	"ginchat/utils"
	// "time"
)

func main() {

  utils.InitConfig()
  utils.InitMySQL()
  db := utils.GetDB()
  // 迁移 schema
  db.AutoMigrate(&models.UserBasic{})

  // Create
  // user := &models.UserBasic{}
  // user.Name = "wwb"
  // user.HeartbeatTime = time.Now()
  // user.LoginOutTime = time.Now()
  // user.LoginTime = time.Now()
  // db.Create(user)

  // Read
  // fmt.Println(db.First(user, 1))

  // Update - 将 product 的 price 更新为 200
  // db.Model(user).Update("PassWord", "1234")
  // Update - 更新多个字段
  // db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // 仅更新非零值字段
  // db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

  // Delete - 删除 product
  // db.Delete(&product, 1)
}