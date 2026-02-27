# 合同管理系统

基于 Python FastAPI + MySQL（后端）和 Vue3 + Element Plus（前端）的合同管理系统。

## 功能模块

- 用户权限管理：用户注册、登录、角色管理
- 客户/供应商管理：客户信息增删改查
- 合同管理：合同信息管理、分类管理
- 合同执行跟踪：进度跟踪、付款记录
- 审批流程：合同审批、多级审批
- 到期提醒：合同到期提醒、续期管理
- 统计报表：数据统计分析
- 文档管理：合同文件上传、版本管理

## 技术栈

### 后端
- Python 3.8+
- FastAPI
- SQLAlchemy
- MySQL
- JWT 认证

### 前端
- Vue 3
- Vite
- Element Plus
- Pinia
- Vue Router
- Axios
- ECharts

## 安装

### 后端安装

1. 进入后端目录
```bash
cd AnXin_Contract_Manage
```

2. 安装依赖
```bash
pip install -r requirements.txt
```

3. 配置数据库

编辑 `app/core/config.py` 中的数据库配置：
```python
MYSQL_HOST = "localhost"
MYSQL_PORT = 3306
MYSQL_USER = "root"
MYSQL_PASSWORD = "your_password"
MYSQL_DATABASE = "contract_manage"
```

4. 创建数据库
```bash
mysql -u root -p
CREATE DATABASE contract_manage CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

5. 运行后端服务
```bash
uvicorn main:app --reload
```

后端 API 文档地址: http://localhost:8000/docs

### 前端安装

1. 进入前端目录
```bash
cd frontend
```

2. 安装依赖
```bash
npm install
```

3. 运行前端服务
```bash
npm run dev
```

前端访问地址: http://localhost:3000

4. 构建生产版本
```bash
npm run build
```

## API 端点

### 认证
- POST /api/auth/register - 用户注册
- POST /api/auth/login - 用户登录
- GET /api/auth/users - 获取用户列表
- GET /api/auth/users/{user_id} - 获取用户详情
- PUT /api/auth/users/{user_id} - 更新用户
- DELETE /api/auth/users/{user_id} - 删除用户

### 客户管理
- GET /api/customers - 获取客户列表
- GET /api/customers/{customer_id} - 获取客户详情
- POST /api/customers - 创建客户
- PUT /api/customers/{customer_id} - 更新客户
- DELETE /api/customers/{customer_id} - 删除客户
- GET /api/contract-types - 获取合同类型列表
- POST /api/contract-types - 创建合同类型

### 合同管理
- GET /api/contracts - 获取合同列表
- GET /api/contracts/{contract_id} - 获取合同详情
- POST /api/contracts - 创建合同
- PUT /api/contracts/{contract_id} - 更新合同
- DELETE /api/contracts/{contract_id} - 删除合同
- GET /api/contracts/{contract_id}/executions - 获取合同执行记录
- POST /api/contracts/{contract_id}/executions - 创建执行记录
- GET /api/contracts/{contract_id}/documents - 获取合同文档
- POST /api/contracts/{contract_id}/documents - 上传文档
- DELETE /api/documents/{document_id} - 删除文档

### 审批与提醒
- GET /api/contracts/{contract_id}/approvals - 获取审批记录
- POST /api/contracts/{contract_id}/approvals - 创建审批记录
- PUT /api/approvals/{approval_id} - 更新审批状态
- GET /api/contracts/{contract_id}/reminders - 获取提醒列表
- POST /api/contracts/{contract_id}/reminders - 创建提醒
- POST /api/reminders/{reminder_id}/send - 发送提醒
- GET /api/expiring-contracts - 获取即将到期的合同
- GET /api/statistics - 获取统计数据

## 项目结构

```
AnXin_Contract_Manage/
├── app/                      # 后端目录
│   ├── api/                  # API 路由
│   │   ├── auth.py           # 认证相关
│   │   ├── customer.py       # 客户管理
│   │   ├── contract.py       # 合同管理
│   │   └── approval.py       # 审批与提醒
│   ├── core/                 # 核心配置
│   │   ├── config.py         # 配置文件
│   │   ├── database.py       # 数据库连接
│   │   └── security.py       # 安全相关
│   ├── models/               # 数据库模型
│   │   └── models.py
│   ├── schemas/              # Pydantic 模型
│   │   ├── user.py
│   │   ├── customer.py
│   │   ├── contract.py
│   │   └── approval.py
│   ├── services/             # 业务逻辑
│   │   ├── user_service.py
│   │   ├── customer_service.py
│   │   ├── contract_service.py
│   │   └── approval_service.py
│   └── utils/                # 工具函数
├── frontend/                 # 前端目录
│   ├── src/
│   │   ├── api/              # API 接口
│   │   │   ├── auth.js
│   │   │   ├── customer.js
│   │   │   ├── contract.js
│   │   │   └── approval.js
│   │   ├── components/       # 公共组件
│   │   ├── router/           # 路由配置
│   │   │   └── index.js
│   │   ├── store/            # 状态管理
│   │   │   └── user.js
│   │   ├── utils/            # 工具函数
│   │   │   └── request.js
│   │   ├── views/            # 页面组件
│   │   │   ├── Login.vue
│   │   │   ├── Layout.vue
│   │   │   ├── Dashboard.vue
│   │   │   ├── Contract.vue
│   │   │   ├── Customer.vue
│   │   │   ├── User.vue
│   │   │   ├── Approval.vue
│   │   │   └── Reminder.vue
│   │   ├── App.vue
│   │   └── main.js
│   ├── index.html
│   ├── package.json
│   └── vite.config.js
├── main.py                   # 后端入口
└── requirements.txt          # 后端依赖列表
```

## 前端页面说明

- **登录页面** (Login.vue) - 用户登录
- **布局页面** (Layout.vue) - 主框架布局，包含侧边栏导航
- **仪表盘** (Dashboard.vue) - 数据统计、图表展示、即将到期合同
- **合同管理** (Contract.vue) - 合同的增删改查、状态管理
- **客户管理** (Customer.vue) - 客户/供应商信息管理
- **用户管理** (User.vue) - 用户信息管理、角色分配
- **审批管理** (Approval.vue) - 合同审批流程
- **到期提醒** (Reminder.vue) - 合同到期提醒管理

## 使用说明

1. 启动后端服务：`uvicorn main:app --reload`
2. 启动前端服务：`cd frontend && npm run dev`
3. 访问前端：http://localhost:3000
4. 使用注册功能创建管理员账户，或直接登录测试

## 开发说明

- 后端使用 FastAPI 自动生成 API 文档，访问 http://localhost:8000/docs 查看
- 前端使用 Vue 3 Composition API 开发
- 使用 Pinia 进行状态管理
- 使用 Element Plus 作为 UI 组件库
- 使用 ECharts 进行数据可视化# AnXin_Contract_Manage
