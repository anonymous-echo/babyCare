package service

import (
	"context"
	"strconv"
	"time"

	"github.com/wxlbd/nutri-baby-server/internal/application/dto"
	"github.com/wxlbd/nutri-baby-server/internal/domain/entity"
	"github.com/wxlbd/nutri-baby-server/internal/domain/repository"
	"github.com/wxlbd/nutri-baby-server/pkg/errors"
	"go.uber.org/zap"
)

// StatisticsService 统计服务
type StatisticsService struct {
	babyRepo          repository.BabyRepository
	collaboratorRepo  repository.BabyCollaboratorRepository
	feedingRecordRepo repository.FeedingRecordRepository
	sleepRecordRepo   repository.SleepRecordRepository
	diaperRecordRepo  repository.DiaperRecordRepository
	growthRecordRepo  repository.GrowthRecordRepository
	userRepo          repository.UserRepository
	logger            *zap.Logger
}

// NewStatisticsService 创建统计服务
func NewStatisticsService(
	babyRepo repository.BabyRepository,
	collaboratorRepo repository.BabyCollaboratorRepository,
	feedingRecordRepo repository.FeedingRecordRepository,
	sleepRecordRepo repository.SleepRecordRepository,
	diaperRecordRepo repository.DiaperRecordRepository,
	growthRecordRepo repository.GrowthRecordRepository,
	userRepo repository.UserRepository,
	logger *zap.Logger,
) *StatisticsService {
	return &StatisticsService{
		babyRepo:          babyRepo,
		collaboratorRepo:  collaboratorRepo,
		feedingRecordRepo: feedingRecordRepo,
		sleepRecordRepo:   sleepRecordRepo,
		diaperRecordRepo:  diaperRecordRepo,
		growthRecordRepo:  growthRecordRepo,
		userRepo:          userRepo,
		logger:            logger,
	}
}

// GetBabyStatistics 获取宝宝统计数据
func (s *StatisticsService) GetBabyStatistics(ctx context.Context, babyID, openID string) (*dto.BabyStatisticsResponse, error) {
	// 1. 将字符串 ID 转换为 int64
	babyIDInt64, err := strconv.ParseInt(babyID, 10, 64)
	if err != nil {
		return nil, errors.New(errors.ParamError, "无效的宝宝ID格式")
	}

	// 2. 验证权限（检查用户是否有权访问该宝宝的数据）
	_, err = s.babyRepo.FindByID(ctx, babyIDInt64)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "查询宝宝信息失败", err)
	}

	// 检查权限
	hasAccess, err := s.checkBabyAccess(ctx, babyIDInt64, openID)
	if err != nil {
		return nil, err
	}
	if !hasAccess {
		return nil, errors.New(errors.PermissionDenied, "没有权限访问该宝宝信息")
	}

	// 3. 获取今日统计
	now := time.Now()
	todayStart := getTodayStart(now)
	todayEnd := getTodayEnd(now)

	todayStats, err := s.getTodayStatistics(ctx, babyIDInt64, todayStart.Unix()*1000, todayEnd.Unix()*1000)
	if err != nil {
		s.logger.Error("获取今日统计失败", zap.String("babyId", babyID), zap.Error(err))
		return nil, err
	}

	// 4. 获取本周统计
	weekStart := getWeekStart(now)
	weekEnd := getTodayEnd(now)
	prevWeekStart := weekStart.AddDate(0, 0, -7)
	prevWeekEnd := weekStart.AddDate(0, 0, -1)

	weeklyStats, err := s.getWeeklyStatistics(ctx, babyIDInt64,
		weekStart.Unix()*1000, weekEnd.Unix()*1000,
		prevWeekStart.Unix()*1000, prevWeekEnd.Unix()*1000)
	if err != nil {
		s.logger.Error("获取本周统计失败", zap.String("babyId", babyID), zap.Error(err))
		return nil, err
	}

	return &dto.BabyStatisticsResponse{
		Today:  *todayStats,
		Weekly: *weeklyStats,
	}, nil
}

