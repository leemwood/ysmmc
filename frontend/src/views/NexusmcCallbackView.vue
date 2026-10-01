<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { setNexusmcUser, type NexusmcUser } from '@/lib/nexusmc'
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
      setNexusmcUser(user)
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
  <div class="container-app flex min-h-[60vh] items-center justify-center py-10">
    <Card class="w-full max-w-md">
      <CardContent class="flex flex-col items-center gap-4 p-8 text-center">
        <template v-if="error">
          <p class="font-semibold text-destructive">登录未完成</p>
          <p class="text-sm text-muted-foreground">{{ error }}</p>
          <Button size="sm" class="btn-press focus-ring" @click="router.replace('/')">
            返回资源首页
          </Button>
        </template>
        <template v-else>
          <p class="text-sm text-muted-foreground">正在完成 NexusMC 登录…</p>
        </template>
      </CardContent>
    </Card>
  </div>
</template>
