<template>
  <div class="resources-management">
    <el-card>
      <template #header>
        <div class="clearfix">
          <span>资源管理</span>
        </div>
      </template>
      
      <el-tabs v-model="activeTab" type="card">
        <el-tab-pane label="供货商管理" name="suppliers">
          <div class="tab-content">
            <suppliers-management v-if="activeTab === 'suppliers'" />
          </div>
        </el-tab-pane>
        
        <el-tab-pane label="产品管理" name="products">
          <div class="tab-content">
            <products-management v-if="activeTab === 'products'" />
          </div>
        </el-tab-pane>
        
        <el-tab-pane label="服务器管理" name="servers">
          <div class="tab-content">
            <servers-management v-if="activeTab === 'servers'" />
          </div>
        </el-tab-pane>
        
        <el-tab-pane label="智简魔方管理" name="zjmf">
          <div class="tab-content">
            <zjmf-management v-if="activeTab === 'zjmf'" />
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script>
import SuppliersManagement from './SuppliersManagement.vue';
import ServersManagement from './ServersManagement.vue';
import ZJMFManagement from './ZJMFManagement.vue';
import ProductsManagement from './ProductsManagement.vue';

export default {
  name: 'ResourcesManagement',
  components: {
    SuppliersManagement,
    ServersManagement,
    ZJMFManagement,
    ProductsManagement
  },
  data() {
    return {
      activeTab: 'suppliers'
    };
  },
  mounted() {
    // 初始化时设置默认激活的标签页
    this.activeTab = 'suppliers';
  }
};
</script>

<style scoped>
.resources-management {
  padding: 20px;
}

.tab-content {
  padding: 20px 0;
}

.clearfix {
  text-align: left;
  font-weight: bold;
  font-size: 18px;
}
</style>