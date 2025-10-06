<template>
  <div class="bookmark-list-container">
    <!-- 搜索和操作栏 -->
    <div class="toolbar">
      <a-space size="medium">
        <a-input-search
          v-model="searchParams.keyword"
          placeholder="搜索书签（URL、标题、摘要、内容）"
          style="width: 320px"
          @search="handleSearch"
          @clear="handleSearch"
          allow-clear
        />
        <a-select
          v-model="selectedTags"
          placeholder="选择标签筛选"
          style="width: 240px"
          multiple
          allow-clear
          @change="handleSearch"
        >
          <a-option v-for="tag in allTags" :key="tag.id" :value="tag.name">
            {{ tag.name }} ({{ tag.count }})
          </a-option>
        </a-select>
      </a-space>

      <a-space size="medium">
        <a-button type="primary" @click="showImportModal">
          <template #icon><icon-import /></template>
          导入书签
        </a-button>
        <a-button type="primary" @click="showCreateModal">
          <template #icon><icon-plus /></template>
          创建书签
        </a-button>
        <a-button
          type="primary"
          status="danger"
          :disabled="selectedRowKeys.length === 0"
          @click="handleBatchDelete"
        >
          <template #icon><icon-delete /></template>
          批量删除 ({{ selectedRowKeys.length }})
        </a-button>
      </a-space>
    </div>

    <!-- 书签表格 -->
    <a-table
      :columns="columns"
      :data="bookmarkList"
      :loading="loading"
      :pagination="pagination"
      :row-selection="rowSelection"
      @selection-change="handleSelectionChange"
      @select="handleSelect"
      @select-all="handleSelectAll"
      @page-change="handlePageChange"
      @page-size-change="handlePageSizeChange"
      row-key="id"
    >
      <template #title="{ record }">
        <a :href="record.url" target="_blank" class="bookmark-title">
          {{ record.title }}
        </a>
      </template>

      <template #tags="{ record }">
        <a-space wrap>
          <a-tag v-for="tag in record.tags" :key="tag.id" color="arcoblue">
            {{ tag.name }}
          </a-tag>
        </a-space>
      </template>

      <template #created_at="{ record }">
        {{ formatDate(record.created_at) }}
      </template>

      <template #actions="{ record }">
        <a-space>
          <a-button type="text" size="small" @click="viewContent(record)">
            查看
          </a-button>
          <a-button type="text" size="small" @click="editBookmark(record)">
            编辑
          </a-button>
          <a-button
            type="text"
            size="small"
            status="danger"
            @click="deleteBookmark(record)"
          >
            删除
          </a-button>
        </a-space>
      </template>
    </a-table>

    <!-- 导入对话框 -->
    <ImportModal
      v-model:visible="importModalVisible"
      @success="handleSearch"
    />

    <!-- 创建/编辑对话框 -->
    <BookmarkFormModal
      v-model:visible="formModalVisible"
      :bookmark="currentBookmark"
      :all-tags="allTags"
      @success="handleSearch"
    />

    <!-- 内容查看对话框 -->
    <ContentViewModal
      v-model:visible="contentModalVisible"
      :bookmark-id="currentBookmarkId"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import {
  IconPlus,
  IconDelete,
  IconImport
} from '@arco-design/web-vue/es/icon'
import { getBookmarkList, deleteBookmarks } from '@/api/bookmark'
import { getTagList } from '@/api/tag'
import type { Bookmark, Tag } from '@/types'
import ImportModal from '@/components/ImportModal.vue'
import BookmarkFormModal from '@/components/BookmarkFormModal.vue'
import ContentViewModal from '@/components/ContentViewModal.vue'

// 搜索参数
const searchParams = reactive({
  keyword: '',
  page: 1,
  page_size: 10
})

const selectedTags = ref<string[]>([])
const bookmarkList = ref<Bookmark[]>([])
const allTags = ref<Tag[]>([])
const loading = ref(false)
const total = ref(0)

// 选中的行
const selectedRowKeys = ref<number[]>([])
console.log('初始化 selectedRowKeys:', selectedRowKeys.value)

// 使用简单配置，事件驱动更新，避免把数组快照传进去
const rowSelection = reactive({
  type: 'checkbox',
  showCheckedAll: true,
  onlyCurrent: false
})

function handleSelectionChange(rowKeys: number[]) {
  console.log('[selection-change] rowKeys:', rowKeys)
  selectedRowKeys.value = rowKeys
}

function handleSelect(rowKeys: number[], row: any) {
  console.log('[select] rowKeys:', rowKeys, 'row:', row)
}

function handleSelectAll(checked: boolean) {
  console.log('[select-all] checked:', checked)
}

