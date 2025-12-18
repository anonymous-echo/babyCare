/**
 * 通用工具函数
 */

/**
 * 生成唯一 ID
 */
export function generateId(): string {
  return `${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
}

/**
 * 防抖函数
 */
export function debounce<T extends (...args: any[]) => any>(
  func: T,
  wait: number
): (...args: Parameters<T>) => void {
  let timeout: number | null = null

  return function (this: any, ...args: Parameters<T>) {
    if (timeout) clearTimeout(timeout)

    timeout = setTimeout(() => {
      func.apply(this, args)
    }, wait) as unknown as number
  }
}

/**
 * 节流函数
 */
export function throttle<T extends (...args: any[]) => any>(
  func: T,
  wait: number
): (...args: Parameters<T>) => void {
  let timeout: number | null = null
  let previous = 0

  return function (this: any, ...args: Parameters<T>) {
    const now = Date.now()

    if (now - previous > wait) {
      if (timeout) {
        clearTimeout(timeout)
        timeout = null
      }
      func.apply(this, args)
      previous = now
    } else if (!timeout) {
      timeout = setTimeout(() => {
        func.apply(this, args)
        previous = Date.now()
        timeout = null
      }, wait) as unknown as number
    }
  }
}

/**
 * 深拷贝
 */
export function deepClone<T>(obj: T): T {
  if (obj === null || typeof obj !== 'object') {
    return obj
  }

  if (obj instanceof Date) {
    return new Date(obj.getTime()) as any
  }

  if (obj instanceof Array) {
    const cloneA = [] as any[]
    for (let i = 0; i < obj.length; i++) {
      cloneA[i] = deepClone(obj[i])
    }
    return cloneA as any
  }

  const cloneO: any = {}
  for (const key in obj) {
    if (obj.hasOwnProperty(key)) {
      cloneO[key] = deepClone(obj[key])
    }
  }
  return cloneO
}

/**
 * 单位转换: ml 转 oz
 */
export function mlToOz(ml: number): number {
  return Math.round((ml * 0.033814) * 100) / 100
}

/**
 * 单位转换: oz 转 ml
 */
export function ozToMl(oz: number): number {
  return Math.round((oz * 29.5735) * 100) / 100
}

/**
 * 数字补零
 */
export function padZero(num: number, length: number = 2): string {
  return String(num).padStart(length, '0')
}

/**
 * 获取文件后缀名
 */
export function getFileExtension(filename: string): string {
  const lastDotIndex = filename.lastIndexOf('.')
  return lastDotIndex !== -1 ? filename.slice(lastDotIndex + 1).toLowerCase() : ''
}

/**
 * 判断是否为空对象
 */
export function isEmptyObject(obj: any): boolean {
  return Object.keys(obj).length === 0
}

/**
 * 数组去重
 */
export function unique<T>(arr: T[]): T[] {
  return Array.from(new Set(arr))
}

/**
 * 范围限制
 */
export function clamp(value: number, min: number, max: number): number {
  return Math.min(Math.max(value, min), max)
}

/**
 * 休眠函数
 */
export function sleep(ms: number): Promise<void> {
  return new Promise(resolve => setTimeout(resolve, ms))
}

/**
 * 格式化数字,添加千分位分隔符
 */
export function formatNumber(num: number): string {
  return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',')
}

/**
 * 校验手机号
 */
export function validatePhone(phone: string): boolean {
  return /^1[3-9]\d{9}$/.test(phone)
}

/**
 * 校验邮箱
 */
export function validateEmail(email: string): boolean {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)
}

/**
 * 格式化时长(秒) 为 "n小时n分钟n秒"
 * @param seconds 总秒数
 * @returns 格式化的时长字符串
 */
export function formatDuration(seconds: number): string {
  if (seconds <= 0) {
    return '0秒'
  }

  const hours = Math.floor(seconds / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  const secs = seconds % 60

  const parts: string[] = []
  if (hours > 0) {
    parts.push(`${hours}小时`)
  }
  if (minutes > 0) {
    parts.push(`${minutes}分钟`)
  }
  if (secs > 0 || parts.length === 0) {
    parts.push(`${secs}秒`)
  }

  return parts.join('')
}

/**
 * 格式化时长(秒) 为简短格式 "HH:MM:SS"
 * @param seconds 总秒数
 * @returns 格式化的时长字符串
 */
export function formatDurationShort(seconds: number): string {
  const hours = Math.floor(seconds / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  const secs = seconds % 60

  if (hours > 0) {
    return `${padZero(hours)}:${padZero(minutes)}:${padZero(secs)}`
  }
  return `${padZero(minutes)}:${padZero(secs)}`
}

export function goBack() {
  uni.navigateBack();
}

export function goBackHome() {
  uni.switchTab({
    url: "/pages/index/index",
  });
}
