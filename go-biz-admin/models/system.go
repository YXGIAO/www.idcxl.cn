package models

import (
	"time"
)

type Task struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Type        string     `json:"type"`                            // sync_data, cleanup, notification, etc
	Status      string     `json:"status" gorm:"default:'pending'"` // pending, running, completed, failed
	CreatedBy   uint       `json:"created_by"`                      // 创建任务的用户ID
	ServerID    *uint      `json:"server_id"`                       // 如果任务与特定服务器关联
	Progress    int        `json:"progress"`                        // 任务进度百分比
	StartedAt   *time.Time `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}