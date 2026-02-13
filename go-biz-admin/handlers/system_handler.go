package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"go-biz-admin/models"
	"go-biz-admin/services"
)

type SystemHandler struct {
	systemService *services.SystemService
}

func NewSystemHandler(systemService *services.SystemService) *SystemHandler {
	return &SystemHandler{
		systemService: systemService,
	}
}

// GetSuppliers godoc
// @Summary Get suppliers
// @Description Get suppliers
// @Tags suppliers
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /suppliers [get]
func (h *SystemHandler) GetSuppliers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")

	suppliers, total, err := h.systemService.GetSuppliers(page, limit, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"suppliers": suppliers,
		"total":     total,
	})
}

// GetSupplier godoc
// @Summary Get supplier
// @Description Get supplier by ID
// @Tags suppliers
// @Accept json
// @Produce json
// @Param id path int true "Supplier ID"
// @Success 200 {object} models.Supplier
// @Router /suppliers/{id} [get]
func (h *SystemHandler) GetSupplier(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid supplier ID"})
		return
	}

	supplier, err := h.systemService.GetSupplier(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Supplier not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, supplier)
}

// CreateSupplier godoc
// @Summary Create supplier
// @Description Create a new supplier
// @Tags suppliers
// @Accept json
// @Produce json
// @Param supplier body models.Supplier true "Supplier"
// @Success 200 {object} models.Supplier
// @Router /suppliers [post]
// CreateSupplier godoc
// @Summary 创建供应商(管理员)
// @Description 创建一个新的供应商记录
// @Tags admin_suppliers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param supplier body models.SupplierCreateRequest true "供应商信息"
// @Success 201 {object} models.SupplierResponse
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /admin/v1/supplier [post]
func (h *SystemHandler) CreateSupplier(c *gin.Context) {
	// 验证管理员权限
	if !h.isAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    http.StatusForbidden,
			"message": "权限不足，需要管理员权限",
		})
		return
	}

	// 解析请求参数
	var req models.SupplierCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "参数错误",
			"errors":  []string{err.Error()},
		})
		return
	}

	// 转换为模型对象
	supplier := models.Supplier{
		Type:            req.Type,
		Name:            req.Name,
		Description:     req.Description,
		Status:          req.Status,
		URL:             req.URL,
		Username:        req.Username,
		Contact:         req.Contact,
		Notes:           req.Notes,
		CurrencyCode:    req.CurrencyCode,
		Rate:            req.Rate,
		AutoUpdateRate:  req.AutoUpdateRate,
		ZJMFServerID:    req.ZJMFServerID,
		ZJMFServerGroup: req.ZJMFServerGroup,
		ZJMFHost:        req.ZJMFHost,
		ZJMFApiEndpoint: req.ZJMFApiEndpoint,
		ZJMFApiKey:      req.ZJMFApiKey,
		ZJMFApiSecret:   req.ZJMFApiSecret,
	}

	// 创建供应商
	if err := h.systemService.CreateSupplier(&supplier); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "创建供应商失败",
			"errors":  []string{err.Error()},
		})
		return
	}

	// 转换为响应对象
	response := models.SupplierResponse{
		ID:             supplier.ID,
		Type:           supplier.Type,
		Name:           supplier.Name,
		Description:    supplier.Description,
		Status:         supplier.Status,
		URL:            supplier.URL,
		Username:       supplier.Username,
		Contact:        supplier.Contact,
		Notes:          supplier.Notes,
		CurrencyCode:   supplier.CurrencyCode,
		Rate:           supplier.Rate,
		AutoUpdateRate: supplier.AutoUpdateRate,
		CreatedAt:      supplier.CreatedAt,
		UpdatedAt:      supplier.UpdatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

// UpdateSupplier godoc
// @Summary 更新供应商信息(管理员)
// @Description 更新指定ID的供应商信息
// @Tags admin_suppliers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Supplier ID"
// @Param supplier body models.SupplierUpdateRequest true "供应商更新信息"
// @Success 200 {object} models.SupplierResponse
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /admin/v1/supplier/{id} [put]
func (h *SystemHandler) UpdateSupplier(c *gin.Context) {
	// 验证管理员权限
	if !h.isAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    http.StatusForbidden,
			"message": "权限不足，需要管理员权限",
		})
		return
	}

	// 获取供应商ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "无效的供应商ID",
		})
		return
	}

	// 解析请求参数
	var req models.SupplierUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "参数错误",
			"errors":  []string{err.Error()},
		})
		return
	}

	// 获取现有供应商信息
	existingSupplier, err := h.systemService.GetSupplier(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    http.StatusNotFound,
				"message": "供应商不存在",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "获取供应商信息失败",
			})
		}
		return
	}

	// 更新字段（仅更新非零值字段）
	if req.Type != "" {
		existingSupplier.Type = req.Type
	}
	if req.Name != "" {
		existingSupplier.Name = req.Name
	}
	if req.Description != "" {
		existingSupplier.Description = req.Description
	}
	if req.Status != 0 {
		existingSupplier.Status = req.Status
	}
	if req.URL != "" {
		existingSupplier.URL = req.URL
	}
	if req.Username != "" {
		existingSupplier.Username = req.Username
	}
	if req.Contact != "" {
		existingSupplier.Contact = req.Contact
	}
	if req.Notes != "" {
		existingSupplier.Notes = req.Notes
	}
	if req.CurrencyCode != "" {
		existingSupplier.CurrencyCode = req.CurrencyCode
	}
	if req.Rate != "" {
		existingSupplier.Rate = req.Rate
	}
	if req.AutoUpdateRate != 0 {
		existingSupplier.AutoUpdateRate = req.AutoUpdateRate
	}
	if req.ZJMFServerID != "" {
		existingSupplier.ZJMFServerID = req.ZJMFServerID
	}
	if req.ZJMFServerGroup != "" {
		existingSupplier.ZJMFServerGroup = req.ZJMFServerGroup
	}
	if req.ZJMFHost != "" {
		existingSupplier.ZJMFHost = req.ZJMFHost
	}
	if req.ZJMFApiEndpoint != "" {
		existingSupplier.ZJMFApiEndpoint = req.ZJMFApiEndpoint
	}
	if req.ZJMFApiKey != "" {
		existingSupplier.ZJMFApiKey = req.ZJMFApiKey
	}
	if req.ZJMFApiSecret != "" {
		existingSupplier.ZJMFApiSecret = req.ZJMFApiSecret
	}

	// 更新供应商
	if err := h.systemService.UpdateSupplier(existingSupplier); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "更新供应商失败",
			"errors":  []string{err.Error()},
		})
		return
	}

	// 转换为响应对象
	response := models.SupplierResponse{
		ID:             existingSupplier.ID,
		Type:           existingSupplier.Type,
		Name:           existingSupplier.Name,
		Description:    existingSupplier.Description,
		Status:         existingSupplier.Status,
		URL:            existingSupplier.URL,
		Username:       existingSupplier.Username,
		Contact:        existingSupplier.Contact,
		Notes:          existingSupplier.Notes,
		CurrencyCode:   existingSupplier.CurrencyCode,
		Rate:           existingSupplier.Rate,
		AutoUpdateRate: existingSupplier.AutoUpdateRate,
		CreatedAt:      existingSupplier.CreatedAt,
		UpdatedAt:      existingSupplier.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// DeleteSupplier godoc
// @Summary Delete supplier
// @Description Delete supplier by ID
// @Tags suppliers
// @Accept json
// @Produce json
// @Param id path int true "Supplier ID"
// @Success 200 {object} map[string]string
// @Router /suppliers/{id} [delete]
func (h *SystemHandler) DeleteSupplier(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid supplier ID"})
		return
	}

	err = h.systemService.DeleteSupplier(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Supplier deleted successfully"})
}

// GetSupplierDetail godoc
// @Summary Get supplier detail
// @Description Get full details of a supplier by ID
// @Tags suppliers
// @Accept json
// @Produce json
// @Param id path int true "Supplier ID"
// @Success 200 {object} map[string]interface{}
// @Router /suppliers/{id}/detail [get]
func (h *SystemHandler) GetSupplierDetail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid supplier ID"})
		return
	}

	// 检查用户权限（示例，需根据实际权限系统调整）
	if !h.isAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	supplier, err := h.systemService.GetSupplierDetail(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Supplier not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// 可选：对敏感字段进行脱敏处理
	if token, ok := supplier["token"]; ok {
		supplier["token"] = maskSensitiveData(fmt.Sprintf("%v", token))
	}
	if secret, ok := supplier["secret"]; ok {
		supplier["secret"] = maskSensitiveData(fmt.Sprintf("%v", secret))
	}

	c.JSON(http.StatusOK, gin.H{"supplier": supplier})
}

