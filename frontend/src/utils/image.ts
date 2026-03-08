import { fileApi } from '@/lib/api'

export function getImageUrl(imageId: string | null | undefined, fallbackUrl?: string | null | undefined): string {
  if (imageId) {
    return fileApi.getUrl(imageId)
  }
  if (fallbackUrl) {
    if (fallbackUrl.startsWith('http')) {
      return fallbackUrl
    }
    const isDev = import.meta.env.DEV
    const baseUrl = isDev ? '' : 'https://api.ysmmc.cn'
    return `${baseUrl}${fallbackUrl}`
  }
  return '/placeholder.png'
}

export function getAvatarUrl(avatarId: string | null | undefined, avatarUrl?: string | null | undefined): string {
  return getImageUrl(avatarId, avatarUrl)
}

export function getModelImageUrl(imageId: string | null | undefined, imageUrl?: string | null | undefined): string {
  return getImageUrl(imageId, imageUrl)
}
