---
title: "个人 API：资源上传与发布流程"
source: https://docs.nexusmc.cn/docs/api/personal-api-resource-upload-flow
collected: 2026-09-26
---

# 个人 API：资源上传与发布流程

按两段式流程说明如何先上传文件，再创建资源、设置版本 Tag、局部更新资源或提交新版本。

## 适用范围

这篇文档适合要用脚本发布或更新资源的开发者。完整字段表见[个人 API：资源](/docs/api/personal-api-resources/)，上传端点细节见[个人 API：文件上传](/docs/api/personal-api-uploads/)，文件数组和版本规则见[个人 API：资源文件与版本](/docs/api/personal-api-resource-files/)。

资源发布 API 采用两段式流程：

1.  先调用 `/api/upload` 系列接口，把文件或图片传到站点存储。
2.  再调用资源创建、资源更新或版本接口，提交上传后返回的 URL、文件名和大小。

不要自行拼接 `/uploads/files/...`、域名或文件路径。站点可能把实际文件名改成随机短名，也可能按内部存储规则生成完整 URL；客户端必须使用上传接口响应里的 `url`。

## 需要的权限

| 操作 | 权限 |
| --- | --- |
| 上传文件、图片、直传或分块上传 | `upload:file` |
| 创建资源 | `resource:create` |
| 更新自己的资源或提交版本 | `resource:update:self` |
| 删除自己的个人身份资源 | `resource:delete:self` |

个人 API Token 只代表 Token 所属用户本人，不能用组织身份投稿，也不会继承管理员或版主权限。

## 上传文件

单文件上传：

POST `/api/upload`

需要权限 `upload:file` 请求类型 `multipart/form-data`

```
curl -X POST "https://www.nexusmc.cn/api/upload" \
  -H "Authorization: Bearer avm_pat_xxxxx_xxxxx" \
  -F "file=@./ApolloSupport-1.0.0.jar"
```

```
import { openAsBlob } from 'node:fs'

const form = new FormData()
form.append('file', await openAsBlob('./ApolloSupport-1.0.0.jar'), 'ApolloSupport-1.0.0.jar')

const response = await fetch('https://www.nexusmc.cn/api/upload', {
  method: 'POST',
  headers: {
    Authorization: `Bearer ${process.env.NEXUSMC_API_TOKEN}`,
  },
  body: form,
})

const uploaded = await response.json()
```

```
import os
import requests

with open('./ApolloSupport-1.0.0.jar', 'rb') as file:
    response = requests.post(
        'https://www.nexusmc.cn/api/upload',
        headers={'Authorization': f"Bearer {os.environ['NEXUSMC_API_TOKEN']}"},
        files={'file': ('ApolloSupport-1.0.0.jar', file)},
    )

response.raise_for_status()
uploaded = response.json()
```

典型响应：

```
{
  "url": "/uploads/files/uuhUuBftBa0P530_VsJ64.jar",
  "filename": "ApolloSupport-1.0.0.jar",
  "size": 204800,
  "sha256": "64位sha256",
  "sha1": "40位sha1"
}
```

字段映射：

| 上传响应字段 | 资源字段 | 说明 |
| --- | --- | --- |
| `url` | `files[].url` 或兼容字段 `fileUrl` | 服务端最终可访问地址，必须原样使用。 |
| `filename` | `files[].fileName` 或兼容字段 `fileName` | 用户原始文件名，也是下载时应展示的文件名。 |
| `size` | `files[].fileSize` 或兼容字段 `fileSize` | 文件大小，单位字节。 |
| `sha256` | `files[].sha256` 或兼容字段 `fileSha256` | SHA-256 校验值。 |
| `sha1` | `files[].sha1` 或兼容字段 `fileSha1` | SHA-1 校验值。 |

上传响应里的 `url` 可能是 `/uploads/...` 相对路径，也可能是站点返回的完整 URL。两种情况都按响应值提交，不需要客户端判断站点内部如何保存文件。

## 创建资源

POST `/api/resources`

需要权限 `resource:create` 请求类型 `application/json`

推荐使用 `files` 数组提交站内托管文件：

