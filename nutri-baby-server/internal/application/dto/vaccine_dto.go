package dto

// PaginationRequest 通用分页请求参数
type PaginationRequest struct {
	Page     *int `json:"page" form:"page"`
	PageSize *int `json:"pageSize" form:"pageSize"`
}

// GetPageWithDefault 获取页码，使用默认值 1
func (p *PaginationRequest) GetPageWithDefault() int {
	if p == nil || p.Page == nil || *p.Page < 1 {
		return 1
	}
	return *p.Page
}

// GetPageSizeWithDefault 获取分页大小，使用默认值 10，最大 100
func (p *PaginationRequest) GetPageSizeWithDefault() int {
	if p == nil || p.PageSize == nil || *p.PageSize < 1 {
		return 10
	}
	if *p.PageSize > 100 {
		return 100
	}
	return *p.PageSize
}


// VaccineRecordDTO 疫苗接种记录DTO
type VaccineRecordDTO struct {
	RecordID    string  `json:"recordId"`
	BabyID      string  `json:"babyId"`
	PlanID      string  `json:"planId"`
	VaccineType string  `json:"vaccineType"`
	VaccineName string  `json:"vaccineName"`
	DoseNumber  int     `json:"doseNumber"`
	VaccineDate int64   `json:"vaccineDate"`
	Hospital    string  `json:"hospital"`
	BatchNumber *string `json:"batchNumber,omitempty"`
	Doctor      *string `json:"doctor,omitempty"`
	Reaction    *string `json:"reaction,omitempty"`
	Note        *string `json:"note,omitempty"`
	CreateBy    string  `json:"createBy"`
	CreateTime  int64   `json:"createTime"`
}

// CreateVaccineRecordRequest 创建疫苗记录请求
type CreateVaccineRecordRequest struct {
	PlanID      string  `json:"planId" binding:"required"`      // 疫苗计划ID
	VaccineType string  `json:"vaccineType" binding:"required"` // 疫苗类型
	VaccineName string  `json:"vaccineName" binding:"required"` // 疫苗名称
	DoseNumber  int     `json:"doseNumber" binding:"required"`  // 剂次
	VaccineDate int64   `json:"vaccineDate" binding:"required"` // 接种日期
	BatchNumber *string `json:"batchNumber" binding:"required"` // 批号
	Doctor      *string `json:"doctor" binding:"required"`      // 医生
	Reaction    *string `json:"reaction" binding:"required"`    // 反应
	Note        *string `json:"note" binding:"required"`        // 备注
}

// VaccineReminderDTO 疫苗提醒DTO
type VaccineReminderDTO struct {
	ReminderID    string `json:"reminderId"`
	BabyID        string `json:"babyId"`
	BabyName      string `json:"babyName,omitempty"`
	PlanID        string `json:"planId"`
	VaccineName   string `json:"vaccineName"`
	DoseNumber    int    `json:"doseNumber"`
	ScheduledDate int64  `json:"scheduledDate"`
	Status        string `json:"status"`
	DaysUntilDue  int    `json:"daysUntilDue"`
	ReminderSent  bool   `json:"reminderSent"`
	CreateTime    int64  `json:"createTime"`
}

// VaccineStatisticsDTO 疫苗统计DTO
type VaccineStatisticsDTO struct {
	Total         int64              `json:"total"`
	Completed     int64              `json:"completed"`
	Pending       int64              `json:"pending"`
	Overdue       int64              `json:"overdue"`
	Percentage    int                `json:"percentage"`
	NextVaccine   *NextVaccineDTO    `json:"nextVaccine,omitempty"`
	RecentRecords []VaccineRecordDTO `json:"recentRecords"`
}

// NextVaccineDTO 下一个待接种疫苗
type NextVaccineDTO struct {
	VaccineName   string `json:"vaccineName"`
	DoseNumber    int    `json:"doseNumber"`
	ScheduledDate int64  `json:"scheduledDate"`
	DaysUntilDue  int    `json:"daysUntilDue"`
}

// ===================================================================
// 新DTO: BabyVaccineSchedule 相关 (合并计划和记录)
// ===================================================================

