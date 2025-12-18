<template>
  <view class="ai-chart-container">
    <view class="chart-header" v-if="title || subtitle">
      <text class="chart-title" v-if="title">{{ title }}</text>
      <text class="chart-subtitle" v-if="subtitle">{{ subtitle }}</text>
    </view>

    <view class="chart-content">
      <canvas
        :canvas-id="chartId"
        :id="chartId"
        class="ai-chart-canvas"
        @touchstart="touchStart"
        @touchmove="touchMove"
        @touchend="touchEnd"
      />
    </view>

    <view class="chart-legend" v-if="showLegend && series.length">
      <view
        class="legend-item"
        v-for="(item, index) in series"
        :key="index"
        @tap="toggleSeries(index)"
      >
        <view
          class="legend-color"
          :style="{ backgroundColor: item.color || defaultColors[index] }"
        />
        <text class="legend-text">{{ item.name }}</text>
      </view>
    </view>

    <view class="chart-actions" v-if="showActions">
      <nut-button
        type="primary"
        size="small"
        @tap="refreshChart"
      >
        刷新
      </nut-button>
      <nut-button
        size="small"
        @tap="exportChart"
      >
        导出
      </nut-button>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useUChart } from '@/composables/useUChart'
import type { AISeries, AIChartData } from '@/types/ai'

interface Props {
  chartId: string
  data: AIChartData
  type?: 'line' | 'column' | 'radar' | 'pie'
  title?: string
  subtitle?: string
  showLegend?: boolean
  showActions?: boolean
  height?: number
  width?: number
  animation?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  type: 'line',
  showLegend: true,
  showActions: false,
  height: 300,
  width: 750,
  animation: true
})

const emit = defineEmits(['refresh', 'export', 'seriesToggle'])

// 默认颜色配置
const defaultColors = [
  '#7dd3a2', // 绿色
  '#52c41a', // 深绿色
  '#ffa940', // 橙色
  '#1890ff', // 蓝色
  '#ff6b6b', // 红色
  '#722ed1', // 紫色
  '#13c2c2', // 青色
  '#faad14'  // 黄色
]

// 图表实例
const { initChart, updateChart, destroyChart, getChartData } = useUChart()

// 状态
const chartInstance = ref<any>(null)
const series = ref<AISeries[]>([])
const categories = ref<string[]>([])

// 图表配置
const chartConfig = ref({
  type: props.type,
  canvasId: props.chartId,
  canvas2d: false,
  background: '#ffffff',
  pixelRatio: 1,
  categories: categories.value,
  series: series.value,
  animation: props.animation,
  width: props.width,
  height: props.height,
  dataLabel: true,
  dataPointShape: true,
  legend: {
    show: props.showLegend
  },
  extra: {
    line: {
      type: 'straight',
      width: 2
    },
    column: {
      width: 20
    },
    radar: {
      gridType: 'polygon',
      gridColor: '#cccccc',
      gridCount: 5,
      axisLine: true,
      axisLineColor: '#cccccc'
    },
    pie: {
      border: true,
      borderWidth: 2,
      borderColor: '#ffffff'
    }
  }
})

// 初始化图表
const initChartInstance = async () => {
  await nextTick()

  try {
    chartInstance.value = await initChart(chartConfig.value)

    if (chartInstance.value) {
      // 添加触摸事件支持
      chartInstance.value.addEventListener('touchstart', handleTouchStart)
      chartInstance.value.addEventListener('touchmove', handleTouchMove)
      chartInstance.value.addEventListener('touchend', handleTouchEnd)
    }
  } catch (error) {
    console.error('初始化AI图表失败:', error)
  }
}

// 更新图表数据
const updateChartData = () => {
  if (!props.data || !chartInstance.value) return

  categories.value = props.data.categories || []
  series.value = props.data.series.map((s, index) => ({
    ...s,
    color: s.color || defaultColors[index % defaultColors.length]
  }))

  const newConfig = {
    ...chartConfig.value,
    categories: categories.value,
    series: series.value
  }

  try {
    updateChart(chartInstance.value, newConfig)
  } catch (error) {
    console.error('更新AI图表数据失败:', error)
  }
}

// 触摸事件处理
const touchStart = (e: TouchEvent) => {
  if (chartInstance.value) {
    chartInstance.value.showToolTip && chartInstance.value.showToolTip(e)
  }
}

