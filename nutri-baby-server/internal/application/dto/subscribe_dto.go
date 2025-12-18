package dto

// ======================== 订阅授权相关 ========================

// SubscribeAuthDTO 单个订阅授权记录
type SubscribeAuthDTO struct {
	TemplateID   string `json:"templateId" binding:"required"`
	TemplateType string `json:"templateType" binding:"required"`
	Status       string `json:"status" binding:"required,oneof=accept reject"`
}

// SubscribeAuthRequest 批量上传授权记录请求
type SubscribeAuthRequest struct {
	Records []SubscribeAuthDTO `json:"records" binding:"required,min=1,dive"`
}

// SubscribeAuthResponse 授权上传响应
type SubscribeAuthResponse struct {
	SuccessCount int `json:"successCount"`
	FailedCount  int `json:"failedCount"`
}

// ======================== 订阅状态查询 ========================

// SubscriptionItem 单个订阅项
type SubscriptionItem struct {
	TemplateType  string `json:"templateType"`
	TemplateID    string `json:"templateId"`
	Status        string `json:"status"`
	SubscribeTime int64  `json:"subscribeTime"`
	ExpireTime    int64  `json:"expireTime,omitempty"`
}

// SubscribeStatusResponse 订阅状态响应
type SubscribeStatusResponse struct {
	Subscriptions []SubscriptionItem `json:"subscriptions"`
}

// ======================== 取消订阅 ========================

// CancelSubscriptionRequest 取消订阅请求
type CancelSubscriptionRequest struct {
	TemplateType string `json:"templateType" binding:"required"`
}

// ======================== 消息发送 ========================

// SendMessageRequest 发送消息请求(内部使用)
type SendMessageRequest struct {
	OpenID       string         `json:"openid" binding:"required"`
	TemplateType string         `json:"templateType" binding:"required"`
	TemplateID   string         `json:"templateId" binding:"required"`
	Data         map[string]any `json:"data" binding:"required"`
	Page         string         `json:"page"`
}

// QueueMessageRequest 加入队列请求
type QueueMessageRequest struct {
	OpenID        string         `json:"openid" binding:"required"`
	TemplateType  string         `json:"templateType" binding:"required"`
	Data          map[string]any `json:"data" binding:"required"`
	Page          string         `json:"page"`
	ScheduledTime int64          `json:"scheduledTime"` // Unix timestamp
}

// ======================== 消息发送日志 ========================

// MessageLogItem 消息日志项
type MessageLogItem struct {
	ID           uint   `json:"id"`
	TemplateType string `json:"templateType"`
	SendStatus   string `json:"sendStatus"`
	ErrMsg       string `json:"errmsg,omitempty"`
	SendTime     int64  `json:"sendTime,omitempty"`
	CreatedAt    int64  `json:"createdAt"`
}

// MessageLogsResponse 消息日志响应
type MessageLogsResponse struct {
	Logs  []MessageLogItem `json:"logs"`
	Total int64            `json:"total"`
}
