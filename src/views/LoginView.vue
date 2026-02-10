<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { supabase } from '../supabase/client'
import { useUserStore } from '../stores/user'
import { useHead } from '@vueuse/head'
import { Mail, Lock, Loader2, ArrowRight, UserPlus, LogIn, ShieldAlert } from 'lucide-vue-next'

const router = useRouter()
const userStore = useUserStore()

useHead({
  title: '登录 / 注册 - YSM 模型站',
  meta: [
    { name: 'description', content: '登录 YSM 模型站，上传和管理你的 Minecraft 模型。' }
  ]
})

const isLogin = ref(true)
const email = ref('')
const password = ref('')
// const username = ref('')
const loading = ref(false)
const errorMsg = ref('')

const handleSubmit = async () => {
  loading.value = true
  errorMsg.value = ''

  try {
    if (isLogin.value) {
      const { error } = await supabase.auth.signInWithPassword({
        email: email.value,
        password: password.value
      })
      if (error) throw error
    } else {
      // Registration: Email & Password only
      const { error } = await supabase.auth.signUp({
        email: email.value,
        password: password.value
      })
      if (error) throw error
      
      // No manual profile creation here, rely on trigger to set random username
      if (!error) {
        alert('注册成功！请前往您的邮箱确认验证邮件，然后登录。')
        isLogin.value = true
        return
      }
    }
    
    await userStore.fetchUser()
    router.push('/')
  } catch (e: any) {
    console.error('Auth error:', e)
    // Translate common Supabase errors to Chinese
    if (e.message.includes('Invalid login credentials')) {
      errorMsg.value = '邮箱或密码错误'
    } else if (e.message.includes('Email not confirmed')) {
      errorMsg.value = '请先前往邮箱确认您的注册邮件'
    } else if (e.message.includes('User already registered')) {
      errorMsg.value = '该邮箱已被注册'
    } else if (e.message.includes('Password should be at least')) {
      errorMsg.value = '密码长度不能少于6位'
    } else if (e.message.includes('Error sending confirmation email')) {
      errorMsg.value = '发送验证邮件失败。请联系管理员检查邮件服务配置。'
    } else {
      errorMsg.value = e.message || '操作失败，请稍后重试'
    }
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="container auth-container">
    <div class="auth-card">
      <div class="auth-header">
        <div class="auth-logo">
          <LogIn v-if="isLogin" :size="24" />
          <UserPlus v-else :size="24" />
        </div>
        <h1 class="auth-title">{{ isLogin ? '欢迎回来' : '开启创意之旅' }}</h1>
        <p class="auth-subtitle">{{ isLogin ? '登录以管理您的模型和收藏' : '加入 YSM 社区，分享您的精彩作品' }}</p>
      </div>
      
      <form @submit.prevent="handleSubmit" class="auth-form">
        <div class="form-group">
          <label for="email">邮箱地址</label>
          <div class="input-wrapper">
            <Mail class="input-icon" :size="18" aria-hidden="true" />
            <input 
              id="email" 
              v-model="email" 
              type="email" 
              class="input" 
              required 
              placeholder="name@example.com"
              aria-required="true"
            >
          </div>
        </div>

        <div class="form-group">
          <div class="label-row">
            <label for="password">密码</label>
            <router-link v-if="isLogin" to="/reset-password" class="forgot-link">忘记密码？</router-link>
          </div>
          <div class="input-wrapper">
            <Lock class="input-icon" :size="18" aria-hidden="true" />
            <input 
              id="password" 
              v-model="password" 
              type="password" 
              class="input" 
              required 
              placeholder="至少 6 位字符"
              aria-required="true"
            >
          </div>
        </div>

        <div v-if="errorMsg" class="error-msg" role="alert" aria-live="polite">
          <ShieldAlert :size="18" aria-hidden="true" /> {{ errorMsg }}
        </div>

        <button type="submit" class="btn btn--primary auth-submit" :disabled="loading" :aria-busy="loading">
          <template v-if="loading">
            <Loader2 class="animate-spin" :size="20" aria-hidden="true" /> 请稍候...
          </template>
          <template v-else>
            {{ isLogin ? '立即登录' : '注册账号' }}
            <ArrowRight :size="18" aria-hidden="true" />
          </template>
        </button>
      </form>

      <div class="auth-footer">
        <p>
          {{ isLogin ? '还没有账号？' : '已经有账号了？' }}
          <button @click="isLogin = !isLogin" class="auth-switch-btn">
            {{ isLogin ? '立即注册' : '返回登录' }}
          </button>
        </p>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@use '../styles/themes/variables' as *;

.auth-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: calc(100vh - 4.5rem);
  padding: $spacing-xl $spacing-md;
  background: radial-gradient(circle at top right, rgba($color-primary, 0.05), transparent),
              radial-gradient(circle at bottom left, rgba($color-secondary, 0.05), transparent);
}

.auth-card {
  background: var(--color-bg-white);
  padding: $spacing-3xl;
  border-radius: $radius-2xl;
  border: 1px solid var(--color-border);
  box-shadow: $shadow-xl;
  width: 100%;
  max-width: 450px;
  animation: authFadeIn 0.6s cubic-bezier(0.16, 1, 0.3, 1);
  backdrop-filter: blur(10px);
}

@keyframes authFadeIn {
  from { opacity: 0; transform: translateY(30px) scale(0.95); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}

.auth-header {
  text-align: center;
  margin-bottom: $spacing-2xl;

  .auth-logo {
    width: 60px;
    height: 60px;
    background: linear-gradient(135deg, var(--color-primary) 0%, var(--color-secondary) 100%);
    border-radius: $radius-xl;
    margin: 0 auto $spacing-lg;
    display: flex;
    align-items: center;
    justify-content: center;
    color: white;
    box-shadow: 0 8px 16px rgba($color-primary, 0.2);
    transition: transform $transition-base;

    &:hover {
      transform: rotate(10deg) scale(1.1);
    }
  }

  .auth-title {
    font-size: 2rem;
    font-weight: 800;
    color: var(--color-text-main);
    margin-bottom: $spacing-xs;
    letter-spacing: -0.025em;
  }

  .auth-subtitle {
    color: var(--color-text-muted);
    font-size: 1rem;
    line-height: 1.5;
  }
}

.auth-form {
  display: flex;
  flex-direction: column;
  gap: $spacing-xl;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: $spacing-sm;

  label {
    font-size: 0.9rem;
    font-weight: 700;
    color: var(--color-text-main);
    margin-left: 2px;
  }

  .label-row {
    display: flex;
    justify-content: space-between;
    align-items: center;

    .forgot-link {
      font-size: 0.8rem;
      color: var(--color-primary);
      text-decoration: none;
      font-weight: 600;
      
      &:hover {
        text-decoration: underline;
      }
    }
  }
}

.input-wrapper {
  position: relative;
  display: flex;
  align-items: center;

  .input-icon {
    position: absolute;
    left: $spacing-md;
    color: var(--color-text-muted);
    pointer-events: none;
    transition: all $transition-base;
    z-index: 1;
  }

  .input {
    width: 100%;
    padding-left: 3.25rem;
    height: 3.5rem;
    border-radius: $radius-xl;
    background: var(--color-bg-light);
    border: 2px solid transparent;
    font-size: 1rem;
    transition: all $transition-base;
    
    &::placeholder {
      color: var(--color-text-muted);
      opacity: 0.6;
    }

    &:focus {
      background: white;
      border-color: var(--color-primary);
      box-shadow: 0 0 0 4px rgba($color-primary, 0.1);
      outline: none;

      & + .input-icon {
        color: var(--color-primary);
        transform: scale(1.1);
      }
    }
  }
}

.error-msg {
  background-color: #fef2f2;
  border: 1px solid #fee2e2;
  color: #dc2626;
  padding: $spacing-md;
  border-radius: $radius-lg;
  font-size: 0.9rem;
  display: flex;
  align-items: center;
  gap: $spacing-sm;
  font-weight: 500;
  animation: shake 0.4s cubic-bezier(.36,.07,.19,.97) both;
}

@keyframes shake {
  10%, 90% { transform: translate3d(-1px, 0, 0); }
  20%, 80% { transform: translate3d(2px, 0, 0); }
  30%, 50%, 70% { transform: translate3d(-4px, 0, 0); }
  40%, 60% { transform: translate3d(4px, 0, 0); }
}

.auth-submit {
  height: 3.75rem;
  font-size: 1.1rem;
  font-weight: 700;
  margin-top: $spacing-sm;
  gap: $spacing-md;
  border-radius: $radius-xl;
  box-shadow: 0 4px 12px rgba($color-primary, 0.2);

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 20px rgba($color-primary, 0.3);
  }

  &:active {
    transform: translateY(0);
  }

  .animate-spin {
    animation: spin 1s linear infinite;
  }
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.auth-footer {
  margin-top: $spacing-2xl;
  text-align: center;
  font-size: 0.95rem;
  color: var(--color-text-muted);

  .auth-switch-btn {
    background: none;
    border: none;
    color: var(--color-primary);
    font-weight: 800;
    cursor: pointer;
    padding: $spacing-xs $spacing-sm;
    border-radius: $radius-md;
    transition: all $transition-base;
    
    &:hover {
      background: rgba($color-primary, 0.05);
      color: var(--color-primary-hover);
    }
  }
}

@media (max-width: 480px) {
  .auth-card {
    padding: $spacing-xl;
    box-shadow: none;
    border: none;
    background: transparent;
  }
}
</style>
