package handlers

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"go-biz-admin/config"
	"go-biz-admin/models"
	"go-biz-admin/services"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// UserHandler 处理用户相关HTTP请求
type UserHandler struct {
	service *services.UserService
}

// NewUserHandler 创建新的UserHandler实例
func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	var credentials struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.Login(credentials.Username, credentials.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Register 用户注册
func (h *UserHandler) Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Register(&user); err != nil {
		log.Println("Register error:", err)
		if err.Error() == "邮箱已被注册" || err.Error() == "用户名已被注册" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// CreateAdminAccount 创建管理员账户
func (h *UserHandler) CreateAdminAccount(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置为管理员角色
	user.Role = "admin"

	// 验证必要字段
	if user.Username == "" || user.Email == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名、邮箱和密码不能为空"})
		return
	}

	// 检查用户名或邮箱是否已存在
	var existingUser models.User
	if err := config.DB.Where("username = ? OR email = ?", user.Username, user.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "用户名或邮箱已存在"})
		return
	}

	// 如果没有提供姓名，则使用用户名作为姓名
	if user.Name == "" {
		user.Name = user.Username
	}

	// 生成UUID
	uuidBytes := make([]byte, 16)
	_, err := rand.Read(uuidBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成UUID失败"})
		return
	}
	// 设置版本和变体
	uuidBytes[6] = (uuidBytes[6] & 0x0f) | 0x40 // 版本4
	uuidBytes[8] = (uuidBytes[8] & 0x3f) | 0x80 // 变体

	// 转换为字符串格式
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		uuidBytes[0:4],
		uuidBytes[4:6],
		uuidBytes[6:8],
		uuidBytes[8:10],
		uuidBytes[10:])
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}
	user.Password = string(hashedPassword)

	// 设置默认值
	if user.Status == 0 {
		user.Status = 1
	}

	// 设置默认头像
	if user.Avatar == "" {
		user.Avatar = "https://via.placeholder.com/100x100/409eff/white?text=头像"
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "管理员账户创建成功",
		"user_id":   user.ID,
		"custom_id": user.CustomID,
	})
}

// GetUsers 获取用户列表
func (h *UserHandler) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")
	keyword := c.Query("keyword")

	users, total, err := h.service.GetUsers(page, limit, status, keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"users": users,
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

// GetUser 获取单个用户
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	user, err := h.service.GetUser(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    user,
	})
}

// UpdateUser 更新用户
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUser, err := h.service.UpdateUser(uint(id), &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

// DeleteUser 删除用户
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.service.DeleteUser(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户删除成功"})
}

// GetUserProfile 获取用户资料
func (h *UserHandler) GetUserProfile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	profile, err := h.service.GetUserProfile(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户资料不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    profile,
	})
}

// UpdateUserProfile 更新用户资料
func (h *UserHandler) UpdateUserProfile(c *gin.Context) {
	var profile models.UserProfile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateUserProfile(&profile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// UpdateUserAvatar 更新用户头像
func (h *UserHandler) UpdateUserAvatar(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req struct {
		Avatar string `json:"avatar" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取用户
	user, err := h.service.GetUser(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 更新头像
	user.Avatar = req.Avatar
	_, err = h.service.UpdateUser(uint(id), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "头像更新成功", "avatar": req.Avatar})
}

// VerifyIdentity 提交实名认证
func (h *UserHandler) VerifyIdentity(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req struct {
		RealName string `json:"real_name" binding:"required"`
		IDCard   string `json:"id_card" binding:"required"`
		Phone    string `json:"phone"`
		Address  string `json:"address"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.VerifyIdentity(uint(id), req.RealName, req.IDCard, req.Phone, req.Address); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "实名认证提交成功"})
}

// GetCurrentUser 获取当前登录用户信息
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	// 从上下文中获取用户ID（通常在JWT中间件中设置）
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	// 根据用户ID获取用户信息
	user, err := h.service.GetUser(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    user,
	})
}
