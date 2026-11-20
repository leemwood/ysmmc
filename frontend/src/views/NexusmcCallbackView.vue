<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { CircleAlert, LoaderCircle } from 'lucide-vue-next'
import { type NexusmcUser } from '@/lib/nexusmc'
import { useNexusmcStore } from '@/stores/nexusmc'

const route = useRoute()
const router = useRouter()
const store = useNexusmcStore()

const error = ref('')
const errorMessages: Record<string, string> = {
  access_denied: '你取消了 NexusMC 授权。',
  state_mismatch: '登录状态校验失败，请重新发起登录。',
  missing_code: '未收到 NexusMC 的授权码，请重新登录。',
  token_exchange: '换取 NexusMC 令牌失败，请稍后重试。',
  userinfo: '获取 NexusMC 用户信息失败，请稍后重试。',
  unauthorized_client: 'NexusMC 侧的 OAuth 应用尚未通过审核或已停用，请联系 NexusMC 管理员确认后重试。',
  invalid_redirect_uri: '回调地址与 NexusMC 登记的不一致，请联系站点管理员。',
}

onMounted(() => {
  const rawUser = route.query.user
  if (typeof rawUser === 'string' && rawUser) {
    try {
      const user = JSON.parse(rawUser) as NexusmcUser
      if (!user.sub && !user.username) throw new Error('invalid user payload')
      store.completeLogin(user)
      // 默认回首页；从受保护页面（如 /me）发起登录时回跳原目标
      void router.replace(store.consumeReturnTo() || '/')
      return
    } catch {
      error.value = '登录数据解析失败，请重新登录。'
      return
    }
  }

  const err = typeof route.query.error === 'string' ? route.query.error : ''
  error.value = errorMessages[err] || 'NexusMC 登录未完成，请重试。'
})
</script>

<template>
  <div class="container-app flex min-h-[58vh] items-center justify-center py-8">
    <section class="site-panel w-full max-w-md p-6 text-center sm:p-9" :role="error ? 'alert' : 'status'">
      <div class="status-icon mx-auto">
        <CircleAlert v-if="error" class="size-6" aria-hidden="true" />
        <LoaderCircle v-else class="size-6 animate-spin" aria-hidden="true" />
      </div>
      <h1 class="mt-4 text-xl font-bold">{{ error ? '登录未完成' : '正在完成 NexusMC 登录…' }}</h1>
      <template v-if="error">
        <p class="mt-2 text-sm text-muted-foreground">{{ error }}</p>
        <button type="button" class="site-button site-button-primary mt-5" @click="router.replace('/')">返回模型库</button>
      </template>
    </section>
  </div>
</template>
