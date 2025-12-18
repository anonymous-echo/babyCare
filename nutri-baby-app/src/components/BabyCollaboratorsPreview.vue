/**
 * 宝宝列表卡片组件 - 显示协作者信息
 * 展示协作者简要信息，点击可进入管理页面
 */

<!-- 在宝宝列表卡片中添加协作者显示区域 -->
<template>
  <!-- ... 现有的头部信息 ... -->

  <!-- 新增：协作者简览区域 -->
  <view class="collaborators-preview" @click.stop="goToCollaborators">
    <view class="collaborators-header">
      <wd-icon name="people" size="16" color="#32dc6e" />
      <text class="collaborators-count">
        亲友团
        ({{ collaboratorsPreview.total }})
      </text>
      <wd-icon name="arrow-right" size="12" color="#999" class="arrow" />
    </view>

    <!-- 协作者简要列表 -->
    <view class="collaborators-list" v-if="collaboratorsPreview.items.length > 0">
      <view
        v-for="(collaborator, index) in collaboratorsPreview.items"
        :key="index"
        class="collaborator-item"
      >
        <!-- 协作者头像 -->
        <image
          v-if="collaborator.avatarUrl"
          :src="collaborator.avatarUrl"
          class="collaborator-avatar"
          mode="aspectFill"
        />
        <view v-else class="collaborator-avatar-placeholder">
          {{ collaborator.nickName.charAt(0) }}
        </view>

        <!-- 协作者信息 -->
        <view class="collaborator-info">
          <text class="collaborator-name">{{ collaborator.nickName }}</text>
          <text class="collaborator-role">
            <text v-if="isCurrentUser(collaborator) && !collaborator.relationship" class="set-relationship-mini" @click.stop="handleSetMyRelationship">
              <wd-icon name="edit" size="10" />
              设置关系
            </text>
            <template v-else>
              {{ collaborator.relationship || '未设置关系' }}
              <text class="role-separator">·</text>
              {{ getRoleLabel(collaborator.role) }}
              <text v-if="collaborator.isCreator" class="creator-badge">
                创建者
              </text>
            </template>
          </text>
        </view>
      </view>

      <!-- 更多按钮 -->
      <view
        v-if="collaboratorsPreview.hasMore"
        class="collaborator-more"
      >
        <text class="more-text">+{{ collaboratorsPreview.more }}</text>
      </view>
    </view>

    <!-- 空状态 -->
    <view v-else class="collaborators-empty">
      <view class="empty-content">
        <text class="empty-text">仅你一人</text>
        <text class="set-relationship-hint" @click.stop="$emit('set-relationship')">
          <wd-icon name="edit" size="12" />
          设置关系
        </text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { BabyCollaborator } from '@/types/collaborator';
import { getRoleLabel } from '@/utils/permission';
import { getUserInfo } from '@/store/user';

interface Props {
  babyId: string;
  collaborators?: BabyCollaborator[];
}

const props = withDefaults(defineProps<Props>(), {
  collaborators: () => [],
});

const emit = defineEmits<{
  'go-to-collaborators': [babyId: string];
  'set-relationship': [];
}>();

/**
 * 协作者预览数据
 * 最多显示3个，超过则显示 "+N"
 */
const collaboratorsPreview = computed(() => {
  const items = props.collaborators || [];
  const displayCount = 3;
  const displayed = items.slice(0, displayCount);
  const hasMore = items.length > displayCount;
  const moreCount = items.length - displayCount;

  return {
    total: items.length,
    items: displayed,
    hasMore,
    more: moreCount,
  };
});

const goToCollaborators = () => {
  emit('go-to-collaborators', props.babyId);
};

// 判断是否是当前用户
const isCurrentUser = (collaborator: BabyCollaborator): boolean => {
  const currentUser = getUserInfo();
  return currentUser?.openid === collaborator.openid;
};

// 设置自己的关系
const handleSetMyRelationship = () => {
  emit('set-relationship');
};
</script>

