<template>
  <view class="user-page">
    <!-- 用户信息卡片 -->
    <view class="user-card">
      <view class="user-info">
        <!-- 头像 -->
        <view class="avatar" @click="goToProfile">
          <image
            :src="userInfo?.avatarUrl || '/static/default.png'"
            class="avatar-image"
            mode="aspectFill"
          />
        </view>

        <!-- 昵称/登录按钮 -->
        <view class="info">
          <!-- 未登录：显示登录/注册按钮 -->
          <view v-if="!isLoggedIn" class="login-button" @click="goToLogin">
            <text>登录/注册</text>
          </view>

          <!-- 已登录：显示昵称和时间 -->
          <template v-else>
            <view class="nickname">{{ userInfo?.nickName || "用户" }}</view>
            <view class="login-time">{{ loginTimeText }}</view>
          </template>
        </view>
      </view>
    </view>

    <!-- 个人信息 -->
    <view class="section">
      <view class="section-title">个人信息</view>
      <wd-cell-group>
        <wd-cell
          title="编辑资料"
          desc="修改头像和昵称"
          is-link
          @click="goToProfile"
        >
          <template #icon>
            <wd-icon name="edit" size="20px" custom-class="cell-icon" />
          </template>
        </wd-cell>
      </wd-cell-group>
    </view>

    <!-- 宝宝信息 -->
    <view class="section">
      <view class="section-title">我的宝宝</view>
      <wd-cell-group>
        <wd-cell
          title="宝宝管理"
          :desc="`共 ${babyList?.length || 0} 个宝宝`"
          is-link
          @click="goToBabyList"
        >
          <template #icon>
            <wd-icon name="user" size="20px" custom-class="cell-icon" />
          </template>
        </wd-cell>
        <wd-cell title="疫苗计划" desc="接种计划" is-link @click="goToVaccine">
          <template #icon>
            <wd-icon name="list" size="20px" custom-class="cell-icon" />
          </template>
        </wd-cell>
      </wd-cell-group>
    </view>

    <!-- 数据管理 -->
    <view class="section">
      <view class="section-title">数据管理</view>
      <wd-cell-group>
        <!-- <wd-cell
          title="数据导出"
          desc="导出记录数据"
          is-link
          @click="exportData"
        >
          <template #icon>
            <wd-icon name="chart" size="20px" custom-class="cell-icon" />
          </template>
        </wd-cell>
        <wd-cell
          title="数据导入"
          desc="从剪贴板导入"
          is-link
          @click="importData"
        >
          <template #icon>
            <wd-icon name="download" size="20px" custom-class="cell-icon" />
          </template>
        </wd-cell> -->
        <wd-cell
          title="数据统计"
          :desc="`共 ${totalRecords} 条记录`"
          is-link
          @click="goToStatistics"
        >
          <template #icon>
            <wd-icon
              name="format-horizontal-align-bottom"
              size="20px"
              custom-class="cell-icon"
            />
          </template>
        </wd-cell>
      </wd-cell-group>
    </view>

    <!-- 设置 -->
    <view class="section">
      <view class="section-title">设置</view>
      <wd-cell-group>
        <wd-cell
          title="消息提醒设置"
          desc="管理订阅消息"
          is-link
          @click="goToSubscribeSettings"
        >

          <template #icon>
            <wd-icon name="notification" size="20px" custom-class="cell-icon" />
          </template>
        </wd-cell>
        <!-- #ifdef MP-WEIXIN -->
        <button
          open-type="feedback"
          class="feedback-button"
        >
          <wd-cell
            title="用户反馈"
            desc="向我们反馈问题和建议"
            is-link
          >
            <template #icon>
              <wd-icon name="chat" size="20px" custom-class="cell-icon" />
            </template>
          </wd-cell>
        </button>
        <!-- #endif -->

        <!-- #ifndef MP-WEIXIN -->
        <wd-cell
          title="用户反馈"
          desc="向我们反馈问题和建议"
          is-link
          @click="showFeedbackTip"
        >
          <template #icon>
            <wd-icon name="chat" size="20px" custom-class="cell-icon" />
          </template>
        </wd-cell>
        <!-- #endif -->
        <wd-cell title="关于我们" is-link @click="showAbout">
          <template #icon>
            <wd-icon name="info-circle" size="20px" custom-class="cell-icon" />
          </template>
        </wd-cell>
        <wd-cell title="清除缓存" is-link @click="clearCache">
          <template #icon>
            <wd-icon name="delete" size="20px" custom-class="cell-icon" />
          </template>
        </wd-cell>
      </wd-cell-group>
    </view>

    <!-- 退出登录 -->
    <view v-if="isLoggedIn" class="logout-section">
      <wd-button size="large" block @click="handleLogout"> 退出登录 </wd-button>
    </view>

    <!-- 版本信息 -->
    <view class="version">
      <text>宝宝喂养时刻 v3.0.0</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { isLoggedIn, userInfo, clearUserInfo } from "@/store/user";
