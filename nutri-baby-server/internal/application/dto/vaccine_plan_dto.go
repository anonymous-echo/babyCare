package dto

// VaccinePlanDTO 疫苗计划响应DTO
type VaccinePlanDTO struct {
	PlanID        string `json:"planId"`
	VaccineType   string `json:"vaccineType"`
	VaccineName   string `json:"vaccineName"`
	Description   string `json:"description"`
	AgeInMonths   int    `json:"ageInMonths"`
	DoseNumber    int    `json:"doseNumber"`
	IsRequired    bool   `json:"isRequired"`
	ReminderDays  int    `json:"reminderDays"`
	IsCustom      bool   `json:"isCustom"`                // 是否用户自定义
	TemplateID    string `json:"templateId"`              // 来源模板ID (可选)
	ScheduledDate int64  `json:"scheduledDate,omitempty"` // 预定日期(根据宝宝出生日期计算,可选)
	Status        string `json:"status,omitempty"`        // 状态: pending/completed/overdue (可选)
}

// CreateBabyVaccinePlanRequest 创建宝宝疫苗计划请求
type CreateBabyVaccinePlanRequest struct {
	VaccineType  string `json:"vaccineType" binding:"required"`
	VaccineName  string `json:"vaccineName" binding:"required"`
	Description  string `json:"description"`
	AgeInMonths  int    `json:"ageInMonths" binding:"required,min=0"`
	DoseNumber   int    `json:"doseNumber" binding:"required,min=1"`
	IsRequired   bool   `json:"isRequired"`
	ReminderDays int    `json:"reminderDays" binding:"min=0,max=30"`
}

// UpdateBabyVaccinePlanRequest 更新宝宝疫苗计划请求
type UpdateBabyVaccinePlanRequest struct {
	VaccineName  string `json:"vaccineName"`
	Description  string `json:"description"`
	AgeInMonths  int    `json:"ageInMonths" binding:"min=0"`
	DoseNumber   int    `json:"doseNumber" binding:"min=1"`
	IsRequired   bool   `json:"isRequired"`
	ReminderDays int    `json:"reminderDays" binding:"min=0,max=30"`
}

// InitializeVaccinePlansRequest 初始化疫苗计划请求
type InitializeVaccinePlansRequest struct {
	// 可选:是否强制重新初始化(删除已有计划)
	Force bool `json:"force"`
}

// InitializeVaccinePlansResponse 初始化疫苗计划响应
type InitializeVaccinePlansResponse struct {
	TotalPlans int              `json:"totalPlans"` // 初始化的计划总数
	Plans      []VaccinePlanDTO `json:"plans"`      // 计划列表
	Message    string           `json:"message"`    // 提示信息
}

// GetVaccinePlansResponse 获取疫苗计划列表响应
type GetVaccinePlansResponse struct {
	Plans      []VaccinePlanDTO `json:"plans"`
	Total      int              `json:"total"`
	Completed  int              `json:"completed"`
	Percentage int              `json:"percentage"`
}
