<template>
  <view class="diaper-page">
    <!-- 排泄类型快捷按钮 -->
    <view class="quick-buttons">
      <view class="button-row">
        <wd-button
          type="primary"
          size="large"
          class="type-button"
          @click="quickRecord('pee')"
        >
          <view class="button-content">
            <text class="icon">💧</text>
            <text>小便</text>
          </view>
        </wd-button>

        <wd-button
          type="warning"
          size="large"
          class="type-button"
          @click="quickRecord('poop')"
        >
          <view class="button-content">
            <text class="icon">💩</text>
            <text>大便</text>
          </view>
        </wd-button>
      </view>

      <wd-button type="success" size="large" block @click="quickRecord('both')">
        <view class="button-content">
          <text class="icon">💧💩</text>
          <text>小便+大便</text>
        </view>
      </wd-button>
    </view>

    <!-- 详情区域 -->
    <wd-popup v-model="showDetails" position="bottom" :safe-area-inset-bottom="true">
      <view class="details-section">
        <view class="section-title">详细信息</view>

      <!-- 记录时间选择 -->
      <wd-cell-group>
        <wd-datetime-picker
          label="记录时间"
          v-model="recordDateTime"
          type="datetime"
          :min-date="minDateTime"
          :max-date="maxDateTime"
          @confirm="onDateTimeConfirm"
        />
      </wd-cell-group>

      <!-- 大便详情 (仅大便/两者时显示) -->
      <view v-if="form.type === 'poop' || form.type === 'both'" class="poop-details">
        <!-- 大便颜色 -->
        <view class="detail-item">
          <view class="detail-label">颜色</view>
          <view class="color-selector">
            <view
              v-for="color in poopColors"
              :key="color.value"
              class="color-item"
              :class="{ active: form.poopColor === color.value }"
              @click="form.poopColor = color.value"
            >
              <view
                class="color-circle"
                :style="{ background: color.color }"
              ></view>
              <text class="color-label">{{ color.label }}</text>
            </view>
          </view>
        </view>

        <!-- 大便性状 -->
        <view class="detail-item">
          <view class="detail-label">性状</view>
          <wd-radio-group v-model="form.poopTexture" shape="check" inline>
            <wd-radio
              v-for="texture in poopTextures"
              :key="texture.value"
              :value="texture.value"
            >
              {{ texture.label }}
            </wd-radio>
          </wd-radio-group>
        </view>
      </view>

      <!-- 备注 -->
      <wd-cell-group>

        <!-- 备注 -->
        <wd-cell title="备注">
          <wd-textarea
            v-model="form.note"
            placeholder="有什么需要记录的吗?"
            :max-length="200"
            :rows="2"
          />
        </wd-cell>
      </wd-cell-group>

      <!-- 提交按钮 -->
      <view class="submit-button">
        <wd-button type="primary" size="large" block @click="handleSubmit">
          {{ isEditing ? '更新记录' : '保存记录' }}
        </wd-button>
      </view>
      </view>
    </wd-popup>

    <!-- 时间选择器弹窗 (独立于表单外) -->
  </view>
</template>

<script setup lang="ts">
import { ref, computed } from "vue";
import { onLoad } from "@dcloudio/uni-app";
import { currentBabyId, getCurrentBaby } from "@/store/baby";
import { getUserInfo } from "@/store/user";
import type { DiaperType, PoopColor, PoopTexture } from "@/types";

// 直接调用 API 层
import * as diaperApi from "@/api/diaper";

// 编辑模式相关
const editId = ref<string>("");
const isEditing = computed(() => !!editId.value);

// 表单数据
const form = ref<{
  type: DiaperType;
  poopColor: PoopColor | undefined;
  poopTexture: PoopTexture | undefined;
  note: string;
  }>({
  type: "pee",
  poopColor: undefined,
  poopTexture: undefined,
  note: "",
});

// 是否显示详情
const showDetails = ref(false);

// 日期时间选择器
const recordDateTime = ref(new Date().getTime()); // 记录时间,初始为当前时间戳
const minDateTime = ref(
  new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).getTime()
); // 最小: 30天前
const maxDateTime = ref(new Date().getTime()); // 最大: 当前时间

// 确认日期时间选择
const onDateTimeConfirm = ({ value }: { value: number }) => {
  recordDateTime.value = value;
  console.log("[Diaper] 记录时间已更改为:", new Date(value));
};


// 格式化记录时间显示
const formatRecordTime = (timestamp: number): string => {
  const date = new Date(timestamp);
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const day = String(date.getDate()).padStart(2, "0");
  const hours = String(date.getHours()).padStart(2, "0");
  const minutes = String(date.getMinutes()).padStart(2, "0");
  return `${year}-${month}-${day} ${hours}:${minutes}`;
};

// 大便颜色选项
const poopColors = [
  { value: "yellow", label: "黄色", color: "#FFD700" },
  { value: "green", label: "绿色", color: "#90EE90" },
  { value: "brown", label: "棕色", color: "#8B4513" },
  { value: "black", label: "黑色", color: "#000000" },
  { value: "red", label: "红色", color: "#FF6347" },
  { value: "white", label: "白色", color: "#F0F0F0" },
] as const;

// 大便性状选项
const poopTextures = [
  { value: "watery", label: "稀水状" },
  { value: "loose", label: "稀软" },
  { value: "paste", label: "糊状" },
  { value: "soft", label: "软便" },
  { value: "formed", label: "成形" },
  { value: "hard", label: "硬结" },
] as const;

