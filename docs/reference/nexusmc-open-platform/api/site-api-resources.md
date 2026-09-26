---
title: "站点 API：资源"
source: https://docs.nexusmc.cn/docs/api/site-api-resources
collected: 2026-09-26
---

# 站点 API：资源

读取公开资源的列表、分类、筛选项和详情。

以下接口均需要 `resource:read:public`。只返回已审核、公开、无访问限制且未归档的资源。

## 列表

GET `/api/site/v1/resources?page={page}&pageSize={pageSize}`

需要权限 `resource:read:public` 认证方式 `X-API-Key`

支持站内资源页已有的 `sort`、`platform`、`category`、`subCategory`、`sideSupport`、`thirdLevel`、`version`、`tag`、`collection`、`officialTags` 与 `customCategorySelections` 筛选参数。

响应沿用站内资源列表结构，`pagination.total` 是符合筛选条件的总数量。

```
curl "https://www.nexusmc.cn/api/site/v1/resources?platform=java&category=plugin&page=1&pageSize=20" \
  -H "X-API-Key: $NEXUSMC_SITE_API_KEY"
```

## 分类

GET `/api/site/v1/resources/categories?platform={platform}`

需要权限 `resource:read:public` 认证方式 `X-API-Key`

可选 `platform`。返回资源根分类及子分类，供列表请求中的 `category` 和 `subCategory` 使用。

## 筛选

GET `/api/site/v1/resources/filters?platform={platform}&groupId={groupId}`

需要权限 `resource:read:public` 认证方式 `X-API-Key`

返回平台、分类、版本组、游戏版本、官方标签组、自定义分类组、收录集和可用排序方式。`platform` 与可选的 `groupId` 会同时影响分类和版本候选项。

## 详情

GET `/api/site/v1/resources/{idOrSlug}`

需要权限 `resource:read:public` 认证方式 `X-API-Key`

`idOrSlug` 可使用资源 ID、公开 ID 或 slug。详情保留站内公开资源详情结构，不会返回下载口令、后台编辑信息或访问受限内容。
