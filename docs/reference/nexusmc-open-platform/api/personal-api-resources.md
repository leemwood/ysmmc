---
title: "个人 API：资源"
source: https://docs.nexusmc.cn/docs/api/personal-api-resources
collected: 2026-09-26
---

# 个人 API：资源

查询、创建和更新资源的接口，以及版本 Tag 和每个请求可直接使用的完整参数表。

## 接口一览

| 方法 | 路径 | 权限 | 用途 |
| --- | --- | --- | --- |
| GET | `/api/resources/version-tags` | `无需 Token` | 读取当前可用的资源版本 Tag |
| GET | `/api/resources/{id}/edit/core` | `resource:read:self` | 读取自己的资源核心编辑数据 |
| GET | `/api/resources/{id}/edit/dependencies` | `resource:read:self` | 读取资源的依赖数据 |
| GET | `/api/resources/{id}/edit/relations` | `resource:read:self` | 读取资源的关联关系 |
| GET | `/api/resources/{id}/edit/tutorials` | `resource:read:self` | 读取资源关联的教程帖 |
| GET | `/api/resources/{id}/edit/documentation` | `resource:read:self` | 读取资源关联的文档帖 |
| GET | `/api/resources/{id}/edit/permissions` | `resource:read:self` | 读取当前用户对资源的编辑权限 |
| GET | `/api/users/me/resources` | `resource:read:self` | 分页查询自己的资源 |
| POST | `/api/resources` | `resource:create` | 创建资源或保存草稿 |
| PUT / PATCH | `/api/resources/{id}` | `resource:update:self` | 更新自己的资源 |
| PUT | `/api/resources/{id}/archive` | `resource:update:self` | 归档或取消归档自己的资源 |
| DELETE | `/api/resources/{id}` | `resource:delete:self` | 删除自己的资源 |

`{id}` 推荐使用列表或创建响应返回的稳定资源 ID，也兼容当前或历史 slug、公开编号和公开链接中的路径段。正文 `content` 的结构见 [投稿正文 TipTap 快速接入](/docs/api/tiptap-content-format/)，文件项和版本 Tag 见 [资源文件与版本](/docs/api/personal-api-resource-files/)。

## 查询自己的资源列表

GET `/api/users/me/resources`

需要权限 `resource:read:self`

```
curl "https://www.nexusmc.cn/api/users/me/resources?page=1&pageSize=20&status=draft" \
  -H "Authorization: Bearer avm_pat_xxxxx_xxxxx" \
  -H "Accept: application/json"
```

```
const url = new URL('https://www.nexusmc.cn/api/users/me/resources')
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

const resources = await response.json()
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
    params={'page': 1, 'pageSize': 20, 'status': 'draft'},
)
response.raise_for_status()
resources = response.json()
```

| 查询参数 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `page` | `number` | 否 | 页码，从 `1` 开始。 |
| `pageSize` | `number` | 否 | 每页数量。 |
| `status` | `string` | 否 | 按审核或发布状态筛选。 |
| `q` | `string` | 否 | 按标题等内容搜索。 |
| `sort` | `string` | 否 | 排序方式。 |

响应的 `resources[]` 每项都包含 `id`、`slug` 和 `path`。更新或发布版本时直接使用 `resources[].id`，无需查看网页 URL。

读取单项编辑数据：

GET `/api/resources/{id}/edit/core`

需要权限 `resource:read:self`

```
curl "https://www.nexusmc.cn/api/resources/RESOURCE_ID/edit/core" \
  -H "Authorization: Bearer avm_pat_xxxxx_xxxxx" \
  -H "Accept: application/json"
```

```
const response = await fetch('https://www.nexusmc.cn/api/resources/RESOURCE_ID/edit/core', {
  headers: {
    Authorization: `Bearer ${process.env.NEXUSMC_API_TOKEN}`,
    Accept: 'application/json',
  },
})

const resource = await response.json()
```

