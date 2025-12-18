<template>
  <view class="statistics-page">
    <!-- 时间范围选择 -->
    <view class="time-range">
      <wd-tabs v-model="timeRange">
        <wd-tab title="本周" name="week" />
        <wd-tab title="本月" name="month" />
      </wd-tabs>
    </view>

    <!-- 未登录提示 -->
    <view v-if="!isLoggedIn" class="guest-tip">
      <image src="/static/monitoring.svg" class="tip-icon" />
      <text class="tip-text">登录后查看数据</text>
    </view>

    <!-- 喂养统计 -->
    <view class="stat-section">
      <view class="section-header">
        <image class="icon-img" src="/static/breastfeeding.svg" mode="aspectFit" />
        <text class="title">喂养统计</text>
      </view>

      <view class="stat-cards">
        <view class="stat-card">
          <view class="card-label">奶瓶奶量</view>
          <view class="card-value">{{ feedingStats.totalMilk }}ml</view>
        </view>
        <view class="stat-card">
          <view class="card-label">喂养次数</view>
          <view class="card-value">{{ feedingStats.count }}次</view>
        </view>
        <view class="stat-card">
          <view class="card-label">日均奶量</view>
          <view class="card-value">{{ feedingStats.avgMilk }}ml</view>
        </view>
      </view>

      <!-- 每日奶量柱状图 -->
      <view class="daily-chart">
        <view class="chart-title">每日奶瓶奶量趋势</view>
        <canvas 
          type="2d"
          canvas-id="feedingChart" 
          id="feedingChart" 
          class="chart-canvas"
          @touchstart.stop="touchFeeding"
          @touchmove.stop.prevent="moveFeeding"
          @touchend.stop="touchEndFeeding"
        ></canvas>
      </view>
    </view>

    <!-- 睡眠统计 -->
    <view class="stat-section">
      <view class="section-header">
        <image class="icon-img" src="/static/moon_stars.svg" mode="aspectFit" />
        <text class="title">睡眠统计</text>
      </view>

      <view class="stat-cards">
        <view class="stat-card">
          <view class="card-label">总时长</view>
          <view class="card-value">{{ sleepStats.totalDurationFormatted }}</view>
        </view>
        <view class="stat-card">
          <view class="card-label">睡眠次数</view>
          <view class="card-value">{{ sleepStats.count }}次</view>
        </view>
        <view class="stat-card">
          <view class="card-label">日均时长</view>
          <view class="card-value">{{ sleepStats.avgDurationFormatted }}</view>
        </view>
      </view>

      <!-- 睡眠质量分析 -->
      <view class="sleep-quality">
        <view class="quality-title">睡眠质量分析</view>
        <view class="quality-content">
          <view class="quality-item">
            <text class="quality-label">最长单次睡眠:</text>
            <text class="quality-value">{{ sleepStats.longestSleep }}分</text>
          </view>
          <view class="quality-item">
            <text class="quality-label">平均单次时长:</text>
            <text class="quality-value">{{ sleepStats.avgSingleSleep }}分</text>
          </view>
          <view class="quality-item">
            <text class="quality-label">夜间睡眠:</text>
            <text class="quality-value">{{ sleepStats.nightSleepCount }}次 ({{ formatDurationToTimeString(Math.round(sleepStats.nightSleepHours * 60)) }})</text>
          </view>
          <view class="quality-item">
            <text class="quality-label">小睡:</text>
            <text class="quality-value">{{ sleepStats.napCount }}次 ({{ formatDurationToTimeString(Math.round(sleepStats.napHours * 60)) }})</text>
          </view>
          <view v-if="sleepStats.recommendation" class="quality-recommendation">
            <image src="/static/lightbulb_yellow.svg" class="recommendation-icon" />
            <text class="recommendation-text">{{ sleepStats.recommendation }}</text>
          </view>
        </view>
      </view>
    </view>

    <!-- 排泄统计 -->
    <view class="stat-section">
      <view class="section-header">
        <image class="icon-img" src="/static/baby_changing_station.svg" mode="aspectFit" />
        <text class="title">排泄统计</text>
      </view>

      <view class="stat-cards">
        <view class="stat-card">
          <view class="card-label">换尿布</view>
          <view class="card-value">{{ diaperStats.total }}次</view>
        </view>
        <view class="stat-card">
          <view class="card-label">小便</view>
          <view class="card-value">{{ diaperStats.wet }}次</view>
        </view>
        <view class="stat-card">
          <view class="card-label">大便</view>
          <view class="card-value">{{ diaperStats.dirty }}次</view>
        </view>
      </view>
    </view>

    <!-- 成长统计 -->
    <view v-if="growthStats.hasData" class="stat-section">
      <view class="section-header">
        <image class="icon-img" src="/static/monitoring.svg" mode="aspectFit" />
        <text class="title">成长统计</text>
      </view>

      <!-- 最新数据 -->
      <view class="stat-cards">
        <view v-if="growthStats.latestHeight" class="stat-card">
          <view class="card-label">最新身高</view>
          <view class="card-value">{{ growthStats.latestHeight }}cm</view>
        </view>
        <view v-if="growthStats.latestWeight" class="stat-card">
          <view class="card-label">最新体重</view>
          <view class="card-value">{{ growthStats.latestWeight }}g</view>
        </view>
        <view v-if="growthStats.latestHead" class="stat-card">
          <view class="card-label">最新头围</view>
          <view class="card-value">{{ growthStats.latestHead }}cm</view>
        </view>
      </view>

      <!-- 成长曲线 -->
      <view class="growth-charts">
        <!-- 身高曲线 -->
        <view v-if="growthStats.heightData.length > 0" class="chart-container">
          <view class="chart-title">身高趋势 (cm)</view>
          <canvas 
            type="2d"
            canvas-id="heightChart" 
            id="heightChart" 
            class="chart-canvas"
            @touchstart.stop="touchHeight"
            @touchmove.stop.prevent="moveHeight"
            @touchend.stop="touchEndHeight"
          ></canvas>
        </view>

        <!-- 体重曲线 -->
        <view v-if="growthStats.weightData.length > 0" class="chart-container">
          <view class="chart-title">体重趋势 (g)</view>
          <canvas 
            type="2d"
            canvas-id="weightChart" 
            id="weightChart" 
            class="chart-canvas"
            @touchstart.stop="touchWeight"
            @touchmove.stop.prevent="moveWeight"
            @touchend.stop="touchEndWeight"
          ></canvas>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, watch, nextTick } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { isLoggedIn } from '@/store/user'