// isAdmin 检查用户是否为管理员（示例函数，需根据实际权限系统实现）
func (h *SystemHandler) isAdmin(c *gin.Context) bool {
	// 实际实现中应从JWT或会话中获取用户角色
	return true
}

// maskSensitiveData 对敏感数据进行脱敏处理
func maskSensitiveData(data string) string {
	if len(data) <= 4 {
		return "****"
	}
	return data[:2] + "****" + data[len(data)-2:]
}

// GetProducts godoc
// @Summary Get products
// @Description Get products
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /products [get]
func (h *SystemHandler) GetProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")
	category := c.Query("category")

	var supplierID *uint
	if supplierIDStr := c.Query("supplier_id"); supplierIDStr != "" {
		id, err := strconv.Atoi(supplierIDStr)
		if err == nil {
			supplierID = uintPtr(uint(id))
		}
	}

	products, total, err := h.systemService.GetProducts(page, limit, status, supplierID, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"total":    total,
	})
}

// AdminGetSupplierDetail godoc
// @Summary Get supplier detail (admin)
// @Description Get supplier details by ID for admin
// @Tags admin_suppliers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Supplier ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/v1/supplier/{id} [get]
func (h *SystemHandler) AdminGetSupplierDetail(c *gin.Context) {
	// 验证管理员权限
	if !h.isAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	// 获取供应商ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid supplier ID"})
		return
	}

	// 查询供应商
	supplier, err := h.systemService.GetSupplier(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Supplier not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// 组织响应数据
	response := gin.H{
		"supplier": gin.H{
			"id":               supplier.ID,
			"type":             supplier.Type,
			"name":             supplier.Name,
			"url":              supplier.URL,
			"username":         supplier.Username,
			"token":            supplier.ZJMFApiKey,    // 使用智简魔方的API密钥作为token
			"secret":           supplier.ZJMFApiSecret, // 使用智简魔方的API密钥作为secret
			"contact":          supplier.Contact,
			"notes":            supplier.Notes,
			"currency_code":    supplier.CurrencyCode,
			"rate":             supplier.Rate,
			"auto_update_rate": supplier.AutoUpdateRate,
			"rate_update_time": supplier.RateUpdateTime,
		},
	}

	c.JSON(http.StatusOK, response)
}

