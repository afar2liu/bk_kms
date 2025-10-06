import request from '@/utils/request'
import type { ApiResponse, Tag } from '@/types'

/**
 * 获取标签列表
 */
export function getTagList(name?: string) {
  return request.get<ApiResponse<Tag[]>>('/api/v1/tags', {
    params: name ? { name } : undefined
  })
}

/**
 * 更新标签（重命名）
 */
export function updateTag(id: number, name: string) {
  return request.put<ApiResponse>(`/api/v1/tag/${id}`, { name })
}
