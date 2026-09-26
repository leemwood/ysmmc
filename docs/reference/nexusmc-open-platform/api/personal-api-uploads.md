---
title: "个人 API：文件上传"
source: https://docs.nexusmc.cn/docs/api/personal-api-uploads
collected: 2026-09-26
---

# 个人 API：文件上传

单文件、多文件、图片、直传和分块上传端点的完整参数、响应与调用示例。

## 接入要求

-   权限：`upload:file`
-   认证：`Authorization: Bearer <token>`
-   普通上传使用 `multipart/form-data`，不要手动设置 boundary。
-   上传成功后，把响应中的站内 URL 交给资源或其他内容接口，不要自行拼接 `/uploads/...` 路径。
-   上传入口统一使用 `/api/upload`。旧测试期使用过的 `/api/transfer` 不再作为公开个人 API 路径。

## 上传响应字段

文件上传成功后通常返回：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `url` | `string` | 服务端最终可访问地址。提交资源文件时映射到 `files[].url`。 |
| `filename` | `string` | 原始文件名。提交资源文件时映射到 `files[].fileName`。 |
| `size` | `number` | 文件字节数。提交资源文件时映射到 `files[].fileSize`。 |
| `sha256` | `string` | SHA-256 校验值。普通上传、多文件上传和分块上传会返回。 |
| `sha1` | `string` | SHA-1 校验值。普通上传、多文件上传和分块上传会返回。 |

直传确认成功后返回 `url`、`filename`、`size`、`key`、`bucketId`，不会返回 `sha256` / `sha1`。如果你的客户端必须保存哈希，请改用普通上传或分块上传。

## 查询当前上传限制

GET `/api/upload/limits`

需要权限 `upload:file`

实际限额由站点配置和账号等级共同决定，会随配置变化。脚本不要写死体积上限，上传前先读一次：

| 字段 | 说明 |
| --- | --- |
| `maxFileSizeMb` | 当前账号单文件上限（MB）。 |
| `maxImageSizeMb` | 单张图片上限（MB）。 |
| `maxTotalSizeMb` | 单次请求总大小上限（MB）。 |
| `chunkThresholdMb` | 超过该体积建议改用分块上传（MB）。 |
| `maxMultiFiles` | 多文件接口一次最多上传几个文件。 |
| `allowedFileTypes` | 允许的文件扩展名。 |
| `allowedImageTypes` | 允许的图片扩展名。 |
| `userLevel` | 当前账号等级。 |
| `levelLimitsEnabled` | 是否启用按等级区分限额。 |
| `siteMaxFileSizeMb` | 站点级单文件上限（MB）。 |
| `levelMaxSizeMb` | 当前等级对应的上限（MB）。 |

该接口限流为每分钟 60 次，不需要每次上传都调用。

## 新手先照这个流程

上传一个资源文件时，推荐按这个顺序做：

1.  先尝试直传：调用 `POST /api/upload/direct/init`。
2.  如果初始化成功，按返回的 `uploadUrl`、`method`、`headers` 上传文件，再调用 `POST /api/upload/direct/complete` 确认。
3.  如果初始化失败、当前环境不支持直传、临时上传地址过期，降级调用 `POST /api/upload` 普通上传。
4.  如果普通上传因为文件太大失败，改用分块上传。

上传完成后，只使用最终响应里的 `url`、`filename`、`size`。不要使用本地文件路径，也不要自己拼接上传目录。

伪代码：

```
try direct/init
  if success:
    upload file to uploadUrl
    result = direct/complete
  else:
    result = POST /api/upload

submit result.url / result.filename / result.size to resource API
```

## 推荐方式：直传

直传适合默认文件上传路径，尤其是资源文件。客户端先向站点申请一个临时上传地址，再按响应要求上传到临时地址，最后回到站点确认上传结果。

如果当前站点配置、网络环境或文件参数不适合直传，本接口会返回错误。客户端不要把这当成致命失败，应自动降级到普通上传或分块上传。

### 1\. 初始化直传

POST `/api/upload/direct/init`

需要权限 `upload:file` 请求类型 `application/json`

| 参数 | 位置 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- | --- |
| `filename` | JSON | string | 是 | 原始文件名。 |
| `size` | JSON | number | 是 | 文件总字节数。 |
| `mimetype` | JSON | string | 否 | 文件 MIME 类型，默认 `application/octet-stream`。 |
| `type` | JSON | string | 否 | 当前仅支持 `files`。 |

