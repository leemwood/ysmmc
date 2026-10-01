<script setup lang="ts">
import { computed } from 'vue'
import { Avatar, AvatarImage } from '@/components/ui/avatar'
import { useNexusmcStore } from '@/stores/nexusmc'

type AvatarSize = 'xs' | 'sm' | 'md' | 'lg'

const props = withDefaults(defineProps<{ size?: AvatarSize }>(), { size: 'md' })

// 各使用点原有的尺寸组合，保持视觉一致
const SIZES: Record<AvatarSize, { avatar: string; text: string }> = {
  xs: { avatar: 'h-8 w-8', text: 'text-xs' },
  sm: { avatar: 'h-8 w-8', text: 'text-sm' },
  md: { avatar: 'h-9 w-9', text: 'text-sm' },
  lg: { avatar: 'h-12 w-12', text: 'text-lg' },
}

const store = useNexusmcStore()
const user = computed(() => store.user)
const initial = computed(() => (user.value?.username || 'N').slice(0, 1).toUpperCase())
</script>

<template>
  <Avatar :class="SIZES[props.size].avatar">
    <AvatarImage
      v-if="user?.avatar"
      :src="user.avatar"
      :alt="user?.username || 'NexusMC 用户'"
    />
    <span
      v-else
      class="flex h-full w-full items-center justify-center rounded-full bg-muted font-medium"
      :class="SIZES[props.size].text"
    >
      {{ initial }}
    </span>
  </Avatar>
</template>
