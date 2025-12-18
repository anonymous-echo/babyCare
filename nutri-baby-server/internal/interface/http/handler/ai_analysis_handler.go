package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/wxlbd/nutri-baby-server/internal/application/service"
	"github.com/wxlbd/nutri-baby-server/internal/domain/entity"
	"github.com/wxlbd/nutri-baby-server/pkg/errors"
	"github.com/wxlbd/nutri-baby-server/pkg/response"
)

// AIAnalysisHandler AI分析处理器
type AIAnalysisHandler struct {
	aiAnalysisService service.AIAnalysisService
	logger            *zap.Logger
}

// NewAIAnalysisHandler 创建AI分析处理器
func NewAIAnalysisHandler(
	aiAnalysisService service.AIAnalysisService,
	logger *zap.Logger,
) *AIAnalysisHandler {
	return &AIAnalysisHandler{
		aiAnalysisService: aiAnalysisService,
		logger:            logger,
	}
}

// AnalyzeWithTools 使用工具调用进行AI分析
// @Summary 使用工具调用进行AI分析
// @Description 使用增强的工具调用机制进行AI分析
// @Tags 增强AI分析
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param request body service.CreateAnalysisRequest true "分析请求"
// @Success 200 {object} response.Response{data=service.AnalysisResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/ai/enhanced/analysis [post]
func (h *AIAnalysisHandler) CreateAnalysis(c *gin.Context) {
	var req service.CreateAnalysisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMessage(c, 1001, "参数错误: "+err.Error())
		return
	}

	// 验证权限
	if err := h.checkPermission(c, req.BabyID); err != nil {
		response.Error(c, err)
		return
	}

	result, err := h.aiAnalysisService.CreateAnalysis(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("使用工具调用创建AI分析任务失败",
			zap.Error(err),
			zap.Int64("baby_id", req.BabyID),
			zap.String("analysis_type", string(req.AnalysisType)),
		)
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}

// GenerateDailyTipsWithTools 使用工具调用生成每日建议
// @Summary 使用工具调用生成每日建议
// @Description 使用增强的工具调用机制生成AI育儿建议
// @Tags 增强AI分析
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param baby_id query int true "宝宝ID"
// @Param date query string false "日期 (YYYY-MM-DD)，默认为今天"
// @Success 200 {object} response.Response{data=service.DailyTipsResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/ai/enhanced/daily-tips [post]
func (h *AIAnalysisHandler) GenerateDailyTips(c *gin.Context) {
	babyID, err := strconv.ParseInt(c.Query("baby_id"), 10, 64)
	if err != nil {
		response.ErrorWithMessage(c, 1001, "无效的宝宝ID")
		return
	}

	dateStr := c.Query("date")
	var date time.Time
	if dateStr != "" {
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			response.ErrorWithMessage(c, 1001, "无效的日期格式")
			return
		}
	} else {
		date = time.Now()
	}

	// 验证权限
	if err := h.checkPermission(c, babyID); err != nil {
		response.Error(c, errors.ErrPermissionDenied)
		return
	}

	result, err := h.aiAnalysisService.GenerateDailyTips(c.Request.Context(), strconv.FormatInt(babyID, 10), date)
	if err != nil {
		h.logger.Error("使用工具调用生成每日建议失败",
			zap.Error(err),
			zap.Int64("baby_id", babyID),
			zap.Time("date", date),
		)
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}

// ProcessPendingAnalysesWithTools 使用工具调用处理待分析任务
// @Summary 使用工具调用处理待分析任务
// @Description 使用增强的工具调用机制处理待分析的AI任务
// @Tags 增强AI分析
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/ai/enhanced/process-pending [post]
func (h *AIAnalysisHandler) ProcessPendingAnalyses(c *gin.Context) {
	if err := h.aiAnalysisService.ProcessPendingAnalyses(c.Request.Context()); err != nil {
		h.logger.Error("使用工具调用处理待分析AI任务失败", zap.Error(err))
		response.Error(c, err)
		return
	}
	response.Success(c, gin.H{"message": "处理完成"})
}

// TestToolCalling 测试工具调用功能
// @Summary 测试工具调用功能
// @Description 测试AI模型的工具调用能力
// @Tags 增强AI分析
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param baby_id query int true "宝宝ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/ai/enhanced/test-tools [get]
func (h *AIAnalysisHandler) TestToolCalling(c *gin.Context) {
	babyID, err := strconv.ParseInt(c.Query("baby_id"), 10, 64)
	if err != nil {
		response.ErrorWithMessage(c, 1001, "无效的宝宝ID")
		return
	}

	// 验证权限
	if err := h.checkPermission(c, babyID); err != nil {
		response.Error(c, err)
		return
	}

	// 创建一个测试分析请求
	req := &service.CreateAnalysisRequest{
		BabyID:       babyID,
		AnalysisType: entity.AIAnalysisTypeFeeding,
		StartDate:    service.CustomTime{Time: time.Now().AddDate(0, 0, -7)},
		EndDate:      service.CustomTime{Time: time.Now()},
	}

	result, err := h.aiAnalysisService.CreateAnalysis(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("工具调用测试失败",
			zap.Error(err),
			zap.Int64("baby_id", babyID),
		)
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"message": "工具调用测试成功",
		"result":  result,
	})
}

