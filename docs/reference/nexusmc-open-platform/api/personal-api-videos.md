---
title: "个人 API：视频"
source: https://docs.nexusmc.cn/docs/api/personal-api-videos
collected: 2026-09-26
---

# 个人 API：视频

查询、创建和更新视频投稿，包含来源链接、分区、标签和访问控制参数。

## 接口一览

| 方法 | 路径 | 权限 | 用途 |
| --- | --- | --- | --- |
| GET | `/api/videos/{id}/edit` | `video:read:self` | 读取自己的视频编辑数据 |
| GET | `/api/users/me/videos` | `video:read:self` | 分页查询自己的视频 |
| POST | `/api/videos` | `video:create` | 创建视频投稿 |
| PUT | `/api/videos/{id}` | `video:update:self` | 更新自己的视频投稿 |

`{id}` 推荐使用列表或创建响应返回的稳定视频 ID，也兼容当前或历史 slug、公开编号和公开链接中的路径段。

## 查询自己的视频

GET `/api/users/me/videos`

需要权限 `video:read:self`

```
curl "https://www.nexusmc.cn/api/users/me/videos?page=1&pageSize=12&status=pending" \
  -H "Authorization: Bearer avm_pat_xxxxx_xxxxx" \
  -H "Accept: application/json"
```

```
const url = new URL('https://www.nexusmc.cn/api/users/me/videos')
url.search = new URLSearchParams({
  page: '1',
  pageSize: '12',
  status: 'pending',
}).toString()

const response = await fetch(url, {
  headers: {
    Authorization: `Bearer ${process.env.NEXUSMC_API_TOKEN}`,
    Accept: 'application/json',
  },
})

const videos = await response.json()
```

```
import os
import requests

response = requests.get(
    'https://www.nexusmc.cn/api/users/me/videos',
    headers={
        'Authorization': f"Bearer {os.environ['NEXUSMC_API_TOKEN']}",
        'Accept': 'application/json',
    },
    params={'page': 1, 'pageSize': 12, 'status': 'pending'},
)
response.raise_for_status()
videos = response.json()
```

| 查询参数 | 类型 | 必填 | 默认值 | 说明 |
| --- | --- | --- | --- | --- |
| `page` | `number` | 否 | `1` | 页码。 |
| `pageSize` | `number` | 否 | `12` | 每页数量，1～50。 |
| `status` | `string` | 否 | — | 按状态筛选。 |
| `q` | `string` | 否 | 空 | 按标题或简介搜索。 |
| `sort` | `string` | 否 | `updated` | `updated` 或 `created`。 |

响应的 `videos[]` 每项都包含 `id`、`slug` 和 `path`。更新视频时直接使用 `videos[].id`，无需查看网页 URL。

读取单项编辑数据：

GET `/api/videos/{id}/edit`

需要权限 `video:read:self`

```
curl "https://www.nexusmc.cn/api/videos/VIDEO_ID/edit" \
  -H "Authorization: Bearer avm_pat_xxxxx_xxxxx" \
  -H "Accept: application/json"
```

```
const response = await fetch('https://www.nexusmc.cn/api/videos/VIDEO_ID/edit', {
  headers: {
    Authorization: `Bearer ${process.env.NEXUSMC_API_TOKEN}`,
    Accept: 'application/json',
  },
})

const video = await response.json()
```

```
import os
import requests

response = requests.get(
    'https://www.nexusmc.cn/api/videos/VIDEO_ID/edit',
    headers={
        'Authorization': f"Bearer {os.environ['NEXUSMC_API_TOKEN']}",
        'Accept': 'application/json',
    },
)
response.raise_for_status()
video = response.json()
```

## 创建视频投稿

POST `/api/videos`

需要权限 `video:create` 请求类型 `application/json`

### 完整请求参数

| 参数 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `sourceUrl` | `string` | 是 | Bilibili、腾讯视频或抖音的视频链接。 |
| `categoryId` | `string` | 否 | 视频分区 ID；省略时使用排序最靠前的默认分区。 |
| `title` | `string` | 否 | 标题；省略时尝试读取来源元数据，再使用来源兜底标题。 |
| `description` | `string` | 否 | 视频简介。 |
| `customSlug` | `string` | 否 | 自定义路径。 |
| `organizationId` | `string | null` | 否 | 个人 Token 不允许指定组织；请省略或传 `null`。 |
| `coverImage` | `string | null` | 否 | 手动封面 URL；省略时尝试读取来源封面。 |
| `tags` | `array<string>` | 否 | 标签，去重后最多保留 10 个。 |
| `visibility` | `string` | 否 | `public`、`partial` 或 `private`。 |
| `visibilityAllowUserIds` | `array<string>` | 否 | `partial` 可见用户 ID，最多 100 个。 |
| `accessMode` | `string` | 否 | `none`、`password`、`reply` 或 `login`。 |
| `accessPassword` | `string` | 否 | `password` 模式的访问密码，4～64 字符。 |