// getTodayStatistics 获取今日统计
func (s *StatisticsService) getTodayStatistics(ctx context.Context, babyID int64, startTime, endTime int64) (*dto.TodayStatistics, error) {
	stats := &dto.TodayStatistics{
		Feeding: dto.TodayFeedingStats{},
		Sleep:   dto.TodaySleepStats{},
		Diaper:  dto.TodayDiaperStats{},
		Growth:  dto.TodayGrowthStats{},
	}

	// 1. 喂养统计
	feedingStats, err := s.getTodayFeedingStats(ctx, babyID, startTime, endTime)
	if err != nil {
		return nil, err
	}
	stats.Feeding = *feedingStats

	// 2. 睡眠统计
	sleepStats, err := s.getTodaySleepStats(ctx, babyID, startTime, endTime)
	if err != nil {
		return nil, err
	}
	stats.Sleep = *sleepStats

	// 3. 换尿布统计
	diaperStats, err := s.getTodayDiaperStats(ctx, babyID, startTime, endTime)
	if err != nil {
		return nil, err
	}
	stats.Diaper = *diaperStats

	// 4. 成长统计（最新的成长记录）
	growthStats, err := s.getTodayGrowthStats(ctx, babyID)
	if err != nil {
		return nil, err
	}
	stats.Growth = *growthStats

	return stats, nil
}

// getTodayFeedingStats 获取今日喂养统计
func (s *StatisticsService) getTodayFeedingStats(ctx context.Context, babyID int64, startTime, endTime int64) (*dto.TodayFeedingStats, error) {
	stats := &dto.TodayFeedingStats{}

	// 获取今日所有喂养记录（用于统计今日的母乳、奶瓶等）
	todayRecords, err := s.feedingRecordRepo.GetDailyStats(ctx, babyID, startTime, endTime)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "查询喂养记录失败", err)
	}

	for _, record := range todayRecords {
		stats.TotalCount += int(record.TotalCount)
		switch record.FeedingType {
		case entity.FeedingTypeBreast:
			stats.BreastCount += int(record.TotalCount)
		case entity.FeedingTypeBottle:
			stats.BottleCount += int(record.TotalCount)
			stats.BottleMl += record.TotalAmount
		case entity.FeedingTypeFood:
			stats.FoodCount += int(record.TotalCount)
		}
	}

	// 查询该宝宝所有的喂养记录，获取最新的一条
	latestRecord, err := s.feedingRecordRepo.FindLatestRecord(ctx, babyID)
	if err != nil {
		s.logger.Info("查询最新喂养记录失败", zap.Int64("babyId", babyID), zap.Error(err))
	}
	if latestRecord != nil {
		stats.LastFeedingTime = &latestRecord.Time
	}
	return stats, nil
}

// getTodaySleepStats 获取今日睡眠统计
func (s *StatisticsService) getTodaySleepStats(ctx context.Context, babyID int64, startTime, endTime int64) (*dto.TodaySleepStats, error) {
	records, _, err := s.sleepRecordRepo.FindByBabyID(ctx, babyID, startTime, endTime, 1, 100)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "查询睡眠记录失败", err)
	}

	stats := &dto.TodaySleepStats{
		SessionCount: len(records),
	}

	for _, record := range records {
		if record.Duration != nil {
			// ⚠️ 注意：Duration 存储为秒，需要转换为分钟
			// 为避免整数除法精度丢失，使用向上取整
			// 例如：9秒 / 60 = 0分钟（整数除法），但应该返回最小1分钟
			// 所以采用向上取整的方式：(duration + 59) / 60
			minutes := int((*record.Duration + 59) / 60) // +59用于向上取整
			stats.TotalMinutes += minutes
		}
	}

	// 如果有多条记录，获取最后一条的时长
	if len(records) > 0 && records[len(records)-1].Duration != nil {
		stats.LastSleepMinutes = int((*records[len(records)-1].Duration + 59) / 60) // +59用于向上取整
	}

	return stats, nil
}

// getTodayDiaperStats 获取今日换尿布统计
func (s *StatisticsService) getTodayDiaperStats(ctx context.Context, babyID int64, startTime, endTime int64) (*dto.TodayDiaperStats, error) {
	records, _, err := s.diaperRecordRepo.FindByBabyID(ctx, babyID, startTime, endTime, 1, 100)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "查询换尿布记录失败", err)
	}

	stats := &dto.TodayDiaperStats{
		TotalCount: len(records),
	}

	for _, record := range records {
		if record.Type == "pee" || record.Type == "both" {
			stats.PeeCount++
		}
		if record.Type == "poop" || record.Type == "both" {
			stats.PoopCount++
		}
	}

	return stats, nil
}

