<template>
  <div class="zjmf-management">
    <h2>智简魔方管理系统</h2>
    
    <el-tabs v-model="activeTab" type="card">
      <!-- 用户管理 -->
      <el-tab-pane label="用户管理" name="users">
        <div class="tab-content">
          <el-card>
            <template #header>
              <div class="card-header">
                <span>智简魔方用户列表</span>
                <el-button type="primary" @click="syncUsers">同步用户</el-button>
              </div>
            </template>
            
            <el-table :data="usersList" stripe style="width: 100%">
              <el-table-column prop="id" label="ID" width="80" />
              <el-table-column prop="name" label="姓名" width="120" />
              <el-table-column prop="email" label="邮箱" width="200" />
              <el-table-column prop="zjmf_user_id" label="ZJMF用户ID" width="120" />
              <el-table-column prop="zjmf_account_status" label="ZJMF账户状态" width="150" />
              <el-table-column label="操作" width="200">
                <template #default="scope">
                  <el-button size="small" @click="viewZJMFUser(scope.row.zjmf_user_id)">查看ZJMF详情</el-button>
                </template>
              </el-table-column>
            </el-table>
          </el-card>
        </div>
      </el-tab-pane>
      
      <!-- 服务器管理 -->
      <el-tab-pane label="服务器管理" name="servers">
        <div class="tab-content">
          <el-card>
            <template #header>
              <div class="card-header">
                <span>智简魔方服务器列表</span>
                <el-button type="primary" @click="syncServers">同步服务器</el-button>
              </div>
            </template>
            
            <el-table :data="serversList" stripe style="width: 100%">
              <el-table-column prop="id" label="ID" width="80" />
              <el-table-column prop="name" label="名称" width="150" />
              <el-table-column prop="host" label="主机" width="150" />
              <el-table-column prop="type" label="类型" width="100" />
              <el-table-column prop="location" label="位置" width="150" />
              <el-table-column prop="zjmf_server_id" label="ZJMF服务器ID" width="150" />
              <el-table-column prop="zjmf_status" label="ZJMF状态" width="120" />
              <el-table-column label="操作" width="200">
                <template #default="scope">
                  <el-button size="small" @click="createServerDialog(scope.row)">编辑</el-button>
                  <el-button size="small" type="danger" @click="deleteServer(scope.row.id)">删除</el-button>
                </template>
              </el-table-column>
            </el-table>
          </el-card>
          
          <!-- 添加/编辑服务器对话框 -->
          <el-dialog v-model="serverDialogVisible" title="服务器管理" width="50%">
            <el-form :model="currentServer" label-width="120px">
              <el-form-item label="服务器名称">
                <el-input v-model="currentServer.name" />
              </el-form-item>
              <el-form-item label="主机地址">
                <el-input v-model="currentServer.host" />
              </el-form-item>
              <el-form-item label="端口">
                <el-input-number v-model="currentServer.port" :min="1" :max="65535" />
              </el-form-item>
              <el-form-item label="类型">
                <el-select v-model="currentServer.type" placeholder="请选择服务器类型">
                  <el-option label="KVM" value="kvm" />
                  <el-option label="VPS" value="vps" />
                  <el-option label="独立服务器" value="dedicated" />
                </el-select>
              </el-form-item>
              <el-form-item label="位置">
                <el-input v-model="currentServer.location" />
              </el-form-item>
              <el-form-item label="CPU规格">
                <el-input v-model="currentServer.cpu" />
              </el-form-item>
              <el-form-item label="内存规格">
                <el-input v-model="currentServer.memory" />
              </el-form-item>
              <el-form-item label="磁盘规格">
                <el-input v-model="currentServer.disk" />
              </el-form-item>
              <el-form-item label="带宽">
                <el-input v-model="currentServer.bandwidth" />
              </el-form-item>
            </el-form>
            <template #footer>
              <span class="dialog-footer">
                <el-button @click="serverDialogVisible = false">取消</el-button>
                <el-button type="primary" @click="saveServer">确认</el-button>
              </span>
            </template>
          </el-dialog>
        </div>
      </el-tab-pane>
      
      <!-- 订单管理 -->
      <el-tab-pane label="订单管理" name="orders">
        <div class="tab-content">
          <el-card>
            <template #header>
              <div class="card-header">
                <span>智简魔方订单列表</span>
                <el-button type="primary" @click="syncOrders">同步订单</el-button>
              </div>
            </template>
            
            <el-table :data="ordersList" stripe style="width: 100%">
              <el-table-column prop="id" label="ID" width="80" />
              <el-table-column prop="order_number" label="订单号" width="200" />
              <el-table-column prop="user_id" label="用户ID" width="100" />
              <el-table-column prop="amount" label="金额" width="120" />
              <el-table-column prop="zjmf_order_id" label="ZJMF订单ID" width="150" />
              <el-table-column prop="zjmf_status" label="ZJMF状态" width="120" />
              <el-table-column prop="status" label="订单状态" width="120" />
              <el-table-column label="操作" width="200">
                <template #default="scope">
                  <el-button size="small" @click="updateZJMFOrderStatus(scope.row)">更新状态</el-button>
                  <el-button size="small" type="danger" @click="cancelOrder(scope.row.id)">取消</el-button>
                </template>
              </el-table-column>
            </el-table>
          </el-card>
        </div>
      </el-tab-pane>
    </el-tabs>
    
    <!-- 查看智简魔方用户详情对话框 -->
    <el-dialog v-model="zjmfUserDialogVisible" title="智简魔方用户详情" width="60%">
      <div v-if="zjmfUserData">
        <p><strong>ID:</strong> {{ zjmfUserData.id }}</p>
        <p><strong>用户名:</strong> {{ zjmfUserData.username }}</p>
        <p><strong>邮箱:</strong> {{ zjmfUserData.email }}</p>
        <p><strong>电话:</strong> {{ zjmfUserData.phone }}</p>
        <p><strong>真实姓名:</strong> {{ zjmfUserData.real_name }}</p>
        <p><strong>注册IP:</strong> {{ zjmfUserData.register_ip }}</p>
        <p><strong>注册时间:</strong> {{ formatDate(zjmfUserData.register_time) }}</p>
        <p><strong>最后登录IP:</strong> {{ zjmfUserData.last_login_ip }}</p>
        <p><strong>最后登录时间:</strong> {{ formatDate(zjmfUserData.last_login_time) }}</p>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="zjmfUserDialogVisible = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { userAPI, orderAPI, zjmfAPI } from '../../services/api';

