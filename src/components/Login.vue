<template>
  <div class="login-container">
    <div class="login-wrapper">
      <div class="login-header">
        <h1>登录</h1>
        <p>还没有账号？<a href="/register">立即注册</a></p>
      </div>
      
      <div class="login-form">
        <div class="input-group">
          <input 
            type="text" 
            v-model="username" 
            placeholder="用户名/手机号/邮箱" 
            class="input-field"
          >
        </div>
        
        <div class="input-group">
          <input 
            type="password" 
            v-model="password" 
            placeholder="密码" 
            class="input-field"
          >
        </div>
        
        <div class="remember-group">
          <label class="remember-checkbox">
            <input type="checkbox" v-model="rememberMe">
            <span class="checkmark"></span>
            记住我
          </label>
          <a href="#" class="forgot-password">忘记密码？</a>
        </div>
        
        <button class="login-btn" :disabled="loading" @click="handleLogin">
          <span v-if="!loading">登录</span>
          <span v-else>登录中...</span>
        </button>
        <div class="error-message" v-if="errorMessage">{{ errorMessage }}</div>
        
        <div class="other-login">
          <span class="other-login-text">其他登录方式</span>
          <div class="other-login-icons">
            <a href="#" class="login-icon qq-icon">QQ</a>
            <a href="#" class="login-icon wechat-icon">微信</a>
            <a href="#" class="login-icon microsoft-icon">微软</a>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { authAPI, userAPI, currentUserAPI } from '../services/api'

export default {
  name: 'Login',
  data() {
    return {
      username: '',
      password: '',
      rememberMe: false,
      loading: false,
      errorMessage: ''
    }
  },
  methods: {
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
    
    async handleLogin() {
      this.errorMessage = ''
      if (!this.username || !this.password) {
        this.errorMessage = '请输入用户名/邮箱/手机号和密码'
        return
      }

      this.loading = true
      try {
        const res = await authAPI.login({ username: this.username, password: this.password })
        // 成功返回应包含 token
        const token = res && res.token ? res.token : null
        if (token) {
          localStorage.setItem('token', token)
          
          // 获取用户信息并存储角色
          let userRole = 'user' // 默认角色
          
          try {
            // 使用新的API获取当前用户信息
            const userData = await currentUserAPI.getCurrentUser()
            console.log('Current User Data:', userData) // 调试信息
            
            // 从响应中获取用户角色
            const user = userData.data || userData
            if (user && user.role) {
              userRole = user.role
            } else {
              // 如果无法获取用户角色，抛出错误
              throw new Error('无法获取用户角色信息')
            }
          } catch (err) {
            console.error('获取用户角色失败:', err)
            // 如果无法获取用户角色，不应继续登录
            this.errorMessage = '获取用户角色失败，请重试或联系管理员';
            return; // 停止登录流程
          }
          
          // 设置用户角色
          localStorage.setItem('userRole', userRole)
          console.log('Stored userRole:', userRole) // 调试信息
          
          // 跳转到首页或路由 query 中的 redirect
          const redirect = this.$route.query.redirect || '/'
          this.$router.push(redirect)
        } else {
          this.errorMessage = '登录失败：未收到令牌'
        }
      } catch (err) {
        console.error('Login error:', err)
        if (err && err.message === 'Network Error') {
          this.errorMessage = '网络错误：无法连接到后端，请确认后端服务在 http://localhost:8081 正在运行'
        } else if (err && err.response && err.response.data) {
          const data = err.response.data
          this.errorMessage = data.error || data.message || data.msg || JSON.stringify(data)
        } else if (err && err.message) {
          this.errorMessage = err.message
        } else {
          this.errorMessage = '登录失败，请重试'
        }
      } finally {
        this.loading = false
      }
    }
  }
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 20px;
}

.login-wrapper {
  width: 100%;
  max-width: 400px;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 8px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
  padding: 40px;
  box-sizing: border-box;
}

.login-header {
  text-align: center;
  margin-bottom: 30px;
}

.login-header h1 {
  font-size: 24px;
  color: #222;
  margin-bottom: 10px;
  font-weight: 600;
}

.login-header p {
  color: #999;
  font-size: 14px;
}

.login-header a {
  color: #409eff;
  text-decoration: none;
}

.input-group {
  margin-bottom: 20px;
}

.input-field {
  width: 100%;
  padding: 14px 16px;
  border: 1px solid #e5e9ef;
  border-radius: 4px;
  font-size: 16px;
  box-sizing: border-box;
  transition: border-color 0.2s;
  background: rgba(255, 255, 255, 0.8);
}

.input-field:focus {
  outline: none;
  border-color: #409eff;
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.1);
}

.remember-group {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  font-size: 14px;
}

.remember-checkbox {
  display: flex;
  align-items: center;
  cursor: pointer;
  user-select: none;
  color: #666;
}

.remember-checkbox input {
  display: none;
}

.checkmark {
  display: inline-block;
  width: 16px;
  height: 16px;
  border: 1px solid #ccd0d7;
  border-radius: 2px;
  margin-right: 6px;
  position: relative;
  transition: all 0.2s;
}

.remember-checkbox input:checked + .checkmark {
  background-color: #409eff;
  border-color: #409eff;
}

.remember-checkbox input:checked + .checkmark::after {
  content: '';
  position: absolute;
  left: 4px;
  top: 1px;
  width: 4px;
  height: 8px;
  border: solid white;
  border-width: 0 2px 2px 0;
  transform: rotate(45deg);
}

.forgot-password {
  color: #409eff;
  text-decoration: none;
}

.login-btn {
  width: 100%;
  padding: 14px;
  background-color: #409eff;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 16px;
  cursor: pointer;
  transition: background-color 0.2s;
  margin-bottom: 20px;
}

.login-btn:hover {
  background-color: #66b1ff;
}

.error-message {
  color: #e74c3c;
  text-align: center;
  margin-top: 12px;
}

.other-login {
  text-align: center;
  border-top: 1px solid #e5e9ef;
  padding-top: 24px;
}

.other-login-text {
  display: inline-block;
  position: relative;
  top: -10px;
  background: rgba(255, 255, 255, 0.9);
  color: #999;
  font-size: 14px;
  padding: 0 10px;
}

.other-login-icons {
  display: flex;
  justify-content: center;
  gap: 20px;
  margin-top: 10px;
}

.login-icon {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  color: white;
  font-weight: bold;
  text-decoration: none;
  transition: all 0.2s;
  font-size: 12px;
}

.qq-icon {
  background: #409eff;
}

.wechat-icon {
  background: #00c800;
}

.microsoft-icon {
  background: #444;
}

.qq-icon:hover,
.wechat-icon:hover,
.microsoft-icon:hover {
  opacity: 0.8;
  transform: scale(1.1);
}

@media (max-width: 768px) {
  .login-wrapper {
    padding: 30px 20px;
  }
}
</style>