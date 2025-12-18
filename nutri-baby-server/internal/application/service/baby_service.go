package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/wxlbd/nutri-baby-server/internal/application/dto"
	"github.com/wxlbd/nutri-baby-server/internal/domain/entity"
	"github.com/wxlbd/nutri-baby-server/internal/domain/repository"
	"github.com/wxlbd/nutri-baby-server/pkg/errors"
	"github.com/wxlbd/nutri-baby-server/pkg/utils"
	"go.uber.org/zap"
)

// BabyService 宝宝服务 (去家庭化架构)
type BabyService struct {
	babyRepo               repository.BabyRepository
	collaboratorRepo       repository.BabyCollaboratorRepository
	invitationRepo         repository.BabyInvitationRepository
	userRepo               repository.UserRepository
	vaccineScheduleService *VaccineScheduleService
	wechatService          *WechatService
	logger                 *zap.Logger
}

// NewBabyService 创建宝宝服务
func NewBabyService(
	babyRepo repository.BabyRepository,
	collaboratorRepo repository.BabyCollaboratorRepository,
	invitationRepo repository.BabyInvitationRepository,
	userRepo repository.UserRepository,
	vaccineScheduleService *VaccineScheduleService,
	wechatService *WechatService,
	logger *zap.Logger,
) *BabyService {
	return &BabyService{
		babyRepo:               babyRepo,
		collaboratorRepo:       collaboratorRepo,
		invitationRepo:         invitationRepo,
		userRepo:               userRepo,
		vaccineScheduleService: vaccineScheduleService,
		wechatService:          wechatService,
		logger:                 logger,
	}
}

// CreateBaby 创建宝宝
func (s *BabyService) CreateBaby(ctx context.Context, openID string, req *dto.CreateBabyRequest) (*dto.BabyDTO, error) {
	// 验证日期格式
	if _, err := time.Parse(time.DateOnly, req.BirthDate); err != nil {
		return nil, errors.New(errors.ParamError, "出生日期格式错误，应为YYYY-MM-DD")
	}

	// 获取用户信息以获取UserID
	user, err := s.userRepo.FindByOpenID(ctx, openID)
	if err != nil {
		return nil, err
	}

	// 创建宝宝实体 (ID由snowflake自动生成)
	baby := &entity.Baby{
		Name:      req.Name,
		Nickname:  req.Nickname,
		Gender:    req.Gender,
		BirthDate: req.BirthDate,
		AvatarURL: req.AvatarURL,
		UserID:    user.ID,
	}

	// 创建宝宝
	if err := s.babyRepo.Create(ctx, baby); err != nil {
		return nil, err
	}

	// 创建者自动成为管理员
	creator := &entity.BabyCollaborator{
		BabyID:     baby.ID,
		UserID:     user.ID,
		Role:       "admin",
		AccessType: "permanent",
	}

	if err := s.collaboratorRepo.Create(ctx, creator); err != nil {
		return nil, err
	}

	// 如果指定了复制协作者,则批量复制
	if req.CopyCollaboratorsFrom != "" {
		sourceBabyID, err := strconv.ParseInt(req.CopyCollaboratorsFrom, 10, 64)
		if err == nil { // 只在转换成功时执行复制
			if err := s.copyCollaborators(ctx, sourceBabyID, baby.ID, openID); err != nil {
				// 记录错误但不影响创建宝宝
				// logger.Warn("Failed to copy collaborators", zap.Error(err))
			}
		}
	}
	// 初始化疫苗计划
	if err := s.vaccineScheduleService.InitializeSchedulesForBaby(ctx, strconv.FormatInt(baby.ID, 10), openID); err != nil {
		// 记录错误但不影响创建宝宝
		s.logger.Error("初始化疫苗计划失败", zap.Error(err))
	}

	// 如果该用户没有其他宝宝,直接将该宝宝设置为默认宝宝
	if err := s.setDefaultBabyIfNeeded(ctx, openID, baby.ID); err != nil {
		s.logger.Error("设置默认宝宝失败", zap.Error(err))
	}

	return &dto.BabyDTO{
		BabyID:     strconv.FormatInt(baby.ID, 10),
		Name:       baby.Name,
		Nickname:   baby.Nickname,
		Gender:     baby.Gender,
		BirthDate:  baby.BirthDate,
		AvatarURL:  baby.AvatarURL,
		CreatorID:  strconv.FormatInt(baby.UserID, 10),
		CreateTime: baby.CreatedAt,
		UpdateTime: baby.UpdatedAt,
	}, nil
}

