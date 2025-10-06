import request from '@/utils/request'
import type {
  ApiResponse,
  PageData,
  Bookmark,
  BookmarkListParams,
  CreateBookmarkRequest,
  UpdateBookmarkRequest,
  BookmarkContent
} from '@/types'

/**
 * 获取书签列表
 */
export function getBookmarkList(params: BookmarkListParams) {
  return request.get<ApiResponse<PageData<Bookmark>>>('/api/v1/bookmarks', { params })
}

/**
 * 创建书签
 */
export function createBookmark(data: CreateBookmarkRequest) {
  return request.post<ApiResponse>('/api/v1/bookmark', data)
}

/**
 * 更新书签
 */
export function updateBookmark(data: UpdateBookmarkRequest) {
  return request.put<ApiResponse>('/api/v1/bookmarks', data)
}

/**
 * 删除书签
 */
export function deleteBookmarks(ids: number[]) {
  return request.delete<ApiResponse>('/api/v1/bookmark', { data: ids })
}

/**
 * 获取书签内容
 */
export function getBookmarkContent(id: number) {
  return request.get<ApiResponse<BookmarkContent>>(`/api/v1/bookmark/${id}/content`)
}

/**
 * 导入书签（SSE 流式响应）
 * @param file HTML 文件
 * @param generateTag 是否自动生成分类标签
 * @param createArchive 是否创建归档
 * @param onProgress 进度回调
 */
export async function importBookmarks(
  file: File,
  generateTag: boolean,
  createArchive: boolean,
  onProgress: (event: any) => void
): Promise<void> {
  const formData = new FormData()
  formData.append('bookmark_file', file)
  formData.append('generate_tag', generateTag.toString())
  formData.append('create_archive', createArchive.toString())

  const url = `/api/v1/bookmarks/import`
  const token = localStorage.getItem('bk_kms_token')

  const response = await fetch(url, {
    method: 'POST',
    headers: {
      Authorization: `Bearer ${token}`
    },
    body: formData
  })

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`)
  }

  const reader = response.body?.getReader()
  if (!reader) {
    throw new Error('无法读取响应流')
  }

  const decoder = new TextDecoder()

  while (true) {
    const { value, done } = await reader.read()
    if (done) break

    const chunk = decoder.decode(value)
    const lines = chunk.split('\n')

    for (const line of lines) {
      if (line.startsWith('data: ')) {
        const jsonStr = line.substring(6)
        try {
          const data = JSON.parse(jsonStr)
          onProgress(data)
        } catch (e) {
          console.error('JSON parse error:', e)
        }
      }
    }
  }
}
