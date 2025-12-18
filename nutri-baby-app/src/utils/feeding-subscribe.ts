/**
 * 喂养订阅消息管理工具
 *
 * 功能:
 * 1. 一次性申请三个喂养相关的订阅消息权限 (母乳、奶瓶、时间过长提醒)
 * 2. 管理申请状态,避免频繁骚扰用户
 * 3. 记录和导出申请结果
 */

import {
  getAuthStatus,
  requestSubscribeMessage,
  recordGuideShown,
  dismissGuideForever,
  getAuthRecord
} from '@/store/subscribe'
import type { SubscribeMessageType, SubscribeAuthStatus } from '@/types'
import { StorageKeys, getStorage, setStorage } from '@/utils/storage'

/**
 * 喂养订阅消息类型
 */
const FEEDING_MESSAGE_TYPES: SubscribeMessageType[] = [
  'breast_feeding_reminder',      // 母乳喂养提醒
  'bottle_feeding_reminder',      // 奶瓶喂养提醒
  'feeding_duration_alert',       // 喂奶时间过长提醒
]

/**
 * 喂养订阅申请记录
 */
interface FeedingSubscribeRecord {
  lastRequestTime: number        // 上次申请时间
  requestCount: number           // 申请次数
  statusMap: Record<string, 'accept' | 'reject'> // 各消息的授权状态
  isDismissedForever: boolean    // 是否永久关闭
}

/**
 * 获取喂养订阅申请记录
 */
function getFeedingSubscribeRecord(): FeedingSubscribeRecord {
  const record = getStorage<FeedingSubscribeRecord>(StorageKeys.FEEDING_SUBSCRIBE_RECORD)
  if (record) {
    return record
  }
  return {
    lastRequestTime: 0,
    requestCount: 0,
    statusMap: {},
    isDismissedForever: false,
  }
}

/**
 * 保存喂养订阅申请记录
 */
function saveFeedingSubscribeRecord(record: FeedingSubscribeRecord) {
  setStorage(StorageKeys.FEEDING_SUBSCRIBE_RECORD, record)
}


/**
 * 一次性申请所有喂养订阅消息
 *
 * 这会触发微信官方的多选界面,用户可以:
 * - 为每个消息单独选择接受或拒绝
 * - 勾选"总是保持以上选择"
 * - 点击确认或取消
 *
 * @returns Promise<Map<SubscribeMessageType, 'accept' | 'reject'>>
 */
export async function requestAllFeedingSubscribeMessages(): Promise<Map<SubscribeMessageType, 'accept' | 'reject'>> {
  console.log('[FeedingSubscribe] 开始申请喂养订阅消息:', FEEDING_MESSAGE_TYPES)

  try {
    uni.showLoading({
      title: '请求授权中...',
      mask: true,
    })

    // 一次性请求所有三个消息
    const results = await requestSubscribeMessage(FEEDING_MESSAGE_TYPES)

    uni.hideLoading()

    console.log('[FeedingSubscribe] 申请结果:', results)
    console.log('[FeedingSubscribe] 申请结果详情:', {
      breast_feeding_reminder: results.get('breast_feeding_reminder'),
      bottle_feeding_reminder: results.get('bottle_feeding_reminder'),
      feeding_duration_alert: results.get('feeding_duration_alert')
    })

    // 记录申请
    recordFeedingSubscribeRequest(results)

    // 显示结果提示
    showResultToast(results)

    return results
  } catch (error: any) {
    uni.hideLoading()
    console.error('[FeedingSubscribe] 申请失败:', error)
    console.error('[FeedingSubscribe] 申请失败详情:', error.message || JSON.stringify(error))
    uni.showToast({
      title: '授权失败,请检查网络或重试',
      icon: 'none',
    })
    throw error
  }
}

/**
 * 记录申请结果
 */
