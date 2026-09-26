---
title: "站点 API：全站搜索"
source: https://docs.nexusmc.cn/docs/api/site-api-search
collected: 2026-09-26
---

# 站点 API：全站搜索

使用站点 API 搜索公开资源、帖子、找服玩和视频。

全站搜索面向官网、机器人和服务端索引程序，只返回已审核、公开且无访问限制的资源、帖子、找服玩和视频。需要站点 API 密钥，并授予 `search:read:public`。

GET `/api/site/v1/search`

需要权限 `search:read:public` 认证方式 `X-API-Key`

## 请求

```
curl 'https://www.nexusmc.cn/api/site/v1/search?q=fabric&type=resources&sort=relevance&page=1&pageSize=20' \
  -H "X-API-Key: $NEXUSMC_SITE_API_KEY" \
  -H 'Accept: application/json'
```

| 参数 | 必填 | 说明 |
| --- | --- | --- |
| `q` | 是 | 搜索词，最长 100 个字符。 |
| `type` | 否 | `all`、`resources`、`posts`、`players`、`videos`，默认 `all`。 |
| `sort` | 否 | `relevance`、`latest`、`popular`，默认 `relevance`。 |
| `page` | 否 | 从 1 开始，默认 `1`。 |
| `pageSize` | 否 | 每页 1-50 条，默认 `20`。 |
| `platform` | 否 | 资源或找服玩的平台筛选值。 |
| `category` | 否 | 资源、找服玩或视频分类筛选值。 |
| `boardId` | 否 | 帖子版块 ID。 |
| `tag` | 否 | 标签筛选值。 |

`type=all` 时，`items` 是四类结果合并后的当前页，`groups` 保留各类型总数。调用方如需稳定分页和单一内容结构，建议分别请求具体 `type`。

## 响应

```
{
  "schemaVersion": "1",
  "items": [
    {
      "id": "resource-id",
      "type": "resource",
      "slug": "fabric-example",
      "title": "Fabric 示例资源",
      "excerpt": "公开资源摘要",
      "createdAt": "2026-09-04T08:00:00.000Z",
      "updatedAt": "2026-09-04T08:30:00.000Z"
    }
  ],
  "groups": { "resources": 1, "posts": 0, "players": 0, "videos": 0 },
  "pagination": { "page": 1, "pageSize": 20, "total": 1, "totalPages": 1 },
  "appliedFilters": { "q": "fabric", "type": "resources", "sort": "relevance", "platform": "", "category": "", "boardId": "", "tag": "" },
  "generatedAt": "2026-09-04T08:30:00.000Z"
}
```

结果中的 `id` 和 `slug` 可用于拼接官网详情链接。搜索不会返回草稿、私密帖子、付费内容、归档资源或审核资料；无结果时返回空数组和正常分页对象，不是错误。

| 状态码 | code | 说明 |
| --- | --- | --- |
| `400` | `SEARCH_QUERY_REQUIRED` | 缺少 `q` 或搜索词为空。 |
| `400` | `SEARCH_TYPE_INVALID` | `type` 不在允许范围内。 |
| `401` | `SITE_API_KEY_INVALID` | 密钥无效、已停用、已撤销或已过期。 |
| `403` | `SITE_API_KEY_SCOPE_REQUIRED` | 密钥没有 `search:read:public`。 |
