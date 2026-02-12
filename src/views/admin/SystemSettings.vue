<template>
  <div class="system-settings">
    <el-card>
      <template #header>
        <div class="clearfix">
          <span>系统设置</span>
        </div>
      </template>
      
      <el-tabs v-model="activeTab" type="card">
        <el-tab-pane label="基础设置" name="basic">
          <el-form :model="settingsForm" :rules="settingsRules" ref="settingsForm" label-width="150px">
            <el-form-item label="网站名称" prop="site_name">
              <el-input v-model="settingsForm.site_name" placeholder="网站名称"></el-input>
            </el-form-item>
            
            <el-form-item label="网站描述" prop="site_description">
              <el-input v-model="settingsForm.site_description" type="textarea" placeholder="网站描述"></el-input>
            </el-form-item>
            
            <el-form-item label="联系邮箱" prop="contact_email">
              <el-input v-model="settingsForm.contact_email" placeholder="联系邮箱"></el-input>
            </el-form-item>
            
            <el-form-item label="支持邮箱" prop="support_email">
              <el-input v-model="settingsForm.support_email" placeholder="支持邮箱"></el-input>
            </el-form-item>
            
            <el-form-item label="默认时区" prop="default_timezone">
              <el-select v-model="settingsForm.default_timezone" placeholder="选择时区">
                <el-option label="亚洲/上海" value="Asia/Shanghai"></el-option>
                <el-option label="亚洲/东京" value="Asia/Tokyo"></el-option>
                <el-option label="美国/纽约" value="America/New_York"></el-option>
                <el-option label="欧洲/伦敦" value="Europe/London"></el-option>
              </el-select>
            </el-form-item>
          </el-form>
        </el-tab-pane>
        
        <el-tab-pane label="系统设置" name="system">
          <el-form :model="settingsForm" :rules="settingsRules" ref="settingsForm" label-width="150px">
            <el-form-item label="维护模式" prop="maintenance_mode">
              <el-switch
                v-model="settingsForm.maintenance_mode"
                active-value="true"
                inactive-value="false"
                active-text="开启"
                inactive-text="关闭">
              </el-switch>
              <div class="setting-description">开启后普通用户将无法访问网站</div>
            </el-form-item>
            
            <el-form-item label="注册开关" prop="registration_enabled">
              <el-switch
                v-model="settingsForm.registration_enabled"
                active-value="true"
                inactive-value="false"
                active-text="开启"
                inactive-text="关闭">
              </el-switch>
              <div class="setting-description">控制是否允许新用户注册</div>
            </el-form-item>
            
            <el-form-item label="邮件通知" prop="email_notifications">
              <el-switch
                v-model="settingsForm.email_notifications"
                active-value="true"
                inactive-value="false"
                active-text="开启"
                inactive-text="关闭">
              </el-switch>
              <div class="setting-description">开启系统邮件通知功能</div>
            </el-form-item>
            
            <el-form-item label="最大上传大小(MB)" prop="max_upload_size">
              <el-input-number v-model.number="settingsForm.max_upload_size" :min="1" :max="1024" :step="1"></el-input-number>
              <div class="setting-description">限制文件上传的最大大小（MB）</div>
            </el-form-item>
            
            <el-form-item label="会话超时(分钟)" prop="session_timeout">
              <el-input-number v-model.number="settingsForm.session_timeout" :min="1" :max="1440" :step="1"></el-input-number>
              <div class="setting-description">用户会话超时时间（分钟）</div>
            </el-form-item>
          </el-form>
        </el-tab-pane>
        
        <el-tab-pane label="站务设置" name="affairs">
          <el-form :model="settingsForm" :rules="settingsRules" ref="settingsForm" label-width="150px">
            <el-alert title="站务设置功能将在后续版本中实现" type="info" :closable="false" show-icon>
              <p>在此处将可以配置：</p>
              <ul style="margin-top: 10px;">
                <li>公告管理</li>
                <li>站内信设置</li>
                <li>积分规则</li>
                <li>用户等级体系</li>
                <li>内容审核策略</li>
              </ul>
            </el-alert>
          </el-form>
        </el-tab-pane>
      </el-tabs>
      
      <div style="margin-top: 20px; text-align: right;">
        <el-button type="primary" @click="saveSettings">保存设置</el-button>
        <el-button @click="resetSettings">重置</el-button>
      </div>
    </el-card>
  </div>
