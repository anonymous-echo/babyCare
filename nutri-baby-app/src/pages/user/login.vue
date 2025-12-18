<template>
  <view class="login-container">
    <!-- 背景图片 -->
    <image class="login-bg-image" src="/static/login-background.png" mode="aspectFill" />

    <view class="login-content">
      <!-- Logo 和标题 -->
      <view class="logo-section">
        <view class="logo-wrapper">
          <image class="logo" src="/static/logo.png" mode="aspectFit" />
        </view>
        <text class="app-name">宝宝喂养时刻</text>
        <text class="app-desc">记录宝宝的每一刻</text>
      </view>

      <!-- 登录按钮 -->
      <view class="login-actions">
        <!-- 隐私协议复选框 -->
        <view class="privacy-checkbox">
          <view class="checkbox-wrapper">
            <view
              class="checkbox"
              :class="{ 'checkbox-checked': agreePrivacy }"
              @click="agreePrivacy = !agreePrivacy"
            >
              <text v-if="agreePrivacy" class="checkbox-icon">✓</text>
            </view>
            <view class="checkbox-label">
              我已阅读并同意
              <text class="link" @click.stop="showPrivacy = true">《隐私政策》</text>
            </view>
          </view>
        </view>

        <wd-button
          type="primary"
          size="large"
          block
          :loading="loading"
          :disabled="!agreePrivacy || loading"
          @click="handleLogin"
          class="login-btn"
          :class="{ 'login-btn-disabled': !agreePrivacy }"
        >
          <text v-if="!loading">微信一键登录</text>
          <text v-else>登录中...</text>
        </wd-button>
      </view>

      <!-- 隐私政策弹窗 -->
      <view v-if="showPrivacy" class="privacy-modal">
        <view class="privacy-modal-overlay" @click="showPrivacy = false"></view>
        <view class="privacy-modal-content">
          <view class="privacy-modal-header">
            <text class="privacy-modal-title">隐私政策</text>
            <text class="privacy-modal-close" @click="showPrivacy = false">×</text>
          </view>
          <scroll-view class="privacy-modal-body" scroll-y="true">
            <view class="privacy-text">
              <text class="privacy-section-title">1. 隐私信息收集</text>
              <text class="privacy-section-content">
                宝宝喂养时刻应用会收集您和宝宝的以下信息用于提供服务：
              </text>
              <text class="privacy-item">• 微信账号信息（昵称、头像、OpenID）</text>
              <text class="privacy-item">• 宝宝基本信息（姓名、性别、出生日期）</text>
              <text class="privacy-item">• 喂养、睡眠、成长等育儿记录</text>
              <text class="privacy-item">• 协作者信息（用于共同管理）</text>

              <text class="privacy-section-title">2. 信息使用范围</text>
              <text class="privacy-section-content">
                收集的信息仅用于：
              </text>
              <text class="privacy-item">• 提供育儿记录和数据统计功能</text>
              <text class="privacy-item">• 发送喂养提醒和疫苗提醒通知</text>
              <text class="privacy-item">• 改进应用功能和用户体验</text>
              <text class="privacy-item">• 维护数据安全和防止欺诈</text>

              <text class="privacy-section-title">3. 信息保护</text>
              <text class="privacy-section-content">
                我们采取以下措施保护您的信息安全：
              </text>
              <text class="privacy-item">• 使用加密传输保护数据</text>
              <text class="privacy-item">• 限制员工访问您的个人信息</text>
              <text class="privacy-item">• 定期检查和更新安全措施</text>

              <text class="privacy-section-title">4. 第三方信息共享</text>
              <text class="privacy-section-content">
                我们不会向第三方出售、出租或以其他方式共享您的个人信息，除非：
              </text>
              <text class="privacy-item">• 获得您的明确同意</text>
              <text class="privacy-item">• 法律或监管要求</text>
              <text class="privacy-item">• 微信官方平台需要（仅用于账号认证）</text>

              <text class="privacy-section-title">5. 您的权利</text>
              <text class="privacy-section-content">
                您有权：
              </text>
              <text class="privacy-item">• 访问您的个人信息</text>
              <text class="privacy-item">• 更正或删除您的信息</text>
              <text class="privacy-item">• 撤销授权和同意</text>
              <text class="privacy-item">• 注销账号</text>

              <text class="privacy-section-title">6. 政策更新</text>
              <text class="privacy-section-content">
                我们可能随时更新此隐私政策。如有重大更改，我们将以明显方式告知您。持续使用本应用表示您同意最新版本的政策。
              </text>

              <text class="privacy-section-title">7. 联系我们</text>
              <text class="privacy-section-content">
                如对隐私政策有任何疑问，请通过以下方式联系我们：
              </text>
              <text class="privacy-item">• 应用内反馈</text>
              <text class="privacy-item">• 微信客服</text>
            </view>
          </scroll-view>
          <view class="privacy-modal-footer">
            <wd-button
              type="primary"
              size="large"
              block
              @click="showPrivacy = false"
            >
              我已阅读
            </wd-button>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { wxLogin, isLoggedIn } from '@/store/user'
