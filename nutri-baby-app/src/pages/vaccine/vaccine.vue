<template>
  <view>
    <wd-navbar title="疫苗日程" left-text="返回" right-text="设置" left-arrow fixed placeholder safeAreaInsetTop>
      <template #capsule>
        <wd-navbar-capsule @back="goBack" @back-home="goBackHome" />
      </template>
    </wd-navbar>
    <view class="vaccine-page">
      <!-- 疫苗完成度 -->
      <view v-if="currentBaby" class="progress-card">
        <view class="card-header">
          <image src="/static/progress_activity.svg" class="header-icon" />
          <text class="header-title">疫苗接种进度</text>
        </view>

        <view class="progress-info">
          <view class="progress-bar-container">
            <view class="progress-bar">
              <view
                class="progress-fill"
                :style="{ width: completionStats.completionRate + '%' }"
              ></view>
            </view>
            <text class="progress-text">
              {{ completionStats.completed + completionStats.skipped }} /
              {{ completionStats.total }} ({{
                completionStats.completionRate
              }}%)
            </text>
          </view>
        </view>
      </view>

      <!-- 即将到期提醒 (基于待接种的日程计算) -->
      <view
        v-if="upcomingSchedules && upcomingSchedules.length > 0"
        class="reminders-section"
      >
        <view class="section-title">
          <image src="/static/recent.svg" class="section-icon" />
          近期待接种 ({{ upcomingSchedules.length }})
        </view>

        <view class="reminder-list">
          <view
            v-for="schedule in upcomingSchedules"
            :key="schedule.scheduleId"
            class="reminder-item"
          >
            <view class="reminder-content">
              <view class="vaccine-name">
                {{ schedule.vaccineName }} (第{{ schedule.doseNumber }}针)
              </view>
              <view class="vaccine-date">
                建议月龄: {{ schedule.ageInMonths }}个月
              </view>
            </view>
            <view class="reminder-action">
              <wd-button
                size="small"
                type="primary"
                @click.stop="handleRecordVaccine(schedule)"
              >
                记录接种
              </wd-button>
            </view>
          </view>
        </view>
      </view>

      <!-- 疫苗计划列表 -->
      <view class="plan-section">
        <view class="section-header">
          <view class="section-title">
            <image src="/static/calendar_month.svg" class="section-icon" />
            疫苗日程
          </view>
          <wd-button size="small" @click="showAddDialog = true">
            添加自定义计划
          </wd-button>
        </view>

        <wd-tabs v-model="activeTab">
          <wd-tab title="全部" pane-key="all" name="all" />
          <wd-tab title="已完成" pane-key="completed" name="completed" />
          <wd-tab title="未完成" pane-key="pending" name="pending" />
          <wd-tab title="已跳过" pane-key="skipped" name="skipped" />
        </wd-tabs>

        <view class="plan-list">
          <view
            v-for="schedule in vaccineSchedules"
            :key="schedule.scheduleId"
            class="plan-item"
            :class="{ completed: schedule.vaccinationStatus === 'completed' }"
          >
            <view class="plan-header">
              <view class="plan-name">
                <text class="required-badge" v-if="schedule.isRequired"
                  >必打</text
                >
                <text class="custom-badge" v-if="isCustomPlan(schedule)"
                  >自定义</text
                >
                {{ schedule.vaccineName }}
              </view>
              <view class="plan-header-right">
                <text class="plan-age">{{ schedule.ageInMonths }}个月</text>
                <view
                  class="plan-actions"
                  v-if="schedule.vaccinationStatus === 'pending'"
                >
                  <wd-button
                    size="small"
                    type="info"
                    @click.stop="handleEditPlan(schedule)"
                  >
                    编辑
                  </wd-button>
                  <wd-button
                    v-if="isCustomPlan(schedule)"
                    size="small"
                    type="danger"
                    @click.stop="handleDeletePlan(schedule)"
                  >
                    删除
                  </wd-button>
                </view>
              </view>
            </view>

            <view class="plan-detail">
              <text class="plan-dose">第{{ schedule.doseNumber }}针</text>
              <text v-if="schedule.description" class="plan-desc">{{
                schedule.description
              }}</text>
            </view>

            <view
              v-if="schedule.vaccinationStatus === 'completed'"
              class="plan-record"
            >
              <image src="/static/check-icon.svg" class="completed-icon" />
              <text class="completed-text">已接种</text>
              <text class="completed-date">
                {{ formatDate(schedule.vaccineDate || 0, "YYYY-MM-DD") }}
              </text>
              <text v-if="schedule.hospital" class="hospital-info">
                {{ schedule.hospital }}
              </text>
            </view>

            <view
              v-else-if="schedule.vaccinationStatus === 'skipped'"
              class="plan-record"
            >
              <image src="/static/skip-icon.svg" class="skipped-icon" />
              <text class="skipped-text">已跳过</text>
            </view>

            <view v-else class="plan-action">
              <wd-button
                size="small"
                type="primary"
                @click="handleRecordBySchedule(schedule)"
              >
                记录接种
              </wd-button>
              <wd-button
                size="small"
                type="info"
                @click="handleSkipSchedule(schedule)"
              >
                跳过
              </wd-button>
            </view>
          </view>

          <!-- 加载更多组件 -->
          <wd-loadmore
            :state="loadMoreState"
            @reload="loadMore"
            loading-text="加载中..."
            finished-text="没有更多了"
            error-text="加载失败，点击重试"
          />
        </view>
      </view>

      <!-- 接种记录对话框 -->
      <wd-popup
        v-model="showRecordDialog"
        position="bottom"
        custom-style="height: 75%"
        round
        closeable
      >
        <wd-cell-group title="记录疫苗接种" border>
          <wd-input
            v-model="recordForm.vaccineName"
            placeholder="疫苗名称"
            readonly
            label="疫苗名称"
          />
          <wd-datetime-picker
            v-model="recordForm.vaccineDate"
            type="date"
            label="接种日期"
          />
          <wd-input
            v-model="recordForm.hospital"
            placeholder="接种医院"
            label="接种医院*"
            required
          />
          <wd-input
            v-model="recordForm.batchNumber"
            placeholder="疫苗批号"
            label="疫苗批号"
          />
          <wd-input
            v-model="recordForm.doctor"
            placeholder="接种医生"
            label="接种医生"
          />
          <wd-textarea
            v-model="recordForm.reaction"
            placeholder="不良反应"
            label="不良反应"
            auto-height
          />
          <wd-textarea
            v-model="recordForm.note"
            placeholder="备注"
            label="备注"
            auto-height
          />
        </wd-cell-group>
        <view class="dialog-footer">
          <wd-button type="primary" size="large" @click="handleSaveRecord">
            保存
          </wd-button>
          <wd-button type="info" size="large" @click="showRecordDialog = false">
            取消
          </wd-button>
        </view>
      </wd-popup>

      <!-- 添加/编辑疫苗计划对话框 -->
      <wd-popup
        v-model="showAddDialog"
        position="bottom"
        :style="{ height: '80%' }"
        round
        closeable
      >
        <wd-cell-group :title="isEdit ? '编辑疫苗计划' : '添加疫苗计划'" border>
          <wd-form ref="formRef">
            <wd-input
              label="疫苗名称"
              v-model="planForm.vaccineName"
              placeholder="疫苗名称"
              required
            />
            <wd-input
              label="疫苗类型"
              v-model="planForm.vaccineType"
              placeholder="例如: HepB, BCG, DTaP"
              required
            />
            <wd-input
              type="number"
              label="接种月龄"
              v-model.number="planForm.ageInMonths"
              placeholder="接种月龄"
              required
            />
            <wd-input
              type="number"
              label="剂次"
              v-model.number="planForm.doseNumber"
              placeholder="剂次"
              required
            />
            <wd-input
              type="number"
              label="提醒天数"
              v-model.number="planForm.reminderDays"
              placeholder="提醒天数"
            />
            <wd-cell
              title="是否必打"
              title-width="100px"
              prop="switchVal"
              center
            >
              <view style="text-align: left; padding-left: 46rpx">
                <wd-switch v-model="planForm.isRequired" />
              </view>
            </wd-cell>
            <wd-textarea
              label="疫苗说明"
              v-model="planForm.description"
              placeholder="疫苗说明"
              :max-length="200"
            />
            <view class="dialog-footer">
              <wd-button
                type="primary"
                size="large"
                @click="handleSubmitPlan"
                block
              >
                {{ isEdit ? "保存" : "添加" }}
              </wd-button>
              <wd-button
                type="info"
                size="large"
                @click="handleCancelPlan"
                block
              >
                取消
              </wd-button>
            </view>
          </wd-form>
        </wd-cell-group>
      </wd-popup>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from "vue";
