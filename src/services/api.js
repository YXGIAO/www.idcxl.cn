import axios from 'axios';

// 创建axios实例
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8081/api';
const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000, // 增加超时时间
});

// 请求拦截器
apiClient.interceptors.request.use(
  config => {
    // 在发送请求之前做些什么，比如添加认证token
    const token = localStorage.getItem('token');
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`;
      config.headers['Content-Type'] = 'application/json';
    }
    return config;
  },
  error => {
    // 对请求错误做些什么
    return Promise.reject(error);
  }
);

// 响应拦截器
apiClient.interceptors.response.use(
  response => {
    // 对响应数据做点什么
    // 如果后端返回的是标准格式，直接返回数据部分
    if (response.data && typeof response.data === 'object') {
      return response.data;
    }
    return response;
  },
  error => {
    // 对响应错误做点什么
    console.error('API Error:', error);
    if (error.response && error.response.status === 401) {
      // 如果是401未授权错误，跳转到登录页
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// 用户管理API
export const userAPI = {
  // 获取用户列表
  getUsers: (params) => apiClient.get('/users', { params }),
  
  // 获取单个用户
  getUser: (id) => apiClient.get(`/users/${id}`),
  
  // 创建用户
  createUser: (data) => apiClient.post('/users', data),
  
  // 更新用户
  updateUser: (id, data) => apiClient.put(`/users/${id}`, data),
  
  // 更新用户头像
  updateUserAvatar: (id, avatarUrl) => apiClient.put(`/users/${id}/avatar`, { avatar: avatarUrl }),
  
  // 删除用户
  deleteUser: (id) => apiClient.delete(`/users/${id}`),
  
  // 获取用户资料
  getUserProfile: (id) => apiClient.get(`/users/${id}/profile`),
  
  // 更新用户资料
  updateUserProfile: (id, data) => apiClient.put(`/users/${id}/profile`, data),
  
  // 提交实名认证
  verifyIdentity: (id, data) => apiClient.post(`/users/${id}/verify-identity`, data),
};

// 当前用户API
export const currentUserAPI = {
  // 获取当前登录用户信息
  getCurrentUser: () => apiClient.get('/user'),
};

// 产品管理API
export const productAPI = {
  // 获取产品列表
  getProducts: (params) => apiClient.get('/products', { params }),
  
  // 获取单个产品
  getProduct: (id) => apiClient.get(`/products/${id}`),
  
  // 创建产品
  createProduct: (data) => {
    return apiClient.post('/products', data);
  },
  
  // 更新产品
  updateProduct: (id, data) => {
    return apiClient.put(`/products/${id}`, data);
  },
  
  // 删除产品
  deleteProduct: (id) => apiClient.delete(`/products/${id}`),
};

// 订单管理API
export const orderAPI = {
  // 获取产品订单列表
  getProductOrders: (params) => apiClient.get('/orders/products', { params }),
  
  // 获取单个产品订单
  getProductOrder: (id) => apiClient.get(`/orders/products/${id}`),
  
  // 创建产品订单
  createProductOrder: (data) => apiClient.post('/orders/products', data),
  
  // 更新产品订单
  updateProductOrder: (id, data) => apiClient.put(`/orders/products/${id}`, data),
  
  // 获取续费订单列表
  getRenewalOrders: (params) => apiClient.get('/orders/renewals', { params }),
  
  // 获取单个续费订单
  getRenewalOrder: (id) => apiClient.get(`/orders/renewals/${id}`),
  
  // 创建续费订单
  createRenewalOrder: (data) => apiClient.post('/orders/renewals', data),
  
  // 更新续费订单
  updateRenewalOrder: (id, data) => apiClient.put(`/orders/renewals/${id}`, data),
};

// 财务管理API
export const financeAPI = {
  // 获取交易流水列表
  getTransactions: (params) => apiClient.get('/finance/transactions', { params }),
  
  // 获取单个交易流水
  getTransaction: (id) => apiClient.get(`/finance/transactions/${id}`),
  
  // 获取账单列表
  getBills: (params) => apiClient.get('/finance/bills', { params }),
  
  // 获取单个账单
  getBill: (id) => apiClient.get(`/finance/bills/${id}`),
  
  // 支付账单
  payBill: (id, paymentMethod) => apiClient.put(`/finance/bills/${id}/pay`, {}, { params: { method: paymentMethod } }),
  
  // 创建账单
  createBill: (data) => apiClient.post('/finance/bills', data),
  
  // 更新账单
  updateBill: (id, data) => apiClient.put(`/finance/bills/${id}`, data),
};

// 工单管理API
export const ticketAPI = {
  // 获取工单列表
  getTickets: (params) => apiClient.get('/tickets', { params }),
  
  // 获取单个工单
  getTicket: (id) => apiClient.get(`/tickets/${id}`),
  
  // 创建工单
  createTicket: (data) => {
    return apiClient.post('/tickets', data);
  },
  
  // 更新工单
  updateTicket: (id, data) => {
    return apiClient.put(`/tickets/${id}`, data);
  },
  
  // 获取工单统计
  getTicketStats: () => apiClient.get('/tickets/stats'),
};

// 供应商管理API
export const supplierAPI = {
  // 获取供应商列表
  getSuppliers: (params) => apiClient.get('/suppliers', { params }),
  
  // 获取单个供应商
  getSupplier: (id) => apiClient.get(`/suppliers/${id}`),
  
  // 创建供应商
  createSupplier: (data) => {
    return apiClient.post('/suppliers', data);
  },
  
  // 更新供应商
  updateSupplier: (id, data) => {
    return apiClient.put(`/suppliers/${id}`, data);
  },
  
  // 删除供应商
  deleteSupplier: (id) => apiClient.delete(`/suppliers/${id}`),
  
  // 同步供应商信息
  syncSupplierInfo: (id) => apiClient.post(`/suppliers/${id}/sync`),
  
  // 从供应商同步产品
  syncProductsFromSupplier: (supplierId) => apiClient.post(`/suppliers/${supplierId}/sync-products`),
  
  // 获取产品列表
  getProducts: (params) => apiClient.get('/products', { params }),
};

// 服务器管理API
export const serverAPI = {
  // 获取服务器列表
  getServers: (params) => apiClient.get('/servers', { params }),
  
  // 获取单个服务器
  getServer: (id) => apiClient.get(`/servers/${id}`),
  
  // 创建服务器
  createServer: (data) => {
    return apiClient.post('/servers', data);
  },
  
  // 更新服务器
  updateServer: (id, data) => {
    return apiClient.put(`/servers/${id}`, data);
  },
  
  // 删除服务器
  deleteServer: (id) => apiClient.delete(`/servers/${id}`),
};

// 任务管理API
export const taskAPI = {
  // 获取任务列表
  getTasks: (params) => apiClient.get('/tasks', { params }),
  
  // 获取单个任务
  getTask: (id) => apiClient.get(`/tasks/${id}`),
  
  // 创建任务
  createTask: (data) => {
    return apiClient.post('/tasks', data);
  },
  
  // 更新任务
  updateTask: (id, data) => {
    return apiClient.put(`/tasks/${id}`, data);
  },
  
  // 运行任务
  runTask: (id) => apiClient.post(`/tasks/${id}/run`),
};

// 系统设置API
export const systemAPI = {
  // 获取系统设置
  getSystemSettings: () => apiClient.get('/system/settings'),
  
  // 更新系统设置
  updateSystemSettings: (data) => apiClient.put('/system/settings', data),
  
  // 获取系统统计信息
  getSystemStats: () => apiClient.get('/system/stats'),
};

// 智简魔方相关API
export const zjmfAPI = {
  // 获取智简魔方供应商列表
  getZJMFSuppliers: (params) => apiClient.get('/zjmf/suppliers', { params }),
  
  // 获取智简魔方用户详情
  getZJMFUser: (userId) => apiClient.get(`/zjmf/user/${userId}`),
  
  // 在智简魔方系统中创建服务器
  createZJMFServer: (serverData) => apiClient.post('/zjmf/server', serverData),
  
  // 同步智简魔方服务器到本地系统
  syncZJMFServers: () => apiClient.post('/zjmf/sync-servers'),
  
  // 获取智简魔方供应商信息
  getZJMFInfo: (supplierData) => {
    // 发送POST请求到后端，让后端调用智简魔方API
    return apiClient.post('/zjmf/info', supplierData);
  },
  
  // 获取智简魔方产品信息
  getZJMFProducts: (supplierData) => {
    // 发送POST请求到后端，让后端调用智简魔方API获取产品信息
    return apiClient.post('/zjmf/products', supplierData);
  },
};

// 上游产品管理API
export const upstreamProductAPI = {
  // 获取上游产品列表
  getUpstreamProducts: (params) => apiClient.get('/admin/v1/upstream/host', { params }),
  
  // 获取单个上游产品详情
  getUpstreamProduct: (id) => apiClient.get(`/admin/v1/upstream/host/${id}`),
  
  // 创建上游产品
  createUpstreamProduct: (data) => apiClient.post('/admin/v1/upstream/host', data),
  
  // 更新上游产品
  updateUpstreamProduct: (id, data) => apiClient.put(`/admin/v1/upstream/host/${id}`, data),
  
  // 删除上游产品
  deleteUpstreamProduct: (id) => apiClient.delete(`/admin/v1/upstream/host/${id}`),
  
  // 获取状态选项
  getStatusOptions: () => apiClient.get('/admin/v1/upstream/host/status-options'),
  
  // 获取付款周期选项
  getBillingCycleOptions: () => apiClient.get('/admin/v1/upstream/host/billing-cycle-options'),
};

// 认证API
export const authAPI = {
  // 用户登录
  login: (credentials) => apiClient.post('/auth/login', credentials),
  
  // 用户注册
  register: (userData) => apiClient.post('/auth/register', userData),
};