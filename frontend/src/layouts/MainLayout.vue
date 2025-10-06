<template>
  <a-layout class="main-layout">
    <a-layout-header class="header">
      <div class="header-left">
        <h1 class="logo">📚 书签知识管理系统</h1>
      </div>
      <div class="header-right">
        <a-space size="large">
          <span class="username">{{ userStore.userInfo?.username }}</span>
          <a-button type="text" @click="handleLogout">
            <template #icon><icon-export /></template>
            退出登录
          </a-button>
        </a-space>
      </div>
    </a-layout-header>

    <a-layout>
      <a-layout-sider
        :width="220"
        :style="{ background: 'var(--color-bg-2)' }"
      >
        <a-menu
          :selected-keys="selectedKeys"
          :style="{ width: '100%', height: '100%' }"
          @menu-item-click="handleMenuClick"
        >
          <a-menu-item key="bookmarks">
            <template #icon><icon-book /></template>
            书签列表
          </a-menu-item>
          <a-menu-item key="tags">
            <template #icon><icon-tags /></template>
            标签管理
          </a-menu-item>
        </a-menu>
      </a-layout-sider>

      <a-layout-content class="content">
        <router-view />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Modal } from '@arco-design/web-vue'
import {
  IconExport,
  IconBook,
  IconTags
} from '@arco-design/web-vue/es/icon'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const selectedKeys = computed(() => {
  const name = route.name as string
  if (name === 'BookmarkList') return ['bookmarks']
  if (name === 'TagManage') return ['tags']
  return []
})

function handleMenuClick(key: string) {
  if (key === 'bookmarks') {
    router.push({ name: 'BookmarkList' })
  } else if (key === 'tags') {
    router.push({ name: 'TagManage' })
  }
}

function handleLogout() {
  Modal.confirm({
    title: '确认退出',
    content: '确定要退出登录吗？',
    onOk: () => {
      userStore.logout()
      router.push({ name: 'Login' })
    }
  })
}
</script>

<style scoped>
.main-layout {
  height: 100vh;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 24px;
  background: var(--color-bg-1);
  border-bottom: 1px solid var(--color-border-2);
}

.header-left {
  display: flex;
  align-items: center;
}

.logo {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: var(--color-text-1);
}

.header-right {
  display: flex;
  align-items: center;
}

.username {
  color: var(--color-text-2);
  font-size: 14px;
}

.content {
  background: var(--color-bg-1);
  overflow-y: auto;
}
</style>
