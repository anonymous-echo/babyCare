<template>
  <view class="index-container">
    <wd-navbar
      fixed
      placeholder
      title="首页"
      left-arrow
      safe-area-inset-top
      custom-style="background: transparent; border: none;"
    >
      <template #left>
        <view
          v-if="currentBaby"
          class="baby-profile-header"
          @click="goToBabyList"
        >
          <view class="avatar-wrapper">
            <image
              v-if="currentBaby.avatarUrl"
              :src="resolveImageUrl(currentBaby.avatarUrl)"
              mode="aspectFill"
              class="avatar-img"
            />
            <image
              v-else
              src="/static/default.png"
              mode="aspectFill"
              class="avatar-img"
            />
          </view>
          <view class="baby-meta">
            <text class="name">{{ currentBaby.name }}</text>
            <text class="age">{{ babyAge }}</text>
          </view>
          <wd-icon name="caret-down-small" size="14" color="#64748B" />
        </view>
        <view v-else class="empty-baby-header" @click="handleAddBaby">
          <view class="plus-circle">
            <wd-icon name="plus" size="18" color="#FFFFFF" />
          </view>
          <text class="label">添加宝宝</text>
        </view>
      </template>
    </wd-navbar>

    <view class="page-body">
      <!-- 疫苗提醒 - 浮动设计 -->
      <view
        v-if="upcomingVaccines.length > 0"
        class="vaccine-alert-card premium-shadow"
        @click="goToVaccine"
      >
        <view class="alert-content">
          <view class="alert-icon">
            <wd-icon name="notification" size="20" color="#FA8C8C" />
          </view>
          <view class="alert-text">
            <text class="title">疫苗接种提醒</text>
            <text class="desc">{{ upcomingVaccines[0] }}</text>
          </view>
          <wd-icon name="arrow-right" size="16" color="#94A3B8" />
        </view>
      </view>

      <!-- 核心概览卡片 - 梦幻渐变 -->
      <view class="hero-overview-card premium-shadow">
        <view class="hero-bg"></view>
        <view class="hero-content">
          <view class="last-feeding">
            <text class="label">距离上次喂养</text>
            <text class="time">{{ lastFeedingTime }}</text>
          </view>
          <view class="quick-btn" @click="handleFeeding">
            <text>去记录</text>
            <wd-icon name="chevron-right" size="14" />
          </view>
        </view>
      </view>

      <!-- AI 智能建议 -->
      <view class="section-container">
        <view class="section-header">
          <text class="title">智护建议</text>
          <text class="more" @click="handleAIAnalysis">详细分析</text>
        </view>
        <DailyTipsCard
          :tips="todayTips"
          :max-display="3"
          @tip-click="handleTipClick"
        />
      </view>

      <!-- 快捷操作网格 -->
      <view class="action-dock premium-shadow">
        <view class="action-item" @click="handleFeeding">
          <view class="icon-box feeding"
            ><image src="/static/breastfeeding.svg"
          /></view>
          <text>喂养</text>
        </view>
        <view class="action-item" @click="handleSleep">
          <view class="icon-box sleep"
            ><image src="/static/moon_stars.svg"
          /></view>
          <text>睡眠</text>
        </view>
        <view class="action-item" @click="handleDiaper">
          <view class="icon-box diaper"
            ><image src="/static/blanket.svg"
          /></view>
          <text>尿布</text>
        </view>
        <view class="action-item" @click="handleGrowth">
          <view class="icon-box growth"
            ><image src="/static/monitoring.svg"
          /></view>
          <text>成长</text>
        </view>
      </view>

      <!-- 今日统计 - 模块化 -->
      <view class="section-container">
        <view class="section-header">
          <text class="title">今日动态</text>
        </view>
        <view class="stats-matrix">
          <view class="matrix-card feeding">
            <view class="card-inner">
              <text class="label">喂养总量</text>
              <view class="main-val">
                <text class="num">{{ todayStats.totalMilk }}</text>
                <text class="unit">ml</text>
              </view>
              <text class="sub"
                >{{ todayStats.breastfeedingCount }}次母乳 /
                {{ todayStats.bottleFeedingCount }}次奶瓶</text
              >
            </view>
          </view>

          <view class="matrix-card sleep">
            <view class="card-inner">
              <text class="label">睡眠时长</text>
              <view class="main-val">
                <text class="num">{{
                  formatSleepDuration(todayStats.sleepDurationMinutes)
                }}</text>
              </view>
              <text class="sub"
                >上次睡了
                {{ formatSleepDuration(todayStats.lastSleepMinutes) }}</text
              >
            </view>
          </view>

          <view class="matrix-card diaper">
            <view class="card-inner">
              <text class="label">换尿布</text>
              <view class="main-val">
                <text class="num">{{ todayStats.diaperCount }}</text>
                <text class="unit">次</text>
              </view>
              <text class="sub"
                >尿尿{{ todayStats.peeCount }} 粑粑{{
                  todayStats.poopCount
                }}</text
              >
            </view>
          </view>

          <view class="matrix-card growth">
            <view class="card-inner">
              <text class="label">当前体重</text>
              <view class="main-val">
                <text class="num">{{ todayStats.latestWeight ?? "-" }}</text>
                <text class="unit">g</text>
              </view>
              <text class="sub">本周增长 {{ weeklyStats.weightGain }}g</text>
            </view>
          </view>
        </view>
      </view>

      <!-- 周度回顾 -->
      <view class="section-container">
        <view class="section-header">
          <text class="title">周回顾</text>
        </view>
        <view class="weekly-card glass-card">
          <view class="week-item">
            <text class="val">{{ weeklyStats.feedingCount }}次</text>
            <text class="lab">总喂养</text>
          </view>
          <view class="week-sep"></view>
          <view class="week-item">
            <text class="val">{{
              formatSleepDuration(weeklyStats.sleepMinutes)
            }}</text>
            <text class="lab">总睡眠</text>
          </view>
          <view class="week-sep"></view>
          <view class="week-item">
            <text class="val">{{ weeklyStats.weightGain }}g</text>
            <text class="lab">体重增量</text>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { onShow, onPullDownRefresh } from "@dcloudio/uni-app";
