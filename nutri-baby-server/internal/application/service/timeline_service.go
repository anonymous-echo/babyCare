package service

import (
	"context"
	"sort"
	"strconv"
	"sync"

	"go.uber.org/zap"

	"github.com/wxlbd/nutri-baby-server/internal/application/dto"
	"github.com/wxlbd/nutri-baby-server/internal/domain/repository"
)

// TimelineService 时间线服务
type TimelineService struct {
	*BaseRecordService
	feedingService *FeedingRecordService
	sleepService   *SleepRecordService
	diaperService  *DiaperRecordService
	growthService  *GrowthRecordService
}

// NewTimelineService 创建时间线服务
func NewTimelineService(
	babyRepo repository.BabyRepository,
	collaboratorRepo repository.BabyCollaboratorRepository,
	userRepo repository.UserRepository,
	feedingService *FeedingRecordService,
	sleepService *SleepRecordService,
	diaperService *DiaperRecordService,
	growthService *GrowthRecordService,
	logger *zap.Logger,
) *TimelineService {
	return &TimelineService{
		BaseRecordService: NewBaseRecordService(babyRepo, collaboratorRepo, userRepo, logger),
		feedingService:    feedingService,
		sleepService:      sleepService,
		diaperService:     diaperService,
		growthService:     growthService,
	}
}

