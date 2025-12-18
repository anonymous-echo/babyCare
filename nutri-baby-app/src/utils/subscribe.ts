/**
 * 订阅消息工具函数
 * 用于管理微信一次性订阅消息的授权
 */

import {
  saveSubscribeAuth,
  getSubscribeStatus,
  type SubscribeAuthRecord
} from '@/api/subscribe'

// 订阅消息模板配置
// TODO: 替换为实际的模板ID
// 获取方式: 微信公众平台 -> 功能 -> 订阅消息 -> 选择模板
export const SUBSCRIBE_TEMPLATES = {
  // 喂养提醒
  FEEDING_REMINDER: {
    templateId: 'YOUR_FEEDING_REMINDER_TEMPLATE_ID', // TODO: 替换为实际的模板ID
    templateType: 'breast_feeding_reminder',
    name: '喂养提醒'
  },
  // 疫苗提醒
  VACCINE_REMINDER: {
    templateId: 'YOUR_VACCINE_REMINDER_TEMPLATE_ID', // TODO: 替换为实际的模板ID
    templateType: 'vaccine_reminder',
    name: '疫苗提醒'
  }
} as const

/**
 * 请求订阅授权(智能策略)
 * @param templateType 模板类型
 * @param options 配置选项
 * @returns Promise<boolean> 是否授权成功
 */
export async function requestSubscribeAuth(
  templateType: keyof typeof SUBSCRIBE_TEMPLATES,
  options: {
    force?: boolean // 是否强制询问(跳过检查)
    silentOnReject?: boolean // 用户拒绝时是否静默处理
  } = {}
): Promise<boolean> {
  try {
    const template = SUBSCRIBE_TEMPLATES[templateType]
    if (!template) {
      console.error('[Subscribe] Invalid template type:', templateType)
      return false
    }

    // 检查是否需要询问授权(避免频繁打扰用户)
    if (!options.force) {
      const shouldAsk = await shouldRequestAuth(templateType)
      if (!shouldAsk) {
        console.log('[Subscribe] Skip requesting auth (recently asked or has available auth)')
        return false
      }
    }

    // 调用微信订阅消息API
    const res = await uni.requestSubscribeMessage({
      tmplIds: [template.templateId]
    })

    // 检查授权结果
    const authStatus = res[template.templateId]
    const isAccepted = authStatus === 'accept'

    if (isAccepted) {
      // 上报后端保存授权记录
      await uploadSubscribeAuth([{
        templateId: template.templateId,
        templateType: template.templateType,
        status: 'accept'
      }])

      console.log('[Subscribe] Authorization succeeded:', template.name)

      // 记录本次询问时间
      setLastAskTime(templateType)
      return true
    } else {
      if (!options.silentOnReject) {
        console.log('[Subscribe] User rejected:', template.name, authStatus)
      }
      // 记录本次询问时间(即使拒绝也记录,避免频繁打扰)
      setLastAskTime(templateType)
      return false
    }
  } catch (error: any) {
    if (error.errMsg?.includes('requestSubscribeMessage:fail cancel')) {
      // 用户取消授权
      if (!options.silentOnReject) {
        console.log('[Subscribe] User cancelled authorization')
      }
    } else {
      console.error('[Subscribe] Request failed:', error)
    }
    return false
  }
}

/**
 * 判断是否应该请求授权(智能策略)
 * @param templateType 模板类型
 * @returns Promise<boolean>
 */
async function shouldRequestAuth(templateType: string): Promise<boolean> {
  // 1. 检查是否有可用的授权
  try {
    const hasAuth = await checkAvailableAuth(templateType)
    if (hasAuth) {
      return false // 已有可用授权,不需要再问
    }
  } catch (error) {
    console.error('[Subscribe] Check available auth failed:', error)
  }

  // 2. 检查距离上次询问是否超过24小时
  const lastAskTime = getLastAskTime(templateType)
  if (lastAskTime) {
    const hoursSinceLastAsk = (Date.now() - lastAskTime) / (1000 * 60 * 60)
    if (hoursSinceLastAsk < 24) {
      return false // 24小时内已询问过,避免频繁打扰
    }
  }

  return true
}

/**
 * 检查用户是否有可用的授权
 * @param templateType 模板类型
 * @returns Promise<boolean>
 */
async function checkAvailableAuth(templateType: string): Promise<boolean> {
  try {
    const res = await getSubscribeStatus()
    const subscription = res.subscriptions.find(s => s.templateType === templateType)
    return subscription?.status === 'available'
  } catch (error) {
    console.error('[Subscribe] Check auth failed:', error)
    return false
  }
}

/**
 * 上报授权记录到后端(内部使用)
 * @param records 授权记录列表
 */
async function uploadSubscribeAuth(records: SubscribeAuthRecord[]): Promise<void> {
  try {
    await saveSubscribeAuth(records)
  } catch (error) {
    console.error('[Subscribe] Save auth failed:', error)
    throw error
  }
}

/**
 * 获取上次询问授权的时间
 * @param templateType 模板类型
 */
function getLastAskTime(templateType: string): number | null {
  const key = `last_subscribe_ask_${templateType}`
  const value = uni.getStorageSync(key)
  return value ? parseInt(value) : null
}

/**
 * 记录本次询问授权的时间
 * @param templateType 模板类型
 */
function setLastAskTime(templateType: string): void {
  const key = `last_subscribe_ask_${templateType}`
  uni.setStorageSync(key, Date.now().toString())
}

/**
 * 批量请求多个订阅授权
 * @param templateTypes 模板类型列表
 * @returns Promise<{ success: string[]; failed: string[] }>
 */
export async function requestMultipleSubscribeAuth(
  templateTypes: Array<keyof typeof SUBSCRIBE_TEMPLATES>
): Promise<{ success: string[]; failed: string[] }> {
  const templates = templateTypes
    .map(type => SUBSCRIBE_TEMPLATES[type])
    .filter(Boolean)

  if (templates.length === 0) {
    return { success: [], failed: [] }
  }

  try {
    // 微信一次最多请求3个订阅消息
    const tmplIds = templates.slice(0, 3).map(t => t.templateId)

    const res = await uni.requestSubscribeMessage({ tmplIds })

    const records: Array<{ templateId: string; templateType: string; status: 'accept' | 'reject' }> = []
    const success: string[] = []
    const failed: string[] = []

    templates.forEach(template => {
      const status = res[template.templateId]
      const isAccepted = status === 'accept'

      records.push({
        templateId: template.templateId,
        templateType: template.templateType,
        status: isAccepted ? 'accept' : 'reject'
      })

      if (isAccepted) {
        success.push(template.templateType)
      } else {
        failed.push(template.templateType)
      }
    })

    // 上报后端
    if (records.length > 0) {
      await uploadSubscribeAuth(records)
    }

    return { success, failed }
  } catch (error) {
    console.error('[Subscribe] Request multiple auth failed:', error)
    return { success: [], failed: templateTypes.map(t => SUBSCRIBE_TEMPLATES[t].templateType) }
  }
}
