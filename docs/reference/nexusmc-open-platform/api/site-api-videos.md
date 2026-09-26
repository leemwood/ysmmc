---
title: "站点 API：视频区"
source: https://docs.nexusmc.cn/docs/api/site-api-videos
collected: 2026-09-26
---

# 站点 API：视频区

读取公开视频的列表、分类、筛选项和详情。

以下接口均需要 `video:read:public`，并且只返回已审核、公开、无访问限制的视频。

## 列表

GET `/api/site/v1/videos?page={page}&pageSize={pageSize}`

需要权限 `video:read:public` 认证方式 `X-API-Key`

支持 `page`、`pageSize`、`sort`、`category` 与 `collection`。响应中的 `pagination.total` 为当前筛选条件下的视频总数；可用排序为 `latest`、`hot`、`discuss`。

## 分类

GET `/api/site/v1/videos/categories`

需要权限 `video:read:public` 认证方式 `X-API-Key`

返回视频分类以及每个分类的公开视频数量。

## 筛选

GET `/api/site/v1/videos/filters`

需要权限 `video:read:public` 认证方式 `X-API-Key`

返回分类、收录集与排序方式，供第三方界面构建筛选控件。

## 详情

GET `/api/site/v1/videos/{idOrSlug}`

需要权限 `video:read:public` 认证方式 `X-API-Key`

`idOrSlug` 可使用视频 ID、公开 ID 或 slug。该请求不会增加站内播放/浏览统计。
