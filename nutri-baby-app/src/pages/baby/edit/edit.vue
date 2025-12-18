<template>
  <view class="page">
    <wd-message-box />
    <wd-toast />

    <wd-form ref="formRef" :model="formData" :rules="formRules">
      <!-- 基础信息 -->
      <wd-cell-group custom-class="group" title="基础信息" border>
        <!-- 宝宝头像 -->
        <wd-cell title="宝宝头像" title-width="200rpx">
          <view class="avatar-section">
            <view class="avatar-preview">
              <!-- 用户上传的头像 -->
              <image
                v-if="formData.avatarUrl"
                :src="formData.avatarUrl"
                mode="aspectFill"
              />
              <!-- 默认头像 -->
              <image
                v-else
                src="@/static/default.png"
                mode="aspectFill"
              />
            </view>
            <wd-button size="small" class="avatar-btn" @click="chooseAvatar">
              <wd-icon name="photograph" size="16" />
              {{ formData.avatarUrl ? '更换头像' : '选择头像' }}
            </wd-button>
          </view>
        </wd-cell>

        <!-- 宝宝姓名 -->
        <wd-input
          label="宝宝姓名"
          label-width="200rpx"
          :maxlength="20"
          show-word-limit
          prop="name"
          required
          clearable
          v-model="formData.name"
          placeholder="请输入宝宝姓名"
        />

        <!-- 小名昵称 -->
        <wd-input
          label="小名昵称"
          label-width="200rpx"
          :maxlength="20"
          show-word-limit
          clearable
          v-model="formData.nickname"
          placeholder="请输入小名或昵称（可选）"
        />

        <!-- 性别 -->
        <wd-cell title="性别" title-width="200rpx" prop="gender" center>
          <view style="text-align: left">
            <wd-radio-group v-model="formData.gender" inline>
              <wd-radio value="male">
                <text>👦 男孩</text>
              </wd-radio>
              <wd-radio value="female">
                <text>👧 女孩</text>
            </wd-radio>
            </wd-radio-group>
          </view>
        </wd-cell>

        <!-- 出生日期 -->
        <wd-datetime-picker
          label="出生日期"
          label-width="200rpx"
          placeholder="请选择出生日期"
          prop="birthDate"
          type="date"
          @confirm="handleDateConfirm"
        />
      </wd-cell-group>
    </wd-form>

    <!-- 底部按钮 -->
    <view class="button-container">
      <wd-button type="primary" size="large" @click="handleSubmit" block :loading="isSubmitting">
        {{ isEdit ? '保存更改' : '添加宝宝' }}
      </wd-button>
      <wd-button v-if="isEdit" plain size="large" @click="handleCancel" block>
        取消
      </wd-button>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { formatDate } from '@/utils/date'
import { uploadFile } from '@/utils/request'

// 直接调用 API 层
import * as babyApi from '@/api/baby'
import * as vaccineApi from '@/api/vaccine'

// 导入 store 以更新宝宝列表
import { fetchBabyDetail } from '@/store/baby'

// 表单数据
const formData = ref({
  name: '',
  nickname: '',
  gender: 'male' as 'male' | 'female',
  birthDate: '',
  avatarUrl: '',
})

// 表单验证规则
const formRules = ref({
  name: [
    { required: true, message: '请输入宝宝姓名', errorType: 'message' },
    { validator: (val: string) => val.trim().length > 0, message: '宝宝姓名不能为空', errorType: 'message' },
  ],
  nickname: [],
  gender: [
    { required: true, message: '请选择宝宝性别', errorType: 'message' },
  ],
  birthDate: [
    { required: true, message: '请选择出生日期', errorType: 'message' },
  ],
})

// 是否为编辑模式
const isEdit = ref(false)
const editId = ref('')

// 日期选择器
const selectedDate = ref()

// 提交状态
const isSubmitting = ref(false)

// 表单 ref
const formRef = ref<any>(null)

// 页面加载
onMounted(async () => {
  // 获取页面参数
  const pages = getCurrentPages()
  const currentPage = pages[pages.length - 1] as any
  const options = currentPage.options || {}

  if (options.id) {
    // 编辑模式
    isEdit.value = true
    editId.value = options.id

    try {
      const baby = await babyApi.apiFetchBabyDetail(options.id)
      if (baby) {
        formData.value = {
          name: baby.name,
          nickname: baby.nickname || '',
          gender: baby.gender,
          birthDate: baby.birthDate,
          avatarUrl: baby.avatarUrl || '',
        }

        // 设置选中的日期
        selectedDate.value = new Date(baby.birthDate)
      }
    } catch (error) {
      console.error('加载宝宝信息失败:', error)
      uni.showToast({
        title: '加载数据失败',
        icon: 'none'
      })
    }
  }
})

