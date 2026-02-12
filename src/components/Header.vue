<template>
  <header class="navbar">
    <div class="nav-container">
      <div class="logo">
        <img src="https://via.placeholder.com/40x40/409eff/white?text=仙林云" alt="仙林云" />
        <h1>仙林云</h1>
      </div>
      <nav class="nav-menu">
        <ul>
          <li><a href="/" @click.prevent="navigateTo('/')">首页</a></li>
          <li class="dropdown">
            <a href="#products" class="dropdown-btn" @click.prevent="scrollToSection('products')">服务器&云产品</a>
            <ul class="dropdown-content">
              <li><a href="#game-cloud" @click.prevent="scrollToSection('game-cloud')">游戏云</a></li>
              <li><a href="#net-bar-cloud" @click.prevent="scrollToSection('net-bar-cloud')">网吧云</a></li>
              <li><a href="#cloud-server" @click.prevent="scrollToSection('cloud-server')">云服务器</a></li>
              <li><a href="#cdn" @click.prevent="scrollToSection('cdn')">CDN</a></li>
              <li><a href="#virtual-host" @click.prevent="scrollToSection('virtual-host')">虚拟主机</a></li>
            </ul>
          </li>
          <li><a href="#community" @click.prevent="scrollToSection('community')">社区</a></li>
          <li><a href="#about" @click.prevent="scrollToSection('about')">关于仙林云</a></li>
        </ul>
      </nav>
      <div class="nav-auth">
        <div v-if="isLoggedIn" class="user-profile" @click="toggleUserMenu" ref="userMenuRef">
          <div class="user-info">
            <div class="user-name">{{ userName }}</div>
            <div class="user-email">{{ userEmail }}</div>
            <div class="user-role" v-if="userRole">{{ userRole === 'admin' ? '管理员' : userRole }}</div>
          </div>
          <div class="avatar">
            <img :src="userAvatar" alt="头像" />
          </div>
          <!-- 用户菜单 -->
          <div v-show="showUserMenu" class="user-dropdown-menu">
            <a href="/user-profile">个人资料</a>
            <a href="/admin" v-if="userRole === 'admin'">管理后台</a>
            <a href="#" @click.prevent="handleLogout">退出登录</a>
          </div>
        </div>
        <a v-else href="/login" class="login-btn">登录/注册</a>
      </div>
    </div>
  </header>
</template>

<script>
import { authAPI, userAPI, currentUserAPI } from '../services/api'

export default {
  name: 'Header',
  data() {
    return {
      isLoggedIn: false,
      userAvatar: 'https://via.placeholder.com/40x40/409eff/white?text=头像', // 默认头像
      userName: '用户',
      userEmail: '邮箱',
      userRole: 'user', // 默认角色为user
      showUserMenu: false
    }
  },
  mounted() {
    this.checkAuthStatus()
    // 监听全局点击事件，用于关闭用户菜单
    document.addEventListener('click', this.closeUserMenu)
  },
  beforeDestroy() {
    document.removeEventListener('click', this.closeUserMenu)
  },
  methods: {
    navigateTo(path) {
      if (this.$route.path !== path) {
        this.$router.push(path);
      }
    },
    scrollToSection(section) {
      // 如果在首页，滚动到相应部分
      if (this.$route.path === '/') {
        const element = document.getElementById(section);
        if (element) {
          element.scrollIntoView({ behavior: 'smooth' });
        }
      } else {
        // 如果不在首页，先导航到首页再滚动
        this.$router.push({ path: '/', hash: `#${section}` });
      }
    },
    checkAuthStatus() {
      const token = localStorage.getItem('token')
      this.isLoggedIn = !!token
      if (this.isLoggedIn) {
        this.fetchUserInfo()
      }
    },
    async fetchUserInfo() {
      try {
        // 直接通过API获取当前用户信息
        try {
          const response = await currentUserAPI.getCurrentUser()
          const userData = response.data || response
          
          if (userData) {
            this.userAvatar = userData.avatar || 'https://via.placeholder.com/40x40/409eff/white?text=头像'
            this.userName = userData.username || '用户'
            this.userEmail = userData.email || ''
            this.userRole = userData.role || 'user'  // 获取用户角色
            console.log('Header: 用户信息获取成功', { id: userData.id, role: this.userRole })
          } else {
            throw new Error('用户数据为空')
          }
        } catch (apiError) {
          console.error('Header: 获取当前用户信息API调用失败:', apiError)
          // 发生错误时，清除用户信息并退出登录
          this.handleLogout();
          alert('获取用户信息失败，请重新登录');
        }
      } catch (error) {
        console.error('获取用户信息失败:', error)
        this.userName = '用户'
        this.userEmail = '获取失败'
        this.userRole = 'user'
      }
    },
    parseJwt(token) {
      try {
        if (!token || typeof token !== 'string') {
          console.error('无效的token格式')
          return {}
        }
        
        const parts = token.split('.')
        if (parts.length !== 3) {
          console.error('JWT格式不正确，应该包含3个部分')
          return {}
        }
        
        const base64Url = parts[1]
        if (!base64Url) {
          console.error('JWT payload部分为空')
          return {}
        }
        
        const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/')
        // 处理base64填充
        const paddedBase64 = base64.padEnd(base64.length + (4 - base64.length % 4) % 4, '=')
        
        try {
          const decoded = atob(paddedBase64)
          const jsonPayload = decodeURIComponent(decoded.split('').map(function(c) {
            return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2)
          }).join(''))
          
          return JSON.parse(jsonPayload)
        } catch (decodeError) {
          console.error('JWT解码失败:', decodeError)
          // 尝试直接解析base64
          try {
            const jsonPayload = atob(paddedBase64)
            return JSON.parse(jsonPayload)
          } catch (finalError) {
            console.error('JWT最终解析失败:', finalError)
            return {}
          }
        }
      } catch (e) {
        console.error('解析JWT失败:', e)
        return {}
      }
    },
    toggleUserMenu(event) {
      event.stopPropagation() // 阻止事件冒泡
      this.showUserMenu = !this.showUserMenu
    },
    closeUserMenu(event) {
      // 检查点击是否在用户菜单外
      if (this.$refs.userMenuRef && !this.$refs.userMenuRef.contains(event.target)) {
        this.showUserMenu = false
      }
    },
    handleLogout() {
      localStorage.removeItem('token')
      this.isLoggedIn = false
      this.userAvatar = 'https://via.placeholder.com/40x40/409eff/white?text=头像'
      this.userName = '用户'
      this.userEmail = '邮箱'
      this.userRole = 'user'
      this.$router.push('/')
      // 刷新页面以确保所有组件状态更新
      window.location.reload()
    }
  }
}
</script>

