<template>
  <view>
    <!-- 导航栏 -->
    <wd-navbar fixed placeholder title="睡眠记录" left-arrow safeAreaInsetTop>
      <template #capsule>
        <wd-navbar-capsule @back="goBack" @back-home="goBackHome" />
      </template>
    </wd-navbar>

    <view class="sleep-page">
      <!-- 当前状态卡片 -->
      <view class="status-card">
        <view v-if="ongoingRecord" class="status-content">
          <view class="status-icon">💤</view>
          <view class="status-info">
            <text class="status-text">宝宝正在睡觉</text>
            <view class="sleep-duration">
              <text class="duration">{{ formattedTime }}</text>
              <text class="duration-label">已睡眠</text>
            </view>
          </view>
        </view>
        <view v-else class="status-content">
          <view class="status-icon">👀</view>
          <view class="status-info">
            <text class="status-text">宝宝醒着</text>
            <text class="status-subtitle">点击下方开始记录睡眠</text>
          </view>
        </view>
      </view>

      <!-- 睡眠类型选择 -->
      <view v-if="!ongoingRecord" class="sleep-type-card">
        <view class="card-header">
          <text class="card-title">睡眠类型</text>
        </view>
        <wd-radio-group v-model="sleepType">
          <wd-radio value="nap">小睡</wd-radio>
          <wd-radio value="night">夜间长睡</wd-radio>
        </wd-radio-group>
      </view>

      <!-- 操作按钮 -->
      <view class="action-buttons">
        <wd-button
          v-if="!ongoingRecord"
          type="primary"
          size="large"
          block
          @click="startSleep"
        >
          <view class="button-content">
            <text class="icon">💤</text>
            <text>开始睡觉</text>
          </view>
        </wd-button>

        <wd-button v-else type="success" size="large" block @click="endSleep">
          <view class="button-content">
            <text class="icon">✓</text>
            <text>宝宝醒了</text>
          </view>
        </wd-button>
      </view>

      <!-- 快速补记睡眠卡片 -->
      <view v-if="!ongoingRecord" class="quick-record-card">
        <view class="card-header">
          <text class="card-title">快速补记</text>
        </view>
        <wd-button
          type="info"
          size="large"
          block
          @click="showQuickRecordModal = true"
        >
          <view class="button-content">
            <text class="icon">⏰</text>
            <text>补记历史睡眠</text>
          </view>
        </wd-button>
      </view>

      <!-- 最近记录卡片 -->
      <view v-if="lastRecord && !ongoingRecord" class="last-record-card">
        <view class="card-header">
          <text class="card-title">上次睡眠</text>
        </view>
        <view class="record-item">
          <view class="record-left">
            <view class="record-label">{{
              lastRecord.type === "nap" ? "小睡" : "夜间长睡"
            }}</view>
            <text class="record-time">{{ formatRecordTime(lastRecord) }}</text>
          </view>
          <view class="record-right">
            <text class="duration-text">{{
              formatDuration(lastRecord.duration || 0)
            }}</text>
          </view>
        </view>
      </view>
    </view>

    <!-- 快速补记睡眠对话框 -->
    <wd-popup
      v-model="showQuickRecordModal"
      position="bottom"
      :safe-area-inset-bottom="true"
    >
      <view class="quick-record-modal">
        <view class="modal-header">
          <view class="modal-title">补记睡眠</view>
          <view class="close-btn" @click="showQuickRecordModal = false">✕</view>
        </view>

        <view class="quick-record-form">
          <!-- 睡眠类型 -->
          <view class="form-item">
            <view class="form-label">睡眠类型</view>
            <wd-radio-group v-model="quickRecord.type" shape="button">
              <wd-radio :value="'nap'">小睡</wd-radio>
              <wd-radio :value="'night'">夜间长睡</wd-radio>
            </wd-radio-group>
          </view>

          <!-- 开始时间 -->
          <wd-datetime-picker
            label="开始时间"
            size="large"
            v-model="quickRecord.startTime"
            :min-date="minDateTime"
            :max-date="maxDateTime"
          />
          <!-- 结束时间 -->
          <wd-datetime-picker
            label="结束时间"
            size="large"
            v-model="quickRecord.endTime"
            :min-date="minDateTime"
            :max-date="maxDateTime"
          />

          <!-- 操作按钮 -->
          <view class="modal-actions">
            <wd-button
              size="large"
              @click="showQuickRecordModal = false"
              type="info"
            >
              取消
            </wd-button>
            <wd-button
              type="primary"
              size="large"
              @click="handleQuickSleepConfirm"
            >
              {{ isEditing ? "更新记录" : "确定" }}
            </wd-button>
          </view>
        </view>
      </view>
    </wd-popup>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from "vue";