import { currentBaby } from '@/store/baby'
import { formatDate } from '@/utils/date'
import uCharts from '@qiun/ucharts'
import { useUChart, columnChartPreset, lineChartPreset } from '@/composables/useUChart'

// 直接调用 API 层
import * as feedingApi from '@/api/feeding'
import * as sleepApi from '@/api/sleep'
import * as diaperApi from '@/api/diaper'
import * as growthApi from '@/api/growth'
import * as statisticsApi from '@/api/statistics'

// 图表实例
let feedingChartInstance: any = null
let heightChartInstance: any = null
let weightChartInstance: any = null

// 时间范围
const timeRange = ref<string>('week')

// 格式化睡眠时长为 X时Y分
const formatDurationToTimeString = (minutes: number): string => {
  if (minutes <= 0) return '0分'

  const hours = Math.floor(minutes / 60)
  const remainingMinutes = minutes % 60

  if (hours === 0) {
    return `${remainingMinutes}分`
  } else if (remainingMinutes === 0) {
    return `${hours}时`
  } else {
    return `${hours}时${remainingMinutes}分`
  }
}

// 初始化图表
const {
  chartData: feedingChartData,
  chartOpts: feedingChartOpts,
  updateChartData: updateFeedingChart
} = useUChart('column', columnChartPreset())

const {
  chartData: heightChartData,
  chartOpts: heightChartOpts,
  updateChartData: updateHeightChart
} = useUChart('line', lineChartPreset())

const {
  chartData: weightChartData,
  chartOpts: weightChartOpts,
  updateChartData: updateWeightChart
} = useUChart('line', lineChartPreset())

