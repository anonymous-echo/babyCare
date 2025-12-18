<template>
  <view>
    <wd-navbar
      title="邀请亲友"
      left-arrow
      safeAreaInsetTop
      @click-left="handleBack"
    >
      <template #capsule>
        <wd-navbar-capsule @back="handleBack" @back-home="handleBackHome" />
      </template>
    </wd-navbar>
    <view class="invite-page">
      <!-- 顶部信息卡片 -->
      <view class="header-card">
        <view class="baby-info">
          <view class="baby-avatar">
            <!-- 宝宝头像 -->
            <image
              v-if="babyAvatarUrl"
              :src="babyAvatarUrl"
              mode="aspectFill"
            />
            <!-- 默认头像 -->
            <image
              v-else
              src="@/static/default.png"
              mode="aspectFill"
            />
          </view>
          <view class="baby-detail">
            <text class="baby-name">{{ babyName }}</text>
            <text class="baby-desc">邀请亲友共同记录成长</text>
          </view>
        </view>
      </view>

      <!-- 设置表单 -->
      <wd-cell-group>
        <wd-cell title="与宝宝的关系" is-link @click="showRelationshipPicker = true">
          <text>{{ selectedRelationship || '请选择' }}</text>
        </wd-cell>

        <wd-cell title="亲友角色">
          <wd-radio-group v-model="selectedRole" inline>
            <wd-radio value="editor">编辑者</wd-radio>
            <wd-radio value="viewer">查看者</wd-radio>
          </wd-radio-group>
        </wd-cell>

        <wd-cell title="访问权限">
          <wd-radio-group v-model="accessType" inline>
            <wd-radio value="permanent">永久</wd-radio>
            <wd-radio value="temporary">临时</wd-radio>
          </wd-radio-group>
        </wd-cell>

        <wd-datetime-picker
          v-if="accessType === 'temporary'"
          label="到期时间"
          v-model="expiresDateValue"
          type="datetime"
          :min-date="minDate"
          :max-date="maxDate"
          @confirm="onDateTimeConfirm"
          @cancel="showDatetimePickerModal = false"
        />
      </wd-cell-group>

      <!-- 角色说明 -->
      <view class="role-tips">
        <text class="tip-icon">ℹ️</text>
        <text class="tip-text" v-if="selectedRole === 'editor'">
          编辑者可以记录和编辑所有数据
        </text>
        <text class="tip-text" v-else> 查看者只能查看数据，不能编辑 </text>
      </view>

      <!-- 生成按钮 -->
      <view class="button-wrapper">
        <wd-button
          type="primary"
          size="large"
          block
          @click="handleGenerateQRCode"
          :loading="generating"
        >
          {{ generating ? "生成中..." : "生成邀请二维码" }}
        </wd-button>
      </view>

      <!-- 二维码展示 -->
      <view v-if="qrcodeUrl" class="qrcode-section">
        <wd-card>
          <view class="qrcode-wrapper">
            <image :src="qrcodeUrl" class="qrcode-image" mode="aspectFit" />
          </view>

          <view class="qrcode-footer">
            <text class="footer-text">长按识别二维码或保存到相册</text>
            <wd-button type="success" size="small" @click="saveQRCode">
              保存到相册
            </wd-button>
          </view>
        </wd-card>
      </view>
    </view>

    <!-- 关系选择弹窗 -->
    <wd-popup
      v-model="showRelationshipPicker"
      position="bottom"
      custom-style="height: auto; padding: 0"
      safe-area-inset-bottom
    >
      <view class="relationship-popup">
        <view class="popup-header">
          <text class="popup-title">选择与宝宝的关系</text>
          <wd-icon name="close" @click="showRelationshipPicker = false" />
        </view>

        <!-- 自定义输入 -->
        <view class="custom-input-section">
          <wd-input
            v-model="customRelationship"
            placeholder="或输入自定义关系"
            clearable
          />
        </view>

        <!-- 预设选项 -->
        <view class="preset-options">
          <view
            v-for="option in relationshipOptions"
            :key="option.value"
            class="option-item"
            :class="{ active: selectedRelationship === option.value }"
            @click="selectRelationship(option.value)"
          >
            <text>{{ option.label }}</text>
          </view>
        </view>

        <!-- 确认按钮 -->
        <view class="popup-footer">
          <wd-button
            type="primary"
            size="large"
            block
            @click="confirmRelationship"
          >
            确认
          </wd-button>
        </view>
      </view>
    </wd-popup>
  </view>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { onLoad } from "@dcloudio/uni-app";