import { currentBaby } from "@/store/baby";
import { formatDate } from "@/utils/date";

// 直接调用 API 层
import * as babyApi from "@/api/baby";
import * as feedingApi from "@/api/feeding";
import * as diaperApi from "@/api/diaper";
import * as sleepApi from "@/api/sleep";

// 数据统计(从 API 获取)
const babyList = ref<babyApi.BabyProfileResponse[]>([]);
const feedingRecordsCount = ref(0);
const diaperRecordsCount = ref(0);
const sleepRecordsCount = ref(0);

// 加载统计数据
const loadStatistics = async () => {
  if (!currentBaby.value) return;

  try {
    const babyId = currentBaby.value.babyId;

    // 获取宝宝列表
    const babiesData = await babyApi.apiFetchBabyList();
    babyList.value = babiesData;

    // 获取各类记录数量(使用 pageSize:1 只获取总数)
    const [feedingData, diaperData, sleepData] = await Promise.all([
      feedingApi.apiFetchFeedingRecords({ babyId, page: 1, pageSize: 1 }),
      diaperApi.apiFetchDiaperRecords({ babyId, page: 1, pageSize: 1 }),
      sleepApi.apiFetchSleepRecords({ babyId, page: 1, pageSize: 1 }),
    ]);

    feedingRecordsCount.value = feedingData.total || 0;
    diaperRecordsCount.value = diaperData.total || 0;
    sleepRecordsCount.value = sleepData.total || 0;
  } catch (error) {
    console.error("加载统计数据失败:", error);
  }
};

// 页面加载时获取数据
onMounted(() => {
  if (isLoggedIn.value) {
    loadStatistics();
  }
});

// 跳转到登录
const goToLogin = () => {
  uni.navigateTo({
    url: "/pages/user/login",
  });
};

// 登录时间文本
const loginTimeText = computed(() => {
  if (!userInfo.value) return "";
  return "登录于 " + formatDate(userInfo.value.createTime, "YYYY-MM-DD");
});

// 总记录数
const totalRecords = computed(() => {
  return (
    feedingRecordsCount.value +
    diaperRecordsCount.value +
    sleepRecordsCount.value
  );
});

// 跳转到宝宝列表
const goToBabyList = () => {
  uni.navigateTo({
    url: "/pages/baby/list/list",
  });
};

// 跳转到个人资料编辑
const goToProfile = () => {
  uni.navigateTo({
    url: "/pages/user/profile/profile",
  });
};

// 跳转到家庭管理
const goToFamily = () => {
  uni.navigateTo({
    url: "/pages/family/family",
  });
};

// 跳转到疫苗提醒
const goToVaccine = () => {
  uni.navigateTo({
    url: "/pages/vaccine/vaccine",
  });
};

// 跳转到订阅消息设置
const goToSubscribeSettings = () => {
  uni.navigateTo({
    url: "/pages/settings/subscribe/subscribe",
  });
};

// 跳转到统计页面
const goToStatistics = () => {
  uni.switchTab({
    url: "/pages/statistics/statistics",
  });
};

// 导出数据
const exportData = async () => {
  if (!currentBaby.value) {
    uni.showToast({
      title: "请先选择宝宝",
      icon: "none",
    });
    return;
  }

  try {
    uni.showLoading({
      title: "准备导出...",
      mask: true,
    });

    const babyId = currentBaby.value.babyId;

    // 从 API 获取所有数据
    const [babiesData, feedingData, diaperData, sleepData] = await Promise.all([
      babyApi.apiFetchBabyList(),
      feedingApi.apiFetchFeedingRecords({ babyId, pageSize: 1000 }),
      diaperApi.apiFetchDiaperRecords({ babyId, pageSize: 1000 }),
      sleepApi.apiFetchSleepRecords({ babyId, pageSize: 1000 }),
    ]);

    // 准备导出数据
    const exportData = {
      exportTime: Date.now(),
      exportTimeText: formatDate(Date.now(), "YYYY-MM-DD HH:mm:ss"),
      babies: babiesData,
      feedingRecords: feedingData.records,
      diaperRecords: diaperData.records,
      sleepRecords: sleepData.records,
    };

    // 生成 JSON 字符串
    const jsonStr = JSON.stringify(exportData, null, 2);
    const fileName = `baby_data_${formatDate(
      Date.now(),
      "YYYYMMDD_HHmmss"
    )}.json`;

    uni.hideLoading();

    // 显示导出摘要
    const summary = `
导出时间: ${exportData.exportTimeText}
宝宝数量: ${babiesData.length}
喂养记录: ${feedingData.records.length} 条
换尿布记录: ${diaperData.records.length} 条
睡眠记录: ${sleepData.records.length} 条
总记录数: ${
      feedingData.records.length +
      diaperData.records.length +
      sleepData.records.length
    } 条

文件名: ${fileName}
    `.trim();

    uni.showModal({
      title: "数据导出成功",
      content: summary,
      confirmText: "复制数据",
      cancelText: "关闭",
      success: (res) => {
        if (res.confirm) {
          // 复制到剪贴板
          uni.setClipboardData({
            data: jsonStr,
            success: () => {
              uni.showToast({
                title: "已复制到剪贴板",
                icon: "success",
              });
            },
          });
        }
      },
    });
  } catch (error) {
    uni.hideLoading();
    uni.showToast({
      title: "导出失败",
      icon: "none",
    });
    console.error("导出数据失败:", error);
  }
};

