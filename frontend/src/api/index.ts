import api from './client';
import type { ApiResponse, PaginatedResponse, Model, Announcement } from '../types';

export const modelApi = {
  list: (page = 1, pageSize = 12, search = '') =>
    api.get<ApiResponse<PaginatedResponse<Model>>>('/models', {
      params: { page, page_size: pageSize, search },
    }),
  getById: (id: string) =>
    api.get<ApiResponse<{ model: Model; is_favorited: boolean; favorite_count: number }>>(`/models/${id}`),
};

export const announcementApi = {
  list: () => api.get<ApiResponse<Announcement[]>>('/announcements'),
};
