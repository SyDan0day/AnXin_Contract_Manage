# 合同管理系统

基于 Go + Gin + MySQL（后端）和 Vue3 + Element Plus（前端）的合同管理系统。

## 目录

- [功能模块](#功能模块)
- [技术栈](#技术栈)
- [快速开始](#快速开始)
  - [手动部署](#手动部署)
  - [Docker 部署](#docker-部署)
- [API 认证](#api-认证)
- [API 端点](#api-端点)
- [项目结构](#项目结构)
- [前端页面说明](#前端页面说明)
- [环境变量](#环境变量)
- [Docker 操作](#docker-操作)
- [数据备份与恢复](#数据备份与恢复)
- [安全更新 (最新版本)](#安全更新-最新版本)
- [安全说明](#安全说明)
- [常见问题](#常见问题)
- [变更日志](#变更日志)

## 功能模块

- **用户权限管理**：用户注册、登录、角色管理（管理员、经理、普通用户）
- **客户/供应商管理**：客户信息增删改查、客户分类、信用等级
- **合同管理**：合同信息管理、合同分类管理、合同状态跟踪
- **合同执行跟踪**：进度跟踪、付款记录，执行阶段管理
- **审批流程**：合同审批、多级审批、审批记录查询
- **状态变更审批**：关键状态变更（归档、终止、执行中、待付款）需管理员审批
- **合同生命周期**：完整的合同状态变更历史记录
- **合同归档**：已完成合同归档管理
- **到期提醒**：合同到期提醒、续期管理、提醒通知
- **统计报表**：数据统计分析、图表展示
- **文档管理**：合同文件上传、版本管理
- **合同类型管理**：合同类型分类管理
- **待办提示**：侧边栏菜单红点提示待办事项

## 技术栈

### 后端
- Go 1.21+
- Gin Web Framework（高性能 HTTP 框架）
- GORM（ORM 库）
- MySQL 8.0
- JWT（用户认证）
- bcrypt（密码加密）

### 前端
- Vue 3（渐进式前端框架）
- Vite（构建工具）
- Element Plus（UI 组件库）
- Pinia（状态管理）
- Vue Router（路由管理）
- Axios（HTTP 客户端）
- ECharts（数据可视化）

## 快速开始

### 手动部署

手动部署适合开发环境或不想使用 Docker 的用户。

#### 前置要求

- Go 1.21 或更高版本
- MySQL 8.0 或更高版本
- Node.js 16+ 和 npm
- Git

#### 后端部署步骤

1. **克隆项目**

```bash
git clone <repository-url>
cd AnXin_Contract_Manage
```

2. **安装 Go 依赖**

```bash
go mod download
```

3. **配置环境变量**

复制环境变量示例文件并修改配置：

```bash
cp .env.example .env
```

编辑 `.env` 文件，修改数据库连接信息：

```env
APP_NAME=合同管理系统
APP_VERSION=1.0.0

MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_USER=root
MYSQL_PASSWORD=your_password
MYSQL_DATABASE=contract_manage

SECRET_KEY=your-secret-key-change-in-production
JWT_ALGORITHM=HS256
ACCESS_TOKEN_EXPIRE_MINUTES=30

UPLOAD_DIR=uploads
```

4. **创建数据库**

```bash
mysql -u root -p -e "CREATE DATABASE contract_manage CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
```

如果使用 Docker 安装 MySQL：

```bash
docker run -d --name mysql \
  -e MYSQL_ROOT_PASSWORD=root123456 \
  -e MYSQL_DATABASE=contract_manage \
  -p 3306:3306 \
  mysql:8.0
```

5. **运行后端服务**

```bash
go run main.go
```

首次运行会自动创建数据库表。后端服务默认运行在 http://localhost:8000

#### 前端部署步骤

1. **进入前端目录**

```bash
cd frontend
```

2. **安装依赖**

```bash
npm install
```

3. **配置 API 地址（如需要）**

编辑 `vite.config.js` 配置后端 API 地址：

```javascript
export default defineConfig({
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8000',
        changeOrigin: true
      }
    }
  }
})
```

4. **运行开发服务器**

```bash
npm run dev
```

前端服务默认运行在 http://localhost:3000

5. **构建生产版本**

```bash
npm run build
```

构建完成后，静态文件会生成在 `dist` 目录，可以部署到任意 Web 服务器。

### Docker 部署

Docker 部署适合快速启动或生产环境。

#### 前置要求

- Docker 20.10+
- Docker Compose 2.0+

#### 部署步骤

1. **进入项目目录**

```bash
cd AnXin_Contract_Manage
```

2. **配置环境变量**

```bash
cp .env.example .env
```

3. **启动服务**

```bash
docker-compose up -d
```

首次启动会：
- 创建 MySQL 容器并初始化数据库
- 构建并启动后端容器
- 构建并启动前端容器

4. **访问系统**

- 前端：http://localhost
- 后端 API：http://localhost:8000
- MySQL：localhost:3306

#### 查看服务状态

```bash
docker-compose ps
```

#### 查看日志

```bash
# 查看所有服务日志
docker-compose logs -f

# 查看后端日志
docker-compose logs -f backend

# 查看前端日志
docker-compose logs -f frontend

# 查看数据库日志
docker-compose logs -f mysql
```

## API 认证

除 `/api/auth/register` 和 `/api/auth/login` 外，所有 API 需要 JWT 认证。

### 认证流程

1. **注册用户**（可选，已有用户可跳过）

```bash
curl -X POST http://localhost:8000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "email": "admin@example.com",
    "password": "password123",
    "full_name": "管理员"
  }'
```

2. **登录获取 Token**

```bash
curl -X POST http://localhost:8000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "password123"
  }'
```

响应：

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "bearer"
}
```

3. **在请求中携带 Token**

```bash
curl -X GET http://localhost:8000/api/auth/users \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Token 过期时间

默认 Token 有效期为 30 分钟，可在环境变量中修改 `ACCESS_TOKEN_EXPIRE_MINUTES`。

### 默认管理员账号

程序首次启动时会自动创建管理员账号，方便首次登录系统。

| 用户名 | 密码 | 角色 | 说明 |
|--------|------|------|------|
| admin | Admin@12345 | admin | 超级管理员，拥有所有权限 |
| auditadmin | Audit@12345 | audit_admin | 审计管理员，负责审计日志管理 |

**安全警告**：默认密码仅用于首次登录，生产环境必须修改！

如需修改管理员账号信息，可在 `.env` 中配置：

```env
# 超级管理员
ADMIN_USERNAME=admin
ADMIN_PASSWORD=StrongPassword123!
ADMIN_EMAIL=admin@example.com

# 审计管理员
AUDIT_ADMIN_USERNAME=auditadmin
AUDIT_ADMIN_PASSWORD=AuditPassword123!
AUDIT_ADMIN_EMAIL=audit@example.com
```

> 注意：如果管理员账号已存在，则不会重复创建。生产环境部署前请务必修改默认密码。

## API 端点

### 公共端点（无需认证）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/auth/register | 用户注册 |
| POST | /api/auth/login | 用户登录 |
| GET | / | 服务信息 |
| GET | /health | 健康检查 |

### 需要认证的端点

#### 用户管理（需要管理员权限）

| 方法 | 路径 | 说明 | 权限要求 |
|------|------|------|----------|
| GET | /api/auth/users | 获取用户列表 | 管理员 |
| GET | /api/auth/users/:user_id | 获取用户详情 | 管理员 |
| PUT | /api/auth/users/:user_id | 更新用户信息 | 管理员 |
| DELETE | /api/auth/users/:user_id | 删除用户 | 管理员 |

**注意**：用户管理 API 需要管理员权限（`admin` 角色），普通用户访问会返回 403 Forbidden 错误。

请求示例：

```bash
# 获取用户列表
curl -X GET http://localhost:8000/api/auth/users?skip=0&limit=100 \
  -H "Authorization: Bearer <token>"

# 更新用户
curl -X PUT http://localhost:8000/api/auth/users/1 \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "新名字",
    "email": "newemail@example.com"
  }'
```

#### 客户管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/customers | 获取客户列表 |
| GET | /api/customers/:customer_id | 获取客户详情 |
| POST | /api/customers | 创建客户 |
| PUT | /api/customers/:customer_id | 更新客户 |
| DELETE | /api/customers/:customer_id | 删除客户 |
| GET | /api/contract-types | 获取合同类型列表 |
| POST | /api/contract-types | 创建合同类型 |
| PUT | /api/contract-types/:type_id | 更新合同类型 |
| DELETE | /api/contract-types/:type_id | 删除合同类型 |

请求示例：

```bash
# 创建客户
curl -X POST http://localhost:8000/api/customers \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "客户名称",
    "code": "C001",
    "type": "customer",
    "contact_person": "张三",
    "contact_phone": "13800138000",
    "contact_email": "zhangsan@example.com"
  }'
```

#### 合同管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/contracts | 获取合同列表 |
| GET | /api/contracts/:contract_id | 获取合同详情 |
| POST | /api/contracts | 创建合同 |
| PUT | /api/contracts/:contract_id | 更新合同 |
| PUT | /api/contracts/:contract_id/status | 直接更新合同状态 |
| POST | /api/contracts/:contract_id/status-change | 申请状态变更（需要审批的状态） |
| GET | /api/contracts/:contract_id/status-change | 获取状态变更申请记录 |
| DELETE | /api/contracts/:contract_id | 删除合同 |
| GET | /api/contracts/:contract_id/lifecycle | 获取合同生命周期记录 |
| GET | /api/contracts/:contract_id/executions | 获取执行记录 |
| POST | /api/contracts/:contract_id/executions | 创建执行记录 |
| DELETE | /api/executions/:execution_id | 删除执行记录 |
| GET | /api/contracts/:contract_id/documents | 获取文档列表 |
| POST | /api/contracts/:contract_id/documents | 上传文档 |
| DELETE | /api/documents/:document_id | 删除文档 |

**合同状态说明：**
- `draft` - 草稿
- `pending` - 待审批
- `active` - 已生效
- `completed` - 已完成
- `terminated` - 已终止
- `archived` - 已归档

**需要审批的状态变更：**
- `archived` (归档)
- `terminated` (终止)

请求示例：

```bash
# 创建合同
curl -X POST http://localhost:8000/api/contracts \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "contract_no": "CT20240101",
    "title": "采购合同",
    "customer_id": 1,
    "contract_type_id": 1,
    "amount": 100000,
    "currency": "CNY",
    "status": "draft",
    "sign_date": "2024-01-01",
    "start_date": "2024-01-01",
    "end_date": "2024-12-31"
  }'

# 申请状态变更（需要审批的状态）
curl -X POST http://localhost:8000/api/contracts/1/status-change \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "to_status": "archived",
    "reason": "合同已完成，申请归档"
  }'

# 响应示例（需要审批时）：
{
  "direct": false,
  "request": {
    "id": 1,
    "contract_id": 1,
    "from_status": "completed",
    "to_status": "archived",
    "reason": "合同已完成，申请归档",
    "status": "pending"
  }
}

# 响应示例（直接变更时）：
{
  "direct": true,
  "contract": { ... }
}
```

#### 审批管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/contracts/:contract_id/approvals | 获取审批记录 |
| POST | /api/contracts/:contract_id/approvals | 创建审批记录 |
| PUT | /api/approvals/:approval_id | 更新审批状态 |
| GET | /api/pending-approvals | 获取待审批列表 |
| GET | /api/pending-status-changes | 获取待审批状态变更列表 |
| POST | /api/status-change-requests/:request_id/approve | 审批通过状态变更 |
| POST | /api/status-change-requests/:request_id/reject | 拒绝状态变更 |
| GET | /api/notifications/count | 获取待办事项数量 |

**通知数量返回字段：**
- `pendingApprovals` - 待审批合同数量
- `pendingStatusChanges` - 待审批状态变更数量
- `expiringContracts` - 即将到期合同数量
- `total` - 总计

#### 审计日志管理（需要管理员权限）

| 方法 | 路径 | 说明 | 权限要求 |
|------|------|------|----------|
| GET | /api/audit-logs | 获取审计日志列表 | 管理员 |
| DELETE | /api/audit-logs/:id | 删除单个审计日志 | 管理员 |
| POST | /api/audit-logs/batch-delete | 批量删除审计日志 | 管理员 |
| GET | /api/audit-logs/export | 导出审计日志 | 管理员 |

**注意**：审计日志 API 需要管理员权限（`admin` 角色），用于系统安全审计和合规性检查。

#### 生命周期事件类型

| 事件类型 | 说明 |
|----------|------|
| created | 合同创建 |
| submitted | 提交审批 |
| approved | 审批通过 |
| rejected | 审批拒绝 |
| activated | 合同生效 |
| progress | 执行进度更新 |
| payment | 付款记录 |
| completed | 合同完成 |
| terminated | 合同终止 |
| archived | 合同归档 |
| status_changed | 状态变更 |

请求示例：

```bash
# 创建审批
curl -X POST http://localhost:8000/api/contracts/1/approvals \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "pending",
    "comment": "请审批"
  }'

# 更新审批状态
curl -X PUT http://localhost:8000/api/approvals/1 \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "approved",
    "comment": "同意"
  }'
```

#### 提醒管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/contracts/:contract_id/reminders | 获取提醒列表 |
| POST | /api/contracts/:contract_id/reminders | 创建提醒 |
| POST | /api/reminders/:reminder_id/send | 发送提醒 |
| GET | /api/expiring-contracts | 获取即将到期合同 |
| GET | /api/statistics | 获取统计数据 |

## 项目结构

```
AnXin_Contract_Manage/
├── config/                    # 配置模块
│   └── config.go              # 环境配置加载
├── handlers/                  # HTTP 处理器
│   ├── auth.go                # 认证相关（登录、注册、用户管理）
│   ├── customer.go            # 客户管理
│   ├── contract.go            # 合同管理（含生命周期、状态变更）
│   └── approval.go            # 审批与提醒
├── middleware/                # 中间件
│   ├── auth.go                # JWT 认证中间件
│   ├── security.go            # 安全中间件（XSS防护、速率限制）
│   ├── error_handler.go       # 全局错误处理中间件
│   └── validator.go           # 输入验证中间件
├── models/                    # 数据模型
│   └── models.go              # GORM 模型定义（含生命周期、状态变更请求）
├── services/                  # 业务逻辑层
│   ├── user_service.go        # 用户服务
│   ├── customer_service.go    # 客户服务
│   ├── contract_service.go    # 合同服务（含生命周期、归档、状态变更）
│   └── approval_service.go    # 审批服务
├── frontend/                  # 前端项目
│   ├── src/
│   │   ├── api/               # API 接口定义
│   │   ├── components/        # 公共组件
│   │   ├── router/            # 路由配置
│   │   ├── store/             # 状态管理
│   │   ├── utils/             # 工具函数
│   │   └── views/             # 页面组件
│   │       ├── Login.vue      # 登录页面
│   │       ├── Register.vue   # 注册页面
│   │       ├── Layout.vue     # 布局组件（含待办提示）
│   │       ├── Dashboard.vue  # 仪表盘
│   │       ├── Contract.vue   # 合同管理
│   │       ├── ContractDetail.vue  # 合同详情（含生命周期）
│   │       ├── Customer.vue   # 客户管理
│   │       ├── User.vue       # 用户管理
│   │       ├── Approval.vue   # 审批管理
│   │       └── Reminder.vue   # 到期提醒
│   ├── package.json
│   └── vite.config.js
├── main.go                    # 后端入口文件
├── go.mod                     # Go 模块定义
├── go.sum                     # Go 依赖校验
├── .env.example               # 环境变量示例
├── .env                       # 环境变量（需手动创建）
├── init.sql                   # 数据库初始化脚本
├── docker-compose.yml         # Docker Compose 配置
├── Dockerfile                 # 后端 Docker 构建文件
└── README.md                  # 项目文档
```

## 前端页面说明

| 页面 | 文件 | 说明 |
|------|------|------|
| 登录 | Login.vue | 用户登录，验证用户名密码，保存 Token |
| 注册 | Register.vue | 用户注册，填写基本信息 |
| 布局 | Layout.vue | 主框架布局，包含侧边栏导航和顶部栏，支持待办事项红点提示 |
| 仪表盘 | Dashboard.vue | 数据统计、图表展示、即将到期合同 |
| 合同管理 | Contract.vue | 合同增删改查、状态管理 |
| 合同详情 | ContractDetail.vue | 合同详细信息，包含执行跟踪、文档管理、审批记录、生命周期时间线、状态变更、归档操作 |
| 客户管理 | Customer.vue | 客户/供应商信息管理，包含合同类型管理 |
| 用户管理 | User.vue | 用户信息管理、角色分配、用户注册 |
| 审批管理 | Approval.vue | 合同审批流程、审批历史、状态变更审批 |
| 到期提醒 | Reminder.vue | 合同到期提醒管理 |

## 环境变量

| 变量名 | 说明 | 默认值 | 必填 |
|--------|------|--------|------|
| APP_NAME | 应用名称 | 合同管理系统 | 否 |
| APP_VERSION | 版本号 | 1.0.0 | 否 |
| GO_ENV | 运行环境 (development/production) | development | 否 |
| GIN_MODE | Gin 运行模式 (debug/release/test) | debug | 否 |
| MYSQL_HOST | MySQL 主机地址 | localhost | 是 |
| MYSQL_PORT | MySQL 端口 | 3306 | 是 |
| MYSQL_USER | MySQL 用户名 | root | 是 |
| MYSQL_PASSWORD | MySQL 密码 | - | 是 |
| MYSQL_DATABASE | 数据库名称 | contract_manage | 是 |
| SECRET_KEY | JWT 签名密钥 | - | 是 |
| JWT_ALGORITHM | JWT 算法 | HS256 | 否 |
| ACCESS_TOKEN_EXPIRE_MINUTES | Token 过期时间(分钟) | 30 | 否 |
| UPLOAD_DIR | 文件上传目录 | uploads | 否 |
| ADMIN_USERNAME | 超级管理员用户名 | admin | 否 |
| ADMIN_PASSWORD | 超级管理员密码 | Admin@12345 | 否 |
| ADMIN_EMAIL | 超级管理员邮箱 | admin@example.com | 否 |
| AUDIT_ADMIN_USERNAME | 审计管理员用户名 | auditadmin | 否 |
| AUDIT_ADMIN_PASSWORD | 审计管理员密码 | Audit@12345 | 否 |
| AUDIT_ADMIN_EMAIL | 审计管理员邮箱 | audit@example.com | 否 |

### SECRET_KEY 安全建议

- 使用至少 32 位的随机字符串
- 不要使用默认值
- 生产环境定期更换
- 可以使用以下命令生成：
```bash
openssl rand -base64 32
```

### 密码策略要求

系统实施了严格的密码策略：

1. **密码长度**：至少 8 位字符
2. **密码复杂度**：必须包含以下字符类型中的至少三种：
   - 大写字母 (A-Z)
   - 小写字母 (a-z)
   - 数字 (0-9)
   - 特殊字符 (!@#$%^&* 等)

3. **密码示例**：
   - ✅ 合规密码：`Contract@2024`、`Admin#Pass123`、`Test@12345`
   - ❌ 不合规密码：`admin123`（只有2种字符类型）、`12345678`（只有数字）

4. **管理员密码**：
   - 首次部署必须修改默认管理员密码
   - 生产环境禁止使用默认密码 `Admin@12345`
   - 审计管理员密码 `Audit@12345` 也必须修改

### 环境变量安全

1. **生产环境检测**：
   - 设置 `GO_ENV=production` 或 `GIN_MODE=release` 会启用生产环境安全策略
   - 生产环境下禁止使用默认的 JWT 密钥和管理员密码

2. **Docker 环境配置**：
   - 使用环境变量文件或 Docker secrets 管理敏感信息
   - 不要在 `docker-compose.yml` 中硬编码密码

## Docker 操作

### 常用命令

```bash
# 启动所有服务
docker-compose up -d

# 停止所有服务
docker-compose down

# 重启服务
docker-compose restart

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f

# 重建并启动
docker-compose up -d --build

# 进入后端容器
docker exec -it contract_backend sh

# 进入前端容器
docker exec -it contract_frontend sh

# 进入 MySQL 容器
docker exec -it contract_mysql mysql -u contract_user -p contract_manage
```

### 修改配置后重启

修改 `docker-compose.yml` 或环境变量后，需要重建容器：

```bash
docker-compose down
docker-compose up -d --build
```

### 端口说明

| 服务 | 端口 | 说明 |
|------|------|------|
| 前端 | 80 | HTTP 服务 |
| 后端 | 8000 | API 服务 |
| MySQL | 3306 | 数据库服务 |

## 数据备份与恢复

### 备份数据

```bash
# 备份整个数据库
docker exec contract_mysql mysqldump -u contract_user -pcontract123 contract_manage > backup_$(date +%Y%m%d).sql
```

### 恢复数据

```bash
# 从备份文件恢复
cat backup_20240101.sql | docker exec -i contract_mysql mysql -u contract_user -pcontract123 contract_manage
```

### 数据持久化

Docker Compose 配置了 MySQL 数据卷 `mysql_data`，数据会持久化存储。即使删除容器，重新启动后数据依然存在。

如需完全清除数据：

```bash
docker-compose down -v
```

## 安全更新 (最新版本)

本版本对系统进行了全面的安全加固，修复了多项安全漏洞：

### 已修复的安全问题

1. **认证授权加固**
   - 用户管理 API 和审计日志 API 添加了管理员权限检查
   - 只有管理员才能访问敏感管理功能，防止越权访问

2. **JWT 密钥安全**
   - 生产环境强制要求更换默认 JWT 密钥
   - 生产环境强制要求更换默认管理员密码
   - 通过环境变量 `GO_ENV` 或 `GIN_MODE` 检测生产环境

3. **文件上传安全增强**
   - **移除HTML/HTM文件类型**：从白名单中移除`.html`和`.htm`，防止XSS攻击
   - **双重扩展名检测**：阻止`test.php.jpg`等恶意文件名格式
   - **文件内容验证**：扫描文件内容，检测PHP标签、JavaScript代码和WebShell关键字
   - **严格文件类型白名单**：只允许安全的文档、图片和文本文件类型

4. **SQL 注入防护**
   - 所有数据库查询使用参数化查询（GORM `?` 占位符）
   - 添加输入参数验证和边界检查
   - LIKE查询中的通配符（`%`、`_`）转义，防止LIKE注入
   - 状态参数白名单验证，防止非法状态值注入

5. **XSS 防护修复**
   - 修复XSS保护中间件的响应处理问题
   - 确保所有响应内容（包括HTML和JSON）都能正确写回客户端
   - 自动转义所有JSON响应中的字符串值
   - 实现递归清理嵌套JSON数据

6. **输入验证强化**
   - 所有查询参数添加验证：长度限制、格式检查、边界值检查
   - 日期参数格式验证（YYYY-MM-DD格式）
   - 关键词参数长度限制和特殊字符过滤
   - 分页参数边界检查（最小值、最大值限制）

7. **敏感信息保护**
   - 从 Docker 镜像中移除 `.env` 文件
   - Docker Compose 使用环境变量引用，不再硬编码密码
   - 支持通过环境变量设置敏感配置

8. **密码策略加强**
   - 最小密码长度从 6 位增加到 8 位
   - 添加密码复杂度要求：必须包含大写字母、小写字母、数字和特殊字符中的至少三种
   - 恢复严格的密码策略要求（至少3种字符类型）

9. **错误处理改进**
   - 添加全局错误处理中间件
   - 统一的错误响应格式
   - 添加 panic 恢复和堆栈跟踪记录
   - 提供用户友好的错误信息

### 安全最佳实践

1. **生产环境部署前必须**
   - 修改默认 JWT 密钥（至少 32 位随机字符串）
   - 修改默认管理员和审计管理员密码
   - 使用 HTTPS 部署
   - 配置防火墙规则

2. **定期维护**
   - 定期更新 JWT 密钥
   - 定期备份数据库
   - 监控系统日志和安全事件

3. **密码策略**
   - 密码长度至少 8 位
   - 必须包含大写字母、小写字母、数字和特殊字符中的至少三种
   - 禁止使用常见弱密码

### 安全测试验证

1. **文件上传安全测试**
   - PHP文件上传被阻止（返回"不支持的文件类型: php"）
   - 双重扩展名文件被阻止（返回"检测到恶意文件名格式"）
   - HTML文件上传被阻止（返回"不支持的文件类型: html"）
   - 合法PDF文件上传成功

2. **SQL注入测试**
   - 参数化查询防止SQL注入
   - LIKE通配符转义防止LIKE注入
   - 状态参数白名单验证

3. **XSS防护测试**
   - API调试页面正常显示（Content-Length: 59360）
   - JSON响应中的字符串值自动转义
   - HTML响应正确返回内容

## 安全说明

- **密码加密**：用户密码使用 bcrypt 加密存储，不可逆
- **认证机制**：除登录注册外，所有 API 需要有效的 JWT Token
- **Token 时效**：Token 默认 30 分钟过期，需要重新登录
- **权限控制**：管理员权限 API 需要 `admin` 角色
- **文件上传安全**：
  - 文件类型白名单验证（PDF、Word、Excel、图片、文本）
  - 双重扩展名检测（阻止test.php.jpg等恶意文件名）
  - 文件内容恶意代码扫描（PHP标签、JavaScript、WebShell关键字）
  - 路径遍历防护
- **SQL 注入防护**：
  - 参数化查询（GORM `?` 占位符）
  - 输入参数验证和边界检查
  - LIKE通配符转义
  - 状态参数白名单验证
- **XSS 防护**：
  - 自动转义所有JSON响应中的字符串值
  - 修复XSS保护中间件响应处理问题
  - 确保所有响应内容正确写回客户端
- **输入验证**：
  - 所有查询参数添加验证（长度、格式、边界值）
  - 日期参数格式验证（YYYY-MM-DD）
  - 关键词参数长度限制和特殊字符过滤
- **错误处理**：统一的错误响应，隐藏内部实现细节
- **生产环境建议**：
  - 修改默认的 SECRET_KEY
  - 修改默认的管理员密码
  - 使用 HTTPS 部署
  - 配置防火墙规则
  - 定期备份数据库
  - 定期更新 JWT 密钥

## 常见问题

### 0. 安全相关常见问题

#### Q: 为什么提示"SECRET_KEY cannot be the default value in production environment"？
**A**: 在生产环境中（`GO_ENV=production` 或 `GIN_MODE=release`），系统禁止使用默认的 JWT 密钥。请修改 `.env` 文件中的 `SECRET_KEY` 为至少 32 位的随机字符串。

#### Q: 为什么上传文件时提示"不支持的文件类型"？
**A**: 系统实施了文件类型白名单验证，只允许上传特定类型的文件（PDF、Word、Excel、图片、文本）。HTML文件已被移除以防止XSS攻击。请确保上传的文件扩展名在允许列表中。

#### Q: 为什么上传文件时提示"检测到恶意文件名格式"？
**A**: 系统检测到双重扩展名（如test.php.jpg），这是恶意文件上传的常见手法。请使用安全的文件名格式，避免使用.php、.jsp、.asp等脚本扩展名。

#### Q: 上传的文件会被检查内容吗？
**A**: 是的，系统会扫描文件内容，检测PHP标签、JavaScript代码和WebShell关键字。如果检测到恶意内容，文件上传将被拒绝。

#### Q: 为什么访问用户管理 API 返回 403 Forbidden？
**A**: 用户管理 API 需要管理员权限。只有 `admin` 角色的用户才能访问。请使用管理员账号登录。

#### Q: 密码有什么要求？
**A**: 密码长度至少 8 位，必须包含大写字母、小写字母、数字和特殊字符中的至少三种。

### 1. 启动失败，提示数据库连接失败

检查：
- MySQL 服务是否启动
- `.env` 中的数据库配置是否正确
- 数据库用户是否有权限访问数据库

### WSL 连接 Windows MySQL

如果后端运行在 WSL 中，数据库在 Windows 上，配置如下：

1. 获取 Windows IP（在 WSL 中执行）：
   ```bash
   cat /etc/resolv.conf
   ```

2. 在 `.env` 中配置：
   ```env
   MYSQL_HOST=172.x.x.x  # 从上面获取的 IP
   ```

3. 确保 Windows 防火墙允许 MySQL 端口（3306）

### 2. 前端无法访问后端 API

检查：
- 后端服务是否正常运行（http://localhost:8000/health）
- 前端 `vite.config.js` 中的代理配置是否正确
- 防火墙是否允许对应端口

### 3. Token 过期后怎么办

Token 过期后，前端会自动跳转到登录页面，需要重新登录。

### 4. 如何修改管理员权限

数据库中修改用户的 `role` 字段为 `admin`：

```sql
UPDATE users SET role = 'admin' WHERE username = 'admin';
```

### 5. 上传文件大小限制

默认限制 10MB，如需修改请编辑 `docker-compose.yml` 或 Nginx 配置。

### 6. 如何查看后端 API 文档

后端未集成 Swagger，可参考本文档的 API 端点说明。

## 变更日志

### v1.2.0 (安全增强版本) - 2026-03-23

#### 安全修复
- ✅ 修复API调试页面白页问题：修复XSS保护中间件响应处理，确保所有响应正确写回
- ✅ 增强SQL注入防护：添加输入参数验证、边界检查、LIKE通配符转义
- ✅ 增强文件上传安全：移除HTML/HTM文件类型、添加双重扩展名检测、文件内容验证
- ✅ 修复XSS保护中间件：确保HTML和JSON响应都能正确处理
- ✅ 恢复严格密码策略：密码必须包含至少3种字符类型
- ✅ 增强输入验证：所有查询参数添加验证（长度、格式、边界值）

#### 新增安全功能
- 🛡️ 文件内容恶意代码检测（PHP标签、JavaScript、WebShell关键字）
- 🛡️ 双重扩展名检测（阻止test.php.jpg等恶意文件名）
- 🛡️ SQL查询参数验证和边界检查
- 🛡️ 日期参数格式验证（YYYY-MM-DD）
- 🛡️ 关键词参数长度限制和特殊字符过滤

#### 改进
- 🔧 优化XSS保护中间件，修复响应内容丢失问题
- 🔧 完善输入验证，覆盖所有API端点
- 🔧 更新默认管理员密码，使用更安全的默认值
- 🔧 完善安全测试文档和验证步骤

### v1.1.0 (安全加固版本) - 2026-03-23

#### 安全修复
- ✅ 修复认证授权缺陷：用户管理和审计日志 API 添加管理员权限检查
- ✅ 修复 JWT 密钥风险：生产环境强制更换默认密钥和密码
- ✅ 修复文件上传安全：添加文件类型白名单和路径遍历防护
- ✅ 修复 XSS 防护不完整：添加全局 XSS 防护中间件
- ✅ 修复敏感信息泄露：从 Docker 镜像中移除 .env 文件，使用环境变量引用
- ✅ 修复输入验证不足：加强密码策略，最小长度 8 位，要求复杂度
- ✅ 修复错误处理不完善：添加全局错误处理中间件，统一错误响应格式

#### 新增功能
- 🔒 添加 `middleware/error_handler.go` 全局错误处理中间件
- 🔒 添加 `middleware/security.go` 中的 XSS 防护功能
- 🔒 添加环境变量 `GO_ENV` 和 `GIN_MODE` 用于生产环境检测
- 🔒 添加审计管理员账号配置 `AUDIT_ADMIN_*`

#### 改进
- 🔧 重构文件上传函数，增强安全性
- 🔧 优化错误日志记录和用户友好错误信息
- 🔧 更新 Docker 配置，提高安全性
- 🔧 完善文档，添加安全最佳实践

### v1.0.0 (初始版本) - 2024-01-01

#### 核心功能
- 🎉 用户权限管理（注册、登录、角色管理）
- 🎉 客户/供应商管理
- 🎉 合同全生命周期管理
- 🎉 审批流程管理
- 🎉 文档管理和文件上传
- 🎉 统计报表和提醒功能
- 🎉 前后端分离架构

## 许可证

MIT License
