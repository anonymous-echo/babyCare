// AI分析相关类型定义

/**
 * AI分析类型
 */
export type AIAnalysisType = 'feeding' | 'sleep' | 'growth' | 'health' | 'behavior'

/**
 * AI分析状态
 */
export type AIAnalysisStatus = 'pending' | 'analyzing' | 'completed' | 'failed'

/**
 * AI分析实体
 */
export interface AIAnalysis {
  id: number
  baby_id: number
  analysis_type: AIAnalysisType
  status: AIAnalysisStatus
  start_date: string
  end_date: string
  input_data?: string
  result?: AIAnalysisResult
  score?: number
  insights?: string[]
  alerts?: string[]
  created_at: string
  updated_at: string
}

/**
 * AI分析结果
 */
export interface AIAnalysisResult {
  analysis_id: number
  baby_id: number
  analysis_type: AIAnalysisType
  score: number
  insights: AIInsight[]
  alerts: AIAlert[]
  patterns: AIPattern[]
  predictions: AIPrediction[]
  metadata?: Record<string, any>
  user_friendly?: UserFriendlyResult // 新增用户友好结果
}

/**
 * AI洞察建议
 */
export interface AIInsight {
  type: string
  title: string
  description: string
  priority: 'high' | 'medium' | 'low'
  category: string
}

/**
 * AI异常警告
 */
export interface AIAlert {
  level: 'critical' | 'warning' | 'info'
  type: string
  title: string
  description: string
  suggestion: string
  timestamp: string
}

/**
 * AI识别模式
 */
export interface AIPattern {
  pattern_type: string
  description: string
  confidence: number
  frequency: string
  time_range: TimeRange
}

/**
 * AI预测结果
 */
export interface AIPrediction {
  prediction_type: string
  value: string
  confidence: number
  time_frame: string
  reason: string
}

/**
 * 时间范围
 */
export interface TimeRange {
  start: string
  end: string
}

/**
 * 用户友好的分析结果
 */
export interface UserFriendlyResult {
  overall_summary: string // 总体评价
  score_explanation: string // 评分说明
  key_highlights: UserFriendlyHighlight[] // 关键亮点
  improvement_areas: UserFriendlyImprovement[] // 改进建议
  next_step_actions: UserFriendlyAction[] // 下一步行动
  encouraging_words: string // 鼓励话语
}

/**
 * 用户友好的亮点
 */
export interface UserFriendlyHighlight {
  title: string // 标题
  description: string // 描述
  icon: string // 图标建议
}

/**
 * 用户友好的改进建议
 */
export interface UserFriendlyImprovement {
  area: string // 改进领域
  issue: string // 问题描述
  suggestion: string // 具体建议
  priority: 'high' | 'medium' | 'low' // 优先级
  difficulty: 'easy' | 'medium' | 'hard' // 实施难度
}

/**
 * 用户友好的行动建议
 */
export interface UserFriendlyAction {
  action: string // 行动内容
  timeline: string // 时间安排
  benefit: string // 预期收益
  how_to: string // 具体做法
}

/**
 * 每日建议
 */
export interface DailyTips {
  tips: DailyTip[]
  generated_at: string
  expired_at: string
}

/**
 * 单个建议
 */
export interface DailyTip {
  id: string
  icon: string
  title: string
  description: string
  type: 'feeding' | 'sleep' | 'health' | 'growth' | 'behavior' | 'general'
  priority: 'high' | 'medium' | 'low'
  action_url?: string
}

/**
 * 分析统计
 */
export interface AnalysisStats {
  total_analyses: number
  completed_analyses: number
  failed_analyses: number
  average_score?: number
  analysis_type_counts: Record<AIAnalysisType, number>
  recent_analyses: AIAnalysis[]
}

/**
 * 创建分析请求
 */
export interface CreateAnalysisRequest {
  baby_id: number
  analysis_type: AIAnalysisType
  start_date: string
  end_date: string
  options?: Record<string, any>
}

/**
 * 分析响应
 */
export interface AnalysisResponse {
  analysis_id: number
  status: AIAnalysisStatus
  result?: AIAnalysisResult
  created_at: string
}

/**
 * 分析状态响应（用于轮询）
 */
export interface AnalysisStatusResponse {
  analysis_id: string
  status: AIAnalysisStatus
  progress: number // 进度百分比 0-100
  message: string  // 状态描述
  updated_at: string
}

/**
 * 批量分析响应
 */
export interface BatchAnalysisResponse {
  analyses: AnalysisResponse[]
  total_count: number
  completed_count: number
  failed_count: number
}

/**
 * 每日建议响应
 */
export interface DailyTipsResponse {
  tips: DailyTip[]
  generated_at: string
  expired_at: string
}

/**
 * 分析统计响应
 */
export interface AnalysisStatsResponse {
  total_analyses: number
  completed_analyses: number
  failed_analyses: number
  average_score?: number
  analysis_type_counts: Record<string, number>
  recent_analyses: AnalysisResponse[]
}

/**
 * 分析查询参数
 */
export interface AIAnalysisParams {
  baby_id?: number
  analysis_type?: AIAnalysisType
  start_date?: string
  end_date?: string
  status?: AIAnalysisStatus
  limit?: number
  offset?: number
}

/**
 * AI图表数据
 */
export interface AIChartData {
  categories: string[]
  series: AISeries[]
  title?: string
  subtitle?: string
}

/**
 * AI图表系列
 */
export interface AISeries {
  name: string
  data: number[]
  color?: string
  type?: 'line' | 'column' | 'radar' | 'pie'
}

/**
 * 喂养模式分析数据
 */
export interface FeedingPatternData {
  regularity: number  // 规律性评分
  adequacy: number    // 适量性评分
  timeliness: number  // 及时性评分
  diversity: number   // 多样性评分
}

/**
 * 睡眠质量分析数据
 */
export interface SleepQualityData {
  continuity: number  // 连续性评分
  duration: number    // 时长评分
  regularity: number  // 规律性评分
  depth: number       // 深度评分
}

/**
 * 成长发育分析数据
 */
export interface GrowthAssessmentData {
  height_percentile: number
  weight_percentile: number
  head_circumference_percentile: number
  growth_velocity: number
  development_milestone: number
}

/**
 * 健康风险评估
 */
export interface HealthRiskAssessment {
  overall_risk: 'low' | 'medium' | 'high'
  risk_factors: string[]
  recommendations: string[]
  monitoring_items: string[]
}

/**
 * AI分析配置
 */
export interface AIAnalysisConfig {
  enabled: boolean
  model: string
  timeout: number
  max_retries: number
  cache_enabled: boolean
  cache_ttl: number
}

/**
 * AI功能开关
 */
export interface AIFeatureFlags {
  analysis_enabled: boolean
  daily_tips_enabled: boolean
  real_time_monitoring: boolean
  predictive_analytics: boolean
  personalized_recommendations: boolean
}

/**
 * AI错误类型
 */
export interface AIError {
  code: string
  message: string
  details?: any
  timestamp: string
}

/**
 * AI分析状态更新事件
 */
export interface AIAnalysisStatusEvent {
  analysis_id: number
  status: AIAnalysisStatus
  progress?: number
  message?: string
  timestamp: string
}

/**
 * 实时AI监控数据
 */
export interface RealTimeAIMonitoring {
  current_status: 'normal' | 'warning' | 'critical'
  key_metrics: Record<string, number>
  recent_alerts: AIAlert[]
  recommendations: string[]
  last_updated: string
}