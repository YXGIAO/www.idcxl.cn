# Go-Biz-Admin 后台管理系统 PRD

## 1. 产品概述

### 1.1 产品简介

Go-Biz-Admin 是一个基于 Go + Vue.js 技术栈的企业级后台管理系统，主要用于管理服务器资源、供应商、订单、财务等业务。系统支持集成智简魔方（ZJMF）等第三方供应商平台，实现产品信息的自动同步和统一管理。

### 1.2 产品定位

- **目标用户**：IDC服务商、云服务器代理商、IT资源管理团队
- **核心价值**：统一管理多供应商资源，自动化同步产品信息，简化订单和财务管理流程
- **应用场景**：云服务器管理、VPS销售、主机产品代理等业务场景

### 1.3 技术架构

```
前端：Vue.js + Element Plus
后端：Go + Gin + GORM
数据库：MySQL
第三方集成：智简魔方API
```

---

## 2. 功能模块

### 2.1 用户管理模块

#### 2.1.1 用户注册/登录
- 用户注册功能（邮箱/手机号）
- 用户登录功能（支持JWT认证）
- 密码找回功能

#### 2.1.2 用户信息管理
- 查看用户列表
- 创建管理员账户
- 查看用户详情
- 用户信息编辑
- 用户状态管理（启用/禁用）

#### 2.1.3 权限管理
- 基于角色的权限控制（RBAC）
- 管理员权限与普通用户权限区分
- JWT Token 认证机制

---

### 2.2 供应商管理模块

#### 2.2.1 供应商基础信息
- 供应商名称、描述
- 接口类型（手动/智简魔方/v10）
- 联系方式
- 状态管理（活跃/非活跃/暂停）

#### 2.2.2 供应商接口配置
- 接口地址配置
- API 密钥管理
- API 密钥加密存储

#### 2.2.3 智简魔方集成
- 自动配置功能（一键连接测试）
- API 密钥认证
- 供应商信息同步（余额、产品数量等）
- 产品信息同步（自动从智简魔方获取产品列表）

#### 2.2.4 财务信息
- 账户余额显示
- 货币汇率配置
- 汇率自动更新
- 产品统计（正常数量、总数量）

---

### 2.3 上游产品管理模块

#### 2.3.1 产品信息展示
- 产品名称和标识
- 产品状态（未付款/开通中/已开通/已暂停/已删除/开通失败）
- 产品基础信息展示

#### 2.3.2 财务信息
- 首次付款金额
- 续费金额
- 付款周期（月付/年付等）
- 计费周期时间

#### 2.3.3 时间信息
- 开通时间（Unix 时间戳）
- 到期时间（Unix 时间戳）
- 到期提醒功能

#### 2.3.4 用户信息
- 用户 ID 和用户名
- 公司名称
- 邮箱和手机号

#### 2.3.5 配置信息
- IP 数量
- 产品配置详情
- 上游产品 ID

#### 2.3.6 产品操作
- 产品列表查询（支持关键词、供应商、状态、时间筛选）
- 产品详情查看
- 产品创建/编辑/删除
- 产品导入（从供应商同步）

---

### 2.4 产品管理模块

#### 2.4.1 产品目录管理
- 产品分类
- 产品层级结构

#### 2.4.2 产品信息管理
- 产品基本信息
- 产品价格配置
- 产品规格参数
- 产品描述

#### 2.4.3 产品操作
- 创建产品
- 编辑产品
- 删除产品
- 产品上架/下架

---

### 2.5 订单管理模块

#### 2.5.1 产品订单
- 订单列表查询
- 订单详情查看
- 创建产品订单
- 更新订单状态
- 订单筛选（按时间、状态、用户等）

#### 2.5.2 续费订单
- 续费订单列表
- 创建续费订单
- 续费金额计算
- 续费状态管理

#### 2.5.3 订单流程
- 订单创建 → 待支付 → 已支付 → 开通中 → 已开通 → 到期
- 订单取消流程
- 订单退款流程

---

### 2.6 财务管理模块

#### 2.6.1 交易记录
- 交易列表查询
- 交易详情查看
- 交易类型（收入/支出）
- 交易时间筛选

#### 2.6.2 账单管理
- 账单列表查询
- 账单详情查看
- 创建账单
- 更新账单信息
- 账单支付状态
- 账单到期提醒

