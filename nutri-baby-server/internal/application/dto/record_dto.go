package dto

// ============ 喂养 Detail 强类型定义 ============

// FeedingSession 单次喂养会话(多段式记录)
type FeedingSession struct {
	Side      string `json:"side"`                // left, right, both
	StartTime int64  `json:"startTime"`           // 开始时间戳(毫秒)
	EndTime   *int64 `json:"endTime,omitempty"`   // 结束时间戳(毫秒,可选)
	Duration  int    `json:"duration"`            // 时长(秒)
}

// BreastFeedingDetail 母乳喂养详情
type BreastFeedingDetail struct {
	Type          string            `json:"type"`                    // "breast" 固定值
	Side          string            `json:"side"`                    // left, right, both - 主要喂养侧(向后兼容)
	Duration      int               `json:"duration"`                // 总时长(秒)
	LeftDuration  *int              `json:"leftDuration,omitempty"`  // 左侧总时长(秒)
	RightDuration *int              `json:"rightDuration,omitempty"` // 右侧总时长(秒)
	Sessions      []FeedingSession  `json:"sessions,omitempty"`      // 多段式记录(可选,新功能)
}

// BottleFeedingDetail 奶瓶喂养详情
type BottleFeedingDetail struct {
	Type       string   `json:"type"`                 // "bottle" 固定值
	BottleType string   `json:"bottleType"`           // formula, breast-milk
	Amount     int64    `json:"amount"`               // 奶量(ml)
	Unit       string   `json:"unit"`                 // ml, oz
	Remaining  *float64 `json:"remaining,omitempty"`  // 剩余量(可选)
}

// FoodFeedingDetail 辅食详情
type FoodFeedingDetail struct {
	Type     string  `json:"type"`              // "food" 固定值
	FoodName string  `json:"foodName"`          // 辅食名称
	Note     *string `json:"note,omitempty"`    // 备注(接受程度、过敏反应等)
}

// FeedingDetail 喂养详情(向后兼容的全能结构体，用于数据库JSONB存储)
// 注意：这个结构体仅用于数据库存储和反序列化，业务逻辑应使用上面的强类型结构体
type FeedingDetail struct {
	// 类型标识
	Type string `json:"type"` // breast, bottle, food

	// 母乳喂养相关
	Side          string           `json:"side,omitempty"`          // left, right, both
	Duration      int              `json:"duration,omitempty"`      // 总时长(秒)
	LeftDuration  *int             `json:"leftDuration,omitempty"`  // 左侧时长(秒)
	RightDuration *int             `json:"rightDuration,omitempty"` // 右侧时长(秒)
	Sessions      []FeedingSession `json:"sessions,omitempty"`      // 多段式记录

	// 奶瓶喂养相关
	BottleType string   `json:"bottleType,omitempty"` // formula, breast-milk
	Amount     int64    `json:"amount,omitempty"`     // 奶量
	Unit       string   `json:"unit,omitempty"`       // ml, oz
	Remaining  *float64 `json:"remaining,omitempty"`  // 剩余量

	// 辅食相关
	FoodName string `json:"foodName,omitempty"` // 辅食名称

	// 通用
	Note *string `json:"note,omitempty"` // 备注

	// 向后兼容旧数据
	BreastSide  string `json:"breastSide,omitempty"`  // 旧字段：兼容 breastSide
	LeftTime    int    `json:"leftTime,omitempty"`    // 旧字段：兼容 leftTime
	RightTime   int    `json:"rightTime,omitempty"`   // 旧字段：兼容 rightTime
	FormulaType string `json:"formulaType,omitempty"` // 旧字段：兼容 formulaType
}

