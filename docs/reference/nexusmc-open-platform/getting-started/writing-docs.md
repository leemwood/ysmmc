---
title: "使用 MDX 编写文档"
source: https://docs.nexusmc.cn/docs/getting-started/writing-docs
collected: 2026-09-26
---

# 使用 MDX 编写文档

面向文档维护者，说明如何新增、组织、编写和验证 docs 内容。

## 适用对象

这篇文档面向维护 `apps/docs/src/content/docs` 的人。你可以用它来新增说明、补充投稿流程、更新 API 文档或调整资源规范。

文档内容使用 MDX 编写，但日常写作尽量保持 Markdown 风格。只有确实需要交互组件或特殊展示时，再引入 MDX 组件。

TIP

优先用 Markdown 写清楚内容；只有当提示、折叠、代码组能降低阅读成本时，再使用增强语法。

## 文件放在哪里

当前文档内容位于：

```
apps/docs/src/content/docs
```

推荐按分类放置：

```
docs/
├─ api/
├─ getting-started/
├─ operations/
├─ site-info/
└─ standards/
```

文件名应使用英文小写和短横线：

```
resource-publish-spec.mdx
post-submission-flow.mdx
personal-api-tokens.mdx
```

不要使用空格、中文文件名或含义不清的编号文件名。

## Frontmatter

每篇文档必须有 frontmatter。

```
---
title: 文档标题
description: 一句说明这篇文档解决什么问题
section: 投稿流程
order: 10
---
```

字段说明：

| 字段 | 必填 | 说明 |
| --- | --- | --- |
| `title` | 是 | 页面标题，也会出现在文档目录中。 |
| `description` | 否 | 页面摘要，建议填写。 |
| `section` | 是 | 左侧导航分组，如 `快速开始`、`API`、`投稿流程`、`资源规范`。 |
| `order` | 否 | 分组内排序，数字越小越靠前。 |
| `draft` | 否 | 设为 `true` 时不公开生成。 |
| `locale` | 否 | 多语言预留，如 `en`。 |
| `translationKey` | 否 | 多语言同篇文档关联键。 |

## 写作流程

1.  明确读者是谁：普通用户、投稿者、资源作者、脚本开发者，还是文档维护者。
2.  明确用户要完成的任务。
3.  先写结论和适用范围，再写细节。
4.  用表格列字段、状态和错误原因。
5.  用代码块放请求、响应、命令或配置。
6.  结尾补“常见失败原因”或“检查清单”。
7.  本地运行构建，确认页面能生成。

IMPORTANT

写文档时先确认内容是否已实现。规划中、未上线或只在本地分支存在的能力，不要写成用户可直接使用的功能。

## 推荐结构

流程类文档建议使用：

```
## 适用范围

## 状态流转

## 标准流程

## 必填信息

## 编辑或更新规则

## 常见失败原因

## 提交前检查清单
```

规范类文档建议使用：

```
## 适用范围

## 填写要求

## 推荐写法

## 不推荐写法

## 示例

## 检查清单
```

API 类文档建议使用：

```
## Base URL

## 认证方式

## 响应格式

## 接口总览

## 请求参数

## 请求示例

## 响应示例

## 常见失败原因

## 疑难解答
```

## Markdown 写法约定

### 标题

标题层级从 `##` 开始。页面标题已经由 frontmatter 的 `title` 提供，不要在正文第一行再写一个重复的 `# 标题`。

### 表格

字段说明、状态说明、错误原因优先用表格。

```
| 字段 | 必填 | 说明 |
| --- | --- | --- |
| `title` | 是 | 标题。 |
```

### 代码块

请求、响应、命令和配置使用 fenced code block，并标注语言。

```
{
  "error": "Unauthorized"
}
```

### 代码组

使用 `code-group` 可以把多个代码块放进同一组，用 tab 切换查看。代码块标题写在语言后的方括号里。

NOTE

API 文档里优先用代码组放 `curl`、`fetch` 或 SDK 示例；配置文档里可以放 JavaScript、TypeScript、YAML 等不同写法。

```
::: code-group

```js [config.js]
const config = {
  // ...
}

export default config
```

```ts [config.ts]
import type { UserConfig } from 'vitepress'

const config: UserConfig = {
  // ...
}

export default config
```

:::
```

```
const config = {
  // ...
}

export default config
```

```
import type { UserConfig } from 'vitepress'

const config: UserConfig = {
  // ...
}

export default config
```

### 提示容器

使用 `:::` 可以写提示容器，支持 `note`、`tip`、`important`、`warning`、`caution`。`info` 会按 `note` 处理，`danger` 会按 `caution` 处理。

推荐使用方式：

| 类型 | 适合内容 |
| --- | --- |
| `note` | 补充说明、背景、适用范围。 |
| `tip` | 推荐做法、效率建议。 |
| `important` | 必须遵守的规则。 |
| `warning` | 容易失败或需要注意的限制。 |
| `caution` | 可能造成安全、审核或数据风险的行为。 |

