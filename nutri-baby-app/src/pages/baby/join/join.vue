<template>
<view>
<wd-navbar title="邀请详情" left-text="返回" right-text="设置" left-arrow safeAreaInsetTop fixed>
  <template #capsule>
    <wd-navbar-capsule @back="handleBack" @back-home="goToHome" height="auto" />
  </template>
</wd-navbar>
  <view class="join-container">
    <!-- 内容区域（留出导航栏高度） -->
    <view class="content-wrapper" :style="{ paddingTop: navbarTotalHeight + 'rpx' }">
      <!-- 加载状态 -->
      <view v-if="loading" class="loading-wrapper">
        <view class="loading-spinner"></view>
        <text class="loading-text">加载中...</text>
      </view>

      <!-- 邀请信息展示 -->
      <view v-else-if="invitationInfo" class="content">
        <!-- 宝宝信息卡片 -->
        <view class="baby-card">
          <view class="baby-avatar">
            <image v-if="invitationInfo.babyAvatar" :src="invitationInfo.babyAvatar" mode="aspectFill" />
            <text v-else class="avatar-placeholder">👶</text>
          </view>

          <view class="baby-info">
            <view class="baby-name">{{ invitationInfo.babyName }}</view>
            <view class="inviter-info">
              <text class="inviter-name">{{ invitationInfo.inviterName }}</text>
              <text class="invite-text">邀请你一起记录宝宝的成长</text>
            </view>
          </view>
        </view>

        <!-- 权限信息 -->
        <view class="permission-card">
          <view class="card-title">协作权限</view>
          <view class="permission-list">
            <view class="permission-item">
              <text class="label">协作角色:</text>
              <text class="value">{{ roleText }}</text>
            </view>
            <view class="permission-item">
              <text class="label">访问权限:</text>
              <text class="value">{{ accessTypeText }}</text>
            </view>
            <view v-if="invitationInfo.expiresAt" class="permission-item">
              <text class="label">权限过期:</text>
              <text class="value">{{ formatExpireTime }}</text>
            </view>
          </view>

          <!-- 权限说明 -->
          <view class="permission-desc">
            <text v-if="role === 'admin'">管理员可管理宝宝信息、邀请/移除协作者</text>
            <text v-else-if="role === 'editor'">编辑者可记录和编辑所有数据</text>
            <text v-else>查看者仅可查看数据,不能编辑</text>
          </view>
        </view>

        <!-- 操作按钮 -->
        <view class="actions">
          <wd-button type="primary" size="large" @click="handleJoin">
            确认加入
          </wd-button>
          <wd-button type="info" size="large" @click="handleCancel">
            取消
          </wd-button>
        </view>

        <!-- 温馨提示 -->
        <view class="tips">
          <view class="tip-title">温馨提示</view>
          <view class="tip-item">• 加入后可与家人共同记录宝宝的成长</view>
          <view class="tip-item">• 所有协作者的记录将实时同步</view>
          <view class="tip-item">• 请谨慎选择协作者,保护宝宝隐私</view>
        </view>
      </view>

      <!-- 错误状态 -->
      <view v-else class="error-wrapper">
        <text class="error-icon">⚠️</text>
        <text class="error-text">{{ errorMessage }}</text>
        <wd-button type="primary" @click="handleBack">返回</wd-button>
      </view>
    </view>
  </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { onLoad, onShow } from '@dcloudio/uni-app'
import { joinBabyCollaboration } from '@/store/collaborator'
import { apiGetInvitationByCode } from '@/api/baby'
import { checkLoginStatus } from '@/store/user'
import { StorageKeys } from '@/utils/storage'

// 导航栏相关
const statusBarHeight = ref(0)
const navbarTotalHeight = ref(0)
const navbarContentHeight = ref(88) // 导航栏内容高度（rpx）
const menuButtonWidth = ref(0) // 胶囊按钮宽度（rpx）
const menuButtonHeight = ref(0) // 胶囊按钮高度（rpx）
const menuButtonTop = ref(0) // 胶囊按钮顶部距离（px）

// 页面参数
const babyId = ref('')
const token = ref('')
const shortCode = ref('') // 新增短码参数

// 页面状态
const loading = ref(true)
const invitationInfo = ref<any>(null)
const errorMessage = ref('')
const role = ref('')

