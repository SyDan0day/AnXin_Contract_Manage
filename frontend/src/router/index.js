/**
 * 路由配置模块
 * 定义应用的路由规则和导航守卫
 * 
 * 路由结构：
 * - /login, /register: 无需认证的公共路由
 * - /: 主布局路由，包含子路由（需认证）
 */

import { createRouter, createWebHistory } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/store/user'

/**
 * 路由配置数组
 * 每个路由对象包含：
 * - path: URL路径
 * - name: 路由名称（用于编程式导航）
 * - component: 路由组件
 * - meta: 元数据（标题、权限等）
 */
const routes = [
  // ==================== 公共路由 ====================
  
  /**
   * 登录页面
   * 无需认证即可访问
   */
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { title: '登录' }
  },
  
  /**
   * 注册页面
   * 无需认证即可访问
   */
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/Register.vue'),
    meta: { title: '注册' }
  },
  
  // ==================== 受保护的布局路由 ====================
  
  /**
   * 主布局路由
   * 所有需要认证的页面都在此布局下
   * redirect: 默认重定向到仪表盘
   * requiresAuth: true 表示需要登录才能访问
   */
  {
    path: '/',
    component: () => import('@/views/Layout.vue'),
    redirect: '/dashboard',
    meta: { requiresAuth: true },
    
    // 子路由（嵌套在 Layout 中）
    children: [
      /**
       * 仪表盘/首页
       * 显示系统概览、统计图表、即将到期合同等
       */
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard.vue'),
        meta: { title: '仪表盘', icon: 'DataAnalysis' }
      },
      
      /**
       * 合同列表页面
       * 支持合同的增删改查、状态管理
       */
      {
        path: 'contracts',
        name: 'Contracts',
        component: () => import('@/views/Contract.vue'),
        meta: { title: '合同管理', icon: 'Document' }
      },
      
      /**
       * 合同详情页面
       * hidden: true 表示不在侧边栏菜单显示
       * 显示合同详细信息、执行跟踪、文档、生命周期等
       */
      {
        path: 'contracts/:id',
        name: 'ContractDetail',
        component: () => import('@/views/ContractDetail.vue'),
        meta: { title: '合同详情', hidden: true }
      },
      
      /**
       * 客户管理页面
       * 管理客户和供应商信息
       */
      {
        path: 'customers',
        name: 'Customers',
        component: () => import('@/views/Customer.vue'),
        meta: { title: '客户管理', icon: 'User' }
      },
      
      /**
       * 用户管理页面
       * 仅管理员可访问
       * roles: ['admin'] 表示需要 admin 角色
       */
      {
        path: 'users',
        name: 'Users',
        component: () => import('@/views/User.vue'),
        meta: { title: '用户管理', icon: 'UserFilled', roles: ['admin'] }
      },
      
      /**
       * 审批管理页面
       * 处理合同审批和状态变更审批
       * 审批角色可见：销售负责人、技术负责人、财务负责人、管理员
       */
      {
        path: 'approvals',
        name: 'Approvals',
        component: () => import('@/views/Approval.vue'),
        meta: { 
          title: '审批管理', 
          icon: 'Check',
          roles: ['admin', 'sales_manager', 'tech_leader', 'finance_leader']
        }
      },
      
      /**
       * 到期提醒页面
       * 显示即将到期的合同列表
       * 所有角色可见（销售只能看到自己的）
       */
      {
        path: 'reminders',
        name: 'Reminders',
        component: () => import('@/views/Reminder.vue'),
        meta: { 
          title: '到期提醒', 
          icon: 'Bell',
          roles: ['admin', 'contract_manager', 'audit_admin', 'sales_manager', 'tech_leader', 'finance_leader', 'user']
        }
      },
      
      /**
       * 审计日志页面
       * 仅管理员和审计管理员可访问
       */
      {
        path: 'audit',
        name: 'Audit',
        component: () => import('@/views/Audit.vue'),
        meta: { title: '审计日志', icon: 'Document', roles: ['admin', 'audit_admin'] }
      }
    ]
  }
]

/**
 * 创建路由实例
 * createWebHistory(): 使用 HTML5 History API，适合 SPA 应用
 */
const router = createRouter({
  history: createWebHistory(),
  routes
})

/**
 * 检查用户是否有权限访问该路由
 * @param {Array} roles - 路由允许的角色列表
 * @param {string} userRole - 用户角色
 * @returns {boolean} 是否有权限
 */
function hasPermission(roles, userRole) {
  // 没有设置角色要求，默认放行
  if (!roles || roles.length === 0) {
    return true
  }
  
  // 用户未登录，无权限
  if (!userRole) {
    return false
  }
  
  // 检查用户角色是否在允许列表中
  return roles.includes(userRole)
}

/**
 * 全局前置守卫
 * 在路由跳转前进行权限验证
 * 
 * @param {Route} to - 目标路由对象
 * @param {Route} from - 当前路由对象
 * @param {Function} next - 确认导航的函数
 */
router.beforeEach((to, from, next) => {
  // 获取用户状态
  const userStore = useUserStore()
  const userRole = userStore.userInfo?.role
  
  // 情况1: 目标路由需要认证但用户未登录
  // 重定向到登录页
  if (to.meta.requiresAuth && !userStore.token) {
    next('/login')
  }
  // 情况2: 用户已登录但访问登录/注册页
  // 重定向到首页，防止重复登录
  else if ((to.path === '/login' || to.path === '/register') && userStore.token) {
    next('/')
  }
  // 情况3: 检查角色权限
  else if (to.meta.roles && !hasPermission(to.meta.roles, userRole)) {
    // 无权限，跳转到首页
    ElMessage.warning('您没有权限访问该页面')
    next('/')
  }
  // 情况4: 其他情况，正常放行
  else {
    next()
  }
})

// 导出路由实例，供 main.js 使用
export default router
