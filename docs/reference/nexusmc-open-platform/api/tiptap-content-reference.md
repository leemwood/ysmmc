---
title: "TipTap 节点与属性速查"
source: https://docs.nexusmc.cn/docs/api/tiptap-content-reference
collected: 2026-09-26
---

# TipTap 节点与属性速查

列出个人 API 正文常用 TipTap 节点、标记和关键 attrs，方便同步脚本直接对照生成 JSON。

## 常用节点总表

| 能力 | 节点 / 标记 | 位置 | 说明 |
| --- | --- | --- | --- |
| 普通段落 | `paragraph` | block | 最常见的正文块。 |
| 标题 | `heading` | block | 用 `attrs.level` 表示 1 到 6 级标题。 |
| 引用 | `blockquote` | block | 内部通常包含段落。 |
| 分割线 | `horizontalRule` | block | 无正文内容。 |
| 无序列表 | `bulletList` / `listItem` | block | 列表项内部包含段落或其他块。 |
| 有序列表 | `orderedList` / `listItem` | block | 可配合 `attrs.start`。 |
| 任务列表 | `taskList` / `taskItem` | block | `taskItem.attrs.checked` 表示是否完成。 |
| 代码块 | `codeBlock` | block | 文本内容放在子 `text` 节点。 |
| 图片 | `image` | block | `attrs.src` 是图片 URL。 |
| 表格 | `table` / `tableRow` / `tableHeader` / `tableCell` | block | 单元格内通常包含段落。 |
| @提及 | `mention` | inline | `attrs.id` 是用户 ID，`attrs.label` 是用户名。 |
| 表情 | `emoji` | inline | 用于站内表情。 |
| 折叠块 | `collapse` | block | `attrs.summary` 是折叠标题。 |
| Mermaid 流程图 | `mermaidDiagram` | block | `attrs.code` 存流程图源码。 |
| 数学公式 | `math` | inline / block | `attrs.inline` 区分行内和块级。 |
| 动态图 | `dynamicImage` | block | 支持刷新按钮和加载时刷新。 |
| Bilibili 嵌入视频 | `bilibiliVideo` | block | 使用 `attrs.bvid` 和 `attrs.page`。 |
| 链接卡片 | `linkCardInline` / `linkCardBlock` | inline / block | 用于富链接预览。 |

## 常用文本标记

| 能力 | mark | 关键 attrs |
| --- | --- | --- |
| 粗体 | `bold` | 无 |
| 斜体 | `italic` | 无 |
| 删除线 | `strike` | 无 |
| 行内代码 | `code` | 无 |
| 链接 | `link` | `href` |
| 下划线 | `underline` | 无 |
| 高亮 | `highlight` | `color` |
| 上标 | `superscript` | 无 |
| 下标 | `subscript` | 无 |
| 字色 / 字号 / 字体 | `textStyle` | `color`、`fontSize`、`fontFamily` |

## 代码块属性

| 属性 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `language` | `string` | 否 | 代码语言，例如 `json`、`js`、`ts`、`bash`、`diff`。 |
| `highlightLines` | `string` | 否 | 高亮行，格式如 `1`、`1,3`、`2-5`。 |
| `showLineNumbers` | `boolean` | 否 | 是否显示行号。 |
| `lineNumberStart` | `number` | 否 | 行号起始值，默认 `1`。 |
| `diffLanguage` | `string` | 否 | `language` 为 `diff` 时，指定 diff 内部语言。 |

## 图片属性

| 属性 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `src` | `string` | 是 | 图片 URL。站内图片通常来自上传接口返回值。 |
| `alt` | `string` | 否 | 替代文本。 |
| `title` | `string` | 否 | 图片标题。 |
| `width` | `string` | 否 | 显示宽度，例如 `720px` 或 `100%`。 |
| `height` | `string` | 否 | 显示高度。 |
| `style` | `string` | 否 | 兼容已有内容的样式字段。 |

## 表格结构

| 层级 | 节点 | 内容 |
| --- | --- | --- |
| 表格 | `table` | 包含多个 `tableRow`。 |
| 行 | `tableRow` | 包含多个 `tableHeader` 或 `tableCell`。 |
| 表头单元格 | `tableHeader` | 通常包含一个 `paragraph`。 |
| 普通单元格 | `tableCell` | 通常包含一个 `paragraph`。 |

## 特殊块属性

| 节点 | 属性 | 类型 | 说明 |
| --- | --- | --- | --- |
| `collapse` | `summary` | `string` | 折叠标题。 |
| `collapse` | `open` | `boolean` | 默认是否展开。 |
| `mermaidDiagram` | `code` | `string` | Mermaid 源码。当前推荐使用 `flowchart` / `graph` 流程图。 |
| `math` | `latex` | `string` | LaTeX 内容。 |
| `math` | `inline` | `boolean` | `true` 为行内公式，`false` 为块级公式。 |
| `dynamicImage` | `src` | `string` | 动态图片 URL。 |
| `dynamicImage` | `refreshOnLoad` | `boolean` | 页面加载时是否刷新。 |
| `dynamicImage` | `showRefreshButton` | `boolean` | 是否显示刷新按钮。 |
| `bilibiliVideo` | `bvid` | `string` | Bilibili 视频 BV 号。 |
| `bilibiliVideo` | `page` | `number` | 分 P 页码，默认 `1`。 |
| `linkCardInline` / `linkCardBlock` | `url` | `string` | 目标链接。 |
| `linkCardInline` / `linkCardBlock` | `title` | `string` | 卡片标题。 |
| `linkCardInline` / `linkCardBlock` | `description` | `string` | 卡片描述。 |
| `linkCardInline` / `linkCardBlock` | `image` | `string` | 预览图 URL。 |
| `linkCardInline` / `linkCardBlock` | `siteName` | `string` | 站点名。 |
| `linkCardInline` / `linkCardBlock` | `favicon` | `string` | 站点图标 URL。 |
| `linkCardInline` / `linkCardBlock` | `variant` | `string` | `inline` 或 `normal`。 |
