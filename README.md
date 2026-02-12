# 仙林云 - idc业务管理系统

<div align="center">

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.19-blue.svg)](https://golang.org/)
[![Vue.js](https://img.shields.io/badge/vue.js-3.x-brightgreen.svg)](https://vuejs.org/)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)]()

一个基于Go + Vue3的现代化业务管理系统，专为IDC服务商设计，集成智简魔方API，提供完整的供应商管理、产品同步、订单处理等功能。

</div>

## 🌟 项目特色

- **全栈现代化架构**：Go Gin后端 + Vue3前端 + Element Plus UI
- **智能供应商对接**：深度集成智简魔方IDC系统API
- **自动化产品同步**：一键同步上游产品信息到本地系统
- **完善的业务流程**：用户管理、订单处理、财务管理、工单系统
- **响应式设计**：适配各种设备屏幕尺寸

## 🏗️ 技术架构

### 后端技术栈
- **语言**：Go 1.19+
- **框架**：Gin Web Framework
- **ORM**：GORM
- **数据库**：MySQL 8.0+
- **定时任务**：内置cron调度器

### 前端技术栈
- **框架**：Vue 3 Composition API
- **UI库**：Element Plus
- **构建工具**：Vite
- **HTTP客户端**：Axios
- **路由**：Vue Router 4

## 🚀 快速开始

### 环境要求

- Go 1.19 或更高版本
- Node.js 16.0 或更高版本
- MySQL 8.0 或更高版本
- Git

### 安装步骤

#### 1. 克隆项目

```bash
git clone https://github.com/yourusername/go-biz-admin.git
cd go-biz-admin
```

#### 2. 后端配置

```bash
cd go-biz-admin
# 复制环境配置文件
cp .env.example .env

# 编辑配置文件，设置数据库连接信息
vim .env

# 安装Go依赖
go mod tidy

# 启动后端服务
go run cmd/main.go
```

#### 3. 前端配置

```bash
# 返回项目根目录
cd ..

# 安装前端依赖
npm install

# 启动开发服务器
npm run dev
```

### 环境变量配置

创建 `.env` 文件在 `go-biz-admin` 目录下：

```env
# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=go_biz_admin
DB_CHARSET=utf8mb4

# 服务器配置
SERVER_PORT=8081
JWT_SECRET=your_jwt_secret_key

# 智简魔方API配置
ZJMF_API_BASE_URL=https://your-zjmf-domain.com
ZJMF_API_KEY=your_api_key
ZJMF_API_SECRET=your_api_secret
```

## 📁 项目结构

```
go-biz-admin/
├── api/           # API路由定义
├── cmd/           # 程序入口
├── config/        # 配置管理
├── cron/          # 定时任务
├── handlers/      # 请求处理器
├── middleware/    # 中间件
├── migrations/    # 数据库迁移
├── models/        # 数据模型
├── services/      # 业务逻辑
├── .env           # 环境配置
├── go.mod         # Go模块文件
└── go.sum         # Go依赖校验

src/               # 前端源码
├── components/    # Vue组件
├── services/      # API服务
├── views/         # 页面视图
├── App.vue        # 根组件
├── main.js        # 入口文件
└── router.js      # 路由配置
```

## 🎯 核心功能

### 1. 供应商管理
- 🔄 智简魔方API自动对接
- 📊 供应商信息统一管理
- ⚡ 产品信息一键同步
- 🔧 多种接口类型支持

### 2. 产品管理
- 📦 上游产品自动同步
- 🏷️ 产品状态实时监控
- 🔍 产品搜索和筛选
- 📈 库存和价格管理

### 3. 订单系统
- 🛒 产品订单全流程管理
- 🔁 续费订单自动化处理
- 💰 支付状态实时跟踪
- 📊 订单统计和报表

### 4. 财务管理
- 💳 交易记录管理
- 🧾 账单生成和管理
- 💵 支付处理集成
- 📈 财务数据分析

### 5. 工单系统
- 🎫 工单创建和分配
- 💬 实时消息沟通
- 📊 工单状态跟踪
- 📈 工单统计分析

### 6. 用户管理
- 👥 用户账户管理
- 🔐 权限控制系统
- 📋 用户资料维护
- 📊 用户行为分析

## 🔧 开发指南

### API文档

后端API采用RESTful设计，主要端点包括：

```
# 认证相关
POST /api/auth/login      # 用户登录
POST /api/auth/register   # 用户注册

# 供应商管理
GET  /api/suppliers       # 获取供应商列表
POST /api/suppliers       # 创建供应商
POST /api/suppliers/:id/sync-products  # 同步产品

# 产品管理
GET  /api/products        # 获取产品列表
POST /api/products        # 创建产品

# 订单管理
GET  /api/orders/products # 获取产品订单
POST /api/orders/products # 创建产品订单
```

### 数据库设计

主要数据表包括：
- `users` - 用户信息表
- `suppliers` - 供应商信息表
- `upstream_products` - 上游产品表
- `orders` - 订单表
- `bills` - 账单表
- `tickets` - 工单表

### 代码规范

**Go后端**：
- 遵循Go官方代码风格
- 使用Goland或VSCode开发
- 单元测试覆盖率 > 80%

**Vue前端**：
- 使用Composition API
- 遵循ESLint规范
- 组件化开发模式

## 🐛 常见问题

### 1. 数据库连接失败
```
检查.env文件中的数据库配置是否正确
确认MySQL服务是否正常运行
验证数据库用户权限设置
```

### 2. API同步失败
```
确认智简魔方API配置正确
检查网络连接和防火墙设置
验证API密钥和密钥是否有效
```

### 3. 前端页面空白
```
确认Node.js版本符合要求
重新安装前端依赖
检查浏览器控制台错误信息
```

## 🤝 贡献指南

欢迎提交Issue和Pull Request！

1. Fork项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启Pull Request

## 📄 许可证

本项目采用MIT许可证 - 查看 [LICENSE](LICENSE) 文件了解更多详情。

## 📞 联系方式

- 邮箱: yangxiaogiao@163.com

## 🙏 致谢

感谢以下开源项目的支持：

- [Gin Framework](https://github.com/gin-gonic/gin)
- [GORM](https://gorm.io/)
- [Vue.js](https://vuejs.org/)
- [Element Plus](https://element-plus.org/)
- [Vite](https://vitejs.dev/)

---

<p align="center">
  Made with ❤️ by xianlincloud Team
</p>