import { onLoad, onShow } from "@dcloudio/uni-app";
import { currentBabyId, currentBaby } from "@/store/baby";
import { getUserInfo } from "@/store/user";
import { formatDate, formatDuration } from "@/utils/date";
import { goBack, goBackHome, padZero } from "@/utils/common";
import {
  StorageKeys,
  getStorage,
  setStorage,
  removeStorage,
} from "@/utils/storage";
import type { SleepRecord } from "@/types";

// 直接调用 API 层
import * as sleepApi from "@/api/sleep";

// 编辑模式相关
const editId = ref<string>("");
const isEditing = computed(() => !!editId.value);

// 临时睡眠记录类型
interface TempSleepRecording {
  babyId: string;
  type: "nap" | "night";
  startTime: number; // 开始时间戳(毫秒)
}

const minDateTime = ref(
  Date.parse(currentBaby.value?.birthDate || "2015-01-01")
); // 最小: 出生日期或2015-01-01
const maxDateTime = ref(new Date().getTime()); // 最大: 当前时间

// 睡眠类型
const sleepType = ref<"nap" | "night">("nap");

// 进行中的睡眠记录
const ongoingRecord = ref<SleepRecord | null>(null);

// 最后一次睡眠记录
const lastRecord = ref<SleepRecord | null>(null);

// 快速补记相关
const showQuickRecordModal = ref(false);

const quickRecord = ref({
  type: "nap" as "nap" | "night",
  startTime: new Date().getTime(),
  endTime: new Date().getTime(),
});

// 快速补记时间格式化
const formatQuickTime = (date: Date): string => {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const day = String(date.getDate()).padStart(2, "0");
  const hours = String(date.getHours()).padStart(2, "0");
  const minutes = String(date.getMinutes()).padStart(2, "0");
  return `${year}-${month}-${day} ${hours}:${minutes}`;
};

// 处理快速补记睡眠
const handleQuickSleepConfirm = async () => {
  const user = getUserInfo();
  if (!user) {
    uni.showToast({
      title: "请先登录",
      icon: "none",
    });
    return;
  }

  if (!currentBaby.value) {
    uni.showToast({
      title: "请先选择宝宝",
      icon: "none",
    });
    return;
  }

  // 验证时间
  if (quickRecord.value.startTime >= quickRecord.value.endTime) {
    uni.showToast({
      title: "开始时间必须早于结束时间",
      icon: "none",
    });
    return;
  }

  try {
    const elapsedSeconds = Math.floor(
      (quickRecord.value.endTime - quickRecord.value.startTime) / 1000
    );

    if (isEditing.value) {
      // 更新模式
      await sleepApi.apiUpdateSleepRecord(editId.value, {
        babyId: currentBabyId.value,
        sleepType: quickRecord.value.type,
        startTime: quickRecord.value.startTime,
        endTime: quickRecord.value.endTime,
        duration: elapsedSeconds,
      });

      uni.showToast({
        title: "更新成功",
        icon: "success",
      });
    } else {
      // 创建模式
      await sleepApi.apiCreateSleepRecord({
        babyId: currentBabyId.value,
        sleepType: quickRecord.value.type,
        startTime: quickRecord.value.startTime,
        endTime: quickRecord.value.endTime,
        duration: elapsedSeconds,
      });

      uni.showToast({
        title: "保存成功",
        icon: "success",
      });
    }

    showQuickRecordModal.value = false;

    // 重置表单
    quickRecord.value = {
      type: "nap",
      startTime: 0,
      endTime: 0,
    };

    setTimeout(() => {
      uni.navigateBack();
    }, 1000);
  } catch (error: any) {
    console.error("保存快速补记睡眠失败:", error);
    uni.showToast({
      title: error.message || "保存失败",
      icon: "none",
    });
  }
};

// 定时器相关
const timerRunning = ref(false);
const startTime = ref(0); // 开始时间戳 (毫秒)
const timerTrigger = ref(0); // 用于触发视图更新的虚拟响应式值
const tempRecordCheckDone = ref(false); // 防止重复检测临时记录
let timerInterval: number | null = null;

