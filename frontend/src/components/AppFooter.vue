<script setup lang="ts">
</script>

<template>
  <footer class="border-t bg-muted/30">
    <div class="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
      <div class="flex flex-col items-center justify-between gap-6 md:flex-row">
        <div class="flex flex-col items-center md:items-start gap-2">
          <p class="text-base font-semibold gradient-text">
            YSM 模型站 (YSMMC)
          </p>
          <p class="text-sm text-muted-foreground text-center md:text-left">
            非营利性公益模型分享平台，基于 GPL-3.0 协议开源
          </p>
        </div>
        <div class="flex flex-col sm:flex-row items-center gap-4 sm:gap-6">
          <a 
            href="https://qm.qq.com/q/SUKmYH7RyW" 
            target="_blank" 
            rel="noopener noreferrer"
            class="text-sm text-muted-foreground hover:text-primary transition-colors"
          >
            用户交流群
          </a>
          <p class="text-xs text-muted-foreground">
            &copy; {{ new Date().getFullYear() }} YSMMC. All rights reserved.
          </p>
        </div>
      </div>
    </div>
  </footer>
</template>
