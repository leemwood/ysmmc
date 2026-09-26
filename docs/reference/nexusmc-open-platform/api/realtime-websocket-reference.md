---
title: "WebSocket：全部命令与事件"
source: https://docs.nexusmc.cn/docs/api/realtime-websocket-reference
collected: 2026-09-26
---

# WebSocket：全部命令与事件

原生 Socket.IO 实时接口的命令参数、scope、主题、事件和错误响应。

本文仅描述原生 Socket.IO `/realtime/v1`。正向、反向连接和凭据选择先看[WebSocket 连接](/docs/api/realtime-websocket/)；OneBot JSON 动作见[OneBot 接口与请求](/docs/api/open-platform-onebot-reference/)。所有服务端事件均为统一信封：

```
{
  "id": "evt_xxx",
  "type": "message.created",
  "ts": "2026-09-04T08:30:00.000Z",
  "requestId": "req-1",
  "payload": {}
}
```

## 权限

WebSocket Key 必须包含 `bot:realtime:connect`，可选权限为 `bot:chat:subscribe`、`bot:chat:read`、`bot:chat:write`、`bot:chat:moderate`、`bot:dm:read`、`bot:dm:write`、`bot:notification:read` 和 `bot:content:subscribe`。连接认证与 `ready` 字段见[WebSocket 连接](/docs/api/realtime-websocket/)。

## 全部客户端命令

使用 `socket.emit(命令名, 参数)` 发送。`requestId` 可附在命令参数中，用于关联响应；仅部分写入命令会在同一连接中缓存该 ID 的结果。其他命令不能仅凭 `requestId` 假定幂等，调用方应按响应和业务对象核对重试结果。

| 命令 | 必要参数 / 主要可选参数 | 所需 scope | 响应 |
| --- | --- | --- | --- |
| `ping` | `requestId?` | 连接权限 | `pong` |
| `subscribe`、`unsubscribe` | `topic` | `bot:chat:subscribe`；特定主题另有权限 | `command.accepted` |
| `ack` | `cursor`，非负安全整数 | 连接权限 | `command.accepted` |
| `resume` | `cursor`，`limit?` | 连接权限 | 先重放事件，再返回 `command.accepted` |
| `room.apply` | `roomId`，`message?` | `bot:chat:write`、`bot:chat:subscribe` | `command.accepted`，含申请 ID |
| `room.join` | `inviteCode` | `bot:chat:write`、`bot:chat:subscribe` | `command.accepted`，含群聊 ID |
| `room.leave` | `roomId` | `bot:chat:write` | `command.accepted` |
| `room.member.role` | `roomId`、`userId`、`role`（`admin` 或 `member`） | `bot:chat:moderate`，且机器人为群管理 | `command.accepted` |
| `room.member.remove` | `roomId`、`userId` | 同上 | `command.accepted` |
| `room.member.mute` | `roomId`、`userId`，`durationMs?`、`reason?` | 同上 | `command.accepted` |
| `room.member.unmute` | `roomId`、`userId` | 同上 | `command.accepted` |
| `message.send` | `roomId`、`content` | `bot:chat:write`，且为未禁言群成员 | `command.accepted` 或待审 `command.rejected` |
| `message.read` | `roomId`，`limit?`；或 `messageId` 标记已读 | `bot:chat:read`，且为群成员 | `message.batch` 或 `command.accepted` |
| `message.edit` | `messageId`、`content` | `bot:chat:write`，且为发送者 | `command.accepted` |
| `message.recall` | `messageId` | `bot:chat:write`，且为发送者 | `command.accepted` |
| `conversation.message.send` | `conversationId`、`content` | `bot:dm:write`，且为会话参与者 | `command.accepted` 或待审 `command.rejected` |
| `conversation.message.read` | `conversationId`，`limit?` | `bot:dm:read`，且为会话参与者 | `conversation.message.batch` |
| `conversation.message.edit` | `messageId`、`content` | `bot:dm:write`，且为发送者 | `command.accepted` |
| `conversation.message.recall` | `messageId` | `bot:dm:write`，且为发送者 | `command.accepted` |
| `typing.start`、`typing.stop` | `roomId` | `bot:chat:write`，且为群成员 | `command.accepted` |

`content` 可为非空文本或下文所述的互动消息对象。发送、编辑和撤回仍经过站内隐私、成员、禁言与审核规则；编辑限发送后 2 分钟，撤回限 3 分钟。群消息读取与私信读取的 `limit` 会限定在 1 到 100。

## 主题订阅

```
socket.emit('subscribe', { requestId: 'sub-1', topic: 'room:ROOM_ID' });
socket.emit('unsubscribe', { requestId: 'sub-2', topic: 'room:ROOM_ID' });
```

| 主题 | 用途 | 额外限制 |
| --- | --- | --- |
| `room:<id>` | 群聊消息和成员事件 | 机器人必须是群成员 |
| `conversation:<id>` | 私聊消息事件 | 需要 `bot:dm:read`，且机器人必须是会话参与者 |
| `user:<id>` | 用户在线状态 | 需要 `bot:chat:read`，且只能订阅机器人所属开发者或机器人自身 |
| `notification:<id>` | 通知事件 | 需要 `bot:notification:read`，且只能访问所属开发者或机器人自身 |
| `content:resource`、`content:post`、`content:player`、`content:video`、`content:all` | 公开内容变化 | 需要 `bot:content:subscribe` |

## 断线恢复

`ready` 会返回 `latestCursor`。机器人处理事件后发送确认，重连并完成主题订阅后使用 `resume` 补拉遗漏事件：

```
socket.emit('ack', { cursor: 1024 });
socket.emit('resume', { cursor: 1024, limit: 200 });
```

