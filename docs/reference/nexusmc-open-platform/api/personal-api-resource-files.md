---
title: "个人 API：资源文件与版本"
source: https://docs.nexusmc.cn/docs/api/personal-api-resource-files
collected: 2026-09-26
---

# 个人 API：资源文件与版本

files 子结构、版本 Tag、主文件、Loader 继承、版本提交和文件审核切换机制。

## 先理解文件工作流

站点托管文件需要两步：

1.  调用[文件上传接口](/docs/api/personal-api-uploads/)取得 `url`、`filename`、`size` 和哈希。
2.  把上传结果组装成 `files`，提交到资源创建、资源更新或版本接口。

不要提交本机路径。`C:\\build\\mod.jar` 或 `./dist/mod.jar` 对服务端没有意义。

## `files` 完整参数表

| 参数 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `url` | `string` | 是 | 文件地址。站点托管文件必须使用上传接口返回的 URL。兼容别名：`fileUrl`。 |
| `fileName` | `string` | 是 | 显示文件名。兼容别名：`filename`。 |
| `fileSize` | `number` | 否 | 文件大小，单位字节。兼容别名：`size`。 |
| `isPrimary` | `boolean` | 否 | 是否为默认下载文件。每个版本最终只保留一个主文件；兼容别名：`isMain`、`primary`。 |
| `subcategoryIds` | `array<string>` | 否 | 该文件支持的 Loader/子分类 ID。兼容别名：`loaderIds`、`subCategoryIds`、`loaders`。 |
| `gameVersions` | `array<string>` | 否 | 该文件支持的游戏版本。兼容别名：`mcVersions`。 |
| `extractCode` | `string` | 否 | 外部网盘提取码。站点本地文件通常不需要。 |
| `sha256` | `string` | 否 | 64 位 SHA-256。兼容别名：`fileSha256`。 |
| `sha1` | `string` | 否 | 40 位 SHA-1。兼容别名：`fileSha1`。 |
| `createdAt` | `string` | 否 | 文件时间；通常由服务端生成。 |

推荐统一使用左侧列出的规范字段，不要在同一个客户端里混用别名。

## 主文件与 Loader 继承

`isPrimary` 只决定默认下载文件，不会强制主文件使用资源自身的 Loader。

每个文件分别按以下顺序解析：

1.  使用文件自身非空的 `subcategoryIds`；为空时检查 `loaderIds` 等兼容别名。
2.  文件完全没有 Loader 时，才回退到资源级 `subCategory`。
3.  使用文件自身非空的 `gameVersions`；为空时检查 `mcVersions`。
4.  文件完全没有游戏版本时，才回退到请求级或资源级 `mcVersions`。

```
{
  "url": "/uploads/files/build.jar",
  "fileName": "build.jar",
  "isPrimary": true,
  "loaderIds": ["neoforge"],
  "gameVersions": ["1.21.1"]
}
```

上面的主文件只标记为 `neoforge`。即使客户端同时生成空的 `subcategoryIds: []`，服务端也会采用非空的 `loaderIds`。

历史版本

较早创建且没有保存主文件级 Loader 的版本无法自动还原原始选择，读取时仍会用资源级 Loader 兜底。需要纠正时，请重新提交包含完整 `files` 的版本。

## 本地下载与外链下载

| `downloadType` | 应提交的字段 | 行为 |
| --- | --- | --- |
| `local` | `files` | 使用站点上传结果；`extractCode` 会被忽略或清空。 |
| `external` | `fileUrl`，可选 `extractCode` | 使用外部 HTTP/HTTPS 下载地址。 |

从本地切换到外链或从外链切换到本地时，必须在同一个请求中同时更新 `downloadType` 和对应文件字段。

## 通过资源更新接口提交版本

PATCH `/api/resources/{id}`

需要权限 `resource:update:self` 请求类型 `application/json`