// 选择头像
const chooseAvatar = () => {
  uni.chooseImage({
    count: 1,
    sizeType: ['compressed'],
    sourceType: ['album', 'camera'],
    success: async (res) => {
      const tempFilePath = res.tempFilePaths[0]
      if (!tempFilePath) return

      try {
        // 显示上传中提示
        uni.showLoading({
          title: '上传中...',
          mask: true
        })

        // 调用上传接口
        const uploadResult: any = await uploadFile({
          filePath: tempFilePath,
          name: 'file',
          formData: {
            type: 'baby_avatar',
            related_id: isEdit.value ? editId.value : ''
          }
        })

        // 解析响应数据
        if (uploadResult.code === 0) {
          formData.value.avatarUrl = uploadResult.data.url
          uni.showToast({
            title: '上传成功',
            icon: 'success'
          })
        } else {
          throw new Error(uploadResult.message || '上传失败')
        }
      } catch (error: any) {
        console.error('上传头像失败:', error)
        uni.showToast({
          title: error.message || '上传失败',
          icon: 'none'
        })
      } finally {
        uni.hideLoading()
      }
    }
  })
}

// 日期确认
const handleDateConfirm = (val: any) => {
  const { value } = val
  console.log('selectedValue:', value)
  formData.value.birthDate = formatDate(value, 'YYYY-MM-DD')
}

// 提交表单
const handleSubmit = async () => {
  try {
    // 验证表单
    const valid = await formRef.value?.validate()
    if (!valid) {
      return
    }

    isSubmitting.value = true

    if (isEdit.value) {
      // 更新
      await babyApi.apiUpdateBaby(editId.value, {
        name: formData.value.name,
        nickname: formData.value.nickname,
        gender: formData.value.gender,
        birthDate: formData.value.birthDate,
        avatarUrl: formData.value.avatarUrl,
      })

      // 同步 store 中的宝宝数据，确保列表页面能看到最新信息
      try {
        await fetchBabyDetail(editId.value)
        console.log('[BabyEdit] 宝宝信息已同步到 store')
      } catch (error) {
        console.warn('[BabyEdit] 同步宝宝信息失败:', error)
        // 同步失败不影响用户体验，继续返回
      }

      uni.showToast({
        title: '更新成功',
        icon: 'success'
      })

      setTimeout(() => {
        uni.navigateBack()
      }, 1000)
    } else {
      // 添加
      const newBaby = await babyApi.apiCreateBaby({
        name: formData.value.name,
        nickname: formData.value.nickname,
        gender: formData.value.gender,
        birthDate: formData.value.birthDate,
        avatarUrl: formData.value.avatarUrl,
      })

      // ✨ 为新宝宝自动获取疫苗计划
      console.log('[BabyEdit] 为新宝宝获取疫苗计划:', newBaby.name)

      try {
        // 从服务器获取该宝宝的疫苗计划
        await vaccineApi.apiFetchVaccinePlans(newBaby.babyId)

        // 显示友好的提示
        uni.showModal({
          title: '✅ 宝宝添加成功',
          content: `已为 ${newBaby.name} 自动生成国家免疫规划疫苗计划和接种提醒，可在"疫苗管理"页面查看详情。`,
          showCancel: false,
          confirmText: '好的',
          success: () => {
            // 跳转到首页
            uni.reLaunch({
              url: '/pages/index/index'
            })
          }
        })
      } catch (vaccineError) {
        console.error('获取疫苗计划失败:', vaccineError)
        // 即使疫苗计划获取失败,宝宝添加仍然成功
        uni.showToast({
          title: '宝宝添加成功',
          icon: 'success'
        })
        setTimeout(() => {
          uni.reLaunch({
            url: '/pages/index/index'
          })
        }, 1500)
      }
    }
  } catch (error: any) {
    console.error('保存失败:', error)
    uni.showToast({
      title: error.message || '保存失败',
      icon: 'none'
    })
  } finally {
    isSubmitting.value = false
  }
}

// 取消编辑
const handleCancel = () => {
  uni.navigateBack()
}
</script>

<style lang="scss" scoped>
@import '@/styles/colors.scss';

// ===== 页面布局 =====
.page {
  min-height: 100vh;
  background: $gradient-bg-light;
  padding-top: 20rpx;
  padding-bottom: 120rpx; // 为底部按钮预留空间
}

:deep(.wd-form) {
  padding: 0;
}