// 获取时间范围 - 改为过去7天/30天
const getTimeRange = () => {
  const now = Date.now()
  const days = timeRange.value === 'week' ? 7 : 30
  // 计算开始时间：今天往前推N天的00:00:00
  const startDate = new Date(now - (days - 1) * 24 * 60 * 60 * 1000)
  startDate.setHours(0, 0, 0, 0)
  return { start: startDate.getTime(), end: now }
}

// 统计数据(从按日统计 API 获取)
const dailyStats = ref<statisticsApi.DailyStatsResponse>({})

// 加载统计数据
const loadRecords = async () => {
  if (!currentBaby.value) return

  const babyId = currentBaby.value.babyId
  const { start, end } = getTimeRange()

  try {
    const dailyStatsData = await statisticsApi.apiFetchDailyStats({ babyId, startDate: start, endDate: end })
    dailyStats.value = dailyStatsData.data || {}
  } catch (error) {
    console.error('加载统计数据失败:', error)
    uni.showToast({
      title: '加载数据失败',
      icon: 'none'
    })
  }
}

// 监听时间范围变化,重新加载数据
watch(timeRange, async () => {
  await loadRecords()
  await nextTick()
  drawCharts()
})

// 喂养统计
const feedingStats = computed(() => {
  if (!currentBaby.value || !dailyStats.value.feeding) {
    return {
      totalMilk: 0,
      count: 0,
      avgMilk: 0,
      dailyData: [],
      maxDaily: 0,
    }
  }

  let totalMilk = 0
  let totalCount = 0
  const dailyMap = new Map<string, number>()

  // 使用按日统计数据
  dailyStats.value.feeding?.forEach(item => {
    // 只统计奶瓶喂养的奶量，母乳喂养不计入
    if (item.feedingType === 'bottle') {
      totalMilk += item.totalAmount
      totalCount += item.totalCount
      
      // 按日期统计（转换日期格式）
      const date = new Date(item.date)
      const dateStr = formatDate(date.getTime(), 'MM-DD')
      dailyMap.set(dateStr, item.totalAmount)
    }
  })

  // 生成每日数据
  const days = timeRange.value === 'week' ? 7 : 30
  const dailyData = []
  let maxDaily = 0

  for (let i = days - 1; i >= 0; i--) {
    const date = new Date(Date.now() - i * 24 * 60 * 60 * 1000)
    const dateStr = formatDate(date.getTime(), 'MM-DD')
    const amount = Math.round(dailyMap.get(dateStr) || 0)

    dailyData.push({
      label: i === 0 ? '今' : formatDate(date.getTime(), 'DD'),
      amount,
    })

    if (amount > maxDaily) maxDaily = amount
  }

  return {
    totalMilk: Math.round(totalMilk),
    count: totalCount,
    avgMilk: days > 0 ? Math.round(totalMilk / days) : 0,
    dailyData,
    maxDaily,
  }
})

