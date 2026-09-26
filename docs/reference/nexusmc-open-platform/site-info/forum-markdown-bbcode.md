---
title: "论坛正文语法支持"
source: https://docs.nexusmc.cn/docs/site-info/forum-markdown-bbcode
collected: 2026-09-26
---

# 论坛正文语法支持

说明论坛 Markdown 编辑器支持的常用 Markdown、兼容 BBCode 范围，以及从 Discuz、XenForo 等论坛迁移正文时的注意事项。

## 基本原则

论坛正文默认以 Markdown 为主。你可以直接写常见 Markdown，也可以在 Markdown 编辑器中混用一部分 BBCode，系统会在解析正文前把受支持的 BBCode 转成等价的 Markdown 或 TipTap 内容结构。

这个兼容层主要服务于历史论坛内容迁移，例如从 Discuz、XenForo 或其他 BBCode 论坛复制帖子到 NexusMC。它不是完整的 Discuz / XenForo 渲染器，因此建议把迁移后的内容再预览一遍，确认图片、链接、引用和代码块效果正确。

TIP

新内容优先使用 Markdown；迁移旧内容时，可以先粘贴 BBCode，再逐步清理不兼容的标签。

## 常用 Markdown

Markdown 编辑器支持日常发帖所需的基础写法，包括段落、标题、强调、列表、引用、链接、图片和代码块。

| 用途 | 写法示例 |
| --- | --- |
| 标题 | `## 二级标题` |
| 加粗 | `**重要内容**` |
| 斜体 | `*强调内容*` |
| 删除线 | `~~已废弃内容~~` |
| 链接 | `[NexusMC](https://example.com)` |
| 图片 | `![](https://example.com/image.png)` |
| 引用 | `> 引用内容` |
| 无序列表 | `- 列表项` |
| 有序列表 | `1. 列表项` |
| 行内代码 | `` `server.properties` `` |

代码块可以使用围栏语法，并标注语言：

```
```js
const name = 'NexusMC';
console.log(name);
```
```

## 代码块增强

论坛代码块在 Markdown 编辑器和高级 TipTap 编辑器中使用同一套扩展。行号默认关闭，需要时在语言后添加 `:line-numbers`。

| 功能 | 写法 |
| --- | --- |
| 高亮指定行 | ` ```js{1,4-6} ` |
| 开启行号 | ` ```ts:line-numbers ` |
| 指定起始行号 | ` ```ts:line-numbers=2 ` |
| diff 代码并指定内部语言 | ` ```diff lang="js" ` |

示例：

```
```ts:line-numbers=2 {1,3}
export const site = 'NexusMC';
console.log(site);
export default site;
```
```

diff 代码块可以继续标注实际代码语言，用于让新增行和删除行内部也按对应语言高亮：

```
```diff lang="js":line-numbers {1,3}
+const enabled = true;
 console.log(enabled);
-const enabled = false;
```
```

## 支持的 BBCode

以下 BBCode 可以直接粘贴到 Markdown 编辑器中。系统会尽量转换为 Markdown 或站内扩展能识别的结构。