import { isLoggedIn, fetchUserInfo } from "@/store/user";
import { currentBaby, fetchBabyList } from "@/store/baby";
import { aiStore } from "@/store/ai";
import { formatRelativeTime, calculateAge, formatDate } from "@/utils/date";
import { resolveImageUrl } from "@/utils/assets";
import DailyTipsCard from "@/components/DailyTipsCard.vue";

// 临时类型定义，避免导入问题
interface DailyTip {
  id: string;
  title: string;
  description: string;
  type: string;
  priority: "high" | "medium" | "low";
}

// 直接调用 API 层
import * as statisticsApi from "@/api/statistics";
import * as vaccineApi from "@/api/vaccine";

// 喂养订阅消息管理
import { requestAllFeedingSubscribeMessages } from "@/utils/feeding-subscribe";

// ============ 导航栏相关 ============

// 导航栏相关
const statusBarHeight = ref(0); // 状态栏高度（px）
const menuButtonWidth = ref(0); // 胶囊按钮宽度（px）
const menuButtonHeight = ref(0); // 胶囊按钮高度（px）
const menuButtonTop = ref(0); // 胶囊按钮顶部距离（px）

// 宝宝年龄
const babyAge = computed(() => {
  if (!currentBaby.value) return "";
  return calculateAge(currentBaby.value.birthDate);
});

// 跳转到宝宝列表
const goToBabyList = () => {
  uni.navigateTo({
    url: "/pages/baby/list/list",
  });
};

// 添加宝宝
const handleAddBaby = () => {
  uni.navigateTo({
    url: "/pages/baby/edit/edit",
  });
};

// ============ 响应式数据 ============

// 统计数据
const statistics = ref<statisticsApi.BabyStatisticsResponse | null>(null);

// 疫苗提醒数据（最多显示2个即将接种或逾期的疫苗）
const upcomingVaccines = ref<string[]>([]);

// 每日建议数据
const todayTips = computed(() => aiStore.todayTips.value || []);

// ============ 计算属性 ============