// GetUserBabies 获取用户可访问的宝宝列表
func (s *BabyService) GetUserBabies(ctx context.Context, openID string) ([]dto.BabyDTO, error) {
	// 获取用户信息以获取UserID
	user, err := s.userRepo.FindByOpenID(ctx, openID)
	if err != nil {
		return nil, err
	}

	babies, err := s.babyRepo.FindByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	result := make([]dto.BabyDTO, 0, len(babies))
	for _, baby := range babies {
		result = append(result, dto.BabyDTO{
			BabyID:     strconv.FormatInt(baby.ID, 10),
			Name:       baby.Name,
			Nickname:   baby.Nickname,
			Gender:     baby.Gender,
			BirthDate:  baby.BirthDate,
			AvatarURL:  baby.AvatarURL,
			CreatorID:  strconv.FormatInt(baby.UserID, 10),
			CreateTime: baby.CreatedAt,
			UpdateTime: baby.UpdatedAt,
		})
	}

	return result, nil
}

// GetBabyDetail 获取宝宝详情
func (s *BabyService) GetBabyDetail(ctx context.Context, babyID, openID string) (*dto.BabyDTO, error) {
	// 转换babyID from string to int64
	babyIDInt64, err := strconv.ParseInt(babyID, 10, 64)
	if err != nil {
		return nil, errors.New(errors.ParamError, "invalid baby id format")
	}

	// 检查权限
	if err := s.checkPermission(ctx, babyIDInt64, openID); err != nil {
		return nil, err
	}

	baby, err := s.babyRepo.FindByID(ctx, babyIDInt64)
	if err != nil {
		return nil, err
	}

	return &dto.BabyDTO{
		BabyID:     strconv.FormatInt(baby.ID, 10),
		Name:       baby.Name,
		Nickname:   baby.Nickname,
		Gender:     baby.Gender,
		BirthDate:  baby.BirthDate,
		AvatarURL:  baby.AvatarURL,
		CreatorID:  strconv.FormatInt(baby.UserID, 10),
		CreateTime: baby.CreatedAt,
		UpdateTime: baby.UpdatedAt,
	}, nil
}

// UpdateBaby 更新宝宝信息
func (s *BabyService) UpdateBaby(ctx context.Context, babyID, openID string, req *dto.UpdateBabyRequest) error {
	// 转换babyID from string to int64
	babyIDInt64, err := strconv.ParseInt(babyID, 10, 64)
	if err != nil {
		return errors.New(errors.ParamError, "invalid baby id format")
	}

	// 获取用户信息以获取UserID
	user, err := s.userRepo.FindByOpenID(ctx, openID)
	if err != nil {
		return err
	}

	// 检查编辑权限
	canEdit, err := s.collaboratorRepo.CanEdit(ctx, babyIDInt64, user.ID)
	if err != nil {
		return err
	}
	if !canEdit {
		return errors.New(errors.PermissionDenied, "您没有权限编辑该宝宝信息")
	}

	baby, err := s.babyRepo.FindByID(ctx, babyIDInt64)
	if err != nil {
		return err
	}

	// 更新字段
	if req.Name != "" {
		baby.Name = req.Name
	}
	if req.Nickname != "" {
		baby.Nickname = req.Nickname
	}
	if req.Gender != "" {
		baby.Gender = req.Gender
	}
	if req.BirthDate != "" {
		if _, err := time.Parse("2006-01-02", req.BirthDate); err != nil {
			return errors.New(errors.ParamError, "出生日期格式错误，应为YYYY-MM-DD")
		}
		baby.BirthDate = req.BirthDate
	}
	if req.AvatarURL != "" {
		baby.AvatarURL = req.AvatarURL
	}

	return s.babyRepo.Update(ctx, baby)
}

