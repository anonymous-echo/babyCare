<template>
  <view class="baby-list-container">
    <wd-navbar
      title="切换宝宝"
      left-text="返回"
      left-arrow
      fixed
      placeholder
      safe-area-inset-top
      @click-left="goBack"
    />

    <scroll-view class="baby-scroll" scroll-y>
      <view class="baby-list-content">
        <view
          v-for="baby in babyList"
          :key="baby.babyId"
          class="baby-premium-card"
          :class="{ active: baby.babyId === currentBabyId }"
          @click="handleSelectBaby(baby.babyId)"
        >
          <view class="card-bg-gradient"></view>

          <view class="card-main">
            <view class="avatar-box premium-shadow">
              <image
                :src="resolveImageUrl(baby.avatarUrl) || '/static/default.png'"
                mode="aspectFill"
                class="avatar"
              />
              <view v-if="baby.gender" class="gender-tag" :class="baby.gender">
                <text>{{ baby.gender === "male" ? "♂" : "♀" }}</text>
              </view>
            </view>

            <view class="info-box">
              <view class="name-row">
                <text class="name">{{ baby.name }}</text>
                <text
                  v-if="baby.babyId === userInfo?.defaultBabyId"
                  class="default-badge"
                  >默认</text
                >
              </view>
              <text class="age-text">{{ calculateAge(baby.birthDate) }}</text>

              <!-- 快捷状态展示 -->
              <view class="collaborator-row">
                <BabyCollaboratorsPreview
                  :baby-id="baby.babyId"
                  :collaborators="getCollaborators(baby.babyId)"
                  @click.stop
                />
              </view>
            </view>

            <view class="selection-indicator">
              <view class="check-orb">
                <wd-icon name="check-bold" size="14" color="#FFF" />
              </view>
            </view>
          </view>

          <!-- 卡片操作区 -->
          <view class="card-actions-dock" @click.stop>
            <view
              class="action-item"
              @click="handleInvite(baby.babyId, baby.name)"
            >
              <wd-icon name="share" size="16" />
              <text>邀请</text>
            </view>
            <view class="action-divider"></view>
            <view class="action-item" @click="handleEdit(baby.babyId)">
              <wd-icon name="edit-1" size="16" />
              <text>编辑</text>
            </view>
            <view
              v-if="baby.babyId !== userInfo?.defaultBabyId"
              class="action-divider"
            ></view>
            <view
              v-if="baby.babyId !== userInfo?.defaultBabyId"
              class="action-item special"
              @click="handleSetDefault(baby.babyId, baby.name)"
            >
              <wd-icon name="star" size="16" />
              <text>设为默认</text>
            </view>
          </view>
        </view>

        <!-- 添加卡片 -->
        <view class="add-baby-trigger premium-shadow" @click="handleAdd">
          <view class="add-pulsar">
            <wd-icon name="add" size="32" color="#7BD3A2" />
          </view>
          <text class="add-label">添加新成员</text>
        </view>
      </view>

      <view v-if="babyList.length === 0" class="empty-state-v2">
        <wd-status-tip image="content" description="还没有添加宝宝哦" />
      </view>
    </scroll-view>

    <!-- 关系设置弹窗 (Redesigned) -->
    <wd-popup
      v-model="relationshipDialog.show"
      position="bottom"
      round
      safe-area-inset-bottom
    >
      <view class="premium-popup-content">
        <view class="popup-header">
          <text class="popup-title">设置您的身份</text>
          <wd-icon
            name="close"
            size="24"
            @click="relationshipDialog.show = false"
          />
        </view>

        <view class="relationship-grid">
          <view
            v-for="opt in relationshipOptions"
            :key="opt.value"
            class="rel-pill"
            :class="{
              active: relationshipDialog.selectedRelationship === opt.value,
            }"
            @click="selectRelationship(opt.value)"
          >
            {{ opt.label }}
          </view>
        </view>

        <view class="custom-rel-input">
          <wd-input
            v-model="relationshipDialog.customRelationship"
            placeholder="或输入其他身份（如：干妈、表姐）"
            no-border
          />
        </view>

        <view class="popup-footer">
          <wd-button block round type="primary" @click="confirmRelationship"
            >确认并保存</wd-button
          >
        </view>
      </view>
    </wd-popup>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import {
  currentBabyId,
  setCurrentBaby,
  getCollaborators,
  setCollaborators,
} from "@/store/baby";
import { userInfo, setDefaultBaby, getUserInfo } from "@/store/user";
import { calculateAge } from "@/utils/date";
import { goBack } from "@/utils/common";
import { resolveImageUrl } from "@/utils/assets";
import BabyCollaboratorsPreview from "@/components/BabyCollaboratorsPreview.vue";
import { updateFamilyMember } from "@/store/collaborator";
import * as babyApi from "@/api/baby";
import * as collaboratorApi from "@/api/collaborator";

