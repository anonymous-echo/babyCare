<template>
  <view class="detail-page">
    <wd-navbar
      :title="isEditing ? '编辑喂养记录' : '喂养记录详情'"
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
        <wd-card title="备注信息">
          <wd-cell-group border>
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
          <wd-cell title="喂养类型" :value="feedingTypeText" />
          <wd-cell title="喂养时间" :value="formattedFeedingTime" />
          <wd-cell
            v-if="record.actualCompleteTime"
            title="完成时间"
            :value="formattedCompleteTime"
          />
        </wd-cell-group>
      </wd-card>

      <!-- 详细信息卡片 -->
      <wd-card title="详细信息">
        <wd-cell-group border>
          <!-- 母乳喂养 -->
          <template v-if="record.feedingType === 'breast' && record.detail.type === 'breast'">
            <wd-cell title="喂养侧" :value="breastSideText" />
            <wd-cell title="时长" :value="durationText" />
            <wd-cell
              v-if="record.detail.leftDuration"
              title="左侧时长"
              :value="formatDuration(record.detail.leftDuration)"
            />
            <wd-cell
              v-if="record.detail.rightDuration"
              title="右侧时长"
              :value="formatDuration(record.detail.rightDuration)"
            />
          </template>

          <!-- 奶瓶喂养 -->
          <template v-if="record.feedingType === 'bottle' && record.detail.type === 'bottle'">
            <wd-cell title="奶类型" :value="bottleTypeText" />
            <wd-cell title="奶量" :value="`${record.amount} ${record.detail.unit || 'ml'}`" />
            <wd-cell
              v-if="record.detail.remaining"
              title="剩余量"
              :value="`${record.detail.remaining} ${record.detail.unit || 'ml'}`"
            />
          </template>

          <!-- 辅食 -->
          <template v-if="record.feedingType === 'food' && record.detail.type === 'food'">
            <wd-cell title="辅食名称" :value="record.detail.foodName || '未知'" />
          </template>

          <!-- 备注 -->
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
import { apiGetFeedingRecordById, apiDeleteFeedingRecord, apiUpdateFeedingRecord } from '@/api/feeding'
import type { FeedingRecordResponse } from '@/api/feeding'
import { formatDate } from '@/utils/date'
import { formatDuration } from '@/utils/common'

// 页面参数
const recordId = ref('')
const loading = ref(true)
const saving = ref(false)
const isEditing = ref(false)
const record = ref<FeedingRecordResponse | null>(null)

// 编辑表单(只允许编辑备注)
const editForm = ref({
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
    record.value = await apiGetFeedingRecordById(recordId.value)
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
const feedingTypeText = computed(() => {
  if (!record.value) return ''
  const map: Record<string, string> = {
    breast: '母乳喂养',
    bottle: '奶瓶喂养',
    food: '辅食',
  }
  return map[record.value.feedingType] || record.value.feedingType
})

const breastSideText = computed(() => {
  if (!record.value || record.value.feedingType !== 'breast') return ''
  const detail = record.value.detail
  if (detail.type !== 'breast') return ''
  const side = detail.side
  const map: Record<string, string> = {
    left: '左侧',
    right: '右侧',
    both: '双侧',
  }
  return map[side] || side
})

const bottleTypeText = computed(() => {
  if (!record.value || record.value.feedingType !== 'bottle') return ''
  const detail = record.value.detail
  if (detail.type !== 'bottle') return ''
  const type = detail.bottleType
  const map: Record<string, string> = {
    formula: '配方奶',
    'breast-milk': '母乳',
  }
  return map[type] || type
})

const durationText = computed(() => {
  if (!record.value) return ''
  return formatDuration(record.value.duration || 0)
})

const formattedFeedingTime = computed(() => {
  if (!record.value) return ''
  return formatDate(record.value.feedingTime, 'YYYY-MM-DD HH:mm')
})

const formattedCompleteTime = computed(() => {
  if (!record.value || !record.value.actualCompleteTime) return ''
  return formatDate(record.value.actualCompleteTime, 'YYYY-MM-DD HH:mm')
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

  // 构建更新数据(只更新备注)
  const updateData: any = {}

  if (editForm.value.note) {
    updateData.note = editForm.value.note
  }

  saving.value = true
  try {
    await apiUpdateFeedingRecord(recordId.value, updateData)
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
          await apiDeleteFeedingRecord(recordId.value)
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