```
curl -X POST "https://www.nexusmc.cn/api/upload/direct/init" \
  -H "Authorization: Bearer avm_pat_xxxxx_xxxxx" \
  -H "Content-Type: application/json" \
  -d '{
    "filename": "build.jar",
    "size": 102400,
    "mimetype": "application/java-archive",
    "type": "files"
  }'
```

```
const response = await fetch('https://www.nexusmc.cn/api/upload/direct/init', {
  method: 'POST',
  headers: {
    Authorization: `Bearer ${process.env.NEXUSMC_API_TOKEN}`,
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    filename: 'build.jar',
    size: 102400,
    mimetype: 'application/java-archive',
    type: 'files',
  }),
})

const init = await response.json()
```

```
import os
import requests

response = requests.post(
    'https://www.nexusmc.cn/api/upload/direct/init',
    headers={
        'Authorization': f"Bearer {os.environ['NEXUSMC_API_TOKEN']}",
        'Content-Type': 'application/json',
    },
    json={
        'filename': 'build.jar',
        'size': 102400,
        'mimetype': 'application/java-archive',
        'type': 'files',
    },
)
response.raise_for_status()
init = response.json()
```

响应：

```
{
  "uploadUrl": "https://upload-gateway.example.com/presigned-url",
  "method": "PUT",
  "headers": {
    "Content-Type": "application/java-archive",
    "Content-Disposition": "attachment; filename=\"build.jar\"; filename*=UTF-8''build.jar"
  },
  "expiresIn": 900,
  "key": "uploads/files/random.jar",
  "url": "https://www.nexusmc.cn/uploads/files/random.jar",
  "bucketId": "upload_target_id",
  "filename": "build.jar",
  "size": 102400,
  "mimetype": "application/java-archive",
  "maxFileSizeMb": 100
}
```

保存初始化响应中的 `uploadUrl`、`method`、`headers`、`key`、`bucketId`、`filename`、`size`，下一步要原样使用。

### 2\. 上传到临时地址

使用初始化响应中的 `method`、`uploadUrl` 和 `headers` 原样请求临时上传地址。这个请求不使用个人 API Token。

```
curl -X PUT "$UPLOAD_URL" \
  -H "Content-Type: application/java-archive" \
  -H "Content-Disposition: attachment; filename=\"build.jar\"; filename*=UTF-8''build.jar" \
  --data-binary "@./build.jar"
```

```
import { openAsBlob } from 'node:fs'

const file = await openAsBlob('./build.jar', { type: 'application/java-archive' })
const response = await fetch(init.uploadUrl, {
  method: init.method,
  headers: init.headers,
  body: file,
})

if (!response.ok) {
  throw new Error(await response.text())
}
```

```
import requests

with open('./build.jar', 'rb') as file:
    response = requests.request(
        init['method'],
        init['uploadUrl'],
        headers=init['headers'],
        data=file,
    )

response.raise_for_status()
```

这一步成功后，还不能把初始化响应当作最终上传结果提交给资源接口。必须继续调用确认接口。

### 3\. 确认直传

POST `/api/upload/direct/complete`

需要权限 `upload:file` 请求类型 `application/json`

| 参数 | 位置 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- | --- |
| `bucketId` | JSON | string | 是 | 初始化响应返回的上传确认标识。 |
| `key` | JSON | string | 是 | 初始化响应返回的文件标识。 |
| `filename` | JSON | string | 是 | 原始文件名。 |
| `size` | JSON | number | 是 | 文件总字节数，必须和实际上传文件大小一致。 |

```
curl -X POST "https://www.nexusmc.cn/api/upload/direct/complete" \
  -H "Authorization: Bearer avm_pat_xxxxx_xxxxx" \
  -H "Content-Type: application/json" \
  -d '{
    "bucketId": "upload_target_id",
    "key": "uploads/files/random.jar",
    "filename": "build.jar",
    "size": 102400
  }'
```

```
const response = await fetch('https://www.nexusmc.cn/api/upload/direct/complete', {
  method: 'POST',
  headers: {
    Authorization: `Bearer ${process.env.NEXUSMC_API_TOKEN}`,
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    bucketId: init.bucketId,
    key: init.key,
    filename: init.filename,
    size: init.size,
  }),
})

const uploaded = await response.json()
```

