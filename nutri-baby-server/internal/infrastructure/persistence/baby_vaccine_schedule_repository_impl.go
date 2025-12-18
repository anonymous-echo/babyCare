package persistence

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/wxlbd/nutri-baby-server/internal/domain/entity"
	"github.com/wxlbd/nutri-baby-server/internal/domain/repository"
	"github.com/wxlbd/nutri-baby-server/pkg/errors"
)

type babyVaccineScheduleRepositoryImpl struct {
	db                      *gorm.DB
	vaccinePlanTemplateRepo repository.VaccinePlanTemplateRepository
}

// NewBabyVaccineScheduleRepository 创建宝宝疫苗接种日程仓储实例
func NewBabyVaccineScheduleRepository(
	db *gorm.DB,
	vaccinePlanTemplateRepo repository.VaccinePlanTemplateRepository,
) repository.BabyVaccineScheduleRepository {
	return &babyVaccineScheduleRepositoryImpl{
		db:                      db,
		vaccinePlanTemplateRepo: vaccinePlanTemplateRepo,
	}
}

// Create 创建疫苗接种日程
func (r *babyVaccineScheduleRepositoryImpl) Create(ctx context.Context, schedule *entity.BabyVaccineSchedule) error {
	if schedule.ID == 0 {
		schedule.ID = 0 // ID由snowflake生成
	}

	err := r.db.WithContext(ctx).Create(schedule).Error
	if err != nil {
		return errors.Wrap(errors.DatabaseError, "创建疫苗接种日程失败", err)
	}
	return nil
}

// FindByID 根据ID查找日程
func (r *babyVaccineScheduleRepositoryImpl) FindByID(ctx context.Context, scheduleID int64) (*entity.BabyVaccineSchedule, error) {
	var schedule entity.BabyVaccineSchedule
	err := r.db.WithContext(ctx).
		Where("id = ?", scheduleID).
		First(&schedule).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(errors.NotFound, "疫苗接种日程不存在")
		}
		return nil, errors.Wrap(errors.DatabaseError, "查询疫苗接种日程失败", err)
	}
	return &schedule, nil
}

// FindByBabyID 查找宝宝的所有疫苗接种日程（分页）
func (r *babyVaccineScheduleRepositoryImpl) FindByBabyID(ctx context.Context, babyID int64, page, pageSize int) ([]*entity.BabyVaccineSchedule, error) {
	var schedules []*entity.BabyVaccineSchedule
	err := r.db.WithContext(ctx).
		Where("baby_id = ?", babyID).
		Order("age_in_months ASC, dose_number ASC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&schedules).Error

	if err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "查询宝宝疫苗接种日程失败", err)
	}
	return schedules, nil
}

// FindByBabyIDAndStatus 根据状态查找宝宝的疫苗接种日程（分页）
func (r *babyVaccineScheduleRepositoryImpl) FindByBabyIDAndStatus(ctx context.Context, babyID int64, status string, page, pageSize int) (int64, []*entity.BabyVaccineSchedule, error) {
	var schedules []*entity.BabyVaccineSchedule
	var total int64
	query := r.db.Model(&entity.BabyVaccineSchedule{}).WithContext(ctx).
		Where("baby_id = ?", babyID).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if status != "" {
				db.Where("vaccination_status = ?", status)
			}
			return db
		})
	err := query.Count(&total).Error
	if err != nil {
		return 0, nil, errors.Wrap(errors.DatabaseError, "查询疫苗接种日程总数失败", err)
	}
	err = query.Order("age_in_months ASC, dose_number ASC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&schedules).Error

	if err != nil {
		return 0, nil, errors.Wrap(errors.DatabaseError, "查询宝宝疫苗接种日程失败", err)
	}
	return total, schedules, nil
}

// Update 更新日程
func (r *babyVaccineScheduleRepositoryImpl) Update(ctx context.Context, schedule *entity.BabyVaccineSchedule) error {
	err := r.db.WithContext(ctx).
		Where("id = ?", schedule.ID).
		Updates(schedule).Error

	if err != nil {
		return errors.Wrap(errors.DatabaseError, "更新疫苗接种日程失败", err)
	}
	return nil
}

// Delete 删除日程(软删除)
func (r *babyVaccineScheduleRepositoryImpl) Delete(ctx context.Context, scheduleID int64) error {
	err := r.db.WithContext(ctx).
		Model(&entity.BabyVaccineSchedule{}).
		Where("id = ?", scheduleID).
		Delete(&entity.BabyVaccineSchedule{}).Error

	if err != nil {
		return errors.Wrap(errors.DatabaseError, "删除疫苗接种日程失败", err)
	}
	return nil
}

// BatchCreate 批量创建日程(从模板初始化)
func (r *babyVaccineScheduleRepositoryImpl) BatchCreate(ctx context.Context, schedules []*entity.BabyVaccineSchedule) error {
	if len(schedules) == 0 {
		return nil
	}

	err := r.db.WithContext(ctx).Create(&schedules).Error
	if err != nil {
		return errors.Wrap(errors.DatabaseError, "批量创建疫苗接种日程失败", err)
	}
	return nil
}