// GetProduct godoc
// @Summary Get product
// @Description Get product by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} models.UpstreamProduct
// @Router /products/{id} [get]
func (h *SystemHandler) GetProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	product, err := h.systemService.GetProduct(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, product)
}

// CreateProduct godoc
// @Summary Create product
// @Description Create a new product
// @Tags products
// @Accept json
// @Produce json
// @Param product body models.UpstreamProduct true "Product"
// @Success 200 {object} models.UpstreamProduct
// @Router /products [post]
func (h *SystemHandler) CreateProduct(c *gin.Context) {
	var product models.UpstreamProduct
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.systemService.CreateProduct(&product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// UpdateProduct godoc
// @Summary Update product
// @Description Update product by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body models.UpstreamProduct true "Product"
// @Success 200 {object} models.UpstreamProduct
// @Router /products/{id} [put]
func (h *SystemHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product models.UpstreamProduct
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 保留ID以更新特定记录
	product.ID = uint(id)
	err = h.systemService.UpdateProduct(&product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct godoc
// @Summary Delete product
// @Description Delete product by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]string
// @Router /products/{id} [delete]
func (h *SystemHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	err = h.systemService.DeleteProduct(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

// GetServers godoc
// @Summary Get servers
// @Description Get servers
// @Tags servers
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /servers [get]
func (h *SystemHandler) GetServers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")
	serverType := c.Query("type")

	servers, total, err := h.systemService.GetServers(page, limit, status, serverType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"servers": servers,
		"total":   total,
	})
}

// GetServer godoc
// @Summary Get server
// @Description Get server by ID
// @Tags servers
// @Accept json
// @Produce json
// @Param id path int true "Server ID"
// @Success 200 {object} models.Server
// @Router /servers/{id} [get]
func (h *SystemHandler) GetServer(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid server ID"})
		return
	}

	server, err := h.systemService.GetServer(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, server)
}

// CreateServer godoc
// @Summary Create server
// @Description Create a new server
// @Tags servers
// @Accept json
// @Produce json
// @Param server body models.Server true "Server"
// @Success 200 {object} models.Server
// @Router /servers [post]
func (h *SystemHandler) CreateServer(c *gin.Context) {
	var server models.Server
	if err := c.ShouldBindJSON(&server); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.systemService.CreateServer(&server)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, server)
}

// UpdateServer godoc
// @Summary Update server
// @Description Update server by ID
// @Tags servers
// @Accept json
// @Produce json
// @Param id path int true "Server ID"
// @Param server body models.Server true "Server"
// @Success 200 {object} models.Server
// @Router /servers/{id} [put]
func (h *SystemHandler) UpdateServer(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid server ID"})
		return
	}

	var server models.Server
	if err := c.ShouldBindJSON(&server); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 保留ID以更新特定记录
	server.ID = uint(id)
	err = h.systemService.UpdateServer(&server)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, server)
}

