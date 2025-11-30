<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { supabase } from '../supabase/client'
import { useUserStore } from '../stores/user'
import { Download, Calendar, Tag, Trash2, User as UserIcon, Edit2, Heart } from 'lucide-vue-next'
import { useHead } from '@vueuse/head'
import type { Model } from '../types'
import { marked } from 'marked'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const model = ref<Model | null>(null)
const loading = ref(true)
const downloading = ref(false)
const isFavorited = ref(false)
const favoriteLoading = ref(false)

useHead({
  title: computed(() => model.value ? `${model.value.title} - YSM 模型` : '加载中... - YSM 模型站'),
  meta: [
    { 
      name: 'description', 
      content: computed(() => model.value?.description || '查看此 Minecraft YSM 模型的详细信息和预览。') 
    },
    {
      name: 'keywords',
      content: computed(() => model.value?.tags ? `Minecraft, YSM, ${model.value.tags.join(', ')}` : 'Minecraft, YSM')
    }
  ]
})

const isOwner = computed(() => {
  return userStore.user && model.value && userStore.user.id === model.value.user_id
})

const renderedDescription = computed(() => {
  if (!model.value?.description) return '暂无描述'
  return marked(model.value.description)
})

const fetchModel = async () => {
  loading.value = true
  const { data, error } = await supabase
    .from('models')
    .select('*, profiles(username, avatar_url)')
    .eq('id', route.params.id)
    .single()

  if (error) {
    console.error(error)
    router.push('/') // Redirect if not found
  } else {
    model.value = data
    checkFavoriteStatus()
  }
  loading.value = false
}

const checkFavoriteStatus = async () => {
  if (!userStore.user || !model.value) return
  const { data } = await supabase
    .from('favorites')
    .select('id')
    .eq('user_id', userStore.user.id)
    .eq('model_id', model.value.id)
    .single()
  isFavorited.value = !!data
}

const toggleFavorite = async () => {
  if (!userStore.user) {
    router.push('/login')
    return
  }
  if (!model.value || favoriteLoading.value) return

  favoriteLoading.value = true
  try {
    if (isFavorited.value) {
      const { error } = await supabase
        .from('favorites')
        .delete()
        .eq('user_id', userStore.user.id)
        .eq('model_id', model.value.id)
      if (error) throw error
      isFavorited.value = false
    } else {
      const { error } = await supabase
        .from('favorites')
        .insert({ user_id: userStore.user.id, model_id: model.value.id })
      if (error) throw error
      isFavorited.value = true
    }
  } catch (e) {
    console.error('Toggle favorite failed', e)
  } finally {
    favoriteLoading.value = false
  }
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString()
}

const handleDownload = async () => {
  if (!model.value) return
  
  downloading.value = true
  
  try {
    // Get download URL
    const { data } = await supabase.storage
      .from('models')
      .createSignedUrl(model.value.file_path, 60) // 60 seconds valid

    if (data?.signedUrl) {
      // Increment download count via RPC if exists, or manual update
      try {
        await supabase.rpc('increment_downloads', { row_id: model.value.id })
      } catch (rpcError) {
        // Fallback if RPC doesn't exist or fails
        const { data: current } = await supabase.from('models').select('downloads').eq('id', model.value.id).single()
        if (current) {
            await supabase.from('models').update({ downloads: current.downloads + 1 }).eq('id', model.value.id)
        }
      }
        
      // Update local count
      model.value.downloads++

      // Trigger download
      window.open(data.signedUrl, '_blank')
    }
  } catch (e) {
    console.error('Download failed', e)
  } finally {
    downloading.value = false
  }
}

const handleDelete = async () => {
  if (!confirm('你确定要删除这个模型吗？')) return

  try {
    // Delete file from storage
    // Note: We should probably delete the file, but for simplicity just deleting the record is often enough if we have cascading deletes or clean up scripts.
    // But let's try to delete the file.
    if (model.value?.file_path) {
        await supabase.storage.from('models').remove([model.value.file_path])
    }
    if (model.value?.image_url) {
        // Extract path from URL if possible, or just ignore for now as image_url is public URL
        // Ideally we store image_path separately.
    }

    const { error } = await supabase
      .from('models')
      .delete()
      .eq('id', model.value?.id)

    if (error) throw error
    
    router.push('/')
  } catch (e) {
    console.error('Delete failed', e)
    alert('删除模型失败')
  }
}

