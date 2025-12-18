<template>
  <view>
    <wd-navbar
      title="成长记录"
      left-text="返回"
      left-arrow
      safeAreaInsetTop
      placeholder
      fixed
    >
      <template #capsule>
        <wd-navbar-capsule @back="goBack" @back-home="goBackHome" />
      </template>
    </wd-navbar>
    <view class="growth-container">
      <!-- 页面内容 -->
      <view class="growth-content">
        <!-- 最新数据卡片 -->
        <view v-if="latestRecord" class="latest-card">
          <view class="card-title">最新记录</view>
          <view class="data-grid">
            <view v-if="latestRecord.height" class="data-item">
              <view class="data-icon">📏</view>
              <view class="data-value">{{ latestRecord.height }}</view>
              <view class="data-label">身高(cm)</view>
            </view>
            <view v-if="latestRecord.weight" class="data-item">
              <view class="data-icon">⚖️</view>
              <view class="data-value">{{ latestRecord.weight }}</view>
              <view class="data-label">体重(g)</view>
            </view>
            <view v-if="latestRecord.headCircumference" class="data-item">
              <view class="data-icon">📐</view>
              <view class="data-value">{{
                latestRecord.headCircumference
              }}</view>
              <view class="data-label">头围(cm)</view>
            </view>
          </view>
          <view class="record-time">
            记录于
            {{ formatDate(latestRecord.measureTime, "YYYY-MM-DD HH:mm") }}
          </view>
        </view>

        <!-- 添加记录按钮 -->
        <view class="add-section">
          <wd-button
            type="primary"
            size="large"
            block
            @click="showAddDialog = true"
          >
            + 添加成长记录
          </wd-button>
        </view>

        <!-- 历史记录列表 -->
        <view class="records-section">
          <view class="section-title">历史记录</view>

          <view v-if="recordList.length === 0" class="empty-state">
            <wd-status-tip description="暂无成长记录" />
          </view>

          <view v-else class="record-list">
            <view
              v-for="record in recordList"
              :key="record.recordId"
              class="record-item"
            >
              <view class="record-header">
                <view class="record-date">
                  {{ formatDate(record.measureTime, "YYYY-MM-DD") }}
                </view>
                <wd-button
                  size="small"
                  type="info"
                  @click="handleDelete(record.recordId)"
                >
                  删除
                </wd-button>
              </view>

              <view class="record-data">
                <view v-if="record.height" class="data-row">
                  <text class="data-label">身高:</text>
                  <text class="data-value">{{ record.height }} cm</text>
                </view>
                <view v-if="record.weight" class="data-row">
                  <text class="data-label">体重:</text>
                  <text class="data-value">{{ record.weight }} g</text>
                </view>
                <view v-if="record.headCircumference" class="data-row">
                  <text class="data-label">头围:</text>
                  <text class="data-value"
                    >{{ record.headCircumference }} cm</text
                  >
                </view>
                <view v-if="record.note" class="data-row">
                  <text class="data-label">备注:</text>
                  <text class="data-value">{{ record.note }}</text>
                </view>
              </view>
            </view>
          </view>
        </view>
      </view>

      <!-- 添加记录对话框 -->
      <wd-popup
        v-model="showAddDialog"
        custom-style="height: 60%"
        position="bottom"
        round
        closeable
      >
        <view class="dialog-content">
          <view class="dialog-title">{{
            isEditing ? "编辑成长记录" : "添加成长记录"
          }}</view>
          <wd-form ref="form" :model="formData">
            <wd-cell-group border>
              <wd-input
                v-model="formData.height"
                placeholder="请输入身高..."
                label="身高"
                required
                type="number"
                inputmode="numeric"
                ><template #suffix>厘米</template></wd-input
              >
              <wd-input
                v-model="formData.weight"
                placeholder="请输入体重..."
                label="体重"
                required
                type="number"
                inputmode="numeric"
                ><template #suffix>克</template></wd-input
              >
              <wd-input
                v-model="formData.headCircumference"
                placeholder="请输入头围..."
                label="头围"
                required
                type="number"
                inputmode="numeric"
                ><template #suffix>厘米</template></wd-input
              >
              <!-- 备注 -->
              <wd-textarea
                v-model="formData.note"
                placeholder="请输入备注..."
                label="备注"
                rows="3"
              ></wd-textarea>
            </wd-cell-group>
            <view class="dialog-footer">
              <wd-button
                type="info"
                size="large"
                block
                @click="showAddDialog = false"
              >
                取消
              </wd-button>
              <wd-button
                type="primary"
                size="large"
                block
                @click="handleSubmit"
              >
                {{ isEditing ? "更新" : "保存" }}
              </wd-button>
            </view>
          </wd-form>
        </view>
      </wd-popup>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { onLoad } from "@dcloudio/uni-app";
