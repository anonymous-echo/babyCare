<template>
  <view class="ucharts-wrapper" :class="{ 'is-loading': isLoading }">
    <!-- 加载中指示器 -->
    <view v-if="isLoading" class="ucharts-loading">
      <view class="loading-spinner"></view>
      <text class="loading-text">{{ loadingText }}</text>
    </view>

    <!-- Canvas 容器 -->
    <canvas
      :id="canvasId"
      :canvas-id="canvasId"
      :class="['ucharts-canvas', `ucharts-canvas-${canvasId}`]"
      :style="{ width: width, height: height }"
      @click="handleCanvasClick"
    />

    <!-- 错误提示 -->
    <view v-if="error" class="ucharts-error">
      <text class="error-text">{{ error }}</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, nextTick, onBeforeUnmount } from 'vue'
import uCharts from '@qiun/ucharts'
import type { ChartType, ChartData, ChartOptions } from '@/composables/useUChart'

interface Props {
  /**
   * Canvas ID，必须唯一
   */
  canvasId: string

  /**
   * 图表类型
   */
  chartType: ChartType

  /**
   * 图表数据
   */
  chartData: ChartData

  /**
   * 图表配置选项
   */
  chartOptions?: ChartOptions

  /**
   * Canvas 宽度，默认 100%
   */
  width?: string

  /**
   * Canvas 高度，默认 400rpx
   */
  height?: string

  /**
   * Canvas 类型，H5 不支持 '2d'，留空即可
   */
  canvasType?: string

  /**
   * 加载中文本
   */
  loadingText?: string

  /**
   * 是否正在加载
   */
  isLoading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  width: '100%',
  height: '400rpx',
  canvasType: '',
  loadingText: '加载中...',
  isLoading: false
})

const emit = defineEmits<{
  /**
   * Canvas 完成绘制事件
   */
  complete: [detail: any]

  /**
   * Canvas 点击事件
   */
  click: [event: any]

  /**
   * 错误事件
   */
  error: [error: string]
}>()

// 图表实例
const chartInstance = ref<any>(null)

// 错误信息
const error = ref<string>('')

// 是否已初始化
const isInitialized = ref(false)

/**
 * 初始化图表 - H5 环境使用简化方式
 */
const initChart = async () => {
  try {
    await nextTick()

    console.log('[UChart] 开始初始化图表:', props.canvasId)

    // 获取系统信息判断平台
    const systemInfo = uni.getSystemInfoSync()
    const pixelRatio = systemInfo.pixelRatio || 1

    console.log('[UChart] 平台信息:', {
      platform: systemInfo.platform,
      uniPlatform: systemInfo.uniPlatform
    })

    // H5 环境使用传统 canvas 方式
    // #ifdef H5
    const canvas = document.getElementById(props.canvasId) as HTMLCanvasElement
    if (!canvas) {
      error.value = 'Canvas 元素未找到'
      emit('error', error.value)
      console.error('[UChart] Canvas 元素未找到:', props.canvasId)
      return
    }

    const ctx = canvas.getContext('2d')
    if (!ctx) {
      error.value = 'Canvas 上下文获取失败'
      emit('error', error.value)
      return
    }

    // 获取 canvas 尺寸
    const rect = canvas.getBoundingClientRect()
    const width = rect.width
    const height = rect.height

    console.log('[UChart] H5 Canvas 信息:', { width, height, pixelRatio })

    // 设置 canvas 尺寸
    canvas.width = width * pixelRatio
    canvas.height = height * pixelRatio
    canvas.style.width = width + 'px'
    canvas.style.height = height + 'px'
    ctx.scale(pixelRatio, pixelRatio)

    // 创建 uCharts 实例
    chartInstance.value = new uCharts({
      type: props.chartType,
      context: ctx,
      width: width,
      height: height,
      categories: props.chartData.categories || [],
      series: props.chartData.series || [],
      animation: true,
      pixelRatio: pixelRatio,
      ...props.chartOptions
    })

    isInitialized.value = true
    error.value = ''
    console.log('[UChart] 图表初始化成功 (H5)')
    emit('complete', { instance: chartInstance.value })
    // #endif

    // #ifndef H5
    // 小程序环境使用 SelectorQuery
    const selectQuery = uni.createSelectorQuery() as any
    selectQuery.select(`.ucharts-canvas-${props.canvasId}` as any)
      .fields({ node: true, size: true })
      .exec((res: any) => {
        console.log('[UChart] 小程序 Canvas 查询结果:', res)

        if (!res || !res[0]) {
          error.value = 'Canvas 节点查询返回为空'
          emit('error', error.value)
          console.error('[UChart] Canvas 节点查询返回为空', props.canvasId)
          return
        }

        if (!res[0].node) {
          error.value = `Canvas 节点不存在，选择器: .ucharts-canvas-${props.canvasId}`
          emit('error', error.value)
          console.error('[UChart] Canvas 节点获取失败，详情:', res[0])
          return
        }

        const canvas = res[0].node
        const ctx = canvas.getContext('2d')

        if (!ctx) {
          error.value = '获取 Canvas 上下文失败'
          emit('error', error.value)
          console.error('[UChart] Canvas 上下文获取失败')
          return
        }

        console.log('[UChart] 小程序 Canvas 信息:', {
          width: res[0].width,
          height: res[0].height,
          pixelRatio
        })

        // 设置 canvas 尺寸
        canvas.width = res[0].width * pixelRatio
        canvas.height = res[0].height * pixelRatio
        ctx.scale(pixelRatio, pixelRatio)

        // 创建 uCharts 实例
        chartInstance.value = new uCharts({
          type: props.chartType,
          context: ctx,
          width: res[0].width,
          height: res[0].height,
          categories: props.chartData.categories || [],
          series: props.chartData.series || [],
          animation: true,
          pixelRatio: pixelRatio,
          ...props.chartOptions
        } as any)

        isInitialized.value = true
        error.value = ''
        console.log('[UChart] 图表初始化成功 (小程序)')
        emit('complete', { instance: chartInstance.value })
      })
    // #endif
  } catch (err: any) {
    error.value = err?.message || '图表初始化失败'
    emit('error', error.value)
    console.error('[UChart] 初始化失败:', err)
  }
}

