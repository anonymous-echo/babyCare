/**
 * 协作者类型定义
 * 职责: 定义所有与协作者相关的类型
 */

/**
 * 宝宝协作者信息 (来自后端API)
 */
export interface BabyCollaborator {
  openid: string;
  nickName: string;
  avatarUrl?: string;
  role: 'admin' | 'editor' | 'viewer';
  accessType: 'permanent' | 'temporary';
  expiresAt?: number;
  joinedAt: number;
  isCreator: boolean;
  relationship?: string; // 与宝宝的关系，如：爸爸、妈妈、祖母等
}

/**
 * 当前用户对宝宝的权限信息
 */
export interface MyPermission {
  babyId: string;
  role: 'admin' | 'editor' | 'viewer';
  accessType: 'permanent' | 'temporary';
  expiresAt?: number;
  joinedAt: number;
}

/**
 * 角色权限定义
 */
export interface RolePermissions {
  role: 'admin' | 'editor' | 'viewer';
  label: string;
  description: string;
  permissions: string[];
}

/**
 * 邀请请求
 */
export interface InviteCollaboratorRequest {
  targetOpenid?: string;
  targetPhone?: string;
  role: 'admin' | 'editor' | 'viewer';
  expiresAt?: number;
  accessType: 'permanent' | 'temporary';
  relationship?: string; // 与宝宝的关系
}

/**
 * 邀请响应
 */
export interface InviteCollaboratorResponse {
  babyId: string;
  shortCode: string;
  token: string;
  qrCodeUrl: string;
  expiresAt?: number;
}

/**
 * 权限枚举
 */
export const ROLE_PERMISSIONS: Record<string, RolePermissions> = {
  admin: {
    role: 'admin',
    label: '管理员',
    description: '拥有完全控制权',
    permissions: [
      'viewBaby',
      'viewRecords',
      'viewCollaborators',
      'addRecord',
      'editRecord',
      'deleteRecord',
      'editBaby',
      'deleteBaby',
      'inviteCollaborator',
      'removeCollaborator',
      'updateCollaboratorRole',
    ],
  },
  editor: {
    role: 'editor',
    label: '编辑者',
    description: '可查看和编辑记录',
    permissions: [
      'viewBaby',
      'viewRecords',
      'viewCollaborators',
      'addRecord',
      'editRecord',
      'deleteRecord',
      'inviteCollaborator',
    ],
  },
  viewer: {
    role: 'viewer',
    label: '查看者',
    description: '仅可查看数据',
    permissions: ['viewBaby', 'viewRecords', 'viewCollaborators'],
  },
};

/**
 * 权限访问类型
 */
export interface AccessType {
  type: 'permanent' | 'temporary';
  label: string;
}

export const ACCESS_TYPES: Record<string, AccessType> = {
  permanent: {
    type: 'permanent',
    label: '永久有效',
  },
  temporary: {
    type: 'temporary',
    label: '有期限',
  },
};

/**
 * 计算权限是否过期
 */
export function isPermissionExpired(permission: MyPermission): boolean {
  if (permission.accessType === 'permanent') {
    return false;
  }
  if (!permission.expiresAt) {
    return false;
  }
  return Date.now() > permission.expiresAt * 1000;
}

/**
 * 检查是否有特定权限
 */
export function hasPermission(
  permission: MyPermission | undefined,
  action: string
): boolean {
  if (!permission) return false;
  if (isPermissionExpired(permission)) return false;

  const rolePermissions = ROLE_PERMISSIONS[permission.role];
  if (!rolePermissions) return false;

  return rolePermissions.permissions.includes(action);
}

/**
 * 获取权限显示文本
 */
export function getPermissionText(permission: MyPermission): string {
  if (isPermissionExpired(permission)) {
    return '权限已过期';
  }

  const roleLabel = ROLE_PERMISSIONS[permission.role]?.label || permission.role;

  if (permission.accessType === 'temporary' && permission.expiresAt) {
    const expiresDate = new Date(permission.expiresAt * 1000);
    return `${roleLabel} · 有效期至 ${expiresDate.toLocaleDateString('zh-CN')}`;
  }

  return `${roleLabel} · 永久有效`;
}

/**
 * 格式化协作者加入时间
 */
export function formatJoinedDate(timestamp: number): string {
  const date = new Date(timestamp * 1000);
  return date.toLocaleDateString('zh-CN');
}
