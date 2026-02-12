package services

import (
	"fmt"
	"math/rand"
	"time"

	"go-biz-admin/config"
	"go-biz-admin/models"
)

// OrderService 订单服务
type OrderService struct{}

// NewOrderService 创建订单服务
func NewOrderService() *OrderService {
	return &OrderService{}
}

// GetProductOrders 获取产品订单列表
func (s *OrderService) GetProductOrders(page, limit int, status string, userID *uint) ([]models.ProductOrder, int, error) {
	var orders []models.ProductOrder
	var total int64

	db := config.DB.Model(&models.ProductOrder{}).Preload("User").Preload("Product")

	if status != "" {
		db = db.Where("status = ?", status)
	}

	if userID != nil {
		db = db.Where("user_id = ?", *userID)
	}

	db.Count(&total)

	offset := (page - 1) * limit
	if err := db.Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, int(total), nil
}

// GetProductOrder 获取单个产品订单
func (s *OrderService) GetProductOrder(id uint) (*models.ProductOrder, error) {
	var order models.ProductOrder
	err := config.DB.Preload("User").Preload("Product").First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// CreateProductOrder 创建产品订单
func (s *OrderService) CreateProductOrder(order *models.ProductOrder) error {
	// 检查关联的产品是否存在
	var product models.Product
	if err := config.DB.First(&product, order.ProductID).Error; err != nil {
		return fmt.Errorf("产品不存在")
	}
	order.ProductName = product.Name

	// 设置默认状态
	if order.Status == "" {
		order.Status = "pending"
	}

	// 生成订单号
	order.OrderNumber = generateOrderNumber()

	// 设置时间戳
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	return config.DB.Create(order).Error
}

// UpdateProductOrder 更新产品订单
func (s *OrderService) UpdateProductOrder(order *models.ProductOrder) error {
	order.UpdatedAt = time.Now()
	return config.DB.Save(order).Error
}

// GetRenewalOrders 获取续费订单列表
func (s *OrderService) GetRenewalOrders(page, limit int, status string, userID *uint) ([]models.RenewalOrder, int, error) {
	var orders []models.RenewalOrder
	var total int64

	db := config.DB.Model(&models.RenewalOrder{}).Preload("User").Preload("Product")

	if status != "" {
		db = db.Where("status = ?", status)
	}

	if userID != nil {
		db = db.Where("user_id = ?", *userID)
	}

	db.Count(&total)

	offset := (page - 1) * limit
	if err := db.Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, int(total), nil
}

// GetRenewalOrder 获取单个续费订单
func (s *OrderService) GetRenewalOrder(id uint) (*models.RenewalOrder, error) {
	var order models.RenewalOrder
	err := config.DB.Preload("User").Preload("Product").First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// CreateRenewalOrder 创建续费订单
func (s *OrderService) CreateRenewalOrder(order *models.RenewalOrder) error {
	// 检查关联的产品是否存在
	var product models.Product
	if err := config.DB.First(&product, order.ProductID).Error; err != nil {
		return fmt.Errorf("产品不存在")
	}
	order.ProductName = product.Name

	// 设置默认状态
	if order.Status == "" {
		order.Status = "pending"
	}

	// 生成订单号
	order.OrderNumber = generateOrderNumber()

	// 设置时间戳
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	return config.DB.Create(order).Error
}

// UpdateRenewalOrder 更新续费订单
func (s *OrderService) UpdateRenewalOrder(order *models.RenewalOrder) error {
	order.UpdatedAt = time.Now()
	return config.DB.Save(order).Error
}

// generateOrderNumber 生成订单号
func generateOrderNumber() string {
	now := time.Now()
	return fmt.Sprintf("ORD%s%06d", now.Format("20060102"), rand.Intn(999999))
}