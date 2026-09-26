---
title: "OneBot 11/12：接口与请求"
source: https://docs.nexusmc.cn/docs/api/open-platform-onebot-reference
collected: 2026-09-26
---

# OneBot 11/12：接口与请求

分版本列出当前网关支持的动作、JSON 请求与响应、事件帧和站点扩展。

本文是本站 OneBot 网关**当前已实现的接口**，分别给出 OneBot 11 和 OneBot 12 的调用方式。先阅读[接入概念与机器人绑定](/docs/api/open-platform-bots/)以及[正反向 WebSocket 连接](/docs/api/realtime-websocket/)。动作通过正向或反向 OneBot WebSocket 发送 JSON 文本帧；[OneBot HTTP WebHook](/docs/api/open-platform-webhooks/)只接收事件，不提供 HTTP 动作调用入口。

## 通用请求与响应

两种版本都接受 `{ "action": "动作名", "params": {}, "echo": "请求标识" }`。`echo` 可省略；提供后响应原样返回，用于匹配并发请求。响应格式：

```
{
  "status": "ok",
  "retcode": 0,
  "data": { "message_id": "MESSAGE_ID" },
  "message": "",
  "wording": "",
  "echo": "send-1"
}
```

失败时 `status` 为 `failed`、`data` 为 `null`，`message` / `wording` 带错误码；扩展错误详情可能在 `nexusmc_error`。例如内容进入人工审核时，返回 `MESSAGE_REVIEW_REQUIRED`，`nexusmc_error.messageId` 保留站内消息 ID。没有绑定机器人身份的连接只能调用 `get_login_info` / `get_self_info`、`get_status`、`get_version_info` 和 `ping`；发消息等动作返回身份错误。

## OneBot 11

本站 OneBot 11 的用户 `user_id` 与 `self_id` 为数字 UID，群 `group_id` 和消息 `message_id` 是站内字符串 ID。只接受纯数字 ID 的客户端需自行建立群和消息 ID 映射。

### 动作

| 动作 | 关键 `params` | 结果与限制 |
| --- | --- | --- |
| `get_login_info` | `{}` | 当前机器人数字 `user_id`、名称；未绑定时为站点数据身份。 |
| `get_status`、`get_version_info`、`ping` | `{}` | 连接状态、网关版本或保活响应。 |
| `send_private_msg` | `user_id`、`message`；或 `nexusmc_conversation_id` | 给用户发私信；需 `bot:dm:write` 并遵守对方隐私设置。 |
| `send_group_msg` | `group_id`、`message` | 向机器人已加入且未禁言的群发消息；需 `bot:chat:write`。 |
| `send_msg` | `message_type: "private"` + `user_id`，或 `message_type: "group"` + `group_id`；以及 `message` | 按类型发送。 |
| `delete_msg` | `message_id` | 撤回机器人自己的消息；需在撤回时限内。 |
| `get_msg` | `message_id` | 读取可访问的群或私信消息。 |
| `get_group_msg_history` | `group_id`，`count?` | 返回 `data.messages`；需群成员身份和读取权限。 |

群消息请求示例：

```
{
  "action": "send_group_msg",
  "params": { "group_id": "ROOM_ID", "message": "机器人上线了" },
  "echo": "group-1"
}
```

私信请求示例（`user_id` 是数字 UID）：

```
{
  "action": "send_private_msg",
  "params": { "user_id": 12345, "message": "你好" },
  "echo": "private-1"
}
```

### 事件

连接建立后首先收到 `post_type: "meta_event"`、`meta_event_type: "lifecycle"`、`sub_type: "connect"`。订阅 `message.group` 的群消息示例：

```
{
  "time": 1780000000,
  "self_id": -12345,
  "post_type": "message",
  "message_type": "group",
  "group_id": "ROOM_ID",
  "message_id": "MESSAGE_ID",
  "user_id": 67890,
  "message": "你好",
  "raw_message": "你好",
  "sender": { "user_id": 67890, "nickname": "用户" },
  "nexusmc_content": "你好"
}
```

私信事件的 `message_type` 为 `private`，增加 `nexusmc_conversation_id`，没有 `group_id`。公开内容、编辑、撤回和互动通知为 `post_type: "notice"`、`notice_type: "nexusmc"`；真实事件名在 `sub_type` 与 `nexusmc_event_type`。

