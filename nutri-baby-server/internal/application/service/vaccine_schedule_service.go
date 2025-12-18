package service

import (
	"context"
	"sort"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/wxlbd/nutri-baby-server/internal/application/dto"
	"github.com/wxlbd/nutri-baby-server/internal/domain/entity"
	"github.com/wxlbd/nutri-baby-server/internal/domain/repository"
	"github.com/wxlbd/nutri-baby-server/pkg/errors"
)

// VaccineScheduleService 疫苗接种日程服务(新)
type VaccineScheduleService struct {
	scheduleRepo     repository.BabyVaccineScheduleRepository
	babyRepo         repository.BabyRepository
	collaboratorRepo repository.BabyCollaboratorRepository
	userRepository   repository.UserRepository
	logger           *zap.Logger
}

// NewVaccineScheduleService 创建疫苗接种日程服务实例
func NewVaccineScheduleService(
	scheduleRepo repository.BabyVaccineScheduleRepository,
	babyRepo repository.BabyRepository,
	collaboratorRepo repository.BabyCollaboratorRepository,
	userRepository repository.UserRepository,
	logger *zap.Logger,
) *VaccineScheduleService {
	return &VaccineScheduleService{
		scheduleRepo:     scheduleRepo,
		babyRepo:         babyRepo,
		collaboratorRepo: collaboratorRepo,
		userRepository:   userRepository,
		logger:           logger,
	}
}

// GetVaccineSchedules 获取宝宝的疫苗接种日程列表（支持状态过滤和分页）
// 职责：仅返回日程列表，不包含统计信息
func (s *VaccineScheduleService) GetVaccineSchedules(ctx context.Context, request *dto.GetVaccineScheduleListRequest) (int64, []dto.VaccineScheduleDTO, error) {
	// 1. 验证权限
	if err := s.checkPermission(ctx, request.BabyID, request.OpenID); err != nil {
		return 0, nil, err
	}

	// 转换 babyID 从 string 到 int64
	babyIDInt64, err := strconv.ParseInt(request.BabyID, 10, 64)
	if err != nil {
		return 0, nil, errors.New(errors.ParamError, "invalid baby id format")
	}

	// 2. 如果指定了状态，验证并按状态查询（带分页）
	// 验证状态值
	if request.Status != "" {
		if request.Status != entity.VaccinationStatusPending &&
			request.Status != entity.VaccinationStatusCompleted &&
			request.Status != entity.VaccinationStatusSkipped {
			return 0, nil, errors.New(errors.ParamError, "无效的状态值")
		}
	}

	total, schedules, err := s.scheduleRepo.FindByBabyIDAndStatus(ctx, babyIDInt64, request.Status, request.GetPageWithDefault(), request.GetPageSizeWithDefault())
	if err != nil {
		return 0, nil, err
	}

	// 转换为DTO
	scheduleDTOs := make([]dto.VaccineScheduleDTO, 0, len(schedules))
	for _, schedule := range schedules {
		scheduleDTOs = append(scheduleDTOs, s.toScheduleDTO(schedule))
	}
	return total, scheduleDTOs, nil
}

