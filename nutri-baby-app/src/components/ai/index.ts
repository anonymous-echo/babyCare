// AI组件统一导出

export { default as AIInsightCard } from './AIInsightCard.vue'
export { default as AIAlertCard } from './AIAlertCard.vue'
export { default as AIScoreCard } from './AIScoreCard.vue'

// 导出类型
export type {
  AIInsight,
  AIAlert,
  DailyTip,
  UserFriendlyResult
} from '@/types/ai'