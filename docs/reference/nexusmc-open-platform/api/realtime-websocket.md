---
title: "WebSocket：正向与反向连接"
source: https://docs.nexusmc.cn/docs/api/realtime-websocket
collected: 2026-09-26
---

# WebSocket：正向与反向连接

区分原生 Socket.IO 和 OneBot WebSocket，配置连接方向、地址和密钥。

本站有两套 WebSocket 接入方式。原生实时 API 使用 Socket.IO，客户端连接本站 `/realtime/v1`；OneBot 11/12 使用普通 WebSocket，可选择客户端连接本站（正向），或由本站 worker 连接你的服务（反向）。两套协议的密钥和消息格式不同，不能混用。

| 接入 | 方向 | 协议与地址 | 凭据 |
| --- | --- | --- | --- |
| 原生实时 API | 客户端连接本站 | Socket.IO `https://www.nexusmc.cn/realtime/v1` | 机器人的 WebSocket Key |
| OneBot 正向 | 客户端连接本站 | 普通 WebSocket `wss://www.nexusmc.cn/onebot/v1/connections/<连接 ID>` | 该 OneBot 连接的密钥 |
| OneBot 反向 | 本站连接你的服务 | 连接配置中的公开 `ws://` 或 `wss://` 地址 | 同一个 OneBot 连接的密钥，由本站在握手时发送 |

连接前，账号需要正常、邮箱已验证、通过开发者考试，并同意当前版本的开放平台协议。可在[开放平台控制台](https://www.nexusmc.cn/open-platform/console)创建机器人、WebSocket Key 和 OneBot 连接。连接密钥只在创建或轮换时完整显示；不要把它放在公开网页代码或日志中。

## 原生 Socket.IO

创建机器人时授予 `bot:realtime:connect`，再创建 WebSocket Key。下例订阅公开资源，机器人还需 `bot:chat:subscribe` 和 `bot:content:subscribe`。客户端使用 `socket.io-client`，通过 `auth.token` 或 `Authorization: Bearer` 认证：

```
import { io } from 'socket.io-client';

const socket = io('https://www.nexusmc.cn/realtime/v1', {
  auth: { token: process.env.NEXUSMC_WEBSOCKET_KEY },
  transports: ['websocket'],
});

socket.on('connect_error', (error) => console.error(error.message));
socket.on('ready', (event) => console.log(event.payload.scopes, event.payload.latestCursor));
socket.emit('subscribe', { requestId: 'sub-1', topic: 'content:resource' });
socket.on('content.changed', (event) => console.log(event.payload));
```

服务端先发送 `ready`，其中包含协议名、机器人 ID、scope、心跳间隔和最近事件游标。认证失败时连接错误为 `BOT_AUTH_REQUIRED` 或 `BOT_AUTH_INVALID`。此入口只有客户端连本站的模式；它不接收普通 WebSocket 的 OneBot JSON 帧，也没有反向 Socket.IO 连接。

原生 `subscribe`、消息、群管理、私信、`ack` 和 `resume` 等所有命令与事件见[WebSocket 接口参考](/docs/api/realtime-websocket-reference/)。

## OneBot 正向 WebSocket

在控制台创建 `transport: "forward_websocket"` 的连接，选择 `onebot11` 或 `onebot12`，可选绑定机器人身份。创建响应中的 `forwardUrl` 是相对路径，需拼接本站主域名。Node.js `ws` 客户端示例：

```
import WebSocket from 'ws';

const ws = new WebSocket('wss://www.nexusmc.cn/onebot/v1/connections/CONNECTION_ID', {
  headers: { Authorization: `Bearer ${process.env.NEXUSMC_ONEBOT_SECRET}` },
});

ws.on('message', (frame) => console.log(JSON.parse(frame.toString())));
ws.on('error', console.error);
```

握手失败返回 HTTP `401`。建立连接后先收到 OneBot 生命周期事件；事件帧和动作请求使用普通 JSON。完整的 OneBot 11/12 格式见[OneBot 接口与请求](/docs/api/open-platform-onebot-reference/)。

## OneBot 反向 WebSocket

创建 `transport: "reverse_websocket"` 的连接，并将 `endpoint` 设为你的公开 `ws://` 或 `wss://` 地址。本站 worker 主动建立连接，握手时发送 `Authorization: Bearer <连接密钥>`；你的服务需要验证该密钥，再接收事件或发送动作。示例服务端：

```
import { createServer } from 'node:http';
import { WebSocketServer } from 'ws';

const server = createServer();
const wss = new WebSocketServer({ noServer: true });
server.on('upgrade', (request, socket, head) => {
  if (request.headers.authorization !== `Bearer ${process.env.NEXUSMC_ONEBOT_SECRET}`) {
    socket.end('HTTP/1.1 401 Unauthorized\r\nConnection: close\r\n\r\n');
    return;
  }
  wss.handleUpgrade(request, socket, head, (client) => {
    client.on('message', (frame) => console.log(JSON.parse(frame.toString())));
  });
});
server.listen(8080);
```

通过公开、安全的目标地址提供服务；本站会拒绝内网、回环等不安全目标。worker 会在断线后重连，待投递事件由后台队列处理。正向与反向的 OneBot JSON 帧相同，区别只是发起连接的一方。OneBot HTTP WebHook 是单向 HTTP 推送，见[Webhook 接口](/docs/api/open-platform-webhooks/)。

OneBot 连接的机器人绑定、事件筛选、管理和入群申请见[OneBot 11/12 接入概念](/docs/api/open-platform-bots/)。
