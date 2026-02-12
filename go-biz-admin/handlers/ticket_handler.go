package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-biz-admin/models"
	"go-biz-admin/services"
)

// TicketHandler 工单处理器
type TicketHandler struct {
	service *services.TicketService
}

// NewTicketHandler 创建工单处理器
func NewTicketHandler(service *services.TicketService) *TicketHandler {
	return &TicketHandler{service: service}
}

// GetTickets 获取工单列表
func (h *TicketHandler) GetTickets(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")
	category := c.Query("category")
	userIdStr := c.Query("user_id")
	
	var userId *uint
	if userIdStr != "" {
		uid, err := strconv.Atoi(userIdStr)
		if err == nil {
			temp := uint(uid)
			userId = &temp
		}
	}

	tickets, total, err := h.service.GetTickets(page, limit, status, category, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tickets": tickets,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

// GetTicket 获取单个工单
func (h *TicketHandler) GetTicket(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	ticket, err := h.service.GetTicket(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "工单不存在"})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

// CreateTicket 创建工单
func (h *TicketHandler) CreateTicket(c *gin.Context) {
	var ticket models.Ticket
	if err := c.ShouldBindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateTicket(&ticket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ticket)
}

// UpdateTicket 更新工单
func (h *TicketHandler) UpdateTicket(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var ticket models.Ticket
	if err := c.ShouldBindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticket.ID = uint(id)
	if err := h.service.UpdateTicket(&ticket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

// GetTicketStats 获取工单统计
func (h *TicketHandler) GetTicketStats(c *gin.Context) {
	stats, err := h.service.GetTicketStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}