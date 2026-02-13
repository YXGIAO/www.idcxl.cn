package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"time"

	"go-biz-admin/config"
	"go-biz-admin/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService 用户服务
type UserService struct{}

// NewUserService 创建用户服务
func NewUserService() *UserService {
	return &UserService{}
}

const JWTSecret = "your-secret-key-change-in-production"

// 生成随机UUID
func generateUUID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	bytes[6] = (bytes[6] & 0x0f) | 0x40 // 版本4
	bytes[8] = (bytes[8] & 0x3f) | 0x80 // 变体

	uuid := make([]byte, 36)
	hex.Encode(uuid[0:8], bytes[0:4])
	uuid[8] = '-'
	hex.Encode(uuid[9:13], bytes[4:6])
	uuid[13] = '-'
	hex.Encode(uuid[14:18], bytes[6:8])
	uuid[18] = '-'
	hex.Encode(uuid[19:23], bytes[8:10])
	uuid[23] = '-'
	hex.Encode(uuid[24:], bytes[10:])

	return string(uuid), nil
}

// 验证邮箱格式
func validateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// AuthenticateUser 用户认证
func (s *UserService) AuthenticateUser(username, password string) (string, error) {
	var user models.User

	// 先按 username 或 email 查找
	result := config.DB.Where("username = ? OR email = ?", username, username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 如果未找到，再尝试按手机号查找（手机号保存在 user profile）
			var profile models.UserProfile
			profileResult := config.DB.Where("phone = ?", username).First(&profile)
			if profileResult.Error != nil && errors.Is(profileResult.Error, gorm.ErrRecordNotFound) {
				// 未找到任何匹配
				return "", fmt.Errorf("用户不存在: %s", username)
			} else if profileResult.Error != nil {
				// 数据库错误
				return "", fmt.Errorf("查询用户资料失败: %v", profileResult.Error)
			}
			
			// 通过 profile.UserID 加载用户
			if err := config.DB.First(&user, profile.UserID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return "", fmt.Errorf("用户不存在: %d", profile.UserID)
				}
				return "", fmt.Errorf("查询用户失败: %v", err)
			}
		} else {
			// 数据库错误
			return "", fmt.Errorf("查询用户失败: %v", result.Error)
		}
	}

	// 验证密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("密码错误")
	}

	// 生成JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// CreateUser 创建用户
func (s *UserService) CreateUser(user *models.User) error {
	// 验证邮箱格式
	if !validateEmail(user.Email) {
		return errors.New("邮箱格式不正确")
	}

	// 检查邮箱是否已存在
	var existingUser models.User
	result := config.DB.Where("email = ?", user.Email).First(&existingUser)
	if result.Error == nil {
		return errors.New("邮箱已被注册")
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("查询邮箱失败: %v", result.Error)
	}

	// 检查用户名是否已存在
	result = config.DB.Where("username = ?", user.Username).First(&existingUser)
	if result.Error == nil {
		return errors.New("用户名已被注册")
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("查询用户名失败: %v", result.Error)
	}

	// 如果没有提供姓名，则使用用户名作为姓名
	if user.Name == "" {
		user.Name = user.Username
	}

	// 生成UUID
	uuid, err := generateUUID()
	if err != nil {
		return err
	}
	user.UUID = uuid

	// 获取当前最大自定义ID并加1
	var maxCustomIDUser models.User
	config.DB.Table("users").Order("custom_id DESC").Limit(1).Find(&maxCustomIDUser)
	if maxCustomIDUser.CustomID > 0 {
		user.CustomID = maxCustomIDUser.CustomID + 1
	} else {
		// 如果没有用户，则从10000开始
		user.CustomID = 10000
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// 设置默认值
	if user.Role == "" {
		user.Role = "user"
	}
	if user.Status == 0 { // 修改：Status是int类型，不是字符串
		user.Status = 1 // 修改：启用状态为1
	}

	return config.DB.Create(user).Error
}

// GetUsers 获取用户列表
func GetUsers(page, limit int, status, keyword string) ([]models.User, int, error) {
	var users []models.User
	var total int64

	db := config.DB.Model(&models.User{})

	if status != "" {
		statusInt := 1 // 默认启用
		if status == "active" || status == "1" {
			statusInt = 1
		} else if status == "inactive" || status == "0" {
			statusInt = 0
		}
		db = db.Where("status = ?", statusInt)
	}

	if keyword != "" {
		db = db.Where("name LIKE ? OR email LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	db.Count(&total)

	offset := (page - 1) * limit
	if err := db.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, int(total), nil
}

// GetUser 获取单个用户
func GetUser(id uint) (*models.User, error) {
	var user models.User
	result := config.DB.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("用户不存在: %d", id)
		}
		return nil, result.Error
	}
	
	// 尝试加载用户资料，但不作为必需条件
	var profile models.UserProfile
	resultProfile := config.DB.Where("user_id = ?", id).First(&profile)
	if resultProfile.Error != nil {
		if !errors.Is(resultProfile.Error, gorm.ErrRecordNotFound) {
			// 记录非"记录未找到"的错误
			fmt.Printf("查询用户资料时发生错误: %v\n", resultProfile.Error)
		}
		// 如果用户资料不存在，仍返回用户信息但不设置Profile
		return &user, nil
	}
	
	// 只有在成功获取资料时才赋值
	user.Profile = profile
	return &user, nil
}

// UpdateUser 更新用户
func UpdateUser(user *models.User) error {
	return config.DB.Save(user).Error
}

// DeleteUser 删除用户
func DeleteUser(id uint) error {
	return config.DB.Delete(&models.User{}, id).Error
}

// GetUserProfile 获取用户资料
func GetUserProfile(userID uint) (*models.UserProfile, error) {
	var profile models.UserProfile
	result := config.DB.Where("user_id = ?", userID).First(&profile)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("用户资料不存在: %d", userID)
		}
		return nil, result.Error
	}
	return &profile, nil
}

