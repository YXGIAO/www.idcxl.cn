package models

import (
	"time"
)

// Supplier 供应商模型
type Supplier struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	Type            string    `json:"type" gorm:"default:'default'"` // 供应商类型: default/whmcs/finance
	Name            string    `json:"name" gorm:"not null;type:varchar(255)"`
	Description     string    `json:"description" gorm:"type:text"`
	Status          int       `json:"status" gorm:"default:1"` // 0:禁用, 1:启用
	URL             string    `json:"url" gorm:"type:varchar(255)"`                     // 供应商系统链接地址
	Username        string    `json:"username" gorm:"type:varchar(100)"`                // 登录用户名
	Contact         string    `json:"contact" gorm:"type:varchar(100)"`                 // 联系方式
	Notes           string    `json:"notes" gorm:"type:text"`                           // 备注信息
	CurrencyCode    string    `json:"currency_code" gorm:"default:'CNY'"` // 货币代码
	Rate            string    `json:"rate" gorm:"default:'1.0000'"`       // 汇率
	AutoUpdateRate  int       `json:"auto_update_rate" gorm:"default:0"`  // 是否自动更新汇率
	RateUpdateTime  int64     `json:"rate_update_time"`        // 汇率更新时间戳
	
	// 智简魔方专用字段
	ZJMFServerID    string    `json:"zjmf_server_id" gorm:"type:varchar(100)"`          
	ZJMFServerGroup string    `json:"zjmf_server_group" gorm:"type:varchar(100)"`       
	ZJMFHost        string    `json:"zjmf_host" gorm:"type:varchar(100)"`               
	ZJMFApiEndpoint string    `json:"zjmf_api_endpoint" gorm:"type:varchar(255)"`       
	ZJMFApiKey      string    `json:"zjmf_api_key" gorm:"type:varchar(255)"`            
	ZJMFApiSecret   string    `json:"zjmf_api_secret" gorm:"type:varchar(255)"`         
	ZJMFStatus      string    `json:"zjmf_status" gorm:"type:varchar(50)"`             
	ZJMFAccount     string    `json:"zjmf_account" gorm:"type:varchar(100)"`            
	
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// TableName sets the table name for Supplier
func (Supplier) TableName() string {
	return "suppliers"
}

// SupplierRateItem 供应商汇率项
type SupplierRateItem struct {
	ID             uint   `json:"id"`
	Name           string `json:"name" gorm:"type:varchar(255)"`
	Type           string `json:"type" gorm:"type:varchar(50)"`
	CurrencyCode   string `json:"currency_code" gorm:"type:varchar(10)"`
	Rate           string `json:"rate" gorm:"type:varchar(20)"`
	AutoUpdateRate int    `json:"auto_update_rate"`
	RateUpdateTime int64  `json:"rate_update_time"`
	UpdatedAt      string `json:"updated_at"`
}

// SupplierRateResponse 供应商汇率响应
type SupplierRateResponse struct {
	Code    int                `json:"code"`
	Message string             `json:"message"`
	Data    []SupplierRateItem `json:"data"`
}

// SupplierCreateRequest 供应商创建请求
type SupplierCreateRequest struct {
	Type            string `json:"type" binding:"required,oneof=default whmcs finance"`
	Name            string `json:"name" binding:"required,max=100"`
	Description     string `json:"description" binding:"max=500"`
	Status          int    `json:"status" binding:"oneof=0 1"`
	URL             string `json:"url" binding:"required,url"`
	Username        string `json:"username" binding:"required"`
	Contact         string `json:"contact" binding:"required,email"`
	Notes           string `json:"notes"`
	CurrencyCode    string `json:"currency_code" binding:"required,len=3"`
	Rate            string `json:"rate" binding:"required"`
	AutoUpdateRate  int    `json:"auto_update_rate" binding:"oneof=0 1"`
	
	// 智简魔方专用字段
	ZJMFServerID    string `json:"zjmf_server_id,omitempty"`
	ZJMFServerGroup string `json:"zjmf_server_group,omitempty"`
	ZJMFHost        string `json:"zjmf_host,omitempty"`
	ZJMFApiEndpoint string `json:"zjmf_api_endpoint,omitempty"`
	ZJMFApiKey      string `json:"zjmf_api_key" binding:"required"`
	ZJMFApiSecret   string `json:"zjmf_api_secret" binding:"required"`
}

// SupplierUpdateRequest 供应商更新请求
type SupplierUpdateRequest struct {
	Type            string `json:"type" binding:"omitempty,oneof=default whmcs finance"`
	Name            string `json:"name" binding:"omitempty,max=100"`
	Description     string `json:"description" binding:"omitempty,max=500"`
	Status          int    `json:"status" binding:"omitempty,oneof=0 1"`
	URL             string `json:"url" binding:"omitempty,url"`
	Username        string `json:"username" binding:"omitempty"`
	Contact         string `json:"contact" binding:"omitempty,email"`
	Notes           string `json:"notes" binding:"omitempty"`
	CurrencyCode    string `json:"currency_code" binding:"omitempty,len=3"`
	Rate            string `json:"rate" binding:"omitempty"`
	AutoUpdateRate  int    `json:"auto_update_rate" binding:"omitempty,oneof=0 1"`
	
	// 智简魔方专用字段
	ZJMFServerID    string `json:"zjmf_server_id,omitempty"`
	ZJMFServerGroup string `json:"zjmf_server_group,omitempty"`
	ZJMFHost        string `json:"zjmf_host,omitempty"`
	ZJMFApiEndpoint string `json:"zjmf_api_endpoint,omitempty"`
	ZJMFApiKey      string `json:"zjmf_api_key" binding:"omitempty"`
	ZJMFApiSecret   string `json:"zjmf_api_secret" binding:"omitempty"`
}

// SupplierStatusResponse 供应商状态响应
type SupplierStatusResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    StatusData  `json:"data"`
}

// StatusData 状态详情
type StatusData struct {
	SupplierID     int    `json:"supplier_id"`
	Status         string `json:"status"` 
	LastCheckTime  string `json:"last_check_time"`
	ResponseTimeMS int    `json:"response_time_ms"`
	Error          string `json:"error,omitempty"`
}

// SupplierResponse 供应商响应
type SupplierResponse struct {
	ID              uint      `json:"id"`
	Type            string    `json:"type"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Status          int       `json:"status"`
	URL             string    `json:"url"`
	Username        string    `json:"username"`
	Contact         string    `json:"contact"`
	Notes           string    `json:"notes"`
	CurrencyCode    string    `json:"currency_code"`
	Rate            string    `json:"rate"`
	AutoUpdateRate  int       `json:"auto_update_rate"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}