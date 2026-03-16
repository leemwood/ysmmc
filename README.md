# YSM 模型站 (YSMMC)

非营利性公益模型分享平台，仅供学习交流使用。

## 技术栈

### 后端
- **语言**: Go 1.23+
- **框架**: [Gin Web Framework](https://gin-gonic.com/)
- **ORM**: [GORM](https://gorm.io/)
- **数据库**: PostgreSQL 14+
- **认证**: JWT (JSON Web Tokens) with Refresh Token support
- **依赖**: `godotenv`, `crypto`, `uuid`

### 前端
- **框架**: Vue 3.5+ (Composition API)
- **构建工具**: Vite 6+
- **样式**: [Tailwind CSS 4](https://tailwindcss.com/)
- **组件库**: [Radix Vue](https://www.radix-vue.com/) + [Lucide Vue Next](https://lucide.dev/)
- **状态管理**: Pinia 3+
- **路由**: Vue Router 5+
- **网络请求**: Axios

## 功能特性

### 👤 用户系统
- 注册、登录及登出
- 基于 JWT 的身份验证（支持 Token 刷新）
- 邮箱验证及邮箱更换流程
- 找回密码及重置密码
- 个人资料管理（支持头像及基本信息修改，需审核）
- 公开个人主页展示

### 📦 模型管理
- 模型上传、编辑及删除
- **多版本支持**: 每个模型可以拥有多个版本，支持设置当前主版本
- **多图展示**: 支持上传多张展示图，并可自定义排序
- 模型收藏系统
- 模型下载统计
- 搜索与分类列表

### 🛡️ 管理员功能
- **仪表盘**: 实时查看系统统计数据
- **审核系统**: 审核新发布的模型、模型更新以及用户资料修改
- **用户管理**: 调整用户角色（管理员/普通用户）、封禁或解封违规用户
- **公告系统**: 发布、编辑及删除全站公告
- **超级管理员**: 拥有设置和移除其他管理员的最高权限

## 快速开始

### 环境要求
- Go 1.23+
- Node.js 20+
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
# 编辑 .env 配置数据库连接及 JWT_SECRET (必填)
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

## 环境变量说明 (backend/.env)

| 变量名 | 说明 | 示例 |
| :--- | :--- | :--- |
| `DB_*` | 数据库连接配置 | - |
| `JWT_SECRET` | JWT 签名密钥 (建议 32 位以上) | `your-secret-key` |
| `MAX_FILE_SIZE` | 最大上传限制 (字节) | `104857600` (100MB) |
| `ENABLE_DATE_PARTITION` | 上传文件是否按日期分区存储 | `true` |
| `ALLOWED_ORIGINS` | CORS 允许跨域的域名 | `http://localhost:5173` |

## 项目结构

```
ysmmc/
├── backend/                 # Go 后端
│   ├── cmd/server/         # 程序入口
│   ├── internal/           # 内部代码
│   │   ├── config/        # 配置管理
│   │   ├── handler/       # HTTP 处理器 (Controller)
│   │   ├── middleware/    # 中间件 (Auth, CORS, RateLimit)
│   │   ├── model/         # GORM 数据模型
│   │   ├── repository/    # 数据库操作层
│   │   ├── router/        # 路由定义
│   │   └── service/       # 业务逻辑层
│   ├── migrations/        # SQL 迁移文件
│   └── templates/         # 邮件等 HTML 模板
│
└── frontend/               # Vue 前端
    ├── src/
    │   ├── assets/        # 静态资源
    │   ├── components/    # 公共组件 & UI 组件
    │   ├── lib/           # 工具函数 & API 封装
    │   ├── router/        # 路由配置
    │   ├── stores/        # Pinia 状态管理
    │   ├── types/         # TypeScript 类型定义
    │   ├── utils/         # 通用工具类
    │   └── views/         # 页面视图
    └── package.json
```

## 交流群

点击链接加入群聊【YSM 免费模型站用户群】：[https://qm.qq.com/q/SUKmYH7RyW](https://qm.qq.com/q/SUKmYH7RyW)

## 许可证

本项目基于 [GPL-3.0](LICENSE) 协议开源。
