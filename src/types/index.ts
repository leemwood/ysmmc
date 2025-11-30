export interface Profile {
  id: string
  username: string
  avatar_url: string | null
  bio: string | null
  role: 'user' | 'admin'
  profile_status: 'approved' | 'pending_review'
  pending_changes: {
    username?: string
    bio?: string
    avatar_url?: string
  } | null
  created_at: string
}

export interface Model {
  id: string
  user_id: string
  title: string
  description: string | null
  file_path: string
  image_url: string | null
  tags: string[]
  is_public: boolean
  status: 'pending' | 'approved' | 'rejected'
  update_status?: 'idle' | 'pending_review'
  pending_changes?: {
    title?: string
    description?: string
    tags?: string[]
    file_path?: string
    image_url?: string
    is_public?: boolean
  } | null
  downloads: number
  created_at: string
  profiles?: Profile
}

export interface Favorite {
  id: string
  user_id: string
  model_id: string
  created_at: string
  models?: Model
}