// 角色文本映射
const roleTextMap: Record<string, string> = {
  admin: '管理员',
  editor: '编辑者',
  viewer: '查看者',
}

// 访问权限文本映射
const accessTypeTextMap: Record<string, string> = {
  permanent: '永久有效',
  temporary: '临时权限',
}

// 计算属性
const roleText = computed(() => roleTextMap[role.value] || '编辑者')
const accessTypeText = computed(() => {
  return invitationInfo.value?.accessType
    ? accessTypeTextMap[invitationInfo.value.accessType]
    : '永久有效'
})

const formatExpireTime = computed(() => {
  if (!invitationInfo.value?.expiresAt) return ''
  const date = new Date(invitationInfo.value.expiresAt)
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
})

// 初始化导航栏高度
onMounted(() => {
  const systemInfo = uni.getSystemInfoSync()
  statusBarHeight.value = systemInfo.statusBarHeight || 0

  // 获取胶囊按钮位置信息（微信小程序）
  // #ifdef MP-WEIXIN
  try {
    const menuButton = uni.getMenuButtonBoundingClientRect()
    if (menuButton) {
      // 胶囊按钮的宽度和高度（保持 px，与导航栏样式中使用 rpx 统一处理）
      menuButtonWidth.value = menuButton.width // px
      menuButtonHeight.value = menuButton.height // px
      menuButtonTop.value = menuButton.top // px（状态栏下的距离）

      // 计算导航栏内容高度（rpx）
      // navbar-content 的高度应该 = 胶囊按钮的高度
      navbarContentHeight.value = Math.round(menuButton.height * 2) // 胶囊高度转为 rpx

      // 总高度 = 状态栏高度 × 2（px→rpx） + 胶囊顶部距离 × 2（px→rpx） + 导航栏内容高度
      navbarTotalHeight.value =
        Math.round(statusBarHeight.value * 2) +
        Math.round(menuButton.top) 
      console.log('[Join Navbar] Capsule aligned:', {
        statusBarHeight: statusBarHeight.value,
        menuButtonTop: menuButtonTop.value,
        menuButtonWidth: menuButton.width,
        menuButtonHeight: menuButton.height,
        menuButtonBottom: menuButton.top + menuButton.height,
        navbarContentHeight: navbarContentHeight.value,
        navbarTotalHeight: navbarTotalHeight.value,
      })
    }
  } catch (e) {
    console.warn('[Join Navbar] Failed to get menu button info, using defaults', e)
    // 使用默认值
    navbarContentHeight.value = 88
    navbarTotalHeight.value = Math.round(statusBarHeight.value * 2)
  }
  // #endif
})

// 返回首页
function goToHome() {
  console.log('[Join] 点击首页图标，跳转到首页')
  uni.switchTab({
    url: '/pages/index/index',
  })
}

// 页面加载
onLoad((options) => {
  console.log('Join page onLoad with options:', options)

  // 支持两种参数格式：
  // 1. 扫码进入: ?code=ABC123
  // 2. 分享链接: ?babyId=xxx&token=xxx
  if (options?.code) {
    shortCode.value = options.code
  }

  if (options?.babyId) {
    babyId.value = options.babyId
  }
  if (options?.token) {
    token.value = options.token
  }

  // 加载邀请信息
  loadInvitationInfo()
})

// 页面显示（从登录页返回时会触发）
onShow(() => {
  console.log('Join page onShow')

  // 检查是否有待加入的邀请（从登录页返回）
  const autoJoin = uni.getStorageSync(StorageKeys.AUTO_JOIN_AFTER_LOGIN)

  if (autoJoin && checkLoginStatus()) {
    console.log('Auto join after login:', autoJoin)

    // 清除标记
    uni.removeStorageSync(StorageKeys.AUTO_JOIN_AFTER_LOGIN)

    // 恢复邀请信息
    babyId.value = autoJoin.babyId
    token.value = autoJoin.token
    invitationInfo.value = autoJoin.invitationInfo
    role.value = autoJoin.role
    loading.value = false

    // 自动执行加入操作
    setTimeout(() => {
      handleJoin()
    }, 500)
  }
})

