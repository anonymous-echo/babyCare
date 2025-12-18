package entity

import (
	"database/sql/driver"
	"encoding/json"

	"gorm.io/plugin/soft_delete"
)

// 喂养类型常量 (与前端保持一致)
const (
	FeedingTypeBreast = "breast" // 母乳喂养
	FeedingTypeBottle = "bottle" // 奶瓶喂养
	FeedingTypeFood   = "food"   // 辅食
)

// FeedingRecord 喂养记录实体
type FeedingRecord struct {
	ID              int64         `gorm:"primaryKey;column:id" json:"id"`                                    // 雪花ID主键
	BabyID          int64         `gorm:"column:baby_id;index" json:"babyId"`                                // 宝宝ID (引用Baby.ID)
	Time            int64         `gorm:"column:time;index" json:"time"`                                     // 记录时间(毫秒时间戳)
	FeedingType     string        `gorm:"column:feeding_type;type:varchar(16);not null" json:"feedingType"`  // 喂养类型: breast/bottle/food
	Amount          int64         `gorm:"column:amount" json:"amount,omitempty"`                             // 奶量(ml)，bottle类型时使用
	Duration        int           `gorm:"column:duration" json:"duration,omitempty"`                         // 时长(秒)，breast类型时使用
	Detail          FeedingDetail `gorm:"column:detail;type:jsonb" json:"detail"`                            // 完整详情(向后兼容)
	CreatedBy       int64         `gorm:"column:created_by" json:"createdBy"`                                // 创建者用户ID (引用User.ID)
	CreatedByName   string        `gorm:"column:created_by_name;type:varchar(64)" json:"createdByName"`      // 冗余:创建者昵称
	CreatedByAvatar string        `gorm:"column:created_by_avatar;type:varchar(512)" json:"createdByAvatar"` // 冗余:创建者头像

	// 提醒相关字段
	ActualCompleteTime *int64 `gorm:"column:actual_complete_time" json:"actualCompleteTime,omitempty"` // 实际喂养完成时间戳(毫秒)
	ReminderInterval   *int   `gorm:"column:reminder_interval" json:"reminderInterval,omitempty"`      // 提醒间隔(分钟)
	NextReminderTime   *int64 `gorm:"column:next_reminder_time" json:"nextReminderTime,omitempty"`     // 下次提醒时间戳(毫秒)
	ReminderSent       bool   `gorm:"column:reminder_sent;default:false;index" json:"reminderSent"`    // 是否已发送提醒
	ReminderTime       *int64 `gorm:"column:reminder_time" json:"reminderTime,omitempty"`              // 提醒发送时间戳(毫秒)

	CreatedAt int64                 `gorm:"column:created_at;autoCreateTime:milli" json:"createdAt"`                 // 创建时间(毫秒时间戳)
	UpdatedAt int64                 `gorm:"column:updated_at;autoUpdateTime:milli" json:"updatedAt"`                 // 更新时间(毫秒时间戳)
	DeletedAt soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;index;type:bigint;default:0" json:"-"` // 软删除(毫秒时间戳)
}

// TableName 指定表名
func (FeedingRecord) TableName() string {
	return "feeding_records"
}

// FeedingDetail 喂养详情(使用interface{}存储不同类型)
type FeedingDetail map[string]any

// Scan 实现sql.Scanner接口
func (f *FeedingDetail) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, f)
}

// Value 实现driver.Valuer接口
func (f FeedingDetail) Value() (driver.Value, error) {
	return json.Marshal(f)
}

// SleepRecord 睡眠记录实体
type SleepRecord struct {
	ID              int64                 `gorm:"primaryKey;column:id" json:"id"`                                    // 雪花ID主键
	BabyID          int64                 `gorm:"column:baby_id;index" json:"babyId"`                                // 宝宝ID (引用Baby.ID)
	StartTime       int64                 `gorm:"column:start_time;index" json:"startTime"`                          // 开始时间(毫秒时间戳)
	EndTime         *int64                `gorm:"column:end_time" json:"endTime"`                                    // 结束时间(毫秒时间戳)
	Duration        *int                  `gorm:"column:duration" json:"duration"`                                   // 时长(秒)
	Type            string                `gorm:"column:type;type:varchar(16)" json:"type"`                          // nap, night
	CreatedBy       int64                 `gorm:"column:created_by" json:"createdBy"`                                // 创建者用户ID (引用User.ID)
	CreatedByName   string                `gorm:"column:created_by_name;type:varchar(64)" json:"createdByName"`      // 冗余:创建者昵称
	CreatedByAvatar string                `gorm:"column:created_by_avatar;type:varchar(512)" json:"createdByAvatar"` // 冗余:创建者头像
	CreatedAt       int64                 `gorm:"column:created_at;autoCreateTime:milli" json:"createdAt"`           // 创建时间(毫秒时间戳)
	UpdatedAt       int64                 `gorm:"column:updated_at;autoUpdateTime:milli" json:"updatedAt"`           // 更新时间(毫秒时间戳)
	DeletedAt       soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;index" json:"-"`                 // 软删除(毫秒时间戳)
}