// DeleteBaby 删除宝宝
func (s *BabyService) DeleteBaby(ctx context.Context, babyID, openID string) error {
	// 转换babyID from string to int64
	babyIDInt64, err := strconv.ParseInt(babyID, 10, 64)
	if err != nil {
		return errors.New(errors.ParamError, "invalid baby id format")
	}

	// 获取用户信息以获取UserID
	user, err := s.userRepo.FindByOpenID(ctx, openID)
	if err != nil {
		return err
	}

	// 只有管理员可以删除宝宝
	isAdmin, err := s.collaboratorRepo.IsAdmin(ctx, babyIDInt64, user.ID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return errors.New(errors.PermissionDenied, "只有管理员可以删除宝宝")
	}

	return s.babyRepo.Delete(ctx, babyIDInt64)
}

// GetCollaborators 获取宝宝的协作者列表
func (s *BabyService) GetCollaborators(ctx context.Context, babyID, openID string) ([]dto.CollaboratorDTO, error) {
	// 转换babyID from string to int64
	babyIDInt64, err := strconv.ParseInt(babyID, 10, 64)
	if err != nil {
		return nil, errors.New(errors.ParamError, "invalid baby id format")
	}

	// 检查权限
	if err := s.checkPermission(ctx, babyIDInt64, openID); err != nil {
		return nil, err
	}

	collaborators, err := s.collaboratorRepo.FindByBabyID(ctx, babyIDInt64)
	if err != nil {
		return nil, err
	}

	result := make([]dto.CollaboratorDTO, 0, len(collaborators))
	for _, collab := range collaborators {
		// 检查关联的User是否被加载
		if collab.User == nil {
			s.logger.Warn("协作者关联的User未加载,跳过该协作者",
				zap.Int64("babyID", babyIDInt64),
				zap.Int64("userID", collab.UserID),
			)
			continue
		}

		result = append(result, dto.CollaboratorDTO{
			OpenID:       collab.User.OpenID,
			NickName:     collab.User.NickName,
			AvatarURL:    collab.User.AvatarURL,
			Role:         collab.Role,
			Relationship: collab.Relationship,
			AccessType:   collab.AccessType,
			ExpiresAt:    collab.ExpiresAt,
			JoinTime:     collab.CreatedAt,
		})
	}

	return result, nil
}

