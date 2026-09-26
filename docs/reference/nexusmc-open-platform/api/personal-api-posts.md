---
title: "个人 API：帖子"
source: https://docs.nexusmc.cn/docs/api/personal-api-posts
collected: 2026-09-26
---

# 个人 API：帖子

查询、创建和更新帖子，以及普通帖、问答帖和投票帖的完整参数。

## 接口一览

| 方法 | 路径 | 权限 | 用途 |
| --- | --- | --- | --- |
| GET | `/api/posts/{id}/edit/core` | `post:read:self` | 读取自己的帖子核心编辑数据 |
| GET | `/api/posts/{id}/edit/event` | `post:read:self` | 读取帖子上的活动信息 |
| GET | `/api/posts/{id}/edit/permissions` | `post:read:self` | 读取当前用户对帖子的编辑权限 |
| GET | `/api/users/me/posts` | `post:read:self` | 分页查询自己的帖子 |
| POST | `/api/posts` | `post:create` | 创建帖子或保存草稿 |
| PUT | `/api/posts/{id}` | `post:update:self` | 更新自己的帖子 |

`{id}` 推荐使用列表或创建响应返回的稳定帖子 ID，也兼容当前或历史 slug、公开编号和公开链接中的路径段。帖子正文 `content` 推荐使用 [投稿正文 TipTap 快速接入](/docs/api/tiptap-content-format/)。

## 查询自己的帖子

GET `/api/users/me/posts`

需要权限 `post:read:self`

```
curl "https://www.nexusmc.cn/api/users/me/posts?page=1&pageSize=20&status=draft" \
  -H "Authorization: Bearer avm_pat_xxxxx_xxxxx" \
  -H "Accept: application/json"
```

```
const url = new URL('https://www.nexusmc.cn/api/users/me/posts')
url.search = new URLSearchParams({
  page: '1',
  pageSize: '20',
  status: 'draft',
}).toString()

const response = await fetch(url, {
  headers: {
    Authorization: `Bearer ${process.env.NEXUSMC_API_TOKEN}`,
    Accept: 'application/json',
  },
})

const posts = await response.json()
```

```
import os
import requests

response = requests.get(
    'https://www.nexusmc.cn/api/users/me/posts',
    headers={
        'Authorization': f"Bearer {os.environ['NEXUSMC_API_TOKEN']}",
        'Accept': 'application/json',
    },
    params={'page': 1, 'pageSize': 20, 'status': 'draft'},
)
response.raise_for_status()
posts = response.json()
```

| 查询参数 | 类型 | 必填 | 默认值 | 说明 |
| --- | --- | --- | --- | --- |
| `page` | `number` | 否 | `1` | 页码，从 1 开始。 |
| `pageSize` | `number` | 否 | `20` | 每页数量，1～50。 |
| `status` | `string` | 否 | — | 按状态筛选。 |
| `q` | `string` | 否 | 空 | 按标题搜索。 |
| `sort` | `string` | 否 | `updated` | `updated` 或 `created`。 |

响应的 `posts[]` 每项都包含 `id`、`slug` 和 `path`。更新帖子时直接使用 `posts[].id`，无需查看网页 URL。

读取单项编辑数据：

GET `/api/posts/{id}/edit/core`

需要权限 `post:read:self`

```
curl "https://www.nexusmc.cn/api/posts/POST_ID/edit/core" \
  -H "Authorization: Bearer avm_pat_xxxxx_xxxxx" \
  -H "Accept: application/json"
```

```
const response = await fetch('https://www.nexusmc.cn/api/posts/POST_ID/edit/core', {
  headers: {
    Authorization: `Bearer ${process.env.NEXUSMC_API_TOKEN}`,
    Accept: 'application/json',
  },
})

const post = await response.json()
```

