<template>
  <view class="detail-page">
    <wd-navbar
      :title="isEditing ? '编辑睡眠记录' : '睡眠记录详情'"
      left-arrow
      safeAreaInsetTop
      @click-left="handleBack"
    >
      <template #capsule>
        <wd-navbar-capsule @back="handleBack" @back-home="handleBackHome" />
      </template>
    </wd-navbar>

    <view v-if="loading" class="loading">
      <wd-loading />
    </view>

    <view v-else-if="record" class="content">
      <!-- 编辑模式 -->
      <template v-if="isEditing">
        <wd-card title="睡眠信息">
          <wd-cell-group border>
            <wd-cell title="睡眠类型">
              <wd-radio-group v-model="editForm.sleepType">
                <wd-radio value="nap">小睡</wd-radio>
                <wd-radio value="night">夜间睡眠</wd-radio>
              </wd-radio-group>
            </wd-cell>
            <wd-textarea
              v-model="editForm.note"
              label="备注"
              placeholder="请输入备注"
              :maxlength="200"
            />
          </wd-cell-group>
        </wd-card>

        <!-- 编辑模式按钮 -->
        <view class="actions">
          <wd-button type="success" block @click="handleSave" :loading="saving">
            保存
          </wd-button>
          <wd-button type="default" block @click="cancelEdit">取消</wd-button>
        </view>
      </template>

      <!-- 查看模式 -->
      <template v-else>
      <!-- 基本信息卡片 -->
      <wd-card title="基本信息">
        <wd-cell-group border>
          <wd-cell title="睡眠类型" :value="sleepTypeText" />
          <wd-cell title="开始时间" :value="formattedStartTime" />
          <wd-cell v-if="record.endTime" title="结束时间" :value="formattedEndTime" />
          <wd-cell title="睡眠时长" :value="durationText" />
        </wd-cell-group>
      </wd-card>

      <!-- 详细信息卡片 -->
      <wd-card title="详细信息">
        <wd-cell-group border>
          <wd-cell v-if="record.note" title="备注" :value="record.note" />
        </wd-cell-group>
      </wd-card>

      <!-- 其他信息卡片 -->
      <wd-card title="其他信息">
        <wd-cell-group border>
          <wd-cell title="创建时间" :value="formattedCreateTime" />
          <wd-cell title="创建人" :value="record.createBy" />
        </wd-cell-group>
      </wd-card>

      <!-- 操作按钮 -->
      <view class="actions">
        <wd-button type="primary" block @click="handleEdit">编辑记录</wd-button>
        <wd-button type="error" block @click="handleDelete">删除记录</wd-button>
      </view>
      </template>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { apiGetSleepRecordById, apiDeleteSleepRecord, apiUpdateSleepRecord } from '@/api/sleep'
import type { SleepRecordResponse } from '@/api/sleep'
import { formatDate } from '@/utils/date'
import { formatDuration } from '@/utils/common'

// 页面参数
const recordId = ref('')
const loading = ref(true)
const saving = ref(false)
const isEditing = ref(false)
const record = ref<SleepRecordResponse | null>(null)

// 编辑表单
const editForm = ref({
  sleepType: 'nap' as 'nap' | 'night',
  note: '',
})

// 页面加载
onLoad((options) => {
  if (options?.id) {
    recordId.value = options.id
    loadRecord()
  }
})

// 加载记录详情
async function loadRecord() {
  if (!recordId.value) return

  loading.value = true
  try {
    record.value = await apiGetSleepRecordById(recordId.value)
  } catch (error: any) {
    uni.showToast({
      title: error.message || '加载失败',
      icon: 'none',
    })
    setTimeout(() => {
      uni.navigateBack()
    }, 1500)
  } finally {
    loading.value = false
  }
}

// 计算属性
const sleepTypeText = computed(() => {
  if (!record.value) return ''
  const map: Record<string, string> = {
    nap: '小睡',
    night: '夜间睡眠',
  }
  return map[record.value.sleepType] || record.value.sleepType
})


const durationText = computed(() => {
  if (!record.value) return ''
  return formatDuration(record.value.duration || 0)
})

const formattedStartTime = computed(() => {
  if (!record.value) return ''
  return formatDate(record.value.startTime, 'YYYY-MM-DD HH:mm')
})

const formattedEndTime = computed(() => {
  if (!record.value || !record.value.endTime) return ''
  return formatDate(record.value.endTime, 'YYYY-MM-DD HH:mm')
})

const formattedCreateTime = computed(() => {
  if (!record.value) return ''
  return formatDate(record.value.createTime, 'YYYY-MM-DD HH:mm')
})

// 编辑记录
function handleEdit() {
  if (!record.value) return

  // 填充编辑表单
  editForm.value = {
    sleepType: record.value.sleepType,
    note: record.value.note || '',
  }

  isEditing.value = true
}

// 取消编辑
function cancelEdit() {
  isEditing.value = false
}

// 保存编辑
async function handleSave() {
  if (!recordId.value) return

  // 构建更新数据
  const updateData: any = {
    sleepType: editForm.value.sleepType,
  }

  if (editForm.value.note) {
    updateData.note = editForm.value.note
  }

  saving.value = true
  try {
    await apiUpdateSleepRecord(recordId.value, updateData)
    uni.showToast({
      title: '保存成功',
      icon: 'success',
    })

    // 重新加载数据
    await loadRecord()
    isEditing.value = false
  } catch (error: any) {
    uni.showToast({
      title: error.message || '保存失败',
      icon: 'none',
    })
  } finally {
    saving.value = false
  }
}

// 删除记录
async function handleDelete() {
  uni.showModal({
    title: '确认删除',
    content: '确定要删除这条记录吗？',
    success: async (res) => {
      if (res.confirm) {
        try {
          await apiDeleteSleepRecord(recordId.value)
          uni.showToast({
            title: '删除成功',
            icon: 'success',
          })
          setTimeout(() => {
            uni.navigateBack()
          }, 1500)
        } catch (error: any) {
          uni.showToast({
            title: error.message || '删除失败',
            icon: 'none',
          })
        }
      }
    },
  })
}

function handleBack() {
  if (isEditing.value) {
    uni.showModal({
      title: '提示',
      content: '编辑内容尚未保存,确定要离开吗?',
      success: (res) => {
        if (res.confirm) {
          uni.navigateBack()
        }
      },
    })
  } else {
    uni.navigateBack()
  }
}

function handleBackHome() {
  if (isEditing.value) {
    uni.showModal({
      title: '提示',
      content: '编辑内容尚未保存,确定要离开吗?',
      success: (res) => {
        if (res.confirm) {
          uni.switchTab({
            url: '/pages/index/index',
          })
        }
      },
    })
  } else {
    uni.switchTab({
      url: '/pages/index/index',
    })
  }
}
</script>

<style lang="scss" scoped>
.detail-page {
  min-height: 100vh;
  background-color: #f5f5f5;
}

.loading {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 60vh;
}

.content {
  padding: 20rpx;
}

.actions {
  margin-top: 40rpx;
  display: flex;
  flex-direction: column;
  gap: 20rpx;
  padding: 0 20rpx 40rpx;
}
</style>
