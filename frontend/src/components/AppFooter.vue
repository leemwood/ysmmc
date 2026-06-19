<script setup lang="ts">
import { Github, MessageCircle, ArrowUp } from 'lucide-vue-next'

function scrollToTop() {
  window.scrollTo({ top: 0, behavior: 'smooth' })
}
</script>

<template>
  <footer class="border-t bg-muted/30">
    <div class="container-app py-10">
      <!-- 品牌与社交 -->
      <div class="flex flex-col items-center justify-between gap-6 md:flex-row md:items-start">
        <div class="flex flex-col items-center gap-2 text-center md:items-start md:text-left">
          <p class="text-lg font-semibold gradient-text">
            YSM 模型站 (YSMMC)
          </p>
          <p class="max-w-sm text-sm text-muted-foreground">
            非营利性公益模型分享平台，基于 GPL-3.0 协议开源
          </p>
        </div>

        <div class="flex items-center gap-3">
          <a
            href="#"
            target="_blank"
            rel="noopener noreferrer"
            aria-label="GitHub"
            class="focus-ring btn-press inline-flex h-10 items-center gap-2 rounded-full bg-muted/50 px-4 text-sm text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
          >
            <Github class="h-4 w-4" />
            <span>GitHub</span>
          </a>
          <a
            href="https://qm.qq.com/q/SUKmYH7RyW"
            target="_blank"
            rel="noopener noreferrer"
            aria-label="QQ 用户交流群"
            class="focus-ring btn-press inline-flex h-10 items-center gap-2 rounded-full bg-muted/50 px-4 text-sm text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
          >
            <MessageCircle class="h-4 w-4" />
            <span>QQ</span>
          </a>
        </div>
      </div>

      <!-- 版权与返回顶部 -->
      <div class="mt-8 flex flex-col items-center justify-between gap-3 border-t pt-6 text-xs text-muted-foreground sm:flex-row sm:text-sm">
        <p>&copy; {{ new Date().getFullYear() }} YSMMC. All rights reserved.</p>
        <button
          type="button"
          class="focus-ring btn-press inline-flex items-center gap-1.5 rounded-md px-2 py-1 transition-colors hover:bg-muted hover:text-foreground"
          @click="scrollToTop"
        >
          <ArrowUp class="h-3.5 w-3.5" />
          返回顶部
        </button>
      </div>
    </div>
  </footer>
</template>