## OneBot 12

本站 OneBot 12 的 `self.user_id` 和用户 `user_id` 是站内字符串用户 ID，`self.platform` 固定为 `nexusmc`。群和消息仍使用站内字符串 ID。文本消息使用文本段数组；当前发送动作只处理文本段，互动按钮对象请使用 `nexusmc_content` 扩展字段。

### 动作

| 动作 | 关键 `params` | 结果与限制 |
| --- | --- | --- |
| `get_self_info` | `{}` | 当前机器人字符串 `user_id`、名称；未绑定时为站点数据身份。 |
| `get_status`、`get_version_info`、`ping` | `{}` | 连接状态、网关版本或保活响应。 |
| `send_message` | `detail_type: "private"` + `user_id`，或 `detail_type: "group"` + `group_id`；以及 `message` | 发送私信或群消息，受相应 scope、隐私与成员规则约束。 |
| `delete_message` | `message_id` | 撤回机器人自己的消息。 |
| `get_message` | `message_id` | 读取可访问的群或私信消息。 |
| `get_group_message_history` | `group_id`，`count?` | 返回 `data.messages`。 |

群消息请求示例：

```
{
  "action": "send_message",
  "params": {
    "detail_type": "group",
    "group_id": "ROOM_ID",
    "message": [{ "type": "text", "data": { "text": "机器人上线了" } }]
  },
  "echo": "group-1"
}
```

私信请求示例（`user_id` 是字符串用户 ID）：

```
{
  "action": "send_message",
  "params": {
    "detail_type": "private",
    "user_id": "USER_ID",
    "message": [{ "type": "text", "data": { "text": "你好" } }]
  },
  "echo": "private-1"
}
```

### 事件

连接建立后首先收到 `type: "meta"`、`detail_type: "connect"`。群消息示例：

```
{
  "id": "EVENT_ID",
  "time": 1780000000,
  "type": "message",
  "detail_type": "group",
  "sub_type": "",
  "self": { "platform": "nexusmc", "user_id": "BOT_USER_ID" },
  "group_id": "ROOM_ID",
  "message_id": "MESSAGE_ID",
  "user_id": "SENDER_USER_ID",
  "message": [{ "type": "text", "data": { "text": "你好" } }],
  "alt_message": "你好",
  "nexusmc_content": "你好"
}
```

私信事件的 `detail_type` 为 `private`，增加 `nexusmc_conversation_id`。公开内容、编辑、撤回和互动通知为 `type: "notice"`、`detail_type: "nexusmc"`，具体事件名在 `sub_type` 与 `nexusmc_event_type`。

## 共同的事件订阅与站点扩展

两种版本的可选公开事件均为 `resource.created`、`resource.updated`、`resource.comment.created`、`post.created`、`post.updated`、`post.comment.created`、`user.activity.created`、`user.resource.created`、`user.post.created`、`user.video.created`。绑定机器人后还可选 `message.private`、`message.group`、`interaction.created`、`notification.created`、`bot.mention`。在连接的 `events` / `filters` 中配置订阅；无需向 OneBot WebSocket 额外发送订阅动作。`message.private`、`message.group` 的编辑与撤回事件沿用对应父事件的订阅，收到 `notice` 帧。过滤字段与示例见[Webhook 接口](/docs/api/open-platform-webhooks/)。

站点命令使用 `nexusmc.` 前缀，后缀为[原生 WebSocket 命令名](/docs/api/realtime-websocket-reference/)，`params` 与原生命令参数一致。例如申请入群与读取私信历史：

```
{ "action": "nexusmc.room.apply", "params": { "roomId": "ROOM_ID", "message": "机器人用途" }, "echo": "apply-1" }
```

```
{ "action": "nexusmc.conversation.message.read", "params": { "conversationId": "CONVERSATION_ID", "limit": 50 }, "echo": "history-1" }
```

发送互动按钮消息时，将符合[原生互动消息格式](/docs/api/realtime-websocket-reference/)的对象放入 `params.nexusmc_content`。当前网关还接受部分跨版本别名，例如 `get_msg` / `get_message`、`delete_msg` / `delete_message`；新客户端建议使用本页对应版本的动作名。上述动作仍受机器人 scope、群成员资格、私信隐私、内容扫描和消息时限限制；网关未实现的标准动作返回 `ACTION_UNSUPPORTED`。