onMounted(() => {
  fetchModel()
})
</script>

<template>
  <div class="container">
    <div v-if="loading" class="loading">
      正在加载模型详情...
    </div>

    <div v-else-if="model" class="model-detail">
      <div class="model-header">
        <div class="header-content">
          <h1>{{ model.title }}</h1>
          <div class="author-info">
            <div class="avatar-placeholder">
              <img v-if="model.profiles?.avatar_url" :src="model.profiles.avatar_url" class="avatar-img" />
              <UserIcon v-else :size="20" />
            </div>
            <RouterLink :to="`/user/${model.user_id}`" class="author-link">{{ model.profiles?.username || '未知' }}</RouterLink>
            <span class="separator">•</span>
            <span class="date">
              <Calendar :size="14" /> {{ formatDate(model.created_at) }}
            </span>
          </div>
        </div>
        
        <div class="header-actions">
          <button 
            v-if="isOwner" 
            @click="router.push(`/model/${model.id}/edit`)" 
            class="btn btn--secondary"
          >
            <Edit2 :size="18" />
            编辑
          </button>
          <button 
            v-if="isOwner" 
            @click="handleDelete" 
            class="btn btn--danger"
          >
            <Trash2 :size="18" />
            删除
          </button>
          <button 
            @click="toggleFavorite" 
            class="btn" 
            :class="isFavorited ? 'btn--danger' : 'btn--secondary'"
            :disabled="favoriteLoading"
          >
            <Heart :size="18" :fill="isFavorited ? 'currentColor' : 'none'" />
            {{ isFavorited ? '已收藏' : '收藏' }}
          </button>
          <button 
            @click="handleDownload" 
            class="btn btn--primary"
            :disabled="downloading"
          >
            <Download :size="18" />
            {{ downloading ? '下载中...' : '下载模型' }}
          </button>
        </div>
      </div>

      <div class="model-content">
        <div class="main-column">
          <div class="preview-image">
            <img :src="model.image_url || 'https://via.placeholder.com/800x600?text=暂无预览'" :alt="model.title">
          </div>
          
          <div class="description-card">
            <h2>描述</h2>
            <div class="description-text markdown-body" v-html="renderedDescription"></div>
          </div>
        </div>

        <div class="sidebar">
          <div class="sidebar-card">
            <h3>统计</h3>
            <div class="stat-item">
              <span class="label">下载量</span>
              <span class="value">{{ model.downloads }}</span>
            </div>
            <div class="stat-item">
              <span class="label">许可</span>
              <span class="value">公开</span>
            </div>
          </div>

          <div v-if="model.tags && model.tags.length > 0" class="sidebar-card">
            <h3>标签</h3>
            <div class="tags-list">
              <span v-for="tag in model.tags" :key="tag" class="tag">
                <Tag :size="12" /> {{ tag }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@use '../styles/themes/variables' as *;

.model-detail {
  padding-bottom: $spacing-2xl;
}

.model-header {
  display: flex;
  flex-direction: column;
  gap: $spacing-md;
  margin-bottom: $spacing-xl;
  
  @media (min-width: 768px) {
    flex-direction: row;
    justify-content: space-between;
    align-items: flex-start;
  }
}

.header-content {
  h1 {
    font-size: 2rem;
    font-weight: 800;
    color: var(--color-text-main);
    margin-bottom: $spacing-sm;
  }
}

.author-info {
  display: flex;
  align-items: center;
  gap: $spacing-sm;
  color: var(--color-text-muted);
  font-size: 0.875rem;
}

.avatar-placeholder {
  width: 32px;
  height: 32px;
  background-color: #e0e7ff;
  color: var(--color-primary);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

.avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.author-link {
  font-weight: 500;
  color: var(--color-text-main);
  &:hover {
    color: var(--color-primary);
    text-decoration: underline;
  }
}

.separator {
  color: var(--color-border);
}

.date {
  display: flex;
  align-items: center;
  gap: 4px;
}

.header-actions {
  display: flex;
  gap: $spacing-md;
}

.model-content {
  display: grid;
  gap: $spacing-xl;
  
  @media (min-width: 768px) {
    grid-template-columns: 1fr 300px;
  }
}

.preview-image {
  border-radius: $radius-lg;
  overflow: hidden;
  border: 1px solid var(--color-border);
  background-color: #f3f4f6;
  margin-bottom: $spacing-xl;
  
  img {
    width: 100%;
    height: auto;
    display: block;
  }
}

.description-card {
  background-color: white;
  padding: $spacing-xl;
  border-radius: $radius-lg;
  border: 1px solid var(--color-border);
  
  h2 {
    font-size: 1.25rem;
    font-weight: 600;
    margin-bottom: $spacing-md;
    color: var(--color-text-main);
  }
}

.description-text {
  color: var(--color-text-main);
  line-height: 1.6;
}

.markdown-body {
  :deep(h1), :deep(h2), :deep(h3), :deep(h4), :deep(h5), :deep(h6) {
    margin-top: 1.5rem;
    margin-bottom: 1rem;
    font-weight: 600;
    line-height: 1.25;
    color: var(--color-text-main);
  }

  :deep(h1) { font-size: 1.5rem; }
  :deep(h2) { font-size: 1.25rem; border-bottom: 1px solid var(--color-border); padding-bottom: 0.5rem; }
  :deep(h3) { font-size: 1.125rem; }
  
  :deep(p) {
    margin-bottom: 1rem;
  }

  :deep(ul), :deep(ol) {
    margin-bottom: 1rem;
    padding-left: 2rem;
  }

  :deep(li) {
    margin-bottom: 0.25rem;
  }

  :deep(a) {
    color: var(--color-primary);
    text-decoration: underline;
    
    &:hover {
      color: var(--color-primary-hover);
    }
  }

  :deep(code) {
    background-color: #f3f4f6;
    padding: 0.2rem 0.4rem;
    border-radius: 0.25rem;
    font-family: monospace;
    font-size: 0.875rem;
  }

  :deep(pre) {
    background-color: #1f2937;
    color: #f9fafb;
    padding: 1rem;
    border-radius: 0.5rem;
    overflow-x: auto;
    margin-bottom: 1rem;
    
    code {
      background-color: transparent;
      padding: 0;
      color: inherit;
    }
  }

  :deep(blockquote) {
    border-left: 4px solid var(--color-border);
    padding-left: 1rem;
    margin-left: 0;
    margin-bottom: 1rem;
    color: var(--color-text-muted);
  }
  
  :deep(img) {
    max-width: 100%;
    border-radius: 0.5rem;
    margin: 1rem 0;
  }
}

.sidebar {
  display: flex;
  flex-direction: column;
  gap: $spacing-lg;
}

.sidebar-card {
  background-color: white;
  padding: $spacing-lg;
  border-radius: $radius-lg;
  border: 1px solid var(--color-border);
  
  h3 {
    font-size: 1rem;
    font-weight: 600;
    margin-bottom: $spacing-md;
    color: var(--color-text-main);
  }
}

.stat-item {
  display: flex;
  justify-content: space-between;
  padding: $spacing-sm 0;
  border-bottom: 1px solid var(--color-border);
  
  &:last-child {
    border-bottom: none;
  }
  
  .label {
    color: var(--color-text-muted);
  }
  
  .value {
    font-weight: 500;
    color: var(--color-text-main);
  }
}

.tags-list {
  display: flex;
  flex-wrap: wrap;
  gap: $spacing-sm;
}

.tag {
  display: flex;
  align-items: center;
  gap: 4px;
  background-color: #f3f4f6;
  color: var(--color-text-muted);
  padding: 4px 8px;
  border-radius: $radius-md;
  font-size: 0.75rem;
}

.loading {
  text-align: center;
  padding: $spacing-2xl;
  color: var(--color-text-muted);
}
</style>
