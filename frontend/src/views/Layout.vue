<template>
  <el-container class="layout-container">
    <el-aside width="260px">
      <div class="sidebar">
        <div class="logo-area">
          <div class="logo-icon">
            <svg viewBox="0 0 32 32" fill="none">
              <rect width="32" height="32" rx="8" fill="url(#sidebarLogo)"/>
              <path d="M8 10h16v3H8zM8 15h12v3H8zM8 20h16v3H8z" fill="white" opacity="0.9"/>
              <defs>
                <linearGradient id="sidebarLogo" x1="0" y1="0" x2="32" y2="32">
                  <stop stop-color="#6366F1"/>
                  <stop offset="1" stop-color="#8B5CF6"/>
                </linearGradient>
              </defs>
            </svg>
          </div>
          <span class="logo-text">安心合同</span>
        </div>
        
        <el-menu
          :default-active="activeMenu"
          router
          class="sidebar-menu"
        >
          <el-menu-item index="/dashboard">
            <el-icon><Odometer /></el-icon>
            <span>仪表盘</span>
          </el-menu-item>
          <el-menu-item index="/contracts">
            <el-icon><Document /></el-icon>
            <span>合同管理</span>
            <el-badge v-if="notificationCounts.expiringContracts > 0" :value="notificationCounts.expiringContracts" :max="99" class="menu-badge" />
          </el-menu-item>
          <el-menu-item index="/customers">
            <el-icon><OfficeBuilding /></el-icon>
            <span>客户管理</span>
          </el-menu-item>
          <el-menu-item index="/approvals">
            <el-icon><Checked /></el-icon>
            <span>审批管理</span>
            <el-badge v-if="notificationCounts.pendingApprovals + notificationCounts.pendingStatusChanges > 0" :value="notificationCounts.pendingApprovals + notificationCounts.pendingStatusChanges" :max="99" class="menu-badge" />
          </el-menu-item>
          <el-menu-item index="/reminders">
            <el-icon><Bell /></el-icon>
            <span>到期提醒</span>
            <el-badge v-if="notificationCounts.expiringContracts > 0" :value="notificationCounts.expiringContracts" :max="99" class="menu-badge" />
          </el-menu-item>
          <el-menu-item index="/users">
            <el-icon><UserFilled /></el-icon>
            <span>用户管理</span>
          </el-menu-item>
        </el-menu>
        
        <div class="sidebar-footer">
          <div class="user-card">
            <el-avatar :size="36" class="user-avatar">
              {{ userStore.userInfo?.username?.charAt(0)?.toUpperCase() }}
            </el-avatar>
            <div class="user-info">
              <div class="user-name">{{ userStore.userInfo?.username }}</div>
              <div class="user-role">管理员</div>
            </div>
          </div>
        </div>
      </div>
    </el-aside>
    
    <el-container>
      <el-header>
        <div class="header-left">
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item v-if="currentRoute">{{ currentRoute }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        
        <div class="header-right">
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
                <el-dropdown-item command="profile">
                  <el-icon><User /></el-icon>
                  个人设置
                </el-dropdown-item>
                <el-dropdown-item divided command="logout">
                  <el-icon><SwitchButton /></el-icon>
                  退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      
      <el-main>
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { computed, ref, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/store/user'
import { ElMessageBox } from 'element-plus'
import { getNotificationCounts } from '@/api/approval'
import { 
  Odometer, Document, OfficeBuilding, Checked, Bell, 
  UserFilled, User, ArrowDown, SwitchButton 
} from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const notificationCounts = ref({
  pendingApprovals: 0,
  pendingStatusChanges: 0,
  expiringContracts: 0,
  total: 0
})

let notificationTimer = null

const loadNotifications = async () => {
  try {
    const counts = await getNotificationCounts()
    notificationCounts.value = counts
  } catch (error) {
    console.error('Failed to load notifications:', error)
  }
}

onMounted(() => {
  loadNotifications()
  notificationTimer = setInterval(loadNotifications, 30000)
})

onUnmounted(() => {
  if (notificationTimer) {
    clearInterval(notificationTimer)
  }
})

const activeMenu = computed(() => route.path)

const routeNames = {
  '/dashboard': '仪表盘',
  '/contracts': '合同管理',
  '/customers': '客户管理',
  '/approvals': '审批管理',
  '/reminders': '到期提醒',
  '/users': '用户管理'
}

const currentRoute = computed(() => routeNames[route.path])

const handleCommand = (command) => {
  if (command === 'profile') {
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
    ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(() => {
      userStore.logout()
      router.push('/login')
    })
  }
}
</script>

<style scoped>
.layout-container {
  height: 100vh;
}

.el-aside {
  background: white;
  box-shadow: 2px 0 12px rgba(0, 0, 0, 0.04);
}

.sidebar {
  height: 100%;
  display: flex;
  flex-direction: column;
}

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

.sidebar-menu {
  flex: 1;
  border-right: none;
  padding: 12px 0;
}

:deep(.el-menu-item) {
  height: 48px;
  margin: 4px 12px;
  border-radius: 12px;
  color: #64748B;
  font-weight: 500;
  transition: all 0.2s ease;
  position: relative;
}

.menu-badge {
  position: absolute;
  right: 60px;
  top: 8px;
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

:deep(.el-menu-item .el-icon) {
  margin-right: 12px;
  font-size: 18px;
}

.sidebar-footer {
  padding: 16px;
  border-top: 1px solid #F1F5F9;
}

.user-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: #F8FAFC;
  border-radius: 12px;
}

.user-avatar {
  background: linear-gradient(135deg, #6366F1 0%, #8B5CF6 100%);
  color: white;
  font-weight: 600;
  font-size: 14px;
}

.user-info {
  flex: 1;
  min-width: 0;
}

.user-name {
  font-size: 14px;
  font-weight: 600;
  color: #1E293B;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.user-role {
  font-size: 12px;
  color: #94A3B8;
}

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

:deep(.el-breadcrumb__item) {
  font-size: 14px;
}

:deep(.el-breadcrumb__inner) {
  color: #94A3B8;
}

:deep(.el-breadcrumb__item:last-child .el-breadcrumb__inner) {
  color: #1E293B;
  font-weight: 500;
}

.header-right {
  display: flex;
  align-items: center;
}

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

.el-main {
  background: #F8FAFC;
  padding: 24px;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

:deep(.el-dropdown-menu__item) {
  padding: 10px 20px;
  font-size: 14px;
}

:deep(.el-dropdown-menu__item .el-icon) {
  margin-right: 8px;
}
</style>