const touchMove = (e: TouchEvent) => {
  if (chartInstance.value) {
    chartInstance.value.scrollTooltip && chartInstance.value.scrollTooltip(e)
  }
}

const touchEnd = (e: TouchEvent) => {
  // 触摸结束处理
}

const handleTouchStart = (e: Event) => {
  touchStart(e as TouchEvent)
}

const handleTouchMove = (e: Event) => {
  touchMove(e as TouchEvent)
}

const handleTouchEnd = (e: Event) => {
  touchEnd(e as TouchEvent)
}

// 切换系列显示
const toggleSeries = (index: number) => {
  if (!chartInstance.value) return

  try {
    const currentSeries = series.value[index]
    if (currentSeries) {
      // 切换显示状态
      currentSeries.show = !currentSeries.show
      updateChartData()
      emit('seriesToggle', index, currentSeries.show)
    }
  } catch (error) {
    console.error('切换系列失败:', error)
  }
}

// 刷新图表
const refreshChart = () => {
  updateChartData()
  emit('refresh')
}

// 导出图表
const exportChart = () => {
  try {
    if (chartInstance.value) {
      // 这里可以实现图表导出功能
      console.log('导出图表数据:', getChartData())
      emit('export', getChartData())
    }
  } catch (error) {
    console.error('导出图表失败:', error)
  }
}

// 获取图表数据
const getChartData = () => {
  return {
    categories: categories.value,
    series: series.value,
    title: props.title,
    subtitle: props.subtitle,
    type: props.type
  }
}

// 监听数据变化
watch(() => props.data, updateChartData, { deep: true })

// 生命周期
onMounted(() => {
  initChartInstance()
})

onUnmounted(() => {
  if (chartInstance.value) {
    // 移除事件监听器
    chartInstance.value.removeEventListener('touchstart', handleTouchStart)
    chartInstance.value.removeEventListener('touchmove', handleTouchMove)
    chartInstance.value.removeEventListener('touchend', handleTouchEnd)

    // 销毁图表实例
    destroyChart(chartInstance.value)
  }
})

// 暴露方法
defineExpose({
  refreshChart,
  exportChart,
  getChartData,
  updateChartData
})
</script>

<style lang="scss" scoped>
.ai-chart-container {
  background: #ffffff;
  border-radius: 12rpx;
  padding: 24rpx;
  margin: 16rpx 0;
  box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.04);

  .chart-header {
    margin-bottom: 24rpx;
    text-align: center;

    .chart-title {
      display: block;
      font-size: 32rpx;
      font-weight: 600;
      color: #333333;
      margin-bottom: 8rpx;
    }

    .chart-subtitle {
      display: block;
      font-size: 24rpx;
      color: #999999;
    }
  }

  .chart-content {
    position: relative;
    margin-bottom: 24rpx;
  }

  .ai-chart-canvas {
    width: 100%;
    height: 300rpx;
  }

  .chart-legend {
    display: flex;
    flex-wrap: wrap;
    justify-content: center;
    gap: 24rpx;
    margin-bottom: 24rpx;

    .legend-item {
      display: flex;
      align-items: center;
      gap: 8rpx;
      padding: 8rpx 16rpx;
      background: #f5f5f5;
      border-radius: 8rpx;
      cursor: pointer;
      transition: all 0.2s;

      &:active {
        transform: scale(0.95);
      }

      .legend-color {
        width: 16rpx;
        height: 16rpx;
        border-radius: 50%;
      }

      .legend-text {
        font-size: 24rpx;
        color: #666666;
      }
    }
  }

  .chart-actions {
    display: flex;
    justify-content: center;
    gap: 16rpx;
  }
}

// 暗色模式适配
@media (prefers-color-scheme: dark) {
  .ai-chart-container {
    background: #1a1a1a;
    color: #ffffff;

    .chart-header {
      .chart-title {
        color: #ffffff;
      }

      .chart-subtitle {
        color: #cccccc;
      }
    }

    .chart-legend {
      .legend-item {
        background: #2a2a2a;

        .legend-text {
          color: #cccccc;
        }
      }
    }
  }
}
</style>

<style>
/* 全局图表样式 */
.qiun-columns {
  background: transparent !important;
}

.qiun-line {
  background: transparent !important;
}

.qiun-radar {
  background: transparent !important;
}

.qiun-pie {
  background: transparent !important;
}
</style>