```
curl -X PATCH "https://www.nexusmc.cn/api/resources/RESOURCE_ID" \
  -H "Authorization: Bearer avm_pat_xxxxx_xxxxx" \
  -H "Content-Type: application/json" \
  -d '{
    "publishVersion": true,
    "version": "1.1.0",
    "versionTag": "beta",
    "versionTitle": "NeoForge 构建",
    "changelog": "支持 1.21.1。",
    "files": [{
      "url": "/uploads/files/build.jar",
      "fileName": "build.jar",
      "fileSize": 204800,
      "isPrimary": true,
      "loaderIds": ["neoforge"],
      "gameVersions": ["1.21.1"]
    }]
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
    publishVersion: true,
    version: '1.1.0',
    versionTag: 'beta',
    versionTitle: 'NeoForge 构建',
    changelog: '支持 1.21.1。',
    files: [{
      url: '/uploads/files/build.jar',
      fileName: 'build.jar',
      fileSize: 204800,
      isPrimary: true,
      loaderIds: ['neoforge'],
      gameVersions: ['1.21.1'],
    }],
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
        'publishVersion': True,
        'version': '1.1.0',
        'versionTag': 'beta',
        'versionTitle': 'NeoForge 构建',
        'changelog': '支持 1.21.1。',
        'files': [{
            'url': '/uploads/files/build.jar',
            'fileName': 'build.jar',
            'fileSize': 204800,
            'isPrimary': True,
            'loaderIds': ['neoforge'],
            'gameVersions': ['1.21.1'],
        }],
    },
)
response.raise_for_status()
updated = response.json()
```

非草稿资源的请求满足任一条件时，会创建版本记录：

-   `publishVersion` 为 `true`
-   请求包含 `files`、`fileUrl`、`fileName`、`fileSize`、`fileSha256`、`fileSha1`、`additionalFiles` 或 `extractCode`

### 版本相关参数

| 参数 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `publishVersion` | `boolean` | 否 | 传 `true` 明确提交新版本。 |
| `version` | `string` | 否 | 新版本号，最长 50 字符。 |
| `versionTag` | `string` | 否 | 新版本 Tag；必须使用 `/api/resources/version-tags` 返回的启用 `key`。省略时沿用资源当前 Tag。 |
| `versionTitle` | `string` | 否 | 新版本标题，最长 200 字符。 |
| `changelog` | `string` 或 `object` | 否 | 更新日志；对象可使用 TipTap 文档。 |
| `downloadType` | `string` | 否 | `local` 或 `external`。 |
| `files` | `array` | 否 | 完整文件列表，推荐。 |
| `fileUrl` | `string` | 否 | 传统主文件或外链字段。 |
| `fileName` | `string` | 否 | 传统主文件名。 |
| `fileSize` | `number` | 否 | 传统主文件字节数。 |
| `fileSha256` | `string` | 否 | 传统主文件 SHA-256。 |
| `fileSha1` | `string` | 否 | 传统主文件 SHA-1。 |
| `additionalFiles` | `array` | 否 | 旧客户端附加文件字段。新客户端使用 `files`。 |
| `extractCode` | `string` | 否 | 外链提取码。 |
| `mcVersions` | `array<string>` | 否 | 版本级游戏版本列表。 |

需要审核的版本只有在审核通过后才切换资源主文件。

## 独立版本接口

POST `/api/resources/{id}/versions`

需要权限 `resource:update:self` 请求类型 `application/json`

```
curl -X POST "https://www.nexusmc.cn/api/resources/RESOURCE_ID/versions" \
  -H "Authorization: Bearer avm_pat_xxxxx_xxxxx" \
  -H "Content-Type: application/json" \
  -d '{
    "version": "1.1.0",
    "versionTag": "releases",
    "title": "NeoForge 构建",
    "changelog": "支持 1.21.1。",
    "downloadType": "local",
    "files": [{
      "url": "/uploads/files/build.jar",
      "fileName": "build.jar",
      "fileSize": 204800,
      "isPrimary": true,
      "loaderIds": ["neoforge"],
      "gameVersions": ["1.21.1"]
    }]
  }'
```

