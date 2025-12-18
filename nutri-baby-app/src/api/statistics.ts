import { request } from '@/utils/request'
import type { ApiResponse } from '@/types'
// ============ 统计 DTO ============

// 今日喂养统计
export interface TodayFeedingStats {
  breastCount: number  // 母乳喂养次数
  bottleMl: number     // 奶瓶总毫升数
  totalCount: number   // 总喂养次数
  lastFeedingTime?: number // 最后一次喂养时间戳(毫秒)
}

// 今日睡眠统计
export interface TodaySleepStats {
  totalMinutes: number   // 总睡眠分钟数
  lastSleepMinutes: number // 上次睡眠分钟数
  sessionCount: number   // 睡眠次数
}

// 今日换尿布统计
export interface TodayDiaperStats {
  totalCount: number // 总换尿布次数
  wetCount: number   // 尿湿次数
  dirtyCount: number // 排便次数
}

// 今日成长统计
export interface TodayGrowthStats {
  latestWeight?: number            // 最新体重 (kg)
  latestHeight?: number            // 最新身高 (cm)
  latestHeadCircumference?: number // 最新头围 (cm)
}

// 今日统计
export interface TodayStatistics {
  feeding: TodayFeedingStats
  sleep: TodaySleepStats
  diaper: TodayDiaperStats
  growth: TodayGrowthStats
}

// 本周喂养统计
export interface WeeklyFeedingStats {
  totalCount: number  // 本周总喂养次数
  trend: number       // 趋势对比（与上周的差异）
  avgPerDay: number   // 日均喂养次数
}

// 本周睡眠统计
export interface WeeklySleepStats {
  totalMinutes: number  // 本周总睡眠分钟数
  trend: number       // 趋势对比（与上周的分钟数差异）
  avgPerDay: number   // 日均睡眠分钟数
}

// 本周成长统计
export interface WeeklyGrowthStats {
  weightGain: number       // 周内体重增长 (kg)
  heightGain: number       // 周内身高增长 (cm)
  weekStartWeight?: number // 周初体重 (kg)
}

// 本周统计
export interface WeeklyStatistics {
  feeding: WeeklyFeedingStats
  sleep: WeeklySleepStats
  growth: WeeklyGrowthStats
}

// 完整统计响应
export interface BabyStatisticsResponse {
  today: TodayStatistics
  weekly: WeeklyStatistics
}

// ============ 按日统计 DTO ============

// 每日喂养统计项
export interface DailyFeedingStatsItem {
  date: string          // 日期，格式 YYYY-MM-DD
  feedingType: string   // 喂养类型：breast/bottle/food
  totalCount: number    // 总次数
  totalAmount: number   // 总量（ml）
  totalDuration: number // 总时长（秒）
}

// 每日睡眠统计项
export interface DailySleepStatsItem {
  date: string          // 日期，格式 YYYY-MM-DD
  totalDuration: number // 总时长（秒）
  totalCount: number    // 总次数
}

// 每日排泄统计项
export interface DailyDiaperStatsItem {
  date: string       // 日期，格式 YYYY-MM-DD
  diaperType: string // 排泄类型：pee/poop/both
  totalCount: number // 总次数
}

// 每日成长统计项
export interface DailyGrowthStatsItem {
  date: string              // 日期，格式 YYYY-MM-DD
  latestHeight?: number     // 最新身高（cm）
  latestWeight?: number     // 最新体重（g）
  latestHeadCircumference?: number // 最新头围（cm）
  recordCount: number       // 当日记录数
}

// 按日统计响应
export interface DailyStatsResponse {
  feeding?: DailyFeedingStatsItem[] // 喂养统计
  sleep?: DailySleepStatsItem[]     // 睡眠统计
  diaper?: DailyDiaperStatsItem[]   // 排泄统计
  growth?: DailyGrowthStatsItem[]   // 成长统计
}

// ============ API 接口 ============

/**
 * 获取宝宝统计数据
 * @param babyId 宝宝ID
 * @returns 统计数据响应
 */
export const apiFetchBabyStatistics = (babyId: string): Promise<ApiResponse<BabyStatisticsResponse>> => {
  return request({
    method: 'GET',
    url: `/babies/${babyId}/statistics`,
  })
}

/**
 * 获取按日统计数据
 * @param babyId 宝宝ID
 * @param startDate 开始日期（毫秒时间戳）
 * @param endDate 结束日期（毫秒时间戳）
 * @param types 统计类型，逗号分隔：feeding,sleep,diaper,growth，默认全部
 * @returns 按日统计数据响应
 */
export const apiFetchDailyStats = (params: {
  babyId: string
  startDate: number
  endDate: number
  types?: string
}): Promise<ApiResponse<DailyStatsResponse>> => {
  return request({
    method: 'GET',
    url: `/babies/${params.babyId}/daily-stats`,
    data: {
      startDate: params.startDate,
      endDate: params.endDate,
      types: params.types || 'feeding,sleep,diaper,growth'
    }
  })
}
