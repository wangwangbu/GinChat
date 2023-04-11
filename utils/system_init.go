package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db *gorm.DB
	Red *redis.Client
)

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

func InitRedis() {
	Red = redis.NewClient(&redis.Options{
		Addr: viper.GetString("redis.addr"),
		Password: viper.GetString("redis.password"),
		DB: viper.GetInt("redis.DB"),
		PoolSize: viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConns"),
	})
	pong, err := Red.Ping(Red.Context()).Result()
	if err != nil {
		fmt.Println("init redis error ...", err)
	} else {
		fmt.Println("Redis inited success ...", pong)
	}
}

const (
	PublishKey = "websocket"
)

// Publish 发布消息到 Redis
func Publish(ctx context.Context, channel string, msg string) error {
	var err error
	err = Red.Publish(ctx, channel, msg).Err()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Publish...", msg)
	return err
}

// Subscribe 订阅 Redis 消息
func Subscribe(ctx context.Context, channel string) (string, error) {
	sub := Red.Subscribe(ctx, channel)
	msg, err := sub.ReceiveMessage(ctx)
	fmt.Println("Subscribe...", msg.Payload)
	return msg.Payload, err
}

func GetDB() *gorm.DB {
	return db
}