---
title: "本站 OAuth2 接入"
source: https://docs.nexusmc.cn/docs/api/site-oauth2-provider
collected: 2026-09-26
---

# 本站 OAuth2 接入

说明第三方站点如何把 NexusMC 作为 OAuth2 身份提供方接入，包括应用创建、授权码流程、接口地址、作用域和常见失败原因。

## 适用范围

这篇文档面向需要把 NexusMC 账号作为登录方式接入自己站点或工具的开发者。

当前开放的是标准 OAuth2 授权码流程，适合：

-   第三方网站登录
-   独立工具或桌面程序的网页登录回调
-   需要读取 NexusMC 用户基础资料或邮箱资料的服务

服务端应用可以使用传统 `client_secret`，也可以改用更适合生产环境的 [`private_key_jwt` 公钥认证](/docs/api/open-platform-public-key-auth/)。

NOTE

接入前需要先在开放平台创建并通过审核的“本站 OAuth2 应用”。所有应用都会获得 `client_id`；使用 `client_secret_basic` 的应用另有 `client_secret`，使用 `private_key_jwt` 的应用则需要先登记公钥。

## Base URL

生产环境：

```
https://www.nexusmc.cn/api/oauth2
```

## 你需要先准备什么

在第三方站点接入前，需要先让 NexusMC 管理员在后台“后台管理 -> 站点 API 管理 -> OAuth 应用”里创建应用，并提供下面几项信息：

| 字段 | 说明 |
| --- | --- |
| `client_id` | 应用公开标识，例如 `0123abcd...` |
| `client_secret` | 完整应用密钥，例如 `avm_oac_<client_id>_<secret>`，只会完整展示一次 |
| `redirect_uri` | 你的回调地址，必须与后台登记值完全一致 |
| `scope` | 空格分隔的授权范围，完整列表见下方 |

IMPORTANT

`client_secret` 请使用后台创建应用时显示的完整值，例如 `avm_oac_<client_id>_<secret>`。不要拆分、截断或自行重新拼接。服务端当前兼容裸 secret 是为了照顾历史数据，新接入不要依赖这个兼容行为。

WARNING

`redirect_uri` 必须逐字匹配。只要协议、域名、路径、查询串有一处不同，授权请求就会被拒绝。

## 当前可申请的作用域

| 作用域 | 含义 |
| --- | --- |
| `user:basic` | 读取用户 ID、UID、用户名、slug、头像、角色 |
| `user:email` | 读取邮箱地址与邮箱验证状态 |
| `user:content:read` | 读取用户本人的四类投稿摘要 |
| `user:interaction:read` | 读取用户本人的收藏和点赞记录 |
| `user:notification:read` | 读取用户本人的通知和未读数量 |
| `user:notification:write` | 将用户本人的通知标记为已读 |

如果请求时不传 `scope`，服务端会按该应用后台已勾选的默认作用域发放授权。

## 授权码流程

完整流程如下：

1.  第三方站点把用户重定向到 NexusMC 的 `/authorize`
2.  用户在 NexusMC 登录
3.  用户在授权确认页同意或拒绝
4.  NexusMC 把用户带回你的 `redirect_uri`，并附带 `code`
5.  你的服务端调用 `/token` 换取 `access_token`
6.  你的服务端再调用 `/userinfo` 读取用户资料

IMPORTANT

`client_secret` 只能放在你自己的服务端使用，不要放到浏览器端、移动端包体或公开仓库。

## 发起授权

GET `/api/oauth2/authorize`

请求参数：

| 参数 | 必填 | 说明 |
| --- | --- | --- |
| `response_type` | 是 | 固定传 `code` |
| `client_id` | 是 | 管理员创建应用后提供 |
| `redirect_uri` | 是 | 必须与后台登记值完全一致 |
| `scope` | 否 | 空格分隔的作用域列表 |
| `state` | 否 | 建议传，用于防止回调串改 |

请求示例：

```
https://www.nexusmc.cn/api/oauth2/authorize?response_type=code&client_id=your_client_id&redirect_uri=https%3A%2F%2Fexample.com%2Fauth%2Fnexusmc%2Fcallback&scope=user%3Abasic%20user%3Aemail&state=csrf_token_here
```

授权成功后，NexusMC 会跳回：

```
https://example.com/auth/nexusmc/callback?code=AUTH_CODE&state=csrf_token_here
```

如果用户拒绝授权，则会跳回：

```
https://example.com/auth/nexusmc/callback?error=access_denied&state=csrf_token_here
```

## 换取 Access Token

POST `/api/oauth2/token`

请求类型 `application/x-www-form-urlencoded`

请求参数：

