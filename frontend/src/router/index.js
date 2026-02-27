import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/store/user'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { title: '登录' }
  },
  {
    path: '/',
    component: () => import('@/views/Layout.vue'),
    redirect: '/dashboard',
    meta: { requiresAuth: true },
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard.vue'),
        meta: { title: '仪表盘', icon: 'DataAnalysis' }
      },
      {
        path: 'contracts',
        name: 'Contracts',
        component: () => import('@/views/Contract.vue'),
        meta: { title: '合同管理', icon: 'Document' }
      },
      {
        path: 'customers',
        name: 'Customers',
        component: () => import('@/views/Customer.vue'),
        meta: { title: '客户管理', icon: 'User' }
      },
      {
        path: 'users',
        name: 'Users',
        component: () => import('@/views/User.vue'),
        meta: { title: '用户管理', icon: 'UserFilled', roles: ['admin'] }
      },
      {
        path: 'approvals',
        name: 'Approvals',
        component: () => import('@/views/Approval.vue'),
        meta: { title: '审批管理', icon: 'Check' }
      },
      {
        path: 'reminders',
        name: 'Reminders',
        component: () => import('@/views/Reminder.vue'),
        meta: { title: '到期提醒', icon: 'Bell' }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  
  if (to.meta.requiresAuth && !userStore.token) {
    next('/login')
  } else if (to.path === '/login' && userStore.token) {
    next('/')
  } else {
    next()
  }
})

export default router