<template>
  <wd-popup
    v-model="visible"
    position="bottom"
    :closeable="true"
    round
    :safe-area-inset-bottom="true"
    @close="handleClose"
  >
    <view class="subscribe-guide">
      <!-- 头部图标 -->
      <view class="guide-header">
        <image v-if="template?.icon" :src="template.icon" class="guide-icon" mode="aspectFit" />
        <view v-else class="guide-icon-placeholder">🔔</view>
      </view>

      <!-- 标题和描述 -->
      <view class="guide-content">
        <text class="guide-title">{{ template?.title || '消息提醒' }}</text>
        <text class="guide-description">{{ description }}</text>

        <!-- 场景化提示 -->
        <view v-if="contextMessage" class="guide-context">
          <text class="context-message">{{ contextMessage }}</text>
        </view>

        <!-- 微信授权说明 -->
        <view class="wechat-auth-notice">
          <text class="notice-text">📱 点击下方按钮后,将跳转到微信官方授权页面</text>
        </view>
      </view>

      <!-- 操作按钮 -->
      <view class="guide-actions">
        <wd-button size="large" type="primary" class="btn-confirm" @click="handleConfirm">
          {{ confirmText }}
        </wd-button>
        <wd-button size="large" type="default" class="btn-cancel" @click="handleDismiss">
          {{ dismissText }}
        </wd-button>
      </view>

      <!-- 不再提示选项 -->
      <view class="guide-footer" @click="handleNeverShow">
        <text class="never-show-text">不再提示</text>
      </view>
    </view>
  </wd-popup>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import type { SubscribeMessageType, SubscribeMessageTemplate } from '@/types'
import {
  getTemplateConfig,
  recordGuideShown,
  dismissGuideForever,
  requestSubscribeMessage,
} from '@/store/subscribe'

interface Props {
  /** 消息类型 */
  type: SubscribeMessageType
  /** 是否显示 */
  modelValue: boolean
  /** 场景化提示文案 */
  contextMessage?: string
  /** 自定义描述 */
  customDescription?: string
  /** 确认按钮文案 */
  confirmText?: string
  /** 取消按钮文案 */
  dismissText?: string
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void
  (e: 'confirm', result: 'accept' | 'reject'): void
  (e: 'dismiss'): void
}

const props = withDefaults(defineProps<Props>(), {
  confirmText: '立即开启',
  dismissText: '暂不需要',
})

const emit = defineEmits<Emits>()

const visible = ref(false)
const template = ref<SubscribeMessageTemplate>()

// 监听 modelValue 变化
watch(
  () => props.modelValue,
  (newValue) => {
    if (newValue) {
      // ✅ 移除二次检查,因为调用方已经检查过 authStatus 了
      // 对于一次性订阅消息,每次都需要显示,不需要在组件内部再次检查

      // 加载模板配置
      template.value = getTemplateConfig(props.type)
      if (!template.value) {
        console.error(`[SubscribeGuide] 未找到模板配置: ${props.type}`)
        emit('update:modelValue', false)
        return
      }

      visible.value = true
      recordGuideShown(props.type)
    } else {
      visible.value = false
    }
  },
  { immediate: true }
)

// 描述文案
const description = computed(() => {
  return props.customDescription || template.value?.description || '开启消息提醒,不错过重要时刻'
})

/** 关闭弹窗 */
function handleClose() {
  visible.value = false
  emit('update:modelValue', false)
}

/** 确认开启 */
async function handleConfirm() {
  try {
    uni.showLoading({ title: '请求授权中...' })

    const results = await requestSubscribeMessage([props.type])
    const result = results.get(props.type)

    uni.hideLoading()

    if (result === 'accept') {
      uni.showToast({
        title: '开启成功',
        icon: 'success',
      })
      emit('confirm', 'accept')
    } else {
      uni.showToast({
        title: '您拒绝了授权',
        icon: 'none',
      })
      emit('confirm', 'reject')
    }

    handleClose()
  } catch (error: any) {
    uni.hideLoading()
    console.error('[SubscribeGuide] 授权失败:', error)
    uni.showToast({
      title: '授权失败',
      icon: 'none',
    })
    handleClose()
  }
}

/** 暂不需要 */
function handleDismiss() {
  emit('dismiss')
  handleClose()
}

/** 不再提示 */
function handleNeverShow() {
  uni.showModal({
    title: '确认操作',
    content: '确定不再显示该提示?您仍可在"设置"中手动开启提醒',
    success: (res) => {
      if (res.confirm) {
        dismissGuideForever(props.type)
        uni.showToast({
          title: '已关闭提示',
          icon: 'success',
        })
        handleClose()
      }
    },
  })
}
</script>

<style lang="scss" scoped>
.subscribe-guide {
  padding: 40rpx;
  background-color: #fff;
  border-radius: 32rpx 32rpx 0 0;
}

.guide-header {
  display: flex;
  justify-content: center;
  margin-bottom: 32rpx;

  .guide-icon {
    width: 120rpx;
    height: 120rpx;
  }

  .guide-icon-placeholder {
    width: 120rpx;
    height: 120rpx;
    font-size: 80rpx;
    line-height: 120rpx;
    text-align: center;
  }
}

.guide-content {
  text-align: center;
  margin-bottom: 32rpx;

  .guide-title {
    display: block;
    font-size: 36rpx;
    font-weight: 600;
    color: #1a1a1a;
    margin-bottom: 16rpx;
  }

  .guide-description {
    display: block;
    font-size: 28rpx;
    color: #666;
    line-height: 1.6;
  }

  .guide-context {
    margin-top: 24rpx;
    padding: 20rpx 24rpx;
    background-color: #fff7e6;
    border-radius: 12rpx;
    border-left: 4rpx solid #ffa940;

    .context-message {
      font-size: 26rpx;
      color: #d46b08;
      line-height: 1.5;
    }
  }

  .wechat-auth-notice {
    margin-top: 24rpx;
    padding: 16rpx 20rpx;
    background-color: #f0f9ff;
    border-radius: 12rpx;
    border: 1rpx solid #91d5ff;

    .notice-text {
      font-size: 24rpx;
      color: #0958d9;
      line-height: 1.5;
    }
  }
}

.guide-actions {
  display: flex;
  flex-direction: column;
  gap: 24rpx;
  margin-bottom: 24rpx;

  .btn-cancel,
  .btn-confirm {
    width: 100%;
    border-radius: 16rpx;
    font-size: 30rpx;
  }
}

.guide-footer {
  text-align: center;
  padding: 16rpx 0;

  .never-show-text {
    font-size: 26rpx;
    color: #999;
    text-decoration: underline;
  }
}
</style>
