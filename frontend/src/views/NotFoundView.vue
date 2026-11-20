<script setup lang="ts">
import { ArrowLeft, Box } from 'lucide-vue-next'
</script>

<template>
  <div class="container-app flex min-h-[58vh] items-center justify-center py-8">
    <section class="site-panel w-full max-w-md p-6 text-center sm:p-10">
      <div class="status-icon mx-auto"><Box class="size-6" aria-hidden="true" /></div>
      <h1 class="mt-4 text-2xl font-bold">页面未找到</h1>
      <p class="mt-2 text-sm text-muted-foreground">这个地址可能已失效。你可以从模型库继续浏览。</p>
      <RouterLink to="/" class="site-button site-button-primary mt-5">
        <ArrowLeft class="size-4" aria-hidden="true" /> 返回模型库
      </RouterLink>
    </section>
  </div>
</template>
