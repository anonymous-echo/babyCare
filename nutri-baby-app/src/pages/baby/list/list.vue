<template>
  <view class="baby-list-page">
    <!-- 头部 -->

    <!-- 宝宝列表 -->
    <view class="baby-list">
      <view
        v-for="baby in babyList"
        :key="baby.babyId"
        class="baby-card"
        :class="{
          active: baby.babyId === currentBabyId,
          'is-default': baby.babyId === userInfo?.defaultBabyId,
        }"
      >
        <!-- 默认标签 -->
        <view
          v-if="baby.babyId === userInfo?.defaultBabyId"
          class="default-badge"
        >
          <wd-icon name="star-on" size="12" color="#ff9800" />
          <text>默认</text>
        </view>

        <!-- 卡片头部 - 点击切换宝宝 -->
        <view class="card-header" @click="handleSelectBaby(baby.babyId)">
          <!-- 头像 -->
          <view class="baby-avatar">
            <image
              v-if="baby.avatarUrl"
              :src="baby.avatarUrl"
              mode="aspectFill"
            />
            <image v-else src="/static/default.png" mode="aspectFill" />
          </view>

          <!-- 信息 -->
          <view class="baby-info">
            <view class="name-row">
              <text class="baby-name">{{ baby.name }}</text>
              <text v-if="baby.nickname" class="nickname">{{
                baby.nickname
              }}</text>
            </view>
            <view class="baby-meta">
              <text class="gender">{{
                baby.gender === "male" ? "👦 男宝" : "👧 女宝"
              }}</text>
              <text class="divider">|</text>
              <text class="age">{{ calculateAge(baby.birthDate) }}</text>
            </view>
          </view>

          <!-- 选中标记 -->
          <view v-if="baby.babyId === currentBabyId" class="check-icon">
            <wd-icon name="check-circle-fill" size="24" color="#fa2c19" />
          </view>
        </view>

        <!-- 分割线 -->
        <view class="divider-line" />

        <!-- 操作按钮区域 -->
        <view class="card-actions" @click.stop>
          <!-- 邀请协作按钮（全宽） -->
          <wd-button
            size="small"
            plain
            type="primary"
            class="full-width-btn"
            @click="handleInvite(baby.babyId, baby.name)"
          >
            <wd-icon name="share" size="14" />
            邀请协作
          </wd-button>

          <!-- 设为默认按钮（全宽，仅当非默认宝宝时显示） -->
          <wd-button
            v-if="baby.babyId !== userInfo?.defaultBabyId"
            size="small"
            plain
            type="warning"
            class="full-width-btn"
            @click="handleSetDefault(baby.babyId, baby.name)"
          >
            <wd-icon name="star" size="14" />
            设为默认
          </wd-button>

          <!-- 协作者预览组件 -->
          <BabyCollaboratorsPreview
            :baby-id="baby.babyId"
            :collaborators="getCollaborators(baby.babyId)"
            @go-to-collaborators="() => handleGoToCollaborators(baby.babyId, baby.name)"
            @set-relationship="() => handleSetRelationship(baby.babyId, baby.name)"
          />

          <!-- 编辑和删除按钮（并排，各占50%） -->
          <view class="action-row">
            <wd-button
              size="small"
              plain
              type="warning"
              @click="handleEdit(baby.babyId)"
            >
              <wd-icon name="edit" size="14" />
              编辑
            </wd-button>
            <wd-button
              size="small"
              plain
              type="danger"
              @click="handleDelete(baby.babyId)"
            >
              <wd-icon name="delete-thin" size="14" />
              删除
            </wd-button>
          </view>
        </view>
      </view>

      <!-- 空状态 -->
      <wd-status-tip
        v-if="babyList.length === 0"
        description="还没有添加宝宝"
        image="empty"
      >
        <template #description>
          <text class="empty-text">还没有添加宝宝哦~</text>
        </template>
      </wd-status-tip>
    </view>

    <!-- 关系设置弹窗 -->
    <wd-popup
      v-model="relationshipDialog.show"
      position="bottom"
      custom-style="height: auto; padding: 0"
      safe-area-inset-bottom
    >
      <view class="relationship-popup">
        <view class="popup-header">
          <text class="popup-title">设置与{{ relationshipDialog.babyName }}的关系</text>
          <wd-icon name="close" @click="relationshipDialog.show = false" />
        </view>

        <!-- 自定义输入 -->
        <view class="custom-input-section">
          <wd-input
            v-model="relationshipDialog.customRelationship"
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
            :class="{ active: relationshipDialog.selectedRelationship === option.value }"
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

    <!-- 添加按钮 -->
    <view class="add-button">
      <wd-button type="primary" size="large" block @click="handleAdd">
        <wd-icon name="plus" size="18" />
        添加宝宝
      </wd-button>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import { currentBabyId, setCurrentBaby, getCollaborators, setCollaborators } from "@/store/baby";
