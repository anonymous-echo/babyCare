<template>
  <view>
    <wd-navbar fixed placeholder title="首页" left-arrow safeAreaInsetTop>
      <template #left>
        <view
          v-if="currentBaby"
          class="baby-info"
          @click="goToBabyList"
          :style="{
            maxWidth: '360rpx',
            height: menuButtonHeight * 2 + 'rpx',
          }"
        >
          <view class="baby-content">
            <view class="baby-avatar">
              <image
                v-if="currentBaby.avatarUrl"
                :src="currentBaby.avatarUrl"
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
            <view class="baby-text">
              <text class="baby-name">{{ currentBaby.name }}</text>
              <text class="baby-age">{{ babyAge }}</text>
            </view>
            <wd-icon name="right" size="12" color="#7f8c8d" class="arrow-icon" />
          </view>
        </view>
        <!-- 没有宝宝时显示录入按钮 -->
        <view
          v-else
          class="add-baby-button"
          @click="handleAddBaby"
          :style="{
            maxWidth: '360rpx',
            height: menuButtonHeight * 2 + 'rpx',
          }"
        >
          <view class="button-content">
            <wd-icon name="plus" size="18" color="#32dc6e" class="plus-icon" />
            <text class="button-text">添加宝宝</text>
          </view>
        </view>
      </template>
    </wd-navbar>
    <view class="index-page">
      <!-- 页面内容 -->
      <view class="page-content">
        <!-- 疫苗提醒横幅 -->
        <view
          v-if="upcomingVaccines.length > 0"
          class="vaccine-banner"
          @click="goToVaccine"
        >
          <wd-notice-bar
            prefix="warn-bold"
            :text="upcomingVaccines"
            :delay="3"
            custom-class="space"
          >
            <template #suffix>
              <wd-icon name="arrow-right" size="12" color="#2c3e50" />
            </template>
          </wd-notice-bar>
        </view>

        <!-- AI每日建议 -->
        <DailyTipsCard 
          :tips="todayTips" 
          :max-display="10"
          @tip-click="handleTipClick"
        />

        <!-- 距离上次喂养时间提示 - 显眼卡片 -->
        <view class="last-feeding-card">
          <view class="feeding-left">
            <view class="feeding-icon">
              <image
                src="/static/breastfeeding.svg"
                mode="aspectFill"
                style="width: 40rpx; height: 40rpx"
              />
            </view>
            <view class="feeding-info">
              <text class="feeding-label"
                >距离上次喂养 {{ lastFeedingTime }}</text
              >
            </view>
          </view>
          <view class="feeding-action">
            <view class="action-btn" @click="handleFeeding">
              <wd-icon name="arrow-right" size="12" color="#2c3e50" />
            </view>
          </view>
        </view>
        <!-- 快捷操作 - 4 列网格 -->
        <view class="quick-actions">
          <view class="action-title">快捷操作</view>
          <view class="action-grid">
            <view class="action-card action-feeding" @click="handleFeeding">
              <view class="action-icon-wrapper">
                <image
                  src="/static/breastfeeding.svg"
                  mode="aspectFill"
                  class="action-icon"
                />
              </view>
              <text class="action-label">记录喂养</text>
            </view>
            <view class="action-card action-sleep" @click="handleSleep">
              <view class="action-icon-wrapper">
                <image
                  src="/static/moon_stars.svg"
                  mode="aspectFill"
                  class="action-icon"
                />
              </view>
              <text class="action-label">记录睡眠</text>
            </view>
            <view class="action-card action-diaper" @click="handleDiaper">
              <view class="action-icon-wrapper">
                <image
                  src="/static/blanket.svg"
                  mode="aspectFill"
                  class="action-icon"
                />
              </view>
              <text class="action-label">记录尿布</text>
            </view>
            <view class="action-card action-growth" @click="handleGrowth">
              <view class="action-icon-wrapper">
                <image
                  src="/static/monitoring.svg"
                  mode="aspectFill"
                  class="action-icon"
                />
              </view>
              <text class="action-label">记录成长</text>
            </view>
            <!-- <view class="action-card action-ai" @click="handleAIAnalysis">
              <view class="action-icon-wrapper">
                <image
                  src="/static/smart_toy.svg"
                  mode="aspectFill"
                  class="action-icon"
                />
              </view>
              <text class="action-label">AI分析</text>
            </view> -->
          </view>
        </view>
        <!-- 今日数据概览 - 2x2 网格 -->
        <view class="today-stats">
          <view class="card-header">
            <text class="card-title">今日概览</text>
            <text class="card-subtitle">实时数据</text>
          </view>
          <view class="stats-grid">
            <!-- 喂养统计 -->
            <view class="stat-card">
              <view class="stat-header">
                <image
                  src="/static/breastfeeding.svg"
                  mode="aspectFill"
                  class="stat-icon"
                />
                <text class="stat-card-title">喂养统计</text>
              </view>
              <view class="stat-main">
                <text class="stat-value"
                  >母乳 {{ todayStats.breastfeedingCount }} 次</text
                >
                <text class="stat-sub">奶瓶 {{ todayStats.bottleFeedingCount }} 次 / {{ todayStats.totalMilk }}ml</text>
              </view>
            </view>

            <!-- 睡眠时长 -->
            <view class="stat-card">
              <view class="stat-header">
                <image
                  class="stat-icon"
                  src="/static/moon_stars.svg"
                  mode="aspectFill"
                />
                <text class="stat-card-title">总睡眠时长</text>
              </view>
              <view class="stat-main">
                <text class="stat-value">{{ formatSleepDuration(todayStats.sleepDurationMinutes) }}</text>
                <text class="stat-sub"
                  >上次睡眠: {{ formatSleepDuration(todayStats.lastSleepMinutes) }}</text
                >
              </view>
            </view>

            <!-- 换尿布次数 -->
            <view class="stat-card">
              <view class="stat-header">
                <image
                  src="/static/baby_changing_station.svg"
                  mode="aspectFill"
                  class="stat-icon"
                />
                <text class="stat-card-title">换尿布次数</text>
              </view>
              <view class="stat-main">
                <text class="stat-value">{{ todayStats.diaperCount }} 次</text>
                <text class="stat-sub"
                  >小便: {{ todayStats.peeCount }}, 大便:
                  {{ todayStats.poopCount }}</text
                >
              </view>
            </view>

            <!-- 体重 -->
            <view class="stat-card">
              <view class="stat-header">
                <image
                  src="/static/weight.svg"
                  mode="aspectFill"
                  class="stat-icon"
                />
                <text class="stat-card-title">体重</text>
              </view>
              <view class="stat-main">
                <text class="stat-value"
                  >{{ todayStats.latestWeight ?? "-" }} 克</text
                >
                <text class="stat-sub"
                  >{{ weeklyStats.weightGain >= 0 ? "↑" : "↓" }}
                  {{ Math.abs(weeklyStats.weightGain) }}克</text
                >
              </view>
            </view>
          </view>
        </view>

        <!-- 本周概览卡片 -->
        <view class="weekly-overview">
          <view class="card-header">
            <text class="card-title">本周概览</text>
            <text class="card-subtitle">过去 7 天数据</text>
          </view>
          <view class="overview-grid">
            <view class="overview-item">
              <view class="overview-label">总喂养次数</view>
              <text class="overview-value">{{ weeklyStats.feedingCount }} 次</text>
              <text
                class="overview-trend"
                :class="weeklyStats.feedingTrend >= 0 ? 'up' : 'down'"
              >
                {{ weeklyStats.feedingTrend >= 0 ? "↑" : "↓" }}
                {{ Math.abs(weeklyStats.feedingTrend) }}
              </text>
            </view>
            <view class="overview-item">
              <view class="overview-label">总睡眠时长</view>
              <text class="overview-value">{{ formatSleepDuration(weeklyStats.sleepMinutes) }}</text>
              <text
                class="overview-trend"
                :class="weeklyStats.sleepTrend >= 0 ? 'up' : 'down'"
              >
                {{ weeklyStats.sleepTrend >= 0 ? "↑" : "↓" }}
                {{ formatSleepTrend(weeklyStats.sleepTrend) }}
              </text>
            </view>
            <view class="overview-item">
              <view class="overview-label">体重增长</view>
              <text class="overview-value"
                >{{ weeklyStats.weightGain }} 克</text
              >
              <text
                class="overview-trend"
                :class="weeklyStats.weightGain >= 0 ? 'up' : 'down'"
              >
                {{ weeklyStats.weightGain >= 0 ? "↑" : "↓" }}
                {{ Math.abs(weeklyStats.weightGain) }}
              </text>
            </view>
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
import {
  formatRelativeTime,
  calculateAge,
  formatDate,
} from "@/utils/date";
import DailyTipsCard from "@/components/DailyTipsCard.vue";