| 参数 | 必填 | 说明 |
| --- | --- | --- |
| `grant_type` | 是 | 首次换 Token 时传 `authorization_code` |
| `client_id` | 是 | 应用 `client_id` |
| `client_secret` | 按鉴权方式 | `client_secret_basic` 应用传后台显示的完整应用密钥；`private_key_jwt` 应用不传 |
| `client_assertion_type` | 按鉴权方式 | `private_key_jwt` 应用固定传 `urn:ietf:params:oauth:client-assertion-type:jwt-bearer` |
| `client_assertion` | 按鉴权方式 | `private_key_jwt` 应用使用已登记公钥对应私钥签发的短时 JWT |
| `code` | 是 | 授权回调带回的授权码 |
| `redirect_uri` | 是 | 必须与发起授权时一致 |

```
curl -X POST "https://www.nexusmc.cn/api/oauth2/token" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=authorization_code" \
  -d "client_id=your_client_id" \
  -d "client_secret=avm_oac_your_client_id_your_secret" \
  -d "code=AUTH_CODE" \
  -d "redirect_uri=https://example.com/auth/nexusmc/callback"
```

```
const body = new URLSearchParams({
  grant_type: 'authorization_code',
  client_id: process.env.NEXUSMC_CLIENT_ID,
  // 填后台显示的完整 client_secret，例如 avm_oac_<client_id>_<secret>
  client_secret: process.env.NEXUSMC_CLIENT_SECRET,
  code,
  redirect_uri: 'https://example.com/auth/nexusmc/callback',
})

const response = await fetch('https://www.nexusmc.cn/api/oauth2/token', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/x-www-form-urlencoded',
  },
  body,
})

const token = await response.json()
```

```
import requests

response = requests.post(
    'https://www.nexusmc.cn/api/oauth2/token',
    headers={'Content-Type': 'application/x-www-form-urlencoded'},
    data={
        'grant_type': 'authorization_code',
        'client_id': 'your_client_id',
        'client_secret': 'avm_oac_your_client_id_your_secret',
        'code': 'AUTH_CODE',
        'redirect_uri': 'https://example.com/auth/nexusmc/callback',
    },
)

response.raise_for_status()
token = response.json()
```

响应示例：

```
{
  "access_token": "avm_oat_xxxxx_xxxxx",
  "refresh_token": "avm_ort_xxxxx_xxxxx",
  "token_type": "Bearer",
  "expires_in": 7200,
  "scope": "user:basic user:email"
}
```

## 刷新 Token

POST `/api/oauth2/token`

请求类型 `application/x-www-form-urlencoded`

刷新时把 `grant_type` 改为 `refresh_token`：

| 参数 | 必填 | 说明 |
| --- | --- | --- |
| `grant_type` | 是 | 固定传 `refresh_token` |
| `client_id` | 是 | 应用 `client_id` |
| `client_secret` | 按鉴权方式 | `client_secret_basic` 应用传后台显示的完整应用密钥；`private_key_jwt` 应用不传 |
| `client_assertion_type` | 按鉴权方式 | `private_key_jwt` 应用固定传 `urn:ietf:params:oauth:client-assertion-type:jwt-bearer` |
| `client_assertion` | 按鉴权方式 | `private_key_jwt` 应用使用已登记公钥对应私钥签发的短时 JWT |
| `refresh_token` | 是 | 上次换到的刷新令牌 |

```
curl -X POST "https://www.nexusmc.cn/api/oauth2/token" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=refresh_token" \
  -d "client_id=your_client_id" \
  -d "client_secret=avm_oac_your_client_id_your_secret" \
  -d "refresh_token=avm_ort_xxxxx_xxxxx"
```

```
const refreshBody = new URLSearchParams({
  grant_type: 'refresh_token',
  client_id: 'your_client_id',
  client_secret: 'avm_oac_your_client_id_your_secret',
  refresh_token: 'avm_ort_xxxxx_xxxxx',
})

const refreshResponse = await fetch('https://www.nexusmc.cn/api/oauth2/token', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/x-www-form-urlencoded',
  },
  body: refreshBody,
})

if (!refreshResponse.ok) {
  throw new Error(await refreshResponse.text())
}

const refreshed = await refreshResponse.json()
```

```
import requests

refresh_response = requests.post(
    'https://www.nexusmc.cn/api/oauth2/token',
    headers={'Content-Type': 'application/x-www-form-urlencoded'},
    data={
        'grant_type': 'refresh_token',
        'client_id': 'your_client_id',
        'client_secret': 'avm_oac_your_client_id_your_secret',
        'refresh_token': 'avm_ort_xxxxx_xxxxx',
    },
)

refresh_response.raise_for_status()
refreshed = refresh_response.json()
```

