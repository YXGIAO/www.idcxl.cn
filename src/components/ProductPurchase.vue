<template>
  <div class="product-purchase-container">
    <el-row :gutter="20">
      <!-- 产品基本信息 -->
      <el-col :span="16">
        <el-card>
          <template #header>
            <div class="product-header">
              <h2>{{ product.name }}</h2>
              <el-tag :type="product.status === 'active' ? 'success' : 'info'">
                {{ product.status === 'active' ? '可购买' : '不可用' }}
              </el-tag>
            </div>
          </template>

          <div class="product-description">
            <p>{{ product.description }}</p>
          </div>

          <div class="product-specs">
            <el-descriptions :column="2" border>
              <el-descriptions-item label="基础价格">¥{{ product.price.toFixed(2) }}/月</el-descriptions-item>
              <el-descriptions-item label="CPU">最大{{ product.max_cpu }}核</el-descriptions-item>
              <el-descriptions-item label="内存">最大{{ product.max_memory }}GB</el-descriptions-item>
              <el-descriptions-item label="存储">最大{{ product.max_storage }}GB</el-descriptions-item>
              <el-descriptions-item label="带宽">最大{{ product.max_bandwidth }}Mbps</el-descriptions-item>
            </el-descriptions>
          </div>
        </el-card>
      </el-col>

      <!-- 配置选择表格 -->
      <el-col :span="8">
        <el-card>
          <template #header>
            <h3>配置选项</h3>
          </template>

          <el-form :model="orderForm" label-position="top">
            <el-form-item label="CPU (核)">
              <el-input-number v-model="orderForm.cpu" :min="1" :max="product.max_cpu" @change="calculatePrice"/>
            </el-form-item>

            <el-form-item label="内存 (GB)">
              <el-input-number v-model="orderForm.memory" :min="1" :max="product.max_memory" @change="calculatePrice"/>
            </el-form-item>

            <el-form-item label="存储 (GB)">
              <el-input-number v-model="orderForm.storage" :min="20" :max="product.max_storage" @change="calculatePrice"/>
            </el-form-item>

            <el-form-item label="带宽 (Mbps)">
              <el-input-number v-model="orderForm.bandwidth" :min="1" :max="product.max_bandwidth" @change="calculatePrice"/>
            </el-form-item>

            <el-form-item label="购买时长 (月)">
              <el-select v-model="orderForm.period" @change="calculatePrice">
                <el-option label="1个月" :value="1"></el-option>
                <el-option label="3个月" :value="3"></el-option>
                <el-option label="6个月" :value="6"></el-option>
                <el-option label="12个月" :value="12"></el-option>
              </el-select>
            </el-form-item>

            <el-form-item label="支付方式">
              <el-select v-model="orderForm.paymentMethod">
                <el-option label="支付宝" value="alipay"></el-option>
                <el-option label="微信支付" value="wechat"></el-option>
                <el-option label="银行转账" value="bank"></el-option>
              </el-select>
            </el-form-item>

            <el-divider/>

            <div class="price-summary">
              <h3>价格详情</h3>
              <div class="price-details">
                <p>基础价格: ¥{{ product.price.toFixed(2) }}</p>
                <p>配置加价: ¥{{ calculatedPrice.configPrice.toFixed(2) }}</p>
                <p>总计: ¥{{ calculatedPrice.totalPrice.toFixed(2) }}</p>
                <p>时长: {{ orderForm.period }}个月</p>
                <p class="total">总价: ¥{{ (calculatedPrice.totalPrice * orderForm.period).toFixed(2) }}</p>
              </div>
            </div>

            <el-button type="primary" size="large" @click="submitOrder" :loading="submitting">
              立即购买
            </el-button>
          </el-form>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import { productAPI, orderAPI } from '@/services/api';

export default {
  name: 'ProductPurchase',
  data() {
    return {
      product: {
        id: null,
        name: '',
        description: '',
        price: 0,
        max_cpu: 16,
        max_memory: 64,
        max_storage: 1000,
        max_bandwidth: 100,
        status: 'active'
      },
      orderForm: {
        cpu: 2,
        memory: 4,
        storage: 100,
        bandwidth: 5,
        period: 1,
        paymentMethod: 'alipay',
        notes: ''
      },
      calculatedPrice: {
        basePrice: 0,
        configPrice: 0,
        totalPrice: 0
      },
      submitting: false
    };
  },
  created() {
    this.loadProduct();
  },
  methods: {
    async loadProduct() {
      try {
        const productId = this.$route.params.id;
        const response = await productAPI.getProduct(productId);
        this.product = response.data || response;
        this.calculatedPrice.basePrice = this.product.price;
        this.calculatePrice();
      } catch (error) {
        this.$message.error('获取产品信息失败');
        console.error(error);
      }
    },
    calculatePrice() {
      // 简单价格计算逻辑
      const cpuPrice = this.orderForm.cpu * 100;
      const memoryPrice = this.orderForm.memory * 50;
      const storagePrice = this.orderForm.storage * 0.5;
      const bandwidthPrice = this.orderForm.bandwidth * 10;
      
      this.calculatedPrice.configPrice = cpuPrice + memoryPrice + storagePrice + bandwidthPrice;
      this.calculatedPrice.totalPrice = this.product.price + this.calculatedPrice.configPrice;
    },
    async submitOrder() {
      this.submitting = true;
      try {
        const orderData = {
          product_id: this.product.id,
          cpu: this.orderForm.cpu,
          memory: this.orderForm.memory,
          storage: this.orderForm.storage,
          bandwidth: this.orderForm.bandwidth,
          period: this.orderForm.period,
          payment_method: this.orderForm.paymentMethod,
          notes: this.orderForm.notes,
          amount: this.calculatedPrice.totalPrice * this.orderForm.period
        };
        
        await orderAPI.createProductOrder(orderData);
        this.$message.success('订单创建成功');
        this.$router.push('/user/orders');
      } catch (error) {
        this.$message.error('下单失败: ' + (error.message || '网络错误'));
        console.error(error);
      } finally {
        this.submitting = false;
      }
    }
  }
};
</script>

<style scoped>
.product-purchase-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 2rem;
}

.product-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.product-description {
  margin-bottom: 2rem;
  line-height: 1.6;
}

.price-summary {
  margin: 1.5rem 0;
}

.price-details {
  margin: 1rem 0;
}

.price-details p {
  margin: 0.5rem 0;
}

.price-details .total {
  font-size: 1.2rem;
  font-weight: bold;
  margin-top: 1rem;
  border-top: 1px solid #eee;
  padding-top: 0.5rem;
}

.el-button {
  width: 100%;
  margin-top: 1rem;
}
</style>
