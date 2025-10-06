<template>
  <div class="tag-manage-container">
    <a-card title="标签管理">
      <a-spin :loading="loading">
        <a-space wrap size="large">
          <div
            v-for="tag in tagList"
            :key="tag.id"
            class="tag-item"
          >
            <a-tag
              :color="getTagColor(tag.count || 0)"
              size="large"
            >
              <template #icon>
                <icon-tag />
              </template>
              <span class="tag-name">
                {{ tag.name }}
              </span>
              <a-badge
                :count="tag.count || 0"
                :max-count="999"
                style="margin-left: 8px"
              />
            </a-tag>
            <a-button
              type="text"
              size="small"
              class="edit-btn"
              @click="startEdit(tag)"
            >
              <template #icon>
                <icon-edit />
              </template>
            </a-button>
          </div>
        </a-space>

        <a-empty v-if="tagList.length === 0" description="暂无标签" />
      </a-spin>
    </a-card>

    <!-- 编辑标签对话框 -->
    <a-modal
      v-model:visible="editModalVisible"
      title="重命名标签"
      @ok="handleUpdate"
      @cancel="cancelEdit"
    >
      <a-form :model="editForm" layout="vertical">
        <a-form-item label="标签名称" required>
          <a-input
            v-model="editForm.name"
            placeholder="请输入新的标签名称"
            allow-clear
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { IconTag, IconEdit } from '@arco-design/web-vue/es/icon'
import { getTagList, updateTag } from '@/api/tag'
import type { Tag } from '@/types'

const tagList = ref<Tag[]>([])
const loading = ref(false)

const editModalVisible = ref(false)
const editForm = ref({
  id: 0,
  name: ''
})

// 获取标签列表
async function fetchTagList() {
  loading.value = true
  try {
    const { data } = await getTagList()
    if (data.data) {
      tagList.value = data.data.sort((a, b) => (b.count || 0) - (a.count || 0))
    }
  } catch (error) {
    Message.error('获取标签列表失败')
  } finally {
    loading.value = false
  }
}

// 根据使用次数获取标签颜色
function getTagColor(count: number) {
  if (count >= 10) return 'red'
  if (count >= 5) return 'orangered'
  if (count >= 3) return 'orange'
  return 'arcoblue'
}

// 开始编辑
function startEdit(tag: Tag) {
  editForm.value = {
    id: tag.id,
    name: tag.name
  }
  editModalVisible.value = true
}

// 取消编辑
function cancelEdit() {
  editForm.value = {
    id: 0,
    name: ''
  }
}

// 更新标签
async function handleUpdate() {
  if (!editForm.value.name.trim()) {
    Message.warning('标签名称不能为空')
    return
  }

  try {
    await updateTag(editForm.value.id, editForm.value.name)
    Message.success('重命名成功')
    editModalVisible.value = false
    fetchTagList()
  } catch (error) {
    Message.error('重命名失败')
  }
}


onMounted(() => {
  fetchTagList()
})
</script>

<style scoped>
.tag-manage-container {
  padding: 20px;
}

.tag-item {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.tag-name {
  user-select: none;
}

.edit-btn {
  opacity: 0;
  transition: opacity 0.2s;
  color: var(--color-text-3);
}

.tag-item:hover .edit-btn {
  opacity: 1;
}

.edit-btn:hover {
  color: rgb(var(--primary-6));
}
</style>