```
curl -X POST "https://www.nexusmc.cn/api/resources" \
  -H "Authorization: Bearer avm_pat_xxxxx_xxxxx" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "ApolloSupport",
    "description": "示例资源",
    "content": {
      "type": "doc",
      "content": [
        {
          "type": "paragraph",
          "content": [{ "type": "text", "text": "资源正文。" }]
        }
      ]
    },
    "platform": "java",
    "category": "plugin",
    "downloadType": "local",
    "version": "1.0.0",
    "files": [
      {
        "url": "/uploads/files/uuhUuBftBa0P530_VsJ64.jar",
        "fileName": "ApolloSupport-1.0.0.jar",
        "fileSize": 204800,
        "isPrimary": true
      }
    ]
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
    title: 'ApolloSupport',
    description: '示例资源',
    content: {
      type: 'doc',
      content: [{
        type: 'paragraph',
        content: [{ type: 'text', text: '资源正文。' }],
      }],
    },
    platform: 'java',
    category: 'plugin',
    downloadType: 'local',
    version: '1.0.0',
    files: [{
      url: uploaded.url,
      fileName: uploaded.filename,
      fileSize: uploaded.size,
      isPrimary: true,
    }],
  }),
})

const resource = await response.json()
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
        'title': 'ApolloSupport',
        'description': '示例资源',
        'content': {
            'type': 'doc',
            'content': [{
                'type': 'paragraph',
                'content': [{'type': 'text', 'text': '资源正文。'}],
            }],
        },
        'platform': 'java',
        'category': 'plugin',
        'downloadType': 'local',
        'version': '1.0.0',
        'files': [{
            'url': uploaded['url'],
            'fileName': uploaded['filename'],
            'fileSize': uploaded['size'],
            'isPrimary': True,
        }],
    },
)
response.raise_for_status()
resource = response.json()
```

旧客户端仍可提交 `fileUrl`、`fileName`、`fileSize`，但新脚本建议统一使用 `files`，这样多个 Loader、多个游戏版本或附加文件更清楚。

## 局部更新资源

PATCH `/api/resources/{id}`

需要权限 `resource:update:self` 请求类型 `application/json`

`PATCH` 适合自动化脚本只同步部分字段，例如教程帖、文档帖、外部文档链接或简介：

```
curl -X PATCH "https://www.nexusmc.cn/api/resources/RESOURCE_ID" \
  -H "Authorization: Bearer avm_pat_xxxxx_xxxxx" \
  -H "Content-Type: application/json" \
  -d '{
    "tutorialPostIds": ["tutorial_post_id"],
    "documentationPostRefs": [
      {
        "id": "getting-started",
        "title": "快速开始",
        "items": [
          { "id": "doc_install_post_id", "type": "install" },
          { "id": "doc_usage_post_id", "type": "usage" }
        ]
      }
    ],
    "documentationUrl": "https://docs.example.com/apollo-support"
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
    tutorialPostIds: ['tutorial_post_id'],
    documentationPostRefs: [{
      id: 'getting-started',
      title: '快速开始',
      items: [
        { id: 'doc_install_post_id', type: 'install' },
        { id: 'doc_usage_post_id', type: 'usage' },
      ],
    }],
    documentationUrl: 'https://docs.example.com/apollo-support',
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
        'tutorialPostIds': ['tutorial_post_id'],
        'documentationPostRefs': [{
            'id': 'getting-started',
            'title': '快速开始',
            'items': [
                {'id': 'doc_install_post_id', 'type': 'install'},
                {'id': 'doc_usage_post_id', 'type': 'usage'},
            ],
        }],
        'documentationUrl': 'https://docs.example.com/apollo-support',
    },
)
response.raise_for_status()
updated = response.json()
```

数组字段通常按整体替换处理。要清空时传空数组；要保留旧值时不要发送该字段。

## 通过更新接口提交新版本

非草稿资源的更新请求满足任一条件时，会创建版本记录：

-   `publishVersion` 为 `true`
-   请求包含 `files`、`fileUrl`、`fileName`、`fileSize`、`fileSha256`、`fileSha1`、`additionalFiles` 或 `extractCode`

新版本需要先上传文件，再把上传结果写入更新请求：

```
curl -X PATCH "https://www.nexusmc.cn/api/resources/RESOURCE_ID" \
  -H "Authorization: Bearer avm_pat_xxxxx_xxxxx" \
  -H "Content-Type: application/json" \
  -d '{
    "publishVersion": true,
    "version": "1.1.0",
    "versionTag": "beta",
    "versionTitle": "兼容新版游戏",
    "changelog": "修复若干问题并更新资源文件。",
    "downloadType": "local",
    "files": [
      {
        "url": "/uploads/files/new-build-random.jar",
        "fileName": "ApolloSupport-1.1.0.jar",
        "fileSize": 307200,
        "isPrimary": true,
        "gameVersions": ["1.21.4"]
      }
    ]
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
    versionTitle: '兼容新版游戏',
    changelog: '修复若干问题并更新资源文件。',
    downloadType: 'local',
    files: [{
      url: '/uploads/files/new-build-random.jar',
      fileName: 'ApolloSupport-1.1.0.jar',
      fileSize: 307200,
      isPrimary: true,
      gameVersions: ['1.21.4'],
    }],
  }),
})

const versionUpdate = await response.json()
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
        'versionTitle': '兼容新版游戏',
        'changelog': '修复若干问题并更新资源文件。',
        'downloadType': 'local',
        'files': [{
            'url': '/uploads/files/new-build-random.jar',
            'fileName': 'ApolloSupport-1.1.0.jar',
            'fileSize': 307200,
            'isPrimary': True,
            'gameVersions': ['1.21.4'],
        }],
    },
)
response.raise_for_status()
version_update = response.json()
```

