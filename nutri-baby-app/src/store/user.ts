/**
 * 用户状态管理
 * 职责: 状态管理 + 本地计算,API 调用委托给 api 层
 *
 * ⚠️ 向后兼容: 所有导出函数的签名保持不变,页面组件无需修改
 */
import { ref } from 'vue'
import type { UserInfo } from '@/types'
import { StorageKeys, getStorage, setStorage, removeStorage } from '@/utils/storage'
import * as authApi from '@/api/auth'

// ============ 状态定义 ============

// 用户信息 - 延迟初始化
const userInfo = ref<UserInfo | null>(null)

// Token - 延迟初始化
const token = ref<string | null>(null)

// 是否已登录
const isLoggedIn = ref<boolean>(false)

// 是否为新用户 (首次登录且无宝宝)
const isNewUser = ref<boolean>(false)

// 初始化标记
let isInitialized = false

/**
 * 初始化用户状态 - 从本地存储恢复
 *
 * ⚠️ 注意: 小程序环境不支持 setTimeout 延迟初始化
 * 必须在 App.vue 启动时立即调用
 */
export function initialize() {
  if (!isInitialized) {
    userInfo.value = getStorage<UserInfo>(StorageKeys.USER_INFO) || null
    token.value = getStorage<string>(StorageKeys.TOKEN) || null
    isLoggedIn.value = !!userInfo.value && !!token.value
    isInitialized = true

    console.log('[User Store] initialized:', {
      hasUserInfo: !!userInfo.value,
      hasToken: !!token.value,
      isLoggedIn: isLoggedIn.value
    })
  }
}

// 延迟初始化函数 (保留向后兼容性)
function initializeIfNeeded() {
  initialize()
}

// ============ 本地操作函数 ============

/**
 * 设置用户信息
 *
 * ⚠️ 向后兼容: 函数签名保持不变
 */
export function setUserInfo(info: UserInfo) {
  userInfo.value = info
  isLoggedIn.value = true
  setStorage(StorageKeys.USER_INFO, info)
}

/**
 * 设置 Token
 *
 * ⚠️ 向后兼容: 函数签名保持不变
 */
export function setToken(newToken: string) {
  token.value = newToken
  setStorage(StorageKeys.TOKEN, newToken)
}

/**
 * 获取用户信息
 *
 * ⚠️ 向后兼容: 函数签名保持不变
 */
export function getUserInfo() {
  initializeIfNeeded()
  return userInfo.value
}

/**
 * 获取 Token
 *
 * ⚠️ 向后兼容: 函数签名保持不变
 */
export function getToken() {
  return token.value
}

/**
 * 清除用户信息(退出登录)
 *
 * ⚠️ 向后兼容: 函数签名保持不变
 */
export function clearUserInfo() {
  userInfo.value = null
  isLoggedIn.value = false
  token.value = null
  removeStorage(StorageKeys.USER_INFO)
  removeStorage(StorageKeys.TOKEN)

  // 清除邀请码缓存，防止退出登录后重复跳转到加入页面
  removeStorage(StorageKeys.PENDING_INVITE_CODE)
  removeStorage(StorageKeys.AUTO_JOIN_AFTER_LOGIN)
  removeStorage(StorageKeys.CURRENT_BABY_ID)
  removeStorage(StorageKeys.BABY_LIST)
}

/**
 * 检查登录状态
 *
 * ⚠️ 向后兼容: 函数签名保持不变
 */
export function checkLoginStatus(): boolean {
  return isLoggedIn.value && !!token.value
}

/**
 * 获取是否为新用户
 *
 * ⚠️ 向后兼容: 函数签名保持不变
 */
export function getIsNewUser() {
  return isNewUser.value
}

/**
 * 设置是否为新用户
 *
 * ⚠️ 向后兼容: 函数签名保持不变
 */
export function setIsNewUser(value: boolean) {
  isNewUser.value = value
}

// ============ API 调用函数(委托给 api 层) ============

/**
 * 微信登录 - 集成后端 API (去家庭化架构)
 *
 * API: POST /auth/wechat-login
 *
 * ⚠️ 向后兼容: 函数签名保持不变
 */
