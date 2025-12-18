package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/wxlbd/nutri-baby-server/internal/application/service"
	"github.com/wxlbd/nutri-baby-server/pkg/response"
)

// StatisticsHandler 统计处理器
type StatisticsHandler struct {
	statisticsService *service.StatisticsService
}

// NewStatisticsHandler 创建统计处理器
func NewStatisticsHandler(statisticsService *service.StatisticsService) *StatisticsHandler {
	return &StatisticsHandler{
		statisticsService: statisticsService,
	}
}

// GetBabyStatistics 获取宝宝统计数据
// @Router /v1/babies/:babyId/statistics [get]
func (h *StatisticsHandler) GetBabyStatistics(c *gin.Context) {
	babyID := c.Param("babyId")
	openID := c.GetString("openid")

	statistics, err := h.statisticsService.GetBabyStatistics(c.Request.Context(), babyID, openID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, statistics)
}
