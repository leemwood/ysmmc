---
title: "站点 API：活动"
source: https://docs.nexusmc.cn/docs/api/site-api-activities
collected: 2026-09-26
---

# 站点 API：活动

获取公开活动、活动状态和关联帖子。

活动接口把启用活动功能的公开版块中的帖子事件整理为统一数据。只返回公开、已审核且无访问限制的关联帖子，需要 `activity:read:public`。

GET `/api/site/v1/activities`

需要权限 `activity:read:public` 认证方式 `X-API-Key`

GET `/api/site/v1/activities/active`

需要权限 `activity:read:public` 认证方式 `X-API-Key`

```
curl 'https://www.nexusmc.cn/api/site/v1/activities?status=upcoming&page=1&pageSize=20' \
  -H "X-API-Key: $NEXUSMC_SITE_API_KEY"
```

## 参数

| 参数 | 说明 |
| --- | --- |
| `status` | `upcoming`（默认）、`active`、`ended` 或 `all`。`/active` 路径默认使用 `active`。 |
| `page` / `pageSize` | 分页参数，页码从 1 开始，单页最多 50 条。 |
| `boardId` | 限定一个启用活动的帖子版块。 |
| `from` / `to` | ISO 8601 日期时间窗口。 |

`upcoming` 要求活动状态为 `scheduled` 且开始时间在当前时间之后；`active` 要求已开始且未结束；`ended` 按结束时间早于当前时间判断。`all` 会在每条项目上返回计算后的 `status`。

## 响应

```
{
  "schemaVersion": "1",
  "items": [
    {
      "id": "event-id",
      "title": "新版本讨论活动",
      "startsAt": "2026-09-05T12:00:00.000Z",
      "endsAt": "2026-09-06T12:00:00.000Z",
      "status": "upcoming",
      "post": {
        "id": "post-id",
        "title": "活动公告",
        "slug": "event-announcement",
        "board": { "id": "board-id", "name": "公告" },
        "author": { "id": "user-id", "username": "NexusMC" }
      }
    }
  ],
  "pagination": { "page": 1, "pageSize": 20, "total": 1, "totalPages": 1 },
  "generatedAt": "2026-09-04T08:30:00.000Z"
}
```

日期参数格式错误返回 `400 ACTIVITY_DATE_INVALID`，状态值错误返回 `400 ACTIVITY_STATUS_INVALID`。活动接口不会暴露未审核帖子、私密帖子或没有活动能力的版块。
