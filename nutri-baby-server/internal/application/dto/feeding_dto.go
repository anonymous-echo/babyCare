package dto

// CreateFeedingRecordRequest 创建喂养记录请求
type CreateFeedingRecordRequest struct {
	BabyID             string         `json:"babyId" binding:"required"`
	FeedingType        string         `json:"feedingType" binding:"required,oneof=breast bottle food"`
	Amount             *int64         `json:"amount"`
	Duration           *int           `json:"duration"`
	Detail             map[string]any `json:"detail" binding:"required"`
	Note               *string        `json:"note"`
	FeedingTime        int64          `json:"feedingTime" binding:"required"`
	ActualCompleteTime *int64         `json:"actualCompleteTime"` // 实际喂养完成时间戳(毫秒)，用于准确计算提醒时间
	ReminderInterval   *int           `json:"reminderInterval"`   // 提醒间隔(分钟)
}

// FeedingRecordResponse 喂养记录响应
type FeedingRecordResponse struct {
	RecordID           string         `json:"recordId"`
	BabyID             string         `json:"babyId"`
	FeedingType        string         `json:"feedingType"`
	Amount             *int64         `json:"amount,omitempty"`
	Duration           *int           `json:"duration,omitempty"`
	Detail             map[string]any `json:"detail"`
	Note               *string        `json:"note,omitempty"`
	FeedingTime        int64          `json:"feedingTime"`
	ActualCompleteTime *int64         `json:"actualCompleteTime,omitempty"` // 实际喂养完成时间戳(毫秒)
	ReminderInterval   *int           `json:"reminderInterval,omitempty"`
	NextReminderTime   *int64         `json:"nextReminderTime,omitempty"`
	CreateBy           string         `json:"createBy"`
	CreateByName       string         `json:"createByName"`
	CreateByAvatar     string         `json:"createByAvatar"`
	CreateTime         int64          `json:"createTime"`
	UpdateTime         int64          `json:"updateTime"`
}

// UpdateFeedingRecordRequest 更新喂养记录请求
// 所有字段使用指针类型，支持部分更新（只更新非nil字段）
type UpdateFeedingRecordRequest struct {
	FeedingType        *string        `json:"feedingType,omitempty" binding:"omitempty,oneof=breast bottle food"`
	Amount             *int64         `json:"amount,omitempty"`
	Duration           *int           `json:"duration,omitempty"`
	Detail             map[string]any `json:"detail,omitempty"`
	Note               *string        `json:"note,omitempty"`
	FeedingTime        *int64         `json:"feedingTime,omitempty"`
	ActualCompleteTime *int64         `json:"actualCompleteTime,omitempty"`
	ReminderInterval   *int           `json:"reminderInterval,omitempty"`
}

// FeedingRecordsListResponse 喂养记录列表响应
type FeedingRecordsListResponse struct {
	Records  []*FeedingRecordResponse `json:"records"`
	Total    int64                    `json:"total"`
	Page     int                      `json:"page"`
	PageSize int                      `json:"pageSize"`
}
