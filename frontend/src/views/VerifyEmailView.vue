<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { authApi } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Loader2, CheckCircle, XCircle, Home } from 'lucide-vue-next'
import { RouterLink } from 'vue-router'
import AuthLayout from '@/components/layout/AuthLayout.vue'

const route = useRoute()
const router = useRouter()

const loading = ref(true)
const success = ref(false)
const message = ref('')

onMounted(async () => {
  const token = route.query.token as string

  if (!token) {
    loading.value = false
    success.value = false
    message.value = '无效的验证链接'
    return
  }

  try {
    const response = await authApi.verifyEmail(token)
    success.value = true
    message.value = response.data.message || '邮箱验证成功！'

    setTimeout(() => {
      router.push('/login')
    }, 3000)
  } catch (error: any) {
    success.value = false
    message.value = error.response?.data?.message || '验证失败，请重试'
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <AuthLayout title="邮箱验证">
    <div class="text-center space-y-4">
      <Loader2 v-if="loading" class="mx-auto h-12 w-12 animate-spin text-primary" />
      <CheckCircle v-else-if="success" class="mx-auto h-12 w-12 text-green-500" />
      <XCircle v-else class="mx-auto h-12 w-12 text-destructive" />

      <p class="text-sm font-medium">
        {{ loading ? '正在验证...' : (success ? '验证成功' : '验证失败') }}
      </p>

      <Alert v-if="!loading" :variant="success ? 'default' : 'destructive'">
        <AlertDescription>{{ message }}</AlertDescription>
      </Alert>

      <p v-if="success" class="text-sm text-muted-foreground">
        即将跳转到登录页面...
      </p>

      <RouterLink v-if="!loading && !success" to="/">
        <Button variant="outline" class="btn-press">
          <Home class="mr-2 h-4 w-4" />
          返回首页
        </Button>
      </RouterLink>
    </div>
  </AuthLayout>
</template>
