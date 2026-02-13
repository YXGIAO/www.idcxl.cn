<template>
  <div class="upstream-product-detail">
    <el-card>
      <template #header>
        <div class="clearfix">
          <span>上游产品详情</span>
          <div style="float: right;">
            <el-button type="primary" size="small" @click="goBack">
              <el-icon><ArrowLeft /></el-icon>
              返回列表
            </el-button>
            <el-button type="primary" size="small" @click="refreshData" :loading="loading">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
          </div>
        </div>
      </template>
      
      <!-- 加载状态 -->
      <div v-if="loading" class="loading-container">
        <el-skeleton :rows="10" animated />
      </div>
      
      <!-- 错误状态 -->
      <div v-else-if="error" class="error-container">
        <el-alert
          :title="error"
          type="error"
          show-icon
          :closable="false"
        />
        <el-button type="primary" @click="fetchProductDetail" style="margin-top: 20px;">
          重试
        </el-button>
      </div>
      
      <!-- 产品详情内容 -->
      <div v-else class="product-detail-content">
        <el-descriptions
          :column="2"
          border
          title="产品基本信息"
          size="large"
        >
          <el-descriptions-item label="产品ID">{{ product.id }}</el-descriptions-item>
          <el-descriptions-item label="上游产品ID">{{ product.upstream_host_id || '-' }}</el-descriptions-item>
          <el-descriptions-item label="产品标识">{{ product.name || '-' }}</el-descriptions-item>
          <el-descriptions-item label="商品名称">{{ product.product_name || '-' }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="getStatusType(product.status)" size="small">
              {{ getStatusText(product.status) }}
            </el-tag>
          </el-descriptions-item>
        </el-descriptions>
        
        <el-descriptions
          :column="2"
          border
          title="财务信息"
          size="large"
          style="margin-top: 20px;"
        >
          <el-descriptions-item label="首次付款金额">
            <span class="amount">{{ formatAmount(product.first_payment_amount) }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="续费金额">
            <span class="amount">{{ formatAmount(product.renew_amount) }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="付款周期">{{ product.billing_cycle || '-' }}</el-descriptions-item>
          <el-descriptions-item label="计费周期名称">{{ product.billing_cycle_name || '-' }}</el-descriptions-item>
          <el-descriptions-item label="计费周期时间">{{ product.billing_cycle_time || '-' }}</el-descriptions-item>
        </el-descriptions>
        
        <el-descriptions
          :column="2"
          border
          title="时间信息"
          size="large"
          style="margin-top: 20px;"
        >
          <el-descriptions-item label="开通时间">
            {{ formatTimestamp(product.active_time) }}
          </el-descriptions-item>
          <el-descriptions-item label="到期时间">
            {{ formatTimestamp(product.due_time) }}
          </el-descriptions-item>
          <el-descriptions-item label="创建时间">
            {{ formatTime(product.created_at) }}
          </el-descriptions-item>
          <el-descriptions-item label="更新时间">
            {{ formatTime(product.updated_at) }}
          </el-descriptions-item>
        </el-descriptions>
        
        <el-descriptions
          :column="2"
          border
          title="用户信息"
          size="large"
          style="margin-top: 20px;"
        >
          <el-descriptions-item label="用户ID">{{ product.client_id || '-' }}</el-descriptions-item>
          <el-descriptions-item label="用户名">{{ product.username || '-' }}</el-descriptions-item>
          <el-descriptions-item label="公司名称">{{ product.company || '-' }}</el-descriptions-item>
          <el-descriptions-item label="邮箱">{{ product.email || '-' }}</el-descriptions-item>
          <el-descriptions-item label="手机号">
            {{ product.phone_code || '+86' }}{{ product.phone || '-' }}
          </el-descriptions-item>
        </el-descriptions>
        
        <el-descriptions
          :column="1"
          border
          title="配置信息"
          size="large"
          style="margin-top: 20px;"
        >
          <el-descriptions-item label="IP数量">{{ product.ip_num || 0 }}</el-descriptions-item>
          <el-descriptions-item label="基础信息">
            <div class="base-info">{{ product.base_info || '-' }}</div>
          </el-descriptions-item>
        </el-descriptions>
        
        <el-descriptions
          :column="2"
          border
          title="供应商信息"
          size="large"
          style="margin-top: 20px;"
          v-if="product.supplier"
        >
          <el-descriptions-item label="供应商ID">{{ product.supplier_id || '-' }}</el-descriptions-item>
          <el-descriptions-item label="供应商名称">{{ product.supplier.name || '-' }}</el-descriptions-item>
        </el-descriptions>
      </div>
    </el-card>
  </div>
</template>

<script>
import { upstreamProductAPI } from '@/services/api';
import { ElMessage } from 'element-plus';
import { ArrowLeft, Refresh } from '@element-plus/icons-vue';
import { useRouter } from 'vue-router';

export default {
  name: 'UpstreamProductDetail',
  components: {
    ArrowLeft,
    Refresh,
  },
  data() {
    return {
      loading: false,
      error: '',
      product: {},
      productId: null,
    };
  },
  created() {
    // 从路由参数获取产品ID
    this.productId = this.$route.params.id;
    if (this.productId) {
      this.fetchProductDetail();
    } else {
      this.error = '产品ID不能为空';
    }
  },
  methods: {
    // 获取产品详情
    async fetchProductDetail() {
      if (!this.productId) {
        this.error = '产品ID不能为空';
        return;
      }
      
      this.loading = true;
      this.error = '';
      
      try {
        const response = await upstreamProductAPI.getUpstreamProduct(this.productId);
        
        // 处理响应数据
        if (response && response.host) {
          this.product = response.host;
        } else {
          this.error = '获取产品详情失败：数据格式错误';
          ElMessage.error('获取产品详情失败：数据格式错误');
        }
      } catch (error) {
        console.error('获取产品详情失败:', error);
        this.error = '获取产品详情失败：' + (error.message || '网络错误');
        ElMessage.error('获取产品详情失败：' + (error.message || '网络错误'));
      } finally {
        this.loading = false;
      }
    },
    
    // 刷新数据
    refreshData() {
      this.fetchProductDetail();
    },
    
    // 返回列表页
    goBack() {
      this.$router.push('/admin/upstream-products');
    },
    
    // 获取状态类型
    getStatusType(status) {
      switch (status) {
        case 'Active': return 'success';
        case 'Pending': return 'warning';
        case 'Unpaid': return 'danger';
        case 'Suspended': return 'info';
        case 'Deleted': return 'info';
        case 'Failed': return 'danger';
        default: return 'info';
      }
    },
    
    // 获取状态文本
    getStatusText(status) {
      switch (status) {
        case 'Active': return '已开通';
        case 'Pending': return '开通中';
        case 'Unpaid': return '未付款';
        case 'Suspended': return '已暂停';
        case 'Deleted': return '已删除';
        case 'Failed': return '开通失败';
        default: return status;
      }
    },
    
    // 格式化金额
    formatAmount(amount) {
      if (!amount) return '-';
      // 如果已经是字符串格式，直接返回
      if (typeof amount === 'string') {
        return amount;
      }
      // 如果是数字，格式化为两位小数
      if (typeof amount === 'number') {
        return '¥' + amount.toFixed(2);
      }
      return amount;
    },
    
    // 格式化时间戳（Unix时间戳，单位秒）
    formatTimestamp(timestamp) {
      if (!timestamp) return '-';
      // 将秒转换为毫秒
      const date = new Date(timestamp * 1000);
      return date.toLocaleString();
    },
    
    // 格式化时间（ISO字符串）
    formatTime(timeString) {
      if (!timeString) return '-';
      try {
        const date = new Date(timeString);
        return date.toLocaleString();
      } catch (error) {
        return timeString;
      }
    },
  },
};
</script>

<style scoped>
.upstream-product-detail {
  padding: 20px;
}

.loading-container {
  padding: 40px 0;
  text-align: center;
}

.error-container {
  padding: 40px 0;
  text-align: center;
}

.product-detail-content {
  padding: 10px 0;
}

.amount {
  font-weight: bold;
  color: #e6a23c;
}

.base-info {
  white-space: pre-wrap;
  word-break: break-word;
  line-height: 1.6;
  padding: 8px;
  background-color: #f5f7fa;
  border-radius: 4px;
  border: 1px solid #ebeef5;
}

.el-descriptions {
  margin-bottom: 20px;
}

.el-tag {
  margin-right: 5px;
}
</style>