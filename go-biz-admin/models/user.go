package models

import (
	"time"
	"log"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// 注意：不使用软删除，所以不包含 deleted_at 字段
	Name              string      `json:"name" gorm:"not null"`
	Email             string      `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	Password          string      `json:"password" gorm:"not null"`
	Role              string      `json:"role" gorm:"default:'user'"`
	Phone             string      `json:"phone" gorm:"type:varchar(50)"`
	Address           string      `json:"address" gorm:"type:varchar(255)"`
	Status            int         `json:"status" gorm:"default:1"`                    // 0:禁用, 1:启用
	ZJMFUserID        int         `json:"zjmf_user_id"`                               // 智简魔方用户ID
	ZJMFClientID      string      `json:"zjmf_client_id" gorm:"type:varchar(255)"`   // 智简魔方客户端ID
	ZJMFApiKey        string      `json:"zjmf_api_key" gorm:"type:varchar(255)"`     // 智简魔方API密钥
	ZJMFAccountStatus string      `json:"zjmf_account_status" gorm:"type:varchar(50)"` // 智简魔方账户状态
	Username          string      `json:"username" gorm:"type:varchar(100);uniqueIndex"` // 用户名
	CustomID          uint        `json:"custom_id" gorm:"uniqueIndex"`               // 自定义ID
	UUID              string      `json:"uuid" gorm:"type:varchar(100);uniqueIndex"`  // UUID
	RealNameAuth      *bool       `json:"real_name_auth"`                             // 是否实名认证
	Avatar            string      `json:"avatar" gorm:"type:varchar(255)"`            // 用户头像
	Profile           UserProfile `json:"profile,omitempty" gorm:"foreignKey:UserID"` // 关联用户资料
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

type UserProfile struct {
	ID         uint       `json:"id" gorm:"primaryKey"`
	UserID     uint       `json:"user_id" gorm:"not null;uniqueIndex"`
	RealName   string     `json:"real_name" gorm:"type:varchar(100)"`
	IDCard     string     `json:"id_card" gorm:"type:varchar(50)"`
	Phone      string     `json:"phone" gorm:"type:varchar(20)"`
	Address    string     `json:"address" gorm:"type:varchar(255)"`
	AuthStatus string     `json:"auth_status" gorm:"default:'pending'"` // pending, approved, rejected
	AuthAt     *time.Time `json:"auth_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&User{},
		&UserProfile{},
		&Supplier{},
		&Server{},
		&Product{},
		&Task{},
		&Order{},
		&ProductOrder{},
		&RenewalOrder{},
		&Transaction{},
		&Bill{},
		&Ticket{},
		&UpstreamProduct{},
	)
	if err != nil {
		log.Fatalf("自动迁移数据库表失败: %v", err)
	}

	log.Println("数据库表结构自动迁移完成!")
}