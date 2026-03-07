# YSM 模型站后端服务

基于 Go + Gin + GORM + PostgreSQL 的模型分享平台后端服务。

## 目录结构

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # 服务入口
├── internal/
│   ├── config/                  # 配置管理
│   ├── database/                # 数据库连接和迁移
│   ├── handler/                 # HTTP 处理器
│   ├── middleware/              # 中间件
│   ├── model/                   # 数据模型
│   ├── repository/              # 数据访问层
│   ├── router/                  # 路由定义
│   └── service/                 # 业务逻辑层
├── migrations/                  # 数据库迁移脚本
├── pkg/                         # 公共包
│   ├── auth/                    # 认证相关
│   ├── email/                   # 邮件服务
│   ├── response/                # 响应格式
│   └── utils/                   # 工具函数
├── templates/                   # 邮件模板
├── .env.example                 # 环境变量示例
├── go.mod
└── go.sum
```

## 快速开始

### 环境要求

- Go 1.21+
- PostgreSQL 14+

### 安装依赖

```bash
cd backend
go mod download
```

### 配置环境变量

```bash
cp .env.example .env
# 编辑 .env 文件，设置必要的环境变量
```

### 必要的环境变量

| 变量名 | 说明 | 示例 |
|--------|------|------|
| `DB_PASSWORD` | 数据库密码 | `your_secure_password` |
| `JWT_SECRET` | JWT 密钥（至少 32 字符） | `your_very_long_jwt_secret_key` |

### 运行服务

```bash
go run cmd/server/main.go
```

服务将在 `http://localhost:8080` 启动。

## API 路由

### 公开路由

| 路由 | 方法 | 说明 |
|------|------|------|
| `/api/auth/register` | POST | 用户注册 |
| `/api/auth/login` | POST | 用户登录 |
| `/api/auth/refresh` | POST | 刷新 Token |
| `/api/auth/forgot-password` | POST | 忘记密码 |
| `/api/auth/reset-password` | POST | 重置密码 |
| `/api/auth/verify` | GET | 邮箱验证 |
| `/api/models` | GET | 模型列表 |
| `/api/models/:id` | GET | 模型详情 |
| `/api/models/:id/download` | POST | 模型下载 |
| `/api/users/:id` | GET | 用户公开信息 |
| `/api/announcements` | GET | 公告列表 |
| `/health` | GET | 健康检查 |

### 认证路由

| 路由 | 方法 | 说明 |
|------|------|------|
| `/api/auth/me` | GET | 当前用户信息 |
| `/api/auth/logout` | POST | 登出 |
| `/api/auth/change-email` | POST | 更改邮箱 |
| `/api/users/me` | GET/PUT | 个人信息管理 |
| `/api/users/me/password` | PUT | 修改密码 |
| `/api/models` | POST | 创建模型 |
| `/api/models/:id` | PUT/DELETE | 更新/删除模型 |
| `/api/models/:id/favorite` | POST/DELETE/GET | 收藏管理 |
| `/api/favorites` | GET | 收藏列表 |
| `/api/upload/model` | POST | 上传模型文件 |
| `/api/upload/image` | POST | 上传图片 |

### 管理员路由

| 路由 | 方法 | 说明 |
|------|------|------|
| `/api/admin/stats` | GET | 统计数据 |
| `/api/admin/models/pending` | GET | 待审核模型 |
| `/api/admin/models/:id/approve` | PUT | 批准模型 |
| `/api/admin/models/:id/reject` | PUT | 拒绝模型 |
| `/api/admin/users` | GET | 用户列表 |
| `/api/admin/users/:id/role` | PUT | 更新用户角色 |
| `/api/admin/users/:id/ban` | PUT | 封禁用户 |
| `/api/admin/users/:id/unban` | PUT | 解封用户 |
| `/api/admin/announcements` | POST/PUT/DELETE | 公告管理 |

## 安全特性

### 认证授权

- JWT Token 认证
- Access Token 有效期 24 小时
- Refresh Token 有效期 7 天，存储在数据库中
- 支持主动吊销 Refresh Token

### 速率限制

| 接口 | 限制 |
|------|------|
| 登录 | 每 IP 每分钟 5 次 |
| 注册 | 每 IP 每小时 3 次 |
| 忘记密码 | 每 IP 每小时 3 次 |

### 文件上传安全

- 文件类型白名单
- 文件魔数（Magic Number）验证
- 文件名长度限制
- UUID 文件名防止覆盖
- 磁盘空间检查

### 输入验证

- XSS 过滤
- 邮箱格式验证
- 用户名格式验证（支持中文）
- 密码长度验证

## 存储服务

### 目录结构

```
uploads/
├── models/          # 模型文件
│   └── 2026-03/    # 按日期分区（可选）
├── images/          # 图片文件
│   └── 2026-03/    # 按日期分区（可选）
└── temp/            # 临时文件
```

### 配置项

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `UPLOAD_PATH` | `./uploads` | 上传根目录 |
| `MAX_FILE_SIZE` | `104857600` | 最大文件大小（100MB） |
| `MAX_DISK_USAGE` | `90` | 最大磁盘使用率百分比 |
| `ENABLE_DATE_PARTITION` | `false` | 是否按日期分目录 |

## 测试

### 运行所有测试

```bash
go test ./... -v -cover
```

### 运行特定测试

```bash
go test ./pkg/auth/... -v
go test ./internal/service/... -v
```

### 测试覆盖率目标

| 模块 | 目标覆盖率 |
|------|-----------|
| pkg/auth | 90% |
| pkg/utils | 85% |
| internal/middleware | 80% |
| internal/service | 75% |

## 环境变量完整列表

```env
# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=ysmmc

# JWT 配置
JWT_SECRET=your-jwt-secret-key-at-least-32-characters
JWT_EXPIRE_HOURS=24
JWT_REFRESH_EXPIRE_DAYS=7

# 服务器配置
SERVER_PORT=8080
GIN_MODE=debug

# 文件上传配置
UPLOAD_PATH=./uploads
MAX_FILE_SIZE=104857600
MAX_DISK_USAGE=90
ENABLE_DATE_PARTITION=true

# SMTP 配置（可选）
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USER=your_email@example.com
SMTP_PASSWORD=your_smtp_password
SMTP_FROM=YSM模型站 <your_email@example.com>

# 前端 URL
FRONTEND_URL=http://localhost:5173

# CORS 允许的域名（逗号分隔）
ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000
```

## 开发指南

### 添加新的 API

1. 在 `internal/model/` 定义数据模型
2. 在 `internal/repository/` 创建数据访问层
3. 在 `internal/service/` 创建业务逻辑层
4. 在 `internal/handler/` 创建 HTTP 处理器
5. 在 `internal/router/router.go` 注册路由

### 代码规范

- 使用简体中文注释
- 遵循 Go 命名规范
- 错误信息使用中文
- 敏感信息不记录到日志

## 更新日志

### 2026-03-06

**安全修复**
- 移除硬编码凭证
- 修改默认管理员密码机制
- 收紧 CORS 配置
- 修复文件路径泄露
- 添加速率限制
- 修复权限提升风险
- 完善 Token 过期机制
- 增强文件上传安全
- 添加输入验证

**存储优化**
- 创建存储服务层
- 启动时目录初始化
- 跨平台磁盘空间检查
- 支持按日期分目录存储
- 临时文件清理机制

## 许可证

MIT License