| BBCode | 说明 | 转换结果 |
| --- | --- | --- |
| `[b]文本[/b]` | 加粗 | Markdown 加粗 |
| `[i]文本[/i]` | 斜体 | Markdown 斜体 |
| `[s]文本[/s]` / `[del]文本[/del]` | 删除线 | Markdown 删除线 |
| `[u]文本[/u]` | 下划线 | 站内下划线标记 |
| `[sub]2[/sub]` | 下标 | 站内下标标记 |
| `[sup]2[/sup]` | 上标 | 站内上标标记 |
| `[color=red]文本[/color]` | 字色 | 站内文本样式 |
| `[size=3]文本[/size]` | 字号 | 转为安全的字号值 |
| `[font=Arial]文本[/font]` | 字体 | 站内文本样式 |
| `[highlight]文本[/highlight]` | 高亮 | 站内高亮标记 |
| `[bgcolor=yellow]文本[/bgcolor]` / `[backcolor=yellow]文本[/backcolor]` | 背景色 | 转为高亮标记 |
| `[url]https://example.com[/url]` | 链接 | Markdown 链接 |
| `[url=https://example.com]文本[/url]` | 带标题链接 | Markdown 链接 |
| `[email]admin@example.com[/email]` | 邮箱链接 | `mailto:` 链接 |
| `[email=admin@example.com]联系管理员[/email]` | 带标题邮箱链接 | `mailto:` 链接 |
| `[img]https://example.com/a.png[/img]` | 图片 | Markdown 图片 |
| `[img=120,80]...[/img]` | 带宽高图片 | Markdown 图片附带尺寸 |
| `[grid]...图片...[/grid]` | Discourse 图片网格 | 站内图片轮播块 |
| `[quote]内容[/quote]` | 引用 | Markdown 引用块 |
| `[quote=Alice]内容[/quote]` | 带作者引用 | 引用块并保留作者说明 |
| `[spoiler]内容[/spoiler]` / `[ispoiler]内容[/ispoiler]` | 剧透 | 行内剧透或折叠块 |
| `[hide]内容[/hide]` 等隐藏类标签 | 旧论坛权限隐藏 | 折叠占位块，并提示需要手动设置权限 |
| `[list][*]A[*]B[/list]` | 无序列表 | Markdown 列表 |
| `[list=1][*]A[*]B[/list]` | 有序列表 | Markdown 有序列表 |
| `[code]...[/code]` | 代码块 | Markdown 围栏代码块 |
| `[code=js]...[/code]` | 带语言代码块 | Markdown 围栏代码块 |
| `[icode]server.properties[/icode]` | 行内代码 | Markdown 行内代码 |
| `[pre]...[/pre]` | 预格式文本 | Markdown 围栏代码块 |
| `[plain]...[/plain]` / `[noparse]...[/noparse]` | 不解析内容 | 转义为普通文本 |
| `[h2]标题[/h2]` / `[heading=2]标题[/heading]` | 标题 | Markdown 标题 |
| `[hr]` | 分隔线 | Markdown 分隔线 |
| `[center]内容[/center]` | 居中 | 站内对齐语法 |
| `[left]内容[/left]` / `[right]内容[/right]` / `[justify]内容[/justify]` | 对齐 | 站内对齐语法 |
| `[p=30,2,left]内容[/p]` | Discuz 段落参数 | 尽量保留对齐，缩进参数会降级 |
| `[fly]内容[/fly]` / `[float=left]内容[/float]` | 跑马灯/浮动 | 保留纯文本内容 |
| `[indent]内容[/indent]` | 缩进 | 引用块 |
| `[table][tr][th]A[/th][/tr][/table]` | 表格 | Markdown 表格 |
| `[media=bilibili]BV...[/media]` / `[media=x,500,375]bili:BV...[/media]` / `[bilibili]BV...[/bilibili]` | Bilibili 视频 | 站内 Bilibili 视频块 |
| `[media=youtube]...[/media]` / `[youtube]...[/youtube]` | YouTube 视频 | 外链 |
| `[video]https://...[/video]` / `[audio]https://...[/audio]` / `[Dplayer]...[/Dplayer]` / `[html5video]...[/html5video]` / `[flash]...[/flash]` | 音视频链接 | 外链 |
| `[xigua]...[/xigua]` | 西瓜视频 | 外链或普通占位文本 |
| `[wyy]歌曲 ID[/wyy]` | 网易云音乐 | 外链 |
| `[user]Alice[/user]` | 用户提及 | 普通 `@用户名` 文本 |
| `[attach]123[/attach]` | 附件占位 | 普通附件文本 |
| `[qq]123456[/qq]` | QQ 联系方式 | 普通文本 |
| `[postbg]...[/postbg]` | 旧站帖子背景 | 引用占位 |
| `[password]...[/password]` | 旧站密码帖 | 权限占位提示 |
| `[page]` / `[index]...[/index]` / `[md]...[/md]` | Discuz 分页/索引/Markdown 包裹 | 分隔线、列表或直接保留 Markdown |
| `[toanchor=id]跳转[/toanchor]` | 锚点跳转 | Markdown 锚点链接 |
| `[anchor=id]标题[/anchor]` | 锚点位置 | 加粗文本占位 |

示例：

```
[quote=Alice]
[b]旧论坛公告[/b]
请查看 [url=https://example.com]项目主页[/url]。
[/quote]

[code=js]
console.log('[b]这里不会被当作加粗[/b]');
[/code]
```

CAUTION

`[hide]`、`[hidereply]`、`[hidethanks]` 等旧论坛隐藏标签只会转成“隐藏内容迁移占位”折叠块。它们不会自动继承旧站的回复可见、积分可见或付费可见规则。需要真实权限控制时，请使用站内发布表单的可见性设置。

## HTML 兼容

从旧论坛复制内容时，正文中可能夹带少量 HTML。当前会处理一部分常见、低风险的表现型 HTML：

| HTML | 处理方式 |
| --- | --- |
| `<br>` | 转为换行 |
| `<u>` | 转为 `[u]` |
| `<sub>` | 转为 `[sub]` |
| `<sup>` | 转为 `[sup]` |
| `<center>` | 转为居中对齐 |
| 带 `align` 或 `text-align` 的 `p` / `div` | 转为站内对齐语法 |

其它复杂 HTML、脚本、样式或论坛私有标签不应依赖自动兼容。迁移后如果出现未识别文本，建议手动改写成 Markdown 或站内编辑器提供的内容块。

## 迁移建议

1.  先把旧帖原文粘贴到 Markdown 编辑器。
2.  使用预览确认标题、引用、列表、图片、链接和代码块。
3.  对未识别的论坛私有标签手动改写。
4.  对长期维护的教程、资源说明和公告，尽量整理为纯 Markdown。

CAUTION

通过个人 API 提交正文时，`content` 仍然建议传 TipTap JSON 结构。普通字符串不会被服务端当作完整 Markdown 或 BBCode 自动渲染。API 正文格式见 [投稿正文 TipTap 快速接入](/docs/api/tiptap-content-format/)。
