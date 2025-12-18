/**
 * 订阅消息 API 接口封装
 */

import { request } from '@/utils/request'

/**
 * 订阅授权记录
 */
export interface SubscribeAuthRecord {
  templateId: string
  templateType: string
  status: 'accept' | 'reject'
}

/**
 * 订阅授权响应
 */
export interface SubscribeAuthResponse {
  successCount: number
  failedCount: number
}

/**
 * 订阅项
 */
export interface SubscriptionItem {
  templateType: string
  templateId: string
  status: 'available' | 'used' | 'expired'
  subscribeTime: number // Unix timestamp (秒)
  expireTime?: number // Unix timestamp (秒)
}

/**
 * 订阅状态响应
 */
export interface SubscribeStatusResponse {
  subscriptions: SubscriptionItem[]
}

/**
 * 消息日志项
 */
export interface MessageLogItem {
  id: number
  templateType: string
  sendStatus: 'success' | 'failed' | 'pending'
  errmsg?: string
  sendTime?: number // Unix timestamp (秒)
  createdAt: number // Unix timestamp (秒)
}

/**
 * 消息日志响应
 */
export interface MessageLogsResponse {
  logs: MessageLogItem[]
  total: number
}

/**
 * 上报订阅授权记录
 * @param records 授权记录列表
 * @returns Promise<SubscribeAuthResponse>
 */
export async function saveSubscribeAuth(
  records: SubscribeAuthRecord[]
): Promise<SubscribeAuthResponse> {
  return request.post<SubscribeAuthResponse>('/api/v1/subscribe/authorize', {
    records
  })
}

/**
 * 获取用户订阅状态
 * @returns Promise<SubscribeStatusResponse>
 */
export async function getSubscribeStatus(): Promise<SubscribeStatusResponse> {
  return request.get<SubscribeStatusResponse>('/api/v1/subscribe/status')
}

/**
 * 取消订阅
 * @param templateType 模板类型
 * @returns Promise<void>
 */
export async function cancelSubscription(templateType: string): Promise<void> {
  return request.post('/api/v1/subscribe/cancel', {
    templateType
  })
}

/**
 * 获取消息发送日志
 * @param page 页码 (从1开始)
 * @param pageSize 每页数量
 * @returns Promise<MessageLogsResponse>
 */
export async function getMessageLogs(
  page: number = 1,
  pageSize: number = 20
): Promise<MessageLogsResponse> {
  const offset = (page - 1) * pageSize
  return request.get<MessageLogsResponse>('/api/v1/subscribe/logs', {
    offset,
    limit: pageSize
  })
}

/**
 * 喂养提醒请求数据
 */
export interface FeedingReminderRequest {
  babyId: string
  feedingType: 'breast' | 'bottle'
  lastFeedingTime: number // 上次喂养时间(Unix时间戳,毫秒)
  lastFeedingDuration?: number // 上次喂养时长(秒)
  lastFeedingAmount?: number // 上次喂养量(ml)
  lastFeedingBottleType?: 'formula' | 'breast-milk'
  templateType: 'breast_feeding_reminder' | 'bottle_feeding_reminder'
}

/**
 * 发送喂养提醒消息
 * @param data 喂养提醒数据
 * @returns Promise<void>
 */
export async function sendFeedingReminder(
  data: FeedingReminderRequest
): Promise<void> {
  return request.post('/api/v1/subscribe/feeding-reminder', data)
}

/**
 * 发送长时间喂养警告
 * @param babyId 宝宝ID
 * @param startTime 开始时间(Unix时间戳,毫秒)
 * @param duration 已喂养时长(秒)
 * @returns Promise<void>
 */
export async function sendFeedingDurationAlert(
  babyId: string,
  startTime: number,
  duration: number
): Promise<void> {
  return request.post('/api/v1/subscribe/feeding-duration-alert', {
    babyId,
    startTime,
    duration
  })
}
