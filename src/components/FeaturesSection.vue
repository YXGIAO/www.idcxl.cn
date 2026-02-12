<template>
  <section id="features" class="features">
    <div class="container">
      <h2 class="section-title">为什么选择仙林云</h2>
      <p class="section-subtitle">我们致力于为您提供稳定、高效、便捷的云服务体验</p>
      
      <div class="features-grid">
        <div class="feature-card">
          <div class="feature-icon">
            <span>⚡</span>
          </div>
          <h3>高速稳定</h3>
          <p>采用先进的虚拟化技术和优质硬件，确保服务的高性能和高稳定性，让您和您的用户享受流畅体验。</p>
        </div>
        
        <div class="feature-card">
          <div class="feature-icon">
            <span>💰</span>
          </div>
          <h3>高性价比</h3>
          <p>云服务器15元/月起，虚拟主机8元/月起，高防游戏云25元/月起，价格为同行40%。</p>
        </div>
        
        <div class="feature-card">
          <div class="feature-icon">
            <span>🔧</span>
          </div>
          <h3>一键部署</h3>
          <p>支持预装1Panel/宝塔面板，一键式操作，三分钟即可开启建站！游戏云预装各类服务端，开箱即用。</p>
        </div>
        
        <div class="feature-card">
          <div class="feature-icon">
            <span>🛡️</span>
          </div>
          <h3>安全可靠</h3>
          <p>采用先进的安全措施和数据备份机制，确保用户数据的安全和服务的稳定运行。</p>
        </div>
        
        <div class="feature-card">
          <div class="feature-icon">
            <span>🌐</span>
          </div>
          <h3>全球部署</h3>
          <p>提供国内、香港、海外等多种服务器选择，满足不同地区业务需求，网络覆盖全球。</p>
        </div>
        
        <div class="feature-card">
          <div class="feature-icon">
            <span>👥</span>
          </div>
          <h3>专业服务</h3>
          <p>15000人技术交流社区，专业客服团队，随时为您解决各种技术问题。</p>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
export default {
  name: 'FeaturesSection'
}
</script>

<style scoped>
.features {
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

.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 2rem;
}

.feature-card {
  background: rgba(255, 255, 255, 0.7);  /* 半透明背景 */
  backdrop-filter: blur(10px);  /* 毛玻璃效果 */
  border-radius: 8px;
  padding: 2rem;
  text-align: center;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.05);
  transition: transform 0.3s;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.feature-card:hover {
  transform: translateY(-5px);
}

.feature-icon {
  font-size: 3rem;
  margin-bottom: 1.5rem;
}

.feature-card h3 {
  font-size: 1.4rem;
  margin-bottom: 1rem;
  color: #333;
}

.feature-card p {
  color: #666;
  line-height: 1.6;
}
</style>