```
import os
import requests

response = requests.get(
    'https://www.nexusmc.cn/api/resources/RESOURCE_ID/edit/core',
    headers={
        'Authorization': f"Bearer {os.environ['NEXUSMC_API_TOKEN']}",
        'Accept': 'application/json',
    },
)
response.raise_for_status()
resource = response.json()
```

### 其他编辑数据接口

`core` 之外的编辑数据按主题拆成了独立接口，权限同为 `resource:read:self`，响应各自只包一层：

| 接口 | 响应字段 |
| --- | --- |
| `GET /api/resources/{id}/edit/dependencies` | `dependencies` |
| `GET /api/resources/{id}/edit/relations` | `relations` |
| `GET /api/resources/{id}/edit/tutorials` | `tutorialPosts` |
| `GET /api/resources/{id}/edit/documentation` | `documentationPosts`、`documentationSections` |
| `GET /api/resources/{id}/edit/permissions` | `viewerPermissions` |

聚合接口 `GET /api/resources/{id}/edit` 已停止服务，会返回 `410` 和 `RESOURCE_EDIT_AGGREGATE_RETIRED`，请改用上表中的接口。

## 获取可用版本 Tag

GET `/api/resources/version-tags`

该接口无需 Token，返回后台当前启用的版本 Tag。客户端在显示发布或更新选项前应先读取本接口，不要假设后台始终只保留默认三项。

```
curl "https://www.nexusmc.cn/api/resources/version-tags" \
  -H "Accept: application/json"
```

```
[
  { "key": "alpha", "label": "开发版", "enabled": true, "sortOrder": 0 },
  { "key": "beta", "label": "测试版", "enabled": true, "sortOrder": 1 },
  { "key": "releases", "label": "正式版", "enabled": true, "sortOrder": 2 }
]
```

提交时传 `key`，例如 `"versionTag": "beta"`。Tag 不存在或已停用时，写接口返回 `400` 和 `所选版本标签不存在或已停用`。

## 创建资源

POST `/api/resources`

需要权限 `resource:create` 请求类型 `application/json`

### 完整请求参数

