<template>
  <view class="daily-tips-card" v-if="tips && tips.length > 0">
    <!-- 卡片头部 -->
    <view class="tips-header">
      <view class="header-left">
        <view class="ai-icon-wrapper">
          <image
            src="/static/smart_toy.svg"
            mode="aspectFill"
            class="ai-icon"
          />
        </view>
        <view class="header-text">
          <text class="header-title">今日建议</text>
          <text class="header-subtitle">个性化智能推荐</text>
        </view>
      </view>
      <!-- 移除右侧文案区域 -->
    </view>

    <!-- 建议内容 - 横向滚动 -->
    <scroll-view scroll-x class="tips-scroll" show-scrollbar="false">
      <view class="tips-container">
        <view
          v-for="(tip, index) in displayTips"
          :key="index"
          class="tip-card clickable"
          @click="handleTipClick(tip)"
        >
          <view class="tip-header">
            <text class="tip-title">{{ tip.title }}</text>
            <!-- 添加点击指示器 -->
            <view class="click-indicator">
              <wd-icon name="arrow-right" size="14" color="#999" />
            </view>
          </view>
          <text class="tip-description">{{ tip.description }}</text>
        </view>
      </view>
    </scroll-view>

    <!-- 空状态 -->
    <view class="tips-empty" v-if="!tips || tips.length === 0">
      <view class="empty-icon">💡</view>
      <text class="empty-text">暂无今日建议</text>
      <text class="empty-subtext">AI正在为您准备个性化建议</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'

// 临时类型定义，避免导入问题
interface DailyTip {
  id: string
  title: string
  description: string
  type: string
  priority: 'high' | 'medium' | 'low'
}

// Props
interface Props {
  tips?: DailyTip[]
  maxDisplay?: number
}

const props = withDefaults(defineProps<Props>(), {
  maxDisplay: 3
})

// Emits
const emit = defineEmits<{
  tipClick: [tip: DailyTip]
}>()

// 显示的建议数量
const displayTips = computed(() => {
  if (!props.tips) return []
  return props.tips.slice(0, props.maxDisplay)
})


// 处理建议点击
const handleTipClick = (tip: DailyTip) => {
  emit('tipClick', tip)
}
</script>

<style lang="scss" scoped>
@import '@/styles/colors.scss';

.daily-tips-card {
  background: $color-bg-primary;
  border: 1rpx solid $color-border-primary;
  border-radius: $radius-lg;
  padding: $spacing-lg $spacing-md;
  margin-bottom: $spacing-2xl;
}

.tips-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: $spacing-lg;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12rpx;
}

.ai-icon-wrapper {
  width: 48rpx;
  height: 48rpx;
  border-radius: $radius-full;
  background: $color-primary-lighter;
  display: flex;
  align-items: center;
  justify-content: center;
}

.ai-icon {
  width: 44rpx;
  height: 44rpx;
}


.header-text {
  display: flex;
  flex-direction: column;
  gap: 2rpx;
}

.header-title {
  font-size: 28rpx;
  font-weight: $font-weight-semibold;
  color: $color-text-primary;
}

.header-subtitle {
  font-size: 22rpx;
  color: $color-text-secondary;
}

.tips-scroll {
  width: 100%;
  white-space: nowrap;
}

.tips-container {
  display: flex;
  gap: $spacing-md;
  padding: 0 $spacing-md;
}

.tip-card {
  min-width: 280rpx;
  max-width: 320rpx;
  background: $color-bg-secondary;
  border: 1rpx solid $color-border-primary;
  border-radius: $radius-md;
  padding: $spacing-lg $spacing-md;
  display: flex;
  flex-direction: column;
  gap: $spacing-md;
  transition: all $transition-base;
  position: relative;
  overflow: hidden;
  cursor: pointer;

  // 添加微妙的渐变效果
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 2rpx;
    background: linear-gradient(90deg, transparent, $color-primary-light, transparent);
    opacity: 0;
    transition: opacity 0.3s ease;
  }

  // 悬停效果
  &:hover {
    background: rgba($color-primary, 0.05);

    &::before {
      opacity: 0.6;
    }
  }

  // 按压效果
  &:active {
    transform: scale(0.98);
    box-shadow: inset 0 2rpx 8rpx rgba(0, 0, 0, 0.1);
  }
}

.tip-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8rpx;
}

.tip-title {
  font-size: 28rpx;
  color: $color-text-primary;
  font-weight: 350;
  text-align: left;
  line-height: 1.4;
  flex: 1;
}

// 点击指示器
.click-indicator {
  opacity: 0.6;
  transition: all 0.2s ease;
  margin-left: 8rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  height: 28rpx; // 与标题行高匹配，确保垂直居中

  .tip-card:hover & {
    opacity: 1;
    transform: translateX(2rpx);
  }
}

.tip-description {
  font-size: 24rpx;
  color: $color-text-secondary;
  line-height: 1.5;
  width: 100%;
  height: 72rpx;
  overflow: hidden;
  word-wrap: break-word;
  white-space: normal;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  text-overflow: ellipsis;
}

.tips-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60rpx 0;
  text-align: center;
}

.empty-icon {
  font-size: 48rpx;
  margin-bottom: $spacing-md;
  opacity: 0.6;
}

.empty-text {
  font-size: 26rpx;
  color: $color-text-secondary;
  margin-bottom: 8rpx;
}

.empty-subtext {
  font-size: 22rpx;
  color: $color-text-tertiary;
}

// 暗色模式适配
@media (prefers-color-scheme: dark) {
  .daily-tips-card {
    background: #1a1a1a;
    border-color: #333333;
  }

  .tip-card {
    background: #2a2a2a;
    border-color: #444444;

    &:hover {
      background: rgba($color-primary, 0.1);
    }
  }
}
</style>
