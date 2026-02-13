package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"

	"go-biz-admin/api"
	"go-biz-admin/config"
	"go-biz-admin/cron"
	"go-biz-admin/models"
	"go-biz-admin/services"
)

func main() {
	// 加载.env文件
	loadEnvFile()

	// 初始化配置和数据库
	config.InitDB()

	// 自动迁移数据库表结构
	models.AutoMigrate(config.DB)

	// 初始化服务
	services.InitServices()

	// 启动定时任务
	go cron.StartCronJobs()

	// 创建路由
	r := api.SetupRouter()

	// 从环境变量获取端口号，默认为8080
	port := getEnvOrDefault("PORT", "8081")

	// 输出启动信息
	fmt.Printf("服务器启动于 :%s\n", port)
	fmt.Printf("定时任务执行于 %s\n", time.Now().Format("2006-01-02 15:04:05"))

	// 启动服务器
	if err := r.Run(":" + port); err != nil {
		log.Printf("服务器启动失败: %v", err)
	}
}

// loadEnvFile 加载.env配置文件
func loadEnvFile() {
	// 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		log.Printf("获取当前工作目录失败: %v", err)
		return
	}

	// 构建.env文件路径
	envPath := filepath.Join(wd, ".env")

	// 检查.env文件是否存在
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		log.Printf("警告: .env文件不存在于 %s，将使用系统环境变量", envPath)
		return
	}

	// 加载.env文件
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("加载.env文件失败: %v", err)
		return
	}

	log.Printf("成功加载.env文件: %s", envPath)

	// 打印关键配置信息用于调试
	log.Printf("数据库配置 - HOST: %s, USER: %s, NAME: %s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"))
}

// getEnvOrDefault 获取环境变量，如果不存在则返回默认值
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
