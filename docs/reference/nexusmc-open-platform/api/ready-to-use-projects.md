---
title: "可直接使用 API 的项目"
source: https://docs.nexusmc.cn/docs/api/ready-to-use-projects
collected: 2026-09-26
---

# 可直接使用 API 的项目

汇总已经接入 NexusMC API、可以直接用于自动发布或同步内容的工具与项目。

## 项目目录

这里汇总已经封装 NexusMC API、无需从头编写请求代码即可使用的项目。每个工具都有独立文档页，用于记录安装方式、所需权限、配置示例和使用限制。

| 项目 | 类型 | 适用场景 | 文档 |
| --- | --- | --- | --- |
| Upload to NexusMC | GitHub Action | 从 GitHub Actions 创建或更新资源、发布版本、上传资源文件，并把 GitHub Markdown 转换成 TipTap。 | [查看使用文档](/docs/api/upload-to-nexusmc/) |
| NexusMCPublisher | Gradle 插件 | 在 Gradle 构建中上传 JAR 等产物，并为已有 NexusMC 资源发布单文件或多文件版本。 | [查看使用文档](/docs/api/nexusmc-publisher-gradle/) |

## 收录范围

本目录优先收录满足以下条件的项目：

-   已经公开发布，用户可以直接安装或引用。
-   明确使用 NexusMC 公开 API，而不是模拟网页操作。
-   说明所需 Token scope，不要求用户提供超出功能范围的权限。
-   提供可核对的源码、版本或维护入口。
-   能完成一项明确任务，例如自动发布资源、同步版本或转换正文格式。

NOTE

这些项目可能由第三方维护。使用前请检查项目源码、当前版本和所需权限；NexusMC API 的字段与服务端规则仍以本站 API 文档为准。

## 接入前准备

无论使用哪个工具，都建议先完成以下准备：

1.  解锁个人 API Token，并确认账号具备对应的发布资格。
2.  为每个自动化项目创建独立 Token。
3.  只授予工具实际需要的最小 scope。
4.  把 Token 保存在对应平台的 secret 管理功能中，不要写入源码或公开配置。
5.  第一次运行时使用草稿或测试内容验证结果。

如果需要自己编写接入程序，请从[用户个人 API：接入总览](/docs/api/personal-api-tokens/)开始。