#### 2.6.3 财务报表
- 收入统计
- 支出统计
- 利润分析
- 趋势分析

---

### 2.7 工单管理模块

#### 2.7.1 工单创建
- 工单类型（技术支持/账单问题/其他）
- 工单优先级
- 工单描述
- 附件上传

#### 2.7.2 工单处理
- 工单列表查询
- 工单详情查看
- 工单回复
- 工单状态管理（待处理/处理中/已解决/已关闭）
- 工单分配

#### 2.7.3 工单统计
- 工单数量统计
- 处理时效分析
- 工单类型分布

---

### 2.8 服务器管理模块

#### 2.8.1 服务器信息
- 服务器名称
- 主机地址和端口
- 服务器类型（KVM/VPS/独立服务器）
- 登录凭证（用户名/密码）
- 服务器位置

#### 2.8.2 配置信息
- CPU 规格
- 内存规格
- 磁盘规格
- 带宽信息
- IP 数量

#### 2.8.3 服务器状态
- 服务器状态（活跃/非活跃/维护中）
- 智简魔方服务器状态同步
- 服务器备注信息

#### 2.8.4 服务器操作
- 创建服务器
- 编辑服务器
- 删除服务器
- 同步智简魔方服务器列表

---

### 2.9 任务管理模块

#### 2.9.1 任务列表
- 任务列表查询
- 任务详情查看
- 任务状态（待执行/执行中/已完成/失败）
- 任务执行日志

#### 2.9.2 任务操作
- 创建任务
- 编辑任务
- 删除任务
- 手动执行任务
- 任务结果查看

#### 2.9.3 定时任务
- 定时任务配置
- 任务执行频率设置
- 任务历史记录

---

### 2.10 系统设置模块

#### 2.10.1 基础设置
- 系统名称
- 系统Logo
- 系统描述

#### 2.10.2 API 配置
- 智简魔方API地址
- API 密钥配置
- 接口超时设置

#### 2.10.3 数据库配置
- 数据库连接配置
- 数据库备份设置

#### 2.10.4 其他设置
- 系统日志级别
- 文件上传限制
- 系统维护模式

---

## 3. 非功能性需求

### 3.1 性能要求
- 页面加载时间 < 2秒
- API 响应时间 < 500ms
- 支持并发用户数 > 100
- 数据库查询优化

### 3.2 安全要求
- 密码加密存储（bcrypt）
- JWT Token 认证
- API 密钥加密存储
- SQL 注入防护
- XSS 攻击防护
- CSRF 防护

### 3.3 可用性要求
- 系统可用性 > 99.9%
- 数据备份机制
- 错误日志记录
- 友好的错误提示

### 3.4 兼容性要求
- 支持 Chrome、Firefox、Edge 等主流浏览器
- 支持响应式设计
- 支持 Windows、Linux、macOS 服务器

---

## 4. 数据模型

### 4.1 用户表 (users)
- id: 主键
- username: 用户名
- password: 密码（加密）
- email: 邮箱
- role: 角色（admin/user）
- status: 状态
- created_at/updated_at: 时间戳

### 4.2 供应商表 (suppliers)
- id: 主键
- type: 类型（default/whmcs/finance）
- name: 名称
- url: 供应商地址
- username: 用户名
- api_key: API 密钥
- api_secret: API 密钥
- status: 状态
- balance: 账户余额
- currency_code: 货币代码
- rate: 汇率
- zjmf_*: 智简魔方专用字段

### 4.3 上游产品表 (upstream_products)
- id: 主键
- name: 产品标识
- product_name: 商品名称
- status: 状态
- first_payment_amount: 首次付款金额
- renew_amount: 续费金额
- billing_cycle: 付款周期
- active_time: 开通时间
- due_time: 到期时间
- upstream_host_id: 上游产品 ID
- client_id: 用户 ID
- supplier_id: 供应商 ID

### 4.4 服务器表 (servers)
- id: 主键
- name: 名称
- host: 主机地址
- port: 端口
- type: 类型
- username: 用户名
- password: 密码
- location: 位置
- status: 状态
- cpu/memory/disk/bandwidth: 配置信息
- supplier_id: 供应商 ID
- zjmf_*: 智简魔方专用字段

### 4.5 产品订单表 (product_orders)
- id: 主键
- user_id: 用户 ID
- product_id: 产品 ID
- order_no: 订单号
- amount: 金额
- status: 状态
- created_at/updated_at: 时间戳