// UpdateVaccineSchedule 更新疫苗接种日程(记录接种或跳过)
func (s *VaccineScheduleService) UpdateVaccineSchedule(
	ctx context.Context,
	babyID, scheduleID, openID string,
	req *dto.UpdateVaccineScheduleRequest,
) error {
	// 1. 验证权限
	if err := s.checkPermission(ctx, babyID, openID); err != nil {
		return err
	}

	// 转换 IDs from string to int64
	babyIDInt64, err := strconv.ParseInt(babyID, 10, 64)
	if err != nil {
		return errors.New(errors.ParamError, "invalid baby id format")
	}

	scheduleIDInt64, err := strconv.ParseInt(scheduleID, 10, 64)
	if err != nil {
		return errors.New(errors.ParamError, "invalid schedule id format")
	}

	// 2. 查询日程是否存在
	schedule, err := s.scheduleRepo.FindByID(ctx, scheduleIDInt64)
	if err != nil {
		return err
	}

	// 3. 验证日程是否属于该宝宝
	if schedule.BabyID != babyIDInt64 {
		return errors.New(errors.PermissionDenied, "无权操作该疫苗接种日程")
	}

	// 4. 检查是否已完成
	if schedule.IsCompleted() {
		return errors.New(errors.Conflict, "该疫苗接种日程已完成,无法重复记录")
	}

	// 获取用户信息
	user, err := s.getUserInfo(ctx, openID)
	if err != nil {
		return err
	}

	// 5. 根据请求的状态处理
	if req.VaccinationStatus == entity.VaccinationStatusSkipped {
		// 5a. 标记为跳过
		err = s.scheduleRepo.MarkAsSkipped(ctx, scheduleIDInt64)
		if err != nil {
			return err
		}

		return nil
	}

	// 5b. 标记为已完成 (completed)
	if req.VaccinationStatus == entity.VaccinationStatusCompleted {
		// 验证必填字段
		if req.VaccineDate == 0 {
			return errors.New(errors.ParamError, "接种日期不能为空")
		}
		if req.Hospital == "" {
			return errors.New(errors.ParamError, "接种医院不能为空")
		}

		err = s.scheduleRepo.MarkAsCompleted(
			ctx,
			scheduleIDInt64,
			req.VaccineDate,
			req.Hospital,
			req.BatchNumber,
			req.Doctor,
			req.Reaction,
			req.Note,
			user.ID,
			user.NickName,
			user.AvatarURL,
		)
		if err != nil {
			return err
		}

		return nil
	}

	// 6. 无效的状态
	return errors.New(errors.ParamError, "无效的接种状态")
}

// CreateCustomSchedule 创建自定义疫苗接种日程
func (s *VaccineScheduleService) CreateCustomSchedule(
	ctx context.Context,
	babyID, openID string,
	req *dto.CreateVaccineScheduleRequest,
) error {
	// 1. 验证权限
	if err := s.checkPermission(ctx, babyID, openID); err != nil {
		return err
	}

	// 转换 babyID from string to int64
	babyIDInt64, err := strconv.ParseInt(babyID, 10, 64)
	if err != nil {
		return errors.New(errors.ParamError, "invalid baby id format")
	}

	// 2. 获取宝宝信息以计算 scheduled_date
	baby, err := s.babyRepo.FindByID(ctx, babyIDInt64)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return errors.New(errors.BabyNotFound, "宝宝不存在")
		}
		return errors.Wrap(errors.DatabaseError, "查询宝宝信息失败", err)
	}

	// 获取用户信息
	user, err := s.getUserInfo(ctx, openID)
	if err != nil {
		return err
	}

	// 3. 解析出生日期并计算预定接种日期
	birthTime, err := time.Parse("2006-01-02", baby.BirthDate)
	if err != nil {
		return errors.Wrap(errors.ParamError, "解析出生日期失败", err)
	}

	// 计算预定接种日期 = 出生日期 + 接种月龄
	scheduledDate := birthTime.AddDate(0, req.AgeInMonths, 0).UnixMilli()

	// 4. 创建日程实体
	var templateID *int64
	if req.TemplateID != nil {
		templateIDInt64, err := strconv.ParseInt(*req.TemplateID, 10, 64)
		if err == nil {
			templateID = &templateIDInt64
		}
	}

	schedule := &entity.BabyVaccineSchedule{
		BabyID:            babyIDInt64,
		TemplateID:        templateID,
		VaccineType:       req.VaccineType,
		VaccineName:       req.VaccineName,
		Description:       req.Description,
		AgeInMonths:       req.AgeInMonths,
		DoseNumber:        req.DoseNumber,
		IsRequired:        req.IsRequired,
		ReminderDays:      req.ReminderDays,
		IsCustom:          true, // 标记为自定义
		VaccinationStatus: entity.VaccinationStatusPending,
		ScheduledDate:     scheduledDate, // 设置计划接种日期
		ReminderSent:      false,         // 初始化提醒状态
		CompletedBy:       &user.ID,
		CompletedByName:   &user.NickName,
		CompletedByAvatar: &user.AvatarURL,
	}

	// 5. 保存日程
	return s.scheduleRepo.Create(ctx, schedule)
}

