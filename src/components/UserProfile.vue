<template>
  <div class="user-profile">
    <el-container>
      <el-header class="header">
        <h2>用户账户管理</h2>
      </el-header>

      <el-main class="main-content">
        <el-row :gutter="20">
          <el-col :span="8">
            <div class="profile-sidebar">
              <div class="avatar-section">
                <div class="avatar-upload-container">
                  <el-upload
                    class="avatar-uploader"
                    action="#"
                    :show-file-list="false"
                    :http-request="customUpload"
                    :before-upload="beforeAvatarUpload"
                  >
                    <img v-if="user.avatar" :src="user.avatar" class="avatar" />
                    <i v-else class="el-icon-user-solid avatar-uploader-icon"></i>
                  </el-upload>
                  <input 
                    type="file" 
                    id="avatar-input" 
                    accept="image/*" 
                    @change="onFileSelected"
                    ref="avatarInput"
                    style="display: none"
                  />
                  <div class="avatar-actions">
                    <el-button size="small" type="primary" @click="$refs.avatarInput.click()">
                      选择图片
                    </el-button>
                    <el-button 
                      size="small" 
                      type="success" 
                      :disabled="!tempAvatar"
                      @click="uploadAvatar"
                    >
                      {{ uploading ? '上传中...' : '保存头像' }}
                    </el-button>
                  </div>
                  <div v-if="uploadMessage" class="upload-message" :class="uploadMessageType">
                    {{ uploadMessage }}
                  </div>
                </div>
                <p class="username">{{ user.username || '未设置用户名' }}</p>
                <p class="user-id">ID: {{ user.id || '---' }}</p>
              </div>

              <el-menu
                class="profile-menu"
                :default-active="activeMenu"
                @select="handleMenuSelect"
              >
                <el-menu-item index="profile">个人资料</el-menu-item>
                <el-menu-item index="security">安全设置</el-menu-item>
                <el-menu-item index="orders">我的订单</el-menu-item>
                <el-menu-item index="billing">账单信息</el-menu-item>
                <el-menu-item index="tickets">我的工单</el-menu-item>
              </el-menu>
            </div>
          </el-col>

          <el-col :span="16">
            <div class="profile-content">
              <!-- 个人资料部分 -->
              <div v-if="activeMenu === 'profile'" class="profile-section">
                <h3>个人资料</h3>
                <el-form :model="user" label-width="100px" class="profile-form">
                  <el-form-item label="用户名">
                    <el-input v-model="user.username" :disabled="!editingProfile"></el-input>
                  </el-form-item>
                  <el-form-item label="邮箱">
                    <el-input v-model="user.email" :disabled="!editingProfile"></el-input>
                  </el-form-item>
                  <el-form-item label="手机号">
                    <el-input v-model="user.phone" :disabled="!editingProfile"></el-input>
                  </el-form-item>
                  <el-form-item label="真实姓名">
                    <el-input v-model="user.realName" :disabled="!editingProfile"></el-input>
                  </el-form-item>
                  <el-form-item label="身份证号">
                    <el-input v-model="user.idCard" :disabled="!editingProfile"></el-input>
                  </el-form-item>
                  <el-form-item label="实名认证">
                    <el-tag 
                      :type="user.verified ? 'success' : 'warning'"
                      size="small"
                    >
                      {{ user.verified ? '已认证' : '未认证' }}
                    </el-tag>
                  </el-form-item>
                  <el-form-item>
                    <el-button 
                      v-if="!editingProfile" 
                      type="primary" 
                      @click="editingProfile = true"
                    >
                      编辑资料
                    </el-button>
                    <div v-else>
                      <el-button type="primary" @click="saveProfile">保存</el-button>
                      <el-button @click="cancelEdit">取消</el-button>
                    </div>
                  </el-form-item>
                </el-form>
              </div>

              <!-- 安全设置部分 -->
              <div v-if="activeMenu === 'security'" class="security-section">
                <h3>安全设置</h3>
                <el-form label-width="120px" class="security-form">
                  <el-form-item label="登录密码">
                    <el-input 
                      v-model="passwordForm.currentPassword" 
                      type="password" 
                      placeholder="当前密码"
                    ></el-input>
                    <el-input 
                      v-model="passwordForm.newPassword" 
                      type="password" 
                      placeholder="新密码"
                      style="margin-top: 10px;"
                    ></el-input>
                    <el-input 
                      v-model="passwordForm.confirmPassword" 
                      type="password" 
                      placeholder="确认新密码"
                      style="margin-top: 10px;"
                    ></el-input>
                    <el-button 
                      type="primary" 
                      @click="changePassword" 
                      style="margin-top: 10px;"
                    >
                      修改密码
                    </el-button>
                  </el-form-item>
                  
                  <el-form-item label="绑定邮箱">
                    <div class="security-item">
                      <span>{{ user.email || '未绑定' }}</span>
                      <el-button size="small" type="text">修改</el-button>
                    </div>
                  </el-form-item>
                  
                  <el-form-item label="绑定手机">
                    <div class="security-item">
                      <span>{{ user.phone || '未绑定' }}</span>
                      <el-button size="small" type="text">修改</el-button>
                    </div>
                  </el-form-item>
                </el-form>
              </div>

              <!-- 我的订单部分 -->
              <div v-if="activeMenu === 'orders'" class="orders-section">
                <h3>我的订单</h3>
                <el-tabs v-model="orderTab" type="card">
                  <el-tab-pane label="产品订单" name="product">
                    <el-table :data="productOrders" style="width: 100%">
                      <el-table-column prop="id" label="订单ID" width="100"></el-table-column>
                      <el-table-column prop="productName" label="产品名称"></el-table-column>
                      <el-table-column prop="amount" label="金额" width="100">
                        <template #default="scope">¥{{ scope.row.amount }}</template>
                      </el-table-column>
                      <el-table-column prop="status" label="状态" width="100">
                        <template #default="scope">
                          <el-tag 
                            :type="getStatusType(scope.row.status)"
                            size="small"
                          >
                            {{ scope.row.status }}
                          </el-tag>
                        </template>
                      </el-table-column>
                      <el-table-column label="操作" width="150">
                        <template slot-scope="scope">
                          <el-button size="mini" @click="viewOrder(scope.row)">查看</el-button>
                        </template>
                      </el-table-column>
                    </el-table>
                    <el-pagination
                      @size-change="handleSizeChange"
                      @current-change="handleCurrentChange"
                      :current-page="currentPage"
                      :page-sizes="[5, 10, 20]"
                      :page-size="pageSize"
                      layout="total, sizes, prev, pager, next, jumper"
                      :total="totalOrders">
                    </el-pagination>
                  </el-tab-pane>
                  <el-tab-pane label="续费订单" name="renewal">
                    <el-table :data="renewalOrders" style="width: 100%">
                      <el-table-column prop="id" label="订单ID" width="100"></el-table-column>
                      <el-table-column prop="productName" label="产品名称"></el-table-column>
                      <el-table-column prop="amount" label="金额" width="100">
                        <template #default="scope">¥{{ scope.row.amount }}</template>
                      </el-table-column>
                      <el-table-column prop="status" label="状态" width="100">
                        <template #default="scope">
                          <el-tag 
                            :type="getStatusType(scope.row.status)"
                            size="small"
                          >
                            {{ scope.row.status }}
                          </el-tag>
                        </template>
                      </el-table-column>
                      <el-table-column label="操作" width="150">
                        <template slot-scope="scope">
                          <el-button size="mini" @click="viewOrder(scope.row)">查看</el-button>
                        </template>
                      </el-table-column>
                    </el-table>
                  </el-tab-pane>
                </el-tabs>
              </div>

              <!-- 账单信息部分 -->
              <div v-if="activeMenu === 'billing'" class="billing-section">
                <h3>账单信息</h3>
                <el-table :data="bills" style="width: 100%">
                  <el-table-column prop="id" label="账单ID" width="100"></el-table-column>
                  <el-table-column prop="description" label="描述"></el-table-column>
                  <el-table-column prop="amount" label="金额" width="100">
                    <template #default="scope">¥{{ scope.row.amount }}</template>
                  </el-table-column>
                  <el-table-column prop="status" label="状态" width="100">
                    <template #default="scope">
                      <el-tag 
                        :type="scope.row.status === 'paid' ? 'success' : 'warning'"
                        size="small"
                      >
                        {{ scope.row.status === 'paid' ? '已支付' : '待支付' }}
                      </el-tag>
                    </template>
                  </el-table-column>
                  <el-table-column prop="createdAt" label="创建时间" width="150"></el-table-column>
                  <el-table-column label="操作" width="150">
                    <template slot-scope="scope">
                      <el-button 
                        size="mini" 
                        @click="payBill(scope.row)" 
                        :disabled="scope.row.status === 'paid'"
                      >
                        {{ scope.row.status === 'paid' ? '已支付' : '支付' }}
                      </el-button>
                    </template>
                  </el-table-column>
                </el-table>
              </div>

              <!-- 我的工单部分 -->
              <div v-if="activeMenu === 'tickets'" class="tickets-section">
                <h3>我的工单</h3>
                <el-table :data="tickets" style="width: 100%">
                  <el-table-column prop="id" label="工单ID" width="100"></el-table-column>
                  <el-table-column prop="title" label="标题"></el-table-column>
                  <el-table-column prop="category" label="分类" width="100"></el-table-column>
                  <el-table-column prop="status" label="状态" width="100">
                    <template #default="scope">
                      <el-tag 
                        :type="getTicketStatusType(scope.row.status)"
                        size="small"
                      >
                        {{ scope.row.status }}
                      </el-tag>
                    </template>
                  </el-table-column>
                  <el-table-column prop="createdAt" label="创建时间" width="150"></el-table-column>
                  <el-table-column label="操作" width="150">
                    <template slot-scope="scope">
                      <el-button size="mini" @click="viewTicket(scope.row)">查看</el-button>
                    </template>
                  </el-table-column>
                </el-table>
              </div>
            </div>
          </el-col>
        </el-row>
      </el-main>
    </el-container>
  </div>
