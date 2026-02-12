<template>
  <div class="admin-app">
    <!-- 管理员专用头部导航 -->
    <header class="admin-header">
      <div class="header-content">
        <!-- logo -->
        <div class="logo">
          <img src="https://via.placeholder.com/40x40/409eff/white?text=仙林云" alt="仙林云" />
          <span class="site-name">仙林云 - 管理后台</span>
        </div>
        
        <!-- 导航 -->
        <nav class="admin-nav">
          <ul>
            <li><router-link to="/admin/dashboard" class="nav-link" active-class="active">仪表盘</router-link></li>
            <li><router-link to="/admin/users" class="nav-link" active-class="active">用户管理</router-link></li>
            <li><router-link to="/admin/products" class="nav-link" active-class="active">产品管理</router-link></li>
            <li><router-link to="/admin/orders" class="nav-link" active-class="active">订单管理</router-link></li>
            <li><router-link to="/admin/finance" class="nav-link" active-class="active">财务管理</router-link></li>
            <li><router-link to="/admin/tickets" class="nav-link" active-class="active">工单管理</router-link></li>
            <li><router-link to="/admin/resources" class="nav-link" active-class="active">资源管理</router-link></li>
            <li><router-link to="/admin/tasks" class="nav-link" active-class="active">任务管理</router-link></li>
            <li><router-link to="/admin/system-settings" class="nav-link" active-class="active">系统设置</router-link></li>
          </ul>
        </nav>
        
        <!-- 用户信息 -->
        <div class="user-info">
          <span class="welcome-text">欢迎,</span>
          <span class="username">{{ userName }}</span>
          <el-dropdown>
            <span class="el-dropdown-link">
              <img src="https://via.placeholder.com/30x30/409eff/white?text=头" alt="头像" class="avatar" />
              <i class="el-icon-arrow-down el-icon--right"></i>
            </span>
            <template #dropdown>
              <el-dropdown-item @click="gohome">首页</el-dropdown-item>
              <el-dropdown-item @click="goToProfile">个人资料</el-dropdown-item>
              <el-dropdown-item @click="logout">退出登录</el-dropdown-item>
            </template>
          </el-dropdown>
        </div>
      </div>
    </header>

    <!-- 管理员内容区域 -->
    <main class="admin-main">
      <router-view />
    </main>

    <!-- 管理员专用底部 -->
    <footer class="admin-footer">
      <div class="footer-content">
        <p>&copy; 2026 仙林云管理系统. 版权所有.</p>
        <p>管理员: {{ userName }} | 当前版本: v1.0.5</p>
      </div>
    </footer>
  </div>
</template>

<script>
import { ElDropdown, ElDropdownItem } from 'element-plus';
import { ArrowDown } from '@element-plus/icons-vue';

export default {
  name: 'AdminApp',
  components: {
    ElDropdown,
    ElDropdownItem,
    ArrowDown
  },
  data() {
    return {
      userName: '管理员',
      userRole: 'admin'
    };
  },
  async mounted() {
    // 获取当前用户信息并显示用户名
    try {
      const token = localStorage.getItem('token');
      if (token) {
        const response = await fetch('/api/user', {
          method: 'GET',
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
          }
        });
        
        if (response.ok) {
          const result = await response.json();
          const userData = result.data || result;
          this.userName = userData.name || userData.username || '管理员';
          this.userRole = userData.role || 'user';
        }
      }
    } catch (error) {
      console.error('获取用户信息失败:', error);
    }
  },
  methods: {
    logout() {
      // 清除本地存储的token
      localStorage.removeItem('token');
      localStorage.removeItem('userRole');
      // 跳转到登录页
      this.$router.push('/login');
    },
    goToProfile() {
      this.$router.push('/profile');
    },
    gohome() {
      this.$router.push('/');
    }
    
  }
};
</script>

<style scoped>
.admin-app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.admin-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
  position: sticky;
  top: 0;
  z-index: 1000;
}

