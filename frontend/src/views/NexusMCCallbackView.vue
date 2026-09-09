<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Loader2 } from 'lucide-vue-next'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const error = ref('')
const message = ref('')

onMounted(async () => {
  const queryError = route.query.error as string | undefined
  const bindSuccess = route.query.bind === 'success'

  if (queryError) {
    error.value = decodeURIComponent(queryError) === 'access_denied'
      ? '您拒绝了 NexusMC 授权'
      : `NexusMC 登录失败：${decodeURIComponent(queryError)}`
    return
  }

  if (bindSuccess) {
    message.value = 'NexusMC 账号绑定成功'
    setTimeout(() => router.push('/profile'), 1000)
    return
  }

  // 登录回调：token 通过 URL fragment 下发（避免进入服务器日志）
  const hash = new URLSearchParams(window.location.hash.slice(1))
  const accessToken = hash.get('access_token')
  const refreshToken = hash.get('refresh_token')

  if (!accessToken || !refreshToken) {
    error.value = 'NexusMC 回调参数缺失'
    return
  }

  authStore.setTokens(accessToken, refreshToken)
  await authStore.fetchUser()
  if (!authStore.user) {
    error.value = '登录态获取失败，请重试'
    return
  }
  router.push('/')
})
</script>

<template>
  <div class="mx-auto max-w-md px-4 py-16">
    <Card>
      <CardHeader>
        <CardTitle>NexusMC 登录</CardTitle>
      </CardHeader>
      <CardContent class="space-y-4">
        <div v-if="!error && !message" class="flex items-center gap-2 text-muted-foreground">
          <Loader2 class="h-4 w-4 animate-spin" />
          正在完成 NexusMC 授权登录…
        </div>
        <Alert v-if="message" class="bg-green-500/10 text-green-600 border-green-500/20">
          <AlertDescription>{{ message }}</AlertDescription>
        </Alert>
        <Alert v-if="error" variant="destructive">
          <AlertDescription>{{ error }}</AlertDescription>
        </Alert>
        <RouterLink to="/login" class="block text-center text-sm text-muted-foreground hover:text-foreground">
          返回登录
        </RouterLink>
      </CardContent>
    </Card>
  </div>
</template>
