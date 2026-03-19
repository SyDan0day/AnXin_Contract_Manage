<template>
  <div class="register-container">
    <div class="register-left">
      <div class="brand-section">
        <div class="brand-logo">
          <svg viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
            <rect width="48" height="48" rx="12" fill="url(#regLogoGrad)"/>
            <path d="M14 16h20v4H14zM14 24h16v4H14zM14 32h20v4H14z" fill="white" opacity="0.9"/>
            <circle cx="36" cy="36" r="8" fill="white" opacity="0.2"/>
            <defs>
              <linearGradient id="regLogoGrad" x1="0" y1="0" x2="48" y2="48">
                <stop stop-color="#6366F1"/>
                <stop offset="1" stop-color="#8B5CF6"/>
              </linearGradient>
            </defs>
          </svg>
        </div>
        <h1 class="brand-title">安信合同</h1>
        <p class="brand-subtitle">智能合同管理解决方案</p>
      </div>
      
      <div class="benefits-list">
        <div class="benefit-item">
          <div class="benefit-icon">
            <el-icon><CircleCheck /></el-icon>
          </div>
          <div class="benefit-text">
            <h3>安全可靠</h3>
            <p>企业级数据加密保护</p>
          </div>
        </div>
        <div class="benefit-item">
          <div class="benefit-icon">
            <el-icon><DataLine /></el-icon>
          </div>
          <div class="benefit-text">
            <h3>高效管理</h3>
            <p>全流程数字化跟踪</p>
          </div>
        </div>
        <div class="benefit-item">
          <div class="benefit-icon">
            <el-icon><Connection /></el-icon>
          </div>
          <div class="benefit-text">
            <h3>随时随地</h3>
            <p>云端同步随时访问</p>
          </div>
        </div>
      </div>
    </div>
    
    <div class="register-right">
      <div class="register-card">
        <div class="register-header">
          <h2>创建账号</h2>
          <p>开始使用安信合同管理系统</p>
        </div>
        
        <el-form ref="registerFormRef" :model="registerForm" :rules="registerRules" size="large">
          <div class="form-row">
            <el-form-item prop="username">
              <el-input 
                v-model="registerForm.username" 
                placeholder="用户名"
                :prefix-icon="User"
                clearable
              />
            </el-form-item>
            <el-form-item prop="full_name">
              <el-input 
                v-model="registerForm.full_name" 
                placeholder="姓名"
                :prefix-icon="UserFilled"
                clearable
              />
            </el-form-item>
          </div>
          
          <el-form-item prop="email">
            <el-input 
              v-model="registerForm.email" 
              placeholder="邮箱地址"
              :prefix-icon="Message"
              clearable
            />
          </el-form-item>
          
          <el-form-item prop="phone">
            <el-input 
              v-model="registerForm.phone" 
              placeholder="手机号（可选）"
              :prefix-icon="Phone"
              clearable
            />
          </el-form-item>
          
          <el-form-item prop="password">
            <el-input
              v-model="registerForm.password"
              type="password"
              placeholder="设置密码（至少6位）"
              :prefix-icon="Lock"
              show-password
              clearable
            />
          </el-form-item>
          
          <el-form-item prop="confirmPassword">
            <el-input
              v-model="registerForm.confirmPassword"
              type="password"
              placeholder="确认密码"
              :prefix-icon="Lock"
              show-password
              clearable
            />
          </el-form-item>
          
          <el-form-item>
            <el-button 
              type="primary" 
              :loading="loading" 
              class="register-btn"
              @click="handleRegister"
            >
              <el-icon v-if="!loading"><Plus /></el-icon>
              {{ loading ? '注册中...' : '创 建 账 号' }}
            </el-button>
          </el-form-item>
        </el-form>
        
        <div class="login-prompt">
          <span>已有账号？</span>
          <router-link to="/login">立即登录</router-link>
        </div>
      </div>
      
      <div class="register-footer">
        <p>注册即表示同意 <a href="#">服务条款</a> 和 <a href="#">隐私政策</a></p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { 
  User, UserFilled, Message, Phone, Lock, Plus,
  CircleCheck, DataLine, Connection 
} from '@element-plus/icons-vue'
import { register } from '@/api/auth'

const router = useRouter()
const registerFormRef = ref(null)
const loading = ref(false)

const registerForm = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
  full_name: '',
  phone: '',
  department: ''
})

const validatePass = (rule, value, callback) => {
  if (value === '') {
    callback(new Error('请输入密码'))
  } else {
    if (registerForm.confirmPassword !== '') {
      registerFormRef.value.validateField('confirmPassword')
    }
    callback()
  }
}

