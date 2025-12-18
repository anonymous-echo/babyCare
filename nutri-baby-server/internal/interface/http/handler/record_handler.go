package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/wxlbd/nutri-baby-server/internal/application/dto"
	"github.com/wxlbd/nutri-baby-server/internal/application/service"
	"github.com/wxlbd/nutri-baby-server/pkg/response"
)

// RecordHandler 记录处理器
type RecordHandler struct {
	feedingService  *service.FeedingRecordService
	sleepService    *service.SleepRecordService
	diaperService   *service.DiaperRecordService
	growthService   *service.GrowthRecordService
	timelineService *service.TimelineService
}

// NewRecordHandler 创建记录处理器
func NewRecordHandler(
	feedingService *service.FeedingRecordService,
	sleepService *service.SleepRecordService,
	diaperService *service.DiaperRecordService,
	growthService *service.GrowthRecordService,
	timelineService *service.TimelineService,
) *RecordHandler {
	return &RecordHandler{
		feedingService:  feedingService,
		sleepService:    sleepService,
		diaperService:   diaperService,
		growthService:   growthService,
		timelineService: timelineService,
	}
}

// CreateFeedingRecord 创建喂养记录
// @Router /feeding-records [post]
func (h *RecordHandler) CreateFeedingRecord(c *gin.Context) {
	var req dto.CreateFeedingRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMessage(c, 1001, "参数错误: "+err.Error())
		return
	}

	openID := c.GetString("openid")

	record, err := h.feedingService.CreateFeedingRecord(c.Request.Context(), openID, &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, record)
}

// GetFeedingRecords 获取喂养记录列表
// @Router /feeding-records [get]
func (h *RecordHandler) GetFeedingRecords(c *gin.Context) {
	query := h.parseRecordQuery(c)
	openID := c.GetString("openid")

	records, total, err := h.feedingService.GetFeedingRecords(c.Request.Context(), openID, query)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"records":  records,
		"total":    total,
		"page":     query.Page,
		"pageSize": query.PageSize,
	})
}

// CreateSleepRecord 创建睡眠记录
// @Router /sleep-records [post]
func (h *RecordHandler) CreateSleepRecord(c *gin.Context) {
	var req dto.CreateSleepRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMessage(c, 1001, "参数错误: "+err.Error())
		return
	}

	openID := c.GetString("openid")

	record, err := h.sleepService.CreateSleepRecord(c.Request.Context(), openID, &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, record)
}

// GetSleepRecords 获取睡眠记录列表
// @Router /sleep-records [get]
func (h *RecordHandler) GetSleepRecords(c *gin.Context) {
	query := h.parseRecordQuery(c)
	openID := c.GetString("openid")

	records, total, err := h.sleepService.GetSleepRecords(c.Request.Context(), openID, query)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"records":  records,
		"total":    total,
		"page":     query.Page,
		"pageSize": query.PageSize,
	})
}

// CreateDiaperRecord 创建尿布记录
// @Router /diaper-records [post]
func (h *RecordHandler) CreateDiaperRecord(c *gin.Context) {
	var req dto.CreateDiaperRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMessage(c, 1001, "参数错误: "+err.Error())
		return
	}

	openID := c.GetString("openid")

	record, err := h.diaperService.CreateDiaperRecord(c.Request.Context(), openID, &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, record)
}

// GetDiaperRecords 获取尿布记录列表
// @Router /diaper-records [get]
func (h *RecordHandler) GetDiaperRecords(c *gin.Context) {
	query := h.parseRecordQuery(c)
	openID := c.GetString("openid")

	records, total, err := h.diaperService.GetDiaperRecords(c.Request.Context(), openID, query)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"records":  records,
		"total":    total,
		"page":     query.Page,
		"pageSize": query.PageSize,
	})
}

// CreateGrowthRecord 创建生长记录
// @Router /growth-records [post]
func (h *RecordHandler) CreateGrowthRecord(c *gin.Context) {
	var req dto.CreateGrowthRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMessage(c, 1001, "参数错误: "+err.Error())
		return
	}

	openID := c.GetString("openid")

	record, err := h.growthService.CreateGrowthRecord(c.Request.Context(), openID, &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, record)
}

// GetGrowthRecords 获取生长记录列表
// @Router /growth-records [get]
func (h *RecordHandler) GetGrowthRecords(c *gin.Context) {
	query := h.parseRecordQuery(c)
	openID := c.GetString("openid")

	records, total, err := h.growthService.GetGrowthRecords(c.Request.Context(), openID, query)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"records":  records,
		"total":    total,
		"page":     query.Page,
		"pageSize": query.PageSize,
	})
}

// GetTimeline 获取时间线记录
// @Router /timeline [get]
func (h *RecordHandler) GetTimeline(c *gin.Context) {
	var query dto.TimelineQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.ErrorWithMessage(c, 1001, "参数错误: "+err.Error())
		return
	}

	openID := c.GetString("openid")

	result, err := h.timelineService.GetTimeline(c.Request.Context(), openID, &query)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}