// 今日数据统计
const todayStats = computed(() => {
  const today = statistics.value?.today;
  if (!today) {
    return {
      breastfeedingCount: 0,
      bottleFeedingCount: 0,
      totalMilk: 0,
      sleepDuration: 0,
      sleepDurationMinutes: 0,
      lastSleepMinutes: 0,
      diaperCount: 0,
      peeCount: 0,
      poopCount: 0,
      latestWeight: null,
    };
  }

  return {
    // 喂养相关
    breastfeedingCount: today.feeding?.breastCount || 0, // 母乳次数
    bottleFeedingCount:
      (today.feeding?.totalCount || 0) - (today.feeding?.breastCount || 0), // 奶瓶次数
    totalMilk: today.feeding?.bottleMl || 0, // 奶瓶总毫升数
    // 睡眠相关
    sleepDuration: (today.sleep?.totalMinutes || 0) * 60, // 转换为秒，兼容 formatDuration
    sleepDurationMinutes: today.sleep?.totalMinutes || 0, // 保留分钟数用于显示
    lastSleepMinutes: today.sleep?.lastSleepMinutes || 0,
    // 尿布相关
    diaperCount: today.diaper?.totalCount ?? 0,
    wetCount: today.diaper?.wetCount ?? 0,
    dirtyCount: today.diaper?.dirtyCount ?? 0,
    // 成长相关
    latestWeight: today.growth?.latestWeight || null,
  };
});

// 距上次喂奶时间
const lastFeedingTime = computed(() => {
  if (!statistics.value?.today?.feeding?.lastFeedingTime) {
    return "-";
  }
  return formatRelativeTime(statistics.value.today.feeding.lastFeedingTime!);
});

// 格式化睡眠时间为 X时Y分
const formatSleepDuration = (minutes: number): string => {
  if (minutes <= 0) return "0分";

  const hours = Math.floor(minutes / 60);
  const remainingMinutes = minutes % 60;

  if (hours === 0) {
    return `${remainingMinutes}分`;
  } else if (remainingMinutes === 0) {
    return `${hours}时`;
  } else {
    return `${hours}时${remainingMinutes}分`;
  }
};

// 格式化睡眠趋势为 ±X时Y分
const formatSleepTrend = (minutes: number): string => {
  if (minutes === 0) return "0分";

  const absMinutes = Math.abs(minutes);
  const hours = Math.floor(absMinutes / 60);
  const remainingMinutes = absMinutes % 60;

  let result = minutes > 0 ? "+" : "-";

  if (hours === 0) {
    result += `${remainingMinutes}分`;
  } else if (remainingMinutes === 0) {
    result += `${hours}时`;
  } else {
    result += `${hours}时${remainingMinutes}分`;
  }

  return result;
};

// ============ 本周概览数据 ============

// 本周统计数据
const weeklyStats = computed(() => {
  if (!statistics.value) {
    return {
      feedingCount: 0,
      feedingTrend: 0,
      sleepMinutes: 0,
      sleepTrend: 0,
      weightGain: 0,
    };
  }

  const weekly = statistics.value.weekly;
  return {
    feedingCount: weekly.feeding.totalCount,
    feedingTrend: weekly.feeding.trend,
    sleepMinutes: weekly.sleep.totalMinutes,
    sleepTrend: weekly.sleep.trend,
    weightGain: weekly.growth.weightGain,
  };
});

// 最近体重
const latestWeight = computed(() => {
  // 示例：从今日数据中获取，实际应从成长记录中获取
  return "7.5";
});

// 页面加载 (仅在首次挂载时执行)
onMounted(() => {
  console.log("[Index] onMounted");
  // 初始化导航栏
  initializeNavbar();
});

