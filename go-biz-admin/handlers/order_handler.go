package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-biz-admin/models"
	"go-biz-admin/services"
)

// OrderHandler 订单处理器
type OrderHandler struct {
	service *services.OrderService
}

// NewOrderHandler 创建订单处理器
func NewOrderHandler(service *services.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

// GetProductOrders 获取产品订单列表
func (h *OrderHandler) GetProductOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")
	userIdStr := c.Query("user_id")
	
	var userId *uint
	if userIdStr != "" {
		uid, err := strconv.Atoi(userIdStr)
		if err == nil {
			temp := uint(uid)
			userId = &temp
		}
	}

	orders, total, err := h.service.GetProductOrders(page, limit, status, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"total":  total,
		"page":   page,
		"limit":  limit,
	})
}

// GetProductOrder 获取单个产品订单
func (h *OrderHandler) GetProductOrder(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	order, err := h.service.GetProductOrder(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// CreateProductOrder 创建产品订单
func (h *OrderHandler) CreateProductOrder(c *gin.Context) {
	var order models.ProductOrder
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateProductOrder(&order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// UpdateProductOrder 更新产品订单
func (h *OrderHandler) UpdateProductOrder(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var order models.ProductOrder
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order.ID = uint(id)
	if err := h.service.UpdateProductOrder(&order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// GetRenewalOrders 获取续费订单列表
func (h *OrderHandler) GetRenewalOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")
	userIdStr := c.Query("user_id")
	
	var userId *uint
	if userIdStr != "" {
		uid, err := strconv.Atoi(userIdStr)
		if err == nil {
			temp := uint(uid)
			userId = &temp
		}
	}

	orders, total, err := h.service.GetRenewalOrders(page, limit, status, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"total":  total,
		"page":   page,
		"limit":  limit,
	})
}

// GetRenewalOrder 获取单个续费订单
func (h *OrderHandler) GetRenewalOrder(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	order, err := h.service.GetRenewalOrder(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// CreateRenewalOrder 创建续费订单
func (h *OrderHandler) CreateRenewalOrder(c *gin.Context) {
	var order models.RenewalOrder
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateRenewalOrder(&order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// UpdateRenewalOrder 更新续费订单
func (h *OrderHandler) UpdateRenewalOrder(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var order models.RenewalOrder
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order.ID = uint(id)
	if err := h.service.UpdateRenewalOrder(&order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}