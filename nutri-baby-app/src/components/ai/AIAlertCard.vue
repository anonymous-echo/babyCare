<template>
  <view class="ai-alert-card">
    <view class="alert-header">
      <text class="alert-title">⚠️ 健康关注事项</text>
      <view class="alert-count" v-if="alerts.length">
        <text class="count-text">{{ alerts.length }}</text>
      </view>
    </view>

    <view class="alert-list">
      <view
        class="alert-item"
        v-for="(alert, index) in displayedAlerts"
        :key="index"
        :class="`alert-${alert.level}`"
        @tap="handleAlertClick(alert)"
      >
        <view class="alert-icon">{{ getAlertIcon(alert.level) }}</view>
        <view class="alert-content">
          <view class="alert-main">
            <text class="alert-title-text">{{ alert.title }}</text>
            <text class="alert-level">{{ getAlertLevelText(alert.level) }}</text>
          </view>
          <text class="alert-description">{{ alert.description }}</text>
          <view class="alert-suggestion" v-if="alert.suggestion">
            <text class="suggestion-label">建议：</text>
            <text class="suggestion-text">{{ alert.suggestion }}</text>
          </view>
        </view>
        <view class="alert-time">
          <text class="time-text">{{ formatTime(alert.timestamp) }}</text>
        </view>
      </view>
    </view>

    <view class="alert-actions" v-if="alerts.length > maxDisplay">
      <nut-button
        type="primary"
        size="small"
        plain
        @tap="toggleShowAll"
      >
        {{ showAll ? '收起' : `查看全部 (${alerts.length})` }}
      </nut-button>
    </view>

    <view class="alert-empty" v-if="!alerts.length">
      <view class="empty-icon">✅</view>
      <text class="empty-text">暂无健康关注事项</text>
      <text class="empty-subtext">继续保持良好的育儿习惯</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import type { AIAlert } from '@/types/ai'

interface Props {
  alerts: AIAlert[]
  maxDisplay?: number
  showActions?: boolean
  compact?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  maxDisplay: 3,
  showActions: true,
  compact: false
})

const emit = defineEmits(['alertClick', 'showAll'])

const showAll = ref(false)

// 计算显示的警告
const displayedAlerts = computed(() => {
  if (showAll.value) {
    return props.alerts
  }
  return props.alerts.slice(0, props.maxDisplay)
})

// 获取警告图标
const getAlertIcon = (level: string): string => {
  const iconMap: Record<string, string> = {
    critical: '🚨',
    warning: '⚠️',
    info: 'ℹ️'
  }
  return iconMap[level] || '⚠️'
}

// 获取警告级别文本
const getAlertLevelText = (level: string): string => {
  const textMap: Record<string, string> = {
    critical: '严重',
    warning: '警告',
    info: '提示'
  }
  return textMap[level] || level
}

