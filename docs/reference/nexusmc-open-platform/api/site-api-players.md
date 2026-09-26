---
title: "站点 API：找服玩"
source: https://docs.nexusmc.cn/docs/api/site-api-players
collected: 2026-09-26
---

# 站点 API：找服玩

读取公开找服玩的列表、分类、筛选项和详情。

以下接口均需要 `player:read:public`。只返回已审核、公开、无访问限制且未归档的找服玩项目。

## 列表

GET `/api/site/v1/players?page={page}&pageSize={pageSize}`

需要权限 `player:read:public` 认证方式 `X-API-Key`

支持 `page`、`pageSize`、`sort`、`platform`、`category`、`version`、`onlineMode`、`q` 与 `collection`。响应中的 `pagination.total` 为符合当前筛选条件的总数。

`onlineMode` 使用字符串 `true` 或 `false`。可用排序为 `latest`、`updated`、`views`、`rating`、`score`。

## 分类

GET `/api/site/v1/players/categories?platform={platform}`

需要权限 `player:read:public` 认证方式 `X-API-Key`

可选 `platform`，返回当前平台可用的找服玩分类。

## 筛选

GET `/api/site/v1/players/filters?platform={platform}`

需要权限 `player:read:public` 认证方式 `X-API-Key`

返回平台、分类、版本、收录集、在线模式和排序方式。传入 `platform` 会同步限制分类与版本候选项。

## 详情

GET `/api/site/v1/players/{idOrSlug}`

需要权限 `player:read:public` 认证方式 `X-API-Key`

`idOrSlug` 可使用找服玩 ID、公开 ID 或 slug。该请求不会增加站内浏览统计，适合第三方同步与定时查询。
