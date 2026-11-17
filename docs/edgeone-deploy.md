# EdgeOne 部署与 NexusMC 对接说明

本仓库支持以「静态前端 + EdgeOne Cloud Functions (Go)」的形态整体部署到腾讯云 EdgeOne：

```
ysmmc/
├── frontend/            # Vue 3 + Vite 静态站（构建产物 frontend/dist）
├── cloud-functions/     # EdgeOne Go 云函数（Gin，框架模式，入口 api.go）
│   ├── api.go           # 入口：入口文件名为 api.go，函数挂载在站点 /api 前缀下
│   ├── nexusmc/         # NexusMC 站点 API / OAuth2 客户端（仅服务端使用密钥）
│   └── go.mod
└── edgeone.json         # EdgeOne 构建与 SPA 回退配置
```

原有的 Go 后端（`backend/`，api.ysmmc.cn）不受影响；EdgeOne 上的云函数只负责 NexusMC
对接，不承担本站自身业务。

## 云函数提供的接口

入口 `api.go` 挂载在 `/api` 前缀下，即线上访问路径为 `/api/nexusmc/...`：

| 路由 | 说明 |
| --- | --- |
| `GET /api/nexusmc/health` | 存活检查，返回凭据配置状态 |
| `GET /api/nexusmc/config` | 前端启动时探测凭据是否已配置 |
| `GET /api/nexusmc/catalog` | 代理 NexusMC 站点 API 资源筛选目录（缓存 10 分钟） |
| `GET /api/nexusmc/resources` | 代理资源列表，筛选参数白名单透传（缓存 60 秒） |
| `GET /api/nexusmc/resource?id=` | 代理资源详情 |
| `GET /api/nexusmc/search?q=` | 代理全站搜索（仅 resources 类型，需密钥开通 search 权限） |
| `GET /api/nexusmc/auth/start` | 发起 NexusMC OAuth2 授权（写 state cookie 后 302） |
| `GET /api/nexusmc/auth/callback` | OAuth2 回调：state 校验 → 换 token → 拉取 userinfo → 写 httpOnly token cookie → 302 回前端 `/nexusmc/callback` |
| `POST /api/nexusmc/auth/logout` | 清除 token cookie（前端退出登录时调用） |
| `GET /api/nexusmc/me/content` | 我的投稿（需登录，type/status/page/pageSize） |
| `GET /api/nexusmc/me/favorites` · `/me/likes` | 我的收藏 / 点赞（需登录） |
| `GET /api/nexusmc/me/notifications` · `/unread-count` | 通知列表 / 未读数（需登录） |
| `POST /api/nexusmc/me/notifications/read-all` · `/:id/read` | 标记已读（需登录） |

密钥（站点 API Key、OAuth client_secret）只保存在 EdgeOne 环境变量中，浏览器永远
接触不到——这是 NexusMC 开放平台协议的要求（见
docs/reference/nexusmc-open-platform/api/site-api-keys.md 与
site-oauth2-provider.md）。

## NexusMC 侧准备

1. **站点 API 密钥**：需已认证的用户在 NexusMC「开放平台控制台 → 站点 API」创建，
   绑定 OAuth 应用并勾选业务域：`resource:read:public`（必须）、
   `search:read:public`（搜索功能需要）。
2. **OAuth2 应用**：在 NexusMC 后台「站点 API 管理 → OAuth 应用」创建，登记回调地址：

   ```
   https://<你的 EdgeOne 域名>/api/nexusmc/auth/callback
   ```

   `redirect_uri` 必须与登记值逐字一致。 scopes 建议 `user:basic user:email`
   加上个人主页需要的 `user:content:read user:interaction:read
   user:notification:read user:notification:write`（缺少这些权限时个人主页
   接口会报 403 `OAUTH_SCOPE_REQUIRED`，页面会引导重新授权）。
   默认请求的 scope 可用环境变量 `NEXUSMC_OAUTH_SCOPES` 覆盖。

## EdgeOne 环境变量

在 Makers/EdgeOne Pages 控制台「项目设置 → 环境变量」中配置：

