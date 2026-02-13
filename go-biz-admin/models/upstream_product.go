package models

import (
	"time"
	"gorm.io/gorm"
)

// UpstreamProduct 上游产品模型（对应智简魔方等上游供应商的产品实例）
type UpstreamProduct struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"` // 软删除功能
	
	// 产品基本信息
	Name         string `json:"name" gorm:"type:varchar(100);not null;comment:产品标识"`
	ProductName  string `json:"product_name" gorm:"type:varchar(200);comment:商品名称"`
	Status       string `json:"status" gorm:"type:varchar(50);default:'Active';comment:状态:Unpaid未付款,Pending开通中,Active已开通,Suspended已暂停,Deleted已删除,Failed开通失败"`
	
	// 财务信息
	FirstPaymentAmount string `json:"first_payment_amount" gorm:"type:varchar(50);comment:首次付款金额"`
	RenewAmount        string `json:"renew_amount" gorm:"type:varchar(50);comment:续费金额"`
	BillingCycle       string `json:"billing_cycle" gorm:"type:varchar(50);comment:付款周期:monthly,yearly等"`
	BillingCycleName   string `json:"billing_cycle_name" gorm:"type:varchar(50);comment:模块计费周期名称"`
	BillingCycleTime   string `json:"billing_cycle_time" gorm:"type:varchar(50);comment:模块计费周期时间"`
	
	// 时间信息（Unix时间戳）
	ActiveTime int64 `json:"active_time" gorm:"comment:开通时间(Unix时间戳)"`
	DueTime    int64 `json:"due_time" gorm:"comment:到期时间(Unix时间戳)"`
	
	// 上游产品信息
	UpstreamHostID int `json:"upstream_host_id" gorm:"comment:上游产品ID"`
	
	// 用户信息
	ClientID   string `json:"client_id" gorm:"type:varchar(100);not null;comment:用户ID"`
	Username   string `json:"username" gorm:"type:varchar(100);comment:用户名"`
	Company    string `json:"company" gorm:"type:varchar(200);comment:公司名称"`
	Email      string `json:"email" gorm:"type:varchar(100);comment:邮箱"`
	PhoneCode  string `json:"phone_code" gorm:"type:varchar(20);comment:国际电话区号"`
	Phone      string `json:"phone" gorm:"type:varchar(50);comment:手机号"`
	
	// 配置信息
	IpNum    int    `json:"ip_num" gorm:"comment:IP数量"`
	BaseInfo string `json:"base_info" gorm:"type:varchar(1024);comment:产品基础信息"`
	
	// 供应商信息
	SupplierID uint `json:"supplier_id" gorm:"comment:供应商ID"`
	Supplier   Supplier `json:"supplier" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// TableName 指定表名
func (UpstreamProduct) TableName() string {
	return "upstream_products"
}

// ProductListRequest 产品列表请求参数
type ProductListRequest struct {
	Keywords     string `form:"keywords"`
	SupplierID   int    `form:"supplier_id"`
	BillingCycle string `form:"billing_cycle"`
	Status       string `form:"status"`
	StartTime    int64  `form:"start_time"`
	EndTime      int64  `form:"end_time"`
	Page         int    `form:"page" default:"1"`
	Limit        int    `form:"limit" default:"10"`
	OrderBy      string `form:"orderby" default:"id"`
	Sort         string `form:"sort" default:"asc"`
}

// ProductListResponse 产品列表响应
type ProductListResponse struct {
	List  []UpstreamProduct `json:"list"`
	Count int               `json:"count"`
}

// ProductDetailResponse 产品详情响应
type ProductDetailResponse struct {
	Host UpstreamProduct `json:"host"`
}