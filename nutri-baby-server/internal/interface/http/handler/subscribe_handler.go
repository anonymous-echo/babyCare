package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/wxlbd/nutri-baby-server/internal/application/dto"
	"github.com/wxlbd/nutri-baby-server/internal/application/service"
	"github.com/wxlbd/nutri-baby-server/pkg/errors"
	"github.com/wxlbd/nutri-baby-server/pkg/response"
)

// SubscribeHandler 订阅消息处理器
type SubscribeHandler struct {
	subscribeService *service.SubscribeService
	logger           *zap.Logger
}

// NewSubscribeHandler 创建订阅消息处理器
func NewSubscribeHandler(
	subscribeService *service.SubscribeService,
	logger *zap.Logger,
) *SubscribeHandler {
	return &SubscribeHandler{
		subscribeService: subscribeService,
		logger:           logger,
	}
}

// UploadAuth 上传订阅授权记录
// @Summary 上传订阅授权记录
// @Description 前端用户授权后,上传授权结果到后端
// @Tags Subscribe
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body dto.SubscribeAuthRequest true "授权记录"
// @Success 200 {object} response.Response{data=dto.SubscribeAuthResponse}
// @Router /api/v1/subscribe/auth [post]
func (h *SubscribeHandler) UploadAuth(c *gin.Context) {
	// 获取用户openid
	openid, exists := c.Get("openid")
	if !exists {
		response.Error(c, errors.ErrUnauthorized)
		return
	}

	var req dto.SubscribeAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request", zap.Error(err))
		response.Error(c, errors.ErrParamInvalid)
		return
	}

	result, err := h.subscribeService.SaveSubscribeAuth(c.Request.Context(), openid.(string), req.Records)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}

// GetStatus 获取用户订阅状态
// @Summary 获取用户订阅状态
// @Description 查询用户所有的订阅消息状态
// @Tags Subscribe
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} response.Response{data=dto.SubscribeStatusResponse}
// @Router /api/v1/subscribe/status [get]
func (h *SubscribeHandler) GetStatus(c *gin.Context) {
	openid, exists := c.Get("openid")
	if !exists {
		response.Error(c, errors.ErrUnauthorized)
		return
	}

	result, err := h.subscribeService.GetUserSubscriptions(c.Request.Context(), openid.(string))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}

// Cancel 取消订阅
// @Summary 取消订阅
// @Description 用户取消某个类型的订阅消息
// @Tags Subscribe
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body dto.CancelSubscriptionRequest true "取消订阅请求"
// @Success 200 {object} response.Response
// @Router /api/v1/subscribe/cancel [delete]
func (h *SubscribeHandler) Cancel(c *gin.Context) {
	openid, exists := c.Get("openid")
	if !exists {
		response.Error(c, errors.ErrUnauthorized)
		return
	}
	fmt.Println(openid)
	var req dto.CancelSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request", zap.Error(err))
		response.Error(c, errors.ErrParamInvalid)
		return
	}

	//err := h.subscribeService.CancelSubscription(c.Request.Context(), openid.(string), req.TemplateType)
	//if err != nil {
	//	response.Error(c, err)
	//	return
	//}

	response.Success(c, gin.H{"message": "订阅已取消"})
}

// GetLogs 获取消息发送日志
// @Summary 获取消息发送日志
// @Description 查询用户的消息发送历史记录
// @Tags Subscribe
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param offset query int false "偏移量" default(0)
// @Param limit query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=dto.MessageLogsResponse}
// @Router /api/v1/subscribe/logs [get]
func (h *SubscribeHandler) GetLogs(c *gin.Context) {
	openid, exists := c.Get("openid")
	if !exists {
		response.Error(c, errors.ErrUnauthorized)
		return
	}

	offset := 0
	limit := 20

	if offsetParam := c.Query("offset"); offsetParam != "" {
		if _, err := fmt.Sscanf(offsetParam, "%d", &offset); err != nil {
			h.logger.Warn("Invalid offset parameter", zap.Error(err))
		}
	}

	if limitParam := c.Query("limit"); limitParam != "" {
		if _, err := fmt.Sscanf(limitParam, "%d", &limit); err != nil {
			h.logger.Warn("Invalid limit parameter", zap.Error(err))
		}
	}

	// 限制最大返回数量
	if limit > 100 {
		limit = 100
	}

	result, err := h.subscribeService.GetMessageLogs(c.Request.Context(), openid.(string), offset, limit)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}
