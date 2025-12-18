<template>
  <view class="ai-score-card">
    <view class="score-header">
      <view class="header-left">
        <text class="score-icon">{{ getScoreIcon() }}</text>
        <text class="score-title">{{ title }}</text>
      </view>
      <view class="header-right">
        <text class="score-label">综合评分</text>
      </view>
    </view>

    <view class="score-content">
      <view class="score-display">
        <view class="score-circle">
          <circle-progress
            :percent="score"
            :stroke-width="12"
            :stroke-color="getScoreColor()"
            :trail-color="'#f0f0f0'"
            :width="200"
            :radius="90"
            show-percent
            percent-text-size="48"
            percent-text-color="#333333"
          />
        </view>

        <view class="score-level">
          <text class="level-text">{{ getScoreLevel() }}</text>
          <text class="level-desc">{{ getScoreDescription() }}</text>
        </view>
      </view>

      <view class="score-details" v-if="details && details.length">
        <view
          class="detail-item"
          v-for="(detail, index) in details"
          :key="index"
        >
          <view class="detail-icon">{{ getDetailIcon(detail.type) }}</view>
          <view class="detail-content">
            <text class="detail-name">{{ detail.name }}</text>
            <view class="detail-score">
              <text class="score-value">{{ detail.score }}</text>
              <text class="score-max">/100</text>
            </view>
          </view>
          <view class="detail-progress">
            <view
              class="progress-bar"
              :style="{ width: `${detail.score}%`, backgroundColor: getDetailColor(detail.score) }"
            />
          </view>
        </view>
      </view>
    </view>

    <view class="score-actions" v-if="showActions">
      <nut-button
        type="primary"
        size="small"
        @tap="handleRefresh"
      >
        重新分析
      </nut-button>
      <nut-button
        size="small"
        @tap="handleDetail"
      >
        查看详情
      </nut-button>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import CircleProgress from '@/components/common/CircleProgress.vue'

interface DetailItem {
  type: string
  name: string
  score: number
}

interface Props {
  title: string
  score: number
  details?: DetailItem[]
  showActions?: boolean
  size?: 'small' | 'medium' | 'large'
}

const props = withDefaults(defineProps<Props>(), {
  details: () => [],
  showActions: false,
  size: 'medium'
})

const emit = defineEmits(['refresh', 'detail'])

// 获取评分图标
const getScoreIcon = (): string => {
  if (props.score >= 90) return '🏆'
  if (props.score >= 80) return '⭐'
  if (props.score >= 70) return '👍'
  if (props.score >= 60) return '😐'
  return '⚠️'
}

// 获取评分颜色
const getScoreColor = (): string => {
  if (props.score >= 90) return '#52c41a'  // 优秀 - 绿色
  if (props.score >= 80) return '#7dd3a2'  // 良好 - 浅绿
  if (props.score >= 70) return '#ffa940'  // 一般 - 橙色
  if (props.score >= 60) return '#faad14'  // 及格 - 黄色
  return '#ff4d4f'  // 不及格 - 红色
}

// 获取评分等级
const getScoreLevel = (): string => {
  if (props.score >= 90) return '优秀'
  if (props.score >= 80) return '良好'
  if (props.score >= 70) return '一般'
  if (props.score >= 60) return '及格'
  return '需改善'
}

// 获取评分描述
const getScoreDescription = (): string => {
  if (props.score >= 90) return '表现非常出色'
  if (props.score >= 80) return '整体表现良好'
  if (props.score >= 70) return '基本符合预期'
  if (props.score >= 60) return '需要关注改进'
  return '需要重点关注'
}

// 获取详情图标
const getDetailIcon = (type: string): string => {
  const iconMap: Record<string, string> = {
    feeding: '🍼',
    sleep: '😴',
    growth: '📈',
    health: '❤️',
    behavior: '🧠',
    regularity: '⏰',
    adequacy: '✅',
    timeliness: '⚡',
    diversity: '🌈'
  }
  return iconMap[type] || '📊'
}

// 获取详情颜色
const getDetailColor = (score: number): string => {
  if (score >= 90) return '#52c41a'
  if (score >= 80) return '#7dd3a2'
  if (score >= 70) return '#ffa940'
  if (score >= 60) return '#faad14'
  return '#ff4d4f'
}

// 处理刷新
const handleRefresh = () => {
  emit('refresh')
}

// 处理详情
const handleDetail = () => {
  emit('detail')
}

// 计算属性
const circleSize = computed(() => {
  switch (props.size) {
    case 'small': return 150
    case 'large': return 250
    default: return 200
  }
})

const progressSize = computed(() => {
  switch (props.size) {
    case 'small': return 6
    case 'large': return 16
    default: return 12
  }
})
</script>