// 格式化时间显示 - 基于开始时间戳计算
const formattedTime = computed(() => {
  // 依赖 timerTrigger 以触发定期更新
  timerTrigger.value; // 访问此值以建立依赖关系

  if (!timerRunning.value || startTime.value === 0) {
    return "00:00:00";
  }

  const elapsedSeconds = Math.floor((Date.now() - startTime.value) / 1000);
  const hours = Math.floor(elapsedSeconds / 3600);
  const minutes = Math.floor((elapsedSeconds % 3600) / 60);
  const seconds = elapsedSeconds % 60;

  return `${padZero(hours)}:${padZero(minutes)}:${padZero(seconds)}`;
});

// 保存临时睡眠记录到本地
const saveTempRecord = () => {
  const tempRecord: TempSleepRecording = {
    babyId: currentBabyId.value,
    type: sleepType.value,
    startTime: startTime.value,
  };
  setStorage(StorageKeys.TEMP_SLEEP_RECORDING, tempRecord);
  console.log("临时记录已保存:", tempRecord);
};

// 清除临时睡眠记录
const clearTempRecord = () => {
  removeStorage(StorageKeys.TEMP_SLEEP_RECORDING);
  tempRecordCheckDone.value = false; // 重置标志，允许下次检测
  console.log("临时记录已清除");
};

// 恢复临时睡眠记录
const restoreTempRecord = (tempRecord: TempSleepRecording) => {
  const user = getUserInfo();
  if (!user) {
    return;
  }

  sleepType.value = tempRecord.type;
  startTime.value = tempRecord.startTime;
  timerRunning.value = true;

  // 创建本地睡眠记录对象以显示计时器
  ongoingRecord.value = {
    id: `temp_${Date.now()}`, // 临时ID
    babyId: tempRecord.babyId,
    startTime: tempRecord.startTime,
    type: tempRecord.type,
    createBy: user.openid,
    createByName: user.nickName,
    createByAvatar: user.avatarUrl,
    createTime: Date.now(),
  };

  // 启动定时器更新显示
  timerInterval = setInterval(() => {
    // 每秒改变 timerTrigger 以触发计算属性重新计算
    timerTrigger.value++;
  }, 1000) as unknown as number;

  console.log(
    "临时记录已恢复, 已过时长:",
    Math.floor((Date.now() - tempRecord.startTime) / 1000),
    "秒"
  );
};

// 开始睡觉
const startSleep = async () => {
  const user = getUserInfo();
  if (!user) {
    uni.showToast({
      title: "请先登录",
      icon: "none",
    });
    return;
  }

  if (!currentBaby.value) {
    uni.showToast({
      title: "请先选择宝宝",
      icon: "none",
    });
    return;
  }

  try {
    // 使用本地时间戳开始计时
    startTime.value = Date.now();
    timerRunning.value = true;

    // 创建本地睡眠记录对象以显示计时器
    ongoingRecord.value = {
      id: `temp_${Date.now()}`, // 临时ID
      babyId: currentBabyId.value,
      startTime: startTime.value,
      type: sleepType.value,
      createBy: user.openid,
      createByName: user.nickName,
      createByAvatar: user.avatarUrl,
      createTime: Date.now(),
    };

    // 保存临时记录到本地存储
    saveTempRecord();

    // 启动定时器以每秒更新视图
    timerInterval = setInterval(() => {
      // 每秒改变 timerTrigger 以触发计算属性重新计算
      timerTrigger.value++;
    }, 1000) as unknown as number;

    uni.showToast({
      title: "开始记录睡眠",
      icon: "success",
    });

    console.log("开始计时");
  } catch (error: any) {
    uni.showToast({
      title: error.message || "开始失败",
      icon: "none",
    });
  }
};

// 结束睡觉
const endSleep = async () => {
  if (!timerRunning.value || startTime.value === 0) {
    return;
  }

  const user = getUserInfo();
  if (!user) {
    uni.showToast({
      title: "请先登录",
      icon: "none",
    });
    return;
  }

  try {
    // 停止计时器
    if (timerInterval) {
      clearInterval(timerInterval);
      timerInterval = null;
    }

    timerRunning.value = false;

    // 计算总时长(秒)
    const elapsedSeconds = Math.floor((Date.now() - startTime.value) / 1000);

    console.log("停止计时,总时长:", elapsedSeconds, "秒");

    // 调用 API 创建睡眠记录
    await sleepApi.apiCreateSleepRecord({
      babyId: currentBabyId.value,
      sleepType: sleepType.value,
      startTime: startTime.value,
      endTime: Date.now(),
      duration: elapsedSeconds, // 添加时长字段
    });

    console.log("睡眠记录保存成功");

    // 清除临时记录和进行中的记录
    clearTempRecord();
    ongoingRecord.value = null;

    uni.showToast({
      title: "保存成功",
      icon: "success",
    });

    setTimeout(() => {
      uni.navigateBack();
    }, 1000);
  } catch (error: any) {
    console.error("保存睡眠记录失败:", error);

    // 如果保存失败,恢复计时器
    timerRunning.value = true;
    timerInterval = setInterval(() => {
      timerTrigger.value++;
    }, 1000) as unknown as number;

    uni.showToast({
      title: error.message || "保存失败",
      icon: "none",
    });
  }
};

