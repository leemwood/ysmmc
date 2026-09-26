---
title: "个人 API：错误与排查"
source: https://docs.nexusmc.cn/docs/api/personal-api-errors
collected: 2026-09-26
---

# 个人 API：错误与排查

个人 API 的状态码、错误结构、鉴权、权限、组织限制、限流和常见请求问题。

## 错误响应结构

多数错误至少包含 `error`，部分可编程错误还会返回 `code` 和附加字段：

```
{
  "error": "当前 API Token 没有此操作权限",
  "code": "API_TOKEN_SCOPE_REQUIRED",
  "requiredScope": "resource:create"
}
```

客户端应优先根据 HTTP 状态码和 `code` 分支处理，把 `error` 用作日志或用户提示。并非每个旧接口都会返回 `code`。

## 状态码速查

| 状态码 | 含义 | 首先检查 |
| --- | --- | --- |
| `200` | 查询或更新成功 | 读取响应 JSON。 |
| `201` | 创建成功 | 保存响应里的内容 ID、URL 或文件信息。 |
| `400` | 请求参数或组合规则不合法 | 字段类型、必填项、枚举值、URL、正文和关联 ID。 |
| `401` | 未认证 | `Authorization: Bearer ...` 是否存在，Token 是否完整、有效或已撤销。 |
| `403` | 已认证但无权操作 | 权限、接口是否接受个人 Token、内容所有权、组织身份限制和账号安全前置条件。 |
| `404` | 内容不存在 | 用“我的内容”列表接口核对 ID、slug 和当前 Token。 |
| `409` | 状态冲突 | 自定义路径、重复数据或当前内容状态。 |
| `413` | 请求体过大 | 文件大小、正文大小和所用上传接口。 |
| `429` | 请求过于频繁 | 限流响应头和站点发布频率限制。 |
| `500` | 服务端异常 | 记录请求时间和响应，请勿立即无限重试。 |

## 鉴权失败

请求头必须使用 Bearer 认证：

```
Authorization: Bearer avm_pat_xxxxx_xxxxx
```

常见错误：

-   把 Token 放在 URL 查询参数里。
-   少写 `Bearer` 或把整段请求头加引号。
-   Token 复制不完整、过期或已在站内撤销。
-   请求发到了网页地址，而不是 `/api/...`。

CAUTION

个人 API Token 相当于账号凭证。不要写入公开仓库、浏览器前端包、构建日志或错误截图；泄露后立即撤销并重建。

## 权限不足

```
{
  "error": "当前 API Token 没有此操作权限",
  "code": "API_TOKEN_SCOPE_REQUIRED",
  "requiredScope": "post:update:self"
}
```

处理步骤：

1.  对照具体业务文档顶部的权限要求。
2.  回到站内 Token 设置检查是否授予对应权限。
3.  权限变更不应靠客户端猜测；缺少权限时停止该操作并提示用户重新授权。

## 接口不接受个人 Token

部分接口只服务于站内登录会话，不开放给个人 API Token：

```
{
  "error": "当前接口不允许使用个人 API Token，请使用站内登录会话",
  "code": "API_TOKEN_ROUTE_FORBIDDEN",
  "requireSessionConfirmation": true
}
```

这类错误与权限勾选无关，扩大 Token 权限不会改变结果。遇到时改用站内登录会话，或改走该能力的公开读取接口。

同是 `403`，两种 `code` 的处理方式不同：`API_TOKEN_SCOPE_REQUIRED` 是权限没勾对，补齐对应权限即可；`API_TOKEN_ROUTE_FORBIDDEN` 是该接口不面向个人 Token，补权限无效。

## 只能操作自己的内容

带 `:self` 的权限仍会校验内容所有权。即使 Token 用户是协作者、组织成员、版主或管理员，个人 API 也不能借这些身份编辑其他作者的内容。

如果返回 `403` 且权限正确，请确认：