### 4.6 续费订单表 (renewal_orders)
- id: 主键
- user_id: 用户 ID
- product_id: 产品 ID
- order_no: 订单号
- renew_amount: 续费金额
- billing_cycle: 计费周期
- status: 状态
- created_at/updated_at: 时间戳

### 4.7 交易记录表 (transactions)
- id: 主键
- user_id: 用户 ID
- type: 类型（income/expense）
- amount: 金额
- description: 描述
- created_at: 创建时间

### 4.8 账单表 (bills)
- id: 主键
- user_id: 用户 ID
- bill_no: 账单号
- amount: 金额
- status: 状态（unpaid/paid/overdue）
- due_date: 到期日期
- created_at/updated_at: 时间戳

### 4.9 工单表 (tickets)
- id: 主键
- user_id: 用户 ID
- title: 标题
- description: 描述
- status: 状态
- priority: 优先级
- created_at/updated_at: 时间戳

### 4.10 任务表 (tasks)
- id: 主键
- name: 任务名称
- type: 任务类型
- status: 状态
- schedule: 调度配置
- last_run_time: 最后执行时间
- next_run_time: 下次执行时间
- created_at/updated_at: 时间戳

---

## 5. API 接口设计

### 5.1 认证接口
- `POST /api/auth/login` - 用户登录
- `POST /api/auth/register` - 用户注册
- `GET /api/user` - 获取当前用户信息

### 5.2 供应商接口
- `GET /api/suppliers` - 获取供应商列表
- `GET /api/suppliers/:id` - 获取供应商详情
- `POST /api/suppliers` - 创建供应商
- `PUT /api/suppliers/:id` - 更新供应商
- `DELETE /api/suppliers/:id` - 删除供应商
- `POST /api/suppliers/:id/sync` - 同步供应商信息
- `POST /api/suppliers/:id/sync-products` - 同步供应商产品

### 5.3 智简魔方接口
- `GET /api/zjmf/suppliers` - 获取智简魔方供应商
- `GET /api/zjmf/user/:id` - 获取用户详情
- `POST /api/zjmf/server` - 创建服务器
- `POST /api/zjmf/sync-servers` - 同步服务器
- `POST /api/zjmf/info` - 获取智简魔方信息
- `POST /api/zjmf/products` - 获取产品列表

### 5.4 产品接口
- `GET /api/products` - 获取产品列表
- `GET /api/products/:id` - 获取产品详情
- `POST /api/products` - 创建产品
- `PUT /api/products/:id` - 更新产品
- `DELETE /api/products/:id` - 删除产品

### 5.5 订单接口
- `GET /api/orders/products` - 获取产品订单
- `GET /api/orders/products/:id` - 获取产品订单详情
- `POST /api/orders/products` - 创建产品订单
- `PUT /api/orders/products/:id` - 更新产品订单
- `GET /api/orders/renewals` - 获取续费订单
- `GET /api/orders/renewals/:id` - 获取续费订单详情
- `POST /api/orders/renewals` - 创建续费订单
- `PUT /api/orders/renewals/:id` - 更新续费订单

### 5.6 财务接口
- `GET /api/finance/transactions` - 获取交易记录
- `GET /api/finance/transactions/:id` - 获取交易详情
- `GET /api/finance/bills` - 获取账单列表
- `GET /api/finance/bills/:id` - 获取账单详情
- `POST /api/finance/bills` - 创建账单
- `PUT /api/finance/bills/:id` - 更新账单
- `PUT /api/finance/bills/:id/pay` - 支付账单

### 5.7 工单接口
- `GET /api/tickets` - 获取工单列表
- `GET /api/tickets/:id` - 获取工单详情
- `POST /api/tickets` - 创建工单
- `PUT /api/tickets/:id` - 更新工单
- `GET /api/tickets/stats` - 获取工单统计

### 5.8 任务接口
- `GET /api/tasks` - 获取任务列表
- `GET /api/tasks/:id` - 获取任务详情
- `POST /api/tasks` - 创建任务
- `PUT /api/tasks/:id` - 更新任务
- `DELETE /api/tasks/:id` - 删除任务
- `POST /api/tasks/:id/run` - 执行任务

### 5.9 系统设置接口
- `GET /api/system/settings` - 获取系统设置
- `PUT /api/system/settings` - 更新系统设置