const babyList = ref<babyApi.BabyProfileResponse[]>([]);
const relationshipDialog = ref({
  show: false,
  babyId: "",
  babyName: "",
  selectedRelationship: "",
  customRelationship: "",
});

const relationshipOptions = [
  { label: "爸爸", value: "爸爸" },
  { label: "妈妈", value: "妈妈" },
  { label: "爷爷", value: "爷爷" },
  { label: "奶奶", value: "奶奶" },
  { label: "外公", value: "外公" },
  { label: "外婆", value: "外婆" },
  { label: "其他亲友", value: "其他亲友" },
];

const loadBabyList = async () => {
  try {
    const data = await babyApi.apiFetchBabyList();
    babyList.value = data;
    await Promise.all(
      data.map(async (baby) => {
        try {
          const collaborators = await collaboratorApi.apiFetchCollaborators(
            baby.babyId,
          );
          setCollaborators(baby.babyId, collaborators);
        } catch (e) {}
      }),
    );
    if (babyList.value.length === 1 && !currentBabyId.value) {
      const baby = babyList.value[0];
      if (baby) setCurrentBaby(baby.babyId);
    }
  } catch (error) {}
};

onMounted(() => loadBabyList());

const handleSelectBaby = (id: string) => {
  setCurrentBaby(id);
  uni.showToast({ title: "切换成功", icon: "success", duration: 1000 });
  setTimeout(() => uni.navigateBack(), 1000);
};

const handleSetDefault = async (id: string, name: string) => {
  try {
    await setDefaultBaby(id);
    uni.showToast({ title: "设置默认成功", icon: "success" });
  } catch (error) {}
};

const handleAdd = () => uni.navigateTo({ url: "/pages/baby/edit/edit" });
const handleInvite = (id: string, name: string) =>
  uni.navigateTo({
    url: `/pages/baby/invite/invite?babyId=${id}&babyName=${encodeURIComponent(name)}`,
  });
const handleEdit = (id: string) =>
  uni.navigateTo({ url: `/pages/baby/edit/edit?id=${id}` });

const selectRelationship = (value: string) => {
  relationshipDialog.value.selectedRelationship = value;
  relationshipDialog.value.customRelationship = "";
};

const confirmRelationship = async () => {
  const { babyId, selectedRelationship, customRelationship } =
    relationshipDialog.value;
  const finalRelationship = customRelationship.trim() || selectedRelationship;
  if (!finalRelationship) return;
  try {
    const currentUser = getUserInfo();
    if (!currentUser?.openid) return;
    await updateFamilyMember(babyId, currentUser.openid, {
      relationship: finalRelationship,
    });
    const collaborators = getCollaborators(babyId) || [];
    const myCollab = collaborators.find((c) => c.openid === currentUser.openid);
    if (myCollab) {
      myCollab.relationship = finalRelationship;
      setCollaborators(babyId, [...collaborators]);
    }
    relationshipDialog.value.show = false;
  } catch (error) {}
};
</script>

<style lang="scss" scoped>
@import "@/styles/colors.scss";

.baby-list-container {
  min-height: 100vh;
  background: $color-bg-secondary;
}

.baby-scroll {
  height: calc(100vh - 100rpx);
}

.baby-list-content {
  padding: 32rpx;
  display: flex;
  flex-direction: column;
  gap: 32rpx;
}

