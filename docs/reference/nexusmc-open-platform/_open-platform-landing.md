---
title: "NexusMC 开放中心（官网落地页）"
source: https://www.nexusmc.cn/open-platform
collected: 2026-09-26
---

> 本页为开放平台官网介绍页，无具体接口定义；具体接口文档见本目录其余各文件（源站：docs.nexusmc.cn）。控制台（/open-platform/console）需登录后使用。

NEXUSMC OPEN PLATFORM
## 开放中心
连接 NexusMC 的社区内容、用户授权与实时会话能力，为官网、机器人和第三方服务构建稳定接入。
进入控制台
阅读开发文档
登录后才能使用开放平台能力。
realtime.nexusmc.cn
01 const socket = io('/realtime/v1')
02 socket.emit('subscribe', {
03   topic: 'room:community'
04 })
05 ready · connected · 42ms
Realtime services operational
PLATFORM CAPABILITIES
## 一处管理三类开放能力
根据业务所需选择授权、数据或实时连接，每项能力都有独立资格与审核边界。
O
用户授权
## OAuth 应用
让第三方应用在用户确认后读取基础身份、投稿摘要和互动数据。
管理 OAuth 应用 →
A
公开数据
## 站点 API
稳定获取公开资源、帖子、找服玩、视频、活动、搜索和排行榜数据。
管理站点 API →
W
实时连接
## WebSocket
连接机器人、群聊、私聊、通知与内容变化事件，并支持断线补偿。
管理 WebSocket →
ACCESS PATH
## 资格清晰，权限按需开放
OAuth 与站点 API 面向已完成个人或机构认证的用户；WebSocket 面向通过开发者考试并同意开放平台协议的开发者。
01
确认所需能力
先区分用户授权、公开数据和实时事件，不申请无关权限。
02
完成对应资格
按控制台提示完成用户认证、等级要求或开发者考试。
03
提交并等待审核
申请通过后再生成凭据，密钥只在创建时完整显示一次。
REALTIME
## 为机器人准备的实时能力
查看 WebSocket 文档 →
群聊与私聊
订阅消息、发送机器人专属互动消息，并处理按钮操作。
通知与提及
接收机器人提及、定向通知和用户互动事件。
内容变化
跟踪资源、帖子、找服玩和视频的细分状态变化。
断线补偿
使用事件游标恢复中断期间错过的可重放事件。
AGREEMENT 2026-open-platform-v1
## 开放平台协议
我确认会遵守开放平台规则，仅申请业务所需权限，妥善保管凭据，并对机器人和第三方服务的行为负责。
登录
登录后才能使用开放平台能力。
READY TO BUILD
## 从一个明确的权限边界开始
打开控制台 浏览文档

