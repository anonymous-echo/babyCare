/**
 * 协作者相关 API 接口
 * 职责: 所有与协作者相关的API调用
 */

import { get, post, put, del } from '@/utils/request';
import type {
  BabyCollaborator,
  InviteCollaboratorRequest,
  InviteCollaboratorResponse,
} from '@/types/collaborator';

/**
 * 获取宝宝的协作者列表
 *
 * @param babyId - 宝宝ID
 * @returns 协作者列表
 */
export async function apiFetchCollaborators(
  babyId: string
): Promise<BabyCollaborator[]> {
  try {
    const response = await get<BabyCollaborator[]>(
      `/babies/${babyId}/collaborators`
    );
    return response.data || [];
  } catch (error) {
    console.error('[Collaborator] 获取协作者列表失败:', error);
    throw error;
  }
}

/**
 * 邀请协作者
 *
 * @param babyId - 宝宝ID
 * @param data - 邀请数据
 * @returns 邀请响应（包含二维码、token等）
 */
export async function apiInviteCollaborator(
  babyId: string,
  data: InviteCollaboratorRequest
): Promise<InviteCollaboratorResponse> {
  try {
    const response = await post<InviteCollaboratorResponse>(
      `/babies/${babyId}/collaborators/invite`,
      data
    );

    if (response.code !== 0) {
      throw new Error(response.message || '邀请失败');
    }

    return response.data;
  } catch (error) {
    console.error('[Collaborator] 邀请协作者失败:', error);
    throw error;
  }
}

/**
 * 移除协作者
 *
 * @param babyId - 宝宝ID
 * @param openid - 协作者openid
 */
export async function apiRemoveCollaborator(
  babyId: string,
  openid: string
): Promise<void> {
  try {
    const response = await del(
      `/babies/${babyId}/collaborators/${openid}`
    );

    if (response.code !== 0) {
      throw new Error(response.message || '移除失败');
    }
  } catch (error) {
    console.error('[Collaborator] 移除协作者失败:', error);
    throw error;
  }
}

/**
 * 更新协作者角色和权限
 *
 * @param babyId - 宝宝ID
 * @param openid - 协作者openid
 * @param role - 新的角色
 * @param expiresAt - 权限过期时间（可选）
 */
export async function apiUpdateCollaboratorRole(
  babyId: string,
  openid: string,
  role: 'admin' | 'editor' | 'viewer',
  expiresAt?: number
): Promise<void> {
  try {
    const payload: any = { role };
    if (expiresAt !== undefined) {
      payload.expiresAt = expiresAt;
    }

    const response = await put(
      `/babies/${babyId}/collaborators/${openid}/role`,
      payload
    );

    if (response.code !== 0) {
      throw new Error(response.message || '更新失败');
    }
  } catch (error) {
    console.error('[Collaborator] 更新协作者角色失败:', error);
    throw error;
  }
}

/**
 * 通过邀请短码获取邀请详情
 *
 * @param shortCode - 邀请短码
 * @returns 邀请详情
 */
export async function apiGetInvitationByCode(shortCode: string) {
  try {
    const response = await get(`/invitations/code/${shortCode}`);

    if (response.code !== 0) {
      throw new Error(response.message || '获取邀请详情失败');
    }

    return response.data;
  } catch (error) {
    console.error('[Collaborator] 获取邀请详情失败:', error);
    throw error;
  }
}

/**
 * 确认加入宝宝（接受邀请）
 *
 * @param babyId - 宝宝ID
 * @param token - 邀请token
 */
export async function apiJoinBaby(
  babyId: string,
  token: string
): Promise<any> {
  try {
    const response = await post('/babies/join', {
      babyId,
      token,
    });

    if (response.code !== 0) {
      throw new Error(response.message || '加入失败');
    }

    return response.data;
  } catch (error) {
    console.error('[Collaborator] 加入宝宝失败:', error);
    throw error;
  }
}

/**
 * 批量邀请协作者 (未来扩展)
 *
 * @param babyId - 宝宝ID
 * @param invitations - 邀请列表
 */
export async function apiBatchInviteCollaborators(
  babyId: string,
  invitations: InviteCollaboratorRequest[]
): Promise<InviteCollaboratorResponse[]> {
  try {
    const response = await post<InviteCollaboratorResponse[]>(
      `/babies/${babyId}/collaborators/batch-invite`,
      { invitations }
    );

    if (response.code !== 0) {
      throw new Error(response.message || '批量邀请失败');
    }

    return response.data || [];
  } catch (error) {
    console.error('[Collaborator] 批量邀请失败:', error);
    throw error;
  }
}
