import { reactive, ref, computed } from 'vue'
import type {
  AIAnalysis,
  AnalysisResponse,
  DailyTipsResponse,
  AnalysisStatsResponse,
  AIAnalysisType,
  AIAnalysisStatus,
  DailyTip
} from '@/types/ai'
import {
  createAIAnalysis,
  getAIAnalysisResult,
  getLatestAIAnalysis,
  getAIAnalysisStats,
  getDailyTips as apiGetDailyTips,
  generateDailyTips as apiGenerateDailyTips,
  pollAnalysisStatus as apiPollAnalysisStatus,
  batchAIAnalysis
} from '@/api/ai'

/**
 * AI分析状态管理
 */
export const useAIStore = () => {
  // 状态定义
  const analyses = reactive<Record<number, AIAnalysis>>({}) // 分析记录映射
  const dailyTips = ref<Record<string, DailyTip[]>>({}) // 每日建议映射 - 使用ref确保响应式
  const stats = ref<AnalysisStatsResponse | null>(null) // 分析统计
  const isAnalyzing = ref(false) // 是否正在分析
  const analyzingIds = reactive<Set<number>>(new Set()) // 正在分析的ID集合
  const currentAnalysis = ref<AIAnalysis | null>(null) // 当前分析
  const backgroundPollingEnabled = ref(true) // 是否启用后台轮询
  const pollingTimers = reactive<Map<number, number>>(new Map()) // 轮询定时器映射

  // 计算属性
  const completedAnalyses = computed(() => {
    return Object.values(analyses).filter(analysis => analysis.status === 'completed')
  })

  const pendingAnalyses = computed(() => {
    return Object.values(analyses).filter(analysis =>
      analysis.status === 'pending' || analysis.status === 'analyzing'
    )
  })

  const failedAnalyses = computed(() => {
    return Object.values(analyses).filter(analysis => analysis.status === 'failed')
  })

  const todayTips = computed(() => {
    const today: string = new Date().toISOString().split('T')[0]
    return dailyTips.value[today] || []
  })

  const hasUnexpiredTips = computed(() => {
    const now = new Date()
    return Object.entries(dailyTips.value).some(([date, tips]) => {
      const tipDate = new Date(date)
      const expiryDate = new Date(tipDate.getTime() + 24 * 60 * 60 * 1000) // 24小时后过期
      return now < expiryDate && tips.length > 0
    })
  })

  // 方法

  /**
   * 创建AI分析任务
   */
  const createAnalysis = async (babyId: number, analysisType: AIAnalysisType, startDate: string, endDate: string) => {
    try {
      isAnalyzing.value = true

      const response = await createAIAnalysis({
        baby_id: babyId,
        analysis_type: analysisType,
        start_date: startDate,
        end_date: endDate
      })

      // 从API响应中提取数据
      const analysisData = response.data

      // 添加到分析记录
      const analysis: AIAnalysis = {
        id: analysisData.analysis_id,
        baby_id: babyId,
        analysis_type: analysisType,
        status: analysisData.status,
        start_date: startDate,
        end_date: endDate,
        created_at: analysisData.created_at,
        updated_at: analysisData.created_at
      }

      analyses[analysis.id] = analysis
      analyzingIds.add(analysis.id)

      // 开始轮询状态
      pollAnalysisStatusInternal(analysis.id)

      return analysis
    } catch (error) {
      console.error('创建AI分析任务失败:', error)
      throw error
    } finally {
      isAnalyzing.value = false
    }
  }

  /**
   * 轮询分析状态（内部方法）
   */
  const pollAnalysisStatusInternal = async (analysisId: number) => {
    try {
      const result = await apiPollAnalysisStatus(analysisId.toString(), (status, progress, message) => {
        // 更新分析状态
        if (analyses[analysisId]) {
          analyses[analysisId].status = status as AIAnalysisStatus
        }
        console.log(`分析${analysisId}状态更新: ${status}, 进度: ${progress}%, 消息: ${message}`)
      })

      // 更新完整分析结果
      if (result.result && analyses[analysisId]) {
        analyses[analysisId] = {
          ...analyses[analysisId],
          status: result.status as AIAnalysisStatus,
          result: result.result,
          score: result.result.score,
          insights: result.result.insights?.map(insight => JSON.stringify(insight)),
          alerts: result.result.alerts?.map(alert => JSON.stringify(alert)),
          updated_at: new Date().toISOString()
        }
      }

      // 从分析中移除
      analyzingIds.delete(analysisId)

      return result
    } catch (error) {
      console.error('轮询分析状态失败:', error)

      // 更新为失败状态
      if (analyses[analysisId]) {
        analyses[analysisId].status = 'failed'
        analyses[analysisId].updated_at = new Date().toISOString()
      }

      analyzingIds.delete(analysisId)
      throw error
    }
  }

  /**
   * 获取分析结果
   */
  const getAnalysisResult = async (analysisId: number): Promise<AIAnalysis> => {
    try {
      const response = await getAIAnalysisResult(analysisId)
      const analysisData = response.data

      // 更新本地缓存
      if (analyses[analysisId]) {
        analyses[analysisId].status = analysisData.status as AIAnalysisStatus
        if (analysisData.result) {
          analyses[analysisId].result = analysisData.result
          analyses[analysisId].score = analysisData.result.score
          analyses[analysisId].insights = analysisData.result.insights?.map(insight => JSON.stringify(insight))
          analyses[analysisId].alerts = analysisData.result.alerts?.map(alert => JSON.stringify(alert))
        }
        analyses[analysisId].updated_at = new Date().toISOString()
      }

      if (!analyses[analysisId]) {
        throw new Error('分析记录不存在')
      }
      return analyses[analysisId]
    } catch (error) {
      console.error('获取分析结果失败:', error)
      throw error
    }
  }

  /**
   * 获取最新分析结果
   */
  const getLatestAnalysis = async (babyId: number, analysisType: AIAnalysisType): Promise<AIAnalysis | null> => {
    try {
      const response = await getLatestAIAnalysis(babyId, analysisType)
      const analysisData = response.data

      if (analysisData.result) {
        const analysis: AIAnalysis = {
          id: analysisData.analysis_id,
          baby_id: babyId,
          analysis_type: analysisType,
          status: analysisData.status as AIAnalysisStatus,
          start_date: '', // 将从result中获取
          end_date: '',   // 将从result中获取
          result: analysisData.result,
          score: analysisData.result.score,
          insights: analysisData.result.insights?.map(insight => JSON.stringify(insight)),
          alerts: analysisData.result.alerts?.map(alert => JSON.stringify(alert)),
          created_at: analysisData.created_at,
          updated_at: analysisData.created_at
        }

        analyses[analysis.id] = analysis
        return analysis
      }

      return null
    } catch (error: any) {
      if (error?.response?.status === 404) {
        return null // 未找到分析结果
      }
      console.error('获取最新分析结果失败:', error)
      throw error
    }
  }

  /**
   * 获取分析统计
   */
  const getAnalysisStats = async (babyId: number): Promise<AnalysisStatsResponse> => {
    try {
      const response = await getAIAnalysisStats(babyId)
      // 正确更新ref值
      stats.value = response.data
      return response.data
    } catch (error) {
      console.error('获取分析统计失败:', error)
      throw error
    }
  }

  /**
   * 生成每日建议
   */
  const generateDailyTips = async (babyId: number, date?: string): Promise<DailyTip[]> => {
    try {
      const response = await apiGenerateDailyTips(babyId, date)

      const targetDate: string = date || new Date().toISOString().split('T')[0]

      // 直接赋值新对象触发ref更新
      dailyTips.value = { ...dailyTips.value, [targetDate]: response.data.tips }
      console.log('✅ 已生成并保存每日建议:', targetDate, response.data.tips.length, '条')

      return response.data.tips
    } catch (error: any) {
      console.error('生成每日建议失败:', error)
      // 如果是404错误，说明还没有数据，返回空数组
      if (error?.response?.status === 404) {
        return []
      }
      throw error
    }
  }

  /**
   * 获取每日建议
   */
  const getDailyTips = async (babyId: number, date?: string): Promise<DailyTip[]> => {
    try {
      const targetDate: string = date || new Date().toISOString().split('T')[0]

      // 如果已有缓存且未过期，直接返回
      const cachedTips = dailyTips.value[targetDate]
      if (cachedTips && cachedTips.length > 0) {
        console.log('✅ 使用缓存的每日建议:', targetDate, cachedTips.length, '条')
        return cachedTips
      }

      // 否则从服务器获取
      try {
        const response = await apiGetDailyTips(babyId, date)

        // 直接赋值新对象触发ref更新
        dailyTips.value = { ...dailyTips.value, [targetDate]: response.data.tips }
        console.log('✅ 已从API获取并保存每日建议:', targetDate, response.data.tips.length, '条')
        console.log('📊 当前dailyTips keys:', Object.keys(dailyTips.value))

        return response.data.tips
      } catch (error: any) {
        // 如果获取失败（404），尝试生成新的每日建议
        if (error?.response?.status === 404) {
          console.log('⚠️ 未找到每日建议，尝试生成新的建议...')
          return await generateDailyTips(babyId, date)
        }
        throw error
      }
    } catch (error) {
      console.error('❌ 获取每日建议失败:', error)
      // 返回空数组而不是抛出错误，避免阻塞页面加载
      return []
    }
  }

  /**
   * 清除分析缓存
   */
  const clearAnalysisCache = (analysisId?: number) => {
    if (analysisId) {
      delete analyses[analysisId]
    } else {
      // 清除所有分析缓存
      Object.keys(analyses).forEach(key => {
        delete analyses[parseInt(key)]
      })
    }
  }

  /**
   * 清除每日建议缓存
   */
  const clearDailyTipsCache = (date?: string) => {
    if (date && dailyTips.value[date]) {
      const newTips = { ...dailyTips.value }
      delete newTips[date]
      dailyTips.value = newTips
    } else if (!date) {
      // 清除所有每日建议缓存
      dailyTips.value = {}
    }
  }

  /**
   * 检查是否有活跃的分析任务
   */
  const hasActiveAnalysis = computed(() => {
    return analyzingIds.size > 0
  })

  /**
   * 获取指定类型的最新分析
   */
  const getLatestAnalysisByType = (analysisType: AIAnalysisType): AIAnalysis | null => {
    const typeAnalyses = Object.values(analyses).filter(
      analysis => analysis.analysis_type === analysisType && analysis.status === 'completed'
    )

    if (typeAnalyses.length === 0) return null

    // 按创建时间排序，返回最新的
    const sorted = typeAnalyses.sort((a, b) =>
      new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
    )

    return sorted[0] || null
  }

  /**
   * 获取分析概览
   */
  const getAnalysisOverview = (babyId: number) => {
    const babyAnalyses = Object.values(analyses).filter(analysis => analysis.baby_id === babyId)

    return {
      total: babyAnalyses.length,
      completed: babyAnalyses.filter(a => a.status === 'completed').length,
      pending: babyAnalyses.filter(a => a.status === 'pending' || a.status === 'analyzing').length,
      failed: babyAnalyses.filter(a => a.status === 'failed').length,
      averageScore: babyAnalyses
        .filter(a => a.score !== undefined)
        .reduce((sum, a) => sum + (a.score || 0), 0) / babyAnalyses.filter(a => a.score !== undefined).length || 0
    }
  }

  /**
   * 批量分析
   */
  const batchAnalyze = async (babyId: number, startDate: string, endDate: string) => {
    try {
      const response = await batchAIAnalysis(babyId, startDate, endDate)
      return response.data // 返回批量分析响应数据
    } catch (error) {
      console.error('批量分析失败:', error)
      throw error
    }
  }

  /**
   * 开始轮询分析状态
   */
  const startPolling = (analysisId: number, onStatusUpdate?: (status: string, progress?: number, message?: string) => void) => {
    // 如果已经在轮询，不重复启动
    if (pollingTimers.has(analysisId)) {
      console.log(`分析${analysisId}已在轮询中`)
      return
    }

    apiPollAnalysisStatus(analysisId.toString(), onStatusUpdate || (() => { }), 30, 2000)
      .then(result => {
        console.log('轮询完成:', result)
        pollingTimers.delete(analysisId)
      })
      .catch(error => {
        console.error('轮询失败:', error)
        pollingTimers.delete(analysisId)
      })

    // 标记为正在轮询
    pollingTimers.set(analysisId, Date.now())
  }

  /**
   * 停止轮询
   */
  const stopPolling = (analysisId: number) => {
    pollingTimers.delete(analysisId)
    analyzingIds.delete(analysisId)
  }

  /**
   * 停止所有轮询
   */
  const stopAllPolling = () => {
    pollingTimers.clear()
    analyzingIds.clear()
  }

  /**
   * 设置后台轮询
   */
  const setBackgroundPolling = (enabled: boolean) => {
    backgroundPollingEnabled.value = enabled
  }

  /**
   * 获取需要关注的事项
   */
  const getAttentionItems = (babyId: number) => {
    const babyAnalyses = Object.values(analyses).filter(analysis =>
      analysis.baby_id === babyId && analysis.status === 'completed'
    )

    const attentionItems: Array<{
      type: string
      title: string
      description: string
      level: string
      analysisType: string
      score?: number
    }> = []

    babyAnalyses.forEach(analysis => {
      // 检查警告
      if (analysis.alerts && analysis.alerts.length > 0) {
        analysis.alerts.forEach(alertStr => {
          try {
            const alert = JSON.parse(alertStr)
            if (alert.level === 'critical' || alert.level === 'warning') {
              attentionItems.push({
                type: 'alert',
                title: alert.title,
                description: alert.description,
                level: alert.level,
                analysisType: analysis.analysis_type
              })
            }
          } catch (e) {
            console.error('解析警告失败:', e)
          }
        })
      }

      // 检查低分分析
      if (analysis.score !== undefined && analysis.score < 60) {
        attentionItems.push({
          type: 'low_score',
          title: `${analysis.analysis_type}分析评分较低`,
          description: `评分为${analysis.score}分，建议关注`,
          level: 'warning',
          analysisType: analysis.analysis_type,
          score: analysis.score
        })
      }
    })

    return attentionItems.sort((a, b) => {
      const levelPriority: Record<string, number> = { critical: 3, warning: 2, info: 1 }
      return (levelPriority[b.level] || 0) - (levelPriority[a.level] || 0)
    })
  }

  return {
    // 状态
    analyses,
    dailyTips,
    stats,
    isAnalyzing,
    analyzingIds,
    currentAnalysis,
    backgroundPollingEnabled,

    // 计算属性
    completedAnalyses,
    pendingAnalyses,
    failedAnalyses,
    todayTips,
    hasUnexpiredTips,
    hasActiveAnalysis,

    // 方法
    createAnalysis,
    pollAnalysisStatus: pollAnalysisStatusInternal,
    getAnalysisResult,
    getLatestAnalysis,
    getAnalysisStats,
    generateDailyTips,
    getDailyTips,
    clearAnalysisCache,
    clearDailyTipsCache,
    getLatestAnalysisByType,
    getAnalysisOverview,
    getAttentionItems,
    batchAnalyze,
    startPolling,
    stopPolling,
    stopAllPolling,
    setBackgroundPolling
  }
}

// 导出单例实例
export const aiStore = useAIStore()