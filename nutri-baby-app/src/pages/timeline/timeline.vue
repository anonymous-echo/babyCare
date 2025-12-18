<template>
  <view class="timeline-page">
    <!-- 固定顶部筛选条 -->
    <view class="filter-fixed-top">
      <!-- 记录类型筛选 -->
      <wd-tabs v-model="recordTypeFilter" swipeable class="type-tabs">
        <wd-tab title="全部" name="all" />
        <wd-tab title="喂养" name="feeding" />
        <wd-tab title="换尿布" name="diaper" />
        <wd-tab title="睡眠" name="sleep" />
        <wd-tab title="成长" name="growth" />
      </wd-tabs>

      <!-- 日期筛选 -->
      <!-- <view class="date-filter"> -->
      <wd-radio-group
        v-model="filterType"
        inline
        shape="button"
        @change="handleDateFilterChange"
      >
        <wd-radio value="today">今天</wd-radio>
        <wd-radio value="week">本周</wd-radio>
        <wd-radio value="month">本月</wd-radio>
        <wd-radio value="custom">自定义</wd-radio>
      </wd-radio-group>
      <wd-datetime-picker
        id="custom-date-picker"
        ref="dateTimePickerRef"
        style="display: none"
        v-model="selectedDateTimestamp"
        @confirm="onDateConfirm"
        :minDate="minDate"
        :maxDate="maxDate"
      >
      </wd-datetime-picker>
      <!-- </view> -->
    </view>

    <!-- 内容区域 -->
    <view class="timeline-list">
      <view v-if="isLoggedIn">
        <view v-if="groupedRecords.length === 0" class="empty-state">
          <wd-status-tip
            :description="emptyDescription"
            tip="当前时间段暂无数据"
          />
        </view>

        <view v-else>
          <view
            v-for="group in groupedRecords"
            :key="group.date"
            class="date-group"
            :data-date="group.date"
          >
            <!-- 日期标题（浮动） -->
            <view class="date-header">{{ group.dateText }}</view>

            <!-- 记录列表 -->
            <view
              v-for="record in group.records"
              :key="record.id"
              class="record-item"
              :class="`record-${record.type}`"
            >
              <!-- 时间轴圆点 -->
              <view class="timeline-dot" :class="`dot-${record.type}`" />
              <view class="timeline-line" />

              <!-- 记录内容 -->
              <wd-card custom-class="record-card">
                <template #title>
                  <view class="record-header">
                    <view class="record-type">
                      <image :src="record.iconUrl" mode="aspectFill" class="type-icon" />
                      <text class="type-name">{{ record.typeName }}</text>
                    </view>
                    <view class="record-meta">
                      <text class="record-time">{{ record.timeText }}</text>
                      <text class="record-creator">
                        {{ record.createName }}
                        <text v-if="record.relationship" class="relationship">({{ record.relationship }})</text>
                      </text>
                    </view>
                  </view>
                </template>

                <!-- 详细信息显示 -->
                <view class="record-details">
                  <view class="detail-line">{{ record.detail }}</view>
                  <!-- 备注信息 -->
                  <view
                    v-if="record.originalRecord.note"
                    class="detail-line note"
                  >
                    <text class="label">备注:</text>
                    <text class="value">{{ record.originalRecord.note }}</text>
                  </view>
                </view>

                <template #footer>
                  <view class="record-actions">
                    <wd-button
                      size="small"
                      type="primary"
                      @click="editRecord(record)"
                    >
                      编辑
                    </wd-button>
                    <wd-button
                      size="small"
                      type="info"
                      @click="deleteRecord(record)"
                    >
                      删除
                    </wd-button>
                  </view>
                </template>
              </wd-card>
            </view>
          </view>

          <!-- 加载更多组件 -->
          <wd-loadmore
            :state="loadMoreState"
            @reload="loadMore"
            loading-text="加载中..."
            finished-text="没有更多了"
            error-text="加载失败，点击重试"
          />
        </view>
      </view>
      <view v-else>
        <wd-status-tip description="请先登录" tip="登录后查看数据..." />
      </view>
    </view>

  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from "vue";
