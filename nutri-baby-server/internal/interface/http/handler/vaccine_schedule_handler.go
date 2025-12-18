package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/wxlbd/nutri-baby-server/pkg/errors"

	"github.com/wxlbd/nutri-baby-server/internal/application/dto"
	"github.com/wxlbd/nutri-baby-server/internal/application/service"
	"github.com/wxlbd/nutri-baby-server/pkg/response"
)

// VaccineScheduleHandler 疫苗接种日程处理器(新)
type VaccineScheduleHandler struct {
	scheduleService *service.VaccineScheduleService
}

// NewVaccineScheduleHandler 创建疫苗接种日程处理器实例
func NewVaccineScheduleHandler(scheduleService *service.VaccineScheduleService) *VaccineScheduleHandler {
	return &VaccineScheduleHandler{
		scheduleService: scheduleService,
	}
}

// GetVaccineSchedules 获取宝宝的疫苗接种日程列表
// @Summary 获取疫苗接种日程列表
// @Tags 疫苗管理
// @Param babyId path string true "宝宝ID"
// @Param status query string false "状态过滤: pending/completed/skipped"
// @Param page query int false "页码，默认1"
// @Param pageSize query int false "每页数量，默认10，最多100"
// @Success 200 {object} []dto.VaccineScheduleDTO
// @Router /babies/{babyId}/vaccine-schedules [get]
func (h *VaccineScheduleHandler) GetVaccineSchedules(c *gin.Context) {
	babyID := c.Param("babyId")
	openID := c.GetString("openid")
	var req dto.GetVaccineScheduleListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ErrorWithMessage(c, errors.ParamError, fmt.Sprintf("参数错误：%s", err))
		return
	}
	req.BabyID = babyID
	req.OpenID = openID

	// 统一调用 GetVaccineSchedules，支持状态过滤和分页
	total, schedules, err := h.scheduleService.GetVaccineSchedules(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"schedules": schedules,
		"total":     total,
		"page":      req.GetPageWithDefault(),
		"pageSize":  req.GetPageSizeWithDefault(),
		"hasMore":   int(total)/req.GetPageSizeWithDefault() > req.GetPageWithDefault(),
	})
}

// UpdateVaccineSchedule 更新疫苗接种日程(记录接种)
// @Summary 记录疫苗接种
// @Tags 疫苗管理
// @Param babyId path string true "宝宝ID"
// @Param scheduleId path string true "日程ID"
// @Param request body dto.UpdateVaccineScheduleRequest true "接种记录"
// @Success 200 {object} response.Response
// @Router /babies/{babyId}/vaccine-schedules/{scheduleId} [put]
func (h *VaccineScheduleHandler) UpdateVaccineSchedule(c *gin.Context) {
	babyID := c.Param("babyId")
	scheduleID := c.Param("scheduleId")
	openID, _ := c.Get("openid")

	var req dto.UpdateVaccineScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMessage(c, errors.ParamError, fmt.Sprintf("参数错误：%s", err))
		return
	}

	err := h.scheduleService.UpdateVaccineSchedule(
		c.Request.Context(),
		babyID,
		scheduleID,
		openID.(string),
		&req,
	)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, "记录疫苗接种成功")
}

// CreateCustomSchedule 创建自定义疫苗接种日程
// @Summary 创建自定义疫苗计划
// @Tags 疫苗管理
// @Param babyId path string true "宝宝ID"
// @Param request body dto.CreateVaccineScheduleRequest true "疫苗计划"
// @Success 200 {object} response.Response
// @Router /babies/{babyId}/vaccine-schedules [post]
func (h *VaccineScheduleHandler) CreateCustomSchedule(c *gin.Context) {
	babyID := c.Param("babyId")
	openID, _ := c.Get("openid")

	var req dto.CreateVaccineScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMessage(c, errors.ParamError, fmt.Sprintf("参数错误：%s", err))
		return
	}

	err := h.scheduleService.CreateCustomSchedule(
		c.Request.Context(),
		babyID,
		openID.(string),
		&req,
	)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, "创建自定义疫苗计划成功")
}

// UpdateScheduleInfo 更新疫苗接种日程基本信息
// @Summary 更新疫苗接种日程基本信息(仅限未完成的日程)
// @Tags 疫苗管理
// @Param babyId path string true "宝宝ID"
// @Param scheduleId path string true "日程ID"
// @Param data body dto.UpdateScheduleInfoRequest true "更新数据"
// @Success 200 {object} response.Response
// @Router /babies/{babyId}/vaccine-schedules/{scheduleId}/info [patch]
func (h *VaccineScheduleHandler) UpdateScheduleInfo(c *gin.Context) {
	babyID := c.Param("babyId")
	scheduleID := c.Param("scheduleId")
	openID, _ := c.Get("openid")

	var req dto.UpdateScheduleInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMessage(c, errors.ParamError, fmt.Sprintf("参数错误：%s", err))
		return
	}

	err := h.scheduleService.UpdateScheduleInfo(
		c.Request.Context(),
		babyID,
		scheduleID,
		openID.(string),
		&req,
	)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, "更新疫苗计划成功")
}

// DeleteSchedule 删除疫苗接种日程
// @Summary 删除疫苗计划(仅限自定义)
// @Tags 疫苗管理
// @Param babyId path string true "宝宝ID"
// @Param scheduleId path string true "日程ID"
// @Success 200 {object} response.Response
// @Router /babies/{babyId}/vaccine-schedules/{scheduleId} [delete]
func (h *VaccineScheduleHandler) DeleteSchedule(c *gin.Context) {
	babyID := c.Param("babyId")
	scheduleID := c.Param("scheduleId")
	openID, _ := c.Get("openid")

	err := h.scheduleService.DeleteSchedule(
		c.Request.Context(),
		babyID,
		scheduleID,
		openID.(string),
	)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, "删除疫苗计划成功")
}

// GetStatistics 获取疫苗接种统计
// @Summary 获取疫苗接种统计
// @Description 获取宝宝的疫苗接种统计数据（纯统计，不包含任何记录列表）
// @Tags 疫苗管理
// @Param babyId path string true "宝宝ID"
// @Success 200 {object} dto.VaccineScheduleStatisticsDTO
// @Router /babies/{babyId}/vaccine-statistics [get]
func (h *VaccineScheduleHandler) GetStatistics(c *gin.Context) {
	babyID := c.Param("babyId")
	openID, _ := c.Get("openid")

	stats, err := h.scheduleService.GetStatistics(c.Request.Context(), babyID, openID.(string))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, stats)
}

// GetReminders 获取疫苗提醒列表
// @Summary 获取疫苗提醒列表
// @Tags 疫苗管理
// @Param babyId path string true "宝宝ID"
// @Success 200 {object} []dto.VaccineReminderDTO
// @Router /babies/{babyId}/vaccine-reminders [get]
func (h *VaccineScheduleHandler) GetReminders(c *gin.Context) {
	babyID := c.Param("babyId")
	openID, _ := c.Get("openid")

	reminders, err := h.scheduleService.GetVaccineReminders(
		c.Request.Context(),
		babyID,
		openID.(string),
	)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"reminders": reminders,
	})
}
