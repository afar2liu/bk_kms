import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios'
import { Message } from '@arco-design/web-vue'
import { useUserStore } from '@/stores/user'

// 创建 axios 实例
const request: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
request.interceptors.request.use(
  (config) => {
    const userStore = useUserStore()
    const token = userStore.token

    // 添加 token
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }

    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  (response: AxiosResponse) => {
    const { data } = response

    // 如果是文件下载或 SSE 流，直接返回
    if (
      response.config.responseType === 'blob' ||
      response.headers['content-type']?.includes('text/event-stream')
    ) {
      return response
    }

    // 处理业务错误
    if (data.code !== undefined && data.code !== 0) {
      Message.error(data.msg || '请求失败')
      return Promise.reject(new Error(data.msg || '请求失败'))
    }

    return response
  },
  (error) => {
    // 处理 HTTP 错误
    if (error.response) {
      const { status, data } = error.response

      switch (status) {
        case 401:
          Message.error('未授权，请重新登录')
          // 清除 token 并跳转到登录页
          const userStore = useUserStore()
          userStore.logout()
          window.location.href = '/login'
          break
        case 403:
          Message.error('拒绝访问')
          break
        case 404:
          Message.error('请求的资源不存在')
          break
        case 500:
          Message.error('服务器错误')
          break
        default:
          Message.error(data?.msg || '请求失败')
      }
    } else if (error.request) {
      Message.error('网络错误，请检查网络连接')
    } else {
      Message.error(error.message || '请求失败')
    }

    return Promise.reject(error)
  }
)

export default request
