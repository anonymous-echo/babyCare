<template>
  <view class="baby-edit-container">
    <wd-navbar
      :title="isEdit ? '编辑资料' : '添加宝宝'"
      left-text="返回"
      left-arrow
      fixed
      placeholder
      safe-area-inset-top
      @click-left="handleCancel"
    />

    <scroll-view class="edit-scroll" scroll-y>
      <view class="edit-content">
        <!-- 头像大卡片 (Premium Style) -->
        <view class="avatar-hero-card premium-shadow">
          <view class="avatar-upload-box" @click="chooseAvatar">
            <image
              :src="resolveImageUrl(formData.avatarUrl) || '/static/default.png'"
              mode="aspectFill"
              class="avatar-img"
            />
            <view class="camera-fab">
              <wd-icon name="camera" size="18" color="#FFF" />
            </view>
          </view>
          <text class="upload-hint">点击设置宝宝头像</text>
        </view>

        <!-- 表单卡片 -->
        <view class="form-card premium-shadow">
          <view class="form-item-premium">
            <text class="label">姓名</text>
            <wd-input
              v-model="formData.name"
              placeholder="请输入真实姓名"
              no-border
            />
          </view>

          <view class="form-divider"></view>

          <view class="form-item-premium">
            <text class="label">昵称 / 小名</text>
            <wd-input
              v-model="formData.nickname"
              placeholder="如何称呼宝宝？"
              no-border
            />
          </view>

          <view class="form-divider"></view>

          <view class="form-item-premium">
            <text class="label">性别</text>
            <view class="gender-picker-row">
              <view
                class="gender-option male"
                :class="{ active: formData.gender === 'male' }"
                @click="formData.gender = 'male'"
              >
                <text class="icon">👦</text>
                <text class="txt">小王子</text>
              </view>
              <view
                class="gender-option female"
                :class="{ active: formData.gender === 'female' }"
                @click="formData.gender = 'female'"
              >
                <text class="icon">👧</text>
                <text class="txt">小公主</text>
              </view>
            </view>
          </view>

          <view class="form-divider"></view>

          <view class="form-item-premium">
            <text class="label">出生日期</text>
            <wd-datetime-picker
              v-model="selectedDateTimestamp"
              type="date"
              no-border
              placeholder="宝宝哪天出生的？"
              @confirm="handleDateConfirm"
            />
          </view>
        </view>
      </view>
    </scroll-view>

    <!-- 悬浮提交按钮 -->
    <view class="safe-bottom-dock">
      <wd-button
        block
        round
        type="primary"
        size="large"
        :loading="isSubmitting"
        @click="handleSubmit"
      >
        {{ isEdit ? "确认并保存" : "立即创建" }}
      </wd-button>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from "vue";
import { formatDate } from "@/utils/date";
import { uploadFile } from "@/utils/request";
import { resolveImageUrl } from "@/utils/assets";
import * as babyApi from "@/api/baby";
import * as vaccineApi from "@/api/vaccine";
import { fetchBabyDetail } from "@/store/baby";

const formData = ref({
  name: "",
  nickname: "",
  gender: "male" as "male" | "female",
  birthDate: "",
  avatarUrl: "",
});

const isEdit = ref(false);
const editId = ref("");
const isSubmitting = ref(false);
const selectedDateTimestamp = ref<number>(Date.now());

onMounted(async () => {
  const pages = getCurrentPages();
  const options = (pages[pages.length - 1] as any).options || {};
  if (options.id) {
    isEdit.value = true;
    editId.value = options.id;
    try {
      const baby = await babyApi.apiFetchBabyDetail(options.id);
      if (baby) {
        formData.value = {
          name: baby.name,
          nickname: baby.nickname || "",
          gender: baby.gender as "male" | "female",
          birthDate: baby.birthDate,
          avatarUrl: baby.avatarUrl || "",
        };
        selectedDateTimestamp.value = new Date(baby.birthDate).getTime();
      }
    } catch (e) {}
  }
});

const chooseAvatar = () => {
  uni.chooseImage({
    count: 1,
    success: async (res) => {
      const tempPath = res.tempFilePaths[0];
      try {
        if (!tempPath) return;
        uni.showLoading({ title: "上传中..." });
        const result: any = await uploadFile({
          filePath: tempPath,
          name: "file",
          formData: { type: "baby_avatar" },
        });
        if (result.code === 0 && result.data?.url) {
          formData.value.avatarUrl = result.data.url as string;
        }
      } catch (e) {
      } finally {
        uni.hideLoading();
      }
    },
  });
};