// 格式化时间
const formatTime = (timestamp: string): string => {
  try {
    const date = new Date(timestamp)
    const now = new Date()
    const diff = now.getTime() - date.getTime()

    // 小于1小时显示"刚刚"
    if (diff < 60 * 60 * 1000) {
      return '刚刚'
    }

    // 小于24小时显示小时数
    if (diff < 24 * 60 * 60 * 1000) {
      const hours = Math.floor(diff / (60 * 60 * 1000))
      return `${hours}小时前`
    }

    // 显示日期
    return date.toLocaleDateString('zh-CN', {
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  } catch (error) {
    return timestamp
  }
}

// 处理警告点击
const handleAlertClick = (alert: AIAlert) => {
  emit('alertClick', alert)
}

// 切换显示全部
const toggleShowAll = () => {
  showAll.value = !showAll.value
  emit('showAll', showAll.value)
}
</script>

<style lang="scss" scoped>
.ai-alert-card {
  background: #ffffff;
  border-radius: 16rpx;
  padding: 24rpx;
  margin: 16rpx 0;
  box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.04);

  .alert-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24rpx;
    padding-bottom: 16rpx;
    border-bottom: 1rpx solid #f0f0f0;

    .alert-title {
      font-size: 32rpx;
      font-weight: 600;
      color: #333333;
    }

    .alert-count {
      background: #ff4757;
      color: #ffffff;
      border-radius: 50%;
      width: 40rpx;
      height: 40rpx;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 24rpx;
      font-weight: 600;
    }
  }

  .alert-list {
    .alert-item {
      display: flex;
      align-items: flex-start;
      padding: 24rpx;
      margin-bottom: 16rpx;
      border-radius: 12rpx;
      border-left: 8rpx solid;
      transition: all 0.2s;

      &:active {
        transform: scale(0.98);
      }

      &.alert-critical {
        background: linear-gradient(135deg, #fff5f5, #ffecec);
        border-left-color: #ff4757;
      }

      &.alert-warning {
        background: linear-gradient(135deg, #fffaf0, #fff5e6);
        border-left-color: #ffa940;
      }

      &.alert-info {
        background: linear-gradient(135deg, #f0f9ff, #e6f7ff);
        border-left-color: #1890ff;
      }

      .alert-icon {
        font-size: 40rpx;
        margin-right: 20rpx;
        margin-top: 4rpx;
      }

      .alert-content {
        flex: 1;
        min-width: 0;

        .alert-main {
          display: flex;
          justify-content: space-between;
          align-items: center;
          margin-bottom: 12rpx;

          .alert-title-text {
            font-size: 28rpx;
            font-weight: 600;
            color: #333333;
            flex: 1;
            margin-right: 16rpx;
          }

          .alert-level {
            font-size: 24rpx;
            padding: 4rpx 12rpx;
            border-radius: 12rpx;
            font-weight: 500;
            color: #ffffff;
            white-space: nowrap;

            .alert-critical & {
              background: #ff4757;
            }

            .alert-warning & {
              background: #ffa940;
            }

            .alert-info & {
              background: #1890ff;
            }
          }
        }

        .alert-description {
          font-size: 26rpx;
          color: #666666;
          line-height: 1.6;
          margin-bottom: 12rpx;
        }

        .alert-suggestion {
          background: rgba(255, 255, 255, 0.8);
          padding: 16rpx;
          border-radius: 8rpx;
          margin-top: 12rpx;

          .suggestion-label {
            font-size: 24rpx;
            color: #1890ff;
            font-weight: 600;
            margin-right: 8rpx;
          }

          .suggestion-text {
            font-size: 24rpx;
            color: #333333;
            line-height: 1.5;
          }
        }
      }

      .alert-time {
        margin-left: 16rpx;

        .time-text {
          font-size: 22rpx;
          color: #999999;
          white-space: nowrap;
        }
      }
    }
  }

  .alert-actions {
    display: flex;
    justify-content: center;
    padding-top: 16rpx;
    border-top: 1rpx solid #f0f0f0;
  }

  .alert-empty {
    text-align: center;
    padding: 60rpx 0;

    .empty-icon {
      font-size: 80rpx;
      margin-bottom: 16rpx;
    }

    .empty-text {
      display: block;
      font-size: 28rpx;
      color: #52c41a;
      font-weight: 600;
      margin-bottom: 8rpx;
    }

    .empty-subtext {
      display: block;
      font-size: 24rpx;
      color: #999999;
    }
  }
}

// 紧凑模式
.compact-mode {
  padding: 20rpx;

  .alert-item {
    padding: 20rpx;
    margin-bottom: 12rpx;

    .alert-icon {
      font-size: 32rpx;
    }

    .alert-content {
      .alert-main {
        .alert-title-text {
          font-size: 26rpx;
        }

        .alert-level {
          font-size: 22rpx;
        }
      }

      .alert-description {
        font-size: 24rpx;
      }

      .alert-suggestion {
        padding: 12rpx;

        .suggestion-label,
        .suggestion-text {
          font-size: 22rpx;
        }
      }
    }

    .alert-time {
      .time-text {
        font-size: 20rpx;
      }
    }
  }
}

// 暗色模式适配
@media (prefers-color-scheme: dark) {
  .ai-alert-card {
    background: #1a1a1a;
    color: #ffffff;

    .alert-header {
      border-bottom-color: #333333;

      .alert-title {
        color: #ffffff;
      }

      .alert-count {
        background: #ff6b6b;
      }
    }

    .alert-list {
      .alert-item {
        &.alert-critical {
          background: linear-gradient(135deg, #2a1a1a, #331a1a);
        }

        &.alert-warning {
          background: linear-gradient(135deg, #2a251a, #332a1a);
        }

        &.alert-info {
          background: linear-gradient(135deg, #1a252a, #1a2a33);
        }

        .alert-content {
          .alert-main {
            .alert-title-text {
              color: #ffffff;
            }
          }

          .alert-description {
            color: #cccccc;
          }

          .alert-suggestion {
            background: rgba(0, 0, 0, 0.3);

            .suggestion-text {
              color: #ffffff;
            }
          }
        }

        .alert-time {
          .time-text {
            color: #999999;
          }
        }
      }
    }

    .alert-actions {
      border-top-color: #333333;
    }

    .alert-empty {
      .empty-text {
        color: #7dd3a2;
      }

      .empty-subtext {
        color: #cccccc;
      }
    }
  }
}
</style>

<style lang="scss">
// 动画效果
@keyframes alertPulse {
  0% {
    opacity: 1;
  }
  50% {
    opacity: 0.8;
  }
  100% {
    opacity: 1;
  }
}

.alert-critical {
  animation: alertPulse 2s infinite;
}

@keyframes alertSlideIn {
  from {
    opacity: 0;
    transform: translateX(-20rpx);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

.ai-alert-card {
  animation: alertSlideIn 0.3s ease-out;
}
</style>

<style lang="scss">
// 响应式布局
@media (max-width: 375px) {
  .ai-alert-card {
    .alert-list {
      .alert-item {
        padding: 20rpx;

        .alert-icon {
          font-size: 36rpx;
          margin-right: 16rpx;
        }

        .alert-content {
          .alert-main {
            flex-direction: column;
            align-items: flex-start;

            .alert-title-text {
              margin-right: 0;
              margin-bottom: 8rpx;
            }
          }
        }

        .alert-time {
          margin-left: 0;
          margin-top: 12rpx;
        }
      }
    }
  }
}
</style>