package handlers

import (
	"net/http"
	"strconv"

	"go-biz-admin/models"
	"go-biz-admin/services"

	"github.com/gin-gonic/gin"
)

// FinanceHandler 财务处理器
type FinanceHandler struct {
	service *services.FinanceService
}

// NewFinanceHandler 创建财务处理器
func NewFinanceHandler(service *services.FinanceService) *FinanceHandler {
	return &FinanceHandler{service: service}
}

// GetTransactions 获取交易流水列表
func (h *FinanceHandler) GetTransactions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	userIdStr := c.Query("user_id")
	orderIdStr := c.Query("order_id")
	transactionType := c.Query("type")
	status := c.Query("status")

	var userId *uint
	if userIdStr != "" {
		uid, err := strconv.Atoi(userIdStr)
		if err == nil {
			temp := uint(uid)
			userId = &temp
		}
	}

	var orderId *uint
	if orderIdStr != "" {
		oid, err := strconv.Atoi(orderIdStr)
		if err == nil {
			temp := uint(oid)
			orderId = &temp
		}
	}

	transactions, total, err := h.service.GetTransactions(page, limit, userId, orderId, transactionType, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
		"total":        total,
		"page":         page,
		"limit":        limit,
	})
}

// GetTransaction 获取单个交易流水
func (h *FinanceHandler) GetTransaction(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	transaction, err := h.service.GetTransaction(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "交易流水不存在"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// GetBills 获取账单列表
func (h *FinanceHandler) GetBills(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	userIdStr := c.Query("user_id")
	status := c.Query("status")

	var userId *uint
	if userIdStr != "" {
		uid, err := strconv.Atoi(userIdStr)
		if err == nil {
			temp := uint(uid)
			userId = &temp
		}
	}

	bills, total, err := h.service.GetBills(page, limit, userId, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"bills": bills,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetBill 获取单个账单
func (h *FinanceHandler) GetBill(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	bill, err := h.service.GetBill(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "账单不存在"})
		return
	}

	c.JSON(http.StatusOK, bill)
}

// PayBill 支付账单
func (h *FinanceHandler) PayBill(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	paymentMethod := c.Query("method")

	if err := h.service.PayBill(uint(id), paymentMethod); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "账单支付成功"})
}

// CreateBill 创建账单
func (h *FinanceHandler) CreateBill(c *gin.Context) {
	var bill models.Bill
	if err := c.ShouldBindJSON(&bill); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateBill(&bill); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, bill)
}

// UpdateBill 更新账单
func (h *FinanceHandler) UpdateBill(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var bill models.Bill
	if err := c.ShouldBindJSON(&bill); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 先获取现有账单
	existingBill, err := h.service.GetBill(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "账单不存在"})
		return
	}

	// 更新字段
	if bill.UserID != 0 {
		existingBill.UserID = bill.UserID
	}
	if bill.Amount != 0 {
		existingBill.Amount = bill.Amount
	}
	existingBill.DueDate = bill.DueDate
	if bill.Status != "" {
		existingBill.Status = bill.Status
	}
	if bill.Description != "" {
		existingBill.Description = bill.Description
	}

	if err := h.service.UpdateBill(existingBill); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, existingBill)
}