// ToBreastFeeding 转换为母乳喂养详情
func (d *FeedingDetail) ToBreastFeeding() *BreastFeedingDetail {
	if d.Type != "breast" {
		return nil
	}

	// 向后兼容处理
	side := d.Side
	if side == "" && d.BreastSide != "" {
		side = d.BreastSide
	}

	duration := d.Duration
	if duration == 0 && (d.LeftTime > 0 || d.RightTime > 0) {
		duration = d.LeftTime + d.RightTime
	}

	return &BreastFeedingDetail{
		Type:          "breast",
		Side:          side,
		Duration:      duration,
		LeftDuration:  d.LeftDuration,
		RightDuration: d.RightDuration,
		Sessions:      d.Sessions,
	}
}

// ToBottleFeeding 转换为奶瓶喂养详情
func (d *FeedingDetail) ToBottleFeeding() *BottleFeedingDetail {
	if d.Type != "bottle" {
		return nil
	}
	return &BottleFeedingDetail{
		Type:       "bottle",
		BottleType: d.BottleType,
		Amount:     d.Amount,
		Unit:       d.Unit,
		Remaining:  d.Remaining,
	}
}

// ToFoodFeeding 转换为辅食详情
func (d *FeedingDetail) ToFoodFeeding() *FoodFeedingDetail {
	if d.Type != "food" {
		return nil
	}
	return &FoodFeedingDetail{
		Type:     "food",
		FoodName: d.FoodName,
		Note:     d.Note,
	}
}

// FromBreastFeeding 从母乳喂养详情创建
func FromBreastFeeding(detail *BreastFeedingDetail) *FeedingDetail {
	return &FeedingDetail{
		Type:          "breast",
		Side:          detail.Side,
		Duration:      detail.Duration,
		LeftDuration:  detail.LeftDuration,
		RightDuration: detail.RightDuration,
		Sessions:      detail.Sessions,
	}
}

// FromBottleFeeding 从奶瓶喂养详情创建
func FromBottleFeeding(detail *BottleFeedingDetail) *FeedingDetail {
	return &FeedingDetail{
		Type:       "bottle",
		BottleType: detail.BottleType,
		Amount:     detail.Amount,
		Unit:       detail.Unit,
		Remaining:  detail.Remaining,
	}
}

// FromFoodFeeding 从辅食详情创建
func FromFoodFeeding(detail *FoodFeedingDetail) *FeedingDetail {
	return &FeedingDetail{
		Type:     "food",
		FoodName: detail.FoodName,
		Note:     detail.Note,
	}
}

// ============ 其他记录 DTO ============


// FeedingRecordDTO 喂养记录DTO
type FeedingRecordDTO struct {
	RecordID           string        `json:"recordId"`
	BabyID             string        `json:"babyId"`
	FeedingType        string        `json:"feedingType"`
	Amount             int64         `json:"amount"`
	Duration           int           `json:"duration"`
	Detail             FeedingDetail `json:"detail"`
	Note               string        `json:"note"`
	FeedingTime        int64         `json:"feedingTime"`
	ActualCompleteTime *int64        `json:"actualCompleteTime,omitempty"` // 实际喂养完成时间戳(毫秒)
	CreateBy           string        `json:"createBy"`
	CreateTime         int64         `json:"createTime"`
}

// CreateSleepRecordRequest 创建睡眠记录请求
type CreateSleepRecordRequest struct {
	BabyID    string `json:"babyId" binding:"required"`
	StartTime int64  `json:"startTime" binding:"required"`
	EndTime   int64  `json:"endTime"`
	Duration  int    `json:"duration"` // 秒
	SleepType string `json:"sleepType" binding:"required,oneof=nap night"` // 睡眠类型：nap(小睡) | night(夜间长睡)
	Note      string `json:"note"`
}

// SleepRecordDTO 睡眠记录DTO
type SleepRecordDTO struct {
	RecordID   string `json:"recordId"`
	BabyID     string `json:"babyId"`
	StartTime  int64  `json:"startTime"`
	EndTime    int64  `json:"endTime"`
	Duration   int    `json:"duration"`
	SleepType  string `json:"sleepType"` // 睡眠类型：nap(小睡) | night(夜间长睡)
	Note       string `json:"note"`
	CreateBy   string `json:"createBy"`
	CreateTime int64  `json:"createTime"`
}

