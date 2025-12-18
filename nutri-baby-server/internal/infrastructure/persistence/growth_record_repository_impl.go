package persistence

import (
	"context"

	"gorm.io/gorm"

	"github.com/wxlbd/nutri-baby-server/internal/domain/entity"
	"github.com/wxlbd/nutri-baby-server/internal/domain/repository"
	"github.com/wxlbd/nutri-baby-server/pkg/errors"
)

// growthRecordRepositoryImpl 生长记录仓储实现
type growthRecordRepositoryImpl struct {
	db *gorm.DB
}

// NewGrowthRecordRepository 创建生长记录仓储
func NewGrowthRecordRepository(db *gorm.DB) repository.GrowthRecordRepository {
	return &growthRecordRepositoryImpl{db: db}
}

func (r *growthRecordRepositoryImpl) Create(ctx context.Context, record *entity.GrowthRecord) error {
	if err := r.db.WithContext(ctx).Create(record).Error; err != nil {
		return errors.Wrap(errors.DatabaseError, "failed to create growth record", err)
	}
	return nil
}

func (r *growthRecordRepositoryImpl) FindByID(ctx context.Context, recordID int64) (*entity.GrowthRecord, error) {
	var record entity.GrowthRecord
	err := r.db.WithContext(ctx).
		Where("id = ?", recordID).
		First(&record).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.ErrRecordNotFound
	}
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "failed to find growth record", err)
	}

	return &record, nil
}

func (r *growthRecordRepositoryImpl) FindByBabyID(
	ctx context.Context,
	babyID int64,
	startTime, endTime int64,
	page, pageSize int,
) ([]*entity.GrowthRecord, int64, error) {
	var records []*entity.GrowthRecord
	var total int64

	query := r.db.WithContext(ctx).
		Model(&entity.GrowthRecord{}).
		Where("baby_id = ?", babyID)

	if startTime > 0 {
		query = query.Where("time >= ?", startTime)
	}
	if endTime > 0 {
		query = query.Where("time <= ?", endTime)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.Wrap(errors.DatabaseError, "failed to count growth records", err)
	}

	offset := (page - 1) * pageSize
	err := query.
		Order("time DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&records).Error

	if err != nil {
		return nil, 0, errors.Wrap(errors.DatabaseError, "failed to find growth records", err)
	}

	return records, total, nil
}

func (r *growthRecordRepositoryImpl) Update(ctx context.Context, record *entity.GrowthRecord) error {
	err := r.db.WithContext(ctx).
		Model(&entity.GrowthRecord{}).
		Where("id = ?", record.ID).
		Updates(record).Error

	if err != nil {
		return errors.Wrap(errors.DatabaseError, "failed to update growth record", err)
	}

	return nil
}

func (r *growthRecordRepositoryImpl) Delete(ctx context.Context, recordID int64) error {
	err := r.db.WithContext(ctx).
		Where("id = ?", recordID).
		Delete(&entity.GrowthRecord{}).Error

	if err != nil {
		return errors.Wrap(errors.DatabaseError, "failed to delete growth record", err)
	}

	return nil
}

func (r *growthRecordRepositoryImpl) FindUpdatedAfter(
	ctx context.Context,
	babyID int64,
	timestamp int64,
) ([]*entity.GrowthRecord, error) {
	var records []*entity.GrowthRecord

	err := r.db.WithContext(ctx).
		Where("baby_id = ? AND updated_at > ?", babyID, timestamp).
		Find(&records).Error

	if err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "failed to find updated growth records", err)
	}

	return records, nil
}

func (r *growthRecordRepositoryImpl) GetLatestRecord(ctx context.Context, babyID int64) (*entity.GrowthRecord, error) {
	var record entity.GrowthRecord
	err := r.db.WithContext(ctx).
		Where("baby_id = ?", babyID).
		Order("time DESC").
		First(&record).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.ErrRecordNotFound
	}
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "failed to get latest growth record", err)
	}

	return &record, nil
}

func (r *growthRecordRepositoryImpl) GetDailyStats(ctx context.Context, babyID int64, startDate, endDate int64) ([]*entity.DailyGrowthItem, error) {
	var records []*entity.DailyGrowthItem
	query := r.db.WithContext(ctx).
		Model(&entity.GrowthRecord{}).
		Select(`
			to_char(to_timestamp(time / 1000), 'YYYY-MM-DD') AS date,
			MAX(height) AS latest_height,
			MAX(weight) AS latest_weight,
			MAX(head_circumference) AS latest_head_circumference,
			COUNT(*) AS record_count`).
		Where("baby_id = ? AND time BETWEEN ? AND ?", babyID, startDate, endDate).
		Group("date").
		Order("date ASC")

	if err := query.Scan(&records).Error; err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "failed to get daily growth stats", err)
	}

	return records, nil
}
