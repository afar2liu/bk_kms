// API 响应基础类型
export interface ApiResponse<T = any> {
  code: number
  msg: string
  data?: T
}

// 分页数据
export interface PageData<T> {
  rows: T[]
  total: number
}

// 标签
export interface Tag {
  id: number
  name: string
  count?: number
}

// 书签
export interface Bookmark {
  id: number
  url: string
  title: string
  excerpt: string
  author: string
  is_archive: boolean
  created_at: number
  updated_at: number
  tags: Tag[]
}

// 书签内容
export interface BookmarkContent {
  id: number
  url: string
  title: string
  html: string
  created_at: number
  update_at: number
}

// 书签列表请求参数
export interface BookmarkListParams {
  keyword?: string
  tags?: string
  page: number
  page_size?: number
}

// 创建书签请求
export interface CreateBookmarkRequest {
  url: string
  title: string
  excerpt: string
  tags: { name: string }[]
  create_archive: boolean
}

// 更新书签请求
export interface UpdateBookmarkRequest {
  id: number
  url: string
  title: string
  excerpt: string
  author: string
  tags: { name: string }[]
  create_archive: boolean
}

// 登录请求
export interface LoginRequest {
  username: string
  pwd: string
  captcha: string
  captcha_id: string
}

// 登录响应
export interface LoginResponse {
  id: number
  username: string
  token: string
}

// 导入进度事件
export interface ImportProgressEvent {
  type: 'progress' | 'success' | 'error' | 'complete'
  message: string
  current: number
  total: number
  url?: string
}

// 用户信息
export interface UserInfo {
  id: number
  username: string
  owner: boolean
}