```
import os
import requests

response = requests.post(
    'https://www.nexusmc.cn/api/upload/direct/complete',
    headers={
        'Authorization': f"Bearer {os.environ['NEXUSMC_API_TOKEN']}",
        'Content-Type': 'application/json',
    },
    json={
        'bucketId': init['bucketId'],
        'key': init['key'],
        'filename': init['filename'],
        'size': init['size'],
    },
)
response.raise_for_status()
uploaded = response.json()
```

确认成功后返回可提交给资源接口的文件信息：

```
{
  "url": "https://www.nexusmc.cn/uploads/files/random.jar",
  "filename": "build.jar",
  "size": 102400,
  "key": "uploads/files/random.jar",
  "bucketId": "upload_target_id"
}
```

把这里的 `url`、`filename`、`size` 映射到资源文件字段。直传确认响应才是最终上传结果。

## 降级方式：普通单文件上传

普通上传适合直传初始化失败、客户端暂时不支持直传逻辑、或只想用最简单方式接入的场景。

POST `/api/upload`

需要权限 `upload:file` 请求类型 `multipart/form-data`

| 参数 | 位置 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- | --- |
| `file` | multipart | binary | 是 | 要上传的文件。一次请求只传一个。 |

```
curl -X POST "https://www.nexusmc.cn/api/upload" \
  -H "Authorization: Bearer avm_pat_xxxxx_xxxxx" \
  -H "Accept: application/json" \
  -F "file=@./build.jar"
```

```
import { openAsBlob } from 'node:fs'

const form = new FormData()
form.append('file', await openAsBlob('./build.jar'), 'build.jar')

const response = await fetch('https://www.nexusmc.cn/api/upload', {
  method: 'POST',
  headers: {
    Authorization: `Bearer ${process.env.NEXUSMC_API_TOKEN}`,
    Accept: 'application/json',
  },
  body: form,
})

const uploaded = await response.json()
```

```
import os
import requests

with open('./build.jar', 'rb') as file:
    response = requests.post(
        'https://www.nexusmc.cn/api/upload',
        headers={
            'Authorization': f"Bearer {os.environ['NEXUSMC_API_TOKEN']}",
            'Accept': 'application/json',
        },
        files={'file': ('build.jar', file)},
    )

response.raise_for_status()
uploaded = response.json()
```

响应：

```
{
  "url": "/uploads/files/abc123.jar",
  "filename": "build.jar",
  "size": 102400,
  "sha256": "64位sha256",
  "sha1": "40位sha1"
}
```

## 上传多个文件

多个资源文件最稳妥的方式是逐个执行“优先直传，失败降级普通上传”的流程。这样每个文件都能拿到自己的 `url`、`filename`、`size`，后续提交到资源 `files` 数组。

如果你只需要简单批量上传，也可以使用多文件普通上传：

POST `/api/upload/multiple`

需要权限 `upload:file` 请求类型 `multipart/form-data`

| 参数 | 位置 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- | --- |
| `file` | multipart | binary\[\] | 是 | 重复使用同名字段提交多个文件。 |

```
curl -X POST "https://www.nexusmc.cn/api/upload/multiple" \
  -H "Authorization: Bearer avm_pat_xxxxx_xxxxx" \
  -F "file=@./fabric.jar" \
  -F "file=@./neoforge.jar"
```

```
import { openAsBlob } from 'node:fs'

const form = new FormData()
form.append('file', await openAsBlob('./fabric.jar'), 'fabric.jar')
form.append('file', await openAsBlob('./neoforge.jar'), 'neoforge.jar')

const response = await fetch('https://www.nexusmc.cn/api/upload/multiple', {
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

with open('./fabric.jar', 'rb') as fabric, open('./neoforge.jar', 'rb') as neoforge:
    response = requests.post(
        'https://www.nexusmc.cn/api/upload/multiple',
        headers={'Authorization': f"Bearer {os.environ['NEXUSMC_API_TOKEN']}"},
        files=[
            ('file', ('fabric.jar', fabric)),
            ('file', ('neoforge.jar', neoforge)),
        ],
    )

response.raise_for_status()
uploaded = response.json()
```

响应格式：

```
{
  "files": [
    {
      "url": "/uploads/files/fabric-random.jar",
      "filename": "fabric.jar",
      "size": 102400,
      "sha256": "64位sha256",
      "sha1": "40位sha1"
    }
  ]
}
```

创建资源版本时，建议把其中每一项都转换为资源 `files` 元素，不要丢弃非主文件。

## 上传图片

POST `/api/upload/image`

