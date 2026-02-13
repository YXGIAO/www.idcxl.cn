<template>
  <section id="products" class="products">
    <div class="container">
      <h2 class="section-title">我们的产品</h2>
      <p class="section-subtitle">全方位的产品布局，满足不同用户的多样化需求</p>
      
      <div class="products-grid">
        <div class="product-card" v-for="product in products" :key="product.id">
          <div class="product-icon">
            <img :src="product.image || 'https://via.placeholder.com/80x80/007bff/white?text='+product.name.substring(0,2)" :alt="product.name">
          </div>
          <h3>{{ product.name }}</h3>
          <p class="price">¥{{ product.price.toFixed(2) }}/月</p>
          <p>{{ product.description }}</p>
          <el-button type="primary" size="small" @click="$router.push('/product/' + product.id)">
            立即购买
          </el-button>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
import { productAPI } from '@/services/api';

export default {
  name: 'ProductsSection',
  data() {
    return {
      products: []
    };
  },
  async created() {
    try {
      const response = await productAPI.getProducts({});
      this.products = response.products || response;
    } catch (error) {
      console.error('获取产品列表失败:', error);
    }
  }
};
</script>

<style scoped>
.products {
  padding: 5rem 0;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 2rem;
}

.section-title {
  font-size: 2rem;
  text-align: center;
  margin-bottom: 1rem;
  color: #333;
}

.section-subtitle {
  text-align: center;
  color: #666;
  margin-bottom: 3rem;
  font-size: 1.1rem;
}

.products-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 2rem;
  margin-top: 2rem;
}

.product-card {
  background: rgba(255, 255, 255, 0.7);  /* 半透明背景 */
  backdrop-filter: blur(10px);  /* 毛玻璃效果 */
  border-radius: 8px;
  padding: 2rem;
  text-align: center;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.05);
  transition: transform 0.3s, box-shadow 0.3s;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.product-card:hover {
  transform: translateY(-10px);
  box-shadow: 0 15px 30px rgba(0, 0, 0, 0.1);
}

.product-icon {
  margin-bottom: 1.5rem;
}

.product-icon img {
  width: 80px;
  height: 80px;
  object-fit: contain;
}

.product-card h3 {
  font-size: 1.4rem;
  margin-bottom: 1rem;
  color: #333;
}

.product-card p {
  color: #666;
  margin-bottom: 1.5rem;
  line-height: 1.6;
}

.learn-more {
  color: #007bff;
  text-decoration: none;
  font-weight: 500;
  display: inline-block;
  transition: color 0.3s;
}

.learn-more:hover {
  color: #0056b3;
}
</style>