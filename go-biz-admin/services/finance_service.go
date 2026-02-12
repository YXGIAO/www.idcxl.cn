package services

import (
	"errors"
	"fmt"
	"time"

	"go-biz-admin/config"
	"go-biz-admin/models"
)

// FinanceService 财务服务
type FinanceService struct{}

// NewFinanceService 创建财务服务
func NewFinanceService() *FinanceService {
	return &FinanceService{}
}

// GetTransactions 获取交易流水列表
func (s *FinanceService) GetTransactions(page, limit int, userID, orderID *uint, transactionType, status string) ([]models.Transaction, int, error) {
	var transactions []models.Transaction
	var total int64

	db := config.DB.Model(&models.Transaction{})
	
	if userID != nil {
		db = db.Where("user_id = ?", *userID)
	}
	
	if orderID != nil {
		db = db.Where("order_id = ?", *orderID)
	}
	
	if transactionType != "" {
		db = db.Where("type = ?", transactionType)
	}
	
	if status != "" {
		db = db.Where("status = ?", status)
	}
	
	db.Count(&total)
	
	offset := (page - 1) * limit
	if err := db.Offset(offset).Limit(limit).Preload("User").Find(&transactions).Error; err != nil {
		return nil, 0, err
	}

	return transactions, int(total), nil
}

// GetTransaction 获取单个交易流水
func (s *FinanceService) GetTransaction(id uint) (*models.Transaction, error) {
	var transaction models.Transaction
	err := config.DB.Preload("User").First(&transaction, id).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

// GetBills 获取账单列表
func (s *FinanceService) GetBills(page, limit int, userID *uint, status string) ([]models.Bill, int, error) {
	var bills []models.Bill
	var total int64

	db := config.DB.Model(&models.Bill{})
	
	if userID != nil {
		db = db.Where("user_id = ?", *userID)
	}
	
	if status != "" {
		db = db.Where("status = ?", status)
	}
	
	db.Count(&total)
	
	offset := (page - 1) * limit
	if err := db.Offset(offset).Limit(limit).Preload("User").Find(&bills).Error; err != nil {
		return nil, 0, err
	}

	return bills, int(total), nil
}

// GetBill 获取单个账单
func (s *FinanceService) GetBill(id uint) (*models.Bill, error) {
	var bill models.Bill
	err := config.DB.Preload("User").First(&bill, id).Error
	if err != nil {
		return nil, err
	}
	return &bill, nil
}

// PayBill 支付账单
func (s *FinanceService) PayBill(id uint, paymentMethod string) error {
	var bill models.Bill
	if err := config.DB.First(&bill, id).Error; err != nil {
		return errors.New("账单不存在")
	}

	if bill.Status == "paid" {
		return errors.New("账单已支付")
	}

	// 更新账单状态
	bill.Status = "paid"
	bill.UpdatedAt = time.Now()
	
	if err := config.DB.Save(&bill).Error; err != nil {
		return err
	}

	// 创建交易记录
	transaction := models.Transaction{
		UserID:    bill.UserID,
		OrderID:   bill.OrderID,
		OrderType: "bill_payment",
		Amount:    bill.Amount,
		Type:      "expense",
		Status:    "completed",
		TransactionNumber: generateTransactionNumber(),
		Description: fmt.Sprintf("支付账单 #%d", bill.ID),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	return config.DB.Create(&transaction).Error
}

// CreateTransaction 创建交易记录
func (s *FinanceService) CreateTransaction(transaction *models.Transaction) error {
	transaction.TransactionNumber = generateTransactionNumber()
	if transaction.Status == "" {
		transaction.Status = "completed"
	}
	
	return config.DB.Create(transaction).Error
}

// CreateBill 创建账单
func (s *FinanceService) CreateBill(bill *models.Bill) error {
	if bill.Status == "" {
		bill.Status = "unpaid"
	}
	bill.IssueDate = time.Now()
	
	return config.DB.Create(bill).Error
}

// UpdateBill 更新账单
func (s *FinanceService) UpdateBill(bill *models.Bill) error {
	bill.UpdatedAt = time.Now()
	return config.DB.Save(bill).Error
}

// 生成交易号
func generateTransactionNumber() string {
	now := time.Now()
	return fmt.Sprintf("TXN%d%02d%02d%02d%02d%02d%06d",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second(),
		time.Now().Nanosecond()/1000)
}