| 参数 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `title` | `string` | 是 | 标题，1～200 字符。 |
| `content` | `object | string` | 是 | 正文，推荐传 TipTap JSON；最大约 200,000 字符。 |
| `authorContribution` | `string` | 否 | 作者贡献说明，最多 500 字符。 |
| `category` | `string` | 是 | 资源分类 ID 或分类值。 |
| `description` | `string` | 否 | 简介，最多 1,000 字符。 |
| `platform` | `string` | 否 | 平台，最多 20 字符。 |
| `subCategory` | `string | array<string>` | 否 | Loader/子分类。 |
| `sideSupport` | `string` | 否 | 客户端、服务端等运行侧支持。 |
| `mcVersions` | `array<string>` | 否 | 支持的 Minecraft 版本。 |
| `gameVersions` | `array<string>` | 否 | `mcVersions` 的等价别名；两者同时传入时以 `mcVersions` 为准。 |
| `tags` | `array<string>` | 否 | 自定义标签。 |
| `officialTags` | `array<string>` | 否 | 官方标签。 |
| `customCategorySelections` | `object` | 否 | 自定义分类选择，形如 `{ "分组 ID": ["值"] }`；每个分组最多 20 项，每项最多 80 字符。 |
| `fileUrl` | `string | null` | 否 | 外链下载地址；本地托管优先使用 `files`。 |
| `extractCode` | `string | null` | 否 | 外部网盘提取码，最多 100 字符。 |
| `fileSize` | `number` | 否 | 兼容的单文件大小，单位字节。 |
| `fileName` | `string` | 否 | 兼容的单文件名，最多 200 字符。 |
| `fileSha256` | `string` | 否 | 单文件 SHA-256，最多 64 字符。 |
| `fileSha1` | `string` | 否 | 单文件 SHA-1，最多 40 字符。 |
| `files` | `array<object>` | 否 | 版本文件列表。下方列出常用子字段，完整规则见[资源文件与版本](/docs/api/personal-api-resource-files/)。 |
| `dependencies` | `array` | 否 | 资源依赖关系。 |
| `relations` | `array` | 否 | 创建时写入的其他资源关系。 |
| `tutorialPostIds` | `array<string>` | 否 | 关联教程帖 ID。 |
| `documentationPostRefs` | `array` | 否 | 关联文档帖引用。 |
| `coverImage` | `string | null` | 否 | 封面图片 URL。 |
| `downloadType` | `string` | 否 | `local` 或 `external`。 |
| `additionalFiles` | `array` | 否 | 兼容的附加文件列表。 |
| `version` | `string` | 否 | 初始版本号，最多 50 字符。 |
| `versionTag` | `string` | 否 | 初始版本 Tag，最多 32 字符；必须是当前启用的 `key`，省略时使用 `releases`。 |
| `thirdLevelCats` | `string | array<string>` | 否 | 三级分类。 |
| `sourceType` | `string` | 否 | 来源类型，最多 20 字符。 |
| `repositoryUrl` | `string | null` | 否 | 源代码仓库 URL。 |
| `communicationUrl` | `string | null` | 否 | 交流渠道 URL。 |
| `documentationUrl` | `string | null` | 否 | 外部文档 URL。 |
| `isPaid` | `boolean` | 否 | 是否为付费资源。 |
| `attributionType` | `string` | 否 | 归属类型，如原创或转载。 |
| `originalResourceUrl` | `string | null` | 否 | 转载时的原资源 URL；转载资源必须填写。 |
| `originalAuthorName` | `string | null` | 否 | 原作者名称，最多 120 字符。 |
| `originalAuthorLabel` | `string | null` | 否 | 详情页显示的作者称谓，最多 30 字符，例如“主要作者”或“上游作者”；留空时按归属类型自动显示。 |
| `originalAuthorUrl` | `string | null` | 否 | 原作者主页 URL。 |
| `visibility` | `string` | 否 | `public`、`partial` 或 `private`。 |
| `visibilityAllowUserIds` | `array<string>` | 否 | `partial` 可见用户 ID，最多 100 个。 |
| `accessMode` | `string` | 否 | `none`、`password`、`reply` 或 `login`。 |
| `accessPassword` | `string` | 否 | `password` 模式的访问密码，4～64 字符。 |
| `isDraft` | `boolean` | 否 | `true` 保存草稿；否则进入正常发布流程。 |
| `galleryMode` | `string` | 否 | `auto` 或 `manual`。 |
| `galleryImages` | `array<string>` | 否 | 图库图片 URL。 |
| `contributors` | `array` | 否 | 贡献者列表；转载资源会忽略此项。 |
| `customSlug` | `string` | 否 | 自定义路径，最多 120 字符。 |
| `organizationId` | `string | null` | 否 | 个人 Token 不允许指定组织；请省略或传 `null`。 |
| `licenseId` | `string | null` | 否 | 许可证 ID。 |

### `files[]` 常用子字段

创建资源时如果要同时提交站内托管文件，`files` 里的每一项至少要能说明文件地址和文件名。

| 子字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `url` | `string` | 是 | 文件地址。站点托管文件使用上传接口返回的 `url`。 |
| `fileName` | `string` | 是 | 显示文件名。上传响应里的 `filename` 要映射到这里。 |
| `fileSize` | `number` | 否 | 文件大小，单位字节。上传响应里的 `size` 要映射到这里。 |
| `isPrimary` | `boolean` | 否 | 是否作为默认下载文件。多个文件里最终只保留一个主文件。 |
| `subcategoryIds` | `array<string>` | 否 | 这个文件支持的 Loader/子分类。 |
| `gameVersions` | `array<string>` | 否 | 这个文件支持的 Minecraft 版本。 |
| `sha256` | `string` | 否 | 64 位 SHA-256。 |
| `sha1` | `string` | 否 | 40 位 SHA-1。 |