// 临时类型定义，避免导入问题
interface DailyTip {
  id: string
  title: string
  description: string
  type: string
  priority: 'high' | 'medium' | 'low'
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
  if (!statistics.value) {
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

  const today = statistics.value.today;
  return {
    // 喂养相关
    breastfeedingCount: today.feeding.breastCount, // 母乳次数
    bottleFeedingCount: today.feeding.totalCount - today.feeding.breastCount, // 奶瓶次数
    totalMilk: today.feeding.bottleMl, // 奶瓶总毫升数
    // 睡眠相关
    sleepDuration: today.sleep.totalMinutes * 60, // 转换为秒，兼容 formatDuration
    sleepDurationMinutes: today.sleep.totalMinutes, // 保留分钟数用于显示
    lastSleepMinutes: today.sleep.lastSleepMinutes,
    // 尿布相关
    diaperCount: today.diaper.totalCount,
    peeCount: today.diaper.peeCount,
    poopCount: today.diaper.poopCount,
    // 成长相关
    latestWeight: today.growth.latestWeight,
  };
});

// 距上次喂奶时间
const lastFeedingTime = computed(() => {
  if (!statistics.value?.today?.feeding?.lastFeedingTime) {
    return "-";
  }
  return formatRelativeTime(statistics.value.today.feeding.lastFeedingTime);
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
      vaccineApi.apiFetchVaccineReminders({
        babyId,
      }).catch(error => {
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
        r.status === "upcoming" || r.status === "due" || r.status === "overdue"
    );
    upcomingVaccines.value = filtered.map((r: vaccineApi.VaccineReminderResponse) =>
      `${r.vaccineName} ${r.doseNumber ? `（第${r.doseNumber}针）` : ""} ${
        vaccineApi.VaccineReminderStatusMap[r.status]
      }，应于 ${formatDate(r.scheduledDate, "YYYY-MM-DD")}接种`
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
@import '@/styles/colors.scss';

// ===== 设计系统变量 =====
$spacing: 20rpx; // 统一间距

// ===== 导航栏样式 =====
.navbar-wrapper {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  background: $color-bg-primary;
  z-index: 9999;
}

.navbar-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20rpx; // 左右边距
  // 高度由内联样式动态设置
}

// 左侧宝宝信息 - 对齐胶囊位置
.baby-info {
  display: flex;
  align-items: center;
  justify-content: flex-start;
  flex-shrink: 0;
  min-width: 200rpx;
  // 宽高由内联样式动态设置
}

.baby-content {
  display: flex;
  align-items: center;
  gap: 8rpx;
  padding: 6rpx 16rpx 6rpx 6rpx;
  background: $color-bg-secondary;
  border-radius: 40rpx;
  height: 100%;
  max-width: 100%;
}

.baby-avatar {
  width: 52rpx;
  height: 52rpx;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
}

.avatar-img {
  width: 100%;
  height: 100%;
}

.baby-text {
  display: flex;
  flex-direction: column;
  gap: 2rpx;
  flex: 1;
  min-width: 0;
  max-width: 200rpx;
}

.baby-name {
  font-size: 26rpx;
  font-weight: $font-weight-medium;
  color: $color-text-primary;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  line-height: 1.3;
}

.baby-age {
  font-size: 22rpx;
  color: $color-text-secondary;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  line-height: 1.2;
}

.arrow-icon {
  flex-shrink: 0;
  margin-left: 2rpx;
}

.add-baby-hint {
  padding: 16rpx 32rpx;
  background: $color-bg-secondary;
  border-radius: 40rpx;
  font-size: 24rpx;
  color: $color-text-secondary;
}

// 添加宝宝按钮 - 对齐胶囊位置
.add-baby-button {
  display: flex;
  align-items: center;
  justify-content: flex-start;
  flex-shrink: 0;
  min-width: 200rpx;
  // 宽高由内联样式动态设置
}

.add-baby-button .button-content {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8rpx;
  padding: 6rpx 16rpx;
  background: $color-bg-secondary;
  border-radius: 40rpx;
  height: 100%;
  max-width: 100%;
  transition: all $transition-slow;

  &:active {
    background: $color-primary-lighter;
    transform: scale(0.95);
  }
}

.plus-icon {
  flex-shrink: 0;
}

.button-text {
  font-size: 26rpx;
  font-weight: $font-weight-medium;
  color: $color-text-primary;
  white-space: nowrap;
}

// 中间标题 - 居中显示
.navbar-title {
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
  font-size: 34rpx; // 标准导航栏标题大小 (17px = 34rpx)
  font-weight: 600;
  color: $color-text-primary;
  pointer-events: none;
}

// 右侧占位符（与胶囊等宽）
.navbar-right {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  flex-shrink: 0;
  // 宽高由内联样式动态设置
}

.index-page {
  // padding-top 由内联样式动态设置
  min-height: 100vh;
  background: $color-bg-secondary;
  padding-top: 12rpx;
}

// 页面内容区域 - 修复布局
.page-content {
  // 顶部由内联样式动态设置 (导航栏总高度 + 间距)
  padding-left: 16rpx;
  padding-right: 16rpx;
  padding-bottom: $spacing;
  padding-top: 12rpx;

  // 为 tabBar 预留空间（env(safe-area-inset-bottom) 处理全面屏底部安全区）
  margin-bottom: calc(100rpx + env(safe-area-inset-bottom));
}

// ============ 距离上次喂养卡片 ============
.last-feeding-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: $color-bg-primary;
  border: 1rpx solid $color-border-primary;
  border-left: 4rpx solid $color-primary;
  border-radius: $radius-lg;
  padding: $spacing-lg $spacing-md;
  margin-bottom: $spacing-2xl;
  box-shadow: $shadow-primary-sm;
}

.feeding-left {
  display: flex;
  align-items: center;
  gap: 12rpx;
  flex: 1;
}

.feeding-icon {
  width: 48rpx;
  height: 48rpx;
  border-radius: $radius-full;
  background: $color-primary-lighter;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.feeding-info {
  display: flex;
  flex-direction: column;
  gap: 4rpx;
  flex: 1;
}

.feeding-label {
  font-size: 24rpx;
  color: $color-text-secondary;
  font-weight: $font-weight-medium;
}

.feeding-time {
  font-size: 36rpx;
  font-weight: $font-weight-bold;
  color: $color-text-primary;
  line-height: 1.2;
}

.feeding-action {
  display: flex;
  align-items: center;
  flex-shrink: 0;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 6rpx;
  padding: 10rpx 14rpx;
  background: $color-bg-secondary;
  border: 1rpx solid $color-border-primary;
  border-radius: $radius-xl;
  transition: all $transition-base;

  &:active {
    background: $color-primary-lighter;
    transform: scale(0.98);
  }
}

// 疫苗提醒横幅
.vaccine-banner {
  padding: 20rpx 0;
}

.banner-left {
  display: flex;
  align-items: center;
  gap: 12rpx;
  flex: 1;
  min-width: 0;
}

.banner-icon {
  font-size: 32rpx;
  font-weight: 700;
  color: $color-text-secondary;
  flex-shrink: 0;
}

.banner-text {
  display: flex;
  flex-direction: column;
  gap: 4rpx;
  flex: 1;
  min-width: 0;

  text {
    font-size: 32rpx;
    color: $color-text-primary;
    line-height: 1.3;
    word-break: break-word;
  }
}

// 今日数据卡片
.today-stats {
  margin-bottom: $spacing-2xl;
  background: $color-bg-primary;
  border: 1rpx solid $color-border-primary;
  border-radius: $radius-lg;
  padding: $spacing-lg $spacing-md;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16rpx;
}

// 新的数据卡片样式
.stat-card {
  background: $color-bg-secondary;
  border-radius: $radius-md;
  padding: $spacing-lg $spacing-md;
  display: flex;
  flex-direction: column;
  gap: $spacing-md;
}

.stat-icon {
  width: 40rpx;
  height: 40rpx;
}

.stat-header {
  display: flex;
  align-items: center;
  gap: 8rpx;
  text-align: center;
}

.stat-card-title {
  font-size: 28rpx;
  color: $color-text-primary;
  font-weight: 350;
  text-align: left;
  padding-bottom: 0;
  padding-left: 4rpx;
}

.stat-main {
  display: flex;
  flex-direction: column;
  gap: 6rpx;
}

.stat-value {
  font-size: 34rpx;
  font-weight: 450;
  color: $color-primary-light;
  line-height: 1.3;
  text-align: center;
}

.stat-sub {
  padding-top: 0;
  font-size: 24rpx;
  color: $color-text-secondary;
  text-align: center;
}

// ============ 本周概览卡片 ============
.weekly-overview {
  background: $color-bg-primary;
  border: 1rpx solid $color-border-primary;
  border-radius: $radius-lg;
  padding: $spacing-lg $spacing-md;
  margin-bottom: $spacing-2xl;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16rpx;
}

.card-title {
  font-size: 32rpx;
  font-weight: $font-weight-semibold;
  color: $color-text-primary;
}

.card-subtitle {
  font-size: 22rpx;
  color: $color-text-secondary;
}

.overview-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12rpx;
}

.overview-item {
  background: $color-bg-secondary;
  border-radius: $radius-md;
  padding: $spacing-md 8rpx;
  text-align: center;
  display: flex;
  flex-direction: column;
  gap: $spacing-md;
}

.overview-label {
  font-size: 22rpx;
  color: $color-text-secondary;
  font-weight: $font-weight-medium;
}

.overview-value {
  font-size: 36rpx;
  font-weight: 450;
  color: $color-primary-light;
  line-height: 1.2;
}

.overview-trend {
  font-size: 22rpx;
  font-weight: $font-weight-semibold;
  color: $color-text-secondary;
}

// ============ 宝宝成长进度卡片 ============
.growth-milestone {
  background: $color-bg-primary;
  border: 1rpx solid $color-border-primary;
  border-radius: $radius-lg;
  padding: $spacing-lg $spacing-md;
  margin-bottom: $spacing-2xl;
}

.milestone-content {
  display: flex;
  flex-direction: column;
  gap: 16rpx;
}

.milestone-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12rpx 0;
  border-bottom: 1rpx solid $color-divider;

  &:last-of-type {
    border-bottom: none;
  }
}

.milestone-label {
  display: flex;
  align-items: center;
  gap: 8rpx;
}

.label-icon {
  font-size: 28rpx;
}

.label-text {
  font-size: 26rpx;
  color: $color-text-secondary;
  font-weight: $font-weight-medium;
}

.milestone-value {
  font-size: 28rpx;
  font-weight: $font-weight-semibold;
  color: $color-text-primary;
}

.progress-bar {
  width: 100%;
  height: 8rpx;
  background: $color-bg-disabled;
  border-radius: $radius-xs;
  overflow: hidden;
  margin-top: 8rpx;
}

.progress-fill {
  height: 100%;
  background: $gradient-primary-secondary;
  transition: width $transition-slow;
  border-radius: $radius-xs;
}

.progress-text {
  text-align: center;
  font-size: 20rpx;
  color: $color-text-secondary;
  font-weight: $font-weight-medium;
  margin-top: 8rpx;
}

// 快捷操作
.quick-actions {
  margin-bottom: $spacing;
}

.action-title {
  font-size: 28rpx;
  font-weight: $font-weight-semibold;
  color: $color-text-primary;
  text-align: center;
  margin-bottom: $spacing-lg;
}

.action-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12rpx;
}

// 快捷操作卡片
.action-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: $spacing-md;
  padding: $spacing-md 12rpx;
  background: $color-bg-primary;
  border: 1rpx solid $color-border-primary;
  border-radius: $radius-lg;
  transition: all $transition-slow;

  &:active {
    transform: scale(0.95);
    box-shadow: $shadow-md;
  }
}

.action-icon {
  width: 44rpx;
  height: 44rpx;
}

.action-icon-wrapper {
  width: 48rpx;
  height: 48rpx;
  border-radius: $radius-full;
  background: $color-primary-lighter;
  display: flex;
  align-items: center;
  justify-content: center;
}

.action-label {
  font-size: 22rpx;
  font-weight: $font-weight-medium;
  color: $color-text-primary;
  text-align: center;
  line-height: 1.2;
}


</style>