<style lang="scss" scoped>
.ai-score-card {
  background: #ffffff;
  border-radius: 16rpx;
  padding: 24rpx;
  margin: 16rpx 0;
  box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.04);

  .score-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24rpx;

    .header-left {
      display: flex;
      align-items: center;

      .score-icon {
        font-size: 36rpx;
        margin-right: 12rpx;
      }

      .score-title {
        font-size: 32rpx;
        font-weight: 600;
        color: #333333;
      }
    }

    .header-right {
      .score-label {
        font-size: 24rpx;
        color: #999999;
      }
    }
  }

  .score-content {
    .score-display {
      display: flex;
      flex-direction: column;
      align-items: center;
      margin-bottom: 32rpx;

      .score-circle {
        margin-bottom: 24rpx;
      }

      .score-level {
        text-align: center;

        .level-text {
          display: block;
          font-size: 32rpx;
          font-weight: 600;
          color: #333333;
          margin-bottom: 8rpx;
        }

        .level-desc {
          display: block;
          font-size: 26rpx;
          color: #666666;
        }
      }
    }

    .score-details {
      .detail-item {
        display: flex;
        align-items: center;
        margin-bottom: 20rpx;
        padding: 16rpx;
        background: #f8f9fa;
        border-radius: 12rpx;

        .detail-icon {
          font-size: 32rpx;
          margin-right: 16rpx;
          width: 40rpx;
          text-align: center;
        }

        .detail-content {
          flex: 1;
          display: flex;
          justify-content: space-between;
          align-items: center;

          .detail-name {
            font-size: 26rpx;
            color: #333333;
          }

          .detail-score {
            display: flex;
            align-items: baseline;

            .score-value {
              font-size: 28rpx;
              font-weight: 600;
              color: #333333;
            }

            .score-max {
              font-size: 22rpx;
              color: #999999;
              margin-left: 4rpx;
            }
          }
        }

        .detail-progress {
          width: 100rpx;
          height: 8rpx;
          background: #e8e8e8;
          border-radius: 4rpx;
          margin-left: 16rpx;
          overflow: hidden;

          .progress-bar {
            height: 100%;
            border-radius: 4rpx;
            transition: width 0.3s ease;
          }
        }
      }
    }
  }

  .score-actions {
    display: flex;
    justify-content: center;
    gap: 16rpx;
    padding-top: 24rpx;
    border-top: 1rpx solid #f0f0f0;
  }
}

// 不同尺寸
.score-small {
  padding: 20rpx;

  .score-header {
    margin-bottom: 20rpx;

    .score-title {
      font-size: 28rpx;
    }
  }

  .score-content {
    .score-display {
      margin-bottom: 24rpx;

      .score-level {
        .level-text {
          font-size: 28rpx;
        }

        .level-desc {
          font-size: 24rpx;
        }
      }
    }

    .score-details {
      .detail-item {
        padding: 12rpx;
        margin-bottom: 16rpx;

        .detail-icon {
          font-size: 28rpx;
        }

        .detail-content {
          .detail-name {
            font-size: 24rpx;
          }

          .detail-score {
            .score-value {
              font-size: 26rpx;
            }
          }
        }

        .detail-progress {
          width: 80rpx;
          height: 6rpx;
        }
      }
    }
  }
}

.score-large {
  padding: 32rpx;

  .score-header {
    margin-bottom: 32rpx;

    .score-title {
      font-size: 36rpx;
    }
  }

  .score-content {
    .score-display {
      margin-bottom: 40rpx;

      .score-level {
        .level-text {
          font-size: 36rpx;
        }

        .level-desc {
          font-size: 28rpx;
        }
      }
    }

    .score-details {
      .detail-item {
        padding: 20rpx;
        margin-bottom: 24rpx;

        .detail-icon {
          font-size: 36rpx;
        }

        .detail-content {
          .detail-name {
            font-size: 28rpx;
          }

          .detail-score {
            .score-value {
              font-size: 32rpx;
            }
          }
        }

        .detail-progress {
          width: 120rpx;
          height: 10rpx;
        }
      }
    }
  }
}

// 暗色模式适配
@media (prefers-color-scheme: dark) {
  .ai-score-card {
    background: #1a1a1a;
    color: #ffffff;

    .score-header {
      .score-title {
        color: #ffffff;
      }

      .score-label {
        color: #cccccc;
      }
    }

    .score-content {
      .score-display {
        .score-level {
          .level-text {
            color: #ffffff;
          }

          .level-desc {
            color: #cccccc;
          }
        }
      }

      .score-details {
        .detail-item {
          background: #2a2a2a;

          .detail-content {
            .detail-name {
              color: #ffffff;
            }

            .detail-score {
              .score-value {
                color: #ffffff;
              }
            }
          }

          .detail-progress {
            background: #333333;
          }
        }
      }
    }

    .score-actions {
      border-top-color: #333333;
    }
  }
}
</style>

<style lang="scss">
// 圆环进度条样式
.circle-progress {
  display: inline-block;
  position: relative;

  .progress-text {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    font-weight: 600;
    color: #333333;
  }
}

// 动画效果
@keyframes scoreSlideIn {
  from {
    opacity: 0;
    transform: translateY(20rpx);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.ai-score-card {
  animation: scoreSlideIn 0.4s ease-out;
}

@keyframes progressFill {
  from {
    width: 0;
  }
  to {
    width: var(--progress-width);
  }
}

.detail-progress {
  .progress-bar {
    animation: progressFill 0.6s ease-out;
  }
}
</style>

<style lang="scss">
// 响应式布局
@media (max-width: 375px) {
  .ai-score-card {
    padding: 20rpx;

    .score-content {
      .score-details {
        .detail-item {
          flex-direction: column;
          align-items: flex-start;

          .detail-progress {
            width: 100%;
            margin-left: 0;
            margin-top: 12rpx;
          }
        }
      }
    }
  }
}
</style>

<!-- 圆环进度条组件 -->
<script lang="ts">
// 这里应该引入实际的圆环进度条组件
// 由于NutUI没有直接的圆环进度条，这里用占位符
// 实际项目中可以使用第三方组件或自定义实现
interface CircleProgressProps {
  percent: number
  strokeWidth?: number
  strokeColor?: string
  trailColor?: string
  width?: number
  radius?: number
  showPercent?: boolean
  percentTextSize?: number
  percentTextColor?: string
}
</script>