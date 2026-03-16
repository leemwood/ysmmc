<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardHeader, CardTitle, CardContent, CardDescription } from '@/components/ui/card'
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
const success = ref('')
const loading = ref(false)

async function handleSubmit() {
  error.value = ''
  success.value = ''
  loading.value = true

  try {
    if (isLogin.value) {
      await authStore.login(email.value, password.value)
      const redirect = route.query.redirect as string
      router.push(redirect || '/')
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
      if (password.value.length < 6) {
        error.value = '密码至少需要6个字符'
        loading.value = false
        return
      }
      await authStore.register(email.value, password.value, username.value)
      await authStore.login(email.value, password.value)
      success.value = '注册成功！验证邮件已发送到您的邮箱，请查收。'
      const redirect = route.query.redirect as string
      router.push(redirect || '/')
    }
  } catch (err: any) {
    error.value = err.response?.data?.message || '操作失败，请重试'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-[calc(100vh-8rem)] flex items-center justify-center px-4 py-8 sm:py-12 relative overflow-hidden">
    <div class="absolute inset-0 -z-10">
      <div class="absolute top-1/4 left-1/4 w-96 h-96 bg-primary/5 rounded-full blur-3xl"></div>
      <div class="absolute bottom-1/4 right-1/4 w-96 h-96 bg-primary/10 rounded-full blur-3xl"></div>
    </div>

    <Card class="w-full max-w-md animate-slide-up relative">
      <CardHeader class="text-center space-y-2">
        <CardTitle class="text-2xl sm:text-3xl">
          {{ isLogin ? '欢迎回来' : '创建账号' }}
        </CardTitle>
        <CardDescription>
          {{ isLogin ? '登录您的 YSM 模型站账号' : '注册成为 YSM 模型站用户' }}
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form @submit.prevent="handleSubmit" class="space-y-4">
          <div v-if="error" class="rounded-md bg-destructive/10 p-3 text-sm text-destructive animate-fade-in">
            {{ error }}
          </div>
          
          <div v-if="success" class="rounded-md bg-green-500/10 p-3 text-sm text-green-600 animate-fade-in">
            {{ success }}
          </div>

          <div class="space-y-2">
            <Label for="email">邮箱</Label>
            <Input
              id="email"
              v-model="email"
              type="email"
              placeholder="请输入邮箱"
              required
              autocomplete="email"
            />
          </div>

          <div v-if="!isLogin" class="space-y-2">
            <Label for="username">用户名</Label>
            <Input
              id="username"
              v-model="username"
              type="text"
              placeholder="请输入用户名（至少2个字符）"
              required
              autocomplete="username"
            />
          </div>

          <div class="space-y-2">
            <Label for="password">密码</Label>
            <Input
              id="password"
              v-model="password"
              type="password"
              placeholder="请输入密码（至少6个字符）"
              required
              autocomplete="current-password"
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
              autocomplete="new-password"
            />
          </div>

          <Button type="submit" class="w-full btn-press" :disabled="loading">
            <LogIn v-if="isLogin" class="mr-2 h-4 w-4" />
            <UserPlus v-else class="mr-2 h-4 w-4" />
            {{ loading ? '处理中...' : (isLogin ? '登录' : '注册') }}
          </Button>

          <div class="flex flex-col sm:flex-row items-center justify-between gap-2 text-sm">
            <button
              type="button"
              class="text-primary hover:underline transition-colors"
              @click="isLogin = !isLogin"
            >
              {{ isLogin ? '没有账号？去注册' : '已有账号？去登录' }}
            </button>
            <RouterLink
              v-if="isLogin"
              to="/reset-password"
              class="text-muted-foreground hover:text-foreground transition-colors"
            >
              忘记密码？
            </RouterLink>
          </div>
        </form>
      </CardContent>
    </Card>
  </div>
</template>
