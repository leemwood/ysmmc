import axios, { type AxiosInstance, type AxiosResponse } from 'axios'
import type { ApiResponse, LoginResponse, User, Model, ModelVersion, ModelImage, PaginatedResponse, Announcement, Favorite } from '@/types'

const isDev = import.meta.env.DEV

const api: AxiosInstance = axios.create({
  baseURL: isDev ? '/api' : 'https://api.ysmmc.cn/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('access_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export const authApi = {
  register: (data: { email: string; password: string; username: string }) =>
    api.post<ApiResponse<{ id: string; email: string; username: string }>>('/auth/register', data),

  login: (data: { email: string; password: string }) =>
    api.post<ApiResponse<LoginResponse>>('/auth/login', data),

  refresh: (refreshToken: string) =>
    api.post<ApiResponse<{ access_token: string; refresh_token: string; expires_in: number }>>('/auth/refresh', {
      refresh_token: refreshToken,
    }),

  forgotPassword: (email: string) =>
    api.post<ApiResponse<null>>('/auth/forgot-password', { email }),

  resetPassword: (token: string, password: string) =>
    api.post<ApiResponse<null>>('/auth/reset-password', { token, password }),

  verifyEmail: (token: string) =>
    api.get<ApiResponse<null>>('/auth/verify', { params: { token } }),

  changeEmail: (newEmail: string) =>
    api.post<ApiResponse<null>>('/auth/change-email', { new_email: newEmail }),

  verifyEmailChange: (token: string) =>
    api.get<ApiResponse<null>>('/auth/verify-email-change', { params: { token } }),

  me: () => api.get<ApiResponse<User>>('/auth/me'),
}

export const userApi = {
  getMe: () => api.get<ApiResponse<User>>('/users/me'),

  updateMe: (data: { username?: string; bio?: string; avatar_url?: string; avatar_id?: string }) =>
    api.put<ApiResponse<User>>('/users/me', data),

  changePassword: (data: { old_password: string; new_password: string }) =>
    api.put<ApiResponse<null>>('/users/me/password', data),

  getById: (id: string) => api.get<ApiResponse<User>>(`/users/${id}`),

  getUserModels: (id: string, page = 1, pageSize = 10) =>
    api.get<ApiResponse<PaginatedResponse<Model>>>(`/users/${id}/models`, {
      params: { page, page_size: pageSize },
    }),
}

export const modelApi = {
  list: (page = 1, pageSize = 12, search = '') =>
    api.get<ApiResponse<PaginatedResponse<Model>>>('/models', {
      params: { page, page_size: pageSize, search },
    }),

  getById: (id: string) =>
    api.get<ApiResponse<{ model: Model; is_favorited: boolean; favorite_count: number }>>(`/models/${id}`),

  create: (data: {
    title: string
    description?: string
    file_path: string
    file_size: number
    image_url?: string
    image_id?: string
    tags: string[]
    is_public: boolean
  }) => api.post<ApiResponse<Model>>('/models', data),

  update: (id: string, data: Partial<{
    title: string
    description: string
    file_path: string
    file_size: number
    image_url: string
    image_id: string
    tags: string[]
    is_public: boolean
  }>) => api.put<ApiResponse<Model>>(`/models/${id}`, data),

  delete: (id: string) => api.delete<ApiResponse<null>>(`/models/${id}`),

  download: (id: string) => api.post<ApiResponse<{ file_path: string; file_name: string }>>(`/models/${id}/download`),

  addFavorite: (id: string) => api.post<ApiResponse<null>>(`/models/${id}/favorite`),

  removeFavorite: (id: string) => api.delete<ApiResponse<null>>(`/models/${id}/favorite`),

  checkFavorite: (id: string) => api.get<ApiResponse<{ is_favorited: boolean }>>(`/models/${id}/favorite`),
}

export const modelVersionApi = {
  list: (modelId: string) =>
    api.get<ApiResponse<ModelVersion[]>>(`/models/${modelId}/versions`),

  get: (modelId: string, versionId: string) =>
    api.get<ApiResponse<ModelVersion>>(`/models/${modelId}/versions/${versionId}`),

  create: (modelId: string, data: {
    version_number: string
    description?: string
    file_path: string
    file_size: number
    image_id?: string
    image_url?: string
    changelog?: string
  }) => api.post<ApiResponse<ModelVersion>>(`/models/${modelId}/versions`, data),

  update: (modelId: string, versionId: string, data: {
    description?: string
    image_id?: string
    image_url?: string
    changelog?: string
  }) => api.put<ApiResponse<ModelVersion>>(`/models/${modelId}/versions/${versionId}`, data),

  setCurrent: (modelId: string, versionId: string) =>
    api.put<ApiResponse<null>>(`/models/${modelId}/versions/${versionId}/current`),

  delete: (modelId: string, versionId: string) =>
    api.delete<ApiResponse<null>>(`/models/${modelId}/versions/${versionId}`),

  download: (modelId: string, versionId: string) =>
    api.post<ApiResponse<{ download_url: string; file_name: string }>>(`/models/${modelId}/versions/${versionId}/download`),
}

export const modelImageApi = {
  list: (modelId: string) =>
    api.get<ApiResponse<ModelImage[]>>(`/models/${modelId}/images`),

  add: (modelId: string, fileId: string) =>
    api.post<ApiResponse<ModelImage>>(`/models/${modelId}/images`, { file_id: fileId }),

  delete: (modelId: string, fileId: string) =>
    api.delete<ApiResponse<null>>(`/models/${modelId}/images/${fileId}`),

  updateOrder: (modelId: string, images: { file_id: string; sort_order: number }[]) =>
    api.put<ApiResponse<null>>(`/models/${modelId}/images/order`, { images }),
}

export const favoriteApi = {
  list: (page = 1, pageSize = 12) =>
    api.get<ApiResponse<PaginatedResponse<Favorite>>>('/favorites', {
      params: { page, page_size: pageSize },
    }),
}

export const announcementApi = {
  list: () => api.get<ApiResponse<Announcement[]>>('/announcements'),

  listAll: (page = 1, pageSize = 20) =>
    api.get<ApiResponse<PaginatedResponse<Announcement>>>('/announcements/all', {
      params: { page, page_size: pageSize },
    }),

  getById: (id: string) => api.get<ApiResponse<Announcement>>(`/announcements/${id}`),
}

export const adminApi = {
  getStats: () =>
    api.get<
      ApiResponse<{
        total_users: number
        total_models: number
        pending_models: number
        total_downloads: number
      }>
    >('/admin/stats'),

  getSuperAdmin: () =>
    api.get<ApiResponse<User>>('/admin/super-admin'),

  listAllModels: (page = 1, pageSize = 20, status = '', search = '') =>
    api.get<ApiResponse<PaginatedResponse<Model>>>('/admin/models', {
      params: { page, page_size: pageSize, status, search },
    }),

  listPendingModels: (page = 1, pageSize = 20) =>
    api.get<ApiResponse<PaginatedResponse<Model>>>('/admin/models/pending', {
      params: { page, page_size: pageSize },
    }),

  listPendingUpdates: (page = 1, pageSize = 20) =>
    api.get<ApiResponse<PaginatedResponse<Model>>>('/admin/models/pending-updates', {
      params: { page, page_size: pageSize },
    }),

  approveModel: (id: string) => api.put<ApiResponse<null>>(`/admin/models/${id}/approve`),

  rejectModel: (id: string, reason: string) =>
    api.put<ApiResponse<null>>(`/admin/models/${id}/reject`, { reason }),

  deleteModel: (id: string) =>
    api.delete<ApiResponse<null>>(`/admin/models/${id}`),

  listUsers: (page = 1, pageSize = 20) =>
    api.get<ApiResponse<PaginatedResponse<User>>>('/admin/users', {
      params: { page, page_size: pageSize },
    }),

  updateUserRole: (id: string, role: string) =>
    api.put<ApiResponse<null>>(`/admin/users/${id}/role`, { role }),

  setAdmin: (id: string) =>
    api.put<ApiResponse<null>>(`/admin/users/${id}/admin`),

  removeAdmin: (id: string) =>
    api.delete<ApiResponse<null>>(`/admin/users/${id}/admin`),

  banUser: (id: string, reason: string) =>
    api.put<ApiResponse<null>>(`/admin/users/${id}/ban`, { reason }),

  unbanUser: (id: string) =>
    api.put<ApiResponse<null>>(`/admin/users/${id}/unban`),

  listPendingProfiles: (page = 1, pageSize = 20) =>
    api.get<ApiResponse<PaginatedResponse<User>>>('/admin/profiles/pending', {
      params: { page, page_size: pageSize },
    }),

  approveProfile: (id: string) => api.put<ApiResponse<null>>(`/admin/profiles/${id}/approve`),

  rejectProfile: (id: string) => api.put<ApiResponse<null>>(`/admin/profiles/${id}/reject`),

  createAnnouncement: (data: { title: string; content: string }) =>
    api.post<ApiResponse<Announcement>>('/admin/announcements', data),

  updateAnnouncement: (id: string, data: { title?: string; content?: string; is_active?: boolean }) =>
    api.put<ApiResponse<Announcement>>(`/admin/announcements/${id}`, data),

  deleteAnnouncement: (id: string) => api.delete<ApiResponse<null>>(`/admin/announcements/${id}`),
}

export const uploadApi = {
  uploadModel: async (file: File): Promise<AxiosResponse<ApiResponse<{ file_path: string; file_name: string; file_size: number }>>> => {
    const formData = new FormData()
    formData.append('file', file)
    return api.post('/upload/model', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
  },

  uploadImage: async (file: File, category = 'model_image'): Promise<AxiosResponse<ApiResponse<{ id: string; file_id: string; file_name: string; size: number }>>> => {
    const formData = new FormData()
    formData.append('file', file)
    formData.append('category', category)
    return api.post('/upload/image', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
  },
}

export const fileApi = {
  getUrl: (fileId: string) => {
    const isDev = import.meta.env.DEV
    const baseUrl = isDev ? '' : 'https://api.ysmmc.cn'
    return `${baseUrl}/api/files/${fileId}`
  },

  get: (id: string) => api.get(`/files/${id}`, { responseType: 'blob' }),

  upload: async (file: File, category = 'general'): Promise<AxiosResponse<ApiResponse<{ id: string; name: string; mime_type: string; size: number; category: string }>>> => {
    const formData = new FormData()
    formData.append('file', file)
    formData.append('category', category)
    return api.post('/files', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
  },

  delete: (id: string) => api.delete<ApiResponse<null>>(`/files/${id}`),
}

export default api
