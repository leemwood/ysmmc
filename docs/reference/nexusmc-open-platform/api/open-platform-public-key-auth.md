---
title: "开放平台：公钥鉴权"
source: https://docs.nexusmc.cn/docs/api/open-platform-public-key-auth
collected: 2026-09-26
---

# 开放平台：公钥鉴权

使用 P-256 公钥为 OAuth private\_key\_jwt 和站点 API HTTP 请求签名。

公钥鉴权适合部署在可信服务端的 OAuth 应用、机器人和数据同步程序。私钥只保存在调用方，NexusMC 开放中心只接收公钥。即使数据库中的公开凭据标识泄露，攻击者也不能据此伪造请求。

当前支持两种用途：

| 用途 | 开放中心标识 | 使用位置 |
| --- | --- | --- |
| OAuth 客户端认证 | `oauth_client_auth` | `POST /api/oauth2/token` 的 `private_key_jwt` |
| 站点 API 请求签名 | `site_api_request` | `/api/site/v1/*` 的 HTTP Message Signature |

第一版只接受 `P-256` EC 公钥。每个 OAuth 应用最多保留 10 把未撤销公钥，可以在轮换期间同时启用新旧两把。

## 生成密钥

```
openssl ecparam -name prime256v1 -genkey -noout -out nexusmc-private.pem
openssl ec -in nexusmc-private.pem -pubout -out nexusmc-public.pem
```

将 `nexusmc-public.pem` 上传到开放中心，私钥文件必须留在自己的密钥管理系统或服务端。不要上传以 `BEGIN EC PRIVATE KEY` 或 `BEGIN PRIVATE KEY` 开头的内容。

上传后记录 NexusMC 返回的 `kid`。同一把公钥无论以 SPKI PEM 还是公开 JWK 上传，都会得到相同的 RFC 7638 SHA-256 指纹。

## OAuth private\_key\_jwt

在开放中心为应用添加用途为“OAuth 客户端认证”的公钥，再把应用认证方式切换为“公钥签名”。切换后，`/token` 不再接受该应用原来的 `client_secret`。

客户端断言必须满足：

-   JWT header：`alg=ES256`、`typ=JWT`、`kid=<开放中心密钥编号>`
-   `iss` 与 `sub`：都等于应用 `client_id`
-   `aud`：逐字等于 `https://www.nexusmc.cn/api/oauth2/token`
-   `iat`：当前时间，服务端允许少量时钟偏差
-   `exp`：晚于 `iat` 且有效期不超过 120 秒
-   `jti`：每次随机生成，UTF-8 长度不超过 128 字节，不能重复使用

Node.js 示例：

```
import { createPrivateKey, randomUUID, sign } from 'node:crypto'
import { readFile } from 'node:fs/promises'

const clientId = process.env.NEXUSMC_CLIENT_ID
const keyId = process.env.NEXUSMC_KEY_ID
const privateKey = createPrivateKey(await readFile('./nexusmc-private.pem'))
const now = Math.floor(Date.now() / 1000)
const encode = (value) => Buffer.from(JSON.stringify(value)).toString('base64url')
const header = encode({ alg: 'ES256', typ: 'JWT', kid: keyId })
const payload = encode({
  iss: clientId,
  sub: clientId,
  aud: 'https://www.nexusmc.cn/api/oauth2/token',
  iat: now,
  exp: now + 90,
  jti: randomUUID(),
})
const signingInput = `${header}.${payload}`
const signature = sign('sha256', Buffer.from(signingInput), {
  key: privateKey,
  dsaEncoding: 'ieee-p1363',
}).toString('base64url')
const assertion = `${signingInput}.${signature}`

const body = new URLSearchParams({
  grant_type: 'authorization_code',
  client_id: clientId,
  client_assertion_type: 'urn:ietf:params:oauth:client-assertion-type:jwt-bearer',
  client_assertion: assertion,
  code: process.env.NEXUSMC_AUTHORIZATION_CODE,
  redirect_uri: 'https://example.com/oauth/callback',
})

const response = await fetch('https://www.nexusmc.cn/api/oauth2/token', {
  method: 'POST',
  headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
  body,
})
```

