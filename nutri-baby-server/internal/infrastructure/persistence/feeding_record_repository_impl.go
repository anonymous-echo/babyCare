package persistence

import (
	"context"

	"gorm.io/gorm"

	"github.com/wxlbd/nutri-baby-server/internal/domain/entity"
	"github.com/wxlbd/nutri-baby-server/internal/domain/repository"
	"github.com/wxlbd/nutri-baby-server/pkg/errors"
)

// feedingRecordRepositoryImpl 喂养记录仓储实现
type feedingRecordRepositoryImpl struct {
	db *gorm.DB
}

// FindLatestRecord 查询宝宝最新的一条喂养记录
func (r *feedingRecordRepositoryImpl) FindLatestRecord(ctx context.Context, babyID int64) (*entity.FeedingRecord, error) {
	var record entity.FeedingRecord
	err := r.db.WithContext(ctx).
		Where("baby_id = ?", babyID).
		Order("time DESC").
		First(&record).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.ErrRecordNotFound
		}
		return nil, errors.Wrap(errors.DatabaseError, "failed to get last record", err)
	}

	return &record, nil
}

// NewFeedingRecordRepository 创建喂养记录仓储
func NewFeedingRecordRepository(db *gorm.DB) repository.FeedingRecordRepository {
	return &feedingRecordRepositoryImpl{db: db}
}

func (r *feedingRecordRepositoryImpl) GetDailyStats(ctx context.Context, babyID int64, startDate, endDate int64) ([]*entity.DailyFeedingItem, error) {
	var records []*entity.DailyFeedingItem
	query := r.db.WithContext(ctx).
		Model(&entity.FeedingRecord{}).
		Select(`
        to_char(to_timestamp(time / 1000), 'YYYY-MM-DD') AS date,
        feeding_type,
        COUNT(*) AS total_count,
        COALESCE(SUM(amount), 0) AS total_amount,
        COALESCE(SUM(duration), 0) AS total_duration`).
		Where("baby_id = ? AND time BETWEEN ? AND ?", babyID, startDate, endDate).
		Group("date, feeding_type").
		Order("date ASC")

	if err := query.Scan(&records).Error; err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "failed to get daily stats", err)
	}

	return records, nil
}

func (r *feedingRecordRepositoryImpl) Create(ctx context.Context, record *entity.FeedingRecord) error {
	if err := r.db.WithContext(ctx).Create(record).Error; err != nil {
		return errors.Wrap(errors.DatabaseError, "failed to create feeding record", err)
	}
	return nil
}

func (r *feedingRecordRepositoryImpl) FindByID(ctx context.Context, recordID int64) (*entity.FeedingRecord, error) {
	var record entity.FeedingRecord
	err := r.db.WithContext(ctx).
		Where("id = ?", recordID).
		First(&record).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.ErrRecordNotFound
	}
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "failed to find feeding record", err)
	}

	return &record, nil
}

