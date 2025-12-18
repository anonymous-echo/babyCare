/**
 * 换尿布记录 API 接口
 * 职责: 纯 API 调用,无状态,无副作用
 */
import { get, post, put, del } from '@/utils/request'

// ============ 类型定义 ============

/**
 * API 响应: 换尿布记录详情
 */
export interface DiaperRecordResponse {
  recordId: string
  babyId: string
  diaperType: 'pee' | 'poop' | 'both'
  pooColor?: string
  pooTexture?: string
  note?: string
  changeTime: number
  createBy: string
  createTime: number
}

/**
 * API 响应: 换尿布记录列表
 */
export interface DiaperRecordsListResponse {
  records: DiaperRecordResponse[]
  total: number
  page: number
  pageSize: number
}

type DiaperType = 'pee' | 'poop' | 'both'

/**
 * API 请求: 创建换尿布记录
 */
export interface CreateDiaperRecordRequest {
  babyId: string
  diaperType: DiaperType
  pooColor?: string
  pooTexture?: string
  note?: string
  changeTime: number
}

// ============ API 函数 ============

/**
 * 获取换尿布记录列表
 *
 * @param params 查询参数
 * @returns Promise<DiaperRecordsListResponse>
 */
export async function apiFetchDiaperRecords(params: {
  babyId: string
  startTime?: number
  endTime?: number
  page?: number
  pageSize?: number
}): Promise<DiaperRecordsListResponse> {
  const response = await get<DiaperRecordsListResponse>('/diaper-records', params)
  return response.data || { records: [], total: 0, page: 1, pageSize: 10 }
}

/**
 * 创建换尿布记录
 *
 * @param data 创建请求数据
 * @returns Promise<DiaperRecordResponse>
 */
export async function apiCreateDiaperRecord(
  data: CreateDiaperRecordRequest
): Promise<DiaperRecordResponse> {
  const response = await post<DiaperRecordResponse>('/diaper-records', data)
  if (!response.data) {
    throw new Error(response.message || '创建换尿布记录失败')
  }
  return response.data
}

/**
 * 获取单条换尿布记录
 *
 * @param recordId 记录ID
 * @returns Promise<DiaperRecordResponse>
 */
export async function apiGetDiaperRecordById(
  recordId: string
): Promise<DiaperRecordResponse> {
  const response = await get<DiaperRecordResponse>(`/diaper-records/${recordId}`)
  if (!response.data) {
    throw new Error(response.message || '获取换尿布记录失败')
  }
  return response.data
}

/**
 * 更新换尿布记录
 *
 * @param recordId 记录ID
 * @param data 更新数据
 * @returns Promise<void>
 */
export async function apiUpdateDiaperRecord(
  recordId: string,
  data: Partial<CreateDiaperRecordRequest>
): Promise<void> {
  await put(`/diaper-records/${recordId}`, data)
}

/**
 * 删除换尿布记录
 *
 * @param recordId 记录ID
 * @returns Promise<void>
 */
export async function apiDeleteDiaperRecord(recordId: string): Promise<void> {
  await del(`/diaper-records/${recordId}`)
}
