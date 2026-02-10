<script setup lang="ts">
import LoadingSpinner from '../components/LoadingSpinner.vue'
import { ref, onMounted, computed } from 'vue'
import { supabase } from '../supabase/client'
import { Check, X, User, Box, Search, Trash2, Shield, ShieldAlert, Edit, Users, FolderTree } from 'lucide-vue-next'
import { useHead } from '@vueuse/head'
import type { Model, Profile } from '../types'

useHead({
  title: '管理员控制台 - YSM 模型站',
  meta: [
    { name: 'robots', content: 'noindex, nofollow' }
  ]
})

const activeTab = ref<'audit_models' | 'audit_profiles' | 'manage_models' | 'manage_users'>('audit_models')

const pendingModels = ref<Model[]>([])
const pendingProfiles = ref<Profile[]>([])
const allModels = ref<Model[]>([])
const allProfiles = ref<Profile[]>([])
const loading = ref(true)

// Search filters
const modelSearch = ref('')
const userSearch = ref('')

const filteredModels = computed(() => {
  if (!modelSearch.value) return allModels.value
  const query = modelSearch.value.toLowerCase()
  return allModels.value.filter(m => 
    m.title.toLowerCase().includes(query) || 
    m.profiles?.username?.toLowerCase().includes(query)
  )
})

const filteredProfiles = computed(() => {
  if (!userSearch.value) return allProfiles.value
  const query = userSearch.value.toLowerCase()
  return allProfiles.value.filter(p => 
    p.username.toLowerCase().includes(query) || 
    p.id.toLowerCase().includes(query)
  )
})

const fetchPendingModels = async () => {
  loading.value = true
  const { data, error } = await supabase
    .from('models')
    .select('*, profiles(username)')
    .or('status.eq.pending,update_status.eq.pending_review')
    .order('created_at', { ascending: false })

  if (error) console.error(error)
  else pendingModels.value = data as any
  loading.value = false
}

const fetchPendingProfiles = async () => {
  loading.value = true
  const { data, error } = await supabase
    .from('profiles')
    .select('*')
    .eq('profile_status', 'pending_review')
    .order('updated_at', { ascending: false })

  if (error) console.error(error)
  else pendingProfiles.value = data as any
  loading.value = false
}

const fetchAllModels = async () => {
  loading.value = true
  const { data, error } = await supabase
    .from('models')
    .select('*, profiles(username)')
    .order('created_at', { ascending: false })

  if (error) console.error(error)
  else allModels.value = data as any
  loading.value = false
}

const fetchAllProfiles = async () => {
  loading.value = true
  const { data, error } = await supabase
    .from('profiles')
    .select('*')
    .order('created_at', { ascending: false })

  if (error) console.error(error)
  else allProfiles.value = data as any
  loading.value = false
}

const handleApproveModel = async (model: Model) => {
  // Case 1: New Model Approval
  if (model.status === 'pending') {
    const { error } = await supabase
      .from('models')
      .update({ status: 'approved' })
      .eq('id', model.id)

    if (!error) {
      pendingModels.value = pendingModels.value.filter(m => m.id !== model.id)
    }
  } 
  // Case 2: Model Update Approval
  else if (model.update_status === 'pending_review' && model.pending_changes) {
    const changes = model.pending_changes
    const updates: any = {
      update_status: 'idle',
      pending_changes: null,
      updated_at: new Date().toISOString()
    }
    
    // Apply changes
    if (changes.title) updates.title = changes.title
    if (changes.description) updates.description = changes.description
    if (changes.tags) updates.tags = changes.tags
    if (changes.file_path) updates.file_path = changes.file_path
    if (changes.image_url) updates.image_url = changes.image_url
    if (changes.is_public !== undefined) updates.is_public = changes.is_public

    const { error } = await supabase
      .from('models')
      .update(updates)
      .eq('id', model.id)

    if (!error) {
      pendingModels.value = pendingModels.value.filter(m => m.id !== model.id)
    }
  }
}