import { userInfo, setDefaultBaby, getUserInfo } from "@/store/user";
import { calculateAge } from "@/utils/date";
import BabyCollaboratorsPreview from "@/components/BabyCollaboratorsPreview.vue";
import { updateFamilyMember } from "@/store/collaborator";

// 直接调用 API 层
import * as babyApi from "@/api/baby";
import * as collaboratorApi from "@/api/collaborator";

// 宝宝列表(从 API 获取)
const babyList = ref<babyApi.BabyProfileResponse[]>([]);

// 关系设置弹窗状态
const relationshipDialog = ref({
  show: false,
  babyId: '',
  babyName: '',
  selectedRelationship: '',
  customRelationship: '',
});

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

// 加载宝宝列表
const loadBabyList = async () => {
  try {
    const data = await babyApi.apiFetchBabyList();
    babyList.value = data;

    // 并行加载所有宝宝的协作者信息
    await Promise.all(
      data.map(async (baby) => {
        try {
          const collaborators = await collaboratorApi.apiFetchCollaborators(baby.babyId);
          setCollaborators(baby.babyId, collaborators);
        } catch (error) {
          console.warn(`[BabyList] 加载宝宝 ${baby.babyId} 的协作者失败:`, error);
          // 协作者加载失败不影响宝宝列表显示
        }
      })
    );

    // 如果只有一个宝宝且没有选中任何宝宝,默认选中这个宝宝
    if (babyList.value.length === 1 && !currentBabyId.value) {
      const firstBaby = babyList.value[0];
      if (firstBaby) {
        setCurrentBaby(firstBaby.babyId);
        console.log("[BabyList] 自动选中唯一的宝宝:", firstBaby.name);
      }
    }
  } catch (error) {
    console.error("[BabyList] 加载宝宝列表失败:", error);
    uni.showToast({
      title: "加载失败",
      icon: "none",
    });
  }
};

// 页面加载时初始化
onMounted(() => {
  loadBabyList();
});

// 选择宝宝
const handleSelectBaby = (id: string) => {
  setCurrentBaby(id);
  console.log("[BabyList] 切换宝宝:", id);
  uni.showToast({
    title: "已切换",
    icon: "success",
    duration: 1000,
  });

  // 延迟返回首页
  setTimeout(() => {
    uni.navigateBack();
  }, 1000);
};

// 设置为默认宝宝
const handleSetDefault = async (id: string, name: string) => {
  try {
    await setDefaultBaby(id);
    console.log("[BabyList] 设置默认宝宝:", name);
  } catch (error) {
    console.error("[BabyList] 设置默认宝宝失败:", error);
  }
};

// 添加宝宝
const handleAdd = () => {
  uni.navigateTo({
    url: "/pages/baby/edit/edit",
  });
};

// 邀请协作者
const handleInvite = (id: string, name: string) => {
  uni.navigateTo({
    url: `/pages/baby/invite/invite?babyId=${id}&babyName=${encodeURIComponent(
      name
    )}`,
  });
};

// 进入协作者管理页面
const handleGoToCollaborators = (babyId: string, babyName: string) => {
  uni.navigateTo({
    url: `/pages/baby/collaborators/collaborators?babyId=${babyId}&babyName=${encodeURIComponent(
      babyName
    )}`,
  });
};

// 编辑宝宝
const handleEdit = (id: string) => {
  uni.navigateTo({
    url: `/pages/baby/edit/edit?id=${id}`,
  });
};