// ===== 表单分组 =====
:deep(.wd-cell-group) {
  background: $color-bg-primary;
  border: 1rpx solid $color-border-primary;
  border-radius: $radius-lg;
  margin: 0 16rpx 24rpx;
  overflow: hidden;
  box-shadow: $shadow-sm;

  &:first-of-type {
    margin-top: 12rpx;
  }
}

// ===== 分组标题 =====
:deep(.wd-cell-group__title) {
  padding: 16rpx 24rpx 12rpx !important;
  font-size: 24rpx;
  font-weight: $font-weight-bold;
  color: $color-text-primary;
  background: linear-gradient(135deg, rgba(50, 220, 110, 0.05) 0%, $color-bg-primary 20%);
}

// ===== Cell 单元格 =====
:deep(.wd-cell) {
  padding: 16rpx 24rpx;
  background: $color-bg-primary;
  // border-bottom: 1rpx solid $color-border-primary;
  transition: background $transition-base;

  &:last-child {
    border-bottom: none;
  }

  &:active {
    background: $color-bg-secondary;
  }
}

:deep(.wd-cell__title) {
  font-size: 24rpx;
  color: $color-text-primary;
  font-weight: $font-weight-medium;
}

:deep(.wd-cell__value) {
  font-size: 24rpx;
  color: $color-text-secondary;
}

// ===== 输入框 =====
:deep(.wd-input) {
  padding: 16rpx 24rpx;
  background: $color-bg-primary;
  // border-bottom: 1rpx solid $color-border-primary;

  &:last-child {
    // border-bottom: none;
  }
}

:deep(.wd-input__label) {
  font-size: 24rpx;
  color: $color-text-primary;
  font-weight: $font-weight-medium;
}

:deep(.wd-input__control) {
  font-size: 26rpx;
  color: $color-text-secondary;

  &::placeholder {
    color: $color-text-tertiary;
  }
}

// ===== 按钮 =====
:deep(.wd-button) {
  font-size: 28rpx;
  font-weight: $font-weight-medium;
  border-radius: $radius-md;
  height: 88rpx;
  transition: all $transition-base;

  &.is-primary {
    background: $color-primary;
    color: white;
    box-shadow: $shadow-primary-md;

    &:active {
      background: darken($color-primary, 10%);
      transform: scale(0.98);
    }
  }

  &.is-plain {
    border: 2rpx solid $color-border-primary;
    background: $color-bg-primary;
    color: $color-text-primary;

    &:active {
      background: $color-bg-secondary;
      transform: scale(0.98);
    }
  }
}

// ===== 按钮容器 =====
.button-container {
  display: flex;
  flex-direction: column;
  gap: 12rpx;
  padding: 24rpx;
  background: $color-bg-primary;
  border-top: 1rpx solid $color-border-primary;
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  z-index: 100;
  box-shadow: 0 -2rpx 8rpx rgba(0, 0, 0, 0.05);

  :deep(.wd-button) {
    width: 100%;
  }
}

// ===== 头像上传区域 =====
.avatar-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16rpx;
  padding: 24rpx 0;
}

.avatar-preview {
  width: 160rpx;
  height: 160rpx;
  border-radius: $radius-full;
  overflow: hidden;
  box-shadow: $shadow-md;
  background: $color-bg-secondary;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2rpx solid $color-border-primary;

  image {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}

:deep(.wd-button) {
  &.avatar-btn {
    width: 200rpx;
    height: 64rpx;
    font-size: 24rpx;
    border-radius: $radius-md;
    background: $color-primary;
    color: white;
    box-shadow: $shadow-primary-sm;

    &:active {
      transform: scale(0.96);
      background: darken($color-primary, 10%);
    }
  }
}

// ===== 广播框和单选框 =====
:deep(.wd-radio-group) {
  display: flex;
  gap: 24rpx;
  flex-wrap: wrap;
}

:deep(.wd-radio) {
  font-size: 24rpx;
  color: $color-text-primary;

  &.is-checked {
    color: $color-primary;
  }
}

// ===== 日期选择器 =====
:deep(.wd-datetime-picker) {
  padding: 16rpx 24rpx;
  background: $color-bg-primary;
  border-bottom: 1rpx solid $color-border-primary;

  &:last-child {
    border-bottom: none;
  }
}

:deep(.wd-datetime-picker__label) {
  font-size: 24rpx;
  color: $color-text-primary;
  font-weight: $font-weight-medium;
}

:deep(.wd-datetime-picker__placeholder) {
  color: $color-text-tertiary;
}
</style>