// 初始化导航栏
const initializeNavbar = () => {
  // 获取系统信息
  const systemInfo = uni.getSystemInfoSync();
  statusBarHeight.value = systemInfo.statusBarHeight || 0;

  // 获取胶囊按钮信息（仅微信小程序）
  // #ifdef MP-WEIXIN
  try {
    const menuButton = uni.getMenuButtonBoundingClientRect();
    if (menuButton) {
      // 胶囊按钮的宽度和高度（保持 px，与导航栏样式中使用 rpx 统一处理）
      menuButtonWidth.value = menuButton.width; // px
      menuButtonHeight.value = menuButton.height; // px
      menuButtonTop.value = menuButton.top; // px（状态栏下的距离）

      console.log("[Index] 胶囊对齐:", {
        statusBarHeight: statusBarHeight.value,
        menuButtonTop: menuButtonTop.value,
        menuButtonWidth: menuButton.width,
        menuButtonHeight: menuButton.height,
        menuButtonBottom: menuButton.top + menuButton.height,
      });
    }
  } catch (e) {
    console.warn("[Index] 获取胶囊信息失败，使用默认高度", e);
    // 使用默认值
    menuButtonWidth.value = 88; // 默认宽度
    menuButtonHeight.value = 32; // 默认高度
  }
  // #endif
};

// 页面显示 (每次页面显示时执行,包括 switchTab)
onShow(async () => {
  console.log("[Index] onShow - 开始检查登录和宝宝信息");

  // 检查登录和宝宝信息
  await checkLoginAndBaby();
});

// 计算页面内容的 padding-top
// 已改为计算属性 pageContentPaddingTop，无需手动计算

// 检查登录和宝宝信息
const checkLoginAndBaby = async () => {
  console.log("[Index] checkLoginAndBaby - 登录状态:", isLoggedIn.value);

  // 1. 检查登录状态
  if (!isLoggedIn.value) {
    console.log("[Index] 未登录，显示游客模式");
    // ✅ 未登录时不强制跳转，显示游客模式提示
    // 游客模式：用户可以浏览首页，但无法查看真实数据
    return;
  }

  try {
    // 2. 获取用户信息
    await fetchUserInfo();

    // 3. 获取宝宝列表
    const babies = await fetchBabyList();

    console.log("[Index] 宝宝列表:", babies);
    console.log("[Index] 当前宝宝:", currentBaby.value);

    // 4. 有宝宝,加载今日数据
    if (currentBaby.value) {
      await loadTodayData();
    }
  } catch (error) {
    console.error("[Index] 获取用户/宝宝信息失败:", error);
    uni.showToast({
      title: "加载数据失败",
      icon: "none",
    });
  }
};

type LoadTodayDataOptions = {
  preserveData?: boolean;
  pullDown?: boolean;
};

// 加载今日数据
const loadTodayData = async (options: LoadTodayDataOptions = {}) => {
  if (!currentBaby.value) {
    if (options.pullDown) {
      uni.hideNavigationBarLoading();
      uni.stopPullDownRefresh();
    }
    return;
  }

  const babyId = currentBaby.value.babyId;

  if (options.pullDown) {
    uni.showNavigationBarLoading();
  } else {
    uni.showLoading({ title: "加载中", mask: false });
  }

  try {
    if (!options.preserveData) {
      // 清空旧数据
      statistics.value = null;
      upcomingVaccines.value = [];
    }

    // 并行加载统计数据和疫苗提醒（重要数据）
    const [statisticsResponse, vaccineRemindersResponse] = await Promise.all([
      statisticsApi.apiFetchBabyStatistics(babyId),
      vaccineApi
        .apiFetchVaccineReminders({
          babyId,
        })
        .catch((error) => {
          console.error("加载疫苗提醒失败:", error);
          return { reminders: [], total: 0 };
        }),
    ]);

    // 处理统计数据
    statistics.value = statisticsResponse.data;

    // 处理疫苗提醒：筛选出 upcoming、due、overdue 的记录，最多显示2个
    const reminders = vaccineRemindersResponse.reminders || [];
    const filtered = reminders.filter(
      (r: vaccineApi.VaccineReminderResponse) =>
        r.status === "upcoming" || r.status === "due" || r.status === "overdue",
    );
    upcomingVaccines.value = filtered.map(
      (r: vaccineApi.VaccineReminderResponse) =>
        `${r.vaccineName} ${r.doseNumber ? `（第${r.doseNumber}针）` : ""} ${
          vaccineApi.VaccineReminderStatusMap[r.status]
        }，应于 ${formatDate(r.scheduledDate, "YYYY-MM-DD")}接种`,
    );

    // 异步加载每日建议（非阻塞）
    loadDailyTipsAsync(babyId);

    if (options.pullDown) {
      uni.showToast({
        title: "刷新成功",
        icon: "success",
        duration: 1200,
      });
    }
  } catch (error) {
    console.error("[Index] 加载数据失败:", error);
    if (options.pullDown) {
      uni.showToast({
        title: "刷新失败",
        icon: "none",
        duration: 1500,
      });
    }
    // 不显示错误提示，静默失败
  } finally {
    if (options.pullDown) {
      uni.hideNavigationBarLoading();
      uni.stopPullDownRefresh();
    } else {
      uni.hideLoading();
    }
  }
};