| 变量 | 必填 | 说明 |
| --- | --- | --- |
| `NEXUSMC_SITE_API_KEY` | 展示资源必填 | 站点 API 密钥（`avm_sak_...` 完整值） |
| `NEXUSMC_OAUTH_CLIENT_ID` | 登录必填 | OAuth 应用 client_id |
| `NEXUSMC_OAUTH_CLIENT_SECRET` | 登录必填 | OAuth 应用完整密钥（`avm_oac_...`，不要截断） |
| `NEXUSMC_OAUTH_REDIRECT_URI` | 建议 | 回调地址；不设则按请求 Host 自动推导，正式环境建议显式设置 |
| `NEXUSMC_FRONTEND_ORIGIN` | 可选 | 回调后跳回的前端地址；不设则用请求自身来源 |
| `NEXUSMC_OAUTH_SCOPES` | 可选 | 默认 `user:content:read user:interaction:read user:notification:read user:notification:write`；应用未登记这些权限时可用它改回空串（个人主页将不可用但登录不受影响） |
| `NEXUSMC_RESOURCE_PLATFORM` | 可选 | 列表默认平台筛选（如 `java`） |
| `NEXUSMC_RESOURCE_CATEGORY` | 可选 | 列表默认分类筛选（如 YSM 模型对应的分类值，取自 catalog 接口） |
| `NEXUSMC_SITE_ORIGIN` | 可选 | NexusMC 站点地址，默认 `https://www.nexusmc.cn` |

未配置凭据时接口返回 503 `not_configured`。缺少凭据、密钥失效或业务权限不足时，
页面会提示服务暂不可用，并提供 NexusMC 资源库入口。

### 资源接口返回 401 / 403

`/health` 的 `site_api: true` 和 `/config` 的 `configured: true` 仅说明环境变量非空，
不代表密钥已通过上游验证。检查 `/api/nexusmc/resources` 或 `/api/nexusmc/catalog`
响应体中的 `code`：

- `401 SITE_API_KEY_INVALID`：NexusMC 拒绝了服务端站点密钥。到 NexusMC 开放平台
  检查密钥是否启用、到期或被撤销；必要时创建或轮换密钥，将完整值更新到 EdgeOne
  对应部署环境的 `NEXUSMC_SITE_API_KEY`，然后重新部署云函数。完整密钥仅在创建或
  轮换时显示一次，请勿使用被遮盖的展示值、OAuth 应用密钥或用户 token。
- `403 SITE_API_KEY_SCOPE_REQUIRED`：站点密钥缺少业务权限。为列表、目录和详情开启
  `resource:read:public`，为搜索开启 `search:read:public`。

站点密钥错误无法通过用户重新登录解决。密钥只在服务端配置，不要写进前端、仓库、
日志或聊天；部署完成后再次请求上述接口确认返回 200。代码会去除密钥两端的空白，
但无法修复已过期、已撤销或填写错误的密钥。

## 部署步骤

1. 仓库推送到 GitHub/Gitee 后，在 EdgeOne 控制台「导入 Git 仓库」创建项目；
2. 构建配置由根目录 `edgeone.json` 提供：构建 `frontend/`、输出
   `frontend/dist`、SPA 回退 `/* → /index.html`；`cloud-functions/` 会被自动识别
   并以 Go 1.26 运行时交叉编译部署；
3. 在项目设置中配置上面的环境变量；
4. 绑定自定义域名（如需中国大陆加速，域名必须已完成 ICP 备案），并在 NexusMC
   OAuth 应用中把登记的回调地址更新为绑定域名下的
   `/api/nexusmc/auth/callback`。

## 本地开发

```bash
# 终端 1：本地运行云函数（需安装 EdgeOne CLI：npm i -g edgeone）
edgeone makers dev          # 端口以 CLI 输出为准，假设 8787

# 终端 2：前端开发服务器（vite 已把 /api/nexusmc 代理到云函数）
NEXUSMC_DEV_TARGET=http://127.0.0.1:8787 pnpm --dir frontend dev
```

不启动云函数时，`/nexusmc` 页面仍会渲染，并提示资源加载失败。

## 已知限制

- NexusMC 未公开资源列表的完整字段定义，`cloud-functions/nexusmc` 对常见字段名
  做了宽松映射（title/name、cover/cover_image_url 等），并为每个条目补充
  `page_url`（`https://www.nexusmc.cn/resources/{slug}`）；待拿到真实密钥联调后
  可再收紧。
- 全站搜索需要站点 API 密钥额外开通 `search:read:public` 业务域，否则返回 403。
- OAuth access token 登录成功后以 httpOnly cookie（`nexusmc_token`，7 天，
  仅 `/api/nexusmc` 路径可读）保存在浏览器侧，由云函数代理访问 `/me/*` 时附带；
  令牌不写入 localStorage，前端 JS 接触不到。未保存 refresh token，过期后
  个人主页接口返回 401，页面引导重新登录。