// parseRecordQuery 解析记录查询参数
func (h *RecordHandler) parseRecordQuery(c *gin.Context) *dto.RecordListQuery {
	query := &dto.RecordListQuery{
		BabyID:    c.Query("babyId"),
		StartTime: parseInt64(c.Query("startTime")),
		EndTime:   parseInt64(c.Query("endTime")),
	}

	return query
}

// parseInt 解析整数
func parseInt(s string, defaultValue int) int {
	if s == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return val
}

// parseInt64 解析int64
func parseInt64(s string) int64 {
	if s == "" {
		return 0
	}
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return val
}

// GetFeedingRecordById 获取单条喂养记录
// @Router /feeding-records/:id [get]
func (h *RecordHandler) GetFeedingRecordById(c *gin.Context) {
	recordID := c.Param("id")
	openID := c.GetString("openid")

	record, err := h.feedingService.GetFeedingRecordById(c.Request.Context(), openID, recordID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, record)
}

// UpdateFeedingRecord 更新喂养记录
// @Router /feeding-records/:id [put]
func (h *RecordHandler) UpdateFeedingRecord(c *gin.Context) {
	recordID := c.Param("id")
	openID := c.GetString("openid")

	var req dto.UpdateFeedingRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMessage(c, 1001, "参数错误: "+err.Error())
		return
	}

	record, err := h.feedingService.UpdateFeedingRecord(c.Request.Context(), openID, recordID, &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, record)
}

// DeleteFeedingRecord 删除喂养记录
// @Router /feeding-records/:id [delete]
func (h *RecordHandler) DeleteFeedingRecord(c *gin.Context) {
	recordID := c.Param("id")
	openID := c.GetString("openid")

	if err := h.feedingService.DeleteFeedingRecord(c.Request.Context(), openID, recordID); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, nil)
}

// GetSleepRecordById 获取单条睡眠记录
// @Router /sleep-records/:id [get]
func (h *RecordHandler) GetSleepRecordById(c *gin.Context) {
	recordID := c.Param("id")
	openID := c.GetString("openid")

	record, err := h.sleepService.GetSleepRecordById(c.Request.Context(), openID, recordID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, record)
}

// UpdateSleepRecord 更新睡眠记录
// @Router /sleep-records/:id [put]
func (h *RecordHandler) UpdateSleepRecord(c *gin.Context) {
	recordID := c.Param("id")
	openID := c.GetString("openid")

	var req dto.UpdateSleepRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMessage(c, 1001, "参数错误: "+err.Error())
		return
	}

	record, err := h.sleepService.UpdateSleepRecord(c.Request.Context(), openID, recordID, &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, record)
}

// DeleteSleepRecord 删除睡眠记录
// @Router /sleep-records/:id [delete]
func (h *RecordHandler) DeleteSleepRecord(c *gin.Context) {
	recordID := c.Param("id")
	openID := c.GetString("openid")

	if err := h.sleepService.DeleteSleepRecord(c.Request.Context(), openID, recordID); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, nil)
}

// GetDiaperRecordById 获取单条尿布记录
// @Router /diaper-records/:id [get]
func (h *RecordHandler) GetDiaperRecordById(c *gin.Context) {
	recordID := c.Param("id")
	openID := c.GetString("openid")

	record, err := h.diaperService.GetDiaperRecordById(c.Request.Context(), openID, recordID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, record)
}

// UpdateDiaperRecord 更新尿布记录
// @Router /diaper-records/:id [put]
func (h *RecordHandler) UpdateDiaperRecord(c *gin.Context) {
	recordID := c.Param("id")
	openID := c.GetString("openid")

	var req dto.UpdateDiaperRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMessage(c, 1001, "参数错误: "+err.Error())
		return
	}

	record, err := h.diaperService.UpdateDiaperRecord(c.Request.Context(), openID, recordID, &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, record)
}

// DeleteDiaperRecord 删除尿布记录
// @Router /diaper-records/:id [delete]
func (h *RecordHandler) DeleteDiaperRecord(c *gin.Context) {
	recordID := c.Param("id")
	openID := c.GetString("openid")

	if err := h.diaperService.DeleteDiaperRecord(c.Request.Context(), openID, recordID); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, nil)
}

// GetGrowthRecordById 获取单条成长记录
// @Router /growth-records/:id [get]
func (h *RecordHandler) GetGrowthRecordById(c *gin.Context) {
	recordID := c.Param("id")
	openID := c.GetString("openid")

	record, err := h.growthService.GetGrowthRecordById(c.Request.Context(), openID, recordID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, record)
}

// UpdateGrowthRecord 更新成长记录
// @Router /growth-records/:id [put]
func (h *RecordHandler) UpdateGrowthRecord(c *gin.Context) {
	recordID := c.Param("id")
	openID := c.GetString("openid")

	var req dto.UpdateGrowthRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMessage(c, 1001, "参数错误: "+err.Error())
		return
	}

	record, err := h.growthService.UpdateGrowthRecord(c.Request.Context(), openID, recordID, &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, record)
}

// DeleteGrowthRecord 删除成长记录
// @Router /growth-records/:id [delete]
func (h *RecordHandler) DeleteGrowthRecord(c *gin.Context) {
	recordID := c.Param("id")
	openID := c.GetString("openid")

	if err := h.growthService.DeleteGrowthRecord(c.Request.Context(), openID, recordID); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, nil)
}