```
const response = await fetch('https://www.nexusmc.cn/api/resources/RESOURCE_ID/versions', {
  method: 'POST',
  headers: {
    Authorization: `Bearer ${process.env.NEXUSMC_API_TOKEN}`,
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    version: '1.1.0',
    versionTag: 'releases',
    title: 'NeoForge 构建',
    changelog: '支持 1.21.1。',
    downloadType: 'local',
    files: [{
      url: '/uploads/files/build.jar',
      fileName: 'build.jar',
      fileSize: 204800,
      isPrimary: true,
      loaderIds: ['neoforge'],
      gameVersions: ['1.21.1'],
    }],
  }),
})

const version = await response.json()
```

```
import os
import requests

response = requests.post(
    'https://www.nexusmc.cn/api/resources/RESOURCE_ID/versions',
    headers={
        'Authorization': f"Bearer {os.environ['NEXUSMC_API_TOKEN']}",
        'Content-Type': 'application/json',
    },
    json={
        'version': '1.1.0',
        'versionTag': 'releases',
        'title': 'NeoForge 构建',
        'changelog': '支持 1.21.1。',
        'downloadType': 'local',
        'files': [{
            'url': '/uploads/files/build.jar',
            'fileName': 'build.jar',
            'fileSize': 204800,
            'isPrimary': True,
            'loaderIds': ['neoforge'],
            'gameVersions': ['1.21.1'],
        }],
    },
)
response.raise_for_status()
version = response.json()
```

| 参数 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `version` | `string` | 是 | 新版本号，1 到 50 字符。 |
| `versionTag` | `string` | 否 | 新版本 Tag；必须是当前启用的 `key`，省略时使用 `releases`。 |
| `title` | `string` | 否 | 版本标题，最长 200 字符。 |
| `changelog` | `string` 或 `object` | 否 | 更新日志。 |
| `downloadType` | `string` | 否 | `local` 或 `external`；省略时沿用资源设置。 |
| `files` | `array` | 否 | 完整文件列表，支持每个文件独立设置 Loader 和游戏版本。 |
| `fileUrl` | `string` | 否 | 传统主文件或外链地址。 |
| `fileName` | `string` | 否 | 传统主文件名，最长 200 字符。 |
| `fileSize` | `number` | 否 | 传统主文件大小。 |
| `fileSha256` | `string` | 否 | 传统主文件 SHA-256。 |
| `fileSha1` | `string` | 否 | 传统主文件 SHA-1。 |
| `extractCode` | `string` | 否 | 外链提取码，最长 100 字符。 |
| `additionalFiles` | `array` | 否 | 旧客户端附加文件列表。 |
| `mcVersions` | `array<string>` | 否 | 版本级游戏版本；最终至少保留一项。 |

该接口同样执行文件级优先、资源级兜底的继承规则。

## 读取版本

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| `GET` | `/api/resources/{id}/versions` | 获取可见版本列表和各版本 `files`，支持查询参数筛选。 |
| `GET` | `/api/resources/{id}/versions/{versionId}` | 获取已通过审核的单个版本详情。 |
| `GET` | `/api/resources/{id}/files` | 获取按版本分组的已审核下载文件。 |

读取版本列表：

```
curl "https://www.nexusmc.cn/api/resources/RESOURCE_ID/versions?platform=java&mcVersion=1.21.1&loader=fabric&status=approved" \
  -H "Accept: application/json"
```

```
const params = new URLSearchParams({ platform: 'java', mcVersion: '1.21.1', loader: 'fabric', status: 'approved' })
const response = await fetch(`https://www.nexusmc.cn/api/resources/RESOURCE_ID/versions?${params}`)

const versions = await response.json()
```

```
import requests

response = requests.get(
    'https://www.nexusmc.cn/api/resources/RESOURCE_ID/versions',
    params={'platform': 'java', 'mcVersion': '1.21.1', 'loader': 'fabric', 'status': 'approved'},
)
response.raise_for_status()
versions = response.json()
```

读取单个版本：

```
curl "https://www.nexusmc.cn/api/resources/RESOURCE_ID/versions/VERSION_ID" \
  -H "Accept: application/json"
