<template>
  <div class="admin-dashboard-overview">
    <!-- 统计卡片区域 -->
    <div class="dashboard-stats">
      <el-row :gutter="20">
        <el-col :span="6">
          <el-card class="stat-card" shadow="hover">
            <div class="stat-content">
              <div class="stat-icon">
                <User />
              </div>
              <div class="stat-info">
                <div class="stat-number">1,248</div>
                <div class="stat-label">总用户数</div>
              </div>
            </div>
          </el-card>
        </el-col>
        
        <el-col :span="6">
          <el-card class="stat-card" shadow="hover">
            <div class="stat-content">
              <div class="stat-icon">
                <ShoppingCart />
              </div>
              <div class="stat-info">
                <div class="stat-number">3,562</div>
                <div class="stat-label">订单总数</div>
              </div>
            </div>
          </el-card>
        </el-col>
        
        <el-col :span="6">
          <el-card class="stat-card" shadow="hover">
            <div class="stat-content">
              <div class="stat-icon">
                <Money />
              </div>
              <div class="stat-info">
                <div class="stat-number">¥24,890</div>
                <div class="stat-label">本月收入</div>
              </div>
            </div>
          </el-card>
        </el-col>
        
        <el-col :span="6">
          <el-card class="stat-card" shadow="hover">
            <div class="stat-content">
              <div class="stat-icon">
                <Document />
              </div>
              <div class="stat-info">
                <div class="stat-number">128</div>
                <div class="stat-label">待处理工单</div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 主要内容区域 -->
    <el-row :gutter="20" class="dashboard-main">
      <!-- 数据图表 -->
      <el-col :span="16">
        <el-card class="chart-card" shadow="hover">
          <template #header>
            <span class="card-header">业务概览</span>
          </template>
          <div class="chart-placeholder">
            <div style="height: 300px; display: flex; align-items: center; justify-content: center; color: #888;">
              收入趋势图表（可集成ECharts）
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 系统状态 -->
      <el-col :span="8">
        <el-card class="system-card" shadow="hover">
          <template #header>
            <span class="card-header">系统状态</span>
          </template>
          <div class="system-status">
            <div class="status-item">
              <div class="status-icon" :style="{ backgroundColor: '#67C23A' }"></div>
              <div class="status-info">
                <div class="status-title">数据库</div>
                <div class="status-value">正常</div>
              </div>
            </div>
            <div class="status-item">
              <div class="status-icon" :style="{ backgroundColor: '#67C23A' }"></div>
              <div class="status-info">
                <div class="status-title">API服务</div>
                <div class="status-value">正常</div>
              </div>
            </div>
            <div class="status-item">
              <div class="status-icon" :style="{ backgroundColor: '#E6A23C' }"></div>
              <div class="status-info">
                <div class="status-title">存储空间</div>
                <div class="status-value">78% 使用</div>
              </div>
            </div>
            <div class="status-item">
              <div class="status-icon" :style="{ backgroundColor: '#67C23A' }"></div>
              <div class="status-info">
                <div class="status-title">网络服务</div>
                <div class="status-value">正常</div>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 最近活动 -->
      <el-col :span="16">
        <el-card class="activity-card" shadow="hover">
          <template #header>
            <span class="card-header">最近活动</span>
          </template>
          <div class="activity-list">
            <div class="activity-item">
              <div class="activity-icon">
                <User />
              </div>
              <div class="activity-content">
                <div class="activity-title">新用户注册</div>
                <div class="activity-desc">杨小iao 注册了账号</div>
                <div class="activity-time">2分钟前</div>
              </div>
            </div>
            <div class="activity-item">
              <div class="activity-icon">
                <ShoppingCart />
              </div>
              <div class="activity-content">
                <div class="activity-title">新订单创建</div>
                <div class="activity-desc">订单 #ORD20240124001 已创建</div>
                <div class="activity-time">15分钟前</div>
              </div>
            </div>
            <div class="activity-item">
              <div class="activity-icon">
                <Document />
              </div>
              <div class="activity-content">
                <div class="activity-title">工单更新</div>
                <div class="activity-desc">工单 #TCK20240124001 状态已更新为"处理中"</div>
                <div class="activity-time">30分钟前</div>
              </div>
            </div>
            <div class="activity-item">
              <div class="activity-icon">
                <Setting />
              </div>
              <div class="activity-content">
                <div class="activity-title">系统配置更新</div>
                <div class="activity-desc">服务器配置已更新</div>
                <div class="activity-time">1小时前</div>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 快捷操作 -->
      <el-col :span="8">
        <el-card class="quick-actions-card" shadow="hover">
          <template #header>
            <span class="card-header">快捷操作</span>
          </template>
          <div class="quick-actions">
            <el-button type="primary" size="large" @click="goTo('/admin/users')">添加用户</el-button>
            <el-button type="success" size="large" @click="goTo('/admin/orders')">查看订单</el-button>
            <el-button type="info" size="large" @click="goTo('/admin/settings')">系统设置</el-button>
            <el-button type="warning" size="large" @click="goTo('/admin/tickets')">处理工单</el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import { User, ShoppingCart, Money, Document, Setting } from '@element-plus/icons-vue'

export default {
  name: 'DashboardOverview',
  components: {
    User,
    ShoppingCart,
    Money,
    Document,
    Setting
  },
  methods: {
    goTo(path) {
      this.$router.push(path);
    }
  }
}
</script>

<style scoped>
.admin-dashboard-overview {
  padding: 20px;
  background-color: #f5f7fa;
}

.dashboard-stats {
  margin-bottom: 20px;
}

.stat-card {
  height: 100px;
}

.stat-content {
  display: flex;
  align-items: center;
  height: 100%;
}

.stat-icon {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  background-color: #e6f7ff;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #1890ff;
  margin-right: 16px;
}

.stat-info {
  flex: 1;
}

.stat-number {
  font-size: 24px;
  font-weight: 600;
  color: #333;
}

.stat-label {
  font-size: 14px;
  color: #999;
}

.dashboard-main {
  margin-bottom: 20px;
}

.card-header {
  font-weight: 600;
  color: #333;
}

.activity-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.activity-item {
  display: flex;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid #eee;
}

.activity-item:last-child {
  border-bottom: none;
}

.activity-icon {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  background-color: #f0f9ff;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #409eff;
  margin-right: 12px;
  flex-shrink: 0;
}

.activity-content {
  flex: 1;
}

.activity-title {
  font-weight: 500;
  color: #333;
  margin-bottom: 4px;
}

.activity-desc {
  font-size: 13px;
  color: #666;
  margin-bottom: 4px;
}

.activity-time {
  font-size: 12px;
  color: #999;
}

.system-status {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.status-item {
  display: flex;
  align-items: center;
}

.status-icon {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  margin-right: 12px;
  flex-shrink: 0;
}

.status-info {
  flex: 1;
}

.status-title {
  font-size: 14px;
  color: #666;
}

.status-value {
  font-size: 13px;
  color: #333;
  font-weight: 500;
}

.quick-actions {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.chart-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 260px;
  color: #888;
}
</style>