// 页面卸载时清除计时器
onUnmounted(() => {
  if (timerInterval) {
    clearInterval(timerInterval);
  }
});

// 页面加载时检测 editId 参数
onLoad((options) => {
  if (options?.editId) {
    editId.value = options.editId;
    loadSleepRecord(options.editId);
  }
});

// 加载睡眠记录数据(用于编辑模式)
const loadSleepRecord = async (recordId: string) => {
  try {
    const apiResponse = await sleepApi.apiGetSleepRecordById(recordId);

    // 转换API响应为前端类型
    const record = sleepApi.transformSleepRecordResponse(apiResponse);

    // 设置睡眠类型
    sleepType.value = record.type || 'nap';

    // 设置快速补记表单的开始和结束时间
    quickRecord.value = {
      type: record.type || 'nap',
      startTime: record.startTime,
      endTime: record.endTime || Date.now(),
    };

    // 直接打开快速补记模态框
    showQuickRecordModal.value = true;

    console.log("已加载记录数据:", record);
  } catch (error: any) {
    console.error("加载记录失败:", error);
    uni.showToast({
      title: error.message || "加载记录失败",
      icon: "none",
    });
    setTimeout(() => {
      uni.navigateBack();
    }, 1500);
  }
};

// 页面加载
onMounted(() => {
  if (!currentBaby.value) {
    uni.showToast({
      title: "请先选择宝宝",
      icon: "none",
    });
    setTimeout(() => {
      uni.navigateBack();
    }, 1500);
    return;
  }

  // 编辑模式下不检查临时记录
  if (!isEditing.value) {
    checkTempRecord();
  }
});

// 页面显示时也检测(从其他页面返回)
onShow(() => {
  // 每次页面显示时重置检测标志，允许再次检测
  tempRecordCheckDone.value = false;
  checkTempRecord();
});

// 监听睡眠类型变化,如果正在计时则更新临时记录
watch(
  () => sleepType.value,
  () => {
    if (timerRunning.value && startTime.value > 0) {
      saveTempRecord();
      console.log("睡眠类型已更改,临时记录已更新");
    }
  }
);

// 检测并处理临时睡眠记录
const checkTempRecord = () => {
  // 如果已经在计时,不重复检测
  if (timerRunning.value) {
    return;
  }

  // 如果已经检测过本次，不再重复检测（防止 onMounted 和 onShow 重复调用）
  if (tempRecordCheckDone.value) {
    return;
  }

  const tempRecord = getStorage<TempSleepRecording>(
    StorageKeys.TEMP_SLEEP_RECORDING
  );

  if (!tempRecord) {
    tempRecordCheckDone.value = true; // 标记已检测
    return;
  }

  // 检查临时记录是否属于当前宝宝
  if (tempRecord.babyId !== currentBabyId.value) {
    console.log("临时记录不属于当前宝宝,已忽略");
    tempRecordCheckDone.value = true; // 标记已检测
    return;
  }

  // 标记已检测（在显示弹窗前）
  tempRecordCheckDone.value = true;

  // 计算已过时长
  const elapsedSeconds = Math.floor((Date.now() - tempRecord.startTime) / 1000);
  const hours = Math.floor(elapsedSeconds / 3600);
  const minutes = Math.floor((elapsedSeconds % 3600) / 60);
  const seconds = elapsedSeconds % 60;

  console.log(
    "检测到临时记录,已过时长:",
    `${hours}小时${minutes}分${seconds}秒`
  );

  // 弹窗询问用户
  uni.showModal({
    title: "未完成的睡眠记录",
    content: `检测到您之前有一次未完成的${
      tempRecord.type === "nap" ? "小睡" : "夜间长睡"
    }记录,已过 ${hours} 小时 ${minutes} 分钟 ${seconds} 秒,是否继续?`,
    confirmText: "继续",
    cancelText: "重新开始",
    success: (res) => {
      if (res.confirm) {
        // 用户选择继续
        console.log("用户选择继续临时记录");
        // 恢复临时记录
        restoreTempRecord(tempRecord);
      } else {
        // 用户选择重新开始
        console.log("用户选择重新开始,清除临时记录");
        clearTempRecord();
      }
    },
  });
};

