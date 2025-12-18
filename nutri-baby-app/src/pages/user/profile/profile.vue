<template>
  <view>
    <wd-navbar
      title="编辑资料"
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
    <view class="profile-page">
      <wd-message-box />
      <wd-toast />

      <!-- 头像编辑 -->
      <wd-cell-group custom-class="group" title="个人信息" border>
        <wd-cell title="头像" title-width="200rpx">
          <view class="avatar-section">
            <button open-type="chooseAvatar" @chooseavatar="onChooseAvatar">
              <image
                :src="formData.avatarUrl || '/static/default.png'"
                class="avatar-preview"
                mode="aspectFill"
              />
            </button>
          </view>
        </wd-cell>

        <!-- 昵称编辑 -->
        <wd-cell title="昵称" title-width="200rpx">
          <view class="nickname-section">
            <!-- 微信原生昵称输入框 -->
            <input
              type="nickname"
              class="nickname-input"
              v-model="formData.nickName"
              placeholder="请输入昵称"
              maxlength="20"
              @blur="onNicknameBlur"
            />
          </view>
        </wd-cell>

        <!-- 用户ID(只读) -->
        <wd-cell title="用户ID" title-width="200rpx">
          <text class="user-id">{{
            formData.openid ? formData.openid.slice(0, 12) + "..." : "-"
          }}</text>
        </wd-cell>
      </wd-cell-group>

      <!-- 底部按钮 -->
      <view class="footer">
        <wd-button
          type="primary"
          size="large"
          block
          :loading="isSubmitting"
          @click="handleSubmit"
        >
          保存更改
        </wd-button>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { getUserInfo, setUserInfo } from "@/store/user";
import { uploadFile } from "@/utils/request";
import * as authApi from "@/api/auth";
import { goBack, goBackHome } from "@/utils/common";

// 表单数据
const formData = ref({
  openid: "",
  nickName: "",
  avatarUrl: "",
});

// 提交状态
const isSubmitting = ref(false);

// 页面加载
onMounted(() => {
  const user = getUserInfo();
  if (user) {
    formData.value = {
      openid: user.openid,
      nickName: user.nickName || "",
      avatarUrl: user.avatarUrl || "",
    };
  }
});

// 微信头像选择器
const onChooseAvatar = async (e: any) => {
  console.log("[Profile] 选择微信头像:", e.detail.avatarUrl);

  // 使用封装的uploadFile函数
  const uploadResult: any = await uploadFile({
    filePath: e.detail.avatarUrl,
    name: "file",
    formData: {
      type: "user_avatar",
    },
  });
  console.log("[Profile] 上传结果:", uploadResult);
  // 解析响应数据
  if (uploadResult.code === 0) {
    formData.value.avatarUrl = uploadResult.data.url;

    uni.showToast({
      title: "上传成功",
      icon: "success",
    });
  }
  uni.showToast({
    title: "头像已更新",
    icon: "success",
    duration: 1500,
  });
};

// 上传本地图片
const uploadLocalImage = () => {
  uni.chooseImage({
    count: 1,
    sizeType: ["compressed"],
    sourceType: ["album", "camera"],
    success: async (res) => {
      const tempFilePath = res.tempFilePaths[0];
      if (!tempFilePath) return;

      try {
        // 显示上传中提示
        uni.showLoading({
          title: "上传中...",
          mask: true,
        });

        // 使用封装的uploadFile函数
        const uploadResult: any = await uploadFile({
          filePath: tempFilePath,
          name: "file",
          formData: {
            type: "user_avatar",
          },
        });

        // 解析响应数据
        if (uploadResult.code === 0) {
          formData.value.avatarUrl = uploadResult.data.url;
          isSubmitting.value = true;

          // 调用更新接口
          const result = await authApi.apiUpdateUserInfo({
            nickName: formData.value.nickName,
            avatarUrl: formData.value.avatarUrl,
          });
          if (result) {
            // 更新本地状态
            const user = getUserInfo();
            if (user) {
              setUserInfo({
                ...user,
                nickName: formData.value.nickName,
                avatarUrl: formData.value.avatarUrl,
              });
            }
          }
        } else {
          throw new Error(uploadResult.message || "上传失败");
        }
      } catch (error: any) {
        console.error("[Profile] 上传头像失败:", error);
        uni.showToast({
          title: error.message || "上传失败",
          icon: "none",
        });
      } finally {
        uni.hideLoading();
      }
    },
    fail: (err) => {
      console.error("[Profile] 选择图片失败:", err);
      uni.showToast({
        title: "选择图片失败",
        icon: "none",
      });
    },
  });
};

// 昵称输入完成
const onNicknameBlur = () => {
  console.log("[Profile] 昵称已输入:", formData.value.nickName);
};

// 提交表单
const handleSubmit = async () => {
  if (!formData.value.nickName.trim()) {
    uni.showToast({
      title: "请输入昵称",
      icon: "none",
    });
    return;
  }

  try {
    isSubmitting.value = true;

    // 调用更新接口
    await authApi.apiUpdateUserInfo({
      nickName: formData.value.nickName,
      avatarUrl: formData.value.avatarUrl,
    });

    // 更新本地状态
    const user = getUserInfo();
    if (user) {
      setUserInfo({
        ...user,
        nickName: formData.value.nickName,
        avatarUrl: formData.value.avatarUrl,
      });
    }

    uni.showToast({
      title: "保存成功",
      icon: "success",
    });

    setTimeout(() => {
      uni.navigateBack();
    }, 1000);
  } catch (error: any) {
    console.error("[Profile] 保存失败:", error);
    uni.showToast({
      title: error.message || "保存失败",
      icon: "none",
    });
  } finally {
    isSubmitting.value = false;
  }
};
</script>

<style lang="scss" scoped>
.profile-page {
  min-height: 100vh;
  background: #f5f5f5;
  padding-bottom: 100rpx;
}

.avatar-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 24rpx;
  padding: 20rpx 0;
  width: 100%;

  button {
    padding: 0;
    margin: 0;
    background: transparent;
    border: none;
    line-height: 1;

    &::after {
      border: none;
    }
  }
}

.avatar-section-button {
  width: 160rpx;
  height: 160rpx;
  border-radius: 50%;
}

.avatar-preview {
  width: 160rpx;
  height: 160rpx;
  border-radius: 50%;
  border: 4rpx solid #e0e0e0;
  object-fit: cover;
}

.avatar-actions {
  display: flex;
  flex-direction: column;
  gap: 12rpx;
  width: 100%;
}

.avatar-btn {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border: none;
  border-radius: 8rpx;
  padding: 12rpx 24rpx;
  font-size: 28rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8rpx;
  line-height: 1.5;
}

.avatar-btn::after {
  border: none;
}

.nickname-section {
  width: 100%;
}

.nickname-input {
  width: 100%;
  height: 48rpx;
  padding: 12rpx;
  font-size: 28rpx;
  border: none;
  border-radius: 8rpx;
  box-sizing: border-box;
}

.user-id {
  font-size: 24rpx;
  color: #999;
}

.footer {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 24rpx;
  background: white;
  box-shadow: 0 -2rpx 8rpx rgba(0, 0, 0, 0.1);
  z-index: 10;
}
</style>
