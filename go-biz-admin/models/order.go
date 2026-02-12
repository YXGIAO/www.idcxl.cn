package models

import (
	"time"
)

// Order 订单模型
type Order struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	OrderNumber     string    `json:"order_number" gorm:"type:varchar(100);uniqueIndex;not null"`
	UserID          uint      `json:"user_id"`
	Amount          float64   `json:"amount"`
	Status          string    `json:"status" gorm:"default:'pending';type:varchar(50)"` // pending, paid, cancelled, refunded
	PaymentMethod   string    `json:"payment_method" gorm:"type:varchar(50)"`
	Description     string    `json:"description" gorm:"type:text"`
	ZJMFOrderID     string    `json:"zjmf_order_id" gorm:"type:varchar(100)"`      // 智简魔方订单ID
	ZJMFSupplierID  uint      `json:"zjmf_supplier_id"`   // 智简魔方供应商ID
	ZJMFProductID   string    `json:"zjmf_product_id" gorm:"type:varchar(100)"`    // 智简魔方产品ID
	ZJFServerID    string    `json:"zjmf_server_id" gorm:"type:varchar(100)"`     // 智简魔方服务器ID
	ZJFHostingID   string    `json:"zjmf_hosting_id" gorm:"type:varchar(100)"`    // 智简魔方托管ID
	ZJFStatus      string    `json:"zjmf_status" gorm:"type:varchar(50)"`        // 智简魔方订单状态
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Order) TableName() string {
	return "orders"
}

// Product 产品模型
type Product struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"type:varchar(255)"`
	Description string    `json:"description" gorm:"type:text"`
	Price       float64   `json:"price"`
	Status      string    `json:"status" gorm:"default:'active';type:varchar(50)"` // active, inactive, discontinued
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductOrder struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id"`
	ProductID   uint      `json:"product_id"`
	ProductName string    `json:"product_name" gorm:"type:varchar(255)"`
	Amount      float64   `json:"amount"`
	Status      string    `json:"status" gorm:"default:'pending';type:varchar(50)"` // pending, paid, cancelled, completed
	PaymentMethod string  `json:"payment_method" gorm:"type:varchar(50)"`
	OrderNumber string    `json:"order_number" gorm:"type:varchar(100)"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Product     Product   `json:"product" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	User        User      `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type RenewalOrder struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	OriginalOrderID uint    `json:"original_order_id"`
	UserID        uint      `json:"user_id"`
	ProductID     uint      `json:"product_id"`
	ProductName   string    `json:"product_name" gorm:"type:varchar(255)"`
	Amount        float64   `json:"amount"`
	Status        string    `json:"status" gorm:"default:'pending';type:varchar(50)"` // pending, paid, cancelled, completed
	PaymentMethod string    `json:"payment_method" gorm:"type:varchar(50)"`
	OrderNumber   string    `json:"order_number" gorm:"type:varchar(100)"`
	RenewalPeriod int       `json:"renewal_period"` // 续费月数
	ExpireAt      time.Time `json:"expire_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Product       Product   `json:"product" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	User          User      `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}