// getTodayGrowthStats 获取今日成长统计（最新记录）
func (s *StatisticsService) getTodayGrowthStats(ctx context.Context, babyID int64) (*dto.TodayGrowthStats, error) {
	// 获取最新的成长记录
	records, _, err := s.growthRecordRepo.FindByBabyID(ctx, babyID, 0, 9999999999999, 1, 1)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "查询成长记录失败", err)
	}

	stats := &dto.TodayGrowthStats{}

	if len(records) > 0 {
		record := records[0]
		stats.LatestWeight = record.Weight
		stats.LatestHeight = record.Height
		stats.LatestHeadCircumference = record.HeadCircumference
	}

	return stats, nil
}

// getWeeklyStatistics 获取本周统计
func (s *StatisticsService) getWeeklyStatistics(ctx context.Context, babyID int64, weekStart, weekEnd, prevWeekStart, prevWeekEnd int64) (*dto.WeeklyStatistics, error) {
	stats := &dto.WeeklyStatistics{
		Feeding: dto.WeeklyFeedingStats{},
		Sleep:   dto.WeeklySleepStats{},
		Growth:  dto.WeeklyGrowthStats{},
	}

	// 1. 本周喂养统计和趋势
	feedingStats, err := s.getWeeklyFeedingStats(ctx, babyID, weekStart, weekEnd, prevWeekStart, prevWeekEnd)
	if err != nil {
		return nil, err
	}
	stats.Feeding = *feedingStats

	// 2. 本周睡眠统计和趋势
	sleepStats, err := s.getWeeklySleepStats(ctx, babyID, weekStart, weekEnd, prevWeekStart, prevWeekEnd)
	if err != nil {
		return nil, err
	}
	stats.Sleep = *sleepStats

	// 3. 本周成长统计
	growthStats, err := s.getWeeklyGrowthStats(ctx, babyID, weekStart, weekEnd)
	if err != nil {
		return nil, err
	}
	stats.Growth = *growthStats

	return stats, nil
}

// getWeeklyFeedingStats 获取本周喂养统计和趋势
func (s *StatisticsService) getWeeklyFeedingStats(ctx context.Context, babyID int64, weekStart, weekEnd, prevWeekStart, prevWeekEnd int64) (*dto.WeeklyFeedingStats, error) {
	// 本周喂养
	thisWeekRecords, _, err := s.feedingRecordRepo.FindByBabyID(ctx, babyID, weekStart, weekEnd, 1, 100)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "查询本周喂养记录失败", err)
	}

	// 上周喂养
	prevWeekRecords, _, err := s.feedingRecordRepo.FindByBabyID(ctx, babyID, prevWeekStart, prevWeekEnd, 1, 100)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "查询上周喂养记录失败", err)
	}

	thisWeekCount := len(thisWeekRecords)
	prevWeekCount := len(prevWeekRecords)
	trend := thisWeekCount - prevWeekCount

	avgPerDay := 0.0
	if thisWeekCount > 0 {
		avgPerDay = float64(thisWeekCount) / 7.0
	}

	return &dto.WeeklyFeedingStats{
		TotalCount: thisWeekCount,
		Trend:      trend,
		AvgPerDay:  avgPerDay,
	}, nil
}

