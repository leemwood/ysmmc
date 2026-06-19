<script setup lang="ts">
import { useRouter } from 'vue-router'
import { Button } from '@/components/ui/button'
import { Home } from 'lucide-vue-next'

const router = useRouter()

function goHome() {
  router.push({ name: 'home' })
}
</script>

<template>
  <div class="container-app flex min-h-[calc(100vh-4rem)] items-center justify-center px-4 py-12">
    <div class="surface w-full max-w-md animate-fade-in px-6 py-12 text-center">
      <h1 class="text-8xl font-bold tracking-tighter gradient-text sm:text-9xl">404</h1>
      <h2 class="mt-4 text-2xl font-semibold tracking-tight sm:text-3xl">页面未找到</h2>
      <p class="mt-2 text-sm text-muted-foreground sm:text-base">
        您访问的页面不存在或已被移除
      </p>
      <Button class="focus-ring btn-press mt-8" @click="goHome">
        <Home class="mr-2 h-4 w-4" />
        返回首页
      </Button>
    </div>
  </div>
</template>
