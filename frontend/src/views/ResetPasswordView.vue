<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Mail, ArrowLeft } from 'lucide-vue-next'
import { RouterLink } from 'vue-router'

const authStore = useAuthStore()

const email = ref('')
const loading = ref(false)
const submitted = ref(false)
const error = ref('')

async function handleSubmit() {
  error.value = ''
  loading.value = true

  try {
    await authStore.forgotPassword(email.value)
    submitted.value = true
  } catch (err: any) {
    error.value = err.response?.data?.message || '发送失败，请重试'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="mx-auto max-w-md px-4 py-12">
    <Card>
      <CardHeader class="text-center">
        <CardTitle class="text-2xl">重置密码</CardTitle>
      </CardHeader>
      <CardContent>
        <div v-if="submitted" class="text-center space-y-4">
          <Mail class="mx-auto h-12 w-12 text-primary" />
          <p class="text-muted-foreground">
            如果该邮箱已注册，我们已发送重置密码链接到您的邮箱。
          </p>
          <RouterLink to="/login">
            <Button variant="outline">
              <ArrowLeft class="mr-2 h-4 w-4" />
              返回登录
            </Button>
          </RouterLink>
        </div>

        <form v-else @submit.prevent="handleSubmit" class="space-y-4">
          <div v-if="error" class="rounded-md bg-destructive/10 p-3 text-sm text-destructive">
            {{ error }}
          </div>

          <p class="text-sm text-muted-foreground">
            请输入您的注册邮箱，我们将发送重置密码链接。
          </p>

          <div class="space-y-2">
            <Label for="email">邮箱</Label>
            <Input
              id="email"
              v-model="email"
              type="email"
              placeholder="请输入邮箱"
              required
            />
          </div>

          <Button type="submit" class="w-full" :disabled="loading">
            {{ loading ? '发送中...' : '发送重置链接' }}
          </Button>

          <RouterLink to="/login" class="block text-center text-sm text-muted-foreground hover:text-foreground">
            返回登录
          </RouterLink>
        </form>
      </CardContent>
    </Card>
  </div>
</template>