// 睡眠统计
const sleepStats = computed(() => {
  if (!currentBaby.value || !dailyStats.value.sleep) {
    return {
      totalHours: 0,
      count: 0,
      avgHours: 0,
      longestSleep: 0,
      avgSingleSleep: 0,
      nightSleepCount: 0,
      nightSleepHours: 0,
      napCount: 0,
      napHours: 0,
      recommendation: '',
      totalDurationFormatted: '0分',
      avgDurationFormatted: '0分'
    }
  }

  // 使用按日统计数据计算总时长和次数
  let totalSeconds = 0
  let totalCount = 0
  
  dailyStats.value.sleep?.forEach(item => {
    totalSeconds += item.totalDuration
    totalCount += item.totalCount
  })
  
  const totalMinutes = Math.round(totalSeconds / 60)
  const days = timeRange.value === 'week' ? 7 : 30

  // 简化统计（由于按日统计无法获取单次睡眠详情，使用估算）
  const avgSingleSleep = totalCount > 0 ? Math.round(totalMinutes / totalCount) : 0
  
  // 计算宝宝月龄
  const birthDate = new Date(currentBaby.value.birthDate)
  const now = new Date()
  const monthsOld = (now.getFullYear() - birthDate.getFullYear()) * 12 +
                    (now.getMonth() - birthDate.getMonth())

  // 生成建议
  let recommendation = ''
  const dailyHours = totalMinutes / days / 60

  // 根据月龄判断睡眠是否充足
  if (monthsOld < 3) {
    // 0-3个月: 14-17小时
    if (dailyHours < 14) {
      recommendation = '建议增加睡眠时间,新生儿需要14-17小时睡眠'
    } else if (dailyHours > 17) {
      recommendation = '睡眠时间较长,如有异常请咨询医生'
    } else {
      recommendation = '睡眠时间正常,继续保持'
    }
  } else if (monthsOld < 12) {
    // 3-12个月: 12-16小时
    if (dailyHours < 12) {
      recommendation = '建议增加睡眠时间,婴儿需要12-16小时睡眠'
    } else if (dailyHours > 16) {
      recommendation = '睡眠时间较长,注意观察宝宝状态'
    } else {
      recommendation = '睡眠时间正常,继续保持'
    }
  } else {
    // 12个月以上: 11-14小时
    if (dailyHours < 11) {
      recommendation = '建议增加睡眠时间,幼儿需要11-14小时睡眠'
    } else if (dailyHours > 14) {
      recommendation = '睡眠时间较长,可适当增加活动'
    } else {
      recommendation = '睡眠时间正常,继续保持'
    }
  }

  return {
    totalHours: Math.round(totalMinutes / 60 * 10) / 10,
    count: totalCount,
    avgHours: Math.round((totalMinutes / days / 60) * 10) / 10,
    longestSleep: 0, // 按日统计无法获取单次最长睡眠
    avgSingleSleep,
    nightSleepCount: 0, // 按日统计无法区分夜间睡眠
    nightSleepHours: 0,
    napCount: 0,
    napHours: 0,
    recommendation,
    // 添加格式化后的时长字段（X时Y分）
    totalDurationFormatted: formatDurationToTimeString(totalMinutes),
    avgDurationFormatted: formatDurationToTimeString(Math.round(totalMinutes / days))
  }
})

// 排泄统计
const diaperStats = computed(() => {
  if (!currentBaby.value || !dailyStats.value.diaper) {
    return { total: 0, wet: 0, dirty: 0 }
  }

  let wet = 0
  let dirty = 0
  let total = 0

  // 使用按日统计数据
  dailyStats.value.diaper?.forEach(item => {
    total += item.totalCount
    if (item.diaperType === 'pee') {
      wet += item.totalCount
    } else if (item.diaperType === 'poop') {
      dirty += item.totalCount
    } else if (item.diaperType === 'both') {
      wet += item.totalCount
      dirty += item.totalCount
    }
  })

  return {
    total,
    wet,
    dirty,
  }
})

// 成长统计
const growthStats = computed(() => {
  if (!currentBaby.value || !dailyStats.value.growth) {
    return {
      hasData: false,
      latestHeight: 0,
      latestWeight: 0,
      latestHead: 0,
      dates: [],
      heightData: [],
      weightData: [],
      headData: [],
      heightMin: 0,
      heightMax: 0,
      weightMin: 0,
      weightMax: 0
    }
  }

  if (!dailyStats.value.growth || dailyStats.value.growth.length === 0) {
    return {
      hasData: false,
      latestHeight: 0,
      latestWeight: 0,
      latestHead: 0,
      dates: [],
      heightData: [],
      weightData: [],
      headData: [],
      heightMin: 0,
      heightMax: 0,
      weightMin: 0,
      weightMax: 0
    }
  }

  // 使用按日统计数据
  const growthData = dailyStats.value.growth
  const dates: string[] = []
  const heightData: number[] = []
  const weightData: number[] = []
  const headData: number[] = []

  // 最新数据（取最后一条记录）
  const latestRecord = growthData[growthData.length - 1]

  growthData.forEach(record => {
    const date = new Date(record.date)
    dates.push(`${date.getMonth() + 1}/${date.getDate()}`)

    if (record.latestHeight) heightData.push(record.latestHeight)
    if (record.latestWeight) weightData.push(record.latestWeight)
    if (record.latestHeadCircumference) headData.push(record.latestHeadCircumference)
  })

  // 计算最大最小值
  const heightMin = heightData.length > 0 ? Math.min(...heightData) : 0
  const heightMax = heightData.length > 0 ? Math.max(...heightData) : 0
  const weightMin = weightData.length > 0 ? Math.min(...weightData) : 0
  const weightMax = weightData.length > 0 ? Math.max(...weightData) : 0

  return {
    hasData: true,
    latestHeight: latestRecord?.latestHeight || 0,
    latestWeight: latestRecord?.latestWeight || 0,
    latestHead: latestRecord?.latestHeadCircumference || 0,
    dates,
    heightData,
    weightData,
    headData,
    heightMin: Math.floor(heightMin - 2),
    heightMax: Math.ceil(heightMax + 2),
    weightMin: Math.floor(weightMin - 0.5),
    weightMax: Math.ceil(weightMax + 0.5)
  }
})

