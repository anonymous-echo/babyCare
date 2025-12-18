/**
 * 本地存储工具类
 */

const STORAGE_PREFIX = "nutri_baby_";

/**
 * 存储键名 - 去家庭化架构
 */
export const StorageKeys = {
  // 用户相关
  USER_INFO: `${STORAGE_PREFIX}user_info`,
  TOKEN: `${STORAGE_PREFIX}token`,

  // 宝宝相关
  CURRENT_BABY_ID: `${STORAGE_PREFIX}current_baby_id`,
  BABY_LIST: `${STORAGE_PREFIX}baby_list`,

  // 协作者相关 (替代家庭成员)
  COLLABORATORS_PREFIX: `${STORAGE_PREFIX}collaborators_`, // 按宝宝ID存储: collaborators_{babyId}

  // 记录相关
  FEEDING_RECORDS: `${STORAGE_PREFIX}feeding_records`,
  DIAPER_RECORDS: `${STORAGE_PREFIX}diaper_records`,
  SLEEP_RECORDS: `${STORAGE_PREFIX}sleep_records`,
  GROWTH_RECORDS: `${STORAGE_PREFIX}growth_records`,
  OTHER_RECORDS: `${STORAGE_PREFIX}other_records`,

  // 疫苗相关 - 按宝宝 ID 隔离存储
  VACCINE_PLANS_PREFIX: `${STORAGE_PREFIX}vaccine_plans_`, // vaccine_plans_{babyId}
  VACCINE_RECORDS: `${STORAGE_PREFIX}vaccine_records`,
  VACCINE_REMINDERS: `${STORAGE_PREFIX}vaccine_reminders`,

  // @deprecated - 旧的全局疫苗计划存储键,仅用于数据迁移
  VACCINE_PLANS: `${STORAGE_PREFIX}vaccine_plans`,

  // 设置
  SETTINGS: `${STORAGE_PREFIX}settings`,

  // 订阅消息相关
  SUBSCRIBE_AUTH_RECORDS: `${STORAGE_PREFIX}subscribe_auth_records`, // 授权记录
  SUBSCRIBE_GUIDE_RECORDS: `${STORAGE_PREFIX}subscribe_guide_records`, // 引导显示记录
  SUBSCRIBE_REMINDER_CONFIGS: `${STORAGE_PREFIX}subscribe_reminder_configs`, // 提醒配置
  FEEDING_SUBSCRIBE_RECORD: `${STORAGE_PREFIX}feeding_subscribe_record`, // 喂养订阅申请记录
  FEEDING_REMINDER_PREFERENCES: `${STORAGE_PREFIX}feeding_reminder_preferences`, // 喂养提醒用户偏好

  // 临时记录 (未完成的记录)
  TEMP_BREAST_FEEDING: `${STORAGE_PREFIX}temp_breast_feeding`, // 临时母乳喂养记录
  TEMP_SLEEP_RECORDING: `${STORAGE_PREFIX}temp_sleep_recording`, // 临时睡眠记录

  // 离线数据队列
  OFFLINE_QUEUE: `${STORAGE_PREFIX}offline_queue`,

  // 邀请相关（扫码加入流程）
  PENDING_INVITE_CODE: `${STORAGE_PREFIX}pending_invite_code`, // 待处理的邀请短码
  AUTO_JOIN_AFTER_LOGIN: `${STORAGE_PREFIX}auto_join_after_login`, // 登录后自动加入的邀请信息
};

/**
 * 设置存储
 */
export function setStorage<T>(key: string, value: T): void {
  try {
    uni.setStorageSync(key, JSON.stringify(value));
  } catch (e) {
    console.error("setStorage error:", e);
  }
}

/**
 * 获取存储
 */
export function getStorage<T>(key: string): T | null {
  try {
    const value = uni.getStorageSync(key);
    return value ? JSON.parse(value) : null;
  } catch (e) {
    console.error("getStorage error:", e);
    return null;
  }
}

/**
 * 移除存储
 */
export function removeStorage(key: string): void {
  try {
    uni.removeStorageSync(key);
  } catch (e) {
    console.error("removeStorage error:", e);
  }
}

/**
 * 清空所有存储
 */
export function clearStorage(): void {
  try {
    uni.clearStorageSync();
  } catch (e) {
    console.error("clearStorage error:", e);
  }
}

/**
 * 获取存储信息
 */
export function getStorageInfo() {
  try {
    return uni.getStorageInfoSync();
  } catch (e) {
    console.error("getStorageInfo error:", e);
    return null;
  }
}

/**
 * 清理已废弃的家庭相关数据 (用于迁移到去家庭化架构)
 */
export function clearDeprecatedFamilyData(): void {
  try {
    removeStorage(StorageKeys.FAMILY_LIST);
    removeStorage(StorageKeys.CURRENT_FAMILY_ID);
    removeStorage(StorageKeys.FAMILY_MEMBERS);
    removeStorage(StorageKeys.INVITATIONS);
    console.log("Deprecated family data cleared");
  } catch (e) {
    console.error("clearDeprecatedFamilyData error:", e);
  }
}
