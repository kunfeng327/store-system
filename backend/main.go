package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"store/backend/models"
	"store/backend/routes"
)

func main() {
	// 读取配置
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置失败（请确认 backend 目录下有 config.yaml）: %v", err)
	}

	// 连接 MySQL
	dsn := viper.GetString("db.dsn")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("连接 MySQL 失败: %v", err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 自动建表 + 种子数据
	if err := db.AutoMigrate(&models.User{}, &models.Product{}, &models.Review{},
		&models.Order{}, &models.OrderItem{}, &models.BlindBoxRecord{}); err != nil {
		log.Fatalf("建表失败: %v", err)
	}
	models.Seed(db)

	// JWT 配置
	models.JWTSsecret = viper.GetString("jwt.secret")
	models.JWTExpire = time.Duration(viper.GetInt("jwt.expire_hours")) * time.Hour

	r := gin.Default()
	r.Use(CORS())
	routes.Setup(r, db)

	port := viper.GetString("server.port")
	fmt.Printf("商店系统后端启动: http://localhost:%s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

// CORS 允许前端跨域访问
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