import { onReachBottom } from "@dcloudio/uni-app";
import { currentBaby, currentBabyId } from "@/store/baby";
import { userInfo } from "@/store/user";
import { formatDate } from "@/utils/date";
import { shouldShowGuide } from "@/store/subscribe";
import { goBack,goBackHome } from "@/utils/common";

// 直接调用 API 层 (使用新架构)
import * as vaccineApi from "@/api/vaccine";

// Tab状态
const activeTab = ref<"all" | "completed" | "pending" | "skipped">("all");

// 对话框状态
const showRecordDialog = ref(false);
const showAddDialog = ref(false);
const isEdit = ref(false);
const editPlanId = ref("");

// 订阅消息引导状态
const showVaccineGuide = ref(false);

// 疫苗日程数据 (新架构 - 合并计划和记录)
const vaccineSchedules = ref<vaccineApi.VaccineScheduleResponse[]>([]);
const vaccineStats = ref<{
  total: number;
  completed: number;
  pending: number;
  skipped: number;
  completionRate: number;
}>({
  total: 0,
  completed: 0,
  pending: 0,
  skipped: 0,
  completionRate: 0,
});

// 分页相关状态
const currentPage = ref(1);
const pageSize = ref(10);
const isLoadingMore = ref(false);
const hasMore = ref(true);
const totalSchedules = ref(0);

