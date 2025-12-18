package service

import (
	"context"
	"strconv"
	"strings"

	"go.uber.org/zap"

	"github.com/wxlbd/nutri-baby-server/internal/application/dto"
	"github.com/wxlbd/nutri-baby-server/internal/domain/repository"
	"github.com/wxlbd/nutri-baby-server/pkg/errors"
)

// DailyStatsService 按日统计服务
type DailyStatsService struct {
	*BaseRecordService
	feedingRecordRepo repository.FeedingRecordRepository
	sleepRecordRepo   repository.SleepRecordRepository
	diaperRecordRepo  repository.DiaperRecordRepository
	growthRecordRepo  repository.GrowthRecordRepository
}

// NewDailyStatsService 创建按日统计服务
func NewDailyStatsService(
	babyRepo repository.BabyRepository,
	collaboratorRepo repository.BabyCollaboratorRepository,
	userRepo repository.UserRepository,
	feedingRecordRepo repository.FeedingRecordRepository,
	sleepRecordRepo repository.SleepRecordRepository,
	diaperRecordRepo repository.DiaperRecordRepository,
	growthRecordRepo repository.GrowthRecordRepository,
	logger *zap.Logger,
) *DailyStatsService {
	return &DailyStatsService{
		BaseRecordService: NewBaseRecordService(babyRepo, collaboratorRepo, userRepo, logger),
		feedingRecordRepo: feedingRecordRepo,
		sleepRecordRepo:   sleepRecordRepo,
		diaperRecordRepo:  diaperRecordRepo,
		growthRecordRepo:  growthRecordRepo,
	}
}

// GetDailyStats 获取按日统计数据
func (s *DailyStatsService) GetDailyStats(ctx context.Context, openID string, req *dto.DailyStatsRequest) (*dto.DailyStatsResponse, error) {
	// 验证宝宝访问权限
	if err := s.CheckBabyAccess(ctx, req.BabyID, openID); err != nil {
		return nil, err
	}

	// 转换babyID from string to int64
	babyIDInt64, err := strconv.ParseInt(req.BabyID, 10, 64)
	if err != nil {
		return nil, errors.New(errors.ParamError, "无效的宝宝ID格式")
	}

	// 解析统计类型
	types := parseStatsTypes(req.Types)

	response := &dto.DailyStatsResponse{}

	// 获取喂养统计
	if contains(types, "feeding") {
		feedingStats, err := s.getFeedingDailyStats(ctx, babyIDInt64, req.StartDate, req.EndDate)
		if err != nil {
			s.logger.Error("获取喂养按日统计失败", zap.Error(err))
			return nil, err
		}
		response.Feeding = feedingStats
	}

	// 获取睡眠统计
	if contains(types, "sleep") {
		sleepStats, err := s.getSleepDailyStats(ctx, babyIDInt64, req.StartDate, req.EndDate)
		if err != nil {
			s.logger.Error("获取睡眠按日统计失败", zap.Error(err))
			return nil, err
		}
		response.Sleep = sleepStats
	}

	// 获取排泄统计
	if contains(types, "diaper") {
		diaperStats, err := s.getDiaperDailyStats(ctx, babyIDInt64, req.StartDate, req.EndDate)
		if err != nil {
			s.logger.Error("获取排泄按日统计失败", zap.Error(err))
			return nil, err
		}
		response.Diaper = diaperStats
	}

	// 获取成长统计
	if contains(types, "growth") {
		growthStats, err := s.getGrowthDailyStats(ctx, babyIDInt64, req.StartDate, req.EndDate)
		if err != nil {
			s.logger.Error("获取成长按日统计失败", zap.Error(err))
			return nil, err
		}
		response.Growth = growthStats
	}

	return response, nil
}

// getFeedingDailyStats 获取喂养按日统计
func (s *DailyStatsService) getFeedingDailyStats(ctx context.Context, babyID int64, startDate, endDate int64) ([]*dto.DailyFeedingStatsItem, error) {
	records, err := s.feedingRecordRepo.GetDailyStats(ctx, babyID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.DailyFeedingStatsItem, 0, len(records))
	for _, record := range records {
		result = append(result, &dto.DailyFeedingStatsItem{
			Date:          record.Date,
			FeedingType:   record.FeedingType,
			TotalCount:    record.TotalCount,
			TotalAmount:   record.TotalAmount,
			TotalDuration: record.TotalDuration,
		})
	}

	return result, nil
}

// getSleepDailyStats 获取睡眠按日统计
func (s *DailyStatsService) getSleepDailyStats(ctx context.Context, babyID int64, startDate, endDate int64) ([]*dto.DailySleepStatsItem, error) {
	records, err := s.sleepRecordRepo.GetDailyStats(ctx, babyID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.DailySleepStatsItem, 0, len(records))
	for _, record := range records {
		result = append(result, &dto.DailySleepStatsItem{
			Date:          record.Date,
			TotalDuration: record.TotalDuration,
			TotalCount:    record.TotalCount,
		})
	}

	return result, nil
}

// getDiaperDailyStats 获取排泄按日统计
func (s *DailyStatsService) getDiaperDailyStats(ctx context.Context, babyID int64, startDate, endDate int64) ([]*dto.DailyDiaperStatsItem, error) {
	records, err := s.diaperRecordRepo.GetDailyStats(ctx, babyID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.DailyDiaperStatsItem, 0, len(records))
	for _, record := range records {
		result = append(result, &dto.DailyDiaperStatsItem{
			Date:       record.Date,
			DiaperType: record.DiaperType,
			TotalCount: record.TotalCount,
		})
	}

	return result, nil
}

// getGrowthDailyStats 获取成长按日统计
func (s *DailyStatsService) getGrowthDailyStats(ctx context.Context, babyID int64, startDate, endDate int64) ([]*dto.DailyGrowthStatsItem, error) {
	records, err := s.growthRecordRepo.GetDailyStats(ctx, babyID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.DailyGrowthStatsItem, 0, len(records))
	for _, record := range records {
		result = append(result, &dto.DailyGrowthStatsItem{
			Date:                    record.Date,
			LatestHeight:            record.LatestHeight,
			LatestWeight:            record.LatestWeight,
			LatestHeadCircumference: record.LatestHeadCircumference,
			RecordCount:             record.RecordCount,
		})
	}

	return result, nil
}

// parseStatsTypes 解析统计类型
func parseStatsTypes(types string) []string {
	if types == "" {
		return []string{"feeding", "sleep", "diaper", "growth"} // 默认全部
	}
	return strings.Split(strings.ReplaceAll(types, " ", ""), ",")
}

// contains 检查切片是否包含指定元素
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