// InviteCollaborator 邀请协作者 (微信分享/二维码)
// 注意:同一用户对同一宝宝只有一个有效邀请,重复调用会返回已有的邀请
func (s *BabyService) InviteCollaborator(ctx context.Context, babyID, openID string, req *dto.InviteCollaboratorRequest) (*dto.BabyInvitationDTO, error) {
	// 转换babyID from string to int64
	babyIDInt64, err := strconv.ParseInt(babyID, 10, 64)
	if err != nil {
		return nil, errors.New(errors.ParamError, "invalid baby id format")
	}

	// 获取用户信息以获取UserID
	user, err := s.userRepo.FindByOpenID(ctx, openID)
	if err != nil {
		return nil, err
	}

	// 只有管理员和编辑者可以邀请
	canEdit, err := s.collaboratorRepo.CanEdit(ctx, babyIDInt64, user.ID)
	if err != nil {
		return nil, err
	}
	if !canEdit {
		return nil, errors.New(errors.PermissionDenied, "您没有权限邀请协作者")
	}

	// 获取宝宝信息
	baby, err := s.babyRepo.FindByID(ctx, babyIDInt64)
	if err != nil {
		return nil, err
	}

	// 检查是否已经存在未使用的邀请 (邀请码唯一性)
	existingInvitation, err := s.invitationRepo.FindByBabyAndInviter(ctx, babyIDInt64, user.ID)
	if err != nil && !errors.Is(err, errors.ErrNotFound) {
		return nil, err
	}

	// 如果存在未使用的邀请,直接返回已有邀请信息
	if existingInvitation != nil {
		s.logger.Info("邀请码已存在,直接返回已有记录",
			zap.String("babyID", babyID),
			zap.String("inviterID", openID),
			zap.String("shortCode", existingInvitation.ShortCode),
		)

		// 重新生成小程序码的完整URL并返回
		scene := fmt.Sprintf("c=%s", existingInvitation.ShortCode)
		qrcodeURL, errQR := s.wechatService.GenerateQRCode(ctx, scene, "pages/baby/join/join")
		if errQR != nil {
			s.logger.Error("Failed to regenerate QR code", zap.Error(errQR))
			qrcodeURL = ""
		}

		return &dto.BabyInvitationDTO{
			BabyID:      babyID,
			Name:        baby.Name,
			InviterName: user.NickName,
			Role:        existingInvitation.Role,
			ExpiresAt:   existingInvitation.ExpiresAt,
			ShortCode:   existingInvitation.ShortCode,
			QRCodeParams: &dto.QRCodeParams{
				Scene:     scene,
				Page:      "pages/baby/join/join",
				QRCodeURL: qrcodeURL,
			},
		}, nil
	}

	// 如果不存在邀请,创建新邀请
	// 生成临时token
	token := s.generateInvitationToken()

	// 生成6位短码 (用于小程序码scene参数)
	shortCode, err := s.generateUniqueShortCode(ctx)
	if err != nil {
		return nil, errors.Wrap(errors.InternalError, "生成短码失败", err)
	}

	// 创建邀请记录 (ID由snowflake自动生成)
	invitation := &entity.BabyInvitation{
		BabyID:       babyIDInt64,
		UserID:       user.ID,
		Token:        token,
		ShortCode:    shortCode,
		InviteType:   req.InviteType,
		Role:         req.Role,
		Relationship: req.Relationship,
		AccessType:   req.AccessType,
		ExpiresAt:    req.ExpiresAt, // 只保留协作者权限的过期时间
	}

	if err = s.invitationRepo.Create(ctx, invitation); err != nil {
		return nil, err
	}

	// 构建返回信息
	result := &dto.BabyInvitationDTO{
		BabyID:      babyID,
		Name:        baby.Name,
		InviterName: user.NickName,
		Role:        req.Role,
		ExpiresAt:   req.ExpiresAt,
		ShortCode:   shortCode,
	}

	// 二维码参数 - 使用短码避免32字符限制
	scene := fmt.Sprintf("c=%s", shortCode) // 仅8个字符: "c=ABC123"

	// 调用微信服务生成小程序码
	qrcodeURL, err := s.wechatService.GenerateQRCode(ctx, scene, "pages/baby/join/join")
	if err != nil {
		s.logger.Error("Failed to generate QR code", zap.Error(err))
		// 二维码生成失败不影响邀请创建,返回空URL
		qrcodeURL = ""
	}

	result.QRCodeParams = &dto.QRCodeParams{
		Scene:     scene,
		Page:      "pages/baby/join/join",
		QRCodeURL: qrcodeURL,
	}

	return result, nil
}