// 页面加载时检测 editId 参数
onLoad((options) => {
  if (options?.editId) {
    editId.value = options.editId;
    loadDiaperRecord(options.editId);
  }
});

// 加载尿布记录数据
const loadDiaperRecord = async (recordId: string) => {
  try {
    const record = await diaperApi.apiGetDiaperRecordById(recordId);

    // 填充表单
    form.value = {
      type: record.diaperType as DiaperType,
      poopColor: record.pooColor as PoopColor | undefined,
      poopTexture: record.pooTexture as PoopTexture | undefined,
      note: record.note || '',
    };

    // 设置记录时间
    recordDateTime.value = record.changeTime;

    // 打开详情弹窗
    showDetails.value = true;

    console.log('[Diaper] 已加载记录数据:', record);
  } catch (error: any) {
    console.error('[Diaper] 加载记录失败:', error);
    uni.showToast({
      title: error.message || '加载记录失败',
      icon: 'none',
    });
    setTimeout(() => {
      uni.navigateBack();
    }, 1500);
  }
};

// 快速记录
const quickRecord = (type: DiaperType) => {
  const currentBaby = getCurrentBaby();
  if (!currentBaby) {
    uni.showToast({
      title: "请先选择宝宝",
      icon: "none",
    });
    return;
  }

  // 设置记录类型
  form.value.type = type;

  // 重置时间为当前时间 (支持连续记录)
  recordDateTime.value = new Date().getTime();

  // 展开详情区域
  showDetails.value = true;
};

// 保存记录
const saveRecord = async (changeTime?: number) => {
  const user = getUserInfo();
  if (!user) {
    uni.showToast({
      title: "请先登录",
      icon: "none",
    });
    return;
  }

  try {
    // 使用传入的时间或当前表单中的时间
    const finalChangeTime = changeTime ?? recordDateTime.value;

    if (isEditing.value) {
      // 更新模式
      await diaperApi.apiUpdateDiaperRecord(editId.value, {
        babyId: currentBabyId.value,
        diaperType: form.value.type,
        pooColor: form.value.poopColor,
        pooTexture: form.value.poopTexture,
        note: form.value.note || undefined,
        changeTime: finalChangeTime,
      });

      uni.showToast({
        title: "更新成功",
        icon: "success",
      });
    } else {
      // 创建模式
      await diaperApi.apiCreateDiaperRecord({
        babyId: currentBabyId.value,
        diaperType: form.value.type,
        pooColor: form.value.poopColor,
        pooTexture: form.value.poopTexture,
        note: form.value.note || undefined,
        changeTime: finalChangeTime,
      });

      uni.showToast({
        title: "保存成功",
        icon: "success",
      });
    }

    setTimeout(() => {
      uni.navigateBack();
    }, 1000);
  } catch (error: any) {
    console.error("[Diaper] 保存换尿布记录失败:", error);
    uni.showToast({
      title: error.message || "保存失败",
      icon: "none",
    });
  }
};

// 提交记录
const handleSubmit = async () => {
  await saveRecord();
  // 提交成功后关闭弹窗
  showDetails.value = false;
};
</script>

<style lang="scss" scoped>
.diaper-page {
  min-height: 100vh;
  background: #f5f5f5;
  padding: 20rpx;
}

.quick-buttons {
  background: white;
  border-radius: 16rpx;
  padding: 30rpx;
  margin-bottom: 20rpx;
}

.button-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20rpx;
  margin-bottom: 20rpx;
}

.type-button {
  flex: 1;
}

.button-content {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8rpx;
  padding: 8rpx 0;

  .icon {
    font-size: 32rpx;
  }
}

.details-section {
  background: white;
  border-radius: 16rpx 16rpx 0 0;
  padding: 30rpx;
  max-height: 80vh;
  overflow-y: auto;

  .wd-cell-group {
    margin-bottom: 20rpx;
  }
}

.poop-details {
  margin: 20rpx 0;
  padding: 20rpx;
  // background: #f8f8f8;
  border-radius: 12rpx;
}

.detail-item {
  margin-bottom: 24rpx;

  &:last-child {
    margin-bottom: 0;
  }
}

.detail-label {
  font-size: 28rpx;
  font-weight: 500;
  color: #333;
  margin-bottom: 16rpx;
}

.section-title {
  font-size: 32rpx;
  font-weight: bold;
  margin-bottom: 24rpx;
}

.time-value {
  color: #fa2c19;
  font-weight: 500;
  font-size: 28rpx;
}

.time-display {
  display: flex;
  align-items: center;
  gap: 10rpx;
  color: #333;
  font-size: 28rpx;
}

.color-selector {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 12rpx;
}

.color-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6rpx;
  padding: 8rpx;
  border-radius: 8rpx;
  border: 2rpx solid transparent;
  transition: all 0.3s;

  &.active {
    border-color: #fa2c19;
    background: rgba(250, 44, 25, 0.05);
  }
}

.color-circle {
  width: 48rpx;
  height: 48rpx;
  border-radius: 50%;
  border: 2rpx solid #ddd;
}

.color-label {
  font-size: 20rpx;
  color: #666;
}

.texture-list {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12rpx;
}

.submit-button {
  margin-top: 40rpx;
}
</style>
