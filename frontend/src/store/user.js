/**
 * 用户状态管理模块 (Pinia Store)
 * 
 * 使用 Pinia 管理全局用户状态
 * 状态持久化到 localStorage，支持页面刷新后保持登录状态
 * 
 * 存储的数据：
 * - token: JWT 认证令牌
 * - userInfo: 用户基本信息（用户名、角色等）
 */

import { defineStore } from 'pinia'
import { ref } from 'vue'

/**
 * 定义用户 Store
 * 使用 Composition API 风格
 */
export const useUserStore = defineStore('user', () => {
  // ==================== 状态定义 ====================
  
  /**
   * JWT 认证令牌
   * 初始值从 localStorage 读取，页面刷新时保持登录状态
   */
  const token = ref(localStorage.getItem('token') || '')
  
  /**
   * 用户信息对象
   * 包含用户名、邮箱、角色等
   * 初始值从 localStorage 读取并解析为对象
   */
  const userInfo = ref(JSON.parse(localStorage.getItem('userInfo') || 'null'))
  
  // ==================== 方法定义 ====================
  
  /**
   * 设置认证令牌
   * 同时更新内存状态和 localStorage 持久化
   * 
   * @param {string} newToken - 新的 JWT 令牌
   */
  const setToken = (newToken) => {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }
  
  /**
   * 设置用户信息
   * 同时更新内存状态和 localStorage 持久化
   * 
   * @param {Object} info - 用户信息对象
   */
  const setUserInfo = (info) => {
    userInfo.value = info
    localStorage.setItem('userInfo', JSON.stringify(info))
  }
  
  /**
   * 退出登录
   * 清除 token、userInfo 以及 localStorage 中的数据
   */
  const logout = () => {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('userInfo')
  }
  
  // ==================== 导出 ====================
  
  return {
    // 状态
    token,
    userInfo,
    
    // 方法
    setToken,
    setUserInfo,
    logout
  }
})
