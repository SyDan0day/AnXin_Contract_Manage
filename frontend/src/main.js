/**
 * 应用入口文件
 * 
 * 功能：
 * 1. 创建 Vue 3 应用实例
 * 2. 注册插件：
 *    - Pinia (状态管理)
 *    - Vue Router (路由管理)
 *    - Element Plus (UI 组件库)
 * 3. 注册 Element Plus 图标组件
 * 4. 挂载应用到 #app 元素
 */

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'

// 引入 Element Plus 中文语言包
import zhCn from 'element-plus/dist/locale/zh-cn.mjs'

// 引入 Element Plus 所有图标组件
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

// 引入根组件
import App from './App.vue'

// 引入路由配置
import router from './router'

// ==================== 创建应用实例 ====================

// 创建 Vue 应用
const app = createApp(App)

// 创建 Pinia 实例（状态管理）
const pinia = createPinia()

// ==================== 注册全局组件 ====================

// 注册所有 Element Plus 图标为全局组件
// 使用方式: <el-icon><图标名称 /></el-icon>
// 例如: <el-icon><Document /></el-icon>
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

// ==================== 注册插件 ====================

// 使用 Pinia 状态管理
app.use(pinia)

// 使用 Vue Router 路由管理
app.use(router)

// 使用 Element Plus UI 组件库
// locale: zhCn 设置为中文界面
app.use(ElementPlus, { locale: zhCn })

// ==================== 挂载应用 ====================

// 将应用挂载到 index.html 中的 #app 元素
app.mount('#app')
