/**
 * 订阅消息管理 Store
 *
 * 核心功能:
 * 1. 管理订阅消息模板配置
 * 2. 追踪用户授权状态
 * 3. 控制引导显示逻辑(避免频繁骚扰)
 * 4. 提供智能申请时机判断
 * 5. 同步授权记录到后端服务器
 */
import { ref, reactive } from 'vue'
import type {
  SubscribeMessageType,
  SubscribeMessageTemplate,
  SubscribeAuthRecord,
  SubscribeAuthStatus,
  SubscribeReminderConfig,
  SubscribeGuideRecord,
} from '@/types'
import { StorageKeys, getStorage, setStorage } from '@/utils/storage'
import { request } from '@/utils/request'

/**
 * 订阅消息模板配置
 * 模板ID来自微信公众平台
 */
const TEMPLATE_CONFIGS: SubscribeMessageTemplate[] = [
  {
    type: 'vaccine_reminder',
    templateId: 'J6RbROH-yhNdgj2FPwrz4FnzzpITH2KcHV5h9qjcVbI',
    title: '疫苗接种提醒',
    keywords: ['宝宝名称', '疫苗名称', '接种时间', '接种地址', '接种针数'],
    description: '提前3天提醒您带宝宝接种疫苗',
    // icon: '/static/icons/vaccine.png',  // 暂时禁用,使用默认emoji
    priority: 5,
  },
  {
    type: 'pump_reminder',
    templateId: 'ovzn4mwR2i1oOLR8ps6BspHHrxgYPCZ4lqhzvgSW5UU',
    title: '吸奶器吸奶提醒',
    keywords: ['距离上次', '上次时长', '上次总量', '上次时间', '温馨提示'],
    description: '定时提醒您使用吸奶器',
    // icon: '/static/icons/pump.png',
    priority: 2,
  },
  {
    type: 'breast_feeding_reminder',
    templateId: '2JRV0DnOHnasHzzamWFoWGaUxrgW6GY69-eGn4tBFZE',
    title: '母乳喂养提醒',
    keywords: ['上次时间', '距离上次', '上次位置', '温馨提示'],
    description: '智能提醒您该给宝宝哺乳了，不错过每次喂养时间',
    // icon: '/static/icons/breast.png',
    priority: 4,
  },
  {
    type: 'feeding_duration_alert',
    templateId: 'QSMUFLFfRiEVbMRL0DfWNIPoN3kNN8nA0mRpwXyinNw',
    title: '喂奶时间过长提醒',
    keywords: ['开始时间', '持续时长', '温馨提示'],
    description: '当喂奶时长超过建议时间时提醒，保护乳房健康',
    // icon: '/static/icons/alert.png',
    priority: 3,
  },
  {
    type: 'bottle_feeding_reminder',
    templateId: 'ssttSBSWM_IXh5zVOu9GBeuabX8NFcwM2IG-VK-RXNY',
    title: '奶瓶喂养提醒',
    keywords: ['上次时间', '距离上次', '上次容量', '上次类型', '温馨提示'],
    description: '开启提醒，智能提示您宝宝的喂养时间，让喂养更有规律',
    // icon: '/static/icons/bottle.png',
    priority: 4,
  },
]

/**
 * 引导显示间隔配置(毫秒)
 */
const GUIDE_SHOW_INTERVALS = {
  default: 7 * 24 * 60 * 60 * 1000, // 默认7天
  afterReject: 14 * 24 * 60 * 60 * 1000, // 用户拒绝后14天
  afterBan: Infinity, // 永久拒绝后不再显示
}

/**
 * 授权记录 Map
 */
const authRecords = reactive<Map<SubscribeMessageType, SubscribeAuthRecord>>(new Map())

/**
 * 引导显示记录 Map
 */
const guideRecords = reactive<Map<SubscribeMessageType, SubscribeGuideRecord>>(new Map())

/**
 * 提醒配置 Map
 */
const reminderConfigs = reactive<Map<SubscribeMessageType, SubscribeReminderConfig>>(new Map())

/**
 * 初始化标记
 */
let isInitialized = false

/**
 * 延迟初始化 - 从本地存储加载数据
 */