import { onReachBottom, onPullDownRefresh } from "@dcloudio/uni-app";
import { isLoggedIn } from "@/store/user";
import { currentBaby } from "@/store/baby";
import {
  formatDate,
  isToday,
  getTodayStart,
  getWeekStart,
  getMonthStart,
} from "@/utils/date";
import { formatDuration } from "@/utils/common";

// 使用新的时间线聚合 API
import * as timelineApi from "@/api/timeline";
import type { TimelineItem } from "@/api/timeline";
import * as feedingApi from "@/api/feeding";
import * as diaperApi from "@/api/diaper";
import * as sleepApi from "@/api/sleep";

// 日期筛选
const filterType = ref<"today" | "week" | "month" | "custom">("today");
const customStartDate = ref(getTodayStart());
const customEndDate = ref(Date.now());

// 记录类型筛选
const recordTypeFilter = ref<"all" | "feeding" | "diaper" | "sleep" | "growth">(
  "all"
);

// Wot UI 日期选择器相关
const selectedDateTimestamp = ref<number[]>([]);

// 时间线数据(从聚合 API 获取)
const timelineItems = ref<TimelineItem[]>([]);
const totalRecords = ref(0);
const dateTimePickerRef = ref<any>(null);
// 日期选择器的最小和最大日期
// 最小日期为当前宝宝的出生日期
const minDate = ref(Date.parse(currentBaby.value?.birthDate || ""));
const maxDate = ref(Date.now());
// 分页相关
const currentPage = ref(1);
const pageSize = ref(5);
const isLoadingMore = ref(false);
const hasMore = ref(true);
const handleDateFilterChange = ({ value }: { value: "today" | "week" | "month" | "custom" }) => {
  console.log("Date filter changed to:", value);
  if (value === "custom") {
    console.log("Opening custom date picker");
    dateTimePickerRef.value?.open();
  }
  // 重置分页，重新加载数据
  currentPage.value = 1;
  hasMore.value = true;
  loadRecords(true)
};

// 展示用的记录接口
interface TimelineRecord {
  id: string;
  type: "feeding" | "diaper" | "sleep" | "growth";
  time: number;
  iconUrl: string;
  typeName: string;
  timeText: string;
  detail: string;
  originalRecord: any;
  createName: string;   // 创建者昵称
  relationship: string; // 创建者与宝宝的关系
}