// 删除宝宝
const handleDelete = (id: string) => {
  uni.showModal({
    title: "确认删除",
    content: "删除后无法恢复,确定要删除这个宝宝吗?",
    success: async (res) => {
      if (res.confirm) {
        try {
          await babyApi.apiDeleteBaby(id);

          uni.showToast({
            title: "删除成功",
            icon: "success",
          });

          // 重新加载宝宝列表
          await loadBabyList();

          // 如果删除的是当前选中的宝宝,需要清除选中状态
          if (id === currentBabyId.value) {
            setCurrentBaby("");
          }
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

// 设置关系
const handleSetRelationship = (babyId: string, babyName: string) => {
  // 获取当前用户在该宝宝中的关系
  const collaborators = getCollaborators(babyId) || [];
  const currentUser = getUserInfo();
  const myCollaborator = collaborators.find(c => c.openid === currentUser?.openid);
  
  relationshipDialog.value = {
    show: true,
    babyId,
    babyName,
    selectedRelationship: myCollaborator?.relationship || '',
    customRelationship: '',
  };
};

// 选择预设关系
const selectRelationship = (value: string) => {
  relationshipDialog.value.selectedRelationship = value;
  relationshipDialog.value.customRelationship = '';
};

// 确认关系设置
const confirmRelationship = async () => {
  const { babyId, selectedRelationship, customRelationship } = relationshipDialog.value;
  
  // 优先使用自定义输入
  const finalRelationship = customRelationship.trim() || selectedRelationship;
  
  if (!finalRelationship) {
    uni.showToast({
      title: '请选择或输入关系',
      icon: 'none',
    });
    return;
  }
  
  try {
    const currentUser = getUserInfo();
    if (!currentUser?.openid) {
      uni.showToast({
        title: '用户信息异常',
        icon: 'none',
      });
      return;
    }
    
    await updateFamilyMember(babyId, currentUser.openid, {
      relationship: finalRelationship,
    });
    
    // 更新本地数据
    const collaborators = getCollaborators(babyId) || [];
    const myCollaborator = collaborators.find(c => c.openid === currentUser.openid);
    if (myCollaborator) {
      myCollaborator.relationship = finalRelationship;
      setCollaborators(babyId, [...collaborators]);
    }
    
    relationshipDialog.value.show = false;
    
  } catch (error: any) {
    console.error('设置关系失败:', error);
    uni.showToast({
      title: error.message || '设置失败',
      icon: 'none',
    });
  }
};
</script>

<style lang="scss" scoped>
@import '@/styles/colors.scss';
.baby-list-page {
  min-height: 100vh;
  background: $gradient-bg-light;
  padding-bottom: 140rpx;
}

.header {
  background: $color-bg-primary;
  padding: 40rpx 30rpx;
  text-align: center;
  box-shadow: $shadow-sm;
}

.title {
  font-size: 36rpx;
  font-weight: $font-weight-bold;
  color: $color-text-primary;
}

.baby-list {
  padding: 24rpx;
}

/* 卡片样式 */
.baby-card {
  background: $color-bg-primary;
  border-radius: $radius-xl;
  margin-bottom: $spacing-2xl;
  overflow: hidden;
  box-shadow: $shadow-md;
  transition: all $transition-slow;
  position: relative;

  &.active {
    box-shadow: 0 4rpx 20rpx rgba(50, 220, 110, 0.25);
    border: 2px solid $color-primary;
  }

  &.is-default {
    background: linear-gradient(135deg, rgba(50, 220, 110, 0.05) 0%, $color-bg-primary 20%);
  }
}

/* 默认标签 */
.default-badge {
  position: absolute;
  top: 16rpx;
  right: 16rpx;
  background: linear-gradient(135deg, $color-primary 0%, $color-primary-light 100%);
  color: white;
  font-size: 22rpx;
  padding: 8rpx 16rpx;
  border-radius: $radius-xl;
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 6rpx;
  font-weight: $font-weight-bold;
  box-shadow: $shadow-primary-md;
  z-index: 10;

  text {
    line-height: 1;
  }

  .nut-icon {
    line-height: 1;
  }
}

/* 卡片头部 */
.card-header {
  padding: 30rpx;
  display: flex;
  align-items: center;
  cursor: pointer;
  transition: background $transition-base;

  &:active {
    background: rgba(0, 0, 0, 0.02);
  }
}

.baby-avatar {
  width: 120rpx;
  height: 120rpx;
  border-radius: $radius-full;
  overflow: hidden;
  flex-shrink: 0;
  box-shadow: $shadow-md;

  image {
    width: 100%;
    height: 100%;
  }

  .avatar-placeholder {
    width: 100%;
    height: 100%;
    background: $gradient-primary;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 48rpx;
    font-weight: $font-weight-bold;
    color: white;
  }
}

.baby-info {
  flex: 1;
  margin-left: 24rpx;
  overflow: hidden;
}

.name-row {
  display: flex;
  align-items: center;
  gap: 12rpx;
  margin-bottom: 12rpx;
  flex-wrap: wrap;
}

.baby-name {
  font-size: 34rpx;
  font-weight: $font-weight-bold;
  color: $color-text-primary;
  line-height: 1.2;
}

.nickname {
  font-size: 26rpx;
  color: $color-text-secondary;
  background: $color-bg-secondary;
  padding: 4rpx 12rpx;
  border-radius: $radius-md;
  font-weight: $font-weight-normal;
}

.baby-meta {
  font-size: 26rpx;
  color: $color-text-secondary;
  display: flex;
  align-items: center;
  gap: 12rpx;

  .divider {
    color: $color-border-light;
  }

  .gender {
    font-weight: $font-weight-medium;
  }

  .age {
    color: $color-text-secondary;
  }
}

.check-icon {
  margin-left: 16rpx;
  flex-shrink: 0;
  animation: scaleIn $transition-slow;
}

@keyframes scaleIn {
  from {
    transform: scale(0);
  }
  to {
    transform: scale(1);
  }
}

/* 分割线 */
.divider-line {
  height: 1rpx;
  background: linear-gradient(
    90deg,
    transparent 0%,
    $color-divider 50%,
    transparent 100%
  );
  margin: 0 30rpx;
}

/* 操作按钮区域 */
.card-actions {
  padding: $spacing-md 30rpx 30rpx;
  display: flex;
  flex-direction: column;
  gap: $spacing-md;
}

.full-width-btn {
  // 兼容不同组件库渲染类名，保证按钮能占满整行
  :deep(.nut-button),
  :deep(.wd-button) {
    width: 100%;
    height: 64rpx;
    font-size: 26rpx;
    border-radius: $radius-md;
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: center;
    gap: 8rpx;
    transition: all $transition-base;

    &:active {
      transform: scale(0.96);
    }

    // 确保图标和文字垂直居中对齐
    .nut-icon {
      line-height: 1;
      vertical-align: middle;
    }
  }
}

.action-row {
  // 使用两列网格布局，保证两个按钮各占 50% 且与上方全宽按钮保持同宽
  display: grid;
  grid-template-columns: 1fr 1fr;
  column-gap: $spacing-md;
  width: 100%;

  :deep(.nut-button),
  :deep(.wd-button) {
    width: 100%;
    height: 64rpx;
    font-size: 26rpx;
    border-radius: $radius-md;
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: center;
    gap: 8rpx;
    transition: all $transition-base;

    &:active {
      transform: scale(0.96);
    }

    // 确保图标和文字垂直居中对齐
    .nut-icon {
      line-height: 1;
      vertical-align: middle;
    }
  }
}

/* 空状态 */
.empty-text {
  color: $color-text-secondary;
  font-size: 28rpx;
}

/* 添加按钮 */
.add-button {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: $spacing-2xl;
  background: linear-gradient(180deg, transparent 0%, $color-bg-primary 20%);
  backdrop-filter: blur(10rpx);

  :deep(.nut-button) {
    height: 88rpx;
    font-size: 32rpx;
    border-radius: $radius-lg;
    box-shadow: $shadow-primary-md;
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: center;
    gap: $spacing-md;

    &:active {
      transform: scale(0.98);
    }

    // 图标文字对齐
    .nut-icon {
      line-height: 1;
    }
  }
}

// ===== 关系设置弹窗 =====
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