function initializeIfNeeded() {
  if (!isInitialized) {
    // 加载授权记录
    const savedAuthRecords = getStorage<Record<SubscribeMessageType, SubscribeAuthRecord>>(
      StorageKeys.SUBSCRIBE_AUTH_RECORDS
    )
    if (savedAuthRecords) {
      Object.entries(savedAuthRecords).forEach(([type, record]) => {
        authRecords.set(type as SubscribeMessageType, record)
      })
    }

    // 加载引导记录
    const savedGuideRecords = getStorage<Record<SubscribeMessageType, SubscribeGuideRecord>>(
      StorageKeys.SUBSCRIBE_GUIDE_RECORDS
    )
    if (savedGuideRecords) {
      Object.entries(savedGuideRecords).forEach(([type, record]) => {
        guideRecords.set(type as SubscribeMessageType, record)
      })
    }

    // 加载提醒配置
    const savedReminderConfigs = getStorage<Record<SubscribeMessageType, SubscribeReminderConfig>>(
      StorageKeys.SUBSCRIBE_REMINDER_CONFIGS
    )
    if (savedReminderConfigs) {
      Object.entries(savedReminderConfigs).forEach(([type, config]) => {
        reminderConfigs.set(type as SubscribeMessageType, config)
      })
    } else {
      // 初始化默认提醒配置
      TEMPLATE_CONFIGS.forEach((template) => {
        reminderConfigs.set(template.type, {
          type: template.type,
          enabled: false,
          advanceDays: template.type === 'vaccine_reminder' ? 3 : undefined,
          intervalMinutes:
            template.type === 'breast_feeding_reminder'
              ? 180 // 母乳提醒每3小时一次
              : template.type === 'bottle_feeding_reminder'
              ? 180 // 奶瓶提醒每3小时一次
              : undefined,
        })
      })
      saveReminderConfigs()
    }

    isInitialized = true
  }
}

/**
 * 保存授权记录到存储
 */
function saveAuthRecords() {
  const recordsObj: Record<string, SubscribeAuthRecord> = {}
  authRecords.forEach((record, type) => {
    recordsObj[type] = record
  })
  setStorage(StorageKeys.SUBSCRIBE_AUTH_RECORDS, recordsObj)
}

/**
 * 保存引导记录到存储
 */
function saveGuideRecords() {
  const recordsObj: Record<string, SubscribeGuideRecord> = {}
  guideRecords.forEach((record, type) => {
    recordsObj[type] = record
  })
  setStorage(StorageKeys.SUBSCRIBE_GUIDE_RECORDS, recordsObj)
}

/**
 * 保存提醒配置到存储
 */
function saveReminderConfigs() {
  const configsObj: Record<string, SubscribeReminderConfig> = {}
  reminderConfigs.forEach((config, type) => {
    configsObj[type] = config
  })
  setStorage(StorageKeys.SUBSCRIBE_REMINDER_CONFIGS, configsObj)
}

/**
 * 获取模板配置
 */
export function getTemplateConfig(type: SubscribeMessageType): SubscribeMessageTemplate | undefined {
  return TEMPLATE_CONFIGS.find((t) => t.type === type)
}

/**
 * 获取所有模板配置
 */
export function getAllTemplateConfigs(): SubscribeMessageTemplate[] {
  return TEMPLATE_CONFIGS
}

/**
 * 获取授权记录
 */
export function getAuthRecord(type: SubscribeMessageType): SubscribeAuthRecord | undefined {
  initializeIfNeeded()
  return authRecords.get(type)
}

/**
 * 获取授权状态
 */
export function getAuthStatus(type: SubscribeMessageType): SubscribeAuthStatus {
  const record = getAuthRecord(type)
  return record?.status || 'unknown'
}

/**
 * 判断是否应该显示引导(一次性订阅消息版本)
 *
 * 规则:
 * 1. ❌ 不再检查本地 authStatus,因为一次性订阅消息授权会被消耗
 * 2. 如果用户已被Ban(拒绝3次),不显示
 * 3. 如果引导被永久关闭,不显示
 * 4. 如果距离上次显示时间未达到间隔,不显示
 *
 * 注意:对于一次性订阅消息,每次授权只能发送一条消息,发送后授权自动失效。
 * 因此不能依赖本地 authStatus='accept' 来判断是否显示引导,必须每次都询问用户。
 */