const handleRejectModel = async (model: Model) => {
  // Case 1: New Model Rejection
  if (model.status === 'pending') {
    const { error } = await supabase
      .from('models')
      .update({ status: 'rejected' })
      .eq('id', model.id)

    if (!error) {
      pendingModels.value = pendingModels.value.filter(m => m.id !== model.id)
    }
  }
  // Case 2: Model Update Rejection
  else if (model.update_status === 'pending_review') {
    const { error } = await supabase
      .from('models')
      .update({ 
        update_status: 'idle',
        pending_changes: null
      })
      .eq('id', model.id)

    if (!error) {
      pendingModels.value = pendingModels.value.filter(m => m.id !== model.id)
    }
  }
}

const handleDeleteModel = async (id: string) => {
  if (!confirm('确定要永久删除此模型吗？此操作不可撤销。')) return

  const { error } = await supabase
    .from('models')
    .delete()
    .eq('id', id)

  if (error) {
    alert('删除失败: ' + error.message)
  } else {
    allModels.value = allModels.value.filter(m => m.id !== id)
  }
}

const handleApproveProfile = async (profile: Profile) => {
  if (!profile.pending_changes) return

  const changes = profile.pending_changes as any
  const updates: any = {
    profile_status: 'approved',
    pending_changes: null
  }
  
  if (changes.username) updates.username = changes.username
  if (changes.bio) updates.bio = changes.bio
  if (changes.avatar_url) updates.avatar_url = changes.avatar_url

  const { error } = await supabase
    .from('profiles')
    .update(updates)
    .eq('id', profile.id)

  if (!error) {
    pendingProfiles.value = pendingProfiles.value.filter(p => p.id !== profile.id)
  }
}

const handleRejectProfile = async (id: string) => {
  const { error } = await supabase
    .from('profiles')
    .update({ 
      profile_status: 'approved',
      pending_changes: null 
    })
    .eq('id', id)

  if (!error) {
    pendingProfiles.value = pendingProfiles.value.filter(p => p.id !== id)
  }
}

const handleUpdateUserRole = async (id: string, newRole: 'user' | 'admin') => {
  const { error } = await supabase
    .from('profiles')
    .update({ role: newRole })
    .eq('id', id)

  if (error) {
    alert('角色更新失败: ' + error.message)
  } else {
    const user = allProfiles.value.find(p => p.id === id)
    if (user) user.role = newRole
  }
}

const handleDeleteUser = async (id: string) => {
  if (!confirm('确定要删除此用户资料吗？(注意：这不会删除其 Auth 账号，仅删除资料表记录)')) return

  const { error } = await supabase
    .from('profiles')
    .delete()
    .eq('id', id)

  if (error) {
    alert('删除失败: ' + error.message)
  } else {
    allProfiles.value = allProfiles.value.filter(p => p.id !== id)
  }
}

const switchTab = (tab: typeof activeTab.value) => {
  activeTab.value = tab
  if (tab === 'audit_models') fetchPendingModels()
  else if (tab === 'audit_profiles') fetchPendingProfiles()
  else if (tab === 'manage_models') fetchAllModels()
  else if (tab === 'manage_users') fetchAllProfiles()
}

onMounted(() => {
  fetchPendingModels()
})
</script>