// GetTimeline 获取时间线记录
func (s *TimelineService) GetTimeline(ctx context.Context, openID string, query *dto.TimelineQuery) (*dto.TimelineResponse, error) {
	// 检查权限
	if err := s.CheckBabyAccess(ctx, query.BabyID, openID); err != nil {
		return nil, err
	}

	// 设置分页参数（使用统一的默认值：pageSize 默认 10，最大 100）
	page := query.GetPageWithDefault()
	pageSize := query.GetPageSizeWithDefault()

	// 构建查询参数 (不分页,先获取所有数据)
	recordQuery := &dto.RecordListQuery{
		BabyID:    query.BabyID,
		StartTime: query.StartTime,
		EndTime:   query.EndTime,
	}
	// 设置记录查询的分页参数为最大值，确保获取所有数据
	pageVal := 1
	pageSizeVal := 1000
	recordQuery.Page = &pageVal
	recordQuery.PageSize = &pageSizeVal

	// 根据 recordType 决定查询哪些类型
	recordType := query.RecordType
	queryFeeding := recordType == "" || recordType == "feeding"
	querySleep := recordType == "" || recordType == "sleep"
	queryDiaper := recordType == "" || recordType == "diaper"
	queryGrowth := recordType == "" || recordType == "growth"

	// 计算需要查询的类型数量
	queryCount := 0
	if queryFeeding {
		queryCount++
	}
	if querySleep {
		queryCount++
	}
	if queryDiaper {
		queryCount++
	}
	if queryGrowth {
		queryCount++
	}

	// 并发查询所需类型的记录
	var (
		feedingRecords []dto.FeedingRecordDTO
		sleepRecords   []dto.SleepRecordDTO
		diaperRecords  []dto.DiaperRecordDTO
		growthRecords  []dto.GrowthRecordDTO
		wg             sync.WaitGroup
		mu             sync.Mutex
		errs           []error
	)

	wg.Add(queryCount)

	// 查询喂养记录
	if queryFeeding {
		go func() {
			defer wg.Done()
			records, _, err := s.feedingService.GetFeedingRecords(ctx, openID, recordQuery)
			if err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
				s.logger.Warn("获取喂养记录失败", zap.Error(err))
				return
			}
			feedingRecords = records
		}()
	}

	// 查询睡眠记录
	if querySleep {
		go func() {
			defer wg.Done()
			records, _, err := s.sleepService.GetSleepRecords(ctx, openID, recordQuery)
			if err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
				s.logger.Warn("获取睡眠记录失败", zap.Error(err))
				return
			}
			sleepRecords = records
		}()
	}

	// 查询排泄记录
	if queryDiaper {
		go func() {
			defer wg.Done()
			records, _, err := s.diaperService.GetDiaperRecords(ctx, openID, recordQuery)
			if err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
				s.logger.Warn("获取排泄记录失败", zap.Error(err))
				return
			}
			diaperRecords = records
		}()
	}

	// 查询成长记录
	if queryGrowth {
		go func() {
			defer wg.Done()
			records, _, err := s.growthService.GetGrowthRecords(ctx, openID, recordQuery)
			if err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
				s.logger.Warn("获取成长记录失败", zap.Error(err))
				return
			}
			growthRecords = records
		}()
	}

	wg.Wait()

	// 如果所有查询都失败,返回错误
	if len(errs) == queryCount {
		return nil, errs[0]
	}

	// 聚合所有记录为 TimelineItem
	var items []dto.TimelineItem

	// 转换喂养记录
	for _, record := range feedingRecords {
		item := dto.TimelineItem{
			RecordType: "feeding",
			RecordID:   record.RecordID,
			BabyID:     record.BabyID,
			EventTime:  record.FeedingTime,
			Detail:     record,
			CreateBy:   record.CreateBy,
			CreateTime: record.CreateTime,
		}
		s.enrichTimelineItem(ctx, &item)
		items = append(items, item)
	}

	// 转换睡眠记录
	for _, record := range sleepRecords {
		item := dto.TimelineItem{
			RecordType: "sleep",
			RecordID:   record.RecordID,
			BabyID:     record.BabyID,
			EventTime:  record.StartTime,
			Detail:     record,
			CreateBy:   record.CreateBy,
			CreateTime: record.CreateTime,
		}
		s.enrichTimelineItem(ctx, &item)
		items = append(items, item)
	}

	// 转换排泄记录
	for _, record := range diaperRecords {
		item := dto.TimelineItem{
			RecordType: "diaper",
			RecordID:   record.RecordID,
			BabyID:     record.BabyID,
			EventTime:  record.ChangeTime,
			Detail:     record,
			CreateBy:   record.CreateBy,
			CreateTime: record.CreateTime,
		}
		s.enrichTimelineItem(ctx, &item)
		items = append(items, item)
	}

	// 转换成长记录
	for _, record := range growthRecords {
		item := dto.TimelineItem{
			RecordType: "growth",
			RecordID:   record.RecordID,
			BabyID:     record.BabyID,
			EventTime:  record.MeasureTime,
			Detail:     record,
			CreateBy:   record.CreateBy,
			CreateTime: record.CreateTime,
		}
		s.enrichTimelineItem(ctx, &item)
		items = append(items, item)
	}

	// 按 eventTime 倒序排序 (最新的在前面)
	sort.Slice(items, func(i, j int) bool {
		return items[i].EventTime > items[j].EventTime
	})

	// 计算总数
	total := int64(len(items))

	// 内存分页
	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= len(items) {
		items = []dto.TimelineItem{}
	} else {
		if end > len(items) {
			end = len(items)
		}
		items = items[start:end]
	}

	return &dto.TimelineResponse{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// enrichTimelineItem 丰富时间线项的创建者信息
func (s *TimelineService) enrichTimelineItem(ctx context.Context, item *dto.TimelineItem) {
	if item.CreateBy == "" {
		return
	}

	// 将 CreateBy 从 string 转换为 int64 (用户ID)
	userID, err := strconv.ParseInt(item.CreateBy, 10, 64)
	if err != nil {
		s.logger.Warn("解析创建者ID失败", zap.String("createBy", item.CreateBy), zap.Error(err))
		return
	}

	// 获取创建者用户信息
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		s.logger.Warn("获取创建者用户信息失败", zap.Int64("userId", userID), zap.Error(err))
		return
	}

	// 设置创建者昵称
	item.CreateName = user.NickName

	// 获取创建者与宝宝的关系
	babyIDInt64, err := strconv.ParseInt(item.BabyID, 10, 64)
	if err != nil {
		s.logger.Warn("解析宝宝ID失败", zap.String("babyId", item.BabyID), zap.Error(err))
		return
	}

	collaborator, err := s.collaboratorRepo.FindByBabyAndUser(ctx, babyIDInt64, userID)
	if err != nil {
		s.logger.Warn("获取创建者关系失败", zap.Int64("userId", userID), zap.String("babyId", item.BabyID), zap.Error(err))
		return
	}

	if collaborator != nil {
		item.Relationship = collaborator.Relationship
	}
}
