/** * 协作者管理页面框架 * 路径: pages/baby/collaborators/collaborators.vue *
职责: 显示和管理宝宝的所有协作者 */

<template>
  <view class="collaborators-page">
    <!-- 导航栏 -->
    <wd-navbar
      fixed
      safeAreaInsetTop
      placeholder
      :title="`${babyName}的亲友团`"
      left-arrow
    >
      <template #capsule>
        <wd-navbar-capsule @back="goBack" @back-home="goBackHome" />
      </template>
    </wd-navbar>

    <!-- 页面内容 -->
    <view class="page-content">
      <!-- 搜索框 -->
      <view class="search-section">
        <wd-search
          v-model="searchKeyword"
          placeholder="搜索亲友名称"
          clearable
          @search="onSearch"
        />
      </view>

      <!-- 亲友团列表 -->
      <view v-if="filteredCollaborators.length > 0" class="collaborators-list">
        <view
          v-for="collaborator in filteredCollaborators"
          :key="collaborator.openid"
          class="collaborator-card"
        >
          <!-- 协作者头像和基本信息 -->
          <view class="collaborator-header">
            <image
              v-if="collaborator.avatarUrl"
              :src="collaborator.avatarUrl"
              class="avatar"
              mode="aspectFill"
            />
            <view v-else class="avatar-placeholder">
              {{ collaborator.nickName.charAt(0) }}
            </view>

            <view class="info">
              <text class="name">
                {{ collaborator.nickName }}
                <text v-if="isCurrentUser(collaborator)" class="self-badge"
                  >自己</text
                >
              </text>
              <text class="role">
                {{ getRoleLabel(collaborator.role) }}
                <text v-if="collaborator.isCreator" class="creator-badge">
                  创建者
                </text>
              </text>
              <text v-if="collaborator.relationship" class="relationship">
                {{ collaborator.relationship }}
              </text>
            </view>
          </view>

          <!-- 权限信息 -->
          <view class="permission-info">
            <text
              v-if="collaborator.accessType === 'permanent'"
              class="access-type"
            >
              永久权限
            </text>
            <text v-else class="access-type">
              有效期至
              {{
                collaborator.expiresAt
                  ? formatDate(collaborator.expiresAt * 1000, "YYYY-MM-DD")
                  : "未设置"
              }}
            </text>
          </view>

          <!-- 权限显示 -->
          <view class="permission-details">
            <text class="permission-badge" :class="`role-${collaborator.role}`">
              {{ getRoleLabel(collaborator.role) }}
            </text>
          </view>

          <!-- 操作按钮 -->
          <view
            v-if="!collaborator.isCreator && !isCurrentUser(collaborator)"
            class="actions"
          >
            <wd-button
              size="small"
              plain
              type="warning"
              @click="showRoleChangeDialog(collaborator)"
            >
              变更权限
            </wd-button>
            <wd-button
              size="small"
              plain
              type="danger"
              @click="showRemoveConfirm(collaborator)"
            >
              移除
            </wd-button>
          </view>
        </view>
      </view>

      <!-- 空状态 -->
      <view v-else class="empty-state">
        <wd-status-tip description="暂无协作者" image="empty">
          <template #description>
            <text class="tip-text">您可以点击下方按钮邀请亲友</text>
          </template>
        </wd-status-tip>
      </view>
    </view>

    <!-- 邀请按钮 -->
    <view class="invite-button">
      <wd-button type="primary" size="large" block @click="goToInvite">
        <wd-icon name="plus" size="18" />
邀请新亲友
      </wd-button>
    </view>

    <!-- 角色变更弹窗 -->
    <wd-popup
      v-model="showRoleDialog"
      position="bottom"
      custom-style="height: auto; display: flex; flex-direction: column; padding: 0"
      safe-area-inset-bottom
    >
      <view class="role-popup-wrapper">
        <view class="role-popup-content">
          <wd-cell-group
            border
            :title="`编辑 ${currentCollaborator?.nickName} 的信息`"
          >
            <wd-cell title="与宝宝的关系" is-link @click="showRelationshipPicker = true">
              <text>{{ newRelationship || '未设置' }}</text>
            </wd-cell>

            <wd-cell title="亲友角色" title-width="150rpx">
              <view style="text-align: left">
                <wd-radio-group v-model="newRole" inline>
                  <wd-radio value="admin">管理员</wd-radio>
                  <wd-radio value="editor">编辑者</wd-radio>
                  <wd-radio value="viewer">查看者</wd-radio>
                </wd-radio-group>
              </view>
            </wd-cell>
            <wd-cell title="访问权限" title-width="150rpx">
              <view style="text-align: left">
                <wd-radio-group v-model="newAccessType" inline>
                  <wd-radio value="permanent">永久</wd-radio>
                  <wd-radio value="temporary">临时</wd-radio>
                </wd-radio-group>
              </view>
            </wd-cell>
            <wd-datetime-picker
              v-if="newAccessType === 'temporary'"
              v-model="newExpiresAt"
              type="datetime"
              label="到期时间"
            />
          </wd-cell-group>
        </view>
        <!-- 确认和取消按钮 - 固定在底部 -->
        <view class="popup-actions">
          <wd-button size="medium" plain type="warning" @click="cancelRoleChange">
            取消
          </wd-button>
          <wd-button size="medium" type="primary" @click="confirmRoleChange">
            确认
          </wd-button>
        </view>
      </view>
    </wd-popup>

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
            :class="{ active: newRelationship === option.value }"
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
            @click="confirmRelationshipSelection"
          >
            确认
          </wd-button>
        </view>
      </view>
    </wd-popup>
  </view>
