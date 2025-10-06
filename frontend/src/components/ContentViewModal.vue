<template>
  <a-modal
    v-model:visible="modalVisible"
    title="书签内容"
    width="900px"
    :footer="false"
  >
    <a-spin :loading="loading" style="width: 100%">
      <div v-if="content" class="content-container">
        <div class="content-header">
          <h2>{{ content.title }}</h2>
          <a :href="content.url" target="_blank" class="content-url">
            {{ content.url }}
          </a>
          <div class="content-meta">
            <span>创建时间: {{ formatDate(content.created_at) }}</span>
            <span>更新时间: {{ formatDate(content.update_at) }}</span>
          </div>
        </div>
        <a-divider />
        <div class="content-body" v-html="content.html"></div>
      </div>
      <a-empty v-else description="暂无内容" />
    </a-spin>
  </a-modal>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { Message } from '@arco-design/web-vue'
import { getBookmarkContent } from '@/api/bookmark'
import type { BookmarkContent } from '@/types'

interface Props {
  visible: boolean
  bookmarkId: number
}

interface Emits {
  (e: 'update:visible', value: boolean): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const modalVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

const loading = ref(false)
const content = ref<BookmarkContent | null>(null)

// 获取内容
async function fetchContent() {
  if (!props.bookmarkId) return

  loading.value = true
  try {
    const { data } = await getBookmarkContent(props.bookmarkId)
    if (data.data) {
      content.value = data.data
    }
  } catch (error) {
    Message.error('获取内容失败')
  } finally {
    loading.value = false
  }
}

// 格式化日期
function formatDate(timestamp: number) {
  return new Date(timestamp * 1000).toLocaleString('zh-CN')
}

// 监听 visible 变化
watch(
  () => props.visible,
  (visible) => {
    if (visible) {
      fetchContent()
    } else {
      content.value = null
    }
  }
)
</script>

<style scoped>
.content-container {
  max-height: 70vh;
  overflow-y: auto;
}

.content-header h2 {
  margin: 0 0 12px 0;
  font-size: 24px;
  color: var(--color-text-1);
}

.content-url {
  color: rgb(var(--primary-6));
  text-decoration: none;
  word-break: break-all;
}

.content-url:hover {
  text-decoration: underline;
}

.content-meta {
  margin-top: 12px;
  font-size: 12px;
  color: var(--color-text-3);
  display: flex;
  gap: 20px;
}

.content-body {
  padding: 20px 0;
  line-height: 1.8;
}

.content-body :deep(img) {
  max-width: 100%;
  height: auto;
}

.content-body :deep(pre) {
  background: var(--color-fill-2);
  padding: 12px;
  border-radius: 4px;
  overflow-x: auto;
}

.content-body :deep(code) {
  background: var(--color-fill-2);
  padding: 2px 6px;
  border-radius: 2px;
  font-family: 'Courier New', monospace;
}
</style>