// 转换时间线数据为展示格式
const allRecords = computed<TimelineRecord[]>(() => {
  if (!currentBaby.value) return [];

  const records: TimelineRecord[] = [];

  timelineItems.value.forEach((item) => {
    let iconUrl = "";
    let typeName = "";
    let detail = "";

    if (item.recordType === "feeding") {
      const record = item.detail as feedingApi.FeedingRecordResponse;
      iconUrl = "/static/breastfeeding.svg";
      typeName = "喂养";

      if (record.feedingType === "breast") {
        detail = `母乳喂养 ${formatDuration(record.duration || 0)}`;
        const feedingDetail = record.detail;
        if (feedingDetail && feedingDetail.type === "breast") {
          const breastSide = feedingDetail.side;
          if (breastSide === "left") detail += " (左侧)";
          else if (breastSide === "right") detail += " (右侧)";
          else if (breastSide === "both") detail += " (双侧)";
        }
      } else if (record.feedingType === "bottle") {
        const feedingDetail = record.detail;
        if (feedingDetail && feedingDetail.type === "bottle") {
          detail = `奶瓶喂养 ${record.amount}${feedingDetail.unit || "ml"}`;
          detail +=
            feedingDetail.bottleType === "formula" ? " (配方奶)" : " (母乳)";
        } else {
          detail = `奶瓶喂养 ${record.amount}ml`;
        }
      } else {
        const feedingDetail = record.detail;
        if (feedingDetail && feedingDetail.type === "food") {
          detail = `辅食: ${feedingDetail.foodName || "未知"}`;
        } else {
          detail = "辅食";
        }
      }
    } else if (item.recordType === "diaper") {
      const record = item.detail as diaperApi.DiaperRecordResponse;
      iconUrl = "/static/baby_changing_station.svg";
      typeName = "换尿布";

      if (record.diaperType === "pee") detail = "小便";
      else if (record.diaperType === "poop") detail = "大便";
      else detail = "小便+大便";

      if (record.pooColor) detail += ` (${record.pooColor})`;
    } else if (item.recordType === "sleep") {
      const record = item.detail as sleepApi.SleepRecordResponse;
      iconUrl = "/static/moon_stars.svg";
      typeName = "睡眠";

      const duration = record.duration || 0;
      detail = `${
        record.sleepType === "nap" ? "小睡" : "夜间睡眠"
      } ${formatDuration(duration)}`;
    } else if (item.recordType === "growth") {
      iconUrl = "/static/monitoring.svg";
      typeName = "成长";
      const record = item.detail as any;
      const parts: string[] = [];
      if (record.height) parts.push(`身高 ${record.height}cm`);
      if (record.weight) parts.push(`体重 ${record.weight}kg`);
      if (record.headCircumference)
        parts.push(`头围 ${record.headCircumference}cm`);
      detail = parts.join(", ");
    }

    records.push({
      id: item.recordId,
      type: item.recordType,
      time: item.eventTime,
      iconUrl,
      typeName,
      timeText: formatDate(item.eventTime, "HH:mm"),
      detail,
      originalRecord: item.detail,
      createName: item.createName || '',
      relationship: item.relationship || '',
    });
  });

  // 后端已根据 recordType 筛选，前端直接返回
  return records;
});

// 按日期分组
const groupedRecords = computed(() => {
  const groups: {
    date: string;
    dateText: string;
    records: TimelineRecord[];
  }[] = [];

  allRecords.value.forEach((record) => {
    const date = formatDate(record.time, "YYYY-MM-DD");
    let group = groups.find((g) => g.date === date);

    if (!group) {
      let dateText = formatDate(record.time, "MM月DD日");
      if (isToday(record.time)) {
        dateText = "今天 " + dateText;
      }

      group = { date, dateText, records: [] };
      groups.push(group);
    }

    group.records.push(record);
  });

  return groups;
});

// 获取记录类型的显示名称
const getRecordTypeName = (type: string): string => {
  const map: Record<string, string> = {
    feeding: "喂养",
    diaper: "换尿布",
    sleep: "睡眠",
    growth: "成长",
  };
  return map[type] || "记录";
};

// 空状态描述
const emptyDescription = computed(() => {
  if (!isLoggedIn.value) return "登录后查看记录";

  if (timelineItems.value.length === 0) return "当前时间段暂无数据";

  if (allRecords.value.length === 0 && recordTypeFilter.value !== "all") {
    return `当前时间段暂无${getRecordTypeName(recordTypeFilter.value)}记录`;
  }

  return "暂无记录";
});