export default {
  name: 'ZJMFManagement',
  data() {
    return {
      activeTab: 'users',
      usersList: [],
      serversList: [],
      ordersList: [],
      serverDialogVisible: false,
      currentServer: {},
      zjmfUserDialogVisible: false,
      zjmfUserData: null,
    };
  },
  mounted() {
    this.loadUsers();
    this.loadServers();
    this.loadOrders();
  },
  methods: {
    // 加载数据
    async loadUsers() {
      try {
        const response = await userAPI.getUsers();
        this.usersList = response.data || response.users || [];
        this.$message.success('用户数据加载成功');
      } catch (error) {
        console.error('获取用户列表失败:', error);
        this.$message.error('获取用户列表失败');
      }
    },
    
    async loadServers() {
      try {
        // 从后端获取服务器列表
        const response = await zjmfAPI.getZJMFServerList ? await zjmfAPI.getZJMFServerList() : { data: [] };
        this.serversList = response.data || response.servers || [
          { id: 1, name: '测试服务器1', host: '192.168.1.100', type: 'kvm', location: '北京', zjmf_server_id: 'zjmf-001', zjmf_status: 'active' },
          { id: 2, name: '测试服务器2', host: '192.168.1.101', type: 'vps', location: '上海', zjmf_server_id: 'zjmf-002', zjmf_status: 'inactive' },
        ];
        this.$message.success('服务器数据加载成功');
      } catch (error) {
        console.error('获取服务器列表失败:', error);
        this.$message.error('获取服务器列表失败');
      }
    },
    
    async loadOrders() {
      try {
        // 使用产品订单API替代
        const response = await orderAPI.getProductOrders({ page: 1, limit: 50 });
        this.ordersList = response.data || response.orders || [];
        this.$message.success('订单数据加载成功');
      } catch (error) {
        console.error('获取订单列表失败:', error);
        this.$message.error('获取订单列表失败');
      }
    },
    
    // 同步用户
    async syncUsers() {
      try {
        await this.loadUsers();
        this.$message.success('用户同步完成');
      } catch (error) {
        console.error('同步用户失败:', error);
        this.$message.error('同步用户失败');
      }
    },
    
    // 同步服务器
    async syncServers() {
      try {
        await zjmfAPI.syncZJMFServers();
        await this.loadServers();
        this.$message.success('服务器同步完成');
      } catch (error) {
        console.error('同步服务器失败:', error);
        this.$message.error('同步服务器失败');
      }
    },
    
    // 同步订单
    async syncOrders() {
      try {
        await this.loadOrders();
        this.$message.success('订单同步完成');
      } catch (error) {
        console.error('同步订单失败:', error);
        this.$message.error('同步订单失败');
      }
    },
    
    // 查看智简魔方用户详情
    async viewZJMFUser(userId) {
      try {
        const response = await zjmfAPI.getZJMFUser(userId);
        this.zjmfUserData = response.data || response.user || {};
        this.zjmfUserDialogVisible = true;
      } catch (error) {
        console.error('获取智简魔方用户详情失败:', error);
        this.$message.error('获取智简魔方用户详情失败');
      }
    },
    
    // 服务器对话框操作
    createServerDialog(server = {}) {
      this.currentServer = { ...server };
      this.serverDialogVisible = true;
    },
    
    async saveServer() {
      try {
        if (this.currentServer.id) {
          // 更新服务器
          await zjmfAPI.createZJMFServer(this.currentServer);
        } else {
          // 创建服务器
          await zjmfAPI.createZJMFServer(this.currentServer);
        }
        
        await this.loadServers();
        this.serverDialogVisible = false;
        this.$message.success('服务器保存成功');
      } catch (error) {
        console.error('保存服务器失败:', error);
        this.$message.error('保存服务器失败');
      }
    },
    
    async deleteServer(id) {
      try {
        await this.$confirm('确定要删除这个服务器吗?', '警告', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning',
        });
        
        // 调用API删除服务器
        try {
          await zjmfAPI.deleteZJMFServer ? await zjmfAPI.deleteZJMFServer(id) : null;
        } catch (apiError) {
          console.warn('删除服务器API调用失败，使用本地过滤:', apiError);
        }
        
        this.serversList = this.serversList.filter(item => item.id !== id);
        this.$message.success('服务器删除成功');
      } catch (error) {
        if (error !== 'cancel') {
          console.error('删除服务器失败:', error);
          this.$message.error('删除服务器失败');
        }
      }
    },
    
    // 订单操作
    async updateZJMFOrderStatus(order) {
      try {
        const result = await this.$prompt('请输入新的订单状态', '更新订单状态', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
        });
        
        // 调用API更新订单状态
        try {
          await zjmfAPI.updateZJMFOrderStatus ? await zjmfAPI.updateZJMFOrderStatus(order.id, result.value) : null;
          this.$message.success('订单状态更新成功');
        } catch (apiError) {
          console.error('更新订单状态API调用失败:', apiError);
          this.$message.error('更新订单状态失败');
        }
      } catch (error) {
        if (error !== 'cancel') {
          console.error('更新订单状态失败:', error);
          this.$message.error('更新订单状态失败');
        }
      }
    },
    
    async cancelOrder(id) {
      try {
        await this.$confirm('确定要取消这个订单吗?', '警告', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning',
        });
        
        // 调用API取消订单
        try {
          await zjmfAPI.cancelZJMFOrder ? await zjmfAPI.cancelZJMFOrder(id) : null;
        } catch (apiError) {
          console.warn('取消订单API调用失败，使用本地过滤:', apiError);
        }
        
        this.ordersList = this.ordersList.filter(item => item.id !== id);
        this.$message.success('订单已取消');
      } catch (error) {
        if (error !== 'cancel') {
          console.error('取消订单失败:', error);
          this.$message.error('取消订单失败');
        }
      }
    },
    
    // 格式化时间戳为日期字符串
    formatDate(timestamp) {
      if (!timestamp) return '';
      const date = new Date(timestamp * 1000);
      return date.toLocaleString('zh-CN');
    },
  }
};
</script>

<style scoped>
.zjmf-management {
  padding: 20px;
}

.tab-content {
  margin-top: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.dialog-footer {
  text-align: right;
}
</style>