// 对话框状态
const importModalVisible = ref(false)
const formModalVisible = ref(false)
const contentModalVisible = ref(false)
const currentBookmark = ref<Bookmark | null>(null)
const currentBookmarkId = ref(0)

// 表格列定义
const columns = [
  { title: '标题', slotName: 'title', ellipsis: true, tooltip: true },
  { title: '摘要', dataIndex: 'excerpt', ellipsis: true, tooltip: true, width: 200 },
  { title: '作者', dataIndex: 'author', width: 120 },
  { title: '标签', slotName: 'tags', width: 200 },
  { title: '创建时间', slotName: 'created_at', width: 180 },
  { title: '操作', slotName: 'actions', width: 200, fixed: 'right' }
]

// 分页配置
const pagination = computed(() => ({
  current: searchParams.page,
  pageSize: searchParams.page_size,
  total: total.value,
  showTotal: true,
  showPageSize: true
}))

// 获取书签列表
async function fetchBookmarkList() {
  loading.value = true
  try {
    const params = {
      ...searchParams,
      tags: selectedTags.value.join(',')
    }
    const { data } = await getBookmarkList(params)
    if (data.data) {
      bookmarkList.value = data.data.rows
      total.value = data.data.total
      console.log('书签列表加载完成，数量:', data.data.rows.length)
      console.log('第一条书签 ID:', data.data.rows[0]?.id)
    }
  } catch (error) {
    Message.error('获取书签列表失败')
  } finally {
    loading.value = false
  }
}

// 获取标签列表
async function fetchTagList() {
  try {
    const { data } = await getTagList()
    if (data.data) {
      allTags.value = data.data
    }
  } catch (error) {
    console.error('获取标签列表失败', error)
  }
}

// 搜索
function handleSearch() {
  searchParams.page = 1
  fetchBookmarkList()
}

// 分页变化
function handlePageChange(page: number) {
  searchParams.page = page
  fetchBookmarkList()
}

function handlePageSizeChange(pageSize: number) {
  searchParams.page_size = pageSize
  searchParams.page = 1
  fetchBookmarkList()
}

// 显示导入对话框
function showImportModal() {
  importModalVisible.value = true
}

// 显示创建对话框
function showCreateModal() {
  currentBookmark.value = null
  formModalVisible.value = true
}

// 编辑书签
function editBookmark(bookmark: Bookmark) {
  currentBookmark.value = bookmark
  formModalVisible.value = true
}

// 查看内容
function viewContent(bookmark: Bookmark) {
  currentBookmarkId.value = bookmark.id
  contentModalVisible.value = true
}

// 删除书签
function deleteBookmark(bookmark: Bookmark) {
  Modal.confirm({
    title: '确认删除',
    content: `确定要删除书签"${bookmark.title}"吗？`,
    onOk: async () => {
      try {
        await deleteBookmarks([bookmark.id])
        Message.success('删除成功')
        fetchBookmarkList()
      } catch (error) {
        Message.error('删除失败')
      }
    }
  })
}

// 批量删除
function handleBatchDelete() {
  Modal.confirm({
    title: '确认批量删除',
    content: `确定要删除选中的 ${selectedRowKeys.value.length} 个书签吗？`,
    onOk: async () => {
      try {
        await deleteBookmarks(selectedRowKeys.value)
        Message.success('批量删除成功')
        selectedRowKeys.value = []
        fetchBookmarkList()
      } catch (error) {
        Message.error('批量删除失败')
      }
    }
  })
}

// 格式化日期
function formatDate(timestamp: number) {
  return new Date(timestamp * 1000).toLocaleString('zh-CN')
}

// 监控 selectedRowKeys 变化
watch(selectedRowKeys, (newVal, oldVal) => {
  console.log('selectedRowKeys 发生变化:')
  console.log('  旧值:', oldVal)
  console.log('  新值:', newVal)
  console.log('  批量删除按钮状态:', newVal.length === 0 ? '禁用' : '启用')
}, { deep: true })

// 监控 rowSelection 变化
watch(rowSelection, (newVal) => {
  console.log('rowSelection 配置变化:', newVal)
}, { deep: true })

onMounted(() => {
  console.log('组件挂载完成')
  fetchBookmarkList()
  fetchTagList()
})
</script>

<style scoped>
.bookmark-list-container {
  padding: 20px;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
  padding: 16px;
  background: var(--color-bg-2);
  border-radius: 4px;
}

.bookmark-title {
  color: rgb(var(--primary-6));
  text-decoration: none;
}

.bookmark-title:hover {
  text-decoration: underline;
}
</style>
