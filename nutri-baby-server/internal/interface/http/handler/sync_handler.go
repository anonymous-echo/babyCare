package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/wxlbd/nutri-baby-server/pkg/response"
)

// SyncHandler 同步处理器
type SyncHandler struct {
	// TODO: 添加WebSocket同步相关依赖
}

// NewSyncHandler 创建同步处理器
func NewSyncHandler() *SyncHandler {
	return &SyncHandler{}
}

// HandleSync WebSocket同步处理
// @Router /sync [get]
func (h *SyncHandler) HandleSync(c *gin.Context) {
	// TODO: 实现WebSocket同步逻辑
	response.ErrorWithMessage(c, 5001, "WebSocket同步功能待实现")
}
