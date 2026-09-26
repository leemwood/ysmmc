---
title: "NexusMCPublisher Gradle 插件"
source: https://docs.nexusmc.cn/docs/api/nexusmc-publisher-gradle
collected: 2026-09-26
---

# NexusMCPublisher Gradle 插件

在 Gradle 构建中上传单个或多个构建产物，并向已有 NexusMC 资源发布新版本。

## 项目简介

[NexusMCPublisher](https://plugins.gradle.org/plugin/io.github.lijinhong11.nexusmcpublisher) 是一个可从 Gradle Plugin Portal 直接安装的 Gradle 插件，用于把构建产物发布到已有 NexusMC 资源。项目源码位于 [lijinhong11/NexusMCPublisher](https://github.com/lijinhong11/NexusMCPublisher)。

插件使用个人 API 的两段式流程：先通过 `POST /api/upload` 上传每个文件，再通过 `POST /api/resources/{id}/versions` 提交版本和文件信息。项目兼容 Java 8。

## 安装插件

在 `build.gradle.kts` 的 `plugins` 中加入：

```
plugins {
    java
    id("io.github.lijinhong11.nexusmcpublisher") version "1.0.3"
}
```

上面是 Gradle Plugin Portal 在 2026 年 8 月 1 日显示的版本。接入新项目时请在[插件门户](https://plugins.gradle.org/plugin/io.github.lijinhong11.nexusmcpublisher)核对当前版本。

## 最小配置

```
import io.github.lijinhong11.nexusmcpublisher.VersionTag

version = "1.2.3"

nexusMCPublisher {
    resourceId.set("8a86e434-5050-4f4a-a861-5a5fe28873a4")
    versionTag.set(VersionTag.RELEASE)
    versionTitle.set("Minecraft 兼容性更新")
    changelog.set("修复已知问题并更新资源文件。")
    mcVersions.set(listOf("1.21.4"))
}
```

应用 Java 插件后，默认使用 `jar` 任务的输出。需要发布其他文件时可以显式指定：

```
nexusMCPublisher {
    artifact("build/libs/example.jar")
}
```

执行发布任务：

```
./gradlew publishToNexusMC
```

## Token 与权限

不要把 Token 写入构建脚本或提交到仓库。优先通过环境变量提供：

```
export NEXUSMC_API_TOKEN='avm_pat_xxxxx_xxxxx'
./gradlew publishToNexusMC
```

也可以写入当前用户的 `~/.gradle/gradle.properties`，但不能提交该文件：

```
nexusMCToken=avm_pat_xxxxx_xxxxx
```

发布版本需要以下最小权限：

| Token scope | 用途 |
| --- | --- |
| `upload:file` | 上传构建产物。 |
| `resource:update:self` | 为 Token 用户自己的资源发布版本。 |
| `resource:read:self` | 可选；通过“我的资源”接口查询稳定资源 ID。 |

## 查询资源 ID

`resourceId` 建议使用“我的资源”接口返回的稳定 `resources[].id`，不需要从网页 URL 获取：

```
curl "https://www.nexusmc.cn/api/users/me/resources?q=RESOURCE_TITLE&pageSize=10" \
  -H "Authorization: Bearer $NEXUSMC_API_TOKEN" \
  -H "Accept: application/json"
```

站点接口也兼容当前或历史 slug、公开编号和公开链接中的路径段，但 slug 可能随标题或自定义路径修改，不适合长期写死。`resourceId` 只是资源名称或拼写错误时，发布任务会返回 `404 Resource not found`。

## 发布多个文件

每个 `file` 块会独立上传，随后一起写入版本接口的 `files[]`。必须且只能有一个主文件：

```
import io.github.lijinhong11.nexusmcpublisher.VersionTag

nexusMCPublisher {
    resourceId.set("8a86e434-5050-4f4a-a861-5a5fe28873a4")
    versionTag.set(VersionTag.RELEASE)
    mcVersions.set(listOf("1.21.1", "1.21.4"))

    file {
        artifact("build/libs/plugin-paper.jar")
        primary.set(true)
        gameVersions.set(listOf("1.21.4"))
        loaders.set(listOf("paper"))
    }

    file {
        artifact("build/libs/plugin-fabric.jar")
        gameVersions.set(listOf("1.21.1", "1.21.4"))
        loaders.set(listOf("fabric"))
    }
}
```

没有声明 `file {}` 时，插件使用顶层 `artifact`；应用 Java 插件后则默认使用 `jar` 输出，并把它作为唯一主文件。

## 常用配置

| 配置 | 必填 | 默认值 | 说明 |
| --- | --- | --- | --- |
| `resourceId` | 是 | 无 | 已有 NexusMC 资源的稳定 ID。 |
| `token` | 是 | 环境变量或 Gradle 属性 | 个人 API Token。 |
| `version` | 是 | `project.version` | 要发布的版本号。 |
| `versionTag` | 否 | `VersionTag.RELEASE` | `RELEASE`、`BETA` 或 `ALPHA`。 |
| `versionTitle` | 否 | `Version <version>` | 版本标题。 |
| `changelog` | 否 | 空 | 版本更新说明。 |
| `mcVersions` | 否 | 空列表 | 本次发布支持的 Minecraft 版本。 |
| `downloadType` | 否 | `local` | NexusMC 下载类型。 |
| `artifact` | 是 | Java `jar` 输出 | 单文件发布时上传的构建产物。 |
| `file {}` | 否 | 无 | 可重复的多文件声明；存在时替代顶层 `artifact`。 |
| `baseUrl` | 否 | `https://www.nexusmc.cn` | API 基础地址，主要用于测试或私有部署。 |

完整配置和最新行为以项目的[开源仓库](https://github.com/lijinhong11/NexusMCPublisher)与[插件门户](https://plugins.gradle.org/plugin/io.github.lijinhong11.nexusmcpublisher)为准。NexusMC API 字段及服务端规则以[个人 API：资源文件与版本](/docs/api/personal-api-resource-files/)为准。

[返回可直接使用 API 的项目总览](/docs/api/ready-to-use-projects/)。
