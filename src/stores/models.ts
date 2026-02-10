import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Model } from '../types'
import { supabase } from '../supabase/client'

export const useModelStore = defineStore('models', () => {
  const models = ref<Model[]>([])
  const totalCount = ref(0)
  const loading = ref(false)
  const lastFetched = ref<number | null>(null)
  const singleModelCache = ref<Record<string, { data: Model; timestamp: number }>>({})
  const CACHE_DURATION = 5 * 60 * 1000 // 5 minutes cache

  const fetchModels = async (page: number, pageSize: number, searchQuery: string = '', force: boolean = false) => {
    // Basic caching logic: only cache the first page with no search query
    const isFirstPageNoSearch = page === 1 && !searchQuery
    
    if (!force && isFirstPageNoSearch && models.value.length > 0 && lastFetched.value && (Date.now() - lastFetched.value < CACHE_DURATION)) {
      return
    }

    loading.value = true
    
    let query = supabase
      .from('models')
      .select('*, profiles(username)', { count: 'exact' })
      .eq('is_public', true)
      .eq('status', 'approved')
      .order('created_at', { ascending: false })

    if (searchQuery) {
      query = query.ilike('title', `%${searchQuery}%`)
    }

    const from = (page - 1) * pageSize
    const to = from + pageSize - 1
    
    query = query.range(from, to)

    const { data, count, error } = await query

    if (error) {
      console.error('Error fetching models:', error)
    } else {
      if (isFirstPageNoSearch) {
        models.value = data as any
        totalCount.value = count || 0
        lastFetched.value = Date.now()
      }
      return { data, count }
    }
    loading.value = false
  }

  const fetchModelById = async (id: string, force: boolean = false) => {
    if (!force && singleModelCache.value[id] && (Date.now() - singleModelCache.value[id].timestamp < CACHE_DURATION)) {
      return singleModelCache.value[id].data
    }

    loading.value = true
    const { data, error } = await supabase
      .from('models')
      .select('*, profiles(username, avatar_url)')
      .eq('id', id)
      .single()

    loading.value = false

    if (error) {
      console.error('Error fetching model by id:', error)
      return null
    }

    if (data) {
      singleModelCache.value[id] = {
        data: data as any,
        timestamp: Date.now()
      }
    }
    return data as any
  }

  const clearCache = () => {
    models.value = []
    totalCount.value = 0
    lastFetched.value = null
    singleModelCache.value = {}
  }

  return {
    models,
    totalCount,
    loading,
    fetchModels,
    fetchModelById,
    clearCache
  }
})
