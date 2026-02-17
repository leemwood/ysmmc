<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { supabase } from '../supabase/client'
import { useHead } from '@vueuse/head'

useHead({
  title: '设置新密码 - YSM 模型站',
})

const router = useRouter()
const password = ref('')
const confirmPassword = ref('')
const loading = ref(false)
const errorMsg = ref('')
const successMsg = ref('')

onMounted(async () => {
  // Check if we have a session (user clicked the email link)
  const { data: { session } } = await supabase.auth.getSession()
  if (!session) {
    // If no session, redirect to home or login
    // But wait, for update password flow, supabase usually logs the user in via the link
    // If the link is invalid or expired, we might want to show an error.
  }
})

const handleUpdatePassword = async () => {
  if (password.value !== confirmPassword.value) {
    errorMsg.value = '两次输入的密码不一致'
    return
  }

  if (password.value.length < 6) {
    errorMsg.value = '密码长度不能少于6位'
    return
  }

  loading.value = true
  errorMsg.value = ''
  successMsg.value = ''

  try {
    const { error } = await supabase.auth.updateUser({
      password: password.value
    })

    if (error) throw error

    successMsg.value = '密码修改成功！正在跳转...'
    setTimeout(() => {
      router.push('/')
    }, 2000)
  } catch (e: any) {
    console.error('Update password error:', e)
    errorMsg.value = e.message || '密码修改失败，请重试'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="container auth-container">
    <div class="card auth-card">
      <h1 class="auth-title">设置新密码</h1>
      
      <form @submit.prevent="handleUpdatePassword" class="auth-form">
        <div class="form-group">
          <label for="password">新密码</label>
          <input 
            id="password" 
            v-model="password" 
            type="password" 
            class="input" 
            required 
            minlength="6"
            placeholder="请输入新密码"
            aria-required="true"
          >
        </div>

        <div class="form-group">
          <label for="confirmPassword">确认新密码</label>
          <input 
            id="confirmPassword" 
            v-model="confirmPassword" 
            type="password" 
            class="input" 
            required 
            minlength="6"
            placeholder="请再次输入新密码"
            aria-required="true"
          >
        </div>

        <div v-if="errorMsg" class="error-message" role="alert" aria-live="polite">
          <ShieldAlert :size="18" aria-hidden="true" /> {{ errorMsg }}
        </div>

        <div v-if="successMsg" class="success-message" role="status" aria-live="polite">
          <CheckCircle :size="18" aria-hidden="true" /> {{ successMsg }}
        </div>

        <button type="submit" class="btn btn--primary" :disabled="loading" :aria-busy="loading">
          <template v-if="loading">
            <Loader2 class="animate-spin" :size="18" aria-hidden="true" /> 提交中...
          </template>
          <template v-else>
            确认修改
          </template>
        </button>
      </form>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@use '../styles/themes/variables' as *;

.auth-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 60vh;
}

.auth-card {
  width: 100%;
  max-width: 400px;
  padding: $spacing-xl;
}

.auth-title {
  text-align: center;
  margin-bottom: $spacing-xl;
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--color-text-main);
}

.auth-form {
  display: flex;
  flex-direction: column;
  gap: $spacing-lg;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: $spacing-sm;

  label {
    font-weight: 500;
    font-size: 0.875rem;
    color: var(--color-text-main);
  }
}

.error-message {
  color: var(--color-danger);
  font-size: 0.875rem;
  text-align: center;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: $spacing-xs;
}

.success-message {
  color: var(--color-success);
  font-size: 0.875rem;
  text-align: center;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: $spacing-xs;
}

.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
</style>