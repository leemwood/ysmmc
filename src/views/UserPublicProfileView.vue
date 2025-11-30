<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { supabase } from '../supabase/client'
import { User as UserIcon, Package, Heart } from 'lucide-vue-next'
import { useHead } from '@vueuse/head'
import type { Model, Profile } from '../types'

const route = useRoute()
const profile = ref<Profile | null>(null)
const uploadedModels = ref<Model[]>([])
const favoritedModels = ref<Model[]>([])
const loading = ref(true)
const activeTab = ref<'uploads' | 'favorites'>('uploads')

const userId = computed(() => route.params.id as string)

useHead({
  title: computed(() => profile.value ? `${profile.value.username} 的个人主页 - YSM 模型站` : '加载中...'),
})

const fetchUserProfile = async () => {
  loading.value = true
  
  // 1. Fetch Profile
  const { data: profileData, error: profileError } = await supabase
    .from('profiles')
    .select('*')
    .eq('id', userId.value)
    .single()
    
  if (profileError || !profileData) {
    console.error(profileError)
    loading.value = false
    return
  }
  profile.value = profileData

  // 2. Fetch Uploaded Models (Public & Approved)
  const { data: uploadsData } = await supabase
    .from('models')
    .select('*')
    .eq('user_id', userId.value)
    .eq('is_public', true)
    .eq('status', 'approved')
    .order('created_at', { ascending: false })
    
  uploadedModels.value = uploadsData as any || []

  // 3. Fetch Favorited Models
  const { data: favoritesData } = await supabase
    .from('favorites')
    .select('model_id, models(*, profiles(username))')
    .eq('user_id', userId.value)
    .order('created_at', { ascending: false })
    
  if (favoritesData) {
    // Filter out nulls or deleted models just in case
    favoritedModels.value = favoritesData
      .map((f: any) => f.models)
      .filter((m: any) => m && m.status === 'approved' && m.is_public)
  }

  loading.value = false
}

onMounted(() => {
  fetchUserProfile()
})
</script>

<template>
  <div class="container">
    <div v-if="loading" class="loading">
      正在加载用户资料...
    </div>

    <div v-else-if="!profile" class="empty-state">
      <p>未找到该用户</p>
    </div>

    <div v-else class="profile-page">
      <!-- Header -->
      <div class="profile-header">
        <div class="avatar-wrapper">
          <img v-if="profile.avatar_url" :src="profile.avatar_url" class="avatar-img">
          <div v-else class="avatar-placeholder">
            <UserIcon :size="64" />
          </div>
        </div>
        <div class="info">
          <h1>{{ profile.username }}</h1>
          <p class="bio">{{ profile.bio || '这位用户很懒，什么都没写。' }}</p>
          <p class="join-date">加入于 {{ new Date(profile.created_at).toLocaleDateString() }}</p>
        </div>
      </div>

      <!-- Tabs -->
      <div class="tabs">
        <button 
          @click="activeTab = 'uploads'" 
          :class="['tab-btn', activeTab === 'uploads' ? 'active' : '']"
        >
          <Package :size="18" /> 上传的模型 ({{ uploadedModels.length }})
        </button>
        <button 
          @click="activeTab = 'favorites'" 
          :class="['tab-btn', activeTab === 'favorites' ? 'active' : '']"
        >
          <Heart :size="18" /> 收藏的模型 ({{ favoritedModels.length }})
        </button>
      </div>

      <!-- Content -->
      <div class="content-area">
        <div v-if="activeTab === 'uploads'">
          <div v-if="uploadedModels.length === 0" class="empty-tab">
            <Package :size="48" />
            <p>该用户还没有上传任何模型</p>
          </div>
          <div v-else class="model-grid">
            <div 
              v-for="model in uploadedModels" 
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
                <div class="model-meta">
                  <span>{{ new Date(model.created_at).toLocaleDateString() }}</span>
                  <span>•</span>
                  <span>{{ model.downloads }} 下载</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div v-if="activeTab === 'favorites'">
          <div v-if="favoritedModels.length === 0" class="empty-tab">
            <Heart :size="48" />
            <p>该用户还没有收藏任何模型</p>
          </div>
          <div v-else class="model-grid">
            <div 
              v-for="model in favoritedModels" 
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
                <RouterLink :to="`/user/${model.user_id}`" class="model-author">作者: {{ model.profiles?.username || '未知' }}</RouterLink>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@use '../styles/themes/variables' as *;

.profile-header {
  display: flex;
  align-items: center;
  gap: $spacing-xl;
  padding: $spacing-2xl;
  background-color: white;
  border-radius: $radius-lg;
  border: 1px solid var(--color-border);
  margin-bottom: $spacing-xl;
  
  @media (max-width: 600px) {
    flex-direction: column;
    text-align: center;
  }
}

.avatar-wrapper {
  width: 120px;
  height: 120px;
  flex-shrink: 0;
  
  .avatar-img {
    width: 100%;
    height: 100%;
    border-radius: 50%;
    object-fit: cover;
    border: 4px solid white;
    box-shadow: $shadow-md;
  }
  
  .avatar-placeholder {
    width: 100%;
    height: 100%;
    background-color: #e0e7ff;
    color: var(--color-primary);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
  }
}

.info {
  h1 {
    font-size: 2rem;
    font-weight: 800;
    color: var(--color-text-main);
    margin-bottom: $spacing-xs;
  }
  
  .bio {
    font-size: 1.125rem;
    color: var(--color-text-main);
    margin-bottom: $spacing-md;
    line-height: 1.6;
  }
  
  .join-date {
    color: var(--color-text-muted);
    font-size: 0.875rem;
  }
}

.tabs {
  display: flex;
  gap: $spacing-lg;
  border-bottom: 1px solid var(--color-border);
  margin-bottom: $spacing-xl;
}

.tab-btn {
  display: flex;
  align-items: center;
  gap: $spacing-sm;
  padding: $spacing-md $spacing-sm;
  background: none;
  border: none;
  border-bottom: 2px solid transparent;
  font-size: 1rem;
  font-weight: 500;
  color: var(--color-text-muted);
  cursor: pointer;
  transition: all 0.2s;
  
  &:hover {
    color: var(--color-text-main);
  }
  
  &.active {
    color: var(--color-primary);
    border-bottom-color: var(--color-primary);
  }
}

.model-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: $spacing-lg;
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
  font-size: 1rem;
  font-weight: 600;
  color: var(--color-text-main);
  margin-bottom: $spacing-xs;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  transition: color 0.2s;
}

.model-meta {
  font-size: 0.875rem;
  color: var(--color-text-muted);
}

.model-author {
  font-size: 0.875rem;
  color: var(--color-text-muted);
  display: inline-block;
  text-decoration: none;
  transition: color 0.2s;
  
  &:hover {
    color: var(--color-primary);
  }
}

.empty-tab {
  text-align: center;
  padding: $spacing-2xl;
  color: var(--color-text-muted);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: $spacing-md;
}

.loading {
  text-align: center;
  padding: $spacing-2xl;
  color: var(--color-text-muted);
}
</style>
