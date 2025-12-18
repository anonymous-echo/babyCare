package entity

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

// 疫苗提醒状态常量
const (
	ReminderStatusUpcoming  = "upcoming"  // 即将到期
	ReminderStatusDue       = "due"       // 应接种
	ReminderStatusOverdue   = "overdue"   // 已逾期
	ReminderStatusCompleted = "completed" // 已完成
)

// 疫苗接种状态常量
const (
	VaccinationStatusPending   = "pending"   // 未接种
	VaccinationStatusCompleted = "completed" // 已完成
	VaccinationStatusSkipped   = "skipped"   // 跳过/不接种
)

// VaccinePlanTemplate 疫苗计划模板(国家免疫规划标准)
type VaccinePlanTemplate struct {
	ID           int64                 `gorm:"primaryKey;column:id" json:"id"`                                // 雪花ID主键
	VaccineType  string                `gorm:"column:vaccine_type;type:varchar(32);index" json:"vaccineType"` // 疫苗类型
	VaccineName  string                `gorm:"column:vaccine_name;type:varchar(64)" json:"vaccineName"`       // 疫苗名称
	Description  string                `gorm:"column:description;type:text" json:"description"`               // 描述
	AgeInMonths  int                   `gorm:"column:age_in_months" json:"ageInMonths"`                       // 接种年龄(月)
	DoseNumber   int                   `gorm:"column:dose_number" json:"doseNumber"`                          // 接种次数
	IsRequired   bool                  `gorm:"column:is_required" json:"isRequired"`                          // 是否必接种
	ReminderDays int                   `gorm:"column:reminder_days;default:7" json:"reminderDays"`            // 提前多少天提醒
	SortOrder    int                   `gorm:"column:sort_order;default:0" json:"sortOrder"`                  // 排序
	CreatedAt    int64                 `gorm:"column:created_at;autoCreateTime:milli" json:"createdAt"`       // 创建时间(毫秒时间戳)
	UpdatedAt    int64                 `gorm:"column:updated_at;autoUpdateTime:milli" json:"updatedAt"`       // 更新时间(毫秒时间戳)
	DeletedAt    soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;index;default:0" json:"-"`   // 软删除(毫秒时间戳)
}

// TableName 指定表名
func (VaccinePlanTemplate) TableName() string {
	return "vaccine_plan_templates"
}

// BabyVaccineSchedule 宝宝疫苗接种日程(合并计划和记录)
type BabyVaccineSchedule struct {
	// 主键
	ID int64 `gorm:"primaryKey;column:id" json:"id"` // 雪花ID主键

	// 宝宝和计划基础信息
	BabyID     int64  `gorm:"column:baby_id;index;not null" json:"babyId"`    // 宝宝ID (引用Baby.ID)
	TemplateID *int64 `gorm:"column:template_id" json:"templateId,omitempty"` // 来源模板ID (引用VaccinePlanTemplate.ID, 可选)

	// 疫苗基本信息
	VaccineType  string `gorm:"column:vaccine_type;type:varchar(32);not null" json:"vaccineType"` // 疫苗类型
	VaccineName  string `gorm:"column:vaccine_name;type:varchar(64);not null" json:"vaccineName"` // 疫苗名称
	Description  string `gorm:"column:description;type:text" json:"description"`                  // 描述
	AgeInMonths  int    `gorm:"column:age_in_months;not null" json:"ageInMonths"`                 // 接种年龄(月)
	DoseNumber   int    `gorm:"column:dose_number;not null" json:"doseNumber"`                    // 接种次数
	IsRequired   bool   `gorm:"column:is_required;default:true" json:"isRequired"`                // 是否必接种
	ReminderDays int    `gorm:"column:reminder_days;default:7" json:"reminderDays"`               // 提前多少天提醒
	IsCustom     bool   `gorm:"column:is_custom;default:false" json:"isCustom"`                   // 是否用户自定义

	// 接种状态 (关键字段)
	VaccinationStatus string `gorm:"column:vaccination_status;type:varchar(16);not null;default:'pending';index" json:"vaccinationStatus"`
	// 状态值: pending(未接种), completed(已完成), skipped(跳过/不接种)

	// 接种记录信息 (仅在 status='completed' 时有值)
	VaccineDate       *int64  `gorm:"column:vaccine_date" json:"vaccineDate,omitempty"`                                // 实际接种日期(毫秒时间戳)
	Hospital          *string `gorm:"column:hospital;type:varchar(128)" json:"hospital,omitempty"`                     // 接种医院
	BatchNumber       *string `gorm:"column:batch_number;type:varchar(64)" json:"batchNumber,omitempty"`               // 批次号
	Doctor            *string `gorm:"column:doctor;type:varchar(64)" json:"doctor,omitempty"`                          // 接种医生
	Reaction          *string `gorm:"column:reaction;type:text" json:"reaction,omitempty"`                             // 接种反应
	Note              *string `gorm:"column:note;type:text" json:"note,omitempty"`                                     // 备注
	CompletedBy       *int64  `gorm:"column:completed_by" json:"completedBy,omitempty"`                                // 记录接种的用户ID (引用User.ID)
	CompletedByName   *string `gorm:"column:completed_by_name;type:varchar(64)" json:"completedByName,omitempty"`      // 记录者昵称
	CompletedByAvatar *string `gorm:"column:completed_by_avatar;type:varchar(512)" json:"completedByAvatar,omitempty"` // 记录者头像
	CompletedTime     *int64  `gorm:"column:completed_time" json:"completedTime,omitempty"`                            // 接种记录创建时间(毫秒时间戳)

	// 提醒相关字段 (合并自 VaccineReminder)
	ScheduledDate  int64  `gorm:"column:scheduled_date;index" json:"scheduledDate"`        // 计划接种日期(毫秒时间戳)
	ReminderSent   bool   `gorm:"column:reminder_sent;default:false" json:"reminderSent"`  // 是否已发送提醒
	ReminderSentAt *int64 `gorm:"column:reminder_sent_at" json:"reminderSentAt,omitempty"` // 提醒发送时间(毫秒时间戳)

	// 审计字段
	CreatedBy int64                 `gorm:"column:created_by;not null" json:"createdBy"`                       // 创建者用户ID (引用User.ID)
	CreatedAt int64                 `gorm:"column:created_at;autoCreateTime:milli;default:0" json:"createdAt"` // 创建时间(毫秒时间戳)
	UpdatedAt int64                 `gorm:"column:updated_at;autoUpdateTime:milli;default:0" json:"updatedAt"` // 更新时间(毫秒时间戳)
	DeletedAt soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;index;default:0" json:"-"`       // 软删除(毫秒时间戳)

	// 关联
	Template *VaccinePlanTemplate `gorm:"foreignKey:TemplateID;references:ID" json:"template,omitempty"`
	Baby     *Baby                `gorm:"foreignKey:BabyID;references:ID" json:"baby,omitempty"`
}