// 获取 Canvas 上下文
const getCanvasContext = (canvasId: string, callback: (ctx: any, width: number, height: number) => void) => {
  console.log('[Statistics] 开始获取 Canvas 上下文:', canvasId)
  
  const query = uni.createSelectorQuery()
  const selector = query.select(`#${canvasId}`) as any
  selector
    .fields({ node: true, size: true }, null as any)
    .exec((res: any) => {
      console.log('[Statistics] Canvas 查询结果:', canvasId, res)
      
      if (!res || !res[0]) {
        console.error('[Statistics] Canvas 节点未找到:', canvasId)
        return
      }
      
      if (!res[0].node) {
        console.error('[Statistics] Canvas node 属性不存在:', canvasId, res[0])
        return
      }
      
      const canvas = res[0].node
      const ctx = canvas.getContext('2d')
      const dpr = uni.getSystemInfoSync().pixelRatio || 1
      
      console.log('[Statistics] Canvas 信息:', {
        canvasId,
        width: res[0].width,
        height: res[0].height,
        dpr
      })
      
      canvas.width = res[0].width * dpr
      canvas.height = res[0].height * dpr
      ctx.scale(dpr, dpr)
      
      callback(ctx, res[0].width, res[0].height)
    })
}

// 绘制喂养柱状图
const drawFeedingChart = () => {
  console.log('[Statistics] 开始绘制喂养图表')
  console.log('[Statistics] 喂养数据:', feedingStats.value.dailyData)
  
  getCanvasContext('feedingChart', (ctx, width, height) => {
    console.log('[Statistics] 创建喂养图表实例')
    
    const dataLength = feedingStats.value.dailyData.length
    const itemCount = timeRange.value === 'week' ? 7 : 10 // 本周显示7天，本月显示10天
    const enableScroll = dataLength > itemCount // 数据超过单屏数量时启用滚动
    
    const chartData = {
      $this: {},
      type: 'column',
      context: ctx,
      width: width,
      height: height,
      background: '#ffffff',
      categories: feedingStats.value.dailyData.map(d => d.label),
      series: [{
        name: '奶量(ml)',
        data: feedingStats.value.dailyData.map(d => d.amount)
      }],
      animation: true,
      color: ['#7dd3a2'],
      padding: [15, 15, 10, 15] as [number, number, number, number],
      dataLabel: true,
      enableScroll: enableScroll,
      legend: {
        show: false
      },
      xAxis: {
        disableGrid: true,
        itemCount: itemCount,
        scrollShow: true,
        scrollAlign: 'right',
        boundaryGap: 'center',
        fontSize: 11,
        axisLineHeight: 15
      },
      yAxis: {
        gridType: 'dash',
        dashLength: 2,
        data: [{
          min: 0,
          max: null
        }]
      },
      extra: {
        column: {
          type: 'group',
          width: enableScroll ? 15 : 20
        }
      }
    }
    
    console.log('[Statistics] 图表配置:', {
      dataLength,
      itemCount,
      enableScroll,
      chartData
    })
    
    try {
      feedingChartInstance = new uCharts(chartData, () => {
        console.log('[Statistics] 喂养图表绘制完成')
      })
    } catch (error) {
      console.error('[Statistics] 喂养图表创建失败:', error)
    }
  })
}

