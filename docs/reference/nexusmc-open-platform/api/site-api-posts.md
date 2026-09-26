---
title: "站点 API：帖子"
source: https://docs.nexusmc.cn/docs/api/site-api-posts
collected: 2026-09-26
---

# 站点 API：帖子

读取公开帖子的列表、版块分类、筛选项和详情。

以下接口均需要 `post:read:public`，且只读取已审核、公开、无访问限制的帖子。

## 列表

GET `/api/site/v1/posts?boardId={boardId}&page={page}&pageSize={pageSize}`

需要权限 `post:read:public` 认证方式 `X-API-Key`

支持 `boardId`、`page`、`pageSize`、`sort`、`includePinned` 与 `collection`。响应中的 `pagination.total` 为当前筛选的帖子总数。

可用 `sort` 为 `latest`、`replies`、`activity`；`includePinned=false` 可取消置顶优先。

## 分类

GET `/api/site/v1/posts/categories`

需要权限 `post:read:public` 认证方式 `X-API-Key`

返回站内公开版块树及各版块公开帖子数量。列表使用返回的版块 `id` 作为 `boardId`。

## 筛选

GET `/api/site/v1/posts/filters`

需要权限 `post:read:public` 认证方式 `X-API-Key`

返回可筛选版块、可用收录集、置顶开关及排序方式，适合第三方界面在调用列表前构建筛选控件。

## 详情

GET `/api/site/v1/posts/{idOrSlug}`

需要权限 `post:read:public` 认证方式 `X-API-Key`

`idOrSlug` 可使用帖子 ID、公开 ID 或 slug。返回站内公开帖子详情，包括允许公开的正文、版块和活动信息；草稿、私密和受限帖子统一不可读取。
