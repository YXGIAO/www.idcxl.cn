package models

import (
	"time"
)

// Server 服务器模型
type Server struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	Name            string    `json:"name" gorm:"not null;type:varchar(255)"`
	Host            string    `json:"host" gorm:"type:varchar(255)"`
	Port            int       `json:"port" gorm:"default:22"`
	Type            string    `json:"server_type" gorm:"column:type;type:varchar(50)"` // 服务器类型 (kvm, vps, dedicated, etc.)
	Username        string    `json:"username" gorm:"type:varchar(100)"`
	Password        string    `json:"password" gorm:"type:varchar(255)"`
	Location        string    `json:"location" gorm:"type:varchar(255)"`                      // 服务器位置
	Status          string    `json:"status" gorm:"default:active;type:varchar(50)"`  // active, inactive, maintenance
	CPU             string    `json:"cpu" gorm:"type:varchar(100)"`                           // CPU规格
	Memory          string    `json:"memory" gorm:"type:varchar(100)"`                        // 内存规格
	Disk            string    `json:"disk" gorm:"type:varchar(100)"`                          // 磁盘规格
	Bandwidth       string    `json:"bandwidth" gorm:"type:varchar(100)"`                     // 带宽
	IPCount         int       `json:"ip_count"`                                               // IP数量
	Notes           string    `json:"notes" gorm:"type:text"`                                 // 备注
	Description     string    `json:"description" gorm:"type:text"`                           // 描述
	SupplierID      uint      `json:"supplier_id"`                                            // 供应商ID
	ZJMFServerID    string    `json:"zjmf_server_id" gorm:"type:varchar(100)"`                // 智简魔方服务器ID
	ZJMFServerGroup string    `json:"zjmf_server_group" gorm:"type:varchar(100)"`             // 智简魔方服务器组
	ZJMFStatus      string    `json:"zjmf_status" gorm:"type:varchar(50)"`                    // 智简魔方状态
	ZJMFHost        string    `json:"zjmf_host" gorm:"type:varchar(255)"`                     // 智简魔方主机地址
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Server) TableName() string {
	return "servers"
}