package entity

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

// SubscribeRecord 订阅记录实体 (一次性订阅消息机制)
type SubscribeRecord struct {
	ID           int64  `gorm:"primaryKey;column:id" json:"id"`                                                // 雪花ID主键
	UserID       int64  `gorm:"column:user_id;not null;index:idx_user_type" json:"userId"`                     // 用户ID (引用User.ID)
	TemplateID   string `gorm:"column:template_id;size:128;not null" json:"templateId"`                        // 模板ID
	TemplateType string `gorm:"column:template_type;size:32;not null;index:idx_user_type" json:"templateType"` // 模板类型

	// 状态: available(可用), used(已使用), expired(已过期)
	Status string `gorm:"column:status;size:16;not null;default:'available';index" json:"status"`

	AuthorizeTime int64  `gorm:"column:authorize_time;not null" json:"authorizeTime"` // 授权时间(毫秒时间戳)
	UsedTime      *int64 `gorm:"column:used_time" json:"usedTime,omitempty"`          // 使用时间(毫秒时间戳)
	ExpireTime    *int64 `gorm:"column:expire_time" json:"expireTime,omitempty"`      // 过期时间(授权后7天, 毫秒时间戳)

	CreatedAt int64                 `gorm:"column:created_at;autoCreateTime:milli" json:"createdAt"`     // 创建时间(毫秒时间戳)
	UpdatedAt int64                 `gorm:"column:updated_at;autoUpdateTime:milli" json:"updatedAt"`     // 更新时间(毫秒时间戳)
	DeletedAt soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;index;default:0" json:"-"` // 软删除(毫秒时间戳)
}

// TableName 指定表名
func (SubscribeRecord) TableName() string {
	return "subscribe_records"
}

// IsAvailable 检查授权是否可用 (一次性订阅消息机制)
func (s *SubscribeRecord) IsAvailable() bool {
	// 检查状态
	if s.Status != "available" {
		return false
	}

	// 检查是否过期(微信一次性订阅消息有效期为7天)
	if s.ExpireTime != nil && time.Now().UnixMilli() > *s.ExpireTime {
		return false
	}

	return true
}

// MarkAsUsed 标记为已使用
func (s *SubscribeRecord) MarkAsUsed() {
	s.Status = "used"
	now := time.Now().UnixMilli()
	s.UsedTime = &now
}

// MarkAsExpired 标记为已过期
func (s *SubscribeRecord) MarkAsExpired() {
	s.Status = "expired"
}

// MessageSendLog 消息发送日志实体
type MessageSendLog struct {
	ID               int64  `gorm:"primaryKey;column:id" json:"id"`                                            // 雪花ID主键
	UserID           int64  `gorm:"column:user_id;not null;index" json:"userId"`                               // 用户ID (引用User.ID)
	TemplateID       string `gorm:"column:template_id;type:varchar(128);not null" json:"templateId"`           // 模板ID
	Data             string `gorm:"column:data;type:jsonb;not null" json:"data"`                               // JSONB存储
	Page             string `gorm:"column:page;type:varchar(256)" json:"page,omitempty"`                       // 小程序页面路径
	MiniprogramState string `gorm:"column:miniprogram_state;size:32;default:'formal'" json:"miniprogramState"` // 小程序状态
	SendStatus       string `gorm:"column:send_status;type:varchar(16);not null;index" json:"sendStatus"`      // success/failed/pending
	ErrCode          int    `gorm:"column:err_code" json:"errcode,omitempty"`                                  // 错误码
	ErrMsg           string `gorm:"column:err_msg;type:text" json:"errmsg,omitempty"`                          // 错误信息
	SendTime         *int64 `gorm:"column:send_time;index" json:"sendTime,omitempty"`                          // 发送时间(毫秒时间戳)
	CreatedAt        int64  `gorm:"column:created_at;autoCreateTime:milli;default:0" json:"createdAt"`         // 创建时间(毫秒时间戳)
}

// TableName 指定表名
func (MessageSendLog) TableName() string {
	return "message_send_logs"
}

// MessageSendQueue 消息发送队列实体
type MessageSendQueue struct {
	ID            int64  `gorm:"primaryKey;column:id" json:"id"`                                       // 雪花ID主键
	UserID        int64  `gorm:"column:user_id;not null;index" json:"userId"`                          // 用户ID (引用User.ID)
	TemplateID    string `gorm:"column:template_id;type:varchar(128);not null" json:"templateId"`      // 模板ID
	TemplateType  string `gorm:"column:template_type;type:varchar(32);not null" json:"templateType"`   // 模板类型
	Data          string `gorm:"column:data;type:jsonb;not null" json:"data"`                          // JSONB存储
	Page          string `gorm:"column:page;size:256" json:"page,omitempty"`                           // 小程序页面路径
	ScheduledTime int64  `gorm:"column:scheduled_time;not null;index" json:"scheduledTime"`            // 计划发送时间(毫秒时间戳)
	RetryCount    int    `gorm:"column:retry_count;not null;default:0" json:"retryCount"`              // 重试次数
	MaxRetry      int    `gorm:"column:max_retry;not null;default:3" json:"maxRetry"`                  // 最大重试次数
	Status        string `gorm:"column:status;size:16;not null;default:'pending';index" json:"status"` // pending/processing/sent/failed
	ErrorMsg      string `gorm:"column:error_msg;type:text" json:"errorMsg,omitempty"`                 // 错误信息
	CreatedAt     int64  `gorm:"column:created_at;autoCreateTime:milli;default:0" json:"createdAt"`    // 创建时间(毫秒时间戳)
	UpdatedAt     int64  `gorm:"column:updated_at;autoUpdateTime:milli;default:0" json:"updatedAt"`    // 更新时间(毫秒时间戳)
}

// TableName 指定表名
func (MessageSendQueue) TableName() string {
	return "message_send_queue"
}

// CanRetry 判断是否可以重试
func (m *MessageSendQueue) CanRetry() bool {
	return m.RetryCount < m.MaxRetry
}

// IncrementRetry 增加重试次数
func (m *MessageSendQueue) IncrementRetry() {
	m.RetryCount++
}