export function shouldShowGuide(type: SubscribeMessageType): boolean {
  initializeIfNeeded()

  // 检查是否被Ban
  const authStatus = getAuthStatus(type)
  if (authStatus === 'ban') {
    return false // 已被Ban,不再显示
  }

  // 检查引导记录
  const guideRecord = guideRecords.get(type)
  if (!guideRecord) {
    return true // 首次显示
  }

  // 检查是否永久关闭
  if (guideRecord.isDismissedForever) {
    return false
  }

  // 检查时间间隔
  const now = Date.now()
  if (now < guideRecord.nextShowTime) {
    return false // 未到下次显示时间
  }

  return true
}

/**
 * 记录引导显示
 */
export function recordGuideShown(type: SubscribeMessageType) {
  initializeIfNeeded()

  const now = Date.now()
  const authStatus = getAuthStatus(type)

  let interval = GUIDE_SHOW_INTERVALS.default
  if (authStatus === 'reject') {
    interval = GUIDE_SHOW_INTERVALS.afterReject
  }

  const existingRecord = guideRecords.get(type)
  const newRecord: SubscribeGuideRecord = {
    type,
    showCount: (existingRecord?.showCount || 0) + 1,
    lastShowTime: now,
    nextShowTime: now + interval,
    isDismissedForever: false,
  }

  guideRecords.set(type, newRecord)
  saveGuideRecords()
}

/**
 * 永久关闭引导
 */
export function dismissGuideForever(type: SubscribeMessageType) {
  initializeIfNeeded()

  const existingRecord = guideRecords.get(type)
  const newRecord: SubscribeGuideRecord = {
    type,
    showCount: existingRecord?.showCount || 0,
    lastShowTime: existingRecord?.lastShowTime || Date.now(),
    nextShowTime: Infinity,
    isDismissedForever: true,
  }

  guideRecords.set(type, newRecord)
  saveGuideRecords()
}

/**
 * 请求订阅消息授权
 *
 * @param types 要申请的消息类型数组
 * @returns Promise<Map<SubscribeMessageType, 'accept' | 'reject'>>
 */
export async function requestSubscribeMessage(
  types: SubscribeMessageType[]
): Promise<Map<SubscribeMessageType, 'accept' | 'reject'>> {
  initializeIfNeeded()

  const templateIds = types.map((type) => {
    const config = getTemplateConfig(type)
    if (!config) {
      throw new Error(`未找到模板配置: ${type}`)
    }
    return config.templateId
  })

  return new Promise((resolve) => {
    uni.requestSubscribeMessage({
      tmplIds: templateIds,
      success: (res) => {
      console.log("申请订阅的模板 IDs",templateIds)
        console.log('[Subscribe] requestSubscribeMessage success:', res)

        const results = new Map<SubscribeMessageType, 'accept' | 'reject'>()

        types.forEach((type, index) => {
          const templateId = templateIds[index]
          const status = res[templateId]

          let authStatus: 'accept' | 'reject' = 'reject'
          if (status === 'accept') {
            authStatus = 'accept'
          }

          results.set(type, authStatus)
          updateAuthRecord(type, authStatus)
        })

        // 上传授权记录到后端
        uploadAuthRecordsToBackend(results).catch((err) => {
          console.error('[Subscribe] Failed to upload auth records:', err)
        })

        resolve(results)
      },
      fail: (err) => {
        console.error('[Subscribe] requestSubscribeMessage fail:', err)
        console.error('[Subscribe] fail error details:', JSON.stringify(err))

        // 失败视为拒绝
        const results = new Map<SubscribeMessageType, 'accept' | 'reject'>()
        types.forEach((type) => {
          results.set(type, 'reject')
          updateAuthRecord(type, 'reject')
        })

        resolve(results)
      },
    })
  })
}

/**
 * 更新授权记录
 */
function updateAuthRecord(type: SubscribeMessageType, authStatus: 'accept' | 'reject') {
  initializeIfNeeded()

  const now = Date.now()
  const existingRecord = authRecords.get(type)
  const templateConfig = getTemplateConfig(type)

  if (!templateConfig) {
    console.error(`[Subscribe] 未找到模板配置: ${type}`)
    return
  }

  let finalStatus: SubscribeAuthStatus = authStatus
  const rejectCount = (existingRecord?.rejectCount || 0) + (authStatus === 'reject' ? 1 : 0)

  // 拒绝3次后标记为 ban
  if (rejectCount >= 3) {
    finalStatus = 'ban'
  }

  const newRecord: SubscribeAuthRecord = {
    type,
    templateId: templateConfig.templateId,
    status: finalStatus,
    lastRequestTime: now,
    requestCount: (existingRecord?.requestCount || 0) + 1,
    rejectCount,
    acceptCount: (existingRecord?.acceptCount || 0) + (authStatus === 'accept' ? 1 : 0),
    updateTime: now,
  }

  authRecords.set(type, newRecord)
  saveAuthRecords()

  // 如果用户同意,自动启用对应的提醒
  if (authStatus === 'accept') {
    const config = reminderConfigs.get(type)
    if (config) {
      config.enabled = true
      reminderConfigs.set(type, config)
      saveReminderConfigs()
    }
  }
}

