<script setup lang="ts">
import { ref } from 'vue'
import { supabase } from '../supabase/client'
import { useHead } from '@vueuse/head'

useHead({
  title: '重置密码 - YSM 模型站',
})

const email = ref('')
const loading = ref(false)
const message = ref('')
const errorMsg = ref('')
const isSuccess = ref(false)

const handleResetPassword = async () => {
  loading.value = true
  message.value = ''
  errorMsg.value = ''
  isSuccess.value = false

  try {
    // Supabase will send a password reset email to the user
    // The link in the email should redirect to /update-password
    const { error } = await supabase.auth.resetPasswordForEmail(email.value, {
      redirectTo: `${window.location.origin}/update-password`,
    })

    if (error) throw error

    isSuccess.value = true
    message.value = '重置密码邮件已发送，请检查您的邮箱。'
  } catch (e: any) {
    console.error('Reset password error:', e)
    if (e.message.includes('Rate limit exceeded')) {
      errorMsg.value = '请求过于频繁，请稍后再试'
    } else {
      errorMsg.value = e.message || '发送失败，请检查邮箱是否正确'
    }
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="container auth-container">
    <div class="card auth-card">
      <h1 class="auth-title">重置密码</h1>
      
      <div v-if="isSuccess" class="success-message">
        <p>{{ message }}</p>
        <p class="tip">请点击邮件中的链接来设置新密码。</p>
      </div>

      <form v-else @submit.prevent="handleResetPassword" class="auth-form">
        <p class="instruction">请输入您的注册邮箱，我们将向您发送重置密码的链接。</p>

        <div class="form-group">
          <label for="email">邮箱</label>
          <input 
            id="email" 
            v-model="email" 
            type="email" 
            class="input" 
            required 
            placeholder="请输入邮箱地址"
          >
        </div>

        <div v-if="errorMsg" class="error-message">
          {{ errorMsg }}
        </div>

        <button type="submit" class="btn btn--primary" :disabled="loading">
          {{ loading ? '发送中...' : '发送重置邮件' }}
        </button>
        
        <div class="auth-footer">
          <router-link to="/login" class="back-link">返回登录</router-link>
        </div>
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

.instruction {
  color: var(--color-text-muted);
  font-size: 0.875rem;
  margin-bottom: $spacing-sm;
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
}

.success-message {
  text-align: center;
  color: var(--color-success);
  background-color: #ecfdf5;
  padding: $spacing-lg;
  border-radius: $radius-md;
  border: 1px solid #a7f3d0;
  
  .tip {
    margin-top: $spacing-sm;
    font-size: 0.875rem;
    color: var(--color-text-muted);
  }
}

.auth-footer {
  margin-top: $spacing-sm;
  text-align: center;
  
  .back-link {
    color: var(--color-text-muted);
    font-size: 0.875rem;
    text-decoration: none;
    
    &:hover {
      color: var(--color-primary);
      text-decoration: underline;
    }
  }
}
</style>