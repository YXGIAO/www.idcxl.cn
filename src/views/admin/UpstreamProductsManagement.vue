<template>
  <div class="upstream-products-management">
    <el-card>
      <template #header>
        <div class="clearfix">
          <span>上游产品管理</span>
          <el-button type="primary" size="small" style="float: right;" @click="fetchProducts">
            刷新
          </el-button>
        </div>
      </template>
      
      <!-- 搜索表单 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键字">
          <el-input
            v-model="searchForm.keywords"
            placeholder="ID/用户名/邮箱/手机号/商品名称/产品标识"
            clearable
            style="width: 250px;"
          />
        </el-form-item>
        
        <el-form-item label="供应商ID">
          <el-input
            v-model="searchForm.supplier_id"
            placeholder="供应商ID"
            clearable
            style="width: 120px;"
            type="number"
          />
        </el-form-item>
        
        <el-form-item label="付款周期">
          <el-select v-model="searchForm.billing_cycle" placeholder="全部" clearable>
            <el-option label="全部" value=""></el-option>
            <el-option label="月付" value="monthly"></el-option>
            <el-option label="季付" value="quarterly"></el-option>
            <el-option label="半年付" value="halfyearly"></el-option>
            <el-option label="年付" value="yearly"></el-option>
            <el-option label="两年付" value="biennially"></el-option>
            <el-option label="三年付" value="triennially"></el-option>
          </el-select>
        </el-form-item>
        
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="全部" value=""></el-option>
            <el-option label="未付款" value="Unpaid"></el-option>
            <el-option label="开通中" value="Pending"></el-option>
            <el-option label="已开通" value="Active"></el-option>
            <el-option label="已暂停" value="Suspended"></el-option>
            <el-option label="已删除" value="Deleted"></el-option>
            <el-option label="开通失败" value="Failed"></el-option>
          </el-select>
        </el-form-item>
        
        <el-form-item label="时间范围">
          <el-date-picker
            v-model="dateRange"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="x"
            :default-time="[new Date(2000, 1, 1, 0, 0, 0), new Date(2000, 1, 1, 23, 59, 59)]"
          />
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
      
      <!-- 产品列表表格 -->
      <el-table :data="products" style="width: 100%" v-loading="loading" border>
        <el-table-column prop="id" label="ID" width="80" align="center"></el-table-column>
        <el-table-column prop="name" label="产品标识" width="120"></el-table-column>
        <el-table-column prop="product_name" label="商品名称" width="150"></el-table-column>
        
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)" size="small">
              {{ getStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="first_payment_amount" label="首付款" width="100" align="right">
          <template #default="scope">
            {{ formatAmount(scope.row.first_payment_amount) }}
          </template>
        </el-table-column>
        
        <el-table-column prop="renew_amount" label="续费金额" width="100" align="right">
          <template #default="scope">
            {{ formatAmount(scope.row.renew_amount) }}
          </template>
        </el-table-column>
        
        <el-table-column prop="billing_cycle" label="付款周期" width="80" align="center">
          <template #default="scope">
            {{ getBillingCycleText(scope.row.billing_cycle) }}
          </template>
        </el-table-column>
        
        <el-table-column prop="due_time" label="到期时间" width="160"></el-table-column>
        
        <el-table-column prop="username" label="用户名" width="120"></el-table-column>
        <el-table-column prop="company" label="公司名称" width="150"></el-table-column>
        <el-table-column prop="email" label="邮箱" width="180"></el-table-column>
        <el-table-column prop="phone" label="手机号" width="120">
          <template #default="scope">
            {{ scope.row.phone_code }}{{ scope.row.phone }}
          </template>
        </el-table-column>
        
        <el-table-column prop="ip_num" label="IP数量" width="80" align="center"></el-table-column>
        <el-table-column prop="base_info" label="基础信息" min-width="200"></el-table-column>
        
        <el-table-column prop="created_at" label="创建时间" width="160"></el-table-column>
        
        <el-table-column label="操作" width="120" align="center" fixed="right">
          <template #default="scope">
            <el-button type="primary" size="small" @click="handleView(scope.row)">查看</el-button>
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
        style="margin-top: 20px;"
      />
    </el-card>
  </div>
</template>

<script>
import { upstreamProductAPI } from '@/services/api';
import { ElMessage } from 'element-plus';

export default {
  name: 'UpstreamProductsManagement',
  data() {
    return {
      loading: false,
      products: [],
      dateRange: null,
      searchForm: {
        keywords: '',
        supplier_id: '',
        billing_cycle: '',
        status: '',
        start_time: null,
        end_time: null,
        orderby: 'id',
        sort: 'desc'
      },
      pagination: {
        page: 1,
        limit: 10,
        total: 0
      }
    };
  },
  mounted() {
    this.fetchProducts();
  },
  watch: {
    dateRange(newVal) {
      if (newVal && newVal.length === 2) {
        this.searchForm.start_time = Math.floor(newVal[0] / 1000); // 转换为秒
        this.searchForm.end_time = Math.floor(newVal[1] / 1000);
      } else {
        this.searchForm.start_time = null;
        this.searchForm.end_time = null;
      }
    }
  },
  methods: {
    // 获取产品列表
    async fetchProducts() {
      this.loading = true;
      try {
        // 构建请求参数
        const params = {
          page: this.pagination.page,
          limit: this.pagination.limit,
          orderby: this.searchForm.orderby,
          sort: this.searchForm.sort
        };
        
        // 添加筛选参数
        if (this.searchForm.keywords) params.keywords = this.searchForm.keywords;
        if (this.searchForm.supplier_id) params.supplier_id = parseInt(this.searchForm.supplier_id);
        if (this.searchForm.billing_cycle) params.billing_cycle = this.searchForm.billing_cycle;
        if (this.searchForm.status) params.status = this.searchForm.status;
        if (this.searchForm.start_time) params.start_time = this.searchForm.start_time;
        if (this.searchForm.end_time) params.end_time = this.searchForm.end_time;
        
        const response = await upstreamProductAPI.getUpstreamProducts(params);
        
        // 处理响应数据
        if (response && response.code === 200 && response.data) {
          this.products = response.data.list || [];
          this.pagination.total = response.data.count || 0;
          
          // 格式化数据
          this.products = this.products.map(item => {
            return {
              ...item,
              due_time: item.due_time || '-',
              created_at: item.created_at ? new Date(item.created_at).toLocaleString() : '-',
              phone_code: item.phone_code || '+86',
              phone: item.phone || '-',
              base_info: item.base_info || '-'
            };
          });
        } else {
          this.products = [];
          this.pagination.total = 0;
          ElMessage.warning('暂无数据');
        }
      } catch (error) {
        console.error('获取上游产品列表失败:', error);
        ElMessage.error('获取上游产品列表失败: ' + (error.message || '网络错误'));
      } finally {
        this.loading = false;
      }
    },
    
    // 搜索
    handleSearch() {
      this.pagination.page = 1;
      this.fetchProducts();
    },
    
    // 重置搜索条件
    handleReset() {
      this.searchForm = {
        keywords: '',
        supplier_id: '',
        billing_cycle: '',
        status: '',
        start_time: null,
        end_time: null,
        orderby: 'id',
        sort: 'desc'
      };
      this.dateRange = null;
      this.pagination.page = 1;
      this.fetchProducts();
    },
    
    // 查看产品详情
    handleView(product) {
      // 跳转到产品详情页
      this.$router.push(`/admin/upstream-products/${product.id}`);
    },
    
    // 分页大小改变
    handleSizeChange(val) {
      this.pagination.limit = val;
      this.pagination.page = 1;
      this.fetchProducts();
    },
    
    // 当前页改变
    handleCurrentChange(val) {
      this.pagination.page = val;
      this.fetchProducts();
    },
    
    // 获取状态类型
    getStatusType(status) {
      switch (status) {
        case 'Active': return 'success';
        case 'Pending': return 'warning';
        case 'Unpaid': return 'danger';
        case 'Suspended': return 'info';
        case 'Deleted': return 'info';
        case 'Failed': return 'danger';
        default: return 'info';
      }
    },
    
    // 获取状态文本
    getStatusText(status) {
      switch (status) {
        case 'Active': return '已开通';
        case 'Pending': return '开通中';
        case 'Unpaid': return '未付款';
        case 'Suspended': return '已暂停';
        case 'Deleted': return '已删除';
        case 'Failed': return '开通失败';
        default: return status;
      }
    },
    
    // 获取付款周期文本
    getBillingCycleText(cycle) {
      switch (cycle) {
        case 'monthly': return '月付';
        case 'quarterly': return '季付';
        case 'halfyearly': return '半年付';
        case 'yearly': return '年付';
        case 'biennially': return '两年付';
        case 'triennially': return '三年付';
        default: return cycle;
      }
    },
    
    // 格式化金额
    formatAmount(amount) {
      if (!amount) return '-';
      // 如果已经是字符串格式，直接返回
      if (typeof amount === 'string') {
        return amount;
      }
      // 如果是数字，格式化为两位小数
      if (typeof amount === 'number') {
        return '¥' + amount.toFixed(2);
      }
      return amount;
    }
  }
};
</script>

<style scoped>
.upstream-products-management {
  padding: 20px;
}

.search-form {
  margin-bottom: 20px;
}

.el-table {
  margin-top: 20px;
}

.el-pagination {
  margin-top: 20px;
  text-align: right;
}
</style>