### 最小示例

```
curl -X POST 'https://www.nexusmc.cn/api/resources' \
  -H 'Authorization: Bearer avm_pat_xxxxx_xxxxx' \
  -H 'Content-Type: application/json' \
  -d '{
    "title": "示例资源",
    "category": "mod",
    "content": {"type":"doc","content":[]},
    "versionTag": "releases",
    "isDraft": true
  }'
```

```
const response = await fetch('https://www.nexusmc.cn/api/resources', {
  method: 'POST',
  headers: {
    Authorization: `Bearer ${process.env.NEXUSMC_API_TOKEN}`,
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    title: '示例资源',
    category: 'mod',
    content: { type: 'doc', content: [] },
    versionTag: 'releases',
    isDraft: true,
  }),
})

const created = await response.json()
```

```
import os
import requests

response = requests.post(
    'https://www.nexusmc.cn/api/resources',
    headers={
        'Authorization': f"Bearer {os.environ['NEXUSMC_API_TOKEN']}",
        'Content-Type': 'application/json',
    },
    json={
        'title': '示例资源',
        'category': 'mod',
        'content': {'type': 'doc', 'content': []},
        'versionTag': 'releases',
        'isDraft': True,
    },
)
response.raise_for_status()
created = response.json()
```

## 更新资源

PUT / PATCH `/api/resources/{id}`

需要权限 `resource:update:self` 请求类型 `application/json`

`PUT` 和 `PATCH` 在此接口中都按部分更新处理：只修改请求里出现的字段。个人 Token 只能更新本人创建的资源。

NOTE

更新已发布（`approved`）或已拒（`rejected`）的资源会重新进入审核，状态回到 `pending`，审核期间前台仍展示旧内容。只修改归档字段，或同时把 `visibility` 改为 `private`，不会触发重新审核。

### 完整请求参数