import { currentBaby } from "@/store/baby";
import { formatDate } from "@/utils/date";

// 直接调用 API 层
import * as growthApi from "@/api/growth";
import { goBack, goBackHome } from "@/utils/common";

// 编辑模式相关
const editId = ref<string>("");
const isEditing = computed(() => !!editId.value);

// 对话框显示状态
const showAddDialog = ref(false);
const showDatePicker = ref(false);
const selectedDate = ref(new Date());

// 表单数据
const formData = ref({
  height: "",
  weight: "",
  headCircumference: "",
  time: Date.now(),
  note: "",
});

// 成长记录列表(从 API 获取)
const records = ref<growthApi.GrowthRecordResponse[]>([]);

// 最新记录
const latestRecord = computed(() => {
  return records.value.length > 0 ? records.value[0] : null;
});

// 历史记录列表
const recordList = computed(() => {
  return records.value;
});

// 加载成长记录
const loadRecords = async () => {
  if (!currentBaby.value) return;

  try {
    const data = await growthApi.apiFetchGrowthRecords({
      babyId: currentBaby.value.babyId,
      pageSize: 100,
    });
    console.log("成长记录 API 响应:", data); // 调试日志
    records.value = data.records;
  } catch (error) {
    console.error("加载成长记录失败:", error);
  }
};

// 页面加载时检测 editId 参数
onLoad((options) => {
  if (options?.editId) {
    editId.value = options.editId;
    loadGrowthRecord(options.editId);
  }
});

// 加载成长记录数据
const loadGrowthRecord = async (recordId: string) => {
  try {
    const record = await growthApi.apiGetGrowthRecordById(recordId);

    // 填充表单
    formData.value = {
      height: record.height ? String(record.height) : "",
      weight: record.weight ? String(record.weight) : "",
      headCircumference: record.headCircumference
        ? String(record.headCircumference)
        : "",
      time: record.measureTime,
      note: record.note || "",
    };

    // 打开对话框
    showAddDialog.value = true;

    console.log("[Growth] 已加载记录数据:", record);
  } catch (error: any) {
    console.error("[Growth] 加载记录失败:", error);
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
  loadRecords();
});

// 日期选择确认
const handleDateConfirm = ({ value }: any) => {
  console.log("选择的日期:", value);
  formData.value.time = value;
};

// 提交表单
const handleSubmit = async () => {
  if (!currentBaby.value) {
    uni.showToast({
      title: "请先选择宝宝",
      icon: "none",
    });
    return;
  }

  // 验证至少填写一项
  if (
    !formData.value.height &&
    !formData.value.weight &&
    !formData.value.headCircumference
  ) {
    uni.showToast({
      title: "请至少填写一项数据",
      icon: "none",
    });
    return;
  }

  // 验证数据范围
  const height = parseFloat(formData.value.height);
  const weight = parseFloat(formData.value.weight);
  const headCircumference = parseFloat(formData.value.headCircumference);

  if (formData.value.height && (isNaN(height) || height <= 0 || height > 200)) {
    uni.showToast({
      title: "身高数据不合理",
      icon: "none",
    });
    return;
  }

  if (
    formData.value.weight &&
    (isNaN(weight) || weight <= 0 || weight > 200000)
  ) {
    uni.showToast({
      title: "体重数据不合理",
      icon: "none",
    });
    return;
  }

  if (
    formData.value.headCircumference &&
    (isNaN(headCircumference) ||
      headCircumference <= 0 ||
      headCircumference > 100)
  ) {
    uni.showToast({
      title: "头围数据不合理",
      icon: "none",
    });
    return;
  }

  // 添加或更新记录
  try {
    if (isEditing.value) {
      // 更新模式
      await growthApi.apiUpdateGrowthRecord(editId.value, {
        babyId: currentBaby.value.babyId,
        measureTime: formData.value.time,
        height: formData.value.height ? height : undefined,
        weight: formData.value.weight ? weight : undefined,
        headCircumference: formData.value.headCircumference
          ? headCircumference
          : undefined,
        note: formData.value.note || undefined,
      });

      uni.showToast({
        title: "更新成功",
        icon: "success",
      });
    } else {
      // 创建模式
      await growthApi.apiCreateGrowthRecord({
        babyId: currentBaby.value.babyId,
        measureTime: formData.value.time,
        height: formData.value.height ? height : undefined,
        weight: formData.value.weight ? weight : undefined,
        headCircumference: formData.value.headCircumference
          ? headCircumference
          : undefined,
        note: formData.value.note || undefined,
      });

      uni.showToast({
        title: "添加成功",
        icon: "success",
      });
    }

    // 重新加载记录
    await loadRecords();

    // 重置表单
    formData.value = {
      height: "",
      weight: "",
      headCircumference: "",
      time: Date.now(),
      note: "",
    };

    showAddDialog.value = false;
  } catch (error: any) {
    uni.showToast({
      title: error.message || "保存失败",
      icon: "none",
    });
  }
};