// UpdateUserProfile 更新用户资料
func UpdateUserProfile(profile *models.UserProfile) error {
	var existingProfile models.UserProfile
	result := config.DB.Where("user_id = ?", profile.UserID).First(&existingProfile)
	
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 如果资料不存在，则创建
			return config.DB.Create(profile).Error
		}
		return result.Error
	}

	// 如果资料存在，则更新
	profile.ID = existingProfile.ID
	return config.DB.Save(profile).Error
}

// VerifyIdentity 提交实名认证
func (s *UserService) VerifyIdentity(userID uint, realName, idCard, phone, address string) error {
	// 检查用户是否存在
	var user models.User
	result := config.DB.First(&user, userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return fmt.Errorf("查询用户失败: %v", result.Error)
	}

	// 检查是否已有认证资料
	var profile models.UserProfile
	result = config.DB.Where("user_id = ?", userID).First(&profile)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 如果没有资料，则创建
			now := time.Now()
			profile = models.UserProfile{
				UserID:     userID,
				RealName:   realName,
				IDCard:     idCard,
				Phone:      phone,
				Address:    address,
				AuthStatus: "pending",
				CreatedAt:  now,
				UpdatedAt:  now,
			}
			return config.DB.Create(&profile).Error
		}
		return result.Error
	} else {
		// 如果有资料，则更新
		profile.RealName = realName
		profile.IDCard = idCard
		profile.Phone = phone
		profile.Address = address
		profile.AuthStatus = "pending"
		now := time.Now()
		profile.UpdatedAt = now
		return config.DB.Save(&profile).Error
	}
}

// ApproveIdentity 审核实名认证（内部使用）
func ApproveIdentity(userID uint) error {
	var profile models.UserProfile
	result := config.DB.Where("user_id = ?", userID).First(&profile)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("用户资料不存在")
		}
		return fmt.Errorf("查询用户资料失败: %v", result.Error)
	}

	now := time.Now()
	profile.AuthStatus = "approved"
	profile.AuthAt = &now
	profile.UpdatedAt = now

	// 更新用户表中的实名认证状态
	var user models.User
	result = config.DB.First(&user, userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return fmt.Errorf("查询用户失败: %v", result.Error)
	}
	auth := true
	user.RealNameAuth = &auth

	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Save(&profile).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// RejectIdentity 拒绝实名认证（内部使用）
func RejectIdentity(userID uint) error {
	var profile models.UserProfile
	result := config.DB.Where("user_id = ?", userID).First(&profile)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("用户资料不存在")
		}
		return fmt.Errorf("查询用户资料失败: %v", result.Error)
	}

	profile.AuthStatus = "rejected"
	profile.UpdatedAt = time.Now()

	return config.DB.Save(&profile).Error
}

// Login 用户登录
func (s *UserService) Login(username, password string) (string, error) {
	return s.AuthenticateUser(username, password)
}

// Register 用户注册
func (s *UserService) Register(user *models.User) error {
	return s.CreateUser(user)
}

// GetUsersPaginated 获取分页用户列表
func (s *UserService) GetUsers(page, limit int, status, keyword string) ([]models.User, int, error) {
	return GetUsers(page, limit, status, keyword)
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUser(id uint) (*models.User, error) {
	return GetUser(id)
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(id uint, user *models.User) (*models.User, error) {
	user.ID = id
	err := UpdateUser(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(id uint) error {
	return DeleteUser(id)
}

// GetUserProfile 获取用户资料
func (s *UserService) GetUserProfile(userID uint) (*models.UserProfile, error) {
	return GetUserProfile(userID)
}

// UpdateUserProfile 更新用户资料
func (s *UserService) UpdateUserProfile(profile *models.UserProfile) error {
	return UpdateUserProfile(profile)
}