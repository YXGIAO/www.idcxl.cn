<template>
  <div class="suppliers-management">
    <el-card>
      <template #header>
        <div class="clearfix">
          <span>智简魔方供应商管理</span>
          <el-button type="primary" size="small" style="float: right; margin-left: 10px;" @click="showAutoConfigDialog">
            自动配置
          </el-button>
          <el-button type="success" size="small" style="float: right;" @click="showAddDialog">
            添加供应商
          </el-button>
        </div>
      </template>
      
      <!-- 搜索表单 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="供应商状态">
          <el-select v-model="searchForm.status" placeholder="选择状态" clearable>
            <el-option label="全部" value=""></el-option>
            <el-option label="活跃" value="active"></el-option>
            <el-option label="非活跃" value="inactive"></el-option>
            <el-option label="暂停" value="suspended"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="接口类型">
          <el-select v-model="searchForm.interface_type" placeholder="选择接口类型" clearable>
            <el-option label="全部" value=""></el-option>
            <el-option label="手动" value="manual"></el-option>
            <el-option label="智简魔方" value="zjmf"></el-option>
            <el-option label="v10" value="v10"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchSuppliers">查询</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
      
      <!-- 供应商表格 -->
      <el-table :data="suppliers" style="width: 100%" v-loading="loading" border>
        <el-table-column prop="id" label="ID" width="80" sortable></el-table-column>
        <el-table-column prop="name" label="供应商名称" min-width="150">
          <template #default="scope">
            <span style="font-weight: bold;">{{ scope.row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="interface_type" label="接口类型" width="120" align="center">
          <template #default="scope">
            <el-tag :type="getInterfaceTypeTagType(scope.row.interface_type)">
              {{ getInterfaceTypeText(scope.row.interface_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="zjmf_account" label="账户" width="150">
          <template #default="scope">
            <span v-if="scope.row.zjmf_account">{{ scope.row.zjmf_account }}</span>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="contact" label="联系方式" width="150">
          <template #default="scope">
            <span v-if="scope.row.contact">{{ scope.row.contact }}</span>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="balance" label="账户余额" width="120" align="right">
          <template #default="scope">
            <span style="color: #e6a23c; font-weight: bold;">¥{{ (scope.row.balance || 0).toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="product_stats" label="产品统计" width="150">
          <template #default="scope">
            <div>
              <div>正常: {{ scope.row.normal_products_count || 0 }}</div>
              <div>总计: {{ scope.row.total_products_count || 0 }}</div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)" size="small">
              {{ getStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="scope">
            <el-button size="small" type="primary" @click="viewSupplier(scope.row)">查看</el-button>
            <el-button size="small" @click="editSupplier(scope.row)">编辑</el-button>
            <el-button 
              v-if="scope.row.interface_type === 'zjmf'" 
              size="small" 
              type="success" 
              @click="syncSupplierInfo(scope.row)"
            >
              同步信息
            </el-button>
            <el-button size="small" type="danger" @click="deleteSupplier(scope.row)">删除</el-button>
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

    <!-- 添加/编辑供应商对话框 -->
    <el-dialog
      :title="dialogTitle"
      v-model="dialogVisible"
      width="600px"
      :before-close="handleDialogClose"
    >
      <el-form :model="supplierForm" :rules="supplierRules" ref="supplierForm" label-width="120px">
        <!-- 基础信息部分 -->
        <el-divider content-position="left">基础信息</el-divider>
        <el-form-item label="供应商名称" prop="name">
          <el-input v-model="supplierForm.name" placeholder="请输入供应商名称"></el-input>
        </el-form-item>
        
        <el-form-item label="联系方式" prop="contact">
          <el-input v-model="supplierForm.contact" placeholder="请输入联系方式（可选）"></el-input>
        </el-form-item>
        
        <el-form-item label="备注" prop="description">
          <el-input 
            v-model="supplierForm.description" 
            type="textarea" 
            :rows="3"
            placeholder="请输入备注信息（可选）"
          ></el-input>
        </el-form-item>
        
        <!-- 自动化配置部分 -->
        <el-divider content-position="left">自动化配置</el-divider>
        <el-form-item label="接口类型" prop="interface_type">
          <el-select v-model="supplierForm.interface_type" placeholder="请选择接口类型">
            <el-option label="手动" value="manual"></el-option>
            <el-option label="智简魔方" value="zjmf"></el-option>
            <el-option label="v10" value="v10"></el-option>
          </el-select>
        </el-form-item>
        
        <el-form-item label="接口地址" prop="interface_url">
          <el-input v-model="supplierForm.interface_url" placeholder="请输入接口地址"></el-input>
        </el-form-item>
        
        <el-form-item label="用户名" prop="username">
          <el-input v-model="supplierForm.username" placeholder="请输入用户名"></el-input>
        </el-form-item>
        
        <el-form-item label="API密钥" prop="api_key">
          <el-input 
            v-model="supplierForm.api_key" 
            type="password" 
            show-password 
            placeholder="请输入API密钥"
          ></el-input>
        </el-form-item>
        
        <el-form-item label="状态" prop="status">
          <el-select v-model="supplierForm.status" placeholder="选择状态">
            <el-option label="活跃" value="active"></el-option>
            <el-option label="非活跃" value="inactive"></el-option>
            <el-option label="暂停" value="suspended"></el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="cancelDialog">取消</el-button>
          <el-button type="primary" @click="submitSupplierForm">确定</el-button>
        </div>
      </template>
    </el-dialog>
    
    <!-- 自动配置对话框 -->
    <el-dialog
      title="自动配置供应商"
      v-model="autoConfigDialogVisible"
      width="500px"
    >
      <el-steps :active="autoConfigStep" finish-status="success" align-center>
        <el-step title="选择接口类型"></el-step>
        <el-step title="填写配置信息"></el-step>
        <el-step title="测试连接"></el-step>
      </el-steps>
      
      <div style="margin-top: 30px;">
        <!-- 第一步：选择接口类型 -->
        <div v-if="autoConfigStep === 0">
          <el-alert 
            title="请选择要配置的接口类型" 
            type="info" 
            show-icon 
            :closable="false"
            style="margin-bottom: 20px;"
          ></el-alert>
          <el-radio-group v-model="autoConfigForm.interface_type" @change="handleInterfaceTypeChange">
            <el-radio label="zjmf">智简魔方 IDC</el-radio>
            <el-radio label="v10">v10 接口</el-radio>
          </el-radio-group>
        </div>
        
        <!-- 第二步：填写配置信息 -->
        <div v-if="autoConfigStep === 1">
          <el-form :model="autoConfigForm" label-width="100px">
            <el-form-item label="接口地址" required>
              <el-input v-model="autoConfigForm.interface_url" placeholder="请输入接口地址"></el-input>
            </el-form-item>
            <el-form-item label="用户名" required>
              <el-input v-model="autoConfigForm.username" placeholder="请输入用户名"></el-input>
            </el-form-item>
            <el-form-item label="API密钥" required>
              <el-input 
                v-model="autoConfigForm.api_key" 
                type="password" 
                show-password 
                placeholder="请输入API密钥"
              ></el-input>
            </el-form-item>
            <el-form-item label="供应商名称">
              <el-input v-model="autoConfigForm.name" placeholder="请输入供应商名称（可选）"></el-input>
            </el-form-item>
          </el-form>
        </div>
        
        <!-- 第三步：测试连接 -->
        <div v-if="autoConfigStep === 2">
          <el-result 
            :icon="testResult.icon" 
            :title="testResult.title" 
            :sub-title="testResult.subtitle"
          >
            <template #extra>
              <el-button v-if="testResult.success" type="primary" @click="confirmAutoConfig">
                确认添加
              </el-button>
              <el-button v-else @click="autoConfigStep = 1">重新配置</el-button>
            </template>
          </el-result>
        </div>
      </div>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="autoConfigDialogVisible = false">取消</el-button>
          <el-button 
            v-if="autoConfigStep > 0 && autoConfigStep < 2" 
            @click="autoConfigStep--"
          >
            上一步
          </el-button>
          <el-button 
            v-if="autoConfigStep < 1" 
            type="primary" 
            @click="autoConfigStep++"
            :disabled="!autoConfigForm.interface_type"
          >
            下一步
          </el-button>
          <el-button 
            v-if="autoConfigStep === 1" 
            type="primary" 
            @click="testConnection"
            :loading="testLoading"
          >
            测试连接
          </el-button>
        </div>
      </template>
    </el-dialog>
    
    <!-- 供应商详情对话框 -->
    <el-dialog
      title="供应商详情"
      v-model="detailDialogVisible"
      width="600px"
    >
      <el-descriptions :column="1" border v-if="currentSupplier">
        <el-descriptions-item label="ID">{{ currentSupplier.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ currentSupplier.name }}</el-descriptions-item>
        <el-descriptions-item label="接口类型">
          <el-tag :type="getInterfaceTypeTagType(currentSupplier.interface_type)">
            {{ getInterfaceTypeText(currentSupplier.interface_type) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="接口地址">{{ currentSupplier.interface_url || '-' }}</el-descriptions-item>
        <el-descriptions-item label="用户名">{{ currentSupplier.username || '-' }}</el-descriptions-item>
        <el-descriptions-item label="联系方式">{{ currentSupplier.contact || '-' }}</el-descriptions-item>
        <el-descriptions-item label="账户余额">¥{{ (currentSupplier.balance || 0).toFixed(2) }}</el-descriptions-item>
        <el-descriptions-item label="产品统计">
          正常: {{ currentSupplier.normal_products_count || 0 }} / 总计: {{ currentSupplier.total_products_count || 0 }}
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(currentSupplier.status)">
            {{ getStatusText(currentSupplier.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="备注">{{ currentSupplier.description || '-' }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatDate(currentSupplier.created_at) }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ formatDate(currentSupplier.updated_at) }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="detailDialogVisible = false">关闭</el-button>
          <el-button type="primary" @click="editSupplier(currentSupplier)">编辑</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { supplierAPI, zjmfAPI } from '@/services/api';
import { ElMessage, ElMessageBox } from 'element-plus';

export default {
  name: 'SuppliersManagement',
  data() {
    return {
      loading: false,
      testLoading: false,
      suppliers: [],
      searchForm: {
        status: '',
        interface_type: ''
      },
      pagination: {
        page: 1,
        limit: 10,
        total: 0
      },
      // 对话框管理
      dialogVisible: false,
      dialogType: '', // 'add' 或 'edit'
      dialogTitle: '',
      detailDialogVisible: false,
      autoConfigDialogVisible: false,
      autoConfigStep: 0,
      currentSupplier: null,
      
      // 供应商表单
      supplierForm: {
        id: null,
        name: '',
        contact: '',
        description: '',
        interface_type: 'manual',
        interface_url: '',
        username: '',
        api_key: '',
        status: 'active'
      },
      
      // 自动配置表单
      autoConfigForm: {
        interface_type: '',
        interface_url: '',
        username: '',
        api_key: '',
        name: ''
      },
      
      // 测试结果
      testResult: {
        icon: 'info',
        title: '',
        subtitle: '',
        success: false
      },
      
      // 表单验证规则
      supplierRules: {
        name: [
          { required: true, message: '请输入供应商名称', trigger: 'blur' }
        ],
        interface_type: [
          { required: true, message: '请选择接口类型', trigger: 'change' }
        ],
        interface_url: [
          { required: true, message: '请输入接口地址', trigger: 'blur' }
        ],
        username: [
          { required: false, message: '请输入用户名', trigger: 'blur' }
        ],
        api_key: [
          { required: false, message: '请输入API密钥', trigger: 'blur' }
        ],
        contact: [
          { required: false, message: '请输入联系方式', trigger: 'blur' }
        ]
      }
    };
  },
  mounted() {
    this.fetchSuppliers();
  },
  methods: {
    // 对话框管理方法
    showAddDialog() {
      this.resetForm();
      this.dialogType = 'add';
      this.dialogTitle = '添加供应商';
      this.dialogVisible = true;
    },
    
    showAutoConfigDialog() {
      this.resetAutoConfigForm();
      this.autoConfigStep = 0;
      this.autoConfigDialogVisible = true;
    },
    
    editSupplier(supplier) {
      this.supplierForm = { ...supplier };
      this.dialogType = 'edit';
      this.dialogTitle = '编辑供应商';
      this.dialogVisible = true;
    },
    
    viewSupplier(supplier) {
      this.currentSupplier = supplier;
      this.detailDialogVisible = true;
    },
    
    handleDialogClose() {
      this.cancelDialog();
    },
    
    cancelDialog() {
      this.dialogVisible = false;
      this.$nextTick(() => {
        if (this.$refs.supplierForm) {
          this.$refs.supplierForm.clearValidate();
        }
      });
    },
    
    resetForm() {
      this.supplierForm = {
        id: null,
        name: '',
        contact: '',
        description: '',
        interface_type: 'manual',
        interface_url: '',
        username: '',
        api_key: '',
        status: 'active'
      };
      if (this.$refs.supplierForm) {
        this.$refs.supplierForm.clearValidate();
      }
    },
    
    resetAutoConfigForm() {
      this.autoConfigForm = {
        interface_type: '',
        interface_url: '',
        username: '',
        api_key: '',
        name: ''
      };
      this.testResult = {
        icon: 'info',
        title: '',
        subtitle: '',
        success: false
      };
    },
    
    resetSearch() {
      this.searchForm = {
        status: '',
        interface_type: ''
      };
      this.pagination.page = 1;
      this.fetchSuppliers();
    },
    
    // 标准化数据格式
    normalizeFormData(form) {
      // 创建一个新对象，确保只有后端需要的字段
      const normalized = {
        // 必填字段
        type: form.interface_type === 'zjmf' ? 'default' : 'default', // 根据接口类型设置供应商类型
        name: form.name,
        url: form.interface_url || 'http://example.com', // URL字段必填，提供默认值
        username: form.username || form.name || 'supplier_user', // 用户名必填
        contact: form.contact || 'noreply@example.com', // 联系方式必填且必须是邮箱格式
        currency_code: 'CNY', // 货币代码必填
        rate: '1.0000', // 汇率必填
        auto_update_rate: 0, // 自动更新汇率
        
        // 可选字段
        description: form.description || '',
        status: form.status === 'active' ? 1 : 0, // 将字符串状态转换为数字
        notes: form.notes || '',
        
        // 智简魔方专用字段
        zjmf_api_key: form.api_key || '',
        zjmf_api_secret: form.api_secret || form.api_key || '',
        zjmf_api_endpoint: form.interface_url || '',
        zjmf_server_id: form.zjmf_server_id || '',
        zjmf_server_group: form.zjmf_server_group || '',
        zjmf_host: form.zjmf_host || form.interface_url || '',
        zjmf_status: form.zjmf_status || 'active',
        zjmf_account: form.zjmf_account || form.username || ''
      };

      // 如果是编辑模式，添加ID
      if (form.id) {
        normalized.id = form.id;
      }

      return normalized;
    },
    
    // 获取接口类型显示文本
    getInterfaceTypeText(type) {
      switch (type) {
        case 'manual': return '手动';
        case 'zjmf': return '智简魔方';
        case 'v10': return 'v10';
        default: return type;
      }
    },
    
    // 同步上游供应商信息
    async syncSupplierInfo(supplier) {
      if (supplier.interface_type !== 'zjmf') {
        ElMessage.warning('只有智简魔方类型的供应商才能同步信息');
        return;
      }
      
      try {
        // 显示加载状态
        this.loading = true;
        ElMessage.info('正在同步供应商信息...');
        
        // 直接调用后端API进行同步
        const response = await supplierAPI.syncSupplierInfo(supplier.id);
        
        // 更新供应商信息
        const updatedSupplier = {
          ...supplier,
          balance: response.balance || response.data?.balance || 0,
          available_products_count: response.available_products_count || response.data?.available_products_count || 0,
          total_products_count: response.total_products_count || response.data?.total_products_count || 0,
          normal_products_count: response.normal_products_count || response.data?.normal_products_count || 0,
          zjmf_status: response.status || response.data?.status || 'active'
        };
        
        // 刷新列表
        await this.fetchSuppliers();
        
        ElMessage.success('供应商信息同步成功');
      } catch (error) {
        console.error('同步供应商信息失败:', error);
        ElMessage.error('同步供应商信息失败: ' + (error.message || '网络错误'));
      } finally {
        this.loading = false;
      }
    },
    
    // 从供应商同步产品信息
    async syncProductsFromSupplier(supplier) {
      if (supplier.interface_type !== 'zjmf') {
        ElMessage.warning('只有智简魔方类型的供应商才能同步产品');
        return;
      }
      
      try {
        this.loading = true;
        ElMessage.info('正在从供应商同步产品信息...');
        
        // 调用后端API同步产品
        const response = await supplierAPI.syncProductsFromSupplier(supplier.id);
        
        ElMessage.success(`成功同步 ${response.syncedCount || 0} 个产品`);
      } catch (error) {
        console.error('同步产品信息失败:', error);
        ElMessage.error('同步产品信息失败: ' + (error.message || '网络错误'));
      } finally {
        this.loading = false;
      }
    },
    
    async fetchSuppliers() {
      this.loading = true;
      try {
        const params = {
          page: this.pagination.page,
          limit: this.pagination.limit,
          ...this.searchForm
        };
        const response = await supplierAPI.getSuppliers(params);
        
        // 处理响应数据
        if (response && response.suppliers) {
          // 转换后端数据为前端格式
          this.suppliers = response.suppliers.map(item => {
            return {
              ...item,
              id: item.id,
              name: item.name || '',
              interface_type: item.interface_type || item.type || 'manual', // 默认为手动
              interface_url: item.interface_url || item.api_endpoint || '',
              zjmf_account: item.zjmf_account || item.account || '',
              status: item.status === 1 || item.status === 'active' ? 'active' : 'inactive',
              available_products_count: item.available_products_count || 0,
              assigned_products_count: item.assigned_products_count || 0,
              total_products_count: item.total_products_count || 0,
              normal_products_count: item.normal_products_count || 0,
              balance: item.balance || 0,
              description: item.description || ''
            };
          });
          this.pagination.total = response.total || response.suppliers.length;
        } else {
          this.suppliers = [];
          this.pagination.total = 0;
        }
      } catch (error) {
        console.error('获取供货商列表失败:', error);
        ElMessage.error('获取供货商列表失败: ' + (error.message || '网络错误'));
      } finally {
        this.loading = false;
      }
    },
    async deleteSupplier(row) {
      ElMessageBox.confirm(`确定要删除供货商 "${row.name}" 吗？`, '确认删除', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async () => {
        try {
          await supplierAPI.deleteSupplier(row.id);
          ElMessage.success('供货商删除成功');
          this.fetchSuppliers();
        } catch (error) {
          console.error('删除供货商失败:', error);
          ElMessage.error('删除供货商失败: ' + (error.message || '网络错误'));
        }
      }).catch(() => {});
    },
    async submitSupplierForm() {
      this.$refs.supplierForm.validate(async (valid) => {
        if (valid) {
          try {
            // 标准化数据格式
            const normalizedData = this.normalizeFormData(this.supplierForm);
            
            if (this.supplierForm.id) {
              // 更新供货商
              await supplierAPI.updateSupplier(this.supplierForm.id, normalizedData);
              ElMessage.success('供货商更新成功');
            } else {
              // 创建供货商
              await supplierAPI.createSupplier(normalizedData);
              ElMessage.success('供货商创建成功');
            }
            this.dialogVisible = false;
            this.fetchSuppliers();
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
    getStatusType(status) {
      switch (status) {
        case 'active': return 'success';
        case 'inactive': return 'info';
        case 'suspended': return 'danger';
        default: return 'info';
      }
    },
    // 验证账户字段 - 只在接口类型为zjmf或v10时必填
    validateAccount(rule, value, callback) {
      if (this.supplierForm.interface_type === 'zjmf' || this.supplierForm.interface_type === 'v10') {
        if (!value) {
          callback(new Error('供应商账户不能为空'));
        } else {
          callback();
        }
      } else {
        callback(); // 如果不是zjmf或v10类型，则不需要验证账户字段
      }
    },
    
    getStatusText(status) {
      switch (status) {
        case 'active': return '活跃';
        case 'inactive': return '非活跃';
        case 'suspended': return '暂停';
        default: return status;
      }
    },
    handleSizeChange(val) {
      this.pagination.limit = val;
      this.fetchSuppliers();
    },
    handleCurrentChange(val) {
      this.pagination.page = val;
      this.fetchSuppliers();
    },
    
    // 获取接口类型标签样式
    getInterfaceTypeTagType(type) {
      switch (type) {
        case 'manual': return 'info';
        case 'zjmf': return 'success';
        case 'v10': return 'warning';
        default: return 'info';
      }
    },
    
    // 处理接口类型变化
    handleInterfaceTypeChange(value) {
      // 清空相关字段当切换接口类型时
      if (value === 'manual') {
        this.autoConfigForm.username = '';
        this.autoConfigForm.api_key = '';
      }
    },
    
    // 测试连接
    async testConnection() {
      if (!this.autoConfigForm.interface_url || !this.autoConfigForm.username || !this.autoConfigForm.api_key) {
        ElMessage.warning('请填写完整的配置信息');
        return;
      }
      
      this.testLoading = true;
      try {
        // 调用后端API测试连接
        const testData = {
          interface_type: this.autoConfigForm.interface_type,
          interface_url: this.autoConfigForm.interface_url,
          username: this.autoConfigForm.username,
          api_key: this.autoConfigForm.api_key
        };
        
        const response = await zjmfAPI.getZJMFInfo(testData);
        
        // 根据响应判断连接是否成功
        if (response && (response.code === 200 || response.success)) {
          this.testResult = {
            icon: 'success',
            title: '连接成功',
            subtitle: '已成功连接到智简魔方系统',
            success: true
          };
          this.autoConfigStep = 2;
        } else {
          throw new Error(response.message || '连接失败');
        }
      } catch (error) {
        console.error('连接测试失败:', error);
        this.testResult = {
          icon: 'error',
          title: '连接失败',
          subtitle: error.message || '无法连接到指定的接口地址，请检查配置信息',
          success: false
        };
        this.autoConfigStep = 2;
      } finally {
        this.testLoading = false;
      }
    },
    
    // 确认自动配置
    async confirmAutoConfig() {
      try {
        // 准备供应商数据
        const supplierData = {
          name: this.autoConfigForm.name || `智简魔方供应商_${Date.now()}`,
          interface_type: this.autoConfigForm.interface_type,
          interface_url: this.autoConfigForm.interface_url,
          username: this.autoConfigForm.username,
          api_key: this.autoConfigForm.api_key,
          status: 'active'
        };
        
        // 创建供应商
        await supplierAPI.createSupplier(supplierData);
        ElMessage.success('供应商创建成功');
        
        // 关闭对话框并刷新列表
        this.autoConfigDialogVisible = false;
        this.fetchSuppliers();
      } catch (error) {
        console.error('创建供应商失败:', error);
        ElMessage.error('创建供应商失败: ' + (error.message || '网络错误'));
      }
    },
    
    // 格式化日期
    formatDate(dateString) {
      if (!dateString) return '-';
      const date = new Date(dateString);
      return date.toLocaleString('zh-CN');
    }
  }
};
</script>

<style scoped>
.search-form {
  margin-bottom: 20px;
}
</style>