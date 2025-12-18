import { ref, computed, reactive, watch, nextTick } from 'vue'
import type { Ref } from 'vue'

export type ChartType = 'line' | 'column' | 'pie' | 'area' | 'radar' | 'gauge' | 'ring'

/**
 * 图表数据结构
 */
export interface ChartData {
  categories?: string[]
  series: Array<{
    name: string
    data: number[] | Array<{ name: string; value: number }>
  }>
}

/**
 * 图表配置选项
 */
export interface ChartOptions {
  color?: string[]
  padding?: [number, number, number, number]
  fontSize?: number
  legend?: {
    show?: boolean
    position?: 'top' | 'bottom' | 'left' | 'right'
  }
  xAxis?: {
    disableGrid?: boolean
    disabled?: boolean
    boundaryGap?: boolean
  }
  yAxis?: {
    gridType?: 'solid' | 'dash'
    dashLength?: number
  }
  extra?: {
    column?: {
      type?: 'group' | 'stack'
      width?: number
    }
    line?: {
      type?: 'straight' | 'curve'
      width?: number
    }
    pie?: {
      offsetAngle?: number
      labelWidth?: number
    }
    area?: {
      type?: 'straight' | 'curve'
    }
  }
  enableScroll?: boolean
  enableMarkLine?: boolean
  enableAutoScale?: boolean
  tooltips?: boolean
}

/**
 * uCharts 组合式函数
 * @param chartType - 图表类型
 * @param defaultOptions - 默认配置选项
 * @returns 图表管理对象
 */
export function useUChart(
  chartType: ChartType,
  defaultOptions: ChartOptions = {}
) {
  // 图表数据
  const chartData: Ref<ChartData> = ref({
    categories: [],
    series: []
  })

  // 图表配置
  const chartOpts: Ref<ChartOptions> = ref({
    color: ['#7dd3a2', '#52c41a', '#ff7f50', '#4a90e2', '#faad14'],
    padding: [15, 15, 0, 15],
    fontSize: 12,
    legend: {
      show: true,
      position: 'bottom'
    },
    xAxis: {
      disableGrid: true,
      boundaryGap: true
    },
    yAxis: {
      gridType: 'dash',
      dashLength: 2
    },
    enableScroll: false,
    enableMarkLine: false,
    enableAutoScale: true,
    tooltips: true,
    ...defaultOptions
  })

  // 图表加载状态
  const isLoading = ref(false)

  // 图表实例引用
  const chartInstance: Ref<any> = ref(null)

  /**
   * 更新图表数据
   * @param data - 新的图表数据
   * @param animate - 是否启用动画
   */
  const updateChartData = async (data: ChartData, animate = true) => {
    isLoading.value = true
    try {
      if (animate) {
        // 使用 nextTick 确保动画流畅
        await nextTick()
      }
      chartData.value = { ...data }
    } finally {
      isLoading.value = false
    }
  }

  /**
   * 更新图表配置
   * @param options - 新的配置选项
   */
  const updateChartOpts = (options: Partial<ChartOptions>) => {
    chartOpts.value = {
      ...chartOpts.value,
      ...options
    }
  }

  /**
   * 设置图表实例
   * @param instance - uCharts 实例
   */
  const setChartInstance = (instance: any) => {
    chartInstance.value = instance
  }

  /**
   * 重新绘制图表
   */
  const redrawChart = async () => {
    if (chartInstance.value) {
      await nextTick()
      // 触发图表重新绘制（通过修改 chartData 触发）
      chartData.value = { ...chartData.value }
    }
  }

  /**
   * 清空图表数据
   */
  const clearChartData = () => {
    chartData.value = {
      categories: [],
      series: []
    }
  }

  /**
   * 柱状图特定方法：获取柱子宽度
   */
  const getColumnWidth = (totalWidth: number, categoryCount: number) => {
    if (categoryCount === 0) return 0
    return Math.floor(totalWidth / categoryCount / 2)
  }

  /**
   * 折线图特定方法：转换为光滑曲线
   */
  const enableCurveMode = () => {
    updateChartOpts({
      extra: {
        ...chartOpts.value.extra,
        line: { type: 'curve', width: 2 }
      }
    })
  }

  /**
   * 堆叠柱状图配置
   */
  const enableStackColumn = () => {
    updateChartOpts({
      extra: {
        ...chartOpts.value.extra,
        column: { type: 'stack' }
      }
    })
  }

  /**
   * 分组柱状图配置
   */
  const enableGroupColumn = () => {
    updateChartOpts({
      extra: {
        ...chartOpts.value.extra,
        column: { type: 'group' }
      }
    })
  }

  /**
   * 启用图表滚动
   */
  const enableScroll = (enable = true) => {
    updateChartOpts({ enableScroll: enable })
  }

  /**
   * 启用标记线
   */
  const enableMarkLine = (enable = true) => {
    updateChartOpts({ enableMarkLine: enable })
  }

  return {
    // 数据和配置
    chartData,
    chartOpts,
    isLoading,
    chartInstance,

    // 基础方法
    updateChartData,
    updateChartOpts,
    setChartInstance,
    redrawChart,
    clearChartData,

    // 特定图表类型的方法
    getColumnWidth,
    enableCurveMode,
    enableStackColumn,
    enableGroupColumn,
    enableScroll,
    enableMarkLine
  }
}

/**
 * 预设配置：柱状图
 */
export const columnChartPreset = (): ChartOptions => ({
  color: ['#7dd3a2', '#52c41a'],
  padding: [15, 15, 0, 15],
  xAxis: { disableGrid: true },
  yAxis: { gridType: 'dash', dashLength: 2 },
  extra: {
    column: {
      type: 'group',
      width: 30
    }
  }
})

/**
 * 预设配置：折线图
 */
export const lineChartPreset = (): ChartOptions => ({
  color: ['#7dd3a2', '#ff7f50', '#4a90e2'],
  padding: [15, 15, 0, 15],
  legend: { show: true, position: 'bottom' },
  extra: {
    line: {
      type: 'curve',
      width: 2
    }
  }
})

/**
 * 预设配置：饼图
 */
export const pieChartPreset = (): ChartOptions => ({
  color: ['#7dd3a2', '#52c41a', '#ff7f50', '#4a90e2', '#faad14'],
  padding: [5, 5, 5, 5],
  legend: { show: true, position: 'bottom' },
  extra: {
    pie: {
      labelWidth: 15
    }
  }
})

/**
 * 预设配置：区域图
 */
export const areaChartPreset = (): ChartOptions => ({
  color: ['#7dd3a2', '#ff7f50'],
  padding: [15, 15, 0, 15],
  legend: { show: true, position: 'bottom' },
  extra: {
    area: {
      type: 'curve'
    }
  }
})