// 表单数据 (新架构)
const recordForm = ref({
  scheduleId: "",
  vaccineName: "",
  vaccineDate: Date.now(),
  hospital: "",
  batchNumber: "",
  doctor: "",
  reaction: "",
  note: "",
});

// 疫苗计划管理表单数据
const planForm = ref({
  vaccineName: "",
  vaccineType: "",
  ageInMonths: 0,
  doseNumber: 1,
  reminderDays: 7,
  isRequired: true,
  description: "",
});

// 加载疫苗数据 (新架构 - 支持分页)
const loadVaccineData = async (isRefresh: boolean = false) => {
  if (!currentBaby.value) return;

  // 防止重复加载
  if (isLoadingMore.value) {
    console.log("正在加载中，跳过重复请求");
    return;
  }

  // 如果不是刷新且没有更多数据，直接返回
  if (!isRefresh && !hasMore.value) {
    console.log("没有更多数据，跳过加载");
    return;
  }

  // 如果是刷新，重置分页
  if (isRefresh) {
    currentPage.value = 1;
    vaccineSchedules.value = [];
    hasMore.value = true;
  }

  const babyId = currentBaby.value.babyId;
  const pageToLoad = currentPage.value;

  console.log("加载疫苗日程数据", {
    babyId,
    page: pageToLoad,
    isRefresh,
    status: activeTab.value === "all" ? undefined : activeTab.value,
  });

  try {
    isLoadingMore.value = true;

    // 根据 tab 状态构建请求参数
    const status =
      activeTab.value === "completed"
        ? "completed"
        : activeTab.value === "pending"
        ? "pending"
        : activeTab.value === "skipped"
        ? "skipped"
        : undefined;

    // 并行加载日程列表和统计数据
    const [scheduleData, statsData] = await Promise.all([
      vaccineApi.apiFetchVaccineSchedules(babyId, {
        page: pageToLoad,
        pageSize: pageSize.value,
        status: status,
      }),
      // 无论过滤什么状态，统计信息总是返回全量数据
      vaccineApi.apiFetchVaccineScheduleStatistics(babyId),
    ]);

    // 如果是刷新，替换数据；否则追加数据
    if (isRefresh) {
      vaccineSchedules.value = scheduleData.schedules || [];
      // 刷新后，如果有数据，下次加载第2页
      if ((scheduleData.schedules || []).length > 0) {
        currentPage.value = 2;
      }
    } else {
      vaccineSchedules.value.push(...(scheduleData.schedules || []));
      // 加载更多后，如果有数据，页码递增
      if ((scheduleData.schedules || []).length > 0) {
        currentPage.value++;
      }
    }

    // 更新统计数据（全量统计，不受 tab 过滤影响）
    vaccineStats.value = statsData || {
      total: 0,
      completed: 0,
      pending: 0,
      skipped: 0,
      completionRate: 0,
    };

    // 使用 API 返回的 total 作为当前过滤条件下的总数（重要！）
    totalSchedules.value = scheduleData.total || 0;

    // 判断是否还有更多数据（基于当前过滤条件的总数）
    const loadedCount = vaccineSchedules.value.length;
    hasMore.value = loadedCount < totalSchedules.value;

    console.log("疫苗数据加载成功", {
      tab: activeTab.value,
      loadedPage: pageToLoad,
      nextPage: currentPage.value,
      loadedCount,
      totalInCurrentFilter: totalSchedules.value,
      totalAllSchedules: vaccineStats.value.total,
      hasMore: hasMore.value,
      stats: vaccineStats.value,
    });
  } catch (error) {
    console.error("加载疫苗数据失败:", error);
    // 确保即使出错也初始化为空数组
    if (isRefresh) {
      vaccineSchedules.value = [];
    }
    vaccineStats.value = {
      total: 0,
      completed: 0,
      pending: 0,
      skipped: 0,
      completionRate: 0,
    };

    uni.showToast({
      title: "加载数据失败",
      icon: "none",
    });
  } finally {
    isLoadingMore.value = false;
  }
};

