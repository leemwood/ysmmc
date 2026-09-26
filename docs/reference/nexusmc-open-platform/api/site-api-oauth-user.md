---
title: "站点 API：OAuth 用户身份"
source: https://docs.nexusmc.cn/docs/api/site-api-oauth-user
collected: 2026-09-26
---

# 站点 API：OAuth 用户身份

使用已绑定 OAuth 应用的站点 API 密钥读取调用用户身份。

站点 API 默认只代表应用本身。需要知道“当前登录用户是谁”时，调用方必须同时提供站点 API 密钥和 OAuth Access Token。两者必须绑定到同一个 OAuth 应用，服务端不会接受不匹配的凭据。

## 双凭据请求

```
X-API-Key: avm_sak_xxxxx_xxxxx
Authorization: Bearer avm_oat_xxxxx_xxxxx
Accept: application/json
```

站点 API 密钥负责应用识别、业务 scope 和速率限制；OAuth Access Token 负责用户身份和 OAuth scope。Access Token 由[本站 OAuth2 授权码流程](/docs/api/site-oauth2-provider/)签发。

## 获取当前用户

GET `/api/site/v1/me`

需要权限 `user:basic` 认证方式 `X-API-Key + OAuth Bearer`

```
curl "https://www.nexusmc.cn/api/site/v1/me" \
  -H "X-API-Key: $NEXUSMC_SITE_API_KEY" \
  -H "Authorization: Bearer $NEXUSMC_ACCESS_TOKEN" \
  -H "Accept: application/json"
```

响应示例：

```
{
  "schemaVersion": "1",
  "user": {
    "id": "user_uuid",
    "uid": 1024,
    "username": "Steve",
    "slug": "steve",
    "avatar": "/uploads/images/avatar.png",
    "role": "user",
    "status": "active",
    "contribution": 128,
    "emailVerified": true
  },
  "oauth": {
    "clientId": "oauth_client_row_id",
    "scopes": ["user:basic"]
  }
}
```

当前接口只返回基础身份字段，不返回邮箱和敏感资料。投稿、互动和通知接口使用独立 OAuth scope，参见[OAuth 用户数据](/docs/api/site-api-oauth-user-data/)。

## 错误处理

| 状态码 | `code` | 含义 |
| --- | --- | --- |
| `401` | `SITE_API_KEY_INVALID` | 站点 API 密钥无效、停用或过期。 |
| `401` | `OAUTH_ACCESS_TOKEN_INVALID` | OAuth Access Token 无效、撤销或过期。 |
| `403` | `SITE_API_OAUTH_BINDING_REQUIRED` | 站点 API 密钥尚未绑定 OAuth 应用。 |
| `403` | `SITE_API_OAUTH_CLIENT_MISMATCH` | API 密钥和 Access Token 属于不同 OAuth 应用。 |
| `403` | `OAUTH_SCOPE_REQUIRED` | Access Token 未包含 `user:basic`。 |

不要把 Access Token 或站点 API 密钥放入浏览器代码、日志或公开仓库。
