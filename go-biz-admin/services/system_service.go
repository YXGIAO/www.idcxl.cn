package services

import (
	"fmt"
	"time"

	"go-biz-admin/config"
	"go-biz-admin/models"
)

// SystemService 系统服务
type SystemService struct{}

// NewSystemService 创建系统服务
func NewSystemService() *SystemService {
	return &SystemService{}
}

// 系统设置相关
var systemSettings = map[string]string{
	"site_name":        "智简魔方业务管理系统",
	"site_description": "专业的业务管理平台",
	"contact_email":    "admin@example.com",
	"support_email":    "support@example.com",
	"default_timezone": "Asia/Shanghai",
	"maintenance_mode": "false",
}

// GetSystemSettings 获取系统设置
func (s *SystemService) GetSystemSettings() map[string]string {
	return systemSettings
}

// UpdateSystemSettings 更新系统设置
func (s *SystemService) UpdateSystemSettings(settings map[string]string) error {
	for key, value := range settings {
		systemSettings[key] = value
	}
	return nil
}

// 供应商管理相关
// GetSuppliers 获取供应商列表
func (s *SystemService) GetSuppliers(page, limit int, status string) ([]models.Supplier, int, error) {
	var suppliers []models.Supplier
	var total int64

	db := config.DB.Model(&models.Supplier{})

	if status != "" {
		statusInt := 1 // 默认启用
		if status == "active" || status == "1" {
			statusInt = 1
		} else if status == "inactive" || status == "0" {
			statusInt = 0
		}
		db = db.Where("status = ?", statusInt)
	}

	db.Count(&total)

	offset := (page - 1) * limit
	if err := db.Offset(offset).Limit(limit).Find(&suppliers).Error; err != nil {
		return nil, 0, err
	}

	return suppliers, int(total), nil
}

// GetSupplier 获取单个供应商
func (s *SystemService) GetSupplier(id uint) (*models.Supplier, error) {
	var supplier models.Supplier
	err := config.DB.First(&supplier, id).Error
	if err != nil {
		return nil, err
	}
	return &supplier, nil
}

// CreateSupplier 创建供应商
func (s *SystemService) CreateSupplier(supplier *models.Supplier) error {
	if supplier.Status == 0 {
		supplier.Status = 1 // 设置默认状态为启用
	}

	return config.DB.Create(supplier).Error
}

// UpdateSupplier 更新供应商
func (s *SystemService) UpdateSupplier(supplier *models.Supplier) error {
	return config.DB.Save(supplier).Error
}

// DeleteSupplier 删除供应商
func (s *SystemService) DeleteSupplier(id uint) error {
	return config.DB.Delete(&models.Supplier{}, id).Error
}

// GetSupplierDetail 获取供应商详细信息(包含所有字段)
func (s *SystemService) GetSupplierDetail(id uint) (map[string]interface{}, error) {
	var supplier models.Supplier
	err := config.DB.First(&supplier, id).Error
	if err != nil {
		return nil, err
	}

	// 转换为API响应格式
	detail := map[string]interface{}{
		"id":               supplier.ID,
		"type":             supplier.Type,
		"name":             supplier.Name,
		"url":              supplier.URL,
		"username":         supplier.Username,
		"token":            supplier.ZJMFApiKey,
		"secret":           supplier.ZJMFApiSecret,
		"contact":          supplier.Contact,
		"notes":            supplier.Notes,
		"currency_code":    supplier.CurrencyCode,
		"rate":             supplier.Rate,
		"auto_update_rate": supplier.AutoUpdateRate,
		"rate_update_time": supplier.RateUpdateTime,
	}

	return detail, nil
}

// 产品管理相关
// GetProducts 获取产品列表
// GetProducts 获取产品列表
func (s *SystemService) GetProducts(page, limit int, status string, supplierID *uint, category string) ([]models.UpstreamProduct, int, error) {
	var products []models.UpstreamProduct
	var total int64

	db := config.DB.Model(&models.UpstreamProduct{})

	if status != "" {
		db = db.Where("status = ?", status)
	}

	if supplierID != nil {
		db = db.Where("supplier_id = ?", *supplierID)
	}

	if category != "" {
		db = db.Where("category = ?", category)
	}

	db.Count(&total)

	offset := (page - 1) * limit
	if err := db.Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, int(total), nil
}