</template>

<script>
import { systemAPI } from '@/services/api.js';

export default {
  name: 'SystemSettings',
  data() {
    return {
      activeTab: 'basic',
      settingsForm: {
        site_name: '',
        site_description: '',
        contact_email: '',
        support_email: '',
        default_timezone: 'Asia/Shanghai',
        maintenance_mode: 'false',
        registration_enabled: 'true',
        email_notifications: 'true',
        max_upload_size: 10,
        session_timeout: 120
      },
      settingsRules: {
        site_name: [
          { required: true, message: '请输入网站名称', trigger: 'blur' }
        ],
        contact_email: [
          { required: true, message: '请输入联系邮箱', trigger: 'blur' },
          { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
        ],
        support_email: [
          { required: true, message: '请输入支持邮箱', trigger: 'blur' },
          { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
        ],
        max_upload_size: [
          { required: true, message: '请输入最大上传大小', trigger: 'blur' },
          { type: 'number', min: 1, max: 1024, message: '上传大小必须在1-1024MB之间', trigger: 'blur' }
        ],
        session_timeout: [
          { required: true, message: '请输入会话超时时间', trigger: 'blur' },
          { type: 'number', min: 1, max: 1440, message: '会话超时时间必须在1-1440分钟之间', trigger: 'blur' }
        ]
      }
    };
  },
  async mounted() {
    await this.loadSettings();
  },
  methods: {
    async loadSettings() {
      try {
        const response = await systemAPI.getSystemSettings();
        if (response.code === 200 && response.data) {
          // 如果API返回的数据格式不同，调整字段映射
          this.settingsForm = {
            site_name: response.data.site_name || response.data.siteName || '智简魔方业务管理系统',
            site_description: response.data.site_description || response.data.siteDescription || '专业的业务管理平台',
            contact_email: response.data.contact_email || response.data.contactEmail || 'admin@example.com',
            support_email: response.data.support_email || response.data.supportEmail || 'support@example.com',
            default_timezone: response.data.default_timezone || response.data.defaultTimezone || 'Asia/Shanghai',
            maintenance_mode: response.data.maintenance_mode || response.data.maintenanceMode || 'false',
            registration_enabled: response.data.registration_enabled || response.data.registrationEnabled || 'true',
            email_notifications: response.data.email_notifications || response.data.emailNotifications || 'true',
            max_upload_size: response.data.max_upload_size || response.data.maxUploadSize || 10,
            session_timeout: response.data.session_timeout || response.data.sessionTimeout || 120
          };
        } else {
          // 如果获取设置失败，使用默认值
          this.resetToDefaults();
          this.$message.info('未能获取系统设置，使用默认值');
        }
      } catch (error) {
        console.error('获取系统设置失败:', error);
        this.resetToDefaults();
        this.$message.error('获取系统设置失败');
      }
    },
    resetToDefaults() {
      this.settingsForm = {
        site_name: '智简魔方业务管理系统',
        site_description: '专业的业务管理平台',
        contact_email: 'admin@example.com',
        support_email: 'support@example.com',
        default_timezone: 'Asia/Shanghai',
        maintenance_mode: 'false',
        registration_enabled: 'true',
        email_notifications: 'true',
        max_upload_size: 10,
        session_timeout: 120
      };
    },
    async saveSettings() {
      this.$refs.settingsForm.validate(async (valid) => {
        if (valid) {
          try {
            const response = await systemAPI.updateSystemSettings(this.settingsForm);
            if (response.code === 200) {
              this.$message.success('系统设置保存成功');
            } else {
              this.$message.error(response.message || '保存设置失败');
            }
          } catch (error) {
            console.error('保存系统设置失败:', error);
            this.$message.error('保存系统设置失败');
          }
        }
      });
    },
    resetSettings() {
      this.$confirm('确定要重置所有设置为默认值吗？此操作不可撤销。', '确认重置', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        this.resetToDefaults();
        this.$message.info('已重置为默认设置');
      }).catch(() => {
        // 用户取消重置
      });
    }
  }
};
</script>

<style scoped>
.setting-description {
  font-size: 12px;
  color: #909399;
  margin-top: 5px;
}
</style>