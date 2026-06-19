<script setup lang="ts">
import { RouterLink } from 'vue-router'
import type { Model } from '@/types'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Download, ImageOff } from 'lucide-vue-next'
import { getModelImageUrl } from '@/utils/image'

interface Props {
  model: Model
}

const props = defineProps<Props>()

function formatDate(date: string) {
  return new Date(date).toLocaleDateString('zh-CN')
}

function getImageSrc(model: Model): string {
  return getModelImageUrl(model.image_id, model.image_url)
}
</script>

<template>
  <RouterLink :to="`/model/${model.id}`" class="block rounded-xl focus-ring">
    <Card class="overflow-hidden card-hover group active:scale-[0.98] transition-transform duration-200">
      <div class="aspect-[4/3] w-full bg-muted overflow-hidden">
        <img
          v-if="model.image_id || model.image_url"
          :src="getImageSrc(model)"
          :alt="model.title"
          loading="lazy"
          decoding="async"
          class="h-full w-full object-cover transition-transform duration-300 group-hover:scale-105"
        />
        <div v-else class="flex h-full flex-col items-center justify-center gap-2 text-muted-foreground">
          <ImageOff class="h-10 w-10 opacity-50" />
          <span class="text-sm">无预览图</span>
        </div>
      </div>
      <CardContent class="p-4">
        <h3 class="font-semibold line-clamp-1 text-balance group-hover:text-primary transition-colors">
          {{ model.title }}
        </h3>
        <p class="mt-1 text-sm text-muted-foreground line-clamp-2">
          {{ model.description || '暂无描述' }}
        </p>
        <div class="mt-3 flex flex-wrap gap-1.5">
          <Badge
            v-for="tag in model.tags?.slice(0, 3)"
            :key="tag"
            variant="secondary"
            class="rounded-full text-xs"
          >
            {{ tag }}
          </Badge>
        </div>
        <div class="mt-3 flex items-center justify-between text-xs text-muted-foreground">
          <div class="flex items-center gap-3">
            <span class="flex items-center gap-1">
              <Download class="h-3 w-3" />
              {{ model.downloads }}
            </span>
          </div>
          <span>{{ formatDate(model.created_at) }}</span>
        </div>
      </CardContent>
    </Card>
  </RouterLink>
</template>
