# YSM 模型站

非营利性公益模型分享平台，仅供学习交流使用。

## 技术栈

### 后端
- Go 1.21+
- Gin Web Framework
- GORM ORM
- PostgreSQL
- JWT 认证

### 前端
- Vue 3 + TypeScript
- Vite
- Pinia 状态管理
- Vue Router
- Tailwind CSS + shadcn-vue

## 快速开始

### 环境要求
- Go 1.21+
- Node.js 18+
- PostgreSQL 14+

### 数据库配置

```bash
# 创建数据库
createdb ysmmc

# 运行迁移
psql -d ysmmc -f backend/migrations/001_init.sql
```

### 后端启动

```bash
cd backend
cp .env.example .env
# 编辑 .env 配置数据库连接
go mod tidy
go run cmd/server/main.go
```

### 前端启动

```bash
cd frontend
npm install
npm run dev
```

### 访问地址
- 前端: http://localhost:5173
- 后端 API: http://localhost:8080/api

## 默认账号

- 邮箱: admin@ysmmc.local
- 密码: admin123

## 功能特性

- 用户注册/登录
- 邮箱验证
- 密码重置
- 模型上传/下载
- 模型收藏
- 管理员审核
- 公告系统

## 项目结构

```
ysmmc/
├── backend/                 # Go 后端
│   ├── cmd/server/         # 程序入口
│   ├── internal/           # 内部代码
│   │   ├── config/        # 配置管理
│   │   ├── handler/       # HTTP 处理器
│   │   ├── middleware/    # 中间件
│   │   ├── model/         # 数据模型
│   │   ├── repository/    # 数据访问层
│   │   ├── router/        # 路由配置
│   │   └── service/       # 业务逻辑层
│   ├── pkg/               # 公共包
│   ├── migrations/        # 数据库迁移
│   └── templates/         # 模板文件
│
└── frontend/               # Vue 前端
    ├── src/
    │   ├── assets/        # 静态资源
    │   ├── components/    # 组件
    │   ├── lib/           # 工具库
    │   ├── router/        # 路由配置
    │   ├── stores/        # 状态管理
    │   ├── types/         # 类型定义
    │   └── views/         # 页面视图
    └── package.json
```

## 许可证

本项目仅供学习交流使用。
