package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/wxlbd/nutri-baby-server/internal/application/dto"
	"github.com/wxlbd/nutri-baby-server/internal/domain/entity"
	"github.com/wxlbd/nutri-baby-server/internal/domain/repository"
	"github.com/wxlbd/nutri-baby-server/internal/infrastructure/config"
	"go.uber.org/zap"
)

// SchedulerService 定时任务服务
type SchedulerService struct {
	scheduler           *gocron.Scheduler
	vaccineScheduleRepo repository.BabyVaccineScheduleRepository // 新增: 疫苗接种日程仓储
	feedingRecordRepo   repository.FeedingRecordRepository
	userRepo            repository.UserRepository
	babyRepo            repository.BabyRepository             // 新增: 宝宝仓储
	collaboratorRepo    repository.BabyCollaboratorRepository // 协作者仓储
	subscribeService    *SubscribeService
	aiAnalysisService   AIAnalysisService // 新增: AI分析服务
	strategyFactory     *FeedingReminderStrategyFactory
	logger              *zap.Logger
}

// NewSchedulerService 创建定时任务服务
func NewSchedulerService(
	vaccineScheduleRepo repository.BabyVaccineScheduleRepository,
	feedingRecordRepo repository.FeedingRecordRepository,
	userRepo repository.UserRepository,
	babyRepo repository.BabyRepository,
	collaboratorRepo repository.BabyCollaboratorRepository, // 协作者仓储
	subscribeService *SubscribeService,
	aiAnalysisService AIAnalysisService, // 新增: AI分析服务
	cfg *config.Config,
	logger *zap.Logger,
) *SchedulerService {
	// 创建 gocron 调度器，使用本地时区
	scheduler := gocron.NewScheduler(time.Local)

	return &SchedulerService{
		scheduler:           scheduler,
		vaccineScheduleRepo: vaccineScheduleRepo,
		feedingRecordRepo:   feedingRecordRepo,
		userRepo:            userRepo,
		babyRepo:            babyRepo,
		collaboratorRepo:    collaboratorRepo,
		subscribeService:    subscribeService,
		aiAnalysisService:   aiAnalysisService,
		strategyFactory:     NewFeedingReminderStrategyFactory(cfg),
		logger:              logger,
	}
}

// Start 启动定时任务
func (s *SchedulerService) Start() {
	// 启动调度器(用于一次性定时任务)
	s.scheduler.StartAsync()

	// 新增: 每5分钟自动处理一次待分析的AI任务
	_, err := s.scheduler.Every(5).Minutes().Do(s.processAIAnalysisTasks)
	if err != nil {
		s.logger.Error("添加AI分析定时任务失败", zap.Error(err))
	} else {
		s.logger.Info("AI分析自动处理任务已启用 (每5分钟一次)")
	}

	// 新增: 每天凌晨 00:00 自动生成活跃用户的每日建议
	_, err = s.scheduler.Every(1).Day().At("00:00").Do(s.generateDailyTipsForActiveBabies)
	if err != nil {
		s.logger.Error("添加每日建议生成任务失败", zap.Error(err))
	} else {
		s.logger.Info("每日建议自动生成任务已启用 (每天 00:00)")
	}

	s.logger.Info("Scheduler service started with auto-processing enabled")
}

// Stop 停止定时任务
func (s *SchedulerService) Stop() {
	s.scheduler.Stop()
	s.logger.Info("Scheduler service stopped")
}

// processAIAnalysisTasks 处理待分析的AI任务（定时任务回调）
// 每5分钟自动调用一次，批量处理待处理的分析任务
func (s *SchedulerService) processAIAnalysisTasks() {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
	defer cancel()

	if err := s.aiAnalysisService.ProcessPendingAnalyses(ctx); err != nil {
		s.logger.Error("自动处理待分析AI任务失败", zap.Error(err))
		return
	}

	s.logger.Info("自动处理待分析AI任务成功")
}

