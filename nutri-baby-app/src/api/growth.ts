/**
 * 成长记录 API 接口
 * 职责: 纯 API 调用,无状态,无副作用
 */
import { get, post, put, del } from "@/utils/request";

// ============ 类型定义 ============

/**
 * API 响应: 成长记录详情
 */
export interface GrowthRecordResponse {
  recordId: string;
  babyId: string;
  height?: number; // 仅当有值时返回
  weight?: number; // 仅当有值时返回
  headCircumference?: number; // 仅当有值时返回
  note?: string;
  measureTime: number;
  createBy: string;
  createTime: number;
}

/**
 * API 响应: 成长记录列表
 */
export interface GrowthRecordsListResponse {
  records: GrowthRecordResponse[];
  total: number;
  page: number;
  pageSize: number;
}

/**
 * API 请求: 创建成长记录
 */
export interface CreateGrowthRecordRequest {
  babyId: string;
  height?: number;
  weight?: number;
  headCircumference?: number;
  note?: string;
  measureTime: number;
}

// ============ API 函数 ============

/**
 * 获取成长记录列表
 *
 * @param params 查询参数
 * @returns Promise<GrowthRecordsListResponse>
 */
export async function apiFetchGrowthRecords(params: {
  babyId: string;
  startTime?: number;
  endTime?: number;
  page?: number;
  pageSize?: number;
}): Promise<GrowthRecordsListResponse> {
  const response = await get<GrowthRecordsListResponse>(
    "/growth-records",
    params,
  );
  return response.data || { records: [], total: 0, page: 1, pageSize: 10 };
}

/**
 * 创建成长记录
 *
 * @param data 创建请求数据
 * @returns Promise<GrowthRecordResponse>
 */
export async function apiCreateGrowthRecord(
  data: CreateGrowthRecordRequest,
): Promise<GrowthRecordResponse> {
  const response = await post<GrowthRecordResponse>("/growth-records", data);
  if (!response.data) {
    throw new Error(response.message || "创建成长记录失败");
  }
  return response.data;
}

/**
 * 获取单条成长记录
 *
 * @param recordId 记录ID
 * @returns Promise<GrowthRecordResponse>
 */
export async function apiGetGrowthRecordById(
  recordId: string,
): Promise<GrowthRecordResponse> {
  const response = await get<GrowthRecordResponse>(`/growth-records/${recordId}`);
  if (!response.data) {
    throw new Error(response.message || "获取成长记录失败");
  }
  return response.data;
}

/**
 * 更新成长记录
 *
 * @param recordId 记录ID
 * @param data 更新数据
 * @returns Promise<void>
 */
export async function apiUpdateGrowthRecord(
  recordId: string,
  data: Partial<CreateGrowthRecordRequest>,
): Promise<void> {
  await put(`/growth-records/${recordId}`, data);
}

/**
 * 删除成长记录
 *
 * @param recordId 记录ID
 * @returns Promise<void>
 */
export async function apiDeleteGrowthRecord(recordId: string): Promise<void> {
  await del(`/growth-records/${recordId}`);
}
