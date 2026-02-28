<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { authApi } from '@/lib/api'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Loader2, CheckCircle, XCircle } from 'lucide-vue-next'

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
  <div class="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900 px-4">
    <Card class="w-full max-w-md">
      <CardHeader class="text-center">
        <CardTitle class="text-2xl">邮箱验证</CardTitle>
      </CardHeader>
      <CardContent class="flex flex-col items-center space-y-4">
        <div v-if="loading" class="flex flex-col items-center space-y-4">
          <Loader2 class="h-12 w-12 animate-spin text-primary" />
          <p class="text-muted-foreground">正在验证...</p>
        </div>
        
        <div v-else class="flex flex-col items-center space-y-4">
          <CheckCircle v-if="success" class="h-12 w-12 text-green-500" />
          <XCircle v-else class="h-12 w-12 text-red-500" />
          
          <Alert :variant="success ? 'default' : 'destructive'">
            <AlertDescription class="text-center">{{ message }}</AlertDescription>
          </Alert>
          
          <p v-if="success" class="text-sm text-muted-foreground">
            即将跳转到登录页面...
          </p>
          
          <router-link 
            v-else 
            to="/login" 
            class="text-sm text-primary hover:underline"
          >
            返回登录
          </router-link>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