// GetInvitationByShortCode 通过短码获取邀请详情
func (s *BabyService) GetInvitationByShortCode(ctx context.Context, shortCode string) (*dto.InvitationDetailDTO, error) {
	// 查找邀请记录
	invitation, err := s.invitationRepo.FindByShortCode(ctx, shortCode)
	if err != nil {
		return nil, err
	}

	// 获取宝宝信息
	baby, err := s.babyRepo.FindByID(ctx, invitation.BabyID)
	if err != nil {
		return nil, err
	}

	// 获取邀请人信息 (UserID to User)
	inviter, err := s.userRepo.FindByID(ctx, invitation.UserID) // 这里需要通过UserID获取User,但userRepo没有这个方法
	// 为了兼容,我们直接从User关联获取
	// 暂时使用ID查询 - 需要检查是否有GetUserByID方法
	if err != nil {
		return nil, err
	}

	return &dto.InvitationDetailDTO{
		BabyID:      strconv.FormatInt(baby.ID, 10),
		BabyName:    baby.Name,
		BabyAvatar:  baby.AvatarURL,
		InviterName: inviter.NickName,
		Role:        invitation.Role,
		AccessType:  invitation.AccessType,
		ExpiresAt:   invitation.ExpiresAt,
		Token:       invitation.Token,
	}, nil
}

// JoinBaby 加入宝宝协作 (通过微信分享或二维码)
func (s *BabyService) JoinBaby(ctx context.Context, openID string, req *dto.JoinBabyRequest) (*dto.BabyDTO, error) {
	// 转换babyID from string to int64
	babyIDInt64, err := strconv.ParseInt(req.BabyID, 10, 64)
	if err != nil {
		return nil, errors.New(errors.ParamError, "invalid baby id format")
	}

	// 查找邀请记录
	invitation, err := s.invitationRepo.FindByToken(ctx, req.Token)
	if err != nil {
		return nil, err
	}

	if invitation.BabyID != babyIDInt64 {
		return nil, errors.New(errors.ParamError, "邀请参数不匹配")
	}

	// 获取用户信息以获取UserID
	user, err := s.userRepo.FindByOpenID(ctx, openID)
	if err != nil {
		return nil, err
	}

	// 检查用户是否已经是协作者
	existing, err := s.collaboratorRepo.FindByBabyAndUser(ctx, babyIDInt64, user.ID)
	if err != nil && !errors.Is(err, errors.ErrNotFound) {
		return nil, err
	}

	if existing != nil {
		return nil, errors.New(errors.ParamError, "您已经是该宝宝的协作者")
	}

	// 创建亲友团成员记录
	collaborator := &entity.BabyCollaborator{
		BabyID:       invitation.BabyID,
		UserID:       user.ID,
		Role:         invitation.Role,
		Relationship: invitation.Relationship,
		AccessType:   invitation.AccessType,
		ExpiresAt:    invitation.ExpiresAt,
	}

	if err := s.collaboratorRepo.Create(ctx, collaborator); err != nil {
		return nil, err
	}

	// 如果该用户没有其他宝宝,直接将该宝宝设置为默认宝宝
	if err := s.setDefaultBabyIfNeeded(ctx, openID, babyIDInt64); err != nil {
		s.logger.Error("设置默认宝宝失败", zap.Error(err))
	}

	// 返回宝宝信息
	return s.GetBabyDetail(ctx, req.BabyID, openID)
}

// RemoveCollaborator 移除协作者
func (s *BabyService) RemoveCollaborator(ctx context.Context, babyID, openID, targetOpenID string) error {
	// 转换babyID from string to int64
	babyIDInt64, err := strconv.ParseInt(babyID, 10, 64)
	if err != nil {
		return errors.New(errors.ParamError, "invalid baby id format")
	}

	// 获取用户信息以获取UserID
	user, err := s.userRepo.FindByOpenID(ctx, openID)
	if err != nil {
		return err
	}

	// 只有管理员可以移除协作者
	isAdmin, err := s.collaboratorRepo.IsAdmin(ctx, babyIDInt64, user.ID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return errors.New(errors.PermissionDenied, "只有管理员可以移除协作者")
	}

	// 不能移除创建者
	baby, err := s.babyRepo.FindByID(ctx, babyIDInt64)
	if err != nil {
		return err
	}

	// 获取目标用户信息
	targetUser, err := s.userRepo.FindByOpenID(ctx, targetOpenID)
	if err != nil {
		return err
	}

	if baby.UserID == targetUser.ID {
		return errors.New(errors.ParamError, "不能移除创建者")
	}

	return s.collaboratorRepo.Delete(ctx, babyIDInt64, targetUser.ID)
}

