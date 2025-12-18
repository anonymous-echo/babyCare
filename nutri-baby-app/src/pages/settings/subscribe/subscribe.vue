<template>
  <view class="settings-page">
    <!-- 顶部说明 -->
    <view class="page-header">
      <text class="header-title">消息提醒设置</text>
      <text class="header-desc">开启提醒,不错过宝宝的重要时刻</text>
    </view>

    <!-- 提醒列表 -->
    <view class="reminder-list">
      <view
        v-for="template in templates"
        :key="template.type"
        class="reminder-item"
        @click="handleToggleReminder(template)"
      >
        <view class="item-left">
          <image v-if="template.icon" :src="template.icon" class="item-icon" mode="aspectFit" />
          <view v-else class="item-icon-placeholder">{{ getIconEmoji(template.type) }}</view>

          <view class="item-info">
            <text class="item-title">{{ template.title }}</text>
            <text class="item-desc">{{ template.description }}</text>

            <!-- 授权状态提示 -->
            <text v-if="getAuthStatus(template.type) === 'ban'" class="item-status error">
              已拒绝,请在微信设置中手动开启
            </text>
            <text v-else-if="getAuthStatus(template.type) === 'reject'" class="item-status warning">
              暂未授权
            </text>
          </view>
        </view>

        <view class="item-right">
          <wd-switch
            :model="getReminderEnabled(template.type)"
            :disabled="getAuthStatus(template.type) === 'ban'"
            @update:model="(val: boolean) => handleSwitchChange(template, val)"
          />
        </view>
      </view>
    </view>

    <!-- 高级设置(可选) -->
    <view v-if="hasEnabledReminders" class="advanced-settings">
      <text class="section-title">高级设置</text>

      <!-- 疫苗提醒提前天数 -->
      <view v-if="vaccineReminderEnabled" class="setting-item">
        <text class="setting-label">疫苗提醒提前天数</text>
        <wd-input-number
          v-model="vaccineAdvanceDays"
          :min="1"
          :max="7"
          @change="handleVaccineAdvanceDaysChange"
        />
      </view>

      <!-- 喂养提醒间隔 -->
      <view v-if="feedingReminderEnabled" class="setting-item">
        <text class="setting-label">喂养提醒间隔(分钟)</text>
        <wd-input-number
          v-model="feedingIntervalMinutes"
          :min="60"
          :max="360"
          :step="30"
          @change="handleFeedingIntervalChange"
        />
      </view>
    </view>

    <!-- 底部说明 -->
    <view class="footer-note">
      <text class="note-text">
        💡 提示:订阅消息由微信官方管理,您可以在微信的"设置 > 通知 > 订阅消息"中管理所有订阅
      </text>
    </view>

    <!-- 清除授权记录(仅开发调试用) -->
    <!-- <view class="debug-section">
      <wd-button type="warning" size="small" @click="handleClearRecords">
        清除授权记录(调试)
      </wd-button>
    </view> -->
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import type { SubscribeMessageType, SubscribeMessageTemplate } from '@/types'
import {
  getAllTemplateConfigs,
  getAuthStatus as _getAuthStatus,
  getReminderConfig,
  updateReminderConfig,
  enableReminder,
  disableReminder,
  getAllReminderConfigs,
  hasEnabledReminders as _hasEnabledReminders,
} from '@/store/subscribe'
import { StorageKeys, removeStorage } from '@/utils/storage'

const templates = ref<SubscribeMessageTemplate[]>([])
const vaccineAdvanceDays = ref(3)
const feedingIntervalMinutes = ref(180)

onMounted(() => {
  loadTemplates()
  loadAdvancedSettings()
})

/** 加载模板配置 */
function loadTemplates() {
  templates.value = getAllTemplateConfigs().sort((a, b) => b.priority - a.priority)
}

/** 加载高级设置 */
function loadAdvancedSettings() {
  const vaccineConfig = getReminderConfig('vaccine_reminder')
  if (vaccineConfig?.advanceDays) {
    vaccineAdvanceDays.value = vaccineConfig.advanceDays
  }

  const breastConfig = getReminderConfig('breast_feeding_reminder')
  if (breastConfig?.intervalMinutes) {
    feedingIntervalMinutes.value = breastConfig.intervalMinutes
  }
}

/** 获取授权状态 */
function getAuthStatus(type: SubscribeMessageType) {
  return _getAuthStatus(type)
}

/** 获取提醒启用状态 */
function getReminderEnabled(type: SubscribeMessageType) {
  const config = getReminderConfig(type)
  return config?.enabled || false
}

/** 是否有已启用的提醒 */
const hasEnabledReminders = computed(() => {
  return _hasEnabledReminders()
})

/** 疫苗提醒是否启用 */
const vaccineReminderEnabled = computed(() => {
  return getReminderEnabled('vaccine_reminder')
})

/** 喂养提醒是否启用 */
const feedingReminderEnabled = computed(() => {
  return (
    getReminderEnabled('breast_feeding_reminder') || getReminderEnabled('bottle_feeding_reminder')
  )
})

