<template>
  <div class="users-management">
    <el-card>
      <template #header>
        <div class="clearfix">
          <span>用户管理</span>
          <el-button style="float: right; padding: 3px 0" type="text" @click="showAddDialog">添加用户</el-button>
        </div>
      </template>
      
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="用户状态">
          <el-select v-model="searchForm.status" placeholder="选择状态">
            <el-option label="全部" value=""></el-option>
            <el-option label="启用" value="1"></el-option>
            <el-option label="禁用" value="0"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="用户角色">
          <el-select v-model="searchForm.role" placeholder="选择角色">
            <el-option label="全部" value=""></el-option>
            <el-option label="普通用户" value="user"></el-option>
            <el-option label="管理员" value="admin"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchUsers">查询</el-button>
        </el-form-item>
      </el-form>
      
      <el-table :data="users" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80"></el-table-column>
        <el-table-column prop="username" label="用户名" width="120"></el-table-column>
        <el-table-column prop="name" label="真实姓名" width="120"></el-table-column>
        <el-table-column prop="email" label="邮箱" width="200"></el-table-column>
        <el-table-column prop="role" label="角色" width="100">
          <template #default="scope">
            <el-tag :type="getRoleType(scope.row.role)">
              {{ getRoleText(scope.row.role) }}
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
        <el-table-column prop="created_at" label="创建时间" width="160"></el-table-column>
        <el-table-column label="操作" width="250">
          <template #default="scope">
            <el-button size="mini" @click="viewUser(scope.row)">查看</el-button>
            <el-button size="mini" type="primary" @click="showEditDialog(scope.row)">编辑</el-button>
            <el-button size="mini" type="warning" @click="resetPassword(scope.row)">重置密码</el-button>
            <el-button size="mini" type="danger" @click="deleteUser(scope.row)">删除</el-button>
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

    <!-- 用户编辑对话框 -->
    <el-dialog :title="dialogTitle" v-model="dialogVisible" width="50%" :before-close="handleDialogClose">
      <el-form :model="userForm" :rules="userRules" ref="userForm" label-width="120px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="userForm.username" :disabled="dialogType === 'edit'"></el-input>
        </el-form-item>
        <el-form-item label="真实姓名" prop="name">
          <el-input v-model="userForm.name"></el-input>
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="userForm.email"></el-input>
        </el-form-item>
        <el-form-item label="密码" :prop="dialogType === 'add' ? 'password' : ''">
          <el-input v-model="userForm.password" type="password" :placeholder="dialogType === 'add' ? '请输入密码' : '留空则不修改密码'"></el-input>
        </el-form-item>
        <el-form-item label="角色" prop="role">
          <el-select v-model="userForm.role" placeholder="选择角色">
            <el-option label="普通用户" value="user"></el-option>
            <el-option label="管理员" value="admin"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="userForm.status" placeholder="选择状态">
            <el-option label="启用" value="1"></el-option>
            <el-option label="禁用" value="0"></el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="cancelDialog">取 消</el-button>
          <el-button type="primary" @click="submitUserForm">确 定</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { userAPI } from '@/services/api';
import { ElMessage, ElMessageBox } from 'element-plus';