// DeleteServer godoc
// @Summary Delete server
// @Description Delete server by ID
// @Tags servers
// @Accept json
// @Produce json
// @Param id path int true "Server ID"
// @Success 200 {object} map[string]string
// @Router /servers/{id} [delete]
func (h *SystemHandler) DeleteServer(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid server ID"})
		return
	}

	err = h.systemService.DeleteServer(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Server deleted successfully"})
}

// GetTasks godoc
// @Summary Get tasks
// @Description Get tasks
// @Tags tasks
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /tasks [get]
func (h *SystemHandler) GetTasks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")
	taskType := c.Query("type")

	tasks, total, err := h.systemService.GetTasks(page, limit, status, taskType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tasks": tasks,
		"total": total,
	})
}

// GetTask godoc
// @Summary Get task
// @Description Get task by ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} models.Task
// @Router /tasks/{id} [get]
func (h *SystemHandler) GetTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	task, err := h.systemService.GetTask(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, task)
}

// CreateTask godoc
// @Summary Create task
// @Description Create a new task
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body models.Task true "Task"
// @Success 200 {object} models.Task
// @Router /tasks [post]
func (h *SystemHandler) CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.systemService.CreateTask(&task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

// UpdateTask godoc
// @Summary Update task
// @Description Update task by ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param task body models.Task true "Task"
// @Success 200 {object} models.Task
// @Router /tasks/{id} [put]
func (h *SystemHandler) UpdateTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 保留ID以更新特定记录
	task.ID = uint(id)
	err = h.systemService.UpdateTask(&task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

// DeleteTask godoc
// @Summary Delete task
// @Description Delete task by ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} map[string]string
// @Router /tasks/{id} [delete]
func (h *SystemHandler) DeleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	err = h.systemService.DeleteTask(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

// RunTask godoc
// @Summary Run task
// @Description Run task by ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} map[string]string
// @Router /tasks/{id}/run [post]
func (h *SystemHandler) RunTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	err = h.systemService.RunTask(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task executed successfully"})
}

// GetSystemSettings 获取系统设置
func (h *SystemHandler) GetSystemSettings(c *gin.Context) {
	settings := h.systemService.GetSystemSettings()
	c.JSON(http.StatusOK, settings)
}

// UpdateSystemSettings 更新系统设置
func (h *SystemHandler) UpdateSystemSettings(c *gin.Context) {
	var settings map[string]string
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.systemService.UpdateSystemSettings(settings)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Settings updated successfully"})
}

// GetSupplierRates godoc
// @Summary 获取供应商汇率列表(管理员)
// @Description 获取所有供应商的汇率信息
// @Tags admin_suppliers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.SupplierRateResponse
// @Failure 403 {object} ErrorResponse
// @Router /admin/v1/supplier/rate [get]
func (h *SystemHandler) GetSupplierRates(c *gin.Context) {
	// 验证管理员权限
	if !h.isAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    http.StatusForbidden,
			"message": "权限不足，需要管理员权限",
		})
		return
	}

	// 获取所有供应商的汇率信息
	suppliers, err := h.systemService.GetSupplierRates()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "获取供应商汇率失败",
			"errors":  []string{err.Error()},
		})
		return
	}

	// 转换为响应对象
	response := models.SupplierRateResponse{
		Code:    http.StatusOK,
		Message: "获取成功",
	}
	response.Data = make([]models.SupplierRateItem, len(suppliers))
	for i, supplier := range suppliers {
		response.Data[i] = models.SupplierRateItem{
			ID:             supplier.ID,
			Name:           supplier.Name,
			Type:           supplier.Type,
			CurrencyCode:   supplier.CurrencyCode,
			Rate:           supplier.Rate,
			AutoUpdateRate: supplier.AutoUpdateRate,
			RateUpdateTime: supplier.RateUpdateTime,
			UpdatedAt:      supplier.UpdatedAt.Format(time.RFC3339),
		}
	}

	c.JSON(http.StatusOK, response)
}

