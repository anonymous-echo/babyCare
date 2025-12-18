package persistence

import (
	"context"

	"gorm.io/gorm"

	"github.com/wxlbd/nutri-baby-server/internal/domain/entity"
	"github.com/wxlbd/nutri-baby-server/internal/domain/repository"
	"github.com/wxlbd/nutri-baby-server/pkg/errors"
)

// sleepRecordRepositoryImpl 睡眠记录仓储实现
type sleepRecordRepositoryImpl struct {
	db *gorm.DB
}

// NewSleepRecordRepository 创建睡眠记录仓储
func NewSleepRecordRepository(db *gorm.DB) repository.SleepRecordRepository {
	return &sleepRecordRepositoryImpl{db: db}
}

func (r *sleepRecordRepositoryImpl) GetDailyStats(ctx context.Context, babyID int64, startDate, endDate int64) ([]*entity.DailySleepItem, error) {
	var records []*entity.DailySleepItem
	query := r.db.WithContext(ctx).
		Model(&entity.SleepRecord{}).
		Select(`
			to_char(to_timestamp(start_time / 1000), 'YYYY-MM-DD') AS date,
			COALESCE(SUM(duration), 0) AS total_duration,
			COUNT(*) AS total_count`).
		Where("baby_id = ? AND start_time BETWEEN ? AND ?", babyID, startDate, endDate).
		Group("date").
		Order("date ASC")

	if err := query.Find(&records).Error; err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "failed to get daily stats", err)
	}

	return records, nil
}

func (r *sleepRecordRepositoryImpl) Create(ctx context.Context, record *entity.SleepRecord) error {
	if err := r.db.WithContext(ctx).Create(record).Error; err != nil {
		return errors.Wrap(errors.DatabaseError, "failed to create sleep record", err)
	}
	return nil
}

func (r *sleepRecordRepositoryImpl) FindByID(ctx context.Context, recordID int64) (*entity.SleepRecord, error) {
	var record entity.SleepRecord
	err := r.db.WithContext(ctx).
		Where("id = ?", recordID).
		First(&record).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.ErrRecordNotFound
	}
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "failed to find sleep record", err)
	}

	return &record, nil
}

func (r *sleepRecordRepositoryImpl) FindByBabyID(
	ctx context.Context,
	babyID int64,
	startTime, endTime int64,
	page, pageSize int,
) ([]*entity.SleepRecord, int64, error) {
	var records []*entity.SleepRecord
	var total int64

	query := r.db.WithContext(ctx).
		Model(&entity.SleepRecord{}).
		Where("baby_id = ?", babyID)

	if startTime > 0 {
		query = query.Where("start_time >= ?", startTime)
	}
	if endTime > 0 {
		query = query.Where("start_time <= ?", endTime)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.Wrap(errors.DatabaseError, "failed to count sleep records", err)
	}

	offset := (page - 1) * pageSize
	err := query.
		Order("start_time DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&records).Error

	if err != nil {
		return nil, 0, errors.Wrap(errors.DatabaseError, "failed to find sleep records", err)
	}

	return records, total, nil
}

func (r *sleepRecordRepositoryImpl) Update(ctx context.Context, record *entity.SleepRecord) error {
	err := r.db.WithContext(ctx).
		Model(&entity.SleepRecord{}).
		Where("id = ?", record.ID).
		Updates(record).Error

	if err != nil {
		return errors.Wrap(errors.DatabaseError, "failed to update sleep record", err)
	}

	return nil
}

func (r *sleepRecordRepositoryImpl) Delete(ctx context.Context, recordID int64) error {
	err := r.db.WithContext(ctx).
		Where("id = ?", recordID).
		Delete(&entity.SleepRecord{}).Error

	if err != nil {
		return errors.Wrap(errors.DatabaseError, "failed to delete sleep record", err)
	}

	return nil
}

func (r *sleepRecordRepositoryImpl) FindUpdatedAfter(
	ctx context.Context,
	babyID int64,
	timestamp int64,
) ([]*entity.SleepRecord, error) {
	var records []*entity.SleepRecord

	err := r.db.WithContext(ctx).
		Where("baby_id = ? AND updated_at > ?", babyID, timestamp).
		Find(&records).Error

	if err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "failed to find updated sleep records", err)
	}

	return records, nil
}

func (r *sleepRecordRepositoryImpl) FindOngoingSleep(ctx context.Context, babyID int64) (*entity.SleepRecord, error) {
	var record entity.SleepRecord
	err := r.db.WithContext(ctx).
		Where("baby_id = ? AND end_time = 0", babyID).
		First(&record).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // 没有找到进行中的睡眠记录，返回nil而不是错误
	}
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "failed to find ongoing sleep", err)
	}

	return &record, nil
}
