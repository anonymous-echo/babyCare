/**
 * 宝宝数据状态管理
 * 职责: 状态管理 + 本地计算,API 调用委托给 api 层
 *
 * ⚠️ 向后兼容: 所有导出函数的签名保持不变,页面组件无需修改
 */
import { ref, computed } from "vue";
import type { BabyProfile } from "@/types";
import type { BabyCollaborator, MyPermission } from "@/types/collaborator";
import { StorageKeys, getStorage, setStorage } from "@/utils/storage";
import * as babyApi from "@/api/baby";
import { getUserInfo } from "./user";

// ============ 状态定义 ============

// 宝宝列表 - 延迟初始化
const babyList = ref<BabyProfile[]>([]);

// 当前选中的宝宝 ID - 延迟初始化
const currentBabyId = ref<string>("");

// 初始化标记
let isInitialized = false;

// 延迟初始化 - 仅在首次访问时从存储读取
function initializeIfNeeded() {
  if (!isInitialized) {
    babyList.value = getStorage<BabyProfile[]>(StorageKeys.BABY_LIST) || [];
    currentBabyId.value = getStorage<string>(StorageKeys.CURRENT_BABY_ID) || "";
    isInitialized = true;
  }
}

// 当前宝宝信息
const currentBaby = computed(() => {
  initializeIfNeeded(); // 确保数据已加载
  const baby = babyList.value.find((baby) => baby.babyId === currentBabyId.value) || null;
  return baby;
});

// ============ Collaborator 相关状态 ============

// 宝宝的协作者列表 Map<babyId, collaborators[]>
const collaboratorsMap = ref<Map<string, BabyCollaborator[]>>(new Map());

// 当前用户对每个宝宝的权限 Map<babyId, MyPermission>
const myPermissionsMap = ref<Map<string, MyPermission>>(new Map());

// ============ 本地查询函数 ============

/**
 * 获取宝宝列表(本地)
 *
 * ⚠️ 向后兼容: 函数签名保持不变
 */
export function getBabyList() {
  initializeIfNeeded();
  return babyList.value;
}

/**
 * 获取当前宝宝(本地)
 *
 * ⚠️ 向后兼容: 函数签名保持不变
 */
export function getCurrentBaby() {
  return currentBaby.value;
}

// ============ API 调用函数(委托给 api 层) ============

/**
 * 获取用户可访问的宝宝列表 (去家庭化架构)
 *
 * API: GET /babies
 *
 * ⚠️ 向后兼容: 函数签名保持不变
 */