</template>

<script>
import { userAPI, orderAPI, financeAPI, ticketAPI, currentUserAPI } from '@/services/api.js';

export default {
  name: 'UserProfile',
  data() {
    return {
      activeMenu: 'profile',
      editingProfile: false,
      orderTab: 'product',
      user: {
        id: null,
        username: '',
        email: '',
        phone: '',
        realName: '',
        idCard: '',
        verified: false,
        avatar: ''
      },
      passwordForm: {
        currentPassword: '',
        newPassword: '',
        confirmPassword: ''
      },
      tempAvatar: null,
      uploading: false,
      uploadMessage: '',
      uploadMessageType: 'info', // 'info', 'success', 'error'
      productOrders: [],
      renewalOrders: [],
      bills: [],
      tickets: [],
      currentPage: 1,
      pageSize: 5,
      totalOrders: 0
    };
  },
  methods: {
    async loadUserProfile() {
      try {
        // 使用新的API获取当前用户信息
        const response = await currentUserAPI.getCurrentUser();
        const userData = response.data || response;
        
        if(userData) {
          this.user = {
            id: userData.id,
            username: userData.username,
            email: userData.email,
            phone: userData.phone || '',
            realName: userData.name || '',
            idCard: '', // 身份证信息需要从用户资料中获取
            verified: userData.real_name_auth || false,
            avatar: userData.avatar || 'https://via.placeholder.com/100x100/409eff/white?text=头像'
          };
          
          // 如果需要获取更多用户资料，可以调用用户资料API
          if(userData.id) {
            try {
              const profileResponse = await userAPI.getUserProfile(userData.id);
              const profileData = profileResponse.data || profileResponse;
              if(profileData) {
                this.user.realName = profileData.real_name || this.user.realName;
                this.user.idCard = profileData.id_card || '';
                this.user.verified = profileData.auth_status === 'approved';
              }
            } catch (profileError) {
              console.log('无法获取用户资料，使用基本信息');
            }
          }
        }
      } catch (error) {
        console.error('加载用户信息失败:', error);
        this.$message.error('加载用户信息失败');
      }
    },
    async loadOrders() {
      try {
        // 模拟加载订单信息
        this.productOrders = [
          { id: 1, productName: '云服务器套餐A', amount: 299.00, status: 'completed' },
          { id: 2, productName: '云服务器套餐B', amount: 599.00, status: 'pending' },
          { id: 3, productName: '云服务器套餐C', amount: 199.00, status: 'cancelled' }
        ];
        
        this.renewalOrders = [
          { id: 101, productName: '云服务器套餐A续费', amount: 269.10, status: 'completed' },
          { id: 102, productName: '云服务器套餐B续费', amount: 539.10, status: 'processing' }
        ];
      } catch (error) {
        console.error('加载订单信息失败:', error);
        this.$message.error('加载订单信息失败');
      }
    },
    async loadBills() {
      try {
        // 模拟加载账单信息
        this.bills = [
          { id: 1, description: '云服务器套餐A月费', amount: 299.00, status: 'paid', createdAt: '2023-05-01' },
          { id: 2, description: '云服务器套餐B月费', amount: 599.00, status: 'pending', createdAt: '2023-06-01' },
          { id: 3, description: '云服务器套餐C月费', amount: 199.00, status: 'pending', createdAt: '2023-06-15' }
        ];
      } catch (error) {
        console.error('加载账单信息失败:', error);
        this.$message.error('加载账单信息失败');
      }
    },
    async loadTickets() {
      try {
        // 模拟加载工单信息
        this.tickets = [
          { id: 1, title: '服务器连接问题', category: '技术', status: 'open', createdAt: '2023-05-10' },
          { id: 2, title: '账单疑问', category: '财务', status: 'closed', createdAt: '2023-05-15' },
          { id: 3, title: '产品功能建议', category: '产品', status: 'processing', createdAt: '2023-06-01' }
        ];
      } catch (error) {
        console.error('加载工单信息失败:', error);
        this.$message.error('加载工单信息失败');
      }
    },
    handleMenuSelect(index) {
      this.activeMenu = index;
      
      if (index === 'orders') {
        this.loadOrders();
      } else if (index === 'billing') {
        this.loadBills();
      } else if (index === 'tickets') {
        this.loadTickets();
      }
    },
    onFileSelected(event) {
      const file = event.target.files[0];
      if (file) {
        // 验证文件类型
        if (!file.type.match('image.*')) {
          this.uploadMessage = '请选择图片文件';
          this.uploadMessageType = 'error';
          return;
        }

        // 验证文件大小 (最大2MB)
        if (file.size > 2 * 1024 * 1024) {
          this.uploadMessage = '图片大小不能超过 2MB!';
          this.uploadMessageType = 'error';
          return;
        }

        // 显示预览
        const reader = new FileReader();
        reader.onload = (e) => {
          this.user.avatar = e.target.result;
          this.tempAvatar = file;
        };
        reader.readAsDataURL(file);

        this.uploadMessage = '已选择文件，点击保存上传';
        this.uploadMessageType = 'info';
      }
    },
    beforeAvatarUpload(file) {
      // 这个方法现在只做验证，不处理上传
      const isImage = file.type.match('image.*');
      const isLt2M = file.size / 1024 / 1024 < 2;

      if (!isImage) {
        this.$message.error('头像图片只能是 JPG/PNG 格式!');
      } else if (!isLt2M) {
        this.$message.error('头像图片大小不能超过 2MB!');
      }
      
      // 如果验证通过，我们仍然返回false，因为我们使用自定义上传方法
      return isImage && isLt2M;
    },
    customUpload(options) {
      // 使用我们自己的上传逻辑而不是默认的HTTP请求
      this.$refs.avatarInput.value = null;
      const input = document.createElement('input');
      input.type = 'file';
      input.accept = 'image/*';
      input.onchange = this.onFileSelected;
      input.click();
    },
    async uploadAvatar() {
      if (!this.tempAvatar) {
        this.uploadMessage = '请先选择头像';
        this.uploadMessageType = 'error';
        return;
      }

      this.uploading = true;
      try {
        // 模拟上传过程
        await new Promise(resolve => setTimeout(resolve, 1500));
        
        // 这里应该是实际的API调用
        
        this.uploadMessage = '头像更新成功！';
        this.uploadMessageType = 'success';
        
        // 清空临时头像状态
        this.tempAvatar = null;
        
        // 实际项目中，这里会收到服务器返回的真实URL
        // this.user.avatar = response.data.avatarUrl;
      } catch (error) {
        console.error('上传头像失败:', error);
        this.uploadMessage = '上传头像失败，请重试';
        this.uploadMessageType = 'error';
      } finally {
        this.uploading = false;
      }
    },
    saveProfile() {
      // 模拟保存资料
      this.$message.success('资料保存成功');
      this.editingProfile = false;
    },
    cancelEdit() {
      this.editingProfile = false;
    },
    changePassword() {
      if (!this.passwordForm.currentPassword) {
        this.$message.error('请输入当前密码');
        return;
      }
      if (!this.passwordForm.newPassword) {
        this.$message.error('请输入新密码');
        return;
      }
      if (this.passwordForm.newPassword !== this.passwordForm.confirmPassword) {
        this.$message.error('两次输入的密码不一致');
        return;
      }
      
      // 模拟修改密码
      this.$message.success('密码修改成功');
      this.passwordForm = {
        currentPassword: '',
        newPassword: '',
        confirmPassword: ''
      };
    },
    getStatusType(status) {
      switch(status) {
        case 'completed': return 'success';
        case 'pending': return 'warning';
        case 'processing': return 'primary';
        case 'cancelled': return 'info';
        default: return 'info';
      }
    },
    getTicketStatusType(status) {
      switch(status) {
        case 'open': return 'danger';
        case 'processing': return 'warning';
        case 'closed': return 'success';
        case 'resolved': return 'primary';
        default: return 'info';
      }
    },
    viewOrder(order) {
      this.$message.info(`查看订单 ${order.id}`);
    },
    payBill(bill) {
      this.$message.success(`支付账单 ${bill.id} 成功`);
      bill.status = 'paid';
    },
    viewTicket(ticket) {
      this.$message.info(`查看工单 ${ticket.id}`);
    },
    handleSizeChange(val) {
      this.pageSize = val;
    },
    handleCurrentChange(val) {
      this.currentPage = val;
    }
  },
  mounted() {
    this.loadUserProfile();
    this.loadOrders();
    this.loadBills();
    this.loadTickets();
  }
};
</script>

