<template>
  <a-modal
    v-model:visible="modalVisible"
    title="导入书签"
    :footer="false"
    :mask-closable="false"
    width="600px"
    @cancel="handleCancel"
  >
    <div class="import-container">
      <!-- 文件上传 -->
      <a-upload
        v-if="!importing"
        :custom-request="handleUpload"
        accept=".html"
        :show-file-list="false"
        drag
      >
        <template #upload-button>
          <div class="upload-area">
            <div class="upload-icon">
              <icon-upload :size="48" />
            </div>
            <div class="upload-text">
              <p><strong>点击或拖拽文件到此处上传</strong></p>
              <p class="upload-hint">仅支持 Netscape Bookmark 格式的 .html 文件</p>
            </div>
          </div>
        </template>
      </a-upload>

      <!-- 选项 -->
      <div v-if="!importing" class="options">
        <a-space direction="vertical" size="medium">
          <a-checkbox v-model="generateTag">
            自动将文件夹名称作为标签
          </a-checkbox>
          <a-checkbox v-model="createArchive">
            创建归档（自动获取网页内容）
          </a-checkbox>
        </a-space>
      </div>

      <!-- 导入进度 -->
      <div v-if="importing" class="progress-container">
        <a-progress
          :percent="progressPercent"
          :status="progressStatus"
          :stroke-width="8"
        />
        <div class="progress-text">
          {{ progressText }}
        </div>

        <!-- 日志列表 -->
        <div class="log-container">
          <div
            v-for="(log, index) in logs"
            :key="index"
            :class="['log-item', `log-${log.type}`]"
          >
            <span class="log-icon">{{ getLogIcon(log.type) }}</span>
            <span class="log-message">{{ log.message }}</span>
          </div>
        </div>
      </div>

      <!-- 统计信息 -->
      <div v-if="completed" class="summary">
        <a-result status="success" title="导入完成！">
          <template #subtitle>
            <div class="summary-stats">
              <div class="stat-item">
                <div class="stat-value success">{{ stats.success }}</div>
                <div class="stat-label">成功</div>
              </div>
              <div class="stat-item">
                <div class="stat-value skip">{{ stats.skip }}</div>
                <div class="stat-label">跳过</div>
              </div>
              <div class="stat-item">
                <div class="stat-value error">{{ stats.error }}</div>
                <div class="stat-label">失败</div>
              </div>
            </div>
          </template>
          <template #extra>
            <a-button type="primary" @click="handleClose">
              关闭
            </a-button>
          </template>
        </a-result>
      </div>
    </div>
  </a-modal>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { Message } from '@arco-design/web-vue'
import { IconUpload } from '@arco-design/web-vue/es/icon'
import { importBookmarks } from '@/api/bookmark'
import type { ImportProgressEvent } from '@/types'

interface Props {
  visible: boolean
}

