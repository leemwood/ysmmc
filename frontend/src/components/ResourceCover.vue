<script setup lang="ts">
import { ImageOff } from 'lucide-vue-next'

// 资源封面：有图时撑满容器（object-fit: contain，见 main.css 的
// .resource-cover-img），无图时渲染统一的占位提示。外层定高容器由调用方提供。
withDefaults(
  defineProps<{
    cover: string
    title: string
    imgClass?: string
    iconClass?: string
  }>(),
  { imgClass: '', iconClass: 'h-10 w-10' },
)
</script>

<template>
  <img
    v-if="cover"
    :src="cover"
    :alt="title"
    loading="lazy"
    decoding="async"
    class="resource-cover-img"
    :class="imgClass"
  />
  <div v-else class="flex h-full flex-col items-center justify-center gap-2 text-muted-foreground">
    <ImageOff :class="iconClass" class="opacity-50" />
    <span class="text-sm">无预览图</span>
  </div>
</template>
