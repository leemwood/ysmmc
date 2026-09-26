---
title: "个人 API：找服玩"
source: https://docs.nexusmc.cn/docs/api/personal-api-players
collected: 2026-09-26
---

# 个人 API：找服玩

查询、创建和更新找服玩投稿，以及连接地址、可见性和发布校验规则。

## 接口一览

| 方法 | 路径 | 权限 | 用途 |
| --- | --- | --- | --- |
| GET | `/api/players/{id}/edit/core` | `player:read:self` | 读取自己的投稿核心编辑数据 |
| GET | `/api/players/{id}/edit/tutorials` | `player:read:self` | 读取投稿关联的教程帖 |
| GET | `/api/users/me/players` | `player:read:self` | 分页查询自己的投稿 |
| POST | `/api/players` | `player:create` | 创建投稿或保存草稿 |
| PUT | `/api/players/{id}` | `player:update:self` | 更新自己的投稿 |

`{id}` 推荐使用列表或创建响应返回的稳定投稿 ID，也兼容当前或历史 slug、公开编号和公开链接中的路径段。正文 `content` 推荐使用 [投稿正文 TipTap 快速接入](/docs/api/tiptap-content-format/)。

## 查询自己的投稿

GET `/api/users/me/players`

需要权限 `player:read:self`

```
curl "https://www.nexusmc.cn/api/users/me/players?page=1&pageSize=12&status=draft" \
  -H "Authorization: Bearer avm_pat_xxxxx_xxxxx" \
  -H "Accept: application/json"
```

```
const url = new URL('https://www.nexusmc.cn/api/users/me/players')
url.search = new URLSearchParams({
  page: '1',
  pageSize: '12',
  status: 'draft',
}).toString()

const response = await fetch(url, {
  headers: {
    Authorization: `Bearer ${process.env.NEXUSMC_API_TOKEN}`,
    Accept: 'application/json',
  },
})

const players = await response.json()
```

```
import os
import requests

response = requests.get(
    'https://www.nexusmc.cn/api/users/me/players',
    headers={
        'Authorization': f"Bearer {os.environ['NEXUSMC_API_TOKEN']}",
        'Accept': 'application/json',
    },
    params={'page': 1, 'pageSize': 12, 'status': 'draft'},
)
response.raise_for_status()
players = response.json()
```

| 查询参数 | 类型 | 必填 | 默认值 | 说明 |
| --- | --- | --- | --- | --- |
| `page` | `number` | 否 | `1` | 页码。 |
| `pageSize` | `number` | 否 | `12` | 每页数量，1～50。 |
| `status` | `string` | 否 | — | 按状态筛选。 |
| `q` | `string` | 否 | 空 | 按标题或简介搜索。 |
| `sort` | `string` | 否 | `updated` | `updated` 或 `created`。 |

响应的 `players[]` 每项都包含 `id`、`slug` 和 `path`。更新投稿时直接使用 `players[].id`，无需查看网页 URL。

读取单项编辑数据：

GET `/api/players/{id}/edit/core`

需要权限 `player:read:self`

```
curl "https://www.nexusmc.cn/api/players/PLAYER_ID/edit/core" \
  -H "Authorization: Bearer avm_pat_xxxxx_xxxxx" \
  -H "Accept: application/json"
```

```
const response = await fetch('https://www.nexusmc.cn/api/players/PLAYER_ID/edit/core', {
  headers: {
    Authorization: `Bearer ${process.env.NEXUSMC_API_TOKEN}`,
    Accept: 'application/json',
  },
})

const player = await response.json()
```

```
import os
import requests

response = requests.get(
    'https://www.nexusmc.cn/api/players/PLAYER_ID/edit/core',
    headers={
        'Authorization': f"Bearer {os.environ['NEXUSMC_API_TOKEN']}",
        'Accept': 'application/json',
    },
)
response.raise_for_status()
player = response.json()
```

### 其他编辑数据接口

| 接口 | 权限 | 响应 |
| --- | --- | --- |
| `GET /api/players/{id}/edit/tutorials` | `player:read:self` | `tutorialPosts` |

聚合接口 `GET /api/players/{id}/edit` 已停止服务，会返回 `410` 和 `PLAYER_EDIT_AGGREGATE_RETIRED`，请改用 `core` 与 `tutorials` 两个接口。

## 创建找服玩投稿

POST `/api/players`

需要权限 `player:create` 请求类型 `application/json`

### 完整请求参数