function recordFeedingSubscribeRequest(results: Map<SubscribeMessageType, 'accept' | 'reject'>) {
  const record = getFeedingSubscribeRecord()
  const now = Date.now()

  record.lastRequestTime = now
  record.requestCount++

  // 更新每个消息的状态
  for (const type of FEEDING_MESSAGE_TYPES) {
    const status = results.get(type)
    if (status) {
      record.statusMap[type] = status
    }
  }

  saveFeedingSubscribeRecord(record)
  console.log('[FeedingSubscribe] 记录已保存:', record)
}

/**
 * 永久关闭喂养订阅申请弹窗
 */
export function dismissFeedingSubscribeForever() {
  const record = getFeedingSubscribeRecord()
  record.isDismissedForever = true
  record.lastRequestTime = Date.now()
  saveFeedingSubscribeRecord(record)
  console.log('[FeedingSubscribe] 已永久关闭申请弹窗')
}

/**
 * 重置喂养订阅申请状态(用于调试或用户主动重置)
 */
export function resetFeedingSubscribeRequest() {
  const emptyRecord: FeedingSubscribeRecord = {
    lastRequestTime: 0,
    requestCount: 0,
    statusMap: {},
    isDismissedForever: false,
  }
  saveFeedingSubscribeRecord(emptyRecord)
  console.log('[FeedingSubscribe] 已重置申请状态')
}

/**
 * 显示申请结果提示
 */
function showResultToast(results: Map<SubscribeMessageType, 'accept' | 'reject'>) {
  const acceptCount = Array.from(results.values()).filter((s) => s === 'accept').length
  const totalCount = results.size

  if (acceptCount === totalCount) {
    uni.showToast({
      title: '✨ 提醒已全部开启',
      icon: 'success',
      duration: 2000,
    })
  } else if (acceptCount > 0) {
    uni.showToast({
      title: `✅ ${acceptCount}/${totalCount} 个提醒已开启`,
      icon: 'success',
      duration: 2000,
    })
  } else {
    uni.showToast({
      title: '您拒绝了授权',
      icon: 'none',
      duration: 2000,
    })
  }
}

/**
 * 获取喂养订阅的当前状态
 *
 * @returns
 * {
 *   breast_feeding_reminder: { authStatus: 'accept' | 'reject' | 'ban' | 'unknown', lastRequestTime?: number },
 *   bottle_feeding_reminder: { ... },
 *   feeding_duration_alert: { ... }
 * }
 */
export function getFeedingSubscribeStatus(): Record<
  string,
  { authStatus: SubscribeAuthStatus; lastRequestTime?: number }
> {
  const status: Record<string, { authStatus: SubscribeAuthStatus; lastRequestTime?: number }> = {}

  for (const type of FEEDING_MESSAGE_TYPES) {
    const authStatus = getAuthStatus(type)
    const authRecord = getAuthRecord(type)

    status[type] = {
      authStatus,
      lastRequestTime: authRecord?.lastRequestTime,
    }
  }

  return status
}

/**
 * 获取喂养订阅统计信息
 */
export function getFeedingSubscribeStats() {
  const status = getFeedingSubscribeStatus()
  const record = getFeedingSubscribeRecord()

  const acceptCount = Object.values(status).filter((s) => s.authStatus === 'accept').length
  const rejectCount = Object.values(status).filter((s) => s.authStatus === 'reject').length
  const banCount = Object.values(status).filter((s) => s.authStatus === 'ban').length
  const unknownCount = Object.values(status).filter((s) => s.authStatus === 'unknown').length

  return {
    acceptCount,        // 已同意
    rejectCount,        // 已拒绝
    banCount,           // 已禁用(3次拒绝)
    unknownCount,       // 未申请
    totalCount: FEEDING_MESSAGE_TYPES.length,
    requestCount: record.requestCount,
    isDismissedForever: record.isDismissedForever,
    lastRequestTime: record.lastRequestTime,
  }
}

/**
 * 检查用户是否启用了任何喂养提醒
 */
export function hasEnabledFeedingReminders(): boolean {
  const status = getFeedingSubscribeStatus()
  return Object.values(status).some((s) => s.authStatus === 'accept')
}
