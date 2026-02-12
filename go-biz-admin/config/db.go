package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"go-biz-admin/models"  // 导入models包

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

type Config struct {
	DBType     string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSslMode  string
	ZJMF       ZJMFConfig
}

// ZJMFConfig 智简魔方配置结构
type ZJMFConfig struct {
	BaseURL   string
	APIKey    string
	APISecret string
}

var AppConfig Config

// InitConfig 初始化应用程序的数据库配置
// 该函数从环境变量中获取数据库配置信息，如果环境变量不存在则使用默认值
// 配置包括数据库类型、主机地址、端口、用户名、密码、数据库名和SSL模式
func InitConfig() {
	AppConfig = Config{
		DBType:     getEnv("DB_TYPE", "mysql"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),    // MySQL默认端口
		DBUser:     getEnv("DB_USER", "root"),     // MySQL用户名，默认为root
		DBPassword: getEnv("DB_PASSWORD", ""), // MySQL密码，空字符串表示无密码
		DBName:     getEnv("DB_NAME", "gobiz_admin"),     // MySQL 数据库名
		DBSslMode:  getEnv("DB_SSL_MODE", "disable"),
		ZJMF: ZJMFConfig{
			BaseURL:   getEnv("ZJMF_BASE_URL", "http://w2.test.idcsmart.com"),
			APIKey:    getEnv("ZJMF_API_KEY", "your_api_key"),
			APISecret: getEnv("ZJMF_API_SECRET", "your_api_secret"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func InitDB() {
	// 确保配置被初始化
	InitConfig()

	var err error
	
	// 首先尝试连接到MySQL服务器（不指定数据库）
	dsnWithoutDB := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		AppConfig.DBUser, AppConfig.DBPassword, AppConfig.DBHost, AppConfig.DBPort)

	DB, err = gorm.Open(mysql.Open(dsnWithoutDB), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 开启GORM日志
	})
	if err != nil {
		log.Printf("连接MySQL服务器失败: %v\n", err)
		log.Printf("尝试连接到 %s@tcp(%s:%s)\n", AppConfig.DBUser, AppConfig.DBHost, AppConfig.DBPort)
		log.Println("请确保已安装并启动MySQL服务，并且数据库凭据正确")
		os.Exit(1) // 如果数据库服务器连接失败，退出程序
	}

	// 检查并创建数据库
	err = createDatabaseIfNotExists(DB, AppConfig.DBName)
	if err != nil {
		log.Printf("创建数据库失败: %v\n", err)
		os.Exit(1)
	}

	// 关闭当前连接，重新连接到指定数据库
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("获取底层sql.DB失败: %v", err)
	}
	sqlDB.Close()

	// 重新连接到指定数据库
	dsnWithDB := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		AppConfig.DBUser, AppConfig.DBPassword, AppConfig.DBHost, AppConfig.DBPort, AppConfig.DBName)

	DB, err = gorm.Open(mysql.Open(dsnWithDB), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 开启GORM日志
	})
	if err != nil {
		log.Printf("连接数据库 %s 失败: %v\n", AppConfig.DBName, err)
		log.Printf("尝试连接到 %s@tcp(%s:%s)/%s\n", AppConfig.DBUser, AppConfig.DBHost, AppConfig.DBPort, AppConfig.DBName)
		os.Exit(1) // 如果数据库连接失败，退出程序
	}

	log.Printf("数据库 %s 连接成功!\n", AppConfig.DBName)

	sqlDB, err = DB.DB()
	if err != nil {
		log.Fatalf("获取底层sql.DB失败: %v", err)
	}
	
	// 设置连接池
	sqlDB.SetMaxIdleConns(10)              // 设置空闲连接池中的最大连接数
	sqlDB.SetMaxOpenConns(100)             // 设置打开数据库连接的最大数量
	sqlDB.SetConnMaxLifetime(30 * time.Minute) // 设置连接可重用的最长时间

	// 自动迁移数据库表结构
	autoMigrateTables()
}

// createDatabaseIfNotExists 检查并创建数据库（如果不存在）
func createDatabaseIfNotExists(db *gorm.DB, dbName string) error {
	// 检查数据库是否存在
	var count int64
	err := db.Raw("SELECT COUNT(*) FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?", dbName).Count(&count).Error
	if err != nil {
		return fmt.Errorf("检查数据库是否存在时出错: %v", err)
	}

	if count == 0 {
		// 数据库不存在，创建数据库
		log.Printf("数据库 %s 不存在，正在创建...\n", dbName)
		err = db.Exec(fmt.Sprintf("CREATE DATABASE `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;", dbName)).Error
		if err != nil {
			return fmt.Errorf("创建数据库 %s 失败: %v", dbName, err)
		}
		log.Printf("数据库 %s 创建成功!\n", dbName)
	} else {
		log.Printf("数据库 %s 已存在。\n", dbName)
	}

	return nil
}

// autoMigrateTables 自动迁移数据库表结构
func autoMigrateTables() {
	// 导入模型文件中的结构体定义
	// 注意：这里的结构体定义应该在models包中定义
	err := DB.AutoMigrate(
		&models.User{},
		&models.UserProfile{},
		&models.Supplier{},
		&models.Server{},
		&models.Product{},
		&models.Task{},
		&models.Order{},
		&models.ProductOrder{},
		&models.RenewalOrder{},
		&models.Transaction{},
		&models.Bill{},
		&models.Ticket{},
	)
	if err != nil {
		log.Fatalf("自动迁移数据库表失败: %v", err)
	}

	log.Println("数据库表结构自动迁移完成!")
}

// ensureSupplierTableColumns 确保供应商表包含所有必需的列
func ensureSupplierTableColumns() {
	// GORM v2 不需要手动处理此逻辑，AutoMigrate会自动处理
}

// removeContactFieldsFromSuppliers 移除供应商表中的联系人相关字段
func removeContactFieldsFromSuppliers() {
	// GORM v2 不需要手动处理此逻辑
}

// removeColumnIfExists 检查列是否存在，如果存在则删除
func removeColumnIfExists(tableName, columnName string) {
	// GORM v2 不需要手动处理此逻辑
}

// addColumnIfNotExists 检查列是否存在，如果不存在则添加
func addColumnIfNotExists(tableName, columnName, columnDefinition string) {
	// GORM v2 不需要手动处理此逻辑
}