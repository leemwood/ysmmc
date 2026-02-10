<script setup lang="ts">
import { RouterLink, useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { Upload, User as UserIcon, LogOut, Menu, X, Shield } from 'lucide-vue-next'
import { ref } from 'vue'

const userStore = useUserStore()
const router = useRouter()
const isMenuOpen = ref(false)

const handleSignOut = async () => {
  await userStore.signOut()
  router.push('/login')
}
</script>

<template>
  <nav class="navbar" aria-label="主导航">
    <div class="container navbar__container">
      <RouterLink to="/" class="navbar__logo" aria-label="YSM 管理器 首页">
        YSM 管理器
      </RouterLink>

      <div class="navbar__desktop">
        <template v-if="userStore.user">
          <RouterLink to="/upload" class="btn btn--primary">
            <Upload :size="18" aria-hidden="true" />
            上传模型
          </RouterLink>
          <div class="navbar__user">
            <RouterLink v-if="userStore.isAdmin" to="/admin" class="navbar__link" title="管理员控制台" aria-label="管理员控制台">
              <Shield :size="20" aria-hidden="true" />
            </RouterLink>
            <RouterLink to="/profile" class="navbar__link" aria-label="个人中心">
              <UserIcon :size="20" aria-hidden="true" />
              个人中心
            </RouterLink>
            <button @click="handleSignOut" class="btn btn--secondary btn--sm" aria-label="退出登录">
              <LogOut :size="18" aria-hidden="true" />
              退出登录
            </button>
          </div>
        </template>
        <template v-else>
          <RouterLink to="/login" class="btn btn--primary">
            登录 / 注册
          </RouterLink>
        </template>
      </div>

      <button 
        class="navbar__toggle" 
        @click="isMenuOpen = !isMenuOpen"
        :aria-expanded="isMenuOpen"
        aria-controls="mobile-menu"
        :aria-label="isMenuOpen ? '关闭菜单' : '打开菜单'"
      >
        <Menu v-if="!isMenuOpen" aria-hidden="true" />
        <X v-else aria-hidden="true" />
      </button>
    </div>

    <div v-if="isMenuOpen" id="mobile-menu" class="navbar__mobile" role="menu">
      <template v-if="userStore.user">
        <RouterLink to="/upload" class="navbar__mobile-link" @click="isMenuOpen = false" role="menuitem">
          <Upload :size="18" aria-hidden="true" />
          上传模型
        </RouterLink>
        <RouterLink v-if="userStore.isAdmin" to="/admin" class="navbar__mobile-link" @click="isMenuOpen = false" role="menuitem">
          <Shield :size="18" aria-hidden="true" />
          管理员控制台
        </RouterLink>
        <RouterLink to="/profile" class="navbar__mobile-link" @click="isMenuOpen = false" role="menuitem">
          <UserIcon :size="18" aria-hidden="true" />
          个人中心
        </RouterLink>
        <button @click="handleSignOut; isMenuOpen = false" class="navbar__mobile-link" role="menuitem">
          <LogOut :size="18" aria-hidden="true" />
          退出登录
        </button>
      </template>
      <template v-else>
        <RouterLink to="/login" class="navbar__mobile-link" @click="isMenuOpen = false" role="menuitem">
          登录 / 注册
        </RouterLink>
      </template>
    </div>
  </nav>
</template>

<style lang="scss" scoped>
@use '../styles/themes/variables' as *;

.navbar {
  background-color: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(12px);
  border-bottom: 1px solid var(--color-border);
  position: sticky;
  top: 0;
  z-index: 100;
  transition: $transition-base;

  &__container {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 4.5rem;
  }

  &__logo {
    font-size: 1.5rem;
    font-weight: 800;
    color: var(--color-primary);
    text-decoration: none;
    display: flex;
    align-items: center;
    gap: $spacing-sm;
    letter-spacing: -0.025em;

    &::before {
      content: '';
      display: block;
      width: 32px;
      height: 32px;
      background: linear-gradient(135deg, var(--color-primary) 0%, var(--color-primary-hover) 100%);
      border-radius: $radius-md;
    }
  }

  &__desktop {
    display: none;
    align-items: center;
    gap: $spacing-xl;

    @media (min-width: 768px) {
      display: flex;
    }
  }

  &__user {
    display: flex;
    align-items: center;
    gap: $spacing-md;
  }

  &__link {
    display: flex;
    align-items: center;
    gap: $spacing-xs;
    color: var(--color-text-main);
    text-decoration: none;
    font-weight: 500;
    transition: $transition-base;
    padding: $spacing-xs $spacing-sm;
    border-radius: $radius-md;

    &:hover {
      color: var(--color-primary);
      background-color: rgba($color-primary, 0.05);
    }
  }

  &__toggle {
    background: none;
    border: none;
    cursor: pointer;
    color: var(--color-text-main);
    
    @media (min-width: 768px) {
      display: none;
    }
  }

  &__mobile {
    border-top: 1px solid var(--color-border);
    padding: $spacing-md;
    display: flex;
    flex-direction: column;
    gap: $spacing-sm;
    background-color: white;

    &-link {
      display: flex;
      align-items: center;
      gap: $spacing-sm;
      padding: $spacing-sm;
      color: var(--color-text-main);
      border-radius: $radius-md;

      &:hover {
        background-color: var(--color-background);
      }
    }
  }
}

.btn--sm {
  padding: $spacing-xs $spacing-sm;
  font-size: 0.875rem;
}
</style>
