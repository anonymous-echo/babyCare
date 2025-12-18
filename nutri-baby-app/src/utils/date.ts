/**
 * 日期时间工具函数
 */

/**
 * 格式化日期
 * @param timestamp 时间戳
 * @param format 格式字符串
 */
export function formatDate(timestamp: number, format: string = 'YYYY-MM-DD HH:mm:ss'): string {
  const date = new Date(timestamp)

  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')

  return format
    .replace('YYYY', String(year))
    .replace('MM', month)
    .replace('DD', day)
    .replace('HH', hours)
    .replace('mm', minutes)
    .replace('ss', seconds)
}

/**
 * 格式化时间为相对时间
 * @param timestamp 时间戳
 */
export function formatRelativeTime(timestamp: number): string {
  const now = Date.now()
  const diff = now - timestamp
  const seconds = Math.floor(diff / 1000)
  const minutes = Math.floor(seconds / 60)
  const hours = Math.floor(minutes / 60)
  const days = Math.floor(hours / 24)

  if (seconds < 60) {
    return '刚刚'
  } else if (minutes < 60) {
    return `${minutes}分钟前`
  } else if (hours < 24) {
    return `${hours}小时前`
  } else if (days < 7) {
    return `${days}天前`
  } else {
    return formatDate(timestamp, 'MM-DD HH:mm')
  }
}

/**
 * 计算两个时间戳之间的时长
 * @param startTime 开始时间戳
 * @param endTime 结束时间戳
 * @returns 时长对象 { hours, minutes, seconds }
 */
export function calculateDuration(startTime: number, endTime: number): {
  hours: number
  minutes: number
  seconds: number
  totalMinutes: number
} {
  const diff = endTime - startTime
  const seconds = Math.floor(diff / 1000)
  const minutes = Math.floor(seconds / 60)
  const hours = Math.floor(minutes / 60)

  return {
    hours,
    minutes: minutes % 60,
    seconds: seconds % 60,
    totalMinutes: minutes,
  }
}

/**
 * 格式化时长为字符串
 * @param minutes 总分钟数
 */
export function formatDuration(minutes: number): string {
  const hours = Math.floor(minutes / 60)
  const mins = minutes % 60

  if (hours > 0) {
    return `${hours} 小时 ${mins} 分钟`
  } else {
    return `${mins} 分钟`
  }
}

/**
 * 获取今天的开始时间戳
 */
export function getTodayStart(): number {
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  return today.getTime()
}

/**
 * 获取今天的结束时间戳
 */
export function getTodayEnd(): number {
  const today = new Date()
  today.setHours(23, 59, 59, 999)
  return today.getTime()
}

/**
 * 判断时间戳是否为今天
 */
export function isToday(timestamp: number): boolean {
  const today = new Date()
  const target = new Date(timestamp)
  return (
    today.getFullYear() === target.getFullYear() &&
    today.getMonth() === target.getMonth() &&
    today.getDate() === target.getDate()
  )
}

/**
 * 根据出生日期计算年龄
 * @param birthDate 出生日期字符串 YYYY-MM-DD
 */
export function calculateAge(birthDate: string): string {
  const birth = new Date(birthDate)
  const now = new Date()

  const years = now.getFullYear() - birth.getFullYear()
  const months = now.getMonth() - birth.getMonth()
  const days = now.getDate() - birth.getDate()

  let ageYears = years
  let ageMonths = months
  let ageDays = days

  if (ageDays < 0) {
    ageMonths--
    const lastMonth = new Date(now.getFullYear(), now.getMonth(), 0)
    ageDays += lastMonth.getDate()
  }

  if (ageMonths < 0) {
    ageYears--
    ageMonths += 12
  }

  if (ageYears > 0) {
    return `${ageYears}岁${ageMonths}个月`
  } else if (ageMonths > 0) {
    return `${ageMonths}个月${ageDays}天`
  } else {
    return `${ageDays}天`
  }
}

/**
 * 获取本周开始时间
 */
export function getWeekStart(): number {
  const now = new Date()
  const day = now.getDay() || 7 // 周日为0,转为7
  const diff = day - 1
  const weekStart = new Date(now.getFullYear(), now.getMonth(), now.getDate() - diff)
  weekStart.setHours(0, 0, 0, 0)
  return weekStart.getTime()
}

/**
 * 获取本月开始时间
 */
export function getMonthStart(): number {
  const now = new Date()
  const monthStart = new Date(now.getFullYear(), now.getMonth(), 1)
  monthStart.setHours(0, 0, 0, 0)
  return monthStart.getTime()
}