<style lang="scss" scoped>
@import '@/styles/colors.scss';

.collaborators-preview {
  padding: $spacing-md $spacing-lg;
  background: $color-bg-secondary;
  border-radius: $radius-md;
  margin-bottom: $spacing-lg;
  cursor: pointer;
  transition: all $transition-base;

  &:active {
    background: rgba(50, 220, 110, 0.1);
  }
}

.collaborators-header {
  display: flex;
  align-items: center;
  gap: $spacing-sm;
  margin-bottom: $spacing-md;

  .arrow {
    margin-left: auto;
  }
}

.collaborators-count {
  font-size: $font-size-sm;
  color: $color-text-secondary;
  font-weight: $font-weight-medium;
  flex: 1;
}

.collaborators-list {
  display: flex;
  flex-direction: column;
  gap: $spacing-sm;
}

.collaborator-item {
  display: flex;
  align-items: center;
  gap: $spacing-md;
  padding: $spacing-sm 0;
}

.collaborator-avatar {
  width: 32rpx;
  height: 32rpx;
  border-radius: $radius-full;
  flex-shrink: 0;
}

.collaborator-avatar-placeholder {
  width: 32rpx;
  height: 32rpx;
  border-radius: $radius-full;
  background: $gradient-primary;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: $font-size-xs;
  font-weight: $font-weight-bold;
  flex-shrink: 0;
}

.collaborator-info {
  display: flex;
  flex-direction: column;
  gap: 2rpx;
  flex: 1;
  min-width: 0;
}

.collaborator-name {
  font-size: $font-size-sm;
  color: $color-text-primary;
  font-weight: $font-weight-medium;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.collaborator-role {
  font-size: $font-size-xs;
  color: $color-text-secondary;
  display: flex;
  align-items: center;
  gap: 4rpx;
}

.creator-badge {
  display: inline-block;
  background: $color-primary-lighter;
  color: $color-primary;
  padding: 2rpx 6rpx;
  border-radius: 4rpx;
  font-size: 20rpx;
  font-weight: 500;
  white-space: nowrap;
}

.collaborator-more {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: $spacing-sm 0;
  border-top: 1rpx solid $color-border-light;
  margin-top: $spacing-sm;
}

.more-text {
  font-size: $font-size-sm;
  color: $color-text-secondary;
  font-weight: $font-weight-medium;
}

.collaborators-empty {
  padding: $spacing-md 0;
  text-align: center;
  color: $color-text-tertiary;
  font-size: $font-size-sm;
}

.empty-content {
  display: flex;
  flex-direction: column;
  gap: $spacing-sm;
  align-items: center;
}

.empty-text {
  color: $color-text-tertiary;
  font-size: $font-size-sm;
}

.set-relationship-hint {
  color: $color-primary;
  font-size: $font-size-sm;
  font-weight: $font-weight-medium;
  cursor: pointer;
  padding: $spacing-sm $spacing-md;
  border-radius: $radius-md;
  background: rgba(50, 220, 110, 0.1);
  transition: all $transition-base;
  display: flex;
  align-items: center;
  gap: $spacing-xs;
  border: 1rpx solid rgba(50, 220, 110, 0.2);

  &:active {
    background: rgba(50, 220, 110, 0.2);
    transform: scale(0.98);
  }
}

.role-separator {
  color: $color-border-light;
  margin: 0 4rpx;
}

.set-relationship-mini {
  color: $color-primary;
  font-size: $font-size-xs;
  font-weight: $font-weight-medium;
  cursor: pointer;
  padding: 2rpx 6rpx;
  border-radius: 4rpx;
  background: rgba(50, 220, 110, 0.1);
  transition: all $transition-base;
  display: inline-flex;
  align-items: center;
  gap: 2rpx;
  border: 1rpx solid rgba(50, 220, 110, 0.2);

  &:active {
    background: rgba(50, 220, 110, 0.2);
    transform: scale(0.95);
  }
}
</style>