```
import os
import requests

response = requests.get(
    'https://www.nexusmc.cn/api/posts/POST_ID/edit/core',
    headers={
        'Authorization': f"Bearer {os.environ['NEXUSMC_API_TOKEN']}",
        'Accept': 'application/json',
    },
)
response.raise_for_status()
post = response.json()
```

### 其他编辑数据接口

| 接口 | 权限 | 响应 |
| --- | --- | --- |
| `GET /api/posts/{id}/edit/event` | `post:read:self` | `event`，未设置活动时为 `null` |
| `GET /api/posts/{id}/edit/permissions` | `post:read:self` | `canEdit`、`canDelete`、`canEditReportPublicPost`、`canEditLockedPostAsPrivilegedUser`、`canEditPostAsModerator` |

聚合接口 `GET /api/posts/{id}/edit` 已停止服务，会返回 `410` 和 `POST_EDIT_AGGREGATE_RETIRED`，请改用上表中的接口。

## 创建帖子

POST `/api/posts`

需要权限 `post:create` 请求类型 `application/json`

### 完整请求参数

| 参数 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `title` | `string` | 是 | 标题，1～200 字符。 |
| `content` | `object | string` | 是 | 正文，推荐 TipTap JSON；最大约 200,000 字符。 |
| `boardId` | `string` | 是 | 目标板块 ID，最多 50 字符。 |
| `copyrightEnabled` | `boolean` | 否 | 是否显示版权声明。 |
| `copyrightLicense` | `string | null` | 否 | 版权许可证，最多 50 字符。 |
| `isDraft` | `boolean` | 否 | `true` 保存草稿。 |
| `tags` | `array<string>` | 否 | 普通标签，最多 20 个，每项最多 50 字符。 |
| `visibility` | `string` | 否 | `public`、`partial` 或 `private`。 |
| `visibilityAllowUserIds` | `array<string>` | 否 | `partial` 可见用户 ID，最多 100 个。 |
| `accessMode` | `string` | 否 | `none`、`password`、`reply` 或 `login`。 |
| `accessPassword` | `string` | 否 | `password` 模式的访问密码，4～64 字符。 |
| `draftKey` | `string | null` | 否 | 客户端草稿键，最多 80 字符。 |
| `customSlug` | `string | null` | 否 | 自定义路径，最多 120 字符。 |
| `organizationId` | `string | null` | 否 | 个人 Token 不允许指定组织；请省略或传 `null`。 |
| `topicTagIds` | `array<string>` | 否 | 话题标签 ID，最多 5 个，每项最多 80 字符。 |
| `postType` | `string` | 否 | `default`、`qa`、`poll` 或 `reference`。 |
| `referenceTargetType` | `string | null` | 否 | 引用目标类型，`resource` 或 `player`。仅 `reference` 使用。 |
| `referenceTargetId` | `string | null` | 否 | 引用目标 ID，最多 80 字符；目标必须是公开且已审核的资源或找服玩。仅 `reference` 使用。 |
| `qaRewardAmount` | `integer` | 否 | 问答采纳奖励，0～1,000,000。仅 `qa` 使用。 |
| `qaRewardPointSystemId` | `string | null` | 否 | 问答奖励积分类型 ID，最多 80 字符。 |
| `pollOptions` | `array<string>` | 否 | 投票选项，2～20 个，每项最多 120 字符。仅 `poll` 使用。 |
| `pollMode` | `string` | 否 | `single` 或 `multiple`。 |
| `pollMaxChoices` | `integer` | 否 | 最多可选项数，1～20；单选固定为 1。 |
| `pollRevealAt` | `string | null` | 否 | 投票结果公开时间，建议 ISO 8601。 |
| `pollRewardAmount` | `integer` | 否 | 每位投票者奖励，0～1,000,000。 |
| `pollRewardBudget` | `integer` | 否 | 投票奖励总预算，0～1,000,000，不能小于单次奖励。 |
| `pollRewardPointSystemId` | `string | null` | 否 | 投票奖励积分类型 ID，最多 80 字符。 |
| `event` | `object | null` | 否 | 活动信息，与 `postType` 独立设置；子字段见下方说明。 |