| 参数 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `title` | `string` | 否 | 标题，1～200 字符。 |
| `content` | `object | string` | 否 | 正文；推荐传 TipTap JSON。 |
| `authorContribution` | `string` | 否 | 作者贡献说明，最多 500 字符。 |
| `category` | `string` | 否 | 资源分类 ID 或分类值。 |
| `description` | `string` | 否 | 简介，最多 1,000 字符。 |
| `platform` | `string` | 否 | 平台，最多 20 字符。 |
| `subCategory` | `string | array<string>` | 否 | Loader/子分类。 |
| `sideSupport` | `string` | 否 | 运行侧支持。 |
| `mcVersions` | `array<string>` | 否 | 支持的 Minecraft 版本。 |
| `gameVersions` | `array<string>` | 否 | `mcVersions` 的等价别名；两者同时传入时以 `mcVersions` 为准。 |
| `tags` | `array<string>` | 否 | 自定义标签。 |
| `officialTags` | `array<string>` | 否 | 官方标签。 |
| `customCategorySelections` | `object` | 否 | 自定义分类选择，形如 `{ "分组 ID": ["值"] }`；每个分组最多 20 项，每项最多 80 字符。 |
| `fileUrl` | `string | null` | 否 | 外链下载地址。 |
| `extractCode` | `string | null` | 否 | 外部网盘提取码。 |
| `fileSize` | `number` | 否 | 兼容的单文件大小，单位字节。 |
| `fileName` | `string` | 否 | 兼容的单文件名。 |
| `fileSha256` | `string` | 否 | 单文件 SHA-256。 |
| `fileSha1` | `string` | 否 | 单文件 SHA-1。 |
| `files` | `array<object>` | 否 | 新版本文件列表；按上方 `files[]` 子字段提交。数组传入后会替换当前版本文件。 |
| `dependencies` | `array` | 否 | 用本数组整体替换依赖关系。 |
| `tutorialPostIds` | `array<string>` | 否 | 用本数组整体替换教程帖关系。 |
| `documentationPostRefs` | `array` | 否 | 用本数组整体替换文档帖引用。 |
| `coverImage` | `string | null` | 否 | 封面图片 URL。 |
| `downloadType` | `string` | 否 | `local` 或 `external`。 |
| `additionalFiles` | `array` | 否 | 兼容的附加文件列表。 |
| `version` | `string` | 否 | 版本号，最多 50 字符。 |
| `versionTag` | `string` | 否 | 当前资源或新版本的 Tag，最多 32 字符；必须使用当前启用的 `key`。 |
| `publishVersion` | `boolean` | 否 | 是否显式创建版本记录。 |
| `versionTitle` | `string` | 否 | 版本标题，最多 200 字符。 |
| `changelog` | `object | string` | 否 | 版本更新说明。 |
| `thirdLevelCats` | `string | array<string>` | 否 | 三级分类。 |
| `sourceType` | `string` | 否 | 来源类型。 |
| `repositoryUrl` | `string | null` | 否 | 源代码仓库 URL。 |
| `communicationUrl` | `string | null` | 否 | 交流渠道 URL。 |
| `documentationUrl` | `string | null` | 否 | 外部文档 URL。 |
| `isPaid` | `boolean` | 否 | 是否为付费资源。 |
| `attributionType` | `string` | 否 | 归属类型。 |
| `originalResourceUrl` | `string | null` | 否 | 转载时的原资源 URL。 |
| `originalAuthorName` | `string | null` | 否 | 原作者名称。 |
| `originalAuthorLabel` | `string | null` | 否 | 详情页显示的作者称谓，最多 30 字符；传空值可恢复自动称谓。 |
| `originalAuthorUrl` | `string | null` | 否 | 原作者主页 URL。 |
| `visibility` | `string` | 否 | `public`、`partial` 或 `private`。 |
| `visibilityAllowUserIds` | `array<string>` | 否 | `partial` 可见用户 ID，最多 100 个。 |
| `accessMode` | `string` | 否 | `none`、`password`、`reply` 或 `login`。 |
| `accessPassword` | `string` | 否 | 新访问密码；未传时可保留已有密码哈希。 |
| `isDraft` | `boolean` | 否 | 草稿状态。 |
| `archived` | `boolean` | 否 | 是否归档。 |
| `archivedAt` | `string | null` | 否 | 归档时间，建议 ISO 8601。 |
| `archiveNote` | `string | null` | 否 | 归档说明，最多 500 字符。 |
| `galleryMode` | `string` | 否 | `auto` 或 `manual`。 |
| `galleryImages` | `array<string>` | 否 | 图库图片 URL。 |
| `contributors` | `array` | 否 | 用本数组更新贡献者。 |
| `customSlug` | `string` | 否 | 自定义路径。 |
| `organizationId` | `string | null` | 否 | 个人 Token 不允许指定组织；请省略或传 `null`。 |
| `licenseId` | `string | null` | 否 | 许可证 ID。 |

CAUTION

更新接口不接收 `relations`。资源关系请使用对应关系接口，不要把创建请求原样复用到更新请求。

### 更新示例

```
curl -X PATCH 'https://www.nexusmc.cn/api/resources/RESOURCE_ID' \
  -H 'Authorization: Bearer avm_pat_xxxxx_xxxxx' \
  -H 'Content-Type: application/json' \
  -d '{
    "description": "新的简介",
    "mcVersions": ["1.21.1"],
    "tags": ["科技", "自动化"]
  }'
```

```
const response = await fetch('https://www.nexusmc.cn/api/resources/RESOURCE_ID', {
  method: 'PATCH',
  headers: {
    Authorization: `Bearer ${process.env.NEXUSMC_API_TOKEN}`,
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    description: '新的简介',
    mcVersions: ['1.21.1'],
    tags: ['科技', '自动化'],
  }),
})

const updated = await response.json()
```