// UpdateScheduleInfo 更新疫苗接种日程基本信息(仅限未完成的日程)
func (s *VaccineScheduleService) UpdateScheduleInfo(
	ctx context.Context,
	babyID, scheduleID, openID string,
	req *dto.UpdateScheduleInfoRequest,
) error {
	// 1. 验证权限
	if err := s.checkPermission(ctx, babyID, openID); err != nil {
		return err
	}

	// 转换 IDs from string to int64
	babyIDInt64, err := strconv.ParseInt(babyID, 10, 64)
	if err != nil {
		return errors.New(errors.ParamError, "invalid baby id format")
	}

	scheduleIDInt64, err := strconv.ParseInt(scheduleID, 10, 64)
	if err != nil {
		return errors.New(errors.ParamError, "invalid schedule id format")
	}

	// 2. 查询日程是否存在
	schedule, err := s.scheduleRepo.FindByID(ctx, scheduleIDInt64)
	if err != nil {
		return err
	}

	// 3. 验证日程是否属于该宝宝
	if schedule.BabyID != babyIDInt64 {
		return errors.New(errors.PermissionDenied, "无权操作该疫苗接种日程")
	}

	// 4. 检查是否已完成或已跳过（只允许编辑待接种的日程）
	if schedule.VaccinationStatus != entity.VaccinationStatusPending {
		return errors.New(errors.Conflict, "只能编辑待接种状态的疫苗日程")
	}

	// 5. 更新字段（只更新非nil的字段）
	needRecalculateDate := false

	if req.VaccineType != nil {
		schedule.VaccineType = *req.VaccineType
	}
	if req.VaccineName != nil {
		schedule.VaccineName = *req.VaccineName
	}
	if req.Description != nil {
		schedule.Description = *req.Description
	}
	if req.AgeInMonths != nil {
		schedule.AgeInMonths = *req.AgeInMonths
		needRecalculateDate = true // 月龄变化需要重新计算接种日期
	}
	if req.DoseNumber != nil {
		schedule.DoseNumber = *req.DoseNumber
	}
	if req.IsRequired != nil {
		schedule.IsRequired = *req.IsRequired
	}
	if req.ReminderDays != nil {
		schedule.ReminderDays = *req.ReminderDays
	}

	// 6. 如果月龄变化，重新计算 scheduled_date
	if needRecalculateDate {
		baby, err := s.babyRepo.FindByID(ctx, babyIDInt64)
		if err != nil {
			return errors.Wrap(errors.DatabaseError, "查询宝宝信息失败", err)
		}

		birthTime, err := time.Parse("2006-01-02", baby.BirthDate)
		if err != nil {
			return errors.Wrap(errors.ParamError, "解析出生日期失败", err)
		}

		schedule.ScheduledDate = birthTime.AddDate(0, schedule.AgeInMonths, 0).UnixMilli()
	}

	// 7. 更新到数据库
	return s.scheduleRepo.Update(ctx, schedule)
}

// DeleteSchedule 删除疫苗接种日程(仅限自定义日程)
func (s *VaccineScheduleService) DeleteSchedule(ctx context.Context, babyID, scheduleID, openID string) error {
	// 1. 验证权限
	if err := s.checkPermission(ctx, babyID, openID); err != nil {
		return err
	}

	// 转换 IDs from string to int64
	babyIDInt64, err := strconv.ParseInt(babyID, 10, 64)
	if err != nil {
		return errors.New(errors.ParamError, "invalid baby id format")
	}

	scheduleIDInt64, err := strconv.ParseInt(scheduleID, 10, 64)
	if err != nil {
		return errors.New(errors.ParamError, "invalid schedule id format")
	}

	// 2. 查询日程
	schedule, err := s.scheduleRepo.FindByID(ctx, scheduleIDInt64)
	if err != nil {
		return err
	}

	// 3. 验证所属
	if schedule.BabyID != babyIDInt64 {
		return errors.New(errors.PermissionDenied, "无权操作该疫苗接种日程")
	}

	// 4. 检查是否为自定义日程
	if !schedule.IsCustom {
		return errors.New(errors.Conflict, "不能删除系统预设的疫苗接种日程")
	}

	// 5. 检查是否已完成
	if schedule.IsCompleted() {
		return errors.New(errors.Conflict, "不能删除已完成的疫苗接种日程")
	}

	// 6. 删除日程
	return s.scheduleRepo.Delete(ctx, scheduleIDInt64)
}