.header-content {
  max-width: 1400px;
  margin: 0 auto;
  display: flex;
  align-items: center;
  padding: 0 20px;
  height: 60px;
  gap: 16px;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo img {
  width: 40px;
  height: 40px;
  border-radius: 50%;
}

.site-name {
  font-size: 18px;
  font-weight: bold;
}

.admin-nav ul {
  display: flex;
  list-style: none;
  margin: 0;
  padding: 0;
  gap: 10px;
  flex-wrap: wrap;
  justify-content: center;
}

.admin-nav a {
  color: rgba(255, 255, 255, 0.85);
  text-decoration: none;
  padding: 10px 16px;
  border-radius: 6px;
  transition: all 0.3s ease;
  font-weight: 500;
  font-size: 15px;
  position: relative;
  overflow: hidden;
  min-width: 75px;
  text-align: center;
  white-space: nowrap;
}

/* 下拉菜单触发器样式 */
.el-dropdown-link {
  color: rgba(255, 255, 255, 0.85) !important;
  font-weight: 500 !important;
  font-size: 15px !important;
  padding: 10px 16px !important;
  text-decoration: none !important;
  display: inline-flex !important;
  align-items: center !important;
  height: 40px !important;
  line-height: 40px !important;
}

/* 下拉菜单选项样式 */
.el-dropdown-menu__item {
  font-weight: 500;
  font-size: 15px;
  color: rgba(255, 255, 255, 0.85) !important;
  padding: 10px 16px !important;
  min-width: 75px !important;
  text-align: center !important;
  line-height: 1.5 !important;
}

.el-dropdown-menu__item:hover {
  color: white !important;
  background: rgba(255, 255, 255, 0.1) !important;
}

.admin-nav a::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(255, 255, 255, 0.1);
  opacity: 0;
  transform: translateY(100%);
  transition: transform 0.3s ease, opacity 0.3s ease;
  z-index: 0;
}

.admin-nav a:hover {
  color: white;
  text-decoration: none;
}

.admin-nav a:hover::before {
  opacity: 1;
  transform: translateY(0);
}

.admin-nav a.active {
  color: white;
  background: rgba(255, 255, 255, 0.15);
  box-shadow: inset 0 0 10px rgba(255, 255, 255, 0.2);
  font-weight: 600;
}

.admin-nav a.active::after {
  content: '';
  position: absolute;
  bottom: -2px;
  left: 10%;
  width: 80%;
  height: 3px;
  background: linear-gradient(90deg, #409EFF, #66b1ff);
  border-radius: 2px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.welcome-text {
  opacity: 0.8;
}

.username {
  font-weight: 600;
  margin-right: 8px;
}

.avatar {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  vertical-align: middle;
}

.admin-main {
  flex: 1;
  padding: 20px;
  background-color: #f5f7fa;
  min-height: calc(100vh - 120px);
}

.admin-footer {
  background: #2c3e50;
  color: #ecf0f1;
  padding: 20px;
  text-align: center;
}

.footer-content {
  max-width: 1200px;
  margin: 0 auto;
}

.footer-content p {
  margin: 5px 0;
  font-size: 14px;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .header-content {
    flex-wrap: wrap;
    justify-content: center;
  }
  
  .admin-nav {
    order: 0;
    margin: 10px 0;
    width: 100%;
    justify-content: center;
  }
  
  .user-info {
    order: 1;
    margin-left: auto;
  }
  
  .admin-nav ul {
    justify-content: center;
  }
}

@media (max-width: 768px) {
  .header-content {
    flex-direction: column;
    height: auto;
    gap: 10px;
    padding: 10px;
  }
  
  .admin-nav ul {
    flex-wrap: wrap;
    justify-content: center;
    gap: 8px;
  }
  
  .admin-nav a {
    padding: 8px 12px;
    font-size: 14px;
  }
  
  .admin-main {
    padding: 10px;
  }
}</style>
