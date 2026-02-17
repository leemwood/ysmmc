import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { supabase } from '../supabase/client'
import type { User } from '@supabase/supabase-js'
import type { Profile } from '../types'

export const useUserStore = defineStore('user', () => {
  const user = ref<User | null>(null)
  const profile = ref<Profile | null>(null)
  const loading = ref(true)

  const isAdmin = computed(() => profile.value?.role === 'admin')

  const fetchUser = async () => {
    loading.value = true
    const { data: { user: currentUser } } = await supabase.auth.getUser()
    user.value = currentUser
    
    if (currentUser) {
        // Fix: Use maybeSingle() instead of single() to avoid 406 error when no rows returned
        // Also explicitly select fields to ensure format matches expected JSON
        const { data, error } = await supabase
            .from('profiles')
            .select('*')
            .eq('id', currentUser.id)
            .maybeSingle()
            
        if (error) {
            console.error('Error fetching profile:', error)
        }
        profile.value = data
    } else {
        profile.value = null
    }
    
    loading.value = false
  }

  const signOut = async () => {
    await supabase.auth.signOut()
    user.value = null
    profile.value = null
  }

  return {
    user,
    profile,
    isAdmin,
    loading,
    fetchUser,
    signOut
  }
})