export async function fetchBabyList(): Promise<BabyProfile[]> {
  initializeIfNeeded();
  try {
    const apiResponse = await babyApi.apiFetchBabyList();

    // 将 API 响应的字段映射到本地类型
    const babies: BabyProfile[] = apiResponse.map((baby) => ({
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

    babyList.value = babies;
    setStorage(StorageKeys.BABY_LIST, babies);

    // 设置当前宝宝的逻辑优化：
    // 1. 如果已经有选中的宝宝且该宝宝仍在列表中，保持不变（用户手动选择优先）
    // 2. 如果没有选中宝宝或选中的宝宝不在列表中，则自动选择：
    //    a. 优先使用用户设置的默认宝宝
    //    b. 否则选择列表中的第一个宝宝
    const userInfo = getUserInfo();
    const defaultBabyId = userInfo?.defaultBabyId;

    if (babies.length > 0) {
      // 检查当前选中的宝宝是否仍在列表中
      const currentBabyStillExists = currentBabyId.value && babies.some((b) => b.babyId === currentBabyId.value);
      
      if (!currentBabyStillExists) {
        // 当前宝宝不存在，需要重新选择
        if (defaultBabyId && babies.some((b) => b.babyId === defaultBabyId)) {
          setCurrentBaby(defaultBabyId);
        } else {
          const firstBaby = babies[0];
          if (firstBaby) {
            setCurrentBaby(firstBaby.babyId);
          }
        }
      }
      // 如果当前宝宝仍然存在，保持不变
    } else {
      // 没有宝宝，清空选择
      setCurrentBaby("");
    }

    return babies;
  } catch (error: any) {
    console.error("fetch baby list error:", error);
    throw error;
  }
}

/**
 * 获取宝宝详情 (去家庭化架构)
 *
 * API: GET /babies/{babyId}
 *
 * ⚠️ 向后兼容: 函数签名保持不变
 */
export async function fetchBabyDetail(babyId: string): Promise<BabyProfile> {
  try {
    const apiResponse = await babyApi.apiFetchBabyDetail(babyId);

    // 映射字段
    const baby: BabyProfile = {
      babyId: apiResponse.babyId,
      name: apiResponse.name,
      nickname: apiResponse.nickname,
      gender: apiResponse.gender,
      birthDate: apiResponse.birthDate,
      avatarUrl: apiResponse.avatarUrl,
      creatorId: apiResponse.creatorId,
      familyGroup: apiResponse.familyGroup,
      createTime: apiResponse.createTime,
      updateTime: apiResponse.updateTime,
    };

    // 更新本地列表
    const index = babyList.value.findIndex((b) => b.babyId === baby.babyId);
    if (index !== -1) {
      babyList.value[index] = baby;
    } else {
      babyList.value.push(baby);
    }
    setStorage(StorageKeys.BABY_LIST, babyList.value);

    return baby;
  } catch (error: any) {
    console.error("fetch baby detail error:", error);
    throw error;
  }
}

// ============ 增删改操作已移除 ============
// 页面组件应直接调用 API 层:
// - import * as babyApi from '@/api/baby'
// - await babyApi.apiCreateBaby(data)
// - await babyApi.apiUpdateBaby(id, data)
// - await babyApi.apiDeleteBaby(id)
//
// 调用成功后,需要手动刷新宝宝列表:
// - await fetchBabyList()

// ============ 本地操作函数 ============

/**
 * 设置当前宝宝
 *
 * ⚠️ 向后兼容: 函数签名保持不变
 */
export function setCurrentBaby(id: string) {
  currentBabyId.value = id;
  setStorage(StorageKeys.CURRENT_BABY_ID, id);
}

/**
 * 根据 ID 获取宝宝信息(本地)
 *
 * ⚠️ 向后兼容: 函数签名保持不变
 */
export function getBabyById(id: string): BabyProfile | null {
  return babyList.value.find((baby) => baby.babyId === id) || null;
}

/**
 * 清除宝宝数据 (用于登出)
 *
 * ⚠️ 向后兼容: 函数签名保持不变
 */
export function clearBabyData() {
  babyList.value = [];
  currentBabyId.value = "";
  collaboratorsMap.value.clear();
  myPermissionsMap.value.clear();
  setStorage(StorageKeys.BABY_LIST, []);
  setStorage(StorageKeys.CURRENT_BABY_ID, "");
}

// ============ Collaborator 相关方法 ============

/**
 * 设置宝宝的协作者列表
 */
export function setCollaborators(
  babyId: string,
  collaborators: BabyCollaborator[]
): void {
  collaboratorsMap.value.set(babyId, collaborators);
}

/**
 * 获取宝宝的协作者列表
 */
export function getCollaborators(babyId: string): BabyCollaborator[] | undefined {
  return collaboratorsMap.value.get(babyId);
}

/**
 * 设置当前用户对宝宝的权限
 */
export function setMyPermission(
  babyId: string,
  permission: MyPermission
): void {
  myPermissionsMap.value.set(babyId, permission);
}

/**
 * 获取当前用户对宝宝的权限
 */
export function getMyPermission(babyId: string): MyPermission | undefined {
  return myPermissionsMap.value.get(babyId);
}

/**
 * 清除宝宝的协作者数据（当宝宝被删除或权限过期时）
 */
export function clearCollaboratorData(babyId: string): void {
  collaboratorsMap.value.delete(babyId);
  myPermissionsMap.value.delete(babyId);
}

// ============ 导出 ============

// 直接导出 computed 对象，支持响应式
export { babyList, currentBabyId, currentBaby, collaboratorsMap, myPermissionsMap };