/**
 * 处理 Canvas 点击
 */
const handleCanvasClick = (event: any) => {
  emit('click', event)
  
  if (chartInstance.value && chartInstance.value.touchLegend) {
    const touch = event.touches?.[0] || event.changedTouches?.[0]
    if (touch) {
      chartInstance.value.touchLegend(touch)
    }
  }
}

/**
 * 重新绘制图表
 */
const redraw = async () => {
  if (!isInitialized.value) {
    await initChart()
    return
  }

  await nextTick()
  if (chartInstance.value && chartInstance.value.updateData) {
    try {
      chartInstance.value.updateData({
        categories: props.chartData.categories || [],
        series: props.chartData.series || [],
        ...props.chartOptions
      })
    } catch (err: any) {
      error.value = err?.message || '重绘失败'
      emit('error', error.value)
    }
  }
}

/**
 * 监听数据变化，自动更新图表
 */
watch(
  () => [props.chartData, props.chartOptions],
  async () => {
    if (isInitialized.value && props.chartData.series?.length > 0) {
      await redraw()
    }
  },
  { deep: true }
)

// 暴露方法给父组件
defineExpose({
  redraw,
  getChartInstance: () => chartInstance.value
})

onMounted(async () => {
  // 延迟初始化，确保 DOM 已渲染
  // 在微信小程序中，Canvas 需要更长时间才能初始化
  console.log('[UChart] 组件挂载，等待 Canvas 渲染...')
  setTimeout(() => {
    if (props.chartData.series?.length > 0) {
      console.log('[UChart] 开始初始化图表，Canvas ID:', props.canvasId)
      initChart()
    } else {
      console.log('[UChart] 图表数据为空，跳过初始化')
    }
  }, 500)
})

onBeforeUnmount(() => {
  // 清理图表实例
  if (chartInstance.value) {
    chartInstance.value = null
  }
})
</script>

<style scoped lang="scss">
.ucharts-wrapper {
  position: relative;
  width: 100%;
  background: white;
  border-radius: 12rpx;
  overflow: hidden;

  &.is-loading {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 400rpx;
  }
}

.ucharts-canvas {
  display: block;
  width: 100%;
}

.ucharts-loading {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.9);
  z-index: 10;
}

.loading-spinner {
  width: 40rpx;
  height: 40rpx;
  border: 4rpx solid #f0f0f0;
  border-top-color: #7dd3a2;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.loading-text {
  margin-top: 20rpx;
  font-size: 24rpx;
  color: #999;
}

.ucharts-error {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 240, 245, 0.95);
  z-index: 10;
  border: 2rpx solid #ffebee;
}

.error-text {
  font-size: 24rpx;
  color: #d32f2f;
  text-align: center;
  padding: 20rpx;
}
</style>
