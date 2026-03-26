<template>
  <!-- 登录页面容器：左右分栏布局 -->
  <div class="login-container">
    <!-- 左侧品牌展示区域 -->
    <div class="login-left">
      <!-- 品牌信息区 -->
      <div class="brand-section">
        <div class="brand-logo">
          <!-- 项目 Logo -->
          <img src="/log.png" alt="Logo" style="width: 48px; height: 48px;" />
        </div>
        <h1 class="brand-title">安信合同</h1>
        <p class="brand-subtitle">智能合同管理解决方案</p>
      </div>
      
      <!-- 功能特性列表 -->
      <div class="features-list">
        <!-- 特性1：合同全生命周期管理 -->
        <div class="feature-item">
          <div class="feature-icon">
            <el-icon><Document /></el-icon>
          </div>
          <div class="feature-text">
            <h3>合同全生命周期管理</h3>
            <p>从签订到执行，全程数字化跟踪</p>
          </div>
        </div>
        <!-- 特性2：智能到期提醒 -->
        <div class="feature-item">
          <div class="feature-icon">
            <el-icon><Clock /></el-icon>
          </div>
          <div class="feature-text">
            <h3>智能到期提醒</h3>
            <p>提前预警，避免合同逾期风险</p>
          </div>
        </div>
        <!-- 特性3：数据可视化分析 -->
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
    
    <!-- 右侧登录表单区域 -->
    <div class="login-right">
      <!-- 登录卡片 -->
      <div class="login-card">
        <!-- 登录标题 -->
        <div class="login-header">
          <h2>欢迎回来</h2>
          <p>请登录您的账号继续</p>
        </div>
        
        <!-- 登录表单 -->
        <el-form 
          ref="loginFormRef" 
          :model="loginForm" 
          :rules="loginRules" 
          size="large"
        >
          <!-- 用户名输入框 -->
          <el-form-item prop="username">
            <el-input 
              v-model="loginForm.username" 
              placeholder="请输入用户名"
              :prefix-icon="User"
              clearable
            />
          </el-form-item>
          
          <!-- 密码输入框 -->
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
          
          <!-- 记住我和忘记密码 -->
          <div class="login-options">
            <el-checkbox v-model="rememberMe">记住我</el-checkbox>
            <a href="#" class="forgot-link">忘记密码？</a>
          </div>
          
          <!-- 登录按钮 -->
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
        
        <!-- 注册提示 -->
        <div class="register-prompt">
          <span>还没有账号？</span>
          <router-link to="/register">立即注册</router-link>
        </div>
      </div>
      
      <!-- 页脚 -->
      <div class="login-footer">
        <p>© 2024 安信合同管理系统 · 保留所有权利</p>
      </div>
    </div>
  </div>
</template>

<script setup>
/**
 * 登录页面逻辑
 * 
 * 功能：
 * 1. 用户名密码验证
 * 2. 调用登录 API
 * 3. 成功后保存 token 并跳转首页
 */

import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Document, Clock, DataLine, User, Lock, Right } from '@element-plus/icons-vue'

// API 方法
import { login } from '@/api/auth'

// Pinia 用户状态
import { useUserStore } from '@/store/user'

// ==================== 组件状态 ====================

// 路由实例，用于页面跳转
const router = useRouter()

// 用户状态 store
const userStore = useUserStore()

// 表单引用，用于表单验证
const loginFormRef = ref(null)

// 登录按钮加载状态
const loading = ref(false)

// 记住我复选框状态
const rememberMe = ref(false)

// ==================== 表单数据 ====================

/**
 * 登录表单数据
 * username: 用户名
 * password: 密码
 */
const loginForm = reactive({
  username: '',
  password: ''
})

/**
 * 表单验证规则
 */
const loginRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少6位', trigger: 'blur' }
  ]
}

// ==================== 事件处理 ====================

/**
 * 处理登录提交
 * 
 * 流程：
 * 1. 表单验证
 * 2. 调用登录 API
 * 3. 保存 token 和用户信息
 * 4. 跳转到首页
 */
const handleLogin = async () => {
  // 触发表单验证
  await loginFormRef.value.validate(async (valid) => {
    if (valid) {
      // 显示加载状态
      loading.value = true
      try {
        // 调用登录 API
        const res = await login(loginForm)
        
        // 保存 token 到 Pinia store（会自动持久化到 localStorage）
        userStore.setToken(res.access_token)
        
        // 保存用户信息
        userStore.setUserInfo(res.user_info || { 
          username: loginForm.username, 
          role: 'user' 
        })
        
        // 显示欢迎提示
        ElMessage.success({ message: '欢迎回来！', duration: 2000 })
        
        // 跳转到首页
        router.push('/')
      } catch (error) {
        // 登录失败，错误提示由 request.js 统一处理
        console.error('登录失败:', error)
      } finally {
        // 关闭加载状态
        loading.value = false
      }
    }
  })
}
</script>

<style scoped>
/* ==================== 布局样式 ==================== */

/* 全屏登录容器：左右分栏 */
.login-container {
  display: flex;
  min-height: 100vh;
  background: linear-gradient(135deg, #F8FAFC 0%, #E2E8F0 100%);
}

/* 左侧品牌区域：占一半宽度，紫色渐变背景 */
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

/* 装饰性圆形背景 */
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

/* 品牌信息区 */
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

/* 功能特性列表 */
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

/* 右侧登录表单区域 */
.login-right {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding: 40px;
  position: relative;
}

/* 登录卡片 */
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

/* 登录选项 */
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

/* 登录按钮 */
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

/* 注册提示 */
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

/* 页脚 */
.login-footer {
  margin-top: 40px;
}

.login-footer p {
  font-size: 12px;
  color: #94A3B8;
  margin: 0;
}

/* 表单样式覆盖 */
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

/* 响应式：小屏幕隐藏左侧品牌区 */
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
