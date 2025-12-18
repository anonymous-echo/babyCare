/**
 * 宝宝亲友团管理状态 (去家庭化架构)
 *
 * 核心逻辑:
 * - 每个宝宝有独立的亲友团列表
 * - 支持通过微信分享和二维码邀请亲友
 * - 亲友拥有不同的角色权限: admin / editor / viewer
 * - 每个亲友可以设置与宝宝的关系: 爸爸、妈妈、爷爷、奶奶等
 * - 支持永久和临时访问权限
 */
import { ref } from 'vue'
import type { BabyCollaborator } from '@/types'
import { get, post, put, del } from '@/utils/request'

// 当前宝宝的协作者列表 (按需加载)
const collaborators = ref<Map<string, BabyCollaborator[]>>(new Map())

/**
 * 获取宝宝的协作者列表
 *
 * API: GET /babies/{babyId}/collaborators
 * 响应: [ { openid, nickName, avatarUrl, role, accessType, expiresAt?, joinTime } ]
 */
export async function fetchCollaborators(babyId: string): Promise<BabyCollaborator[]> {
  try {
    const response = await get<BabyCollaborator[]>(`/babies/${babyId}/collaborators`)

    if (response.code === 0 && response.data) {
      const collaboratorList = response.data
      collaborators.value.set(babyId, collaboratorList)
      return collaboratorList
    } else {
      throw new Error(response.message || '获取协作者列表失败')
    }
  } catch (error: any) {
    console.error('fetch collaborators error:', error)
    throw error
  }
}

/**
 * 邀请协作者 (生成微信分享或二维码)
 *
 * API: POST /babies/{babyId}/collaborators/invite
 * 请求: { inviteType: 'share' | 'qrcode', role, accessType, expiresAt? }
 * 响应: {
 *   babyId, babyName, inviterName, role,
 *   shareParams?: { title, path, imageUrl },
 *   qrcodeParams?: { scene, page, qrcodeUrl },
 *   expiresAt?, validUntil
 * }
 */
export async function inviteCollaborator(
  babyId: string,
  inviteType: 'share' | 'qrcode',
  role: 'admin' | 'editor' | 'viewer',
  accessType: 'permanent' | 'temporary',
  expiresAt?: number,
  relationship?: string
): Promise<any> {
  try {
    const requestData: any = {
      inviteType,
      role,
      accessType,
    }

    if (accessType === 'temporary' && expiresAt) {
      requestData.expiresAt = expiresAt
    }

    if (relationship) {
      requestData.relationship = relationship
    }

    const response = await post<any>(
      `/babies/${babyId}/collaborators/invite`,
      requestData
    )

    if (response.code === 0 && response.data) {
      uni.showToast({
        title: '邀请生成成功',
        icon: 'success',
      })

      return response.data
    } else {
      throw new Error(response.message || '生成邀请失败')
    }
  } catch (error: any) {
    console.error('invite collaborator error:', error)
    uni.showToast({
      title: error.message || '生成邀请失败',
      icon: 'none',
    })
    throw error
  }
}

/**
 * 通过邀请加入宝宝协作
 *
 * API: POST /babies/join
 * 请求: { babyId, token }
 * 响应: { code, message, data: { babyId, babyName, role } }
 */
export async function joinBabyCollaboration(
  babyId: string,
  token: string
): Promise<{ babyId: string; name: string; role: string }> {
  try {
    const response = await post<any>('/babies/join', {
      babyId,
      token,
    })

    if (response.code === 0 && response.data) {
      uni.showToast({
        title: '加入成功',
        icon: 'success',
      })

      return response.data
    } else {
      throw new Error(response.message || '加入失败')
    }
  } catch (error: any) {
    console.error('join baby collaboration error:', error)
    uni.showToast({
      title: error.message || '加入失败',
      icon: 'none',
    })
    throw error
  }
}

