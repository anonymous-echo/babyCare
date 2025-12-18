package repository

import (
	"context"

	"github.com/wxlbd/nutri-baby-server/internal/domain/entity"
)

// SubscribeRepository 订阅消息仓储接口
type SubscribeRepository interface {
	// ==================== 消息发送队列管理(暂不使用) ====================

	// AddToSendQueue 将消息加入发送队列
	AddToSendQueue(ctx context.Context, queue *entity.MessageSendQueue) error

	// GetPendingMessages 获取待发送的消息(按计划时间排序)
	GetPendingMessages(ctx context.Context, limit int) ([]*entity.MessageSendQueue, error)

	// UpdateQueueStatus 更新队列消息状态
	UpdateQueueStatus(ctx context.Context, id int64, status string, errorMsg string) error

	// IncrementRetryCount 增加重试次数
	IncrementRetryCount(ctx context.Context, id int64) error

	// DeleteQueueMessage 删除队列消息
	DeleteQueueMessage(ctx context.Context, id int64) error

	// ==================== 消息发送日志 ====================

	// CreateSendLog 创建发送日志
	CreateSendLog(ctx context.Context, log *entity.MessageSendLog) error

	// GetSendLogs 获取发送日志(分页)
	GetSendLogs(ctx context.Context, userID int64, offset, limit int) ([]*entity.MessageSendLog, int64, error)

	// GetSendLogsByTemplateType 根据模板类型获取发送日志
	GetSendLogsByTemplateType(ctx context.Context, templateType string, offset, limit int) ([]*entity.MessageSendLog, int64, error)

	// GetRecentFailedLogs 获取最近失败的发送日志
	GetRecentFailedLogs(ctx context.Context, hours int, limit int) ([]*entity.MessageSendLog, error)
}