// generateDailyTipsForActiveBabies 生成活跃用户的每日建议
func (s *SchedulerService) generateDailyTipsForActiveBabies() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	s.logger.Info("开始生成活跃用户的每日建议")

	// 1. 查找活跃宝宝（过去7天有登录的用户）
	activeSince := time.Now().AddDate(0, 0, -7).UnixMilli()
	activeBabies, err := s.babyRepo.FindActiveBabies(ctx, activeSince)
	if err != nil {
		s.logger.Error("查找活跃宝宝失败", zap.Error(err))
		return
	}

	s.logger.Info("找到活跃宝宝", zap.Int("count", len(activeBabies)))

	// 2. 为每个活跃宝宝生成建议
	for _, baby := range activeBabies {
		// 检查context是否取消
		select {
		case <-ctx.Done():
			s.logger.Warn("每日建议生成任务超时或取消")
			return
		default:
		}

		babyIDStr := strconv.FormatInt(baby.ID, 10)
		date := time.Now()

		// GenerateDailyTips 内部会检查是否已存在，如果已存在则直接返回
		// 如果不存在，则调用AI生成
		_, err := s.aiAnalysisService.GenerateDailyTips(ctx, babyIDStr, date)
		if err != nil {
			s.logger.Error("生成每日建议失败",
				zap.String("babyID", babyIDStr),
				zap.Error(err),
			)
		} else {
			s.logger.Info("生成每日建议成功", zap.String("babyID", babyIDStr))
		}

		// 稍微延时，避免瞬间请求过多，给其他请求留出资源
		time.Sleep(500 * time.Millisecond)
	}

	s.logger.Info("活跃用户每日建议生成完成")
}

// CheckVaccineReminders 检查疫苗提醒(使用新的 BabyVaccineSchedule 架构)
func (s *SchedulerService) CheckVaccineReminders() error {
	// ctx := context.Background()

	// TODO: 实现基于 BabyVaccineSchedule 的提醒逻辑
	// 1. 查询所有待接种的日程 (vaccination_status='pending')
	// 2. 根据 scheduled_date 和当前时间计算是否需要发送提醒
	// 3. 调用 GetReminderStatus() 方法获取提醒状态
	// 4. 发送订阅消息提醒
	// 5. 更新 reminder_sent 和 reminder_sent_at 字段

	s.logger.Info("CheckVaccineReminders 暂未实现(待迁移到新架构)")
	return nil
}

// AddFeedingReminderTask 添加喂养提醒一次性定时任务
//
// 在创建喂养记录成功后调用此方法，将在指定的 nextReminderTime 时间自动执行提醒
//
// 参数:
//   - ctx: 上下文
//   - record: 喂养记录实体，必须包含有效的 NextReminderTime
//
// 返回:
//   - jobTag: gocron 任务的标签，可用于后续取消任务
//   - err: 错误信息
func (s *SchedulerService) AddFeedingReminderTask(ctx context.Context, record *entity.FeedingRecord) (string, error) {
	// 检查是否设置了下次提醒时间
	if record.NextReminderTime == nil {
		s.logger.Debug("未设置下次提醒时间，跳过任务添加",
			zap.String("recordID", strconv.FormatInt(record.ID, 10)))
		return "", nil
	}

	// 计算执行时间
	executeTime := time.UnixMilli(*record.NextReminderTime)
	now := time.Now()

	// 如果执行时间已经过期，不添加任务
	if executeTime.Before(now) {
		s.logger.Warn("下次提醒时间已过期，跳过任务添加",
			zap.String("recordID", strconv.FormatInt(record.ID, 10)),
			zap.Time("executeTime", executeTime),
			zap.Time("now", now))
		return "", nil
	}

	// 创建提醒回调函数
	reminderJob := func() {
		// 创建新的上下文用于后台任务
		taskCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := s.executeFeedingReminder(taskCtx, record); err != nil {
			s.logger.Error("执行喂养提醒失败",
				zap.String("recordID", strconv.FormatInt(record.ID, 10)),
				zap.Error(err))
		}
	}

	// 创建任务标签用于识别和取消
	jobTag := fmt.Sprintf("feeding_reminder_%s", strconv.FormatInt(record.ID, 10))

	// 使用 gocron 的一次性任务 API
	// StartAt() 指定任务开始时间, LimitRunsTo(1) 限制只执行一次
	job, err := s.scheduler.Every(1).Second().
		StartAt(executeTime).
		LimitRunsTo(1).
		Tag(jobTag).
		Do(reminderJob)
	if err != nil {
		s.logger.Error("添加喂养提醒任务失败",
			zap.String("recordID", strconv.FormatInt(record.ID, 10)),
			zap.Error(err))
		return "", err
	}

	s.logger.Info("添加喂养提醒任务成功",
		zap.String("recordID", strconv.FormatInt(record.ID, 10)),
		zap.String("jobTag", jobTag),
		zap.String("jobName", job.GetName()),
		zap.Time("executeTime", executeTime),
		zap.Duration("delay", executeTime.Sub(now)))

	return jobTag, nil
}

// CancelFeedingReminderTask 取消喂养提醒任务
//
// 如果用户编辑了喂养记录或取消了提醒，可调用此方法取消已添加的任务
func (s *SchedulerService) CancelFeedingReminderTask(jobTag string) {
	err := s.scheduler.RemoveByTag(jobTag)
	if err != nil {
		s.logger.Warn("取消喂养提醒任务失败",
			zap.String("jobTag", jobTag),
			zap.Error(err))
	} else {
		s.logger.Info("喂养提醒任务已取消", zap.String("jobTag", jobTag))
	}
}

