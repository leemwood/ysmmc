<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { LogIn, UserPlus } from 'lucide-vue-next'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const isLogin = ref(true)
const email = ref('')
const password = ref('')
const username = ref('')
const confirmPassword = ref('')
const error = ref('')
const loading = ref(false)

async function handleSubmit() {
  error.value = ''
  loading.value = true

  try {
    if (isLogin.value) {
      await authStore.login(email.value, password.value)
    } else {
      if (password.value !== confirmPassword.value) {
        error.value = '两次输入的密码不一致'
        loading.value = false
        return
      }
      if (username.value.length < 2) {
        error.value = '用户名至少需要2个字符'
        loading.value = false
        return
      }
      await authStore.register(email.value, password.value, username.value)
      await authStore.login(email.value, password.value)
    }

    const redirect = route.query.redirect as string
    router.push(redirect || '/')
  } catch (err: any) {
    error.value = err.response?.data?.message || '操作失败，请重试'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="mx-auto max-w-md px-4 py-12">
    <Card>
      <CardHeader class="text-center">
        <CardTitle class="text-2xl">
          {{ isLogin ? '登录' : '注册' }}
        </CardTitle>
      </CardHeader>
      <CardContent>
        <form @submit.prevent="handleSubmit" class="space-y-4">
          <div v-if="error" class="rounded-md bg-destructive/10 p-3 text-sm text-destructive">
            {{ error }}
          </div>

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

          <div v-if="!isLogin" class="space-y-2">
            <Label for="username">用户名</Label>
            <Input
              id="username"
              v-model="username"
              type="text"
              placeholder="请输入用户名"
              required
            />
          </div>

          <div class="space-y-2">
            <Label for="password">密码</Label>
            <Input
              id="password"
              v-model="password"
              type="password"
              placeholder="请输入密码"
              required
            />
          </div>

          <div v-if="!isLogin" class="space-y-2">
            <Label for="confirmPassword">确认密码</Label>
            <Input
              id="confirmPassword"
              v-model="confirmPassword"
              type="password"
              placeholder="请再次输入密码"
              required
            />
          </div>

          <Button type="submit" class="w-full" :disabled="loading">
            <LogIn v-if="isLogin" class="mr-2 h-4 w-4" />
            <UserPlus v-else class="mr-2 h-4 w-4" />
            {{ loading ? '处理中...' : (isLogin ? '登录' : '注册') }}
          </Button>

          <div class="flex items-center justify-between text-sm">
            <button
              type="button"
              class="text-primary hover:underline"
              @click="isLogin = !isLogin"
            >
              {{ isLogin ? '没有账号？去注册' : '已有账号？去登录' }}
            </button>
            <RouterLink
              v-if="isLogin"
              to="/reset-password"
              class="text-muted-foreground hover:text-foreground"
            >
              忘记密码？
            </RouterLink>
          </div>
        </form>
      </CardContent>
    </Card>
  </div>
</template>