// 格式化记录时间
const formatRecordTime = (record: SleepRecord) => {
  return formatDate(record.startTime, "MM-DD HH:mm");
};
</script>

<style lang="scss" scoped>
@import "@/styles/colors.scss";

// 设计系统变量
$spacing: 20rpx;
$radius-lg: 16rpx;
$radius-md: 12rpx;

.sleep-page {
  min-height: 100vh;
  background: $color-bg-secondary;
  padding: $spacing;
  padding-top: 12rpx;
  padding-bottom: calc(100rpx + env(safe-area-inset-bottom));
}

// ============ 当前状态卡片 ============
.status-card {
  background: $color-bg-primary;
  border: 1rpx solid $color-border-primary;
  border-radius: $radius-lg;
  padding: 40rpx $spacing;
  margin-bottom: 24rpx;
  box-shadow: $shadow-primary-sm;
}

.status-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20rpx;
  text-align: center;
}

.status-icon {
  font-size: 80rpx;
  line-height: 1;
}

.status-info {
  display: flex;
  flex-direction: column;
  gap: 12rpx;
  align-items: center;
}

.status-text {
  font-size: 36rpx;
  font-weight: 600;
  color: $color-text-primary;
}

.status-subtitle {
  font-size: 24rpx;
  color: $color-text-secondary;
}

.sleep-duration {
  display: flex;
  flex-direction: column;
  gap: 8rpx;
  align-items: center;
  margin-top: 12rpx;
}

.duration {
  font-size: 64rpx;
  font-weight: bold;
  color: $color-primary;
  font-family: "Courier New", monospace;
  line-height: 1;
}

.duration-label {
  font-size: 24rpx;
  color: $color-text-secondary;
  font-weight: 500;
}

// ============ 睡眠类型卡片 ============
.sleep-type-card {
  background: $color-bg-primary;
  border: 1rpx solid $color-border-primary;
  border-radius: $radius-lg;
  padding: $spacing;
  margin-bottom: 24rpx;
  box-shadow: $shadow-primary-sm;
}

.card-header {
  margin-bottom: 20rpx;
}

.card-title {
  font-size: 28rpx;
  font-weight: 600;
  color: $color-text-primary;
}

// ============ 操作按钮 ============
.action-buttons {
  margin-bottom: 24rpx;
}

.button-content {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12rpx;

  .icon {
    font-size: 32rpx;
  }
}

// ============ 快速补记卡片 ============
.quick-record-card {
  background: $color-bg-primary;
  border: 1rpx solid $color-border-primary;
  border-radius: $radius-lg;
  padding: $spacing;
  margin-bottom: 24rpx;
  box-shadow: $shadow-primary-sm;
}

// ============ 上次睡眠卡片 ============
.last-record-card {
  background: $color-bg-primary;
  border: 1rpx solid $color-border-primary;
  border-radius: $radius-lg;
  padding: $spacing;
  box-shadow: $shadow-primary-sm;
}

.record-item {
  background: $color-bg-secondary;
  border-radius: $radius-md;
  padding: 16rpx $spacing;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.record-left {
  display: flex;
  flex-direction: column;
  gap: 8rpx;
}

.record-label {
  font-size: 24rpx;
  color: $color-text-secondary;
  font-weight: 500;
}

.record-time {
  font-size: 26rpx;
  color: $color-text-primary;
  font-weight: 600;
}

.record-right {
  display: flex;
  align-items: center;
}

.duration-text {
  color: $color-primary;
  font-weight: bold;
  font-size: 28rpx;
}

// ============ 模态框样式 ============
.quick-record-modal {
  background: white;
  border-radius: 16rpx 16rpx 0 0;
  padding: 30rpx;
  max-height: 80vh;
  overflow-y: auto;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30rpx;
}

.modal-title {
  font-size: 36rpx;
  font-weight: bold;
  color: #333;
}

.close-btn {
  font-size: 48rpx;
  color: #999;
  line-height: 1;
  cursor: pointer;
}

.quick-record-form {
  padding: 0;
}

.form-item {
  margin-bottom: 30rpx;

  &:last-child {
    margin-bottom: 0;
  }
}

.form-label {
  font-size: 28rpx;
  font-weight: 500;
  margin-bottom: 12rpx;
  color: #333;
}

.modal-actions {
  display: grid;
  grid-template-columns: 1fr 2fr;
  gap: 20rpx;
  margin-top: 40rpx;
}
</style>