// TableName 指定表名
func (BabyVaccineSchedule) TableName() string {
	return "baby_vaccine_schedules"
}

// IsCompleted 判断是否已完成接种
func (s *BabyVaccineSchedule) IsCompleted() bool {
	return s.VaccinationStatus == VaccinationStatusCompleted
}

// IsPending 判断是否待接种
func (s *BabyVaccineSchedule) IsPending() bool {
	return s.VaccinationStatus == VaccinationStatusPending
}

// IsSkipped 判断是否已跳过
func (s *BabyVaccineSchedule) IsSkipped() bool {
	return s.VaccinationStatus == VaccinationStatusSkipped
}

// MarkAsCompleted 标记为已完成
func (s *BabyVaccineSchedule) MarkAsCompleted(vaccineDate int64, hospital string, completedBy int64, completedByName, completedByAvatar string) {
	s.VaccinationStatus = VaccinationStatusCompleted
	s.VaccineDate = &vaccineDate
	s.Hospital = &hospital
	s.CompletedBy = &completedBy
	s.CompletedByName = &completedByName
	s.CompletedByAvatar = &completedByAvatar
	now := time.Now().UnixMilli()
	s.CompletedTime = &now
}

// MarkAsSkipped 标记为跳过
func (s *BabyVaccineSchedule) MarkAsSkipped() {
	s.VaccinationStatus = VaccinationStatusSkipped
}

// GetReminderStatus 获取提醒状态(实时计算)
func (s *BabyVaccineSchedule) GetReminderStatus() string {
	if s.VaccinationStatus != VaccinationStatusPending {
		return ReminderStatusCompleted
	}
	days := s.DaysUntilDue()
	if days > 7 {
		return ReminderStatusUpcoming
	} else if days >= 0 {
		return ReminderStatusDue
	}
	return ReminderStatusOverdue
}

// DaysUntilDue 计算距离应接种日期的天数
func (s *BabyVaccineSchedule) DaysUntilDue() int {
	now := time.Now().UnixMilli()
	diff := s.ScheduledDate - now
	return int(diff / (24 * 60 * 60 * 1000))
}

// MarkReminderSent 标记提醒已发送
func (s *BabyVaccineSchedule) MarkReminderSent() {
	s.ReminderSent = true
	now := time.Now().UnixMilli()
	s.ReminderSentAt = &now
}