/** 获取图标emoji */
function getIconEmoji(type: SubscribeMessageType): string {
  const emojiMap: Record<SubscribeMessageType, string> = {
    vaccine_reminder: '💉',
    breast_feeding_reminder: '🤱',
    bottle_feeding_reminder: '🍼',
    pump_reminder: '🔔',
    feeding_duration_alert: '⏰',
  }
  return emojiMap[type] || '🔔'
}

/** 点击提醒项 */
function handleToggleReminder(template: SubscribeMessageTemplate) {
  const enabled = getReminderEnabled(template.type)
  handleSwitchChange(template, !enabled)
}

/** 开关切换 */
async function handleSwitchChange(template: SubscribeMessageTemplate, value: boolean) {
  if (value) {
    // 启用提醒
    const success = await enableReminder(template.type)
    if (!success) {
      // 恢复开关状态
      const config = getReminderConfig(template.type)
      if (config) {
        config.enabled = false
      }
    }
    // 刷新模板列表以更新 UI
    loadTemplates()
  } else {
    // 禁用提醒
    uni.showModal({
      title: '确认关闭',
      content: `确定关闭"${template.title}"提醒吗?`,
      success: (res) => {
        if (res.confirm) {
          disableReminder(template.type)
          uni.showToast({
            title: '已关闭提醒',
            icon: 'success',
          })
          // 刷新模板列表以更新 UI
          loadTemplates()
        }
      },
    })
  }
}

/** 疫苗提前天数变更 */
function handleVaccineAdvanceDaysChange(value: number) {
  updateReminderConfig('vaccine_reminder', {
    advanceDays: value,
  })
  uni.showToast({
    title: `已设置为提前${value}天`,
    icon: 'success',
  })
}

/** 喂养间隔变更 */
function handleFeedingIntervalChange(value: number) {
  const types: SubscribeMessageType[] = ['breast_feeding_reminder', 'bottle_feeding_reminder']
  types.forEach((type) => {
    if (getReminderEnabled(type)) {
      updateReminderConfig(type, {
        intervalMinutes: value,
      })
    }
  })
  uni.showToast({
    title: `已设置为${value}分钟`,
    icon: 'success',
  })
}

/** 清除授权记录(调试用) */
function handleClearRecords() {
  uni.showModal({
    title: '确认操作',
    content: '确定清除所有授权记录吗?(仅用于调试)',
    success: (res) => {
      if (res.confirm) {
        removeStorage(StorageKeys.SUBSCRIBE_AUTH_RECORDS)
        removeStorage(StorageKeys.SUBSCRIBE_GUIDE_RECORDS)
        removeStorage(StorageKeys.SUBSCRIBE_REMINDER_CONFIGS)
        uni.showToast({
          title: '已清除',
          icon: 'success',
        })
        setTimeout(() => {
          uni.navigateBack()
        }, 1500)
      }
    },
  })
}
</script>

<style lang="scss" scoped>
.settings-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding-bottom: 40rpx;
}

.page-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 60rpx 40rpx 40rpx;
  color: #fff;

  .header-title {
    display: block;
    font-size: 40rpx;
    font-weight: 600;
    margin-bottom: 12rpx;
  }

  .header-desc {
    display: block;
    font-size: 26rpx;
    opacity: 0.9;
  }
}

.reminder-list {
  margin-top: 24rpx;
  background-color: #fff;
  border-radius: 16rpx;
  overflow: hidden;
}

.reminder-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 32rpx 40rpx;
  border-bottom: 1rpx solid #f0f0f0;

  &:last-child {
    border-bottom: none;
  }

  .item-left {
    display: flex;
    align-items: center;
    flex: 1;
    margin-right: 24rpx;

    .item-icon,
    .item-icon-placeholder {
      width: 80rpx;
      height: 80rpx;
      margin-right: 24rpx;
      flex-shrink: 0;
    }

    .item-icon-placeholder {
      font-size: 56rpx;
      line-height: 80rpx;
      text-align: center;
    }

    .item-info {
      flex: 1;

      .item-title {
        display: block;
        font-size: 30rpx;
        font-weight: 500;
        color: #1a1a1a;
        margin-bottom: 8rpx;
      }

      .item-desc {
        display: block;
        font-size: 24rpx;
        color: #999;
        line-height: 1.4;
      }

      .item-status {
        display: block;
        font-size: 22rpx;
        margin-top: 8rpx;

        &.error {
          color: #ff4d4f;
        }

        &.warning {
          color: #faad14;
        }
      }
    }
  }

  .item-right {
    flex-shrink: 0;
  }
}

.advanced-settings {
  margin-top: 24rpx;
  background-color: #fff;
  border-radius: 16rpx;
  padding: 32rpx 40rpx;

  .section-title {
    display: block;
    font-size: 28rpx;
    font-weight: 600;
    color: #1a1a1a;
    margin-bottom: 24rpx;
  }

  .setting-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 24rpx 0;
    border-bottom: 1rpx solid #f0f0f0;

    &:last-child {
      border-bottom: none;
    }

    .setting-label {
      font-size: 28rpx;
      color: #333;
    }
  }
}

.footer-note {
  margin-top: 24rpx;
  padding: 32rpx 40rpx;

  .note-text {
    display: block;
    font-size: 24rpx;
    color: #999;
    line-height: 1.6;
  }
}

.debug-section {
  margin-top: 24rpx;
  padding: 0 40rpx;
}
</style>
