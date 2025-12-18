package dto

// ============ 今日统计 ============

// TodayFeedingStats 今日喂养统计
type TodayFeedingStats struct {
	BreastCount     int    `json:"breastCount"`               // 母乳喂养次数
	BottleMl        int64  `json:"bottleMl"`                  // 奶瓶总毫升数
	BottleCount     int    `json:"bottleCount"`               // 奶瓶喂养次数
	FoodCount       int    `json:"foodCount"`                 // 饮食喂养次数
	TotalCount      int    `json:"totalCount"`                // 总喂养次数
	LastFeedingTime *int64 `json:"lastFeedingTime,omitempty"` // 最后一次喂养时间戳(毫秒)，nil 表示今天无喂养记录
}

// TodaySleepStats 今日睡眠统计
type TodaySleepStats struct {
	TotalMinutes     int `json:"totalMinutes"`     // 总睡眠分钟数
	LastSleepMinutes int `json:"lastSleepMinutes"` // 上次睡眠分钟数
	SessionCount     int `json:"sessionCount"`     // 睡眠次数
}

// TodayDiaperStats 今日换尿布统计
type TodayDiaperStats struct {
	TotalCount int `json:"totalCount"` // 总换尿布次数
	PeeCount   int `json:"peeCount"`   // 小便次数
	PoopCount  int `json:"poopCount"`  // 大便次数
}

// TodayGrowthStats 今日成长统计
type TodayGrowthStats struct {
	LatestWeight            *float64 `json:"latestWeight,omitempty"`            // 最新体重 (kg)
	LatestHeight            *float64 `json:"latestHeight,omitempty"`            // 最新身高 (cm)
	LatestHeadCircumference *float64 `json:"latestHeadCircumference,omitempty"` // 最新头围 (cm)
}

// TodayStatistics 今日统计
type TodayStatistics struct {
	Feeding TodayFeedingStats `json:"feeding"` // 喂养统计
	Sleep   TodaySleepStats   `json:"sleep"`   // 睡眠统计
	Diaper  TodayDiaperStats  `json:"diaper"`  // 换尿布统计
	Growth  TodayGrowthStats  `json:"growth"`  // 成长统计
}

// ============ 本周统计 ============

// WeeklyFeedingStats 本周喂养统计
type WeeklyFeedingStats struct {
	TotalCount int     `json:"totalCount"` // 本周总喂养次数
	Trend      int     `json:"trend"`      // 趋势对比（与上周的差异）
	AvgPerDay  float64 `json:"avgPerDay"`  // 日均喂养次数
}

// WeeklySleepStats 本周睡眠统计
type WeeklySleepStats struct {
	TotalMinutes int     `json:"totalMinutes"` // 本周总睡眠分钟数
	Trend        float64 `json:"trend"`        // 趋势对比（与上周的分钟数差异）
	AvgPerDay    float64 `json:"avgPerDay"`    // 日均睡眠分钟数
}

// WeeklyGrowthStats 本周成长统计
type WeeklyGrowthStats struct {
	WeightGain      float64  `json:"weightGain"`                // 周内体重增长 (kg)
	HeightGain      float64  `json:"heightGain"`                // 周内身高增长 (cm)
	WeekStartWeight *float64 `json:"weekStartWeight,omitempty"` // 周初体重 (kg)
}

// WeeklyStatistics 本周统计
type WeeklyStatistics struct {
	Feeding WeeklyFeedingStats `json:"feeding"` // 喂养统计
	Sleep   WeeklySleepStats   `json:"sleep"`   // 睡眠统计
	Growth  WeeklyGrowthStats  `json:"growth"`  // 成长统计
}

// ============ 完整统计响应 ============

// BabyStatisticsResponse 宝宝统计响应
type BabyStatisticsResponse struct {
	Today  TodayStatistics  `json:"today"`  // 今日统计
	Weekly WeeklyStatistics `json:"weekly"` // 本周统计
}
