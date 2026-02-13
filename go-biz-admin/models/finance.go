package models

import "time"

type Transaction struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id"`
	OrderID     uint      `json:"order_id"`
	OrderType   string    `json:"order_type" gorm:"type:varchar(50)"` // product_order, renewal_order
	Amount      float64   `json:"amount"`
	Type        string    `json:"type" gorm:"type:varchar(20)"` // income, expense
	Status      string    `json:"status" gorm:"default:'completed';type:varchar(20)"` // pending, completed, failed, refunded
	TransactionNumber string `json:"transaction_number" gorm:"type:varchar(100)"`
	Description string    `json:"description" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	User        User      `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Bill struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id"`
	OrderID     uint      `json:"order_id"`
	Amount      float64   `json:"amount"`
	Status      string    `json:"status" gorm:"default:'unpaid';type:varchar(20)"` // unpaid, paid, overdue
	DueDate     time.Time `json:"due_date"`
	IssueDate   time.Time `json:"issue_date"`
	Description string    `json:"description" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	User        User      `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}