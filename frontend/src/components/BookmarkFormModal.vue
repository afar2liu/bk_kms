<template>
  <a-modal
    v-model:visible="modalVisible"
    :title="isEdit ? '编辑书签' : '创建书签'"
    width="700px"
    @ok="handleSubmit"
    @cancel="handleCancel"
  >
    <a-form
      ref="formRef"
      :model="formData"
      :rules="rules"
      layout="vertical"
    >
      <a-form-item field="url" label="URL" required>
        <a-input
          v-model="formData.url"
          placeholder="请输入书签 URL"
          allow-clear
        />
      </a-form-item>

      <a-form-item field="title" label="标题">
        <a-input
          v-model="formData.title"
          placeholder="留空则自动获取"
          allow-clear
        />
      </a-form-item>

      <a-form-item field="excerpt" label="摘要">
        <a-textarea
          v-model="formData.excerpt"
          placeholder="留空则自动获取"
          :auto-size="{ minRows: 3, maxRows: 5 }"
          allow-clear
        />
      </a-form-item>

      <a-form-item v-if="isEdit" field="author" label="作者">
        <a-input
          v-model="formData.author"
          placeholder="作者"
          allow-clear
        />
      </a-form-item>

      <a-form-item field="tags" label="标签">
        <a-select
          v-model="selectedTagNames"
          placeholder="选择或输入标签"
          multiple
          allow-create
          allow-clear
        >
          <a-option
            v-for="tag in allTags"
            :key="tag.id"
            :value="tag.name"
          >
            {{ tag.name }}
          </a-option>
        </a-select>
      </a-form-item>

      <a-form-item field="create_archive" label="归档选项">
        <a-checkbox v-model="formData.create_archive">
          创建归档（自动获取网页内容）
        </a-checkbox>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { Message } from '@arco-design/web-vue'
import { createBookmark, updateBookmark } from '@/api/bookmark'
import type { Bookmark, Tag, CreateBookmarkRequest, UpdateBookmarkRequest } from '@/types'

interface Props {
  visible: boolean
  bookmark: Bookmark | null
  allTags: Tag[]
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

const isEdit = computed(() => !!props.bookmark)

const formRef = ref()
const formData = ref<any>({
  url: '',
  title: '',
  excerpt: '',
  tags: [],
  create_archive: false
})

const selectedTagNames = ref<string[]>([])

const rules = {
  url: [
    { required: true, message: '请输入 URL' },
    { type: 'url', message: '请输入有效的 URL' }
  ]
}

// 监听 bookmark 变化，填充表单
watch(
  () => props.bookmark,
  (bookmark) => {
    if (bookmark) {
      const tags = bookmark.tags || []
      formData.value = {
        id: bookmark.id,
        url: bookmark.url,
        title: bookmark.title,
        excerpt: bookmark.excerpt,
        author: bookmark.author,
        tags: tags.map((t) => ({ name: t.name })),
        create_archive: bookmark.is_archive
      }
      selectedTagNames.value = tags.map((t) => t.name)
    } else {
      resetForm()
    }
  },
  { immediate: true }
)

// 重置表单
function resetForm() {
  formData.value = {
    url: '',
    title: '',
    excerpt: '',
    tags: [],
    create_archive: false
  }
  selectedTagNames.value = []
  formRef.value?.clearValidate()
}

// 提交表单
async function handleSubmit() {
  try {
    await formRef.value?.validate()
  } catch (error) {
    // 验证失败
    return
  }

  // 转换标签格式
  formData.value.tags = selectedTagNames.value.map((name) => ({ name }))

  console.log('提交数据:', formData.value)
  console.log('是否编辑模式:', isEdit.value)

  try {
    if (isEdit.value) {
      await updateBookmark(formData.value as UpdateBookmarkRequest)
      Message.success('更新成功')
    } else {
      await createBookmark(formData.value as CreateBookmarkRequest)
      Message.success('创建成功')
    }
    emit('success')
    modalVisible.value = false
  } catch (error: any) {
    console.error('保存书签失败:', error)
    Message.error(error.message || (isEdit.value ? '更新失败' : '创建失败'))
  }
}

// 取消
function handleCancel() {
  resetForm()
}
</script>
