<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { Search, ChevronLeft, ChevronRight, Plus } from 'lucide-vue-next'
import { useHead } from '@vueuse/head'
import type { Model } from '../types'
import LoadingSpinner from '../components/LoadingSpinner.vue'
import ModelCard from '../components/ModelCard.vue'
import { useModelStore } from '../stores/models'

useHead({
  title: 'YSM 模型站 - 免费下载 Minecraft Yes Steve Model 模型',
  meta: [
    { name: 'description', content: '浏览并下载社区创作的 Minecraft Yes Steve Model (YSM) 玩家模型。包括二次元、动漫、游戏角色等多种风格。' },
    { name: 'keywords', content: 'Minecraft, YSM, Yes Steve Model, 模型下载, 皮肤, 3D模型, MC模组' }
  ]
})

const modelStore = useModelStore()
const localModels = ref<Model[]>([])
const loading = ref(true)
const searchQuery = ref('')
const page = ref(1)
const pageSize = 12
const totalCount = ref(0)

const fetchModels = async (force: boolean = false) => {
  loading.value = true
  
  // Use store for first page no search
  if (page.value === 1 && !searchQuery.value) {
    await modelStore.fetchModels(page.value, pageSize, '', force)
    localModels.value = modelStore.models
    totalCount.value = modelStore.totalCount
  } else {
    // For other cases, fetch directly but still use the store's logic helper
    const result = await modelStore.fetchModels(page.value, pageSize, searchQuery.value, true)
    if (result) {
      localModels.value = result.data as any
      totalCount.value = result.count || 0
    }
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
      
      <div class="search-bar" role="search">
        <input 
          v-model="searchQuery" 
          @keyup.enter="handleSearch" 
          type="text" 
          placeholder="搜索模型..." 
          class="input"
          aria-label="搜索模型"
        >
        <button @click="handleSearch" class="btn btn--primary" aria-label="开始搜索">
          <Search :size="20" aria-hidden="true" />
        </button>
      </div>
    </div>

    <LoadingSpinner v-if="loading" message="正在寻找优质模型..." aria-live="polite" :aria-busy="true" />

    <div v-else-if="localModels.length === 0" class="empty-state" role="alert" aria-live="polite">
      <div class="empty-icon" aria-hidden="true">
        <Plus :size="48" />
      </div>
      <h3>暂无模型</h3>
      <p>未找到相关模型。快来成为第一个上传者吧！</p>
      <RouterLink to="/upload" class="btn btn--primary">
        <Plus :size="20" aria-hidden="true" /> 立即上传
      </RouterLink>
    </div>

    <main v-else class="content-section" aria-live="polite">
      <div class="model-grid">
        <ModelCard 
          v-for="(model, index) in localModels" 
          :key="model.id" 
          :model="model"
          :index="index"
        />
      </div>

      <nav v-if="totalCount > pageSize" class="pagination" role="navigation" aria-label="分页导航">
        <button 
          @click="changePage(page - 1)" 
          :disabled="page === 1" 
          class="btn btn--secondary pagination-btn"
          aria-label="上一页"
        >
          <ChevronLeft :size="20" aria-hidden="true" />
        </button>
        <span class="pagination-info" aria-current="page">第 {{ page }} 页 / 共 {{ Math.ceil(totalCount / pageSize) }} 页</span>
        <button 
          @click="changePage(page + 1)" 
          :disabled="page >= Math.ceil(totalCount / pageSize)" 
          class="btn btn--secondary pagination-btn"
          aria-label="下一页"
        >
          <ChevronRight :size="20" aria-hidden="true" />
        </button>
      </nav>
    </main>
  </div>
</template>

<style lang="scss" scoped>
@use '../styles/themes/variables' as *;

.hero {
  text-align: center;
  padding: $spacing-3xl 0;
  background: radial-gradient(circle at top, rgba($color-primary, 0.1) 0%, transparent 70%);
  
  h1 {
    font-size: 3.5rem;
    font-weight: 900;
    margin-bottom: $spacing-md;
    color: var(--color-text-main);
    letter-spacing: -0.025em;
    background: linear-gradient(135deg, var(--color-text-main) 0%, var(--color-primary) 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
  }
  
  p {
    font-size: 1.25rem;
    color: var(--color-text-muted);
    margin-bottom: $spacing-2xl;
    max-width: 600px;
    margin-left: auto;
    margin-right: auto;
  }
}

.search-bar {
  max-width: 600px;
  margin: 0 auto;
  display: flex;
  gap: $spacing-sm;
  background: white;
  padding: $spacing-xs;
  border-radius: $radius-xl;
  box-shadow: $shadow-lg;
  border: 1px solid var(--color-border);
  transition: $transition-base;

  &:focus-within {
    border-color: var(--color-primary);
    box-shadow: $shadow-xl;
    transform: translateY(-2px);
  }

  .input {
    border: none;
    box-shadow: none;
    padding-left: $spacing-lg;
    font-size: 1.125rem;
    
    &:focus {
      box-shadow: none;
    }
  }

  .btn {
    padding: $spacing-md $spacing-xl;
    border-radius: $radius-lg;
  }
}

.pagination-info {
  font-weight: 500;
  color: var(--color-text-muted);
}

.empty-state {
  text-align: center;
  padding: $spacing-3xl 0;
  
  .empty-icon {
    font-size: 4rem;
    margin-bottom: $spacing-md;
  }

  h3 {
    font-size: 1.5rem;
    font-weight: 700;
    margin-bottom: $spacing-xs;
  }

  p {
    color: var(--color-text-muted);
    margin-bottom: $spacing-xl;
  }
}

.content-section {
  animation: fadeIn 0.5s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.model-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: $spacing-xl;
  padding-bottom: $spacing-3xl;
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
</style>