import { inviteCollaborator } from "@/store/collaborator";
import * as babyApi from "@/api/baby";

// 页面参数
const babyId = ref("");
const babyName = ref("");
const babyAvatarUrl = ref("");

// 表单数据
const selectedRole = ref<"editor" | "viewer">("editor");
const selectedRelationship = ref("");
const customRelationship = ref("");
const accessType = ref<"permanent" | "temporary">("permanent");
const expiresDateValue = ref<number>(Date.now() + 7 * 24 * 60 * 60 * 1000); // 默认7天后
const showDatetimePickerModal = ref(false);
const showRelationshipPicker = ref(false);

// 关系选项
const relationshipOptions = [
  { label: '爸爸', value: '爸爸' },
  { label: '妈妈', value: '妈妈' },
  { label: '爷爷', value: '爷爷' },
  { label: '奶奶', value: '奶奶' },
  { label: '外公', value: '外公' },
  { label: '外婆', value: '外婆' },
  { label: '叔叔', value: '叔叔' },
  { label: '姑姑', value: '姑姑' },
  { label: '舅舅', value: '舅舅' },
  { label: '姨妈', value: '姨妈' },
  { label: '其他亲友', value: '其他亲友' },
];

// 二维码相关
const qrcodeUrl = ref("");
const generating = ref(false);

// 日期选择器范围
const minDate = Date.now();
const maxDate = Date.now() + 365 * 24 * 60 * 60 * 1000;

// 页面加载
onLoad((options) => {
  if (options?.babyId) {
    babyId.value = options.babyId;
  }
  if (options?.babyName) {
    babyName.value = decodeURIComponent(options.babyName);
  }

  // 获取宝宝详情（包括头像）
  if (babyId.value) {
    babyApi.apiFetchBabyDetail(babyId.value)
      .then((baby) => {
        if (baby?.avatarUrl) {
          babyAvatarUrl.value = baby.avatarUrl;
        }
      })
      .catch((error) => {
        console.warn("[Invite] 获取宝宝头像失败:", error);
        // 头像加载失败不影响邀请流程
      });
  }
});

// 日期时间选择确认
function onDateTimeConfirm({ value }: { value: number }) {
  expiresDateValue.value = value;
  showDatetimePickerModal.value = false;
}

// 生成二维码
async function handleGenerateQRCode() {
  if (!babyId.value) {
    uni.showToast({
      title: "宝宝ID不能为空",
      icon: "none",
    });
    return;
  }

  generating.value = true;

  try {
    // 计算过期时间戳
    const expiresAt =
      accessType.value === "temporary" ? expiresDateValue.value : undefined;

    // 调用API生成邀请（二维码方式）
    const invitationData = await inviteCollaborator(
      babyId.value,
      "qrcode",
      selectedRole.value,
      accessType.value,
      expiresAt,
      selectedRelationship.value || undefined,
    );

    const { qrcodeParams } = invitationData;

    if (!qrcodeParams || !qrcodeParams.qrcodeUrl) {
      uni.showToast({
        title: "二维码生成失败",
        icon: "none",
      });
      return;
    }

    // 显示二维码
    qrcodeUrl.value = qrcodeParams.qrcodeUrl;

    uni.showToast({
      title: "二维码生成成功",
      icon: "success",
    });
  } catch (error: any) {
    console.error("Generate QR code error:", error);
    uni.showToast({
      title: error.message || "生成失败",
      icon: "none",
    });
  } finally {
    generating.value = false;
  }
}

