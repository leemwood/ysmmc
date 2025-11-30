# 部署指南 (Cloudflare Pages)

如果您的项目部署在 **Cloudflare Pages** 并且遇到 `401 Unauthorized` 或其他 Supabase 连接错误，通常是因为**环境变量未正确配置**。

由于 `.env` 文件通常被 git 忽略（为了安全），Cloudflare 的构建服务器无法直接读取本地的环境变量。您必须在 Cloudflare Pages 的后台手动添加它们。

## ✅ 解决步骤

1. 登录 [Cloudflare Dashboard](https://dash.cloudflare.com/)。
2. 进入 **Pages** -> 选择您的项目 (`ysmmc`)。
3. 点击 **Settings** (设置) -> **Environment variables** (环境变量)。
4. 点击 **Add variable** (添加变量)，添加以下两个变量（值请参考您本地的 `.env` 文件）：

   | Variable Name (变量名) | Value (值) |
   | ----------------------- | ---------- |
   | `VITE_SUPABASE_URL`     | `您的Supabase项目URL` |
   | `VITE_SUPABASE_ANON_KEY`| `您的Supabase Anon Key` |

   > **注意**：请同时为 **Production** (生产环境) 和 **Preview** (预览环境) 配置这些变量。

5. **关键步骤：重新部署**
   配置完环境变量后，它们不会立即生效。您必须：
   - 转到 **Deployments** (部署) 标签页。
   - 找到最新的部署，点击 **Retry deployment** (重试部署) 或者推一个新的 commit 触发重新构建。
   - **构建过程**（Build）会读取这些变量并将它们打包到前端代码中。

## 🔍 验证方法

部署完成后，打开网站控制台（F12 -> Console）。
- 如果看到 `Supabase URL or Key is missing` 的红色错误，说明环境变量仍然没有读取到。
- 如果没有报错，但请求依然 401，请检查填写的 `ANON_KEY` 是否正确（不要误填成 Service Role Key）。
