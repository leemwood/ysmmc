<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { KeyRound } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const password = ref('')
const confirmPassword = ref('')
const loading = ref(false)
const error = ref('')
const success = ref(false)

onMounted(() => {
  if (!route.query.token) {
    router.push('/login')
  }
})

async function handleSubmit() {
  error.value = ''

  if (password.value !== confirmPassword.value) {
    error.value = '两次输入的密码不一致'
    return
  }

  if (password.value.length < 6) {
    error.value = '密码至少需要6个字符'
    return
  }

  loading.value = true

  try {
    await authStore.resetPassword(route.query.token as string, password.value)
    success.value = true
    setTimeout(() => {
      router.push('/login')
    }, 2000)
  } catch (err: any) {
    error.value = err.response?.data?.message || '重置失败，请重试'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="mx-auto max-w-md px-4 py-12">
    <Card>
      <CardHeader class="text-center">
        <CardTitle class="text-2xl">设置新密码</CardTitle>
      </CardHeader>
      <CardContent>
        <div v-if="success" class="text-center space-y-4">
          <KeyRound class="mx-auto h-12 w-12 text-green-500" />
          <p class="text-muted-foreground">
            密码已重置成功，即将跳转到登录页面...
          </p>
        </div>

        <form v-else @submit.prevent="handleSubmit" class="space-y-4">
          <div v-if="error" class="rounded-md bg-destructive/10 p-3 text-sm text-destructive">
            {{ error }}
          </div>

          <div class="space-y-2">
            <Label for="password">新密码</Label>
            <Input
              id="password"
              v-model="password"
              type="password"
              placeholder="请输入新密码"
              required
            />
          </div>

          <div class="space-y-2">
            <Label for="confirmPassword">确认密码</Label>
            <Input
              id="confirmPassword"
              v-model="confirmPassword"
              type="password"
              placeholder="请再次输入新密码"
              required
            />
          </div>

          <Button type="submit" class="w-full" :disabled="loading">
            {{ loading ? '处理中...' : '重置密码' }}
          </Button>
        </form>
      </CardContent>
    </Card>
  </div>
</template>