| 参数 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `title` | `string` | 是 | 服务器名称。 |
| `description` | `string` | 是 | 简短介绍。 |
| `category` | `string` | 是 | 服务器分类。 |
| `ip` | `string` | 是 | 可解析的公网服务器地址；即使隐藏展示也必须提供。 |
| `onlineMode` | `boolean` | 是 | 是否开启正版验证，必须传布尔值。 |
| `networkEnvironments` | `array<string>` | 是 | 网络环境，至少 1 项。 |
| `content` | `object | string` | 否 | 详细介绍；正式发布时须达到站点当前最小字数。 |
| `platform` | `string` | 否 | 平台，默认 `java`。 |
| `versions` | `array<string>` | 否 | 支持的游戏版本。 |
| `tags` | `array<string>` | 否 | 标签。 |
| `tutorialPostIds` | `array<string>` | 否 | 关联教程帖 ID。 |
| `port` | `number` | 否 | 端口，1～65535；默认 `25565`。 |
| `hideConnectionAddress` | `boolean` | 否 | 是否对访客隐藏连接地址。 |
| `connectionApplyUrl` | `string | null` | 否 | HTTP/HTTPS 连接申请入口。 |
| `website` | `string | null` | 否 | 官方网站。 |
| `qqGroup` | `string | null` | 否 | QQ 群信息。 |
| `discord` | `string | null` | 否 | Discord 信息。 |
| `coverImage` | `string | null` | 否 | 封面图片 URL。 |
| `galleryMode` | `string` | 否 | `auto` 或 `manual`；默认 `auto`。 |
| `galleryImages` | `array<string>` | 否 | `manual` 模式的图库图片 URL。 |
| `onlineCount` | `number` | 否 | 当前在线人数，最小为 0。 |
| `maxPlayers` | `number` | 否 | 最大人数，最小为 0。 |
| `isDraft` | `boolean` | 否 | `true` 保存草稿；否则提交审核。 |
| `visibility` | `string` | 否 | `public`、`partial` 或 `private`。 |
| `visibilityAllowUserIds` | `array<string>` | 否 | `partial` 可见用户 ID，最多 100 个。 |
| `accessMode` | `string` | 否 | `none`、`password`、`reply` 或 `login`。 |
| `accessPassword` | `string` | 否 | `password` 模式的访问密码，4～64 字符。 |
| `customSlug` | `string | null` | 否 | 自定义路径。 |
| `organizationId` | `string | null` | 否 | 个人 Token 不允许指定组织；请省略或传 `null`。 |

### 创建示例

```
curl -X POST 'https://www.nexusmc.cn/api/players' \
  -H 'Authorization: Bearer avm_pat_xxxxx_xxxxx' \
  -H 'Content-Type: application/json' \
  -d '{
    "title": "示例服务器",
    "description": "一个生存服务器",
    "category": "survival",
    "ip": "play.example.com",
    "port": 25565,
    "onlineMode": true,
    "networkEnvironments": ["电信"],
    "isDraft": true
  }'
```

```
const response = await fetch('https://www.nexusmc.cn/api/players', {
  method: 'POST',
  headers: {
    Authorization: `Bearer ${process.env.NEXUSMC_API_TOKEN}`,
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    title: '示例服务器',
    description: '一个生存服务器',
    category: 'survival',
    ip: 'play.example.com',
    port: 25565,
    onlineMode: true,
    networkEnvironments: ['电信'],
    isDraft: true,
  }),
})

const created = await response.json()
```

```
import os
import requests

response = requests.post(
    'https://www.nexusmc.cn/api/players',
    headers={
        'Authorization': f"Bearer {os.environ['NEXUSMC_API_TOKEN']}",
        'Content-Type': 'application/json',
    },
    json={
        'title': '示例服务器',
        'description': '一个生存服务器',
        'category': 'survival',
        'ip': 'play.example.com',
        'port': 25565,
        'onlineMode': True,
        'networkEnvironments': ['电信'],
        'isDraft': True,
    },
)
response.raise_for_status()
created = response.json()
```

## 更新找服玩投稿

PUT `/api/players/{id}`

需要权限 `player:update:self` 请求类型 `application/json`

所有字段均为可选，只修改请求里出现的字段；但更新后的完整数据仍必须满足标题、简介、分类、地址、正版验证和网络环境要求。

NOTE

更新已发布（`approved`）或已拒（`rejected`）的投稿会重新进入审核，状态回到 `pending`，审核期间前台仍展示旧内容。

