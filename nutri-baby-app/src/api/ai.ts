import { request } from '@/utils/request'
import type {
  AIAnalysisParams,
  AnalysisResponse,
  AnalysisStatusResponse,
  BatchAnalysisResponse,
  CreateAnalysisRequest,
  DailyTipsResponse,
  AnalysisStatsResponse,
  AIAnalysisType
} from '@/types/ai'
import type { ApiResponse } from '@/types'

/**
 * AI分析相关API
 */

/**
 * 创建AI分析任务
 */
export const createAIAnalysis = (data: CreateAnalysisRequest): Promise<ApiResponse<AnalysisResponse>> => {
  return request<AnalysisResponse>({
    url: '/ai-analysis',
    method: 'POST',
    data,
    retry: 2, // 失败时重试2次
    retryDelay: 1000 // 重试延迟1秒
  })
}

/**
 * 获取AI分析结果
 */
export const getAIAnalysisResult = (analysisId: number): Promise<ApiResponse<AnalysisResponse>> => {
  return request<AnalysisResponse>({
    url: `/ai-analysis/${analysisId}`,
    method: 'GET'
  })
}

/**
 * 获取分析状态（用于轮询）
 */
export const getAnalysisStatus = (analysisId: string): Promise<ApiResponse<AnalysisStatusResponse>> => {
  return request<AnalysisStatusResponse>({
    url: `/ai-analysis/${analysisId}/status`,
    method: 'GET',
    retry: 1, // 轮询失败时重试1次
    retryDelay: 500,
    showError: false // 轮询失败不显示错误提示
  })
}

/**
 * 获取最新AI分析结果
 */
export const getLatestAIAnalysis = (babyId: number, analysisType?: AIAnalysisType): Promise<ApiResponse<AnalysisResponse>> => {
  const params: any = {}
  if (analysisType) {
    params.type = analysisType
  }
  
  return request<AnalysisResponse>({
    url: `/ai-analysis/baby/${babyId}/latest`,
    method: 'GET',
    data: params
  })
}

/**
 * 批量AI分析
 */
export const batchAIAnalysis = (babyId: number, startDate: string, endDate: string): Promise<ApiResponse<BatchAnalysisResponse>> => {
  return request<BatchAnalysisResponse>({
    url: '/ai-analysis/batch',
    method: 'POST',
    data: {
      baby_id: babyId,
      start_date: startDate,
      end_date: endDate
    },
    retry: 2,
    retryDelay: 1000
  })
}

/**
 * 获取AI分析统计
 */
export const getAIAnalysisStats = (babyId: number, days: number = 30): Promise<ApiResponse<AnalysisStatsResponse>> => {
  return request<AnalysisStatsResponse>({
    url: `/ai-analysis/baby/${babyId}/history`,
    method: 'GET',
    data: {
      days: days
    }
  })
}

/**
 * 生成每日建议
 */
export const generateDailyTips = (babyId: number, date?: string): Promise<ApiResponse<DailyTipsResponse>> => {
  return request<DailyTipsResponse>({
    url: `/ai-analysis/daily-tips/${babyId}/generate`,
    method: 'POST',
    data: date ? { date } : {}
  })
}

/**
 * 获取每日建议
 */
export const getDailyTips = (babyId: number, date?: string): Promise<ApiResponse<DailyTipsResponse>> => {
  return request<DailyTipsResponse>({
    url: `/ai-analysis/daily-tips/${babyId}`,
    method: 'GET',
    data: date ? { date } : undefined
  })
}

/**
 * 轮询分析状态
 */
export const pollAnalysisStatus = async (
  analysisId: string,
  onStatusUpdate: (status: string, progress?: number, message?: string) => void,
  maxAttempts = 30,
  interval = 2000
): Promise<AnalysisResponse> => {
  for (let attempt = 0; attempt < maxAttempts; attempt++) {
    try {
      // 使用专门的状态查询API
      const statusResponse = await getAnalysisStatus(analysisId)
      const statusResult = statusResponse.data

      // 更新状态，包含进度和消息
      onStatusUpdate(statusResult.status, statusResult.progress, statusResult.message)

      // 如果分析完成，获取完整结果
      if (statusResult.status === 'completed') {
        const resultResponse = await getAIAnalysisResult(parseInt(analysisId))
        return resultResponse.data
      }

      // 如果分析失败，抛出错误
      if (statusResult.status === 'failed') {
        throw new Error(statusResult.message || '分析失败')
      }

      // 等待下次轮询
      await new Promise(resolve => setTimeout(resolve, interval))
    } catch (error) {
      console.error(`轮询分析状态失败 (attempt ${attempt + 1}):`, error)

      // 最后一次尝试失败则抛出错误
      if (attempt === maxAttempts - 1) {
        throw error
      }

      // 继续下一次尝试
      await new Promise(resolve => setTimeout(resolve, interval))
    }
  }

  throw new Error('分析超时')
}

/**
 * 获取分析图表数据
 */
