package repository

import (
	"context"

	"github.com/wxlbd/nutri-baby-server/internal/domain/entity"
)

// VaccinePlanTemplateRepository 疫苗计划模板仓储接口(国家免疫规划标准)
type VaccinePlanTemplateRepository interface {
	// FindAll 查找所有模板
	FindAll(ctx context.Context) ([]*entity.VaccinePlanTemplate, error)
	// FindByID 根据ID查找模板
	FindByID(ctx context.Context, templateID int64) (*entity.VaccinePlanTemplate, error)
	// Create 创建模板(系统初始化)
	Create(ctx context.Context, template *entity.VaccinePlanTemplate) error
	// BatchCreate 批量创建模板
	BatchCreate(ctx context.Context, templates []*entity.VaccinePlanTemplate) error
}

// BabyVaccineScheduleRepository 宝宝疫苗接种日程仓储接口(合并计划和记录)
type BabyVaccineScheduleRepository interface {
	// Create 创建疫苗接种日程
	Create(ctx context.Context, schedule *entity.BabyVaccineSchedule) error

	// FindByID 根据ID查找日程
	FindByID(ctx context.Context, scheduleID int64) (*entity.BabyVaccineSchedule, error)

	// FindByBabyID 查找宝宝的所有疫苗接种日程（分页）
	FindByBabyID(ctx context.Context, babyID int64, page, pageSize int) ([]*entity.BabyVaccineSchedule, error)

	// FindByBabyIDAndStatus 根据状态查找宝宝的疫苗接种日程（分页）
	FindByBabyIDAndStatus(ctx context.Context, babyID int64, status string, page, pageSize int) (total int64, schedules []*entity.BabyVaccineSchedule, err error)

	// Update 更新日程
	Update(ctx context.Context, schedule *entity.BabyVaccineSchedule) error

	// Delete 删除日程(软删除)
	Delete(ctx context.Context, scheduleID int64) error

	// BatchCreate 批量创建日程(从模板初始化)
	BatchCreate(ctx context.Context, schedules []*entity.BabyVaccineSchedule) error

	// InitializeFromTemplates 从模板为宝宝初始化疫苗接种日程
	InitializeFromTemplates(ctx context.Context, babyID int64, createdBy int64) error

	// CountByBabyID 统计宝宝的疫苗接种日程总数
	CountByBabyID(ctx context.Context, babyID int64) (int64, error)

	// CountCompletedByBabyID 统计宝宝已完成接种的数量
	CountCompletedByBabyID(ctx context.Context, babyID int64) (int64, error)

	// MarkAsCompleted 标记日程为已完成
	MarkAsCompleted(ctx context.Context, scheduleID int64, vaccineDate int64, hospital string,
		batchNumber, doctor, reaction, note *string, completedBy int64, completedByName, completedByAvatar string) error

	// MarkAsSkipped 标记日程为跳过
	MarkAsSkipped(ctx context.Context, scheduleID int64) error

	// GetStatistics 获取宝宝疫苗接种统计
	GetStatistics(ctx context.Context, babyID int64) (total, completed, pending, skipped int64, err error)
}