需要审核的版本只有在审核通过后才会切换资源主文件。

`versionTag` 必须使用 `GET /api/resources/version-tags` 返回的启用 `key`。通过更新接口发版时省略该字段会沿用资源当前 Tag。

## 独立版本接口

POST `/api/resources/{id}/versions`

需要权限 `resource:update:self` 请求类型 `application/json`

如果你的脚本只负责发版，也可以直接调用独立版本接口：

```
curl -X POST "https://www.nexusmc.cn/api/resources/RESOURCE_ID/versions" \
  -H "Authorization: Bearer avm_pat_xxxxx_xxxxx" \
  -H "Content-Type: application/json" \
  -d '{
    "version": "1.1.0",
    "versionTag": "releases",
    "title": "兼容新版游戏",
    "changelog": "修复若干问题并更新资源文件。",
    "downloadType": "local",
    "files": [
      {
        "url": "/uploads/files/new-build-random.jar",
        "fileName": "ApolloSupport-1.1.0.jar",
        "fileSize": 307200,
        "isPrimary": true
      }
    ],
    "mcVersions": ["1.21.4"]
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
    title: '兼容新版游戏',
    changelog: '修复若干问题并更新资源文件。',
    downloadType: 'local',
    files: [{
      url: '/uploads/files/new-build-random.jar',
      fileName: 'ApolloSupport-1.1.0.jar',
      fileSize: 307200,
      isPrimary: true,
    }],
    mcVersions: ['1.21.4'],
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
        'title': '兼容新版游戏',
        'changelog': '修复若干问题并更新资源文件。',
        'downloadType': 'local',
        'files': [{
            'url': '/uploads/files/new-build-random.jar',
            'fileName': 'ApolloSupport-1.1.0.jar',
            'fileSize': 307200,
            'isPrimary': True,
        }],
        'mcVersions': ['1.21.4'],
    },
)
response.raise_for_status()
version = response.json()
```

独立版本接口省略 `versionTag` 时使用 `releases`。后台可能调整可用 Tag，长期运行的发布脚本应在发版前读取 `/api/resources/version-tags`，不要只在代码里固定默认三项。

## 下载文件名

资源文件的下载展示名来自 `fileName`。上传接口返回的 `filename` 是用户原始文件名，脚本应把它映射到 `fileName`。例如上传 `ApolloSupport-1.0.0.jar`，下载时也应展示 `ApolloSupport-1.0.0.jar`，站点内部保存时使用什么随机文件名不影响下载展示名。

如果需要直接生成附件下载链接，可以使用下载代理：

GET `/api/upload/download?url={fileUrl}&filename={downloadName}`

```
curl -L "https://www.nexusmc.cn/api/upload/download?url=%2Fuploads%2Ffiles%2FuuhUuBftBa0P530_VsJ64.jar&filename=ApolloSupport-1.0.0.jar" \
  -o ApolloSupport-1.0.0.jar
```

```
const url = new URL('https://www.nexusmc.cn/api/upload/download')
url.search = new URLSearchParams({
  url: '/uploads/files/uuhUuBftBa0P530_VsJ64.jar',
  filename: 'ApolloSupport-1.0.0.jar',
}).toString()

const response = await fetch(url)
const file = await response.blob()
```

```
import requests

response = requests.get(
    'https://www.nexusmc.cn/api/upload/download',
    params={
        'url': '/uploads/files/uuhUuBftBa0P530_VsJ64.jar',
        'filename': 'ApolloSupport-1.0.0.jar',
    },
)
response.raise_for_status()
with open('ApolloSupport-1.0.0.jar', 'wb') as file:
    file.write(response.content)
```

`url` 必须是上传接口返回的站点托管 URL。`filename` 用于 `Content-Disposition`，建议传资源文件保存的 `fileName`。

## 常见错误

| 问题 | 原因 | 处理 |
| --- | --- | --- |
| 下载链接指向错误位置 | 客户端自行拼了 `/uploads/...`、域名或文件路径 | 始终使用上传响应的 `url`。 |
| 下载文件名变成随机名 | 把内部文件名当成 `fileName` | 使用上传响应里的 `filename`。 |
| 提交资源后没有切换主文件 | 新版本还在审核中 | 等版本审核通过后再确认公开下载。 |
| 上传接口 403 | Token 缺少 `upload:file` | 给 Token 增加上传权限。 |
| 更新接口 403 | Token 不是资源作者或缺少 `resource:update:self` | 使用资源作者本人的 Token，并检查权限。 |