// VaccineScheduleDTO 疫苗接种日程DTO
type VaccineScheduleDTO struct {
	ScheduleID        string  `json:"scheduleId"`
	BabyID            string  `json:"babyId"`
	TemplateID        *string `json:"templateId,omitempty"`
	VaccineType       string  `json:"vaccineType"`
	VaccineName       string  `json:"vaccineName"`
	Description       string  `json:"description,omitempty"`
	AgeInMonths       int     `json:"ageInMonths"`
	DoseNumber        int     `json:"doseNumber"`
	IsRequired        bool    `json:"isRequired"`
	ReminderDays      int     `json:"reminderDays"`
	IsCustom          bool    `json:"isCustom"`
	VaccinationStatus string  `json:"vaccinationStatus"` // pending, completed, skipped
	// 接种记录信息(仅在 status='completed' 时有值)
	VaccineDate       *int64  `json:"vaccineDate,omitempty"`
	Hospital          *string `json:"hospital,omitempty"`
	BatchNumber       *string `json:"batchNumber,omitempty"`
	Doctor            *string `json:"doctor,omitempty"`
	Reaction          *string `json:"reaction,omitempty"`
	Note              *string `json:"note,omitempty"`
	CompletedBy       *string `json:"completedBy,omitempty"`
	CompletedByName   *string `json:"completedByName,omitempty"`
	CompletedByAvatar *string `json:"completedByAvatar,omitempty"`
	CompletedTime     *int64  `json:"completedTime,omitempty"`
	CreateBy          string  `json:"createBy"`
	CreateTime        int64   `json:"createTime"`
}

// CreateVaccineScheduleRequest 创建疫苗接种日程请求(自定义疫苗)
type CreateVaccineScheduleRequest struct {
	VaccineType  string  `json:"vaccineType" binding:"required"`
	VaccineName  string  `json:"vaccineName" binding:"required"`
	Description  string  `json:"description"`
	AgeInMonths  int     `json:"ageInMonths" binding:"required"`
	DoseNumber   int     `json:"doseNumber" binding:"required"`
	IsRequired   bool    `json:"isRequired"`
	ReminderDays int     `json:"reminderDays"`
	TemplateID   *string `json:"templateId"` // 可选:基于模板创建
}

// UpdateVaccineScheduleRequest 更新疫苗接种日程请求(记录接种)
type UpdateVaccineScheduleRequest struct {
	VaccineDate       int64   `json:"vaccineDate"`       // 接种日期
	Hospital          string  `json:"hospital"`          // 接种医院
	BatchNumber       *string `json:"batchNumber"`       // 批号(可选)
	Doctor            *string `json:"doctor"`            // 医生(可选)
	Reaction          *string `json:"reaction"`          // 不良反应(可选)
	Note              *string `json:"note"`              // 备注(可选)
	VaccinationStatus string  `json:"vaccinationStatus"` // pending, completed, skipped
}

// UpdateScheduleInfoRequest 更新疫苗接种日程基本信息请求(仅限未完成的日程)
type UpdateScheduleInfoRequest struct {
	VaccineType  *string `json:"vaccineType"`  // 疫苗类型
	VaccineName  *string `json:"vaccineName"`  // 疫苗名称
	Description  *string `json:"description"`  // 描述
	AgeInMonths  *int    `json:"ageInMonths"`  // 接种月龄
	DoseNumber   *int    `json:"doseNumber"`   // 剂次
	IsRequired   *bool   `json:"isRequired"`   // 是否必打
	ReminderDays *int    `json:"reminderDays"` // 提醒天数
}

// VaccineScheduleStatisticsDTO 疫苗接种日程统计DTO（纯统计数据）
type VaccineScheduleStatisticsDTO struct {
	Total            int64 `json:"total"`
	Completed        int64 `json:"completed"`
	Pending          int64 `json:"pending"`
	Skipped          int64 `json:"skipped"`
	CompletionRate   int   `json:"completionRate"`
}

type GetVaccineScheduleListRequest struct {
	BabyID string `json:"babyId" uri:"babyId"`
	OpenID string `json:"openId"`
	Status string `json:"status" form:"status"`
	PaginationRequest
}

