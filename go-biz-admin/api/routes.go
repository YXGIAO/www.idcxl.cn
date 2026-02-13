package api

import (
	"io"
	"log"
	"os"
	"time"

	"go-biz-admin/config"
	"go-biz-admin/handlers"
	"go-biz-admin/middleware"
	"go-biz-admin/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter 配置API路由
func SetupRouter() *gin.Engine {
	// 尝试打开日志文件（logs/server.log），如果失败则继续并打印到 stdout
	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Println("创建 logs 目录失败：", err)
	}
	logFile, err := os.OpenFile("logs/server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Println("无法打开日志文件 logs/server.log：", err)
	} else {
		gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)
		gin.DefaultErrorWriter = io.MultiWriter(logFile, os.Stderr)
	}

	r := gin.New()
	// 日志和自定义恢复中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 自定义CORS配置，允许所有来源，允许凭据，允许自定义头部
	configC := cors.DefaultConfig()
	configC.AllowOrigins = []string{"*"}
	configC.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD", "PATCH"}
	configC.AllowHeaders = []string{"*"}
	configC.ExposeHeaders = []string{"Content-Length", "Content-Type", "X-Requested-With"}
	configC.AllowCredentials = true
	configC.MaxAge = 12 * time.Hour
	r.Use(cors.New(configC))

	// 创建服务实例
	userService := services.NewUserService()
	systemService := services.NewSystemService()
	zjmfService := services.NewZJMFService(
		config.AppConfig.ZJMF.BaseURL,
		config.AppConfig.ZJMF.APIKey,
		config.AppConfig.ZJMF.APISecret,
	)
	upstreamProductService := services.NewUpstreamProductService()

	// 创建处理器实例
	userHandler := handlers.NewUserHandler(userService)
	systemHandler := handlers.NewSystemHandler(systemService)
	zjmfHandler := handlers.NewZJMFHandler(zjmfService)
	upstreamProductHandler := handlers.NewUpstreamProductHandler(upstreamProductService)
	orderHandler := handlers.NewOrderHandler(services.NewOrderService())
	financeHandler := handlers.NewFinanceHandler(services.NewFinanceService())
	ticketHandler := handlers.NewTicketHandler(services.NewTicketService())

	// API路由组
	api := r.Group("/api")
	{
		// 用户认证相关路由
		auth := api.Group("/auth")
		{
			auth.POST("/login", userHandler.Login)
			auth.POST("/register", userHandler.Register)
		}

		// 用户管理相关路由
		users := api.Group("/users")
		{
			users.GET("", userHandler.GetUsers)
			users.POST("", userHandler.Register)
			users.POST("/admin", userHandler.CreateAdminAccount)
		}

		// 当前用户信息API
		currentUser := api.Group("/user")
		currentUser.Use(middleware.JWTAuthMiddleware())
		{
			currentUser.GET("", userHandler.GetCurrentUser)
		}

		// 供应商管理相关路由
		suppliers := api.Group("/suppliers")
		{
			suppliers.GET("", systemHandler.GetSuppliers)
			suppliers.GET("/:id", systemHandler.GetSupplier)
			suppliers.POST("", systemHandler.CreateSupplier)
			suppliers.DELETE("/:id", systemHandler.DeleteSupplier)
			suppliers.POST("/:id/sync", systemHandler.SyncSupplierInfo)
			suppliers.POST("/:id/sync-products", systemHandler.SyncProductsFromSupplier)
		}

		// 智简魔方相关路由
		zjmf := api.Group("/zjmf")
		{
			zjmf.GET("/suppliers", zjmfHandler.GetSuppliers)
			zjmf.GET("/user/:id", zjmfHandler.GetUserDetail)
			zjmf.POST("/server", zjmfHandler.CreateServer)
			zjmf.POST("/sync-servers", zjmfHandler.SyncServers)
			zjmf.POST("/info", zjmfHandler.GetInfo)
			zjmf.POST("/products", zjmfHandler.GetProducts)
		}

		// 任务管理相关路由
		tasks := api.Group("/tasks")
		tasks.Use(middleware.JWTAuthMiddleware()) // 需要认证才能访问任务API
		{
			tasks.GET("", systemHandler.GetTasks)
			tasks.GET("/:id", systemHandler.GetTask)
			tasks.POST("", systemHandler.CreateTask)
			tasks.PUT("/:id", systemHandler.UpdateTask)
			tasks.DELETE("/:id", systemHandler.DeleteTask)
			tasks.POST("/:id/run", systemHandler.RunTask)
		}

		// 订单相关路由
		orders := api.Group("/orders")
		orders.Use(middleware.JWTAuthMiddleware()) // 需要认证
		{
			orders.GET("/products", orderHandler.GetProductOrders)
			orders.GET("/products/:id", orderHandler.GetProductOrder)
			orders.POST("/products", orderHandler.CreateProductOrder)
			orders.PUT("/products/:id", orderHandler.UpdateProductOrder)

			orders.GET("/renewals", orderHandler.GetRenewalOrders)
			orders.GET("/renewals/:id", orderHandler.GetRenewalOrder)
			orders.POST("/renewals", orderHandler.CreateRenewalOrder)
			orders.PUT("/renewals/:id", orderHandler.UpdateRenewalOrder)
		}

		// 财务相关路由
		finance := api.Group("/finance")
		finance.Use(middleware.JWTAuthMiddleware()) // 需要认证
		{
			finance.GET("/transactions", financeHandler.GetTransactions)
			finance.GET("/transactions/:id", financeHandler.GetTransaction)
			finance.GET("/bills", financeHandler.GetBills)
			finance.GET("/bills/:id", financeHandler.GetBill)
			finance.PUT("/bills/:id/pay", financeHandler.PayBill)
			finance.POST("/bills", financeHandler.CreateBill)
			finance.PUT("/bills/:id", financeHandler.UpdateBill)
		}

		// 工单相关路由
		tickets := api.Group("/tickets")
		tickets.Use(middleware.JWTAuthMiddleware()) // 需要认证
		{
			tickets.GET("", ticketHandler.GetTickets)
			tickets.GET("/:id", ticketHandler.GetTicket)
			tickets.POST("", ticketHandler.CreateTicket)
			tickets.PUT("/:id", ticketHandler.UpdateTicket)
			tickets.GET("/stats", ticketHandler.GetTicketStats)
		}

		// 系统设置相关路由
		system := api.Group("/system")
		system.Use(middleware.JWTAuthMiddleware()) // 需要认证
		{
			system.GET("/settings", systemHandler.GetSystemSettings)
			system.PUT("/settings", systemHandler.UpdateSystemSettings)
		}

		// 产品相关路由
		products := api.Group("/products")
		products.Use(middleware.JWTAuthMiddleware()) // 需要认证
		{
			products.GET("", systemHandler.GetProducts)
			products.GET("/:id", systemHandler.GetProduct)
			products.POST("", systemHandler.CreateProduct)
			products.PUT("/:id", systemHandler.UpdateProduct)
			products.DELETE("/:id", systemHandler.DeleteProduct)
		}
	}

	// 管理员API路由组
	admin := r.Group("/admin")
	admin.Use(middleware.JWTAuthMiddleware(), middleware.AdminAuthMiddleware())
	{
		v1 := admin.Group("/v1")
		{
			v1.GET("/supplier/:id", systemHandler.AdminGetSupplierDetail)

			// 上游产品管理路由
			upstream := v1.Group("/upstream")
			{
				upstream.GET("/host", upstreamProductHandler.GetProductList)
				upstream.GET("/host/:id", upstreamProductHandler.GetProductByID)
				upstream.POST("/host", upstreamProductHandler.CreateProduct)
				upstream.PUT("/host/:id", upstreamProductHandler.UpdateProduct)
				upstream.DELETE("/host/:id", upstreamProductHandler.DeleteProduct)
				upstream.GET("/host/status-options", upstreamProductHandler.GetStatusOptions)
				upstream.GET("/host/billing-cycle-options", upstreamProductHandler.GetBillingCycleOptions)
			}
		}
	}

	return r
}
