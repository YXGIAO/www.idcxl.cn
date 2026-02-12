<template>
  <div class="products-management">
    <el-card>
      <template #header>
        <div class="clearfix">
          <span>产品管理</span>
          <el-button type="primary" size="small" style="float: right; margin-left: 10px;" @click="syncProductsDialogVisible = true">
            从供应商同步产品
          </el-button>
          <el-button type="success" size="small" style="float: right;" @click="fetchProducts">
            刷新
          </el-button>
        </div>
      </template>
      
      <!-- 搜索表单 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="产品状态">
          <el-select v-model="searchForm.status" placeholder="选择状态" clearable>
            <el-option label="全部" value=""></el-option>
            <el-option label="活跃" value="active"></el-option>
            <el-option label="非活跃" value="inactive"></el-option>
          </el-select>
        </el-form-item>
        
        <el-form-item label="供应商">
          <el-select 
            v-model="searchForm.supplier_id" 
            placeholder="选择供应商" 
            clearable
            filterable
            @change="handleSupplierChange"
          >
            <el-option label="全部供应商" value=""></el-option>
            <el-option 
              v-for="supplier in suppliers" 
              :key="supplier.id" 
              :label="supplier.name" 
              :value="supplier.id"
            >
              <span>{{ supplier.name }}</span>
              <span style="float: right; color: #8492a6; font-size: 13px">
                {{ supplier.type }}
              </span>
            </el-option>
          </el-select>
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="fetchProducts">查询</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
      
      <!-- 产品表格 -->
      <el-table :data="products" style="width: 100%" v-loading="loading" border>
        <el-table-column prop="id" label="ID" width="80" sortable></el-table-column>
        <el-table-column prop="name" label="产品名称" min-width="200">
          <template #default="scope">
            <span style="font-weight: bold;">{{ scope.row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" min-width="250" show-overflow-tooltip></el-table-column>
        <el-table-column prop="supplier_name" label="供应商" width="150">
          <template #default="scope">
            <el-tag 
              v-if="scope.row.supplier_name" 
              type="info" 
              size="small"
            >
              {{ scope.row.supplier_name }}
            </el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="price" label="价格" width="120" align="right">
          <template #default="scope">
            <span style="color: #e6a23c; font-weight: bold;">¥{{ scope.row.price.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)" size="small">
              {{ getStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" sortable></el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="scope">
            <el-button size="small" type="primary" @click="viewProduct(scope.row)">查看</el-button>
            <el-button size="small" @click="editProduct(scope.row)">编辑</el-button>
            <el-button 
              v-if="scope.row.supplier_id"
              size="small" 
              type="warning" 
              @click="syncSingleProduct(scope.row)"
            >
              同步
            </el-button>
            <el-button size="small" type="danger" @click="deleteProduct(scope.row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 分页 -->
      <el-pagination
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        :current-page="pagination.page"
        :page-sizes="[10, 20, 50, 100]"
        :page-size="pagination.limit"
        layout="total, sizes, prev, pager, next, jumper"
        :total="pagination.total"
        style="margin-top: 20px; text-align: center;"
      >
      </el-pagination>
    </el-card>
    
    <!-- 从供应商同步产品对话框 -->
    <el-dialog
      v-model="syncProductsDialogVisible"
      title="从供应商同步产品"
      width="500px"
      :before-close="handleSyncDialogClose"
    >
      <el-form :model="syncForm" label-width="120px">
        <el-form-item label="选择供应商" required>
          <el-select 
            v-model="syncForm.supplier_id" 
            placeholder="请选择供应商" 
            style="width: 100%"
            filterable
          >
            <el-option 
              v-for="supplier in suppliers" 
              :key="supplier.id" 
              :label="supplier.name" 
              :value="supplier.id"
            >
              <span>{{ supplier.name }}</span>
              <span style="float: right; color: #8492a6; font-size: 13px">
                {{ supplier.type }}
              </span>
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="同步方式">
          <el-radio-group v-model="syncForm.sync_type">
            <el-radio label="full">全量同步</el-radio>
            <el-radio label="incremental">增量同步</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="syncProductsDialogVisible = false">取消</el-button>
          <el-button 
            type="primary" 
            @click="syncProductsFromSupplier" 
            :loading="syncLoading"
          >
            开始同步
          </el-button>
        </span>
      </template>
    </el-dialog>
    
    <!-- 产品详情对话框 -->
    <el-dialog
      v-model="productDetailDialogVisible"
      title="产品详情"
      width="600px"
    >
      <el-descriptions :column="1" border v-if="currentProduct">
        <el-descriptions-item label="产品ID">{{ currentProduct.id }}</el-descriptions-item>
        <el-descriptions-item label="产品名称">{{ currentProduct.name }}</el-descriptions-item>
        <el-descriptions-item label="描述">{{ currentProduct.description }}</el-descriptions-item>
        <el-descriptions-item label="供应商">
          <el-tag v-if="currentProduct.supplier_name" type="info">
            {{ currentProduct.supplier_name }}
          </el-tag>
          <span v-else>-</span>
        </el-descriptions-item>
        <el-descriptions-item label="价格">¥{{ currentProduct.price.toFixed(2) }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(currentProduct.status)">
            {{ getStatusText(currentProduct.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ currentProduct.created_at }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ currentProduct.updated_at }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="productDetailDialogVisible = false">关闭</el-button>
          <el-button type="primary" @click="editProduct(currentProduct)">编辑</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { productAPI, supplierAPI } from '@/services/api';
import { ElMessage, ElMessageBox } from 'element-plus';

export default {
  name: 'ProductsManagement',
  data() {
    return {
      loading: false,
      syncLoading: false,
      products: [],
      suppliers: [],
      searchForm: {
        status: '',
        supplier_id: ''
      },
      pagination: {
        page: 1,
        limit: 10,
        total: 0
      },
      // 对话框控制
      syncProductsDialogVisible: false,
      productDetailDialogVisible: false,
      // 表单数据
      syncForm: {
        supplier_id: '',
        sync_type: 'full'
      },
      currentProduct: null
    };
  },
  mounted() {
    this.fetchProducts();
    this.fetchSuppliers();
  },
  methods: {
    // 获取产品列表
    async fetchProducts() {
      this.loading = true;
      try {
        const params = {
          page: this.pagination.page,
          limit: this.pagination.limit,
          ...this.searchForm
        };
        
        // 清理空值参数
        Object.keys(params).forEach(key => {
          if (params[key] === '' || params[key] === null || params[key] === undefined) {
            delete params[key];
          }
        });
        
        const response = await productAPI.getProducts(params);
        
        // 处理响应数据
        if (response && response.products) {
          this.products = response.products.map(item => {
            return {
              ...item,
              id: item.id,
              name: item.name || '',
              description: item.description || '',
              price: item.price || 0,
              status: item.status || 'active',
              supplier_name: item.supplier?.name || '',
              supplier_id: item.supplier_id || null,
              created_at: item.created_at ? new Date(item.created_at).toLocaleString() : '',
              updated_at: item.updated_at ? new Date(item.updated_at).toLocaleString() : ''
            };
          });
          this.pagination.total = response.total || response.products.length;
        } else {
          this.products = [];
          this.pagination.total = 0;
        }
      } catch (error) {
        console.error('获取产品列表失败:', error);
        ElMessage.error('获取产品列表失败: ' + (error.message || '网络错误'));
      } finally {
        this.loading = false;
      }
    },
    
    // 获取供应商列表
    async fetchSuppliers() {
      try {
        const response = await supplierAPI.getSuppliers({
          page: 1,
          limit: 1000, // 获取所有供应商
          status: 'active'
        });
        
        if (response && response.suppliers) {
          this.suppliers = response.suppliers;
        } else {
          this.suppliers = [];
        }
      } catch (error) {
        console.error('获取供应商列表失败:', error);
        ElMessage.error('获取供应商列表失败');
      }
    },
    
    // 供应商选择变化处理
    handleSupplierChange(value) {
      this.searchForm.supplier_id = value;
      this.pagination.page = 1;
      this.fetchProducts();
    },
    
    // 重置搜索
    resetSearch() {
      this.searchForm = {
        status: '',
        supplier_id: ''
      };
      this.pagination.page = 1;
      this.fetchProducts();
    },
    
    // 查看产品详情
    viewProduct(product) {
      this.currentProduct = product;
      this.productDetailDialogVisible = true;
    },
    
    // 编辑产品
    editProduct(product) {
      // 这里可以跳转到编辑页面或打开编辑对话框
      ElMessage.info('编辑功能待实现');
    },
    
    // 删除产品
    async deleteProduct(id) {
      try {
        await ElMessageBox.confirm(
          '确定要删除这个产品吗？此操作不可撤销。',
          '确认删除',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
        );
        
        // 调用删除API
        await productAPI.deleteProduct(id);
        ElMessage.success('产品删除成功');
        this.fetchProducts();
      } catch (error) {
        if (error !== 'cancel') {
          console.error('删除产品失败:', error);
          ElMessage.error('删除产品失败: ' + (error.message || '网络错误'));
        }
      }
    },
    
    // 单个产品同步
    async syncSingleProduct(product) {
      if (!product.supplier_id) {
        ElMessage.warning('该产品没有关联供应商，无法同步');
        return;
      }
      
      try {
        await ElMessageBox.confirm(
          `确定要同步产品 "${product.name}" 吗？`,
          '确认同步',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'info'
          }
        );
        
        this.loading = true;
        // 调用同步API
        await supplierAPI.syncProductsFromSupplier(product.supplier_id);
        ElMessage.success('产品同步成功');
        this.fetchProducts();
      } catch (error) {
        if (error !== 'cancel') {
          console.error('同步产品失败:', error);
          ElMessage.error('同步产品失败: ' + (error.message || '网络错误'));
        }
      } finally {
        this.loading = false;
      }
    },
    
    // 从供应商同步产品
    async syncProductsFromSupplier() {
      if (!this.syncForm.supplier_id) {
        ElMessage.warning('请选择供应商');
        return;
      }
      
      try {
        this.syncLoading = true;
        ElMessage.info('正在从供应商同步产品信息...');
        
        // 调用后端API同步产品
        const response = await supplierAPI.syncProductsFromSupplier(this.syncForm.supplier_id);
        
        ElMessage.success(`成功同步 ${response.syncedCount || 0} 个产品`);
        this.syncProductsDialogVisible = false;
        this.resetSyncForm();
        this.fetchProducts();
      } catch (error) {
        console.error('同步产品信息失败:', error);
        ElMessage.error('同步产品信息失败: ' + (error.message || '网络错误'));
      } finally {
        this.syncLoading = false;
      }
    },
    
    // 重置同步表单
    resetSyncForm() {
      this.syncForm = {
        supplier_id: '',
        sync_type: 'full'
      };
    },
    
    // 处理同步对话框关闭
    handleSyncDialogClose() {
      this.syncProductsDialogVisible = false;
      this.resetSyncForm();
    },
    
    // 状态标签样式
    getStatusType(status) {
      switch (status) {
        case 'active': return 'success';
        case 'inactive': return 'info';
        default: return 'info';
      }
    },
    
    getStatusText(status) {
      switch (status) {
        case 'active': return '活跃';
        case 'inactive': return '非活跃';
        default: return status;
      }
    },
    
    // 分页处理
    handleSizeChange(val) {
      this.pagination.limit = val;
      this.pagination.page = 1;
      this.fetchProducts();
    },
    
    handleCurrentChange(val) {
      this.pagination.page = val;
      this.fetchProducts();
    }
  }
};
</script>

<style scoped>
.products-management {
  padding: 20px;
}

.search-form {
  margin-bottom: 20px;
  padding: 20px;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.clearfix:before,
.clearfix:after {
  display: table;
  content: "";
}

.clearfix:after {
  clear: both
}

.dialog-footer {
  text-align: right;
}

.el-tag {
  margin-right: 5px;
}
</style>