刷新 Token 时使用同样的客户端断言字段，只把业务参数改为 `grant_type=refresh_token` 和 `refresh_token=...`。每次请求都必须创建新的 `jti` 和签名。

## 站点 API 请求签名

创建站点 API 凭据时选择“公钥请求签名”，开放中心会返回公开的 `credentialId`，不会生成静态 API Key。签名使用 RFC 9421 的固定 NexusMC profile。

所有请求必须覆盖以下组件，顺序不能改变：

```
"@method" "@authority" "@path" "@query" "x-nexusmc-credential" "x-nexusmc-date" "x-nexusmc-nonce"
```

有请求体时还必须在末尾覆盖：

```
"content-type" "content-digest"
```

其他约束：

-   `@authority` 固定使用公开 API 域名 `www.nexusmc.cn`，不要使用反向代理内网 Host
-   `created` 与当前时间允许少量偏差，`expires-created` 不超过 120 秒
-   `keyid` 使用开放中心公钥编号，`alg` 固定为 `ecdsa-p256-sha256`
-   `x-nexusmc-date` 使用 ISO 8601 UTC 时间
-   `x-nexusmc-nonce` 使用至少 96 bit 随机值，每次请求都必须不同
-   请求体摘要格式为 `sha-256=:<标准 Base64>:`

以下 Node.js 示例签名一个 GET 请求：

```
import { createPrivateKey, randomBytes, sign } from 'node:crypto'
import { readFile } from 'node:fs/promises'

const url = new URL('https://www.nexusmc.cn/api/site/v1/search?q=stone')
const credentialId = process.env.NEXUSMC_CREDENTIAL_ID
const keyId = process.env.NEXUSMC_KEY_ID
const privateKey = createPrivateKey(await readFile('./nexusmc-private.pem'))
const created = Math.floor(Date.now() / 1000)
const expires = created + 90
const nonce = randomBytes(16).toString('base64url')
const date = new Date(created * 1000).toISOString()
const covered = ['@method', '@authority', '@path', '@query', 'x-nexusmc-credential', 'x-nexusmc-date', 'x-nexusmc-nonce']
const params = `(${covered.map((item) => `"${item}"`).join(' ')})`+
  `;created=${created};expires=${expires};keyid="${keyId}";alg="ecdsa-p256-sha256"`
const base = [
  `"@method": GET`,
  `"@authority": ${url.host}`,
  `"@path": ${url.pathname}`,
  `"@query": ${url.search || '?'}`,
  `"x-nexusmc-credential": ${credentialId}`,
  `"x-nexusmc-date": ${date}`,
  `"x-nexusmc-nonce": ${nonce}`,
  `"@signature-params": ${params}`,
].join('\n')
const signature = sign('sha256', Buffer.from(base), {
  key: privateKey,
  dsaEncoding: 'ieee-p1363',
}).toString('base64')

const response = await fetch(url, {
  headers: {
    'X-NexusMC-Credential': credentialId,
    'X-NexusMC-Date': date,
    'X-NexusMC-Nonce': nonce,
    'Signature-Input': `sig1=${params}`,
    Signature: `sig1=:${signature}:`,
  },
})
```

对于 JSON 请求体，先对实际发送的 UTF-8 字节计算 SHA-256，添加 `Content-Type` 和 `Content-Digest`，再把这两个小写组件追加到 covered 列表和签名基串中。序列化完成后不能再次格式化或修改请求体。

## 轮换与故障处理

1.  先上传新公钥，并勾选与旧公钥相同的用途。
2.  部署新私钥，让新请求开始使用新的 `kid`。
3.  观察新密钥的最近使用时间。
4.  确认所有实例切换后，撤销旧公钥。

`401 SITE_API_KEY_INVALID` 不区分凭据不存在、签名错误、重放、过期或密钥已撤销，避免向外暴露认证细节。调用方不应自动重试同一签名；应重新生成时间、nonce 和签名后再按业务重试策略发起请求。

相关标准：[OAuth JWT 客户端认证（RFC 7523）](https://www.rfc-editor.org/rfc/rfc7523)、[HTTP Message Signatures（RFC 9421）](https://www.rfc-editor.org/rfc/rfc9421)、[JWK Thumbprint（RFC 7638）](https://www.rfc-editor.org/rfc/rfc7638)。
