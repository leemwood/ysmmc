<script setup lang="ts">
import LoadingSpinner from '../components/LoadingSpinner.vue'
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { supabase } from '../supabase/client'
import { useUserStore } from '../stores/user'
import { useModelStore } from '../stores/models'
import { Download, Calendar, Tag, Trash2, Edit2, Heart, Share2, FileCode, CheckCircle, Loader2, ShieldAlert } from 'lucide-vue-next'
import { useHead } from '@vueuse/head'
import type { Model } from '../types'
import { marked } from 'marked'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const modelStore = useModelStore()

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
  const data = await modelStore.fetchModelById(route.params.id as string)

  if (!data) {
    router.push('/') // Redirect if not found
  } else {
    model.value = data
    checkFavoriteStatus()
  }
  loading.value = false
}

const checkFavoriteStatus = async () => {
  if (!userStore.user || !model.value) return
  const { data, error } = await supabase
    .from('favorites')
    .select('id')
    .eq('user_id', userStore.user.id)
    .eq('model_id', model.value.id)
    .maybeSingle()
  
  if (error) {
    console.error('Error checking favorite status:', error)
    return
  }
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
  <div class="container py-xl">
    <LoadingSpinner v-if="loading" message="正在获取模型详情..." />

    <div v-else-if="!model" class="error-container" role="alert" aria-live="assertive">
      <div class="error-icon" aria-hidden="true">
        <ShieldAlert :size="64" />
      </div>
      <h2>模型未找到</h2>
      <p>该模型可能已被删除或不存在。</p>
      <RouterLink to="/" class="btn btn--primary">返回首页</RouterLink>
    </div>

    <div v-else class="detail-wrapper">
      <header class="model-header">
        <div class="header-main">
          <h1 class="model-title">{{ model.title }}</h1>
          <div class="model-meta-info">
            <RouterLink :to="`/user/${model.user_id}`" class="author-tag" :aria-label="`查看作者 ${model.profiles?.username || '未知作者'} 的主页`">
              <div class="author-avatar" role="img" :aria-label="model.profiles?.username?.charAt(0) || '头像'">
                {{ model.profiles?.username?.charAt(0).toUpperCase() || '?' }}
              </div>
              <span>{{ model.profiles?.username || '未知作者' }}</span>
            </RouterLink>
            <span class="meta-divider" aria-hidden="true"></span>
            <span class="meta-date" :aria-label="`发布于 ${formatDate(model.created_at)}`">
              <Calendar :size="14" aria-hidden="true" /> {{ formatDate(model.created_at) }}
            </span>
          </div>
        </div>
        <div v-if="isOwner || userStore.isAdmin" class="header-actions">
          <button @click="router.push(`/model/${model.id}/edit`)" class="btn btn--secondary btn--sm" aria-label="编辑模型">
            <Edit2 :size="16" aria-hidden="true" /> 编辑
          </button>
          <button @click="handleDelete" class="btn btn--danger btn--sm" aria-label="删除模型">
            <Trash2 :size="16" aria-hidden="true" /> 删除
          </button>
        </div>
      </header>

      <div class="model-content">
        <div class="main-column">
          <div class="image-preview-card">
            <img 
              :src="model.image_url || 'https://via.placeholder.com/800x600?text=暂无预览'" 
              :alt="`${model.title} 的预览图`"
              loading="lazy"
              @load="(e) => (e.target as HTMLImageElement).classList.add('loaded')"
            >
          </div>
          
          <div class="description-card">
            <h2>模型描述</h2>
            <div class="description-text markdown-body" v-html="renderedDescription" aria-label="模型详细描述"></div>
          </div>
        </div>

        <aside class="side-column">
          <div class="action-card">
            <button @click="handleDownload" class="btn btn--primary download-btn" :disabled="downloading" :aria-busy="downloading" :aria-label="downloading ? '正在准备下载' : '立即下载模型'">
              <template v-if="downloading">
                <Loader2 class="animate-spin" :size="20" aria-hidden="true" /> 准备中...
              </template>
              <template v-else>
                <Download :size="20" aria-hidden="true" /> 立即下载
              </template>
            </button>
            <div class="secondary-actions">
              <button 
                @click="toggleFavorite" 
                class="btn" 
                :class="isFavorited ? 'btn--danger' : 'btn--secondary'"
                :disabled="favoriteLoading"
                :aria-pressed="isFavorited"
                :aria-label="isFavorited ? '取消收藏' : '收藏模型'"
              >
                <Heart :size="18" :fill="isFavorited ? 'currentColor' : 'none'" aria-hidden="true" />
                {{ isFavorited ? '已收藏' : '收藏' }}
              </button>
              <button class="btn btn--secondary" title="分享" aria-label="分享模型">
                <Share2 :size="18" aria-hidden="true" /> 分享
              </button>
            </div>
          </div>

          <div class="info-card">
            <h3>模型信息</h3>
            <div class="info-list" role="list">
              <div class="info-item" role="listitem">
                <span class="label"><Download :size="14" aria-hidden="true" /> 下载量</span>
                <span class="value" :aria-label="`${model.downloads} 次下载`">{{ model.downloads }} 次</span>
              </div>
              <div class="info-item" role="listitem">
                <span class="label"><FileCode :size="14" aria-hidden="true" /> 文件名</span>
                <span class="value">{{ model.file_path.split('/').pop() }}</span>
              </div>
              <div class="info-item" role="listitem">
                <span class="label"><CheckCircle :size="14" aria-hidden="true" /> 状态</span>
                <span class="value" :class="`status--${model.status}`" :aria-label="`当前状态: ${model.status === 'approved' ? '已发布' : '审核中'}`">
                  {{ model.status === 'approved' ? '已发布' : '审核中' }}
                </span>
              </div>
            </div>
          </div>

          <div v-if="model.tags && model.tags.length > 0" class="info-card">
            <h3>标签</h3>
            <div class="tags-list" role="list" aria-label="模型标签">
              <span v-for="tag in model.tags" :key="tag" class="tag" role="listitem">
                <Tag :size="12" aria-hidden="true" /> {{ tag }}
              </span>
            </div>
          </div>
        </aside>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@use '../styles/themes/variables' as *;

.model-detail {
  padding-bottom: $spacing-2xl;
}

.error-container {
  text-align: center;
  padding: $spacing-3xl 0;
  
  .error-icon {
    font-size: 4rem;
    margin-bottom: $spacing-md;
  }

  h2 {
    font-size: 2rem;
    font-weight: 700;
    margin-bottom: $spacing-md;
  }

  p {
    color: var(--color-text-muted);
    margin-bottom: $spacing-xl;
  }
}

.model-header {
  margin-bottom: $spacing-2xl;
  padding: $spacing-2xl;
  background: white;
  border-radius: $radius-xl;
  border: 1px solid var(--color-border);
  box-shadow: $shadow-md;
  display: flex;
  justify-content: space-between;
  align-items: flex-start;

  .model-title {
    font-size: 2.5rem;
    font-weight: 800;
    margin-bottom: $spacing-md;
    color: var(--color-text-main);
    letter-spacing: -0.025em;
    line-height: 1.2;
  }

  .model-meta-info {
    display: flex;
    align-items: center;
    gap: $spacing-md;
    color: var(--color-text-muted);
    font-size: 0.875rem;

    .author-tag {
      display: flex;
      align-items: center;
      gap: $spacing-sm;
      color: var(--color-text-main);
      text-decoration: none;
      font-weight: 600;
      transition: $transition-base;

      &:hover {
        color: var(--color-primary);
      }
    }

    .author-avatar {
      width: 28px;
      height: 28px;
      background-color: var(--color-primary);
      color: white;
      border-radius: $radius-full;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 0.75rem;
      font-weight: 700;
    }

    .meta-divider {
      width: 4px;
      height: 4px;
      background-color: var(--color-border);
      border-radius: $radius-full;
    }

    .meta-date {
      display: flex;
      align-items: center;
      gap: $spacing-xs;
    }
  }

  .header-actions {
    display: flex;
    gap: $spacing-sm;
  }
}

.model-content {
  display: grid;
  grid-template-columns: 1fr;
  gap: $spacing-2xl;

  @media (min-width: 1024px) {
    grid-template-columns: 2fr 1fr;
  }
}

.main-column {
  display: flex;
  flex-direction: column;
  gap: $spacing-2xl;
}

.image-preview-card {
  background: white;
  border-radius: $radius-xl;
  overflow: hidden;
  border: 1px solid var(--color-border);
  box-shadow: $shadow-lg;
  aspect-ratio: 16 / 9;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f1f5f9;

  img {
    width: 100%;
    height: 100%;
    object-fit: contain;
    transition: $transition-slow;
    opacity: 0;

    &.loaded {
      opacity: 1;
    }
  }
}

.description-card {
  background: white;
  padding: $spacing-2xl;
  border-radius: $radius-xl;
  border: 1px solid var(--color-border);
  box-shadow: $shadow-sm;

  h2 {
    font-size: 1.5rem;
    font-weight: 700;
    margin-bottom: $spacing-lg;
    display: flex;
    align-items: center;
    gap: $spacing-sm;
    color: var(--color-text-main);
    
    &::before {
      content: '';
      width: 4px;
      height: 24px;
      background: var(--color-primary);
      border-radius: $radius-full;
    }
  }

  .description-text {
    color: var(--color-text-main);
    line-height: 1.8;
    font-size: 1.05rem;
  }
}

.side-column {
  display: flex;
  flex-direction: column;
  gap: $spacing-xl;
}

.action-card {
  background: white;
  padding: $spacing-xl;
  border-radius: $radius-xl;
  border: 1px solid var(--color-border);
  box-shadow: $shadow-lg;
  position: sticky;
  top: 6rem;

  .download-btn {
    width: 100%;
    padding: $spacing-md;
    font-size: 1.125rem;
    margin-bottom: $spacing-md;
    height: 3.5rem;
  }

  .secondary-actions {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: $spacing-sm;

    .btn {
      height: 3rem;
    }
  }
}

.info-card {
  background: white;
  padding: $spacing-xl;
  border-radius: $radius-xl;
  border: 1px solid var(--color-border);
  box-shadow: $shadow-sm;

  h3 {
    font-size: 1.125rem;
    font-weight: 700;
    margin-bottom: $spacing-lg;
    color: var(--color-text-main);
    padding-bottom: $spacing-sm;
    border-bottom: 1px solid var(--color-border);
  }

  .info-list {
    display: flex;
    flex-direction: column;
    gap: $spacing-md;
  }

  .info-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 0.935rem;

    .label {
      color: var(--color-text-muted);
      display: flex;
      align-items: center;
      gap: $spacing-xs;
    }

    .value {
      font-weight: 600;
      color: var(--color-text-main);

      &.status--published {
        color: #10b981;
      }
      &.status--pending {
        color: #f59e0b;
      }
    }
  }

  .tags-list {
    display: flex;
    flex-wrap: wrap;
    gap: $spacing-xs;

    .tag {
      display: inline-flex;
      align-items: center;
      gap: 4px;
      padding: 4px 10px;
      background-color: #f1f5f9;
      color: var(--color-text-muted);
      border-radius: $radius-full;
      font-size: 0.75rem;
      font-weight: 600;
      transition: $transition-base;

      &:hover {
        background-color: rgba($color-primary, 0.1);
        color: var(--color-primary);
      }
    }
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