// GetStatistics 获取疫苗接种统计
func (s *VaccineScheduleService) GetStatistics(ctx context.Context, babyID, openID string) (*dto.VaccineScheduleStatisticsDTO, error) {
	// 1. 验证权限
	if err := s.checkPermission(ctx, babyID, openID); err != nil {
		return nil, err
	}

	// 转换 babyID from string to int64
	babyIDInt64, err := strconv.ParseInt(babyID, 10, 64)
	if err != nil {
		return nil, errors.New(errors.ParamError, "invalid baby id format")
	}

	// 2. 获取统计数据
	total, completed, pending, skipped, err := s.scheduleRepo.GetStatistics(ctx, babyIDInt64)
	if err != nil {
		return nil, err
	}

	// 3. 计算完成百分比 (已完成 + 已跳过 都算作已处理)
	percentage := 0
	if total > 0 {
		processed := completed + skipped // 已完成 + 已跳过
		percentage = int(float64(processed) / float64(total) * 100)
	}

	return &dto.VaccineScheduleStatisticsDTO{
		Total:          total,
		Completed:      completed,
		Pending:        pending,
		Skipped:        skipped,
		CompletionRate: percentage,
	}, nil
}

// ===================================================================
// 辅助方法
// ===================================================================

// checkPermission 检查用户是否有权限访问该宝宝的疫苗信息
func (s *VaccineScheduleService) checkPermission(ctx context.Context, babyID, openID string) error {
	// 转换 babyID from string to int64
	babyIDInt64, err := strconv.ParseInt(babyID, 10, 64)
	if err != nil {
		return errors.New(errors.ParamError, "invalid baby id format")
	}

	// 获取用户信息以获取用户ID
	user, err := s.getUserInfo(ctx, openID)
	if err != nil {
		return err
	}

	collaborator, err := s.collaboratorRepo.FindByBabyAndUser(ctx, babyIDInt64, user.ID)
	if err != nil {
		return err
	}
	if collaborator == nil {
		return errors.New(errors.PermissionDenied, "无权访问该宝宝的疫苗信息")
	}

	return nil
}

// getUserInfo 获取用户信息
func (s *VaccineScheduleService) getUserInfo(ctx context.Context, openID string) (*entity.User, error) {
	return s.userRepository.FindByOpenID(ctx, openID)
}

// toScheduleDTO 将实体转换为DTO
func (s *VaccineScheduleService) toScheduleDTO(schedule *entity.BabyVaccineSchedule) dto.VaccineScheduleDTO {
	// 将 ID 转换为字符串
	scheduleID := strconv.FormatInt(schedule.ID, 10)
	babyID := strconv.FormatInt(schedule.BabyID, 10)

	// 转换 TemplateID
	var templateID *string
	if schedule.TemplateID != nil {
		templateIDStr := strconv.FormatInt(*schedule.TemplateID, 10)
		templateID = &templateIDStr
	}

	// 转换 CompletedBy
	var completedBy *string
	if schedule.CompletedBy != nil {
		completedByStr := strconv.FormatInt(*schedule.CompletedBy, 10)
		completedBy = &completedByStr
	}

	return dto.VaccineScheduleDTO{
		ScheduleID:        scheduleID,
		BabyID:            babyID,
		TemplateID:        templateID,
		VaccineType:       schedule.VaccineType,
		VaccineName:       schedule.VaccineName,
		Description:       schedule.Description,
		AgeInMonths:       schedule.AgeInMonths,
		DoseNumber:        schedule.DoseNumber,
		IsRequired:        schedule.IsRequired,
		ReminderDays:      schedule.ReminderDays,
		IsCustom:          schedule.IsCustom,
		VaccinationStatus: schedule.VaccinationStatus,
		VaccineDate:       schedule.VaccineDate,
		Hospital:          schedule.Hospital,
		BatchNumber:       schedule.BatchNumber,
		Doctor:            schedule.Doctor,
		Reaction:          schedule.Reaction,
		Note:              schedule.Note,
		CompletedBy:       completedBy,
		CompletedByName:   schedule.CompletedByName,
		CompletedByAvatar: schedule.CompletedByAvatar,
		CompletedTime:     schedule.CompletedTime,
		CreateBy:          scheduleID, // 创建者是创建这条日程的人
		CreateTime:        schedule.CreatedAt,
	}
}

// ===================================================================
// 新增功能: 初始化和提醒管理 (合并自 VaccineService 和 VaccinePlanService)
// ===================================================================