需要权限 `upload:file` 请求类型 `multipart/form-data`

| 参数 | 位置 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- | --- |
| `file` | multipart | binary | 是 | 图片文件。服务端可能执行格式转换和压缩。 |
| `purpose` | multipart | string | 否 | 图片用途。常用值：`content`、`cover`、`icon`。 |

```
curl -X POST "https://www.nexusmc.cn/api/upload/image" \
  -H "Authorization: Bearer avm_pat_xxxxx_xxxxx" \
  -F "purpose=content" \
  -F "file=@./cover.png"
```

```
import { openAsBlob } from 'node:fs'

const form = new FormData()
form.append('purpose', 'content')
form.append('file', await openAsBlob('./cover.png', { type: 'image/png' }), 'cover.png')

const response = await fetch('https://www.nexusmc.cn/api/upload/image', {
  method: 'POST',
  headers: {
    Authorization: `Bearer ${process.env.NEXUSMC_API_TOKEN}`,
  },
  body: form,
})

const image = await response.json()
```

```
import os
import requests

with open('./cover.png', 'rb') as file:
    response = requests.post(
        'https://www.nexusmc.cn/api/upload/image',
        headers={'Authorization': f"Bearer {os.environ['NEXUSMC_API_TOKEN']}"},
        data={'purpose': 'content'},
        files={'file': ('cover.png', file, 'image/png')},
    )

response.raise_for_status()
image = response.json()
```

典型响应：

```
{
  "url": "/uploads/images/cover_abc123.webp",
  "thumbnailUrl": "/uploads/thumbnails/cover_thumb_abc123.webp",
  "originalSize": 456789,
  "optimizedSize": 123456
}
```

## 查询图片处理状态

GET `/api/upload/image/processing/{jobId}`

需要权限 `upload:file`

上传图片时如果在表单里带 `purpose=content` 和 `processing=background`，响应会改为返回任务而不是最终图片：

```
{
  "processing": true,
  "processingJobId": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

此时缩略图和预览图在后台生成，用返回的 `processingJobId` 轮询本接口：

| 字段 | 说明 |
| --- | --- |
| `jobId` | 任务 ID。 |
| `status` | `processing`、`completed` 或 `failed`。 |
| `result` | `completed` 时返回，包含派生图结果。 |
| `error` | `failed` 时返回，表示已保留原图。 |

`jobId` 只接受 8～64 位的字母、数字和连字符；格式不合法返回 `400`，任务不存在或不属于当前账号返回 `404`。

## 分块上传

大文件按“初始化 -> 上传所有分块 -> 请求完成 -> 查询状态”调用。服务端现在会异步合并分块，所以完成接口可能先返回 `202`，客户端必须支持轮询状态。

推荐使用 `/api/upload/session/...` 这一组路径。旧版 `/api/upload/chunk/...` 仍可使用。

下面示例均使用推荐的 `/api/upload/session/...` 路径。旧客户端如果仍使用 `/api/upload/chunk/...`，只需要把示例中的 URL 替换成对应旧版路径，请求字段保持一致。

### 1\. 初始化

POST `/api/upload/session/init`

需要权限 `upload:file` 请求类型 `application/json`

旧版路径：

POST `/api/upload/chunk/init`

需要权限 `upload:file` 请求类型 `application/json`

| 参数 | 位置 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- | --- |
| `filename` | JSON | string | 是 | 原始文件名。 |
| `size` | JSON | number | 是 | 文件总字节数。 |
| `chunkSize` | JSON | number | 是 | 期望的单块字节数；单块最大 20MB。 |
| `totalChunks` | JSON | integer | 是 | 总分块数，必须等于 `ceil(size / chunkSize)`。 |
| `mimetype` | JSON | string | 否 | 文件 MIME 类型，默认 `application/octet-stream`。 |

```
curl -X POST "https://www.nexusmc.cn/api/upload/session/init" \
  -H "Authorization: Bearer avm_pat_xxxxx_xxxxx" \
  -H "Content-Type: application/json" \
  -d '{
    "filename": "large-build.jar",
    "size": 104857600,
    "chunkSize": 5242880,
    "totalChunks": 20,
    "mimetype": "application/java-archive"
  }'
```

```
const response = await fetch('https://www.nexusmc.cn/api/upload/session/init', {
  method: 'POST',
  headers: {
    Authorization: `Bearer ${process.env.NEXUSMC_API_TOKEN}`,
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    filename: 'large-build.jar',
    size: 104857600,
    chunkSize: 5242880,
    totalChunks: 20,
    mimetype: 'application/java-archive',
  }),
})

