package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	log.Println("Read config success...")
}

func InitMySQL() {
	// 定义日志模板，打印 SQL 语句
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config {
			SlowThreshold: time.Second,	// 慢 SQL 阈值
			LogLevel: logger.Info,	// 级别
			Colorful: true,	// 彩色
		},
	)


	var err error
	db, err = gorm.Open(mysql.Open(viper.GetString("mysql.dsn")), 
	 	&gorm.Config{Logger: newLogger})
	if err != nil {
		panic("failed to connect database")
	}
	log.Println("Open database success...")
}

func GetDB() *gorm.DB {
	return db
}