<script setup lang="ts">
import { RouterLink } from 'vue-router'
import { Download } from 'lucide-vue-next'
import type { Model } from '../types'

defineProps<{
  model: Model
  index: number
}>()
</script>

<template>
  <div 
    class="model-card"
    :style="{ animationDelay: `${index * 50}ms` }"
    role="article"
    :aria-label="`模型: ${model.title}`"
  >
    <RouterLink :to="`/model/${model.id}`" class="model-image-link" :aria-label="`查看 ${model.title} 的详情`">
      <div class="model-image">
        <img 
          :src="model.image_url || 'https://via.placeholder.com/400x300?text=暂无预览'" 
          :alt="`${model.title} 的预览图`"
          loading="lazy"
          @load="(e) => (e.target as HTMLImageElement).classList.add('loaded')"
        >
        <div class="model-overlay" aria-hidden="true">
          <span class="btn btn--primary btn--sm">查看详情</span>
        </div>
      </div>
    </RouterLink>
    <div class="model-info">
      <div class="model-header-row">
        <RouterLink :to="`/model/${model.id}`" class="model-title-link">
          <h3 class="model-title">{{ model.title }}</h3>
        </RouterLink>
      </div>
      <div class="model-footer-row">
        <RouterLink :to="`/user/${model.user_id}`" class="model-author" :aria-label="`作者: ${model.profiles?.username || '未知'}`">
          <div class="author-avatar" aria-hidden="true">
            {{ model.profiles?.username?.charAt(0).toUpperCase() || '?' }}
          </div>
          <span>{{ model.profiles?.username || '未知' }}</span>
        </RouterLink>
        <div class="model-meta">
          <span class="meta-item" :title="`下载量: ${model.downloads}`" :aria-label="`${model.downloads} 次下载`">
            <Download :size="14" aria-hidden="true" /> {{ model.downloads }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@use '../styles/themes/variables' as *;

.model-card {
  background-color: white;
  border: 1px solid var(--color-border);
  border-radius: $radius-xl;
  overflow: hidden;
  transition: $transition-slow;
  display: flex;
  flex-direction: column;
  opacity: 0;
  animation: cardFadeIn 0.5s ease-out forwards;

  &:hover {
    transform: translateY(-8px);
    box-shadow: $shadow-xl;
    border-color: rgba($color-primary, 0.3);

    .model-image img {
      transform: scale(1.05);
    }

    .model-overlay {
      opacity: 1;
    }
  }
}

@keyframes cardFadeIn {
  to { opacity: 1; transform: translateY(0); }
}

.model-image {
  position: relative;
  aspect-ratio: 4/3;
  overflow: hidden;
  background-color: #f1f5f9;

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    transition: $transition-slow;
    opacity: 0;
    
    &.loaded {
      opacity: 1;
    }
  }
}

.model-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: $transition-base;
  backdrop-filter: blur(2px);

  .btn--sm {
    padding: $spacing-xs $spacing-lg;
    font-size: 0.875rem;
  }
}

.model-info {
  padding: $spacing-lg;
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: $spacing-md;
}

.model-header-row {
  .model-title {
    font-size: 1.25rem;
    font-weight: 700;
    color: var(--color-text-main);
    line-height: 1.4;
    transition: color 0.2s;
    
    &:hover {
      color: var(--color-primary);
    }
  }
}

.model-footer-row {
  margin-top: auto;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: $spacing-md;
  border-top: 1px solid var(--color-border);
}

.model-author {
  display: flex;
  align-items: center;
  gap: $spacing-sm;
  color: var(--color-text-muted);
  font-size: 0.875rem;
  font-weight: 500;
  transition: $transition-base;
  text-decoration: none;

  &:hover {
    color: var(--color-primary);
    
    .author-avatar {
      background-color: var(--color-primary);
      color: white;
    }
  }

  .author-avatar {
    width: 24px;
    height: 24px;
    background-color: #e2e8f0;
    border-radius: $radius-full;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 0.75rem;
    font-weight: 700;
    color: var(--color-text-muted);
    transition: $transition-base;
  }
}

.model-meta {
  display: flex;
  gap: $spacing-md;

  .meta-item {
    display: flex;
    align-items: center;
    gap: $spacing-xs;
    color: var(--color-text-muted);
    font-size: 0.875rem;
  }
}

.model-title-link {
  text-decoration: none;
  display: block;
  color: inherit;
}

.model-image-link {
  display: block;
}
</style>