### `event` 活动子字段

传 `null` 可以清除已设置的活动。`enabled` 不为 `false` 时必须提供 `startsAt`。

| 子字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `enabled` | `boolean` | 否 | 传 `false` 时不启用活动，其余子字段会被忽略。 |
| `title` | `string` | 否 | 活动标题，最多 160 字符。 |
| `summary` | `string` | 否 | 活动简介，最多 500 字符。 |
| `startsAt` | `string` | 是 | 活动开始时间，建议 ISO 8601。 |
| `endsAt` | `string | null` | 否 | 活动结束时间，建议 ISO 8601；不能早于 `startsAt`。 |
| `timezone` | `string` | 否 | 时区，最多 80 字符；省略时使用 `Asia/Shanghai`。 |
| `locationName` | `string` | 否 | 活动地点名称，最多 160 字符。 |
| `locationUrl` | `string` | 否 | 活动地点链接，最多 500 字符。 |
| `coverImage` | `string` | 否 | 活动封面 URL，最多 500 字符。 |
| `capacity` | `integer | null` | 否 | 人数上限，0～1,000,000。 |
| `rsvpEnabled` | `boolean` | 否 | 是否开启报名。 |
| `status` | `string` | 否 | `scheduled` 或 `cancelled`；省略时为 `scheduled`。 |

### 普通帖示例

```
curl -X POST 'https://www.nexusmc.cn/api/posts' \
  -H 'Authorization: Bearer avm_pat_xxxxx_xxxxx' \
  -H 'Content-Type: application/json' \
  -d '{
    "title": "Hello NexusMC",
    "boardId": "BOARD_ID",
    "content": {"type":"doc","content":[]},
    "postType": "default"
  }'
```

```
const response = await fetch('https://www.nexusmc.cn/api/posts', {
  method: 'POST',
  headers: {
    Authorization: `Bearer ${process.env.NEXUSMC_API_TOKEN}`,
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    title: 'Hello NexusMC',
    boardId: 'BOARD_ID',
    content: { type: 'doc', content: [] },
    postType: 'default',
  }),
})

const created = await response.json()
```

```
import os
import requests

response = requests.post(
    'https://www.nexusmc.cn/api/posts',
    headers={
        'Authorization': f"Bearer {os.environ['NEXUSMC_API_TOKEN']}",
        'Content-Type': 'application/json',
    },
    json={
        'title': 'Hello NexusMC',
        'boardId': 'BOARD_ID',
        'content': {'type': 'doc', 'content': []},
        'postType': 'default',
    },
)
response.raise_for_status()
created = response.json()
```

## 更新帖子

PUT `/api/posts/{id}`

需要权限 `post:update:self` 请求类型 `application/json`

所有字段均为可选，只修改请求里出现的字段。个人 Token 只能更新本人创建的帖子。

NOTE

更新非草稿状态的帖子会重新进入审核，状态回到 `pending`，审核期间前台仍展示旧内容。草稿只有在请求发布时才会提交审核。

### 完整请求参数