export async function wxLogin(): Promise<UserInfo> {
  return new Promise((resolve, reject) => {
    // 1. 调用 uni.login 获取 code
    uni.login({
      provider: 'weixin',
      success: async (loginRes) => {
        console.log('uni.login success:', loginRes)

        try {
          // 2. 调用后端登录接口
          const response = await authApi.apiWechatLogin({
            code: loginRes.code,
            // 注意: 由于 getUserProfile 已废弃,nickName 和 avatarUrl 可以不传
            // 或者使用头像昵称填写组件获取
          })

          console.log('API login response:', response)

          const { token: newToken, userInfo: user, isNewUser: newUser } = response

          // 3. 保存 Token 和用户信息
          setToken(newToken)
          setUserInfo(user)
          isNewUser.value = newUser

          uni.showToast({
            title: '登录成功',
            icon: 'success',
            duration: 1500,
          })

          resolve(user)
        } catch (error: any) {
          console.error('login API error:', error)
          uni.showToast({
            title: error.message || '登录失败',
            icon: 'none',
          })
          reject(error)
        }
      },
      fail: (err) => {
        console.error('uni.login fail:', err)
        uni.showToast({
          title: '微信登录失败',
          icon: 'none',
        })
        reject(err)
      },
    })
  })
}

/**
 * 刷新 Token
 *
 * API: POST /auth/refresh-token
 *
 * ⚠️ 向后兼容: 函数签名保持不变
 */
export async function refreshToken(): Promise<string> {
  try {
    const response = await authApi.apiRefreshToken()

    const newToken = response.token
    setToken(newToken)
    return newToken
  } catch (error: any) {
    console.error('refresh token error:', error)
    // Token 刷新失败,清除登录状态
    clearUserInfo()
    throw error
  }
}

/**
 * 从服务器获取用户信息
 *
 * API: GET /auth/user-info
 *
 * ⚠️ 向后兼容: 函数签名保持不变
 */
export async function fetchUserInfo(): Promise<UserInfo> {
  try {
    const response = await authApi.apiFetchUserInfo()

    setUserInfo(response)
    return response
  } catch (error: any) {
    console.error('fetch user info error:', error)
    throw error
  }
}

/**
 * 更新用户信息
 *
 * API: PUT /auth/user-info
 *
 * ⚠️ 向后兼容: 函数签名保持不变
 */
export async function updateUserInfo(nickName: string, avatarUrl: string): Promise<UserInfo> {
  try {
    const response = await authApi.apiUpdateUserInfo({ nickName, avatarUrl })

    setUserInfo(response)
    uni.showToast({
      title: '更新成功',
      icon: 'success',
    })

    return response
  } catch (error: any) {
    console.error('update user info error:', error)
    uni.showToast({
      title: error.message || '更新失败',
      icon: 'none',
    })
    throw error
  }
}

/**
 * 设置默认宝宝
 *
 * API: PUT /auth/default-baby
 *
 * ⚠️ 向后兼容: 函数签名保持不变
 */
export async function setDefaultBaby(babyId: string): Promise<void> {
  try {
    await authApi.apiSetDefaultBaby(babyId)

    // 更新本地用户信息中的默认宝宝ID
    if (userInfo.value) {
      userInfo.value.defaultBabyId = babyId
      setStorage(StorageKeys.USER_INFO, userInfo.value)
    }

    uni.showToast({
      title: '设置成功',
      icon: 'success',
    })
  } catch (error: any) {
    console.error('set default baby error:', error)
    uni.showToast({
      title: error.message || '设置失败',
      icon: 'none',
    })
    throw error
  }
}

/**
 * 退出登录
 *
 * ⚠️ 向后兼容: 函数签名保持不变
 */
export async function logout() {
  // 清除本地存储的用户信息和 Token
  clearUserInfo()

  // 跳转到登录页
  uni.reLaunch({
    url: '/pages/user/login',
  })

  uni.showToast({
    title: '已退出登录',
    icon: 'success',
  })
}

// ============ 导出 ============

export { userInfo, isLoggedIn, token, isNewUser }