// 保存二维码
function saveQRCode() {
  if (!qrcodeUrl.value) {
    uni.showToast({
      title: "二维码未生成",
      icon: "none",
    });
    return;
  }

  // 下载二维码图片
  uni.downloadFile({
    url: qrcodeUrl.value,
    success: (res) => {
      if (res.statusCode === 200) {
        uni.saveImageToPhotosAlbum({
          filePath: res.tempFilePath,
          success: () => {
            uni.showToast({
              title: "保存成功",
              icon: "success",
            });
          },
          fail: () => {
            uni.showToast({
              title: "保存失败,请授予相册权限",
              icon: "none",
            });
          },
        });
      }
    },
    fail: (err) => {
      console.error("Download QR code error:", err);
      uni.showToast({
        title: "下载失败",
        icon: "none",
      });
    },
  });
}
function handleBackHome() {
  uni.switchTab({
    url: "/pages/index/index",
  });
}
function handleBack() {
  uni.navigateBack();
}

// 选择预设关系
function selectRelationship(value: string) {
  selectedRelationship.value = value;
  customRelationship.value = ""; // 清空自定义输入
}

// 确认关系选择
function confirmRelationship() {
  // 优先使用自定义输入，否则使用预设选项
  if (customRelationship.value.trim()) {
    selectedRelationship.value = customRelationship.value.trim();
  }
  showRelationshipPicker.value = false;
}
</script>

<style lang="scss" scoped>
@import '@/styles/colors.scss';

// ===== 页面布局 =====
.invite-page {
  min-height: 100vh;
  background: $gradient-bg-light;
  padding: $spacing-lg;
  padding-bottom: 120rpx; // 为按钮预留空间
}

// ===== 顶部卡片 =====
.header-card {
  background: $color-bg-primary;
  border: 2rpx solid $color-border-primary;
  border-radius: $radius-lg;
  padding: $spacing-2xl;
  margin-bottom: $spacing-2xl;
  box-shadow: $shadow-md;
  overflow: hidden;

  .baby-info {
    display: flex;
    align-items: center;
    gap: $spacing-2xl;

    .baby-avatar {
      width: 100rpx;
      height: 100rpx;
      border-radius: $radius-full;
      overflow: hidden;
      flex-shrink: 0;
      box-shadow: $shadow-md;
      border: 2rpx solid $color-border-primary;
      background: $color-bg-secondary;

      image {
        width: 100%;
        height: 100%;
        object-fit: cover;
      }
    }

    .baby-detail {
      flex: 1;
      display: flex;
      flex-direction: column;
      gap: $spacing-sm;

      .baby-name {
        font-size: $font-size-2xl;
        font-weight: $font-weight-bold;
        color: $color-text-primary;
        line-height: 1.3;
      }

      .baby-desc {
        font-size: $font-size-md;
        color: $color-text-secondary;
        line-height: 1.4;
      }
    }
  }
}

// ===== 表单分组 =====
:deep(.wd-cell-group) {
  background: $color-bg-primary;
  border: 1rpx solid $color-border-primary;
  border-radius: $radius-lg;
  margin-bottom: $spacing-2xl;
  overflow: hidden;
  box-shadow: $shadow-sm;
}

:deep(.wd-cell) {
  padding: $spacing-lg $spacing-md;
  background: $color-bg-primary;
  border-bottom: 1rpx solid $color-border-primary;
  transition: background $transition-base;

  &:last-child {
    border-bottom: none;
  }

  &:active {
    background: $color-bg-secondary;
  }
}

:deep(.wd-cell__title) {
  font-size: $font-size-base;
  color: $color-text-primary;
  font-weight: $font-weight-medium;
}

:deep(.wd-cell__value) {
  font-size: $font-size-base;
  color: $color-text-secondary;
}

// ===== 单选框 =====
:deep(.wd-radio-group) {
  display: flex;
  gap: $spacing-xl;
  flex-wrap: wrap;
}

:deep(.wd-radio) {
  font-size: $font-size-base;
  color: $color-text-primary;

  &.is-checked {
    color: $color-primary;
  }
}

// ===== 日期选择器 =====
:deep(.wd-datetime-picker) {
  padding: $spacing-lg $spacing-md;
  background: $color-bg-primary;
  border-bottom: 1rpx solid $color-border-primary;

  &:last-child {
    border-bottom: none;
  }
}

:deep(.wd-datetime-picker__label) {
  font-size: $font-size-base;
  color: $color-text-primary;
  font-weight: $font-weight-medium;
}

