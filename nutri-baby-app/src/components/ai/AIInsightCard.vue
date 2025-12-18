<template>
  <view class="ai-insight-card">
    <view class="insight-header">
      <view class="header-left">
        <text class="insight-icon">{{ getInsightIcon(insight.type) }}</text>
        <text class="insight-title">{{ insight.title }}</text>
      </view>
      <view class="header-right">
        <view class="priority-badge" :class="`priority-${insight.priority}`">
          {{ getPriorityText(insight.priority) }}
        </view>
      </view>
    </view>

    <view class="insight-content">
      <text class="insight-description">{{ insight.description }}</text>

      <view class="insight-category" v-if="insight.category">
        <nut-tag type="primary" size="small">
          {{ insight.category }}
        </nut-tag>
      </view>
    </view>

    <view class="insight-actions" v-if="showActions">
      <nut-button
        type="primary"
        size="small"
        plain
        @tap="handleAction"
      >
        查看详情
      </nut-button>
    </view>
  </view>
</template>

<script setup lang="ts">
import type { AIInsight } from '@/types/ai'

interface Props {
  insight: AIInsight
  showActions?: boolean
  compact?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  showActions: false,
  compact: false
})

const emit = defineEmits(['action', 'detail'])

// 获取洞察图标
const getInsightIcon = (type: string): string => {
  const iconMap: Record<string, string> = {
    feeding: '🍼',
    sleep: '😴',
    growth: '📈',
    health: '❤️',
    behavior: '🧠',
    pattern: '🔍',
    trend: '📊',
    recommendation: '💡',
    warning: '⚠️',
    tip: '✨'
  }
  return iconMap[type] || '💡'
}

// 获取优先级文本
const getPriorityText = (priority: string): string => {
  const textMap: Record<string, string> = {
    high: '高优先级',
    medium: '中等优先级',
    low: '低优先级'
  }
  return textMap[priority] || priority
}

// 处理操作
const handleAction = () => {
  emit('action', props.insight)
  emit('detail', props.insight)
}
</script>

<style lang="scss" scoped>
.ai-insight-card {
  background: #ffffff;
  border-radius: 16rpx;
  padding: 24rpx;
  margin: 16rpx 0;
  box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.04);
  border-left: 8rpx solid #7dd3a2;

  .insight-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 16rpx;

    .header-left {
      display: flex;
      align-items: center;
      flex: 1;

      .insight-icon {
        font-size: 32rpx;
        margin-right: 12rpx;
      }

      .insight-title {
        font-size: 28rpx;
        font-weight: 600;
        color: #333333;
        line-height: 1.4;
      }
    }

    .header-right {
      .priority-badge {
        padding: 4rpx 12rpx;
        border-radius: 12rpx;
        font-size: 20rpx;
        font-weight: 500;
        color: #ffffff;

        &.priority-high {
          background: linear-gradient(135deg, #ff6b6b, #ff4757);
        }

        &.priority-medium {
          background: linear-gradient(135deg, #ffa940, #faad14);
        }

        &.priority-low {
          background: linear-gradient(135deg, #52c41a, #7dd3a2);
        }
      }
    }
  }

  .insight-content {
    margin-bottom: 16rpx;

    .insight-description {
      font-size: 26rpx;
      color: #666666;
      line-height: 1.6;
      margin-bottom: 12rpx;
    }

    .insight-category {
      margin-top: 12rpx;
    }
  }

  .insight-actions {
    display: flex;
    justify-content: flex-end;
    padding-top: 16rpx;
    border-top: 1rpx solid #f0f0f0;
  }
}

// 紧凑模式
.compact-mode {
  padding: 20rpx;
  margin: 12rpx 0;

  .insight-header {
    margin-bottom: 12rpx;

    .insight-title {
      font-size: 26rpx;
    }
  }

  .insight-content {
    margin-bottom: 12rpx;

    .insight-description {
      font-size: 24rpx;
    }
  }
}

// 暗色模式适配
@media (prefers-color-scheme: dark) {
  .ai-insight-card {
    background: #1a1a1a;
    color: #ffffff;

    .insight-header {
      .insight-title {
        color: #ffffff;
      }
    }

    .insight-content {
      .insight-description {
        color: #cccccc;
      }
    }

    .insight-actions {
      border-top-color: #333333;
    }
  }
}
</style>

<style lang="scss">
// 全局样式
.insight-category {
  .nut-tag {
    background: rgba(125, 211, 162, 0.1) !important;
    color: #7dd3a2 !important;
    border-color: rgba(125, 211, 162, 0.3) !important;
  }
}
</style>

<style lang="scss">
// 动画效果
@keyframes insightSlideIn {
  from {
    opacity: 0;
    transform: translateX(-20rpx);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

.ai-insight-card {
  animation: insightSlideIn 0.3s ease-out;
}
</style>