func (r *feedingRecordRepositoryImpl) FindByBabyID(
	ctx context.Context,
	babyID int64,
	startTime, endTime int64,
	page, pageSize int,
) ([]*entity.FeedingRecord, int64, error) {
	var records []*entity.FeedingRecord
	var total int64

	query := r.db.WithContext(ctx).
		Model(&entity.FeedingRecord{}).
		Where("baby_id = ?", babyID)

	// 如果是查询最近一条记录(page=1, pageSize=1),且没有时间范围限制,则只查询未提醒的记录
	// 这是为了支持定时任务查询未提醒的最近喂养记录
	if page == 1 && pageSize == 1 && startTime > 0 && endTime > 0 {
		query = query.Where("reminder_sent = ?", false)
	}

	if startTime > 0 {
		query = query.Where("time >= ?", startTime)
	}
	if endTime > 0 {
		query = query.Where("time <= ?", endTime)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.Wrap(errors.DatabaseError, "failed to count feeding records", err)
	}

	offset := (page - 1) * pageSize
	err := query.
		Order("time DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&records).Error

	if err != nil {
		return nil, 0, errors.Wrap(errors.DatabaseError, "failed to find feeding records", err)
	}

	return records, total, nil
}

// FindByBabyIDAndType 根据宝宝ID和喂养类型查找记录(分页)
func (r *feedingRecordRepositoryImpl) FindByBabyIDAndType(
	ctx context.Context,
	babyID int64,
	feedingType string,
	startTime, endTime int64,
	page, pageSize int,
) ([]*entity.FeedingRecord, int64, error) {
	var records []*entity.FeedingRecord
	var total int64

	query := r.db.WithContext(ctx).
		Model(&entity.FeedingRecord{}).
		Where("baby_id = ? AND feeding_type = ?", babyID, feedingType)

	if startTime > 0 {
		query = query.Where("time >= ?", startTime)
	}
	if endTime > 0 {
		query = query.Where("time <= ?", endTime)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.Wrap(errors.DatabaseError, "failed to count feeding records by type", err)
	}

	offset := (page - 1) * pageSize
	err := query.
		Order("time DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&records).Error

	if err != nil {
		return nil, 0, errors.Wrap(errors.DatabaseError, "failed to find feeding records by type", err)
	}

	return records, total, nil
}

func (r *feedingRecordRepositoryImpl) Update(ctx context.Context, record *entity.FeedingRecord) error {
	err := r.db.WithContext(ctx).
		Model(&entity.FeedingRecord{}).
		Where("id = ?", record.ID).
		Updates(record).Error

	if err != nil {
		return errors.Wrap(errors.DatabaseError, "failed to update feeding record", err)
	}

	return nil
}

func (r *feedingRecordRepositoryImpl) Delete(ctx context.Context, recordID int64) error {
	err := r.db.WithContext(ctx).
		Where("id = ?", recordID).
		Delete(&entity.FeedingRecord{}).Error

	if err != nil {
		return errors.Wrap(errors.DatabaseError, "failed to delete feeding record", err)
	}

	return nil
}

func (r *feedingRecordRepositoryImpl) FindUpdatedAfter(
	ctx context.Context,
	babyID int64,
	timestamp int64,
) ([]*entity.FeedingRecord, error) {
	var records []*entity.FeedingRecord

	err := r.db.WithContext(ctx).
		Where("baby_id = ? AND updated_at > ?", babyID, timestamp).
		Find(&records).Error

	if err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "failed to find updated feeding records", err)
	}

	return records, nil
}

// UpdateReminderStatus 更新提醒状态
func (r *feedingRecordRepositoryImpl) UpdateReminderStatus(
	ctx context.Context,
	recordID int64,
	sent bool,
	reminderTime int64,
) error {
	err := r.db.WithContext(ctx).
		Model(&entity.FeedingRecord{}).
		Where("id = ?", recordID).
		Updates(map[string]interface{}{
			"reminder_sent": sent,
			"reminder_time": reminderTime,
		}).Error

	if err != nil {
		return errors.Wrap(errors.DatabaseError, "failed to update reminder status", err)
	}

	return nil
}

// GetTodayStatsByType 获取今日按类型的统计数据
func (r *feedingRecordRepositoryImpl) GetTodayStatsByType(
	ctx context.Context,
	babyID int64,
	feedingType string,
	todayStart, todayEnd int64,
) (count int64, totalAmount float64, totalDuration int, err error) {
	type Result struct {
		Count         int64
		TotalAmount   float64
		TotalDuration int64
	}

	var result Result
	err = r.db.WithContext(ctx).
		Model(&entity.FeedingRecord{}).
		Select("COUNT(*) as count, COALESCE(SUM(amount), 0) as total_amount, COALESCE(SUM(duration), 0) as total_duration").
		Where("baby_id = ? AND feeding_type = ? AND time >= ? AND time <= ?",
			babyID, feedingType, todayStart, todayEnd).
		Scan(&result).Error

	if err != nil {
		return 0, 0, 0, errors.Wrap(errors.DatabaseError, "failed to get stats by type", err)
	}

	return result.Count, result.TotalAmount, int(result.TotalDuration), nil
}
