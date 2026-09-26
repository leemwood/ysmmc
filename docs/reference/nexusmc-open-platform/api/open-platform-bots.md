---
title: "OneBot 11/12：接入概念"
source: https://docs.nexusmc.cn/docs/api/open-platform-bots
collected: 2026-09-26
---

# OneBot 11/12：接入概念

机器人身份、协议版本、正反向 WebSocket、HTTP WebHook 与群聊准入。

机器人是站内独立账号；OneBot 连接是传输通道。一个机器人可以绑定多个连接，也可以创建**不绑定机器人**的 OneBot 连接，只接收公开站点事件。登录后在[开放平台控制台](https://www.nexusmc.cn/open-platform/console)的“WebSocket”页管理机器人和 WebSocket 连接，在“Webhook”页管理 OneBot HTTP WebHook。具体动作、请求参数和事件帧见[OneBot 接口与请求](/docs/api/open-platform-onebot-reference/)。

## 准入与身份

创建机器人或连接需要账号正常、邮箱已验证、通过开发者考试，并同意当前版本的开放平台协议。创建机器人时至少选择 `bot:realtime:connect`；群聊、私信和通知由相应 `bot:*` scope 控制。OneBot 连接的公开站点事件由连接的 `events` 和 `filters` 决定；原有 Socket.IO 内容主题另需 `bot:content:subscribe`。管理 API 使用登录会话，不接受个人 API Token 或站点 API Key 代替登录。

机器人创建后获得独立的站内用户 ID、数字 UID、名称和头像。WebSocket Key 用于[原有 Socket.IO `/realtime/v1` 接口](/docs/api/realtime-websocket/)；OneBot 连接有各自独立的一次性连接密钥。绑定机器人后，连接按机器人的 scope、群成员资格和会话参与资格执行动作与投递私有事件。连接不绑定机器人时，只能订阅公开站点事件和调用连接状态类动作。

## 管理接口

以下路径均以 `/api/open-platform/bots` 为前缀，使用登录会话：

| 方法 | 路径 | 用途 |
| --- | --- | --- |
| `GET` | `/` | 列出自己拥有或受委托管理的机器人与连接。 |
| `POST` | `/` | 创建机器人；传 `name`、`scopes`，可选 `description`、`avatar`、连接数和速率限制。 |
| `PATCH` | `/:id` | 修改机器人资料、scope、启用状态或限制。 |
| `POST` | `/:id/tokens` | 创建 WebSocket Key；完整 `token` 只返回一次。 |
| `DELETE` | `/:id/tokens/:tokenId` | 撤销 WebSocket Key，并断开使用它的实时连接。 |
| `GET` | `/connections` | 列出 OneBot 连接、协议、传输方式和可选事件。 |
| `POST` | `/connections` | 创建连接，返回一次性 `secret`。 |
| `PATCH` | `/connections/:connectionId` | 修改连接、绑定机器人、事件、过滤条件或启用状态。 |
| `POST` | `/connections/:connectionId/rotate-secret` | 轮换连接密钥，旧密钥失效。 |
| `DELETE` | `/connections/:connectionId` | 删除连接。 |

创建机器人示例请求体：

```
{
  "name": "更新提醒助手",
  "scopes": ["bot:realtime:connect", "bot:chat:subscribe", "bot:chat:read", "bot:chat:write", "bot:content:subscribe"]
}
```

创建连接示例请求体：

```
{
  "name": "生产事件连接",
  "botAppId": "BOT_ID",
  "protocolVersion": "onebot12",
  "transport": "forward_websocket",
  "events": ["resource.updated", "message.group"],
  "filters": {}
}
```

`botAppId` 可省略或设为 `null`。未绑定时 `events` 只能包含公开站点事件；绑定机器人后可另选 `message.private`、`message.group`、`interaction.created`、`notification.created`、`bot.mention`。公开事件及 `resourceId`、`postId`、`userId`、`authorId` 过滤语义见[Webhook 事件说明](/docs/api/open-platform-webhooks/)。可选 `protocolVersion` 为 `onebot11` 或 `onebot12`，默认 `onebot12`。后台可指定机器人或连接的站内管理用户；完整连接密钥仅在创建和轮换时返回。

## 连接方式

| `transport` | 连接方向 | 配置与鉴权 |
| --- | --- | --- |
| `forward_websocket` | 你的客户端连接本站 | 创建后使用返回的 `forwardUrl`，即 `/onebot/v1/connections/<连接 ID>`；用 `Authorization: Bearer <连接密钥>` 握手。 |
| `reverse_websocket` | 本站 worker 连接你的服务 | `endpoint` 填 `ws://` 或 `wss://` 地址；本站在握手时发送 `Authorization: Bearer <连接密钥>`。 |
| `http_webhook` | 本站 worker 向你的服务推送 | `endpoint` 填 `http://` 或 `https://` 地址；只投递事件，不在 HTTP 响应中执行 OneBot 动作。 |

反向 WebSocket 和 HTTP WebHook 的目标地址会经过安全网络目标检查；本机、内网等不安全地址不能作为公开平台的接收端。HTTP WebHook 与[站点 Webhook](/docs/api/open-platform-webhooks/)共用 `X-Nexus-Event`、`X-Nexus-Event-Id`、`X-Nexus-Signature` 头；另外发送 `X-OneBot-Version`。签名仍是对原始请求体使用**该连接自己的密钥**计算 HMAC-SHA256。HTTP 成功返回 `2xx`；失败或超时会记录并重试，单次超时 15 秒。

机器人动作通过本站实时服务执行。反向连接的 worker 默认按本机 `PORT` 连接 `/realtime/v1`；如果 worker 与 API 不在同一主机，部署时设置 `ONEBOT_BRIDGE_URL` 为 API 的内网 HTTP 地址。

## OneBot 版本与边界

连接创建时固定选择 `onebot11` 或 `onebot12`。正向和反向 WebSocket 建立后，先收到相应版本的生命周期事件，再双向收发 JSON；HTTP WebHook 只接收事件，不能通过 HTTP 响应提交动作。需要执行动作时，应创建绑定同一机器人的 WebSocket 连接。

OneBot 11 使用数字 `self_id` / `user_id`，OneBot 12 使用 `self.user_id` 与字符串用户 ID。本站群聊和消息 ID 仍是字符串。事件选择由连接的 `events` 和 `filters` 决定；绑定机器人后私有事件还须满足机器人的 scope、群成员或会话参与权限。当前网关实现的是[接口参考列出的动作和事件](/docs/api/open-platform-onebot-reference/)，不是 OneBot 11 或 12 标准的全部能力。

## 申请加入群聊

开发者可在控制台填写群聊 ID 和申请说明，也可调用 `POST /api/open-platform/bots/:id/room-applications`，请求体为 `{ "roomId": "ROOM_ID", "message": "机器人用途" }`。`GET /api/open-platform/bots/room-applications` 可查看申请状态。群主或群管理通过 `GET /api/messages/rooms/:id/bot-applications` 查看待处理申请，并调用 `POST /api/messages/rooms/:id/bot-applications/:applicationId/review`，请求体 `{ "action": "approve" }` 或 `{ "action": "reject" }`。

批准后机器人以普通成员身份加入。群主可在群成员管理中将其设为管理员；机器人还必须具备 `bot:chat:moderate`，才可执行成员角色、移除和禁言等管理命令。读取群消息需要群成员身份和 `bot:chat:read`；私信需要机器人作为会话参与者和对应的 `bot:dm:*` scope。完整的原有命令、主题与错误码见[WebSocket 接口参考](/docs/api/realtime-websocket-reference/)。
