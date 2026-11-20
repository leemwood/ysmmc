<script setup lang="ts">
import { ArrowUp, ArrowUpRight, Box, Github, MessageCircle } from 'lucide-vue-next'
import { NEXUSMC_SITE } from '@/lib/nexusmc'

function scrollToTop() {
  window.scrollTo({ top: 0 })
}
</script>

<template>
  <footer class="border-t border-border">
    <div class="container-app flex flex-wrap items-center justify-between gap-x-8 gap-y-5 py-7">
      <div>
        <RouterLink to="/" class="inline-flex items-center gap-2 rounded-sm font-bold">
          <Box class="size-[1.125rem] text-primary" aria-hidden="true" /> YSM 模型站
        </RouterLink>
        <p class="mt-1 text-xs text-muted-foreground">非营利模型分享平台 · GPL-3.0 开源</p>
      </div>
      <div class="flex flex-wrap items-center gap-x-5 gap-y-3 text-sm font-medium text-muted-foreground">
        <a href="https://github.com/leemwood/ysmmc" target="_blank" rel="noopener noreferrer" class="inline-flex items-center gap-1.5 rounded-sm hover:text-primary">
          <Github class="size-4" aria-hidden="true" /> GitHub
        </a>
        <a href="https://qm.qq.com/q/SUKmYH7RyW" target="_blank" rel="noopener noreferrer" class="inline-flex items-center gap-1.5 rounded-sm hover:text-primary">
          <MessageCircle class="size-4" aria-hidden="true" /> QQ 交流群
        </a>
        <a :href="NEXUSMC_SITE" target="_blank" rel="noopener noreferrer" class="inline-flex items-center gap-1 rounded-sm hover:text-primary">
          NexusMC <ArrowUpRight class="size-4" aria-hidden="true" />
        </a>
      </div>
      <button type="button" class="site-button site-button-outline site-button-icon" aria-label="返回顶部" @click="scrollToTop">
        <ArrowUp class="size-4" aria-hidden="true" />
      </button>
    </div>
  </footer>
</template>