onPullDownRefresh(async () => {
  if (!isLoggedIn.value || !currentBaby.value) {
    uni.stopPullDownRefresh();
    uni.hideNavigationBarLoading();
    return;
  }
  await loadTodayData({ preserveData: true, pullDown: true });
});

// 跳转到登录
const goToLogin = () => {
  uni.navigateTo({
    url: "/pages/user/login",
  });
};

// 喂养记录（需要检查登录状态）
const handleFeeding = async () => {
  if (!isLoggedIn.value) {
    uni.showModal({
      title: "提示",
      content: "该功能需要登录，是否前往登录？",
      success: (res) => {
        if (res.confirm) {
          goToLogin();
        }
      },
    });
    return;
  }

  if (!currentBaby.value) {
    uni.showToast({
      title: "请先添加宝宝",
      icon: "none",
    });
    return;
  }
  // ✨ 在跳转前申请喂养订阅消息权限
  try {
    console.log("[Index] 检查是否需要申请喂养订阅消息");
    // 申请喂养订阅消息
    await requestAllFeedingSubscribeMessages();
  } catch (error: any) {
    console.error("[Index] 申请订阅消息失败:", error);
    // 静默失败，不影响主功能
  }

  // 申请完成后跳转到喂养记录页面
  uni.navigateTo({
    url: "/pages/record/feeding/feeding",
  });
};

// 换尿布记录
const handleDiaper = () => {
  if (!isLoggedIn.value) {
    uni.showModal({
      title: "提示",
      content: "该功能需要登录，是否前往登录？",
      success: (res) => {
        if (res.confirm) {
          goToLogin();
        }
      },
    });
    return;
  }

  if (!currentBaby.value) {
    uni.showToast({
      title: "请先添加宝宝",
      icon: "none",
    });
    return;
  }
  uni.navigateTo({
    url: "/pages/record/diaper/diaper",
  });
};

// 睡眠记录
const handleSleep = () => {
  if (!isLoggedIn.value) {
    uni.showModal({
      title: "提示",
      content: "该功能需要登录，是否前往登录？",
      success: (res) => {
        if (res.confirm) {
          goToLogin();
        }
      },
    });
    return;
  }

  if (!currentBaby.value) {
    uni.showToast({
      title: "请先添加宝宝",
      icon: "none",
    });
    return;
  }
  uni.navigateTo({
    url: "/pages/record/sleep/sleep",
  });
};

// 成长记录
const handleGrowth = () => {
  if (!isLoggedIn.value) {
    uni.showModal({
      title: "提示",
      content: "该功能需要登录，是否前往登录？",
      success: (res) => {
        if (res.confirm) {
          goToLogin();
        }
      },
    });
    return;
  }

  if (!currentBaby.value) {
    uni.showToast({
      title: "请先添加宝宝",
      icon: "none",
    });
    return;
  }
  uni.navigateTo({
    url: "/pages/record/growth/growth",
  });
};

// 异步加载每日建议（非阻塞式）
const loadDailyTipsAsync = async (babyId: string) => {
  try {
    // 完全静默加载，不输出任何日志
    await aiStore.getDailyTips(parseInt(babyId));
  } catch (error) {
    // 完全静默失败，不输出任何日志，不影响其他功能
    // 即使失败也不影响用户体验
  }
};

