package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go-biz-admin/config"
	"go-biz-admin/models"
)

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ZJMFService 智简魔方服务
type ZJMFService struct {
	BaseURL   string
	APIKey    string
	APISecret string
}

// NewZJMFService 创建智简魔方服务实例
func NewZJMFService(baseURL, apiKey, apiSecret string) *ZJMFService {
	return &ZJMFService{
		BaseURL:   baseURL,
		APIKey:    apiKey,
		APISecret: apiSecret,
	}
}

// makeRequest 发起HTTP请求到智简魔方API
func (s *ZJMFService) makeRequest(method, endpoint string, payload interface{}) ([]byte, error) {
	// 清理endpoint，确保没有多余的斜杠
	cleanEndpoint := endpoint
	if len(cleanEndpoint) > 0 && cleanEndpoint[0] == '/' {
		cleanEndpoint = cleanEndpoint[1:]
	}

	url := s.BaseURL
	if url[len(url)-1] != '/' {
		url += "/"
	}
	url += cleanEndpoint

	fmt.Printf("发起请求: %s %s\n", method, url)

	var req *http.Request
	var err error

	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return nil, err
		}
	}

	// 设置API认证头
	req.Header.Set("Authorization", "Bearer "+s.APIKey)
	req.Header.Set("X-API-Key", s.APIKey)
	req.Header.Set("X-API-Secret", s.APISecret)

	// 添加用户代理
	req.Header.Set("User-Agent", "Go-Biz-Admin/1.0")

	// 添加接受头
	req.Header.Set("Accept", "application/json")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Printf("响应状态码: %d\n", resp.StatusCode)
	fmt.Printf("响应头: Content-Type=%s\n", resp.Header.Get("Content-Type"))

	// 检查重定向
	if resp.StatusCode == 301 || resp.StatusCode == 302 || resp.StatusCode == 307 || resp.StatusCode == 308 {
		location := resp.Header.Get("Location")
		fmt.Printf("检测到重定向到: %s\n", location)
		if location != "" {
			// 递归调用处理重定向
			return s.makeRequest(method, location, payload)
		}
	}

	// 检查是否是HTML错误页面
	contentType := resp.Header.Get("Content-Type")
	if contentType != "" && (contentType == "text/html" || contentType == "text/html; charset=UTF-8") {
		content := string(body)
		if len(content) > 500 {
			content = content[:500]
		}
		fmt.Printf("收到HTML响应，内容: %s\n", content)

		// 检查是否是JavaScript重定向
		if strings.Contains(content, "window.location.href") {
			fmt.Printf("检测到JavaScript重定向，尝试处理...\n")

			// 提取 yxd_token
			var yxdToken string
			tokenRe := regexp.MustCompile(`yxd_token=([a-f0-9]+)`)
			tokenMatches := tokenRe.FindStringSubmatch(content)
			if len(tokenMatches) > 1 {
				yxdToken = tokenMatches[1]
				fmt.Printf("提取到 yxd_token: %s\n", yxdToken)
			}

			// 提取重定向路径
			hrefRe := regexp.MustCompile(`window\.location\.href='([^']+)'`)
			hrefMatches := hrefRe.FindStringSubmatch(content)
			if len(hrefMatches) > 1 {
				redirectPath := hrefMatches[1]
				fmt.Printf("提取到重定向路径: %s\n", redirectPath)

				// 构造带 cookie 的新请求
				if yxdToken != "" {
					cookie := fmt.Sprintf("yxd_token=%s", yxdToken)
					return s.makeRequestWithCookie(method, redirectPath, payload, cookie)
				}
			}

			return nil, fmt.Errorf("无法提取JavaScript重定向信息")
		}

		return nil, fmt.Errorf("received HTML response instead of JSON, status: %d", resp.StatusCode)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// GetUserDetail 获取智简魔方用户详情
func (s *ZJMFService) GetUserDetail(clientID string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("/clients/details?client_id=%s", clientID)
	response, err := s.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var userData map[string]interface{}
	err = json.Unmarshal(response, &userData)
	if err != nil {
		return nil, err
	}

	return userData, nil
}

// CreateServer 在智简魔方中创建服务器
func (s *ZJMFService) CreateServer(serverData map[string]interface{}) (map[string]interface{}, error) {
	response, err := s.makeRequest("POST", "/servers", serverData)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// SyncServers 同步智简魔方服务器列表到本地系统
func (s *ZJMFService) SyncServers() error {
	// 从智简魔方API获取服务器列表
	response, err := s.makeRequest("GET", "/servers", nil)
	if err != nil {
		return err
	}

	var serversData map[string]interface{}
	err = json.Unmarshal(response, &serversData)
	if err != nil {
		return err
	}

	// 解析服务器列表并同步到本地数据库
	serversList, ok := serversData["servers"].([]interface{})
	if !ok {
		return fmt.Errorf("failed to parse servers list from response")
	}

	for _, srv := range serversList {
		serverMap, ok := srv.(map[string]interface{})
		if !ok {
			continue
		}

		// 将服务器数据转换为本地模型
		localServer := models.Server{
			Name:            fmt.Sprintf("%v", serverMap["name"]),
			Host:            fmt.Sprintf("%v", serverMap["host"]),
			Port:            int(serverMap["port"].(float64)),
			Type:            fmt.Sprintf("%v", serverMap["type"]),
			Username:        fmt.Sprintf("%v", serverMap["username"]),
			Password:        "",
			Location:        fmt.Sprintf("%v", serverMap["location"]),
			Status:          fmt.Sprintf("%v", serverMap["status"]),
			CPU:             fmt.Sprintf("%v", serverMap["cpu"]),
			Memory:          fmt.Sprintf("%v", serverMap["memory"]),
			Disk:            fmt.Sprintf("%v", serverMap["disk"]),
			Bandwidth:       fmt.Sprintf("%v", serverMap["bandwidth"]),
			IPCount:         int(serverMap["ip_count"].(float64)),
			Notes:           fmt.Sprintf("%v", serverMap["notes"]),
			Description:     "",
			SupplierID:      0,
			ZJMFServerID:    fmt.Sprintf("%v", serverMap["zjmf_server_id"]),
			ZJMFServerGroup: fmt.Sprintf("%v", serverMap["zjmf_server_group"]),
			ZJMFStatus:      fmt.Sprintf("%v", serverMap["zjmf_status"]),
			ZJMFHost:        fmt.Sprintf("%v", serverMap["zjmf_host"]),
		}

		// 检查服务器是否已存在于本地数据库
		var existingServer models.Server
		result := config.DB.Where("zjmf_server_id = ?", localServer.ZJMFServerID).First(&existingServer)

		if result.Error != nil {
			// 服务器不存在，创建新记录
			if err := config.DB.Create(&localServer).Error; err != nil {
				fmt.Printf("Failed to create server %s: %v\n", localServer.Name, err)
			}
		} else {
			// 服务器存在，更新记录
			existingServer.Name = localServer.Name
			existingServer.Host = localServer.Host
			existingServer.Port = localServer.Port
			existingServer.Type = localServer.Type
			existingServer.Username = localServer.Username
			existingServer.Location = localServer.Location
			existingServer.Status = localServer.Status
			existingServer.CPU = localServer.CPU
			existingServer.Memory = localServer.Memory
			existingServer.Disk = localServer.Disk
			existingServer.Bandwidth = localServer.Bandwidth
			existingServer.IPCount = localServer.IPCount
			existingServer.Notes = localServer.Notes
			existingServer.Description = localServer.Description
			existingServer.SupplierID = localServer.SupplierID
			existingServer.ZJMFStatus = localServer.ZJMFStatus
			existingServer.ZJMFServerGroup = localServer.ZJMFServerGroup
			existingServer.ZJMFHost = localServer.ZJMFHost

			if err := config.DB.Save(&existingServer).Error; err != nil {
				fmt.Printf("Failed to update server %s: %v\n", existingServer.Name, err)
			}
		}
	}

	return nil
}

// SyncServersToSystem 同步智简魔方服务器列表
func (s *ZJMFService) SyncServersToSystem() error {
	return s.SyncServers()
}

// GetBalance 获取智简魔方账户余额
func (s *ZJMFService) GetBalance() (map[string]interface{}, error) {
	endpoint := "/finance/balance" // 假设的余额查询端点
	response, err := s.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var balanceData map[string]interface{}
	err = json.Unmarshal(response, &balanceData)
	if err != nil {
		return nil, err
	}

	return balanceData, nil
}

// GetProducts 获取智简魔方产品列表
func (s *ZJMFService) GetProducts() (map[string]interface{}, error) {
	// 首先尝试认证并获取cookie
	fmt.Println("开始认证流程...")
	cookie, err := s.authenticateAndGetCookie()
	if err != nil {
		fmt.Printf("认证失败: %v\n", err)
		return nil, err
	}
	fmt.Printf("获取到认证信息: %s\n", cookie)

	// 尝试多个可能的产品端点
	endpoints := []string{
		"/api/v1/products",
		"/api/products",
		"/products",
		"/admin/products",
		"/api/product/list",
		"/product/list",
		"/host/products",
		"/api/host/products",
		"/product",
		"/api/product",
	}

	var lastError error

	for _, endpoint := range endpoints {
		fmt.Printf("尝试端点: %s%s\n", s.BaseURL, endpoint)

		// 尝试带cookie的请求
		response, err := s.makeRequestWithCookie("GET", endpoint, nil, cookie)
		if err != nil {
			fmt.Printf("端点 %s 请求失败: %v\n", endpoint, err)
			lastError = err
			continue
		}

		// 检查是否是JavaScript重定向
		if strings.Contains(string(response), "window.location.href") {
			fmt.Printf("端点 %s 返回JavaScript重定向，尝试处理...\n", endpoint)

			// 尝试从重定向中提取实际的产品数据
			processedResponse, err := s.handleJSRedirectResponse(response, endpoint, cookie)
			if err != nil {
				fmt.Printf("处理JavaScript重定向失败: %v\n", err)
				lastError = err
				continue
			}

			response = processedResponse
		}

		// 检查响应是否为有效的JSON
		if len(response) > 0 && response[0] == '{' {
			var productsData map[string]interface{}
			err = json.Unmarshal(response, &productsData)
			if err != nil {
				fmt.Printf("端点 %s JSON解析失败: %v\n", endpoint, err)
				lastError = err
				continue
			}

			fmt.Printf("成功从端点 %s 获取产品数据\n", endpoint)
			return productsData, nil
		} else {
			fmt.Printf("端点 %s 返回非JSON内容: %s\n", endpoint, string(response[:min(len(response), 100)]))
			lastError = fmt.Errorf("invalid response format from endpoint %s", endpoint)
		}
	}

	return nil, fmt.Errorf("所有产品端点都失败了: %v", lastError)
}

// handleJSRedirectResponse 处理JavaScript重定向响应
func (s *ZJMFService) handleJSRedirectResponse(response []byte, originalEndpoint, cookie string) ([]byte, error) {
	content := string(response)

	// 提取重定向路径
	re := regexp.MustCompile(`window\.location\.href='([^']+)'`)
	matches := re.FindStringSubmatch(content)
	if len(matches) > 1 {
		redirectPath := matches[1]
		fmt.Printf("检测到重定向路径: %s\n", redirectPath)

		// 跟随重定向
		finalResponse, err := s.makeRequestWithCookie("GET", redirectPath, nil, cookie)
		if err != nil {
			return nil, fmt.Errorf("跟随重定向失败: %v", err)
		}

		fmt.Printf("重定向后响应长度: %d 字节\n", len(finalResponse))
		return finalResponse, nil
	}

	return nil, fmt.Errorf("无法解析JavaScript重定向路径")
}

// GetSuppliers 获取供应商列表
func (s *ZJMFService) GetSuppliers(params map[string]interface{}) (map[string]interface{}, error) {
	endpoint := "/admin/v1/supplier"

	// 构建查询参数
	queryParams := url.Values{}
	if keywords, ok := params["keywords"].(string); ok && keywords != "" {
		queryParams.Add("keywords", keywords)
	}
	if page, ok := params["page"].(int); ok {
		queryParams.Add("page", strconv.Itoa(page))
	}
	if limit, ok := params["limit"].(int); ok {
		queryParams.Add("limit", strconv.Itoa(limit))
	}
	if orderby, ok := params["orderby"].(string); ok && orderby != "" {
		queryParams.Add("orderby", orderby)
	}
	if sort, ok := params["sort"].(string); ok && sort != "" {
		queryParams.Add("sort", sort)
	}

	// 添加查询参数到URL
	if len(queryParams) > 0 {
		endpoint += "?" + queryParams.Encode()
	}

	response, err := s.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var suppliersData map[string]interface{}
	err = json.Unmarshal(response, &suppliersData)
	if err != nil {
		return nil, err
	}

	return suppliersData, nil
}

// GetSupplierInfo 获取供应商信息（包括余额和产品）
func (s *ZJMFService) GetSupplierInfo() (map[string]interface{}, error) {
	// 并行获取余额和产品信息
	var balanceData map[string]interface{}
	var productsData map[string]interface{}
	var balanceErr, productsErr error

	// 使用goroutine同时获取余额和产品信息
	ch := make(chan bool, 2)
	go func() {
		defer func() { ch <- true }()
		balanceData, balanceErr = s.GetBalance()
	}()
	go func() {
		defer func() { ch <- true }()
		productsData, productsErr = s.GetProducts()
	}()

	// 等待两个goroutine完成
	<-ch
	<-ch

	// 检查是否有错误
	if balanceErr != nil && productsErr != nil {
		return nil, fmt.Errorf("获取余额和产品信息均失败: 余额错误: %v, 产品错误: %v", balanceErr, productsErr)
	} else if balanceErr != nil {
		return nil, fmt.Errorf("获取余额失败: %v", balanceErr)
	} else if productsErr != nil {
		return nil, fmt.Errorf("获取产品信息失败: %v", productsErr)
	}

	// 合并结果
	result := make(map[string]interface{})
	for k, v := range balanceData {
		result[k] = v
	}
	for k, v := range productsData {
		result[k] = v
	}

	return result, nil
}

// authenticateAndGetCookie 执行认证流程并获取cookie
func (s *ZJMFService) authenticateAndGetCookie() (string, error) {
	fmt.Println("开始智简魔方认证流程...")

	// 方法1: 尝试标准API认证
	cookie, err := s.tryStandardAuth()
	if err == nil && cookie != "" {
		fmt.Printf("标准认证成功，获得cookie: %s\n", cookie)
		return cookie, nil
	}
	fmt.Printf("标准认证失败: %v\n", err)

	// 方法2: 从JavaScript重定向提取cookie
	fmt.Println("尝试从JavaScript重定向提取cookie...")
	cookie, err = s.extractCookieFromJSRedirect()
	if err == nil && cookie != "" {
		fmt.Printf("JavaScript重定向认证成功，获得cookie: %s\n", cookie)
		return cookie, nil
	}
	fmt.Printf("JavaScript重定向认证失败: %v\n", err)

	// 方法3: 尝试直接使用API密钥作为认证
	fmt.Println("尝试使用API密钥直接认证...")
	return s.tryDirectAPIAuth()
}

// tryStandardAuth 尝试标准API认证
func (s *ZJMFService) tryStandardAuth() (string, error) {
	// 尝试常见的认证端点
	authEndpoints := []string{
		"/api/auth/login",
		"/auth/login",
		"/api/login",
		"/admin/auth/login",
	}

	for _, endpoint := range authEndpoints {
		fmt.Printf("尝试标准认证端点: %s%s\n", s.BaseURL, endpoint)

		authPayload := map[string]interface{}{
			"username": s.APIKey,
			"password": s.APISecret,
			"api_key":  s.APIKey,
			"secret":   s.APISecret,
		}

		response, err := s.makeRequest("POST", endpoint, authPayload)
		if err != nil {
			fmt.Printf("认证端点 %s 请求失败: %v\n", endpoint, err)
			continue
		}

		// 尝试从响应中提取cookie或token
		// 这里需要根据实际API响应格式调整
		fmt.Printf("认证端点 %s 响应: %s\n", endpoint, string(response[:min(len(response), 200)]))
	}

	return "", fmt.Errorf("所有标准认证端点都失败")
}

// tryDirectAPIAuth 尝试直接API认证
func (s *ZJMFService) tryDirectAPIAuth() (string, error) {
	// 获取yxd_token并构造完整cookie
	fmt.Println("获取yxd_token用于直接认证...")

	// 首先获取yxd_token，处理JavaScript重定向
	yxdToken, err := s.getYxdToken()
	if err != nil {
		fmt.Printf("获取yxd_token失败: %v\n", err)
		// 备用方案：只使用API密钥
		basicCookie := fmt.Sprintf("api_key=%s; api_secret=%s", s.APIKey, s.APISecret)
		fmt.Printf("使用基础API认证: %s\n", basicCookie)
		return basicCookie, nil
	}

	// 构造完整的cookie，包含API密钥信息
	fullCookie := fmt.Sprintf("yxd_token=%s; api_key=%s; api_secret=%s", yxdToken, s.APIKey, s.APISecret)
	fmt.Printf("构造完整认证cookie: %s\n", fullCookie)
	return fullCookie, nil
}

// getYxdToken 获取yxd_token
func (s *ZJMFService) getYxdToken() (string, error) {
	fmt.Println("开始获取yxd_token...")

	// 先尝试直接访问根路径
	response, err := s.makeRequest("GET", "/", nil)
	if err != nil {
		return "", fmt.Errorf("获取根路径失败: %v", err)
	}

	content := string(response)
	fmt.Printf("根路径响应长度: %d 字节\n", len(response))

	// 检查是否是JavaScript重定向
	if strings.Contains(content, "window.location.href") {
		fmt.Println("检测到JavaScript重定向，尝试处理...")
		// 处理JavaScript重定向
		processedResponse, err := s.handleRootJSRedirect(response)
		if err != nil {
			return "", fmt.Errorf("处理根路径JavaScript重定向失败: %v", err)
		}
		content = string(processedResponse)
	}

	// 从内容中提取yxd_token
	if strings.Contains(content, "yxd_token=") {
		re := regexp.MustCompile(`yxd_token=([a-f0-9]+)`)
		matches := re.FindStringSubmatch(content)
		if len(matches) > 1 {
			yxdToken := matches[1]
			fmt.Printf("成功提取yxd_token: %s\n", yxdToken)
			return yxdToken, nil
		}
	}

	return "", fmt.Errorf("无法从响应中提取yxd_token")
}

// handleRootJSRedirect 处理根路径的JavaScript重定向
func (s *ZJMFService) handleRootJSRedirect(response []byte) ([]byte, error) {
	content := string(response)

	// 提取重定向路径
	re := regexp.MustCompile(`window\.location\.href='([^']+)'`)
	matches := re.FindStringSubmatch(content)
	if len(matches) > 1 {
		redirectPath := matches[1]
		fmt.Printf("检测到根路径重定向: %s\n", redirectPath)

		// 跟随重定向
		finalResponse, err := s.makeRequest("GET", redirectPath, nil)
		if err != nil {
			return nil, fmt.Errorf("跟随根路径重定向失败: %v", err)
		}

		fmt.Printf("根路径重定向后响应长度: %d 字节\n", len(finalResponse))
		return finalResponse, nil
	}

	return nil, fmt.Errorf("无法解析根路径JavaScript重定向路径")
}

// extractCookieFromJSRedirect 从JavaScript重定向中提取cookie
func (s *ZJMFService) extractCookieFromJSRedirect() (string, error) {
	// 访问根路径获取JavaScript重定向
	response, err := s.makeRequest("GET", "/", nil)
	if err != nil {
		return "", err
	}

	content := string(response)
	fmt.Printf("根路径响应内容长度: %d 字节\n", len(response))

	// 检查是否是JavaScript重定向
	if strings.Contains(content, "window.location.href") {
		fmt.Println("检测到根路径JavaScript重定向，尝试处理...")
		// 处理JavaScript重定向
		processedResponse, err := s.handleRootJSRedirect(response)
		if err != nil {
			return "", fmt.Errorf("处理根路径JavaScript重定向失败: %v", err)
		}
		content = string(processedResponse)
	}

	// 查找cookie设置代码
	if strings.Contains(content, "yxd_token=") {
		// 使用正则表达式提取cookie值
		re := regexp.MustCompile(`yxd_token=([a-f0-9]+)`)
		matches := re.FindStringSubmatch(content)
		if len(matches) > 1 {
			cookie := fmt.Sprintf("yxd_token=%s", matches[1])
			fmt.Printf("从JavaScript中提取到cookie: %s\n", cookie)

			// 验证cookie是否有效
			isValid, err := s.verifyCookie(cookie)
			if isValid && err == nil {
				fmt.Printf("Cookie验证成功: %s\n", cookie)
				return cookie, nil
			} else {
				fmt.Printf("Cookie验证失败: %v\n", err)
			}
		}
	}

	return "", fmt.Errorf("无法从JavaScript重定向中提取有效cookie")
}

// verifyCookie 验证cookie是否有效
func (s *ZJMFService) verifyCookie(cookie string) (bool, error) {
	fmt.Printf("验证cookie: %s\n", cookie)

	// 尝试访问一个需要认证的端点来验证cookie
	testEndpoints := []string{"/api/test", "/test", "/api/status"}

	for _, endpoint := range testEndpoints {
		response, err := s.makeRequestWithCookie("GET", endpoint, nil, cookie)
		if err != nil {
			continue
		}

		// 如果返回JSON而不是HTML重定向，说明cookie有效
		if len(response) > 0 && response[0] == '{' {
			fmt.Printf("Cookie验证成功，端点 %s 返回JSON数据\n", endpoint)
			return true, nil
		}

		// 检查是否仍然是JavaScript重定向
		if strings.Contains(string(response), "window.location.href") {
			fmt.Printf("端点 %s 仍需要重定向\n", endpoint)
			continue
		}
	}

	return false, fmt.Errorf("所有测试端点都需要重定向")
}

// makeRequestWithCookie 带cookie的请求方法
func (s *ZJMFService) makeRequestWithCookie(method, endpoint string, payload interface{}, cookie string) ([]byte, error) {
	// 清理endpoint，确保没有多余的斜杠
	cleanEndpoint := endpoint
	if len(cleanEndpoint) > 0 && cleanEndpoint[0] == '/' {
		cleanEndpoint = cleanEndpoint[1:]
	}

	url := s.BaseURL
	if url[len(url)-1] != '/' {
		url += "/"
	}
	url += cleanEndpoint

	fmt.Printf("发起带cookie的请求: %s %s\n", method, url)

	var req *http.Request
	var err error

	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return nil, err
		}
	}

	// 设置API认证头
	req.Header.Set("Authorization", "Bearer "+s.APIKey)
	req.Header.Set("X-API-Key", s.APIKey)
	req.Header.Set("X-API-Secret", s.APISecret)

	// 添加用户代理
	req.Header.Set("User-Agent", "Go-Biz-Admin/1.0")

	// 添加接受头
	req.Header.Set("Accept", "application/json")

	// 添加cookie
	req.Header.Set("Cookie", cookie)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Printf("响应状态码: %d\n", resp.StatusCode)
	fmt.Printf("响应头: Content-Type=%s\n", resp.Header.Get("Content-Type"))

	// 检查重定向
	if resp.StatusCode == 301 || resp.StatusCode == 302 || resp.StatusCode == 307 || resp.StatusCode == 308 {
		location := resp.Header.Get("Location")
		fmt.Printf("检测到重定向到: %s\n", location)
		if location != "" {
			// 递归调用处理重定向
			return s.makeRequestWithCookie(method, location, payload, cookie)
		}
	}

	// 检查是否是HTML错误页面
	contentType := resp.Header.Get("Content-Type")
	if contentType != "" && (contentType == "text/html" || contentType == "text/html; charset=UTF-8") {
		content := string(body)
		if len(content) > 500 {
			content = content[:500]
		}
		fmt.Printf("收到HTML响应，内容: %s\n", content)

		// 检查是否是JavaScript重定向
		if strings.Contains(content, "window.location.href") {
			fmt.Printf("检测到JavaScript重定向，尝试处理...\n")

			// 提取 yxd_token
			var yxdToken string
			tokenRe := regexp.MustCompile(`yxd_token=([a-f0-9]+)`)
			tokenMatches := tokenRe.FindStringSubmatch(content)
			if len(tokenMatches) > 1 {
				yxdToken = tokenMatches[1]
				fmt.Printf("提取到 yxd_token: %s\n", yxdToken)
			}

			// 提取重定向路径
			hrefRe := regexp.MustCompile(`window\.location\.href='([^']+)'`)
			hrefMatches := hrefRe.FindStringSubmatch(content)
			if len(hrefMatches) > 1 {
				redirectPath := hrefMatches[1]
				fmt.Printf("提取到重定向路径: %s\n", redirectPath)

				// 构造带 cookie 的新请求
				if yxdToken != "" {
					cookie := fmt.Sprintf("yxd_token=%s", yxdToken)
					return s.makeRequestWithCookie(method, redirectPath, payload, cookie)
				}
			}

			return nil, fmt.Errorf("无法提取JavaScript重定向信息")
		}

		return nil, fmt.Errorf("received HTML response instead of JSON, status: %d", resp.StatusCode)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}