<template>
  <div class="container py-xl">
    <header class="admin-header">
      <div class="header-content">
        <h1>管理员控制台</h1>
        <p>审核用户申请、管理模型库及维护社区秩序</p>
      </div>
      <div class="header-stats" v-if="!loading">
        <div class="stat-item" v-if="pendingModels.length > 0">
          <span class="stat-value">{{ pendingModels.length }}</span>
          <span class="stat-label">待审模型</span>
        </div>
        <div class="stat-item" v-if="pendingProfiles.length > 0">
          <span class="stat-value">{{ pendingProfiles.length }}</span>
          <span class="stat-label">待审资料</span>
        </div>
      </div>
    </header>

    <nav class="tabs-nav" role="tablist" aria-label="管理选项卡">
      <div class="tabs-list">
        <button 
          v-for="tab in [
            { id: 'audit_models', label: '模型审核', icon: FolderTree, count: pendingModels.length },
            { id: 'audit_profiles', label: '资料审核', icon: User, count: pendingProfiles.length },
            { id: 'manage_models', label: '模型管理', icon: Box },
            { id: 'manage_users', label: '用户管理', icon: Users }
          ]"
          :key="tab.id"
          role="tab"
          :aria-selected="activeTab === tab.id"
          :aria-controls="`${tab.id}-panel`"
          :id="`tab-${tab.id}`"
          @click="switchTab(tab.id as any)" 
          :class="['tab-btn', activeTab === tab.id ? 'active' : '']"
        >
          <component :is="tab.icon" :size="18" aria-hidden="true" />
          <span>{{ tab.label }}</span>
          <span v-if="tab.count" class="tab-badge" aria-label="待处理数量">{{ tab.count }}</span>
        </button>
      </div>
    </nav>

    <main class="admin-content" aria-live="polite">
      <LoadingSpinner v-if="loading" message="正在获取数据..." :aria-busy="true" />

      <Transition name="fade" mode="out-in">
        <div :key="activeTab" v-if="!loading" :id="`${activeTab}-panel`" role="tabpanel" :aria-labelledby="`tab-${activeTab}`">
          <!-- Models Audit List -->
          <div v-if="activeTab === 'audit_models'">
            <div v-if="pendingModels.length === 0" class="empty-state" role="alert" aria-live="polite">
              <div class="empty-illustration" aria-hidden="true">
                <Check :size="48" />
              </div>
              <h3>暂无待审核模型</h3>
              <p>所有上传均已处理完毕，休息一下吧！</p>
            </div>

            <div v-else class="audit-grid">
              <article v-for="model in pendingModels" :key="model.id" class="audit-card">
                <div class="card-preview">
                   <img :src="model.image_url || 'https://via.placeholder.com/300x200?text=No+Preview'" :alt="`${model.title} 的预览图`" loading="lazy">
                   <div class="card-type-tag" :class="model.update_status === 'pending_review' ? 'type-update' : 'type-new'">
                     {{ model.update_status === 'pending_review' ? '更新申请' : '新模型' }}
                   </div>
                </div>
                
                <div class="card-body">
                  <header class="card-header">
                    <h3 class="card-title">
                      {{ model.title }}
                      <template v-if="model.update_status === 'pending_review' && model.pending_changes?.title">
                        <ArrowRight :size="16" class="title-arrow" aria-hidden="true" />
                        <span class="new-title">{{ model.pending_changes.title }}</span>
                      </template>
                    </h3>
                    <div class="card-meta">
                      <span class="author">
                        <User :size="14" aria-hidden="true" /> {{ model.profiles?.username || '未知用户' }}
                      </span>
                      <span class="dot" aria-hidden="true">•</span>
                      <span class="date">{{ new Date(model.created_at).toLocaleDateString() }}</span>
                    </div>
                  </header>

                  <div class="card-changes" v-if="model.update_status === 'pending_review' && model.pending_changes">
                     <div v-if="model.pending_changes.description" class="change-item">
                       <span class="change-label">描述更新</span>
                       <p class="change-content">{{ model.pending_changes.description }}</p>
                     </div>
                     <div v-if="model.pending_changes.file_path" class="change-item">
                       <span class="change-badge file-change">
                         <FileCode :size="12" aria-hidden="true" /> 文件已更改
                       </span>
                     </div>
                  </div>

                  <p class="card-desc" v-else>{{ model.description || '暂无描述' }}</p>

                  <div class="card-actions">
                    <button @click="handleApproveModel(model)" class="btn btn--primary btn--sm" :aria-label="`通过模型 ${model.title} 的审核`">
                      <Check :size="16" aria-hidden="true" /> 通过
                    </button>
                    <button @click="handleRejectModel(model)" class="btn btn--outline btn--sm danger" :aria-label="`拒绝模型 ${model.title} 的审核`">
                      <X :size="16" aria-hidden="true" /> 拒绝
                    </button>
                    <router-link :to="`/model/${model.id}`" class="btn btn--ghost btn--sm" :aria-label="`查看模型 ${model.title} 的详情`">
                      查看详情
                    </router-link>
                  </div>
                </div>
              </article>
            </div>
          </div>

          <!-- Profile Audit List -->
          <div v-else-if="activeTab === 'audit_profiles'">
            <div v-if="pendingProfiles.length === 0" class="empty-state" role="alert" aria-live="polite">
              <div class="empty-illustration" aria-hidden="true">
                <Check :size="48" />
              </div>
              <h3>资料审核已清空</h3>
              <p>目前没有需要审核的用户资料更新。</p>
            </div>

            <div v-else class="audit-grid">
              <article v-for="profile in pendingProfiles" :key="profile.id" class="audit-card profile-card">
                <div class="profile-comparison">
                  <div class="avatar-side old">
                    <img :src="profile.avatar_url || 'https://via.placeholder.com/80'" :alt="`${profile.username} 的当前头像`" />
                    <span>当前</span>
                  </div>
                  <ArrowRight :size="20" class="compare-arrow" aria-hidden="true" />
                  <div class="avatar-side new">
                    <img :src="profile.pending_changes?.avatar_url || profile.avatar_url || 'https://via.placeholder.com/80'" :alt="`${profile.username} 的申请头像`" />
                    <span>申请更新</span>
                  </div>
                </div>

                <div class="card-body">
                  <header class="card-header">
                    <h3 class="card-title">
                      {{ profile.username }}
                      <template v-if="profile.pending_changes?.username">
                        <ArrowRight :size="16" class="title-arrow" aria-hidden="true" />
                        <span class="new-title">{{ profile.pending_changes.username }}</span>
                      </template>
                    </h3>
                    <p class="user-id">ID: {{ profile.id }}</p>
                  </header>

                  <div class="card-changes">
                    <div v-if="profile.pending_changes?.bio" class="change-item">
                      <span class="change-label">简介更新</span>
                      <p class="change-content">{{ profile.pending_changes.bio }}</p>
                    </div>
                  </div>

                  <div class="card-actions">
                    <button @click="handleApproveProfile(profile)" class="btn btn--primary btn--sm" aria-label="通过资料审核">
                      <Check :size="16" aria-hidden="true" /> 通过
                    </button>
                    <button @click="handleRejectProfile(profile.id)" class="btn btn--outline btn--sm danger" aria-label="拒绝资料审核">
                      <X :size="16" aria-hidden="true" /> 拒绝
                    </button>
                  </div>
                </div>
              </article>
            </div>
          </div>

          <!-- Model Management -->
          <div v-else-if="activeTab === 'manage_models'">
            <div class="table-toolbar">
              <div class="search-input">
                <Search :size="18" aria-hidden="true" />
                <input v-model="modelSearch" type="text" placeholder="搜索模型名称、作者..." aria-label="搜索模型">
              </div>
            </div>

            <div class="table-container">
              <table class="admin-table">
                <thead>
                  <tr>
                    <th scope="col">预览</th>
                    <th scope="col">模型标题</th>
                    <th scope="col">作者</th>
                    <th scope="col">状态</th>
                    <th scope="col">数据</th>
                    <th scope="col">发布日期</th>
                    <th scope="col" class="text-right">操作</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="model in filteredModels" :key="model.id">
                    <td>
                      <img 
                        :src="model.image_url || 'https://via.placeholder.com/60x40'" 
                        :alt="`${model.title} 的缩略图`"
                        class="table-img"
                        loading="lazy"
                        @load="(e) => (e.target as HTMLImageElement).classList.add('loaded')"
                      >
                    </td>
                    <td>
                      <div class="table-main-info">
                        <span class="main-text">{{ model.title }}</span>
                        <span class="sub-text">ID: {{ model.id.slice(0, 8) }}</span>
                      </div>
                    </td>
                    <td>{{ model.profiles?.username || '未知' }}</td>
                    <td>
                      <span :class="['status-tag', model.status]">
                        {{ model.status === 'approved' ? '已发布' : model.status === 'pending' ? '审核中' : '已拒绝' }}
                      </span>
                    </td>
                    <td>
                      <div class="stats-cell">
                        <span>下载: {{ model.downloads }}</span>
                      </div>
                    </td>
                    <td>{{ new Date(model.created_at).toLocaleDateString() }}</td>
                    <td class="text-right">
                      <div class="table-btns">
                        <router-link :to="`/model/${model.id}/edit`" class="icon-btn" title="编辑" aria-label="编辑模型">
                          <Edit :size="18" aria-hidden="true" />
                        </router-link>
                        <button @click="handleDeleteModel(model.id)" class="icon-btn danger" title="删除" aria-label="删除模型">
                          <Trash2 :size="18" aria-hidden="true" />
                        </button>
                      </div>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>

          <!-- User Management -->
          <div v-else-if="activeTab === 'manage_users'">
            <div class="table-toolbar">
              <div class="search-input">
                <Search :size="18" aria-hidden="true" />
                <input v-model="userSearch" type="text" placeholder="搜索用户名、用户 ID..." aria-label="搜索用户">
              </div>
            </div>

            <div class="table-container">
              <table class="admin-table">
                <thead>
                  <tr>
                    <th scope="col">用户</th>
                    <th scope="col">角色</th>
                    <th scope="col">资料状态</th>
                    <th scope="col">注册日期</th>
                    <th scope="col" class="text-right">操作</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="profile in filteredProfiles" :key="profile.id">
                    <td>
                      <div class="user-cell">
                        <img :src="profile.avatar_url || 'https://via.placeholder.com/40'" :alt="`${profile.username} 的头像`" class="user-avatar">
                        <div class="user-info">
                          <span class="username">{{ profile.username }}</span>
                          <span class="user-id">{{ profile.id }}</span>
                        </div>
                      </div>
                    </td>
                    <td>
                      <span :class="['role-tag', profile.role]">
                        {{ profile.role === 'admin' ? '管理员' : '普通用户' }}
                      </span>
                    </td>
                    <td>
                      <span :class="['status-tag', profile.profile_status === 'approved' ? 'approved' : 'pending']">
                        {{ profile.profile_status === 'approved' ? '正常' : '待审核' }}
                      </span>
                    </td>
                    <td>{{ new Date(profile.created_at).toLocaleDateString() }}</td>
                    <td class="text-right">
                      <div class="table-btns">
                        <button 
                          @click="handleUpdateUserRole(profile.id, profile.role === 'admin' ? 'user' : 'admin')" 
                          class="icon-btn" 
                          :class="profile.role === 'admin' ? 'warning' : 'success'"
                          :title="profile.role === 'admin' ? '降级为普通用户' : '升级为管理员'"
                          :aria-label="profile.role === 'admin' ? '降级为普通用户' : '升级为管理员'"
                        >
                          <component :is="profile.role === 'admin' ? ShieldAlert : Shield" :size="18" aria-hidden="true" />
                        </button>
                        <button @click="handleDeleteUser(profile.id)" class="icon-btn danger" title="删除资料" aria-label="删除用户资料">
                          <Trash2 :size="18" aria-hidden="true" />
                        </button>
                      </div>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </Transition>
    </main>
  </div>
