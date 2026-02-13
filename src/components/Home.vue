<template>
  <div>
    <HeroSection />
    <ProductsSection />
    <FeaturesSection />
    <PricingSection />
  </div>
</template>

<script>
import HeroSection from './HeroSection.vue'
import ProductsSection from './ProductsSection.vue'
import FeaturesSection from './FeaturesSection.vue'
import PricingSection from './PricingSection.vue'

export default {
  name: 'Home',
  components: {
    HeroSection,
    ProductsSection,
    FeaturesSection,
    PricingSection
  },
  mounted() {
    // 检查 URL 是否包含 hash，并滚动到相应部分
    this.$nextTick(() => {
      if (this.$route.hash) {
        const element = document.querySelector(this.$route.hash);
        if (element) {
          element.scrollIntoView({ behavior: 'smooth' });
        }
      }
    });
  },
  watch: {
    '$route.hash'(newHash) {
      // 监听路由 hash 变化并滚动到相应部分
      if (newHash) {
        const element = document.querySelector(newHash);
        if (element) {
          element.scrollIntoView({ behavior: 'smooth' });
        }
      }
    }
  }
}
</script>

<style>
/* Home组件不需要特殊样式，各子组件会处理自己的样式 */
</style>