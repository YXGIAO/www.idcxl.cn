<template>
  <div class="orders-management">
    <el-tabs v-model="activeTab">
      <el-tab-pane label="产品订单" name="product">
        <el-card>
          <template #header>
            <div class="clearfix">
              <span>产品订单管理</span>
            </div>
          </template>
          
          <el-form :inline="true" :model="productSearchForm" class="search-form">
            <el-form-item label="订单状态">
              <el-select v-model="productSearchForm.status" placeholder="选择状态">
                <el-option label="全部" value=""></el-option>
                <el-option label="待支付" value="pending"></el-option>
                <el-option label="已支付" value="paid"></el-option>
                <el-option label="已取消" value="cancelled"></el-option>
                <el-option label="已完成" value="completed"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item label="用户ID">
              <el-input v-model="productSearchForm.user_id" placeholder="用户ID"></el-input>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="fetchProductOrders">查询</el-button>
            </el-form-item>
          </el-form>
          
          <el-table :data="productOrders" style="width: 100%" v-loading="productLoading">
            <el-table-column prop="id" label="ID" width="80"></el-table-column>
            <el-table-column prop="order_number" label="订单号" width="200"></el-table-column>
            <el-table-column prop="user.username" label="用户名" width="120"></el-table-column>
            <el-table-column prop="product.name" label="产品名称"></el-table-column>
            <el-table-column prop="amount" label="金额" width="100">
              <template #default="scope">¥{{ scope.row.amount }}</template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template #default="scope">
                <el-tag :type="getOrderStatusType(scope.row.status)">
                  {{ getOrderStatusText(scope.row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="payment_method" label="支付方式" width="120"></el-table-column>
            <el-table-column prop="created_at" label="创建时间" width="160"></el-table-column>
            <el-table-column label="操作" width="150">
              <template #default="scope">
                <el-button size="mini" @click="viewOrder(scope.row)">查看</el-button>
                <el-button size="mini" type="primary" @click="editOrder(scope.row)">编辑</el-button>
              </template>
            </el-table-column>
          </el-table>
          
          <el-pagination
            @size-change="handleProductSizeChange"
            @current-change="handleProductCurrentChange"
            :current-page="productPagination.page"
            :page-sizes="[10, 20, 50, 100]"
            :page-size="productPagination.limit"
            layout="total, sizes, prev, pager, next, jumper"
            :total="productPagination.total">
          </el-pagination>
        </el-card>
      </el-tab-pane>
      
      <el-tab-pane label="续费订单" name="renewal">
        <el-card>
          <template #header>
            <div class="clearfix">
              <span>续费订单管理</span>
            </div>
          </template>
          
          <el-form :inline="true" :model="renewalSearchForm" class="search-form">
            <el-form-item label="订单状态">
              <el-select v-model="renewalSearchForm.status" placeholder="选择状态">
                <el-option label="全部" value=""></el-option>
                <el-option label="待支付" value="pending"></el-option>
                <el-option label="已支付" value="paid"></el-option>
                <el-option label="已取消" value="cancelled"></el-option>
                <el-option label="已完成" value="completed"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item label="用户ID">
              <el-input v-model="renewalSearchForm.user_id" placeholder="用户ID"></el-input>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="fetchRenewalOrders">查询</el-button>
            </el-form-item>
          </el-form>
          
          <el-table :data="renewalOrders" style="width: 100%" v-loading="renewalLoading">
            <el-table-column prop="id" label="ID" width="80"></el-table-column>
            <el-table-column prop="order_number" label="订单号" width="200"></el-table-column>
            <el-table-column prop="user.username" label="用户名" width="120"></el-table-column>
            <el-table-column prop="product.name" label="产品名称"></el-table-column>
            <el-table-column prop="amount" label="金额" width="100">
              <template #default="scope">¥{{ scope.row.amount }}</template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template #default="scope">
                <el-tag :type="getOrderStatusType(scope.row.status)">
                  {{ getOrderStatusText(scope.row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="renewal_period" label="续费月数" width="100"></el-table-column>
            <el-table-column prop="expire_at" label="到期时间" width="160"></el-table-column>
            <el-table-column prop="created_at" label="创建时间" width="160"></el-table-column>
            <el-table-column label="操作" width="150">
              <template #default="scope">
                <el-button size="mini" @click="viewOrder(scope.row)">查看</el-button>
                <el-button size="mini" type="primary" @click="editOrder(scope.row)">编辑</el-button>
              </template>
            </el-table-column>
          </el-table>
          
          <el-pagination
            @size-change="handleRenewalSizeChange"
            @current-change="handleRenewalCurrentChange"
            :current-page="renewalPagination.page"
            :page-sizes="[10, 20, 50, 100]"
            :page-size="renewalPagination.limit"
            layout="total, sizes, prev, pager, next, jumper"
            :total="renewalPagination.total">
          </el-pagination>
        </el-card>
      </el-tab-pane>
    </el-tabs>

    <!-- 订单详情对话框 -->
    <el-dialog title="订单详情" :visible.sync="orderDialogVisible" width="60%">
      <div v-if="currentOrder">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="订单ID">{{ currentOrder.id }}</el-descriptions-item>
          <el-descriptions-item label="订单号">{{ currentOrder.order_number }}</el-descriptions-item>
          <el-descriptions-item label="用户名">{{ currentOrder.user?.username }}</el-descriptions-item>
          <el-descriptions-item label="产品名称">{{ currentOrder.product?.name }}</el-descriptions-item>
          <el-descriptions-item label="金额">¥{{ currentOrder.amount }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="getOrderStatusType(currentOrder.status)">
              {{ getOrderStatusText(currentOrder.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="支付方式">{{ currentOrder.payment_method }}</el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ currentOrder.created_at }}</el-descriptions-item>
          <el-descriptions-item label="更新时间">{{ currentOrder.updated_at }}</el-descriptions-item>
          <el-descriptions-item label="备注" :span="2">{{ currentOrder.description || '无' }}</el-descriptions-item>
        </el-descriptions>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { orderAPI } from '@/services/api.js';

export default {
  name: 'OrdersManagement',
  data() {
    return {
      activeTab: 'product',
      productOrders: [],
      renewalOrders: [],
      productLoading: false,
      renewalLoading: false,
      productSearchForm: {
        status: '',
        user_id: ''
      },
      renewalSearchForm: {
        status: '',
        user_id: ''
      },
      productPagination: {
        page: 1,
        limit: 10,
        total: 0
      },
      renewalPagination: {
        page: 1,
        limit: 10,
        total: 0
      },
      orderDialogVisible: false,
      currentOrder: null
    };
  },
  mounted() {
    this.fetchProductOrders();
    this.fetchRenewalOrders();
  },
  methods: {
    async fetchProductOrders() {
      this.productLoading = true;
      try {
        const params = {
          page: this.productPagination.page,
          limit: this.productPagination.limit,
          status: this.productSearchForm.status,
          user_id: this.productSearchForm.user_id
        };
        const response = await orderAPI.getProductOrders(params);
        this.productOrders = response.orders;
        this.productPagination.total = response.total;
      } catch (error) {
        this.$message.error('获取产品订单列表失败');
        console.error(error);
      } finally {
        this.productLoading = false;
      }
    },
    async fetchRenewalOrders() {
      this.renewalLoading = true;
      try {
        const params = {
          page: this.renewalPagination.page,
          limit: this.renewalPagination.limit,
          status: this.renewalSearchForm.status,
          user_id: this.renewalSearchForm.user_id
        };
        const response = await orderAPI.getRenewalOrders(params);
        this.renewalOrders = response.orders;
        this.renewalPagination.total = response.total;
      } catch (error) {
        this.$message.error('获取续费订单列表失败');
        console.error(error);
      } finally {
        this.renewalLoading = false;
      }
    },
    getOrderStatusType(status) {
      switch (status) {
        case 'pending': return 'warning';
        case 'paid': return 'success';
        case 'cancelled': return 'info';
        case 'completed': return 'primary';
        default: return 'info';
      }
    },
    getOrderStatusText(status) {
      switch (status) {
        case 'pending': return '待支付';
        case 'paid': return '已支付';
        case 'cancelled': return '已取消';
        case 'completed': return '已完成';
        default: return status;
      }
    },
    viewOrder(order) {
      this.currentOrder = order;
      this.orderDialogVisible = true;
    },
    editOrder(order) {
      this.$message.warning('订单编辑功能暂未开放，如需修改请联系管理员');
    },
    async viewOrder(order) {
      this.currentOrder = order;
      this.orderDialogVisible = true;
    },
    handleProductSizeChange(val) {
      this.productPagination.limit = val;
      this.fetchProductOrders();
    },
    handleProductCurrentChange(val) {
      this.productPagination.page = val;
      this.fetchProductOrders();
    },
    handleRenewalSizeChange(val) {
      this.renewalPagination.limit = val;
      this.fetchRenewalOrders();
    },
    handleRenewalCurrentChange(val) {
      this.renewalPagination.page = val;
      this.fetchRenewalOrders();
    }
  }
};
</script>

<style scoped>
.search-form {
  margin-bottom: 20px;
}
</style>