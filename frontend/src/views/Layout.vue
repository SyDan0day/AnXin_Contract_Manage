<template>
  <!-- 主布局容器：使用 Element Plus 的 Container 组件 -->
  <el-container class="layout-container">
    <!-- 左侧边栏：260px 宽度 -->
    <el-aside width="260px">
      <div class="sidebar">
        <!-- Logo 区域 -->
        <div class="logo-area">
          <div class="logo-icon">
            <img src="/log.png" alt="Logo" style="width: 32px; height: 32px;" />
          </div>
          <span class="logo-text">安信合同</span>
        </div>
        
        <!-- 侧边栏菜单 -->
        <el-menu
          :default-active="activeMenu"
          router
          class="sidebar-menu"
        >
          <!-- 仪表盘 -->
          <el-menu-item index="/dashboard">
            <div class="menu-item-content">
              <el-icon><Odometer /></el-icon>
              <span>仪表盘</span>
            </div>
          </el-menu-item>
          
          <!-- 合同管理（带待办数量提示） -->
          <el-menu-item index="/contracts">
            <div class="menu-item-content">
              <el-icon><Document /></el-icon>
              <span>合同管理</span>
              <!-- 红点提示：显示即将到期合同数量 -->
              <el-badge 
                v-if="notificationCounts.expiringContracts > 0" 
                :value="notificationCounts.expiringContracts" 
                :max="99" 
                class="menu-badge-icon"
                type="warning"
              />
            </div>
          </el-menu-item>
          
          <!-- 客户管理 -->
          <el-menu-item index="/customers">
            <div class="menu-item-content">
              <el-icon><OfficeBuilding /></el-icon>
              <span>客户管理</span>
            </div>
          </el-menu-item>
          
          <!-- 审批管理（仅审批角色可见：销售负责人、技术负责人、财务负责人、管理员） -->
          <el-menu-item v-if="canApproval" index="/approvals">
            <div class="menu-item-content">
              <el-icon><Checked /></el-icon>
              <span>审批管理</span>
              <!-- 红点提示：显示待审批总数 -->
              <el-badge 
                v-if="notificationCounts.pendingApprovals + notificationCounts.pendingStatusChanges > 0" 
                :value="notificationCounts.pendingApprovals + notificationCounts.pendingStatusChanges" 
                :max="99" 
                class="menu-badge-icon"
                type="danger"
              />
            </div>
          </el-menu-item>
          
          <!-- 到期提醒（合同负责人、审批角色、管理员可见） -->
          <el-menu-item v-if="canViewAllContracts" index="/reminders">
            <div class="menu-item-content">
              <el-icon><Bell /></el-icon>
              <span>到期提醒</span>
              <!-- 红点提示 -->
              <el-badge 
                v-if="notificationCounts.expiringContracts > 0" 
                :value="notificationCounts.expiringContracts" 
                :max="99" 
                class="menu-badge-icon"
                type="warning"
              />
            </div>
          </el-menu-item>
          
          <!-- 用户管理（仅管理员可见） -->
          <el-menu-item v-if="isAdmin" index="/users">
            <div class="menu-item-content">
              <el-icon><UserFilled /></el-icon>
              <span>用户管理</span>
            </div>
          </el-menu-item>
          
          <!-- 审计日志（仅管理员和审计管理员可见） -->
          <el-menu-item v-if="isAuditAdmin" index="/audit">
            <div class="menu-item-content">
              <el-icon><Document /></el-icon>
              <span>审计日志</span>
            </div>
          </el-menu-item>
        </el-menu>
        
        <!-- 侧边栏底部：用户信息卡片 -->
        <div class="sidebar-footer">
          <div class="user-card">
            <!-- 用户头像：显示用户名首字母 -->
            <el-avatar :size="36" class="user-avatar">
              {{ userStore.userInfo?.username?.charAt(0)?.toUpperCase() }}
            </el-avatar>
            <div class="user-info">
              <div class="user-name">{{ userStore.userInfo?.username }}</div>
              <div class="user-role">{{ getRoleText(userStore.userInfo?.role) }}</div>
            </div>
          </div>
        </div>
      </div>
    </el-aside>
    
    <!-- 主内容区域 -->
    <el-container>
      <!-- 顶部导航栏 -->
      <el-header>
        <!-- 面包屑导航 -->
        <div class="header-left">
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item v-if="currentRoute">{{ currentRoute }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        
        <!-- 右侧用户操作区 -->
        <div class="header-right">
          <!-- 通知图标 -->
          <div class="header-notifications" @click="handleNotificationClick">
            <el-badge 
              :value="totalNotifications" 
              :hidden="totalNotifications === 0"
              :max="99"
              type="danger"
            >
              <el-icon :size="20"><Bell /></el-icon>
            </el-badge>
          </div>
          
          <!-- 用户下拉菜单 -->
          <el-dropdown @command="handleCommand" trigger="click">
            <div class="user-dropdown">
              <el-avatar :size="32" class="header-avatar">
                {{ userStore.userInfo?.username?.charAt(0)?.toUpperCase() }}
              </el-avatar>
              <span class="username">{{ userStore.userInfo?.username }}</span>
              <el-icon><ArrowDown /></el-icon>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <!-- 个人设置 -->
                <el-dropdown-item command="profile">
                  <el-icon><User /></el-icon>
                  个人设置
                </el-dropdown-item>
                <!-- 退出登录 -->
                <el-dropdown-item divided command="logout">
                  <el-icon><SwitchButton /></el-icon>
                  退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <!-- 通知抽屉 -->
      <el-drawer v-model="notificationDrawerVisible" title="我的通知" direction="rtl" size="400px">
        <div class="notification-header">
          <el-button type="primary" link @click="handleMarkAllRead">全部标为已读</el-button>
        </div>
        <div class="notification-list" v-loading="loadingNotifications">
          <el-empty v-if="notificationList.length === 0" description="暂无新通知" />
          <div 
            v-for="notification in notificationList" 
            :key="notification.id" 
            class="notification-item"
            :class="{ unread: !notification.is_read }"
            @click="handleMarkAsRead(notification)"
          >
            <div class="notification-icon">
              <el-icon v-if="notification.type === 'rejected'" color="#F56C6C"><WarningFilled /></el-icon>
              <el-icon v-else-if="notification.type === 'approved'" color="#67C23A"><CircleCheckFilled /></el-icon>
              <el-icon v-else color="#909399"><Bell /></el-icon>
            </div>
            <div class="notification-content">
              <div class="notification-title">{{ notification.title }}</div>
              <div class="notification-text">{{ notification.content }}</div>
              <div class="notification-time">{{ formatTime(notification.created_at) }}</div>
            </div>
          </div>
        </div>
      </el-drawer>
      
      <!-- 主内容区域：渲染子路由 -->
      <el-main>
        <router-view v-slot="{ Component }">
          <!-- 页面切换动画 -->
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
/**
 * 主布局组件
 * 
 * 功能：
 * 1. 提供统一的页面布局框架
 * 2. 管理侧边栏导航菜单
 * 3. 显示用户信息和通知
 * 4. 处理退出登录
 */

import { computed, ref, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/store/user'
import { ElMessageBox } from 'element-plus'

// API: 获取通知数量
import { getNotificationCounts } from '@/api/approval'
import { getUnreadNotifications, getUnreadCount, markAsRead, markAllAsRead } from '@/api/notification'

// Element Plus 图标
import { 
  Odometer, Document, OfficeBuilding, Checked, Bell, 
  UserFilled, User, ArrowDown, SwitchButton,
  WarningFilled, CircleCheckFilled
} from '@element-plus/icons-vue'

// ==================== 组件状态 ====================

// 路由实例
const router = useRouter()
const route = useRoute()

// 用户状态 store
const userStore = useUserStore()

// ==================== 计算属性 ====================

/**
 * 判断是否为审计管理员
 * admin 和 audit_admin 角色可以查看审计日志
 */
const isAuditAdmin = computed(() => {
  const role = userStore.userInfo?.role
  return role === 'admin' || role === 'audit_admin'
})

/**
 * 判断是否为管理员
 */
const isAdmin = computed(() => {
  const role = userStore.userInfo?.role
  return role === 'admin'
})

/**
 * 判断是否为普通销售（只能看到自己创建的合同）
 */
const isSales = computed(() => {
  const role = userStore.userInfo?.role
  return role === 'user'
})

/**
 * 判断是否为审批角色（销售负责人、技术负责人、财务负责人、管理员）
 * 可以访问审批管理页面
 */
const canApproval = computed(() => {
  const role = userStore.userInfo?.role
  return ['admin', 'sales_manager', 'tech_leader', 'finance_leader'].includes(role)
})

/**
 * 判断是否可以查看所有合同（合同负责人、管理员、审计管理员、审批角色）
 * 可以访问到期提醒页面
 * 销售(user)也可以看到自己创建的即将到期合同
 */
const canViewAllContracts = computed(() => {
  const role = userStore.userInfo?.role
  // 所有角色都可以看到到期提醒（销售可以看到自己的）
  return ['admin', 'contract_manager', 'audit_admin', 'sales_manager', 'tech_leader', 'finance_leader', 'user'].includes(role)
})

/**
 * 获取角色中文名称
 */
const getRoleText = (role) => {
  const roleMap = {
    admin: '管理员',
    sales_manager: '销售负责人',
    tech_leader: '技术负责人',
    finance_leader: '财务负责人',
    contract_manager: '合同负责人',
    user: '销售',
    audit_admin: '审计管理员'
  }
  return roleMap[role] || role
}

const formatTime = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now - date
  const minutes = Math.floor(diff / 60000)
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(diff / 86400000)
  
  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  if (hours < 24) return `${hours}小时前`
  if (days < 7) return `${days}天前`
  
  return date.toLocaleDateString('zh-CN')
}

/**
 * 计算总通知数量
 * = 待审批 + 待处理状态变更 + 即将到期 + 用户通知
 */
const totalNotifications = computed(() => {
  const { pendingApprovals, pendingStatusChanges, expiringContracts, userNotifications } = notificationCounts.value
  return pendingApprovals + pendingStatusChanges + expiringContracts + userNotifications
})

// ==================== 通知数据 ====================

/**
 * 通知数量统计
 * 从后端 API 获取
 */
const notificationCounts = ref({
  pendingApprovals: 0,        // 待审批合同数量
  pendingStatusChanges: 0,  // 待审批状态变更数量
  expiringContracts: 0,      // 即将到期合同数量
  userNotifications: 0,      // 用户通知数量（合同被退回等）
  total: 0                   // 总计
})

/**
 * 定时器引用，用于清理
 */
let notificationTimer = null

/**
 * 加载通知数量
 * 页面加载时调用
 */
const loadNotificationCounts = async () => {
  try {
    const counts = await getNotificationCounts()
    notificationCounts.value = counts
  } catch (error) {
    // 静默处理通知加载失败，不影响用户体验
    console.warn('Failed to load notification counts:', error)
  }
  
  try {
    const userNotifRes = await getUnreadCount()
    notificationCounts.value.userNotifications = userNotifRes.count || 0
  } catch (error) {
    // 静默处理用户通知加载失败
    console.warn('Failed to load user notifications:', error)
  }
}

// ==================== 生命周期 ====================

// 组件挂载时启动定时刷新
onMounted(() => {
  loadNotificationCounts()
  // 每 5 秒刷新一次通知数量（实时检测审批状态）
  notificationTimer = setInterval(loadNotificationCounts, 5000)
})

// 组件卸载时清理定时器
onUnmounted(() => {
  if (notificationTimer) {
    clearInterval(notificationTimer)
  }
})

// ==================== 路由相关 ====================

/**
 * 当前激活的菜单项
 * 根据当前路由路径确定
 */
const activeMenu = computed(() => route.path)

/**
 * 路由名称映射（用于面包屑显示）
 */
const routeNames = {
  '/dashboard': '仪表盘',
  '/contracts': '合同管理',
  '/customers': '客户管理',
  '/approvals': '审批管理',
  '/reminders': '到期提醒',
  '/users': '用户管理'
}

/**
 * 当前路由的中文名称
 */
const currentRoute = computed(() => routeNames[route.path])

// ==================== 事件处理 ====================

const notificationDrawerVisible = ref(false)
const notificationList = ref([])
const loadingNotifications = ref(false)

const loadNotificationList = async () => {
  loadingNotifications.value = true
  try {
    const data = await getUnreadNotifications()
    notificationList.value = data
  } finally {
    loadingNotifications.value = false
  }
}

const handleNotificationClick = async () => {
  await loadNotificationList()
  notificationDrawerVisible.value = true
}

const handleMarkAsRead = async (notification) => {
  if (!notification.is_read) {
    await markAsRead(notification.id)
    notification.is_read = true
    await loadNotificationCounts()
  }
  if (notification.contract_id) {
    notificationDrawerVisible.value = false
    router.push(`/contracts/${notification.contract_id}`)
  }
}

const handleMarkAllRead = async () => {
  await markAllAsRead()
  await loadNotificationList()
  await loadNotificationCounts()
}

/**
 * 处理用户下拉菜单命令
 * @param {string} command - 命令标识
 */
const handleCommand = (command) => {
  if (command === 'profile') {
    // 显示个人设置弹窗
    ElMessageBox.alert(
      `<div style="padding: 10px;">
        <p><strong>用户名：</strong>${userStore.userInfo?.username || '-'}</p>
        <p><strong>邮箱：</strong>${userStore.userInfo?.email || '-'}</p>
        <p><strong>角色：</strong>${userStore.userInfo?.role === 'admin' ? '管理员' : '普通用户'}</p>
      </div>`,
      '个人设置', 
      {
        confirmButtonText: '确定',
        dangerouslyUseHTMLString: true,
      }
    )
  } else if (command === 'logout') {
    // 确认退出登录
    ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(() => {
      // 清除用户状态
      userStore.logout()
      // 跳转到登录页
      router.push('/login')
    })
  }
}
</script>

<style scoped>
/* ==================== 布局样式 ==================== */

/* 全屏容器 */
.layout-container {
  height: 100vh;
}

.el-aside {
  background: white;
  box-shadow: 2px 0 12px rgba(0, 0, 0, 0.04);
}

/* 侧边栏 */
.sidebar {
  height: 100%;
  display: flex;
  flex-direction: column;
}

/* Logo 区域 */
.logo-area {
  height: 64px;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 0 20px;
  border-bottom: 1px solid #F1F5F9;
}

.logo-icon {
  width: 32px;
  height: 32px;
}

.logo-icon svg {
  width: 100%;
  height: 100%;
}

.logo-text {
  font-size: 18px;
  font-weight: 600;
  color: #1E293B;
  letter-spacing: 1px;
}

/* 侧边栏菜单 */
.sidebar-menu {
  flex: 1;
  border-right: none;
  padding: 12px 0;
}

/* 菜单项内容 */
.menu-item-content {
  display: flex;
  align-items: center;
  position: relative;
  width: 100%;
}

.menu-item-content .el-icon {
  margin-right: 12px;
  font-size: 18px;
  position: relative;
}

/* 红点徽章位置 */
.menu-badge-icon {
  position: absolute;
  right: 10px;
  top: -6px;
  transform: translateY(-50%);
}

.menu-badge-icon :deep(.el-badge__content) {
  border: none;
  font-size: 11px;
  padding: 0 5px;
  height: 18px;
  line-height: 18px;
}

/* 菜单项样式 */
:deep(.el-menu-item) {
  height: 48px;
  margin: 4px 12px;
  border-radius: 12px;
  color: #64748B;
  font-weight: 500;
  transition: all 0.2s ease;
  padding: 0 20px !important;
}

:deep(.el-menu-item:hover) {
  background: #F8FAFC;
  color: #1E293B;
}

:deep(.el-menu-item.is-active) {
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.1) 0%, rgba(139, 92, 246, 0.1) 100%);
  color: #6366F1;
}