// 新增：同步供应商信息的处理函数
func (h *SystemHandler) SyncSupplierInfo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid supplier ID"})
		return
	}

	// 从数据库获取供应商信息
	supplier, err := h.systemService.GetSupplier(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Supplier not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// 根据供应商类型调用不同的API获取信息
	var info map[string]interface{}
	// 使用供应商的字段判断是否为智简魔方类型
	if supplier.ZJMFApiEndpoint != "" || strings.Contains(strings.ToLower(supplier.Name), "智简魔方") {
		info, err = h.getZJMFInfo(*supplier)
	} else {
		// 其他类型的供应商可以添加相应的API调用
		info = map[string]interface{}{
			"balance":                  0, // 供应商模型中没有余额字段
			"available_products_count": 0,
			"total_products_count":     0,
			"normal_products_count":    0,
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 注意：供应商模型中没有这些字段，所以暂时不更新数据库
	// 如果需要存储这些信息，需要扩展Supplier模型

	c.JSON(http.StatusOK, info)
}

// 新增：获取智简魔方信息的辅助函数
func (h *SystemHandler) getZJMFInfo(supplier models.Supplier) (map[string]interface{}, error) {
	// 使用ZJMFService来调用API
	zjmfService := services.NewZJMFService(supplier.ZJMFHost, supplier.ZJMFApiKey, supplier.ZJMFApiSecret)

	// 从智简魔方API获取供应商信息（余额和产品）
	info, err := zjmfService.GetSupplierInfo()
	if err != nil {
		return nil, fmt.Errorf("获取智简魔方供应商信息失败: %v", err)
	}

	// 处理并标准化返回的数据
	result := map[string]interface{}{
		"balance":                  extractValue(info, "balance", "data.balance", "result.balance"),
		"available_products_count": extractValue(info, "available_products", "data.available_products", "result.available_products"),
		"total_products_count":     extractValue(info, "total_products", "data.total_products", "result.total_products"),
		"normal_products_count":    extractValue(info, "normal_products", "data.normal_products", "result.normal_products"),
		"status":                   extractValue(info, "status", "data.status", "result.status"),
	}

	// 如果某些值没有找到，使用默认值
	if result["balance"] == nil {
		result["balance"] = 0
	}
	if result["available_products_count"] == nil {
		result["available_products_count"] = 0
	}
	if result["total_products_count"] == nil {
		result["total_products_count"] = 0
	}
	if result["normal_products_count"] == nil {
		result["normal_products_count"] = 0
	}
	if result["status"] == nil {
		result["status"] = "active"
	}

	return result, nil
}

// 辅助函数：从响应中提取值
func extractValue(response map[string]interface{}, keys ...string) interface{} {
	for _, key := range keys {
		parts := strings.Split(key, ".")
		currentValue := interface{}(response) // 将初始值转换为interface{}

		for _, part := range parts {
			if m, ok := currentValue.(map[string]interface{}); ok {
				currentValue = m[part]
				if currentValue == nil {
					break
				}
			} else {
				break
			}
		}

		if currentValue != nil {
			return currentValue
		}
	}
	return nil
}

// 辅助函数：创建uint指针
func uintPtr(u uint) *uint {
	return &u
}

// 辅助函数：将字符串转换为float64
func parseFloat64(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0
	}
	return f
}

// 新增：从供应商同步产品的处理函数
func (h *SystemHandler) SyncProductsFromSupplier(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid supplier ID"})
		return
	}

	// 从数据库获取供应商信息
	supplier, err := h.systemService.GetSupplier(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Supplier not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// 检查供应商类型是否支持同步产品
	if supplier.ZJMFApiEndpoint == "" && !strings.Contains(strings.ToLower(supplier.Name), "智简魔方") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该供应商类型不支持同步产品"})
		return
	}

	// 调用智简魔方API获取产品信息
	zjmfService := services.NewZJMFService(supplier.ZJMFHost, supplier.ZJMFApiKey, supplier.ZJMFApiSecret)
	fmt.Printf("调用智简魔方API: Host=%s, Key=%s\n", supplier.ZJMFHost, supplier.ZJMFApiKey)

	products, err := zjmfService.GetProducts()
	if err != nil {
		fmt.Printf("获取产品信息失败: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("获取供应商产品信息失败: %v", err)})
		return
	}

	fmt.Printf("获取到的产品数据: %+v\n", products)

	// 解析产品数据并保存到本地数据库
	var productsData []map[string]interface{}

	// 尝试多种可能的产品数据格式
	if productList, ok := products["products"].([]interface{}); ok {
		for _, prod := range productList {
			if productMap, ok := prod.(map[string]interface{}); ok {
				productsData = append(productsData, productMap)
			}
		}
	} else if productList, ok := products["data"].([]interface{}); ok {
		// 尝试 data 字段
		for _, prod := range productList {
			if productMap, ok := prod.(map[string]interface{}); ok {
				productsData = append(productsData, productMap)
			}
		}
	} else if productList, ok := products["result"].([]interface{}); ok {
		// 尝试 result 字段
		for _, prod := range productList {
			if productMap, ok := prod.(map[string]interface{}); ok {
				productsData = append(productsData, productMap)
			}
		}
	} else {
		// 如果都不是数组格式，尝试直接使用整个响应
		fmt.Printf("产品数据不是数组格式，原始数据: %+v\n", products)
		// 尝试将整个map作为单个产品
		if len(products) > 0 {
			productsData = append(productsData, products)
		}
	}

	fmt.Printf("解析后的产品数据数量: %d\n", len(productsData))
	if len(productsData) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "没有获取到产品数据", "syncedCount": 0})
		return
	}

	// 使用service层同步产品
	syncedCount, err := h.systemService.SyncProductsFromSupplier(uint(id), productsData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("同步产品失败: %v", err)})
		return
	}

	response := gin.H{
		"message":     fmt.Sprintf("成功同步 %d 个产品", syncedCount),
		"syncedCount": syncedCount,
	}
	c.JSON(http.StatusOK, response)
}
