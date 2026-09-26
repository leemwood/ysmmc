<script setup lang="ts">
import AppNavbar from '@/components/AppNavbar.vue'
import AppFooter from '@/components/AppFooter.vue'
import PWAInstallPrompt from '@/components/PWAInstallPrompt.vue'
import { ToastContainer } from '@/components/ui/toast'
</script>

<template>
  <div class="flex min-h-screen flex-col">
    <AppNavbar />
    <main class="min-h-[calc(100vh-4rem)] flex-1">
      <RouterView v-slot="{ Component }">
        <Transition name="page" mode="out-in">
          <component :is="Component" />
        </Transition>
      </RouterView>
    </main>
    <AppFooter />
    <PWAInstallPrompt />
    <ToastContainer />
  </div>
</template>
