/**
 * 疫苗管理 API 接口
 * 职责: 纯 API 调用,无状态,无副作用
 */
import { get, post, put, del } from "@/utils/request";

// ============ 类型定义 ============

/**
 * 疫苗接种状态 (新架构)
 */
export type VaccinationStatus = "pending" | "completed" | "skipped";

/**
 * 提醒状态
 */
export type VaccineReminderStatus =
  | "upcoming" // 即将到期
  | "due" // 已到期
  | "overdue" // 已逾期
  | "completed"; // 已完成

export const VaccineReminderStatusMap = {
  upcoming: "即将到期",
  due: "已到期",
  overdue: "已逾期",
  completed: "已完成",
};

/**
 * API 响应: 疫苗接种日程详情 (新架构 - 合并计划和记录)
 */
export interface VaccineScheduleResponse {
  scheduleId: string;
  babyId: string;
  templateId?: string; // 模板ID (自定义计划无此字段)
  vaccineType: string;
  vaccineName: string;
  description?: string;
  ageInMonths: number;
  doseNumber: number;
  isRequired: boolean;
  reminderDays: number;
  isCustom: boolean; // 是否为自定义计划

  // 接种状态
  vaccinationStatus: VaccinationStatus; // pending/completed/skipped

  // 接种记录信息 (status='completed' 时有值)
  vaccineDate?: number; // 接种日期时间戳
  hospital?: string;
  batchNumber?: string;
  doctor?: string;
  reaction?: string;
  note?: string;

  // 完成人信息 (status='completed' 时有值)
  completedBy?: string; // 完成人 openid
  completedByName?: string;
  completedByAvatar?: string;
  completedTime?: number; // 完成时间戳

  // 元数据
  createBy: string;
  createTime: number;
  updateTime?: number;
}

/**
 * API 响应: 疫苗接种日程列表 (新架构)
 */
export interface VaccineScheduleListResponse {
  schedules: VaccineScheduleResponse[];
  statistics: {
    total: number;
    completed: number;
    pending: number;
    skipped: number;
    completionRate: number; // 完成率 (0-100)
  };
}

/**
 * API 请求: 更新疫苗接种日程 (记录接种)
 */
export interface UpdateVaccineScheduleRequest {
  vaccinationStatus: "completed" | "skipped"; // 只能更新为已完成或跳过
  vaccineDate?: number; // status='completed' 时必填
  hospital?: string; // status='completed' 时必填
  batchNumber?: string;
  doctor?: string;
  reaction?: string;
  note?: string;
}

/**
 * API 请求: 创建自定义疫苗接种日程
 */
export interface CreateVaccineScheduleRequest {
  vaccineType: string;
  vaccineName: string;
  description?: string;
  ageInMonths: number;
  doseNumber: number;
  isRequired: boolean;
  reminderDays?: number; // 默认 7 天
}

/**
 * API 响应: 疫苗计划详情
 */
export interface VaccinePlanResponse {
  planId: string;
  vaccineName: string;
  vaccineType: string;
  ageInMonths: number;
  doseNumber: number;
  isRequired: boolean;
  reminderDays: number;
  description?: string;
  isCustom?: boolean; // 是否为自定义计划
  templateId?: string; // 模板ID
}

/**
 * API 响应: 疫苗接种记录详情
 */
export interface VaccineRecordResponse {
  recordId: string;
  babyId: string;
  planId: string;
  vaccineName: string;
  vaccineType: string;
  doseNumber: number;
  vaccineDate: number;
  hospital?: string;
  batchNumber?: string;
  doctor?: string;
  reaction?: string;
  note?: string;
  createBy: string;
  createTime: number;
}

/**
 * API 响应: 疫苗提醒详情
 */
export interface VaccineReminderResponse {
  reminderId: string;
  babyId: string;
  planId: string;
  vaccineName: string;
  doseNumber: number;
  scheduledDate: number;
  status: VaccineReminderStatus;
  isSent: boolean;
  sentTime?: number;
  completedTime?: number;
}

/**
 * API 响应: 疫苗计划列表
 */
export interface VaccinePlansListResponse {
  plans: VaccinePlanResponse[];
  total: number;
  completed?: number; // 已完成接种的计划数
  percentage?: number; // 完成百分比
}

/**
 * API 响应: 疫苗接种记录列表
 */
export interface VaccineRecordsListResponse {
  records: VaccineRecordResponse[];
  total: number;
  page: number;
  pageSize: number;
}

/**
 * API 响应: 疫苗提醒列表
 */
export interface VaccineRemindersListResponse {
  reminders: VaccineReminderResponse[];
  total: number;
}