// 绘制身高折线图
const drawHeightChart = () => {
  getCanvasContext('heightChart', (ctx, width, height) => {
    const dataLength = growthStats.value.heightData.length
    const itemCount = 6 // 成长曲线显示6个数据点
    const enableScroll = dataLength > itemCount
    
    heightChartInstance = new uCharts({
      $this: {},
      type: 'line',
      context: ctx,
      width: width,
      height: height,
      background: '#ffffff',
      categories: growthStats.value.dates,
      series: [{
        name: '身高(cm)',
        data: growthStats.value.heightData
      }],
      animation: true,
      color: ['#7dd3a2'],
      padding: [15, 20, 0, 15] as [number, number, number, number],
      enableScroll: enableScroll,
      legend: {
        show: false
      },
      xAxis: {
        disableGrid: false,
        itemCount: itemCount,
        scrollShow: true,
        boundaryGap: 'center'
      },
      yAxis: {
        gridType: 'dash',
        dashLength: 2,
        min: growthStats.value.heightMin,
        max: growthStats.value.heightMax
      },
      extra: {
        line: {
          type: 'curve',
          width: 2
        }
      }
    }, () => {})
  })
}

// 绘制体重折线图
const drawWeightChart = () => {
  getCanvasContext('weightChart', (ctx, width, height) => {
    const dataLength = growthStats.value.weightData.length
    const itemCount = 6 // 成长曲线显示6个数据点
    const enableScroll = dataLength > itemCount
    
    weightChartInstance = new uCharts({
      $this: {},
      type: 'line',
      context: ctx,
      width: width,
      height: height,
      background: '#ffffff',
      categories: growthStats.value.dates,
      series: [{
        name: '体重(kg)',
        data: growthStats.value.weightData
      }],
      animation: true,
      color: ['#52c41a'],
      padding: [15, 20, 0, 15] as [number, number, number, number],
      enableScroll: enableScroll,
      legend: {
        show: false
      },
      xAxis: {
        disableGrid: false,
        itemCount: itemCount,
        scrollShow: true,
        boundaryGap: 'center'
      },
      yAxis: {
        gridType: 'dash',
        dashLength: 2,
        min: growthStats.value.weightMin,
        max: growthStats.value.weightMax
      },
      extra: {
        line: {
          type: 'curve',
          width: 2
        }
      }
    }, () => {})
  })
}

// 清理图表实例
const clearCharts = () => {
  console.log('[Statistics] 清理旧图表实例')
  
  if (feedingChartInstance) {
    try {
      feedingChartInstance.dispose?.()
    } catch (e) {
      console.warn('[Statistics] 清理喂养图表失败:', e)
    }
    feedingChartInstance = null
  }
  
  if (heightChartInstance) {
    try {
      heightChartInstance.dispose?.()
    } catch (e) {
      console.warn('[Statistics] 清理身高图表失败:', e)
    }
    heightChartInstance = null
  }
  
  if (weightChartInstance) {
    try {
      weightChartInstance.dispose?.()
    } catch (e) {
      console.warn('[Statistics] 清理体重图表失败:', e)
    }
    weightChartInstance = null
  }
}

// 绘制所有图表
const drawCharts = async () => {
  console.log('[Statistics] 准备绘制图表')
  console.log('[Statistics] 喂养数据长度:', feedingStats.value.dailyData.length)
  console.log('[Statistics] 身高数据长度:', growthStats.value.heightData.length)
  console.log('[Statistics] 体重数据长度:', growthStats.value.weightData.length)
  
  // 清理旧图表
  clearCharts()
  
  await nextTick()
  
  // 延迟绘制，确保 DOM 已渲染
  setTimeout(() => {
    console.log('[Statistics] 开始延迟绘制')
    
    if (feedingStats.value.dailyData.length > 0) {
      console.log('[Statistics] 绘制喂养图表')
      drawFeedingChart()
    } else {
      console.log('[Statistics] 跳过喂养图表（无数据）')
    }
    
    if (growthStats.value.heightData.length > 0) {
      console.log('[Statistics] 绘制身高图表')
      drawHeightChart()
    } else {
      console.log('[Statistics] 跳过身高图表（无数据）')
    }
    
    if (growthStats.value.weightData.length > 0) {
      console.log('[Statistics] 绘制体重图表')
      drawWeightChart()
    } else {
      console.log('[Statistics] 跳过体重图表（无数据）')
    }
  }, 500)
}

