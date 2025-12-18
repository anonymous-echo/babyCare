package repository

import (
	"context"

	"github.com/wxlbd/nutri-baby-server/internal/domain/entity"
)

// FeedingRecordRepository 喂养记录仓储接口
type FeedingRecordRepository interface {
	// Create 创建记录
	Create(ctx context.Context, record *entity.FeedingRecord) error
	// FindByID 根据ID查找记录
	FindByID(ctx context.Context, recordID int64) (*entity.FeedingRecord, error)
	// FindByBabyID 查找宝宝的喂养记录(分页)
	FindByBabyID(ctx context.Context, babyID int64, startTime, endTime int64, page, pageSize int) ([]*entity.FeedingRecord, int64, error)
	// FindByBabyIDAndType 根据宝宝ID和喂养类型查找记录(分页)
	FindByBabyIDAndType(ctx context.Context, babyID int64, feedingType string, startTime, endTime int64, page, pageSize int) ([]*entity.FeedingRecord, int64, error)
	// Update 更新记录
	Update(ctx context.Context, record *entity.FeedingRecord) error
	// Delete 删除记录
	Delete(ctx context.Context, recordID int64) error
	// FindUpdatedAfter 查找指定时间后更新的记录(用于同步)
	FindUpdatedAfter(ctx context.Context, babyID int64, timestamp int64) ([]*entity.FeedingRecord, error)
	// UpdateReminderStatus 更新提醒状态
	UpdateReminderStatus(ctx context.Context, recordID int64, sent bool, reminderTime int64) error
	// GetTodayStatsByType 获取今日按类型的统计数据
	GetTodayStatsByType(ctx context.Context, babyID int64, feedingType string, todayStart, todayEnd int64) (count int64, totalAmount float64, totalDuration int, err error)
	// FindLatestRecord 查询宝宝最新的一条喂养记录
	FindLatestRecord(ctx context.Context, babyID int64) (*entity.FeedingRecord, error)
	// 获取指定时间范围的每日统计数据
	GetDailyStats(ctx context.Context, babyID int64, startDate, endDate int64) ([]*entity.DailyFeedingItem, error)
}

// SleepRecordRepository 睡眠记录仓储接口
type SleepRecordRepository interface {
	// Create 创建记录
	Create(ctx context.Context, record *entity.SleepRecord) error
	// FindByID 根据ID查找记录
	FindByID(ctx context.Context, recordID int64) (*entity.SleepRecord, error)
	// FindByBabyID 查找宝宝的睡眠记录(分页)
	FindByBabyID(ctx context.Context, babyID int64, startTime, endTime int64, page, pageSize int) ([]*entity.SleepRecord, int64, error)
	// Update 更新记录
	Update(ctx context.Context, record *entity.SleepRecord) error
	// Delete 删除记录
	Delete(ctx context.Context, recordID int64) error
	// FindUpdatedAfter 查找指定时间后更新的记录(用于同步)
	FindUpdatedAfter(ctx context.Context, babyID int64, timestamp int64) ([]*entity.SleepRecord, error)
	// FindOngoingSleep 查找进行中的睡眠记录
	FindOngoingSleep(ctx context.Context, babyID int64) (*entity.SleepRecord, error)
	GetDailyStats(ctx context.Context, babyID int64, startDate, endDate int64) ([]*entity.DailySleepItem, error)
}

// DiaperRecordRepository 换尿布记录仓储接口
type DiaperRecordRepository interface {
	// Create 创建记录
	Create(ctx context.Context, record *entity.DiaperRecord) error
	// FindByID 根据ID查找记录
	FindByID(ctx context.Context, recordID int64) (*entity.DiaperRecord, error)
	// FindByBabyID 查找宝宝的换尿布记录(分页)
	FindByBabyID(ctx context.Context, babyID int64, startTime, endTime int64, page, pageSize int) ([]*entity.DiaperRecord, int64, error)
	// Update 更新记录
	Update(ctx context.Context, record *entity.DiaperRecord) error
	// Delete 删除记录
	Delete(ctx context.Context, recordID int64) error
	// FindUpdatedAfter 查找指定时间后更新的记录(用于同步)
	FindUpdatedAfter(ctx context.Context, babyID int64, timestamp int64) ([]*entity.DiaperRecord, error)
	// GetDailyStats 获取指定时间范围的每日统计数据
	GetDailyStats(ctx context.Context, babyID int64, startDate, endDate int64) ([]*entity.DailyDiaperItem, error)
}

// GrowthRecordRepository 成长记录仓储接口
type GrowthRecordRepository interface {
	// Create 创建记录
	Create(ctx context.Context, record *entity.GrowthRecord) error
	// FindByID 根据ID查找记录
	FindByID(ctx context.Context, recordID int64) (*entity.GrowthRecord, error)
	// FindByBabyID 查找宝宝的成长记录(分页)
	FindByBabyID(ctx context.Context, babyID int64, startTime, endTime int64, page, pageSize int) ([]*entity.GrowthRecord, int64, error)
	// Update 更新记录
	Update(ctx context.Context, record *entity.GrowthRecord) error
	// Delete 删除记录
	Delete(ctx context.Context, recordID int64) error
	// FindUpdatedAfter 查找指定时间后更新的记录(用于同步)
	FindUpdatedAfter(ctx context.Context, babyID int64, timestamp int64) ([]*entity.GrowthRecord, error)
	// GetDailyStats 获取指定时间范围的每日统计数据
	GetDailyStats(ctx context.Context, babyID int64, startDate, endDate int64) ([]*entity.DailyGrowthItem, error)
}