</template>

<style lang="scss" scoped>
@use '../styles/themes/variables' as *;

.admin-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  margin-bottom: $spacing-xl;
  padding-bottom: $spacing-lg;
  border-bottom: 1px solid var(--color-border);
  
  .header-content {
    h1 {
      font-size: 2.25rem;
      font-weight: 800;
      color: var(--color-text-main);
      letter-spacing: -0.02em;
      margin-bottom: $spacing-xs;
    }
    
    p {
      color: var(--color-text-muted);
      font-size: 1.1rem;
    }
  }

  .header-stats {
    display: flex;
    gap: $spacing-md;

    .stat-item {
      display: flex;
      flex-direction: column;
      align-items: center;
      padding: $spacing-sm $spacing-lg;
      background: var(--color-bg-light);
      border-radius: $radius-lg;
      border: 1px solid var(--color-border);
      min-width: 100px;

      .stat-value {
        font-size: 1.5rem;
        font-weight: 700;
        color: var(--color-primary);
      }

      .stat-label {
        font-size: 0.75rem;
        color: var(--color-text-muted);
        text-transform: uppercase;
        letter-spacing: 0.05em;
      }
    }
  }

  @media (max-width: 768px) {
    flex-direction: column;
    align-items: flex-start;
    gap: $spacing-md;

    .header-stats {
      width: 100%;
      .stat-item { flex: 1; }
    }
  }
}