事件只会按当前连接已订阅的主题回放。事件存储保留最近 5000 条，游标超出保留窗口时返回 `CURSOR_EXPIRED`，调用方应重新建立全量同步基线。

## 群聊命令

开发者可在开放平台使用群聊 ID 为自己的机器人提交入群申请，群主或群管理在群聊设置中审核。也可从机器人连接发送 `room.apply`；申请通过后机器人先以普通成员身份加入。只有群主能在成员列表授予机器人群管理角色，且机器人自身还需具备 `bot:chat:moderate` 才能执行群管理命令。

```
socket.emit('room.apply', { requestId: 'apply-1', roomId: 'ROOM_ID', message: '机器人用途' });
socket.emit('room.join', { requestId: 'join-1', inviteCode: 'INVITE_CODE' });
socket.emit('message.send', { requestId: 'send-1', roomId: 'ROOM_ID', content: '你好' });
socket.emit('message.read', { requestId: 'read-1', roomId: 'ROOM_ID', limit: 50 });
socket.emit('message.read', { requestId: 'read-2', roomId: 'ROOM_ID', messageId: 'MESSAGE_ID' });
socket.emit('message.edit', { requestId: 'edit-1', messageId: 'MESSAGE_ID', content: '修改后的内容' });
socket.emit('message.recall', { requestId: 'recall-1', messageId: 'MESSAGE_ID' });
socket.emit('typing.start', { roomId: 'ROOM_ID' });
socket.emit('typing.stop', { roomId: 'ROOM_ID' });
socket.emit('room.leave', { roomId: 'ROOM_ID' });
```

`message.edit` 仅允许发送者在 2 分钟内编辑；`message.recall` 仅允许发送者在 3 分钟内撤回。群管理员命令需要 `bot:chat:moderate`，并遵循站内群主/管理员层级：

```
socket.emit('room.member.role', { roomId: 'ROOM_ID', userId: 'USER_ID', role: 'admin' });
socket.emit('room.member.remove', { roomId: 'ROOM_ID', userId: 'USER_ID' });
socket.emit('room.member.mute', { roomId: 'ROOM_ID', userId: 'USER_ID', durationMs: 3600000, reason: '暂时禁言' });
socket.emit('room.member.unmute', { roomId: 'ROOM_ID', userId: 'USER_ID' });
```

## 私聊命令

机器人必须先作为参与者存在于会话中：

```
socket.emit('conversation.message.send', { conversationId: 'CONVERSATION_ID', content: '你好' });
socket.emit('conversation.message.read', { conversationId: 'CONVERSATION_ID', limit: 50 });
socket.emit('conversation.message.edit', { conversationId: 'CONVERSATION_ID', messageId: 'MESSAGE_ID', content: '修改后的内容' });
socket.emit('conversation.message.recall', { conversationId: 'CONVERSATION_ID', messageId: 'MESSAGE_ID' });
```

私聊消息同样是 2 分钟内可编辑、3 分钟内可撤回。历史读取返回 `conversation.message.batch`，新消息、编辑和撤回分别返回 `conversation.message.created`、`conversation.message.edited` 和 `conversation.message.recalled`。

## 互动按钮消息

机器人可以把 `content` 设置为互动消息，普通用户接口不允许使用该类型：

```
socket.emit('message.send', {
  roomId: 'ROOM_ID',
  content: {
    type: 'interactive',
    text: '请选择操作',
    buttons: [[
      { id: 'accept', label: '接受', action: 'invite.accept' },
      { id: 'reject', label: '拒绝', action: 'invite.reject' }
    ]]
  }
});
```

最多 10 行，每行最多 3 个按钮。用户点击后，机器人会在自己的连接中收到 `interaction.created`，其中包含 `messageId`、`buttonId`、`action`、`value` 和点击用户 ID。

## 事件类型

| 事件 | 触发时机 |
| --- | --- |
| `message.created` / `conversation.message.created` | 群聊或私聊新消息通过扫描后 |
| `message.edited` / `conversation.message.edited` | 消息被发送者编辑 |
| `message.recalled` / `conversation.message.recalled` | 消息被发送者撤回 |
| `message.read` | 群成员标记消息已读 |
| `typing.started` / `typing.stopped` | 用户或机器人输入状态变化 |
| `room.member.joined` / `room.member.left` / `room.member.updated` | 群成员关系变化 |
| `notification.created` | 所属用户收到站内通知 |
| `presence.updated` | 所属用户上线或离线 |
| `content.changed` | 公开且审核通过的资源、帖子或视频发生公开状态变化 |
| `server.status.updated` | 公开找服玩服务器在线状态、在线人数或延迟变化 |
| `interaction.created` | 用户点击机器人互动按钮 |
| `bot.mention` | 公开内容或消息提及机器人 |

## 响应与错误

带 `requestId` 的命令成功返回 `command.accepted`，失败返回 `command.rejected`。常见错误码包括 `SCOPE_REQUIRED`、`TOPIC_INVALID`、`TOPIC_FORBIDDEN`、`ROOM_MEMBERSHIP_REQUIRED`、`CONVERSATION_ACCESS_REQUIRED`、`MESSAGE_PAYLOAD_INVALID`、`MESSAGE_EDIT_WINDOW_EXPIRED`、`MESSAGE_RECALL_WINDOW_EXPIRED`、`ROOM_MANAGER_REQUIRED` 和 `COMMAND_UNSUPPORTED`。

机器人应限制连接数和发送频率，并避免记录完整 WebSocket Key。未识别的命令返回 `error`，其中包含 `COMMAND_UNSUPPORTED`；已识别但校验失败的命令返回 `command.rejected`。
