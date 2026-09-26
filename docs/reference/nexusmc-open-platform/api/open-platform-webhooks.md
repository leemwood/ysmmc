---
title: "Webhook：概念与全部接口"
source: https://docs.nexusmc.cn/docs/api/open-platform-webhooks
collected: 2026-09-26
---

# Webhook：概念与全部接口

站点 Webhook 与 OneBot HTTP WebHook 的创建、事件筛选、投递格式和验签接口。

Webhook 是本站主动向你的 HTTP 接口发送事件，无需持续保持 WebSocket 连接。它是**事件通知**，不是用于读取历史数据的 HTTP 查询接口。成功接收后如需详情，可再调用对应的公开或已授权 API。

[开放平台控制台](https://www.nexusmc.cn/open-platform/console)的“Webhook”页有两种产品：站点 Webhook 使用本站事件 JSON，不绑定机器人；OneBot HTTP WebHook 是一种 OneBot 连接，可选绑定机器人并使用 OneBot 11/12 事件帧。两者管理接口、密钥和请求体不同，接收端应分别验签和解析。

## 开通与创建

登录账号、验证邮箱、通过开发者考试并同意当前版本的开放平台协议后，才能管理 Webhook。此处使用与 WebSocket 相同的准入条件，不要求 OAuth 应用或站点 API Key。管理接口使用**登录会话**；个人 API Token、站点 API Key 和 WebSocket Key 不能代替登录会话管理 Webhook。

在控制台填写名称、接收地址、至少一种事件和可选过滤条件即可创建。接收地址必须是可从本站访问的安全 HTTP(S) 地址；内网、回环等不安全目标会被拒绝。创建成功只显示一次完整密钥，请立即保存在接收端。

站点 Webhook 管理 API 的前缀为 `/api/open-platform`：

| 方法 | 路径 | 用途 |
| --- | --- | --- |
| `GET` | `/webhooks` | 列出自己创建或被指定管理的 Webhook，并返回可选事件。 |
| `POST` | `/webhooks` | 创建，返回 `webhook` 和一次性 `secret`。 |
| `PUT` | `/webhooks/:id` | 修改名称、地址、事件、过滤条件或启用状态。请求体须包含非空 `events`。 |
| `POST` | `/webhooks/:id/rotate-secret` | 轮换密钥，返回新的 `secret`；旧密钥随即失效。 |
| `DELETE` | `/webhooks/:id` | 删除。 |

例如订阅某个资源的更新与评论，创建请求体为：

```
{
  "name": "资源通知",
  "url": "https://example.com/nexusmc/events",
  "events": ["resource.updated", "resource.comment.created"],
  "filters": { "resourceId": "RESOURCE_ID" }
}
```

`POST /api/open-platform/webhooks` 返回 `201`；`secret` 只出现在创建和轮换响应中。编辑后列表只返回 `secretPrefix`。后台也可以指定一个站内用户为管理者，该用户可在自己的控制台查看和管理该 Webhook。

## OneBot HTTP WebHook 管理接口

在“Webhook”页的“OneBot HTTP WebHook”区域创建连接，或使用下面的 API。它们与站点 Webhook 一样需要登录会话和开放平台 WebSocket 准入条件；不能用事件接收密钥管理连接。

| 方法 | 完整路径 | 用途 |
| --- | --- | --- |
| `GET` | `/api/open-platform/bots/connections` | 列出当前账号可管理的全部 OneBot 连接；从 `transport` 区分 HTTP 与 WebSocket。 |
| `POST` | `/api/open-platform/bots/connections` | 创建，`transport` 填 `http_webhook`，返回一次性 `secret`。 |
| `PATCH` | `/api/open-platform/bots/connections/:connectionId` | 修改目标地址、事件、过滤条件、机器人绑定或启用状态。 |
| `POST` | `/api/open-platform/bots/connections/:connectionId/rotate-secret` | 轮换该连接的签名密钥。 |
| `DELETE` | `/api/open-platform/bots/connections/:connectionId` | 删除连接。 |

创建请求体示例：

```
{
  "name": "OneBot 事件接收端",
  "protocolVersion": "onebot12",
  "transport": "http_webhook",
  "endpoint": "https://example.com/onebot/events",
  "botAppId": "BOT_ID",
  "events": ["message.group", "resource.updated"],
  "filters": {}
}
```

`botAppId` 可省略或为 `null`，此时只能选公开事件。绑定机器人后，才可订阅群聊、私信、互动和通知等事件，并受机器人的 scope 与站内访问权限约束。`protocolVersion` 可为 `onebot11` 或 `onebot12`。完整连接规则见[OneBot 11/12 接入概念](/docs/api/open-platform-bots/)。

## 事件与筛选

| 事件 | 含义 |
| --- | --- |
| `resource.created`、`resource.updated` | 公开资源发布、更新。 |
| `resource.comment.created` | 公开资源出现已通过审核的评论。 |
| `post.created`、`post.updated` | 公开帖子发布、更新。 |
| `post.comment.created` | 公开帖子出现已通过审核的评论。 |
| `user.activity.created` | 用户公开内容的更新动态；具体动作见 `data.action` 和 `data.kind`。 |
| `user.resource.created`、`user.post.created`、`user.video.created` | 指定用户发布了对应的公开内容。 |

`filters` 可填 `resourceId`、`postId`、`userId`、`authorId`。每个已填写的条件都必须与事件的同名顶层字段完全相等；不填则不限制。订阅某篇帖子的评论可选择 `post.comment.created` 并填 `postId`；订阅某人发布的资源可选择 `user.resource.created` 并填 `userId`。评论事件的 `authorId` 是**评论作者**，不是资源或帖子的作者；请用 `resourceId` 或 `postId` 锁定被评论对象。不同事件未必包含全部四个 ID，组合过滤时先确认事件字段。

事件在公开且通过审核后进入后台投递队列，不能把“投稿已提交”等同于 `*.created` 已投递。接收端可能因重试再次收到同一个事件，应以 `id` 去重。

## HTTP 投递与验签

接收端收到 `Content-Type: application/json`、`X-Nexus-Event`、`X-Nexus-Event-Id` 和 `X-Nexus-Signature`。请求体形如：

```
{
  "id": "resource:RESOURCE_ID:updated:2026-09-25T08:00:00.000Z",
  "type": "resource.updated",
  "occurredAt": "2026-09-25T08:00:00.000Z",
  "data": { "id": "RESOURCE_ID", "title": "示例资源" },
  "resourceId": "RESOURCE_ID",
  "userId": "USER_ID",
  "authorId": "USER_ID"
}
```

字段随事件变化；未提供的 ID 可能不出现在 JSON 中。签名为 `sha256=` 加上以完整 Webhook 密钥为 key、**原始请求体字节**为消息的 HMAC-SHA256 十六进制结果。先验签，再解析和处理 JSON；不要重新序列化 JSON 后验签。Node.js 示例：

```
import { createHmac, timingSafeEqual } from 'node:crypto';

function verifyWebhook(rawBody, signature, secret) {
  if (!String(signature || '').startsWith('sha256=')) return false;
  const expected = createHmac('sha256', secret).update(rawBody).digest('hex');
  const received = signature.slice('sha256='.length);
  const left = Buffer.from(expected, 'hex');
  const right = Buffer.from(received, 'hex');
  return left.length === right.length && timingSafeEqual(left, right);
}
```

接收端成功处理后返回任意 `2xx`。非 `2xx` 或超时会被记录为失败并尝试重投；单次请求的超时时间为 15 秒。请让接收接口快速确认，再异步处理业务，并按 `X-Nexus-Event-Id` 做幂等去重。需要 OneBot 11/12 消息帧、私信或群聊事件时，请查看[OneBot 接口与请求](/docs/api/open-platform-onebot-reference/)。

## OneBot HTTP 投递

OneBot HTTP WebHook 也向 `endpoint` 发送 JSON `POST`，头部含 `X-OneBot-Version: onebot11` 或 `onebot12`，以及 `X-Nexus-Event`、`X-Nexus-Event-Id`、`X-Nexus-Signature`。验签算法与上文相同，但 HMAC key 必须使用**该 OneBot 连接的密钥**。其 body 是 OneBot `message` 或 `notice` 事件帧，不是上方站点 Webhook 的 `{ id, type, occurredAt, data }` JSON。两种版本的事件示例见[OneBot 接口与请求](/docs/api/open-platform-onebot-reference/)。

HTTP WebHook 只能接收事件；动作请求必须通过绑定同一机器人的正向或反向 OneBot WebSocket 发送。发送失败会记录投递状态并重试，接收端同样应按 `X-Nexus-Event-Id` 去重。