// executeFeedingReminder 执行喂养提醒逻辑
// 向宝宝的所有协作者发送喂养提醒消息
func (s *SchedulerService) executeFeedingReminder(ctx context.Context, record *entity.FeedingRecord) error {
	s.logger.Info("开始执行喂养提醒",
		zap.String("recordID", strconv.FormatInt(record.ID, 10)),
		zap.String("babyID", strconv.FormatInt(record.BabyID, 10)),
		zap.String("feedingType", record.FeedingType))

	// 1. 获取宝宝的所有协作者
	collaborators, err := s.collaboratorRepo.FindByBabyID(ctx, record.BabyID)
	if err != nil {
		s.logger.Error("获取宝宝协作者列表失败",
			zap.Int64("babyID", record.BabyID),
			zap.Error(err))
		return err
	}

	if len(collaborators) == 0 {
		s.logger.Warn("宝宝没有协作者，跳过提醒",
			zap.Int64("babyID", record.BabyID))
		return nil
	}

	// 2. 根据喂养类型获取模板类型
	templateType := s.getTemplateType(record.FeedingType)
	if templateType == "" {
		s.logger.Warn("不支持的喂养类型，无法发送提醒",
			zap.String("feedingType", record.FeedingType))
		return nil
	}

	// 3. 构建提醒消息数据
	strategy, err := s.strategyFactory.GetStrategy(record)
	if err != nil {
		s.logger.Error("获取提醒策略失败", zap.Error(err))
		return err
	}

	lastFeedingTime := time.UnixMilli(record.Time)
	hoursSince := time.Since(lastFeedingTime).Hours()
	messageData := strategy.BuildMessageData(record, lastFeedingTime, hoursSince)

	// 4. 遍历所有协作者，分别发送消息
	var sentCount, failCount int
	for _, collaborator := range collaborators {
		// 跳过已过期的临时协作者
		if collaborator.IsExpired() {
			s.logger.Debug("跳过已过期的临时协作者",
				zap.Int64("userID", collaborator.UserID))
			continue
		}

		// 获取协作者的用户信息
		user, err := s.userRepo.FindByID(ctx, collaborator.UserID)
		if err != nil {
			s.logger.Warn("获取协作者用户信息失败",
				zap.Int64("userID", collaborator.UserID),
				zap.Error(err))
			failCount++
			continue
		}

		// 检查用户是否已授权此提醒
		hasAuth, err := s.subscribeService.CheckAuthorizationStatus(ctx, user.OpenID, templateType)
		if err != nil {
			s.logger.Warn("检查授权状态失败",
				zap.String("openID", user.OpenID),
				zap.Error(err))
			failCount++
			continue
		}

		if !hasAuth {
			s.logger.Debug("用户未授权此提醒，跳过发送",
				zap.String("templateType", templateType),
				zap.String("openID", user.OpenID))
			continue
		}

		// 发送微信订阅消息
		sendReq := &dto.SendMessageRequest{
			OpenID:     user.OpenID,
			TemplateID: strategy.GetTemplateID(),
			Data:       messageData,
			Page:       "pages/record/feeding/feeding",
		}

		if err := s.subscribeService.SendSubscribeMessage(ctx, sendReq); err != nil {
			s.logger.Warn("发送微信消息失败",
				zap.String("openID", user.OpenID),
				zap.Error(err))
			failCount++
			continue
		}

		sentCount++
		s.logger.Debug("向协作者发送提醒成功",
			zap.String("openID", user.OpenID),
			zap.String("role", collaborator.Role))
	}

	// 5. 标记提醒已发送
	now := time.Now().UnixMilli()
	record.ReminderSent = true
	record.ReminderTime = &now

	if err := s.feedingRecordRepo.Update(ctx, record); err != nil {
		s.logger.Error("更新记录状态失败", zap.Error(err))
		return err
	}

	s.logger.Info("喂养提醒发送完成",
		zap.String("recordID", strconv.FormatInt(record.ID, 10)),
		zap.String("templateType", templateType),
		zap.Int("sentCount", sentCount),
		zap.Int("failCount", failCount),
		zap.Int("totalCollaborators", len(collaborators)))

	return nil
}

// getTemplateType 根据喂养类型获取微信订阅消息模板类型
func (s *SchedulerService) getTemplateType(feedingType string) string {
	switch feedingType {
	case entity.FeedingTypeBreast:
		return "breast_feeding_reminder"
	case entity.FeedingTypeBottle:
		return "bottle_feeding_reminder"
	case entity.FeedingTypeFood:
		return "food_feeding_reminder"
	default:
		return ""
	}
}
