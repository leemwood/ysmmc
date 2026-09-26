---
title: "Upload to NexusMC"
source: https://docs.nexusmc.cn/docs/api/upload-to-nexusmc
collected: 2026-09-26
---

# Upload to NexusMC

使用 GitHub Actions 自动创建或更新 NexusMC 资源、上传文件、发布版本并转换 Markdown 正文。

## 项目简介

[Upload to NexusMC](https://github.com/marketplace/actions/upload-to-nexusmc) 适合把 GitHub 仓库的构建和发布流程接到 NexusMC。它遵循官方的两段式资源发布流程：先上传本地文件，再把上传接口返回的文件信息提交给资源 API。

项目源码位于 [mcio-dev/upload-to-nexusmc](https://github.com/mcio-dev/upload-to-nexusmc)。当前可直接引用的主版本为：

```
uses: mcio-dev/upload-to-nexusmc@v1
```

## 支持的操作

-   创建 NexusMC 资源。
-   更新当前 Token 用户自己的资源。
-   通过独立版本接口发布新版本。
-   上传单个文件、多个 Loader 文件和封面图片。
-   根据文件大小自动尝试直传、普通上传和分块上传。
-   把 `README.md` 或 GitHub Release 正文转换成 NexusMC TipTap JSON。

## 准备 Token 和资源 ID

在仓库的 `Settings -> Secrets and variables -> Actions` 中添加：

| 类型 | 名称 | 内容 |
| --- | --- | --- |
| Repository secret | `NEXUSMC_API_TOKEN` | 在 NexusMC 设置页创建的个人 API Token。 |
| Repository variable | `NEXUSMC_RESOURCE_ID` | 已有资源的 ID；只在更新资源或发布版本时需要。 |

按实际操作授予最小权限：

| 操作 | Token scope |
| --- | --- |
| 上传资源文件或封面 | `upload:file` |
| 创建资源 | `resource:create` |
| 更新自己的资源或发布版本 | `resource:update:self` |

IMPORTANT

Token 必须保存为 GitHub Actions secret，不要放进工作流明文、仓库变量、源码或构建日志。一个自动化项目建议使用一个独立 Token，并且只授予它实际需要的 scope。

### 查询资源 ID

不要把资源名当成 `resource_id`，也不需要从地址栏手工截取。使用带 `resource:read:self` 权限的 Token 按标题查询：

```
curl "https://www.nexusmc.cn/api/users/me/resources?q=MittelLib&pageSize=10" \
  -H "Authorization: Bearer $NEXUSMC_API_TOKEN" \
  -H "Accept: application/json"
```

把匹配项的 `resources[].id` 保存到仓库变量 `NEXUSMC_RESOURCE_ID`。响应还会同时返回 `slug` 和 `path`，方便核对是否选中了正确资源。

```
{
  "resources": [
    {
      "id": "8a86e434-5050-4f4a-a861-5a5fe28873a4",
      "slug": "mittelib",
      "path": "/resources/mittelib",
      "title": "MittelLib"
    }
  ]
}
```

自动化配置仍建议使用稳定 `id`。站点接口也接受当前或历史 slug、公开编号和公开链接中的路径段，因此 `mittelib` 确实是该资源 slug 时也能解析；纯资源名称或拼写错误仍会返回 `404 Resource not found`。

## 自动上传并更新资源

下面的工作流会在 GitHub Release 发布后构建项目，上传构建文件，并更新已有的 NexusMC 资源和版本信息：

```
name: Publish to NexusMC

on:
  release:
    types: [published]

jobs:
  publish:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Build
        run: ./gradlew build

      - name: Upload to NexusMC
        uses: mcio-dev/upload-to-nexusmc@v1
        with:
          api_token: ${{ secrets.NEXUSMC_API_TOKEN }}
          resource_id: ${{ vars.NEXUSMC_RESOURCE_ID }}
          file_path: build/libs/MyPlugin-${{ github.ref_name }}.jar
          version: ${{ github.ref_name }}
          version_tag: releases
          version_title: Release ${{ github.ref_name }}
          changelog_markdown: ${{ github.event.release.body }}
          publish_version: true
```

把 `Build` 命令和 `file_path` 改成项目实际使用的构建方式与精确文件路径。`file_path` 不接受模糊匹配，Action 执行时文件必须已经存在。

`resource_id` 存在时，默认的 `auto` 操作会更新已有资源；不提供 `resource_id` 时则会创建资源。生产工作流可以显式设置 `operation: update` 或 `operation: create`，避免配置错误时执行另一种操作。

CAUTION

出现 `NexusMC API request failed with HTTP 404: {"error":"Resource not found"}` 时，文件上传步骤通常已经成功，失败的是随后更新资源或发布版本的请求。先用上面的列表接口确认 `resource_id` 属于当前 Token 用户，再检查仓库变量是否有多余空格、是否误填了资源名称。

## 创建新资源

创建资源时不填写 `resource_id`，并提供标题、分类和正文等必填信息：

```
- name: Create NexusMC resource
  uses: mcio-dev/upload-to-nexusmc@v1
  with:
    api_token: ${{ secrets.NEXUSMC_API_TOKEN }}
    operation: create
    title: My Plugin
    category: plugin
    platform: java
    content_markdown_path: README.md
    version: 1.0.0
    version_tag: releases
    file_path: build/libs/MyPlugin-1.0.0.jar
    is_draft: true
```

建议第一次自动创建时保留 `is_draft: true`，到站内检查标题、正文、分类、文件和版本信息后再正式提交。

## 使用前检查

1.  确认账号已经解锁个人 API Token，并具备对应的创作者发布资格。
2.  确认 Token 只包含当前工作流需要的 scope。
3.  确认构建产物的精确路径与文件名。
4.  确认 `version_tag` 是站点当前启用的 Tag key，例如 `releases` 或 `beta`。
5.  私有仓库的相对图片默认不能被访客读取；正文图片应使用公开地址或 NexusMC 图片上传接口返回的 URL。
6.  先使用草稿或预发布版本验证一次，再接入正式 Release 流程。

完整参数和多文件示例以该项目的 [README](https://github.com/mcio-dev/upload-to-nexusmc/blob/v1/README.md) 为准。API 字段与服务端规则仍以[个人 API：资源上传与发布流程](/docs/api/personal-api-resource-upload-flow/)为准。

[返回可直接使用 API 的项目总览](/docs/api/ready-to-use-projects/)。