// CreateDiaperRecordRequest 创建尿布记录请求
type CreateDiaperRecordRequest struct {
	BabyID     string `json:"babyId" binding:"required"`
	DiaperType string `json:"diaperType" binding:"required,oneof=pee poop both"`
	Note       string `json:"note"`
	ChangeTime int64  `json:"changeTime"` // 毫秒时间戳
}

// DiaperRecordDTO 尿布记录DTO
type DiaperRecordDTO struct {
	RecordID   string `json:"recordId"`
	BabyID     string `json:"babyId"`
	DiaperType string `json:"diaperType"`
	Note       string `json:"note"`
	ChangeTime int64  `json:"changeTime"`
	CreateBy   string `json:"createBy"`
	CreateTime int64  `json:"createTime"`
}

// CreateGrowthRecordRequest 创建生长记录请求
type CreateGrowthRecordRequest struct {
	BabyID            string  `json:"babyId" binding:"required"`
	Height            float64 `json:"height"`            // cm
	Weight            float64 `json:"weight"`            // kg
	HeadCircumference float64 `json:"headCircumference"` // cm
	Note              string  `json:"note"`
	MeasureTime       int64   `json:"measureTime"` // 毫秒时间戳
}

// GrowthRecordDTO 生长记录DTO
type GrowthRecordDTO struct {
	RecordID          string   `json:"recordId"`
	BabyID            string   `json:"babyId"`
	Height            *float64 `json:"height,omitempty"`            // cm, 仅当有值时返回
	Weight            *float64 `json:"weight,omitempty"`            // kg, 仅当有值时返回
	HeadCircumference *float64 `json:"headCircumference,omitempty"` // cm, 仅当有值时返回
	Note              string   `json:"note,omitempty"`
	MeasureTime       int64    `json:"measureTime"`
	CreateBy          string   `json:"createBy"`
	CreateTime        int64    `json:"createTime"`
}

// RecordListQuery 记录列表查询参数
type RecordListQuery struct {
	BabyID    string `form:"babyId"`
	StartTime int64  `form:"startTime"`
	EndTime   int64  `form:"endTime"`
	PaginationRequest
}

// ============ 更新记录 DTO ============

// UpdateSleepRecordRequest 更新睡眠记录请求
// 所有字段使用指针类型，支持部分更新（只更新非nil字段）
type UpdateSleepRecordRequest struct {
	StartTime *int64  `json:"startTime,omitempty"`
	EndTime   *int64  `json:"endTime,omitempty"`
	Duration  *int    `json:"duration,omitempty"`
	SleepType *string `json:"sleepType,omitempty" binding:"omitempty,oneof=nap night"` // 睡眠类型：nap(小睡) | night(夜间长睡)
	Note      *string `json:"note,omitempty"`
}

// UpdateDiaperRecordRequest 更新尿布记录请求
// 所有字段使用指针类型，支持部分更新（只更新非nil字段）
type UpdateDiaperRecordRequest struct {
	DiaperType *string `json:"diaperType,omitempty" binding:"omitempty,oneof=pee poop both"`
	Note       *string `json:"note,omitempty"`
	ChangeTime *int64  `json:"changeTime,omitempty"`
}

// UpdateGrowthRecordRequest 更新生长记录请求
// 所有字段使用指针类型，支持部分更新（只更新非nil字段）
type UpdateGrowthRecordRequest struct {
	Height            *float64 `json:"height,omitempty"`
	Weight            *float64 `json:"weight,omitempty"`
	HeadCircumference *float64 `json:"headCircumference,omitempty"`
	Note              *string  `json:"note,omitempty"`
	MeasureTime       *int64   `json:"measureTime,omitempty"`
}
