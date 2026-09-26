---
title: "站点 API：排行榜"
source: https://docs.nexusmc.cn/docs/api/site-api-rankings
collected: 2026-09-26
---

# 站点 API：排行榜

获取资源、找服玩和视频的公开排行榜。

排行榜接口复用官网排行榜的统计口径，只返回公开内容。先读取过滤器可获得当前支持的榜单类型和筛选项，需要 `ranking:read:public`。

GET `/api/site/v1/rankings/filters`

需要权限 `ranking:read:public` 认证方式 `X-API-Key`

GET `/api/site/v1/rankings`

需要权限 `ranking:read:public` 认证方式 `X-API-Key`

```
curl 'https://www.nexusmc.cn/api/site/v1/rankings?type=resource-downloads&period=month&limit=20' \
  -H "X-API-Key: $NEXUSMC_SITE_API_KEY"
```

## 榜单类型

| `type` | 内容 | 指标 |
| --- | --- | --- |
| `resource-downloads` | 资源下载榜 | 下载 |
| `resource-likes` | 资源点赞榜 | 点赞 |
| `player-peak-online` | 找服玩峰值在线榜 | 峰值在线 |
| `player-play-heat` | 找服玩游玩热度榜 | 热度 |
| `player-likes` | 找服玩点赞榜 | 点赞 |
| `video-plays` | 视频播放榜 | 播放 |

## 参数与响应

`type` 必填；`period` 可选 `week`、`month`、`year`，默认值以当前排行榜查询规则为准。年度、月份和周数可分别用 `year`、`month`、`week` 指定；`platform` 和 `category` 用于内容筛选；`limit` 范围为 5-100。

```
{
  "schemaVersion": "1",
  "ranking": {
    "type": "resource-downloads",
    "typeMeta": { "value": "resource-downloads", "label": "资源下载榜", "sourceType": "resource", "metricLabel": "下载" }
  },
  "range": { "period": "month", "year": 2026, "month": 9, "week": null },
  "items": [
    { "rank": 1, "id": "resource-id", "value": 1280, "metrics": { "downloads": 1280 } }
  ],
  "generatedAt": "2026-09-04T08:30:00.000Z"
}
```

`items` 的展示标题、slug 和详情字段以对应内容域的公开数据为准，客户端应使用 `id` 作为稳定标识。非法 `type`、时间范围或超出 5-100 的 `limit` 会由请求校验拒绝，不会静默改成另一个榜单。