</template>

<script setup lang="ts">
import { ref, computed } from "vue";
import { onLoad } from "@dcloudio/uni-app";
import type { BabyCollaborator } from "@/types/collaborator";
import { getUserInfo } from "@/store/user";
import { getRoleLabel, getRoleDescription } from "@/utils/permission";
import {
  apiFetchCollaborators,
  apiUpdateCollaboratorRole,
  apiRemoveCollaborator,
} from "@/api/collaborator";
import { updateFamilyMember } from "@/store/collaborator";
import { formatDate } from "@/utils/date";
import { goBack, goBackHome } from "@/utils/common";

interface PageState {
  babyId: string;
  babyName: string;
  collaborators: BabyCollaborator[];
  searchKeyword: string;
  loading: boolean;

  // 角色变更相关
  showRoleDialog: boolean;
  currentCollaborator: BabyCollaborator | null;
  newRole: string;
  newRelationship: string;
  customRelationship: string;
  newAccessType: "permanent" | "temporary";
  newExpiresAt: number;
  showDatePicker: boolean;
  showRelationshipPicker: boolean;
}

const state = ref<PageState>({
  babyId: "",
  babyName: "",
  collaborators: [],
  searchKeyword: "",
  loading: false,

  showRoleDialog: false,
  currentCollaborator: null,
  newRole: "editor",
  newRelationship: "",
  customRelationship: "",
  newAccessType: "permanent",
  newExpiresAt: Date.now() + 365 * 24 * 60 * 60 * 1000,
  showDatePicker: false,
  showRelationshipPicker: false,
});

// ============ 页面参数 ============

// onLoad 获取页面参数
let pageParams = { babyId: "", babyName: "" };

onLoad((options: any) => {
  pageParams.babyId = options.babyId || "";
  pageParams.babyName = decodeURIComponent(options.babyName || "");
  state.value.babyId = pageParams.babyId;
  state.value.babyName = pageParams.babyName;
  loadCollaborators();
});

// ============ 计算属性 ============

const babyId = computed(() => state.value.babyId);
const babyName = computed(() => state.value.babyName);
const collaborators = computed(() => state.value.collaborators);
const searchKeyword = computed({
  get: () => state.value.searchKeyword,
  set: (val) => (state.value.searchKeyword = val),
});

const filteredCollaborators = computed(() => {
  if (!searchKeyword.value) {
    return collaborators.value;
  }
  return collaborators.value.filter((c) =>
    c.nickName.toLowerCase().includes(searchKeyword.value.toLowerCase())
  );
});

// 判断协作者是否为当前用户
const isCurrentUser = (collaborator: BabyCollaborator): boolean => {
  const currentUserInfo = getUserInfo();
  return currentUserInfo?.openid === collaborator.openid;
};

const showRoleDialog = computed({
  get: () => state.value.showRoleDialog,
  set: (val) => (state.value.showRoleDialog = val),
});

const currentCollaborator = computed(() => state.value.currentCollaborator);

const newRole = computed({
  get: () => state.value.newRole,
  set: (val) => (state.value.newRole = val),
});

const newAccessType = computed({
  get: () => state.value.newAccessType,
  set: (val) => (state.value.newAccessType = val),
});

const newExpiresAt = computed({
  get: () => state.value.newExpiresAt,
  set: (val) => (state.value.newExpiresAt = val),
});

const showDatePicker = computed({
  get: () => state.value.showDatePicker,
  set: (val) => (state.value.showDatePicker = val),
});

const newRelationship = computed({
  get: () => state.value.newRelationship,
  set: (val) => (state.value.newRelationship = val),
});

const showRelationshipPicker = computed({
  get: () => state.value.showRelationshipPicker,
  set: (val) => (state.value.showRelationshipPicker = val),
});

