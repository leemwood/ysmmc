<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs'
import { LogIn, UserPlus, Eye, EyeOff } from 'lucide-vue-next'
import AuthLayout from '@/components/layout/AuthLayout.vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const activeTab = ref((route.query.mode === 'register') ? 'register' : 'login')
const isLogin = computed(() => activeTab.value === 'login')

const email = ref('')
const password = ref('')
const username = ref('')
const confirmPassword = ref('')
const error = ref('')
const success = ref('')
const loading = ref(false)
const showPassword = ref(false)
const showConfirmPassword = ref(false)

watch(activeTab, () => {
  error.value = ''
  success.value = ''
})

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
  <AuthLayout
    :title="isLogin ? '欢迎回来' : '创建账号'"
    :description="isLogin ? '登录您的 YSM 模型站账号' : '注册成为 YSM 模型站用户'"
  >
    <Tabs v-model="activeTab" class="w-full">
      <TabsList class="grid w-full grid-cols-2 mb-6">
        <TabsTrigger value="login" class="gap-2">
          <LogIn class="h-4 w-4" />
          登录
        </TabsTrigger>
        <TabsTrigger value="register" class="gap-2">
          <UserPlus class="h-4 w-4" />
          注册
        </TabsTrigger>
      </TabsList>

      <TabsContent value="login" class="mt-0">
        <form @submit.prevent="handleSubmit" class="space-y-4">
          <Alert v-if="error" :key="error" variant="destructive" class="animate-shake">
            <AlertDescription>{{ error }}</AlertDescription>
          </Alert>

          <div class="space-y-2">
            <Label for="login-email">邮箱</Label>
            <Input
              id="login-email"
              v-model="email"
              type="email"
              placeholder="请输入邮箱"
              required
              autocomplete="email"
              class="h-11"
              :error="!!error"
            />
          </div>

          <div class="space-y-2">
            <Label for="login-password">密码</Label>
            <Input
              id="login-password"
              v-model="password"
              :type="showPassword ? 'text' : 'password'"
              placeholder="请输入密码"
              required
              autocomplete="current-password"
              class="h-11"
              :error="!!error"
            >
              <template #suffix>
                <button
                  type="button"
                  class="focus-ring rounded-md p-1 text-muted-foreground hover:text-foreground transition-colors"
                  @click="showPassword = !showPassword"
                >
                  <Eye v-if="!showPassword" class="h-4 w-4" />
                  <EyeOff v-else class="h-4 w-4" />
                </button>
              </template>
            </Input>
          </div>

          <Button type="submit" class="w-full btn-press h-11" :loading="loading">
            登录
          </Button>

          <div class="flex items-center justify-center">
            <RouterLink
              to="/reset-password"
              class="text-sm text-muted-foreground hover:text-foreground transition-colors py-2"
            >
              忘记密码？
            </RouterLink>
          </div>
        </form>
      </TabsContent>

      <TabsContent value="register" class="mt-0">
        <form @submit.prevent="handleSubmit" class="space-y-4">
          <Alert v-if="error" :key="error" variant="destructive" class="animate-shake">
            <AlertDescription>{{ error }}</AlertDescription>
          </Alert>

          <Alert
            v-if="success"
            variant="default"
            class="bg-green-500/10 text-green-600 border-green-500/20"
          >
            <AlertDescription>{{ success }}</AlertDescription>
          </Alert>

          <div class="space-y-2">
            <Label for="register-email">邮箱</Label>
            <Input
              id="register-email"
              v-model="email"
              type="email"
              placeholder="请输入邮箱"
              required
              autocomplete="email"
              class="h-11"
              :error="!!error"
            />
          </div>

          <div class="space-y-2">
            <Label for="register-username">用户名</Label>
            <Input
              id="register-username"
              v-model="username"
              type="text"
              placeholder="请输入用户名（至少2个字符）"
              required
              autocomplete="username"
              class="h-11"
              :error="!!error"
            />
          </div>

          <div class="space-y-2">
            <Label for="register-password">密码</Label>
            <Input
              id="register-password"
              v-model="password"
              :type="showPassword ? 'text' : 'password'"
              placeholder="请输入密码（至少6个字符）"
              required
              autocomplete="new-password"
              class="h-11"
              :error="!!error"
            >
              <template #suffix>
                <button
                  type="button"
                  class="focus-ring rounded-md p-1 text-muted-foreground hover:text-foreground transition-colors"
                  @click="showPassword = !showPassword"
                >
                  <Eye v-if="!showPassword" class="h-4 w-4" />
                  <EyeOff v-else class="h-4 w-4" />
                </button>
              </template>
            </Input>
          </div>

          <div class="space-y-2">
            <Label for="register-confirm-password">确认密码</Label>
            <Input
              id="register-confirm-password"
              v-model="confirmPassword"
              :type="showConfirmPassword ? 'text' : 'password'"
              placeholder="请再次输入密码"
              required
              autocomplete="new-password"
              class="h-11"
              :error="!!error"
            >
              <template #suffix>
                <button
                  type="button"
                  class="focus-ring rounded-md p-1 text-muted-foreground hover:text-foreground transition-colors"
                  @click="showConfirmPassword = !showConfirmPassword"
                >
                  <Eye v-if="!showConfirmPassword" class="h-4 w-4" />
                  <EyeOff v-else class="h-4 w-4" />
                </button>
              </template>
            </Input>
          </div>

          <Button type="submit" class="w-full btn-press h-11" :loading="loading">
            注册
          </Button>
        </form>
      </TabsContent>
    </Tabs>
  </AuthLayout>
</template>
