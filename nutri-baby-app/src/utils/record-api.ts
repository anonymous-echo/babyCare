/**
 * 记录相关 API 服务模块
 * 提供统一的记录 CRUD 操作
 */
import { get, post } from './request'
import type { ApiResponse } from '@/types'

/**
 * 通用的分页查询参数接口
 */
export interface RecordQueryParams {
  babyId: string
  startTime?: number
  endTime?: number
  page?: number
  pageSize?: number
}

/**
 * 通用的分页响应接口
 */
export interface PagedResponse<T> {
  records: T[]
  total: number
  page: number
  pageSize: number
}

/**
 * 创建记录
 */
export async function createRecord<T>(
  endpoint: string,
  data: any
): Promise<T> {
  const response = await post<T>(endpoint, data)

  if (response.code === 0 && response.data) {
    return response.data
  } else {
    throw new Error(response.message || '创建记录失败')
  }
}

/**
 * 获取记录列表
 */
export async function fetchRecords<T>(
  endpoint: string,
  params: RecordQueryParams
): Promise<PagedResponse<T>> {
  const response = await get<PagedResponse<T>>(endpoint, params)

  if (response.code === 0 && response.data) {
    return response.data
  } else {
    throw new Error(response.message || '获取记录列表失败')
  }
}