刷新成功后会签发新的 `access_token` 和 `refresh_token`，旧 `refresh_token` 会失效。请保存最新返回值。

## 主动撤销 Token

客户端主动注销或删除本地授权时，可以按 OAuth 2.0 Token Revocation 规范调用撤销接口。服务端会对不存在、已过期或已撤销的 Token 返回相同的成功响应，避免泄露 Token 状态。

POST `/api/oauth2/revoke`

请求类型 `application/x-www-form-urlencoded`

```
curl -X POST "https://www.nexusmc.cn/api/oauth2/revoke" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  --data-urlencode "token=$NEXUSMC_ACCESS_TOKEN" \
  --data-urlencode "token_type_hint=access_token"
```

`token_type_hint` 可传 `access_token` 或 `refresh_token`，也可以省略。撤销 Access Token 时，对应的刷新令牌也会一并失效。

## 获取用户资料

GET `/api/oauth2/userinfo`

认证方式 `Bearer Access Token`

```
curl "https://www.nexusmc.cn/api/oauth2/userinfo" \
  -H "Authorization: Bearer avm_oat_xxxxx_xxxxx" \
  -H "Accept: application/json"
```

```
const response = await fetch('https://www.nexusmc.cn/api/oauth2/userinfo', {
  headers: {
    Authorization: `Bearer ${accessToken}`,
    Accept: 'application/json',
  },
})

const profile = await response.json()
```

```
import requests

response = requests.get(
    'https://www.nexusmc.cn/api/oauth2/userinfo',
    headers={
        'Authorization': f'Bearer {accessToken}',
        'Accept': 'application/json',
    },
)

response.raise_for_status()
profile = response.json()
```

当授权范围为 `user:basic user:email` 时，响应大致如下：

```
{
  "sub": "user_uuid",
  "uid": 1024,
  "username": "Steve",
  "slug": "steve",
  "avatar": "/uploads/images/avatar.png",
  "role": "user",
  "preferred_username": "Steve",
  "email": "steve@example.com",
  "email_verified": true
}
```

如果只有 `user:basic`，则不会返回 `email` 和 `email_verified`。

## 用户撤销授权后的表现

用户可以在 NexusMC 的“设置 -> 账号绑定 -> 已授权第三方站点”里主动撤销授权。

撤销后会发生两件事：

-   旧的 `access_token` 会失效
-   旧的 `refresh_token` 也会失效

这时你的服务端应当把用户带回重新发起授权，而不是无限重试刷新。

## 常见失败原因

| 场景 | 返回 | 原因 | 建议 |
| --- | --- | --- | --- |
| 授权地址参数不对 | `invalid_request` | 缺少 `client_id`、`redirect_uri` 或 `response_type` | 检查拼接参数 |
| 应用不可用 | `unauthorized_client` | 应用不存在、被停用，或密钥错误 | 让管理员检查后台配置 |
| 回调地址不匹配 | `invalid_redirect_uri` | `redirect_uri` 与后台登记值不一致 | 改成完全一致的地址 |
| 作用域不允许 | `invalid_scope` | 请求了应用未开放的作用域 | 调整请求 scope 或后台配置 |
| 用户拒绝授权 | `access_denied` | 用户在确认页点了拒绝 | 提示用户重新发起登录 |
| 换 Token 失败 | `invalid_grant` | 授权码已用过、已过期、回调地址不一致，或授权已撤销 | 重新走一遍授权流程 |
| 客户端认证失败 | `invalid_client` | `client_id`、`client_secret` 或客户端断言不正确 | 检查当前鉴权方式、密钥、公钥及断言声明 |
| Token 无效 | `invalid_token` | Access Token 已过期、被撤销或用户状态不可用 | 重新登录并重新授权 |

## 疑难解答

### `state` 一定要传吗？

强烈建议传。服务端不会替你生成第三方站点自己的 CSRF 校验值，所以应当由你的服务端生成并校验。

### 可以在浏览器前端直接调用 `/token` 吗？

不建议。`client_secret` 和 `private_key_jwt` 的私钥都只能保存在你自己的服务端，不能交给浏览器前端。

### `client_secret` 应该填哪一段？

填后台创建应用后显示的完整字符串，例如：

```
avm_oac_<client_id>_<secret>
```

不要只填最后一段 secret，也不要按 `client_id` 自己拼一个新值。完整密钥只显示一次，丢失后需要重新生成或重建应用。

### `userinfo` 返回的 `sub` 是什么？

`sub` 是 NexusMC 用户的内部唯一 ID，适合作为你自己系统里绑定 NexusMC 账号的稳定主键。

### 用户名会变吗？

`username` 理论上可能变化，所以如果你需要稳定关联，请优先保存 `sub`，展示时再使用最新的 `username`。