.tabs-nav {
  margin-bottom: $spacing-xl;
  background: var(--color-bg-white);
  padding: $spacing-xs;
  border-radius: $radius-xl;
  border: 1px solid var(--color-border);
  box-shadow: $shadow-sm;

  .tabs-list {
    display: flex;
    gap: $spacing-xs;
    flex-wrap: wrap;
  }

  .tab-btn {
    display: flex;
    align-items: center;
    gap: $spacing-sm;
    padding: $spacing-sm $spacing-lg;
    background: transparent;
    border: none;
    border-radius: $radius-lg;
    cursor: pointer;
    font-weight: 600;
    color: var(--color-text-muted);
    transition: all $transition-base;
    position: relative;
    
    &:hover {
      color: var(--color-text-main);
      background: var(--color-bg-light);
    }
    
    &.active {
      color: white;
      background: var(--color-primary);
      box-shadow: 0 4px 12px rgba($color-primary, 0.2);

      .tab-badge {
        background: white;
        color: var(--color-primary);
      }
    }

    .tab-badge {
      display: flex;
      align-items: center;
      justify-content: center;
      min-width: 18px;
      height: 18px;
      padding: 0 5px;
      font-size: 0.7rem;
      font-weight: 700;
      background: var(--color-danger);
      color: white;
      border-radius: $radius-full;
      transition: all $transition-base;
    }
  }
}