// 加载邀请信息
async function loadInvitationInfo() {
  // 优先使用短码方式
  if (shortCode.value) {
    await loadInvitationByShortCode()
  } else if (babyId.value && token.value) {
    await loadInvitationByToken()
  } else {
    errorMessage.value = '邀请链接无效,缺少必要参数'
    loading.value = false
  }
}

// 通过短码加载邀请信息
async function loadInvitationByShortCode() {
  try {
    const response = await apiGetInvitationByCode(shortCode.value)

    console.log('Invitation loaded by short code:', response)

    invitationInfo.value = {
      babyId: response.babyId,
      babyName: response.babyName,
      babyAvatar: response.babyAvatar,
      inviterName: response.inviterName,
      role: response.role,
      accessType: response.accessType,
      expiresAt: response.expiresAt,
    }

    // 保存 babyId 和 token 用于后续加入操作
    babyId.value = response.babyId
    token.value = response.token
    role.value = response.role

    loading.value = false
  } catch (error: any) {
    console.error('Load invitation by short code error:', error)
    errorMessage.value = error.message || '邀请码无效或已过期'
    loading.value = false
  }
}

// 通过 token 加载邀请信息（旧方式，保持兼容）
async function loadInvitationByToken() {
  // 模拟数据（保持原有逻辑）
  setTimeout(() => {
    invitationInfo.value = {
      babyId: babyId.value,
      babyName: '小明',
      babyAvatar: '',
      inviterName: '爸爸',
      role: 'editor',
      accessType: 'permanent',
      expiresAt: null,
    }
    role.value = invitationInfo.value.role
    loading.value = false
  }, 500)
}

// 确认加入
async function handleJoin() {
  if (!babyId.value || !token.value) {
    uni.showToast({
      title: '邀请信息不完整',
      icon: 'none',
    })
    return
  }

  // 检查登录状态
  if (!checkLoginStatus()) {
    console.log('[Join] User not logged in, redirect to login page')
    console.log('[Join] Saving shortCode:', shortCode.value)

    // 保存邀请信息到本地存储
    uni.setStorageSync(StorageKeys.PENDING_INVITE_CODE, shortCode.value)
    uni.setStorageSync(StorageKeys.AUTO_JOIN_AFTER_LOGIN, {
      babyId: babyId.value,
      token: token.value,
      invitationInfo: invitationInfo.value,
      role: role.value,
    })

    // 验证保存是否成功
    const saved = uni.getStorageSync(StorageKeys.PENDING_INVITE_CODE)
    console.log('[Join] Verification - saved PENDING_INVITE_CODE:', saved)

    // 提示并跳转到登录页
    uni.showModal({
      title: '需要登录',
      content: '请先登录后再加入宝宝协作',
      showCancel: false,
      success: () => {
        uni.reLaunch({
          url: '/pages/user/login',
        })
      },
    })

    return
  }

  uni.showLoading({
    title: '加入中...',
  })

  try {
    // 调用加入API
    const result = await joinBabyCollaboration(babyId.value, token.value)

    uni.hideLoading()

    // 加入成功，清除缓存的邀请码，防止后续重复跳转
    console.log('[Join] 加入成功，清除邀请码缓存')
    uni.removeStorageSync(StorageKeys.PENDING_INVITE_CODE)
    uni.removeStorageSync(StorageKeys.AUTO_JOIN_AFTER_LOGIN)

    // 显示成功提示
    uni.showModal({
      title: '加入成功',
      content: `你已成功加入${result.name}的协作团队`,
      showCancel: false,
      success: () => {
        // 跳转到宝宝列表页
        uni.reLaunch({
          url: '/pages/baby/list/list',
        })
      },
    })
  } catch (error: any) {
    uni.hideLoading()
    console.error('join error:', error)
    // 错误已在 joinBabyCollaboration 中处理
  }
}

// 取消加入
function handleCancel() {
  uni.showModal({
    title: '确认取消',
    content: '确定要取消加入吗?',
    success: (res) => {
      if (res.confirm) {
        handleBack()
      }
    },
  })
}

// 返回
function handleBack() {
  console.log('[Join] 取消加入，清除邀请码缓存')

  // 清除缓存的邀请码，防止后续重复跳转
  uni.removeStorageSync(StorageKeys.PENDING_INVITE_CODE)
  uni.removeStorageSync(StorageKeys.AUTO_JOIN_AFTER_LOGIN)

  // 如果是从分享链接进入,返回首页
  uni.reLaunch({
    url: '/pages/index/index',
  })
}
</script>