// 跳转到疫苗页面
const goToVaccine = () => {
  if (!isLoggedIn.value) {
    uni.showModal({
      title: "提示",
      content: "该功能需要登录，是否前往登录？",
      success: (res) => {
        if (res.confirm) {
          goToLogin();
        }
      },
    });
    return;
  }

  if (!currentBaby.value) {
    uni.showToast({
      title: "请先添加宝宝",
      icon: "none",
    });
    return;
  }
  uni.navigateTo({
    url: "/pages/vaccine/vaccine",
  });
};

// 处理每日建议点击
const handleTipClick = (tip: DailyTip) => {
  console.log("点击每日建议:", tip);

  // 显示建议详情弹窗
  uni.showModal({
    title: tip.title,
    content: tip.description,
    showCancel: false,
    confirmText: "知道了",
  });
};

// AI分析
const handleAIAnalysis = () => {
  if (!isLoggedIn?.value) {
    uni.showModal({
      title: "提示",
      content: "该功能需要登录，是否前往登录？",
      success: (res) => {
        if (res.confirm) {
          goToLogin();
        }
      },
    });
    return;
  }

  if (!currentBaby?.value) {
    uni.showToast({
      title: "请先添加宝宝",
      icon: "none",
    });
    return;
  }

  uni.navigateTo({
    url: "/pages/statistics/ai-analysis",
  });
};
</script>

<style lang="scss" scoped>
@import "@/styles/colors.scss";

.index-container {
  min-height: 100vh;
  background: $color-bg-secondary;
}

.page-body {
  padding: $page-gap;
  padding-bottom: calc($page-gap + env(safe-area-inset-bottom) + 100rpx);
}

// ===== Navbar Custom =====
.baby-profile-header {
  display: flex;
  align-items: center;
  gap: 16rpx;
  background: rgba(255, 255, 255, 0.6);
  padding: 8rpx 20rpx 8rpx 8rpx;
  border-radius: 100rpx;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(123, 211, 162, 0.2);

  .avatar-wrapper {
    width: 64rpx;
    height: 64rpx;
    border-radius: 50%;
    overflow: hidden;
    border: 4rpx solid #fff;
    box-shadow: 0 4rpx 12rpx rgba(123, 211, 162, 0.2);

    .avatar-img {
      width: 100%;
      height: 100%;
    }
  }

  .baby-meta {
    display: flex;
    flex-direction: column;

    .name {
      font-size: 26rpx;
      font-weight: 600;
      color: $color-text-primary;
    }
    .age {
      font-size: 20rpx;
      color: $color-text-tertiary;
    }
  }
}