// UpdateCollaboratorRole 更新协作者角色
func (s *BabyService) UpdateCollaboratorRole(ctx context.Context, babyID, openID, targetOpenID, newRole string) error {
	// 转换babyID from string to int64
	babyIDInt64, err := strconv.ParseInt(babyID, 10, 64)
	if err != nil {
		return errors.New(errors.ParamError, "invalid baby id format")
	}

	// 获取用户信息以获取UserID
	user, err := s.userRepo.FindByOpenID(ctx, openID)
	if err != nil {
		return err
	}

	// 只有管理员可以更新角色
	isAdmin, err := s.collaboratorRepo.IsAdmin(ctx, babyIDInt64, user.ID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return errors.New(errors.PermissionDenied, "只有管理员可以更新协作者角色")
	}

	// 不能修改创建者角色
	baby, err := s.babyRepo.FindByID(ctx, babyIDInt64)
	if err != nil {
		return err
	}

	// 获取目标用户信息
	targetUser, err := s.userRepo.FindByOpenID(ctx, targetOpenID)
	if err != nil {
		return err
	}

	if baby.UserID == targetUser.ID {
		return errors.New(errors.ParamError, "不能修改创建者的角色")
	}

	collaborator, err := s.collaboratorRepo.FindByBabyAndUser(ctx, babyIDInt64, targetUser.ID)
	if err != nil {
		return err
	}
	if collaborator == nil {
		return errors.New(errors.NotFound, "协作者不存在")
	}

	collaborator.Role = newRole

	return s.collaboratorRepo.Update(ctx, collaborator)
}

// UpdateFamilyMember 更新亲友团成员信息 (角色和关系)
func (s *BabyService) UpdateFamilyMember(ctx context.Context, babyID, openID, targetOpenID string, req *dto.UpdateFamilyMemberRequest) error {
	// 转换babyID from string to int64
	babyIDInt64, err := strconv.ParseInt(babyID, 10, 64)
	if err != nil {
		return errors.New(errors.ParamError, "invalid baby id format")
	}

	// 获取用户信息以获取UserID
	user, err := s.userRepo.FindByOpenID(ctx, openID)
	if err != nil {
		return err
	}

	// 检查权限：管理员可以更新任何成员信息，用户可以更新自己的关系
	isAdmin, err := s.collaboratorRepo.IsAdmin(ctx, babyIDInt64, user.ID)
	if err != nil {
		return err
	}

	// 如果不是管理员，只能更新自己的关系
	if !isAdmin {
		// 检查是否是更新自己的信息
		targetUser, err := s.userRepo.FindByOpenID(ctx, targetOpenID)
		if err != nil {
			return err
		}

		// 如果不是自己，拒绝访问
		if user.ID != targetUser.ID {
			return errors.New(errors.PermissionDenied, "只有管理员可以更新其他成员信息")
		}

		// 如果是更新自己的信息，只能更新关系，不能更新角色
		if req.Role != "" {
			return errors.New(errors.PermissionDenied, "只有管理员可以修改角色")
		}
	}

	// 获取目标用户信息
	targetUser, err := s.userRepo.FindByOpenID(ctx, targetOpenID)
	if err != nil {
		return err
	}

	collaborator, err := s.collaboratorRepo.FindByBabyAndUser(ctx, babyIDInt64, targetUser.ID)
	if err != nil {
		return err
	}
	if collaborator == nil {
		return errors.New(errors.NotFound, "亲友团成员不存在")
	}

	// 不能修改创建者角色
	baby, err := s.babyRepo.FindByID(ctx, babyIDInt64)
	if err != nil {
		return err
	}

	// 更新角色 (如果提供且不是创建者)
	if req.Role != "" {
		if baby.UserID == targetUser.ID {
			return errors.New(errors.ParamError, "不能修改创建者的角色")
		}
		collaborator.Role = req.Role
	}

	// 更新关系
	if req.Relationship != "" {
		collaborator.Relationship = req.Relationship
	}

	return s.collaboratorRepo.Update(ctx, collaborator)
}

