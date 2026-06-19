<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { KeyRound, Eye, EyeOff } from 'lucide-vue-next'
import AuthLayout from '@/components/layout/AuthLayout.vue'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const password = ref('')
const confirmPassword = ref('')
const loading = ref(false)
const error = ref('')
const success = ref(false)
const showPassword = ref(false)
const showConfirmPassword = ref(false)

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
  <AuthLayout title="设置新密码" description="请输入您的新密码">
    <div v-if="success" class="text-center space-y-4">
      <KeyRound class="mx-auto h-12 w-12 text-green-500" />
      <p class="text-sm text-muted-foreground">
        密码已重置成功，即将跳转到登录页面...
      </p>
    </div>

    <form v-else @submit.prevent="handleSubmit" class="space-y-4">
      <Alert v-if="error" :key="error" variant="destructive" class="animate-shake">
        <AlertDescription>{{ error }}</AlertDescription>
      </Alert>

      <div class="space-y-2">
        <Label for="password">新密码</Label>
        <Input
          id="password"
          v-model="password"
          :type="showPassword ? 'text' : 'password'"
          placeholder="请输入新密码"
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
        <Label for="confirmPassword">确认密码</Label>
        <Input
          id="confirmPassword"
          v-model="confirmPassword"
          :type="showConfirmPassword ? 'text' : 'password'"
          placeholder="请再次输入新密码"
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

      <Button type="submit" class="w-full h-11 btn-press" :loading="loading">
        重置密码
      </Button>
    </form>
  </AuthLayout>
</template>