| 参数 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `title` | `string` | 否 | 标题，1～200 字符。 |
| `content` | `object | string` | 否 | 正文，推荐 TipTap JSON。 |
| `boardId` | `string` | 否 | 目标板块 ID。 |
| `copyrightEnabled` | `boolean` | 否 | 是否显示版权声明。 |
| `copyrightLicense` | `string | null` | 否 | 版权许可证。 |
| `isDraft` | `boolean` | 否 | 草稿状态。 |
| `tags` | `array<string>` | 否 | 普通标签，最多 20 个。 |
| `visibility` | `string` | 否 | `public`、`partial` 或 `private`。 |
| `visibilityAllowUserIds` | `array<string>` | 否 | `partial` 可见用户 ID，最多 100 个。 |
| `accessMode` | `string` | 否 | `none`、`password`、`reply` 或 `login`。 |
| `accessPassword` | `string` | 否 | 新访问密码；未传时可保留已有密码哈希。 |
| `draftKey` | `string | null` | 否 | 客户端草稿键。 |
| `customSlug` | `string | null` | 否 | 自定义路径。 |
| `organizationId` | `string | null` | 否 | 个人 Token 不允许指定组织；请省略或传 `null`。 |
| `topicTagIds` | `array<string>` | 否 | 用本数组整体替换话题标签，最多 5 个。 |
| `postType` | `string` | 否 | `default`、`qa`、`poll` 或 `reference`。 |
| `referenceTargetType` | `string | null` | 否 | 引用目标类型，`resource` 或 `player`。仅 `reference` 使用。 |
| `referenceTargetId` | `string | null` | 否 | 引用目标 ID，最多 80 字符。仅 `reference` 使用。 |
| `qaRewardAmount` | `integer` | 否 | 问答采纳奖励，0～1,000,000。 |
| `qaRewardPointSystemId` | `string | null` | 否 | 问答奖励积分类型 ID。 |
| `pollOptions` | `array<string>` | 否 | 投票选项，2～20 个。 |
| `pollMode` | `string` | 否 | `single` 或 `multiple`。 |
| `pollMaxChoices` | `integer` | 否 | 最多可选项数，1～20。 |
| `pollRevealAt` | `string | null` | 否 | 投票结果公开时间。 |
| `pollRewardAmount` | `integer` | 否 | 每位投票者奖励。 |
| `pollRewardBudget` | `integer` | 否 | 投票奖励总预算。 |
| `pollRewardPointSystemId` | `string | null` | 否 | 投票奖励积分类型 ID。 |
| `event` | `object | null` | 否 | 活动信息；子字段与「创建帖子」中的 `event` 子字段一致，传 `null` 可清除。 |

### 更新示例

```
curl -X PUT 'https://www.nexusmc.cn/api/posts/POST_ID' \
  -H 'Authorization: Bearer avm_pat_xxxxx_xxxxx' \
  -H 'Content-Type: application/json' \
  -d '{
    "title": "更新后的标题",
    "tags": ["教程", "API"]
  }'
```

```
const response = await fetch('https://www.nexusmc.cn/api/posts/POST_ID', {
  method: 'PUT',
  headers: {
    Authorization: `Bearer ${process.env.NEXUSMC_API_TOKEN}`,
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    title: '更新后的标题',
    tags: ['教程', 'API'],
  }),
})

const updated = await response.json()
```

```
import os
import requests

response = requests.put(
    'https://www.nexusmc.cn/api/posts/POST_ID',
    headers={
        'Authorization': f"Bearer {os.environ['NEXUSMC_API_TOKEN']}",
        'Content-Type': 'application/json',
    },
    json={
        'title': '更新后的标题',
        'tags': ['教程', 'API'],
    },
)
response.raise_for_status()
updated = response.json()
```

## 帖子类型规则

| `postType` | 需要关注的字段 | 规则 |
| --- | --- | --- |
| `default` | 无 | 普通帖子；问答和投票字段会被清空或忽略。 |
| `qa` | `qaRewardAmount`、`qaRewardPointSystemId` | 奖励大于 0 时必须能确定积分类型。 |
| `poll` | `pollOptions`、`pollMode`、`pollMaxChoices` | 至少 2 个选项；多选至少可选 2 项。设置奖励时，总预算不能小于单次奖励。 |
| `reference` | `referenceTargetType`、`referenceTargetId` | 必须指向一个公开且已审核的资源或找服玩；目标不可用时返回 400。 |

切换 `postType` 会清理不属于新类型的数据。数组字段如 `tags`、`topicTagIds`、`pollOptions` 传入后会按新值保存；要保留旧值就不要发送。