<style scoped>
.navbar {
  background-color: #fff;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  position: fixed;
  top: 0;
  width: 100%;
  z-index: 1000;
  padding: 15px 0;
}

.nav-container {
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
}

.logo {
  display: flex;
  align-items: center;
}

.logo img {
  height: 40px;
  margin-right: 10px;
}

.logo h1 {
  color: #409eff;
  font-size: 1.8rem;
}

.nav-menu ul {
  display: flex;
  list-style: none;
}

.nav-menu ul li {
  margin-left: 30px;
  position: relative;
}

.nav-menu ul li a {
  text-decoration: none;
  color: #333;
  font-weight: 500;
  transition: color 0.3s;
  padding: 5px 10px;
  display: block;
}

.nav-menu ul li a:hover {
  color: #409eff;
}

/* 下拉菜单样式 */
.dropdown {
  position: relative;
  display: inline-block;
}

.dropdown-content {
  position: absolute;
  background-color: #fff;
  min-width: 140px;
  width: auto; /* 宽度根据内容自适应 */
  box-shadow: 0px 8px 16px rgba(0,0,0,0.1);
  z-index: 1001;
  border-radius: 6px;
  top: 100%;
  left: 0; /* 与触发器左边对齐 */
  opacity: 0;
  visibility: hidden;
  transform: translateY(-6px); /* 仅上移，去除水平平移 */
  transition: opacity 0.18s ease, transform 0.18s ease;
  list-style: none;
  padding: 8px 12px; /* 内边距，但不让每项占满整行 */
  margin: 0;
  box-sizing: border-box;
  display: flex;
  flex-direction: column; /* 强制纵向排列 */
  align-items: flex-start; /* 子项在容器内左对齐 */
  gap: 6px;
}

.dropdown-content li {
  margin: 0;
  padding: 0;
  width: auto; /* 宽度随内容 */
  display: flex;
  justify-content: flex-start; /* 文字靠左 */
}

.dropdown-content a {
  color: #333;
  padding: 8px 12px; /* 悬停背景仅包裹文字 */
  text-decoration: none;
  display: inline-block; /* 背景只跟随文字宽度 */
  transition: background-color 0.18s, color 0.18s;
  box-sizing: border-box;
  text-align: left; /* 文字左对齐 */
  white-space: nowrap; /* 防止中文字符换行 */
}

.dropdown-content li:first-child a {
  border-top-left-radius: 4px;
  border-top-right-radius: 4px;
}

.dropdown-content li:last-child a {
  border-bottom-left-radius: 4px;
  border-bottom-right-radius: 4px;
}

.dropdown-content a:hover {
  background-color: #f5f5f5;
  border-radius: 4px;
}

.dropdown:hover .dropdown-btn,
.dropdown:focus-within .dropdown-btn {
  color: #409eff;
}

.dropdown:hover .dropdown-content,
.dropdown:focus-within .dropdown-content,
.nav-menu ul li:hover > .dropdown-content {
  opacity: 1;
  visibility: visible;
  transform: translateY(0); /* 可见时仅垂直位移到位 */
}

.nav-auth {
  display: flex;
  align-items: center;
}

.login-btn {
  background-color: #409eff;
  color: white !important;
  padding: 8px 16px;
  border-radius: 4px;
  text-decoration: none;
  font-weight: 500;
  transition: background-color 0.3s;
}

.login-btn:hover {
  background-color: #66b1ff;
}

/* 用户信息样式 */
.user-profile {
  display: flex;
  align-items: center;
  cursor: pointer;
  position: relative;
}

.user-info {
  text-align: right;
  margin-right: 10px;
}

.user-email {
  font-size: 12px;
  color: #666;
  margin: 0 0 4px 0;
}

.user-role {
  font-size: 12px;
  color: #409eff;  /* 使用主题色显示角色 */
  margin: 0;
  font-weight: 400;
}

.user-name {
  font-size: 14px;
  color: #333;
  margin: 0;
  font-weight: 500;
}

.avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  overflow: hidden;
  border: 2px solid #409eff;
}

.avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

/* 用户下拉菜单 */
.user-dropdown-menu {
  position: absolute;
  top: 100%;
  right: 0;
  background-color: #fff;
  min-width: 160px;
  box-shadow: 0px 8px 16px rgba(0,0,0,0.1);
  z-index: 1002;
  border-radius: 6px;
  padding: 8px 0;
  margin-top: 8px;
  display: flex;
  flex-direction: column;
}

.user-dropdown-menu a {
  padding: 12px 20px;
  text-decoration: none;
  display: block;
  color: #333;
  font-size: 14px;
}

.user-dropdown-menu a.admin-link {
  color: #e74c3c;
  border-bottom: 1px solid #f0f0f0;
}

.user-dropdown-menu a:hover {
  background-color: #f5f5f5;
}
</style>