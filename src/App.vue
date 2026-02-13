<template>
  <div id="app">
    <Header v-if="$route.path !== '/login' && $route.path !== '/register' && !isAdminRoute" />
    <router-view />
    <!-- 在管理后台页面隐藏底栏（路径以 /admin 开头时隐藏） -->
    <Footer v-if="$route.path !== '/login' && $route.path !== '/register' && !isAdminRoute" />
  </div>
</template>

<script>
import Header from './components/Header.vue'
import Footer from './components/Footer.vue'
import { getCachedImage, revokeAllCachedObjectURLs } from './utils/imageCache'

export default {
  name: 'App',
  components: {
    Header,
    Footer
  },
  computed: {
    isAdminRoute() {
      const p = this.$route && this.$route.path ? this.$route.path : '';
      return p.startsWith('/admin') || p === '/admin';
    }
  },
  mounted() {
    // 设置页面标题
    document.title = '仙林云-新一代云服务提供商';
    // 预加载并使用缓存的全局壁纸，避免在页面切换时重复向远端请求
    (async () => {
      const url = 'https://rba.kanostar.top/adapt'
      try {
        const objectUrl = await getCachedImage(url)
        if (objectUrl) {
          document.body.style.backgroundImage = `url('${objectUrl}')`
          document.body.style.backgroundSize = 'cover'
          document.body.style.backgroundRepeat = 'no-repeat'
          document.body.style.backgroundPosition = 'center center'
        }
      } catch (e) {
        // 如果把远端图片 fetch 为 blob 失败（常见于没有允许跨域 fetch 转 blob 的服务器），
        // 我们回退为直接将远程 URL 设为背景 — 浏览器会用普通的图片请求加载它，能避免背景丢失。
        console.warn('缓存背景加载失败，回退为直接使用远程 URL:', e)
        try {
          document.body.style.backgroundImage = `url('${url}')`
          document.body.style.backgroundSize = 'cover'
          document.body.style.backgroundRepeat = 'no-repeat'
          document.body.style.backgroundPosition = 'center center'
        } catch (e2) {
          console.warn('设置远程背景也失败，保持背景色作为最终回退', e2)
        }
      }
    })()
  },
  beforeUnmount() {
    // 可选：在应用销毁时释放 object URLs（多数 SPA 不会触发）
  }
}
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: 'PingFang SC', 'Microsoft YaHei', sans-serif;
  line-height: 1.6;
  /* 背景图片由 mounted() 通过缓存设置为 object URL，避免重复请求 */
  background-color: #f5f7fa;
  background-size: cover;
  min-height: 100vh;
}

#app {
  min-height: 100vh;
  background: transparent; /* 确保app组件不会覆盖body的背景 */
}
</style>
