package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/wxlbd/nutri-baby-server/internal/application/dto"
	"github.com/wxlbd/nutri-baby-server/internal/application/service"
	"github.com/wxlbd/nutri-baby-server/pkg/response"
)

// DailyStatsHandler 按日统计处理器
type DailyStatsHandler struct {
	dailyStatsService *service.DailyStatsService
}

// NewDailyStatsHandler 创建按日统计处理器
func NewDailyStatsHandler(dailyStatsService *service.DailyStatsService) *DailyStatsHandler {
	return &DailyStatsHandler{
		dailyStatsService: dailyStatsService,
	}
}

// GetDailyStats 获取按日统计数据
// @Router /v1/babies/:babyId/daily-stats [get]
func (h *DailyStatsHandler) GetDailyStats(c *gin.Context) {
	var req dto.DailyStatsRequest

	// 从路径参数获取 babyId
	req.BabyID = c.Param("babyId")

	// 绑定查询参数
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ErrorWithMessage(c, 1001, "参数错误: "+err.Error())
		return
	}

	openID := c.GetString("openid")

	stats, err := h.dailyStatsService.GetDailyStats(c.Request.Context(), openID, &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, stats)
}