// 加载时间线记录 (使用新的聚合 API)
const loadRecords = async (isRefresh: boolean = false, pullDown: boolean = false) => {
  if (!currentBaby.value) return;

  // 防止重复加载
  if (isLoadingMore.value) {
    console.log("[Timeline] 正在加载中，跳过重复请求");
    return;
  }

  // 如果不是刷新且没有更多数据，直接返回
  if (!isRefresh && !hasMore.value) {
    console.log("[Timeline] 没有更多数据，跳过加载");
    return;
  }

  // 如果是刷新，重置分页
  if (isRefresh) {
    currentPage.value = 1;
    timelineItems.value = [];
    hasMore.value = true;
  }

  const babyId = currentBaby.value.babyId;
  const pageToLoad = currentPage.value;

  // 计算时间范围
  let startTime = 0;
  let endTime = Date.now();

  if (filterType.value === "today") {
    startTime = getTodayStart();
  } else if (filterType.value === "week") {
    startTime = getWeekStart();
  } else if (filterType.value === "month") {
    startTime = getMonthStart();
  } else if (filterType.value === "custom") {
    startTime = customStartDate.value;
    endTime = customEndDate.value;
  }

  try {
    isLoadingMore.value = true;
    
    if (pullDown) {
      uni.showNavigationBarLoading();
    } else if (isRefresh) {
      uni.showLoading({ title: "加载中", mask: false });
    }
    
    const response = await timelineApi.apiFetchTimeline({
      babyId,
      startTime,
      endTime,
      recordType: recordTypeFilter.value === 'all' ? '' : recordTypeFilter.value,
      page: pageToLoad,
      pageSize: pageSize.value,
    });

    const newItems = response.data.items || [];
    
    // 如果是刷新，替换数据；否则追加数据
    if (isRefresh) {
      timelineItems.value = newItems;
    } else {
      timelineItems.value.push(...newItems);
    }

    // 页码递增（无论是否有数据，都要递增以避免重复请求同一页）
    currentPage.value = pageToLoad + 1;

    totalRecords.value = response.data.total;

    // 判断是否还有更多数据：返回的数据少于请求的 pageSize，说明已经没有更多了
    hasMore.value = newItems.length >= pageSize.value;

    console.log("[Timeline] 加载数据完成", {
      loadedPage: pageToLoad,
      nextPage: currentPage.value,
      newItemsCount: newItems.length,
      totalLoaded: timelineItems.value.length,
      total: response.data.total,
      hasMore: hasMore.value,
    });
    
    if (pullDown) {
      uni.showToast({
        title: "刷新成功",
        icon: "success",
        duration: 1200,
      });
    }
  } catch (error) {
    console.error("加载时间线失败:", error);
    if (pullDown) {
      uni.showToast({
        title: "刷新失败",
        icon: "none",
        duration: 1500,
      });
    } else {
      uni.showToast({
        title: "加载数据失败",
        icon: "none",
      });
    }
  } finally {
    if (pullDown) {
      uni.hideNavigationBarLoading();
      uni.stopPullDownRefresh();
    } else if (isRefresh) {
      uni.hideLoading();
    }
    isLoadingMore.value = false;
  }
};

// 页面加载
onMounted(() => {
  if (isLoggedIn.value) {
    loadRecords(true);
  }
});

// 监听记录类型变化，重新加载数据
watch(recordTypeFilter, () => {
  currentPage.value = 1;
  hasMore.value = true;
  loadRecords(true);
});

// 下拉刷新
onPullDownRefresh(async () => {
  if (!isLoggedIn.value || !currentBaby.value) {
    uni.stopPullDownRefresh();
    uni.hideNavigationBarLoading();
    return;
  }
  await loadRecords(true, true);
});

// 页面滚动到底部时触发
onReachBottom(() => {
  console.log("[Timeline] onReachBottom 触发", {
    hasMore: hasMore.value,
    isLoadingMore: isLoadingMore.value,
  });

  // loadRecords 内部已经有防重复加载的逻辑
  loadRecords(false);
});

// Wot UI 日期选择器的 confirm 事件处理
const onDateConfirm = ({ value }: { value: number[] }) => {
  console.log("[Timeline] 选择的日期时间戳范围:", value);

  if (!value || value.length === 0 || !value[0]) return;

  // value 是时间戳数组
  const timestamp = value[0];
  const endTimestamp = value[1] || timestamp;
  // 更新时间戳
  selectedDateTimestamp.value = value;

  // 设置当天的起止时间
  customStartDate.value = new Date(timestamp).setHours(0, 0, 0, 0);
  customEndDate.value = new Date(endTimestamp).setHours(23, 59, 59, 999);
  filterType.value = "custom";

  // 重新加载数据（从第一页开始）
  loadRecords(true);
};

