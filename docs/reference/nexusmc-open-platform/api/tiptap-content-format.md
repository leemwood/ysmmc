---
title: "投稿正文 TipTap 快速接入"
source: https://docs.nexusmc.cn/docs/api/tiptap-content-format
collected: 2026-09-26
---

# 投稿正文 TipTap 快速接入

第一次通过个人 API 提交正文时先读这里：哪些业务有 content，应该传对象还是字符串，最小 JSON 长什么样。

## 先看字段边界

个人 API 里不是所有投稿都有富文本正文。接入前先确认你正在写哪个字段：

| 业务 | 富文本正文 | 短简介 | 字段规则 |
| --- | --- | --- | --- |
| 资源 | `content` | `description` | `content` 支持 TipTap JSON；`description` 是短简介。 |
| 帖子 | `content` | 无 | 帖子正文走 `content`，推荐传 TipTap JSON。 |
| 找服玩 | `content` | `description` | `content` 是详细介绍；`description` 是短介绍。 |
| 视频 | 无 | `description` | 视频当前没有 `content`，`description` 是纯文本。 |

IMPORTANT

视频投稿不要传 `content`。如果请求体里带了它，也不会变成视频正文；视频接口当前是来源链接、标题、分区、封面、标签和纯文本简介模型。

## 推荐传对象

接口层允许 `content` 是 `object` 或 `string`，但富文本正文推荐直接传 TipTap JSON 对象：

```
{
  "title": "示例帖子",
  "boardId": "general",
  "content": {
    "type": "doc",
    "content": [
      {
        "type": "paragraph",
        "content": [
          { "type": "text", "text": "这是一段正文。" }
        ]
      }
    ]
  }
}
```

## 字符串不是 Markdown

如果你这样传：

```
{
  "content": "# 标题\n\n- 列表 1\n- 列表 2"
}
```

服务端不会把它当 Markdown 解析。普通字符串会按纯文本兜底成段落。想要标题、列表、表格、图片、代码块等效果，请传 TipTap JSON。

| 入参 | 服务端理解 | 适合场景 |
| --- | --- | --- |
| TipTap JSON 对象 | 结构化富文本 | 推荐。适合正式投稿、同步脚本、复杂排版。 |
| 普通字符串 | 纯文本正文 | 只需要一段纯文本时可用。 |
| Markdown 字符串 | 仍是普通字符串 | 不推荐。API 不负责 Markdown 转 TipTap。 |
| BBCode 字符串 | 仍是普通字符串 | 不推荐。BBCode 兼容是前台编辑器导入能力。 |

前台 Markdown / BBCode 支持范围见[论坛正文语法支持](/docs/site-info/forum-markdown-bbcode/)。

## 最小可用文档

最小 TipTap 正文结构如下：

```
{
  "type": "doc",
  "content": [
    {
      "type": "paragraph",
      "content": [
        { "type": "text", "text": "你好，世界" }
      ]
    }
  ]
}
```

推荐所有正文都使用：

| 层级 | 字段 | 必填 | 说明 |
| --- | --- | --- | --- |
| 根节点 | `type` | 是 | 固定为 `doc`。 |
| 根节点 | `content` | 是 | 子节点数组，可以包含段落、标题、列表、图片等块级节点。 |
| 文本节点 | `type` | 是 | 固定为 `text`。 |
| 文本节点 | `text` | 是 | 实际文字。 |
| 文本节点 | `marks` | 否 | 粗体、链接、行内代码等标记数组。 |

## 正文共同规则

资源、帖子、找服玩三类投稿的 `content` 会参与这些服务端逻辑：

| 逻辑 | 影响 |
| --- | --- |
| 正文字数校验 | 正式投稿会从 TipTap JSON 抽纯文本计算长度。 |
| 摘要提取 | 摘要优先从结构化正文抽取。 |
| 图集提取 | 图片节点会影响首图或图库候选。 |
| 提及通知 | `mention` 节点会被识别为 @提及。 |
| 渲染一致性 | 结构化 JSON 比字符串更不容易出现迁移差异。 |

## 下一步读什么

| 你要做什么 | 继续阅读 |
| --- | --- |
| 查所有常用节点和属性 | [TipTap 节点与属性速查](/docs/api/tiptap-content-reference/) |
| 复制段落、列表、图片、代码块等 JSON 示例 | [TipTap JSON 示例与常见问题](/docs/api/tiptap-content-examples/) |
| 创建资源正文 | [个人 API：资源](/docs/api/personal-api-resources/) |
| 创建帖子正文 | [个人 API：帖子](/docs/api/personal-api-posts/) |
| 创建找服玩正文 | [个人 API：找服玩](/docs/api/personal-api-players/) |