const validatePass2 = (rule, value, callback) => {
  if (value === '') {
    callback(new Error('请再次输入密码'))
  } else if (value !== registerForm.password) {
    callback(new Error('两次输入密码不一致'))
  } else {
    callback()
  }
}

const registerRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  password: [
    { required: true, validator: validatePass, trigger: 'blur' },
    { min: 6, message: '密码长度至少6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, validator: validatePass2, trigger: 'blur' }
  ],
  full_name: [
    { required: true, message: '请输入姓名', trigger: 'blur' }
  ],
  phone: [
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ]
}

const handleRegister = async () => {
  await registerFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const { confirmPassword, ...data } = registerForm
        await register(data)
        ElMessage.success({ message: '注册成功，请登录', duration: 2000 })
        router.push('/login')
      } catch (error) {
        console.error('注册失败:', error)
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped>
.register-container {
  display: flex;
  min-height: 100vh;
  background: linear-gradient(135deg, #F8FAFC 0%, #E2E8F0 100%);
}

.register-left {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 60px;
  background: linear-gradient(135deg, #6366F1 0%, #8B5CF6 100%);
  position: relative;
  overflow: hidden;
}

.register-left::before {
  content: '';
  position: absolute;
  top: -50%;
  right: -50%;
  width: 100%;
  height: 100%;
  background: radial-gradient(circle, rgba(255,255,255,0.1) 0%, transparent 70%);
}

.brand-section {
  position: relative;
  z-index: 1;
  margin-bottom: 60px;
}

.brand-logo {
  width: 64px;
  height: 64px;
  margin-bottom: 24px;
}

.brand-logo svg {
  width: 100%;
  height: 100%;
}

.brand-title {
  font-size: 36px;
  font-weight: 700;
  color: white;
  margin: 0 0 8px;
  letter-spacing: 2px;
}

.brand-subtitle {
  font-size: 16px;
  color: rgba(255, 255, 255, 0.7);
  margin: 0;
}

.benefits-list {
  position: relative;
  z-index: 1;
}

.benefit-item {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  margin-bottom: 32px;
}

.benefit-icon {
  width: 44px;
  height: 44px;
  background: rgba(255, 255, 255, 0.15);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 20px;
  flex-shrink: 0;
}

.benefit-text h3 {
  font-size: 16px;
  font-weight: 600;
  color: white;
  margin: 0 0 4px;
}

.benefit-text p {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.6);
  margin: 0;
}

.register-right {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding: 40px;
  position: relative;
}

.register-card {
  width: 100%;
  max-width: 440px;
  background: white;
  border-radius: 24px;
  padding: 40px;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.08);
}

.register-header {
  text-align: center;
  margin-bottom: 32px;
}

.register-header h2 {
  font-size: 26px;
  font-weight: 600;
  color: #1E293B;
  margin: 0 0 8px;
}

.register-header p {
  font-size: 14px;
  color: #64748B;
  margin: 0;
}

.form-row {
  display: flex;
  gap: 16px;
}

.form-row .el-form-item {
  flex: 1;
}

.register-btn {
  width: 100%;
  height: 48px;
  font-size: 16px;
  font-weight: 500;
  border-radius: 12px;
  background: linear-gradient(135deg, #6366F1 0%, #8B5CF6 100%);
  border: none;
  transition: all 0.3s ease;
}

.register-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(99, 102, 241, 0.35);
}

.login-prompt {
  text-align: center;
  margin-top: 24px;
  padding-top: 24px;
  border-top: 1px solid #E2E8F0;
  color: #64748B;
  font-size: 14px;
}

.login-prompt a {
  color: #6366F1;
  text-decoration: none;
  font-weight: 500;
  margin-left: 4px;
}

.login-prompt a:hover {
  color: #4F46E5;
}

.register-footer {
  margin-top: 24px;
}

.register-footer p {
  font-size: 12px;
  color: #94A3B8;
  margin: 0;
  text-align: center;
}

.register-footer a {
  color: #6366F1;
  text-decoration: none;
}

:deep(.el-input__wrapper) {
  border-radius: 10px;
  padding: 8px 16px;
  box-shadow: 0 0 0 1px #E2E8F0 inset;
  transition: all 0.2s;
}

:deep(.el-input__wrapper:hover) {
  box-shadow: 0 0 0 1px #CBD5E1 inset;
}

:deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 2px rgba(99, 102, 241, 0.2), 0 0 0 1px #6366F1 inset;
}

:deep(.el-form-item) {
  margin-bottom: 16px;
}

@media (max-width: 1024px) {
  .register-left {
    display: none;
  }
  
  .register-right {
    padding: 24px;
  }
  
  .register-card {
    padding: 32px 24px;
  }
  
  .form-row {
    flex-direction: column;
    gap: 0;
  }
}
</style>
