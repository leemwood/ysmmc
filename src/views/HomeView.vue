<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { supabase } from '../supabase/client'
import { Search, Download, Heart, ChevronLeft, ChevronRight } from 'lucide-vue-next'
import { useHead } from '@vueuse/head'
import type { Model } from '../types'

useHead({
  title: 'YSM 模型站 - 免费下载 Minecraft Yes Steve Model 模型',
  meta: [
    { name: 'description', content: '浏览并下载社区创作的 Minecraft Yes Steve Model (YSM) 玩家模型。包括二次元、动漫、游戏角色等多种风格。' },
    { name: 'keywords', content: 'Minecraft, YSM, Yes Steve Model, 模型下载, 皮肤, 3D模型, MC模组' }
  ]
})

const models = ref<Model[]>([])
const loading = ref(true)
const searchQuery = ref('')
const page = ref(1)
const pageSize = 12
const totalCount = ref(0)

const fetchModels = async () => {
  loading.value = true
  
  let query = supabase
    .from('models')
    .select('*, profiles(username)', { count: 'exact' })
    .eq('is_public', true)
    .eq('status', 'approved')
    .order('created_at', { ascending: false })

  if (searchQuery.value) {
    query = query.ilike('title', `%${searchQuery.value}%`)
  }

  const from = (page.value - 1) * pageSize
  const to = from + pageSize - 1
  
  query = query.range(from, to)

  const { data, count, error } = await query

  if (error) {
    console.error(error)
  } else {
    models.value = data as any
    totalCount.value = count || 0
  }
  loading.value = false
}

onMounted(() => {
  fetchModels()
})

const handleSearch = () => {
  page.value = 1
  fetchModels()
}

const changePage = (newPage: number) => {
  if (newPage >= 1 && newPage <= Math.ceil(totalCount.value / pageSize)) {
    page.value = newPage
    fetchModels()
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }
}
</script>

<template>
  <div class="container">
    <div class="hero">
      <h1>发现 YSM 模型</h1>
      <p>浏览并下载社区创作的玩家模型</p>
      
      <div class="search-bar">
        <input 
          v-model="searchQuery" 
          @keyup.enter="handleSearch" 
          type="text" 
          placeholder="搜索模型..." 
          class="input"
        >
        <button @click="handleSearch" class="btn btn--primary">
          <Search :size="20" />
        </button>
      </div>
    </div>

    <div v-if="loading" class="loading">
      正在加载模型...
    </div>

    <div v-else-if="models.length === 0" class="empty-state">
      <p>未找到模型。快来成为第一个上传者吧！</p>
      <RouterLink to="/upload" class="btn btn--primary">上传模型</RouterLink>
    </div>

    <div v-else class="model-grid">
      <div 
        v-for="model in models" 
        :key="model.id" 
        class="model-card"
      >
        <RouterLink :to="`/model/${model.id}`" class="model-image-link">
          <div class="model-image">
            <img :src="model.image_url || 'https://via.placeholder.com/400x300?text=暂无预览'" :alt="model.title">
          </div>
        </RouterLink>
        <div class="model-info">
          <RouterLink :to="`/model/${model.id}`" class="model-title-link">
            <h3 class="model-title">{{ model.title }}</h3>
          </RouterLink>
          <RouterLink :to="`/user/${model.user_id}`" class="model-author">作者：{{ model.profiles?.username || '未知' }}</RouterLink>
          <div class="model-meta">
            <span class="meta-item">
              <Download :size="14" /> {{ model.downloads }}
            </span>
            <!-- <span class="meta-item">
              <Heart :size="14" /> 0
            </span> -->
          </div>
        </div>
      </div>
    </div>

    <div v-if="totalCount > pageSize" class="pagination">
      <button 
        @click="changePage(page - 1)" 
        :disabled="page === 1" 
        class="btn btn--secondary pagination-btn"
      >
        <ChevronLeft :size="20" />
      </button>
      <span class="pagination-info">{{ page }} / {{ Math.ceil(totalCount / pageSize) }}</span>
      <button 
        @click="changePage(page + 1)" 
        :disabled="page >= Math.ceil(totalCount / pageSize)" 
        class="btn btn--secondary pagination-btn"
      >
        <ChevronRight :size="20" />
      </button>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@use '../styles/themes/variables' as *;

.hero {
  text-align: center;
  padding: $spacing-2xl 0;
  
  h1 {
    font-size: 2.5rem;
    font-weight: 800;
    margin-bottom: $spacing-md;
    color: var(--color-text-main);
  }
  
  p {
    font-size: 1.125rem;
    color: var(--color-text-muted);
    margin-bottom: $spacing-xl;
  }
}

.search-bar {
  max-width: 600px;
  margin: 0 auto;
  display: flex;
  gap: $spacing-sm;
}

.model-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: $spacing-lg;
  padding-bottom: $spacing-2xl;
}

.model-card {
  background-color: white;
  border: 1px solid var(--color-border);
  border-radius: $radius-lg;
  overflow: hidden;
  transition: $transition-base;
  display: flex;
  flex-direction: column;

  &:hover {
    transform: translateY(-4px);
    box-shadow: $shadow-lg;
    border-color: var(--color-primary);
  }
}

.model-image {
  aspect-ratio: 4/3;
  background-color: #f3f4f6;
  overflow: hidden;

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}

.model-image-link {
  display: block;
}

.model-info {
  padding: $spacing-md;
}

.model-title-link {
  text-decoration: none;
  display: block;
  color: inherit;
  
  &:hover .model-title {
    color: var(--color-primary);
  }
}

.model-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--color-text-main);
  margin-bottom: $spacing-xs;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  transition: color 0.2s;
}

.model-author {
  font-size: 0.875rem;
  color: var(--color-text-muted);
  margin-bottom: $spacing-md;
  display: inline-block;
  text-decoration: none;
  transition: color 0.2s;
  
  &:hover {
    color: var(--color-primary);
  }
}

.model-meta {
  display: flex;
  gap: $spacing-md;
  color: var(--color-text-muted);
  font-size: 0.875rem;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
}

.loading, .empty-state {
  text-align: center;
  padding: $spacing-2xl;
  color: var(--color-text-muted);
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: $spacing-md;
  margin-top: $spacing-xl;
  padding-bottom: $spacing-xl;
}

.pagination-btn {
  padding: $spacing-xs;
  display: flex;
  align-items: center;
  justify-content: center;
  
  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
}

.pagination-info {
  font-weight: 500;
  color: var(--color-text-muted);
}
</style>
