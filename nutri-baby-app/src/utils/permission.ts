/**
 * 权限检查工具
 * 职责: 提供权限检查和权限相关的工具函数
 */

import type { MyPermission } from '@/types/collaborator';
import { ROLE_PERMISSIONS } from '@/types/collaborator';

/**
 * 检查用户是否可以查看宝宝
 */
export function canViewBaby(
  permission: MyPermission | undefined
): boolean {
  if (!permission) return false;
  if (isPermissionExpired(permission)) return false;
  // 所有角色都可以查看
  return ['admin', 'editor', 'viewer'].includes(permission.role);
}

/**
 * 检查用户是否可以编辑宝宝记录
 */
export function canEditRecords(
  permission: MyPermission | undefined
): boolean {
  if (!permission) return false;
  if (isPermissionExpired(permission)) return false;
  // admin 和 editor 可以编辑
  return ['admin', 'editor'].includes(permission.role);
}

/**
 * 检查用户是否可以管理宝宝 (修改基本信息、删除)
 */
export function canManageBaby(
  permission: MyPermission | undefined
): boolean {
  if (!permission) return false;
  if (isPermissionExpired(permission)) return false;
  // 仅 admin 可以管理
  return permission.role === 'admin';
}

/**
 * 检查用户是否可以管理协作者
 */
export function canManageCollaborators(
  permission: MyPermission | undefined
): boolean {
  return canManageBaby(permission);
}

/**
 * 检查用户是否可以邀请协作者
 */
export function canInviteCollaborators(
  permission: MyPermission | undefined
): boolean {
  if (!permission) return false;
  if (isPermissionExpired(permission)) return false;
  // admin 和 editor 可以邀请
  return ['admin', 'editor'].includes(permission.role);
}

/**
 * 检查特定权限
 */
export function hasPermission(
  permission: MyPermission | undefined,
  action: string
): boolean {
  if (!permission) return false;
  if (isPermissionExpired(permission)) return false;

  const rolePerms = ROLE_PERMISSIONS[permission.role];
  if (!rolePerms) return false;

  return rolePerms.permissions.includes(action);
}

/**
 * 检查权限是否过期
 */
export function isPermissionExpired(permission: MyPermission): boolean {
  if (permission.accessType === 'permanent') {
    return false;
  }
  if (!permission.expiresAt) {
    return false;
  }
  // 后端返回的是秒级时间戳
  return Date.now() > permission.expiresAt * 1000;
}

/**
 * 获取权限过期警告信息
 */
export function getExpirationWarning(permission: MyPermission): string | null {
  if (permission.accessType === 'permanent') {
    return null;
  }
  if (!permission.expiresAt) {
    return null;
  }

  const now = Date.now();
  const expiresAt = permission.expiresAt * 1000;
  const daysLeft = Math.ceil((expiresAt - now) / (1000 * 60 * 60 * 24));

  if (daysLeft < 0) {
    return '权限已过期，请重新申请';
  }

  if (daysLeft === 0) {
    return '权限即将过期（今天），请尽快续期';
  }

  if (daysLeft <= 7) {
    return `权限将在 ${daysLeft} 天后过期`;
  }

  return null;
}

/**
 * 获取权限标签
 */
export function getRoleLabel(role: string): string {
  const rolePerms = ROLE_PERMISSIONS[role];
  return rolePerms?.label || role;
}

/**
 * 获取权限描述
 */
export function getRoleDescription(role: string): string {
  const rolePerms = ROLE_PERMISSIONS[role];
  return rolePerms?.description || '';
}

/**
 * 获取权限详细列表
 */
export function getRolePermissionsList(role: string): string[] {
  const rolePerms = ROLE_PERMISSIONS[role];
  return rolePerms?.permissions || [];
}

/**
 * 权限列表的中文描述
 */
const PERMISSION_DESCRIPTIONS: Record<string, string> = {
  viewBaby: '查看宝宝基本信息',
  viewRecords: '查看记录',
  viewCollaborators: '查看协作者',
  addRecord: '添加新记录',
  editRecord: '编辑记录',
  deleteRecord: '删除记录',
  editBaby: '编辑宝宝信息',
  deleteBaby: '删除宝宝档案',
  inviteCollaborator: '邀请协作者',
  removeCollaborator: '移除协作者',
  updateCollaboratorRole: '修改协作者权限',
};

export function getPermissionDescription(permission: string): string {
  return PERMISSION_DESCRIPTIONS[permission] || permission;
}

/**
 * 格式化权限显示
 */
export function formatPermissionText(permission: MyPermission): string {
  const roleLabel = getRoleLabel(permission.role);

  if (isPermissionExpired(permission)) {
    return `${roleLabel} (已过期)`;
  }

  if (permission.accessType === 'temporary' && permission.expiresAt) {
    const date = new Date(permission.expiresAt * 1000);
    const dateStr = date.toLocaleDateString('zh-CN');
    return `${roleLabel} · 有效期至 ${dateStr}`;
  }

  return `${roleLabel} · 永久有效`;
}

/**
 * 是否显示权限警告
 */
export function shouldShowWarning(permission: MyPermission): boolean {
  const warning = getExpirationWarning(permission);
  return warning !== null;
}

/**
 * 比较两个权限等级 (admin > editor > viewer)
 */
export function compareRoles(
  role1: string,
  role2: string
): -1 | 0 | 1 {
  const roleOrder = { admin: 3, editor: 2, viewer: 1 };
  const order1 = (roleOrder as any)[role1] || 0;
  const order2 = (roleOrder as any)[role2] || 0;

  if (order1 > order2) return 1;
  if (order1 < order2) return -1;
  return 0;
}
