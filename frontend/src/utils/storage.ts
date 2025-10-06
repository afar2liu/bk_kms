/**
 * localStorage 封装
 */

const TOKEN_KEY = 'bk_kms_token'
const USER_INFO_KEY = 'bk_kms_user_info'

export const storage = {
  // Token 相关
  getToken(): string | null {
    const token = localStorage.getItem(TOKEN_KEY)
    if (!token || token === 'undefined' || token === 'null') {
      return null
    }
    return token
  },

  setToken(token: string): void {
    localStorage.setItem(TOKEN_KEY, token)
  },

  removeToken(): void {
    localStorage.removeItem(TOKEN_KEY)
  },

  // 用户信息相关
  getUserInfo(): any {
    const userInfo = localStorage.getItem(USER_INFO_KEY)
    if (!userInfo || userInfo === 'undefined' || userInfo === 'null') {
      return null
    }
    try {
      return JSON.parse(userInfo)
    } catch (e) {
      console.error('Failed to parse user info:', e)
      return null
    }
  },

  setUserInfo(userInfo: any): void {
    localStorage.setItem(USER_INFO_KEY, JSON.stringify(userInfo))
  },

  removeUserInfo(): void {
    localStorage.removeItem(USER_INFO_KEY)
  },

  // 清除所有
  clear(): void {
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(USER_INFO_KEY)
  }
}
