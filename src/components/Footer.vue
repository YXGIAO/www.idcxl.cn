<template>
  <footer class="footer">
    <div class="container">
      <div class="footer-content">
        <div class="footer-section">
          <h3>产品</h3>
          <ul>
            <li><a href="#">云服务器</a></li>
            <li><a href="#">虚拟主机</a></li>
            <li><a href="#">游戏云服务器</a></li>
            <li><a href="#">对象存储</a></li>
            <li><a href="#">云盾CDN</a></li>
            <li><a href="#">域名服务</a></li>
          </ul>
        </div>
        
        <div class="footer-section">
          <h3>服务</h3>
          <ul>
            <li><a href="#">网站搭建</a></li>
            <li><a href="#">游戏服务器</a></li>
            <li><a href="#">论坛社区</a></li>
            <li><a href="#">一键部署</a></li>
            <li><a href="#">技术支持</a></li>
          </ul>
        </div>
        
        <div class="footer-section">
          <h3>支持</h3>
          <ul>
            <li><a href="#">文档中心</a></li>
            <li><a href="#">API文档</a></li>
            <li><a href="#">工单系统</a></li>
            <li><a href="#">在线咨询</a></li>
            <li><a href="#">社区论坛</a></li>
          </ul>
        </div>
        
        <div class="footer-section">
          <h3>关于我们</h3>
          <ul>
            <li><a href="#">公司简介</a></li>
            <li><a href="#">联系我们</a></li>
            <li><a href="#">合作伙伴</a></li>
            <li><a href="#">招聘信息</a></li>
            <li><a href="#">用户协议</a></li>
          </ul>
        </div>
        
        <div class="footer-section contact-section">
          <h3>联系我们</h3>
          <div class="contact-info">
            <p>客服QQ：12345678</p>
            <p>邮箱：support@xianlinyun.com</p>
            <p>微信：xianlinyun_service</p>
            <div class="social-links">
              <a href="#"><img src="https://via.placeholder.com/30x30/007bff/white?text=Q" alt="QQ群"></a>
              <a href="#"><img src="https://via.placeholder.com/30x30/1da1f2/white?text=T" alt="Twitter"></a>
              <a href="#"><img src="https://via.placeholder.com/30x30/1877f2/white?text=F" alt="Facebook"></a>
            </div>
          </div>
        </div>
      </div>
      
      <div class="footer-bottom">
        <div class="copyright">
          <p>&copy; 2023 仙林云 - 新一代云服务提供商. 保留所有权利.</p>
          <p>粤ICP备19052839号-1</p>
        </div>
      </div>
    </div>
  </footer>
</template>

<script>
export default {
  name: 'Footer'
}
</script>

<style scoped>
.footer {
  background-color: rgba(52, 58, 64, 0.5);  /* 50% 透明度 */
  backdrop-filter: blur(10px);  /* 毛玻璃效果 */
  color: #fff;
  padding-top: 3rem;
  position: relative;
  z-index: 10;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 2rem;
}

.footer-content {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 2rem;
  padding-bottom: 2rem;
  border-bottom: 1px solid rgba(73, 80, 87, 0.5);
  margin-bottom: 2rem;
}

.footer-section h3 {
  color: #fff;
  margin-bottom: 1rem;
  font-size: 1.2rem;
}

.footer-section ul {
  list-style: none;
}

.footer-section ul li {
  margin-bottom: 0.5rem;
}

.footer-section ul li a {
  color: #adb5bd;
  text-decoration: none;
  transition: color 0.3s;
}

.footer-section ul li a:hover {
  color: #007bff;
}

.contact-info p {
  margin-bottom: 0.5rem;
}

.social-links {
  display: flex;
  gap: 1rem;
  margin-top: 1rem;
}

.social-links a {
  display: inline-block;
}

.social-links img {
  width: 30px;
  height: 30px;
  border-radius: 50%;
}

.footer-bottom {
  text-align: center;
  padding: 1.5rem 0;
}

.copyright p {
  color: #adb5bd;
  font-size: 0.9rem;
  margin-bottom: 0.3rem;
}
</style>