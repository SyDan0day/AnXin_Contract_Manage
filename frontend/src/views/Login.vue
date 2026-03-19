<template>
  <div class="login-container">
    <div class="login-left">
      <div class="brand-section">
        <div class="brand-logo">
          <svg viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
            <rect width="48" height="48" rx="12" fill="url(#logoGrad)"/>
            <path d="M14 16h20v4H14zM14 24h16v4H14zM14 32h20v4H14z" fill="white" opacity="0.9"/>
            <circle cx="36" cy="36" r="8" fill="white" opacity="0.2"/>
            <defs>
              <linearGradient id="logoGrad" x1="0" y1="0" x2="48" y2="48">
                <stop stop-color="#6366F1"/>
                <stop offset="1" stop-color="#8B5CF6"/>
              </linearGradient>
            </defs>
          </svg>
        </div>
        <h1 class="brand-title">安信合同</h1>
        <p class="brand-subtitle">智能合同管理解决方案</p>
      </div>
      
      <div class="features-list">
        <div class="feature-item">
          <div class="feature-icon">
            <el-icon><Document /></el-icon>
          </div>
          <div class="feature-text">
            <h3>合同全生命周期管理</h3>
            <p>从签订到执行，全程数字化跟踪</p>
          </div>
        </div>
        <div class="feature-item">
          <div class="feature-icon">
            <el-icon><Clock /></el-icon>
          </div>
          <div class="feature-text">
            <h3>智能到期提醒</h3>
            <p>提前预警，避免合同逾期风险</p>
          </div>
        </div>
        <div class="feature-item">
          <div class="feature-icon">
            <el-icon><DataLine /></el-icon>
          </div>
          <div class="feature-text">
            <h3>数据可视化分析</h3>
            <p>多维度统计，洞察业务趋势</p>
          </div>
        </div>
      </div>
    </div>
    
    <div class="login-right">
      <div class="login-card">
        <div class="login-header">
          <h2>欢迎回来</h2>
          <p>请登录您的账号继续</p>
        </div>
        
        <el-form ref="loginFormRef" :model="loginForm" :rules="loginRules" size="large">
          <el-form-item prop="username">
            <el-input 
              v-model="loginForm.username" 
              placeholder="请输入用户名"
              :prefix-icon="User"
              clearable
            />
          </el-form-item>
          <el-form-item prop="password">
            <el-input
              v-model="loginForm.password"
              type="password"
              placeholder="请输入密码"
              :prefix-icon="Lock"
              show-password
              clearable
              @keyup.enter="handleLogin"
            />
          </el-form-item>
          
          <div class="login-options">
            <el-checkbox v-model="rememberMe">记住我</el-checkbox>
            <a href="#" class="forgot-link">忘记密码？</a>
          </div>
          
          <el-form-item>
            <el-button 
              type="primary" 
              :loading="loading" 
              class="login-btn"
              @click="handleLogin"
            >
              <el-icon v-if="!loading"><Right /></el-icon>
              {{ loading ? '登录中...' : '登 录' }}
            </el-button>
          </el-form-item>
        </el-form>
        
        <div class="register-prompt">
          <span>还没有账号？</span>
          <router-link to="/register">立即注册</router-link>
        </div>
      </div>
      
      <div class="login-footer">
        <p>© 2024 安信合同管理系统 · 保留所有权利</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Document, Clock, DataLine, User, Lock, Right } from '@element-plus/icons-vue'
import { login } from '@/api/auth'
import { useUserStore } from '@/store/user'

const router = useRouter()
const userStore = useUserStore()
const loginFormRef = ref(null)
const loading = ref(false)
const rememberMe = ref(false)

const loginForm = reactive({
  username: '',
  password: ''
})

const loginRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少6位', trigger: 'blur' }
  ]
}

const handleLogin = async () => {
  await loginFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const res = await login(loginForm)
        userStore.setToken(res.access_token)
        userStore.setUserInfo(res.user_info || { username: loginForm.username, role: 'user' })
        ElMessage.success({ message: '欢迎回来！', duration: 2000 })
        router.push('/')
      } catch (error) {
        console.error('登录失败:', error)
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped>
.login-container {
  display: flex;
  min-height: 100vh;
  background: linear-gradient(135deg, #F8FAFC 0%, #E2E8F0 100%);
}

.login-left {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 60px;
  background: linear-gradient(135deg, #6366F1 0%, #8B5CF6 100%);
  position: relative;
  overflow: hidden;
}

.login-left::before {
  content: '';
  position: absolute;
  top: -50%;
  right: -50%;
  width: 100%;
  height: 100%;
  background: radial-gradient(circle, rgba(255,255,255,0.1) 0%, transparent 70%);
}

.login-left::after {
  content: '';
  position: absolute;
  bottom: -30%;
  left: -30%;
  width: 80%;
  height: 80%;
  background: radial-gradient(circle, rgba(255,255,255,0.08) 0%, transparent 70%);
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

.features-list {
  position: relative;
  z-index: 1;
}

.feature-item {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  margin-bottom: 32px;
}

.feature-icon {
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

.feature-text h3 {
  font-size: 16px;
  font-weight: 600;
  color: white;
  margin: 0 0 4px;
}

.feature-text p {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.6);
  margin: 0;
}

.login-right {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding: 40px;
  position: relative;
}

.login-card {
  width: 100%;
  max-width: 420px;
  background: white;
  border-radius: 24px;
  padding: 48px;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.08);
}

.login-header {
  text-align: center;
  margin-bottom: 40px;
}

.login-header h2 {
  font-size: 28px;
  font-weight: 600;
  color: #1E293B;
  margin: 0 0 8px;
}

.login-header p {
  font-size: 14px;
  color: #64748B;
  margin: 0;
}

.login-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.forgot-link {
  font-size: 14px;
  color: #6366F1;
  text-decoration: none;
  transition: color 0.2s;
}

.forgot-link:hover {
  color: #4F46E5;
}

.login-btn {
  width: 100%;
  height: 48px;
  font-size: 16px;
  font-weight: 500;
  border-radius: 12px;
  background: linear-gradient(135deg, #6366F1 0%, #8B5CF6 100%);
  border: none;
  transition: all 0.3s ease;
}

.login-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(99, 102, 241, 0.35);
}

.register-prompt {
  text-align: center;
  margin-top: 24px;
  padding-top: 24px;
  border-top: 1px solid #E2E8F0;
  color: #64748B;
  font-size: 14px;
}

.register-prompt a {
  color: #6366F1;
  text-decoration: none;
  font-weight: 500;
  margin-left: 4px;
}

.register-prompt a:hover {
  color: #4F46E5;
}

.login-footer {
  margin-top: 40px;
}

.login-footer p {
  font-size: 12px;
  color: #94A3B8;
  margin: 0;
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
  margin-bottom: 20px;
}

:deep(.el-checkbox__label) {
  color: #64748B;
  font-size: 14px;
}

@media (max-width: 1024px) {
  .login-left {
    display: none;
  }
  
  .login-right {
    padding: 24px;
  }
  
  .login-card {
    padding: 32px 24px;
  }
}
</style>