### 5.10 上游产品接口（管理员）
- `GET /api/admin/v1/upstream/host` - 获取上游产品列表
- `GET /api/admin/v1/upstream/host/:id` - 获取上游产品详情
- `POST /api/admin/v1/upstream/host` - 创建上游产品
- `PUT /api/admin/v1/upstream/host/:id` - 更新上游产品
- `DELETE /api/admin/v1/upstream/host/:id` - 删除上游产品
- `GET /api/admin/v1/upstream/host/status-options` - 获取状态选项
- `GET /api/admin/v1/upstream/host/billing-cycle-options` - 获取计费周期选项

---

## 6. 业务流程

### 6.1 供应商配置流程
```
1. 登录系统
2. 进入供应商管理页面
3. 点击"添加供应商"或"自动配置"
4. 填写供应商信息（名称、接口地址、API密钥等）
5. 测试连接（验证API配置是否正确）
6. 保存供应商信息
7. 同步供应商信息和产品
```

### 6.2 产品订单流程
```
1. 用户选择产品
2. 填写订单信息
3. 提交订单
4. 支付订单
5. 系统自动开通产品（调用供应商API）
6. 订单状态更新为"已开通"
7. 用户收到产品信息
```

### 6.3 续费流程
```
1. 用户查看即将到期的产品
2. 选择续费周期
3. 创建续费订单
4. 支付续费金额
5. 系统更新产品到期时间
6. 发送续费成功通知
```

### 6.4 工单处理流程
```
1. 用户创建工单
2. 系统分配工单
3. 客服处理工单
4. 回复用户
5. 用户确认问题解决
6. 关闭工单
```

---

## 7. 智简魔方集成说明

### 7.1 认证机制
智简魔方 API 使用 JavaScript Cookie 认证：
1. 首次请求返回包含 `yxd_token` 的 HTML 页面
2. 提取 `yxd_token` 并设置到 Cookie
3. 后续请求携带 `yxd_token` Cookie 进行认证
4. 处理 JavaScript 重定向获取实际数据

### 7.2 API 端点
- 产品列表：`/api/v1/products`
- 服务器列表：`/servers`
- 供应商信息：`/finance/balance`

### 7.3 同步机制
- 手动同步：用户点击"同步信息"按钮触发
- 自动同步：通过定时任务定期同步供应商数据
- 增量同步：只同步有变化的数据

---

## 8. 部署说明

### 8.1 环境要求
- Go 1.16+
- Node.js 16+
- MySQL 5.7+
- 操作系统：Windows/Linux/macOS

### 8.2 环境变量配置
```bash
# 数据库配置
DB_TYPE=mysql
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=go_biz_admin

# 智简魔方配置
ZJMF_BASE_URL=https://your-zjmf-domain.com
ZJMF_API_KEY=your_api_key
ZJMF_API_SECRET=your_api_secret
```

### 8.3 部署步骤
1. 安装依赖（Go Modules、npm）
2. 配置环境变量
3. 初始化数据库
4. 启动后端服务
5. 构建前端项目
6. 部署 Nginx 反向代理

---

## 9. 版本规划

### 9.1 当前版本 (v1.0)
- [x] 用户管理
- [x] 供应商管理
- [x] 上游产品管理
- [x] 产品管理
- [x] 订单管理
- [x] 财务管理
- [x] 工单管理
- [x] 任务管理
- [x] 系统设置
- [x] 智简魔方集成

### 9.2 后续规划 (v2.0)
- [ ] 多供应商深度集成
- [ ] 自动化财务报表
- [ ] 邮件/短信通知
- [ ] 数据导入导出
- [ ] 操作日志审计
- [ ] 移动端适配

---

## 10. 附录

### 10.1 术语表
| 术语 | 说明 |
|------|------|
| ZJMF | 智简魔方，第三方供应商平台 |
| Upstream Product | 上游产品，从供应商获取的产品 |
| Billing Cycle | 计费周期，如月付、年付 |
| JWT | JSON Web Token，用于身份认证 |
| RBAC | 基于角色的访问控制 |

### 10.2 更新记录
| 版本 | 日期 | 更新内容 |
|------|------|----------|
| v1.0 | 2026-02-13 | 初始版本PRD |

---

**文档版本**: v1.0
**最后更新**: 2026-02-13
**编写人**: Claude
