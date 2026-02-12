<template>
  <section id="pricing" class="pricing">
    <div class="container">
      <h2 class="section-title">价格方案</h2>
      <p class="section-subtitle">多种配置选择，满足不同需求</p>
      
      <div class="pricing-cards">
        <div class="pricing-card">
          <div class="pricing-header">
            <h3>虚拟主机</h3>
            <div class="price">
              <span class="amount">¥8</span>
              <span class="period">/月起</span>
            </div>
          </div>
          
          <ul class="features">
            <li>适合个人网站、博客</li>
            <li>轻量化配置</li>
            <li>一键安装多种应用</li>
            <li>支持WordPress, Typecho等</li>
            <li>基础安全防护</li>
          </ul>
          
          <button class="select-btn">立即选购</button>
        </div>
        
        <div class="pricing-card featured">
          <div class="popular-tag">最受欢迎</div>
          <div class="pricing-header">
            <h3>云服务器</h3>
            <div class="price">
              <span class="amount">¥15</span>
              <span class="period">/月起</span>
            </div>
          </div>
          
          <ul class="features">
            <li>独立操作系统和硬件配置</li>
            <li>可搭建网站、DNS服务器等</li>
            <li>支持1Panel/宝塔面板</li>
            <li>国内、香港、海外多线路选择</li>
            <li>高防安全防护</li>
          </ul>
          
          <button class="select-btn">立即选购</button>
        </div>
        
        <div class="pricing-card">
          <div class="pricing-header">
            <h3>游戏云服务器</h3>
            <div class="price">
              <span class="amount">¥25</span>
              <span class="period">/月起</span>
            </div>
          </div>
          
          <ul class="features">
            <li>专为游戏优化</li>
            <li>低延迟网络环境</li>
            <li>预装各类游戏服务端</li>
            <li>支持Minecraft等热门游戏</li>
            <li>弹性资源分配</li>
          </ul>
          
          <button class="select-btn">立即选购</button>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
export default {
  name: 'PricingSection'
}
</script>

<style scoped>
.pricing {
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

.pricing-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 2rem;
}

.pricing-card {
  background: rgba(255, 255, 255, 0.7);  /* 半透明背景 */
  backdrop-filter: blur(10px);  /* 毛玻璃效果 */
  border-radius: 8px;
  padding: 2rem;
  text-align: center;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.05);
  position: relative;
  border: 1px solid rgba(255, 255, 255, 0.2);
  transition: transform 0.3s, box-shadow 0.3s;
}

.pricing-card:hover {
  transform: translateY(-10px);
  box-shadow: 0 15px 30px rgba(0, 0, 0, 0.1);
}

.pricing-card.featured {
  border: 2px solid #007bff;
  transform: scale(1.05);
}

.pricing-card.featured:hover {
  transform: scale(1.05) translateY(-10px);
}

.popular-tag {
  position: absolute;
  top: -12px;
  left: 50%;
  transform: translateX(-50%);
  background-color: #ffc107;
  color: #333;
  padding: 0.3rem 1rem;
  border-radius: 20px;
  font-size: 0.8rem;
  font-weight: bold;
}

.pricing-header h3 {
  font-size: 1.5rem;
  margin-bottom: 1rem;
  color: #333;
}

.price {
  margin-bottom: 1.5rem;
}

.amount {
  font-size: 2.5rem;
  font-weight: bold;
  color: #007bff;
}

.period {
  color: #666;
  font-size: 0.9rem;
}

.features {
  list-style: none;
  margin: 1.5rem 0;
  text-align: left;
}

.features li {
  padding: 0.5rem 0;
  border-bottom: 1px solid #f0f0f0;
}

.features li:last-child {
  border-bottom: none;
}

.select-btn {
  background-color: #007bff;
  color: white;
  padding: 0.8rem 1.5rem;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  cursor: pointer;
  width: 100%;
  transition: background-color 0.3s;
}

.select-btn:hover {
  background-color: #0056b3;
}

@media (max-width: 768px) {
  .pricing-card.featured {
    transform: scale(1);
  }
  
  .pricing-card.featured:hover {
    transform: scale(1) translateY(-10px);
  }
}
</style>