package persistence

import (
	"context"

	"gorm.io/gorm"

	"github.com/wxlbd/nutri-baby-server/internal/domain/entity"
	"github.com/wxlbd/nutri-baby-server/internal/domain/repository"
	"github.com/wxlbd/nutri-baby-server/pkg/errors"
)

// diaperRecordRepositoryImpl 尿布记录仓储实现
type diaperRecordRepositoryImpl struct {
	db *gorm.DB
}

// NewDiaperRecordRepository 创建尿布记录仓储
func NewDiaperRecordRepository(db *gorm.DB) repository.DiaperRecordRepository {
	return &diaperRecordRepositoryImpl{db: db}
}

func (r *diaperRecordRepositoryImpl) Create(ctx context.Context, record *entity.DiaperRecord) error {
	if err := r.db.WithContext(ctx).Create(record).Error; err != nil {
		return errors.Wrap(errors.DatabaseError, "failed to create diaper record", err)
	}
	return nil
}

func (r *diaperRecordRepositoryImpl) FindByID(ctx context.Context, recordID int64) (*entity.DiaperRecord, error) {
	var record entity.DiaperRecord
	err := r.db.WithContext(ctx).
		Where("id = ?", recordID).
		First(&record).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.ErrRecordNotFound
	}
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "failed to find diaper record", err)
	}

	return &record, nil
}

func (r *diaperRecordRepositoryImpl) FindByBabyID(
	ctx context.Context,
	babyID int64,
	startTime, endTime int64,
	page, pageSize int,
) ([]*entity.DiaperRecord, int64, error) {
	var records []*entity.DiaperRecord
	var total int64

	query := r.db.WithContext(ctx).
		Model(&entity.DiaperRecord{}).
		Where("baby_id = ?", babyID)

	if startTime > 0 {
		query = query.Where("time >= ?", startTime)
	}
	if endTime > 0 {
		query = query.Where("time <= ?", endTime)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.Wrap(errors.DatabaseError, "failed to count diaper records", err)
	}

	offset := (page - 1) * pageSize
	err := query.
		Order("time DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&records).Error

	if err != nil {
		return nil, 0, errors.Wrap(errors.DatabaseError, "failed to find diaper records", err)
	}

	return records, total, nil
}

func (r *diaperRecordRepositoryImpl) Update(ctx context.Context, record *entity.DiaperRecord) error {
	err := r.db.WithContext(ctx).
		Model(&entity.DiaperRecord{}).
		Where("id = ?", record.ID).
		Updates(record).Error

	if err != nil {
		return errors.Wrap(errors.DatabaseError, "failed to update diaper record", err)
	}

	return nil
}

func (r *diaperRecordRepositoryImpl) Delete(ctx context.Context, recordID int64) error {
	err := r.db.WithContext(ctx).
		Where("id = ?", recordID).
		Delete(&entity.DiaperRecord{}).Error

	if err != nil {
		return errors.Wrap(errors.DatabaseError, "failed to delete diaper record", err)
	}

	return nil
}

func (r *diaperRecordRepositoryImpl) FindUpdatedAfter(
	ctx context.Context,
	babyID int64,
	timestamp int64,
) ([]*entity.DiaperRecord, error) {
	var records []*entity.DiaperRecord

	err := r.db.WithContext(ctx).
		Where("baby_id = ? AND updated_at > ?", babyID, timestamp).
		Find(&records).Error

	if err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "failed to find updated diaper records", err)
	}

	return records, nil
}

func (r *diaperRecordRepositoryImpl) GetDailyStats(ctx context.Context, babyID int64, startDate, endDate int64) ([]*entity.DailyDiaperItem, error) {
	var records []*entity.DailyDiaperItem
	query := r.db.WithContext(ctx).
		Model(&entity.DiaperRecord{}).
		Select(`
			to_char(to_timestamp(time / 1000), 'YYYY-MM-DD') AS date,
			type AS diaper_type,
			COUNT(*) AS total_count`).
		Where("baby_id = ? AND time BETWEEN ? AND ?", babyID, startDate, endDate).
		Group("date, type").
		Order("date ASC")

	if err := query.Scan(&records).Error; err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "failed to get daily diaper stats", err)
	}

	return records, nil
}