// GetAnalysisResult 获取分析结果
func (h *AIAnalysisHandler) GetAnalysisResult(c *gin.Context) {
	analysisID := c.Param("id")
	result, err := h.aiAnalysisService.GetAnalysisResult(c.Request.Context(), analysisID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, result)
}

// GetAnalysisStatus 获取分析状态（用于轮询）
// @Summary 获取分析状态
// @Description 获取指定分析任务的当前状态，用于前端轮询
// @Tags AI分析
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param id path string true "分析任务ID"
// @Success 200 {object} response.Response{data=service.AnalysisStatusResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/ai/analysis/{id}/status [get]
func (h *AIAnalysisHandler) GetAnalysisStatus(c *gin.Context) {
	analysisID := c.Param("id")

	status, err := h.aiAnalysisService.GetAnalysisStatus(c.Request.Context(), analysisID)
	if err != nil {
		h.logger.Error("获取分析状态失败",
			zap.String("analysis_id", analysisID),
			zap.Error(err),
		)
		response.Error(c, err)
		return
	}

	response.Success(c, status)
}

// GetLatestAnalysis 获取最新分析
func (h *AIAnalysisHandler) GetLatestAnalysis(c *gin.Context) {
	babyID := c.Param("babyId")
	analysisType := entity.AIAnalysisType(c.Query("type"))

	id, err := strconv.ParseInt(babyID, 10, 64)
	if err != nil {
		response.ErrorWithMessage(c, 1001, "无效的宝宝ID")
		return
	}

	if err := h.checkPermission(c, id); err != nil {
		response.Error(c, err)
		return
	}

	result, err := h.aiAnalysisService.GetLatestAnalysis(c.Request.Context(), babyID, analysisType)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, result)
}

// GetAnalysisStats 获取分析统计
func (h *AIAnalysisHandler) GetAnalysisStats(c *gin.Context) {
	babyID := c.Param("babyId")
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))

	id, err := strconv.ParseInt(babyID, 10, 64)
	if err != nil {
		response.ErrorWithMessage(c, 1001, "无效的宝宝ID")
		return
	}

	if err := h.checkPermission(c, id); err != nil {
		response.Error(c, err)
		return
	}

	result, err := h.aiAnalysisService.GetAnalysisStats(c.Request.Context(), babyID, days)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, result)
}

// BatchAnalyze 批量分析
func (h *AIAnalysisHandler) BatchAnalyze(c *gin.Context) {
	var req service.BatchAnalysisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMessage(c, 1001, "参数错误: "+err.Error())
		return
	}

	if err := h.checkPermission(c, req.BabyID); err != nil {
		response.Error(c, err)
		return
	}

	result, err := h.aiAnalysisService.BatchAnalyze(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, result)
}

// GetDailyTips 获取每日建议
func (h *AIAnalysisHandler) GetDailyTips(c *gin.Context) {
	babyID := c.Param("babyId")
	dateStr := c.DefaultQuery("date", time.Now().Format("2006-01-02"))

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		response.ErrorWithMessage(c, 1001, "无效的日期格式")
		return
	}

	id, err := strconv.ParseInt(babyID, 10, 64)
	if err != nil {
		response.ErrorWithMessage(c, 1001, "无效的宝宝ID")
		return
	}

	if err := h.checkPermission(c, id); err != nil {
		response.Error(c, err)
		return
	}

	result, err := h.aiAnalysisService.GetDailyTips(c.Request.Context(), babyID, date)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, result)
}

// checkPermission 检查权限
func (h *AIAnalysisHandler) checkPermission(c *gin.Context, babyID int64) error {
	// 从上下文中获取当前用户的openid (由auth中间件设置)
	openid, exists := c.Get("openid")
	if !exists {
		return errors.ErrUnauthorized
	}

	openidStr, ok := openid.(string)
	if !ok || openidStr == "" {
		return errors.ErrUnauthorized
	}

	// 检查用户是否有权限访问该宝宝的数据
	// 这里需要调用权限检查服务，暂时直接返回nil
	// TODO: 实现完整的权限检查逻辑
	_ = openidStr
	_ = babyID
	return nil
}

// RegisterAIAnalysisRoutes 注册AI分析路由
func RegisterAIAnalysisRoutes(router *gin.RouterGroup, handler *AIAnalysisHandler) {
	aiGroup := router.Group("/ai")
	{
		aiGroup.POST("/analysis", handler.CreateAnalysis)
		aiGroup.GET("/analysis/:id/status", handler.GetAnalysisStatus)
		aiGroup.GET("/analysis/:id", handler.GetAnalysisResult)
		aiGroup.GET("/latest/:babyId", handler.GetLatestAnalysis)
		aiGroup.GET("/stats/:babyId", handler.GetAnalysisStats)
		aiGroup.POST("/batch-analyze", handler.BatchAnalyze)
		aiGroup.POST("/daily-tips", handler.GenerateDailyTips)
		aiGroup.GET("/daily-tips/:babyId", handler.GetDailyTips)
		aiGroup.POST("/process-pending", handler.ProcessPendingAnalyses)
		aiGroup.GET("/test-tools", handler.TestToolCalling)
	}
}