.baby-premium-card {
  position: relative;
  background: #fff;
  border-radius: $radius-lg;
  padding: 40rpx;
  overflow: hidden;
  border: 1px solid $color-border-light;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);

  .card-bg-gradient {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 120rpx;
    background: linear-gradient(
      180deg,
      rgba(123, 211, 162, 0.05) 0%,
      transparent 100%
    );
  }

  &.active {
    border-color: $color-primary;
    box-shadow: 0 12rpx 32rpx rgba(123, 211, 162, 0.15);
    .selection-indicator .check-orb {
      transform: scale(1);
      opacity: 1;
    }
  }

  .card-main {
    display: flex;
    gap: 32rpx;
    position: relative;
    z-index: 1;
    margin-bottom: 32rpx;
  }

  .avatar-box {
    position: relative;
    width: 140rpx;
    height: 140rpx;
    border-radius: $radius-full;
    .avatar {
      width: 100%;
      height: 100%;
      border-radius: 50%;
    }
    .gender-tag {
      position: absolute;
      bottom: 0;
      right: 0;
      width: 40rpx;
      height: 40rpx;
      border-radius: 50%;
      background: #fff;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 24rpx;
      border: 1px solid $color-divider;
      &.male {
        color: #8cc7ff;
      }
      &.female {
        color: #ff9ebc;
      }
    }
  }

  .info-box {
    flex: 1;
    .name-row {
      display: flex;
      align-items: center;
      gap: 12rpx;
      margin-bottom: 8rpx;
      .name {
        font-size: 34rpx;
        font-weight: 800;
        color: $color-text-primary;
      }
      .default-badge {
        font-size: 20rpx;
        background: $color-primary-lighter;
        color: $color-primary-dark;
        padding: 2rpx 12rpx;
        border-radius: 100rpx;
        font-weight: 700;
      }
    }
    .age-text {
      font-size: 26rpx;
      color: $color-text-tertiary;
      font-weight: 500;
    }
    .collaborator-row {
      margin-top: 20rpx;
    }
  }

  .selection-indicator {
    .check-orb {
      width: 48rpx;
      height: 48rpx;
      background: $color-primary;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      transform: scale(0.5);
      opacity: 0;
      transition: all 0.3s;
    }
  }
}

.card-actions-dock {
  border-top: 1rpx solid $color-divider;
  padding-top: 24rpx;
  display: flex;
  align-items: center;
  justify-content: space-around;

  .action-item {
    display: flex;
    align-items: center;
    gap: 8rpx;
    font-size: 24rpx;
    font-weight: 600;
    color: $color-text-secondary;
    padding: 12rpx 24rpx;
    border-radius: 100rpx;
    &.special {
      color: $color-warning;
    }
    &:active {
      background: $color-bg-secondary;
    }
  }

  .action-divider {
    width: 1px;
    height: 24rpx;
    background: $color-divider;
  }
}

.add-baby-trigger {
  background: #fff;
  border-radius: $radius-lg;
  height: 200rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16rpx;
  border: 4rpx dashed $color-border-light;

  .add-pulsar {
    animation: pulsar 2s infinite;
  }
  .add-label {
    font-size: 28rpx;
    font-weight: 700;
    color: $color-text-tertiary;
  }
}

@keyframes pulsar {
  0% {
    transform: scale(1);
    opacity: 1;
  }
  50% {
    transform: scale(1.1);
    opacity: 0.8;
  }
  100% {
    transform: scale(1);
    opacity: 1;
  }
}

.premium-popup-content {
  padding: 40rpx;
  .popup-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 40rpx;
    .popup-title {
      font-size: 34rpx;
      font-weight: 800;
      color: $color-text-primary;
    }
  }
}

.relationship-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 16rpx;
  margin-bottom: 32rpx;
  .rel-pill {
    padding: 16rpx 32rpx;
    background: $color-bg-secondary;
    border-radius: $radius-md;
    font-size: 26rpx;
    font-weight: 600;
    color: $color-text-secondary;
    border: 1px solid transparent;
    &.active {
      background: $color-primary-lighter;
      color: $color-primary;
      border-color: $color-primary;
    }
  }
}

.custom-rel-input {
  background: $color-bg-secondary;
  border-radius: $radius-md;
  padding: 12rpx 24rpx;
  margin-bottom: 48rpx;
}

.popup-footer {
  margin-top: 40rpx;
}
</style>