// InitializeFromTemplates 从模板为宝宝初始化疫苗接种日程
func (r *babyVaccineScheduleRepositoryImpl) InitializeFromTemplates(ctx context.Context, babyID, createdBy int64) error {
	// 1. 获取宝宝信息(需要出生日期)
	var baby entity.Baby
	err := r.db.WithContext(ctx).Where("id = ?", babyID).First(&baby).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(errors.NotFound, "宝宝不存在")
		}
		return errors.Wrap(errors.DatabaseError, "查询宝宝信息失败", err)
	}

	// 2. 解析出生日期
	birthTime, err := time.Parse("2006-01-02", baby.BirthDate)
	if err != nil {
		return errors.Wrap(errors.ParamError, "解析出生日期失败", err)
	}

	// 3. 获取所有模板
	templates, err := r.vaccinePlanTemplateRepo.FindAll(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.Wrap(errors.NotFound, "没有找到任何疫苗计划模板", err)
		}
		return errors.Wrap(errors.DatabaseError, "获取疫苗计划模板失败", err)
	}

	if len(templates) == 0 {
		return errors.New(errors.NotFound, "未找到疫苗计划模板")
	}

	// 4. 根据模板创建日程,同时计算 scheduled_date
	schedules := make([]*entity.BabyVaccineSchedule, 0, len(templates))
	for _, template := range templates {
		// 计算预定接种日期 (出生日期 + age_in_months 个月)
		scheduledDate := birthTime.AddDate(0, template.AgeInMonths, 0).UnixMilli()

		schedule := &entity.BabyVaccineSchedule{
			BabyID:            babyID,
			TemplateID:        &template.ID,
			VaccineType:       template.VaccineType,
			VaccineName:       template.VaccineName,
			Description:       template.Description,
			AgeInMonths:       template.AgeInMonths,
			DoseNumber:        template.DoseNumber,
			IsRequired:        template.IsRequired,
			ReminderDays:      template.ReminderDays,
			IsCustom:          false,
			VaccinationStatus: entity.VaccinationStatusPending,
			ScheduledDate:     scheduledDate, // 设置计划接种日期
			ReminderSent:      false,         // 初始化提醒状态
			CreatedBy:         createdBy,
		}
		schedules = append(schedules, schedule)
	}

	// 5. 批量插入
	return r.BatchCreate(ctx, schedules)
}

// CountByBabyID 统计宝宝的疫苗接种日程总数
func (r *babyVaccineScheduleRepositoryImpl) CountByBabyID(ctx context.Context, babyID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.BabyVaccineSchedule{}).
		Where("baby_id = ?", babyID).
		Count(&count).Error

	if err != nil {
		return 0, errors.Wrap(errors.DatabaseError, "统计疫苗接种日程失败", err)
	}
	return count, nil
}

// CountCompletedByBabyID 统计宝宝已完成接种的数量
func (r *babyVaccineScheduleRepositoryImpl) CountCompletedByBabyID(ctx context.Context, babyID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.BabyVaccineSchedule{}).
		Where("baby_id = ? AND vaccination_status = ?", babyID, entity.VaccinationStatusCompleted).
		Count(&count).Error

	if err != nil {
		return 0, errors.Wrap(errors.DatabaseError, "统计已完成接种数量失败", err)
	}
	return count, nil
}

// MarkAsCompleted 标记日程为已完成
func (r *babyVaccineScheduleRepositoryImpl) MarkAsCompleted(
	ctx context.Context,
	scheduleID int64,
	vaccineDate int64,
	hospital string,
	batchNumber, doctor, reaction, note *string,
	completedBy int64,
	completedByName, completedByAvatar string,
) error {
	now := time.Now().UnixMilli()

	updates := map[string]interface{}{
		"vaccination_status":  entity.VaccinationStatusCompleted,
		"vaccine_date":        vaccineDate,
		"hospital":            hospital,
		"completed_by":        completedBy,
		"completed_by_name":   completedByName,
		"completed_by_avatar": completedByAvatar,
		"completed_time":      now,
		"batch_number":        batchNumber,
		"doctor":              doctor,
		"reaction":            reaction,
		"note":                note,
	}

	err := r.db.WithContext(ctx).
		Model(&entity.BabyVaccineSchedule{}).
		Where("id = ?", scheduleID).
		Updates(updates).Error

	if err != nil {
		return errors.Wrap(errors.DatabaseError, "标记疫苗接种日程为已完成失败", err)
	}
	return nil
}

// MarkAsSkipped 标记日程为跳过
func (r *babyVaccineScheduleRepositoryImpl) MarkAsSkipped(ctx context.Context, scheduleID int64) error {
	err := r.db.WithContext(ctx).
		Model(&entity.BabyVaccineSchedule{}).
		Where("id = ?", scheduleID).
		Update("vaccination_status", entity.VaccinationStatusSkipped).Error

	if err != nil {
		return errors.Wrap(errors.DatabaseError, "标记疫苗接种日程为跳过失败", err)
	}
	return nil
}

// GetStatistics 获取宝宝疫苗接种统计
func (r *babyVaccineScheduleRepositoryImpl) GetStatistics(ctx context.Context, babyID int64) (total, completed, pending, skipped int64, err error) {
	// 使用单个查询获取所有统计数据
	type StatusCount struct {
		Status string
		Count  int64
	}

	var statusCounts []StatusCount
	err = r.db.WithContext(ctx).
		Model(&entity.BabyVaccineSchedule{}).
		Select("vaccination_status as status, COUNT(*) as count").
		Where("baby_id = ?", babyID).
		Group("vaccination_status").
		Scan(&statusCounts).Error

	if err != nil {
		return 0, 0, 0, 0, errors.Wrap(errors.DatabaseError, "获取疫苗接种统计失败", err)
	}

	// 解析统计结果
	for _, sc := range statusCounts {
		total += sc.Count
		switch sc.Status {
		case entity.VaccinationStatusCompleted:
			completed = sc.Count
		case entity.VaccinationStatusPending:
			pending = sc.Count
		case entity.VaccinationStatusSkipped:
			skipped = sc.Count
		}
	}

	return total, completed, pending, skipped, nil
}