-   请求中的内容 ID 是否由当前 Token 用户创建。
-   是否误用了组织内容或他人的草稿 ID。
-   是否把列表接口返回的展示 ID 和另一个内容类型的 ID 混用。

## 内容 ID 或 slug 返回 404

四类列表接口会直接返回可用于后续请求的稳定 `id`，无需从网页 URL 获取：

| 内容 | 查询接口 | ID 字段 |
| --- | --- | --- |
| 资源 | `GET /api/users/me/resources` | `resources[].id` |
| 帖子 | `GET /api/users/me/posts` | `posts[].id` |
| 找服玩 | `GET /api/users/me/players` | `players[].id` |
| 视频 | `GET /api/users/me/videos` | `videos[].id` |

内容路径同时接受稳定 ID、当前或历史 slug、公开编号和公开链接中的路径段。仍返回 `404` 时依次检查：

1.  查询值是否只是标题或项目名，而不是响应里的 `id` 或 `slug`。
2.  该内容是否由当前 Token 用户本人创建。
3.  仓库变量或配置值前后是否有空格、引号或换行。
4.  是否把资源 ID 用到了帖子、找服玩或视频接口。

GitHub Actions 出现 `NexusMC API request failed with HTTP 404: {"error":"Resource not found"}` 时，通常是资源更新阶段没有解析到 `resource_id`，与前面的构建或文件上传是否成功无关。

## 创作者发布资格不足

Token 已包含创建 scope 时，正式投稿仍可能返回：

```
{
  "error": "完成开发者考试后，才能进入创作者签约步骤。",
  "code": "CREATOR_EXAM_REQUIRED"
}
```

或：

```
{
  "error": "阅读并同意创作者协议后，将开放资源、帖子、服务器和视频发布能力。",
  "code": "CREATOR_AGREEMENT_REQUIRED"
}
```

这类错误与 Token scope 无关。请使用 Token 所属账号完成 [创作者开通与协议确认](/docs/site-info/creator-onboarding/)，不要通过扩大 Token 权限反复重试。

## 个人 Token 不支持组织身份

```
{
  "error": "个人 API Token 暂不支持组织身份投稿",
  "code": "API_TOKEN_ORGANIZATION_FORBIDDEN"
}
```

创建或更新资源、帖子、找服玩、视频时，不要传非空 `organizationId`。个人 API 创建的内容归 Token 用户本人。

## 参数校验失败

收到 `400` 时按以下顺序检查：

1.  `Content-Type` 是否为 `application/json`，上传接口除外。
2.  请求体是不是合法 JSON，布尔值和数字有没有误传成字符串。
3.  查看当前端点自己的完整参数表，不要把另一个创建或更新接口的请求体原样复用。
4.  检查组合规则，例如 `partial` 可见性必须带用户列表、密码模式必须带密码。
5.  数组字段是否传了数组；要清空通常传 `[]`，要保留就省略字段。

## 限流

触发个人 API 限流时返回 `429`：

```
{
  "error": "请求过于频繁，请稍后再试",
  "code": "PERSONAL_API_RATE_LIMIT_EXCEEDED"
}
```

| 响应头 | 说明 |
| --- | --- |
| `RateLimit-Limit` | 当前窗口允许的请求数。 |
| `RateLimit-Remaining` | 当前窗口剩余请求数。 |
| `RateLimit-Reset` | 距窗口重置的秒数。 |
| `Retry-After` | 建议等待的秒数；存在时优先使用。 |

客户端应采用指数退避并加入少量随机抖动。不要对 `400`、`401`、`403` 无条件重试；这些错误通常需要修改请求或权限。

## 提交排查信息

仍无法解决时，提供以下信息即可，不要提供完整 Token：

-   请求方法和路径，例如 `PATCH /api/resources/{id}`。
-   HTTP 状态码与完整响应 JSON。
-   请求体的脱敏版本。
-   请求时间与时区。
-   Token 只保留前后少量字符用于区分，例如 `avm_pat_abcd…wxyz`。
