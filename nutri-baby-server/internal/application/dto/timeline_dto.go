package dto

// TimelineQuery 时间线查询参数
type TimelineQuery struct {
	BabyID     string `form:"babyId" binding:"required"`
	StartTime  int64  `form:"startTime"`
	EndTime    int64  `form:"endTime"`
	RecordType string `form:"recordType"` // 可选: "feeding" | "sleep" | "diaper" | "growth" | "" (空表示全部)
	PaginationRequest
}

// TimelineItem 时间线记录项
type TimelineItem struct {
	RecordType   string `json:"recordType"` // "feeding" | "sleep" | "diaper" | "growth"
	RecordID     string `json:"recordId"`
	BabyID       string `json:"babyId"`
	EventTime    int64  `json:"eventTime"` // 统一时间戳
	Detail       any    `json:"detail"`    // 具体记录详情
	CreateBy     string `json:"createBy"`
	CreateTime   int64  `json:"createTime"`
	CreateName   string `json:"createName"`   // 创建者昵称
	Relationship string `json:"relationship"` // 创建者与宝宝的关系
}

// TimelineResponse 时间线响应
type TimelineResponse struct {
	Items    []TimelineItem `json:"items"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"pageSize"`
}
