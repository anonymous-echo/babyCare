/**
 * HTTP 请求工具
 */

import type { ApiResponse } from '@/types'
import { StorageKeys, getStorage } from '@/utils/storage'

// API 基础配置
const BASE_URL = import.meta.env.VITE_API_BASE_URL || 'https://api.example.com'
const TIMEOUT = 120000 // 2分钟超时，用于AI分析等长时间请求

/**
 * 请求配置接口
 */
interface RequestConfig {
  url: string
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE'
  data?: any
  header?: Record<string, string>
  timeout?: number
  retry?: number // 重试次数
  retryDelay?: number // 重试延迟（毫秒）
  showError?: boolean // 是否显示错误提示
}

/**
 * 获取请求头
 */
function getHeaders(): Record<string, string> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  }

  // 从本地存储获取 token (使用 getStorage 函数确保正确解析 JSON)
  const token = getStorage<string>(StorageKeys.TOKEN)
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  return headers
}

/**
 * 构建 URL 查询字符串
 */
function buildQueryString(params?: any): string {
  if (!params || Object.keys(params).length === 0) {
    return ''
  }
  const queryParams = Object.entries(params)
    .filter(([_, value]) => value !== undefined && value !== null)
    .map(([key, value]) => `${encodeURIComponent(key)}=${encodeURIComponent(String(value))}`)
    .join('&')
  return queryParams ? `?${queryParams}` : ''
}

/**
 * 延迟函数
 */
function delay(ms: number): Promise<void> {
  return new Promise(resolve => setTimeout(resolve, ms))
}

/**
 * 封装的请求方法（带重试机制）
 */
export function request<T = any>(config: RequestConfig): Promise<ApiResponse<T>> {
  const {
    retry = 0,
    retryDelay = 1000,
    showError = true,
    ...requestConfig
  } = config

  return requestWithRetry<T>(requestConfig, retry, retryDelay, showError)
}

/**
 * 带重试的请求实现
 */
function requestWithRetry<T = any>(
  config: Omit<RequestConfig, 'retry' | 'retryDelay' | 'showError'>,
  retriesLeft: number,
  retryDelay: number,
  showError: boolean,
  attempt: number = 1
): Promise<ApiResponse<T>> {
  return new Promise((resolve, reject) => {
    // 为 GET 请求构建查询字符串
    const finalUrl = config.method === 'GET'
      ? `${BASE_URL}${config.url}${buildQueryString(config.data)}`
      : `${BASE_URL}${config.url}`

    uni.request({
      url: finalUrl,
      method: config.method || 'GET',
      data: config.method === 'GET' ? undefined : config.data,
      header: {
        ...getHeaders(),
        ...config.header,
      },
      timeout: config.timeout || TIMEOUT,
      success: (res) => {
        const data = res.data as ApiResponse<T>

        // 请求成功
        if (res.statusCode === 200) {
          if (data.code === 0) {
            resolve(data)
          } else {
            // 业务错误
            if (showError) {
              uni.showToast({
                title: data.message || '请求失败',
                icon: 'none',
                duration: 2000,
              })
            }
            reject(data)
          }
        } else if (res.statusCode === 401) {
          // token 过期,跳转登录
          if (showError) {
            uni.showToast({
              title: '登录已过期,请重新登录',
              icon: 'none',
            })
          }
          // 清除 token
          uni.removeStorageSync(StorageKeys.TOKEN)
          // 跳转到登录页
          setTimeout(() => {
            uni.reLaunch({
              url: '/pages/user/login',
            })
          }, 1500)
          reject(data)
        } else if (res.statusCode === 404) {
          // 404错误不显示toast，让调用方处理
          reject({ ...data, statusCode: 404 })
        } else {
          // HTTP 错误
          if (showError) {
            uni.showToast({
              title: `请求失败: ${res.statusCode}`,
              icon: 'none',
            })
          }
          reject(data)
        }
      },
      fail: async (err) => {
        console.error(`request error (attempt ${attempt}):`, err)

        // 如果还有重试次数，进行重试
        if (retriesLeft > 0) {
          console.log(`Retrying... (${retriesLeft} retries left)`)
          await delay(retryDelay)
          
          try {
            const result = await requestWithRetry<T>(
              config,
              retriesLeft - 1,
              retryDelay,
              showError,
              attempt + 1
            )
            resolve(result)
          } catch (retryError) {
            reject(retryError)
          }
        } else {
          // 没有重试次数了，显示错误
          if (showError) {
            uni.showToast({
              title: '网络请求失败，请检查网络连接',
              icon: 'none',
              duration: 2000
            })
          }
          reject(err)
        }
      },
    })
  })
}

/**
 * GET 请求
 */
export function get<T = any>(url: string, data?: any): Promise<ApiResponse<T>> {
  return request<T>({
    url,
    method: 'GET',
    data,
  })
}

/**
 * POST 请求
 */
export function post<T = any>(url: string, data?: any): Promise<ApiResponse<T>> {
  return request<T>({
    url,
    method: 'POST',
    data,
  })
}

/**
 * PUT 请求
 */
export function put<T = any>(url: string, data?: any): Promise<ApiResponse<T>> {
  return request<T>({
    url,
    method: 'PUT',
    data,
  })
}

/**
 * DELETE 请求
 */
export function del<T = any>(url: string, data?: any): Promise<ApiResponse<T>> {
  return request<T>({
    url,
    method: 'DELETE',
    data,
  })
}

/**
 * 文件上传
 */
interface UploadConfig {
  filePath: string
  name?: string
  formData?: Record<string, string>
}

export function uploadFile(config: UploadConfig): Promise<any> {
  return new Promise((resolve, reject) => {
    uni.uploadFile({
      url: `${BASE_URL}/upload`,
      filePath: config.filePath,
      name: config.name || 'file',
      formData: config.formData || {},
      header: getHeaders(),
      success: (res) => {
        if (res.statusCode === 200) {
          const data = JSON.parse(res.data)
          resolve(data)
        } else {
          reject(res)
        }
      },
      fail: reject,
    })
  })
}