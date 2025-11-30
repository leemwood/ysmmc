# YSM 模型站 (YSM Model Hub)

一个专注于分享和下载 Minecraft Yes Steve Model (YSM) 玩家模型的社区平台。

## 📖 项目简介

YSM 模型站是一个现代化、响应式的 Web 应用，旨在为 Minecraft 玩家提供一个便捷的平台来浏览、上传和下载 YSM 模型。用户可以注册账户，分享自己的创作，收藏喜欢的模型，并与其他玩家互动。管理员拥有审核机制，确保平台内容的质量。

## ✨ 功能特性

- **模型浏览**：支持分页浏览、关键词搜索模型。
- **模型详情**：展示模型预览图、Markdown 描述、标签、下载量等详细信息。
- **用户系统**：
  - 用户注册与登录。
  - 个人中心：编辑个人资料（头像、昵称、简介）。
  - 公开主页：展示用户上传的模型和收藏列表。
- **模型管理**：
  - 用户可以上传、编辑和删除自己的模型。
  - 支持模型文件 (.ysm) 和预览图上传。
- **收藏功能**：用户可以收藏自己喜欢的模型，并在个人主页查看。
- **管理员后台**：
  - 模型审核流程（批准/拒绝）。
  - 用户资料变更审核。
- **响应式设计**：完美适配桌面端和移动端设备。

## 🛠️ 技术栈

- **前端框架**：[Vue 3](https://vuejs.org/) (Composition API)
- **开发语言**：[TypeScript](https://www.typescriptlang.org/)
- **构建工具**：[Vite](https://vitejs.dev/)
- **状态管理**：[Pinia](https://pinia.vuejs.org/)
- **路由管理**：[Vue Router](https://router.vuejs.org/)
- **样式方案**：SCSS (BEM 命名规范 + 模块化管理)
- **图标库**：[Lucide Vue Next](https://lucide.dev/)
- **后端服务**：[Supabase](https://supabase.com/) (Auth, Database, Storage, RLS)

## 📂 目录结构

```
src/
├── assets/          # 静态资源
├── components/      # 公共组件
├── router/          # 路由配置
├── stores/          # Pinia 状态管理
├── styles/          # 全局样式与主题变量
├── supabase/        # Supabase 客户端与 SQL 定义
├── types/           # TypeScript 类型定义
├── views/           # 页面视图
│   ├── AdminDashboardView.vue   # 管理员看板
│   ├── HomeView.vue             # 首页
│   ├── ModelDetailView.vue      # 模型详情页
│   ├── UserPublicProfileView.vue # 用户公开主页
│   └── ...
├── App.vue          # 根组件
└── main.ts          # 入口文件
```

## 🚀 安装与运行

### 前置要求

- Node.js (推荐 v16+)
- npm 或 pnpm

### 步骤

1. **克隆项目**

   ```bash
   git clone <repository-url>
   cd ysmmc
   ```

2. **安装依赖**

   ```bash
   npm install
   ```

3. **配置环境变量**

   在项目根目录下创建 `.env` 文件，并填入你的 Supabase 配置：

   ```env
   VITE_SUPABASE_URL=your_supabase_url
   VITE_SUPABASE_ANON_KEY=your_supabase_anon_key
   ```

4. **启动开发服务器**

   ```bash
   npm run dev
   ```

5. **构建生产版本**

   ```bash
   npm run build
   ```

## 📝 数据库与存储配置

项目依赖 Supabase 进行数据存储。主要的表结构包括：

- `profiles`: 用户资料表
- `models`: 模型数据表
- `favorites`: 收藏关联表

*注：详细的 SQL 建表语句位于 `src/supabase/` 目录下。*

## 📄 许可证

MIT License