/**
 * 获取提醒配置
 */
export function getReminderConfig(type: SubscribeMessageType): SubscribeReminderConfig | undefined {
  initializeIfNeeded()
  return reminderConfigs.get(type)
}

/**
 * 更新提醒配置
 */
export function updateReminderConfig(type: SubscribeMessageType, config: Partial<SubscribeReminderConfig>) {
  initializeIfNeeded()

  const existingConfig = reminderConfigs.get(type)
  if (!existingConfig) {
    console.error(`[Subscribe] 未找到提醒配置: ${type}`)
    return
  }

  const newConfig = {
    ...existingConfig,
    ...config,
  }

  reminderConfigs.set(type, newConfig)
  saveReminderConfigs()
}

/**
 * 启用提醒(会检查授权状态)
 */
export async function enableReminder(type: SubscribeMessageType): Promise<boolean> {
  initializeIfNeeded()

  const authStatus = getAuthStatus(type)

  // 如果未授权或已拒绝,先申请授权
  if (authStatus !== 'accept') {
    if (authStatus === 'ban') {
      uni.showToast({
        title: '您已多次拒绝,请在设置中手动开启',
        icon: 'none',
        duration: 2000,
      })
      return false
    }

    const results = await requestSubscribeMessage([type])
    const result = results.get(type)

    if (result !== 'accept') {
      return false
    }
  }

  // 启用提醒
  updateReminderConfig(type, { enabled: true })
  return true
}

/**
 * 禁用提醒
 */
export function disableReminder(type: SubscribeMessageType) {
  updateReminderConfig(type, { enabled: false })
}

/**
 * 获取所有提醒配置
 */
export function getAllReminderConfigs(): Map<SubscribeMessageType, SubscribeReminderConfig> {
  initializeIfNeeded()
  return reminderConfigs
}

/**
 * 检查是否有任何已启用的提醒
 */
export function hasEnabledReminders(): boolean {
  initializeIfNeeded()
  return Array.from(reminderConfigs.values()).some((config) => config.enabled)
}

/**
 * 导出响应式状态(用于调试)
 */
export { authRecords, guideRecords, reminderConfigs }

/**
 * 上传授权记录到后端
 * @param authResults 授权结果Map
 */
async function uploadAuthRecordsToBackend(authResults: Map<SubscribeMessageType, 'accept' | 'reject'>) {
  try {
    // 构建请求数据
    const records = Array.from(authResults.entries())
      .filter(([_, status]) => status === 'accept') // 只上传用户同意的记录
      .map(([type, status]) => {
        const templateConfig = getTemplateConfig(type)
        return {
          templateId: templateConfig?.templateId || '',
          templateType: type,
          status: status,
        }
      })

    // 如果没有同意的记录,不发送请求
    if (records.length === 0) {
      console.log('[Subscribe] No accepted records to upload')
      return
    }

    console.log('[Subscribe] Uploading auth records to backend:', records)

    // 调用后端API
    const response = await request({
      url: '/subscribe/auth',
      method: 'POST',
      data: { records },
    })

    if (response.code === 0) {
      console.log('[Subscribe] Auth records uploaded successfully:', response.data)
      // uni.showToast({
      //   title: '订阅设置已同步',
      //   icon: 'success',
      //   duration: 1500,
      // })
    } else {
      console.error('[Subscribe] Failed to upload auth records:', response.message)
    }
  } catch (error) {
    console.error('[Subscribe] Error uploading auth records:', error)
    // 网络错误不影响前端功能,静默失败
  }
}


// 确保在模块加载完成后初始化
if (typeof window !== 'undefined') {
  setTimeout(() => initializeIfNeeded(), 0)
}