export const getAnalysisChartData = (analysisType: AIAnalysisType, data: any) => {
  switch (analysisType) {
    case 'feeding':
      return getFeedingAnalysisChartData(data)
    case 'sleep':
      return getSleepAnalysisChartData(data)
    case 'growth':
      return getGrowthAnalysisChartData(data)
    case 'health':
      return getHealthAnalysisChartData(data)
    default:
      return null
  }
}

/**
 * 获取喂养分析图表数据
 */
const getFeedingAnalysisChartData = (data: any) => {
  if (!data.patterns || !data.patterns.length) return null

  const pattern = data.patterns[0]
  const categories = ['规律性', '适量性', '及时性', '多样性']
  const scores = [
    pattern.regularity || 0,
    pattern.adequacy || 0,
    pattern.timeliness || 0,
    pattern.diversity || 0
  ]

  return {
    categories,
    series: [{
      name: '喂养模式评分',
      data: scores.map(score => Math.round(score * 100)),
      color: '#7dd3a2'
    }],
    title: '喂养模式分析',
    subtitle: '基于AI智能分析的综合评分'
  }
}

/**
 * 获取睡眠分析图表数据
 */
const getSleepAnalysisChartData = (data: any) => {
  if (!data.patterns || !data.patterns.length) return null

  const pattern = data.patterns[0]
  const categories = ['连续性', '时长', '规律性', '深度']
  const scores = [
    pattern.continuity || 0,
    pattern.duration || 0,
    pattern.regularity || 0,
    pattern.depth || 0
  ]

  return {
    categories,
    series: [{
      name: '睡眠质量评分',
      data: scores.map(score => Math.round(score * 100)),
      color: '#52c41a'
    }],
    title: '睡眠质量分析',
    subtitle: '基于AI智能分析的综合评分'
  }
}

/**
 * 获取成长分析图表数据
 */
const getGrowthAnalysisChartData = (data: any) => {
  if (!data.predictions || !data.predictions.length) return null

  const predictions = data.predictions.filter((p: any) => p.prediction_type === 'growth')
  const categories = predictions.map((p: any) => p.time_frame)
  const values = predictions.map((p: any) => parseFloat(p.value) || 0)

  return {
    categories,
    series: [{
      name: '预测值',
      data: values,
      color: '#ff6b6b'
    }],
    title: '成长趋势预测',
    subtitle: '基于AI智能分析的预测结果'
  }
}

/**
 * 获取健康分析图表数据
 */
const getHealthAnalysisChartData = (data: any) => {
  if (!data.alerts || !data.alerts.length) return null

  const alerts = data.alerts
  const levelCounts = {
    critical: alerts.filter((a: any) => a.level === 'critical').length,
    warning: alerts.filter((a: any) => a.level === 'warning').length,
    info: alerts.filter((a: any) => a.level === 'info').length
  }

  return {
    categories: ['严重', '警告', '提示'],
    series: [{
      name: '健康预警',
      data: [levelCounts.critical, levelCounts.warning, levelCounts.info],
      color: '#ff4757',
      type: 'column'
    }],
    title: '健康预警分布',
    subtitle: '基于AI智能分析的风险评估'
  }
}

/**
 * 分析状态文本
 */
export const getAnalysisStatusText = (status: string): string => {
  const statusMap: Record<string, string> = {
    pending: '等待分析',
    analyzing: '分析中...',
    completed: '分析完成',
    failed: '分析失败'
  }
  return statusMap[status] || status
}

/**
 * 分析状态颜色
 */
export const getAnalysisStatusColor = (status: string): string => {
  const colorMap: Record<string, string> = {
    pending: '#ffa940',
    analyzing: '#1890ff',
    completed: '#52c41a',
    failed: '#ff4d4f'
  }
  return colorMap[status] || '#8c8c8c'
}

/**
 * 优先级颜色
 */
export const getPriorityColor = (priority: string): string => {
  const colorMap: Record<string, string> = {
    high: '#ff4d4f',
    medium: '#faad14',
    low: '#52c41a'
  }
  return colorMap[priority] || '#8c8c8c'
}

/**
 * 警告级别颜色
 */
export const getAlertLevelColor = (level: string): string => {
  const colorMap: Record<string, string> = {
    critical: '#ff4d4f',
    warning: '#faad14',
    info: '#1890ff'
  }
  return colorMap[level] || '#8c8c8c'
}

/**
 * 分析类型图标
 */
export const getAnalysisTypeIcon = (type: AIAnalysisType): string => {
  const iconMap: Record<AIAnalysisType, string> = {
    feeding: '🍼',
    sleep: '😴',
    growth: '📈',
    health: '❤️',
    behavior: '🧠'
  }
  return iconMap[type] || '🤖'
}

/**
 * 分析类型名称
 */
export const getAnalysisTypeName = (type: AIAnalysisType): string => {
  const nameMap: Record<AIAnalysisType, string> = {
    feeding: '喂养分析',
    sleep: '睡眠分析',
    growth: '成长分析',
    health: '健康分析',
    behavior: '行为分析'
  }
  return nameMap[type] || '未知分析'
}