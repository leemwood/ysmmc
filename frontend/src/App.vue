<script setup lang="ts">
import AppNavbar from '@/components/AppNavbar.vue'
import AppFooter from '@/components/AppFooter.vue'
import PWAInstallPrompt from '@/components/PWAInstallPrompt.vue'
import { ToastContainer } from '@/components/ui/toast'
</script>

<template>
  <div class="flex min-h-dvh flex-col">
    <a href="#main" class="site-button site-button-primary sr-only focus:not-sr-only focus:fixed focus:top-2 focus:left-2 focus:z-50">跳转到内容</a>
    <AppNavbar />
    <main id="main" tabindex="-1" class="flex-1 focus:outline-none">
      <RouterView />
    </main>
    <AppFooter />
    <PWAInstallPrompt />
    <ToastContainer />
  </div>
</template>
