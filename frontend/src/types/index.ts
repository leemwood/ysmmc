export interface ApiResponse<T> {
  success: boolean;
  message: string;
  data: T;
}

export interface PaginatedResponse<T> {
  items: T[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface User {
  id: string;
  email: string;
  username: string;
  avatar_url?: string;
  bio?: string;
  role: string;
  created_at: string;
}

export interface Model {
  id: string;
  title: string;
  description?: string;
  file_path: string;
  file_size: number;
  image_url?: string;
  tags: string[];
  downloads: number;
  is_public: boolean;
  status: string;
  user_id: string;
  user?: User;
  created_at: string;
}

export interface Announcement {
  id: string;
  title: string;
  content: string;
  is_active: boolean;
  created_at: string;
}
