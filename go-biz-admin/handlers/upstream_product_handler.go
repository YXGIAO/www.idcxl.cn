package handlers

import (
	"net/http"
	"strconv"

	"go-biz-admin/models"
	"go-biz-admin/services"

	"github.com/gin-gonic/gin"
)

// UpstreamProductHandler 上游产品处理器
type UpstreamProductHandler struct {
	service *services.UpstreamProductService
}

// NewUpstreamProductHandler 创建上游产品处理器
func NewUpstreamProductHandler(service *services.UpstreamProductService) *UpstreamProductHandler {
	return &UpstreamProductHandler{service: service}
}

// GetProductList 获取产品列表
// @Summary 获取上游产品列表
// @Description 获取上游产品列表，支持筛选、分页、排序
// @Tags 上游产品管理
// @Accept json
// @Produce json
// @Param keywords query string false "关键字搜索（ID、用户名、邮箱、手机号、商品名称、产品标识）"
// @Param supplier_id query int false "供应商ID"
// @Param billing_cycle query string false "付款周期"
// @Param status query string false "状态：Unpaid未付款、Pending开通中、Active已开通、Suspended已暂停、Deleted已删除、Failed开通失败"
// @Param start_time query int false "开始时间戳（Unix时间戳，单位秒）"
// @Param end_time query int false "结束时间戳（Unix时间戳，单位秒）"
// @Param page query int false "页码，默认1" default(1)
// @Param limit query int false "每页条数，默认10" default(10)
// @Param orderby query string false "排序字段，默认id" default(id)
// @Param sort query string false "排序方式：asc升序、desc降序，默认asc" default(asc)
// @Success 200 {object} models.ProductListResponse
// @Failure 400 {object} map[string]string "参数错误"
// @Failure 500 {object} map[string]string "服务器错误"
// @Router /admin/v1/upstream/host [get]
func (h *UpstreamProductHandler) GetProductList(c *gin.Context) {
	var params models.ProductListRequest

	// 绑定查询参数
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	// 设置默认值
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Limit > 100 {
		params.Limit = 100 // 限制最大每页条数
	}
	if params.OrderBy == "" {
		params.OrderBy = "id"
	}
	if params.Sort == "" {
		params.Sort = "asc"
	}

	// 调用服务获取产品列表
	response, err := h.service.GetProductList(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取产品列表失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    response,
	})
}

// GetProductByID 获取产品详情
// @Summary 获取上游产品详情
// @Description 根据ID获取上游产品详情
// @Tags 上游产品管理
// @Accept json
// @Produce json
// @Param id path int true "产品ID"
// @Success 200 {object} models.ProductDetailResponse
// @Failure 400 {object} map[string]string "参数错误"
// @Failure 404 {object} map[string]string "产品不存在"
// @Failure 500 {object} map[string]string "服务器错误"
// @Router /admin/v1/upstream/host/{id} [get]
func (h *UpstreamProductHandler) GetProductByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的产品ID"})
		return
	}

	product, err := h.service.GetProductByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "产品不存在: " + err.Error()})
		return
	}

	// 按照API文档要求返回 { "host": { ... } } 格式
	response := models.ProductDetailResponse{
		Host: *product,
	}

	c.JSON(http.StatusOK, response)
}

// CreateProduct 创建产品
// @Summary 创建上游产品
// @Description 创建新的上游产品
// @Tags 上游产品管理
// @Accept json
// @Produce json
// @Param product body models.UpstreamProduct true "产品信息"
// @Success 201 {object} models.UpstreamProduct
// @Failure 400 {object} map[string]string "参数错误"
// @Failure 500 {object} map[string]string "服务器错误"
// @Router /admin/v1/upstream/host [post]
func (h *UpstreamProductHandler) CreateProduct(c *gin.Context) {
	var product models.UpstreamProduct
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	if err := h.service.CreateProduct(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建产品失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    201,
		"message": "产品创建成功",
		"data":    product,
	})
}

// UpdateProduct 更新产品
// @Summary 更新上游产品
// @Description 更新已存在的上游产品
// @Tags 上游产品管理
// @Accept json
// @Produce json
// @Param id path int true "产品ID"
// @Param product body models.UpstreamProduct true "产品信息"
// @Success 200 {object} models.UpstreamProduct
// @Failure 400 {object} map[string]string "参数错误"
// @Failure 404 {object} map[string]string "产品不存在"
// @Failure 500 {object} map[string]string "服务器错误"
// @Router /admin/v1/upstream/host/{id} [put]
func (h *UpstreamProductHandler) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的产品ID"})
		return
	}

	var product models.UpstreamProduct
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	if err := h.service.UpdateProduct(uint(id), &product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新产品失败: " + err.Error()})
		return
	}

	// 重新获取更新后的产品
	updatedProduct, _ := h.service.GetProductByID(uint(id))

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "产品更新成功",
		"data":    updatedProduct,
	})
}

// DeleteProduct 删除产品
// @Summary 删除上游产品
// @Description 删除指定的上游产品
// @Tags 上游产品管理
// @Accept json
// @Produce json
// @Param id path int true "产品ID"
// @Success 200 {object} map[string]string "删除成功"
// @Failure 400 {object} map[string]string "参数错误"
// @Failure 404 {object} map[string]string "产品不存在"
// @Failure 500 {object} map[string]string "服务器错误"
// @Router /admin/v1/upstream/host/{id} [delete]
func (h *UpstreamProductHandler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的产品ID"})
		return
	}

	if err := h.service.DeleteProduct(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除产品失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "产品删除成功",
	})
}

// GetStatusOptions 获取状态选项
// @Summary 获取状态选项
// @Description 获取产品状态的可选值
// @Tags 上游产品管理
// @Accept json
// @Produce json
// @Success 200 {array} map[string]string
// @Router /admin/v1/upstream/host/status-options [get]
func (h *UpstreamProductHandler) GetStatusOptions(c *gin.Context) {
	options := h.service.GetStatusOptions()
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    options,
	})
}

// GetBillingCycleOptions 获取付款周期选项
// @Summary 获取付款周期选项
// @Description 获取付款周期的可选值
// @Tags 上游产品管理
// @Accept json
// @Produce json
// @Success 200 {array} map[string]string
// @Router /admin/v1/upstream/host/billing-cycle-options [get]
func (h *UpstreamProductHandler) GetBillingCycleOptions(c *gin.Context) {
	options := h.service.GetBillingCycleOptions()
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    options,
	})
}