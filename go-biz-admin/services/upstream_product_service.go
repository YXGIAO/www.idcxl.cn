package services

import (
	"fmt"
	"time"

	"go-biz-admin/config"
	"go-biz-admin/models"
	"gorm.io/gorm"
)

// UpstreamProductService 上游产品服务
type UpstreamProductService struct{}

// NewUpstreamProductService 创建上游产品服务实例
func NewUpstreamProductService() *UpstreamProductService {
	return &UpstreamProductService{}
}

// GetProductList 获取产品列表（支持筛选、分页、排序）
func (s *UpstreamProductService) GetProductList(params models.ProductListRequest) (*models.ProductListResponse, error) {
	db := config.DB.Model(&models.UpstreamProduct{})

	// 构建查询条件
	db = s.buildQueryConditions(db, params)

	// 获取总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("获取总数失败: %v", err)
	}

	// 分页
	offset := (params.Page - 1) * params.Limit
	if offset < 0 {
		offset = 0
	}

	// 排序
	orderBy := params.OrderBy
	if orderBy == "" {
		orderBy = "id"
	}
	sort := params.Sort
	if sort == "" {
		sort = "asc"
	}
	orderClause := fmt.Sprintf("%s %s", orderBy, sort)

	// 查询数据
	var products []models.UpstreamProduct
	if err := db.Preload("Supplier").
		Order(orderClause).
		Limit(params.Limit).
		Offset(offset).
		Find(&products).Error; err != nil {
		return nil, fmt.Errorf("查询产品列表失败: %v", err)
	}

	// 构建响应
	response := &models.ProductListResponse{
		List:  products,
		Count: int(total),
	}

	return response, nil
}

// buildQueryConditions 构建查询条件
func (s *UpstreamProductService) buildQueryConditions(db *gorm.DB, params models.ProductListRequest) *gorm.DB {
	// 关键字搜索
	if params.Keywords != "" {
		keywords := "%" + params.Keywords + "%"
		db = db.Where("id LIKE ? OR username LIKE ? OR email LIKE ? OR phone LIKE ? OR product_name LIKE ? OR name LIKE ?",
			keywords, keywords, keywords, keywords, keywords, keywords)
	}

	// 供应商ID筛选
	if params.SupplierID > 0 {
		db = db.Where("supplier_id = ?", params.SupplierID)
	}

	// 付款周期筛选
	if params.BillingCycle != "" {
		db = db.Where("billing_cycle = ?", params.BillingCycle)
	}

	// 状态筛选
	if params.Status != "" {
		db = db.Where("status = ?", params.Status)
	}

	// 时间范围筛选
	if params.StartTime > 0 {
		startTime := time.Unix(params.StartTime, 0)
		db = db.Where("created_at >= ?", startTime)
	}
	if params.EndTime > 0 {
		endTime := time.Unix(params.EndTime, 0)
		db = db.Where("created_at <= ?", endTime)
	}

	return db
}

// GetProductByID 根据ID获取产品详情
func (s *UpstreamProductService) GetProductByID(id uint) (*models.UpstreamProduct, error) {
	var product models.UpstreamProduct
	if err := config.DB.Preload("Supplier").First(&product, id).Error; err != nil {
		return nil, fmt.Errorf("产品不存在: %v", err)
	}
	return &product, nil
}

// CreateProduct 创建产品
func (s *UpstreamProductService) CreateProduct(product *models.UpstreamProduct) error {
	// 验证必填字段
	if product.Name == "" {
		return fmt.Errorf("产品标识不能为空")
	}
	if product.ProductName == "" {
		return fmt.Errorf("商品名称不能为空")
	}
	if product.SupplierID == 0 {
		return fmt.Errorf("供应商ID不能为空")
	}

	// 设置默认值
	if product.Status == "" {
		product.Status = "Active"
	}
	if product.CreatedAt.IsZero() {
		product.CreatedAt = time.Now()
	}
	if product.UpdatedAt.IsZero() {
		product.UpdatedAt = time.Now()
	}

	return config.DB.Create(product).Error
}

// UpdateProduct 更新产品
func (s *UpstreamProductService) UpdateProduct(id uint, product *models.UpstreamProduct) error {
	// 检查产品是否存在
	var existingProduct models.UpstreamProduct
	if err := config.DB.First(&existingProduct, id).Error; err != nil {
		return fmt.Errorf("产品不存在: %v", err)
	}

	// 更新字段
	product.ID = id
	product.UpdatedAt = time.Now()

	return config.DB.Save(product).Error
}

// DeleteProduct 删除产品
func (s *UpstreamProductService) DeleteProduct(id uint) error {
	// 检查产品是否存在
	var product models.UpstreamProduct
	if err := config.DB.First(&product, id).Error; err != nil {
		return fmt.Errorf("产品不存在: %v", err)
	}

	return config.DB.Delete(&product).Error
}

// GetStatusOptions 获取状态选项
func (s *UpstreamProductService) GetStatusOptions() []map[string]string {
	return []map[string]string{
		{"value": "Unpaid", "label": "未付款"},
		{"value": "Pending", "label": "开通中"},
		{"value": "Active", "label": "已开通"},
		{"value": "Suspended", "label": "已暂停"},
		{"value": "Deleted", "label": "已删除"},
		{"value": "Failed", "label": "开通失败"},
	}
}

// GetBillingCycleOptions 获取付款周期选项
func (s *UpstreamProductService) GetBillingCycleOptions() []map[string]string {
	return []map[string]string{
		{"value": "monthly", "label": "月付"},
		{"value": "quarterly", "label": "季付"},
		{"value": "halfyearly", "label": "半年付"},
		{"value": "yearly", "label": "年付"},
		{"value": "biennially", "label": "两年付"},
		{"value": "triennially", "label": "三年付"},
	}
}