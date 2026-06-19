<script setup lang="ts">
import { RouterLink } from 'vue-router'
import { Box } from 'lucide-vue-next'
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card'

interface Props {
  title?: string
  description?: string
}

defineProps<Props>()
</script>

<template>
  <div class="min-h-[calc(100vh-4rem)] flex items-center justify-center px-4 py-6 sm:py-12 relative overflow-hidden">
    <div class="absolute inset-0 -z-10">
      <div class="absolute top-1/4 left-1/4 w-64 h-64 sm:w-96 sm:h-96 bg-primary/5 rounded-full blur-3xl"></div>
      <div class="absolute bottom-1/4 right-1/4 w-64 h-64 sm:w-96 sm:h-96 bg-primary/10 rounded-full blur-3xl"></div>
    </div>

    <Card class="w-full max-w-md animate-fade-in relative">
      <CardHeader class="text-center space-y-2 px-4 py-4 sm:px-6 sm:pt-6">
        <RouterLink
          to="/"
          class="inline-flex items-center justify-center gap-2 text-xl font-bold gradient-text mb-1 focus-ring rounded-lg"
        >
          <Box class="h-6 w-6" />
          YSM 模型站
        </RouterLink>
        <CardTitle v-if="title" class="text-xl sm:text-2xl">{{ title }}</CardTitle>
        <CardDescription v-if="description" class="text-xs sm:text-sm">
          {{ description }}
        </CardDescription>
      </CardHeader>
      <CardContent class="px-4 pb-6 sm:px-6">
        <slot />
      </CardContent>
    </Card>
  </div>
</template>
