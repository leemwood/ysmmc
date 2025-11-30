<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { supabase } from '../supabase/client'
import { useUserStore } from '../stores/user'
import { Upload, X, Save, ArrowLeft } from 'lucide-vue-next'
import { useHead } from '@vueuse/head'
import type { Model } from '../types'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

useHead({
  title: '编辑模型 - YSM 模型站',
  meta: [
    { name: 'robots', content: 'noindex, nofollow' }
  ]
})

const model = ref<Model | null>(null)
const loading = ref(true)
const saving = ref(false)
const errorMsg = ref('')
const successMsg = ref('')

// Form Fields
const title = ref('')
const description = ref('')
const tags = ref<string[]>([])
const tagInput = ref('')
const isPublic = ref(true)

// File Uploads
const modelFile = ref<File | null>(null)
const imageFile = ref<File | null>(null)
const imagePreview = ref<string | null>(null)

const currentImage = computed(() => {
  if (imagePreview.value) return imagePreview.value
  if (model.value?.image_url) return model.value.image_url
  return null
})

const fetchModel = async () => {
  loading.value = true
  const { data, error } = await supabase
    .from('models')
    .select('*')
    .eq('id', route.params.id)
    .single()

  if (error || !data) {
    console.error(error)
    router.push('/')
    return
  }

  // Verify ownership
  if (data.user_id !== userStore.user?.id) {
    router.push('/')
    return
  }

  model.value = data as Model
  
  // Populate fields
  // If there are pending changes, maybe warn user?
  // For now, we populate with APPROVED data to base changes on.
  title.value = data.title
  description.value = data.description || ''
  tags.value = data.tags || []
  isPublic.value = data.is_public
  
  loading.value = false
}

const handleImageSelect = (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files && target.files[0]) {
    const file = target.files[0]
    imageFile.value = file
    imagePreview.value = URL.createObjectURL(file)
  }
}

const handleModelSelect = (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files && target.files[0]) {
    modelFile.value = target.files[0]
  }
}

const addTag = () => {
  if (tagInput.value.trim() && !tags.value.includes(tagInput.value.trim())) {
    tags.value.push(tagInput.value.trim())
    tagInput.value = ''
  }
}

const removeTag = (tag: string) => {
  tags.value = tags.value.filter(t => t !== tag)
}

const uploadFile = async (file: File, bucket: string) => {
  const fileExt = file.name.split('.').pop()
  const fileName = `${Math.random()}.${fileExt}`
  const filePath = `${userStore.user?.id}/${fileName}`

  const { error } = await supabase.storage.from(bucket).upload(filePath, file)
  if (error) throw error
  return filePath
}

const getPublicUrl = (bucket: string, path: string) => {
    const { data } = supabase.storage.from(bucket).getPublicUrl(path)
    return data.publicUrl
}

const handleSubmit = async () => {
  if (!userStore.user || !model.value) return

  saving.value = true
  errorMsg.value = ''
  successMsg.value = ''

  try {
    // 1. Prepare Pending Changes Object
    const changes: any = {}
    let hasChanges = false

    if (title.value !== model.value.title) {
      changes.title = title.value
      hasChanges = true
    }
    if (description.value !== (model.value.description || '')) {
      changes.description = description.value
      hasChanges = true
    }
    // Compare arrays simply
    if (JSON.stringify(tags.value) !== JSON.stringify(model.value.tags)) {
      changes.tags = tags.value
      hasChanges = true
    }
    if (isPublic.value !== model.value.is_public) {
      changes.is_public = isPublic.value
      hasChanges = true
    }

    // 2. Handle File Uploads
    if (modelFile.value) {
      const modelPath = await uploadFile(modelFile.value, 'models')
      changes.file_path = modelPath
      hasChanges = true
    }

    if (imageFile.value) {
      const imagePath = await uploadFile(imageFile.value, 'images')
      changes.image_url = getPublicUrl('images', imagePath)
      hasChanges = true
    }

    if (!hasChanges) {
      errorMsg.value = '未检测到任何更改'
      return
    }

    // 3. Update Database
    const { error } = await supabase
      .from('models')
      .update({
        pending_changes: changes,
        update_status: 'pending_review'
      })
      .eq('id', model.value.id)

    if (error) throw error

    successMsg.value = '更新请求已提交，请等待管理员审核。'
    // Optionally redirect after delay
    setTimeout(() => {
      router.push(`/model/${model.value?.id}`)
    }, 2000)

  } catch (e: any) {
    console.error(e)
    errorMsg.value = e.message || '更新提交失败'
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  fetchModel()
})
</script>