// 完成度统计 - 直接使用后端返回的数据
const completionStats = computed(() => {
  if (!currentBaby.value) {
    return {
      total: 0,
      completed: 0,
      pending: 0,
      skipped: 0,
      completionRate: 0,
    };
  }

  return vaccineStats.value;
});

// 近期待接种的日程 (pending状态，按月龄排序，取前3个)
const upcomingSchedules = computed(() => {
  if (!currentBaby.value || !vaccineSchedules.value) {
    return [];
  }

  return vaccineSchedules.value
    .filter((s) => s.vaccinationStatus === "pending")
    .sort((a, b) => a.ageInMonths - b.ageInMonths)
    .slice(0, 3); // 只显示前3个
});

// 显示的日程列表（已按 tab 状态和分页加载，直接显示）
const displaySchedules = computed(() => {
  return (vaccineSchedules.value || []).sort(
    (a, b) => a.ageInMonths - b.ageInMonths
  );
});

// 加载更多的状态计算
const loadMoreState = computed<string>(() => {
  // 如果没有登录或没有选中宝宝，不显示加载状态
  if (!currentBaby.value) return "finished";

  // 如果记录为空，显示完成状态
  if (vaccineSchedules.value.length === 0) return "finished";

  // 根据是否还有更多数据和是否正在加载来返回状态
  if (isLoadingMore.value) return "loading";
  if (!hasMore.value) return "finished";

  // 默认状态
  return "loading";
});

// 监听 tab 变化，重新加载数据
watch(activeTab, async () => {
  console.log("Tab changed to:", activeTab.value);
  await loadVaccineData(true);
});

// 加载更多函数
const loadMore = () => {
  console.log("[Vaccine] loadMore 触发");
  if (hasMore.value && !isLoadingMore.value) {
    loadVaccineData(false);
  }
};

// 页面滚动到底部时触发
onReachBottom(() => {
  console.log("[Vaccine] onReachBottom 触发", {
    hasMore: hasMore.value,
    isLoadingMore: isLoadingMore.value,
  });

  // loadVaccineData 内部已经有防重复加载的逻辑
  loadVaccineData(false);
});

// 处理记录接种(通过日程)
const handleRecordVaccine = (schedule: vaccineApi.VaccineScheduleResponse) => {
  recordForm.value = {
    scheduleId: schedule.scheduleId,
    vaccineName: schedule.vaccineName,
    vaccineDate: Date.now(),
    hospital: "",
    batchNumber: "",
    doctor: "",
    reaction: "",
    note: "",
  };
  showRecordDialog.value = true;
};