interface Emits {
  (e: 'update:visible', value: boolean): void
  (e: 'success'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const modalVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

const generateTag = ref(true)
const createArchive = ref(false)
const importing = ref(false)
const completed = ref(false)
const progressPercent = ref(0)
const progressText = ref('')
const progressStatus = ref<'normal' | 'success' | 'danger'>('normal')

interface LogItem {
  type: string
  message: string
}

const logs = ref<LogItem[]>([])
const stats = ref({
  success: 0,
  skip: 0,
  error: 0
})

// 处理文件上传
async function handleUpload(option: any) {
  const file = option.fileItem.file

  if (!file.name.toLowerCase().endsWith('.html')) {
    Message.error('仅支持 .html 文件')
    return
  }

  // 重置状态
  importing.value = true
  completed.value = false
  logs.value = []
  stats.value = { success: 0, skip: 0, error: 0 }
  progressPercent.value = 0
  progressStatus.value = 'normal'

  try {
    await importBookmarks(file, generateTag.value, createArchive.value, handleProgress)
  } catch (error: any) {
    Message.error(error.message || '导入失败')
    importing.value = false
  }
}

// 处理进度事件
function handleProgress(event: ImportProgressEvent) {
  const { type, message, current, total } = event

  console.log('收到进度事件:', event)

  // 更新进度
  if (total > 0 && current >= 0) {
    const percent = Math.round((current / total) * 100)
    // 确保百分比在 0-100 之间
    progressPercent.value = Math.min(100, Math.max(0, percent))
    progressText.value = `${current} / ${total}`
    console.log('计算进度:', { current, total, percent: progressPercent.value })
  }

  // 添加日志
  logs.value.push({ type, message })

  // 滚动到底部
  setTimeout(() => {
    const container = document.querySelector('.log-container')
    if (container) {
      container.scrollTop = container.scrollHeight
    }
  }, 100)

  // 更新统计（不要累加进度）
  if (type === 'success') {
    stats.value.success++
  } else if (type === 'error') {
    stats.value.error++
  } else if (type === 'progress' && message.includes('跳过')) {
    stats.value.skip++
  }

  // 完成
  if (type === 'complete') {
    completed.value = true
    progressStatus.value = stats.value.error > 0 ? 'danger' : 'success'
    emit('success')
  }
}

// 获取日志图标
function getLogIcon(type: string) {
  const icons: Record<string, string> = {
    progress: '📝',
    success: '✅',
    error: '❌',
    complete: '🎉'
  }
  return icons[type] || '📝'
}

// 关闭对话框
function handleClose() {
  modalVisible.value = false
}

// 取消
function handleCancel() {
  if (importing.value && !completed.value) {
    Message.warning('导入进行中，请等待完成')
    return
  }
  handleClose()
}

// 重置状态
watch(modalVisible, (visible) => {
  if (!visible) {
    setTimeout(() => {
      importing.value = false
      completed.value = false
      logs.value = []
      progressPercent.value = 0
      stats.value = { success: 0, skip: 0, error: 0 }
    }, 300)
  }
})
</script>

<style scoped>
.import-container {
  padding: 20px 0;
}

.upload-area {
  padding: 60px 20px;
  text-align: center;
  border: 2px dashed var(--color-border-2);
  border-radius: 4px;
  transition: all 0.3s;
  cursor: pointer;
}

.upload-area:hover {
  border-color: rgb(var(--primary-6));
  background: var(--color-fill-1);
}

.upload-icon {
  color: rgb(var(--primary-6));
  margin-bottom: 16px;
}

.upload-text p {
  margin: 8px 0;
}

.upload-hint {
  font-size: 12px;
  color: var(--color-text-3);
}

.options {
  margin-top: 20px;
  padding: 16px;
  background: var(--color-fill-1);
  border-radius: 4px;
}

.progress-container {
  margin-top: 20px;
}

.progress-text {
  text-align: center;
  margin: 16px 0;
  font-size: 14px;
  color: var(--color-text-2);
}

.log-container {
  max-height: 300px;
  overflow-y: auto;
  background: var(--color-fill-1);
  border-radius: 4px;
  padding: 12px;
  margin-top: 16px;
}

.log-item {
  padding: 8px 12px;
  margin-bottom: 8px;
  border-radius: 4px;
  font-size: 13px;
  display: flex;
  align-items: center;
  gap: 8px;
  animation: slideIn 0.3s ease;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateX(-10px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

.log-progress {
  background: var(--color-primary-light-1);
  color: rgb(var(--primary-6));
}

.log-success {
  background: var(--color-success-light-1);
  color: rgb(var(--success-6));
}

.log-error {
  background: var(--color-danger-light-1);
  color: rgb(var(--danger-6));
}

.log-complete {
  background: var(--color-purple-light-1);
  color: rgb(var(--purple-6));
  font-weight: 600;
}

.summary {
  margin-top: 20px;
}

.summary-stats {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  margin-top: 20px;
}

.stat-item {
  text-align: center;
  padding: 16px;
  background: var(--color-fill-1);
  border-radius: 8px;
}

.stat-value {
  font-size: 32px;
  font-weight: bold;
  margin-bottom: 8px;
}

.stat-value.success {
  color: rgb(var(--success-6));
}

.stat-value.skip {
  color: rgb(var(--warning-6));
}

.stat-value.error {
  color: rgb(var(--danger-6));
}

.stat-label {
  font-size: 14px;
  color: var(--color-text-2);
}
</style>
