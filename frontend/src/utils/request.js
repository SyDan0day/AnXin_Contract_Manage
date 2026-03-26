/**
 * Axios HTTP 请求封装模块
 * 
 * 功能：
 * 1. 创建 axios 实例，配置基础 URL 和超时时间
 * 2. 请求拦截器：自动添加 JWT Token 到请求头
 * 3. 响应拦截器：统一处理错误响应（如 401 未授权、403 禁止访问等）
 * 
 * 所有 API 请求都通过这个封装模块发送
 */

import axios from 'axios'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/store/user'

/**
 * 创建 axios 实例
 * baseURL: '/api' - 所有请求会自动添加 /api 前缀
 * timeout: 10000 - 请求超时时间 10 秒
 */
const request = axios.create({
  baseURL: '/api',
  timeout: 10000
})

/**
 * 请求拦截器
 * 在请求发送之前执行
 * 
 * 功能：
 * - 从 Pinia store 获取当前用户的 token
 * - 将 token 添加到请求头的 Authorization 字段
 * - 格式: Bearer <token>
 */
request.interceptors.request.use(
  config => {
    // 获取用户 token
    const userStore = useUserStore()
    
    // 如果存在 token，添加到请求头
    if (userStore.token) {
      config.headers.Authorization = `Bearer ${userStore.token}`
    }
    
    return config
  },
  error => {
    // 请求错误处理
    return Promise.reject(error)
  }
)

/**
 * 响应拦截器
 * 在收到响应之后、交给业务代码之前执行
 * 
 * 功能：
 * - 成功响应：直接返回 response.data
 * - 错误响应：根据状态码显示提示信息
 *   - 401: Token 过期，清除登录状态并跳转登录页
 *   - 403: 没有权限
 *   - 404: 请求资源不存在
 *   - 500: 服务器错误
 *   - 其他: 显示错误信息
 */
request.interceptors.response.use(
  response => {
    // 成功响应，直接返回数据部分
    return response.data
  },
  error => {
    // 错误响应处理
    if (error.response) {
      // 服务器返回了错误状态码
      const { status, data } = error.response
      
      switch (status) {
        case 401:
          // Token 过期或无效
          // 清除本地登录状态
          const userStore = useUserStore()
          userStore.logout()
          // 显示过期提示
          ElMessage.error('登录已过期，请重新登录')
          // 跳转到登录页
          window.location.href = '/login'
          break
          
        case 403:
          // 没有权限访问
          ElMessage.error('没有权限访问')
          break
          
        case 404:
          // 请求的资源不存在
          // 如果是通知相关的接口，静默处理
          if (config.url && (config.url.includes('notification'))) {
            console.warn('Notification endpoint not available:', config.url)
          } else {
            ElMessage.error('请求的资源不存在')
          }
          break
          
        case 500:
          // 服务器内部错误
          ElMessage.error('服务器错误')
          break
          
        default:
          // 其他错误，显示后端返回的错误信息
          ElMessage.error(data.error || data.detail || '请求失败')
      }
    } else {
      // 网络错误（如断网、超时等）
      ElMessage.error('网络错误')
    }
    
    // 将错误抛出，方便调用者进行额外处理
    return Promise.reject(error)
  }
)

// 导出封装后的 axios 实例
export default request
