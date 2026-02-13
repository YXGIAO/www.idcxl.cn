 <template>
  <section id="home" class="hero">
    <div class="container">
      <div class="hero-content">
        <div class="text-content">
          <h1>新一代云服务提供商</h1>
          <p>仙林云是国内领先的具有自主知识产权的云计算平台，有云服务器、虚拟主机、游戏云、CDN、云电脑等产品，无论你是新手小白，还是云上大佬，都能给您优秀的体验。</p>
          <div class="hero-buttons">
            <button class="primary-btn">立即选购</button>
            <button class="secondary-btn">了解更多</button>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
export default {
  name: 'HeroSection'
}
</script>

<style scoped>
.hero {
  padding: 10rem 0 5rem;
  min-height: 100vh;
  display: flex;
  align-items: center;
  background-image: url('https://via.placeholder.com/1920x1080/007bff/white?text=仙林云+Cloud');
  /* 使用 cover 确保背景图像覆盖整个区域，同时保持宽高比 */
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  background-attachment: scroll;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 2rem;
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 80vh;
}

.hero-content {
  width: 100%;
  max-width: 600px;
  margin: 0 auto;
  text-align: center;
}

.text-content {
  background: rgba(255, 255, 255, 0.85);  /* 调整半透明背景 */
  backdrop-filter: blur(10px);  /* 毛玻璃效果 */
  border: 1px solid rgba(255, 255, 255, 0.2);  /* 添加边框以增强毛玻璃效果 */
  border-radius: 10px;  /* 添加圆角 */
  padding: 2.5rem;  /* 增加内边距 */
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);  /* 添加阴影 */
}

.hero-content h1 {
  font-size: 2.5rem;
  margin-bottom: 1.5rem;
  color: #333;
  line-height: 1.2;
}

.hero-content p {
  font-size: 1.2rem;
  color: #666;
  margin-bottom: 2rem;
  line-height: 1.8;
}

.hero-buttons {
  display: flex;
  gap: 1rem;
  justify-content: center;
  flex-wrap: wrap;
}

.primary-btn {
  background-color: #007bff;
  color: white;
  padding: 0.8rem 2rem;
  border: none;
  border-radius: 4px;
  font-size: 1.1rem;
  cursor: pointer;
  transition: background-color 0.3s;
}

.primary-btn:hover {
  background-color: #0056b3;
}

.secondary-btn {
  background-color: transparent;
  color: #007bff;
  padding: 0.8rem 2rem;
  border: 2px solid #007bff;
  border-radius: 4px;
  font-size: 1.1rem;
  cursor: pointer;
  transition: all 0.3s;
}

.secondary-btn:hover {
  background-color: #007bff;
  color: white;
}

@media (max-width: 768px) {
  .container {
    padding: 0 1rem;
  }
  
  .hero-content {
    padding: 0 1rem;
  }
  
  .text-content {
    padding: 1.5rem;  /* 移动设备上减少内边距 */
  }
  
  .hero-content h1 {
    font-size: 2rem;  /* 移动设备上减小标题大小 */
  }
  
  .hero-content p {
    font-size: 1rem;  /* 移动设备上减小文字大小 */
  }
}
</style>