package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"go-biz-admin/services"

	"github.com/gin-gonic/gin"
)

// ZJMFHandler 智简魔方处理器
type ZJMFHandler struct {
	service *services.ZJMFService
}

// NewZJMFHandler 创建智简魔方处理器
func NewZJMFHandler(service *services.ZJMFService) *ZJMFHandler {
	return &ZJMFHandler{service: service}
}

// GetUserDetail 获取智简魔方用户详情
func (h *ZJMFHandler) GetUserDetail(c *gin.Context) {
	userID := c.Param("id")

	// 将ID转换为字符串传递给服务
	userData, err := h.service.GetUserDetail(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userData)
}

// CreateServer 在智简魔方中创建服务器
func (h *ZJMFHandler) CreateServer(c *gin.Context) {
	var serverData map[string]interface{}
	if err := c.ShouldBindJSON(&serverData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.CreateServer(serverData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// SyncServers 同步智简魔方服务器
func (h *ZJMFHandler) SyncServers(c *gin.Context) {
	err := h.service.SyncServersToSystem()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("同步失败: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "服务器同步成功"})
}

// GetInfo 获取智简魔方信息
func (h *ZJMFHandler) GetInfo(c *gin.Context) {
	var supplierData map[string]interface{}
	if err := c.ShouldBindJSON(&supplierData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从请求数据中提取必要的参数
	apiEndpoint, _ := supplierData["api_endpoint"].(string)
	apiKey, _ := supplierData["api_key"].(string)
	apiSecret, _ := supplierData["api_secret"].(string)
	if apiEndpoint == "" {
		apiEndpoint, _ = supplierData["interface_url"].(string)
	}
	if apiKey == "" {
		apiKey, _ = supplierData["zjmf_api_key"].(string)
	}
	if apiSecret == "" {
		apiSecret, _ = supplierData["zjmf_api_secret"].(string)
	}

	// 如果缺少必要参数，返回错误
	if apiEndpoint == "" || apiKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "API端点或API密钥不能为空"})
		return
	}

	// 创建一个新的ZJMFService实例用于此请求
	zjmfService := services.NewZJMFService(apiEndpoint, apiKey, apiSecret)

	// 获取供应商信息
	info, err := zjmfService.GetSupplierInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("获取供应商信息失败: %v", err)})
		return
	}

	c.JSON(http.StatusOK, info)
}

// GetSuppliers 获取供应商列表
func (h *ZJMFHandler) GetSuppliers(c *gin.Context) {
	// 获取查询参数
	keywords := c.Query("keywords")
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "20")
	orderby := c.DefaultQuery("orderby", "id")
	sort := c.DefaultQuery("sort", "asc")

	// 转换为适合服务的参数格式
	params := make(map[string]interface{})
	if keywords != "" {
		params["keywords"] = keywords
	}
	if pageNum, err := strconv.Atoi(page); err == nil {
		params["page"] = pageNum
	}
	if limitNum, err := strconv.Atoi(limit); err == nil {
		params["limit"] = limitNum
	}
	if orderby != "" {
		params["orderby"] = orderby
	}
	if sort != "" {
		params["sort"] = sort
	}

	// 调用服务获取供应商列表
	suppliers, err := h.service.GetSuppliers(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("获取供应商列表失败: %v", err)})
		return
	}

	c.JSON(http.StatusOK, suppliers)
}

// GetProducts 获取智简魔方产品信息
func (h *ZJMFHandler) GetProducts(c *gin.Context) {
	var supplierData map[string]interface{}
	if err := c.ShouldBindJSON(&supplierData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从请求数据中提取必要的参数
	apiEndpoint, _ := supplierData["api_endpoint"].(string)
	apiKey, _ := supplierData["api_key"].(string)
	apiSecret, _ := supplierData["api_secret"].(string)
	if apiEndpoint == "" {
		apiEndpoint, _ = supplierData["interface_url"].(string)
	}
	if apiKey == "" {
		apiKey, _ = supplierData["zjmf_api_key"].(string)
	}
	if apiSecret == "" {
		apiSecret, _ = supplierData["zjmf_api_secret"].(string)
	}

	// 如果缺少必要参数，返回错误
	if apiEndpoint == "" || apiKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "API端点或API密钥不能为空"})
		return
	}

	// 创建一个新的ZJMFService实例用于此请求
	zjmfService := services.NewZJMFService(apiEndpoint, apiKey, apiSecret)

	// 获取产品信息
	products, err := zjmfService.GetProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("获取产品信息失败: %v", err)})
		return
	}

	c.JSON(http.StatusOK, products)
}