const session = await response.json()
```

```
import os
import requests

response = requests.post(
    'https://www.nexusmc.cn/api/upload/session/init',
    headers={
        'Authorization': f"Bearer {os.environ['NEXUSMC_API_TOKEN']}",
        'Content-Type': 'application/json',
    },
    json={
        'filename': 'large-build.jar',
        'size': 104857600,
        'chunkSize': 5242880,
        'totalChunks': 20,
        'mimetype': 'application/java-archive',
    },
)
response.raise_for_status()
session = response.json()
```

响应包含 `uploadId`、实际 `chunkSize`、`totalChunks`、`maxFileSizeMb` 和 `chunkThresholdMb`。

### 2\. 上传单块

POST `/api/upload/session/{uploadId}/part/{index}`

需要权限 `upload:file` 请求类型 `multipart/form-data`

旧版路径：

POST `/api/upload/chunk/{uploadId}/{index}`

需要权限 `upload:file` 请求类型 `multipart/form-data`

| 参数 | 位置 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- | --- |
| `uploadId` | path | string | 是 | 初始化响应返回的上传 ID。 |
| `index` | path | integer | 是 | 从 `0` 开始的分块序号。 |
| `chunk` | multipart | binary | 是 | 当前分块内容。大小必须与初始化结果匹配；旧客户端传 `file` 也会被接收。 |

```
curl -X POST "https://www.nexusmc.cn/api/upload/session/UPLOAD_ID/part/0" \
  -H "Authorization: Bearer avm_pat_xxxxx_xxxxx" \
  -F "chunk=@./part-0.bin"
```

```
import { openAsBlob } from 'node:fs'

const form = new FormData()
form.append('chunk', await openAsBlob('./part-0.bin'), 'part-0.bin')

const response = await fetch('https://www.nexusmc.cn/api/upload/session/UPLOAD_ID/part/0', {
  method: 'POST',
  headers: {
    Authorization: `Bearer ${process.env.NEXUSMC_API_TOKEN}`,
  },
  body: form,
})

const part = await response.json()
```

```
import os
import requests

with open('./part-0.bin', 'rb') as chunk:
    response = requests.post(
        'https://www.nexusmc.cn/api/upload/session/UPLOAD_ID/part/0',
        headers={'Authorization': f"Bearer {os.environ['NEXUSMC_API_TOKEN']}"},
        files={'chunk': ('part-0.bin', chunk)},
    )

response.raise_for_status()
part = response.json()
```

每个分块成功后返回该块的 `index` 和 `bytes`。分块可以重试，但完成合并前必须全部成功。

### 3\. 请求完成上传

POST `/api/upload/session/{uploadId}/finish`

需要权限 `upload:file`

旧版路径：

POST `/api/upload/chunk/{uploadId}/complete`

需要权限 `upload:file`

```
curl -X POST "https://www.nexusmc.cn/api/upload/session/UPLOAD_ID/finish" \
  -H "Authorization: Bearer avm_pat_xxxxx_xxxxx" \
  -H "Accept: application/json"
```

```
const response = await fetch('https://www.nexusmc.cn/api/upload/session/UPLOAD_ID/finish', {
  method: 'POST',
  headers: {
    Authorization: `Bearer ${process.env.NEXUSMC_API_TOKEN}`,
    Accept: 'application/json',
  },
})

const finish = await response.json()
```

```
import os
import requests

response = requests.post(
    'https://www.nexusmc.cn/api/upload/session/UPLOAD_ID/finish',
    headers={
        'Authorization': f"Bearer {os.environ['NEXUSMC_API_TOKEN']}",
        'Accept': 'application/json',
    },
)
response.raise_for_status()
finish = response.json()
```

完成接口有两种成功响应：

```
{
  "processing": true,
  "uploadId": "upload_id"
}
```

或在结果仍处于短暂缓存期时直接返回最终文件信息：

```
{
  "url": "/uploads/files/random.jar",
  "filename": "build.jar",
  "size": 102400,
  "sha256": "64位sha256",
  "sha1": "40位sha1"
}
```

如果收到 `processing: true`，继续调用状态接口，不要把该响应当成上传结果提交给资源接口。

### 4\. 查询状态

GET `/api/upload/session/{uploadId}/status`

需要权限 `upload:file`

旧版路径：

GET `/api/upload/chunk/{uploadId}/status`

需要权限 `upload:file`

```
curl "https://www.nexusmc.cn/api/upload/session/UPLOAD_ID/status" \
  -H "Authorization: Bearer avm_pat_xxxxx_xxxxx" \
  -H "Accept: application/json"
