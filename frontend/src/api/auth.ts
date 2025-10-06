import request from '@/utils/request'
import type { ApiResponse, LoginRequest, LoginResponse } from '@/types'

/**
 * 获取验证码
 */
export function getCaptcha() {
  return request.get<ApiResponse<{ captcha_id: string; captcha: string }>>(
    '/api/v1/captcha'
  )
}

/**
 * 用户登录
 */
export function login(data: LoginRequest) {
  return request.post<ApiResponse<LoginResponse>>('/api/v1/auth/login', data)
}
