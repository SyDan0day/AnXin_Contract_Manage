/**
 * Vite 配置文件
 * 
 * Vite 是新一代前端构建工具，比 webpack 更快更简单
 * 
 * 配置说明：
 * 1. plugins: 使用 Vue 插件支持 .vue 文件
 * 2. resolve.alias: 设置路径别名 (@ 代表 src 目录)
 * 3. server: 开发服务器配置
 */

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

export default defineConfig({
  // ==================== 插件配置 ====================
  
  /**
   * Vue 插件
   * 让 Vite 能够处理 .vue 单文件组件
   */
  plugins: [vue()],
  
  // ==================== 路径解析配置 ====================
  
  resolve: {
    alias: {
      /**
       * 路径别名配置
       * '@' 指向 src 目录
       * 
       * 使用示例：
       * - import X from '@/components/X.vue'
       * - import Y from '@/api/y.js'
       * 
       * 等价于：
       * - import X from 'src/components/X.vue'
       * - import Y from 'src/api/y.js'
       */
      '@': resolve(__dirname, 'src')
    }
  },
  
  // ==================== 开发服务器配置 ====================
  
  server: {
    /**
     * 开发服务器端口
     * 访问地址: http://localhost:3000
     */
    port: 3000,
    
    /**
     * 代理配置
     * 将 /api 请求代理到后端服务器
     */
    proxy: {
      /**
       * /api 前缀的请求转发到后端
       * 
       * 示例：
       * 前端请求: GET /api/contracts
       * 实际请求: GET http://localhost:8000/api/contracts
       * 
       * 配置项说明：
       * - target: 后端服务器地址
       * - changeOrigin: 是否修改请求头的 Origin 字段
       */
      '/api': {
        target: 'http://localhost:8000',
        changeOrigin: true
      }
    }
  }
})