```

```
const response = await fetch('https://www.nexusmc.cn/api/upload/session/UPLOAD_ID/status', {
  headers: {
    Authorization: `Bearer ${process.env.NEXUSMC_API_TOKEN}`,
    Accept: 'application/json',
  },
})

const status = await response.json()
```

```
import os
import requests

response = requests.get(
    'https://www.nexusmc.cn/api/upload/session/UPLOAD_ID/status',
    headers={
        'Authorization': f"Bearer {os.environ['NEXUSMC_API_TOKEN']}",
        'Accept': 'application/json',
    },
)
response.raise_for_status()
status = response.json()
```

可能返回：

```
{ "status": "pending" }
```

```
{
  "status": "processing",
  "startedAt": "2026-07-02T05:00:00.000Z",
  "updatedAt": "2026-07-02T05:00:00.000Z"
}
```

```
{
  "status": "completed",
  "result": {
    "url": "/uploads/files/random.jar",
    "filename": "build.jar",
    "size": 102400,
    "sha256": "64位sha256",
    "sha1": "40位sha1"
  },
  "startedAt": "2026-07-02T05:00:00.000Z",
  "updatedAt": "2026-07-02T05:00:03.000Z"
}
```

```
{
  "status": "failed",
  "error": "合并分块上传失败",
  "startedAt": "2026-07-02T05:00:00.000Z",
  "updatedAt": "2026-07-02T05:00:03.000Z"
}
```

状态为 `completed` 时，使用 `result` 里的文件信息。状态为 `failed` 时，展示 `error` 并让用户重新上传。上传会话不存在或已过期时会返回 `410`。

### 5\. 取消上传

DELETE `/api/upload/session/{uploadId}`

需要权限 `upload:file`

旧版路径：

DELETE `/api/upload/chunk/{uploadId}`

需要权限 `upload:file`

```
curl -X DELETE "https://www.nexusmc.cn/api/upload/session/UPLOAD_ID" \
  -H "Authorization: Bearer avm_pat_xxxxx_xxxxx" \
  -H "Accept: application/json"
```

```
const response = await fetch('https://www.nexusmc.cn/api/upload/session/UPLOAD_ID', {
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
    'https://www.nexusmc.cn/api/upload/session/UPLOAD_ID',
    headers={
        'Authorization': f"Bearer {os.environ['NEXUSMC_API_TOKEN']}",
        'Accept': 'application/json',
    },
)
response.raise_for_status()
result = response.json()
```

## 下载代理

GET `/api/upload/download?url={fileUrl}&filename={downloadName}`

该接口用于把站点托管的上传文件作为附件下载，并设置正确的 `Content-Disposition` 文件名。它不需要 `upload:file` 权限，但只接受上传接口返回的站内上传 URL。

| 查询参数 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `url` | string | 是 | 上传接口返回的 `url`，可以是 `/uploads/...` 相对路径或站点返回的完整 URL。 |
| `filename` | string | 否 | 下载时展示的文件名，默认 `download`。 |

```
curl -L "https://www.nexusmc.cn/api/upload/download?url=%2Fuploads%2Ffiles%2Frandom.jar&filename=build.jar" \
  -o build.jar
```

```
const url = new URL('https://www.nexusmc.cn/api/upload/download')
url.search = new URLSearchParams({
  url: '/uploads/files/random.jar',
  filename: 'build.jar',
}).toString()

const response = await fetch(url)
const file = await response.blob()
```

```
import requests

response = requests.get(
    'https://www.nexusmc.cn/api/upload/download',
    params={'url': '/uploads/files/random.jar', 'filename': 'build.jar'},
)
response.raise_for_status()
with open('build.jar', 'wb') as file:
    file.write(response.content)
```

## 上传后怎么使用

| 上传响应字段 | 资源 `files` 字段 |
| --- | --- |
| `url` | `url` |
| `filename` | `fileName` |
| `size` | `fileSize` |
| `sha256` | `sha256` |
| `sha1` | `sha1` |

完整文件结构、主文件和 Loader 规则见[个人 API：资源文件与版本](/docs/api/personal-api-resource-files/)。
