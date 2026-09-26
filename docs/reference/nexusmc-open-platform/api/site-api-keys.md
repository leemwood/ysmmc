---
title: "站点 API：概况"
source: https://docs.nexusmc.cn/docs/api/site-api-keys
collected: 2026-09-26
---

# 站点 API：概况

站点 API 的鉴权方式、功能授权与站点概况接口。

站点 API 面向官网、机器人、数据看板和其他可信服务端程序。管理员在“管理面板 - 站点 API 管理 - API 密钥”创建密钥时，必须先选择一个已创建的 OAuth2 应用进行绑定，再为每把密钥分别勾选可用业务域。

它与[个人 API Token](/docs/api/personal-api-tokens/)不同：站点 API 密钥不代表任何用户，不能投稿、修改内容或访问后台，只能读取已授权的公开站点数据。

## Base URL 与鉴权

```
https://www.nexusmc.cn/api/site/v1
```

每次请求使用传统 API Key 或[公钥请求签名](/docs/api/open-platform-public-key-auth/)鉴权。传统方式如下：

```
X-API-Key: avm_sak_xxxxx_xxxxx
Accept: application/json
```

```
Authorization: ApiKey avm_sak_xxxxx_xxxxx
Accept: application/json
```

完整密钥只会在创建和轮换时显示一次。服务端只保存密钥哈希，管理员也不能再次查看原始值。

新接入的可信服务端建议选择公钥请求签名。签名模式只公开 `credentialId`，私钥始终保留在调用方，并为时间、nonce、请求路径、query 和请求体提供完整性保护。

## 请求速率

管理员可为每把密钥单独设置每分钟请求数，默认值为 `120`。达到上限时接口返回 `429` 和 `SITE_API_KEY_RATE_LIMIT_EXCEEDED`；调用方应读取响应头并在窗口重置前停止重试。

| 响应头 | 说明 |
| --- | --- |
| `RateLimit-Limit` | 当前密钥在一分钟内允许的请求数。 |
| `RateLimit-Remaining` | 当前一分钟窗口内剩余的请求数。 |
| `RateLimit-Reset` | 距离窗口重置的秒数。 |
| `Retry-After` | 收到 `429` 时至少应等待的秒数。 |

## 功能授权

| 业务域 | Scope | 可调用接口 |
| --- | --- | --- |
| 概况 | `site:status:read` | `/overview` |
| [资源](/docs/api/site-api-resources/) | `resource:read:public` | 列表、分类、筛选、详情 |
| [帖子](/docs/api/site-api-posts/) | `post:read:public` | 列表、版块、筛选、详情 |
| [找服玩](/docs/api/site-api-players/) | `player:read:public` | 列表、分类、筛选、详情 |
| [视频区](/docs/api/site-api-videos/) | `video:read:public` | 列表、分类、筛选、详情 |
| 全站搜索 | `search:read:public` | 公开资源、帖子、找服玩、视频搜索 |
| 活动 | `activity:read:public` | 公开活动与关联帖子 |
| 排行榜 | `ranking:read:public` | 资源、找服玩、视频排行榜 |

详细接口：

-   [全站搜索](/docs/api/site-api-search/)
-   [活动](/docs/api/site-api-activities/)
-   [排行榜](/docs/api/site-api-rankings/)
-   [OAuth 用户身份](/docs/api/site-api-oauth-user/)
-   [OAuth 用户数据](/docs/api/site-api-oauth-user-data/)
-   [实时 API：WebSocket 与机器人](/docs/api/realtime-websocket/)

未勾选的业务域会返回 `403` 和 `SITE_API_KEY_SCOPE_REQUIRED`。一把只勾选“资源”的密钥可访问资源的四个接口，但不能读取帖子。

## 筛选目录与固定预设

官网和机器人可以先读取筛选目录，再按目录中的 `value` 调用内容接口：

```
GET /api/site/v1/catalog
GET /api/site/v1/catalog/resources
```

目录会按当前 API Key 的业务域权限过滤，资源目录包含平台、资源分类、子分类、资源属性、游戏版本、官方标签和资源收录等筛选项，帖子目录包含帖子分类和内容收录等筛选项。目录中的选项使用稳定的 `id`、`slug`、`label` 和 `value` 字段。

对于需要长期复用的查询组合，可以使用内置固定预设：

```
GET /api/site/v1/presets
GET /api/site/v1/presets/java-mods/items?page=1&pageSize=20&sort=updated
```

预设会固定自己的筛选条件，只允许定义中 `allowedQuery` 列出的分页和排序参数；调用方不能覆盖平台、分类或板块条件。预设结果与普通内容列表使用相同的公开内容规则和响应字段，并额外返回 `preset` 元数据。

列表响应会增量包含 `schemaVersion`、`domain`、`items`、`pagination`、`filters` 和 `generatedAt`，原有 `resources`、`posts`、`players` 或 `videos` 字段继续保留。筛选值无效时返回 `400 SITE_API_INVALID_FILTER`，不会静默退化为无筛选查询。

## 健康检查与增量同步

机器人可以使用健康检查确认站点 API 可用：

GET`/api/site/v1/health`

需要同步公开内容变化时，使用 `updatedSince` 获取四类内容的更新摘要：

GET`/api/site/v1/changes?updatedSince=2026-09-04T00:00:00.000Z&limit=50`

响应只包含公开、已审核内容的稳定 `id`、`slug`、标题和更新时间。客户端应保存本次响应中的最新 `updatedAt`，下一次作为 `updatedSince` 传回；首次同步未提供该参数时，服务端默认返回最近 24 小时的变化。

## 获取站点概况

GET `/api/site/v1/overview`

需要权限 `site:status:read` 认证方式 `X-API-Key`

```
curl "https://www.nexusmc.cn/api/site/v1/overview" \
  -H "X-API-Key: $NEXUSMC_SITE_API_KEY" \
  -H "Accept: application/json"
```

```
{
  "site": {
    "name": "NexusMC",
    "subtitle": "Minecraft 我的世界中文社区"
  },
  "contentCounts": {
    "resources": 135,
    "posts": 93,
    "players": 18,
    "videos": 12
  },
  "generatedAt": "2026-08-16T06:00:00.000Z"
}
```

旧版 `/status` 入口已退役，当前统一使用 `/overview`。

## 安全边界

四类内容接口只返回已审核、公开且无访问限制的数据；资源和找服玩还会排除已归档项目。接口不会暴露草稿、私密内容、付费或密码内容、审核资料、用户敏感信息和后台日志。

CAUTION

站点 API 密钥只能部署在可信服务端。不要将它放入浏览器、移动端安装包、公开仓库、截图、日志或可下载配置文件。泄露后应立即撤销并创建新密钥。

| 状态码 | `code` | 处理方式 |
| --- | --- | --- |
| `400` | 无固定 code | 检查路径参数、筛选值和分页参数。 |
| `401` | `SITE_API_KEY_INVALID` | 检查请求头、密钥、启用状态、到期和撤销状态。 |
| `403` | `SITE_API_KEY_SCOPE_REQUIRED` | 在管理面板为该密钥勾选响应中要求的业务域。 |
| `429` | `SITE_API_KEY_RATE_LIMIT_EXCEEDED` | 等待 `Retry-After` 指定的时间后再重试。 |

不要对 `401` 或 `403` 自动重试。对网络错误或 `5xx` 使用有限次数的指数退避，且不要在日志中记录完整密钥。