// 编辑记录 - 跳转到对应的添加页面
const editRecord = (record: TimelineRecord) => {
  let url = "";

  switch (record.type) {
    case "feeding":
      url = `/pages/record/feeding/feeding?editId=${record.id}`;
      break;
    case "sleep":
      url = `/pages/record/sleep/sleep?editId=${record.id}`;
      break;
    case "diaper":
      url = `/pages/record/diaper/diaper?editId=${record.id}`;
      break;
    case "growth":
      url = `/pages/record/growth/growth?editId=${record.id}`;
      break;
  }

  if (url) {
    uni.navigateTo({ url });
  }
};

// 删除记录
const deleteRecord = async (record: TimelineRecord) => {
  uni.showModal({
    title: "确认删除",
    content: "确定要删除这条记录吗?",
    success: async (res) => {
      if (res.confirm) {
        try {
          if (record.type === "feeding") {
            await feedingApi.apiDeleteFeedingRecord(record.id);
          } else if (record.type === "diaper") {
            await diaperApi.apiDeleteDiaperRecord(record.id);
          } else if (record.type === "sleep") {
            await sleepApi.apiDeleteSleepRecord(record.id);
          }

          uni.showToast({
            title: "删除成功",
            icon: "success",
          });

          // 重新加载记录
          await loadRecords();
        } catch (error: any) {
          uni.showToast({
            title: error.message || "删除失败",
            icon: "none",
          });
        }
      }
    },
  });
};

// 加载更多状态计算（后端已处理类型筛选，直接使用 hasMore）
const loadMoreState = computed<string>(() => {
  if (!isLoggedIn.value || !currentBaby.value) return "finished";

  if (timelineItems.value.length === 0) return hasMore.value ? "loading" : "finished";

  if (isLoadingMore.value) return "loading";

  if (!hasMore.value) return "finished";

  return "loading";
});

// 加载更多函数
const loadMore = () => {
  console.log("[Timeline] 点击重试加载");
  loadRecords(false);
};

</script>

<style lang="scss" scoped>
.timeline-page {
  min-height: 100vh;
  background: #f6f8f7;
  display: flex;
  flex-direction: column;
}

// ========== 固定顶部筛选条 ==========
.filter-fixed-top {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 100;
  background: white;
  box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.08);
  overflow-x: hidden;
}

.type-tabs {
  background: white;
  border-bottom: 1rpx solid #e8eef5;
}