.empty-baby-header {
  display: flex;
  align-items: center;
  gap: 12rpx;
  background: $color-primary;
  padding: 8rpx 24rpx 8rpx 10rpx;
  border-radius: 100rpx;

  .plus-circle {
    width: 48rpx;
    height: 48rpx;
    background: rgba(255, 255, 255, 0.2);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .label {
    font-size: 24rpx;
    color: #fff;
    font-weight: 500;
  }
}

// ===== Vaccine Alert =====
.vaccine-alert-card {
  background: #fff;
  border-radius: $radius-md;
  margin-bottom: $spacing-xl;
  overflow: hidden;
  border-left: 10rpx solid $color-danger;

  .alert-content {
    display: flex;
    align-items: center;
    padding: $spacing-lg;
    gap: 16rpx;
  }

  .alert-icon {
    width: 80rpx;
    height: 80rpx;
    background: $color-danger-light;
    border-radius: $radius-sm;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .alert-text {
    flex: 1;
    display: flex;
    flex-direction: column;

    .title {
      font-size: 28rpx;
      font-weight: 600;
      color: $color-text-primary;
    }
    .desc {
      font-size: 22rpx;
      color: $color-text-secondary;
      margin-top: 4rpx;
    }
  }
}

// ===== Hero Card =====
.hero-overview-card {
  position: relative;
  height: 240rpx;
  background: $gradient-dream;
  border-radius: $radius-lg;
  margin-bottom: $spacing-xl;
  overflow: hidden;
  display: flex;
  align-items: center;
  padding: 0 40rpx;

  &::before {
    content: "";
    position: absolute;
    top: -20%;
    right: -10%;
    width: 300rpx;
    height: 300rpx;
    background: rgba(255, 255, 255, 0.15);
    border-radius: 50%;
    filter: blur(40px);
  }

  .hero-content {
    width: 100%;
    display: flex;
    justify-content: space-between;
    align-items: center;
    position: relative;
    z-index: 1;
  }

  .last-feeding {
    display: flex;
    flex-direction: column;

    .label {
      font-size: 24rpx;
      color: rgba(30, 41, 59, 0.7);
      margin-bottom: 8rpx;
    }
    .time {
      font-size: 48rpx;
      font-weight: 700;
      color: $color-text-primary;
      letter-spacing: -1rpx;
    }
  }

  .quick-btn {
    background: #fff;
    padding: 16rpx 32rpx;
    border-radius: 100rpx;
    display: flex;
    align-items: center;
    gap: 8rpx;
    box-shadow: 0 8px 20px rgba(0, 0, 0, 0.05);

    text {
      font-size: 24rpx;
      font-weight: 600;
      color: $color-primary-dark;
    }
  }
}

// ===== Section Common =====
.section-container {
  margin-bottom: $spacing-xl;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: $spacing-md;

  .title {
    font-size: 32rpx;
    font-weight: 700;
    color: $color-text-primary;
    position: relative;

    &::after {
      content: "";
      position: absolute;
      bottom: -4rpx;
      left: 0;
      width: 40rpx;
      height: 6rpx;
      background: $color-primary;
      border-radius: 10rpx;
    }
  }

  .more {
    font-size: 24rpx;
    color: $color-text-tertiary;
  }
}

// ===== Action Dock =====
.action-dock {
  background: #fff;
  border-radius: $radius-lg;
  display: flex;
  justify-content: space-around;
  padding: 32rpx 0;
  margin-bottom: $spacing-xl;

  .action-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 12rpx;

    .icon-box {
      width: 100rpx;
      height: 100rpx;
      border-radius: 32rpx;
      display: flex;
      align-items: center;
      justify-content: center;
      image {
        width: 52rpx;
        height: 52rpx;
      }

      &.feeding {
        background: #e9f7f0;
      }
      &.sleep {
        background: #ebf4ff;
      }
      &.diaper {
        background: #fff4e6;
      }
      &.growth {
        background: #f3eeff;
      }
    }

    text {
      font-size: 24rpx;
      font-weight: 500;
      color: $color-text-secondary;
    }
  }
}

// ===== Stats Matrix =====
.stats-matrix {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: $spacing-lg;

  .matrix-card {
    background: #fff;
    border-radius: $radius-md;
    padding: $spacing-lg;
    position: relative;
    overflow: hidden;

    .card-inner {
      position: relative;
      z-index: 1;
      display: flex;
      flex-direction: column;
    }

    .label {
      font-size: 24rpx;
      color: $color-text-tertiary;
      margin-bottom: 12rpx;
    }

    .main-val {
      display: flex;
      align-items: baseline;
      gap: 4rpx;
      margin-bottom: 12rpx;

      .num {
        font-size: 44rpx;
        font-weight: 700;
        color: $color-text-primary;
      }
      .unit {
        font-size: 20rpx;
        color: $color-text-tertiary;
      }
    }

    .sub {
      font-size: 20rpx;
      color: $color-text-tertiary;
    }

    &::after {
      content: "";
      position: absolute;
      top: -20rpx;
      right: -20rpx;
      width: 100rpx;
      height: 100rpx;
      background: rgba(123, 211, 162, 0.05);
      border-radius: 50%;
    }
  }
}

// ===== Weekly Card =====
.weekly-card {
  display: flex;
  align-items: center;
  justify-content: space-around;
  padding: 40rpx 0;
  border-radius: $radius-lg;

  .week-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 8rpx;

    .val {
      font-size: 32rpx;
      font-weight: 700;
      color: $color-primary-dark;
    }
    .lab {
      font-size: 22rpx;
      color: $color-text-tertiary;
    }
  }

  .week-sep {
    width: 1rpx;
    height: 40rpx;
    background: $color-divider;
  }
}
</style>