// 处理记录接种(通过日程 - 别名方法)
const handleRecordBySchedule = (
  schedule: vaccineApi.VaccineScheduleResponse
) => {
  handleRecordVaccine(schedule);
};

// 处理跳过接种
const handleSkipSchedule = async (
  schedule: vaccineApi.VaccineScheduleResponse
) => {
  if (!currentBaby.value) return;

  uni.showModal({
    title: "跳过接种",
    content: `确定要跳过「${schedule.vaccineName}」吗？`,
    success: async (res) => {
      if (res.confirm) {
        try {
          await vaccineApi.apiUpdateVaccineSchedule(
            currentBaby.value!.babyId,
            schedule.scheduleId,
            {
              vaccinationStatus: "skipped",
            }
          );

          uni.showToast({
            title: "已标记为跳过",
            icon: "success",
          });

          // 重新加载数据
          await loadVaccineData(true);
        } catch (error: any) {
          uni.showToast({
            title: error.message || "操作失败",
            icon: "none",
          });
        }
      }
    },
  });
};

// 保存接种记录 (新架构)
const handleSaveRecord = async () => {
  console.log("handleSaveRecord", recordForm.value);
  if (!currentBaby.value || !userInfo.value) {
    uni.showToast({
      title: "请先登录",
      icon: "none",
    });
    return;
  }

  if (!recordForm.value.hospital.trim()) {
    uni.showToast({
      title: "请输入接种医院",
      icon: "none",
    });
    return;
  }

  // 保存前记录当前完成数
  const completedBefore = vaccineStats.value.completed;

  try {
    await vaccineApi.apiUpdateVaccineSchedule(
      currentBaby.value.babyId,
      recordForm.value.scheduleId,
      {
        vaccinationStatus: "completed",
        vaccineDate: recordForm.value.vaccineDate,
        hospital: recordForm.value.hospital.trim(),
        batchNumber: recordForm.value.batchNumber.trim() || undefined,
        doctor: recordForm.value.doctor.trim() || undefined,
        reaction: recordForm.value.reaction.trim() || undefined,
        note: recordForm.value.note.trim() || undefined,
      }
    );

    uni.showToast({
      title: "记录成功",
      icon: "success",
    });

    showRecordDialog.value = false;

    // 重新加载数据
    await loadVaccineData(true);

    // 检查是否是首次添加疫苗记录
    const isFirstRecord = completedBefore === 0;

    // 首次记录后,延迟显示订阅引导
    if (isFirstRecord && shouldShowGuide("vaccine_reminder")) {
      setTimeout(() => {
        showVaccineGuide.value = true;
      }, 1500); // 延迟1.5秒,让用户看到成功提示
    }
  } catch (error: any) {
    uni.showToast({
      title: error.message || "保存失败",
      icon: "none",
    });
  }
};

// 处理订阅消息结果
const handleSubscribeResult = (result: "accept" | "reject") => {
  if (result === "accept") {
    console.log("用户同意订阅疫苗提醒");
  }
};

// 页面加载
onMounted(async () => {
  if (!currentBaby.value) {
    uni.showToast({
      title: "请先选择宝宝",
      icon: "none",
    });
    setTimeout(() => {
      uni.navigateBack();
    }, 1500);
    return;
  }

  // 加载疫苗数据 (刷新模式，重置分页)
  await loadVaccineData(true);
});

// 判断是否为自定义计划
const isCustomPlan = (plan: vaccineApi.VaccineScheduleResponse): boolean => {
  return plan.isCustom;
};

// 编辑疫苗计划
const handleEditPlan = (plan: vaccineApi.VaccineScheduleResponse) => {
  // 检查是否已完成或已跳过
  if (plan.vaccinationStatus !== "pending") {
    uni.showToast({
      title: "只能编辑待接种状态的疫苗日程",
      icon: "none",
    });
    return;
  }

  // 填充表单数据
  isEdit.value = true;
  editPlanId.value = plan.scheduleId;
  planForm.value = {
    vaccineName: plan.vaccineName,
    vaccineType: plan.vaccineType,
    ageInMonths: plan.ageInMonths,
    doseNumber: plan.doseNumber,
    reminderDays: plan.reminderDays,
    isRequired: plan.isRequired,
    description: plan.description || "",
  };
  showAddDialog.value = true;
};