```

```
const response = await fetch('https://www.nexusmc.cn/api/resources/RESOURCE_ID/versions/VERSION_ID')

const version = await response.json()
```

```
import requests

response = requests.get(
    'https://www.nexusmc.cn/api/resources/RESOURCE_ID/versions/VERSION_ID',
)
response.raise_for_status()
version = response.json()
```

读取按版本分组的文件：

```
curl "https://www.nexusmc.cn/api/resources/RESOURCE_ID/files" \
  -H "Accept: application/json"
```

```
const response = await fetch('https://www.nexusmc.cn/api/resources/RESOURCE_ID/files')

const files = await response.json()
```

```
import requests

response = requests.get(
    'https://www.nexusmc.cn/api/resources/RESOURCE_ID/files',
)
response.raise_for_status()
files = response.json()
```

新客户端应优先读取响应中的 `files`。`fileUrl` 和 `additionalFiles` 只用于兼容旧客户端。

版本列表中的每一项还会返回 `versionTag`。`GET /api/resources/{id}/versions` 支持以下可选查询参数；不传时返回全部可见版本。普通用户只可看到已通过审核的版本，作者和有权限的审核者还可看到待审及被拒版本（不含未通过审核的文件信息）。

公开更新检查无需 Token。要用个人 Token 查看**自己**资源的审核进度，Token 需有 `resource:read:self`；`status=pending` 或 `status=rejected` 不能让其他用户看到未公开版本。待审和被拒版本会返回 `status`、`reviewNote`、`createdAt`、`reviewedAt` 等审核信息，但 `files` 为空；单版本详情和 `/files` 始终只返回已审核文件。

| 查询参数 | 匹配方式 |
| --- | --- |
| `platform` | 资源平台，例如 `java`、`bedrock`；资源平台不同则返回空数组。 |
| `mcVersion` | 精确匹配文件的游戏版本，例如 `1.21.1`；旧版本无文件元数据时回退到版本级 `mcVersions`。 |
| `loader` | 精确匹配文件的 `subcategoryIds`，例如 `fabric`；与 `mcVersion` 同时使用时必须由同一个文件满足。 |
| `versionTag` | 精确匹配版本标签；历史空值按 `releases` 处理。 |
| `status` | `approved`、`pending` 或 `rejected`；仍受版本审核可见性约束。更新检查建议传 `approved`。 |

筛选只决定返回哪些版本，不会删减命中版本的 `files`。客户端应在返回的 `files` 中选择对应游戏版本和 Loader 的文件。

需要进一步在本地比较版本号或时间时，使用以下响应字段：

| 筛选条件 | 响应字段 |
| --- | --- |
| 版本 Tag | `versionTag`；历史空值可按 `releases` 处理。 |
| 版本号 | `version` |
| 审核通过时间 | `reviewedAt`；公开列表按它倒序，空值再按 `createdAt` 排序。 |
| 提交时间 | `createdAt`；作者可见的审核列表按它倒序。 |
| 子分类 / Loader | `files[].subcategoryIds` |
| 游戏版本 | `files[].gameVersions`；旧版本可回退读取版本级 `mcVersions`。 |

例如更新检查可以取筛选结果的第一项作为**最近审核通过**的版本，比较其 `id` 与本地保存的版本 ID。`version` 是展示用版本号，服务端不按 SemVer 排序。即使查询时指定了 Loader 和 MC 版本，响应仍含该版本的全部文件；下载前应从 `files` 中再次选择同时匹配两项的文件。

查看本人待审版本：

```
curl "https://www.nexusmc.cn/api/resources/RESOURCE_ID/versions?status=pending" \
  -H "Authorization: Bearer avm_pat_xxxxx_xxxxx" \
  -H "Accept: application/json"
```