.audit-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
  gap: $spacing-xl;

  @media (max-width: 480px) {
    grid-template-columns: 1fr;
  }
}

.audit-card {
  display: flex;
  background: var(--color-bg-white);
  border-radius: $radius-xl;
  border: 1px solid var(--color-border);
  overflow: hidden;
  transition: all $transition-base;
  box-shadow: $shadow-sm;

  &:hover {
    transform: translateY(-4px);
    box-shadow: $shadow-md;
    border-color: rgba($color-primary, 0.3);
  }

  &.profile-card {
    flex-direction: column;
    
    .profile-comparison {
      height: 140px;
      background: linear-gradient(135deg, rgba($color-primary, 0.05) 0%, rgba($color-secondary, 0.05) 100%);
      display: flex;
      align-items: center;
      justify-content: center;
      padding: $spacing-md;

      .avatar-side {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: $spacing-xs;

        img {
          width: 70px;
          height: 70px;
          border-radius: $radius-full;
          object-fit: cover;
          border: 3px solid white;
          box-shadow: $shadow-sm;
        }

        span {
          font-size: 0.7rem;
          font-weight: 600;
          color: var(--color-text-muted);
          text-transform: uppercase;
        }
      }

      .compare-arrow {
        color: var(--color-primary);
        opacity: 0.5;
        margin: 0 $spacing-md;
      }
    }
  }

  .card-preview {
    width: 160px;
    flex-shrink: 0;
    position: relative;

    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }

    .card-type-tag {
      position: absolute;
      top: $spacing-sm;
      left: $spacing-sm;
      padding: 2px 8px;
      border-radius: $radius-sm;
      font-size: 0.7rem;
      font-weight: 700;
      color: white;
      text-transform: uppercase;
      backdrop-filter: blur(4px);

      &.type-new { background: rgba(#10b981, 0.8); }
      &.type-update { background: rgba($color-primary, 0.8); }
    }
  }

  .card-body {
    flex: 1;
    padding: $spacing-lg;
    display: flex;
    flex-direction: column;
    gap: $spacing-md;

    .card-header {
      .card-title {
        font-size: 1.1rem;
        font-weight: 700;
        color: var(--color-text-main);
        margin-bottom: $spacing-xs;
        display: flex;
        align-items: center;
        gap: $spacing-xs;
        flex-wrap: wrap;

        .title-arrow { color: var(--color-primary); opacity: 0.5; }
        .new-title { color: var(--color-primary); }
      }

      .card-meta {
        display: flex;
        align-items: center;
        gap: $spacing-sm;
        font-size: 0.85rem;
        color: var(--color-text-muted);

        .author { display: flex; align-items: center; gap: 4px; }
        .dot { opacity: 0.5; }
      }

      .user-id {
        font-size: 0.8rem;
        color: var(--color-text-muted);
        font-family: monospace;
      }
    }

    .card-changes {
      background: var(--color-bg-light);
      padding: $spacing-sm $spacing-md;
      border-radius: $radius-lg;
      font-size: 0.9rem;

      .change-item {
        &:not(:last-child) { margin-bottom: $spacing-sm; }
        
        .change-label {
          display: block;
          font-size: 0.75rem;
          font-weight: 700;
          color: var(--color-text-muted);
          text-transform: uppercase;
          margin-bottom: 2px;
        }

        .change-content {
          color: var(--color-text-main);
          line-height: 1.4;
        }

        .change-badge {
          display: inline-flex;
          align-items: center;
          gap: 4px;
          padding: 2px 8px;
          background: rgba($color-primary, 0.1);
          color: var(--color-primary);
          border-radius: $radius-sm;
          font-size: 0.75rem;
          font-weight: 600;
        }
      }
    }

    .card-actions {
      margin-top: auto;
      display: flex;
      gap: $spacing-sm;

      .btn { flex: 1; justify-content: center; }
    }
  }
}

.table-toolbar {
  margin-bottom: $spacing-lg;
  display: flex;
  justify-content: flex-end;

  .search-input {
    position: relative;
    width: 320px;

    svg {
      position: absolute;
      left: $spacing-md;
      top: 50%;
      transform: translateY(-50%);
      color: var(--color-text-muted);
      pointer-events: none;
    }

    input {
      width: 100%;
      padding: $spacing-sm $spacing-md $spacing-sm 40px;
      background: var(--color-bg-white);
      border: 1px solid var(--color-border);
      border-radius: $radius-lg;
      font-size: 0.95rem;
      transition: all $transition-base;

      &:focus {
        border-color: var(--color-primary);
        box-shadow: 0 0 0 4px rgba($color-primary, 0.1);
        outline: none;
      }
    }
  }

  @media (max-width: 480px) {
    .search-input { width: 100%; }
  }
}

.table-container {
  background: var(--color-bg-white);
  border-radius: $radius-xl;
  border: 1px solid var(--color-border);
  box-shadow: $shadow-sm;
  overflow: hidden;
  overflow-x: auto;
}

.admin-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
  font-size: 0.95rem;

  th {
    background: var(--color-bg-light);
    padding: $spacing-md $spacing-lg;
    font-weight: 700;
    color: var(--color-text-muted);
    text-transform: uppercase;
    font-size: 0.75rem;
    letter-spacing: 0.05em;
    border-bottom: 1px solid var(--color-border);
  }

  td {
    padding: $spacing-md $spacing-lg;
    border-bottom: 1px solid var(--color-border);
    vertical-align: middle;
  }

  tr:last-child td { border-bottom: none; }

  .table-img {
    width: 60px;
    height: 40px;
    object-fit: cover;
    border-radius: $radius-sm;
    box-shadow: $shadow-sm;
    background-color: #f1f5f9;
    transition: $transition-base;
    opacity: 0;

    &.loaded {
      opacity: 1;
    }
  }

  .table-main-info {
    display: flex;
    flex-direction: column;
    gap: 2px;

    .main-text {
      font-weight: 600;
      color: var(--color-text-main);
    }

    .sub-text {
      font-size: 0.8rem;
      color: var(--color-text-muted);
      font-family: monospace;
    }
  }

  .user-cell {
    display: flex;
    align-items: center;
    gap: $spacing-md;

    .user-avatar {
      width: 40px;
      height: 40px;
      border-radius: $radius-full;
      object-fit: cover;
      box-shadow: $shadow-sm;
    }

    .user-info {
      display: flex;
      flex-direction: column;
      gap: 2px;

      .username { font-weight: 600; color: var(--color-text-main); }
      .user-id { font-size: 0.75rem; color: var(--color-text-muted); font-family: monospace; }
    }
  }

  .status-tag, .role-tag {
    display: inline-flex;
    padding: 2px 10px;
    border-radius: $radius-full;
    font-size: 0.8rem;
    font-weight: 600;

    &.approved, &.success { background: #ecfdf5; color: #059669; }
    &.pending, &.warning { background: #fffbeb; color: #d97706; }
    &.rejected, &.danger { background: #fef2f2; color: #dc2626; }
    &.admin { background: #eef2ff; color: #4f46e5; }
    &.user { background: #f3f4f6; color: #4b5563; }
  }

  .stats-cell {
    font-size: 0.85rem;
    color: var(--color-text-muted);
    font-weight: 500;
  }

  .table-btns {
    display: flex;
    gap: $spacing-xs;
    justify-content: flex-end;
  }
}

.icon-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: $radius-lg;
  border: 1px solid var(--color-border);
  background: var(--color-bg-white);
  color: var(--color-text-muted);
  cursor: pointer;
  transition: all $transition-base;
  
  &:hover {
    color: var(--color-primary);
    border-color: var(--color-primary);
    background: rgba($color-primary, 0.05);
    transform: translateY(-2px);
  }

  &.danger:hover {
    color: var(--color-danger);
    border-color: var(--color-danger);
    background: rgba($color-danger, 0.05);
  }

  &.success:hover {
    color: #059669;
    border-color: #059669;
    background: rgba(#059669, 0.05);
  }

  &.warning:hover {
    color: #d97706;
    border-color: #d97706;
    background: rgba(#d97706, 0.05);
  }
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: $spacing-2xl 0;
  color: var(--color-text-muted);
  text-align: center;

  .empty-illustration {
    width: 80px;
    height: 80px;
    background: var(--color-bg-light);
    border-radius: $radius-full;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: $spacing-lg;
    color: var(--color-primary);
    opacity: 0.5;
  }

  h3 { font-size: 1.25rem; font-weight: 600; color: var(--color-text-main); margin-bottom: $spacing-xs; }
  p { font-size: 1rem; }
}

.text-right { text-align: right; }

// Transition
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
  transform: translateY(10px);
}
</style>
