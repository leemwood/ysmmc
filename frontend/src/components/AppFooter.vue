<script setup lang="ts">
import { Github, MessageCircle, ArrowUp } from 'lucide-vue-next'

function scrollToTop() {
  window.scrollTo({ top: 0, behavior: 'smooth' })
}
</script>

<template>
  <footer class="border-t bg-muted/30">
    <div class="container-app py-8">
      <div class="flex flex-col items-center gap-6 md:flex-row md:justify-between">
        <div class="flex flex-col items-center gap-2 text-center md:items-start md:text-left">
          <p class="text-base font-semibold gradient-text">
            YSM 模型站 (YSMMC)
          </p>
          <p class="text-sm text-muted-foreground">
            非营利性公益模型分享平台，基于 GPL-3.0 协议开源
          </p>
        </div>

        <div class="flex flex-col items-center gap-4 md:flex-row md:gap-6">
          <div class="flex items-center gap-3">
            <a
              href="#"
              target="_blank"
              rel="noopener noreferrer"
              aria-label="GitHub"
              class="inline-flex h-10 w-10 items-center justify-center rounded-full bg-muted/50 text-muted-foreground transition-colors hover:bg-muted hover:text-foreground focus-ring"
            >
              <Github class="h-5 w-5" />
            </a>
            <a
              href="https://qm.qq.com/q/SUKmYH7RyW"
              target="_blank"
              rel="noopener noreferrer"
              aria-label="QQ 用户交流群"
              class="inline-flex h-10 w-10 items-center justify-center rounded-full bg-muted/50 text-muted-foreground transition-colors hover:bg-muted hover:text-foreground focus-ring"
            >
              <MessageCircle class="h-5 w-5" />
            </a>
          </div>
          <p class="text-xs text-muted-foreground">
            &copy; {{ new Date().getFullYear() }} YSMMC. All rights reserved.
          </p>
        </div>
      </div>
    </div>
  </footer>

  <button
    type="button"
    aria-label="返回顶部"
    class="fixed bottom-6 right-6 z-30 inline-flex h-10 w-10 items-center justify-center rounded-full bg-primary text-primary-foreground shadow-lg transition-transform btn-press focus-ring pb-safe md:hidden"
    @click="scrollToTop"
  >
    <ArrowUp class="h-5 w-5" />
  </button>
</template>