const customRelationship = computed({
  get: () => state.value.customRelationship,
  set: (val) => (state.value.customRelationship = val),
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

// ============ 方法 ============

const onBack = () => {
  uni.navigateBack();
};

const onSearch = () => {
  // 搜索已通过计算属性实现
  console.log("[Collaborators] 搜索:", searchKeyword.value);
};

const loadCollaborators = async () => {
  state.value.loading = true;
  try {
    const data = await apiFetchCollaborators(babyId.value);
    state.value.collaborators = data;
    console.log("[Collaborators] 加载协作者列表:", data);
  } catch (error) {
    console.error("[Collaborators] 加载失败:", error);
    uni.showToast({
      title: "加载失败",
      icon: "none",
    });
  } finally {
    state.value.loading = false;
  }
};

const showRoleChangeDialog = (collaborator: BabyCollaborator) => {
  state.value.currentCollaborator = collaborator;
  state.value.newRole = collaborator.role;
  state.value.newRelationship = collaborator.relationship || "";
  state.value.newAccessType = collaborator.accessType;
  state.value.newExpiresAt = collaborator.expiresAt || Date.now();
  state.value.showRoleDialog = true;
};

const onRoleChange = () => {
  console.log("[Collaborators] 角色已变更:", newRole.value);
};

const confirmRoleChange = async () => {
  if (!currentCollaborator.value) return;

  try {
    // 使用新的 updateFamilyMember API 同时更新角色和关系
    await updateFamilyMember(
      babyId.value,
      currentCollaborator.value.openid,
      {
        role: newRole.value as "admin" | "editor" | "viewer",
        relationship: newRelationship.value,
      }
    );

    // 更新本地列表
    const idx = collaborators.value.findIndex(
      (c) => c.openid === currentCollaborator.value?.openid
    );
    if (idx !== -1 && currentCollaborator.value) {
      const collab = collaborators.value[idx];
      if (collab) {
        collab.role = newRole.value as "admin" | "editor" | "viewer";
        collab.relationship = newRelationship.value;
        collab.accessType = newAccessType.value;
      }
    }

    uni.showToast({
      title: "信息已更新",
      icon: "success",
    });

    state.value.showRoleDialog = false;
  } catch (error) {
    console.error("[Collaborators] 更新失败:", error);
    uni.showToast({
      title: "更新失败",
      icon: "none",
    });
  }
};

const cancelRoleChange = () => {
  state.value.showRoleDialog = false;
};

const showRemoveConfirm = async (collaborator: BabyCollaborator) => {
  state.value.currentCollaborator = collaborator;

  // 使用 uni.showModal 显示确认弹窗
  uni.showModal({
    title: "确认移除",
    content: `确定要移除 ${collaborator.nickName} 吗？\n\n移除后，该协作者将无法访问此宝宝的任何数据`,
    confirmText: "确认移除",
    cancelText: "取消",
    confirmColor: "#ff6b6b",
    success: async (res) => {
      if (res.confirm) {
        // 用户点击确认
        await confirmRemove();
      }
      // 用户点击取消，不需要额外处理
    },
  });
};

const confirmRemove = async () => {
  if (!currentCollaborator.value) return;

  try {
    await apiRemoveCollaborator(babyId.value, currentCollaborator.value.openid);

    // 更新本地列表
    state.value.collaborators = collaborators.value.filter(
      (c) => c.openid !== currentCollaborator.value?.openid
    );

    uni.showToast({
      title: "已移除",
      icon: "success",
    });
  } catch (error) {
    console.error("[Collaborators] 移除失败:", error);
    uni.showToast({
      title: "移除失败",
      icon: "none",
    });
  }
};

const goToInvite = () => {
  uni.navigateTo({
    url: `/pages/baby/invite/invite?babyId=${
      babyId.value
    }&babyName=${encodeURIComponent(babyName.value)}`,
  });
};

// 选择预设关系
const selectRelationship = (value: string) => {
  state.value.newRelationship = value;
  state.value.customRelationship = ""; // 清空自定义输入
};

// 确认关系选择
const confirmRelationshipSelection = () => {
  // 优先使用自定义输入，否则使用预设选项
  if (state.value.customRelationship.trim()) {
    state.value.newRelationship = state.value.customRelationship.trim();
  }
  state.value.showRelationshipPicker = false;
};

// ============ 生命周期 ============

// 生命周期由 onLoad 处理参数加载
</script>

<style lang="scss" scoped>
@import "@/styles/colors.scss";

.collaborators-page {
  min-height: 100vh;
  background: $color-bg-secondary;
  padding-bottom: 120rpx;
}

.page-content {
  padding: $spacing-lg;
}

.search-section {
  margin-bottom: $spacing-2xl;
}

.collaborators-list {
  display: flex;
  flex-direction: column;
  gap: $spacing-lg;
}

.collaborator-card {
  background: $color-bg-primary;
  border-radius: $radius-lg;
  padding: $spacing-lg;
  box-shadow: $shadow-md;
}

.collaborator-header {
  display: flex;
  align-items: center;
  gap: $spacing-md;
  margin-bottom: $spacing-md;
}

.avatar {
  width: 60rpx;
  height: 60rpx;
  border-radius: $radius-full;
  flex-shrink: 0;
}

.avatar-placeholder {
  width: 60rpx;
  height: 60rpx;
  border-radius: $radius-full;
  background: $gradient-primary;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: $font-size-lg;
  font-weight: $font-weight-bold;
  flex-shrink: 0;
}

.info {
  display: flex;
  flex-direction: column;
  gap: 4rpx;
  flex: 1;
}

.name {
  font-size: $font-size-md;
  font-weight: $font-weight-semibold;
  color: $color-text-primary;
}

.role {
  font-size: $font-size-sm;
  color: $color-text-secondary;
  display: flex;
  align-items: center;
  gap: 4rpx;
}

.creator-badge {
  display: inline-block;
  background: $color-primary-lighter;
  color: $color-primary;
  padding: 2rpx 8rpx;
  border-radius: 4rpx;
  font-size: 20rpx;
  font-weight: 500;
}

.self-badge {
  display: inline-block;
  background: $color-secondary-lighter;
  color: $color-secondary;
  padding: 2rpx 8rpx;
  border-radius: 4rpx;
  font-size: 20rpx;
  font-weight: 500;
  margin-left: 8rpx;
}

.permission-info {
  padding: $spacing-md 0;
  border-bottom: 1rpx solid $color-border-light;
  margin-bottom: $spacing-md;
}

.access-type {
  font-size: $font-size-sm;
  color: $color-text-secondary;
}

.relationship {
  display: block;
  font-size: $font-size-sm;
  color: $color-text-secondary;
  margin-top: 4rpx;
}

.permission-details {
  margin-bottom: $spacing-md;
  padding: $spacing-md 0;
  border-bottom: 1rpx solid $color-border-light;
}

.permission-badge {
  display: inline-block;
  padding: 6rpx 12rpx;
  border-radius: $radius-sm;
  font-size: $font-size-sm;
  font-weight: $font-weight-medium;

  &.role-admin {
    background: linear-gradient(
      135deg,
      rgba(255, 165, 0, 0.1) 0%,
      rgba(255, 165, 0, 0.05) 100%
    );
    color: #ff9800;
    border: 1rpx solid rgba(255, 165, 0, 0.3);
  }

  &.role-editor {
    background: linear-gradient(
      135deg,
      rgba(50, 220, 110, 0.1) 0%,
      rgba(50, 220, 110, 0.05) 100%
    );
    color: $color-primary;
    border: 1rpx solid rgba(50, 220, 110, 0.3);
  }

  &.role-viewer {
    background: linear-gradient(
      135deg,
      rgba(93, 173, 226, 0.1) 0%,
      rgba(93, 173, 226, 0.05) 100%
    );
    color: $color-secondary;
    border: 1rpx solid rgba(93, 173, 226, 0.3);
  }
}

.actions {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: $spacing-md;
}

.empty-state {
  padding: $spacing-3xl 0;
}

.invite-button {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: $spacing-lg;
  background: linear-gradient(180deg, transparent 0%, $color-bg-primary 20%);

  :deep(.wd-button) {
    height: 88rpx;
    border-radius: $radius-lg;
  }
}

// 弹窗样式
.role-popup-wrapper {
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 0;
}

.role-popup-content {
  flex: 1;
  overflow-y: auto;
  padding: $spacing-lg 0;
}

.role-dialog-content,
.remove-dialog-content {
  padding: $spacing-lg;
}

.target-name {
  display: block;
  font-size: $font-size-md;
  color: $color-text-primary;
  font-weight: $font-weight-semibold;
  margin-bottom: $spacing-lg;
  text-align: center;
}

.role-options,
.access-type-options {
  margin-bottom: $spacing-lg;
}

.label {
  font-size: $font-size-md;
  color: $color-text-primary;
  font-weight: $font-weight-semibold;
  margin-bottom: $spacing-md;
  display: block;
}

.option {
  margin-bottom: $spacing-md;
}

.radio-content {
  display: flex;
  flex-direction: column;
  gap: 4rpx;
  align-items: flex-start;
}

.role-name {
  font-size: $font-size-sm;
  color: $color-text-primary;
  font-weight: $font-weight-medium;
}

.role-desc {
  font-size: $font-size-xs;
  color: $color-text-secondary;
}

.date-picker {
  margin-top: $spacing-md;
  padding: $spacing-md;
  background: $color-bg-secondary;
  border-radius: $radius-md;
}

.popup-actions {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: $spacing-md;
  padding: $spacing-lg;
  border-top: 1rpx solid $color-border-light;
  background: $color-bg-primary;
  flex-shrink: 0;
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