// 删除疫苗计划
const handleDeletePlan = (plan: vaccineApi.VaccineScheduleResponse) => {
  uni.showModal({
    title: "确认删除",
    content: `确定要删除"${plan.vaccineName}"吗?`,
    success: async (res) => {
      if (res.confirm) {
        if (!currentBaby.value) return;
        try {
          await vaccineApi.apiDeleteVaccineSchedule(
            currentBaby.value.babyId,
            plan.scheduleId
          );
          // 重新加载数据
          await loadVaccineData(true);
          uni.showToast({ title: "删除成功", icon: "success" });
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

// 提交疫苗计划表单（添加或编辑）
const handleSubmitPlan = async () => {
  // 验证表单
  if (!planForm.value.vaccineName.trim()) {
    uni.showToast({
      title: "请输入疫苗名称",
      icon: "none",
    });
    return;
  }

  if (!planForm.value.vaccineType.trim()) {
    uni.showToast({
      title: "请输入疫苗类型",
      icon: "none",
    });
    return;
  }

  if (!currentBabyId.value) {
    uni.showToast({
      title: "请先选择宝宝",
      icon: "none",
    });
    return;
  }

  try {
    if (isEdit.value) {
      // 编辑模式 - 更新基本信息
      await vaccineApi.apiUpdateScheduleInfo(
        currentBabyId.value,
        editPlanId.value,
        {
          vaccineType: planForm.value.vaccineType,
          vaccineName: planForm.value.vaccineName,
          description: planForm.value.description,
          ageInMonths: planForm.value.ageInMonths,
          doseNumber: planForm.value.doseNumber,
          isRequired: planForm.value.isRequired,
          reminderDays: planForm.value.reminderDays,
        }
      );
      uni.showToast({ title: "更新成功", icon: "success" });
    } else {
      // 创建模式 - 创建自定义计划
      await vaccineApi.apiCreateCustomSchedule(currentBabyId.value, {
        vaccineType: planForm.value.vaccineType,
        vaccineName: planForm.value.vaccineName,
        description: planForm.value.description,
        ageInMonths: planForm.value.ageInMonths,
        doseNumber: planForm.value.doseNumber,
        isRequired: planForm.value.isRequired,
        reminderDays: planForm.value.reminderDays,
      });
      uni.showToast({ title: "创建成功", icon: "success" });
    }

    // 创建/更新成功后重新加载数据
    await loadVaccineData(true);
    showAddDialog.value = false;
    resetPlanForm();
  } catch (error: any) {
    uni.showToast({
      title: error.message || (isEdit.value ? "更新失败" : "创建失败"),
      icon: "none",
    });
  }
};

// 取消编辑/添加疫苗计划
const handleCancelPlan = () => {
  showAddDialog.value = false;
  resetPlanForm();
};

// 重置疫苗计划表单
const resetPlanForm = () => {
  isEdit.value = false;
  editPlanId.value = "";
  planForm.value = {
    vaccineName: "",
    vaccineType: "",
    ageInMonths: 0,
    doseNumber: 1,
    reminderDays: 7,
    isRequired: true,
    description: "",
  };
};

// 监听对话框关闭，确保关闭时重置表单
watch(showAddDialog, (newVal) => {
  if (!newVal) {
    // 对话框关闭时重置表单
    resetPlanForm();
  }
});
</script>

<style lang="scss" scoped>
.vaccine-page {
  min-height: 100vh;
  background: #f6f8f7;
  padding: 20rpx;
  padding-bottom: 40rpx;
}

.progress-card {
  background: white;
  border: 1rpx solid #cae3d4;
  border-left: 4rpx solid #7dd3a2;
  border-radius: 16rpx;
  padding: 30rpx;
  margin-bottom: 20rpx;
  box-shadow: 0 2rpx 8rpx rgba(125, 211, 162, 0.08);
}

.card-header {
  display: flex;
  align-items: center;
  gap: 12rpx;
  margin-bottom: 20rpx;
}

.header-icon {
  width: 40rpx;
  height: 40rpx;
}

.header-title {
  font-size: 32rpx;
  font-weight: bold;
  color: #333;
}

.progress-bar-container {
  display: flex;
  flex-direction: column;
  gap: 12rpx;
}

.progress-bar {
  height: 16rpx;
  background: #f0f0f0;
  border-radius: 8rpx;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: #7dd3a2;
  transition: width 0.3s;
}

.progress-text {
  font-size: 28rpx;
  text-align: right;
  color: #333;
}

.reminders-section,
.plan-section {
  background: white;
  border: 1rpx solid #cae3d4;
  border-radius: 16rpx;
  padding: 30rpx;
  margin-bottom: 20rpx;
  box-shadow: 0 2rpx 8rpx rgba(125, 211, 162, 0.08);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20rpx;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 8rpx;
  font-size: 32rpx;
  font-weight: bold;
}

.section-icon {
  width: 34rpx;
  height: 34rpx;
}

.reminder-list,
.plan-list {
  display: flex;
  flex-direction: column;
  gap: 16rpx;
}

.reminder-item {
  display: flex;
  align-items: center;
  gap: 20rpx;
  padding: 20rpx;
  background: #f6f8f7;
  border-radius: 12rpx;
  border: 1rpx solid #cae3d4;
}

.reminder-content {
  flex: 1;
}

.vaccine-name {
  font-size: 28rpx;
  font-weight: bold;
  color: #1a1a1a;
  margin-bottom: 8rpx;
}

.vaccine-date {
  font-size: 24rpx;
  color: #666;
  margin-bottom: 8rpx;
}

.plan-item {
  padding: 24rpx;
  background: #f6f8f7;
  border: 1rpx solid #cae3d4;
  border-radius: 12rpx;

  &.completed {
    opacity: 0.6;
  }
}

.plan-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 12rpx;
  margin-bottom: 12rpx;
}

.plan-name {
  font-size: 28rpx;
  font-weight: bold;
  color: #1a1a1a;
  display: flex;
  align-items: flex-start;
  gap: 8rpx;
  flex: 1;
  flex-wrap: wrap;
  word-break: break-word;
}

.plan-header-right {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 12rpx;
  flex-shrink: 0;
}

.plan-actions {
  display: flex;
  gap: 8rpx;
}

.required-badge {
  display: inline-flex;
  padding: 4rpx 8rpx;
  background: #fa2c19;
  color: white;
  font-size: 20rpx;
  border-radius: 4rpx;
  margin-right: 8rpx;
  flex-shrink: 0;
  white-space: nowrap;
}

.custom-badge {
  display: inline-flex;
  padding: 4rpx 8rpx;
  background: #52c41a;
  color: white;
  font-size: 20rpx;
  border-radius: 4rpx;
  flex-shrink: 0;
  white-space: nowrap;
}

.plan-age {
  font-size: 24rpx;
  color: #7dd3a2;
  font-weight: bold;
}

.plan-detail {
  display: flex;
  gap: 20rpx;
  margin-bottom: 12rpx;
  font-size: 24rpx;
  color: #666;
}

.plan-record,
.plan-action {
  display: flex;
  align-items: center;
  gap: 12rpx;
  margin-top: 12rpx;
}

.completed-icon {
  width: 32rpx;
  height: 32rpx;
}

.completed-text {
  font-size: 26rpx;
  color: #52c41a;
  font-weight: bold;
}

.completed-date {
  font-size: 24rpx;
  color: #999;
}

.hospital-info {
  font-size: 22rpx;
  color: #999;
  margin-left: 8rpx;
}

.skipped-icon {
  width: 32rpx;
  height: 32rpx;
}

.skipped-text {
  font-size: 26rpx;
  color: #999;
  font-weight: bold;
}

.dialog-footer {
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  gap: 16rpx;
  padding: 20rpx 30rpx 30rpx 30rpx;
  background: #fff;
  border-top: 1rpx solid #f0f0f0;
  box-shadow: 0 -2rpx 8rpx rgba(0, 0, 0, 0.05);
}

.dialog-footer .wd-button {
  width: 100%;
}
</style>
