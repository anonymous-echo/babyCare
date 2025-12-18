package persistence

import (
	"context"
	"time"

	"github.com/wxlbd/nutri-baby-server/pkg/errors"
	"gorm.io/gorm"

	"github.com/wxlbd/nutri-baby-server/internal/domain/entity"
	"github.com/wxlbd/nutri-baby-server/internal/domain/repository"
)

type subscribeRepositoryImpl struct {
	db *gorm.DB
}

// NewSubscribeRepository 创建订阅仓储实例
func NewSubscribeRepository(db *gorm.DB) repository.SubscribeRepository {
	return &subscribeRepositoryImpl{
		db: db,
	}
}

// ==================== 一次性订阅授权管理 ====================

// CreateSubscribeRecord 创建订阅授权记录(每次授权创建新记录)
func (r *subscribeRepositoryImpl) CreateSubscribeRecord(ctx context.Context, record *entity.SubscribeRecord) error {
	return r.db.WithContext(ctx).Create(record).Error
}

// GetAvailableSubscribeRecord 获取用户可用的授权记录(按授权时间倒序,取最新的一条)
func (r *subscribeRepositoryImpl) GetAvailableSubscribeRecord(ctx context.Context, userID int64, templateType string) (*entity.SubscribeRecord, error) {
	var record entity.SubscribeRecord
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND template_type = ? AND status = ?", userID, templateType, "available").
		Where("expire_time > ?", time.Now().UnixMilli()).
		Order("authorize_time DESC").
		First(&record).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 没有可用授权
		}
		return nil, err
	}

	return &record, nil
}

// GetSubscribeRecord 根据userID和模板类型获取最新的订阅记录
func (r *subscribeRepositoryImpl) GetSubscribeRecord(ctx context.Context, userID int64, templateType string) (*entity.SubscribeRecord, error) {
	var record entity.SubscribeRecord
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND template_type = ?", userID, templateType).
		Order("authorize_time DESC").
		First(&record).Error

	if err != nil {
		return nil, err
	}
	return &record, nil
}

// ListUserSubscriptions 获取用户的所有订阅记录(包括已使用和已过期)
func (r *subscribeRepositoryImpl) ListUserSubscriptions(ctx context.Context, userID int64) ([]*entity.SubscribeRecord, error) {
	var records []*entity.SubscribeRecord
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("authorize_time DESC").
		Find(&records).Error
	return records, err
}

// CountAvailableAuthorizations 统计用户可用的授权数量
func (r *subscribeRepositoryImpl) CountAvailableAuthorizations(ctx context.Context, userID int64, templateType string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.SubscribeRecord{}).
		Where("user_id = ? AND template_type = ? AND status = ?", userID, templateType, "available").
		Where("expire_time > ?", time.Now().UnixMilli()).
		Count(&count).Error
	return count, err
}

// DeleteExpiredRecords 清理过期的授权记录(定期任务使用)
func (r *subscribeRepositoryImpl) DeleteExpiredRecords(ctx context.Context, beforeTime time.Time) (int64, error) {
	result := r.db.WithContext(ctx).
		Where("expire_time < ?", beforeTime).
		Or("status = ? AND used_time < ?", "used", beforeTime).
		Delete(&entity.SubscribeRecord{})

	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}

// ==================== 消息发送队列管理 ====================

func (r *subscribeRepositoryImpl) AddToSendQueue(ctx context.Context, queue *entity.MessageSendQueue) error {
	return r.db.WithContext(ctx).Create(queue).Error
}

func (r *subscribeRepositoryImpl) GetPendingMessages(ctx context.Context, limit int) ([]*entity.MessageSendQueue, error) {
	var messages []*entity.MessageSendQueue
	err := r.db.WithContext(ctx).
		Where("status = ? AND scheduled_time <= ?", "pending", time.Now()).
		Order("scheduled_time ASC").
		Limit(limit).
		Find(&messages).Error
	return messages, err
}

func (r *subscribeRepositoryImpl) UpdateQueueStatus(ctx context.Context, id int64, status string, errorMsg string) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now().UnixMilli(),
	}
	if errorMsg != "" {
		updates["error_msg"] = errorMsg
	}
	return r.db.WithContext(ctx).
		Model(&entity.MessageSendQueue{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *subscribeRepositoryImpl) IncrementRetryCount(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Model(&entity.MessageSendQueue{}).
		Where("id = ?", id).
		UpdateColumn("retry_count", gorm.Expr("retry_count + 1")).Error
}

func (r *subscribeRepositoryImpl) DeleteQueueMessage(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Delete(&entity.MessageSendQueue{}, id).Error
}

// ==================== 消息发送日志 ====================

func (r *subscribeRepositoryImpl) CreateSendLog(ctx context.Context, log *entity.MessageSendLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *subscribeRepositoryImpl) GetSendLogs(ctx context.Context, userID int64, offset, limit int) ([]*entity.MessageSendLog, int64, error) {
	var logs []*entity.MessageSendLog
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.MessageSendLog{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&logs).Error

	return logs, total, err
}

func (r *subscribeRepositoryImpl) GetSendLogsByTemplateType(ctx context.Context, templateType string, offset, limit int) ([]*entity.MessageSendLog, int64, error) {
	var logs []*entity.MessageSendLog
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.MessageSendLog{}).Where("template_type = ?", templateType)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&logs).Error

	return logs, total, err
}

func (r *subscribeRepositoryImpl) GetRecentFailedLogs(ctx context.Context, hours int, limit int) ([]*entity.MessageSendLog, error) {
	var logs []*entity.MessageSendLog
	since := time.Now().Add(-time.Duration(hours) * time.Hour)

	err := r.db.WithContext(ctx).
		Where("send_status = ? AND created_at >= ?", "failed", since).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error

	return logs, err
}
