# 安信合同管理系统 - 快速上手指南

## 目录

- [项目简介](#项目简介)
- [技术栈](#技术栈)
- [项目结构](#项目结构)
- [环境要求](#环境要求)
- [快速启动](#快速启动)
- [开发指南](#开发指南)
- [API 接口说明](#api-接口说明)
- [核心功能流程](#核心功能流程)
- [常见问题](#常见问题)

---

## 项目简介

安信合同管理系统是一个基于 **Go + Gin + MySQL** 后端和 **Vue 3 + Element Plus** 前端的合同全生命周期管理平台。

### 主要功能

- ✅ 用户权限管理（注册、登录、角色管理）
- ✅ 客户/供应商管理
- ✅ 合同全生命周期管理（创建→审批→执行→归档）
- ✅ 审批流程管理（多级审批）
- ✅ 状态变更审批（归档、终止、执行中、待付款）
- ✅ 文档管理（文件上传）
- ✅ 统计报表与数据可视化
- ✅ 合同到期提醒
- ✅ 审计日志

---

## 技术栈

### 前端

| 技术 | 版本 | 说明 |
|------|------|------|
| Vue | 3.3+ | 渐进式前端框架 |
| Vite | 8.0+ | 下一代前端构建工具 |
| Element Plus | 2.4+ | Vue 3 UI 组件库 |
| Pinia | 2.1+ | 状态管理 |
| Vue Router | 4.2+ | 路由管理 |
| Axios | 1.6+ | HTTP 客户端 |
| ECharts | 5.4+ | 数据可视化 |

### 后端

| 技术 | 版本 | 说明 |
|------|------|------|
| Go | 1.21+ | 服务器语言 |
| Gin | - | HTTP Web 框架 |
| GORM | - | ORM 库 |
| MySQL | 8.0+ | 数据库 |
| JWT | - | 身份认证 |

---

## 项目结构

```
AnXin_Contract_Manage/
├── frontend/                      # 前端项目
│   ├── src/
│   │   ├── api/                  # API 接口封装
│   │   │   ├── auth.js           # 认证相关 API
│   │   │   ├── contract.js        # 合同管理 API
│   │   │   ├── customer.js        # 客户管理 API
│   │   │   ├── approval.js        # 审批管理 API
│   │   │   └── audit.js           # 审计日志 API
│   │   ├── views/                # 页面组件
│   │   │   ├── Login.vue          # 登录页
│   │   │   ├── Register.vue       # 注册页
│   │   │   ├── Layout.vue         # 主布局
│   │   │   ├── Dashboard.vue      # 仪表盘
│   │   │   ├── Contract.vue       # 合同列表
│   │   │   ├── ContractDetail.vue # 合同详情
│   │   │   ├── Customer.vue       # 客户管理
│   │   │   ├── User.vue           # 用户管理
│   │   │   ├── Approval.vue        # 审批管理
│   │   │   ├── Reminder.vue       # 到期提醒
│   │   │   └── Audit.vue          # 审计日志
│   │   ├── router/               # 路由配置
│   │   │   └── index.js
│   │   ├── store/                # 状态管理
│   │   │   └── user.js
│   │   ├── utils/                # 工具函数
│   │   │   └── request.js         # Axios 封装
│   │   ├── App.vue               # 根组件
│   │   └── main.js               # 入口文件
│   ├── public/                   # 静态资源
│   ├── vite.config.js            # Vite 配置
│   └── package.json
└── (后端文件...)
```

---

## 环境要求

### 前端环境

- **Node.js**: 16+ (推荐 18+)
- **npm**: 8+ 或 **pnpm**: 8+

### 后端环境

- **Go**: 1.21+
- **MySQL**: 8.0+

---

## 快速启动

### 方式一：使用 Docker（推荐）

```bash
# 克隆项目
git clone <repository-url>
cd AnXin_Contract_Manage

# 复制环境变量文件
cp .env.example .env

# 启动所有服务（后端 + 前端 + MySQL）
docker-compose up -d

# 访问地址
# 前端: http://localhost
# 后端 API: http://localhost:8000
```

### 方式二：手动部署

#### 1. 启动后端

```bash
# 进入项目根目录
cd AnXin_Contract_Manage

# 配置数据库连接（编辑 .env 文件）
vim .env

# 安装 Go 依赖
go mod download

# 启动后端服务
go run main.go
# 服务地址: http://localhost:8000
```

#### 2. 启动前端

```bash
# 进入前端目录
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
# 服务地址: http://localhost:3000
```

---

## 开发指南

### 新增页面

1. **创建页面组件**
   ```vue
   <!-- src/views/MyPage.vue -->
   <template>
     <div>我的页面</div>
   </template>
   
   <script setup>
   // 页面逻辑
   </script>
   ```

2. **配置路由**
   ```javascript
   // src/router/index.js
   {
     path: '/my-page',
     name: 'MyPage',
     component: () => import('@/views/MyPage.vue'),
     meta: { title: '我的页面' }
   }
   ```

### 新增 API

1. **创建 API 文件**
   ```javascript
   // src/api/myModule.js
   import request from '@/utils/request'
   
   export const getMyData = (params) => {
     return request({
       url: '/my-endpoint',
       method: 'get',
       params
     })
   }
   ```

2. **在组件中使用**
   ```vue
   <script setup>
   import { getMyData } from '@/api/myModule'
   
   const data = await getMyData({ id: 1 })
   </script>
   ```

### 新增组件

1. **创建组件文件**
   ```vue
   <!-- src/components/MyComponent.vue -->
   <template>
     <div class="my-component">组件内容</div>
   </template>
   
   <script setup>
   // 组件逻辑
   </script>
   
   <style scoped>
   .my-component { color: red; }
   </style>
   ```

2. **在页面中使用**
   ```vue
   <template>
     <MyComponent />
   </template>
   
   <script setup>
   import MyComponent from '@/components/MyComponent.vue'
   </script>
   ```

---

## API 接口说明

### 认证接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | /api/auth/register | 用户注册 | 否 |
| POST | /api/auth/login | 用户登录 | 否 |
| GET | /api/auth/users | 获取用户列表 | 是(管理员) |
| PUT | /api/auth/users/:id | 更新用户 | 是(管理员) |
| DELETE | /api/auth/users/:id | 删除用户 | 是(管理员) |

### 合同接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | /api/contracts | 获取合同列表 | 是 |
| GET | /api/contracts/:id | 获取合同详情 | 是 |
| POST | /api/contracts | 创建合同 | 是 |
| PUT | /api/contracts/:id | 更新合同 | 是 |
| DELETE | /api/contracts/:id | 删除合同 | 是 |
| PUT | /api/contracts/:id/status | 更新合同状态 | 是 |
| POST | /api/contracts/:id/status-change | 申请状态变更 | 是 |

### 审批接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | /api/pending-approvals | 待审批列表 | 是 |
| PUT | /api/approvals/:id | 更新审批状态 | 是 |
| GET | /api/statistics | 获取统计数据 | 是 |

### 合同状态说明

| 状态值 | 说明 | 是否需要审批 |
|--------|------|-------------|
| draft | 草稿 | 否 |
| pending | 待审批 | 否 |
| active | 已生效 | 是 |
| completed | 已完成 | 否 |
| terminated | 已终止 | 是 |
| archived | 已归档 | 是 |

---

## 核心功能流程

### 1. 用户登录流程

```
用户输入用户名密码
    ↓
提交登录请求 (/api/auth/login)
    ↓
后端验证并返回 JWT Token
    ↓
前端保存 Token 到 localStorage
    ↓
跳转到首页 /dashboard
```

### 2. 合同创建流程

```
点击"新建合同"按钮
    ↓
填写合同表单（标题、客户、金额、日期等）
    ↓
点击"确定"提交
    ↓
调用 POST /api/contracts
    ↓
返回成功后刷新列表
    ↓
显示成功提示
```

### 3. 合同审批流程

```
创建合同（状态：pending）
    ↓
管理员查看待审批列表
    ↓
点击"审批"按钮
    ↓
选择通过/拒绝
    ↓
调用 PUT /api/approvals/:id
    ↓
合同状态变更为 approved/pending
```

### 4. 状态变更审批流程

```
某些状态变更需要审批（archived、terminated）
    ↓
用户申请状态变更
    ↓
调用 POST /api/contracts/:id/status-change
    ↓
状态变为 pending_status_change
    ↓
管理员审批
    ↓
调用 POST /api/status-change-requests/:id/approve 或 reject
    ↓
状态变更生效
```

---

## 常见问题

### Q1: 前端无法请求后端 API

**原因**: 跨域问题或代理配置错误

**解决**:
1. 检查 `vite.config.js` 中的代理配置
2. 确保后端已启动并运行在正确端口

### Q2: Token 过期怎么办

**现象**: 页面自动跳转到登录页

**原因**: JWT Token 默认 30 分钟过期

**解决**: 重新登录即可

### Q3: 如何添加新角色

1. 修改后端代码添加新角色
2. 在前端 `src/store/user.js` 中更新权限判断

### Q4: 文件上传失败

**原因**: 文件类型不在白名单中

**支持的文件类型**: PDF、Word、Excel、图片、文本文件

**不支持**: HTML、PHP、ASP 等脚本文件

### Q5: 如何查看审计日志

**路径**: 只有 admin 和 audit_admin 角色可以在侧边栏看到"审计日志"菜单

---

## 默认账号

| 用户名 | 密码 | 角色 | 说明 |
|--------|------|------|------|
| admin | Admin@12345 | admin | 超级管理员 |
| auditadmin | Audit@12345 | audit_admin | 审计管理员 |

> ⚠️ **安全提示**: 生产环境请务必修改默认密码！

---

## 联系方式

如有问题，请提交 Issue 或联系开发团队。

---

*文档版本: v1.0.0*
*最后更新: 2026-03-24*
