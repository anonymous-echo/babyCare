/**
 * 喂养记录 API 接口
 * 职责: 纯 API 调用,无状态,无副作用
 */
import { get, post, put, del } from '@/utils/request'
import type { FeedingDetail } from '@/types'

// ============ 类型定义 ============

/**
 * API 响应: 喂养记录详情
 */
export interface FeedingRecordResponse {
  recordId: string
  babyId: string
  feedingType: 'breast' | 'bottle' | 'food'
  amount?: number
  duration?: number
  detail: FeedingDetail // 强类型 Detail
  note?: string
  feedingTime: number
  actualCompleteTime?: number // 实际喂养完成时间戳(毫秒)
  createBy: string
  createTime: number
}

/**
 * API 响应: 喂养记录列表
 */
export interface FeedingRecordsListResponse {
  records: FeedingRecordResponse[]
  total: number
  page: number
  pageSize: number
}

/**
 * API 请求: 创建喂养记录
 */
export interface CreateFeedingRecordRequest {
  babyId: string
  feedingType: 'breast' | 'bottle' | 'food'
  amount?: number
  duration?: number
  detail: FeedingDetail // 强类型 Detail
  note?: string
  feedingTime: number
  actualCompleteTime?: number // 实际喂养完成时间戳(毫秒)，用于准确计算提醒时间
  // 新增：用户自定义提醒间隔
  reminderInterval?: number // 提醒间隔(分钟)
}

// ============ API 函数 ============

/**
 * 获取喂养记录列表
 *
 * @param params 查询参数
 * @returns Promise<FeedingRecordsListResponse>
 */
export async function apiFetchFeedingRecords(params: {
  babyId: string
  startTime?: number
  endTime?: number
  page?: number
  pageSize?: number
}): Promise<FeedingRecordsListResponse> {
  const response = await get<FeedingRecordsListResponse>('/feeding-records', params)
  return response.data || { records: [], total: 0, page: 1, pageSize: 10 }
}

/**
 * 创建喂养记录
 *
 * @param data 创建请求数据
 * @returns Promise<FeedingRecordResponse>
 */
export async function apiCreateFeedingRecord(
  data: CreateFeedingRecordRequest
): Promise<FeedingRecordResponse> {
  const response = await post<FeedingRecordResponse>('/feeding-records', data)
  if (!response.data) {
    throw new Error(response.message || '创建喂养记录失败')
  }
  return response.data
}

/**
 * 获取单条喂养记录
 *
 * @param recordId 记录ID
 * @returns Promise<FeedingRecordResponse>
 */
export async function apiGetFeedingRecordById(
  recordId: string
): Promise<FeedingRecordResponse> {
  const response = await get<FeedingRecordResponse>(`/feeding-records/${recordId}`)
  if (!response.data) {
    throw new Error(response.message || '获取喂养记录失败')
  }
  return response.data
}

/**
 * 更新喂养记录
 *
 * @param recordId 记录ID
 * @param data 更新数据
 * @returns Promise<void>
 */
export async function apiUpdateFeedingRecord(
  recordId: string,
  data: Partial<CreateFeedingRecordRequest>
): Promise<void> {
  await put(`/feeding-records/${recordId}`, data)
}

/**
 * 删除喂养记录
 *
 * @param recordId 记录ID
 * @returns Promise<void>
 */
export async function apiDeleteFeedingRecord(recordId: string): Promise<void> {
  await del(`/feeding-records/${recordId}`)
}
