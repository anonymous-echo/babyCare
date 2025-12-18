package entity

// DailyFeedingItem 表示某日某类型喂养的聚合结果
type DailyFeedingItem struct {
	Date          string // 日期，格式 YYYY-MM-DD
	FeedingType   string // 喂养类型
	TotalCount    int64  // 总次数
	TotalAmount   int64  // 总量（ml）
	TotalDuration int64  // 总时长（秒）
}

type DailySleepItem struct {
	Date          string // 日期，格式 YYYY-MM-DD
	TotalDuration int64  // 总时长（秒）
	TotalCount    int64  // 总次数
}

type DailyDiaperItem struct {
	Date       string // 日期，格式 YYYY-MM-DD
	DiaperType string // 排泄类型：pee/poop/both
	TotalCount int64  // 总次数
}

type DailyGrowthItem struct {
	Date                    string // 日期，格式 YYYY-MM-DD
	LatestHeight            *int64 // 最新身高（cm）
	LatestWeight            *int64 // 最新体重（g）
	LatestHeadCircumference *int64 // 最新头围（cm）
	RecordCount             int64  // 当日记录数
}