:deep(.el-menu-item.is-active .el-icon) {
  color: #6366F1;
}

/* 页面切换动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* 下拉菜单样式 */
:deep(.el-dropdown-menu__item) {
  padding: 10px 20px;
  font-size: 14px;
}

:deep(.el-dropdown-menu__item .el-icon) {
  margin-right: 8px;
}

/* 顶部导航栏 */
.el-header {
  background: white;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
}

.header-left {
  display: flex;
  align-items: center;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

/* 通知图标 */
.header-notifications {
  cursor: pointer;
  padding: 8px;
  border-radius: 8px;
  transition: background 0.2s;
}

.header-notifications:hover {
  background: #F8FAFC;
}

.header-notifications .el-icon {
  color: #64748B;
}

.header-notifications:hover .el-icon {
  color: #6366F1;
}

/* 用户下拉 */
.user-dropdown {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 12px;
  border-radius: 10px;
  cursor: pointer;
  transition: background 0.2s;
}

.user-dropdown:hover {
  background: #F8FAFC;
}

.header-avatar {
  background: linear-gradient(135deg, #6366F1 0%, #8B5CF6 100%);
  color: white;
  font-weight: 600;
  font-size: 12px;
}

.username {
  color: #1E293B;
  font-weight: 500;
  font-size: 14px;
}

.notification-header {
  display: flex;
  justify-content: flex-end;
  padding: 12px 16px;
  border-bottom: 1px solid #F1F5F9;
}

.notification-list {
  padding: 8px;
  max-height: calc(100vh - 200px);
  overflow-y: auto;
}

.notification-item {
  display: flex;
  gap: 12px;
  padding: 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.2s;
}

.notification-item:hover {
  background: #F8FAFC;
}

.notification-item.unread {
  background: #F0F9FF;
}

.notification-icon {
  flex-shrink: 0;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #F1F5F9;
  border-radius: 50%;
}

.notification-content {
  flex: 1;
  min-width: 0;
}

.notification-title {
  font-weight: 600;
  font-size: 14px;
  color: #1E293B;
  margin-bottom: 4px;
}

.notification-text {
  font-size: 13px;
  color: #64748B;
  margin-bottom: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.notification-time {
  font-size: 12px;
  color: #94A3B8;
}
</style>
