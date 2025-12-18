package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/wxlbd/nutri-baby-server/internal/application/dto"
	"github.com/wxlbd/nutri-baby-server/internal/application/service"
	"github.com/wxlbd/nutri-baby-server/pkg/response"
)

// BabyHandler 宝宝处理器 (去家庭化架构)
type BabyHandler struct {
	babyService   *service.BabyService
	wechatService *service.WechatService
}

// NewBabyHandler 创建宝宝处理器
func NewBabyHandler(babyService *service.BabyService, wechatService *service.WechatService) *BabyHandler {
	return &BabyHandler{
		babyService:   babyService,
		wechatService: wechatService,
	}
}

// CreateBaby 创建宝宝
// @Router /v1/babies [post]
func (h *BabyHandler) CreateBaby(c *gin.Context) {
	var req dto.CreateBabyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMessage(c, 1001, "参数错误: "+err.Error())
		return
	}

	openID := c.GetString("openid")

	baby, err := h.babyService.CreateBaby(c.Request.Context(), openID, &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, baby)
}

// GetUserBabies 获取用户可访问的宝宝列表
// @Router /v1/babies [get]
func (h *BabyHandler) GetUserBabies(c *gin.Context) {
	openID := c.GetString("openid")

	babies, err := h.babyService.GetUserBabies(c.Request.Context(), openID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, babies)
}

// GetBabyDetail 获取宝宝详情
// @Router /v1/babies/:babyId [get]
func (h *BabyHandler) GetBabyDetail(c *gin.Context) {
	babyID := c.Param("babyId")
	openID := c.GetString("openid")

	baby, err := h.babyService.GetBabyDetail(c.Request.Context(), babyID, openID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, baby)
}

// UpdateBaby 更新宝宝
// @Router /v1/babies/:babyId [put]
func (h *BabyHandler) UpdateBaby(c *gin.Context) {
	babyID := c.Param("babyId")
	openID := c.GetString("openid")

	var req dto.UpdateBabyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMessage(c, 1001, "参数错误: "+err.Error())
		return
	}

	if err := h.babyService.UpdateBaby(c.Request.Context(), babyID, openID, &req); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, nil)
}

// DeleteBaby 删除宝宝
// @Router /v1/babies/:babyId [delete]
func (h *BabyHandler) DeleteBaby(c *gin.Context) {
	babyID := c.Param("babyId")
	openID := c.GetString("openid")

	if err := h.babyService.DeleteBaby(c.Request.Context(), babyID, openID); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, nil)
}

// GetCollaborators 获取宝宝的协作者列表
// @Router /v1/babies/:babyId/collaborators [get]
func (h *BabyHandler) GetCollaborators(c *gin.Context) {
	babyID := c.Param("babyId")
	openID := c.GetString("openid")

	collaborators, err := h.babyService.GetCollaborators(c.Request.Context(), babyID, openID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, collaborators)
}

// InviteCollaborator 邀请协作者
// @Router /v1/babies/:babyId/collaborators/invite [post]
func (h *BabyHandler) InviteCollaborator(c *gin.Context) {
	babyID := c.Param("babyId")
	openID := c.GetString("openid")

	var req dto.InviteCollaboratorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMessage(c, 1001, "参数错误: "+err.Error())
		return
	}

	invitation, err := h.babyService.InviteCollaborator(c.Request.Context(), babyID, openID, &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, invitation)
}

// JoinBaby 加入宝宝协作
// @Router /v1/babies/join [post]
func (h *BabyHandler) JoinBaby(c *gin.Context) {
	openID := c.GetString("openid")

	var req dto.JoinBabyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMessage(c, 1001, "参数错误: "+err.Error())
		return
	}

	baby, err := h.babyService.JoinBaby(c.Request.Context(), openID, &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, baby)
}

// RemoveCollaborator 移除协作者
// @Router /v1/babies/:babyId/collaborators/:openid [delete]
func (h *BabyHandler) RemoveCollaborator(c *gin.Context) {
	babyID := c.Param("babyId")
	targetOpenID := c.Param("openid")
	openID := c.GetString("openid")

	if err := h.babyService.RemoveCollaborator(c.Request.Context(), babyID, openID, targetOpenID); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, nil)
}

// UpdateCollaboratorRole 更新协作者角色
// @Router /v1/babies/:babyId/collaborators/:openid/role [put]
func (h *BabyHandler) UpdateCollaboratorRole(c *gin.Context) {
	babyID := c.Param("babyId")
	targetOpenID := c.Param("openid")
	openID := c.GetString("openid")

	var req struct {
		Role string `json:"role" binding:"required,oneof=admin editor viewer"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMessage(c, 1001, "参数错误: "+err.Error())
		return
	}

	if err := h.babyService.UpdateCollaboratorRole(c.Request.Context(), babyID, openID, targetOpenID, req.Role); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, nil)
}

// UpdateFamilyMember 更新亲友团成员信息
// @Router /v1/babies/:babyId/collaborators/:openid [put]
func (h *BabyHandler) UpdateFamilyMember(c *gin.Context) {
	babyID := c.Param("babyId")
	targetOpenID := c.Param("openid")
	openID := c.GetString("openid")

	var req dto.UpdateFamilyMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMessage(c, 1001, "参数错误: "+err.Error())
		return
	}

	if err := h.babyService.UpdateFamilyMember(c.Request.Context(), babyID, openID, targetOpenID, &req); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, nil)
}

// GenerateInviteQRCode 生成邀请小程序码
// @Router /v1/babies/:babyId/qrcode [get]
func (h *BabyHandler) GenerateInviteQRCode(c *gin.Context) {
	shortCode := c.Query("shortCode") // 从查询参数获取短码

	if shortCode == "" {
		response.ErrorWithMessage(c, 1001, "shortCode参数不能为空")
		return
	}

	// 构建场景值
	scene := "c=" + shortCode

	// 调用微信服务生成小程序码
	qrcodeURL, err := h.wechatService.GenerateQRCode(c.Request.Context(), scene, "pages/baby/join/join")
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"qrcodeUrl": qrcodeURL,
		"scene":     scene,
	})
}

// GetInvitationByShortCode 通过短码获取邀请详情
// @Router /v1/invitations/code/:shortCode [get]
func (h *BabyHandler) GetInvitationByShortCode(c *gin.Context) {
	shortCode := c.Param("shortCode")

	if shortCode == "" {
		response.ErrorWithMessage(c, 1001, "shortCode参数不能为空")
		return
	}

	invitation, err := h.babyService.GetInvitationByShortCode(c.Request.Context(), shortCode)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, invitation)
}
