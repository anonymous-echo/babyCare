/**
 * Store 统一导出
 *
 * 重构后的 Store 架构:
 * - 仅保留全局状态: user, baby
 * - 记录相关数据 (feeding/diaper/sleep/growth/vaccine) 已移除
 * - 页面组件应直接调用 @/api/* 层获取数据
 */

// 用户状态管理 (保留)
export * from './user'

// 宝宝状态管理 (简化版,保留)
// 注: 协作者数据已集成到 baby.ts 中，使用 getCollaborators/setCollaborators/getMyPermission/setMyPermission
export * from './baby'

// 其他模块 (根据需要保留)
// export * from './collaborator'  // 已整合到 baby.ts 中，避免导出冲突
export * from './subscribe'