import { StorageKeys } from '@/utils/storage'

const loading = ref(false)
const agreePrivacy = ref(false)
const showPrivacy = ref(false)


/**
 * 登录后重定向逻辑
 *
 * 流程:
 * 1. 检查是否有待处理的邀请码 -> 跳转到加入页面
 * 2. 否则 -> 跳转到首页
 *
 * 防止无限重定向:
 * 1. 使用 switchTab 代替 reLaunch (因为首页是 tabBar 页面)
 * 2. 设置延迟确保登录状态完全保存
 */
const redirectAfterLogin = () => {
  console.log('[Login] 登录成功,检查重定向目标')

  // 检查是否有待处理的邀请码（从扫码进入但未登录的场景）
  const pendingCode = uni.getStorageSync(StorageKeys.PENDING_INVITE_CODE)
  console.log('[Login] Checking PENDING_INVITE_CODE:', pendingCode)

  if (pendingCode) {
    console.log('[Login] 检测到待处理的邀请码,跳转到加入页面:', pendingCode)

    // 清除缓存（避免重复跳转）
    uni.removeStorageSync(StorageKeys.PENDING_INVITE_CODE)

    // 跳转回加入页面
    uni.reLaunch({
      url: `/pages/baby/join/join?code=${pendingCode}`,
    })

    return
  }

  // 正常跳转到首页
  console.log('[Login] 登录成功,跳转到首页')

  // 对于 tabBar 页面,应该使用 switchTab 而不是 reLaunch
  // 这样可以避免页面重新加载和生命周期冲突
  uni.switchTab({
    url: '/pages/index/index',
    fail: (err) => {
      console.error('[Login] 跳转失败,使用 reLaunch 降级:', err)
      // 如果 switchTab 失败,降级使用 reLaunch
      uni.reLaunch({
        url: '/pages/index/index'
      })
    }
  })
}

// 处理登录 - 必须在用户点击事件中调用
const handleLogin = async () => {
  if (loading.value) return

  loading.value = true

  try {
    // 微信登录
    await wxLogin()

    // 延迟跳转,让用户看到成功提示
    setTimeout(() => {
      redirectAfterLogin()
    }, 1500)

  } catch (error) {
    console.error('登录失败:', error)
  } finally {
    loading.value = false
  }
}
</script>

<style lang="scss" scoped>
.login-container {
  position: relative;
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40rpx 20rpx;
  overflow: hidden;
}

.login-bg-image {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 0;
}

.login-content {
  position: relative;
  z-index: 10;
  width: 100%;
  max-width: 420rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
}

// Logo 部分
.logo-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 120rpx;
  text-align: center;
}