<style lang="scss" scoped>
.join-container {
  min-height: 100vh;
  background: #f6f8f7;
}


.navbar-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20rpx;
  height: 88rpx; // 会被动态样式覆盖
  margin-top: 40rpx;
}

.navbar-left {
  display: flex;
  align-items: center;
  justify-content: left;
  // cursor: pointer;
  transition: opacity 0.3s;

  &:active {
    opacity: 0.7;
  }
}

.navbar-title {
  flex: 1;
  text-align: center;
  font-size: 32rpx;
  font-weight: bold;
  color: white;
  letter-spacing: 2rpx;
}

.navbar-right {
  flex-shrink: 0;
}

// 内容区域样式 - 避免被导航栏遮挡
.content-wrapper {
  padding: 40rpx 20rpx 20rpx;
  padding-top: 50rpx; // 会被动态样式覆盖
}

.loading-wrapper,
.error-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 60vh;
  color: #333;

  .loading-spinner {
    width: 80rpx;
    height: 80rpx;
    border: 6rpx solid rgba(125, 211, 162, 0.3);
    border-top-color: #7dd3a2;
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }

  .loading-text,
  .error-text {
    margin: 20rpx 0;
    font-size: 32rpx;
  }

  .error-icon {
    font-size: 120rpx;
    margin-bottom: 20rpx;
  }
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.content {
  margin-top: 20rpx ;
  .baby-card {
    background: white;
    border-radius: 24rpx;
    padding: 40rpx;
    margin-bottom: 20rpx;
    display: flex;
    align-items: center;
    box-shadow: 0 8rpx 24rpx rgba(0, 0, 0, 0.1);

    .baby-avatar {
      width: 120rpx;
      height: 120rpx;
      border-radius: 60rpx;
      background: linear-gradient(135deg, #ffecd2 0%, #fcb69f 100%);
      display: flex;
      align-items: center;
      justify-content: center;
      margin-right: 24rpx;
      overflow: hidden;

      image {
        width: 100%;
        height: 100%;
      }

      .avatar-placeholder {
        font-size: 60rpx;
      }
    }

    .baby-info {
      flex: 1;

      .baby-name {
        font-size: 36rpx;
        font-weight: bold;
        color: #333;
        margin-bottom: 12rpx;
      }

      .inviter-info {
        font-size: 28rpx;
        color: #666;

        .inviter-name {
          color: #7dd3a2;
          font-weight: 500;
          margin-right: 8rpx;
        }
      }
    }
  }

  .permission-card {
    background: white;
    border-radius: 24rpx;
    padding: 40rpx;
    margin-bottom: 20rpx;
    box-shadow: 0 8rpx 24rpx rgba(0, 0, 0, 0.1);

    .card-title {
      font-size: 32rpx;
      font-weight: bold;
      color: #333;
      margin-bottom: 24rpx;
    }

    .permission-list {
      .permission-item {
        display: flex;
        justify-content: space-between;
        padding: 16rpx 0;
        font-size: 28rpx;

        .label {
          color: #999;
        }

        .value {
          color: #333;
          font-weight: 500;
        }
      }
    }

    .permission-desc {
      margin-top: 20rpx;
      padding: 20rpx;
      background: #f0f9f6;
      border-radius: 12rpx;
      border: 1rpx solid #CAE3D4;
      font-size: 26rpx;
      color: #7dd3a2;
      line-height: 1.6;
    }
  }

  .actions {
    display: flex;
    flex-direction: column;
    gap: 16rpx;
    margin-bottom: 20rpx;
  }

  .tips {
    background: rgba(125, 211, 162, 0.1);
    border-radius: 16rpx;
    padding: 30rpx;
    color: #333;
    border: 1rpx solid rgba(125, 211, 162, 0.2);

    .tip-title {
      font-size: 28rpx;
      font-weight: bold;
      margin-bottom: 16rpx;
      color: #7dd3a2;
    }

    .tip-item {
      font-size: 26rpx;
      line-height: 1.8;
      color: #666;
    }
  }
}
</style>