// GetProduct 获取单个产品
func (s *SystemService) GetProduct(id uint) (*models.UpstreamProduct, error) {
	var product models.UpstreamProduct
	err := config.DB.Preload("Supplier").First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// CreateProduct 创建产品
func (s *SystemService) CreateProduct(product *models.UpstreamProduct) error {
	if product.Status == "" {
		product.Status = "Active"
	}
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	return config.DB.Create(product).Error
}

// UpdateProduct 更新产品
func (s *SystemService) UpdateProduct(product *models.UpstreamProduct) error {
	product.UpdatedAt = time.Now()
	return config.DB.Save(product).Error
}

// DeleteProduct 删除产品
func (s *SystemService) DeleteProduct(id uint) error {
	return config.DB.Delete(&models.UpstreamProduct{}, id).Error
}

// GetSupplierRates 获取所有供应商的汇率信息
func (s *SystemService) GetSupplierRates() ([]models.Supplier, error) {
	var suppliers []models.Supplier
	err := config.DB.Select(
		"id",
		"name",
		"type",
		"currency_code",
		"rate",
		"auto_update_rate",
		"rate_update_time",
		"updated_at",
	).Find(&suppliers).Error
	if err != nil {
		return nil, err
	}
	return suppliers, nil
}

// SyncProductsFromSupplier 从供应商同步产品
func (s *SystemService) SyncProductsFromSupplier(supplierID uint, productsData []map[string]interface{}) (int, error) {
	// 验证供应商是否存在
	_, err := s.GetSupplier(supplierID)
	if err != nil {
		return 0, err
	}

	fmt.Printf("开始同步供应商 %d 的产品，产品数量: %d\n", supplierID, len(productsData))

	syncedCount := 0
	for i, productData := range productsData {
		fmt.Printf("处理第 %d 个产品: %+v\n", i+1, productData)

		// 创建产品模型实例
		localProduct := models.UpstreamProduct{
			Name:        getStringValue(productData, "name", fmt.Sprintf("Product_%d", i+1)),
			ProductName: getStringValue(productData, "product_name", getStringValue(productData, "name", fmt.Sprintf("Product_%d", i+1))),
			Status:      getStringValue(productData, "status", "Active"),
			SupplierID:  supplierID,
		}

		// 检查产品是否已存在于本地数据库
		var existingProduct models.UpstreamProduct
		result := config.DB.Where("name = ? AND supplier_id = ?", localProduct.Name, supplierID).First(&existingProduct)

		if result.Error != nil {
			// 产品不存在，创建新记录
			fmt.Printf("创建新产品: %s\n", localProduct.Name)
			if err := s.CreateProduct(&localProduct); err != nil {
				fmt.Printf("Failed to create product %s: %v\n", localProduct.Name, err)
			} else {
				syncedCount++
			}
		} else {
			// 产品存在，更新记录
			fmt.Printf("更新现有产品: %s\n", existingProduct.Name)
			existingProduct.ProductName = localProduct.ProductName
			existingProduct.Status = localProduct.Status

			if err := s.UpdateProduct(&existingProduct); err != nil {
				fmt.Printf("Failed to update product %s: %v\n", existingProduct.Name, err)
			} else {
				syncedCount++
			}
		}
	}

	fmt.Printf("同步完成，成功处理 %d 个产品\n", syncedCount)
	return syncedCount, nil
}

// 辅助函数：安全获取字符串值
func getStringValue(data map[string]interface{}, key string, defaultValue string) string {
	if value, exists := data[key]; exists {
		if str, ok := value.(string); ok {
			return str
		}
		return fmt.Sprintf("%v", value)
	}
	return defaultValue
}

// 任务管理相关
// GetTasks 获取任务列表
func (s *SystemService) GetTasks(page, limit int, status, taskType string) ([]models.Task, int, error) {
	var tasks []models.Task
	var total int64

	db := config.DB.Model(&models.Task{})

	if status != "" {
		db = db.Where("status = ?", status)
	}

	if taskType != "" {
		db = db.Where("type = ?", taskType)
	}

	db.Count(&total)

	offset := (page - 1) * limit
	if err := db.Offset(offset).Limit(limit).Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	return tasks, int(total), nil
}

// GetTask 获取单个任务
func (s *SystemService) GetTask(id uint) (*models.Task, error) {
	var task models.Task
	err := config.DB.First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// CreateTask 创建任务
func (s *SystemService) CreateTask(task *models.Task) error {
	if task.Status == "" {
		task.Status = "pending"
	}
	task.CreatedAt = time.Now()

	return config.DB.Create(task).Error
}

// UpdateTask 更新任务
func (s *SystemService) UpdateTask(task *models.Task) error {
	task.UpdatedAt = time.Now()
	return config.DB.Save(task).Error
}

// DeleteTask 删除任务
func (s *SystemService) DeleteTask(id uint) error {
	return config.DB.Delete(&models.Task{}, id).Error
}

// RunTask 运行任务
func (s *SystemService) RunTask(id uint) error {
	task := models.Task{}
	if err := config.DB.First(&task, id).Error; err != nil {
		return err
	}

	task.Status = "running"
	now := time.Now()
	task.StartedAt = &now
	config.DB.Save(&task)

	// 模拟任务执行
	// 实际实现中，这里会执行具体的任务逻辑

	task.Status = "completed"
	completedTime := time.Now()
	task.CompletedAt = &completedTime
	return config.DB.Save(&task).Error
}

// 服务器管理相关
// GetServers 获取服务器列表
func (s *SystemService) GetServers(page, limit int, status, serverType string) ([]models.Server, int, error) {
	var servers []models.Server
	var total int64

	db := config.DB.Model(&models.Server{})

	if status != "" {
		db = db.Where("status = ?", status)
	}

	if serverType != "" {
		db = db.Where("type = ?", serverType)
	}

	db.Count(&total)

	offset := (page - 1) * limit
	if err := db.Offset(offset).Limit(limit).Find(&servers).Error; err != nil {
		return nil, 0, err
	}

	return servers, int(total), nil
}

// GetServer 获取单个服务器
func (s *SystemService) GetServer(id uint) (*models.Server, error) {
	var server models.Server
	err := config.DB.First(&server, id).Error
	if err != nil {
		return nil, err
	}
	return &server, nil
}

// CreateServer 创建服务器
func (s *SystemService) CreateServer(server *models.Server) error {
	if server.Status == "" {
		server.Status = "active"
	}
	server.CreatedAt = time.Now()
	server.UpdatedAt = time.Now()

	return config.DB.Create(server).Error
}

// UpdateServer 更新服务器
func (s *SystemService) UpdateServer(server *models.Server) error {
	server.UpdatedAt = time.Now()
	return config.DB.Save(server).Error
}

// DeleteServer 删除服务器
func (s *SystemService) DeleteServer(id uint) error {
	return config.DB.Delete(&models.Server{}, id).Error
}