.logo-wrapper {
  position: relative;
  width: 200rpx;
  height: 200rpx;
  margin-bottom: 40rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo {
  width: 200rpx;
  height: 200rpx;
  border-radius: 100%;
  background: white;
  object-fit: contain;
  box-shadow: 0 8rpx 24rpx rgba(0, 0, 0, 0.15);
}

.app-name {
  font-size: 48rpx;
  font-weight: bold;
  color: white;
  margin-bottom: 16rpx;
  letter-spacing: 2rpx;
}

.app-desc {
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.9);
  font-weight: 500;
}

// 登录操作区域
.login-actions {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 20rpx;
}

.login-btn {
  :deep(.wd-button) {
    background: white !important;
    color: #1a9e4b !important;
    font-weight: 600;
    font-size: 32rpx;
    border-radius: 16rpx;
    height: 88rpx;
    box-shadow: 0 8rpx 24rpx rgba(0, 0, 0, 0.15);

    &:active {
      background: rgba(255, 255, 255, 0.9) !important;
    }
  }
}

// 隐私协议复选框
.privacy-checkbox {
  width: 100%;
  margin-bottom: 20rpx;
  padding: 0 10rpx;
}

.checkbox-wrapper {
  display: flex;
  align-items: center;
  gap: 12rpx;
}

.checkbox {
  width: 36rpx;
  height: 36rpx;
  min-width: 36rpx;
  border: 2rpx solid rgba(255, 255, 255, 0.6);
  border-radius: 6rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;

  &:active {
    border-color: white;
    background: rgba(255, 255, 255, 0.1);
  }

  &.checkbox-checked {
    background: white;
    border-color: white;
  }
}

.checkbox-icon {
  font-size: 24rpx;
  color: #1a9e4b;
  font-weight: bold;
}

.checkbox-label {
  flex: 1;
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.9);
  line-height: 36rpx;

  .link {
    color: white;
    text-decoration: underline;
    font-weight: 500;
  }
}

.login-btn-disabled {
  :deep(.wd-button) {
    background: rgba(255, 255, 255, 0.5) !important;
    color: rgba(26, 158, 75, 0.5) !important;
    cursor: not-allowed;
  }
}

// 隐私政策弹窗
.privacy-modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 1000;
  display: flex;
  align-items: flex-end;
}

.privacy-modal-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.6);
  z-index: 1;
}

.privacy-modal-content {
  position: relative;
  z-index: 2;
  width: 100%;
  height: 80vh;
  max-height: 80%;
  background: white;
  border-radius: 24rpx 24rpx 0 0;
  display: flex;
  flex-direction: column;
  animation: slideUp 0.3s ease;
}

@keyframes slideUp {
  from {
    transform: translateY(100%);
  }
  to {
    transform: translateY(0);
  }
}

.privacy-modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24rpx 24rpx;
  border-bottom: 1rpx solid #f0f0f0;
  min-height: 60rpx;
}

.privacy-modal-title {
  font-size: 32rpx;
  font-weight: bold;
  color: #333;
}

.privacy-modal-close {
  font-size: 48rpx;
  color: #999;
  cursor: pointer;
  line-height: 48rpx;
  width: 48rpx;
  text-align: center;

  &:active {
    color: #666;
  }
}

.privacy-modal-body {
  flex: 1;
  overflow-y: auto;
  padding: 24rpx;
}

.privacy-text {
  display: flex;
  flex-direction: column;
  gap: 16rpx;
  padding: 0 12rpx;
}

.privacy-section-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #333;
  margin-top: 8rpx;
}

.privacy-section-content {
  font-size: 24rpx;
  color: #666;
  line-height: 1.6;
}

.privacy-item {
  font-size: 24rpx;
  color: #666;
  line-height: 1.6;
  margin-left: 20rpx;
}

.privacy-modal-footer {
  padding: 16rpx 24rpx 32rpx;
  border-top: 1rpx solid #f0f0f0;
  background: white;

  :deep(.wd-button) {
    height: 88rpx;
    font-size: 32rpx;
  }
}
</style>