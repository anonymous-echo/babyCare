<template>
  <view class="detail-page">
    <wd-navbar
      :title="isEditing ? '编辑成长记录' : '成长记录详情'"
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
        <wd-card title="测量数据">
          <wd-cell-group border>
            <wd-input
              v-model="editForm.height"
              label="身高 (cm)"
              type="number"
              placeholder="请输入身高"
            />
            <wd-input
              v-model="editForm.weight"
              label="体重 (kg)"
              type="number"
              placeholder="请输入体重"
            />
            <wd-input
              v-model="editForm.headCircumference"
              label="头围 (cm)"
              type="number"
              placeholder="请输入头围"
            />
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
        <!-- 测量数据卡片 -->
        <wd-card title="测量数据">
          <view class="detail-content">
            <view class="detail-item">
              <text class="detail-label">测量时间</text>
              <text class="detail-value">{{ formattedMeasureTime }}</text>
            </view>
            <view v-if="record.height" class="detail-item">
              <text class="detail-label">身高</text>
              <text class="detail-value">{{ record.height }} cm</text>
            </view>
            <view v-if="record.weight" class="detail-item">
              <text class="detail-label">体重</text>
              <text class="detail-value">{{ record.weight }} kg</text>
            </view>
            <view v-if="record.headCircumference" class="detail-item">
              <text class="detail-label">头围</text>
              <text class="detail-value">{{ record.headCircumference }} cm</text>
            </view>
            <view v-if="record.note" class="detail-item detail-item-column">
              <text class="detail-label">备注</text>
              <text class="detail-value">{{ record.note }}</text>
            </view>
          </view>
        </wd-card>

        <!-- 其他信息卡片 -->
        <wd-card title="其他信息">
          <view class="detail-content">
            <view class="detail-item">
              <text class="detail-label">创建时间</text>
              <text class="detail-value">{{ formattedCreateTime }}</text>
            </view>
            <view class="detail-item">
              <text class="detail-label">创建人</text>
              <text class="detail-value">{{ record.createBy }}</text>
            </view>
          </view>
        </wd-card>

        <!-- 查看模式按钮 -->
        <view class="actions">
          <wd-button type="primary" block @click="handleEdit">编辑记录</wd-button>
          <wd-button type="error" block @click="handleDelete">删除记录</wd-button>
        </view>
      </template>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import {
  apiGetGrowthRecordById,
  apiDeleteGrowthRecord,
  apiUpdateGrowthRecord,
} from '@/api/growth'
import type { GrowthRecordResponse } from '@/api/growth'
import { formatDate } from '@/utils/date'

// 页面参数
const recordId = ref('')
const loading = ref(true)
const saving = ref(false)
const isEditing = ref(false)
const record = ref<GrowthRecordResponse | null>(null)

// 编辑表单
const editForm = ref({
  height: '',
  weight: '',
  headCircumference: '',
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
    record.value = await apiGetGrowthRecordById(recordId.value)
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
const formattedMeasureTime = computed(() => {
  if (!record.value) return ''
  return formatDate(record.value.measureTime, 'YYYY-MM-DD HH:mm')
})

const formattedCreateTime = computed(() => {
  if (!record.value) return ''
  return formatDate(record.value.createTime, 'YYYY-MM-DD HH:mm')
})

// 进入编辑模式
function handleEdit() {
  if (!record.value) return

  // 填充编辑表单
  editForm.value = {
    height: record.value.height?.toString() || '',
    weight: record.value.weight?.toString() || '',
    headCircumference: record.value.headCircumference?.toString() || '',
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

  // 构建更新数据（只包含有值的字段）
  const updateData: any = {}

  if (editForm.value.height) {
    updateData.height = parseFloat(editForm.value.height)
  }
  if (editForm.value.weight) {
    updateData.weight = parseFloat(editForm.value.weight)
  }
  if (editForm.value.headCircumference) {
    updateData.headCircumference = parseFloat(editForm.value.headCircumference)
  }
  if (editForm.value.note) {
    updateData.note = editForm.value.note
  }

  // 验证至少有一个字段
  if (Object.keys(updateData).length === 0) {
    uni.showToast({
      title: '请至少填写一项数据',
      icon: 'none',
    })
    return
  }

  saving.value = true
  try {
    await apiUpdateGrowthRecord(recordId.value, updateData)
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
          await apiDeleteGrowthRecord(recordId.value)
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
      content: '编辑内容尚未保存，确定要离开吗？',
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
      content: '编辑内容尚未保存，确定要离开吗？',
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
