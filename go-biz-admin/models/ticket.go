package models

import "time"

type Ticket struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id"`
	Title       string    `json:"title" gorm:"type:varchar(255)"`
	Content     string    `json:"content" gorm:"type:text"` // 使用text类型存储较长的内容
	Category    string    `json:"category" gorm:"type:varchar(50)"` // technical, billing, sales
	Status      string    `json:"status" gorm:"default:'open';type:varchar(50)"` // open, in_progress, resolved, closed
	Priority    string    `json:"priority" gorm:"default:'medium';type:varchar(50)"` // low, medium, high, urgent
	AssignedTo  *uint     `json:"assigned_to"` // 工单处理人ID
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	User        User      `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type TicketStats struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Date      time.Time `json:"date"`
	Total     int       `json:"total"`
	Open      int       `json:"open"`
	InProgress int      `json:"in_progress"`
	Resolved  int       `json:"resolved"`
	Closed    int       `json:"closed"`
	CreatedAt time.Time `json:"created_at"`
}