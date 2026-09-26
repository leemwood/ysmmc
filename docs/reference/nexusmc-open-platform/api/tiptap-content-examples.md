---
title: "TipTap JSON 示例与常见问题"
source: https://docs.nexusmc.cn/docs/api/tiptap-content-examples
collected: 2026-09-26
---

# TipTap JSON 示例与常见问题

提供个人 API 正文常用 JSON 片段、完整正文示例，以及 content 字符串、BBCode、视频字段等常见问题。

## 段落

```
{
  "type": "paragraph",
  "content": [
    { "type": "text", "text": "这里是一段普通正文。" }
  ]
}
```

## 标题

```
{
  "type": "heading",
  "attrs": { "level": 2 },
  "content": [
    { "type": "text", "text": "安装说明" }
  ]
}
```

## 粗体、链接和行内代码

```
{
  "type": "paragraph",
  "content": [
    { "type": "text", "text": "请先阅读 " },
    {
      "type": "text",
      "text": "服务端配置文档",
      "marks": [
        {
          "type": "link",
          "attrs": { "href": "https://example.com/docs" }
        }
      ]
    },
    { "type": "text", "text": "，其中 " },
    {
      "type": "text",
      "text": "version",
      "marks": [{ "type": "code" }]
    },
    { "type": "text", "text": " 字段必须填写。" }
  ]
}
```

## 无序列表

```
{
  "type": "bulletList",
  "content": [
    {
      "type": "listItem",
      "content": [
        {
          "type": "paragraph",
          "content": [{ "type": "text", "text": "支持 Fabric" }]
        }
      ]
    }
  ]
}
```

## 图片

```
{
  "type": "image",
  "attrs": {
    "src": "https://www.nexusmc.cn/uploads/example.png",
    "alt": "功能截图",
    "title": "功能截图"
  }
}
```

## 完整正文示例

```
{
  "type": "doc",
  "content": [
    {
      "type": "heading",
      "attrs": { "level": 2 },
      "content": [{ "type": "text", "text": "功能介绍" }]
    },
    {
      "type": "paragraph",
      "content": [
        { "type": "text", "text": "这是一个适用于 1.20.1 的资源示例，支持多人联机与自定义配置。" }
      ]
    },
    {
      "type": "image",
      "attrs": {
        "src": "https://www.nexusmc.cn/uploads/example-preview.png",
        "alt": "预览图",
        "title": "预览图"
      }
    }
  ]
}
```

## 请求示例

### 资源投稿

```
{
  "title": "示例资源",
  "description": "用于演示 TipTap 正文结构。",
  "platform": "java",
  "category": "resource-pack",
  "content": {
    "type": "doc",
    "content": [
      {
        "type": "heading",
        "attrs": { "level": 2 },
        "content": [{ "type": "text", "text": "安装方法" }]
      },
      {
        "type": "paragraph",
        "content": [{ "type": "text", "text": "将压缩包放入 resourcepacks 文件夹后重启游戏。" }]
      }
    ]
  },
  "downloadType": "external",
  "fileUrl": "https://example.com/download.zip"
}
```

### 帖子投稿

```
{
  "title": "示例帖子",
  "boardId": "general",
  "content": {
    "type": "doc",
    "content": [
      {
        "type": "paragraph",
        "content": [{ "type": "text", "text": "大家好，这是通过 API 发布的一篇富文本帖子。" }]
      }
    ]
  }
}
```

### 找服玩投稿

```
{
  "title": "示例生存服",
  "description": "长期开放，偏原版生存体验。",
  "category": "survival",
  "ip": "play.example.com",
  "onlineMode": true,
  "networkEnvironments": ["mainland"],
  "content": {
    "type": "doc",
    "content": [
      {
        "type": "heading",
        "attrs": { "level": 2 },
        "content": [{ "type": "text", "text": "服务器特色" }]
      }
    ]
  }
}
```

## 常见问题

### `content` 传字符串为什么没有标题和列表效果？

因为普通字符串不会被当 Markdown 渲染。要想出现标题、列表、表格等效果，请传 TipTap JSON。

### API 可以直接传 BBCode 吗？

不建议。BBCode 兼容主要是给前台 Markdown 编辑器迁移旧帖用的导入能力。通过 API 同步旧论坛数据时，最稳妥的方式是先把旧内容转换成 TipTap JSON。

### `content` 可以传 JSON 字符串吗？

可以，接口层接受 `string` 或 `object`。但推荐传对象，避免转义问题。

### 为什么正式投稿时提示正文太短？

服务端会从 TipTap JSON 里提取纯文本再计算字数。空节点、空表格、空图片不会让正文长度变长。

### 我的视频投稿为什么不能传 `content`？

因为当前视频投稿接口没有这个字段。视频区目前是视频源链接、标题和纯文本简介模型，不是富文本正文模型。