// 触摸事件处理
const touchFeeding = (e: any) => {
  console.log('[Statistics] 喂养图表 touchstart', e)
  if (feedingChartInstance) {
    // 开始滚动
    if (feedingChartInstance.scrollStart) {
      feedingChartInstance.scrollStart(e)
    }
    // 显示提示
    if (feedingChartInstance.showToolTip) {
      feedingChartInstance.showToolTip(e)
    }
  }
}

const moveFeeding = (e: any) => {
  if (feedingChartInstance) {
    // 滚动图表
    if (feedingChartInstance.scroll) {
      feedingChartInstance.scroll(e)
    }
    // 更新提示位置
    if (feedingChartInstance.showToolTip) {
      feedingChartInstance.showToolTip(e)
    }
  }
}

const touchEndFeeding = (e: any) => {
  if (feedingChartInstance) {
    // 结束滚动
    if (feedingChartInstance.scrollEnd) {
      feedingChartInstance.scrollEnd(e)
    }
  }
}

const touchHeight = (e: any) => {
  console.log('[Statistics] 身高图表 touchstart', e)
  if (heightChartInstance) {
    // 开始滚动
    if (heightChartInstance.scrollStart) {
      heightChartInstance.scrollStart(e)
    }
    // 显示提示
    if (heightChartInstance.showToolTip) {
      heightChartInstance.showToolTip(e)
    }
  }
}

const moveHeight = (e: any) => {
  if (heightChartInstance) {
    // 滚动图表
    if (heightChartInstance.scroll) {
      heightChartInstance.scroll(e)
    }
    // 更新提示位置
    if (heightChartInstance.showToolTip) {
      heightChartInstance.showToolTip(e)
    }
  }
}

const touchEndHeight = (e: any) => {
  if (heightChartInstance) {
    // 结束滚动
    if (heightChartInstance.scrollEnd) {
      heightChartInstance.scrollEnd(e)
    }
  }
}

const touchWeight = (e: any) => {
  console.log('[Statistics] 体重图表 touchstart', e)
  if (weightChartInstance) {
    // 开始滚动
    if (weightChartInstance.scrollStart) {
      weightChartInstance.scrollStart(e)
    }
    // 显示提示
    if (weightChartInstance.showToolTip) {
      weightChartInstance.showToolTip(e)
    }
  }
}

const moveWeight = (e: any) => {
  if (weightChartInstance) {
    // 滚动图表
    if (weightChartInstance.scroll) {
      weightChartInstance.scroll(e)
    }
    // 更新提示位置
    if (weightChartInstance.showToolTip) {
      weightChartInstance.showToolTip(e)
    }
  }
}

const touchEndWeight = (e: any) => {
  if (weightChartInstance) {
    // 结束滚动
    if (weightChartInstance.scrollEnd) {
      weightChartInstance.scrollEnd(e)
    }
  }
}

// 跳转到AI分析页面
const goToAIAnalysis = () => {
  if (!currentBaby.value) {
    uni.showToast({
      title: '请先选择宝宝',
      icon: 'none'
    })
    return
  }
  
  uni.navigateTo({
    url: '/pages/statistics/ai-analysis'
  })
}

// 初始化页面数据
const initPageData = async () => {
  console.log('[Statistics] 初始化页面数据')
  
  if (!isLoggedIn.value) {
    console.log('[Statistics] 用户未登录')
    return
  }

  if (!currentBaby.value) {
    console.log('[Statistics] 未选择宝宝')
    uni.showToast({
      title: '请先选择宝宝',
      icon: 'none'
    })
    setTimeout(() => {
      uni.navigateBack()
    }, 1500)
    return
  }

  // 加载数据
  await loadRecords()
  
  // 绘制图表
  drawCharts()
}

// 页面加载
onMounted(() => {
  console.log('[Statistics] 页面挂载')
  initPageData()
})

// 页面显示时重新加载数据
onShow(() => {
  console.log('[Statistics] 页面显示')
  initPageData()
})

// 组件卸载时清理图表
onBeforeUnmount(() => {
  console.log('[Statistics] 组件卸载，清理图表')
  clearCharts()
})
</script>

<style lang="scss" scoped>
.statistics-page {
  min-height: 100vh;
  background: #f6f8f7;
  padding-bottom: 40rpx;
}

