<script setup lang="ts">
import { useRouter } from 'vue-router'
import { Button } from '@/components/ui/button'

const router = useRouter()

function goHome() {
  router.push({ name: 'home' })
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50">
    <div class="text-center">
      <h1 class="text-9xl font-bold text-gray-200">404</h1>
      <h2 class="text-2xl font-semibold text-gray-700 mt-4">页面未找到</h2>
      <p class="text-gray-500 mt-2">您访问的页面不存在或已被移除</p>
      <Button class="mt-6" @click="goHome">返回首页</Button>
    </div>
  </div>
</template>
