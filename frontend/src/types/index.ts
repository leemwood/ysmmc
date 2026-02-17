export interface User {
  id: string
  email: string
  username: string
  avatar_url: string | null
  bio: string | null
  role: 'user' | 'admin'
  profile_status: 'approved' | 'pending_review'
  pending_changes: PendingChanges | null
  created_at: string
}

export interface PendingChanges {
  username?: string
  bio?: string
  avatar_url?: string
}

export interface Model {
  id: string
  user_id: string
  title: string
  description: string | null
  file_path: string
  file_size: number
  image_url: string | null
  tags: string[]
  is_public: boolean
  status: 'pending' | 'approved' | 'rejected'
  update_status: 'idle' | 'pending_review'
  pending_changes: ModelPendingChanges | null
  downloads: number
  rejection_reason: string | null
  created_at: string
  updated_at: string
  user?: User
}

export interface ModelPendingChanges {
  title?: string
  description?: string
  tags?: string[]
  file_path?: string
  image_url?: string
  is_public?: boolean
}

export interface Favorite {
  id: string
  user_id: string
  model_id: string
  created_at: string
  model?: Model
}

export interface Announcement {
  id: string
  title: string
  content: string
  is_active: boolean
  created_at: string
}

export interface PaginatedResponse<T> {
  items: T[]
  total: number
  page: number
  page_size: number
  total_pages: number
}

export interface ApiResponse<T> {
  code: number
  message: string
  data?: T
}

export interface LoginResponse {
  access_token: string
  refresh_token: string
  expires_in: number
  user: User
}
