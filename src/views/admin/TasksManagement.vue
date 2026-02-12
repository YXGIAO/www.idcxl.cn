<template>
  <div class="tasks-management">
    <el-card>
      <template #header>
        <div class="clearfix">
          <span>任务管理</span>
          <el-button style="float: right; padding: 3px 0" type="text" @click="showAddDialog">添加任务</el-button>
        </div>
      </template>
      
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="任务状态">
          <el-select v-model="searchForm.status" placeholder="选择状态">
            <el-option label="全部" value=""></el-option>
            <el-option label="待执行" value="pending"></el-option>
            <el-option label="运行中" value="running"></el-option>
            <el-option label="已完成" value="completed"></el-option>
            <el-option label="失败" value="failed"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="任务类型">
          <el-select v-model="searchForm.task_type" placeholder="选择类型">
            <el-option label="全部" value=""></el-option>
            <el-option label="同步" value="sync"></el-option>
            <el-option label="备份" value="backup"></el-option>
            <el-option label="清理" value="cleanup"></el-option>
            <el-option label="其他" value="other"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchTasks">查询</el-button>
        </el-form-item>
      </el-form>
      
      <el-table :data="tasks" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80"></el-table-column>
        <el-table-column prop="name" label="任务名称"></el-table-column>
        <el-table-column prop="type" label="任务类型" width="120">
          <template #default="scope">
            <el-tag :type="getTypeTag(scope.row.type)">
              {{ getTypeText(scope.row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="schedule" label="执行计划" width="150"></el-table-column>
        <el-table-column prop="last_run" label="上次执行" width="160"></el-table-column>
        <el-table-column prop="next_run" label="下次执行" width="160"></el-table-column>
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
            <el-button size="mini" @click="viewTask(scope.row)">查看</el-button>
            <el-button size="mini" type="primary" @click="showEditDialog(scope.row)">编辑</el-button>
            <el-button size="mini" type="warning" @click="runTask(scope.row)">立即执行</el-button>
            <el-button size="mini" type="danger" @click="deleteTask(scope.row)">删除</el-button>
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

    <!-- 任务编辑对话框 -->
    <el-dialog :title="dialogTitle" v-model="dialogVisible" width="50%" :before-close="handleDialogClose">
      <el-form :model="taskForm" :rules="taskRules" ref="taskForm" label-width="120px">
        <el-form-item label="任务名称" prop="name">
          <el-input v-model="taskForm.name"></el-input>
        </el-form-item>
        <el-form-item label="任务类型" prop="type">
          <el-select v-model="taskForm.type" placeholder="选择任务类型">
            <el-option label="同步" value="sync"></el-option>
            <el-option label="备份" value="backup"></el-option>
            <el-option label="清理" value="cleanup"></el-option>
            <el-option label="其他" value="other"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="执行计划" prop="schedule">
          <el-input v-model="taskForm.schedule" placeholder="例如：0 * * * * (每小时执行)"></el-input>
        </el-form-item>
        <el-form-item label="任务描述" prop="description">
          <el-input v-model="taskForm.description" type="textarea"></el-input>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="taskForm.status" placeholder="选择状态">
            <el-option label="启用" value="active"></el-option>
            <el-option label="禁用" value="inactive"></el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="cancelDialog">取 消</el-button>
          <el-button type="primary" @click="submitTaskForm">确 定</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { taskAPI } from '@/services/api';
import { ElMessage, ElMessageBox } from 'element-plus';

export default {
  name: 'TasksManagement',
  data() {
    return {
      loading: false,
      tasks: [],
      searchForm: {
        status: '',
        task_type: ''
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
      taskForm: {
        id: null,
        name: '',
        type: 'sync',
        schedule: '',
        description: '',
        status: 'active'
      },
      taskRules: {
        name: [
          { required: true, message: '请输入任务名称', trigger: 'blur' }
        ],
        type: [
          { required: true, message: '请选择任务类型', trigger: 'change' }
        ],
        schedule: [
          { required: true, message: '请输入执行计划', trigger: 'blur' }
        ]
      }
    };
  },
  mounted() {
    this.fetchTasks();
  },
  methods: {
    // 重构对话框管理方法
    showAddDialog() {
      this.resetForm();
      this.dialogType = 'add';
      this.dialogTitle = '添加任务';
      this.dialogVisible = true;
    },
    
    showEditDialog(row) {
      this.taskForm = { ...row };
      this.dialogType = 'edit';
      this.dialogTitle = '编辑任务';
      this.dialogVisible = true;
    },
    
    handleDialogClose() {
      this.cancelDialog();
    },
    
    cancelDialog() {
      this.dialogVisible = false;
      this.$nextTick(() => {
        if (this.$refs.taskForm) {
          this.$refs.taskForm.clearValidate();
        }
      });
    },
    
    resetForm() {
      this.taskForm = {
        id: null,
        name: '',
        type: 'sync',
        schedule: '',
        description: '',
        status: 'active'
      };
      if (this.$refs.taskForm) {
        this.$refs.taskForm.clearValidate();
      }
    },
    
    // 标准化数据格式
    normalizeFormData(form) {
      return {
        name: form.name,
        type: form.type,
        schedule: form.schedule,
        description: form.description,
        status: form.status
      };
    },
    
    async fetchTasks() {
      this.loading = true;
      try {
        const params = {
          page: this.pagination.page,
          limit: this.pagination.limit,
          ...this.searchForm
        };
        const response = await taskAPI.getTasks(params);
        
        // 处理响应数据
        if (response && response.tasks) {
          this.tasks = response.tasks.map(item => ({
            ...item,
            id: item.id,
            name: item.name || '',
            type: item.type || 'sync',
            schedule: item.schedule || '',
            description: item.description || '',
            status: item.status || 'active',
            last_run: item.last_run || '从未执行',
            next_run: item.next_run || '未安排',
            created_at: item.created_at || ''
          }));
          this.pagination.total = response.total || response.tasks.length;
        } else {
          this.tasks = [];
          this.pagination.total = 0;
        }
      } catch (error) {
        console.error('获取任务列表失败:', error);
        ElMessage.error('获取任务列表失败: ' + (error.message || '网络错误'));
      } finally {
        this.loading = false;
      }
    },
    async deleteTask(row) {
      ElMessageBox.confirm('确定要删除该任务吗?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async () => {
        try {
          await taskAPI.deleteTask(row.id);
          ElMessage.success('任务删除成功');
          this.fetchTasks();
        } catch (error) {
          console.error('删除任务失败:', error);
          ElMessage.error('删除任务失败: ' + (error.message || '网络错误'));
        }
      }).catch(() => {});
    },
    async runTask(row) {
      ElMessageBox.confirm(`确定要立即执行任务 "${row.name}" 吗？`, '确认执行', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async () => {
        try {
          await taskAPI.runTask(row.id);
          ElMessage.success('任务已开始执行');
          this.fetchTasks(); // 刷新列表以显示最新状态
        } catch (error) {
          console.error('执行任务失败:', error);
          ElMessage.error('执行任务失败: ' + (error.message || '网络错误'));
        }
      }).catch(() => {});
    },
    async submitTaskForm() {
      this.$refs.taskForm.validate(async (valid) => {
        if (valid) {
          try {
            // 标准化数据格式
            const normalizedData = this.normalizeFormData(this.taskForm);
            
            if (this.taskForm.id) {
              // 更新任务
              await taskAPI.updateTask(this.taskForm.id, normalizedData);
              ElMessage.success('任务更新成功');
            } else {
              // 创建任务
              await taskAPI.createTask(normalizedData);
              ElMessage.success('任务创建成功');
            }
            this.dialogVisible = false;
            this.fetchTasks();
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
        case 'sync': return '同步';
        case 'backup': return '备份';
        case 'cleanup': return '清理';
        case 'other': return '其他';
        default: return type;
      }
    },
    getTypeTag(type) {
      switch (type) {
        case 'sync': return 'primary';
        case 'backup': return 'success';
        case 'cleanup': return 'warning';
        case 'other': return 'info';
        default: return 'info';
      }
    },
    getStatusType(status) {
      switch (status) {
        case 'active': return 'success';
        case 'inactive': return 'info';
        case 'pending': return 'warning';
        case 'running': return 'primary';
        case 'completed': return 'success';
        case 'failed': return 'danger';
        default: return 'info';
      }
    },
    getStatusText(status) {
      switch (status) {
        case 'active': return '启用';
        case 'inactive': return '禁用';
        case 'pending': return '待执行';
        case 'running': return '运行中';
        case 'completed': return '已完成';
        case 'failed': return '失败';
        default: return status;
      }
    },
    handleSizeChange(val) {
      this.pagination.limit = val;
      this.fetchTasks();
    },
    handleCurrentChange(val) {
      this.pagination.page = val;
      this.fetchTasks();
    },
    viewTask(task) {
      this.$alert(`
        <div><strong>任务ID:</strong> ${task.id}</div>
        <div><strong>名称:</strong> ${task.name}</div>
        <div><strong>类型:</strong> ${this.getTypeText(task.type)}</div>
        <div><strong>执行计划:</strong> ${task.schedule}</div>
        <div><strong>描述:</strong> ${task.description}</div>
        <div><strong>状态:</strong> ${this.getStatusText(task.status)}</div>
        <div><strong>上次执行:</strong> ${task.last_run}</div>
        <div><strong>下次执行:</strong> ${task.next_run}</div>
        <div><strong>创建时间:</strong> ${task.created_at}</div>
        <div><strong>更新时间:</strong> ${task.updated_at}</div>
      `, '任务详情', {
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