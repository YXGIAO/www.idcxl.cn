<template>
  <div class="tickets-management">
    <el-tabs v-model="activeTab">
      <el-tab-pane label="工单列表" name="list">
        <el-card>
          <template #header>
            <div class="clearfix">
              <span>工单列表</span>
            </div>
          </template>
          
          <el-form :inline="true" :model="searchForm" class="search-form">
            <el-form-item label="工单状态">
              <el-select v-model="searchForm.status" placeholder="选择状态">
                <el-option label="全部" value=""></el-option>
                <el-option label="开启" value="open"></el-option>
                <el-option label="处理中" value="in_progress"></el-option>
                <el-option label="已解决" value="resolved"></el-option>
                <el-option label="已关闭" value="closed"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item label="工单分类">
              <el-select v-model="searchForm.category" placeholder="选择分类">
                <el-option label="全部" value=""></el-option>
                <el-option label="技术问题" value="technical"></el-option>
                <el-option label="账务问题" value="billing"></el-option>
                <el-option label="销售咨询" value="sales"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item label="优先级">
              <el-select v-model="searchForm.priority" placeholder="选择优先级">
                <el-option label="全部" value=""></el-option>
                <el-option label="低" value="low"></el-option>
                <el-option label="中" value="medium"></el-option>
                <el-option label="高" value="high"></el-option>
                <el-option label="紧急" value="urgent"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item label="用户ID">
              <el-input v-model="searchForm.user_id" placeholder="用户ID"></el-input>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="fetchTickets">查询</el-button>
            </el-form-item>
          </el-form>
          
          <el-table :data="tickets" style="width: 100%" v-loading="loading">
            <el-table-column prop="id" label="ID" width="80"></el-table-column>
            <el-table-column prop="user.username" label="提交用户" width="120"></el-table-column>
            <el-table-column prop="title" label="标题"></el-table-column>
            <el-table-column prop="category" label="分类" width="100">
              <template #default="scope">
                <el-tag :type="getCategoryType(scope.row.category)">
                  {{ getCategoryText(scope.row.category) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template #default="scope">
                <el-tag :type="getStatusType(scope.row.status)">
                  {{ getStatusText(scope.row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="priority" label="优先级" width="100">
              <template #default="scope">
                <el-tag :type="getPriorityType(scope.row.priority)">
                  {{ getPriorityText(scope.row.priority) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="提交时间" width="160"></el-table-column>
            <el-table-column label="操作" width="200">
              <template #default="scope">
                <el-button size="mini" @click="viewTicket(scope.row)">查看</el-button>
                <el-button size="mini" type="primary" @click="editTicket(scope.row)">处理</el-button>
                <el-button size="mini" type="danger" @click="deleteTicket(scope.row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
          
          <el-pagination
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
            :current-page="pagination.page"
            :page-sizes="[10, 20, 50, 100]"
            :page-size="pagination.limit"
            layout="total, sizes, prev, pager, next, jumper"
            :total="pagination.total">
          </el-pagination>
        </el-card>
      </el-tab-pane>
      
      <el-tab-pane label="工单统计" name="stats">
        <el-card>
          <template #header>
            <div class="clearfix">
              <span>工单统计</span>
            </div>
          </template>
          
          <el-row :gutter="20">
            <el-col :span="6">
              <el-card class="stat-card">
                <div class="stat-number">{{ stats.total }}</div>
                <div class="stat-label">总工单数</div>
              </el-card>
            </el-col>
            <el-col :span="6">
              <el-card class="stat-card">
                <div class="stat-number">{{ stats.open }}</div>
                <div class="stat-label">开启中</div>
              </el-card>
            </el-col>
            <el-col :span="6">
              <el-card class="stat-card">
                <div class="stat-number">{{ stats.in_progress }}</div>
                <div class="stat-label">处理中</div>
              </el-card>
            </el-col>
            <el-col :span="6">
              <el-card class="stat-card">
                <div class="stat-number">{{ stats.resolved }}</div>
                <div class="stat-label">已解决</div>
              </el-card>
            </el-col>
          </el-row>
          
          <el-divider></el-divider>
          
          <div class="chart-container">
            <h3>工单趋势图</h3>
            <div id="ticket-chart" style="height: 400px;"></div>
          </div>
        </el-card>
      </el-tab-pane>
    </el-tabs>

    <!-- 工单详情对话框 -->
    <el-dialog title="工单详情" :visible.sync="ticketDialogVisible" width="60%">
      <div v-if="currentTicket">
        <el-descriptions :column="2" :border="true">
          <el-descriptions-item label="工单ID">{{ currentTicket.id }}</el-descriptions-item>
          <el-descriptions-item label="提交用户">{{ currentTicket.user?.username }}</el-descriptions-item>
          <el-descriptions-item label="标题">{{ currentTicket.title }}</el-descriptions-item>
          <el-descriptions-item label="分类">
            <el-tag :type="getCategoryType(currentTicket.category)">
              {{ getCategoryText(currentTicket.category) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="getStatusType(currentTicket.status)">
              {{ getStatusText(currentTicket.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="优先级">
            <el-tag :type="getPriorityType(currentTicket.priority)">
              {{ getPriorityText(currentTicket.priority) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="提交时间">{{ currentTicket.created_at }}</el-descriptions-item>
          <el-descriptions-item label="更新时间">{{ currentTicket.updated_at }}</el-descriptions-item>
        </el-descriptions>
        
        <div class="ticket-content">
          <h4>工单内容</h4>
          <p>{{ currentTicket.content }}</p>
        </div>
        
        <div class="ticket-actions" style="margin-top: 20px;">
          <el-button @click="updateTicketStatus('in_progress')" :disabled="currentTicket.status === 'in_progress'">标记为处理中</el-button>
          <el-button type="success" @click="updateTicketStatus('resolved')" :disabled="currentTicket.status === 'resolved'">标记为已解决</el-button>
          <el-button type="info" @click="updateTicketStatus('closed')" :disabled="currentTicket.status === 'closed'">标记为已关闭</el-button>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script>
export default {
  name: 'TicketsManagement',
  data() {
    return {
      activeTab: 'list',
      tickets: [],
      loading: false,
      searchForm: {
        status: '',
        category: '',
        priority: '',
        user_id: ''
      },
      pagination: {
        page: 1,
        limit: 10,
        total: 0
      },
      ticketDialogVisible: false,
      currentTicket: null,
      stats: {
        total: 0,
        open: 0,
        in_progress: 0,
        resolved: 0
      }
    };
  },
  mounted() {
    this.fetchTickets();
    this.fetchTicketStats();
  },
  methods: {
    async fetchTickets() {
      this.loading = true;
      try {
        // 这里应该调用API获取工单列表
        // 模拟数据
        this.tickets = [
          {
            id: 1,
            user: { username: 'user1' },
            title: '无法登录账户',
            content: '尝试登录时提示密码错误，但确认密码无误。',
            category: 'technical',
            status: 'open',
            priority: 'high',
            created_at: '2023-12-27 10:00:00',
            updated_at: '2023-12-27 10:00:00'
          },
          {
            id: 2,
            user: { username: 'user2' },
            title: '账单疑问',
            content: '对本月账单金额有疑问，希望能详细说明费用构成。',
            category: 'billing',
            status: 'in_progress',
            priority: 'medium',
            created_at: '2023-12-27 11:00:00',
            updated_at: '2023-12-27 11:30:00'
          }
        ];
        this.pagination.total = 2;
      } catch (error) {
        this.$message.error('获取工单列表失败');
        console.error(error);
      } finally {
        this.loading = false;
      }
    },
    async fetchTicketStats() {
      // 模拟获取工单统计数据
      this.stats = {
        total: 15,
        open: 5,
        in_progress: 3,
        resolved: 6,
        closed: 1
      };
    },
    getCategoryText(category) {
      switch (category) {
        case 'technical': return '技术问题';
        case 'billing': return '账务问题';
        case 'sales': return '销售咨询';
        default: return category;
      }
    },
    getCategoryType(category) {
      switch (category) {
        case 'technical': return 'warning';
        case 'billing': return 'danger';
        case 'sales': return 'primary';
        default: return 'info';
      }
    },
    getStatusText(status) {
      switch (status) {
        case 'open': return '开启';
        case 'in_progress': return '处理中';
        case 'resolved': return '已解决';
        case 'closed': return '已关闭';
        default: return status;
      }
    },
    getStatusType(status) {
      switch (status) {
        case 'open': return 'warning';
        case 'in_progress': return 'primary';
        case 'resolved': return 'success';
        case 'closed': return 'info';
        default: return 'info';
      }
    },
    getPriorityText(priority) {
      switch (priority) {
        case 'low': return '低';
        case 'medium': return '中';
        case 'high': return '高';
        case 'urgent': return '紧急';
        default: return priority;
      }
    },
    getPriorityType(priority) {
      switch (priority) {
        case 'low': return 'info';
        case 'medium': return 'primary';
        case 'high': return 'warning';
        case 'urgent': return 'danger';
        default: return 'info';
      }
    },
    viewTicket(ticket) {
      this.currentTicket = { ...ticket };
      this.ticketDialogVisible = true;
    },
    editTicket(ticket) {
      this.currentTicket = { ...ticket };
      this.ticketDialogVisible = true;
    },
    deleteTicket(ticket) {
      this.$confirm(`确定要删除工单 "${ticket.title}" 吗？`, '确认删除', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        // 这里应该调用API删除工单
        this.$message.success('工单删除成功');
        this.fetchTickets();
      }).catch(() => {
        // 用户取消删除
      });
    },
    updateTicketStatus(newStatus) {
      // 这里应该调用API更新工单状态
      this.currentTicket.status = newStatus;
      this.currentTicket.updated_at = new Date().toISOString().slice(0, 19).replace('T', ' ');
      this.$message.success('工单状态已更新');
    },
    handleSizeChange(val) {
      this.pagination.limit = val;
      this.fetchTickets();
    },
    handleCurrentChange(val) {
      this.pagination.page = val;
      this.fetchTickets();
    }
  }
};
</script>

<style scoped>
.search-form {
  margin-bottom: 20px;
}
.stat-card {
  text-align: center;
}
.stat-number {
  font-size: 24px;
  font-weight: bold;
  color: #409EFF;
}
.stat-label {
  margin-top: 10px;
  color: #909399;
}
.ticket-content {
  margin-top: 20px;
  padding: 15px;
  background-color: #f5f7fa;
  border-radius: 4px;
}
.ticket-actions {
  margin-top: 20px;
}
.chart-container {
  margin-top: 20px;
}
</style>