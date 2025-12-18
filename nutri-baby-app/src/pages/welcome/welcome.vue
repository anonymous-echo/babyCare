<template>
  <view class="welcome-container">
    <!-- 欢迎标题 -->
    <view class="welcome-header">
      <image
        class="logo"
        src="/static/logo.png"
        mode="aspectFit"
      />
      <text class="title">欢迎使用宝宝喂养日志</text>
      <text class="subtitle">记录宝宝成长的每一个精彩瞬间</text>
    </view>

    <!-- 引导卡片 -->
    <view class="guide-card">
      <text class="guide-title">开始使用前,请先选择:</text>

      <view class="options">
        <!-- 创建宝宝选项 -->
        <view class="option-card" @click="handleCreateBaby">
          <view class="option-icon">
            <text class="icon">👶</text>
          </view>
          <view class="option-content">
            <text class="option-title">创建我的宝宝</text>
            <text class="option-desc">为您的宝宝创建成长档案</text>
          </view>
          <view class="option-arrow">
            <text class="arrow">›</text>
          </view>
        </view>

        <!-- 加入协作选项 -->
        <view class="option-card" @click="handleJoinBaby">
          <view class="option-icon">
            <text class="icon">🤝</text>
          </view>
          <view class="option-content">
            <text class="option-title">加入协作</text>
            <text class="option-desc">输入邀请码加入现有宝宝</text>
          </view>
          <view class="option-arrow">
            <text class="arrow">›</text>
          </view>
        </view>
      </view>
    </view>

    <!-- 特性介绍 -->
    <view class="features">
      <view class="feature-item">
        <text class="feature-icon">📊</text>
        <text class="feature-text">全面记录喂养、睡眠、成长数据</text>
      </view>
      <view class="feature-item">
        <text class="feature-icon">👨‍👩‍👧</text>
        <text class="feature-text">支持多人协作共同照护</text>
      </view>
      <view class="feature-item">
        <text class="feature-icon">📱</text>
        <text class="feature-text">数据云端同步,随时随地访问</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { joinBabyCollaboration } from '@/store/collaborator'
import { fetchBabyList } from '@/store/baby'

// 处理创建宝宝
const handleCreateBaby = () => {
  uni.navigateTo({
    url: '/pages/baby/edit/edit'
  })
}

// 处理加入协作 (去家庭化架构)
const handleJoinBaby = () => {
  uni.showModal({
    title: '加入宝宝协作',
    content: '请扫描二维码或点击微信分享链接加入',
    showCancel: true,
    cancelText: '取消',
    confirmText: '手动输入',
    success: (modalRes) => {
      if (modalRes.confirm) {
        // 手动输入邀请码
        uni.showModal({
          title: '输入邀请信息',
          editable: true,
          placeholderText: '格式: babyId,token',
          success: async (res) => {
            if (res.confirm && res.content) {
              try {
                // 解析输入: babyId,token
                const [babyId, token] = res.content.split(',').map(s => s.trim())

                if (!babyId || !token) {
                  throw new Error('格式错误,请输入: babyId,token')
                }

                // 调用加入协作 API
                await joinBabyCollaboration(babyId, token)

                uni.showToast({
                  title: '加入成功',
                  icon: 'success',
                  duration: 2000
                })

                // 刷新宝宝列表
                await fetchBabyList()

                // 跳转到首页
                setTimeout(() => {
                  uni.reLaunch({
                    url: '/pages/index/index'
                  })
                }, 2000)
              } catch (error: any) {
                uni.showToast({
                  title: error.message || '加入失败',
                  icon: 'none',
                  duration: 2000
                })
              }
            }
          }
        })
      }
    }
  })
}

onLoad(() => {
  console.log('[Welcome] 欢迎页面加载 (去家庭化架构)')
})
</script>

<style lang="scss" scoped>
.welcome-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 60rpx 40rpx;
  display: flex;
  flex-direction: column;
}

.welcome-header {
  text-align: center;
  margin-bottom: 80rpx;

  .logo {
    width: 160rpx;
    height: 160rpx;
    margin-bottom: 40rpx;
    border-radius: 50%;
    background-color: rgba(255, 255, 255, 0.2);
  }

  .title {
    display: block;
    font-size: 48rpx;
    font-weight: bold;
    color: #ffffff;
    margin-bottom: 20rpx;
  }

  .subtitle {
    display: block;
    font-size: 28rpx;
    color: rgba(255, 255, 255, 0.9);
  }
}

.guide-card {
  background-color: #ffffff;
  border-radius: 24rpx;
  padding: 40rpx;
  box-shadow: 0 8rpx 32rpx rgba(0, 0, 0, 0.1);
  margin-bottom: 60rpx;

  .guide-title {
    display: block;
    font-size: 32rpx;
    font-weight: 600;
    color: #333333;
    margin-bottom: 40rpx;
  }
}

.options {
  display: flex;
  flex-direction: column;
  gap: 24rpx;
}

.option-card {
  display: flex;
  align-items: center;
  padding: 32rpx 24rpx;
  background-color: #f8f9fa;
  border-radius: 16rpx;
  transition: all 0.3s ease;

  &:active {
    background-color: #e9ecef;
    transform: scale(0.98);
  }

  .option-icon {
    width: 80rpx;
    height: 80rpx;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border-radius: 16rpx;
    margin-right: 24rpx;

    .icon {
      font-size: 48rpx;
    }
  }

  .option-content {
    flex: 1;
    display: flex;
    flex-direction: column;

    .option-title {
      font-size: 32rpx;
      font-weight: 600;
      color: #333333;
      margin-bottom: 8rpx;
    }

    .option-desc {
      font-size: 24rpx;
      color: #999999;
    }
  }

  .option-arrow {
    .arrow {
      font-size: 48rpx;
      color: #cccccc;
      font-weight: 300;
    }
  }
}

.features {
  display: flex;
  flex-direction: column;
  gap: 24rpx;
}

.feature-item {
  display: flex;
  align-items: center;
  padding: 24rpx;
  background-color: rgba(255, 255, 255, 0.15);
  border-radius: 16rpx;
  backdrop-filter: blur(10rpx);

  .feature-icon {
    font-size: 40rpx;
    margin-right: 20rpx;
  }

  .feature-text {
    font-size: 28rpx;
    color: rgba(255, 255, 255, 0.95);
  }
}
</style>