### 创建示例

```
curl -X POST 'https://www.nexusmc.cn/api/videos' \
  -H 'Authorization: Bearer avm_pat_xxxxx_xxxxx' \
  -H 'Content-Type: application/json' \
  -d '{
    "sourceUrl": "https://www.bilibili.com/video/BVxxxxxxxxxx",
    "title": "示例视频",
    "tags": ["教程", "红石"]
  }'
```

```
const response = await fetch('https://www.nexusmc.cn/api/videos', {
  method: 'POST',
  headers: {
    Authorization: `Bearer ${process.env.NEXUSMC_API_TOKEN}`,
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    sourceUrl: 'https://www.bilibili.com/video/BVxxxxxxxxxx',
    title: '示例视频',
    tags: ['教程', '红石'],
  }),
})

const created = await response.json()
```

```
import os
import requests

response = requests.post(
    'https://www.nexusmc.cn/api/videos',
    headers={
        'Authorization': f"Bearer {os.environ['NEXUSMC_API_TOKEN']}",
        'Content-Type': 'application/json',
    },
    json={
        'sourceUrl': 'https://www.bilibili.com/video/BVxxxxxxxxxx',
        'title': '示例视频',
        'tags': ['教程', '红石'],
    },
)
response.raise_for_status()
created = response.json()
```

创建成功返回 `201`，投稿进入 `pending` 审核状态。正式投稿还会检查账号限制、手机绑定和发布频率。

## 更新视频投稿

PUT `/api/videos/{id}`

需要权限 `video:update:self` 请求类型 `application/json`

个人 Token 只能更新本人投稿。更新成功后会重新进入 `pending` 审核状态。

### 完整请求参数

| 参数 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `sourceUrl` | `string` | 否 | 新的视频来源链接；非空时必须来自 Bilibili、腾讯视频或抖音。 |
| `categoryId` | `string` | 否 | 新的视频分区 ID。 |
| `title` | `string` | 否 | 新标题。 |
| `description` | `string` | 否 | 新简介。 |
| `customSlug` | `string` | 否 | 新的自定义路径。 |
| `organizationId` | `string | null` | 否 | 个人 Token 不允许指定组织；请省略。 |
| `coverImage` | `string | null` | 否 | 新封面 URL；传空字符串或 `null` 可清空手动封面。 |
| `tags` | `array<string>` | 否 | 用本数组整体替换标签，去重后最多 10 个。 |
| `visibility` | `string` | 否 | `public`、`partial` 或 `private`。 |
| `visibilityAllowUserIds` | `array<string>` | 否 | `partial` 可见用户 ID。 |
| `accessMode` | `string` | 否 | `none`、`password`、`reply` 或 `login`。 |
| `accessPassword` | `string` | 否 | 新访问密码；未传时可保留已有密码哈希。 |

### 更新示例

```
curl -X PUT 'https://www.nexusmc.cn/api/videos/VIDEO_ID' \
  -H 'Authorization: Bearer avm_pat_xxxxx_xxxxx' \
  -H 'Content-Type: application/json' \
  -d '{
    "title": "更新后的标题",
    "tags": ["教程", "生存"]
  }'
```

```
const response = await fetch('https://www.nexusmc.cn/api/videos/VIDEO_ID', {
  method: 'PUT',
  headers: {
    Authorization: `Bearer ${process.env.NEXUSMC_API_TOKEN}`,
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    title: '更新后的标题',
    tags: ['教程', '生存'],
  }),
})

const updated = await response.json()
```

```
import os
import requests

response = requests.put(
    'https://www.nexusmc.cn/api/videos/VIDEO_ID',
    headers={
        'Authorization': f"Bearer {os.environ['NEXUSMC_API_TOKEN']}",
        'Content-Type': 'application/json',
    },
    json={
        'title': '更新后的标题',
        'tags': ['教程', '生存'],
    },
)
response.raise_for_status()
updated = response.json()
```

## 组合约束

-   `visibility: "partial"` 时必须提供至少一个非本人用户 ID。
-   `accessMode: "password"` 时，首次设置必须提供 4～64 字符的 `accessPassword`。
-   更新来源链接时，服务端会重新解析来源 ID、嵌入地址和元数据；无法识别的链接返回 `400`。