// InitializeSchedulesForBaby 为新宝宝初始化疫苗接种日程(从模板)
// 此方法应该在创建宝宝后立即调用,为宝宝生成完整的疫苗接种计划
func (s *VaccineScheduleService) InitializeSchedulesForBaby(ctx context.Context, babyID, openID string) error {
	// 转换 babyID from string to int64
	babyIDInt64, err := strconv.ParseInt(babyID, 10, 64)
	if err != nil {
		return errors.New(errors.ParamError, "invalid baby id format")
	}

	// 获取用户信息
	user, err := s.getUserInfo(ctx, openID)
	if err != nil {
		return err
	}

	// 1. 验证宝宝是否存在
	baby, err := s.babyRepo.FindByID(ctx, babyIDInt64)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return errors.New(errors.NotFound, "宝宝不存在")
		}
		return errors.Wrap(errors.DatabaseError, "查询宝宝信息失败", err)
	}

	// 2. 检查权限
	if baby.UserID != user.ID {
		return errors.New(errors.PermissionDenied, "无权为该宝宝初始化疫苗日程")
	}

	// 3. 检查是否已初始化
	count, err := s.scheduleRepo.CountByBabyID(ctx, babyIDInt64)
	if err != nil {
		return errors.Wrap(errors.DatabaseError, "查询宝宝信息失败", err)
	}

	if count > 0 {
		// 已初始化,跳过
		return nil
	}

	// 4. 从模板初始化日程(Repository 层会自动计算 scheduled_date)
	err = s.scheduleRepo.InitializeFromTemplates(ctx, babyIDInt64, user.ID)
	if err != nil {
		return errors.Wrap(errors.DatabaseError, "从模板初始化疫苗日程失败", err)
	}

	return nil
}

// GetVaccineReminders 获取疫苗提醒列表
// 只返回需要提醒的疫苗 (due 和 overdue 状态)
func (s *VaccineScheduleService) GetVaccineReminders(
	ctx context.Context,
	babyID, openID string,
) ([]*dto.VaccineReminderDTO, error) {
	// 1. 验证权限
	if err := s.checkPermission(ctx, babyID, openID); err != nil {
		return nil, err
	}

	// 转换 babyID from string to int64
	babyIDInt64, err := strconv.ParseInt(babyID, 10, 64)
	if err != nil {
		return nil, errors.New(errors.ParamError, "invalid baby id format")
	}

	// 2. 查询待接种的日程
	_, schedules, err := s.scheduleRepo.FindByBabyIDAndStatus(ctx, babyIDInt64, entity.VaccinationStatusPending, 1, 1000)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseError, "查询疫苗日程失败", err)
	}

	// 3. 转换为提醒DTO,只返回需要提醒的疫苗 (due/overdue)
	result := make([]*dto.VaccineReminderDTO, 0)
	for _, schedule := range schedules {
		// 实时计算提醒状态
		reminderStatus := schedule.GetReminderStatus()

		// 只返回需要提醒的疫苗 (due 或 overdue)
		// upcoming: 距离接种日期 > 7天,不需要提醒
		// completed: 已接种/已跳过,不需要提醒
		if reminderStatus != entity.ReminderStatusDue && reminderStatus != entity.ReminderStatusOverdue {
			continue
		}

		// 将 ID 转换为字符串
		scheduleID := strconv.FormatInt(schedule.ID, 10)
		babyIDStr := strconv.FormatInt(schedule.BabyID, 10)

		result = append(result, &dto.VaccineReminderDTO{
			ReminderID:    scheduleID, // 使用 scheduleID 作为 reminderID
			BabyID:        babyIDStr,
			PlanID:        scheduleID,
			VaccineName:   schedule.VaccineName,
			DoseNumber:    schedule.DoseNumber,
			ScheduledDate: schedule.ScheduledDate,
			Status:        reminderStatus,
			DaysUntilDue:  schedule.DaysUntilDue(),
			ReminderSent:  schedule.ReminderSent,
			CreateTime:    schedule.CreatedAt,
		})
	}

	// 4. 按 scheduled_date 排序(最近的在前)
	sort.Slice(result, func(i, j int) bool {
		return result[i].ScheduledDate < result[j].ScheduledDate
	})

	return result, nil
}