<style scoped>
.user-profile {
  background-color: #f5f7fa;
  min-height: 100vh;
  padding: 20px;
}

.header {
  background-color: #fff;
  padding: 20px;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
  margin-bottom: 20px;
}

.main-content {
  padding: 0;
}

.profile-sidebar {
  background-color: #fff;
  padding: 20px;
  border-radius: 4px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  height: fit-content;
}

.avatar-section {
  text-align: center;
  margin-bottom: 30px;
}

.avatar-uploader {
  text-align: center;
}

.avatar-uploader .el-upload {
  border: 1px dashed #d9d9d9;
  border-radius: 50%;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  margin: 0 auto 15px;
  transition: all 0.3s;
}

.avatar-uploader .el-upload:hover {
  border-color: #409EFF;
  transform: scale(1.05);
}

.avatar-actions {
  display: flex;
  justify-content: center;
  gap: 10px;
  margin-top: 15px;
}

.upload-message {
  margin-top: 10px;
  padding: 8px 12px;
  border-radius: 4px;
  font-size: 14px;
  text-align: center;
}

.upload-message.info {
  background-color: #e6f7ff;
  color: #1890ff;
  border: 1px solid #91d5ff;
}

.upload-message.success {
  background-color: #f6ffed;
  color: #52c41a;
  border: 1px solid #b7eb8f;
}

.upload-message.error {
  background-color: #fff2f0;
  color: #ff4d4f;
  border: 1px solid #ffccc7;
}

.avatar-uploader-icon {
  font-size: 28px;
  color: #8c939d;
  width: 100px;
  height: 100px;
  line-height: 100px;
  text-align: center;
}

.avatar {
  width: 100px;
  height: 100px;
  display: block;
  border-radius: 50%;
}

.username {
  font-size: 18px;
  font-weight: bold;
  margin: 10px 0 5px 0;
}

.user-id {
  color: #909399;
  margin: 0 0 15px 0;
}

.profile-menu {
  border-right: none;
}

.profile-content {
  background-color: #fff;
  padding: 20px;
  border-radius: 4px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.profile-section, 
.security-section, 
.orders-section, 
.billing-section, 
.tickets-section {
  padding: 10px 0;
}

.profile-form, 
.security-form {
  max-width: 500px;
  margin-top: 20px;
}

.security-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid #eee;
}

.security-item:last-child {
  border-bottom: none;
}

.el-pagination {
  margin-top: 20px;
  text-align: right;
}
</style>