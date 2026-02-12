package services

import (
	"time"

	"go-biz-admin/config"
	"go-biz-admin/models"
)

// TicketService 工单服务
type TicketService struct{}

// NewTicketService 创建工单服务
func NewTicketService() *TicketService {
	return &TicketService{}
}

// GetTickets 获取工单列表
func (s *TicketService) GetTickets(page, limit int, status, category string, userID *uint) ([]models.Ticket, int, error) {
	var tickets []models.Ticket
	var total int64

	db := config.DB.Model(&models.Ticket{})

	if status != "" {
		db = db.Where("status = ?", status)
	}

	if category != "" {
		db = db.Where("category = ?", category)
	}

	if userID != nil {
		db = db.Where("user_id = ?", *userID)
	}

	db.Count(&total)

	offset := (page - 1) * limit
	if err := db.Offset(offset).Limit(limit).Preload("User").Find(&tickets).Error; err != nil {
		return nil, 0, err
	}

	return tickets, int(total), nil
}

// GetTicket 获取单个工单
func (s *TicketService) GetTicket(id uint) (*models.Ticket, error) {
	var ticket models.Ticket
	err := config.DB.Preload("User").First(&ticket, id).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

// CreateTicket 创建工单
func (s *TicketService) CreateTicket(ticket *models.Ticket) error {
	if ticket.Status == "" {
		ticket.Status = "open"
	}
	if ticket.Priority == "" {
		ticket.Priority = "medium"
	}

	ticket.CreatedAt = time.Now()
	ticket.UpdatedAt = time.Now()

	return config.DB.Create(ticket).Error
}

// UpdateTicket 更新工单
func (s *TicketService) UpdateTicket(ticket *models.Ticket) error {
	ticket.UpdatedAt = time.Now()
	return config.DB.Save(ticket).Error
}

// GetTicketStats 获取工单统计
func (s *TicketService) GetTicketStats() (map[string]int, error) {
	var total, open, closed, inProgress int64

	config.DB.Model(&models.Ticket{}).Count(&total)
	config.DB.Model(&models.Ticket{}).Where("status = ?", "open").Count(&open)
	config.DB.Model(&models.Ticket{}).Where("status = ?", "closed").Count(&closed)
	config.DB.Model(&models.Ticket{}).Where("status = ?", "in_progress").Count(&inProgress)

	return map[string]int{
		"total":       int(total),
		"open":        int(open),
		"closed":      int(closed),
		"in_progress": int(inProgress),
	}, nil
}

// UpdateTicketStatus 更新工单状态
func (s *TicketService) UpdateTicketStatus(ticketID uint, status string) error {
	var ticket models.Ticket
	if err := config.DB.First(&ticket, ticketID).Error; err != nil {
		return err
	}

	ticket.Status = status
	ticket.UpdatedAt = time.Now()

	return config.DB.Save(&ticket).Error
}