// getWeeklySleepStats 获取本周睡眠统计和趋势
func (s *StatisticsService) getWeeklySleepStats(ctx context.Context, babyID int64, weekStart, weekEnd, prevWeekStart, prevWeekEnd int64) (*dto.WeeklySleepStats, error) {
	// 本周睡眠
	thisWeekRecords, _, err := s.sleepRecordRepo.FindByBabyID(ctx, babyID, weekStart, weekEnd, 1, 100)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "查询本周睡眠记录失败", err)
	}

	// 上周睡眠
	prevWeekRecords, _, err := s.sleepRecordRepo.FindByBabyID(ctx, babyID, prevWeekStart, prevWeekEnd, 1, 100)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "查询上周睡眠记录失败", err)
	}

	// 计算本周总睡眠分钟数（将秒转换为分钟，使用向上取整避免精度丢失）
	thisWeekMinutes := 0
	for _, record := range thisWeekRecords {
		if record.Duration != nil {
			// ⚠️ Duration 存储为秒，需要转换为分钟
			// 使用向上取整：(duration + 59) / 60
			minutes := int((*record.Duration + 59) / 60)
			thisWeekMinutes += minutes
		}
	}

	// 计算上周总睡眠分钟数（将秒转换为分钟，使用向上取整避免精度丢失）
	prevWeekMinutes := 0
	for _, record := range prevWeekRecords {
		if record.Duration != nil {
			// ⚠️ Duration 存储为秒，需要转换为分钟
			// 使用向上取整：(duration + 59) / 60
			minutes := int((*record.Duration + 59) / 60)
			prevWeekMinutes += minutes
		}
	}

	trend := float64(thisWeekMinutes - prevWeekMinutes)
	avgPerDay := 0.0
	if thisWeekMinutes > 0 {
		avgPerDay = float64(thisWeekMinutes) / 7.0
	}

	return &dto.WeeklySleepStats{
		TotalMinutes: thisWeekMinutes,
		Trend:        roundToOneDecimal(trend),
		AvgPerDay:    roundToOneDecimal(avgPerDay),
	}, nil
}

// getWeeklyGrowthStats 获取本周成长统计
func (s *StatisticsService) getWeeklyGrowthStats(ctx context.Context, babyID int64, weekStart, weekEnd int64) (*dto.WeeklyGrowthStats, error) {
	// 获取一周内的成长记录
	records, _, err := s.growthRecordRepo.FindByBabyID(ctx, babyID, weekStart, weekEnd, 1, 100)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "查询本周成长记录失败", err)
	}

	stats := &dto.WeeklyGrowthStats{
		WeightGain: 0,
		HeightGain: 0,
	}

	if len(records) < 2 {
		return stats, nil
	}

	// 因为数据库按 time DESC 排序，所以第一条是最新的，最后一条是最早的
	latest := records[0]
	oldest := records[len(records)-1]

	if latest.Weight != nil && oldest.Weight != nil {
		stats.WeightGain = roundToOneDecimal(*latest.Weight - *oldest.Weight)
	}

	if latest.Height != nil && oldest.Height != nil {
		stats.HeightGain = roundToOneDecimal(*latest.Height - *oldest.Height)
	}

	if oldest.Weight != nil {
		stats.WeekStartWeight = oldest.Weight
	}

	return stats, nil
}

// checkBabyAccess 检查用户是否有权访问宝宝
func (s *StatisticsService) checkBabyAccess(ctx context.Context, babyID int64, openID string) (bool, error) {
	baby, err := s.babyRepo.FindByID(ctx, babyID)
	if err != nil {
		return false, err
	}

	// 获取创建者用户信息
	user, err := s.userRepo.FindByID(ctx, baby.UserID)
	if err != nil {
		return false, err
	}

	// 检查是否是创建者
	if user.OpenID == openID {
		return true, nil
	}

	// 获取用户ID
	targetUser, err := s.userRepo.FindByOpenID(ctx, openID)
	if err != nil {
		return false, nil
	}

	// 检查是否是协作者
	collaborator, err := s.collaboratorRepo.CheckPermission(ctx, babyID, targetUser.ID)
	if err != nil {
		return false, nil
	}

	if collaborator == nil {
		return false, nil
	}

	// 检查权限是否未过期
	if collaborator.ExpiresAt != nil && time.UnixMilli(*collaborator.ExpiresAt).Before(time.Now()) {
		return false, nil
	}

	return true, nil
}

// Helper functions

// getTodayStart 获取今天的开始时间 (00:00:00)
func getTodayStart(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// getTodayEnd 获取今天的结束时间 (23:59:59)
func getTodayEnd(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 999999999, t.Location())
}

// getWeekStart 获取本周的开始时间 (7天前的今天00:00:00)
func getWeekStart(t time.Time) time.Time {
	sevenDaysAgo := t.AddDate(0, 0, -6)
	year, month, day := sevenDaysAgo.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, sevenDaysAgo.Location())
}

// roundToOneDecimal 四舍五入到小数点后一位
func roundToOneDecimal(f float64) float64 {
	return float64(int(f*10)) / 10.0
}
