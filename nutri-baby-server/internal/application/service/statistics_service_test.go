package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wxlbd/nutri-baby-server/internal/domain/entity"
	"github.com/wxlbd/nutri-baby-server/internal/domain/repository"
	"go.uber.org/zap"
)

// MockFeedingRecordRepository is a mock implementation of repository.FeedingRecordRepository
type MockFeedingRecordRepository struct {
	mock.Mock
}

func (m *MockFeedingRecordRepository) Create(ctx context.Context, record *entity.FeedingRecord) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

func (m *MockFeedingRecordRepository) FindByID(ctx context.Context, recordID int64) (*entity.FeedingRecord, error) {
	args := m.Called(ctx, recordID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.FeedingRecord), args.Error(1)
}

func (m *MockFeedingRecordRepository) FindByBabyID(ctx context.Context, babyID int64, startTime, endTime int64, page, pageSize int) ([]*entity.FeedingRecord, int64, error) {
	args := m.Called(ctx, babyID, startTime, endTime, page, pageSize)
	return args.Get(0).([]*entity.FeedingRecord), args.Get(1).(int64), args.Error(2)
}

func (m *MockFeedingRecordRepository) FindByBabyIDAndType(ctx context.Context, babyID int64, feedingType string, startTime, endTime int64, page, pageSize int) ([]*entity.FeedingRecord, int64, error) {
	args := m.Called(ctx, babyID, feedingType, startTime, endTime, page, pageSize)
	return args.Get(0).([]*entity.FeedingRecord), args.Get(1).(int64), args.Error(2)
}

func (m *MockFeedingRecordRepository) Update(ctx context.Context, record *entity.FeedingRecord) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

func (m *MockFeedingRecordRepository) Delete(ctx context.Context, recordID int64) error {
	args := m.Called(ctx, recordID)
	return args.Error(0)
}

func (m *MockFeedingRecordRepository) FindUpdatedAfter(ctx context.Context, babyID int64, timestamp int64) ([]*entity.FeedingRecord, error) {
	args := m.Called(ctx, babyID, timestamp)
	return args.Get(0).([]*entity.FeedingRecord), args.Error(1)
}

func (m *MockFeedingRecordRepository) UpdateReminderStatus(ctx context.Context, recordID int64, sent bool, reminderTime int64) error {
	args := m.Called(ctx, recordID, sent, reminderTime)
	return args.Error(0)
}

func (m *MockFeedingRecordRepository) GetTodayStatsByType(ctx context.Context, babyID int64, feedingType string, todayStart, todayEnd int64) (count int64, totalAmount float64, totalDuration int, err error) {
	args := m.Called(ctx, babyID, feedingType, todayStart, todayEnd)
	return args.Get(0).(int64), args.Get(1).(float64), args.Get(2).(int), args.Error(3)
}

func (m *MockFeedingRecordRepository) FindLatestRecord(ctx context.Context, babyID int64) (*entity.FeedingRecord, error) {
	args := m.Called(ctx, babyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.FeedingRecord), args.Error(1)
}

func (m *MockFeedingRecordRepository) GetDailyStats(ctx context.Context, babyID int64, startDate, endDate int64) ([]*entity.DailyFeedingItem, error) {
	args := m.Called(ctx, babyID, startDate, endDate)
	return args.Get(0).([]*entity.DailyFeedingItem), args.Error(1)
}

// Ensure MockFeedingRecordRepository implements repository.FeedingRecordRepository
var _ repository.FeedingRecordRepository = (*MockFeedingRecordRepository)(nil)

func TestGetTodayFeedingStats_Accumulation(t *testing.T) {
	// Setup
	mockRepo := new(MockFeedingRecordRepository)
	logger := zap.NewNop()

	// Only mocking what's needed for this test
	service := NewStatisticsService(
		nil, // babyRepo
		nil, // collaboratorRepo
		mockRepo,
		nil, // sleepRecordRepo
		nil, // diaperRecordRepo
		nil, // growthRecordRepo
		nil, // userRepo
		logger,
	)

	ctx := context.Background()
	babyID := int64(123)
	now := time.Now()
	startTime := int64(1000)
	endTime := int64(2000)

	// Mock data: Two records for "bottle" feeding on the same day (simulating split records)
	dailyStats := []*entity.DailyFeedingItem{
		{
			Date:        now.Format("2006-01-02"),
			FeedingType: entity.FeedingTypeBottle,
			TotalCount:  2,
			TotalAmount: 100, // First batch
		},
		{
			Date:        now.Format("2006-01-02"),
			FeedingType: entity.FeedingTypeBottle,
			TotalCount:  3,
			TotalAmount: 150, // Second batch
		},
	}

	mockRepo.On("GetDailyStats", ctx, babyID, startTime, endTime).Return(dailyStats, nil)
	mockRepo.On("FindLatestRecord", ctx, babyID).Return((*entity.FeedingRecord)(nil), nil)

	// Execute
	// Note: getTodayFeedingStats is private, but we are in the same package so we can test it.
	// If it was not exported and we were in a different package (e.g. service_test), we would need to export it or test via public API.
	// Since the file is in package service, we can access it.
	stats, err := service.getTodayFeedingStats(ctx, babyID, startTime, endTime)

	// Verify
	assert.NoError(t, err)
	assert.NotNil(t, stats)

	// Expected: TotalCount = 2 + 3 = 5, BottleCount = 2 + 3 = 5, BottleMl = 100 + 150 = 250
	// Current Buggy Behavior: It likely overwrites, so it might be 3 and 150 (or 2 and 100 depending on order)

	assert.Equal(t, 5, stats.TotalCount, "TotalCount should be accumulated")
	assert.Equal(t, 5, stats.BottleCount, "BottleCount should be accumulated")
	assert.Equal(t, int64(250), stats.BottleMl, "BottleMl should be accumulated")
}
