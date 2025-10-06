<template>
  <div class="login-container">
    <div class="login-box">
      <div class="login-header">
        <h1>📚 书签知识管理系统</h1>
        <p>Bookmark Knowledge Management System</p>
      </div>

      <a-form
        :model="formData"
        :rules="rules"
        @submit="handleLogin"
        layout="vertical"
        class="login-form"
      >
        <a-form-item field="username" label="用户名" hide-label>
          <a-input
            v-model="formData.username"
            placeholder="请输入用户名"
            size="large"
            allow-clear
          >
            <template #prefix>
              <icon-user />
            </template>
          </a-input>
        </a-form-item>

        <a-form-item field="pwd" label="密码" hide-label>
          <a-input-password
            v-model="formData.pwd"
            placeholder="请输入密码"
            size="large"
            allow-clear
          >
            <template #prefix>
              <icon-lock />
            </template>
          </a-input-password>
        </a-form-item>

        <a-form-item field="captcha" label="验证码" hide-label>
          <div class="captcha-wrapper">
            <a-input
              v-model="formData.captcha"
              placeholder="请输入验证码"
              size="large"
              allow-clear
              style="flex: 1"
            >
              <template #prefix>
                <icon-safe />
              </template>
            </a-input>
            <div class="captcha-image" @click="refreshCaptcha">
              <img v-if="captchaImage" :src="captchaImage" alt="验证码" />
              <a-spin v-else />
            </div>
          </div>
        </a-form-item>

        <a-form-item>
          <a-button
            type="primary"
            html-type="submit"
            size="large"
            long
            :loading="loading"
          >
            登录
          </a-button>
        </a-form-item>
      </a-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import { IconUser, IconLock, IconSafe } from '@arco-design/web-vue/es/icon'
import { getCaptcha, login } from '@/api/auth'
import { useUserStore } from '@/stores/user'
import type { LoginRequest } from '@/types'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const formData = ref<LoginRequest>({
  username: '',
  pwd: '',
  captcha: '',
  captcha_id: ''
})

const captchaImage = ref('')
const loading = ref(false)

const rules = {
  username: [{ required: true, message: '请输入用户名' }],
  pwd: [{ required: true, message: '请输入密码' }],
  captcha: [{ required: true, message: '请输入验证码' }]
}

// 获取验证码
async function refreshCaptcha() {
  try {
    const { data } = await getCaptcha()
    if (data.data) {
      formData.value.captcha_id = data.data.captcha_id
      // 后端返回的是 captcha 字段，已经是 base64 data URL 格式
      captchaImage.value = data.data.captcha
    }
  } catch (error) {
    Message.error('获取验证码失败')
  }
}

// 登录
async function handleLogin() {
  if (!formData.value.username || !formData.value.pwd || !formData.value.captcha) {
    return
  }

  loading.value = true
  try {
    const { data } = await login(formData.value)
    if (data.data) {
      // 保存 token 和用户信息
      userStore.setToken(data.data.token)
      userStore.setUserInfo({
        id: data.data.id,
        username: data.data.username,
        owner: true // 默认为 true，后续可以从后端获取
      })

      Message.success('登录成功')

      // 跳转到目标页面或首页
      const redirect = (route.query.redirect as string) || '/bookmarks'
      router.push(redirect)
    }
  } catch (error: any) {
    Message.error(error.message || '登录失败')
    // 刷新验证码
    refreshCaptcha()
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  refreshCaptcha()
})
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.login-box {
  background: white;
  border-radius: 16px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  padding: 48px;
  width: 100%;
  max-width: 420px;
}

.login-header {
  text-align: center;
  margin-bottom: 40px;
}

.login-header h1 {
  font-size: 28px;
  color: #333;
  margin-bottom: 8px;
}

.login-header p {
  font-size: 14px;
  color: #999;
}

.login-form {
  margin-top: 32px;
}

.captcha-wrapper {
  display: flex;
  gap: 12px;
  align-items: center;
}

.captcha-image {
  width: 120px;
  height: 40px;
  border: 1px solid var(--color-border-2);
  border-radius: 4px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  transition: all 0.3s;
}

.captcha-image:hover {
  border-color: rgb(var(--primary-6));
}

.captcha-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
</style>