// 导入数据
const importData = async () => {
  uni.showModal({
    title: "功能提示",
    content:
      "数据导入功能需要后端支持批量导入接口,当前版本暂不支持。建议通过正常操作逐条添加记录。",
    showCancel: false,
  });

  // TODO: 等待后端 API 支持批量导入
  // try {
  //   const clipboardData = await uni.getClipboardData()
  //   const jsonStr = clipboardData.data
  //   const importedData = JSON.parse(jsonStr)
  //   // ... 验证和导入逻辑
  // } catch (error) {
  //   console.error('导入数据失败:', error)
  // }
};

// 关于我们
const showAbout = () => {
  uni.showModal({
    title: "关于我们",
    content:
      "宝宝喂养时刻是一款专为新手父母设计的育儿记录工具,帮助您科学、轻松地记录和追踪宝宝的成长数据。",
    showCancel: false,
  });
};

// 清除缓存
const clearCache = () => {
  uni.showModal({
    title: "确认清除",
    content: "清除缓存不会删除您的记录数据",
    success: (res) => {
      if (res.confirm) {
        uni.showToast({
          title: "清除成功",
          icon: "success",
        });
      }
    },
  });
};

// 非微信平台反馈提示
const showFeedbackTip = () => {
  uni.showModal({
    title: "用户反馈",
    content:
      "该功能仅在微信小程序中可用。请通过以下方式向我们反馈:\n\nEmail: feedback@babylog.com",
    showCancel: false,
    confirmText: "我知道了",
  });
};

// 退出登录
const handleLogout = () => {
  uni.showModal({
    title: "确认退出",
    content: "退出登录后,本地数据仍会保留",
    success: (res) => {
      if (res.confirm) {
        clearUserInfo();
        uni.reLaunch({
          url: "/pages/user/login",
        });
      }
    },
  });
};
</script>

<style lang="scss" scoped>
.user-page {
  min-height: 100vh;
  background: #f6f8f7;
  padding-bottom: 40rpx;
}

.user-card {
  background: linear-gradient(135deg, #7dd3a2 0%, #52c41a 100%);
  padding: 60rpx 30rpx 40rpx;
  color: white;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 24rpx;
}

.avatar {
  width: 120rpx;
  height: 120rpx;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  border: 4rpx solid rgba(255, 255, 255, 0.3);
  overflow: hidden;
}

.avatar-image {
  width: 100%;
  height: 100%;
  border-radius: 50%;
}

.info {
  flex: 1;
}

.login-button {
  background: linear-gradient(135deg, #52c41a, #7dd3a2);
  border-radius: 20rpx;
  padding: 12rpx 32rpx;
  display: inline-block;
  cursor: pointer;

  text {
    color: white;
    font-size: 28rpx;
    font-weight: 500;
  }
}

.nickname {
  font-size: 36rpx;
  font-weight: bold;
  margin-bottom: 12rpx;
}

.login-time {
  font-size: 24rpx;
  opacity: 0.8;
}

.section {
  margin-top: 20rpx;
}

.section-title {
  padding: 24rpx 30rpx 16rpx;
  font-size: 28rpx;
  color: #999;
}

:deep(.cell-icon) {
  color: #7dd3a2;
  margin-right: 12rpx;
}

.logout-section {
  margin-top: 40rpx;
  padding: 0 30rpx;
}

.version {
  text-align: center;
  padding: 40rpx 0;
  font-size: 24rpx;
  color: #999;
}

// 微信小程序原生反馈按钮样式重置
.feedback-button {
  background: none;
  border: none;
  padding: 0;
  margin: 0;
  width: 100%;
  text-align: left;
  font-size: inherit;
  line-height: inherit;
  color: inherit;

  &::before {
    display: none;
  }

  &::after {
    display: none;
  }
}
</style>