:deep(.wd-tabs) {
  .wd-tab__item {
    font-size: 30rpx;
    font-weight: 500;
    padding: 26rpx 0;
    transition: all 0.3s ease;
  }

  .wd-tab__item--active {
    font-weight: 600;
    color: #7dd3a2;
    transform: scale(1.05);
  }

  .wd-tabs__line {
    height: 6rpx;
    border-radius: 3rpx;
    background: linear-gradient(90deg, #7dd3a2, #52c41a);
  }
}

.date-filter {
  background: #f6f8f7;
  display: flex;
  align-items: center;
  gap: 12rpx;
  border-bottom: 1rpx solid #ebebeb;
  min-height: 60rpx;
}

.quick-filters {
  display: flex;
  gap: 12rpx;
  flex: 1;
  align-items: center;
  justify-content: flex-end;
}

// 按钮样式优化
:deep(.wd-button) {
  transition: all 0.25s ease;
}

:deep(.wd-button--small) {
  font-size: 26rpx;
  padding: 0 24rpx;
  height: 60rpx;
  border-radius: 30rpx;
}

:deep(.wd-button--default) {
  background: white;
  color: #666;
  border-color: #e0e0e0;
}

:deep(.wd-button--default:active) {
  background: #f6f8f7;
}

:deep(.wd-button--primary:not(.wd-button--plain)) {
  box-shadow: 0 4rpx 12rpx rgba(125, 211, 162, 0.25);
  transform: translateY(-2rpx);
}

:deep(.wd-button--plain) {
  background: white;
}

// ========== 内容区域 ==========
.timeline-list {
  padding: 20rpx;
  padding-top: 180rpx; // 为固定的顶部预留空间
  padding-bottom: 40rpx;
  min-height: 100vh;
}

.empty-state {
  padding: 120rpx 0;
}

.date-group {
  margin-bottom: 40rpx;
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10rpx);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.date-header {
  width: 100%;
  font-size: 28rpx;
  font-weight: 600;
  color: #333;
  padding: 20rpx 0;
  padding-left: 16rpx;
  background: linear-gradient(to bottom, #f6f8f7 85%, transparent);
  z-index: 5;

  &::before {
    content: "";
    position: absolute;
    left: 0;
    top: 50%;
    transform: translateY(-50%);
    width: 6rpx;
    height: 28rpx;
    background: linear-gradient(180deg, #7dd3a2, #52c41a);
    border-radius: 3rpx;
  }
}

.record-item {
  position: relative;
  padding-left: 60rpx;
  margin-bottom: 20rpx;
  animation: slideIn 0.3s ease;

  &:last-child .timeline-line {
    display: none;
  }
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateX(-20rpx);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

.timeline-dot {
  position: absolute;
  left: 10rpx;
  top: 8rpx;
  width: 24rpx;
  height: 24rpx;
  border-radius: 50%;
  border: 4rpx solid;
  background: white;
  z-index: 2;
  box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.1);

  &.dot-feeding {
    border-color: #7dd3a2;
  }

  &.dot-diaper {
    border-color: #52c41a;
  }

  &.dot-sleep {
    border-color: #1890ff;
  }

  &.dot-growth {
    border-color: #722ed1;
  }
}

.timeline-line {
  position: absolute;
  left: 18rpx;
  top: 32rpx;
  bottom: -20rpx;
  width: 2rpx;
  background: #e8e8e8;
  z-index: 1;
}

// WotUI Card 组件自定义样式
:deep(.record-card) {
  border-radius: 12rpx;
  box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.04);
  transition: transform 0.2s ease, box-shadow 0.2s ease;

  &:active {
    transform: scale(0.98);
    box-shadow: 0 1rpx 4rpx rgba(0, 0, 0, 0.06);
  }
}

.record-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  width: 100%;
}

.record-type {
  display: flex;
  align-items: center;
  gap: 12rpx;
}

.record-meta {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 4rpx;
}

.type-icon {
  width: 32rpx;
  height: 32rpx;
}

.type-name {
  font-size: 28rpx;
  font-weight: bold;
  color: #1a1a1a;
}

.record-time {
  font-size: 24rpx;
  color: #999;
}

.record-creator {
  font-size: 22rpx;
  color: #666;
  display: flex;
  align-items: center;
  gap: 4rpx;
}

.relationship {
  color: #7dd3a2;
  font-weight: 500;
}

.record-details {
  margin-top: 16rpx;
}

.detail-line {
  font-size: 26rpx;
  color: #666;
  line-height: 1.6;
  margin-bottom: 8rpx;

  &.note {
    background: #f6f8f7;
    padding: 12rpx 16rpx;
    border-radius: 8rpx;
    margin-top: 12rpx;

    .label {
      color: #999;
      margin-right: 8rpx;
    }

    .value {
      color: #333;
    }
  }
}

.record-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12rpx;
  margin-top: 16rpx;
}

:deep(.wd-tabs) {
  width: 100% !important;
  height: 80rpx !important;
}
:deep(.wd-datetime-picker .wd-cell) {
  display: none !important;
}
:deep(.wd-radio-group) {
  border-top: 1rpx solid #e8eef5 !important;
  padding: 18rpx 20rpx !important;
}
:deep(.wd-radio.is-button .wd-radio__label) {
  font-size: 22rpx !important;
  min-width: 0 !important;
  padding: 0 12rpx !important;
  width: 100rpx !important;
  height: 42rpx !important;
  line-height: 42rpx;
}
</style>
