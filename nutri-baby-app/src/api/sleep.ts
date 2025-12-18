/**
 * 睡眠记录 API 接口
 * 职责: 纯 API 调用,无状态,无副作用
 */
import { get, post, put, del } from "@/utils/request";

// ============ 类型定义 ============

/**
 * API 响应: 睡眠记录详情
 */
export interface SleepRecordResponse {
  recordId: string;
  babyId: string;
  sleepType: "nap" | "night";  // 睡眠类型
  startTime: number;
  endTime?: number;
  duration?: number;  // 秒数
  note?: string;
  createBy: string;
  createTime: number;
}

/**
 * API 响应: 睡眠记录列表
 */
export interface SleepRecordsListResponse {
  records: SleepRecordResponse[];
  total: number;
  page: number;
  pageSize: number;
}

/**
 * API 请求: 创建睡眠记录
 */
export interface CreateSleepRecordRequest {
  babyId: string;
  sleepType: "nap" | "night";
  startTime: number;
  endTime?: number;
  duration?: number; // 时长(秒)
  quality?: "good" | "fair" | "poor";
  note?: string;
}

/**
 * 将API响应转换为前端SleepRecord类型
 * @param response API响应
 * @returns 前端SleepRecord
 */
export function transformSleepRecordResponse(response: SleepRecordResponse): any {
  return {
    id: response.recordId,
    babyId: response.babyId,
    startTime: response.startTime,
    endTime: response.endTime,
    duration: response.duration ? Math.round(response.duration / 60) : undefined, // 秒转分钟
    type: response.sleepType, // 直接使用sleepType字段
    createBy: response.createBy,
    createByName: '',
    createByAvatar: '',
    createTime: response.createTime,
  };
}

// ============ API 函数 ============

/**
 * 获取睡眠记录列表
 *
 * @param params 查询参数
 * @returns Promise<SleepRecordsListResponse>
 */
export async function apiFetchSleepRecords(params: {
  babyId: string;
  startTime?: number;
  endTime?: number;
  page?: number;
  pageSize?: number;
}): Promise<SleepRecordsListResponse> {
  const response = await get<SleepRecordsListResponse>(
    "/sleep-records",
    params,
  );
  return response.data || { records: [], total: 0, page: 1, pageSize: 10 };
}

/**
 * 创建睡眠记录
 *
 * @param data 创建请求数据
 * @returns Promise<SleepRecordResponse>
 */
export async function apiCreateSleepRecord(
  data: CreateSleepRecordRequest,
): Promise<SleepRecordResponse> {
  const response = await post<SleepRecordResponse>("/sleep-records", data);
  if (!response.data) {
    throw new Error(response.message || "创建睡眠记录失败");
  }
  return response.data;
}

/**
 * 获取单条睡眠记录
 *
 * @param recordId 记录ID
 * @returns Promise<SleepRecordResponse>
 */
export async function apiGetSleepRecordById(
  recordId: string,
): Promise<SleepRecordResponse> {
  const response = await get<SleepRecordResponse>(`/sleep-records/${recordId}`);
  if (!response.data) {
    throw new Error(response.message || "获取睡眠记录失败");
  }
  return response.data;
}

/**
 * 更新睡眠记录
 *
 * @param recordId 记录ID
 * @param data 更新数据
 * @returns Promise<void>
 */
export async function apiUpdateSleepRecord(
  recordId: string,
  data: Partial<CreateSleepRecordRequest>,
): Promise<void> {
  await put(`/sleep-records/${recordId}`, data);
}

/**
 * 删除睡眠记录
 *
 * @param recordId 记录ID
 * @returns Promise<void>
 */
export async function apiDeleteSleepRecord(recordId: string): Promise<void> {
  await del(`/sleep-records/${recordId}`);
}