// TableName 指定表名
func (SleepRecord) TableName() string {
	return "sleep_records"
}

// DiaperRecord 换尿布记录实体
type DiaperRecord struct {
	ID              int64                 `gorm:"primaryKey;column:id" json:"id"`                                    // 雪花ID主键
	BabyID          int64                 `gorm:"column:baby_id;index" json:"babyId"`                                // 宝宝ID (引用Baby.ID)
	Time            int64                 `gorm:"column:time;index" json:"time"`                                     // 记录时间(毫秒时间戳)
	Type            string                `gorm:"column:type;type:varchar(16)" json:"type"`                          // pee, poop, both
	PoopColor       *string               `gorm:"column:poop_color;type:varchar(16)" json:"poopColor"`               // 便便颜色
	PoopTexture     *string               `gorm:"column:poop_texture;type:varchar(16)" json:"poopTexture"`           // 便便质地
	Note            *string               `gorm:"column:note;type:text" json:"note"`                                 // 备注
	CreatedBy       int64                 `gorm:"column:created_by" json:"createdBy"`                                // 创建者用户ID (引用User.ID)
	CreatedByName   string                `gorm:"column:created_by_name;type:varchar(64)" json:"createdByName"`      // 冗余:创建者昵称
	CreatedByAvatar string                `gorm:"column:created_by_avatar;type:varchar(512)" json:"createdByAvatar"` // 冗余:创建者头像
	CreatedAt       int64                 `gorm:"column:created_at;autoCreateTime:milli" json:"createdAt"`           // 创建时间(毫秒时间戳)
	UpdatedAt       int64                 `gorm:"column:updated_at;autoUpdateTime:milli" json:"updatedAt"`           // 更新时间(毫秒时间戳)
	DeletedAt       soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;index;default:0" json:"-"`       // 软删除(毫秒时间戳)
}

// TableName 指定表名
func (DiaperRecord) TableName() string {
	return "diaper_records"
}

// GrowthRecord 成长记录实体
type GrowthRecord struct {
	ID                int64                 `gorm:"primaryKey;column:id" json:"id"`                                    // 雪花ID主键
	BabyID            int64                 `gorm:"column:baby_id;index" json:"babyId"`                                // 宝宝ID (引用Baby.ID)
	Time              int64                 `gorm:"column:time;index" json:"time"`                                     // 记录时间(毫秒时间戳)
	Height            *float64              `gorm:"column:height" json:"height"`                                       // 身高 cm
	Weight            *float64              `gorm:"column:weight" json:"weight"`                                       // 体重 kg
	HeadCircumference *float64              `gorm:"column:head_circumference" json:"headCircumference"`                // 头围 cm
	Note              *string               `gorm:"column:note;type:text" json:"note"`                                 // 备注
	CreatedBy         int64                 `gorm:"column:created_by" json:"createdBy"`                                // 创建者用户ID (引用User.ID)
	CreatedByName     string                `gorm:"column:created_by_name;type:varchar(64)" json:"createdByName"`      // 冗余:创建者昵称
	CreatedByAvatar   string                `gorm:"column:created_by_avatar;type:varchar(512)" json:"createdByAvatar"` // 冗余:创建者头像
	CreatedAt         int64                 `gorm:"column:created_at;autoCreateTime:milli" json:"createdAt"`           // 创建时间(毫秒时间戳)
	UpdatedAt         int64                 `gorm:"column:updated_at;autoUpdateTime:milli" json:"updatedAt"`           // 更新时间(毫秒时间戳)
	DeletedAt         soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;index;default:0" json:"-"`       // 软删除(毫秒时间戳)
}

// TableName 指定表名
func (GrowthRecord) TableName() string {
	return "growth_records"
}
