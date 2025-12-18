/**
 * 数据导出工具
 */

import type { BabyProfile, FeedingRecord, DiaperRecord, SleepRecord } from '@/types'
import { formatDate } from './date'

/**
 * 导出数据接口
 */
export interface ExportData {
  exportTime: number
  exportTimeText: string
  babies: BabyProfile[]
  feedingRecords: FeedingRecord[]
  diaperRecords: DiaperRecord[]
  sleepRecords: SleepRecord[]
}

/**
 * 导出所有数据为 JSON 格式
 */
export function exportAllDataToJSON(data: ExportData): string {
  return JSON.stringify(data, null, 2)
}

/**
 * 保存数据到本地文件
 */
export function saveDataToFile(data: ExportData): Promise<string> {
  return new Promise((resolve, reject) => {
    try {
      const jsonStr = exportAllDataToJSON(data)
      const fileName = `baby_data_${formatDate(Date.now(), 'YYYYMMDD_HHmmss')}.json`

      // 使用 uni.getFileSystemManager 保存文件
      const fs = uni.getFileSystemManager()
      const filePath = `${wx.env.USER_DATA_PATH}/${fileName}`

      fs.writeFile({
        filePath,
        data: jsonStr,
        encoding: 'utf8',
        success: () => {
          resolve(filePath)
        },
        fail: (err) => {
          reject(err)
        }
      })
    } catch (error) {
      reject(error)
    }
  })
}

/**
 * 分享数据文件
 */
export function shareDataFile(filePath: string): Promise<void> {
  return new Promise((resolve, reject) => {
    uni.shareFileMessage({
      filePath,
      success: () => {
        resolve()
      },
      fail: (err) => {
        reject(err)
      }
    })
  })
}

/**
 * 导入数据从 JSON
 */
export function importDataFromJSON(jsonStr: string): ExportData | null {
  try {
    const data = JSON.parse(jsonStr) as ExportData

    // 验证数据结构
    if (!data.babies || !Array.isArray(data.babies)) {
      throw new Error('数据格式错误：缺少 babies 字段')
    }
    if (!data.feedingRecords || !Array.isArray(data.feedingRecords)) {
      throw new Error('数据格式错误：缺少 feedingRecords 字段')
    }
    if (!data.diaperRecords || !Array.isArray(data.diaperRecords)) {
      throw new Error('数据格式错误：缺少 diaperRecords 字段')
    }
    if (!data.sleepRecords || !Array.isArray(data.sleepRecords)) {
      throw new Error('数据格式错误：缺少 sleepRecords 字段')
    }

    return data
  } catch (error) {
    console.error('导入数据失败:', error)
    return null
  }
}

/**
 * 读取文件内容
 */
export function readFileContent(filePath: string): Promise<string> {
  return new Promise((resolve, reject) => {
    const fs = uni.getFileSystemManager()

    fs.readFile({
      filePath,
      encoding: 'utf8',
      success: (res) => {
        resolve(res.data as string)
      },
      fail: (err) => {
        reject(err)
      }
    })
  })
}

/**
 * 生成导出数据摘要
 */
export function generateExportSummary(data: ExportData): string {
  return `
导出时间: ${data.exportTimeText}
宝宝数量: ${data.babies.length}
喂养记录: ${data.feedingRecords.length} 条
换尿布记录: ${data.diaperRecords.length} 条
睡眠记录: ${data.sleepRecords.length} 条
总记录数: ${data.feedingRecords.length + data.diaperRecords.length + data.sleepRecords.length} 条
  `.trim()
}