### 完整请求参数

| 参数 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `title` | `string` | 否 | 服务器名称。 |
| `description` | `string` | 否 | 简短介绍。 |
| `content` | `object | string` | 否 | 详细介绍。 |
| `platform` | `string` | 否 | 平台，如 `java`。 |
| `category` | `string` | 否 | 服务器分类。 |
| `versions` | `array<string>` | 否 | 用本数组整体替换游戏版本。 |
| `networkEnvironments` | `array<string>` | 否 | 用本数组整体替换网络环境，不得为空。 |
| `onlineMode` | `boolean` | 否 | 是否开启正版验证。 |
| `tags` | `array<string>` | 否 | 用本数组整体替换标签。 |
| `tutorialPostIds` | `array<string>` | 否 | 用本数组整体替换关联教程。 |
| `ip` | `string` | 否 | 可解析的公网服务器地址。 |
| `port` | `number` | 否 | 端口，1～65535。 |
| `hideConnectionAddress` | `boolean` | 否 | 是否隐藏连接地址。 |
| `connectionApplyUrl` | `string | null` | 否 | HTTP/HTTPS 连接申请入口。 |
| `website` | `string | null` | 否 | 官方网站。 |
| `qqGroup` | `string | null` | 否 | QQ 群信息。 |
| `discord` | `string | null` | 否 | Discord 信息。 |
| `coverImage` | `string | null` | 否 | 封面图片 URL。 |
| `galleryMode` | `string` | 否 | `auto` 或 `manual`。 |
| `galleryImages` | `array<string>` | 否 | `manual` 模式的新图库；与 `galleryMode` 一起发送。 |
| `onlineCount` | `number` | 否 | 当前在线人数，最小为 0。 |
| `maxPlayers` | `number` | 否 | 最大人数，最小为 0。 |
| `isDraft` | `boolean` | 否 | 草稿传 `true`；将草稿提交审核传 `false`。 |
| `visibility` | `string` | 否 | `public`、`partial` 或 `private`。 |
| `visibilityAllowUserIds` | `array<string>` | 否 | `partial` 可见用户 ID。 |
| `accessMode` | `string` | 否 | `none`、`password`、`reply` 或 `login`。 |
| `accessPassword` | `string` | 否 | 新访问密码；未传时可保留已有密码哈希。 |
| `customSlug` | `string | null` | 否 | 自定义路径。 |
| `organizationId` | `string | null` | 否 | 个人 Token 不允许指定组织；请省略或传 `null`。 |
| `archived` | `boolean` | 否 | 是否归档。 |
| `archivedAt` | `string | null` | 否 | 归档时间，建议 ISO 8601。 |
| `archiveNote` | `string | null` | 否 | 归档说明。 |

### 更新示例

```
curl -X PUT 'https://www.nexusmc.cn/api/players/PLAYER_ID' \
  -H 'Authorization: Bearer avm_pat_xxxxx_xxxxx' \
  -H 'Content-Type: application/json' \
  -d '{
    "title": "更新后的服务器名",
    "tags": ["生存", "公益"]
  }'
```

```
const response = await fetch('https://www.nexusmc.cn/api/players/PLAYER_ID', {
  method: 'PUT',
  headers: {
    Authorization: `Bearer ${process.env.NEXUSMC_API_TOKEN}`,
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    title: '更新后的服务器名',
    tags: ['生存', '公益'],
  }),
})

const updated = await response.json()
```

```
import os
import requests

response = requests.put(
    'https://www.nexusmc.cn/api/players/PLAYER_ID',
    headers={
        'Authorization': f"Bearer {os.environ['NEXUSMC_API_TOKEN']}",
        'Content-Type': 'application/json',
    },
    json={
        'title': '更新后的服务器名',
        'tags': ['生存', '公益'],
    },
)
response.raise_for_status()
updated = response.json()
```

## 组合约束

-   `visibility: "partial"` 时必须提供非空的 `visibilityAllowUserIds`。
-   `accessMode: "password"` 时，首次设置必须提供 4～64 字符的 `accessPassword`。
-   `hideConnectionAddress: true` 时，`connectionApplyUrl`、`qqGroup`、`discord` 至少填写一个。
-   `connectionApplyUrl` 只能是 HTTP 或 HTTPS 地址。
-   正式发布还会检查正文最小字数、账号安全绑定和发布频率；草稿不执行完整发布校验。