// ===== 提示框 =====
.role-tips {
  display: flex;
  align-items: flex-start;
  gap: $spacing-md;
  padding: $spacing-lg;
  background: linear-gradient(135deg, rgba(50, 220, 110, 0.08) 0%, rgba(50, 220, 110, 0.04) 100%);
  border: 1rpx solid $color-border-primary;
  border-left: 4rpx solid $color-primary;
  border-radius: $radius-md;
  margin-top: $spacing-2xl;
  margin-bottom: $spacing-2xl;

  .tip-icon {
    font-size: $font-size-lg;
    flex-shrink: 0;
    line-height: 1.4;
  }

  .tip-text {
    flex: 1;
    font-size: $font-size-sm;
    color: $color-text-primary;
    line-height: 1.6;
    font-weight: $font-weight-medium;
  }
}

// ===== 按钮容器 =====
.button-wrapper {
  margin-top: $spacing-3xl;
  margin-bottom: $spacing-3xl;

  :deep(.wd-button) {
    height: 88rpx;
    font-size: $font-size-lg;
    font-weight: $font-weight-medium;
    border-radius: $radius-md;
    background: $color-primary;
    color: white;
    box-shadow: $shadow-primary-md;
    transition: all $transition-base;

    &:active {
      background: darken($color-primary, 10%);
      transform: scale(0.98);
    }
  }
}

// ===== 二维码区域 =====
.qrcode-section {
  margin-top: $spacing-2xl;
  animation: fadeIn 0.3s ease;

  :deep(.wd-card) {
    background: $color-bg-primary;
    border: 1rpx solid $color-border-primary;
    border-radius: $radius-lg;
    box-shadow: $shadow-md;
  }
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(20rpx);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.qrcode-wrapper {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: $spacing-3xl $spacing-lg;
  background: $color-bg-secondary;

  .qrcode-image {
    width: 400rpx;
    height: 400rpx;
    border-radius: $radius-md;
    background: $color-bg-primary;
    box-shadow: $shadow-md;
    border: 1rpx solid $color-border-primary;
  }
}

.qrcode-footer {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: $spacing-lg;
  padding: $spacing-2xl;
  border-top: 1rpx solid $color-border-primary;
  background: $color-bg-primary;

  .footer-text {
    font-size: $font-size-sm;
    color: $color-text-secondary;
    text-align: center;
  }

  :deep(.wd-button) {
    min-width: 160rpx;
    height: 64rpx;
    font-size: $font-size-base;
    border-radius: $radius-md;
  }
}

// ===== 关系选择弹窗 =====
.relationship-popup {
  background: $color-bg-primary;
  border-radius: $radius-lg $radius-lg 0 0;
  overflow: hidden;

  .popup-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: $spacing-lg $spacing-2xl;
    border-bottom: 1rpx solid $color-border-primary;

    .popup-title {
      font-size: $font-size-lg;
      font-weight: $font-weight-semibold;
      color: $color-text-primary;
    }

    :deep(.wd-icon) {
      font-size: 40rpx;
      color: $color-text-secondary;
      cursor: pointer;
    }
  }

  .custom-input-section {
    padding: $spacing-2xl;
    border-bottom: 1rpx solid $color-border-primary;

    :deep(.wd-input) {
      background: $color-bg-secondary;
      border-radius: $radius-md;
      padding: $spacing-md;
    }
  }

  .preset-options {
    padding: $spacing-lg;
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: $spacing-md;
    max-height: 400rpx;
    overflow-y: auto;

    .option-item {
      padding: $spacing-lg;
      background: $color-bg-secondary;
      border: 2rpx solid $color-border-primary;
      border-radius: $radius-md;
      text-align: center;
      font-size: $font-size-base;
      color: $color-text-primary;
      transition: all $transition-base;
      cursor: pointer;

      &:active {
        transform: scale(0.95);
      }

      &.active {
        background: $color-primary-lighter;
        border-color: $color-primary;
        color: $color-primary;
        font-weight: $font-weight-semibold;
      }
    }
  }

  .popup-footer {
    padding: $spacing-lg $spacing-2xl;
    border-top: 1rpx solid $color-border-primary;

    :deep(.wd-button) {
      height: 88rpx;
      font-size: $font-size-lg;
      border-radius: $radius-md;
    }
  }
}
</style>
