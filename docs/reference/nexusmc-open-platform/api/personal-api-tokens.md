---
title: "用户个人 API：接入总览"
source: https://docs.nexusmc.cn/docs/api/personal-api-tokens
collected: 2026-09-26
---

# 用户个人 API：接入总览

第一次接入个人 API 时先读这里：Base URL、Token 认证、权限、接口导航、响应和限流规则。

## 第一次接入先看什么

如果你刚开始接入，按这个顺序即可：

1.  先通过开发者考试，解锁个人 API Token 管理能力。
2.  在站内设置页创建个人 API Token，并勾选需要的权限。
3.  用本页的认证示例确认 Token 可以访问接口。
4.  如果脚本需要正式投稿，确认当前版本的创作者协议。
5.  进入对应业务文档，直接查看端点、完整参数表和请求示例。

IMPORTANT

Token scope 只表示“这个 Token 可以请求哪些接口”，不能替代账号本身的发布资格。创建帖子、资源、找服玩或视频时，账号仍需完成 [创作者开通与协议确认](/docs/site-info/creator-onboarding/)。

| 你要做什么 | 直接阅读 |
| --- | --- |
| 以站点身份读取公开站点数据 | [站点 API：密钥与公开内容接口](/docs/api/site-api-keys/) |
| 直接使用现成工具接入 API | [可直接使用 API 的项目](/docs/api/ready-to-use-projects/) |
| 上传文件或图片 | [个人 API：文件上传](/docs/api/personal-api-uploads/) |
| 上传资源文件并发布资源 | [个人 API：资源上传与发布流程](/docs/api/personal-api-resource-upload-flow/) |
| 创建、读取或更新资源 | [个人 API：资源](/docs/api/personal-api-resources/) |
| 提交资源文件、版本或设置 Loader | [个人 API：资源文件与版本](/docs/api/personal-api-resource-files/) |
| 按 MC 版本和 Loader 检查资源更新 | [个人 API：资源文件与版本的读取版本章节](/docs/api/personal-api-resource-files/#%E8%AF%BB%E5%8F%96%E7%89%88%E6%9C%AC) |
| 接收实时更新、Webhook 或创建机器人 | [WebSocket 连接](/docs/api/realtime-websocket/)、[Webhook 接口](/docs/api/open-platform-webhooks/)、[OneBot 接入](/docs/api/open-platform-bots/) |
| 创建、读取或更新帖子 | [个人 API：帖子](/docs/api/personal-api-posts/) |
| 创建、读取或更新找服玩投稿 | [个人 API：找服玩](/docs/api/personal-api-players/) |
| 创建、读取或更新视频投稿 | [个人 API：视频](/docs/api/personal-api-videos/) |
| 排查状态码或鉴权失败 | [个人 API：错误与排查](/docs/api/personal-api-errors/) |
| 生成 `content` 富文本对象 | [投稿正文 TipTap 快速接入](/docs/api/tiptap-content-format/) |

NOTE

个人 API 只代表 Token 所属用户本人，不是后台管理接口，也不会继承版主、管理员、协作者或组织身份。

读取公开资源详情、版本、文件元数据和更新摘要时，可以直接使用本页列出的“无需认证”接口，不需要个人 Token 或站点 API Key。查看自己资源的待审或被拒版本时，版本列表可使用带 `resource:read:self` 权限的个人 Token。只有需要读取更广泛的站点内容索引或进行服务端集成时，才使用由管理员单独签发的[站点 API 密钥](/docs/api/site-api-keys/)。不要把个人 Token 当作站点凭据。

## Base URL

```
https://www.nexusmc.cn/api
```

文档中的 `POST /api/posts` 对应完整地址：

```
https://www.nexusmc.cn/api/posts
```

## 认证方式

需要个人身份和 scope 的接口使用 Bearer Token：

```
Authorization: Bearer avm_pat_xxxxx_xxxxx
Accept: application/json
```

发送 JSON 时再增加：

```
Content-Type: application/json
```

快速验证示例：

```
curl "https://www.nexusmc.cn/api/users/me/resources?page=1&pageSize=1" \
  -H "Authorization: Bearer avm_pat_xxxxx_xxxxx" \
  -H "Accept: application/json"
```

```
const response = await fetch('https://www.nexusmc.cn/api/users/me/resources?page=1&pageSize=1', {
  headers: {
    Authorization: `Bearer ${process.env.NEXUSMC_API_TOKEN}`,
    Accept: 'application/json',
  },
})

const result = await response.json()
```

```
import os
import requests

response = requests.get(
    'https://www.nexusmc.cn/api/users/me/resources',
    headers={
        'Authorization': f"Bearer {os.environ['NEXUSMC_API_TOKEN']}",
        'Accept': 'application/json',
    },
    params={'page': 1, 'pageSize': 1},
)
response.raise_for_status()
result = response.json()
```

Token 明文只在创建时显示一次。不要把 Token 放进公开仓库、浏览器前端代码、构建产物或可下载的配置文件。

标记为“无需认证”的公开读取接口不要发送 `Authorization`。这类接口只返回公开、已审核且当前访问者有权查看的数据，不会绕过隐藏、付费或其他内容访问限制。

## 无需认证读取公开资源

自动化程序在发布前应先读取发布元数据，不要通过反复提交错误 payload 来猜分类、标签或版本字段：

GET`/api/resources/publishing-metadata?platform=java`

资源详情和版本读取也无需认证：

```
GET /api/resources/{id}
GET /api/resources/{id}/versions
GET /api/resources/{id}/versions?platform=java&mcVersion=1.21.1&loader=fabric&status=approved
GET /api/resources/{id}/versions/{versionId}
GET /api/resources/{id}/files
```

这里的 `{id}` 支持内部 UUID、公开 slug、公开数字 ID 和 `slug.publicId` 组合标识。版本列表支持按平台、MC 版本、Loader、版本 Tag 和审核状态筛选，参数与同一文件匹配规则见[读取版本](/docs/api/personal-api-resource-files/#%E8%AF%BB%E5%8F%96%E7%89%88%E6%9C%AC)。公开列表按审核通过时间倒序返回，第一个是最近审核通过的版本；它不是版本号字符串的语义化排序。带个人 Token 查询本人待审版本时需要 `resource:read:self`，列表按提交时间倒序。

## 获取自己的内容 ID

更新已有内容前，不需要从浏览器地址栏里猜 ID。四类“我的内容”接口都会返回每项的稳定 `id`、当前 `slug` 和站内 `path`：

| 内容 | 查询接口 | 列表字段 | 所需权限 |
| --- | --- | --- | --- |
| 资源 | `GET /api/users/me/resources` | `resources` | `resource:read:self` |
| 帖子 | `GET /api/users/me/posts` | `posts` | `post:read:self` |
| 找服玩 | `GET /api/users/me/players` | `players` | `player:read:self` |
| 视频 | `GET /api/users/me/videos` | `videos` | `video:read:self` |

例如按标题查询自己的资源：

```
curl "https://www.nexusmc.cn/api/users/me/resources?q=MittelLib&pageSize=10" \
  -H "Authorization: Bearer avm_pat_xxxxx_xxxxx" \
  -H "Accept: application/json"
```

响应中的 `resources[].id` 可直接保存为 GitHub Actions 的 `NEXUSMC_RESOURCE_ID`：

```
{
  "resources": [
    {
      "id": "8a86e434-5050-4f4a-a861-5a5fe28873a4",
      "slug": "mittelib",
      "path": "/resources/mittelib",
      "title": "MittelLib"
    }
  ],
  "pagination": { "page": 1, "pageSize": 10, "total": 1, "totalPages": 1 }
}
```

`id` 是自动化配置的首选值，创建接口的成功响应也会直接返回它。内容接口同时兼容当前或历史 `slug`、公开编号以及完整公开链接中的路径段，但 slug 可能随标题或自定义路径修改，不适合长期写死。

## 权限表

| 权限值 | 允许调用的能力 |
| --- | --- |
| `upload:file` | 上传单文件、多个文件、图片、直传和分块文件。 |
| `resource:read:self` | 读取自己的资源列表和编辑态数据。 |
| `resource:create` | 创建自己的资源。 |
| `resource:update:self` | 更新自己的资源和提交资源版本。 |
| `resource:delete:self` | 删除自己的个人身份资源。 |
| `resource:read:feed` | 读取全站资源更新摘要。 |
| `resource:read:activity` | 读取自己资源的互动摘要。 |
| `post:read:self` | 读取自己的帖子列表和编辑态数据。 |
| `post:create` | 创建自己的帖子。 |
| `post:update:self` | 更新自己的帖子。 |
| `post:read:feed` | 读取全站帖子更新摘要。 |
| `post:read:activity` | 读取自己帖子的互动摘要。 |
| `player:read:self` | 读取自己的找服玩列表和编辑态数据。 |
| `player:create` | 创建自己的找服玩投稿。 |
| `player:update:self` | 更新自己的找服玩投稿。 |
| `player:read:feed` | 读取全站找服玩更新摘要。 |
| `player:read:activity` | 读取自己找服玩投稿的互动摘要。 |
| `video:read:self` | 读取自己的视频列表和编辑态数据。 |
| `video:create` | 创建自己的视频投稿。 |
| `video:update:self` | 更新自己的视频投稿。 |
| `video:read:feed` | 读取全站视频更新摘要。 |
| `video:read:activity` | 读取自己视频投稿的互动摘要。 |

建议一个自动化程序使用一个 Token，并只授予实际需要的最小权限。

NOTE

四个 `*:read:feed` 权限由 `GET /api/users/me/latest-updates` 使用，四个 `*:read:activity` 权限由 `GET /api/users/me/activity` 使用。这两个接口要求同一组的四个权限全部授予，只勾其中一个不会生效。

## 接口导航

### 公开读取与个人内容接口

| 方法 | 路径 | 权限 | 说明 |
| --- | --- | --- | --- |
| GET | `/api/resources/publishing-metadata?platform=java` |
`无需认证`

 | 读取资源发布契约、字段别名和发现接口 |
| GET | `/api/resources/platforms` |

`无需认证`

 | 读取可用资源平台 |
| GET | `/api/resources/categories?platform=java` |

`无需认证`

 | 读取分类树；提交 category 和 subCategory 时使用 filterValue |
| GET | `/api/resources/version-groups` |

`无需认证`

 | 读取版本分组元数据 |
| GET | `/api/resources/game-versions?platform=java` |

`无需认证`

 | 读取可提交到 mcVersions 的游戏版本号 |
| GET | `/api/resources/version-tags` |

`无需认证`

 | 读取 versionTag 可选值 |
| GET | `/api/resources/official-tag-groups` |

`无需认证`

 | 读取官方标签组、标签 ID 和 slug |
| GET | `/api/resources/updates` |

`无需认证`

 | 读取公开资源更新摘要，用于更新检查 |
| GET | `/api/resources/{id}` |

`无需认证`

 | 读取公开、已审核且可见的资源信息 |
| GET | `/api/resources/{id}/versions` |

`无需认证``resource:read:self`

 | 按平台、MC 版本、Loader、Tag、状态筛选；个人 Token 可查看本人待审版本 |
| GET | `/api/resources/{id}/versions/{versionId}` |

`无需认证``resource:read:self`

 | 读取指定的已审核资源版本；个人 Token 不能读取待审版本的文件 |
| GET | `/api/resources/{id}/files` |

`无需认证``resource:read:self`

 | 读取可见资源的已审核版本文件元数据 |
| POST | `/api/upload` |

`upload:file`

 | 上传单个文件 |
| POST | `/api/upload/multiple` |

`upload:file`

 | 上传多个文件 |
| POST | `/api/upload/image` |

`upload:file`

 | 上传图片 |
| POST | `/api/upload/direct/init` |

`upload:file`

 | 初始化直传 |
| POST | `/api/upload/direct/complete` |

`upload:file`

 | 确认直传 |
| POST | `/api/upload/session/init` |

`upload:file`

 | 初始化分块上传 |
| GET | `/api/upload/session/{uploadId}/status` |

`upload:file`

 | 查询分块合并状态 |
| GET | `/api/upload/download` |

`无需认证`

 | 按指定文件名下载站点托管文件 |
| GET | `/api/upload/limits` |

`upload:file`

 | 读取当前账号的实际上传限额与允许类型 |
| GET | `/api/users/me/latest-updates` |

`resource:read:feed``post:read:feed``player:read:feed``video:read:feed`

 | 读取全站最近更新的公开内容与最近互动 |
| GET | `/api/users/me/activity` |

`resource:read:activity``post:read:activity``player:read:activity``video:read:activity`

 | 读取自己内容的互动摘要 |
| GET | `/api/users/{id}/featured` |

`resource:read:self`

 | 读取指定用户的精选内容与精选位上限 |
| GET | `/api/resources/{id}/edit/core` |

`resource:read:self`

 | 读取资源核心编辑数据 |
| POST | `/api/resources` |

`resource:create`

 | 创建资源（写操作，不是读取资源列表） |
| PATCH | `/api/resources/{id}` |

`resource:update:self`

 | 更新自己的资源 |
| POST | `/api/resources/{id}/versions` |

`resource:update:self`

 | 提交待审核的新版本（读取版本请使用同路径的 GET） |
| DELETE | `/api/resources/{id}` |

`resource:delete:self`

 | 删除自己的资源 |
| GET | `/api/posts/{id}/edit/core` |

`post:read:self`

 | 读取帖子核心编辑数据 |
| POST | `/api/posts` |

`post:create`

 | 创建帖子 |
| PUT | `/api/posts/{id}` |

`post:update:self`

 | 更新帖子 |
| GET | `/api/players/{id}/edit/core` |

`player:read:self`

 | 读取找服玩核心编辑数据 |
| POST | `/api/players` |

`player:create`

 | 创建找服玩投稿 |
| PUT | `/api/players/{id}` |

`player:update:self`

 | 更新找服玩投稿 |
| GET | `/api/videos/{id}/edit` |

`video:read:self`

 | 读取视频编辑态 |
| POST | `/api/videos` |

`video:create`

 | 创建视频投稿 |
| PUT | `/api/videos/{id}` |

`video:update:self`

 | 更新视频投稿 |

## 读取全站更新与互动摘要

GET `/api/users/me/latest-updates`

需要权限 `resource:read:feed + post:read:feed + player:read:feed + video:read:feed`

读取全站最近更新的公开内容，按资源、帖子、找服玩、视频分组返回，另外附带最近互动。

| 查询参数 | 说明 |
| --- | --- |
| `resourceLimit`、`postLimit`、`playerLimit`、`videoLimit` | 每类返回条数，默认 8，最多 20。 |
| `interactionLimit` | 最近互动条数，默认 3，最多 8。 |

响应包含 `resources`、`posts`、`players`、`videos` 和 `interactions`。

GET `/api/users/me/activity`

需要权限 `resource:read:activity + post:read:activity + player:read:activity + video:read:activity`

读取自己四类内容的互动摘要，用来判断哪些内容有新评论或回复。

| 查询参数 | 说明 |
| --- | --- |
| `limit` | 每类返回条数，默认 10，最多 30。 |

响应同样按 `resources`、`posts`、`players`、`videos` 分组，每项包含 `id` 和 `interactionSummary`（`commentCount`、`replyCount`、`latestAt`、`comments`、`replies`）。

## 读取用户精选内容

GET `/api/users/{id}/featured`

需要权限 `resource:read:self`

读取指定用户的精选内容与精选位上限，响应包含 `featured`、`limit` 和 `canEdit`。`{id}` 填自己的用户 ID 时，`canEdit` 返回 `true`。

## 正文字段边界

第一次接入时最容易混淆的是 `content` 和 `description`。下面这张表直接按业务列出：

| 业务 | 富文本字段 | 简介字段 | 说明 |
| --- | --- | --- | --- |
| 资源 | `content` | `description` | `content` 支持 TipTap JSON；`description` 是短简介。 |
| 帖子 | `content` | 无 | 帖子正文走 `content`，推荐传 TipTap JSON。 |
| 找服玩 | `content` | `description` | `content` 是详细介绍；`description` 是短介绍。 |
| 视频 | 无 | `description` | 视频当前没有 `content`；简介是纯文本。 |

TipTap JSON 的最小结构、节点和示例见[投稿正文 TipTap 快速接入](/docs/api/tiptap-content-format/)。

## 通用响应规则

成功响应通常直接返回对象或列表，不会统一包在 `data` 字段中。

```
{
  "id": "content_id",
  "title": "示例标题",
  "status": "draft",
  "updatedAt": "2026-06-24T00:00:00.000Z"
}
```

失败响应至少包含 `error`，部分错误还包含 `code`、`requiredScope` 或 `retryAfter`：

```
{
  "error": "当前 API Token 没有此操作权限",
  "code": "API_TOKEN_SCOPE_REQUIRED",
  "requiredScope": "post:create"
}
```

## 限流

响应可能包含：

| 响应头 | 说明 |
| --- | --- |
| `RateLimit-Limit` | 当前窗口允许的请求数。 |
| `RateLimit-Remaining` | 当前窗口剩余请求数。 |
| `RateLimit-Reset` | 距离窗口重置的秒数。 |
| `Retry-After` | 返回 `429` 后至少应等待的秒数。 |

客户端不应把固定额度写死。遇到 `429` 时读取 `Retry-After`，并使用指数退避。

## 通用限制

-   个人 Token 只能读取或修改 Token 所属用户自己的内容。
-   非空 `organizationId` 会被拒绝；个人 Token 暂不支持组织身份。
-   删除仅限 `resource:delete:self` 删除自己的个人身份资源；审核、推荐、封禁和积分调整不属于个人 API。
-   正式投稿可能要求账号先完成手机绑定、强验证或其他站点安全条件。