.guest-tip {
  background: rgba(125, 211, 162, 0.15);
  color: #333;
  padding: 24rpx 30rpx;
  margin: 20rpx;
  border-radius: 12rpx;
  display: flex;
  align-items: center;
  gap: 16rpx;
  border: 1rpx solid rgba(125, 211, 162, 0.3);
}

.tip-icon {
  width: 36rpx;
  height: 36rpx;
}

.tip-text {
  font-size: 28rpx;
  font-weight: 500;
}

.time-range {
  background: white;
}

.ai-analysis-entry {
  margin: 20rpx;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 16rpx;
  padding: 24rpx;
  box-shadow: 0 4rpx 12rpx rgba(102, 126, 234, 0.3);
  cursor: pointer;
  transition: all 0.3s ease;

  &:active {
    transform: scale(0.98);
  }

  .ai-entry-content {
    display: flex;
    justify-content: space-between;
    align-items: center;

    .ai-entry-left {
      display: flex;
      align-items: center;
      gap: 16rpx;

      .ai-entry-icon {
        font-size: 48rpx;
        filter: drop-shadow(0 2rpx 4rpx rgba(0, 0, 0, 0.2));
      }

      .ai-entry-text {
        .ai-entry-title {
          display: block;
          font-size: 32rpx;
          font-weight: 600;
          color: #ffffff;
          margin-bottom: 8rpx;
        }

        .ai-entry-subtitle {
          display: block;
          font-size: 24rpx;
          color: rgba(255, 255, 255, 0.9);
        }
      }
    }

    .ai-entry-right {
      .ai-entry-arrow {
        font-size: 32rpx;
        color: #ffffff;
        font-weight: bold;
      }
    }
  }
}

.stat-section {
  background: white;
  border: 1rpx solid #CAE3D4;
  border-radius: 16rpx;
  margin: 20rpx 20rpx 0;
  padding: 30rpx;
  box-shadow: 0 2rpx 8rpx rgba(125, 211, 162, 0.08);
}

.section-header {
  display: flex;
  align-items: center;
  gap: 12rpx;
  margin-bottom: 24rpx;

  .icon {
    font-size: 40rpx;
  }

  .icon-img {
    width: 48rpx;
    height: 48rpx;
  }

  .title {
    font-size: 32rpx;
    font-weight: bold;
  }
}

.stat-cards {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20rpx;
}

.stat-card {
  background: #f6f8f7;
  border-radius: 12rpx;
  padding: 24rpx;
  text-align: center;
}

.card-label {
  font-size: 24rpx;
  color: #999;
  margin-bottom: 12rpx;
}

.card-value {
  font-size: 32rpx;
  font-weight: bold;
  color: #7dd3a2;
}

.sleep-quality {
  margin-top: 30rpx;
  background: #f0f9f6;
  border-radius: 12rpx;
  padding: 24rpx;
  border: 1rpx solid #CAE3D4;
}

.quality-title {
  font-size: 28rpx;
  font-weight: bold;
  margin-bottom: 20rpx;
  color: #333;
}

.quality-content {
  display: flex;
  flex-direction: column;
  gap: 16rpx;
}

.quality-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 26rpx;
}

.quality-label {
  color: #666;
}

.quality-value {
  color: #333;
  font-weight: 500;
}

.quality-recommendation {
  margin-top: 16rpx;
  padding: 20rpx;
  background: rgba(125, 211, 162, 0.15);
  border-radius: 12rpx;
  display: flex;
  align-items: center;
  gap: 12rpx;
  border: 1rpx solid #CAE3D4;
}

.recommendation-icon {
  width: 32rpx;
  height: 32rpx;
}

.recommendation-text {
  flex: 1;
  font-size: 26rpx;
  color: #333;
  line-height: 1.6;
}


.daily-chart {
  margin-top: 30rpx;
}

.chart-title {
  font-size: 28rpx;
  font-weight: bold;
  margin-bottom: 20rpx;
}

.chart-canvas {
  width: 100%;
  height: 550rpx;
  display: block;
  background-color: #ffffff;
}

.growth-charts {
  margin-top: 30rpx;
}

.chart-container {
  margin-bottom: 40rpx;

  &:last-child {
    margin-bottom: 0;
  }
}
</style>