const handleDateConfirm = ({ value }: { value: number }) => {
  formData.value.birthDate = formatDate(value, "YYYY-MM-DD");
};

const handleSubmit = async () => {
  if (!formData.value.name.trim() || !formData.value.birthDate) {
    uni.showToast({ title: "请完善姓名和生日", icon: "none" });
    return;
  }
  try {
    isSubmitting.value = true;
    if (isEdit.value) {
      await babyApi.apiUpdateBaby(editId.value, formData.value);
      await fetchBabyDetail(editId.value);
      uni.showToast({ title: "已保存", icon: "success" });
      setTimeout(() => uni.navigateBack(), 1000);
    } else {
      const newBaby = await babyApi.apiCreateBaby(formData.value);
      try {
        await vaccineApi.apiFetchVaccinePlans(newBaby.babyId);
      } catch (e) {}
      uni.showToast({ title: "创建成功", icon: "success" });
      setTimeout(() => uni.reLaunch({ url: "/pages/index/index" }), 1000);
    }
  } catch (e) {
  } finally {
    isSubmitting.value = false;
  }
};

const handleCancel = () => uni.navigateBack();
</script>

<style lang="scss" scoped>
@import "@/styles/colors.scss";

.baby-edit-container {
  min-height: 100vh;
  background: $color-bg-secondary;
}

.edit-scroll {
  height: calc(100vh - 160rpx);
}

.edit-content {
  padding: 48rpx 32rpx;
}

.avatar-hero-card {
  background: #fff;
  border-radius: $radius-lg;
  padding: 60rpx 40rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 40rpx;

  .avatar-upload-box {
    position: relative;
    width: 220rpx;
    height: 220rpx;
    margin-bottom: 24rpx;
    .avatar-img {
      width: 100%;
      height: 100%;
      border-radius: 50%;
      border: 8rpx solid $color-primary-lighter;
    }
    .camera-fab {
      position: absolute;
      bottom: 0;
      right: 0;
      width: 72rpx;
      height: 72rpx;
      background: $color-primary;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      border: 4rpx solid #fff;
      box-shadow: 0 4rpx 12rpx rgba(123, 211, 162, 0.4);
    }
  }
  .upload-hint {
    font-size: 24rpx;
    font-weight: 600;
    color: $color-text-tertiary;
  }
}

.form-card {
  background: #fff;
  border-radius: $radius-lg;
  padding: 24rpx 40rpx;
}

.form-item-premium {
  padding: 32rpx 0;
  .label {
    font-size: 24rpx;
    font-weight: 800;
    color: $color-text-secondary;
    margin-bottom: 12rpx;
    display: block;
  }
  :deep(.wd-input),
  :deep(.wd-datetime-picker) {
    background: transparent;
    padding: 0;
    .wd-input__control {
      font-size: 32rpx;
      font-weight: 600;
      color: $color-text-primary;
    }
  }
}

.form-divider {
  height: 1rpx;
  background: $color-divider;
  margin: 0;
}

.gender-picker-row {
  display: flex;
  gap: 24rpx;
  margin-top: 16rpx;

  .gender-option {
    flex: 1;
    height: 100rpx;
    background: $color-bg-secondary;
    border-radius: $radius-md;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 12rpx;
    border: 3rpx solid transparent;
    transition: all 0.3s;

    .icon {
      font-size: 32rpx;
    }
    .txt {
      font-size: 26rpx;
      font-weight: 700;
      color: $color-text-secondary;
    }

    &.male.active {
      background: #ebf4ff;
      border-color: #8cc7ff;
      .txt {
        color: #8cc7ff;
      }
    }
    &.female.active {
      background: #fff0f5;
      border-color: #ff9ebc;
      .txt {
        color: #ff9ebc;
      }
    }
  }
}

.safe-bottom-dock {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 40rpx 48rpx calc(40rpx + env(safe-area-inset-bottom));
  background: linear-gradient(180deg, transparent 0%, #fff 40%);
  z-index: 100;
}
</style>