/**
 * API 请求: 创建疫苗接种记录
 */
export interface CreateVaccineRecordRequest {
  babyId: string;
  planId: string;
  doseNumber: number;
  vaccineName: string;
  vaccineType: string;
  vaccineDate: number;
  hospital?: string;
  batchNumber?: string;
  doctor?: string;
  reaction?: string;
  note?: string;
}

/**
 * API 响应: 疫苗接种统计
 */
export interface VaccineStatsResponse {
  totalPlans: number;
  completedCount: number;
  upcomingCount: number;
  overdueCount: number;
  completionRate: number;
}

// ============ API 函数 ============

/**
 * 获取疫苗计划列表
 *
 * @param babyId 宝宝ID
 * @returns Promise<VaccinePlansListResponse>
 */
export async function apiFetchVaccinePlans(
  babyId: string,
): Promise<VaccinePlansListResponse> {
  const response = await get<VaccinePlansListResponse>(
    `/babies/${babyId}/vaccine-plans`,
  );
  return response.data || { plans: [], total: 0, completed: 0, percentage: 0 };
}

/**
 * 获取宝宝的疫苗接种记录列表
 *
 * ⚠️ 注意: 后端暂未实现该接口,需要添加 GET /v1/babies/:babyId/vaccine-records 路由
 *
 * @param params 查询参数
 * @returns Promise<VaccineRecordsListResponse>
 */
export async function apiFetchVaccineRecords(params: {
  babyId: string;
  page?: number;
  pageSize?: number;
}): Promise<VaccineRecordsListResponse> {
  // TODO: 后端需要实现 GET /v1/babies/:babyId/vaccine-records 接口
  const { babyId, ...queryParams } = params;
  const response = await get<VaccineRecordsListResponse>(
    `/babies/${babyId}/vaccine-records`,
    queryParams,
  );
  return response.data || { records: [], total: 0, page: 1, pageSize: 10 };
}

/**
 * 获取宝宝的疫苗提醒列表
 *
 * @param params 查询参数
 * @returns Promise<VaccineRemindersListResponse>
 */
export async function apiFetchVaccineReminders(params: {
  babyId: string;
}): Promise<VaccineRemindersListResponse> {
  const { babyId, ...queryParams } = params;
  const response = await get<VaccineRemindersListResponse>(
    `/babies/${babyId}/vaccine-reminders`,
    queryParams,
  );
  return response.data || { reminders: [], total: 0 };
}

/**
 * 获取疫苗接种统计
 *
 * @param babyId 宝宝ID
 * @returns Promise<VaccineStatsResponse>
 */
export async function apiFetchVaccineStats(
  babyId: string,
): Promise<VaccineStatsResponse> {
  const response = await get<VaccineStatsResponse>(
    `/babies/${babyId}/vaccine-schedule-statistics`,
  );
  return (
    response.data || {
      totalPlans: 0,
      completedCount: 0,
      upcomingCount: 0,
      overdueCount: 0,
      completionRate: 0,
    }
  );
}

/**
 * 创建疫苗接种记录
 *
 * @param data 创建请求数据
 * @returns Promise<VaccineRecordResponse>
 */
export async function apiCreateVaccineRecord(
  data: CreateVaccineRecordRequest,
): Promise<VaccineRecordResponse> {
  const { babyId, ...recordData } = data;
  const response = await post<VaccineRecordResponse>(
    `/babies/${babyId}/vaccine-records`,
    recordData,
  );
  if (!response.data) {
    throw new Error(response.message || "创建疫苗接种记录失败");
  }
  return response.data;
}

/**
 * 更新疫苗接种记录
 *
 * ⚠️ 注意: 后端暂未实现该接口,需要添加 PUT /v1/vaccine-records/:recordId 路由
 *
 * @param recordId 记录ID
 * @param data 更新数据
 * @returns Promise<void>
 */
export async function apiUpdateVaccineRecord(
  recordId: string,
  data: Partial<CreateVaccineRecordRequest>,
): Promise<void> {
  // TODO: 后端需要实现 PUT /v1/vaccine-records/:recordId 接口
  await put(`/vaccine-records/${recordId}`, data);
}

/**
 * 删除疫苗接种记录
 *
 * ⚠️ 注意: 后端暂未实现该接口,需要添加 DELETE /v1/vaccine-records/:recordId 路由
 *
 * @param recordId 记录ID
 * @returns Promise<void>
 */
export async function apiDeleteVaccineRecord(recordId: string): Promise<void> {
  // TODO: 后端需要实现 DELETE /v1/vaccine-records/:recordId 接口
  await del(`/vaccine-records/${recordId}`);
}