/**
 * 移除协作者 (仅 admin 可操作)
 *
 * API: DELETE /babies/{babyId}/collaborators/{openid}
 * 响应: { code, message, data: null }
 */
export async function removeCollaborator(
  babyId: string,
  openid: string
): Promise<boolean> {
  try {
    const response = await del(`/babies/${babyId}/collaborators/${openid}`)

    if (response.code === 0) {
      // 从本地缓存中移除
      const list = collaborators.value.get(babyId)
      if (list) {
        const index = list.findIndex(c => c.openid === openid)
        if (index !== -1) {
          list.splice(index, 1)
          collaborators.value.set(babyId, [...list])
        }
      }

      uni.showToast({
        title: '移除成功',
        icon: 'success',
      })

      return true
    } else {
      throw new Error(response.message || '移除协作者失败')
    }
  } catch (error: any) {
    console.error('remove collaborator error:', error)
    uni.showToast({
      title: error.message || '移除失败',
      icon: 'none',
    })
    throw error
  }
}

/**
 * 更新协作者角色 (仅 admin 可操作)
 *
 * API: PUT /babies/{babyId}/collaborators/{openid}/role
 * 请求: { role: 'admin' | 'editor' | 'viewer' }
 * 响应: { code, message, data: null }
 */
export async function updateCollaboratorRole(
  babyId: string,
  openid: string,
  role: 'admin' | 'editor' | 'viewer'
): Promise<boolean> {
  try {
    const response = await put(
      `/babies/${babyId}/collaborators/${openid}/role`,
      { role }
    )

    if (response.code === 0) {
      // 更新本地缓存
      const list = collaborators.value.get(babyId)
      if (list) {
        const collaborator = list.find(c => c.openid === openid)
        if (collaborator) {
          collaborator.role = role
          collaborators.value.set(babyId, [...list])
        }
      }

      uni.showToast({
        title: '角色更新成功',
        icon: 'success',
      })

      return true
    } else {
      throw new Error(response.message || '更新角色失败')
    }
  } catch (error: any) {
    console.error('update collaborator role error:', error)
    uni.showToast({
      title: error.message || '更新失败',
      icon: 'none',
    })
    throw error
  }
}

/**
 * 更新亲友团成员信息
 * - 管理员可以更新任何成员的角色和关系
 * - 普通用户可以更新自己的关系
 *
 * API: PUT /babies/{babyId}/collaborators/{openid}
 * 请求: { role?: 'admin' | 'editor' | 'viewer', relationship?: string }
 * 响应: { code, message, data: null }
 */
export async function updateFamilyMember(
  babyId: string,
  openid: string,
  data: { role?: 'admin' | 'editor' | 'viewer'; relationship?: string }
): Promise<boolean> {
  try {
    const response = await put(
      `/babies/${babyId}/collaborators/${openid}`,
      data
    )

    if (response.code === 0) {
      // 更新本地缓存
      const list = collaborators.value.get(babyId)
      if (list) {
        const collaborator = list.find(c => c.openid === openid)
        if (collaborator) {
          if (data.role) {
            collaborator.role = data.role
          }
          if (data.relationship !== undefined) {
            collaborator.relationship = data.relationship
          }
          collaborators.value.set(babyId, [...list])
        }
      }

      uni.showToast({
        title: '更新成功',
        icon: 'success',
      })

      return true
    } else {
      throw new Error(response.message || '更新失败')
    }
  } catch (error: any) {
    console.error('update family member error:', error)
    uni.showToast({
      title: error.message || '更新失败',
      icon: 'none',
    })
    throw error
  }
}

/**
 * 获取本地缓存的协作者列表
 */
export function getCollaborators(babyId: string): BabyCollaborator[] {
  return collaborators.value.get(babyId) || []
}

/**
 * 清除协作者数据 (用于登出或切换宝宝)
 */
export function clearCollaborators(babyId?: string) {
  if (babyId) {
    collaborators.value.delete(babyId)
  } else {
    collaborators.value.clear()
  }
}

export { collaborators }
