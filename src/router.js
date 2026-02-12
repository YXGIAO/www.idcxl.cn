import { createRouter, createWebHistory } from 'vue-router';
import Home from './components/Home.vue';
import Login from './components/Login.vue';
import Register from './components/Register.vue';
import UserProfile from './components/UserProfile.vue';
import ProductPurchase from './components/ProductPurchase.vue';
import AdminApp from './admin.vue';
import UsersManagement from './views/admin/UsersManagement.vue';
import OrdersManagement from './views/admin/OrdersManagement.vue';
import FinanceManagement from './views/admin/FinanceManagement.vue';
import TicketsManagement from './views/admin/TicketsManagement.vue';
import SystemSettings from './views/admin/SystemSettings.vue';
import ProductsManagement from './views/admin/ProductsManagement.vue';
import ServersManagement from './views/admin/ServersManagement.vue';
import SuppliersManagement from './views/admin/SuppliersManagement.vue';
import TasksManagement from './views/admin/TasksManagement.vue';
import ZJMFManagement from './views/admin/ZJMFManagement.vue'; // 智简魔方管理组件
import ResourcesManagement from './views/admin/ResourcesManagement.vue'; // 资源管理组件
import DashboardOverview from './views/admin/DashboardOverview.vue'; // 仪表盘组件
import UpstreamProductsManagement from './views/admin/UpstreamProductsManagement.vue'; // 上游产品管理组件
import UpstreamProductDetail from './views/admin/UpstreamProductDetail.vue'; // 上游产品详情组件

// 定义路由
const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/login',
    name: 'Login',
    component: Login
  },
  {
    path: '/register',
    name: 'Register',
    component: Register
  },
  {
    path: '/admin',
    name: 'AdminApp',
    component: AdminApp,
    meta: { requiresAuth: true, requiresAdmin: true },
    children: [
      {
        path: 'users',
        name: 'UsersManagement',
        component: UsersManagement,
        meta: { requiresAuth: true, requiresAdmin: true },
      },
      {
        path: 'orders',
        name: 'OrdersManagement',
        component: OrdersManagement,
        meta: { requiresAuth: true, requiresAdmin: true },
      },
      {
        path: 'finance',
        name: 'FinanceManagement',
        component: FinanceManagement,
        meta: { requiresAuth: true, requiresAdmin: true },
      },
      {
        path: 'tickets',
        name: 'TicketsManagement',
        component: TicketsManagement,
        meta: { requiresAuth: true, requiresAdmin: true },
      },
      {
        path: 'system-settings',
        name: 'SystemSettings',
        component: SystemSettings,
        meta: { requiresAuth: true, requiresAdmin: true },
      },
      {
        path: 'products',
        name: 'ProductsManagement',
        component: ProductsManagement,
        meta: { requiresAuth: true, requiresAdmin: true },
      },
      {
        path: 'servers',
        name: 'ServersManagement',
        component: ServersManagement,
        meta: { requiresAuth: true, requiresAdmin: true },
      },
      {
        path: 'suppliers',
        name: 'SuppliersManagement',
        component: SuppliersManagement,
        meta: { requiresAuth: true, requiresAdmin: true },
      },
      {
        path: 'tasks',
        name: 'TasksManagement',
        component: TasksManagement,
        meta: { requiresAuth: true, requiresAdmin: true },
      },
      {
        path: 'zjmf',
        name: 'ZJMFManagement',
        component: ZJMFManagement,
        meta: { requiresAuth: true, requiresAdmin: true },
      },
      {
        path: 'resources',
        name: 'ResourcesManagement',
        component: ResourcesManagement,
        meta: { requiresAuth: true, requiresAdmin: true },
      },
      {
        path: 'dashboard',
        name: 'DashboardOverview',
        component: DashboardOverview,
        meta: { requiresAuth: true, requiresAdmin: true },
      },
      {
        path: 'upstream-products',
        name: 'UpstreamProductsManagement',
        component: UpstreamProductsManagement,
        meta: { requiresAuth: true, requiresAdmin: true },
      },
      {
        path: 'upstream-products/:id',
        name: 'UpstreamProductDetail',
        component: UpstreamProductDetail,
        meta: { requiresAuth: true, requiresAdmin: true },
      },
    ],
  },
  // 重定向 /admin 到 /admin/dashboard
  {
    path: '/admin', // 匹配 /admin 路径
    redirect: { name: 'DashboardOverview' },
  },
  {
    path: '/profile',
    name: 'UserProfile',
    component: UserProfile,
    meta: { requiresAuth: true }
  },
  {
    path: '/product/:id',
    name: 'ProductPurchase',
    component: ProductPurchase,
    meta: { requiresAuth: true }
  }
];

// 创建路由实例
const router = createRouter({
  history: createWebHistory(),
  routes
});

// 简单的fetch包装器来获取用户信息
async function checkAdminAccess(to, next) {
  try {
    const token = localStorage.getItem('token');
    if (!token) {
      throw new Error('No token found');
    }
    
    // 使用相对路径，配合Vite代理配置
    const response = await fetch('/api/user', {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    });
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    const result = await response.json();
    const userData = result.data || result;
    const userRole = userData.role || 'user';
    
    console.log('[路由守卫] 实时获取的用户角色:', userRole);
    
    // 更新本地存储的角色信息
    localStorage.setItem('userRole', userRole);
    
    if (userRole === 'admin') {
      console.log('[路由守卫] 管理员用户，允许访问:', to.path);
      next();
    } else {
      console.log('[路由守卫] 非管理员用户，重定向到首页');
      alert('您没有权限访问管理后台，请使用管理员账户登录！');
      next('/');
    }
  } catch (error) {
    console.error('[路由守卫] 获取用户信息失败:', error);
    alert('获取用户信息失败，请重新登录');
    localStorage.removeItem('token');
    localStorage.removeItem('userRole');
    next('/login');
  }
}

// 路由守卫
router.beforeEach(async (to, from, next) => {
  const token = localStorage.getItem('token');
  
  console.log('[路由守卫] 当前路径:', to.path);
  console.log('[路由守卫] 需要认证:', to.meta.requiresAuth);
  console.log('[路由守卫] 需要管理员:', to.meta.requiresAdmin);
  console.log('[路由守卫] Token存在:', !!token);
  
  if (to.meta.requiresAuth && !token) {
    console.log('[路由守卫] 未登录，重定向到登录页');
    next('/login');
  } else if (to.meta.requiresAdmin) {
    // 对于需要管理员权限的路由，重新获取用户角色信息
    await checkAdminAccess(to, next);
  } else {
    console.log('[路由守卫] 允许访问:', to.path);
    next();
  }
});

export default router;