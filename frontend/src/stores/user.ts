import { defineStore } from 'pinia'
import { ref } from 'vue'
import { storage } from '@/utils/storage'
import type { UserInfo } from '@/types'

export const useUserStore = defineStore('user', () => {
  const token = ref<string>(storage.getToken() || '')
  const userInfo = ref<UserInfo | null>(storage.getUserInfo())

  // 设置 token
  function setToken(newToken: string) {
    token.value = newToken
    storage.setToken(newToken)
  }

  // 设置用户信息
  function setUserInfo(info: UserInfo) {
    userInfo.value = info
    storage.setUserInfo(info)
  }

  // 登出
  function logout() {
    token.value = ''
    userInfo.value = null
    storage.clear()
  }

  // 是否已登录
  function isLoggedIn() {
    return !!token.value
  }

  return {
    token,
    userInfo,
    setToken,
    setUserInfo,
    logout,
    isLoggedIn
  }
})