/**
 * 标记疫苗提醒已发送
 *
 * ⚠️ 注意: 后端暂未实现该接口,需要添加 POST /v1/vaccine-reminders/:reminderId/mark-sent 路由
 *
 * @param reminderId 提醒ID
 * @returns Promise<void>
 */
export async function apiMarkReminderSent(reminderId: string): Promise<void> {
  // TODO: 后端需要实现 POST /v1/vaccine-reminders/:reminderId/mark-sent 接口
  await post(`/vaccine-reminders/${reminderId}/mark-sent`, {});
}

// ============ 新架构 API 函数 (合并计划和记录) ============

/**
 * 获取宝宝的疫苗接种日程列表 (新架构)
 *
 * @param babyId 宝宝ID
 * @param params 可选: 查询参数 (status/page/pageSize)
 * @returns Promise<VaccineScheduleListResponse>
 */
export async function apiFetchVaccineSchedules(
  babyId: string,
  params?: {
    status?: VaccinationStatus;
    page?: number;
    pageSize?: number;
  },
): Promise<{
  schedules: VaccineScheduleResponse[];
  total: number;
  page?: number;
  pageSize?: number;
  hasMore?: boolean;
  statistics?: any;
}> {
  // 过滤掉 undefined 的参数
  const queryParams = params ? Object.fromEntries(
    Object.entries(params).filter(([_, v]) => v !== undefined)
  ) : {};

  console.log("apiFetchVaccineSchedules 调用，babyId:", babyId, "queryParams:", queryParams);

  const response = await get<any>(
    `/babies/${babyId}/vaccine-schedules`,
    queryParams,
  );

  console.log("apiFetchVaccineSchedules 响应 response:", response);
  console.log("response.data:", response.data);

  // 后端返回的数据结构: { schedules, total, page, pageSize, hasMore }
  // 需要转换为前端期望的结构
  const result = {
    schedules: response.data?.schedules || [],
    total: response.data?.total || 0,
    page: response.data?.page,
    pageSize: response.data?.pageSize,
    hasMore: response.data?.hasMore,
    statistics: {
      total: 0,
      completed: 0,
      pending: 0,
      skipped: 0,
      completionRate: 0,
    },
  };

  console.log("apiFetchVaccineSchedules 返回结果:", result);

  return result;
}

/**
 * 更新疫苗接种日程 (记录接种或跳过)
 *
 * @param babyId 宝宝ID
 * @param scheduleId 日程ID
 * @param data 更新数据
 * @returns Promise<void>
 */
export async function apiUpdateVaccineSchedule(
  babyId: string,
  scheduleId: string,
  data: UpdateVaccineScheduleRequest,
): Promise<void> {
  await put(`/babies/${babyId}/vaccine-schedules/${scheduleId}`, data);
}

/**
 * 创建自定义疫苗接种日程
 *
 * @param babyId 宝宝ID
 * @param data 自定义计划数据
 * @returns Promise<void>
 */
export async function apiCreateCustomSchedule(
  babyId: string,
  data: CreateVaccineScheduleRequest,
): Promise<void> {
  await post(`/babies/${babyId}/vaccine-schedules`, data);
}

/**
 * 删除疫苗接种日程 (仅限自定义)
 *
 * @param babyId 宝宝ID
 * @param scheduleId 日程ID
 * @returns Promise<void>
 */
export async function apiDeleteVaccineSchedule(
  babyId: string,
  scheduleId: string,
): Promise<void> {
  await del(`/babies/${babyId}/vaccine-schedules/${scheduleId}`);
}

/**
 * 更新疫苗接种日程基本信息 (仅限未完成的日程)
 *
 * @param babyId 宝宝ID
 * @param scheduleId 日程ID
 * @param data 更新数据
 * @returns Promise<void>
 */
export async function apiUpdateScheduleInfo(
  babyId: string,
  scheduleId: string,
  data: Partial<CreateVaccineScheduleRequest>,
): Promise<void> {
  await put(`/babies/${babyId}/vaccine-schedules/${scheduleId}/info`, data);
}

/**
 * 获取疫苗接种统计 (新架构)
 *
 * @param babyId 宝宝ID
 * @returns Promise<VaccineScheduleStatistics>
 */
export async function apiFetchVaccineScheduleStatistics(
  babyId: string,
): Promise<VaccineScheduleListResponse["statistics"]> {
  const response = await get<VaccineScheduleListResponse["statistics"]>(
    `/babies/${babyId}/vaccine-schedule-statistics`,
  );
  return (
    response.data || {
      total: 0,
      completed: 0,
      pending: 0,
      skipped: 0,
      completionRate: 0,
    }
  );
}
