---
title: "站点 API：OAuth 用户数据"
source: https://docs.nexusmc.cn/docs/api/site-api-oauth-user-data
collected: 2026-09-26
---

# 站点 API：OAuth 用户数据

读取 OAuth 用户本人的投稿、收藏、点赞和通知，并将通知标记为已读。

这些接口需要同时携带绑定到同一个 OAuth 应用的站点 API Key 和 OAuth Access Token：

```
X-API-Key: avm_sak_xxxxx_xxxxx
Authorization: Bearer avm_oat_xxxxx_xxxxx
Accept: application/json
```

## 权限范围

| OAuth scope | 能力 |
| --- | --- |
| `user:content:read` | 读取用户本人的资源、帖子、视频和找服玩投稿摘要。 |
| `user:interaction:read` | 读取用户本人对公开内容的收藏和点赞。 |
| `user:notification:read` | 读取用户本人的通知和未读数量。 |
| `user:notification:write` | 将用户本人的通知标记为已读。 |

应用必须在授权请求中显式申请所需 scope，用户也会在授权确认页看到每项用途。

## 通用分页

列表接口接受 `page` 和 `pageSize`，默认分别为 `1` 和 `20`，`pageSize` 最大为 `50`。响应统一包含：

```
{
  "schemaVersion": "1",
  "items": [],
  "pagination": { "page": 1, "pageSize": 20, "total": 0, "totalPages": 0 },
  "generatedAt": "2026-09-05T00:00:00.000Z"
}
```

## 我的投稿

GET `/api/site/v1/me/content`

需要权限 `user:content:read` 认证方式 `X-API-Key + OAuth Bearer`

| 参数 | 必填 | 说明 |
| --- | --- | --- |
| `type` | 否 | `resource`、`post`、`video` 或 `player`，默认 `resource`。 |
| `status` | 否 | `all`、`approved`、`pending` 或 `rejected`，默认 `all`。 |
| `page` | 否 | 页码。 |
| `pageSize` | 否 | 每页数量，最大 50。 |

```
curl "https://www.nexusmc.cn/api/site/v1/me/content?type=post&status=pending" \
  -H "X-API-Key: $NEXUSMC_SITE_API_KEY" \
  -H "Authorization: Bearer $NEXUSMC_ACCESS_TOKEN"
```

响应只包含投稿摘要、审核状态、统计和时间，不返回正文、访问密码、审核内部备注或其他敏感字段。

## 我的收藏与点赞

GET `/api/site/v1/me/favorites`

需要权限 `user:interaction:read` 认证方式 `X-API-Key + OAuth Bearer`

GET `/api/site/v1/me/likes`

需要权限 `user:interaction:read` 认证方式 `X-API-Key + OAuth Bearer`

两者均接受 `type`、`page` 和 `pageSize`。接口只返回仍处于已审核、公开、无访问限制状态的目标；目标转为私密、受限或归档后不会继续出现在结果里。

## 我的通知

GET `/api/site/v1/me/notifications`

需要权限 `user:notification:read` 认证方式 `X-API-Key + OAuth Bearer`

可传 `unread=true` 仅看未读，或 `unread=false` 仅看已读。

GET `/api/site/v1/me/notifications/unread-count`

需要权限 `user:notification:read` 认证方式 `X-API-Key + OAuth Bearer`

```
{
  "schemaVersion": "1",
  "count": 3,
  "generatedAt": "2026-09-05T00:00:00.000Z"
}
```

## 标记通知已读

POST `/api/site/v1/me/notifications/:id/read`

需要权限 `user:notification:write` 认证方式 `X-API-Key + OAuth Bearer`

POST `/api/site/v1/me/notifications/read-all`

需要权限 `user:notification:write` 认证方式 `X-API-Key + OAuth Bearer`

单条操作会同时校验通知 ID 和 OAuth 用户 ID；其他用户的通知统一按不存在处理。`read-all` 也只更新当前 OAuth 用户的数据。本阶段不开放通知删除。

## 错误码

| 状态码 | `code` | 含义 |
| --- | --- | --- |
| `400` | `SITE_API_INVALID_QUERY` | `type`、`status`、分页或 `unread` 参数无效。 |
| `401` | `SITE_API_KEY_INVALID` | 站点 API Key 无效。 |
| `401` | `OAUTH_ACCESS_TOKEN_INVALID` | OAuth Token 无效、过期或已撤销。 |
| `403` | `SITE_API_OAUTH_CLIENT_MISMATCH` | 两个凭据不属于同一个 OAuth 应用。 |
| `403` | `OAUTH_SCOPE_REQUIRED` | OAuth Token 缺少当前接口所需 scope。 |
| `404` | `NOTIFICATION_NOT_FOUND` | 通知不存在或不属于当前用户。 |
