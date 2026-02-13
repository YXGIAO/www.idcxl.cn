<template>
  <div class="servers-management">
    <el-card>
      <template #header>
        <div class="clearfix">
          <span>服务器管理</span>
          <el-button style="float: right; padding: 3px 0" type="text" @click="showAddDialog">添加服务器</el-button>
        </div>
      </template>
      
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="服务器状态">
          <el-select v-model="searchForm.status" placeholder="选择状态">
            <el-option label="全部" value=""></el-option>
            <el-option label="活跃" value="active"></el-option>
            <el-option label="非活跃" value="inactive"></el-option>
            <el-option label="维护中" value="maintenance"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="服务器类型">
          <el-select v-model="searchForm.server_type" placeholder="选择类型">
            <el-option label="全部" value=""></el-option>
            <el-option label="VPS" value="vps"></el-option>
            <el-option label="独立服务器" value="dedicated"></el-option>
            <el-option label="云服务器" value="cloud"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchServers">查询</el-button>
        </el-form-item>
      </el-form>
      
      <el-table :data="servers" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80"></el-table-column>
        <el-table-column prop="name" label="服务器名称"></el-table-column>
        <el-table-column prop="host" label="主机地址"></el-table-column>
        <el-table-column prop="port" label="端口" width="100"></el-table-column>
        <el-table-column prop="username" label="用户名"></el-table-column>
        <el-table-column prop="server_type" label="类型" width="120">
          <template #default="scope">
            <el-tag :type="getTypeTag(scope.row.server_type)">
              {{ getTypeText(scope.row.server_type) }}
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
        <el-table-column prop="location" label="位置" width="120"></el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="160"></el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="scope">
            <el-button size="mini" @click="viewServer(scope.row)">查看</el-button>
            <el-button size="mini" type="primary" @click="showEditDialog(scope.row)">编辑</el-button>
            <el-button size="mini" type="danger" @click="deleteServer(scope.row)">删除</el-button>
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

    <!-- 服务器编辑对话框 -->
    <el-dialog :title="dialogTitle" v-model="dialogVisible" width="50%" :before-close="handleDialogClose">
      <el-form :model="serverForm" :rules="serverRules" ref="serverForm" label-width="120px">
        <el-form-item label="服务器名称" prop="name">
          <el-input v-model="serverForm.name"></el-input>
        </el-form-item>
        <el-form-item label="主机地址" prop="host">
          <el-input v-model="serverForm.host"></el-input>
        </el-form-item>
        <el-form-item label="端口" prop="port">
          <el-input-number v-model="serverForm.port" :min="1" :max="65535"></el-input-number>
        </el-form-item>
        <el-form-item label="用户名" prop="username">
          <el-input v-model="serverForm.username"></el-input>
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="serverForm.password" type="password"></el-input>
        </el-form-item>
        <el-form-item label="服务器类型" prop="server_type">
          <el-select v-model="serverForm.server_type" placeholder="选择服务器类型">
            <el-option label="VPS" value="vps"></el-option>
            <el-option label="独立服务器" value="dedicated"></el-option>
            <el-option label="云服务器" value="cloud"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="位置" prop="location">
          <el-input v-model="serverForm.location"></el-input>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="serverForm.status" placeholder="选择状态">
            <el-option label="活跃" value="active"></el-option>
            <el-option label="非活跃" value="inactive"></el-option>
            <el-option label="维护中" value="maintenance"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="serverForm.description" type="textarea"></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="cancelDialog">取 消</el-button>
          <el-button type="primary" @click="submitServerForm">确 定</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { serverAPI } from '@/services/api';
import { ElMessage, ElMessageBox } from 'element-plus';