<template>
  <div class="container">
    <div class="edit-header">
      <button class="btn btn--ghost" @click="router.back()">
        <ArrowLeft :size="20" /> 返回
      </button>
      <h1>编辑模型</h1>
    </div>

    <div v-if="loading" class="loading">
      正在加载...
    </div>

    <div v-else-if="model" class="card edit-form">
      <form @submit.prevent="handleSubmit">
        <div v-if="model.update_status === 'pending_review'" class="alert alert-warning">
          此模型已有正在审核中的更新请求。提交新请求将覆盖之前的申请。
        </div>

        <div class="form-section">
          <h2>模型文件与预览</h2>
          
          <div class="form-group">
            <label>更新模型文件 (可选)</label>
            <div class="file-input-wrapper">
              <input 
                type="file" 
                accept=".ysm,.zip,.rar,.7z" 
                @change="handleModelSelect" 
                class="file-input"
              >
              <div class="file-placeholder" :class="{ 'has-file': modelFile }">
                <Upload :size="24" />
                <span>{{ modelFile ? modelFile.name : '点击上传新版本 (保持不变则留空)' }}</span>
              </div>
            </div>
          </div>

          <div class="form-group">
            <label>更新预览图片 (可选)</label>
            <div class="file-input-wrapper">
              <input 
                type="file" 
                accept="image/*" 
                @change="handleImageSelect" 
                class="file-input"
              >
              <div class="file-placeholder" :class="{ 'has-file': imageFile }">
                <img v-if="currentImage" :src="currentImage" class="image-preview" />
                <div v-else class="placeholder-content">
                  <Upload :size="24" />
                  <span>点击上传新图片</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="form-section">
          <h2>基本信息</h2>
          
          <div class="form-group">
            <label for="title">标题</label>
            <input 
              id="title" 
              v-model="title" 
              type="text" 
              class="input" 
              required 
            >
          </div>

          <div class="form-group">
            <label for="description">描述 (支持 Markdown)</label>
            <textarea 
              id="description" 
              v-model="description" 
              class="input textarea" 
              rows="8"
            ></textarea>
          </div>

          <div class="form-group">
            <label>标签</label>
            <div class="tags-input">
              <div v-for="tag in tags" :key="tag" class="tag">
                {{ tag }}
                <button type="button" @click="removeTag(tag)"><X :size="14" /></button>
              </div>
              <input 
                v-model="tagInput" 
                @keydown.enter.prevent="addTag" 
                type="text" 
                class="input-ghost" 
                placeholder="添加标签 (按回车)"
              >
            </div>
          </div>

          <div class="form-group checkbox-group">
            <label class="checkbox-label">
              <input type="checkbox" v-model="isPublic">
              <span>公开此模型</span>
            </label>
          </div>
        </div>

        <div v-if="errorMsg" class="error-message">
          {{ errorMsg }}
        </div>
        
        <div v-if="successMsg" class="success-message">
          {{ successMsg }}
        </div>

        <div class="form-actions">
          <button type="button" class="btn btn--secondary" @click="router.back()">取消</button>
          <button type="submit" class="btn btn--primary" :disabled="saving">
            {{ saving ? '提交中...' : '提交更新申请' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@use '../styles/themes/variables' as *;

.edit-header {
  display: flex;
  align-items: center;
  gap: $spacing-md;
  margin-bottom: $spacing-xl;
  
  h1 {
    font-size: 2rem;
    font-weight: 700;
    color: var(--color-text-main);
  }
}

.btn--ghost {
  background: none;
  border: none;
  color: var(--color-text-muted);
  padding: 0;
  
  &:hover {
    color: var(--color-text-main);
  }
}

.edit-form {
  padding: $spacing-xl;
}

.form-section {
  margin-bottom: $spacing-xl;
  padding-bottom: $spacing-xl;
  border-bottom: 1px solid var(--color-border);

  &:last-child {
    border-bottom: none;
    padding-bottom: 0;
    margin-bottom: 0;
  }

  h2 {
    font-size: 1.25rem;
    font-weight: 600;
    margin-bottom: $spacing-md;
    color: var(--color-text-main);
  }
}

.form-group {
  margin-bottom: $spacing-md;

  label {
    display: block;
    font-weight: 500;
    margin-bottom: $spacing-sm;
    color: var(--color-text-main);
  }
}

.file-input-wrapper {
  position: relative;
  height: 120px;
  border: 2px dashed var(--color-border);
  border-radius: $radius-md;
  overflow: hidden;
  transition: $transition-base;

  &:hover {
    border-color: var(--color-primary);
    background-color: #f9fafb;
  }
}

.file-input {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  opacity: 0;
  cursor: pointer;
  z-index: 10;
}

.file-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: var(--color-text-muted);
  gap: $spacing-sm;

  &.has-file {
    color: var(--color-primary);
    border-color: var(--color-primary);
  }
}

.image-preview {
  height: 100%;
  width: 100%;
  object-fit: cover;
}

.placeholder-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: $spacing-sm;
}

.textarea {
  resize: vertical;
}

.tags-input {
  display: flex;
  flex-wrap: wrap;
  gap: $spacing-sm;
  padding: $spacing-sm;
  border: 1px solid var(--color-border);
  border-radius: $radius-md;
  background-color: white;

  &:focus-within {
    border-color: var(--color-primary);
    box-shadow: 0 0 0 2px rgba($color-primary, 0.2);
  }
}

.tag {
  display: flex;
  align-items: center;
  gap: $spacing-xs;
  background-color: #e0e7ff;
  color: var(--color-primary);
  padding: 2px 8px;
  border-radius: $radius-full;
  font-size: 0.875rem;

  button {
    background: none;
    border: none;
    cursor: pointer;
    color: inherit;
    display: flex;
  }
}

.input-ghost {
  border: none;
  outline: none;
  flex: 1;
  min-width: 120px;
  font-size: 0.875rem;
}

.checkbox-group {
  margin-top: $spacing-md;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: $spacing-sm;
  cursor: pointer;
  user-select: none;

  input {
    width: 1.25rem;
    height: 1.25rem;
  }
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: $spacing-md;
  margin-top: $spacing-xl;
}

.error-message {
  color: var(--color-danger);
  margin-top: $spacing-md;
  text-align: center;
}

.success-message {
  color: var(--color-secondary);
  margin-top: $spacing-md;
  text-align: center;
}

.alert {
  padding: $spacing-md;
  border-radius: $radius-md;
  margin-bottom: $spacing-lg;
}

.alert-warning {
  background-color: #fff7ed;
  color: #c2410c;
  border: 1px solid #ffedd5;
}

.loading {
  text-align: center;
  padding: $spacing-2xl;
  color: var(--color-text-muted);
}
</style>