export default {
  name: 'UsersManagement',
  data() {
    return {
      loading: false,
      users: [],
      searchForm: {
        status: '',
        role: ''
      },
      pagination: {
        page: 1,
        limit: 10,
        total: 0
      },
      // 重构对话框管理逻辑
      dialogVisible: false,
      dialogType: '', // 'add' 或 'edit'
      dialogTitle: '',
      userForm: {
        id: null,
        username: '',
        name: '',
        email: '',
        password: '',
        role: 'user',
        status: 1
      },
      userRules: {
        username: [
          { required: true, message: '请输入用户名', trigger: 'blur' },
          { min: 3, max: 20, message: '用户名长度应在3-20个字符之间', trigger: 'blur' }
        ],
        name: [
          { required: true, message: '请输入真实姓名', trigger: 'blur' }
        ],
        email: [
          { required: true, message: '请输入邮箱地址', trigger: 'blur' },
          { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
        ],
        password: [
          { required: true, message: '请输入密码', trigger: 'blur' },
          { min: 6, message: '密码长度至少6位', trigger: 'blur' }
        ]
      }
    };
  },
  mounted() {
    this.fetchUsers();
  },
  methods: {
    // 重构对话框管理方法
    showAddDialog() {
      this.resetForm();
      this.dialogType = 'add';
      this.dialogTitle = '添加用户';
      this.dialogVisible = true;
    },
    
    showEditDialog(row) {
      // 创建用户数据副本
      this.userForm = {
        ...row,
        password: '', // 编辑时不显示原密码
        status: String(row.status) // 确保状态为字符串格式
      };
      this.dialogType = 'edit';
      this.dialogTitle = '编辑用户';
      this.dialogVisible = true;
    },
    
    handleDialogClose() {
      this.cancelDialog();
    },
    
    cancelDialog() {
      this.dialogVisible = false;
      this.$nextTick(() => {
        if (this.$refs.userForm) {
          this.$refs.userForm.clearValidate();
        }
      });
    },
    
    resetForm() {
      this.userForm = {
        id: null,
        username: '',
        name: '',
        email: '',
        password: '',
        role: 'user',
        status: 1
      };
      if (this.$refs.userForm) {
        this.$refs.userForm.clearValidate();
      }
    },
    
    // 标准化数据格式
    normalizeFormData(form) {
      const normalized = {
        username: form.username,
        name: form.name,
        email: form.email,
        role: form.role,
        status: Number(form.status) // 确保状态为数字格式
      };
      
      // 只有在添加用户或密码不为空时才包含密码
      if (this.dialogType === 'add' || form.password) {
        normalized.password = form.password;
      }
      
      return normalized;
    },
    
    async fetchUsers() {
      this.loading = true;
      try {
        const params = {
          page: this.pagination.page,
          limit: this.pagination.limit,
          ...this.searchForm
        };
        const response = await userAPI.getUsers(params);
        
        // 处理响应数据
        if (response && response.users) {
          this.users = response.users.map(item => ({
            ...item,
            id: item.id,
            username: item.username || '',
            name: item.name || '',
            email: item.email || '',
            role: item.role || 'user',
            status: Number(item.status) || 1, // 确保状态为数字
            created_at: item.created_at || ''
          }));
          this.pagination.total = response.total || response.users.length;
        } else {
          this.users = [];
          this.pagination.total = 0;
        }
      } catch (error) {
        console.error('获取用户列表失败:', error);
        ElMessage.error('获取用户列表失败: ' + (error.message || '网络错误'));
      } finally {
        this.loading = false;
      }
    },
    async deleteUser(row) {
      ElMessageBox.confirm(`确定要删除用户 "${row.username}" 吗？`, '确认删除', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async () => {
        try {
          await userAPI.deleteUser(row.id);
          ElMessage.success('用户删除成功');
          this.fetchUsers();
        } catch (error) {
          console.error('删除用户失败:', error);
          ElMessage.error('删除用户失败: ' + (error.message || '网络错误'));
        }
      }).catch(() => {});
    },
    async resetPassword(row) {
      ElMessageBox.prompt('请输入新密码', '重置密码', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        inputType: 'password',
        inputPattern: /^.{6,}$/,
        inputErrorMessage: '密码长度至少6位'
      }).then(async ({ value }) => {
        try {
          await userAPI.updateUser(row.id, { password: value });
          ElMessage.success('密码重置成功');
        } catch (error) {
          console.error('重置密码失败:', error);
          ElMessage.error('密码重置失败: ' + (error.message || '网络错误'));
        }
      }).catch(() => {});
    },
    async submitUserForm() {
      this.$refs.userForm.validate(async (valid) => {
        if (valid) {
          try {
            // 标准化数据格式
            const normalizedData = this.normalizeFormData(this.userForm);
            
            if (this.userForm.id) {
              // 更新用户
              await userAPI.updateUser(this.userForm.id, normalizedData);
              ElMessage.success('用户更新成功');
            } else {
              // 创建用户
              await userAPI.createUser(normalizedData);
              ElMessage.success('用户创建成功');
            }
            this.dialogVisible = false;
            this.fetchUsers();
          } catch (error) {
            console.error('操作失败:', error);
            let errorMessage = '操作失败';
            if (error.response) {
              errorMessage += `: ${error.response.status} - ${error.response.data?.message || JSON.stringify(error.response.data)}`;
            } else {
              errorMessage += `: ${error.message}`;
            }
            ElMessage.error(errorMessage);
          }
        }
      });
    },
    getRoleText(role) {
      switch (role) {
        case 'admin': return '管理员';
        case 'user': return '普通用户';
        default: return role;
      }
    },
    getRoleType(role) {
      switch (role) {
        case 'admin': return 'danger';
        case 'user': return 'primary';
        default: return 'info';
      }
    },
    getStatusType(status) {
      switch (Number(status)) {
        case 1: return 'success';
        case 0: return 'info';
        default: return 'info';
      }
    },
    getStatusText(status) {
      switch (Number(status)) {
        case 1: return '启用';
        case 0: return '禁用';
        default: return status;
      }
    },
    handleSizeChange(val) {
      this.pagination.limit = val;
      this.fetchUsers();
    },
    handleCurrentChange(val) {
      this.pagination.page = val;
      this.fetchUsers();
    },
    viewUser(user) {
      this.$alert(`
        <div><strong>用户ID:</strong> ${user.id}</div>
        <div><strong>用户名:</strong> ${user.username}</div>
        <div><strong>真实姓名:</strong> ${user.name}</div>
        <div><strong>邮箱:</strong> ${user.email}</div>
        <div><strong>角色:</strong> ${this.getRoleText(user.role)}</div>
        <div><strong>状态:</strong> ${this.getStatusText(user.status)}</div>
        <div><strong>创建时间:</strong> ${user.created_at}</div>
        <div><strong>更新时间:</strong> ${user.updated_at}</div>
      `, '用户详情', {
        dangerouslyUseHTMLString: true
      });
    }
  }
};
</script>

<style scoped>
.search-form {
  margin-bottom: 20px;
}
</style>