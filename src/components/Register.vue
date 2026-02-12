<template>
  <div class="register-container">
    <div class="register-wrapper">
      <div class="register-header">
        <h1>注册</h1>
        <p>已有账号？<a href="/login">去登录</a></p>
      </div>

      <div class="register-form">
        <div class="input-group">
          <input
            type="text"
            v-model="username"
            placeholder="用户名"
            class="input-field"
          />
        </div>

        <div class="input-group">
          <input
            type="email"
            v-model="email"
            placeholder="邮箱"
            class="input-field"
          />
        </div>

        <div class="input-group">
          <input
            type="password"
            v-model="password"
            placeholder="密码（至少6位）"
            class="input-field"
          />
        </div>

        <div class="input-group">
          <input
            type="password"
            v-model="confirmPassword"
            placeholder="确认密码"
            class="input-field"
          />
        </div>

        <button class="register-btn" @click="handleRegister">注册</button>

        <div class="note" v-if="errorMessage">{{ errorMessage }}</div>
        <div class="success" v-if="successMessage">{{ successMessage }}</div>
      </div>
    </div>
  </div>
</template>

<script>
import { authAPI } from '../services/api'

export default {
  name: 'Register',
  data() {
    return {
      username: '',
      email: '',
      password: '',
      confirmPassword: '',
      errorMessage: '',
      successMessage: ''
    }
  },
  methods: {
    async handleRegister() {
      this.errorMessage = ''
      this.successMessage = ''

      if (!this.username || !this.email || !this.password || !this.confirmPassword) {
        this.errorMessage = '请填写所有字段'
        return
      }

      if (this.password.length < 6) {
        this.errorMessage = '密码长度至少6位'
        return
      }

      if (this.password !== this.confirmPassword) {
        this.errorMessage = '两次输入的密码不一致'
        return
      }

      try {
        const payload = {
          username: this.username,
          email: this.email,
          password: this.password
        }
        await authAPI.register(payload)
        this.successMessage = '注册成功，正在跳转到登录页...'
        setTimeout(() => {
          this.$router.push({ path: '/login' })
        }, 1000)
      } catch (err) {
        // axios interceptor returns error; try to extract message
        if (err && err.response && err.response.data && err.response.data.message) {
          this.errorMessage = err.response.data.message
        } else if (err && err.response && err.response.status === 409) {
          // 处理邮箱或用户名已被注册的错误
          const errorMsg = err.response.data.error || '邮箱或用户名已被注册'
          if (errorMsg.includes('邮箱')) {
            this.errorMessage = '邮箱已被注册'
          } else if (errorMsg.includes('用户名')) {
            this.errorMessage = '用户名已被注册'
          } else {
            this.errorMessage = errorMsg
          }
        } else if (err && err.message) {
          this.errorMessage = err.message
        } else {
          this.errorMessage = '注册失败，请重试'
        }
      }
    }
  }
}
</script>

<style scoped>
.register-container {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 20px;
}

.register-wrapper {
  width: 100%;
  max-width: 420px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 8px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.08);
  padding: 36px;
  box-sizing: border-box;
}

.register-header {
  text-align: center;
  margin-bottom: 20px;
}

.register-header h1 {
  font-size: 22px;
  color: #222;
  margin-bottom: 6px;
  font-weight: 600;
}

.register-header a {
  color: #409eff;
  text-decoration: none;
}

.input-group { margin-bottom: 16px; }

.input-field {
  width: 100%;
  padding: 12px 14px;
  border: 1px solid #e5e9ef;
  border-radius: 4px;
  font-size: 15px;
  box-sizing: border-box;
}

.input-field:focus { outline: none; border-color: #409eff; }

.register-btn {
  width: 100%;
  padding: 12px;
  background-color: #409eff;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 15px;
  cursor: pointer;
}

.note { color: #e74c3c; margin-top: 12px; text-align: center; }
.success { color: #2ecc71; margin-top: 12px; text-align: center; }

@media (max-width: 768px) {
  .register-wrapper { padding: 24px; }
}
</style>