// checkPermission 检查用户是否有权限访问宝宝
func (s *BabyService) checkPermission(ctx context.Context, babyIDInt64 int64, openID string) error {
	// 获取用户信息以获取UserID
	user, err := s.userRepo.FindByOpenID(ctx, openID)
	if err != nil {
		return err
	}

	isCollaborator, err := s.collaboratorRepo.IsCollaborator(ctx, babyIDInt64, user.ID)
	if err != nil {
		return err
	}
	if !isCollaborator {
		return errors.New(errors.PermissionDenied, "您没有权限访问该宝宝")
	}
	return nil
}

// copyCollaborators 复制协作者列表到新宝宝
func (s *BabyService) copyCollaborators(ctx context.Context, sourceBabyIDInt64, targetBabyIDInt64 int64, openID string) error {
	// 检查源宝宝的权限
	if err := s.checkPermission(ctx, sourceBabyIDInt64, openID); err != nil {
		return err
	}

	// 获取源宝宝的协作者
	sourceCollaborators, err := s.collaboratorRepo.FindByBabyID(ctx, sourceBabyIDInt64)
	if err != nil {
		return err
	}

	// 获取用户信息以获取UserID
	user, err := s.userRepo.FindByOpenID(ctx, openID)
	if err != nil {
		return err
	}

	// 创建新的协作者列表(排除创建者,因为已经添加)
	newCollaborators := make([]*entity.BabyCollaborator, 0)

	for _, collab := range sourceCollaborators {
		if collab.UserID == user.ID {
			continue // 跳过创建者
		}

		newCollaborators = append(newCollaborators, &entity.BabyCollaborator{
			BabyID:       targetBabyIDInt64,
			UserID:       collab.UserID,
			Role:         collab.Role,
			Relationship: collab.Relationship,
			AccessType:   collab.AccessType,
			ExpiresAt:    collab.ExpiresAt,
		})
	}

	return s.collaboratorRepo.BatchCreate(ctx, newCollaborators)
}

// generateInvitationCode 生成邀请码(已废弃,保留兼容)
func (s *BabyService) generateInvitationCode() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// generateInvitationToken 生成邀请token(用于微信分享和二维码)
func (s *BabyService) generateInvitationToken() string {
	bytes := make([]byte, 32) // 64位十六进制字符串
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// generateUniqueShortCode 生成唯一短码
// 循环尝试生成直到找到唯一的短码
func (s *BabyService) generateUniqueShortCode(ctx context.Context) (string, error) {
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		// 生成6位短码
		shortCode, err := utils.GenerateShortCode()
		if err != nil {
			return "", err
		}

		// 检查短码是否已存在
		_, err = s.invitationRepo.FindByShortCode(ctx, shortCode)
		if err != nil {
			// 检查是否是 NotFound 错误
			if appErr, ok := err.(*errors.AppError); ok && appErr.Code == errors.NotFound {
				// 短码不存在,可以使用
				return shortCode, nil
			}
			// 其他错误直接返回
			return "", err
		}

		// 短码已存在(没有错误),继续下一次尝试
	}

	return "", errors.New(errors.InternalError, "无法生成唯一短码,请稍后重试")
}

// setDefaultBabyIfNeeded 如果用户没有其他宝宝,则将该宝宝设置为默认宝宝
func (s *BabyService) setDefaultBabyIfNeeded(ctx context.Context, openID string, babyID int64) error {
	// 获取用户信息
	user, err := s.userRepo.FindByOpenID(ctx, openID)
	if err != nil {
		return err
	}

	// 如果用户已有默认宝宝,则不需要设置
	if user.DefaultBabyID != 0 {
		return nil
	}

	// 设置为默认宝宝
	return s.userRepo.UpdateDefaultBabyID(ctx, openID, babyID)
}