export default {
  name: 'ServersManagement',
  data() {
    return {
      loading: false,
      servers: [],
      searchForm: {
        status: '',
        server_type: ''
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
      serverForm: {
        id: null,
        name: '',
        host: '',
        port: 22,
        username: '',
        password: '',
        server_type: 'vps',
        location: '',
        status: 'active',
        description: ''
      },
      serverRules: {
        name: [
          { required: true, message: '请输入服务器名称', trigger: 'blur' }
        ],
        host: [
          { required: true, message: '请输入主机地址', trigger: 'blur' }
        ],
        port: [
          { required: true, message: '请输入端口号', trigger: 'blur' },
          { type: 'number', min: 1, max: 65535, message: '端口号必须在1-65535之间', trigger: 'blur' }
        ],
        username: [
          { required: true, message: '请输入用户名', trigger: 'blur' }
        ]
      }
    };
  },
  mounted() {
    this.fetchServers();
  },
  methods: {
    // 重构对话框管理方法
    showAddDialog() {
      this.resetForm();
      this.dialogType = 'add';
      this.dialogTitle = '添加服务器';
      this.dialogVisible = true;
    },
    
    showEditDialog(row) {
      this.serverForm = { ...row };
      // 不显示密码，需要重新输入
      this.serverForm.password = '';
      this.dialogType = 'edit';
      this.dialogTitle = '编辑服务器';
      this.dialogVisible = true;
    },
    
    handleDialogClose() {
      this.cancelDialog();
    },
    
    cancelDialog() {
      this.dialogVisible = false;
      this.$nextTick(() => {
        if (this.$refs.serverForm) {
          this.$refs.serverForm.clearValidate();
        }
      });
    },
    
    resetForm() {
      this.serverForm = {
        id: null,
        name: '',
        host: '',
        port: 22,
        username: '',
        password: '',
        server_type: 'vps',
        location: '',
        status: 'active',
        description: ''
      };
      if (this.$refs.serverForm) {
        this.$refs.serverForm.clearValidate();
      }
    },
    
    // 标准化数据格式
    normalizeFormData(form) {
      return {
        name: form.name,
        host: form.host,
        port: Number(form.port),
        username: form.username,
        password: form.password,
        type: form.server_type, // 注意：后端可能使用type字段而不是server_type
        location: form.location,
        status: form.status,
        description: form.description,
        notes: form.notes || form.description // 可能后端使用notes字段
      };
    },
    
    async fetchServers() {
      this.loading = true;
      try {
        const params = {
          page: this.pagination.page,
          limit: this.pagination.limit,
          ...this.searchForm
        };
        const response = await serverAPI.getServers(params);
        
        // 处理响应数据
        if (response && response.servers) {
          this.servers = response.servers.map(item => ({
            ...item,
            id: item.id,
            name: item.name || '',
            host: item.host || '',
            port: Number(item.port) || 22,
            username: item.username || '',
            server_type: item.type || item.server_type || 'vps', // 处理可能的字段名差异
            location: item.location || '',
            status: item.status || 'active',
            description: item.description || item.notes || '',
            created_at: item.created_at || ''
          }));
          this.pagination.total = response.total || response.servers.length;
        } else {
          this.servers = [];
          this.pagination.total = 0;
        }
      } catch (error) {
        console.error('获取服务器列表失败:', error);
        ElMessage.error('获取服务器列表失败: ' + (error.message || '网络错误'));
      } finally {
        this.loading = false;
      }
    },
    async deleteServer(row) {
      ElMessageBox.confirm('确定要删除该服务器吗?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async () => {
        try {
          await serverAPI.deleteServer(row.id);
          ElMessage.success('服务器删除成功');
          this.fetchServers();
        } catch (error) {
          console.error('删除服务器失败:', error);
          ElMessage.error('删除服务器失败: ' + (error.message || '网络错误'));
        }
      }).catch(() => {});
    },
    async submitServerForm() {
      this.$refs.serverForm.validate(async (valid) => {
        if (valid) {
          try {
            // 标准化数据格式
            const normalizedData = this.normalizeFormData(this.serverForm);
            
            if (this.serverForm.id) {
              // 更新服务器
              await serverAPI.updateServer(this.serverForm.id, normalizedData);
              ElMessage.success('服务器更新成功');
            } else {
              // 创建服务器
              await serverAPI.createServer(normalizedData);
              ElMessage.success('服务器创建成功');
            }
            this.dialogVisible = false;
            this.fetchServers();
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
    getTypeText(type) {
      switch (type) {
        case 'vps': return 'VPS';
        case 'dedicated': return '独立服务器';
        case 'cloud': return '云服务器';
        default: return type;
      }
    },
    getTypeTag(type) {
      switch (type) {
        case 'vps': return 'primary';
        case 'dedicated': return 'warning';
        case 'cloud': return 'success';
        default: return 'info';
      }
    },
    getStatusType(status) {
      switch (status) {
        case 'active': return 'success';
        case 'inactive': return 'info';
        case 'maintenance': return 'warning';
        default: return 'info';
      }
    },
    getStatusText(status) {
      switch (status) {
        case 'active': return '活跃';
        case 'inactive': return '非活跃';
        case 'maintenance': return '维护中';
        default: return status;
      }
    },
    handleSizeChange(val) {
      this.pagination.limit = val;
      this.fetchServers();
    },
    handleCurrentChange(val) {
      this.pagination.page = val;
      this.fetchServers();
    },
    viewServer(server) {
      this.$alert(`
        <div><strong>服务器ID:</strong> ${server.id}</div>
        <div><strong>名称:</strong> ${server.name}</div>
        <div><strong>主机地址:</strong> ${server.host}</div>
        <div><strong>端口:</strong> ${server.port}</div>
        <div><strong>用户名:</strong> ${server.username}</div>
        <div><strong>类型:</strong> ${this.getTypeText(server.server_type)}</div>
        <div><strong>状态:</strong> ${this.getStatusText(server.status)}</div>
        <div><strong>位置:</strong> ${server.location}</div>
        <div><strong>描述:</strong> ${server.description}</div>
        <div><strong>创建时间:</strong> ${server.created_at}</div>
        <div><strong>更新时间:</strong> ${server.updated_at}</div>
      `, '服务器详情', {
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