```
::: note
这是一条普通说明。
:::

::: tip
这是一条建议。
:::

::: important
这是需要优先关注的信息。
:::

::: warning
这是风险提醒。
:::

::: caution
这是高风险或负面影响提醒。
:::
```

NOTE

这是一条普通说明。

TIP

这是一条建议。

IMPORTANT

这是需要优先关注的信息。

WARNING

这是风险提醒。

CAUTION

这是高风险或负面影响提醒。

### 自定义提示标题

标题可以写在类型后，也可以使用方括号。

```
::: note 自定义标题
这里是带自定义标题的说明。
:::

::: caution[STOP]
这里是带自定义标题的风险提醒。
:::
```

自定义标题

这里是带自定义标题的说明。

STOP

这里是带自定义标题的风险提醒。

### GitHub Alert

也支持 GitHub 风格的 Alert 写法。

```
> [!NOTE]
> 这里是 Note。

> [!TIP]
> 这里是 Tip。

> [!IMPORTANT]
> 这里是 Important。

> [!WARNING]
> 这里是 Warning。

> [!CAUTION]
> 这里是 Caution。
```

NOTE

这里是 Note。

TIP

这里是 Tip。

IMPORTANT

这里是 Important。

WARNING

这里是 Warning。

CAUTION

这里是 Caution。

### 折叠容器

使用 `details` 创建可折叠内容。标题不写时默认为 `Details`。

适合折叠的内容包括长示例、失败原因对照表、检查清单、FAQ 答案和完整模板。

```
::: details 点击查看代码
这里可以放普通 Markdown。

```js
console.log('Hello, NexusMC!')
```
:::
```

点击查看代码

这里可以放普通 Markdown。

```
console.log('Hello, NexusMC!')
```

### 剧透文本

使用 `:spoiler[...]` 隐藏一段文字。鼠标悬停或键盘聚焦后会显示。

```
这段内容 :spoiler[默认隐藏 **但仍支持加粗**]。
```

这段内容 默认隐藏 **但仍支持加粗**。

### 链接

站内链接使用绝对路径，方便静态生成后稳定访问：

```
[用户个人 API](/docs/api/personal-api-tokens/)
```

### 状态和字段

状态、字段名、路径、scope、命令使用行内代码：

```
资源提交后进入 `pending`，审核通过后变为 `approved`。
```

## 内容准确性要求

写文档时要区分三类信息：

| 类型 | 写法 |
| --- | --- |
| 已实现 | 可以直接写“支持”。 |
| 暂不支持 | 明确写“当前不支持”。 |
| 规划中 | 不要写进用户操作文档；必要时标成“后续规划”。 |

不要把服务端尚未实现的字段、状态或接口写成已开放能力。

WARNING

如果不确定接口、字段或状态是否真实存在，先查服务端实现或问维护者，不要靠猜测补文档。

当前常用状态是：

-   `draft`
-   `pending`
-   `approved`
-   `rejected`

不要用 `published` 代替 `approved`。

## 多语言预留

默认语言是简体中文，默认路径保持 `/docs/...`。

如果后续新增英文文档，可以在 frontmatter 中声明：

```
locale: en
translationKey: getting-started/overview
```

英文文档建议放到：

```
apps/docs/src/content/docs/en/...
```

## 本地验证

写完文档后运行：

```
npm run build
```

如果只是本地预览：

```
npm run dev
```

构建通过后，确认：

-   页面能生成。
-   文档出现在正确分组。
-   左侧顺序符合 `order`。
-   表格和代码块在手机宽度下不会撑出横向页面。
-   站内链接能打开。

## 部署与缓存

文档站是 Astro 服务端输出。部署后，HTML 页面不应该被浏览器、反代或 CDN 长时间强缓存；建议使用 `Cache-Control: public, max-age=0, must-revalidate`。`/_astro/` 下的 CSS、JS 等带 hash 的构建资源可以使用 `public, max-age=31536000, immutable`。

如果访问文档时页面内容还在，但样式和交互全失效，并且开发者工具里出现旧的 `/_astro/*.css` 或 `/_astro/*.js` 资源 `404`，通常是旧 HTML 被缓存了。旧 HTML 引用了上一次构建的 hash 资源，而服务器上只保留了新构建产物。处理方式是刷新 HTML 缓存、检查反代/CDN 的页面缓存规则，并确认应用响应带有正确的 `Cache-Control`。

## 新文档模板

```
---
title: 文档标题
description: 一句说明这篇文档解决什么问题。
section: 投稿流程
order: 10
---

## 适用范围

说明这篇文档适合谁阅读、解决什么问题。

## 核心流程

1. 第一步。
2. 第二步。
3. 第三步。

## 关键字段

| 字段 | 说明 |
| --- | --- |
| `example` | 示例字段。 |

## 常见失败原因

| 场景 | 原因 | 处理方式 |
| --- | --- | --- |
|  |  |  |

## 检查清单

-
-
-
```