```
import os
import requests

response = requests.patch(
    'https://www.nexusmc.cn/api/resources/RESOURCE_ID',
    headers={
        'Authorization': f"Bearer {os.environ['NEXUSMC_API_TOKEN']}",
        'Content-Type': 'application/json',
    },
    json={
        'description': '新的简介',
        'mcVersions': ['1.21.1'],
        'tags': ['科技', '自动化'],
    },
)
response.raise_for_status()
updated = response.json()
```

## 数组字段的更新方式

`dependencies`、`tutorialPostIds`、`documentationPostRefs`、`galleryImages` 等数组字段传入后通常会整体替换旧值。要清空时传空数组；要保留时不要发送该字段。

## 归档资源

PUT `/api/resources/{id}/archive`

需要权限 `resource:update:self` 请求类型 `application/json`

归档只改变展示状态，不会删除资源，也不会触发重新审核。

| 参数 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `archived` | `boolean` | 否 | `true` 归档；`false` 取消归档。 |
| `archivedAt` | `string | null` | 否 | 归档时间，建议 ISO 8601，最多 80 字符。 |
| `archiveNote` | `string | null` | 否 | 归档说明，最多 500 字符。 |

响应返回 `id`、`slug`、`publicId`、`title`、`status`、`archivedAt`、`archiveNote`、`updatedAt` 和 `organizationId`。

CAUTION

归档状态变化且资源属于 Token 用户本人时，后端会要求强验证。未提前在站内完成验证会返回 `403`，`code` 为 `SECURITY_SETUP_REQUIRED`（账号未绑定二步验证、通行密钥或双层密码）或 `STRONG_ACTION_VERIFICATION_REQUIRED`（需要在站内重新完成安全验证）。这条接口不适合完全无人值守的流水线。

```
curl -X PUT 'https://www.nexusmc.cn/api/resources/RESOURCE_ID/archive' \
  -H 'Authorization: Bearer avm_pat_xxxxx_xxxxx' \
  -H 'Content-Type: application/json' \
  -d '{
    "archived": true,
    "archiveNote": "停止维护"
  }'
```

```
const response = await fetch('https://www.nexusmc.cn/api/resources/RESOURCE_ID/archive', {
  method: 'PUT',
  headers: {
    Authorization: `Bearer ${process.env.NEXUSMC_API_TOKEN}`,
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    archived: true,
    archiveNote: '停止维护',
  }),
})

const result = await response.json()
```

```
import os
import requests

response = requests.put(
    'https://www.nexusmc.cn/api/resources/RESOURCE_ID/archive',
    headers={
        'Authorization': f"Bearer {os.environ['NEXUSMC_API_TOKEN']}",
        'Content-Type': 'application/json',
    },
    json={
        'archived': True,
        'archiveNote': '停止维护',
    },
)
response.raise_for_status()
result = response.json()
```

## 删除资源

DELETE `/api/resources/{id}`

需要权限 `resource:delete:self`

个人 Token 只能删除 Token 所属用户以个人身份创建的资源，不能删除组织身份投稿，也不会继承协作者、版主或管理员权限。

```
curl -X DELETE 'https://www.nexusmc.cn/api/resources/RESOURCE_ID' \
  -H 'Authorization: Bearer avm_pat_xxxxx_xxxxx'
```

```
const response = await fetch('https://www.nexusmc.cn/api/resources/RESOURCE_ID', {
  method: 'DELETE',
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

response = requests.delete(
    'https://www.nexusmc.cn/api/resources/RESOURCE_ID',
    headers={
        'Authorization': f"Bearer {os.environ['NEXUSMC_API_TOKEN']}",
        'Accept': 'application/json',
    },
)
response.raise_for_status()
result = response.json()
```

GitHub Actions 等流水线如果需要自动发布和清理资源，建议单独创建一个 Token，只勾选实际需要的 `upload:file`、`resource:create`、`resource:update:self`、`resource:delete:self`。

提交新文件或发布版本前，请继续阅读[资源文件与版本](/docs/api/personal-api-resource-files/)，其中单独列出了 `files` 每个子字段、主文件 Loader 规则和版本接口。