// 删除记录
const handleDelete = async (id: string) => {
  uni.showModal({
    title: "确认删除",
    content: "确定要删除这条成长记录吗?",
    success: async (res) => {
      if (res.confirm) {
        try {
          await growthApi.apiDeleteGrowthRecord(id);
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
</script>

<style lang="scss" scoped>
.growth-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.growth-content {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  background: #f6f8f7;
  padding: 20rpx;
}

.latest-card {
  background: linear-gradient(135deg, #7dd3a2 0%, #52c41a 100%);
  border-radius: 16rpx;
  padding: 30rpx;
  margin-bottom: 20rpx;
  color: white;
  flex-shrink: 0;
}

.card-title {
  font-size: 32rpx;
  font-weight: bold;
  margin-bottom: 20rpx;
}

.data-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20rpx;
  margin-bottom: 20rpx;
}

.data-item {
  text-align: center;
}

.data-icon {
  font-size: 40rpx;
  margin-bottom: 8rpx;
}

.data-value {
  font-size: 36rpx;
  font-weight: bold;
  margin-bottom: 4rpx;
}

.data-label {
  font-size: 24rpx;
  opacity: 0.9;
}

.record-time {
  font-size: 24rpx;
  opacity: 0.8;
  text-align: center;
}

.add-section {
  margin-bottom: 20rpx;
  flex-shrink: 0;
}

.records-section {
  background: white;
  border: 1rpx solid #cae3d4;
  border-radius: 16rpx;
  padding: 30rpx;
  flex: 1;
  overflow-y: auto;
  box-shadow: 0 2rpx 8rpx rgba(125, 211, 162, 0.08);
}

.section-title {
  font-size: 32rpx;
  font-weight: bold;
  margin-bottom: 20rpx;
}

.empty-state {
  padding: 80rpx 0;
}

.record-list {
  display: flex;
  flex-direction: column;
  gap: 20rpx;
}

.record-item {
  background: white;
  border: 1rpx solid #cae3d4;
  border-radius: 12rpx;
  padding: 24rpx;
  box-shadow: 0 2rpx 8rpx rgba(125, 211, 162, 0.08);
}

.record-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16rpx;
}

.record-date {
  font-size: 28rpx;
  font-weight: bold;
  color: #1a1a1a;
}

.record-data {
  display: flex;
  flex-direction: column;
  gap: 12rpx;
}

.data-row {
  display: flex;
  justify-content: space-between;
  font-size: 26rpx;

  .data-label {
    color: #666;
  }

  .data-value {
    color: #1a1a1a;
    font-weight: 500;
  }
}

.dialog-content {
  padding: 24rpx;
  display: flex;
  flex-direction: column;
  max-height: 100%;
}

.dialog-title {
  font-size: 36rpx;
  font-weight: bold;
  text-align: center;
  margin-bottom: 24rpx;
}

.form-section {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  min-height: 0;
}

.form-item {
  margin-bottom: 0;
  flex: 0 1 auto;
}

.form-label {
  display: flex;
  align-items: center;
  gap: 8rpx;
  font-size: 26rpx;
  font-weight: bold;
  margin-bottom: 8rpx;

  .icon {
    font-size: 28rpx;
  }
}

.dialog-footer {
  display: flex;
  flex-direction: column;
  gap: 16rpx;
  margin-top: 16rpx;
}
</style>
