<template>
  <div class="finance-management">
    <el-tabs v-model="activeTab">
      <el-tab-pane label="交易流水" name="transactions">
        <el-card>
          <template #header>
            <div class="clearfix">
              <span>交易流水管理</span>
            </div>
          </template>
          
          <el-form :inline="true" :model="transactionSearchForm" class="search-form">
            <el-form-item label="交易类型">
              <el-select v-model="transactionSearchForm.type" placeholder="选择类型">
                <el-option label="全部" value=""></el-option>
                <el-option label="收入" value="income"></el-option>
                <el-option label="支出" value="expense"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item label="状态">
              <el-select v-model="transactionSearchForm.status" placeholder="选择状态">
                <el-option label="全部" value=""></el-option>
                <el-option label="完成" value="completed"></el-option>
                <el-option label="待处理" value="pending"></el-option>
                <el-option label="失败" value="failed"></el-option>
                <el-option label="已退款" value="refunded"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item label="用户ID">
              <el-input v-model="transactionSearchForm.user_id" placeholder="用户ID"></el-input>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="fetchTransactions">查询</el-button>
            </el-form-item>
          </el-form>
          
          <el-table :data="transactions" style="width: 100%" v-loading="transactionLoading">
            <el-table-column prop="id" label="ID" width="80"></el-table-column>
            <el-table-column prop="transaction_number" label="交易号" width="200"></el-table-column>
            <el-table-column prop="user.username" label="用户名" width="120"></el-table-column>
            <el-table-column prop="order_type" label="订单类型" width="120">
              <template #default="scope">
                <el-tag :type="getOrderTypeTag(scope.row.order_type)">
                  {{ getOrderTypeText(scope.row.order_type) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="amount" label="金额" width="120">
              <template #default="scope">
                <span :class="scope.row.type === 'income' ? 'income' : 'expense'">
                  {{ scope.row.type === 'income' ? '+' : '-' }}¥{{ scope.row.amount }}
                </span>
              </template>
            </el-table-column>
            <el-table-column prop="type" label="类型" width="100">
              <template #default="scope">
                <el-tag :type="scope.row.type === 'income' ? 'success' : 'danger'">
                  {{ scope.row.type === 'income' ? '收入' : '支出' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template #default="scope">
                <el-tag :type="getTransactionStatusType(scope.row.status)">
                  {{ getTransactionStatusText(scope.row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="description" label="描述"></el-table-column>
            <el-table-column prop="created_at" label="创建时间" width="160"></el-table-column>
          </el-table>
          
          <el-pagination
            @size-change="handleTransactionSizeChange"
            @current-change="handleTransactionCurrentChange"
            :current-page="transactionPagination.page"
            :page-sizes="[10, 20, 50, 100]"
            :page-size="transactionPagination.limit"
            layout="total, sizes, prev, pager, next, jumper"
            :total="transactionPagination.total">
          </el-pagination>
        </el-card>
      </el-tab-pane>
      
      <el-tab-pane label="账单管理" name="bills">
        <el-card>
          <template #header>
            <div class="clearfix">
              <span>账单管理</span>
              <el-button style="float: right; padding: 3px 0" type="text" @click="createBill">创建账单</el-button>
            </div>
          </template>
          
          <el-form :inline="true" :model="billSearchForm" class="search-form">
            <el-form-item label="状态">
              <el-select v-model="billSearchForm.status" placeholder="选择状态">
                <el-option label="全部" value=""></el-option>
                <el-option label="未支付" value="unpaid"></el-option>
                <el-option label="已支付" value="paid"></el-option>
                <el-option label="已逾期" value="overdue"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item label="用户ID">
              <el-input v-model="billSearchForm.user_id" placeholder="用户ID"></el-input>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="fetchBills">查询</el-button>
            </el-form-item>
          </el-form>
          
          <el-table :data="bills" style="width: 100%" v-loading="billLoading">
            <el-table-column prop="id" label="ID" width="80"></el-table-column>
            <el-table-column prop="user.username" label="用户名" width="120"></el-table-column>
            <el-table-column prop="amount" label="金额" width="120">
              <template #default="scope">¥{{ scope.row.amount }}</template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template #default="scope">
                <el-tag :type="getBillStatusType(scope.row.status)">
                  {{ getBillStatusText(scope.row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="due_date" label="到期时间" width="160"></el-table-column>
            <el-table-column prop="issue_date" label="开票时间" width="160"></el-table-column>
            <el-table-column prop="description" label="描述"></el-table-column>
            <el-table-column prop="created_at" label="创建时间" width="160"></el-table-column>
            <el-table-column label="操作" width="150">
              <template #default="scope">
                <el-button size="mini" @click="viewBill(scope.row)">查看</el-button>
                <el-button size="mini" type="primary" @click="editBill(scope.row)" :disabled="scope.row.status === 'paid'">编辑</el-button>
                <el-button size="mini" type="success" @click="payBill(scope.row)" :disabled="scope.row.status === 'paid'">标记支付</el-button>
              </template>
            </el-table-column>
          </el-table>
          
          <el-pagination
            @size-change="handleBillSizeChange"
            @current-change="handleBillCurrentChange"
            :current-page="billPagination.page"
            :page-sizes="[10, 20, 50, 100]"
            :page-size="billPagination.limit"
            layout="total, sizes, prev, pager, next, jumper"
            :total="billPagination.total">
          </el-pagination>
        </el-card>
      </el-tab-pane>
    </el-tabs>

    <!-- 账单编辑对话框 -->
    <el-dialog :title="billDialogTitle" :visible.sync="billDialogVisible" width="50%">
      <el-form :model="billForm" :rules="billRules" ref="billForm" label-width="100px">
        <el-form-item label="用户ID" prop="user_id">
          <el-input v-model="billForm.user_id" :disabled="!!billForm.id"></el-input>
        </el-form-item>
        <el-form-item label="金额" prop="amount">
          <el-input-number v-model="billForm.amount" :min="0" :precision="2" :step="1"></el-input-number>
        </el-form-item>
        <el-form-item label="到期时间" prop="due_date">
          <el-date-picker
            v-model="billForm.due_date"
            type="datetime"
            placeholder="选择到期时间">
          </el-date-picker>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="billForm.description" type="textarea"></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="billDialogVisible = false">取 消</el-button>
          <el-button type="primary" @click="submitBillForm">确 定</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { financeAPI } from '@/services/api';

export default {
  name: 'FinanceManagement',
  data() {
    return {
      activeTab: 'transactions',
      transactions: [],
      bills: [],
      transactionLoading: false,
      billLoading: false,
      transactionSearchForm: {
        type: '',
        status: '',
        user_id: ''
      },
      billSearchForm: {
        status: '',
        user_id: ''
      },
      transactionPagination: {
        page: 1,
        limit: 10,
        total: 0
      },
      billPagination: {
        page: 1,
        limit: 10,
        total: 0
      },
      billDialogVisible: false,
      billDialogTitle: '',
      billForm: {
        id: null,
        user_id: '',
        amount: 0,
        due_date: '',
        description: ''
      },
      billRules: {
        user_id: [
          { required: true, message: '请输入用户ID', trigger: 'blur' }
        ],
        amount: [
          { required: true, message: '请输入金额', trigger: 'blur' }
        ],
        due_date: [
          { required: true, message: '请选择到期时间', trigger: 'change' }
        ]
      }
    };
  },
  mounted() {
    this.fetchTransactions();
    this.fetchBills();
  },
  methods: {
    async fetchTransactions() {
      this.transactionLoading = true;
      try {
        // 这里应该调用API获取交易流水列表
        // 模拟数据
        this.transactions = [
          {
            id: 1,
            transaction_number: 'TXN202312270001',
            user: { username: 'user1' },
            order_type: 'product_order',
            amount: 99.00,
            type: 'income',
            status: 'completed',
            description: '购买基础云服务器',
            created_at: '2023-12-27 10:00:00'
          },
          {
            id: 2,
            transaction_number: 'TXN202312270002',
            user: { username: 'user2' },
            order_type: 'renewal_order',
            amount: 990.00,
            type: 'income',
            status: 'completed',
            description: '续费专业云服务器',
            created_at: '2023-12-27 11:00:00'
          }
        ];
        this.transactionPagination.total = 2;
      } catch (error) {
        this.$message.error('获取交易流水列表失败');
        console.error(error);
      } finally {
        this.transactionLoading = false;
      }
    },
    async fetchBills() {
      this.billLoading = true;
      try {
        // 这里应该调用API获取账单列表
        // 模拟数据
        this.bills = [
          {
            id: 1,
            user: { username: 'user1' },
            amount: 99.00,
            status: 'paid',
            due_date: '2023-12-31 23:59:59',
            issue_date: '2023-12-01 10:00:00',
            description: '12月份服务费',
            created_at: '2023-12-01 10:00:00'
          },
          {
            id: 2,
            user: { username: 'user2' },
            amount: 299.00,
            status: 'unpaid',
            due_date: '2024-01-05 23:59:59',
            issue_date: '2023-12-06 10:00:00',
            description: '1月份服务费',
            created_at: '2023-12-06 10:00:00'
          }
        ];
        this.billPagination.total = 2;
      } catch (error) {
        this.$message.error('获取账单列表失败');
        console.error(error);
      } finally {
        this.billLoading = false;
      }
    },
    getOrderTypeText(type) {
      switch (type) {
        case 'product_order': return '产品订单';
        case 'renewal_order': return '续费订单';
        case 'bill_payment': return '账单支付';
        default: return type;
      }
    },
    getOrderTypeTag(type) {
      switch (type) {
        case 'product_order': return 'primary';
        case 'renewal_order': return 'warning';
        case 'bill_payment': return 'success';
        default: return 'info';
      }
    },
    getTransactionStatusType(status) {
      switch (status) {
        case 'completed': return 'success';
        case 'pending': return 'warning';
        case 'failed': return 'danger';
        case 'refunded': return 'info';
        default: return 'info';
      }
    },
    getTransactionStatusText(status) {
      switch (status) {
        case 'completed': return '完成';
        case 'pending': return '待处理';
        case 'failed': return '失败';
        case 'refunded': return '已退款';
        default: return status;
      }
    },
    getBillStatusType(status) {
      switch (status) {
        case 'paid': return 'success';
        case 'unpaid': return 'warning';
        case 'overdue': return 'danger';
        default: return 'info';
      }
    },
    getBillStatusText(status) {
      switch (status) {
        case 'paid': return '已支付';
        case 'unpaid': return '未支付';
        case 'overdue': return '已逾期';
        default: return status;
      }
    },
    createBill() {
      this.billForm = {
        id: null,
        user_id: '',
        amount: 0,
        due_date: '',
        description: ''
      };
      this.billDialogTitle = '创建账单';
      this.billDialogVisible = true;
    },
    editBill(bill) {
      this.billForm = { ...bill };
      this.billDialogTitle = '编辑账单';
      this.billDialogVisible = true;
    },
    viewBill(bill) {
      this.$alert(`
        <div><strong>账单ID:</strong> ${bill.id}</div>
        <div><strong>用户名:</strong> ${bill.user.username}</div>
        <div><strong>金额:</strong> ¥${bill.amount}</div>
        <div><strong>状态:</strong> ${this.getBillStatusText(bill.status)}</div>
        <div><strong>到期时间:</strong> ${bill.due_date}</div>
        <div><strong>开票时间:</strong> ${bill.issue_date}</div>
        <div><strong>描述:</strong> ${bill.description || '无'}</div>
        <div><strong>创建时间:</strong> ${bill.created_at}</div>
      `, '账单详情', {
        dangerouslyUseHTMLString: true
      });
    },
    payBill(bill) {
      this.$confirm(`确定要将账单 "${bill.id}" 标记为已支付吗？`, '确认支付', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async () => {
        try {
          await financeAPI.payBill(bill.id, 'manual');
          this.$message.success('账单已标记为已支付');
          this.fetchBills();
        } catch (error) {
          this.$message.error('支付失败');
          console.error(error);
        }
      }).catch(() => {
        // 用户取消支付
      });
    },
    async submitBillForm() {
      this.$refs.billForm.validate(async (valid) => {
        if (valid) {
          try {
            if (this.billForm.id) {
              // 编辑账单
              await financeAPI.updateBill(this.billForm.id, this.billForm);
              this.$message.success('账单更新成功');
            } else {
              // 新增账单
              await financeAPI.createBill(this.billForm);
              this.$message.success('账单创建成功');
            }
            this.billDialogVisible = false;
            this.fetchBills();
          } catch (error) {
            this.$message.error('操作失败');
            console.error(error);
          }
        }
      });
    },
    handleTransactionSizeChange(val) {
      this.transactionPagination.limit = val;
      this.fetchTransactions();
    },
    handleTransactionCurrentChange(val) {
      this.transactionPagination.page = val;
      this.fetchTransactions();
    },
    handleBillSizeChange(val) {
      this.billPagination.limit = val;
      this.fetchBills();
    },
    handleBillCurrentChange(val) {
      this.billPagination.page = val;
      this.fetchBills();
    }
  }
};
</script>

<style scoped>
.search-form {
  margin-bottom: 20px;
}
.income {
  color: #67C23A;
}
.expense {
  color: #F56C6C;
}
</style>