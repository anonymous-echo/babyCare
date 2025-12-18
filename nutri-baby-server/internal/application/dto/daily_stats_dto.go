package dto

// ============ 按日统计 DTO ============

// DailyFeedingStatsItem 每日喂养统计项
type DailyFeedingStatsItem struct {
	Date          string `json:"date"`          // 日期，格式 YYYY-MM-DD
	FeedingType   string `json:"feedingType"`   // 喂养类型：breast/bottle/food
	TotalCount    int64  `json:"totalCount"`    // 总次数
	TotalAmount   int64  `json:"totalAmount"`   // 总量（ml）
	TotalDuration int64  `json:"totalDuration"` // 总时长（秒）
}

// DailySleepStatsItem 每日睡眠统计项
type DailySleepStatsItem struct {
	Date          string `json:"date"`          // 日期，格式 YYYY-MM-DD
	TotalDuration int64  `json:"totalDuration"` // 总时长（秒）
	TotalCount    int64  `json:"totalCount"`    // 总次数
}

// DailyDiaperStatsItem 每日排泄统计项
type DailyDiaperStatsItem struct {
	Date       string `json:"date"`       // 日期，格式 YYYY-MM-DD
	DiaperType string `json:"diaperType"` // 排泄类型：pee/poop/both
	TotalCount int64  `json:"totalCount"` // 总次数
}

// DailyGrowthStatsItem 每日成长统计项
type DailyGrowthStatsItem struct {
	Date                    string `json:"date"`                    // 日期，格式 YYYY-MM-DD
	LatestHeight            *int64 `json:"latestHeight"`            // 最新身高（cm）
	LatestWeight            *int64 `json:"latestWeight"`            // 最新体重（g）
	LatestHeadCircumference *int64 `json:"latestHeadCircumference"` // 最新头围（cm）
	RecordCount             int64  `json:"recordCount"`             // 当日记录数
}

// DailyStatsRequest 按日统计请求
type DailyStatsRequest struct {
	BabyID    string `form:"babyId" binding:"required"`    // 宝宝ID
	StartDate int64  `form:"startDate" binding:"required"` // 开始日期（毫秒时间戳）
	EndDate   int64  `form:"endDate" binding:"required"`   // 结束日期（毫秒时间戳）
	Types     string `form:"types"`                        // 统计类型，逗号分隔：feeding,sleep,diaper,growth，默认全部
}

// DailyStatsResponse 按日统计响应
type DailyStatsResponse struct {
	Feeding []*DailyFeedingStatsItem `json:"feeding,omitempty"` // 喂养统计
	Sleep   []*DailySleepStatsItem   `json:"sleep,omitempty"`   // 睡眠统计
	Diaper  []*DailyDiaperStatsItem  `json:"diaper,omitempty"`  // 排泄统计
	Growth  []*DailyGrowthStatsItem  `json:"growth,omitempty"`  // 成长统计
}
