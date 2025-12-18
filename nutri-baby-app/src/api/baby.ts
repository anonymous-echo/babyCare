/**
 * 宝宝管理 API 接口
 * 职责: 纯 API 调用,无状态,无副作用
 */
import { get, post, put, del } from "@/utils/request";

// ============ 类型定义 ============

/**
 * API 响应: 宝宝详情
 */
export interface BabyResponse {
  babyId: string;
  name: string;
  nickname?: string;
  gender: "male" | "female";
  birthDate: string;
  avatarUrl?: string;
  creatorId: string;
  familyGroup?: string;
  height?: number;
  weight?: number;
  createTime: number;
  updateTime: number;
}

/**
 * 前端显示: 宝宝档案 (转换后的字段名)
 */
export interface BabyProfileResponse {
  babyId: string;
  name: string;
  nickname?: string;
  gender: "male" | "female";
  birthDate: string;
  avatarUrl?: string;
  creatorId: string;
  familyGroup?: string;
  createTime: number;
  updateTime: number;
}

/**
 * API 请求: 创建宝宝
 */
export interface CreateBabyRequest {
  name: string;
  gender: "male" | "female";
  birthDate: string;
  nickname?: string;
  avatarUrl?: string;
  familyGroup?: string;
  copyCollaboratorsFrom?: string;
}

/**
 * API 请求: 更新宝宝
 */
export interface UpdateBabyRequest {
  name?: string;
  nickname?: string;
  gender?: "male" | "female";
  birthDate?: string;
  avatarUrl?: string;
  familyGroup?: string;
  height?: number;
  weight?: number;
}

/**
 * API 响应: 生成邀请二维码
 */
export interface QRCodeResponse {
  qrcodeUrl: string;
  scene: string;
}

/**
 * API 响应: 邀请详情
 */
export interface InvitationDetailResponse {
  babyId: string;
  babyName: string;
  babyAvatar?: string;
  inviterName: string;
  role: string;
  accessType: string;
  expiresAt?: number;
  validUntil: number;
  token: string;
}

// ============ API 函数 ============

/**
 * 获取用户可访问的宝宝列表
 *
 * @returns Promise<BabyProfileResponse[]> - 返回转换后的宝宝列表 (babyName -> name)
 */
export async function apiFetchBabyList(): Promise<BabyProfileResponse[]> {
  const response = await get<BabyResponse[]>("/babies");
  const babies = response.data || [];

  return babies.map((baby) => ({
    babyId: baby.babyId,
    name: baby.name,
    nickname: baby.nickname,
    gender: baby.gender,
    birthDate: baby.birthDate,
    avatarUrl: baby.avatarUrl,
    creatorId: baby.creatorId,
    familyGroup: baby.familyGroup,
    createTime: baby.createTime,
    updateTime: baby.updateTime,
  }));
}

/**
 * 获取宝宝详情
 *
 * @param babyId 宝宝ID
 * @returns Promise<BabyResponse>
 */
export async function apiFetchBabyDetail(
  babyId: string,
): Promise<BabyResponse> {
  const response = await get<BabyResponse>(`/babies/${babyId}`);
  if (!response.data) {
    throw new Error(response.message || "获取宝宝详情失败");
  }
  return response.data;
}

/**
 * 创建宝宝
 *
 * @param data 创建请求数据
 * @returns Promise<BabyResponse>
 */
export async function apiCreateBaby(
  data: CreateBabyRequest,
): Promise<BabyResponse> {
  const response = await post<BabyResponse>("/babies", data);
  if (response.code !== 0) {
    throw new Error(response.message || "添加宝宝失败");
  }
  return response.data;
}

/**
 * 更新宝宝信息
 *
 * @param babyId 宝宝ID
 * @param data 更新数据
 * @returns Promise<void>
 */
export async function apiUpdateBaby(
  babyId: string,
  data: UpdateBabyRequest,
): Promise<void> {
  const response = await put(`/babies/${babyId}`, data);
  if (response.code !== 0) {
    throw new Error(response.message || "更新宝宝信息失败");
  }
}

/**
 * 删除宝宝
 *
 * @param babyId 宝宝ID
 * @returns Promise<void>
 */
export async function apiDeleteBaby(babyId: string): Promise<void> {
  const response = await del(`/babies/${babyId}`);
  if (response.code !== 0) {
    throw new Error(response.message || "删除宝宝失败");
  }
}

/**
 * 生成邀请二维码
 *
 * @param babyId 宝宝ID
 * @param shortCode 邀请短码
 * @returns Promise<QRCodeResponse>
 */
export async function apiGenerateQRCode(
  babyId: string,
  shortCode: string,
): Promise<QRCodeResponse> {
  const response = await get<QRCodeResponse>(
    `/babies/${babyId}/qrcode?shortCode=${shortCode}`,
  );
  if (!response.data) {
    throw new Error(response.message || "生成二维码失败");
  }
  return response.data;
}

/**
 * 通过短码获取邀请详情
 *
 * @param shortCode 邀请短码
 * @returns Promise<InvitationDetailResponse>
 */
export async function apiGetInvitationByCode(
  shortCode: string,
): Promise<InvitationDetailResponse> {
  const response = await get<InvitationDetailResponse>(
    `/invitations/code/${shortCode}`,
  );
  if (!response.data) {
    throw new Error(response.message || "获取邀请详情失败");
  }
  return response.data;
}
